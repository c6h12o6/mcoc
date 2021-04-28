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
}

var playerMax = 5

type PlayerChamp struct {
	player string
	champ  Champ
}

var bg1_2 = map[string][]Champ{
	"sugar": []Champ{
		Champ{Apocalypse, 6, 2, 0},
		Champ{NickFury, 6, 3, 1},
		Champ{Void, 6, 3, 160},
		Champ{Dragonman, 6, 2, 20},
		Champ{CosmicGhostRider, 6, 3, 20},
		Champ{Terrax, 5, 5, 200},
		Champ{Mojo, 5, 5, 200},
		Champ{Guillotine2099, 6, 2, 0},
		Champ{Guardian, 5, 5, 200},
		Champ{DoctorDoom, 5, 5, 20},
		Champ{Punisher2099, 6, 2, 0},
	},
	"dhdhqqq": []Champ{
		Champ{Longshot, 6, 3, 60},
		Champ{Killmonger, 5, 5, 200},
		Champ{BlackWidowDeadlyOrigins, 6, 3, 0},
		Champ{Mojo, 5, 5, 89},
		Champ{Void, 5, 5, 200},
		Champ{EmmaFrost, 5, 5, 20},
		Champ{Sentinel, 5, 5, 20},
		Champ{Korg, 5, 5, 20},
		Champ{Mysterio, 5, 5, 20},
		Champ{Tigra, 6, 2, 0},
		Champ{Guillotine2099, 6, 2, 0},
		Champ{NickFury, 5, 5, 1},
	},
	"TomJenks": []Champ{
		Champ{Thing, 6, 3, 200},
		Champ{HitMonkey, 6, 3, 20},
		Champ{DoctorDoom, 5, 5, 20},
		Champ{ProfessorX, 5, 5, 20},
		Champ{Sasquatch, 5, 5, 20},
		Champ{SpiderHam, 5, 5, 20},
		Champ{RedGuardian, 5, 5, 20},
		Champ{Apocalypse, 5, 5, 20},
		Champ{Guillotine2099, 5, 5, 20},
		Champ{MoleMan, 5, 5, 20},
		Champ{NickFury, 5, 5, 20},
		Champ{Iceman, 5, 5, 200},
		Champ{Havok, 5, 5, 20},
		Champ{Tigra, 5, 5, 20},
		Champ{Domino, 6, 2, 20},
	},
	"LivingArtiface": []Champ{
		Champ{SilverSurfer, 5, 5, 20},
		Champ{DoctorDoom, 5, 5, 20},
		Champ{SpiderManStealth, 5, 5, 20},
		Champ{CaptainMarvelMovie, 5, 5, 20},
		Champ{Namor, 5, 5, 20},
		Champ{Void, 5, 5, 20},
		Champ{ProfessorX, 5, 5, 20},
		Champ{Thing, 5, 5, 20},
		Champ{SymbioteSupreme, 5, 5, 20},
		Champ{Havok, 5, 5, 20},
		Champ{NickFury, 5, 5, 20},
		Champ{Magneto, 5, 5, 20},
		Champ{Warlock, 5, 5, 20},
		Champ{CosmicGhostRider, 6, 2, 0},
		Champ{Venom, 5, 5, 20},
		Champ{Apocalypse, 6, 2, 0},
		Champ{Guillotine2099, 6, 3, 0},
		Champ{Hyperion, 5, 5, 20},
		Champ{Dragonman, 5, 5, 0},
	},
	"Yves": []Champ{
		Champ{DoctorDoom, 5, 5, 20},
		Champ{CaptainMarvelMovie, 5, 5, 200},
		Champ{SymbioteSupreme, 5, 5, 200},
		Champ{Aegon, 5, 5, 200},
		Champ{RedHulk, 5, 5, 200},
		Champ{BlackWidowDeadlyOrigins, 6, 3, 0},
		Champ{Sentinel, 5, 5, 20},
		Champ{Sasquatch, 6, 2, 20},
		Champ{Medusa, 5, 5, 20},
	},
	"Nino": []Champ{
		Champ{BlackWidowClaireVoyant, 5, 5, 20},
		Champ{Thing, 5, 5, 200},
		Champ{Aegon, 5, 5, 200},
		Champ{OmegaRed, 5, 5, 200},
		Champ{Medusa, 5, 5, 20},
		Champ{Tigra, 5, 5, 20},
		Champ{Havok, 5, 5, 20},
		Champ{SpiderHam, 5, 5, 20},
		Champ{EmmaFrost, 5, 5, 20},
		Champ{Warlock, 5, 5, 20},
	},
	"Timzo": []Champ{
		Champ{Aegon, 5, 5, 20},
		Champ{Void, 5, 5, 200},
		Champ{Guillotine2099, 5, 5, 0},
		Champ{Magneto, 5, 5, 0},
		Champ{Hyperion, 5, 5, 0},
		Champ{Venom, 5, 5, 0},
		Champ{Sentinel, 5, 4, 20},
		Champ{IronManInfinityWar, 5, 4, 20},
		Champ{Annihilus, 5, 4, 20},
		Champ{HitMonkey, 5, 4, 20},
	},
	"Cantona": []Champ{
		Champ{Mojo, 5, 5, 20},
		Champ{Korg, 5, 5, 20},
		Champ{Void, 6, 2, 20},
		Champ{EmmaFrost, 5, 5, 20},
		Champ{MisterSinister, 6, 2, 20},
		Champ{SymbioteSupreme, 5, 5, 20},
		Champ{Apocalypse, 6, 2, 0},
		Champ{Quake, 5, 5, 0},
		Champ{SpiderGwen, 5, 5, 20},
	},
	"MaltLicker": []Champ{
		Champ{NickFury, 5, 5, 20},
		Champ{OmegaRed, 5, 5, 20},
		Champ{Domino, 5, 5, 20},
		Champ{Warlock, 5, 5, 20},
		Champ{Medusa, 5, 5, 20},
		Champ{SymbioteSupreme, 5, 5, 20},
		Champ{BlackWidowClaireVoyant, 5, 5, 0},
		Champ{Korg, 5, 4, 20},
		Champ{EbonyMaw, 5, 4, 20},
		Champ{Annihilus, 5, 4, 20},
		Champ{Darkhawk, 5, 4, 20},
		Champ{Iceman, 5, 4, 20},
	},
	"Spickster": []Champ{
		Champ{DoctorDoom, 5, 5, 20},
		Champ{Mephisto, 5, 5, 20},
		Champ{IronManInfinityWar, 5, 5, 20},
		Champ{Morningstar, 5, 5, 20},
		Champ{Quake, 5, 5, 20},
		Champ{Sentinel, 5, 5, 20},
		Champ{Magneto, 5, 5, 20},
		Champ{CosmicGhostRider, 5, 5, 20},
		Champ{Iceman, 5, 5, 20},
		Champ{HumanTorch, 5, 5, 0},
		Champ{Phoenix, 6, 1, 20},
	},
}

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
		for _, c := range playerChamps {
			if _, ok := diversity[c.champ]; ok {
				continue
			}

			reducedChamps = append(reducedChamps, c)
		}
		if len(reducedChamps) < playerMax {
			recordMemo2(memoKey, nil, 0, fmt.Errorf("No valid teams"), callArgs)
			ch <- memoItem2{pds: nil, score: 0, err: fmt.Errorf("No valid teams")}
			return nil, 0, fmt.Errorf("No valid teams")
		}
		combos := teamCombinations(playerMax, reducedChamps)
		combos = combos[:int(math.Min(float64(len(combos)), 5))]
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

