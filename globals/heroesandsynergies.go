package globals

//go:generate stringer -type=HeroVal

import "regexp"
import "strings"
import "fmt"

// Enumerated type that prints out its name instead of its integer value
type HeroVal int

const (
    Abomination HeroVal = iota
    AgentVenom
    AntMan
    ArchAngel
    Beast 
    BlackBolt
    BlackPanther
    BlackPantherCivilWar
    BlackWidow
    Cable
    CaptainAmerica
    CaptainAmericaWWII
    CaptainMarvel
    Carnage
    Colossus
    Crossbones
    Cyclops
    CyclopsNew
    DareDevil
    DareDevilClassic
    Deadpool
    DeadpoolX
    DoctorStrange
    DoctorVoodoo
    Dormammu
    Drax
    Elektra
    Electro
    Falcon
    Gambit
    Gamora
    GhostRider
    GreenGoblin
    Groot
    Guillotine
    Gwenpool
    Hawkeye
    HowardTheDuck
    Hulk
    Hulkbuster
    Hyperion
    Iceman
    IronFist
    IronFistImmortal
    IronMan
    IronPatriot
    JoeFixit
    Juggernaut
    KamalaKhan
    Kang
    Karnak
    Loki
    LukeCage
    Magik
    Magneto
    MagnetoNow
    MilesMorales
    MoonKnight
    Mordo
    MsMarvel
    Nightcrawler
    OldManLogan
    Phoenix
    Psylocke
    Punisher
    Quake
    RedHulk
    Rhino
    RocketRaccoon
    Rogue
    Ronan
    ScarletWitch
    SheHulk
    SpiderGwen
    Spiderman
    SpidermanSymbiote
    SpidermanStark
    StarLord
    Storm
    SuperiorIronMan
    Thor
    ThorJaneFoster
    Ultron
    UnstoppableColossus
    Venom
    VenomPool
    Vision
    VisionAgeOfUltron
    Vulture
    WarMachine
    WinterSoldier
    Wolverine
    X23
    YellowJacket
    Yondu
    Medusa
    Kingpin
    Hood
    CivilWarrior
    Punisher2099
    Angela
    DoctorOctopus
    // This should always be last
    MaxHeroVal
)

// Struct to represent a set of synergies
type Synergy struct {
    Vals []HeroVal
    Desc string
}

// Struct associating a HeroVal with a set of synergies
type Hero struct {
    Name HeroVal
    Synergies []Synergy
}

// Struct representing a team and their associated weight (defaults to synergies score
type TeamInfo struct {
    Team []Hero
    Count int
}

var abominationsynergies = []Synergy {
    {[]HeroVal{Rhino}, "Friends"},
    {[]HeroVal{Hulk}, "Nemesis"},
    {[]HeroVal{SheHulk}, "Rivals"},
    {[]HeroVal{RedHulk}, "Rivals"},
}

var agentvenomsynergies = []Synergy {
    {[]HeroVal{RedHulk, Groot}, "Teammates"},
    {[]HeroVal{Spiderman}, "Idol"},
    {[]HeroVal{SpidermanSymbiote}, "Family"},
    {[]HeroVal{Venom, VenomPool}, "Family"},
}

var angelasynergies = []Synergy {
    {[]HeroVal{RocketRaccoon, StarLord}, "Friends"},
    {[]HeroVal{Gamora, Groot}, "Friends"},
    {[]HeroVal{Thor}, "Family"},
    {[]HeroVal{Loki}, "Family"},
}

var archangelsynergies = []Synergy {
    {[]HeroVal{Phoenix, Beast}, "Mutant Agenda 3"},
    {[]HeroVal{Iceman, Colossus}, "Mutant Agenda 3"},
    {[]HeroVal{Psylocke}, "Romance 3"},
    {[]HeroVal{GhostRider, BlackWidow}, "Teammates3"},
}

var antmansynergies = []Synergy {
    {[]HeroVal{YellowJacket}, "Nemesis"},
    {[]HeroVal{IronMan}, "Teammates"},
    {[]HeroVal{Spiderman}, "Teammates"},
    {[]HeroVal{Hulk}, "Friends"},
}

var beastsynergies = []Synergy {
                {[]HeroVal{Gambit}, "Friends3"},
                {[]HeroVal{IronPatriot}, "Friends3"},
                {[]HeroVal{BlackPanther, SuperiorIronMan}, "Masterminds3"},
                {[]HeroVal{Nightcrawler, Colossus}, "MutantAgenda3"},
}
var blackboltsynergies = []Synergy {
            {[]HeroVal{KamalaKhan, Ronan}, "CosmicSupremacy"},
            {[]HeroVal{Spiderman, CyclopsNew}, "Friends3"},
            {[]HeroVal{Kang}, "Nemesis2"},
            {[]HeroVal{Hulk}, "Rivals3"},
}

