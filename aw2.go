package main

import . "../mcoc/globals"

//import "github.com/satori/go.uuid"
import "fmt"
import "strings"
import "sort"
import "time"
import "sync"
import "math"
import "github.com/bradfitz/slice"
import "math/rand"

type Champ struct {
	champ HeroVal
	stars int32
	level int32
	sig   int32
  lockedNode int
}

func NewChamp(hv HeroVal, stars int32, level int32, sig int32) Champ {
  return Champ{hv, stars, level, sig, 0}
}

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
func teamCombinations(teamsize int, roster []Champ) []Defenders {
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
			score += champScore(roster[idx])
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

func champValue(a Champ) float32 {
	if a.stars == 6 && a.level == 3 {
		return 15 + float32(a.sig)/200
	} else if a.stars == 6 && a.level == 2 {
		return 9 + float32(a.sig)/200
	} else if a.stars == 5 && a.level == 5 {
		return 8 + float32(a.sig)/200
	} else if a.stars == 6 && a.level == 1 {
		return 6 + float32(a.sig)/200
	} else if a.stars == 5 && a.level == 4 {
		return 5 + float32(a.sig)/200
	} else {
		return 1
	}
}

func champScore(a Champ) float32 {
	if a.champ == Empty {
		return 0.0
	}
	return champValue(a) + float32(GetLevel(a.champ))
}

func (c Champ) String() string {
  return fmt.Sprintf("%v (%v/%v)", c.champ.String(), c.stars, c.level)
}
var memoLock sync.Mutex
var memoCount int
var totalCount int
var totalCalls int
var tryingCount int
var trying2Count int

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
	//fmt.Printf("%v\n", key)
	return key
}

func recordMemo2(memoKey string, pds []PlayerDefenders, score float32, err error, callArgs Defenders) {
	memoLock.Lock()
	memo2[memoKey] = memoItem2{pds: pds, score: score, err: err, callArgs: callArgs}
	memoLock.Unlock()
}

