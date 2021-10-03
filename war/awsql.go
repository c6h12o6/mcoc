package war

import . "github.com/c6h12o6/mcoc/globals"
import "github.com/c6h12o6/mcoc/smp"

//import "github.com/satori/go.uuid"
import "fmt"
import "strings"
import "sort"
import "time"
import "sync"
import "math"
import "github.com/bradfitz/slice"
import "math/rand"
import "os"
import pb "github.com/c6h12o6/mcoc/proto"

import (
	"database/sql"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var playerMax = 5
var problemChild map[string]int
var problemChildLock sync.Mutex

var sqlPassword = os.Getenv("CLOUD_SQL_PASSWORD")

// combinations is a helper function for creating all possible combinations of
// values from "iterable" in groups of "r"
func combinations(iterable []int, r int) [][]int {
	var ret [][]int

	pool := iterable
	n := len(pool)

	if r > n {
		return [][]int{}
	}

	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}

	result := make([]int, r)
	for i, el := range indices {
		result[i] = pool[el]
	}

	tmp := make([]int, len(result))
	copy(tmp, result)
	ret = append(ret, tmp)

	for {
		i := r - 1
		for ; i >= 0 && indices[i] == i+n-r; i -= 1 {
		}

		if i < 0 {
			break
		}

		indices[i] += 1
		for j := i + 1; j < r; j += 1 {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i += 1 {
			result[i] = pool[indices[i]]
		}

		tmp := make([]int, len(result))
		copy(tmp, result)
		ret = append(ret, tmp)
	}

	return ret
}

// teamCombinations creates all combinations of all posssible champions
// into groups of "teamsize" and populates their synergy count
func teamCombinations(teamsize int, roster []Champ, md int, suicides bool) []Defenders {
	// TODO maybe only do this once
	var indices []int
	var teams []Defenders

	// Seriously, go does not have a way to create a slice of 1-n
	// so you have to do this. gross
	for ii := 0; ii < len(roster); ii++ {
		indices = append(indices, ii)
	}

	// Now get the combinations you need, in the form of slices of ints
	teamindices := combinations(indices, teamsize)

	// Turn those slices of integers into slices of heroes
	for _, teamnos := range teamindices {
		team := Defenders{}
		var score float32
		for _, idx := range teamnos {
			team.Champs = append(team.Champs, roster[idx])
			score += ChampScoreWithMasteries(roster[idx], md, suicides)
		}
		team.score = score
		// that has a synergy count of 0. They can't help our count
		teams = append(teams, team)
	}

	// return the list in order of count, high to low
	slice.Sort(teams, func(i, j int) bool {
		return teams[i].score > teams[j].score
	})

	return teams
}

var memoLock sync.Mutex
var memoCount int

type Defenders struct {
	player string
	Champs []Champ
	score  float32
}

type PlayerDefenders struct {
	Player    string
	Defenders Defenders
}

func (d *Defenders) String() string {
	var ret []string
	for _, c := range d.Champs {
		ret = append(ret, c.String())
	}
	ret = append(ret, fmt.Sprintf("%v", d.score))
	return strings.Join(ret, ",")
}

func copyDiversity(d map[HeroVal]bool) map[HeroVal]bool {
	ret := map[HeroVal]bool{}
	for k, v := range d {
		ret[k] = v
	}
	return ret
}

type memoItem2 struct {
	pds      []PlayerDefenders
	score    float32
	err      error
	callArgs Defenders
}

var memo2 = map[string]memoItem2{}

func getMemoKey(diversity map[HeroVal]bool, players []string) string {
	var ret []string
	for h, _ := range diversity {
		ret = append(ret, h.String())
	}
	sort.Strings(ret)
	key := strings.Join(ret, ",") + "," + strings.Join(players, ",")
	return key
}

func recordMemo2(memoKey string, pds []PlayerDefenders, score float32, err error, callArgs Defenders) {
	memoLock.Lock()
	memo2[memoKey] = memoItem2{pds: pds, score: score, err: err, callArgs: callArgs}
	memoLock.Unlock()
}

func findBestBG(ch chan memoItem2, diversity map[HeroVal]bool, roster map[string][]Champ, players []string, callArgs Defenders) ([]PlayerDefenders, float32, error) {
	best := []PlayerDefenders{}
	var bestScore float32
	newCh := make(chan memoItem2)
	var calls int

	memoKey := getMemoKey(diversity, players)

	memoLock.Lock()
	mi, ok := memo2[memoKey]
	memoLock.Unlock()
	if ok {
		memoCount++

		ch <- memoItem2{pds: mi.pds, score: mi.score, err: mi.err, callArgs: callArgs}
		return mi.pds, mi.score, mi.err
	}

	if len(players) != 0 {
		p := players[0]
		playerChamps := roster[p]
		//mastery := masteryMap[p]
		mastery := Player{}

		var reducedChamps []Champ
		var lockedChamps []Champ
		var lockedChampScore float32

		var debug bool

		for _, c := range playerChamps {
			// all locked nodes should also already be in the diversity map
			if c.LockedNode != 0 {
				lockedChamps = append(lockedChamps, c)
				lockedChampScore += ChampScoreWithMasteries(c, mastery.MD, mastery.Suicides)
			}
			if _, ok := diversity[c.Champ]; ok {
				continue
			}

			reducedChamps = append(reducedChamps, c)
		}
		if len(reducedChamps) < playerMax {
			recordMemo2(memoKey, nil, 0, fmt.Errorf("No valid teams"), callArgs)
			ch <- memoItem2{pds: nil, score: 0, err: fmt.Errorf("No valid teams")}

      problemChildLock.Lock()
      problemChild[p] += 1
      problemChildLock.Unlock()

			return nil, 0, fmt.Errorf("No valid teams")
		}
		combos := teamCombinations(playerMax-len(lockedChamps), reducedChamps, mastery.MD, mastery.Suicides)
		combos = combos[:int(math.Min(float64(len(combos)), 5))]
		newCombos := make([]Defenders, len(combos))

		for ii, d := range combos {
			d.score += lockedChampScore
			newCombos[ii] = Defenders{
				Champs: append(d.Champs, lockedChamps...),
				score:  d.score + lockedChampScore,
			}
		}
		combos = newCombos
		if debug {
			for _, d := range combos {
				fmt.Printf("combos %v\n", d.String())
			}
		}
		for _, d := range combos {
			newDiversity := copyDiversity(diversity)
			for _, champ := range d.Champs {
				newDiversity[champ.Champ] = true
			}
			calls++
			go findBestBG(newCh, newDiversity, roster, players[1:], d)
		}
		for ; calls > 0; calls-- {
			select {
			case mi := <-newCh:
				result, newScore, err, ca := mi.pds, mi.score, mi.err, mi.callArgs

				if err != nil {
					continue
				}
				if newScore+ca.score > bestScore {
					bestScore = newScore + ca.score
					best = append(result, PlayerDefenders{Player: p, Defenders: ca})
				}
			}
		}
	}
	recordMemo2(memoKey, best, bestScore, nil, callArgs)
  if bestScore == 0 {
    problemChildLock.Lock()
    defer problemChildLock.Unlock()
    ch <- memoItem2{pds: nil, score: bestScore, err: fmt.Errorf("No valid mappings found: problemChild map: %v", problemChild), callArgs: callArgs}
    return nil, 0, fmt.Errorf("No valid mappings found: problemChild map: %v", problemChild)
  }
	ch <- memoItem2{pds: best, score: bestScore, err: nil, callArgs: callArgs}

	return best, bestScore, nil
}

func findBestBGHelper(ch chan memoItem2, diversity map[HeroVal]bool, roster map[string][]Champ, players []string, callArgs Defenders) ([]PlayerDefenders, float32, error) {
  // TODO make this thread safe
  problemChildLock.Lock()
  problemChild = map[string]int{}
  problemChildLock.Unlock()
	for _, p := range players {
		for _, c := range roster[p] {
			if c.LockedNode != 0 {
				diversity[c.Champ] = true
			}
		}
	}
	fmt.Printf("starting diversity; %v\n", diversity)
	return findBestBG(ch, diversity, roster, players, callArgs)
}

func removeFromRemainingChamps(r []HeroVal, tbr HeroVal) []HeroVal {
	var ret []HeroVal
	for _, h := range r {
		if h == tbr {
			continue
		}
		ret = append(ret, h)
	}
	return ret
}

func copyOccupiedNodes(o map[int]Champ) map[int]Champ {
	ret := map[int]Champ{}
	for k, v := range o {
		ret[k] = v
	}
	return ret
}

func Contains(n []HeroVal, c HeroVal) bool {
	for _, h := range n {
		if c == h {
			return true
		}
	}
	return false
}

func assignChamps(preferences map[int][]HeroVal, occupiedNodes map[int]Champ, remainingChamps []Champ, skippedChamps []Champ) (map[Champ]int, int, error) {
	bestScore := 0
	bestMap := map[Champ]int{}

	if len(remainingChamps) == 0 {
		for n := 55; n > 0 && len(skippedChamps) > 0; n-- {
			if _, ok := occupiedNodes[n]; !ok {
				bestMap[skippedChamps[0]] = n
				skippedChamps = skippedChamps[1:]
			}
		}

		return bestMap, 0, nil
	}

	assigned := false
	c := remainingChamps[0]

	if ChampValue(c) < 7 {
		skippedChamps = append(skippedChamps, c)
		return assignChamps(preferences, occupiedNodes, remainingChamps[1:], skippedChamps)
	}

	var count int
	for n := 55; n > 0; n-- {
		if _, ok := occupiedNodes[n]; ok {
			continue
		}

		if Contains(Nodes[n], c.Champ) {
			count++
			assigned = true
			newOccupiedNodes := copyOccupiedNodes(occupiedNodes)
			newOccupiedNodes[n] = c

			result, score, err := assignChamps(preferences, newOccupiedNodes, remainingChamps[1:], skippedChamps)
			if err == nil && score+1 > bestScore {
				bestScore = score + 1
				bestMap = result
				bestMap[c] = n
			}
		}
		if count > 1 {
			break
		}
	}
	if !assigned {
		skippedChamps = append(skippedChamps, c)
		return assignChamps(preferences, occupiedNodes, remainingChamps[1:], skippedChamps)
	}

	return bestMap, bestScore, nil
}

func getMapPreferences(alliance int) {
}

func assignChampsHelper(bg map[string][]Champ, occupiedNodes map[int]Champ, remainingChamps []Champ, skippedChamps []Champ) (map[Champ]int, int, error) {
  fmt.Printf("assignChampsHelper!!!!!!!!!!!!!!!!!!!!!!\n");
	locked := map[Champ]int{}
	for _, champlist := range bg {
		for _, c := range champlist {
			if c.LockedNode != 0 {
				fmt.Printf("%v is locked\n", c.LockedNode)
				occupiedNodes[c.LockedNode] = c
				locked[c] = c.LockedNode
			}
		}
	}

	nodeMap := map[int]int{}
	champMap := map[HeroVal]int{}
	reverseChampMap := map[int]Champ{}
	var nodePersonId int

	var champPreference []int
	for n := 55; n > 0; n-- {
		if n == 1 || n == 3 || n == 5 || n == 6 || n == 8 {
			continue
		}
		if _, ok := occupiedNodes[n]; !ok {
			champPreference = append(champPreference, nodePersonId)
			nodeMap[nodePersonId] = n
			nodePersonId++
		}
	}
	fmt.Printf("%v\n", nodeMap)

	champPersonId := 0
	var champPeople []smp.Person
	for idx, c := range remainingChamps {
		if c.LockedNode != 0 {
			continue
		}
		champMap[c.Champ] = champPersonId
		reverseChampMap[champPersonId] = c
		champPersonId++
		champPeople = append(champPeople, smp.Person{ID: idx, Prefers: champPreference})
	}
	fmt.Printf("%v\n", champMap)

	fmt.Printf("%v %v\n", len(nodeMap), len(champMap))

	var nodePeople []smp.Person
  AWPreferences, err := getWarMap(1)
  if err != nil {
    fmt.Printf("error getting war map: %v", err)
    return nil, 0, err
  }

	for idx := 0; idx < len(nodeMap); idx++ {
		fmt.Printf("idx: %v, node %v\n", idx, nodeMap[idx])

		var preferences []int
		seen := map[HeroVal]bool{}

		for _, h := range AWPreferences[nodeMap[idx]] {
			// If the champ isn't an option on this roster, skip it
			if _, ok := champMap[h]; !ok {
				continue
			}

			// Handle legacy code
			if h == MaxHeroVal {
				break
			}
			seen[h] = true
			preferences = append(preferences, champMap[h])
		}
		// fill in from least powerful to most powerful so we don't steal away a more powerful
		// champ from a node theyre good on
		for ii := len(remainingChamps) - 1; ii >= 0; ii-- {
			c := remainingChamps[ii]
			if _, ok := seen[c.Champ]; !ok && c.LockedNode == 0 {
				preferences = append(preferences, champMap[c.Champ])
			}
		}

		fmt.Printf("preferences: %v %v\n", len(preferences), preferences)
		nodePeople = append(nodePeople, smp.Person{ID: idx, Prefers: preferences})
	}

	smp.StageMarriage(nodePeople, champPeople, len(champPeople))

	ret := map[Champ]int{}
	for i := 0; i < len(champPeople); i++ {
		ret[reverseChampMap[i]] = nodeMap[champPeople[i].Partner.ID]
	}
	for k, v := range locked {
		ret[k] = v
	}

	return ret, 1, nil
}

func permutations(arr []string) [][]string {
	var helper func([]string, int)
	res := [][]string{}

	helper = func(arr []string, n int) {
		if n == 1 {
			tmp := make([]string, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func Insert(sorted []Champ, champ Champ) []Champ {
	i := sort.Search(len(sorted), func(i int) bool { return ChampScore(sorted[i]) < ChampScore(champ) })
	sorted = append(sorted, Champ{})
	copy(sorted[i+1:], sorted[i:])
	sorted[i] = champ
	return sorted
}

func getBg(bg int) (map[string][]Champ, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@cloudsql(homeproject:us-east1:champdb)/champdb", sqlPassword))
	defer db.Close()

	ret := map[string][]Champ{}

	if err != nil {
		return ret, err
	}

	rows, err := db.Query(`select id, suicides, mystic_dispersion, name from players
                         where alliance = ? AND BG = ?`, 1, bg)
	if err != nil {
		return ret, err
	}

	for rows.Next() {
		var player Player
		err = rows.Scan(&player.Id, &player.Suicides, &player.MD, &player.Name)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		fmt.Printf("%+v\n", player)
		champRows, err := db.Query(`select heroval, stars, herorank, signature, locked from champ
                                where player = ? and deleted = 0`, player.Id)
		var champs []Champ
		for champRows.Next() {
			var c Champ
			err = champRows.Scan(&c.Champ, &c.Stars, &c.Level, &c.Sig, &c.LockedNode)
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			champs = append(champs, c)
		}
		ret[player.Name] = champs
	}

	return ret, nil
}

func getWarMap(alliance int) (map[int][]HeroVal, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@cloudsql(homeproject:us-east1:champdb)/champdb", sqlPassword))
	defer db.Close()

	ret := map[int][]HeroVal{}

	if err != nil {
		return ret, err
	}

	rows, err := db.Query(`select node, champ1, champ2, champ3, champ4, champ5 from nodes where alliance = ?`, alliance)
	if err != nil {
		return ret, err
	}

  type node struct {
    node int
    champ1 int
    champ2 int
    champ3 int
    champ4 int
    champ5 int
  }

	for rows.Next() {
    n := node{}
		err = rows.Scan(&n.node, &n.champ1, &n.champ2, &n.champ3, &n.champ4, &n.champ5)
		if err != nil {
			fmt.Printf("Error retrieving nodes: %v\n", err)
			continue
		}

    ret[n.node] = []HeroVal{
      HeroVal(n.champ1),
      HeroVal(n.champ2),
      HeroVal(n.champ3),
      HeroVal(n.champ4),
      HeroVal(n.champ5)}
  }

  fmt.Printf("full map is %v", ret)

	return ret, nil
}

type PlayerChamp struct {
  Player string
  Champ pb.Champ
}

type Algo2Memo struct {
  pref int
  node int
}

var algo2Memo map[Algo2Memo]string

func getBestChampOptions(hv HeroVal, rosters map[string][]Champ, playerCount map[string]int, diversity map[string]bool) ([]PlayerChamp, error) {
  //var bestChamp PlayerChamp
  var bestScore float32

  type PlayerChampScore struct {
    pc PlayerChamp
    score float32
  }
  var options []PlayerChampScore

  factor := float32(0.8)

  err := fmt.Errorf("Unable to find %v on an open roster", hv)

  for p, champList := range rosters {
    if playerCount[p] == 5 {
      continue
    }
    for _, c := range champList {
      if _, ok := diversity[c.Champ.String()]; !ok {
        if hv == Empty || c.Champ == hv {
          if score := ChampValue(c); score > factor * bestScore {
            pc := PlayerChamp{Player: p, Champ: pb.Champ{
                ChampName: c.Champ.String(),
                Stars: c.Stars,
                Rank: c.Level,
                Sig: c.Sig,
                LockedNode: int32(c.LockedNode),
            }}
            options = append(options, PlayerChampScore{pc, score})

            if score > bestScore {
              //bestChamp = pc
              bestScore = score
              err = nil
            }
            break
          }
        }
      }
    }
  }
  finalOptions := []PlayerChamp{}
  for _, opt := range options {
    fmt.Printf("\t%v %v/%v Sig: %v  Ratio: %v Score: %v\n", opt.pc.Champ.ChampName, opt.pc.Champ.Stars, opt.pc.Champ.Rank, opt.pc.Champ.Sig,
        opt.score / bestScore, opt.score)
    if opt.score / bestScore >= factor {
      finalOptions = append(finalOptions, opt.pc)
    }
  }
  fmt.Printf("pc for sugar is %v", playerCount["sugar"])
  return finalOptions, err
}

func getBestChamp(hv HeroVal, rosters map[string][]Champ, playerCount map[string]int, diversity map[string]bool) (PlayerChamp, error) {
  var bestChamp PlayerChamp
  var bestScore float32

  err := fmt.Errorf("Unable to find %v on an open roster", hv)

  for p, champList := range rosters {
    if playerCount[p] == 5 {
      continue
    }
    for _, c := range champList {
      if _, ok := diversity[c.Champ.String()]; !ok {
        if hv == Empty || c.Champ == hv {
          if score := ChampScore(c); score > bestScore {
            bestChamp = PlayerChamp{Player: p, Champ: pb.Champ{
                ChampName: c.Champ.String(),
                Stars: c.Stars,
                Rank: c.Level,
                Sig: c.Sig,
                LockedNode: int32(c.LockedNode),
            }}
            bestScore = score
            err = nil
            break
          }
        }
      }
    }
  }
  fmt.Printf("returning best champ: %v %v \n", bestChamp, bestScore)
  return bestChamp, err
}

func runAlgo2(bg map[string][]Champ, defenderPreferences map[int][]HeroVal, playerCount map[string]int, solved map[int]PlayerChamp, diversity map[string]bool, maxDepth int) (map[int]PlayerChamp, int, error) {


  for pref := 0; pref < maxDepth; pref++ {
    // 5 is a special case, where we fill in the rest of the slots
    for nodeNo := 55; nodeNo > 0; nodeNo-- {
      if _, ok := solved[nodeNo]; ok {
        continue
      }
      var preference HeroVal
      if pref == 5 {
        fmt.Printf("---------FILLING IN THE REST %v\n", diversity[Empty.String()])
        preference = Empty
      } else {
        preference = defenderPreferences[nodeNo][pref]
      }

      if pref != 5 && preference == Empty {
        continue
      }
      if _, ok := diversity[preference.String()]; ok {
        fmt.Printf("Skipping %v at node %v because it's already assigned\n", preference, nodeNo)
        continue
      }
      pcs := []PlayerChamp{}
      mi := Algo2Memo{pref: pref, node: nodeNo}
      if pref == 5 {
        // Dont try all options when we're just filling in the map
        pc, err := getBestChamp(preference, bg, playerCount, diversity)
        if err != nil {
          fmt.Printf("Cant fill %v: %v\n", nodeNo, err)
          continue
        }
        pcs = append(pcs, pc)
      } else {
        pcsTmp, err := getBestChampOptions(preference, bg, playerCount, diversity)
        if err != nil {
          fmt.Printf("Cant fill %v: %v\n", nodeNo, err)
          continue
        }
        smallestRoster := 9999999
        var selectedPc PlayerChamp
        for _, pc := range pcsTmp {
          if len(bg[pc.Player]) < smallestRoster {
            smallestRoster = len(bg[pc.Player])
            selectedPc = pc
          }
        }
        fmt.Printf("selecting %v for node %v: %v", selectedPc.Player, nodeNo, selectedPc.Champ.ChampName)
        pcs = []PlayerChamp{selectedPc}
        /*
        pcsTmp, err := getBestChampOptions(preference, bg, playerCount, diversity)
        if err != nil {
          fmt.Printf("Cant fill %v: %v\n", nodeNo, err)
          continue
        }

        if p, ok := algo2Memo[mi]; ok {
          fmt.Printf("trying to pick %v's %v\n", p, defenderPreferences[nodeNo][pref])
          //var best PlayerChamp
          for _, pc := range pcsTmp {
            if pc.Player == p && playerCount[p] != 5 {
              pcs = []PlayerChamp{pc}
              break
            }
          }
          if len(pcs) == 0 {
            fmt.Printf("$$$$$$$$$$$$$$$$$$$$$$$$$$ didn't expect to get here\n")
          }
        }
        if len(pcs) == 0 {
          pcs = pcsTmp
        }
        */
      }
      fmt.Printf("-- PCS is %v\n", pcs)
      if len(pcs) == 1 {
        pc := pcs[0]
        fmt.Printf("playerCount for %v is %v\n", pc.Player, playerCount[pc.Player])
        playerCount[pc.Player]++
        pc.Champ.AssignedNode = int32(nodeNo)
        solved[nodeNo] = PlayerChamp{Player: pc.Player, Champ: pc.Champ}
        diversity[pc.Champ.ChampName] = true
      } else {
        var bestResult map[int]PlayerChamp
        for _, pc := range pcs {

          // copy some stuff
          newSolved := map[int]PlayerChamp{}
          for k, v := range solved {
            newSolved[k] = v
          }

          newPlayerCount := map[string]int{}
          for k, v := range playerCount {
            playerCount[k] = v
          }

          newDiversity := map[string]bool{}
          for k, v := range diversity {
            newDiversity[k] = v
          }

          newPlayerCount[pc.Player]++
          pc.Champ.AssignedNode = int32(nodeNo)
          newSolved[nodeNo] = PlayerChamp{Player: pc.Player, Champ: pc.Champ}
          newDiversity[pc.Champ.ChampName] = true

          fmt.Printf("running subtree for %v from %v\n", pc.Champ.ChampName, pc.Player)
          result, _, err := runAlgo2(bg, defenderPreferences, newPlayerCount, newSolved, newDiversity, 5)
          if err != nil {
            fmt.Printf("wtf: %v\n", err)
            continue
          }
          playerChampsOnNodes := 0
          for _, pcTmp := range result {
            if pcTmp.Player == pc.Player {
              playerChampsOnNodes++
            }
          }
          bestResult = result
        }

        if bestResult == nil {
          panic(fmt.Sprintf("Got jack squat for node %v pref %v: %v\n", nodeNo, pref, defenderPreferences[nodeNo][pref])) 
        }
        fmt.Printf("Saving that %v is best for %v at %v", bestResult[nodeNo].Player, bestResult[nodeNo].Champ.ChampName, nodeNo)
        algo2Memo[mi] = bestResult[nodeNo].Player

        solved[nodeNo] = bestResult[nodeNo]
        diversity[solved[nodeNo].Champ.ChampName] = true
        playerCount[solved[nodeNo].Player]++
      }
    }
  }

  return solved, 0, nil
}

func run(bg map[string][]Champ) ([]PlayerDefenders, error) {
	t := time.Now()

	var players []string
	for p, _ := range bg {
		players = append(players, p)
	}

	playerPermutations := permutations(players)
	permD := time.Now().Sub(t)
	fmt.Printf("%v\n", permD)

	var bestResult []PlayerDefenders
	var bestScore float32
	tStart := time.Now()
	for time.Now().Sub(tStart) < time.Minute {
		memo2 = map[string]memoItem2{}
    fmt.Printf("Len playerpermutaitons: %v\n", len(playerPermutations))
		playerList := playerPermutations[rand.Intn(len(playerPermutations))]
		fmt.Printf("\tTrying %v\n", playerList)

		t = time.Now()
		ch := make(chan memoItem2)
		go findBestBGHelper(ch, map[HeroVal]bool{}, bg, playerList, Defenders{})
		select {
		case mi := <-ch:
			result, score, _ := mi.pds, mi.score, mi.err

			d := time.Now().Sub(t)
			if len(result) == len(players) {
				for _, pd := range result {
					fmt.Printf("%s: %s\n", pd.Player, pd.Defenders.String())
				}
				fmt.Printf("Took %v for %v combos\n", d, len(result))
				fmt.Printf("Score: %v\n", score)

				if score > bestScore {
					bestResult = result
					bestScore = score
				}
			}
		}
	}

  if bestScore == 0 {
    problemChildLock.Lock()
    defer problemChildLock.Unlock()
    return nil, fmt.Errorf("No valid mappings found: problemChild map: %v", problemChild)
  }


	var allChamps []Champ
	for _, pd := range bestResult {
		for _, c := range pd.Defenders.Champs {
			allChamps = Insert(allChamps, c)
		}
		fmt.Printf("%s: %s\n", pd.Player, pd.Defenders.String())
	}
	fmt.Printf("Score: %v\n", bestScore)
	fmt.Printf("Allchamps: %v\n", allChamps)
  fmt.Printf("------------------------------------------- JBF2\n")

	result, _, _ := assignChampsHelper(bg, map[int]Champ{}, allChamps, []Champ{})
	fmt.Printf("result length: %v\n", len(result))
	var maplines []string
	assigned := map[int]Champ{}
	for c, n := range result {
		assigned[n] = c
		maplines = append(maplines, fmt.Sprintf("%02d: %v", n, c))
	}
	sort.Strings(maplines)
	fmt.Print(strings.Join(maplines, "\n"))

	var unplacedChampsStrings []string
	var unplacedChamps []Champ
	for _, c := range allChamps {
		_, ok := result[c]
		if !ok {
			unplacedChampsStrings = append(unplacedChampsStrings, c.String())
			unplacedChamps = append(unplacedChamps, c)
		}
	}
	fmt.Printf("\n")
	fmt.Print(strings.Join(unplacedChampsStrings, ","))

	for _, pd := range bestResult {
		var output []string
		output = append(output, fmt.Sprintf("%s: ", pd.Player))
		for _, c := range pd.Defenders.Champs {
      c.AssignedNode = int32(result[c])
			output = append(output, c.Champ.String())
			output = append(output, fmt.Sprintf("(%v) ", result[c]))
		}
		fmt.Printf("%v\n", strings.Join(output, ""))
	}

	return bestResult, nil
}

type Player struct {
	Id       int
	Suicides bool
	MD       int
	Name     string
}


func algo2Helper(roster map[string][]Champ) ([]PlayerDefenders, error) {


  algo2Memo = map[Algo2Memo]string{}
  defenderPreferences, err := getWarMap(1)
  if err != nil {
    return nil, err
  }

	assignments, _, err := runAlgo2(roster, defenderPreferences, map[string]int{}, map[int]PlayerChamp{}, map[string]bool{}, 6)
  if err != nil {
    fmt.Printf("algo2 failed: %v\n", err)
    return nil, err
  }

  tmp := map[string]*Defenders{}
  //for n, c := range assignments {
  for idx := 55; idx > 0; idx-- {
    c := assignments[idx]
    fmt.Printf("%v: %v (%v, %v/%v)\n", idx, c.Champ.ChampName, c.Player, c.Champ.Stars, c.Champ.Rank)

    if _, ok := tmp[c.Player]; !ok {
      tmp[c.Player] = &Defenders{player: c.Player}
    }

    oldStyleChamp := Champ {
        Champ: NameToValue(c.Champ.ChampName),
        Level: c.Champ.Rank,
        Stars: c.Champ.Stars,
        Sig: c.Champ.Sig,
        AssignedNode: c.Champ.AssignedNode,
    }

    tmp[c.Player].Champs = append(tmp[c.Player].Champs, oldStyleChamp)
    tmp[c.Player].score += ChampValue(oldStyleChamp)
  }

  for p, d := range tmp {
    fmt.Printf("%v:", p)
    for _, c := range d.Champs {
      fmt.Printf(" %v %v/%v (%v)", c.Champ, c.Stars, c.Level, c.AssignedNode)
    }
    fmt.Printf("\n");
  }
  return nil, fmt.Errorf("Stubbed out")
}
func BestWarDefense(alliance int, bg int) ([]PlayerDefenders, error) {
	//writeBg(bg1)
  fmt.Printf("----------------------------------------- JBF\n")
	bgRoster, err := getBg(bg)
	if err != nil {
		log.Fatal(err)
	}
	return algo2Helper(bgRoster)
}