var blackpanthersynergies = []Synergy {
    {[]HeroVal{Deadpool}, "Enemies3"},
    {[]HeroVal{IronFist, IronFistImmortal},  "Rivals3"},
    {[]HeroVal{Storm}, "Romance3"},
}

var bpcwsynergies = []Synergy {
    {[]HeroVal{VisionAgeOfUltron, BlackWidow}, "Friends3"},
    {[]HeroVal{WinterSoldier}, "Nemesis3"},
    {[]HeroVal{AntMan}, "Rivals3"},
    {[]HeroVal{Hawkeye, BlackPanther}, "SkillDomination3"},
}

var blackwidowsynergies = []Synergy {
    {[]HeroVal{CaptainMarvel, MsMarvel}, "Friends3"},
    {[]HeroVal{WinterSoldier}, "Romance3"},
    {[]HeroVal{Hawkeye}, "Romance3"},
    {[]HeroVal{Hulk, Hulkbuster}, "Avengers3"},
}

var cablesynergies = []Synergy {
    {[]HeroVal{Deadpool, DeadpoolX}, "Enemies3"},
    {[]HeroVal{Cyclops, CyclopsNew}, "Family3"},
    {[]HeroVal{Phoenix}, "Family3"},
    {[]HeroVal{Rogue}, "Teammates3"},
}

var captainamericasynergies = []Synergy {
    {[]HeroVal{Spiderman, SpidermanSymbiote}, "Friends3"},
    {[]HeroVal{WinterSoldier}, "Friends3"},
    {[]HeroVal{SuperiorIronMan}, "Enemies3"},
    {[]HeroVal{IronMan}, "Enemies3"},
}

var captainamericawwiisynergies = []Synergy {
    {[]HeroVal{WinterSoldier}, "Friends3"},
    {[]HeroVal{Wolverine}, "Friends3"},
    {[]HeroVal{Guillotine}, "Friends3"},
}


var captainmarvelsynergies = []Synergy {
    {[]HeroVal{CaptainAmerica}, "Friends3"},
    {[]HeroVal{Gamora}, "Friends3"},
    {[]HeroVal{IronMan}, "Friends3"},
    {[]HeroVal{Wolverine}, "Romance3"},
}

var civilwarriorsynergies = []Synergy {
    {[]HeroVal{Falcon}, "Friends"},
    {[]HeroVal{WinterSoldier}, "Friends"},
    {[]HeroVal{Guillotine}, "Teammates"},
    {[]HeroVal{IronMan, Hulkbuster}, "Rivals"},
}

var colossussynergies = []Synergy {
    {[]HeroVal{Wolverine, OldManLogan}, "Friends3"},
    {[]HeroVal{Magik}, "Magik"},
    {[]HeroVal{Juggernaut}, "Enemies3"},
}

var crossbonessynergies = []Synergy {
    {[]HeroVal{CaptainAmerica, CaptainAmericaWWII}, "Enemies3"},
    {[]HeroVal{Falcon}, "Enemies3"},
    {[]HeroVal{WinterSoldier}, "Rivals3"},
    {[]HeroVal{BlackWidow}, "Rivals3"},
}

var cyclopssynergies = []Synergy {
    {[]HeroVal{Storm}, "Teammates"},
    {[]HeroVal{Magneto}, "Nemesis"},
    {[]HeroVal{Colossus, Wolverine}, "Mutant Agenda"},
    {[]HeroVal{Phoenix}, "Romance"},
}

var cyclopsnewsynergies = []Synergy {
    {[]HeroVal{Magneto}, "Nemesis"},
    {[]HeroVal{Wolverine}, "Rivals"},
    {[]HeroVal{Colossus, Storm}, "Mutant Agenda"},
}

var daredevilsynergies = []Synergy {
    {[]HeroVal{Elektra}, "Romance"},
    {[]HeroVal{LukeCage}, "Teammates"},
    {[]HeroVal{Punisher}, "Rivals"},
}

var daredevilclassicsynergies = []Synergy {
    {[]HeroVal{BlackWidow}, "Romance"},
    {[]HeroVal{SuperiorIronMan}, "Rivals"},
    {[]HeroVal{Elektra}, "Romance"},
}
var deadpoolxsynergies = []Synergy {
    {[]HeroVal{MoonKnight}, "Rivals"},
    {[]HeroVal{Deadpool}, "Rivals"},
    {[]HeroVal{MagnetoNow}, "Friends"},
}

var draxsynergies = []Synergy {
    {[]HeroVal{Gamora}, "Rivals"},
    {[]HeroVal{StarLord}, "Teammates"},
    {[]HeroVal{Ronan}, "Enemies"},
    {[]HeroVal{AgentVenom}, "Teammates"},
}

var dormammusynergies = []Synergy {
    {[]HeroVal{Mordo}, "Dark Empowerment 3"},
    {[]HeroVal{DoctorVoodoo}, "Enemies 3"},
    {[]HeroVal{Hood}, "Inseparable 3"},
    {[]HeroVal{DoctorStrange}, "Nemesis 3"},
}

