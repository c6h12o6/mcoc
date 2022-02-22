package globals

import "fmt"
import "github.com/c6h12o6/mcoc/smp"

type Champ struct {
	Champ      HeroVal
	Stars      int32
	Level      int32
	Sig        int32
	LockedNode int
  AssignedNode int32
}

func NewChamp(hv HeroVal, stars int32, level int32, sig int32) Champ {
	return Champ{hv, stars, level, sig, 0, 0}
}

func sigValue(a Champ) float32 {
  r := GetDefensiveSigRelevance(a.Champ)
  switch r {
  case UN:
    return 1.0
  case NA:
    if a.Sig == 0 { return 0.75 } else { return 1.0 }
  case HS:
    if a.Sig == 0 { return 0.75 } else { return  float32(0.75) + 0.00125 * float32(a.Sig)}
  default:
    fmt.Printf("Could not find %v\n", a)
  }
  return 0.01
}

func ChampValue(a Champ) float32 {
  if a.Stars == 6 && a.Level == 4 {
    return 150 * sigValue(a)
  } else if a.Stars == 6 && a.Level == 3 {
		return 130 * sigValue(a)
	} else if a.Stars == 6 && a.Level == 2 {
		return 110 * sigValue(a)
	} else if a.Stars == 5 && a.Level == 5 {
		return 100 * sigValue(a)
	} else if a.Stars == 6 && a.Level == 1 {
		return 65 * sigValue(a)
	} else if a.Stars == 5 && a.Level == 4 {
		return 60 + sigValue(a)
	} else {
		return 1
	}
}

func ChampScoreWithMasteries(a Champ, md int, suicides bool) float32 {
	if a.Champ == Empty {
		return 0.0
	}
  score := ChampValue(a) * float32(GetLevel(a.Champ))
  if _, ok := Mystics[a.Champ]; ok {
    score *= 1.0 + (0.02 * float32(md))
  }
  if suicides {
    score *= 0.9
  }
  return score
}

var ChampPreference = []int{
  55, 54, 53, 52, 51, 50, 49, 48, 47, 46,
  45, 44, 43, 42, 41, 40, 39, 38, 37, 36,
  35, 34, 33, 32, 31, 30, 29, 28, 27, 26,
  25, 24, 23, 22, 21, 20, 19, 18, 17, 16,
  15, 14, 13, 12, 11, 10,  9,  8,  7,  6,
   5,  4,  3,  2,  1}

func ChampPerson(id int) smp.Person {
  return smp.Person{
    ID: id,
    Prefers: ChampPreference,
  }
}

func ChampScore(a Champ) float32 {
  return ChampScoreWithMasteries(a, 0, false)
}

func (c Champ) String() string {
	return fmt.Sprintf("%v (%v/%v) %v", c.Champ.String(), c.Stars, c.Level)
}