//var first = map[string]bool{}

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
		if memoCount%100 == 0 {
			//fmt.Printf(".")
		}

		ch <- memoItem2{pds: mi.pds, score: mi.score, err: mi.err, callArgs: callArgs}
		return mi.pds, mi.score, mi.err
	}

	if len(players) != 0 {
		p := players[0]
		/*
		   if players[len(players) - 1] == "Cantona" {
		     fmt.Printf("player %v\n", p)
		   }
		*/
		playerChamps := roster[p]
		//fmt.Printf("playerChamps: %v", playerChamps)
		var reducedChamps []Champ
    var lockedChamps []Champ
    var lockedChampScore float32

    var debug bool
    //if p == "sugar" { debug = true }

		for _, c := range playerChamps {
      // all locked nodes should also already be in the diversity map
      if c.lockedNode != 0 {
        lockedChamps = append(lockedChamps, c)
        lockedChampScore += champScore(c)
      }
			if _, ok := diversity[c.champ]; ok {
				continue
			}
      
			reducedChamps = append(reducedChamps, c)
		}
		if len(reducedChamps) < playerMax {
			recordMemo2(memoKey, nil, 0, fmt.Errorf("No valid teams"), callArgs)
			ch <- memoItem2{pds: nil, score: 0, err: fmt.Errorf("No valid teams")}
      //fmt.Printf("Failing due to %v %v %v\n", p, reducedChamps, diversity)
			return nil, 0, fmt.Errorf("No valid teams")
		}
		combos := teamCombinations(playerMax-len(lockedChamps), reducedChamps)
		combos = combos[:int(math.Min(float64(len(combos)), 5))]
    newCombos := make([]Defenders, len(combos))

    for ii, d := range combos {
      d.score += lockedChampScore
      newCombos[ii] = Defenders{
        champs: append(d.champs, lockedChamps...),
        score: d.score + lockedChampScore,
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
				newDiversity[champ.champ] = true
			}
			calls++
			go findBestBG(newCh, newDiversity, roster, players[1:], d)
		}
		for ; calls > 0; calls-- {
			select {
			case mi := <-newCh:
				result, newScore, err, ca := mi.pds, mi.score, mi.err, mi.callArgs
				//fmt.Printf("result: %v ca: %v\n", result, ca)
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
	//fmt.Printf("returning at end\n")
	return best, bestScore, nil
}

func findBestBGHelper(ch chan memoItem2, diversity map[HeroVal]bool, roster map[string][]Champ, players []string, callArgs Defenders) ([]PlayerDefenders, float32, error) {
  for _, p := range players {
    for _, c := range roster[p] {
      if c.lockedNode != 0 {
        diversity[c.champ] = true
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
  for k,v := range o {
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
  debug := false
  if c.champ == CosmicGhostRider {
    //debug = true
  }
  if champValue(c) < 7 {
    skippedChamps = append(skippedChamps, c)
    return assignChamps(occupiedNodes, remainingChamps[1:], skippedChamps)
  }
  if debug {fmt.Printf("%v\n", c)}
  var count int
  for n := 55; n > 0; n-- {
    if xx, ok := occupiedNodes[n]; ok {
      if debug {fmt.Printf("skipping %v %v\n", n, xx)}
      continue
    }
    //idx := sort.Search(len(Nodes[n]), func (i int) bool { if Nodes[n][i] == c.champ { return true } else { return false }})

    if Contains(Nodes[n], c.champ) {
      count++
      assigned = true
      newOccupiedNodes := copyOccupiedNodes(occupiedNodes)
      newOccupiedNodes[n] = c
      //fmt.Printf("Found %v at node %v index %v\n", c, n, idx)
      result, score, err := assignChamps(newOccupiedNodes, remainingChamps[1:], skippedChamps)
      if err == nil && score + 1 > bestScore {
        bestScore = score + 1
        bestMap = result
        bestMap[c] = n
      }
    } else {
      if debug {fmt.Printf("Could not find %v in node %v\n", c, n)}
    }
    if count > 1 {
      if debug {fmt.Printf("breaking here=====================\n")}
      //fmt.Printf("%v\n", len(remainingChamps))
      break
    }
  }
  if !assigned {
    //fmt.Printf("Skipping assignment of %v (%v)\n", c, len(remainingChamps))
    skippedChamps = append(skippedChamps, c)
    return assignChamps(occupiedNodes, remainingChamps[1:], skippedChamps)
  }


/*
  node := Nodes[startNode]
  for idx, h := range node {
    if h == MaxHeroVal {
    }

    remaining := removeFromRemainingChamps(remainingChamps, h)
    result, score, err := assignNode(startNode - 1, remaining)
    
    if err != nil && score > bestScore {
      bestScore = score + 1
      result[n] = h
      bestMap = result
    }
  }
  
  */
  return bestMap, bestScore, nil
}


func assignChampsHelper(bg map[string][]Champ, occupiedNodes map[int]Champ, remainingChamps []Champ, skippedChamps []Champ) (map[Champ]int, int, error) {
  locked := map[Champ]int{}
  for _, champlist := range bg {
    for _, c := range champlist {
      if c.lockedNode != 0 {
        occupiedNodes[c.lockedNode] = c
        locked[c] = c.lockedNode
      }
    }
  }

  result, score, err := assignChamps(occupiedNodes, remainingChamps, skippedChamps)
  fmt.Printf("adding back locked champs %v\n", locked)
  for k,v := range locked {
    result[k] = v
  }
  return result, score, err
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
  i := sort.Search(len(sorted), func(i int) bool { return champScore(sorted[i]) < champScore(champ)})
  sorted = append(sorted, Champ{})
  copy(sorted[i+1:], sorted[i:])
  sorted[i] = champ
  return sorted
}


func run(bg map[string][]Champ) {
	ch := make(chan memoItem2)
	t := time.Now()
	//combos := teamCombinations(5, bg1_2["sugar"], "sugar")
	//result, score, err := findBestBG(ch, map[HeroVal]bool{}, bg1_2, []string{"sugar", "dhdhqqq", "TomJenks", "LivingArtiface"})
    /*
	players := []string{
		"sugar",
		"dhdhqqq",
		"TomJenks",
		"LivingArtiface",
		"Yves",
		"Nino",
		"Timzo",
		"Cantona",
		"Spickster",
		"MaltLicker"}
  players := []string{
    "Easy",
    "Wayne",
    "MarjorieZ",
    "Emodiva",
    "Aaron",
    "Mike-781",
    "Spliffy",
    "WebSlinger",
    "Basher",
    "Wellsz",
  }
  */
  var players []string
  for p,_ := range bg {
    players = append(players, p)
  }

	playerPermutations := permutations(players)
	permD := time.Now().Sub(t)
	fmt.Printf("%v\n", permD)

	var bestResult []PlayerDefenders
	var bestScore float32
	for n := 0; n < 100; n++ {
		memo2 = map[string]memoItem2{}
		playerList := playerPermutations[rand.Intn(len(playerPermutations))]
		fmt.Printf("\tTrying %v\n", playerList)

		t = time.Now()
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

  result, _, _ := assignChampsHelper(bg1, map[int]Champ{}, allChamps, []Champ{})
  fmt.Printf("result length: %v\n", len(result))
  var maplines []string
  assigned := map[int]Champ{}
  for c, n := range result {
    assigned[n] = c
    maplines = append(maplines, fmt.Sprintf("%02d: %v", n, c))
  }
  sort.Strings(maplines)
  fmt.Print(strings.Join(maplines, "\n"))

    var unplacedChampsStrings  []string
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
        output = append(output, c.champ.String())
        output = append(output, fmt.Sprintf("(%v) ", result[c]))
      }
      fmt.Printf("%v\n", strings.Join(output, ""))
    }

    /*
    result2, _, _ := assignChamps(assigned, unplacedChamps)
    fmt.Printf("result length: %v\n", len(result))
    var maplines2 []string
    assigned := map[int]Champ{}
    for c, n := range result2 {
      assigned[n] = c
      maplines2= append(maplines2, fmt.Sprintf("--%02d: %v", n, c))
    }
    sort.Strings(maplines2)
    fmt.Print(strings.Join(maplines2, "\n"))
  }
  */


	/*
	  var players []string
	  for p, _ := range bg1_2 {
	    players = append(players, p)
	  }
	  fmt.Printf("players: %v\n", players)
	  players = []string{"sugar", "dhdhqqq", "TomJenks", "LivingArtiface", "Yves", "Nino"}

	  ch := make(chan BGScore)
	  t := time.Now()
	  battleGroup := map[string]Defenders{}
	  //score := findBestRoster(ch, battleGroup, players)
	  var score BGScore
	  go findBestRoster(ch, battleGroup, players)
	  select {
	  case result := <-ch:
	    score = result
	  }

	  //bestNodes, score := result.nodes, result.score
	  d := time.Now().Sub(t)

	  //PrintNodes(bestNodes)

	  fmt.Printf("%v\n", score)
	  fmt.Printf("Took %v\n", d)
	  //fmt.Printf("Trying %v Trying 2 %v\n", tryingCount, trying2Count)
	*/
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
		NewChamp(DoctorDoom, 5, 5, 20),
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
		Champ{DoctorDoom, 5, 5, 20, 55},
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
    NewChamp(SilverSurfer, 5,5, 20),
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
    NewChamp(NickFury, 5, 5, 20),
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
    NewChamp(Colossus, 5,5 , 20),
    NewChamp(RedGuardian, 5, 4, 20),
    NewChamp(CullObsidian, 6, 1, 20),
    NewChamp(Nightcrawler, 6, 1, 20),
    NewChamp(Sabretooth, 6, 1, 0),
    NewChamp(Mordo, 6, 1, 0),
    NewChamp(MilesMorales, 6, 1, 0),
    NewChamp(OmegaRed, 5, 5, 20),
  },
}


  run(bg1)
}