var drstrangesynergies = []Synergy {
    {[]HeroVal{X23, Thor}, "Friends3"},
    {[]HeroVal{Spiderman}, "Friends3"},
    {[]HeroVal{ScarletWitch}, "Teammates3"},
    {[]HeroVal{BlackBolt}, "Teammates3"},
}

var droctopussynergies = []Synergy {
    {[]HeroVal{Electro}, "Particle Physics"},
    {[]HeroVal{AntMan}, "Biochemistry"},
    {[]HeroVal{Vulture}, "Engineering"},
    {[]HeroVal{SpidermanStark}, "Nemeis"},
}

var drvoodoosynergies = []Synergy {
    {[]HeroVal{WinterSoldier}, "Friends3"},
    {[]HeroVal{Hood}, "Enemies3"},
    {[]HeroVal{DoctorStrange}, "Rivals3"},
    {[]HeroVal{Rogue}, "Teammates3"},
}

var electrosynergies = []Synergy {
    {[]HeroVal{Spiderman}, "Nemesis"},
    {[]HeroVal{Rhino}, "Friends"},
    {[]HeroVal{MilesMorales}, "Enemies"},
    {[]HeroVal{Venom}, "Teammates"},
}

var elektrasynergies = []Synergy {
    {[]HeroVal{DareDevil, DareDevilClassic}, "Romance"},
    {[]HeroVal{BlackWidow}, "Rivals"},
    {[]HeroVal{Wolverine}, "Friends"},
    {[]HeroVal{Deadpool, DeadpoolX}, "Teammates"},
}

var falconsynergies = []Synergy {
    {[]HeroVal{CaptainAmerica, CaptainAmericaWWII}, "Friends"},
    {[]HeroVal{WarMachine, BlackWidow}, "Enemies"},
    {[]HeroVal{AntMan, Hawkeye}, "Teammates"},
    {[]HeroVal{VisionAgeOfUltron, BlackPantherCivilWar}, "Enemies"},
}


var gambitsynergies = []Synergy {
    {[]HeroVal{X23}, "Friends3"},
    {[]HeroVal{Magneto}, "Enemies3"},
    {[]HeroVal{Nightcrawler, Beast}, "Teammates3"},
}

var gamorasynergies = []Synergy {
    {[]HeroVal{AgentVenom}, "Teammates3"},
    {[]HeroVal{Drax}, "Rivals3"},
    {[]HeroVal{SheHulk}, "Friends3"},
    {[]HeroVal{StarLord}, "Friends3"},
}

var ghostridersynergies = []Synergy {
    {[]HeroVal{Punisher}, "Rivals3"},
    {[]HeroVal{X23}, "Teammates3"},
    {[]HeroVal{Deadpool, Elektra}, "Teammates3"},
}

var grootsynergies = []Synergy {
    {[]HeroVal{RocketRaccoon}, "Inseparable3"},
    {[]HeroVal{StarLord}, "Friends3"},
    {[]HeroVal{Drax}, "Teammates2"},
    {[]HeroVal{Gamora}, "Teammates2"},
}

var guillotinesynergies = []Synergy {
    {[]HeroVal{CaptainAmericaWWII}, "Teammates"},
    {[]HeroVal{Magik}, "Rivals"},
    {[]HeroVal{Venom}, "Nemesis"},
    {[]HeroVal{BlackPanther}, "Friends"},
}

var gwenpoolsynergies = []Synergy {
    {[]HeroVal{Deadpool, DeadpoolX}, "IDOL"},
    {[]HeroVal{HowardTheDuck}, "Friends"},
    {[]HeroVal{ThorJaneFoster}, "Enemies"},
    {[]HeroVal{SpidermanSymbiote}, "Teammates"},
}

var hawkeyesynergies = []Synergy {
    {[]HeroVal{MoonKnight}, "Friends"},
    {[]HeroVal{IronMan}, "Friends"},
    {[]HeroVal{Hulk}, "Friends"},
    {[]HeroVal{ScarletWitch}, "Romance"},
}

var hoodsynergies = []Synergy {
    {[]HeroVal{JoeFixit}, "Crime Bosses"},
    {[]HeroVal{DoctorVoodoo, Punisher}, "Enemies"},
    {[]HeroVal{IronPatriot, Loki}, "Teammates"},
    {[]HeroVal{Dormammu}, "Dark Empowerment"},
}

var howardsynergies = []Synergy {
    {[]HeroVal{SheHulk}, "Friends"},
    {[]HeroVal{RocketRaccoon}, "Friends"},
    {[]HeroVal{MsMarvel}, "Teammates"},
}

var hulksynergies = []Synergy {
    {[]HeroVal{Hawkeye}, "Friends3"},
    {[]HeroVal{Abomination}, "Enemies3"},
    {[]HeroVal{Thor}, "Rivals3"},
}

