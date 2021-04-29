package globals

import "fmt"
type Champ struct {
	Champ      HeroVal
	Stars      int32
	Level      int32
	Sig        int32
	LockedNode int
}

func NewChamp(hv HeroVal, stars int32, level int32, sig int32) Champ {
	return Champ{hv, stars, level, sig, 0}
}

func ChampValue(a Champ) float32 {
	if a.Stars == 6 && a.Level == 3 {
		return 15 + float32(a.Sig)/200
	} else if a.Stars == 6 && a.Level == 2 {
		return 9 + float32(a.Sig)/200
	} else if a.Stars == 5 && a.Level == 5 {
		return 8 + float32(a.Sig)/200
	} else if a.Stars == 6 && a.Level == 1 {
		return 6 + float32(a.Sig)/200
	} else if a.Stars == 5 && a.Level == 4 {
		return 5 + float32(a.Sig)/200
	} else {
		return 1
	}
}

func ChampScore(a Champ) float32 {
	if a.Champ == Empty {
		return 0.0
	}
	return ChampValue(a) + float32(GetLevel(a.Champ))
}

func (c Champ) String() string {
	return fmt.Sprintf("%v (%v/%v)", c.Champ.String(), c.Stars, c.Level)
}
