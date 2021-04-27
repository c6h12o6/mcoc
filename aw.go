package main

import . "../mcoc/globals"
//import "github.com/satori/go.uuid"
import "fmt"
import "strings"
import "sort"
import "time"
import "sync"
import "runtime"

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
    for h, _ := range d.champs {
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
  champs map[HeroVal]Champ
}


func copyDefenders(d Defenders) Defenders {
  ret := map[HeroVal]Champ{}
  for h, c := range d.champs {
    ret[h] = c
  }
  return Defenders{champs: ret}
}

func copyBattleGroup(r map[string]Defenders)  map[string]Defenders {
  ret := map[string]Defenders{}
  for p, d := range r {
    ret[p] = copyDefenders(d)
  }
  return ret
}

type BGScore struct {
  score float32
  battleGroup map[string]Defenders
}

func (bg BGScore) String() string {
  return fmt.Sprintf("--\n%v\n%v\n--\n", bg.score, FormatBattleGroup(bg.battleGroup))
}

func (bg BGScore) Merge(start map[string]Defenders, players []string) BGScore {
  if bg.score == 0 {
    return bg
  }
  ret := copyBattleGroup(bg.battleGroup)
  for p, d := range start {
    for h, c := range d.champs {
      //fmt.Printf("%v %v %v\n%v\n", p, h, c, ret)
      ret[p].champs[h] = c
    }
  }

  var score float32
  for _, d := range ret {
    for _, c := range d.champs {
      score += champScore(c)
    }
  }

  return BGScore{score, ret}
}


var memo2 =map[string]BGScore{}

func recordMemo2(start map[string]Defenders, score BGScore) {
  var hashKeys []string
  for _, d := range start {
    for h, _ := range d.champs {
      hashKeys = append(hashKeys, h.String())
    }
  }
  sort.Strings(hashKeys)
  key := strings.Join(hashKeys, ",")

  memoLock.Lock()
  memo2[key] = score
  memoLock.Unlock()
}

func findBestRoster(battleGroup map[string]Defenders, players []string) BGScore {
  diversity := map[HeroVal]bool{}
  var currentScore float32
  bestScore := BGScore{score: 0}
  var hashKeys []string

  for _, d := range battleGroup {
    for h, c := range d.champs {
      diversity[h] = true
      currentScore += champScore(c)
      hashKeys = append(hashKeys, h.String())
    }
  }

  if len(diversity) == len(bg1) * playerMax {
    //fmt.Printf("Done!\n")
    return BGScore{score: currentScore, battleGroup: battleGroup}
  }

  sort.Strings(hashKeys)
  key := strings.Join(hashKeys, ",")
  memoLock.Lock()
  m, ok := memo2[key]
  memoLock.Unlock()

  if ok {
    memoCount++
    if memoCount % 1000 == 0 {
      fmt.Printf(".")
    }
    return m.Merge(battleGroup, players)
  }

  for _, player := range players {
    champset := bg1[player]
    // Players can only have 5 champs on the map
    if _, ok := battleGroup[player]; ok {
      if len(battleGroup[player].champs) == playerMax {
        //fmt.Printf("continuing away from player %v\n", player)
        continue
      }
    } else {
      battleGroup[player] = Defenders{}
    }

    for h, c := range champset {
      // Don't add it if it's already in the battle group
      if _, ok := diversity[h]; ok {
        continue
      }

      newBG := copyBattleGroup(battleGroup)
      newBG[player].champs[h] = c

      if player == "sugar" {
        fmt.Printf("adding %v\n", h)
        fmt.Printf("trying\n%v\n", FormatBattleGroup(newBG))
        fmt.Printf("best\n%v\n", FormatBattleGroup(bestScore.battleGroup))
      }

      newScore := findBestRoster(newBG, players)
      //fmt.Printf("result %v\n", newScore)
      if newScore.score > bestScore.score {
        //fmt.Printf("setting new score: %v", newScore)
        bestScore = newScore
      }
    }
    if bestScore.score == 0 {
      //fmt.Printf("nothing workable. f\n")
      recordMemo2(battleGroup, bestScore)
      return bestScore
    }
  }
  //fmt.Printf("Returning %v\n", bestScore)
  recordMemo2(battleGroup, bestScore)
  return bestScore
}


func main() {
  var players []string
  for p, _ := range bg1 {
    players = append(players, p)
  }
  fmt.Printf("players: %v\n", players)
  players = []string{"sugar", "dhdhqqq", "TomJenks", "LivingArtiface", "Yves", "Nino"}

  t := time.Now()
  battleGroup := map[string]Defenders{}
  score := findBestRoster(battleGroup, players)
  //bestNodes, score := result.nodes, result.score
  d := time.Now().Sub(t)

  //PrintNodes(bestNodes)

  fmt.Printf("%v\n", score)
  fmt.Printf("Took %v\n", d)
  //fmt.Printf("Trying %v Trying 2 %v\n", tryingCount, trying2Count)
}