var hulkbustersynergies = []Synergy {
    {[]HeroVal{Ultron}, "Enemies3"},
    {[]HeroVal{IronMan}, "Family3"},
    {[]HeroVal{SuperiorIronMan}, "Family3"},
    {[]HeroVal{YellowJacket, Hulk}, "Rivals3"},
}

var hyperionsynergies = []Synergy {
    {[]HeroVal{Thor}, "Friends"},
    {[]HeroVal{IronMan}, "Enemies"},
    {[]HeroVal{DoctorStrange}, "Enemies"},
}

var icemansynergies = []Synergy {
    {[]HeroVal{Hyperion, Magneto}, "Enemies"},
    {[]HeroVal{GhostRider, BlackWidow}, "Teammates"},
    {[]HeroVal{Spiderman}, "Rivals"},
}

var ironfistsynergies = []Synergy {
    {[]HeroVal{BlackPanther}, "Friends"},
    {[]HeroVal{DoctorStrange}, "Friends"},
    {[]HeroVal{Wolverine}, "Friends"},
    {[]HeroVal{LukeCage, SheHulk}, "Heroes For Hire"},
}
var ironmansynergies = []Synergy {
    {[]HeroVal{CaptainAmerica}, "Friends3"},
    {[]HeroVal{Ultron}, "Enemies3"},
    {[]HeroVal{ThorJaneFoster}, "Teammates 3"},
    {[]HeroVal{Thor, WarMachine}, "Teammates 3"},
}

var ironpatriotsynergies = []Synergy {
    {[]HeroVal{IronMan}, "Friends3"},
    {[]HeroVal{Spiderman}, "Enemies3"},
    {[]HeroVal{CaptainAmerica}, "Rivals3"},
}

var joefixitsynergies = []Synergy {
    {[]HeroVal{Wolverine}, "Friends3"},
    {[]HeroVal{MoonKnight}, "Enemies3"},
    {[]HeroVal{MsMarvel}, "Enemies3"},
    {[]HeroVal{Hulk}, "Nemesis3"},
}

var juggernautsynergies = []Synergy {
    {[]HeroVal{Colossus, UnstoppableColossus}, "Enemies3"},
    {[]HeroVal{Hulk}, "Enemies3"},
    {[]HeroVal{DoctorStrange}, "Nemesis 3"},
}

var karnaksynergies = []Synergy {
    {[]HeroVal{Magneto}, "Enemies3"},
    {[]HeroVal{BlackBolt}, "Family3"},
    {[]HeroVal{Beast}, "Teammates3"},
    {[]HeroVal{CaptainMarvel, MsMarvel}, "Teammates3"},
}

var lokisynergies = []Synergy {
    {[]HeroVal{Hulk, RedHulk, ThorJaneFoster}, "Enemies3"},
    {[]HeroVal{Thor}, "Enemies3"},
    {[]HeroVal{Magneto, MagnetoNow}, "Masterminds2"},
}

var lukecagesynergies = []Synergy {
    {[]HeroVal{Rhino}, "Enemies3"},
    {[]HeroVal{IronFist, IronFistImmortal}, "Heroes for Hire 3"},
    {[]HeroVal{DareDevilClassic}, "Teammates2"},
    {[]HeroVal{Juggernaut, IronPatriot}, "Thunderbolts 2"},
}

var milesmoralessynergies = []Synergy {
    {[]HeroVal{Electro}, "Enemies2"},
    {[]HeroVal{IronPatriot}, "Enemies3"},
    {[]HeroVal{Venom}, "Enemies3"},
    {[]HeroVal{SpiderGwen}, "Teammates3"},
}

var moonknightsynergies = []Synergy {
    {[]HeroVal{Spiderman}, "Friends3"},
    {[]HeroVal{IronPatriot}, "Enemies3"},
    {[]HeroVal{Deadpool, DeadpoolX}, "Rivals3"},
    {[]HeroVal{DareDevilClassic}, "Teammates3"},
}

var msmarvelsynergies = []Synergy {
    {[]HeroVal{IronMan}, "Teammates"},
    {[]HeroVal{Hulk}, "Teammates"},
    {[]HeroVal{CaptainAmerica}, "Friends"},
    {[]HeroVal{Thor, ThorJaneFoster}, "Teammates"},
}

var magiksynergies = []Synergy {
    {[]HeroVal{Colossus, UnstoppableColossus}, "Family3"},
    {[]HeroVal{Storm}, "Friends3"},
    {[]HeroVal{CyclopsNew, Guillotine}, "Teammates3"},
    {[]HeroVal{Juggernaut}, "Enemies3"},
}

var magnetosynergies = []Synergy {
    {[]HeroVal{CyclopsNew, Cyclops}, "Nemesis"},
    {[]HeroVal{Storm}, "Rivals"},
    {[]HeroVal{Wolverine}, "Enemies"},
    {[]HeroVal{Magik}, "Friends"},
}

var mordosynergies = []Synergy {
    {[]HeroVal{DoctorStrange}, "Friends3"},
    {[]HeroVal{Thor}, "Enemies3"},
    {[]HeroVal{Falcon, Abomination}, "Enemies2"},
    {[]HeroVal{DoctorStrange}, "Rivals3"},
}

