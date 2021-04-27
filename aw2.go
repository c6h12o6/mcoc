package main

import . "../mcoc/globals"
//import "github.com/satori/go.uuid"
import "fmt"
import "strings"
import "sort"
import "time"
import "sync"
import "runtime"
import "github.com/bradfitz/slice" 

type Champ struct {
  champ HeroVal
  stars int32
  level int32
  sig   int32
}

var playerMax = 5

type PlayerChamp struct {
  player string
  champ Champ
}

var bg1_2 = map[string][]Champ {
  "sugar": []Champ{
    Champ{Apocalypse, 6, 2, 0},
    Champ{NickFury, 6, 3, 1},
    Champ{Void, 6, 3, 160},
    Champ{Dragonman, 6, 2, 20},
    Champ{CosmicGhostRider, 6, 3, 20},
    Champ{Terrax, 5, 5, 200},
    Champ{Mojo, 5, 5, 200},
    //Champ{Guillotine2099, 6, 2, 0},
    //Champ{Guardian, 5, 5, 200},
  },
  "dhdhqqq": []Champ{
    Champ{Longshot, 6, 3, 60},
    Champ{Killmonger, 5, 5, 200},
    Champ{BlackWidowDeadlyOrigins, 6, 3, 0},
    Champ{Mojo, 5, 5, 89},
    //Void: Champ{Void, 5, 5, 200},
    Champ{EmmaFrost, 5, 5, 20},
    Champ{Sentinel, 5, 5, 20},
    Champ{Korg, 5, 5, 20},
    Champ{Mysterio, 5, 5, 20},
    Champ{Tigra, 6, 2, 0},
    Champ{Guillotine2099, 6, 2, 0},
    //NickFury: Champ{NickFury, 5, 5, 1},
  }, 
  "TomJenks": []Champ {
    Champ{Thing, 6, 3, 200},
    Champ{HitMonkey, 6, 3, 20},
    Champ{DoctorDoom, 5, 5, 20},
    Champ{ProfessorX, 5, 5, 20},
    Champ{Sasquatch, 5, 5, 20},
    Champ{SpiderHam, 5, 5, 20},
    //Champ{RedGuardian, 5, 5, 20},
    //Champ{Apocalypse, 5, 5, 20},
    Champ{Guillotine2099, 5, 5, 20},
    Champ{MoleMan, 5, 5, 20},
    //NickFury: Champ{NickFury, 5, 5, 20},
    Champ{Iceman, 5, 5, 200},
    //Havok: Champ{Havok,5, 5, 20},
    Champ{Tigra, 5, 5, 20},
    Champ{Domino, 6, 2, 20},
  },
  "LivingArtiface": []Champ{
    Champ{SilverSurfer, 5, 5, 20},
    Champ{DoctorDoom, 5, 5, 20},
    //Champ{SpiderManStealth, 5, 5, 20},
    //Champ{CaptainMarvelMovie, 5, 5, 20},
    //Namor: Champ{Namor, 5, 5, 20},
    Champ{Void, 5, 5, 20},
    Champ{ProfessorX, 5, 5, 20},
    //Thing: Champ{Thing, 5, 5, 20},
    Champ{SymbioteSupreme, 5, 5, 20},
    Champ{Havok, 5, 5, 20},
    //NickFury: Champ{NickFury, 5, 5, 20},
    Champ{Magneto, 5, 5, 20},
    Champ{Warlock, 5, 5, 20},
    Champ{CosmicGhostRider, 6, 2, 0},
    Champ{Venom, 5, 5, 20},
    //Apocalypse: Champ{Apocalypse, 6, 2, 0},
    Champ{Guillotine2099, 6, 3, 0},
    Champ{Hyperion, 5, 5, 20},
    //Dragonman: Champ{Dragonman, 5, 5, 0},
  }, 
}

