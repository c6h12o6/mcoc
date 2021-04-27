package main

import . "../mcoc/globals"
//import "github.com/satori/go.uuid"
import "fmt"
import "strings"
import "sort"
import "time"
import "sync"

type Champ struct {
  champ HeroVal
  stars int32
  level int32
  sig   int32
}

var playerMax = 3

type PlayerChamp struct {
  player string
  champ Champ
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
  },
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
  },

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

func recordMemo(diversity map[HeroVal]bool, newChamp HeroVal, score float32, nodes []PlayerChamp) memoItem {
  hashKeys :=  []string{string(newChamp)}
  var depth int
  for n := len(nodes) - 1; n >= 0; n-- {
    h := nodes[n]
    if h.champ.champ == MaxHeroVal {
      depth = n
      break
    }
    if h.champ.champ == Empty {
      continue
    }

    hashKeys = append(hashKeys, fmt.Sprintf("%s", h.champ.champ))
  }
  sort.Strings(hashKeys)
  m := memoItem{score, copyNodes(nodes)}
  key := fmt.Sprintf("%v:%v", depth, strings.Join(hashKeys, ","))
  memoLock.Lock()
  memo[key] = m
  memoLock.Unlock()
  fmt.Printf("Recording %v\n", key)
  return m
}

func findBestNodes(ch chan memoItem, nodes []PlayerChamp) {
  diversity := map[HeroVal]bool{}
  playerCounts := map[string]int{}
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
      hashKeys = append(hashKeys, fmt.Sprintf("%s", nodes[n].champ.champ))
      playerCounts[nodes[n].player]++
      continue
    }
    totalCount++
    if totalCount % 1000000 == 0 {
      fmt.Printf("%v %v %v %v\n%v\n", totalCount, n, len(memo), totalCalls, nodes)
    }
    sort.Strings(hashKeys)
    key := fmt.Sprintf("%v:%v", n, strings.Join(hashKeys, ","))
    memoLock.Lock()
    m, ok := memo[key]
    memoLock.Unlock()

    fmt.Printf("Checking %v %v\n", key, ok)
    if ok {
      memoCount++
      if memoCount % 1000 == 0 {
        fmt.Printf(".")
      }
      ch <- m
      return
    }
    bestScore := currentScore

/*
    for ; totalCalls > 500000;  {
      waits++
      if waits % 60 == 0 {
        fmt.Printf("! %v\n", waits)
      }
      time.Sleep(time.Minute)
    }
    */
    //fmt.Printf("diversity %v\n", diversity)

    for _, hero := range Nodes[n+1] {
      // Find the best one of this champ in the BG
      if diversity[hero] {
        //fmt.Printf("Skipping %s\n", hero)
        continue
      }
      for player, champset := range bg1 {
        // Players can only have 5 champs on the map
        if playerCounts[player] == playerMax {
          //fmt.Printf("%v is tapped out\n", player)
          continue
        }
        //fmt.Printf("player: %v, %v\n", player, champset)
        if hero != MaxHeroVal {
          //fmt.Printf("nope\n")
          //fmt.Printf("!!: %v\n", champset[hero])

          if c, ok := champset[hero]; ok {
            //fmt.Printf("hi\n")
            nodes[n] = PlayerChamp{player: player, champ: c}

            calls++
            totalCalls++
            tryingCount++
            if totalCalls % 10000 == 0 {
              //fmt.Printf("Trying  %v at node %v\n", nodes[n].PrettyPrint(), n)
            }

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

        /*else {
          for h, c := range champset {
            if diversity[h] {
              continue
            }

            nodes[n] = PlayerChamp{player: player, champ: c}
            //u, _ := uuid.NewV4()
            //fmt.Printf("%s Trying2  %v at node %v\n%v", u, nodes[n].PrettyPrint(), n, nodes)
            calls++
            totalCalls++
            if totalCalls % 10000 == 0 {
              //fmt.Printf("Trying2  %v at node %v\n", nodes[n].PrettyPrint(), n)
            }
            trying2Count++
            go findBestNodes(newCh, copyNodes(nodes))
          }
        }
        */
      }
      /*
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
      */
    }
    //fmt.Printf("Returning at node %v with score %v. %v\n", n, bestScore, len(bestNodes))
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

    mi := recordMemo(diversity, bestNodes[n].champ.champ, bestScore, bestNodes)
    ch <- mi
    return
    //return bestNodes, bestScore
  }
  ch <- memoItem{score: currentScore, nodes: nodes}
  //return nodes, currentScore
}

func main() {
  nodes := make([]PlayerChamp, len(Nodes))
  for n := len(Nodes) - 1; n >= 0; n-- {
    nodes[n] = PlayerChamp{champ: Champ{champ: MaxHeroVal}}
  }

  ch := make(chan memoItem)
  t := time.Now()
  go findBestNodes(ch, nodes)
  result := <-ch
  bestNodes, score := result.nodes, result.score
  d := time.Now().Sub(t)

  for n := 0; n < len(bestNodes); n++ {
    fmt.Printf("%v:\t%s\n", n+1, bestNodes[n].PrettyPrint())
  }
  fmt.Printf("%v\n", score)
  fmt.Printf("Took %v\n", d)
  fmt.Printf("Trying %v Trying 2 %v\n", tryingCount, trying2Count)
}