var nightcrawlersynergies = []Synergy {
    {[]HeroVal{Beast}, "Friends"},
    {[]HeroVal{Juggernaut}, "Enemies"},
    {[]HeroVal{CyclopsNew, Cyclops}, "Teammates"},
    {[]HeroVal{X23}, "Rivals"},
}

var oldmanlogansynergies = []Synergy {
    {[]HeroVal{Hawkeye}, "Friends3"}, 
    {[]HeroVal{Wolverine}, "Enemies3"},
    {[]HeroVal{Hulk}, "Nemesis3"},
}

var phoenixsynergies = []Synergy {
    {[]HeroVal{Storm}, "Friends3"},
    {[]HeroVal{Beast, Nightcrawler}, "Teammates"},
    {[]HeroVal{Gamora}, "Teammates"},
    {[]HeroVal{Cyclops, Wolverine}, "It's Complicated"},
}

var punishersynergies = []Synergy {
    {[]HeroVal{Spiderman}, "Rivals3"},
    {[]HeroVal{DareDevil, DareDevilClassic}, "Rivals3"},
    {[]HeroVal{Rhino}, "Teammates3"},
}

var punisher2099synergies = []Synergy {
    {[]HeroVal{Punisher}, "CrossTraining"},
    {[]HeroVal{CaptainAmerica}, "Loyal Minister"},
    {[]HeroVal{Thor}, "Believer"},
}

var quakesynergies = []Synergy {
    {[]HeroVal{IronPatriot, Crossbones}, "Enemies3"},
    {[]HeroVal{Hawkeye}, "ShieldAgents3"},
    {[]HeroVal{BlackWidow}, "ShieldClearance10"}, 
    {[]HeroVal{Karnak, BlackBolt}, "Rivals3"},
}

var redhulksynergies = []Synergy {
    {[]HeroVal{Elektra, AgentVenom}, "Thunderbolts"},
    {[]HeroVal{Abomination}, "Enemies"},
    {[]HeroVal{Hulk}, "Nemesis"},
    {[]HeroVal{X23}, "Teammates"},
}

var rhinosynergies = []Synergy {
    {[]HeroVal{Abomination}, "Friends"},
    {[]HeroVal{Electro}, "Teammates"},
    {[]HeroVal{Spiderman, SpiderGwen}, "Enemies"},
    {[]HeroVal{Punisher}, "Friends"},
}

var rocketsynergies = []Synergy {
    {[]HeroVal{StarLord, Groot}, "Friends"},
    {[]HeroVal{Ronan}, "Enemies"},
    {[]HeroVal{Drax}, "Teammates"},
    {[]HeroVal{Gamora}, "Teammates"},
}

var ronansynergies = []Synergy {
    {[]HeroVal{BlackBolt}, "Rivals"},
    {[]HeroVal{IronMan}, "Enemies"},
    {[]HeroVal{Gamora}, "Rivals"},
    {[]HeroVal{Hulk}, "Enemies"},
}

var roguesynergies = []Synergy {
    {[]HeroVal{Nightcrawler}, "Family"},
    {[]HeroVal{Gambit}, "Romance"},
    {[]HeroVal{MsMarvel}, "Rivals"},
    {[]HeroVal{Deadpool}, "Mutant Agenda"},
}

var scarletwitchsynergies = []Synergy {
    {[]HeroVal{CaptainMarvel, MsMarvel}, "Friends3"},
    {[]HeroVal{Phoenix}, "Enemies3"},
    {[]HeroVal{Vision}, "Romance3"},
    {[]HeroVal{AntMan}, "Teammates3"},
}

var shehulksynergies = []Synergy {
    {[]HeroVal{DareDevilClassic, DareDevil}, "LegalDefense3"},
    {[]HeroVal{Hulk}, "Family3"},
    {[]HeroVal{SuperiorIronMan}, "Romance3"},
    {[]HeroVal{MsMarvel}, "Teammates3"},
}

var spidergwensynergies = []Synergy {
    {[]HeroVal{Rhino}, "Enemies3"},
    {[]HeroVal{DareDevilClassic}, "Enemies3"},
    {[]HeroVal{Spiderman}, "Romance3"},
    {[]HeroVal{Punisher}, "Rivals3"},
}

var spidermansynergies = []Synergy {
    {[]HeroVal{CaptainAmerica}, "Friends"},
    {[]HeroVal{Electro}, "Enemies"},
    {[]HeroVal{Hawkeye}, "Friends"},
    {[]HeroVal{Wolverine}, "Friends"},
}

var spidermansymbiotesynergies = []Synergy {
    {[]HeroVal{Electro}, "Enemies3"},
    {[]HeroVal{AgentVenom}, "Family3"},
    {[]HeroVal{Storm}, "Family3"},
}

