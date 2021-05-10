// arenasynergies is functionality devoted to maximizing the number
// of 3 champion teams in arena of Marvel's Contest of Champions
package main

import "fmt"
import "../mcoc/globals"
import "strings"
import "github.com/bradfitz/slice"
import "time"
import "encoding/gob"
import "bytes"
import "io/ioutil"
import "os"
import "sync"


// Variables for memoization

// Whether to use memoization
const memoize = true

// Whether to attempt to initialize memoization from disk
var initmemo = false

// memoization mapping
var memo map[string](string)

// uh... its a synchronization lock.
var lock sync.Mutex

// Some variables for tracking metrics
var memos int
var trials int
var skipped int
var depthcount []int

// combinations is a helper function for creating all possible combinations of
// values from "iterable" in groups of "r"
func Combinations(iterable []int, r int) [][]int {
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
func teamCombinations(teamsize int) []globals.TeamInfo {
    // TODO maybe only do this once
    var indices []int
    var teams []globals.TeamInfo

    // Seriously, go does not have a way to create a slice of 1-n
    // so you have to do this. gross
    for ii := 0; ii < len(globals.MyHeroes); ii++ {
        indices = append(indices, ii)
    }

    // Now get the combinations you need, in the form of slices of ints
    teamindices := combinations(indices, teamsize)

    // Turn those slices of integers into slices of heroes
    for _, teamnos := range teamindices {
        var team []globals.Hero
        for _, idx := range teamnos {
            team = append(team, globals.MyHeroes[idx])
        }
        
        // Count up some synergies and dont bother including anything in the result
        // that has a synergy count of 0. They can't help our count
        syncount := globals.SynergyCount(team)
        if syncount > 0 {
            teams = append(teams, globals.TeamInfo{team, syncount})
        }
    }

    // return the list in order of count, high to low
    slice.Sort(teams, func(i, j int) bool {
        return teams[i].Count > teams[j].Count
    })

    return teams
}

// contains checks to see whether a given slice of heroes has a given hero value in it
func contains(team []globals.Hero, item globals.HeroVal) bool {
    for _, teammate := range(team) {
        if teammate.Name == item {
            return true
        }
    }
    
    return false
}

// removeTeams takes a slices of TeamInfo objects and removes any elements that 
// have *any* of the heroes in rteam. These are the remaining viable teams after
// rteam has been removed. Also returns the total synergy count sum for optimizations
func removeTeams(allteams []globals.TeamInfo, rteam globals.TeamInfo) ([]globals.TeamInfo, int) {
    var ret []globals.TeamInfo
    var retcount int

    for _, checkteam := range allteams {
        good := true
        for _, hero := range checkteam.Team {
            if contains(rteam.Team, hero.Name) {
                good = false
                break
            }
        }
        if good {
            ret = append(ret, checkteam)
            retcount += checkteam.Count
        } 
    }

    return ret, retcount
}

// arenaTeams recursively tries pulling all* possible combinations of 3 man teams, returning a slice
// of TeamInfo objects representing the set of teams that maximizes the total number of synergies 
// across all teams. * Due time constraints, we do not actually try all combinations. Running
// all combinations would be testing something on the order of 10^64 options. We make some intelligent
// choices for performance sake, however, it is possible that there is a set of teams that *could*
// have a larger number of synergies, though it's unlikely.
func arenaTeams(remainingTIs []globals.TeamInfo, depth int, marked bool, ch chan []globals.TeamInfo) ([]globals.TeamInfo, int) {
    lock.Lock()
    if len(depthcount) <= depth {
        depthcount = append(depthcount, 0)
    }
    lock.Unlock()

    if len(remainingTIs) == 0 {
        trials++
        if trials % 1000000 == 0 {
            //fmt.Printf("Tried %d\n", trials)
        }
        ch <- make([]globals.TeamInfo, 0)
        return make([]globals.TeamInfo, 0), 0
    }

    // check the memo table
    memostr := globals.FormatTeamInfos(remainingTIs)
    if memoize {
        lock.Lock()
        val, ok := memo[memostr] 
        lock.Unlock()
        if ok {
            ret := globals.DeserializeTeamInfos(val)
            memos += 1
            syns := 0
            for _, ti := range ret {
                syns += ti.Count
            }
            ch <- ret
            return ret, syns
        }
    }
    var oldbestcount int
    var oldbestteams []globals.TeamInfo
    var allchans []chan []globals.TeamInfo

    // iterate through all the teams remaining, seeing what it would look like
    // if you included testteam in the final answer
    for idx, testteam := range remainingTIs {
        remaining, remainingCount := removeTeams(remainingTIs, testteam) 

        if remainingCount + testteam.Count <= oldbestcount {
            skipped += len(remaining)
            continue
        }

        newch := make(chan []globals.TeamInfo)
        allchans = append(allchans, newch)

        //bestteams, syncount := arenaTeams(remaining, depth+1, marked)
        go arenaTeams(remaining, depth+1, marked, newch)

        // These are the optimizations for time. At the highest levels of recursion.
        // we have to significantly limit the number of teams that we can try.
        // This is probably ok because the highest teams in the list have the largest
        // number of synergies anyway.
        // As we get deeper into the recursion, we can try more options.
        if depth > 10 {
            if idx > 10 {
                break
            }
        } else {
            if idx > 2  {
                break
            }
        }
    }
        
    // Collect all the results
    results := make([][]globals.TeamInfo, 0)
    for idx, checkchan := range allchans {
        tmpthing := <-checkchan
        tmpthing = append(tmpthing, remainingTIs[idx])
        results = append(results, tmpthing)
    }

    // Find the result with the greatest number of synergies
    for _, result := range results {
        syncount := 0
        for _, team := range result {
            syncount += globals.SynergyCount(team.Team)
        }

        if syncount > oldbestcount {
            oldbestcount = syncount
            oldbestteams = result
            
            lock.Lock()
            depthcount[depth] = depthcount[depth]+1
            lock.Unlock()

            if depth < 2 {
                //fmt.Printf("Depth %d\nNewBest\n%s\n", depth, globals.FormatTeamInfos(oldbestteams))
            }
        }
    }

    // store the result for memoization
    if memoize {
        lock.Lock()
        memo[memostr] = globals.FormatTeamInfos(oldbestteams)
        lock.Unlock()
    }

    ch <- oldbestteams
    return oldbestteams, oldbestcount
}

func writeMemoToFile(mapping map[string]string, filename string) { 
    // Save off the map to disk
    b := new(bytes.Buffer)

    e := gob.NewEncoder(b)
    // Encoding the map
    err := e.Encode(mapping)
    if err != nil {
        panic(err)
    }
    ioutil.WriteFile(filename, b.Bytes(), 0644)
}

func loadMemoFromFile(filename string) {
    _, err := os.Stat(filename)
    if err == nil {
        mapbytes, err := ioutil.ReadFile(filename)
        if err != nil {
            panic(err)
        }

        b := bytes.NewBuffer(mapbytes)
        d := gob.NewDecoder(b)
        
        decodedMap := make(map[string]string, 0)
        // Decoding the serialized data
        err = d.Decode(&decodedMap)
        if err != nil {
            panic(err)
        }

        // Ta da! It is a map!
        for k, v := range decodedMap {
            memo[k] = v
        }
        fmt.Printf("Map has %d elements\n", len(memo))
    }
}


func main() {
    t := time.Now()

    memo = make(map[string]string, 0)
    
    // Restore the serialized memo table
    if memoize && initmemo {
        files, _ := ioutil.ReadDir("./")
        for _, f := range files {
            fname := f.Name()
            if strings.HasPrefix(fname, "map") && strings.HasSuffix(fname, ".strings") {
                loadMemoFromFile(fname)
            }
        }
    }

    // get all the teams
    teams := teamCombinations(3)
    fmt.Println(len(teams))

    //Actually solve the problem at hand
    newch := make(chan []globals.TeamInfo)
    go arenaTeams(teams, 0, false, newch)
    bestteam := <- newch
    bestcount := 0
    
    for _, team := range bestteam {
        bestcount += globals.SynergyCount(team.Team)
    }

    for _, teaminfo := range bestteam {
        fmt.Printf("%s: %d\n", globals.FormatTeam(teaminfo.Team), teaminfo.Count)
    }
    fmt.Printf("Total %d\n", bestcount)
    fmt.Printf("Memos %d\n", memos)
    fmt.Printf("Skipped %d\n", skipped)
    fmt.Printf("Depthcount: %v\n", depthcount)

    delta := time.Now().Sub(t)
    fmt.Println(delta)

    // Save off the memoization
    t = time.Now()
    if memoize {
        // Write out the memo files 500,000 at a time
        // Anything bigger and gob barfs
        idx := 1
        tmpmemo := make(map[string]string, 0)
        for k, v := range memo {
            if idx % 500000 == 0 {
                writeMemoToFile(tmpmemo, fmt.Sprintf("map%d.strings2", idx/500000))
                tmpmemo = make(map[string]string, 0)
            }
            tmpmemo[k] = v
            idx += 1
        }

        writeMemoToFile(tmpmemo, fmt.Sprintf("map%d.strings2", idx/500000 + 1))
        fmt.Printf("Waiting for all writes to finish\n")
        fmt.Printf("Done\n")
    }

    delta = time.Now().Sub(t)
    fmt.Println(delta)
}
