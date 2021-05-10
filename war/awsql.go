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

import (
		"log"
		"database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
)

var playerMax = 5

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
	ch <- memoItem2{pds: best, score: bestScore, err: nil, callArgs: callArgs}

	return best, bestScore, nil
}

func findBestBGHelper(ch chan memoItem2, diversity map[HeroVal]bool, roster map[string][]Champ, players []string, callArgs Defenders) ([]PlayerDefenders, float32, error) {
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

func assignChamps(occupiedNodes map[int]Champ, remainingChamps []Champ, skippedChamps []Champ) (map[Champ]int, int, error) {
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
		return assignChamps(occupiedNodes, remainingChamps[1:], skippedChamps)
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

			result, score, err := assignChamps(newOccupiedNodes, remainingChamps[1:], skippedChamps)
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
		return assignChamps(occupiedNodes, remainingChamps[1:], skippedChamps)
	}

	return bestMap, bestScore, nil
}

func assignChampsHelper(bg map[string][]Champ, occupiedNodes map[int]Champ, remainingChamps []Champ, skippedChamps []Champ) (map[Champ]int, int, error) {
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
    if n == 1 || n == 3 || n == 4 || n == 5 || n == 6 {
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
  for idx := 0; idx < len(nodeMap); idx++ {
    fmt.Printf("idx: %v, node %v\n", idx, nodeMap[idx])

    var preferences []int
    seen := map[HeroVal]bool{}

    for _, h := range Nodes[nodeMap[idx]] {
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
    for ii := len(remainingChamps)-1; ii >= 0; ii-- {
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
                                where player = ?`, player.Id)
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

func run(bg map[string][]Champ) []PlayerDefenders {
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
	for ; time.Now().Sub(tStart) < time.Minute; {
		memo2 = map[string]memoItem2{}
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

	var allChamps []Champ
	for _, pd := range bestResult {
		for _, c := range pd.Defenders.Champs {
			allChamps = Insert(allChamps, c)
		}
		fmt.Printf("%s: %s\n", pd.Player, pd.Defenders.String())
	}
	fmt.Printf("Score: %v\n", bestScore)
	fmt.Printf("Allchamps: %v\n", allChamps)

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
			output = append(output, c.Champ.String())
			output = append(output, fmt.Sprintf("(%v) ", result[c]))
		}
		fmt.Printf("%v\n", strings.Join(output, ""))
	}

  return bestResult
}

type Player struct {
  Id int
  Suicides bool
  MD int
  Name string
}

func BestWarDefense(alliance int, bg int) []PlayerDefenders {
  //writeBg(bg1)
  bgRoster, err := getBg(bg)
  if err != nil {
    log.Fatal(err)
  }
	return run(bgRoster)
}

