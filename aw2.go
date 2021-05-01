package main

import . "../mcoc/globals"
import "../mcoc/smp"

//import "github.com/satori/go.uuid"
import "fmt"
import "strings"
import "sort"
import "time"
import "sync"
import "math"
import "github.com/bradfitz/slice"
import "math/rand"

var playerMax = 5

type PlayerChamp struct {
	player string
	champ  Champ
}

var bg1 map[string][]Champ
var bg2 map[string][]Champ
var bg3 map[string][]Champ

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
			team.champs = append(team.champs, roster[idx])
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
	champs []Champ
	score  float32
}

type PlayerDefenders struct {
	player    string
	defenders Defenders
}

func (d *Defenders) String() string {
	var ret []string
	for _, c := range d.champs {
		ret = append(ret, c.String())
	}
	ret = append(ret, fmt.Sprintf("%v", d.score))
	return strings.Join(ret, ",")
}

func copyDefenders(d Defenders) Defenders {
	ret := Defenders{score: d.score}
	for _, c := range d.champs {
		ret.champs = append(ret.champs, c)
	}
	return ret
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
    mastery := masteryMap[p]

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
				champs: append(d.champs, lockedChamps...),
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
			for _, champ := range d.champs {
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
					best = append(result, PlayerDefenders{player: p, defenders: ca})
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

  /*
	result, score, err := assignChamps(occupiedNodes, remainingChamps, skippedChamps)
	fmt.Printf("adding back locked champs %v\n", locked)
	for k, v := range locked {
		result[k] = v
	}
	return result, score, err
  */
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

func run(bg map[string][]Champ) {
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
					fmt.Printf("%s: %s\n", pd.player, pd.defenders.String())
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
		for _, c := range pd.defenders.champs {
			allChamps = Insert(allChamps, c)
		}
		fmt.Printf("%s: %s\n", pd.player, pd.defenders.String())
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
		output = append(output, fmt.Sprintf("%s: ", pd.player))
		for _, c := range pd.defenders.champs {
			output = append(output, c.Champ.String())
			output = append(output, fmt.Sprintf("(%v) ", result[c]))
		}
		fmt.Printf("%v\n", strings.Join(output, ""))
	}
}

type masteries struct {
  Suicides bool
  MD int
}

var masteryMap = map[string]masteries {
  "sugar": masteries{false, 5},
  "TomJenks": masteries{false, 3},
  "Nino": masteries{true, 0},
  "Easy": masteries{true, 0},
  "dhdhqqq": masteries{false, 3},
  "Spickster": masteries{true, 0},
  "Emodiva": masteries{false, 4},
  "LivingArtiface": masteries{false, 2},
  "Cantona": masteries{false, 0},
  "MaltLicker": masteries{false, 0},
  "Marjoriez": masteries{false, 3},
  "Webslinger": masteries{false, 2},
  "Mike-781": masteries{true, 0},
}

func main() {
	bg1 = map[string][]Champ{
		"sugar": []Champ{
			//NewChamp(Apocalypse, 6, 2, 0),
			Champ{NickFury, 6, 3, 1, 51},
			NewChamp(Void, 6, 3, 160),
			NewChamp(Dragonman, 6, 2, 20),
			NewChamp(CosmicGhostRider, 6, 3, 20),
			NewChamp(Terrax, 5, 5, 200),
			NewChamp(Mojo, 5, 5, 200),
			NewChamp(Guillotine2099, 6, 2, 0),
			NewChamp(Guardian, 5, 5, 200),
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(Punisher2099, 6, 2, 0),
		},
		"dhdhqqq": []Champ{
			Champ{Longshot, 6, 3, 60, 54},
			Champ{Killmonger, 5, 5, 200, 49},
			NewChamp(BlackWidowDeadlyOrigins, 6, 3, 0),
			NewChamp(Mojo, 5, 5, 89),
			NewChamp(Void, 5, 5, 200),
			NewChamp(EmmaFrost, 5, 5, 20),
			NewChamp(Sentinel, 5, 5, 20),
			NewChamp(Korg, 5, 5, 20),
			NewChamp(Mysterio, 5, 5, 20),
			NewChamp(Tigra, 6, 2, 0),
			NewChamp(Guillotine2099, 6, 2, 0),
			//NewChamp(NickFury, 5, 5, 1),
		},
		"TomJenks": []Champ{
			Champ{Thing, 6, 3, 200, 53},
			NewChamp(HitMonkey, 6, 3, 20),
			Champ{DoctorDoom, 6, 3, 20, 55},
			NewChamp(ProfessorX, 5, 5, 20),
			NewChamp(Sasquatch, 5, 5, 20),
			NewChamp(SpiderHam, 5, 5, 20),
			NewChamp(RedGuardian, 5, 5, 20),
			NewChamp(Apocalypse, 5, 5, 20),
			NewChamp(Guillotine2099, 5, 5, 20),
			NewChamp(MoleMan, 5, 5, 20),
			NewChamp(NickFury, 5, 5, 20),
			NewChamp(Iceman, 5, 5, 200),
			NewChamp(Havok, 5, 5, 20),
			NewChamp(Tigra, 5, 5, 20),
			NewChamp(Domino, 6, 2, 20),
		},
		"LivingArtiface": []Champ{
			NewChamp(SilverSurfer, 5, 5, 20),
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(SpiderManStealth, 5, 5, 20),
			NewChamp(CaptainMarvelMovie, 5, 5, 20),
			NewChamp(Namor, 5, 5, 20),
			NewChamp(Void, 5, 5, 20),
			NewChamp(ProfessorX, 5, 5, 20),
			NewChamp(Thing, 5, 5, 20),
			NewChamp(SymbioteSupreme, 5, 5, 20),
			NewChamp(Havok, 5, 5, 20),
			NewChamp(NickFury, 5, 5, 20),
			NewChamp(Magneto, 5, 5, 20),
			NewChamp(Warlock, 5, 5, 20),
			NewChamp(CosmicGhostRider, 6, 2, 0),
			NewChamp(Venom, 5, 5, 20),
			NewChamp(Apocalypse, 6, 2, 0),
			NewChamp(Guillotine2099, 6, 3, 0),
			NewChamp(Hyperion, 5, 5, 20),
			NewChamp(Dragonman, 5, 5, 0),
		},
		"Yves": []Champ{
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(CaptainMarvelMovie, 5, 5, 200),
			NewChamp(SymbioteSupreme, 5, 5, 200),
			NewChamp(Aegon, 5, 5, 200),
			NewChamp(RedHulk, 5, 5, 200),
			NewChamp(BlackWidowDeadlyOrigins, 6, 3, 0),
			NewChamp(Sentinel, 5, 5, 20),
			NewChamp(Sasquatch, 6, 2, 20),
			NewChamp(Medusa, 5, 5, 20),
		},
		"Nino": []Champ{
			NewChamp(BlackWidowClaireVoyant, 5, 5, 20),
			NewChamp(Thing, 5, 5, 200),
			NewChamp(Aegon, 5, 5, 200),
			NewChamp(OmegaRed, 5, 5, 200),
			NewChamp(Medusa, 5, 5, 20),
			NewChamp(Tigra, 5, 5, 20),
			NewChamp(Havok, 5, 5, 20),
			NewChamp(SpiderHam, 5, 5, 20),
			NewChamp(EmmaFrost, 5, 5, 20),
			NewChamp(Warlock, 5, 5, 20),
		},
		"Timzo": []Champ{
			NewChamp(Aegon, 5, 5, 20),
			NewChamp(Void, 5, 5, 200),
			NewChamp(Guillotine2099, 5, 5, 0),
			NewChamp(Magneto, 5, 5, 0),
			NewChamp(Hyperion, 5, 5, 0),
			NewChamp(Venom, 5, 5, 0),
			NewChamp(Sentinel, 5, 4, 20),
			NewChamp(IronManInfinityWar, 5, 4, 20),
			NewChamp(Annihilus, 5, 4, 20),
			NewChamp(HitMonkey, 5, 4, 20),
		},
		"Cantona": []Champ{
			NewChamp(Mojo, 5, 5, 20),
			NewChamp(Korg, 5, 5, 20),
			NewChamp(Void, 6, 2, 20),
			NewChamp(EmmaFrost, 5, 5, 20),
			NewChamp(MisterSinister, 6, 2, 20),
			NewChamp(SymbioteSupreme, 5, 5, 20),
			NewChamp(Apocalypse, 6, 2, 0),
			NewChamp(Quake, 5, 5, 0),
			NewChamp(SpiderGwen, 5, 5, 20),
		},
		"MaltLicker": []Champ{
			NewChamp(NickFury, 5, 5, 20),
			NewChamp(OmegaRed, 5, 5, 20),
			NewChamp(Domino, 5, 5, 20),
			NewChamp(Warlock, 5, 5, 20),
			NewChamp(Medusa, 5, 5, 20),
			NewChamp(SymbioteSupreme, 5, 5, 20),
			NewChamp(BlackWidowClaireVoyant, 5, 5, 0),
			NewChamp(Korg, 5, 4, 20),
			NewChamp(EbonyMaw, 5, 4, 20),
			NewChamp(Annihilus, 5, 4, 20),
			NewChamp(Darkhawk, 5, 4, 20),
			NewChamp(Iceman, 5, 4, 20),
		},
		"Spickster": []Champ{
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(Mephisto, 5, 5, 20),
			NewChamp(IronManInfinityWar, 5, 5, 20),
			NewChamp(Morningstar, 5, 5, 20),
			NewChamp(Quake, 5, 5, 20),
			NewChamp(Sentinel, 5, 5, 20),
			NewChamp(Magneto, 5, 5, 20),
			NewChamp(CosmicGhostRider, 5, 5, 20),
			NewChamp(Iceman, 5, 5, 20),
			NewChamp(HumanTorch, 5, 5, 0),
			NewChamp(Phoenix, 6, 1, 20),
		},
	}

	bg3 = map[string][]Champ{
		"Easy": []Champ{
			NewChamp(Void, 6, 3, 200),
			NewChamp(DoctorDoom, 6, 3, 200),
			NewChamp(NickFury, 6, 3, 20),
			NewChamp(Warlock, 5, 5, 200),
			//NewChamp(HumanTorch, 6, 3, 0),
			NewChamp(Magneto, 5, 5, 20),
			NewChamp(Terrax, 6, 2, 20),
			NewChamp(SpiderHam, 5, 5, 20),
			NewChamp(Domino, 5, 5, 20),
			NewChamp(IronManInfinityWar, 5, 5, 20),
			NewChamp(ProfessorX, 5, 5, 20),
			NewChamp(Hyperion, 5, 5, 20),
			NewChamp(Apocalypse, 6, 2, 0),
			NewChamp(Dragonman, 6, 2, 0),
			NewChamp(Havok, 6, 2, 0),
		},
		"Wayne": []Champ{
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(Thing, 5, 5, 20),
			NewChamp(Guillotine2099, 5, 5, 20),
			NewChamp(Korg, 5, 5, 20),
			NewChamp(Quake, 5, 5, 20),
			NewChamp(Void, 5, 5, 200),
			NewChamp(Magik, 5, 5, 20),
			NewChamp(Iceman, 5, 5, 20),
			NewChamp(Hyperion, 5, 5, 20),
			NewChamp(HumanTorch, 5, 5, 20),
			/*
			   NewChamp(Morningstar, 5, 4, 20),
			   NewChamp(Mephisto, 6, 1, 20),*/
		},
		"MarjorieZ": []Champ{
			NewChamp(Hyperion, 5, 5, 20),
			NewChamp(Void, 5, 5, 20),
			NewChamp(Medusa, 5, 5, 20),
			NewChamp(Sentinel, 5, 5, 20),
			NewChamp(Venom, 5, 5, 20),
			NewChamp(ArchAngel, 5, 5, 20),
			NewChamp(EmmaFrost, 5, 4, 20),
			NewChamp(Namor, 5, 4, 20),
			NewChamp(DoctorVoodoo, 5, 4, 20),
			NewChamp(Falcon, 5, 5, 20),
		},
		"Emodiva": []Champ{
			NewChamp(Guillotine2099, 6, 3, 20),
			NewChamp(NickFury, 6, 3, 20),
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(SpiderManStealth, 5, 5, 20),
			NewChamp(CaptainMarvelMovie, 5, 5, 20),
			NewChamp(Warlock, 5, 5, 20),
			NewChamp(Darkhawk, 5, 5, 200),
			NewChamp(Magneto, 5, 5, 20),
			NewChamp(Sasquatch, 5, 5, 20),
			NewChamp(Guardian, 5, 5, 20),
			NewChamp(Killmonger, 5, 5, 20),
			NewChamp(MisterFantastic, 6, 3, 20),
		},
		"Aaron": []Champ{
			NewChamp(Sunspot, 6, 3, 20),
			NewChamp(Mysterio, 5, 5, 20),
			NewChamp(Aegon, 5, 5, 200),
			NewChamp(Crossbones, 5, 5, 20),
			NewChamp(AbominationImmortal, 5, 5, 20),
			NewChamp(Terrax, 5, 5, 20),
			NewChamp(ElsaBloodstone, 5, 5, 20),
			NewChamp(HitMonkey, 5, 5, 20),
			NewChamp(Nightcrawler, 6, 2, 20),
		},
		"Mike-781": []Champ{
			NewChamp(SilverSurfer, 5, 5, 20),
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(ProfessorX, 5, 5, 20),
			NewChamp(Mojo, 5, 5, 20),
			NewChamp(Thing, 5, 5, 20),
			NewChamp(ManThing, 5, 5, 20),
			NewChamp(Mysterio, 5, 5, 20),
			NewChamp(Havok, 5, 5, 20),
			NewChamp(Killmonger, 5, 5, 20),
			NewChamp(Magneto, 5, 5, 20),
			NewChamp(Medusa, 5, 5, 20),
			NewChamp(Korg, 6, 2, 20),
			NewChamp(InvisibleWoman, 6, 2, 20),
		},
		"Spliffy": []Champ{
			NewChamp(Thing, 6, 3, 20),
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(Apocalypse, 5, 5, 20),
			NewChamp(Havok, 5, 5, 20),
			NewChamp(Medusa, 5, 5, 20),
			NewChamp(SpidermanStark, 5, 5, 20),
			NewChamp(Sentinel, 5, 5, 20),
			NewChamp(Longshot, 5, 5, 20),
			NewChamp(ProfessorX, 5, 5, 20),
			NewChamp(BlackWidowDeadlyOrigins, 6, 2, 0),
			//NewChamp(Punisher2099, 6, 1, 20),
			NewChamp(Warlock, 5, 5, 20),
			NewChamp(SymbioteSupreme, 5, 5, 20),
			NewChamp(Domino, 5, 5, 20),
		},
		"WebSlinger": []Champ{
			NewChamp(Void, 5, 5, 200),
			NewChamp(Hyperion, 5, 5, 20),
			NewChamp(Domino, 5, 5, 20),
			NewChamp(HumanTorch, 5, 5, 20),
			NewChamp(SpiderManStealth, 5, 5, 20),
			//NewChamp(NickFury, 5, 5, 20),
			/*
			   NewChamp(Hulkbuster, 5, 4, 20),
			   NewChamp(SpiderGwen, 5, 4, 20),
			   NewChamp(VisionAarkus, 6, 1, 0),*/
		},
		"Basher": []Champ{
			NewChamp(DoctorDoom, 5, 5, 20),
			NewChamp(SpiderManStealth, 5, 5, 20),
			NewChamp(NickFury, 5, 5, 20),
			NewChamp(Void, 5, 5, 20),
			NewChamp(OmegaRed, 5, 5, 20),
			NewChamp(Mysterio, 5, 5, 20),
			NewChamp(VenomTheDuck, 5, 5, 20),
			NewChamp(Hyperion, 5, 5, 20),
			NewChamp(Falcon, 5, 5, 20),
			NewChamp(ThorRagnarok, 6, 1, 20),
			NewChamp(Magik, 5, 5, 20),
		},
		"Wellsz": []Champ{
			NewChamp(Hyperion, 5, 5, 20),
			NewChamp(Colossus, 5, 5, 20),
			NewChamp(RedGuardian, 5, 4, 20),
			NewChamp(CullObsidian, 6, 1, 20),
			NewChamp(Nightcrawler, 6, 1, 20),
			NewChamp(Sabretooth, 6, 1, 0),
			NewChamp(Mordo, 6, 1, 0),
			NewChamp(MilesMorales, 6, 1, 0),
			NewChamp(OmegaRed, 5, 5, 20),
		},
	}

	run(bg3)
}