var spidermanstarksynergies = []Synergy {
    {[]HeroVal{IronMan, Hulkbuster}, "Knowledge Share"},
    {[]HeroVal{Vulture}, "Avengers Tryout"},
    {[]HeroVal{KamalaKhan, MilesMorales}, "ContestNoobs"},
}

var starlordsynergies = []Synergy { 
    {[]HeroVal{RocketRaccoon, Groot}, "Friends3"},
    {[]HeroVal{Drax}, "Teammates3"},
    {[]HeroVal{Gamora}, "Teammates3"},
}

var stormsynergies = []Synergy {
    {[]HeroVal{Magik}, "Friends3"},
    {[]HeroVal{Magneto, MagnetoNow}, "Enemies3"},
    {[]HeroVal{BlackPanther}, "Romance3"},
    {[]HeroVal{Cyclops, CyclopsNew}, "Teammates3"},
}

var superiorironmansynergies = []Synergy {
    {[]HeroVal{DareDevilClassic}, "Rivals"},
    {[]HeroVal{CaptainAmerica}, "Enemies"},
    {[]HeroVal{Thor}, "Teammates"},
}

var thorsynergies = []Synergy {
    {[]HeroVal{DoctorStrange}, "Friends3"},
    {[]HeroVal{IronMan}, "Friends3"},
    {[]HeroVal{Juggernaut}, "Enemies3"},
}

var thorjanefostersynergies = []Synergy {
    {[]HeroVal{BlackWidow}, "Friends"},
    {[]HeroVal{Thor}, "Romance"},
    {[]HeroVal{Vision}, "Teammates"},
    {[]HeroVal{JoeFixit}, "Teammates"},
}

var ultronsynergies = []Synergy {
    {[]HeroVal{ScarletWitch}, "Friends3"},
    {[]HeroVal{BlackWidow}, "Enemies3"},
    {[]HeroVal{IronMan, SuperiorIronMan}, "Family3"},
}

var unstoppablecolossussynergies = []Synergy {
    {[]HeroVal{CyclopsNew}, "Teammates 3"},
    {[]HeroVal{Juggernaut}, "Rivals 3"},
    {[]HeroVal{Magik}, "Family 3"},
    {[]HeroVal{Wolverine, OldManLogan}, "Friends 3"},
}

var venomsynergies = []Synergy {
    {[]HeroVal{Spiderman}, "Nemesis"},
    {[]HeroVal{Electro}, "Rivals"},
    {[]HeroVal{SpidermanSymbiote}, "Family"},
    {[]HeroVal{JoeFixit}, "Teammates"},
}

var venompoolsynergies = []Synergy {
    {[]HeroVal{Venom}, "Inseparable3"},
    {[]HeroVal{Deadpool}, "Friends3"},
    {[]HeroVal{DeadpoolX}, "Friends3"},
    {[]HeroVal{AgentVenom, SpidermanSymbiote}, "Family3"},
}

var visionsynergies = []Synergy {
    {[]HeroVal{ScarletWitch}, "Romance"},
    {[]HeroVal{IronMan}, "Teammates"},
    {[]HeroVal{Magneto}, "Enemies"},
}

var visionaousynergies = []Synergy {
    {[]HeroVal{IronMan}, "Family"},
    {[]HeroVal{ScarletWitch}, "Enemies"},
    {[]HeroVal{Ultron}, "Nemesis"},
}

var warmachinesynergies = []Synergy {
    {[]HeroVal{BlackPanther}, "Enemies"},
    {[]HeroVal{Hulkbuster}, "Friends"},
    {[]HeroVal{BlackWidow}, "Teammates"},
    {[]HeroVal{Hawkeye}, "Enemies"},
}
var wintersoldiersynergies = []Synergy {
    {[]HeroVal{Wolverine}, "Friends3"},
    {[]HeroVal{CaptainAmericaWWII}, "Friends3"},
    {[]HeroVal{CaptainAmerica}, "Teammates3"},
}

var wolverinesynergies = []Synergy {
    {[]HeroVal{CyclopsNew, Cyclops}, "Rivals"},
    {[]HeroVal{Magneto}, "Enemies"},
    {[]HeroVal{CaptainAmerica, CaptainAmericaWWII}, "Friends"},
}

var x23synergies = []Synergy {
    {[]HeroVal{Wolverine, OldManLogan}, "Family3"},
    {[]HeroVal{AgentVenom}, "Teammates3"},
    {[]HeroVal{RedHulk}, "Teammates3"},
}

var yellowjacketsynergies = []Synergy {
    {[]HeroVal{AntMan}, "Nemesis3"},
    {[]HeroVal{SuperiorIronMan}, "Idol3"},
    {[]HeroVal{Ultron}, "Rivals3"},
    {[]HeroVal{JoeFixit}, "Rivals3"},
}

var yondusynergies = []Synergy {
    {[]HeroVal{Beast, Nightcrawler}, "ItAintEasy3"},
    {[]HeroVal{RocketRaccoon}, "Friends3"},
    {[]HeroVal{Ronan}, "Enemies3"},
    {[]HeroVal{StarLord}, "Rivals3"},
}