func assignChamps(occupiedNodes map[int]Champ, remainingChamps []Champ) (map[Champ]int, int, error) {
  if len(remainingChamps) == 0 {
    return map[Champ]int{}, 0, nil
  }

  bestScore := 0
  bestMap := map[Champ]int{}

  assigned := false
  c := remainingChamps[0]
  for n := 55; n > 0; n-- {
    if _, ok := occupiedNodes[n]; ok {
      continue
    }
    idx := sort.Search(len(Nodes[n]), func (i int) bool { if Nodes[n][i] == c.champ { return true } else { return false }})
    if idx == len(Nodes[n]) {
      //fmt.Printf("Could not find %v in node %v\n", c, n)
    } else {
      assigned = true
      newOccupiedNodes := copyOccupiedNodes(occupiedNodes)
      newOccupiedNodes[n] = c
      //fmt.Printf("Found %v at node %v index %v\n", c, n, idx)
      result, score, err := assignChamps(newOccupiedNodes, remainingChamps[1:])
      if err == nil && score + 1 > bestScore {
        bestScore = score + 1
        bestMap = result
        bestMap[c] = n
      }
    }
  }
  if !assigned {
    return assignChamps(occupiedNodes, remainingChamps[1:])
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


func main() {

	ch := make(chan memoItem2)
	t := time.Now()
	//combos := teamCombinations(5, bg1_2["sugar"], "sugar")
	//result, score, err := findBestBG(ch, map[HeroVal]bool{}, bg1_2, []string{"sugar", "dhdhqqq", "TomJenks", "LivingArtiface"})
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
		go findBestBG(ch, map[HeroVal]bool{}, bg1_2, playerList, Defenders{})
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

  result, _, _ := assignChamps(map[int]Champ{}, allChamps)
  fmt.Printf("result length: %v\n", len(result))
  var maplines []string
  for c, n := range result {
    maplines = append(maplines, fmt.Sprintf("%02d: %v", n, c))
  }
  sort.Strings(maplines)
  fmt.Print(strings.Join(maplines, "\n"))
  var unplacedChamps []string
  for _, c := range allChamps {
    _, ok := result[c]
    if !ok {
      unplacedChamps = append(unplacedChamps, c.String())
    }
  }
  fmt.Printf("\n"))
  fmt.Print(strings.Join(unplacedChamps, ","))


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