var bg1 = map[string]map[HeroVal]Champ {
  "sugar": map[HeroVal]Champ{
    Apocalypse: Champ{Apocalypse, 6, 2, 0},
    NickFury: Champ{NickFury, 6, 3, 1},
    Void: Champ{Void, 6, 3, 160},
    Dragonman: Champ{Dragonman, 6, 2, 20},
    CosmicGhostRider: Champ{CosmicGhostRider, 6, 3, 20},
    Terrax: Champ{Terrax, 5, 5, 200},
    Mojo: Champ{Mojo, 5, 5, 200},
    Guillotine2099: Champ{Guillotine2099, 6, 2, 0},
    Guardian: Champ{Guardian, 5, 5, 200},
  }, 
  "dhdhqqq": map[HeroVal]Champ{
    Longshot: Champ{Longshot, 6, 3, 60},
    Killmonger: Champ{Killmonger, 5, 5, 200},
    BlackWidowDeadlyOrigins: Champ{BlackWidowDeadlyOrigins, 6, 3, 0},
    Mojo: Champ{Mojo, 5, 5, 89},
    //Void: Champ{Void, 5, 5, 200},
    EmmaFrost: Champ{EmmaFrost, 5, 5, 20},
    Sentinel: Champ{Sentinel, 5, 5, 20},
    Korg: Champ{Korg, 5, 5, 20},
    Mysterio: Champ{Mysterio, 5, 5, 20},
    Tigra: Champ{Tigra, 6, 2, 0},
    Guillotine2099: Champ{Guillotine2099, 6, 2, 0},
    //NickFury: Champ{NickFury, 5, 5, 1},
  }, /*
  "TomJenks": map[HeroVal]Champ{
    Thing: Champ{Thing, 6, 3, 200},
    HitMonkey: Champ{HitMonkey, 6, 3, 20},
    DoctorDoom: Champ{DoctorDoom, 5, 5, 20},
    ProfessorX: Champ{ProfessorX, 5, 5, 20},
    Sasquatch: Champ{Sasquatch, 5, 5, 20},
    SpiderHam: Champ{SpiderHam, 5, 5, 20},
    //RedGuardian: Champ{RedGuardian, 5, 5, 20},
    //Apocalypse: Champ{Apocalypse, 5, 5, 20},
    Guillotine2099: Champ{Guillotine2099, 5, 5, 20},
    MoleMan: Champ{MoleMan, 5, 5, 20},
    //NickFury: Champ{NickFury, 5, 5, 20},
    Iceman: Champ{Iceman, 5, 5, 200},
    //Havok: Champ{Havok,5, 5, 20},
    Tigra: Champ{Tigra, 5, 5, 20},
    Domino: Champ{Domino, 6, 2, 20},
  },
  "LivingArtiface": map[HeroVal]Champ{
    SilverSurfer: Champ{SilverSurfer, 5, 5, 20},
    DoctorDoom: Champ{DoctorDoom, 5, 5, 20},
    SpiderManStealth: Champ{SpiderManStealth, 5, 5, 20},
    CaptainMarvelMovie: Champ{CaptainMarvelMovie, 5, 5, 20},
    //Namor: Champ{Namor, 5, 5, 20},
    Void: Champ{Void, 5, 5, 20},
    ProfessorX: Champ{ProfessorX, 5, 5, 20},
    //Thing: Champ{Thing, 5, 5, 20},
    SymbioteSupreme: Champ{SymbioteSupreme, 5, 5, 20},
    Havok: Champ{Havok, 5, 5, 20},
    //NickFury: Champ{NickFury, 5, 5, 20},
    Magneto: Champ{Magneto, 5, 5, 20},
    Warlock: Champ{Warlock, 5, 5, 20},
    CosmicGhostRider: Champ{CosmicGhostRider, 6, 2, 0},
    Venom: Champ{Venom, 5, 5, 20},
    //Apocalypse: Champ{Apocalypse, 6, 2, 0},
    Guillotine2099: Champ{Guillotine2099, 6, 3, 0},
    Hyperion: Champ{Hyperion, 5, 5, 20},
    //Dragonman: Champ{Dragonman, 5, 5, 0},
  },
  "Yves": map[HeroVal]Champ{
    DoctorDoom: Champ{DoctorDoom, 5, 5, 20},
    CaptainMarvelMovie: Champ{CaptainMarvelMovie, 5, 5, 200},
    SymbioteSupreme: Champ{SymbioteSupreme, 5, 5, 200}, 
    Aegon: Champ{Aegon, 5, 5, 200},
    RedHulk: Champ{RedHulk, 5, 5, 200},
    BlackWidowDeadlyOrigins: Champ{BlackWidowDeadlyOrigins, 6, 3, 0},
    Sentinel: Champ{Sentinel, 5, 5, 20},
    Sasquatch: Champ{Sasquatch, 6, 2, 20},
    Medusa: Champ{Medusa, 5, 5, 20},
  },
  "Nino": map[HeroVal]Champ {
    BlackWidowClaireVoyant: Champ{BlackWidowClaireVoyant, 5, 5, 20},
    //Thing: Champ{Thing, 5, 5, 200},
    Aegon: Champ{Aegon, 5, 5, 200},
    OmegaRed: Champ{OmegaRed, 5, 5, 200},
    Medusa: Champ{Medusa, 5, 5, 20},
    Tigra: Champ{Tigra, 5, 5, 20},
    Havok: Champ{Havok, 5, 5, 20},
    SpiderHam: Champ{SpiderHam, 5, 5, 20},
    EmmaFrost: Champ{EmmaFrost, 5, 5, 20},
    Warlock: Champ{Warlock, 5, 5, 20},
  },*/

}

 // combinations is a helper function for creating all possible combinations of
 // values from "iterable" in groups of "r"
 func combinations(iterable []int, r int) [][]int {                                             
     var ret [][]int
 
     pool := iterable
     n := len(pool)
 
     if r > n {
         return [][]int {}
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
    return 10 + float32(a.sig)/200
  } else if a.stars == 6 && a.level == 2 {
    return 9 + float32(a.sig)/200
  } else if a.stars == 5 && a.level == 5 {
    return 8 + float32(a.sig)/200
  } else if a.stars == 6 && a.level == 1 {
    return 9 + float32(a.sig)/200
  } else if a.stars == 5 && a.level == 4 {
    return 6 + float32(a.sig)/200
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

func (a Champ) GreaterThan(b Champ) bool {
  if b.champ == MaxHeroVal {
    return true
  }

  delta := champValue(a) - champValue(b)
  if delta > 0 {
    return true
  }
  return false
}

func (pc PlayerChamp) PrettyPrint() string {
  if pc.champ.champ == MaxHeroVal {
    return ""
  }

  return fmt.Sprintf("%s %s", pc.player, pc.champ.champ)
}

func (c Champ) String() string {
  return fmt.Sprintf("%s (%v/%v)", c.champ.String(), c.stars, c.level)
}

func FormatNodes(nodes []PlayerChamp) string {
  s := make([]string, len(nodes))
  for n := 0; n < len(nodes); n++ {
    s[n] = fmt.Sprintf("%v:\t%s", n+1, nodes[n].PrettyPrint())
  }
  return strings.Join(s, "\n")
}

func PrintNodes(nodes []PlayerChamp) {
  fmt.Printf(FormatNodes(nodes) + "\n")
}

func FormatBattleGroup(bg map[string]Defenders) string {
  var entries []string
  for p, d := range bg {
    var subentry []string
    for _, h := range d.champs {
      subentry = append(subentry, h.String())
    }
    entries = append(entries, fmt.Sprintf("%v: %v", p, strings.Join(subentry, ",")))
  }
  return strings.Join(entries, "\n")
}

func PrintBattleGroup(bg map[string]Defenders) {
  fmt.Printf(FormatBattleGroup(bg) + "\n")
}

func (pc PlayerChamp) String() string {
  return pc.PrettyPrint()
}

func copyNodes(nodes []PlayerChamp) []PlayerChamp {
  ret := make([]PlayerChamp, len(nodes))
  copy(ret, nodes)
  return ret
}

type memoItem  struct {
  score float32
  nodes []PlayerChamp
}

var memoLock sync.Mutex
var memo = map[string]memoItem{}
var memoCount int
var totalCount int
var totalCalls int
var tryingCount int
var trying2Count int

func recordMemo(depth int, newChamp HeroVal, score float32, nodes []PlayerChamp, tail bool) memoItem {
  //hashKeys :=  []string{newChamp.String()}
  var hashKeys []string
  for n := len(nodes) - 1; n >= 0; n-- {
    h := nodes[n]
    if h.champ.champ == MaxHeroVal {
      depth = n
      break
    }
    if h.champ.champ == Empty {
      continue
    }

    hashKeys = append(hashKeys, h.champ.champ.String())
  }
  sort.Strings(hashKeys)
  m := memoItem{score, copyNodes(nodes)}
  key := fmt.Sprintf("%v:%v", depth, strings.Join(hashKeys, ","))
  memoLock.Lock()
  memo[key] = m
  memoLock.Unlock()
  //fmt.Printf("Recording %v %v\n", key, tail)
  /*
  if depth == 0 {
    fmt.Printf("===\nDepth was 0: %v\n===\n", FormatNodes(nodes))
  }
  */
  return m
}

func findBestNodes(ch chan memoItem, nodes []PlayerChamp) {
  diversity := map[HeroVal]bool{}
  playerCounts := map[string]int{}
  var champCount int
  var hashKeys []string
  var currentScore float32
  bestNodes := make([]PlayerChamp, len(nodes))
  //var waits int

  newCh := make(chan memoItem)
  var calls int

  for n := len(Nodes) - 1; n >= 0; n-- {
    // Consider the previous mapping locked in
    if nodes[n].champ.champ != MaxHeroVal {
      diversity[nodes[n].champ.champ] = true
      currentScore += champScore(nodes[n].champ)
      bestNodes[n] = nodes[n]
      playerCounts[nodes[n].player]++
      if nodes[n].champ.champ != Empty {
        champCount++
        hashKeys = append(hashKeys, nodes[n].champ.champ.String())
      }
      continue
    }
    //fmt.Printf("----\n%v\n%v\n%v\n", FormatNodes(nodes), playerCounts, champCount)
    if champCount == len(bg1) * playerMax {
      //fmt.Printf("Donezo\n")
      mi := recordMemo(n, Empty, currentScore, nodes, false)
      ch <- mi
      return
    }

    totalCount++
    if totalCount % 1000000 == 0 {
      fmt.Printf("%v %v %v %v\n%v\n%v\n", totalCount, n, len(memo), totalCalls, FormatNodes(nodes), runtime.NumGoroutine())
    }
    sort.Strings(hashKeys)
    key := fmt.Sprintf("%v:%v", n, strings.Join(hashKeys, ","))
    memoLock.Lock()
    m, ok := memo[key]
    memoLock.Unlock()

    //fmt.Printf("Checking %v %v\n", key, ok)
    if ok {
      memoCount++
      if memoCount % 10000 == 0 {
        fmt.Printf(".")
      }
      ch <- m
      return
    }
    bestScore := currentScore

    for _, hero := range Nodes[n+1] {
      // Find the best one of this champ in the BG
      if diversity[hero] {
        continue
      }
      for player, champset := range bg1 {
        // Players can only have 5 champs on the map
        if playerCounts[player] == playerMax {
          continue
        }
        if hero != MaxHeroVal {

          if c, ok := champset[hero]; ok {
           // fmt.Printf("hi\n")
            nodes[n] = PlayerChamp{player: player, champ: c}

            calls++
            totalCalls++
            tryingCount++

            go findBestNodes(newCh, copyNodes(nodes))
          }
        } else {
          // If any champ can go there it can also be empty
          nodes[n] = PlayerChamp{player: "empty", champ: Champ{Empty, 0, 0, 0}}
          calls++
          totalCalls++
          trying2Count++
          go findBestNodes(newCh, copyNodes(nodes))
        }
      }
    }

    for ; calls > 0; calls-- {
      select {
        case result := <-newCh:
          totalCalls--
          newNodes, newScore := result.nodes, result.score
          if newScore > bestScore {
            bestNodes = newNodes
            bestScore = newScore
          }
      }
    }

    mi := recordMemo(n, bestNodes[n].champ.champ, bestScore, bestNodes, true)
    ch <- mi
    return
    //return bestNodes, bestScore
  }
  ch <- memoItem{score: currentScore, nodes: nodes}
  //return nodes, currentScore
}
type Defenders struct {
  player string
  champs []Champ
  score float32
}

type PlayerDefenders struct {
  player string
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
  for k,v := range d {
    ret[k] = v
  }
  return ret
}

type memoItem2 struct {
  pds []PlayerDefenders
  score float32
  err error
  callArgs Defenders
}

var memo2 = map[string]memoItem2 {}

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

var first = map[string]bool{}

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
    if memoCount % 100 == 0 {
      fmt.Printf(".")
    }
    return mi.pds, mi.score, mi.err
  }

  if len(players) != 0 {
    p := players[0]
    //fmt.Printf("player %v\n", p)
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
    if _, ok := first[p]; !ok {
      first[p] = true
      fmt.Printf("%v %v %v\n", p, len(reducedChamps), len(combos))
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
        if newScore + ca.score > bestScore {
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

func main() {

  ch := make(chan memoItem2)
  t := time.Now()
  //combos := teamCombinations(5, bg1_2["sugar"], "sugar")
  //result, score, err := findBestBG(ch, map[HeroVal]bool{}, bg1_2, []string{"sugar", "dhdhqqq", "TomJenks", "LivingArtiface"})
  go findBestBG(ch, map[HeroVal]bool{}, bg1_2, []string{"sugar", "dhdhqqq", "TomJenks", "LivingArtiface"}, Defenders{})
  select {
  case mi := <-ch:
    result, score, _ := mi.pds, mi.score, mi.err

    d := time.Now().Sub(t)
    for _, pd := range result {
      fmt.Printf("%s: %s\n", pd.player, pd.defenders.String())
    }
    fmt.Printf("Took %v for %v combos\n", d, len(result))
    fmt.Printf("Score: %v\n", score)
  }
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