var Heroes = []Hero {
    { Abomination, abominationsynergies },
    { AgentVenom, agentvenomsynergies },
    { ArchAngel, archangelsynergies },
    { Angela, angelasynergies },
    { AntMan, antmansynergies },
    { Beast, beastsynergies },
    { BlackBolt, blackboltsynergies },
    { BlackPanther, blackpanthersynergies },
    { BlackPantherCivilWar, bpcwsynergies },
    { BlackWidow, blackwidowsynergies },
    { Cable, cablesynergies },
    { CaptainAmerica, captainamericasynergies},
    { CaptainAmericaWWII, captainamericawwiisynergies },
    { CaptainMarvel, captainmarvelsynergies},
    { CivilWarrior, civilwarriorsynergies },
    { Colossus, colossussynergies },
    { Crossbones, crossbonessynergies },
    { Cyclops, cyclopssynergies},
    { CyclopsNew, cyclopsnewsynergies },
    { DareDevil, daredevilsynergies },
    { DareDevilClassic, daredevilclassicsynergies },
    { DeadpoolX, deadpoolxsynergies },
    { Drax, draxsynergies },
    { DoctorStrange, drstrangesynergies },
    { DoctorVoodoo, drvoodoosynergies },
    { DoctorOctopus, droctopussynergies },
    { Dormammu, dormammusynergies },
    { Electro, electrosynergies },
    { Elektra, elektrasynergies },
    { Falcon, falconsynergies },
    { Gambit, gambitsynergies },
    { Gamora, gamorasynergies },
    { GhostRider, ghostridersynergies },
    { Groot, grootsynergies },
    { Guillotine, guillotinesynergies },
    { Gwenpool, gwenpoolsynergies },
    { Hawkeye, hawkeyesynergies },
    { Hood, hoodsynergies },
    { HowardTheDuck, howardsynergies },
    { Hulk, hulksynergies },
    { Hulkbuster, hulkbustersynergies },
    { Hyperion, hyperionsynergies },
    { Iceman, icemansynergies },
    { IronFist, ironfistsynergies },
    { IronMan, ironmansynergies },
    { IronPatriot, ironpatriotsynergies },
    { JoeFixit, joefixitsynergies },
    { Juggernaut, juggernautsynergies },
    { Karnak, karnaksynergies },
    { Loki, lokisynergies },
    { LukeCage, lukecagesynergies },
    { Magik, magiksynergies },
    { Magneto, magnetosynergies },
    { MilesMorales, milesmoralessynergies },
    { MoonKnight, moonknightsynergies },
    { Mordo, mordosynergies },
    { MsMarvel, msmarvelsynergies },
    { Nightcrawler, nightcrawlersynergies },
    { OldManLogan, oldmanlogansynergies },
    { Phoenix, phoenixsynergies },
    { Punisher, punishersynergies },
    { Punisher2099, punisher2099synergies },
    { Quake, quakesynergies },
    { RedHulk, redhulksynergies },
    { Rhino, rhinosynergies },
    { RocketRaccoon, rocketsynergies },
    { Rogue, roguesynergies },
    { Ronan, ronansynergies },
    { ScarletWitch, scarletwitchsynergies },
    { SheHulk, shehulksynergies },
    { SpiderGwen, spidergwensynergies },
    { Spiderman, spidermansynergies },
    { SpidermanSymbiote, spidermansymbiotesynergies },
    { SpidermanStark, spidermanstarksynergies },
    { StarLord, starlordsynergies },
    { Storm, stormsynergies },
    { SuperiorIronMan, superiorironmansynergies },
    { Thor, thorsynergies },
    { ThorJaneFoster, thorjanefostersynergies },
    { Ultron, ultronsynergies },
    { UnstoppableColossus, unstoppablecolossussynergies },
    { Venom, venomsynergies },
    { VenomPool, venompoolsynergies },
    { Vision, visionsynergies },
    { VisionAgeOfUltron, visionaousynergies },
    { WarMachine, warmachinesynergies },
    { WinterSoldier, wintersoldiersynergies },
    { Wolverine, wolverinesynergies },
    { X23, x23synergies },
    { YellowJacket, yellowjacketsynergies },
    { Yondu, yondusynergies },
}

var MyHeroes = []Hero {
    { ArchAngel, archangelsynergies},
    { Beast, beastsynergies},
    { BlackBolt, blackboltsynergies },
    //{ BlackPanther, blackpanthersynergies },
    { BlackPantherCivilWar, bpcwsynergies },
    { BlackWidow, blackwidowsynergies },
    //{ Cable, cablesynergies },
    { CaptainAmerica, captainamericasynergies},
    { CaptainAmericaWWII, captainamericawwiisynergies },
    { CaptainMarvel, captainmarvelsynergies},
    { Colossus, colossussynergies },
    { DoctorStrange, drstrangesynergies },
    { DoctorVoodoo, drvoodoosynergies },
    { Dormammu, dormammusynergies },
    { Gambit, gambitsynergies },
    { Gamora, gamorasynergies },
    { GhostRider, ghostridersynergies },
    { Groot, grootsynergies },
    { Hulk, hulksynergies },
    { Hulkbuster, hulkbustersynergies },
    { IronMan, ironmansynergies },
    { IronPatriot, ironpatriotsynergies },
    { JoeFixit, joefixitsynergies },
    { Karnak, karnaksynergies },
    { Loki, lokisynergies },
    { LukeCage, lukecagesynergies },
    { MilesMorales, milesmoralessynergies },
    { MoonKnight, moonknightsynergies },
    { Mordo, mordosynergies },
    { OldManLogan, oldmanlogansynergies },
    { Punisher, punishersynergies },
    { Quake, quakesynergies },
    { ScarletWitch, scarletwitchsynergies },
    { SheHulk, shehulksynergies },
    { SpiderGwen, spidergwensynergies },
    { SpidermanSymbiote, spidermansymbiotesynergies },
    { StarLord, starlordsynergies },
    //{ Storm, stormsynergies },
    { Thor, thorsynergies },
    { Ultron, ultronsynergies },
    { VenomPool, venompoolsynergies },
    { WinterSoldier, wintersoldiersynergies },
    { X23, x23synergies },
    { YellowJacket, yellowjacketsynergies },
    { Yondu, yondusynergies },
}

func contains(team []Hero, item HeroVal) bool {
    for _, teammate := range(team) {
        if teammate.Name == item {
            return true
        }
    }
    
    return false
}

// mapping of hero names to hero values
var namesToValues map[string](HeroVal)

// mapping of hero structures by hero values
var heroByValue map[HeroVal]Hero

// initialize the namesToValues map
func initMappings() {
    namesToValues = make(map[string]HeroVal, 0)
    for ii := Abomination; ii < MaxHeroVal; ii++ {
        namesToValues[fmt.Sprintf("%v", ii)] = ii;
    }
}

// NameToValue takes a string representing a hero value and returns that value
func NameToValue(name string) HeroVal {
    if namesToValues == nil {
        initMappings()
    }

    return namesToValues[name]
}

// HeroByValue returns a Hero structure given a HeroVal
func HeroByValue(val HeroVal) Hero {
    if heroByValue == nil {
        heroByValue = make(map[HeroVal]Hero, len(Heroes))

        for _, hero := range Heroes {
            heroByValue[hero.Name] = hero
        }
    }

    ret, ok := heroByValue[val]
    if !ok {
        panic(fmt.Sprintf("%v(%d) not found", val, val))
    }

    return ret
}

// SynergyCount returns the number of synergies on a given team that are satisfied
func SynergyCount(team []Hero) int {
    var count int

    for _, hero := range team {
        for _, syns := range hero.Synergies {
            for _, synval := range syns.Vals {
                if contains(team, synval) {
                    count++
                    break
                }
            }
        }
    }

    return count
}

// FormatTeam prints a slice of heroes prettily
func FormatTeam(team []Hero) string {
    var arr []string

    for _, hero := range team {
        arr = append(arr, fmt.Sprintf("%s", hero.Name))
    }

    return strings.Join(arr, ", ")
}

// FormatTeamInfo prints a TeamInfo structure prettily
func FormatTeamInfo(team TeamInfo) string {
    return FormatTeam(team.Team)
}

// FormatTeamInfos prints a slice of TeamInfo structures prettily
func FormatTeamInfos(teams []TeamInfo) string {
    var ret string
    for idx, team := range teams {
        ret += fmt.Sprintf("%d ", idx) + FormatTeamInfo(team) + "\n"
    }
    
    return ret
}

// DeserializeTeamInfo is basically the reverse of FormatTeamInfo
func DeserializeTeamInfo(teamstr string) TeamInfo {
    var ret TeamInfo
    team := make([]Hero, 0)

    re := regexp.MustCompile("(\\S+), (\\S+), (\\S+)")
    results := re.FindStringSubmatch(teamstr)

    team = append(team, HeroByValue(NameToValue(results[1])))
    team = append(team, HeroByValue(NameToValue(results[2])))
    team = append(team, HeroByValue(NameToValue(results[3])))

    ret.Count = SynergyCount(team)
    ret.Team = team

    return ret
}

// DeserializeTeamInfos is basically the reverse of FormatTeamInfos
func DeserializeTeamInfos(fullstr string) []TeamInfo {
    ret := make([]TeamInfo, 0)

    teamlines := strings.Split(fullstr, "\n")

    for _, teamstr := range teamlines {
        split := strings.SplitN(teamstr, " ", 2)
        if len(split) != 2 {
            continue
        }
        teamstrstrip := split[1]

        ti := DeserializeTeamInfo(teamstrstrip)
        ret = append(ret, ti)
    }

    return ret
}

