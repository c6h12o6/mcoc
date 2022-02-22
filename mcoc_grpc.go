package main

import . "github.com/c6h12o6/mcoc/globals"
import "github.com/c6h12o6/mcoc/war"

import "fmt"
import "os"
import "strconv"
import "log"
import "context"

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

  //pbtools "github.com/protocolbuffers/protobuf-go"

	//"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/examples/data"

	//"github.com/golang/protobuf/proto"

  "time"

  "sync"
  "github.com/google/uuid"
  "strings"

	pb "github.com/c6h12o6/mcoc/proto"
)

var sqlPassword = os.Getenv("CLOUD_SQL_PASSWORD")
var sqlHost = os.Getenv("SQL_HOST")

type mcocServer struct {
	pb.UnimplementedMcocServiceServer
	db *sql.DB
}

var defenseLock sync.Mutex

func (s *mcocServer) createPlayer(email string) (int, error) {
  fmt.Printf("create player for %v\n", email)
  _, err := s.db.Exec("Insert into players (suicides, mystic_dispersion, alliance, bg, name, email, access_level) values (?, ?, ?, ?, ?, ?, ?)", 0, 0, 0, 4, email, email, 0)
  if err != nil { return 0, err }

  fmt.Printf("inserted\n")
  pid, err := s.getPlayerIdFromEmail(email) 
  return pid, err
}


func (s *mcocServer) createSession(email string) (int, string, error) {
  playerId, err := s.getPlayerIdFromEmail(email)
  if (err != nil) {
    if strings.Contains(err.Error(), "no rows in result set") {
      fmt.Printf("No user found\n")
      playerId, err = s.createPlayer(email)
      if err != nil {
        return 0, "", err
      }
    } else {
      return 0, "", err
    }
  }
  sessionId := uuid.New().String()
  now := time.Now().UTC().Format("2006-01-02 03:04:05")
	_, err = s.db.Exec("Insert into sessions (id, player_id, created) values(?, ?, ?)", sessionId, playerId, now)
  if err != nil {
    return 0, "", err
  }
	return playerId, sessionId, err
}

func (s *mcocServer) CreateAlliance(ctx context.Context, req *pb.CreateAllianceRequest) (*pb.CreateAllianceResponse, error) {
  now := time.Now().UTC().Format("2006-01-02 03:04:05")
  reference := uuid.New().String()

	result, err := s.db.Exec("Insert into alliances (name, created, reference) values(?, ?, ?)", req.Name, now, reference)
  if err != nil {
    return nil, err
  }

  allianceId, err := result.LastInsertId()
  if err != nil {
    return nil, err
  }
  fmt.Printf("Created alliance %v\n", allianceId)
  _, err = s.db.Exec("Update players set alliance=?, bg=1 where id=?", allianceId, req.PlayerId)
  if err != nil {
    return nil, err
  }

  for idx := 1; idx <= 55; idx++ {
    _, err = s.db.Exec("Insert into nodes (alliance, node) values(?, ?)", allianceId, idx)
    if err != nil {
      return nil, err
    }
  }

  return &pb.CreateAllianceResponse{AllianceId: int32(allianceId)}, nil
}

func (s *mcocServer) GetAllianceInfo(ctx context.Context, req *pb.GetAllianceInfoRequest) (*pb.GetAllianceInfoResponse, error) {
  resp := pb.GetAllianceInfoResponse{}

	err := s.db.QueryRow("Select name, reference from alliances where id like ?", req.AllianceId).Scan(
    &resp.Name,
    &resp.Reference)
  if err != nil {
    return nil, err
  }

  return &resp, nil
}

func (s *mcocServer) GetPlayerInfo(ctx context.Context, req *pb.GetPlayerInfoRequest) (*pb.GetPlayerInfoResponse, error) {
  resp := pb.GetPlayerInfoResponse{}

  fmt.Printf("Get Player id from session id: %v\n", req.SessionId)
  pid, err := s.getPlayerIdFromSession(req.SessionId)
  if err != nil {
    return nil, err
  }
  fmt.Printf("Getting info for player id: %v\n", pid)
	err = s.db.QueryRow("Select name, suicides, mystic_dispersion from players where id like ?", pid).Scan(
    &resp.Name,
    &resp.Suicides,
    &resp.MysticDispersion,
  )
  if err != nil {
    return nil, err
  }
  fmt.Printf("resp.Suicides was %v %v\n", resp.Suicides, resp.MysticDispersion)

  return &resp, nil
}

func (s *mcocServer) SetPlayerInfo(ctx context.Context, req *pb.SetPlayerInfoRequest) (*pb.SetPlayerInfoResponse, error) {
  resp := pb.SetPlayerInfoResponse{}

  fmt.Printf("Set Player id from session id: %v\n", req.SessionId)
  pid, err := s.getPlayerIdFromSession(req.SessionId)
  if err != nil {
    return nil, err
  }
  fmt.Printf("Setting info for player id: %v\n", pid)
  _, err = s.db.Exec("Update players Set name=?, suicides=?, mystic_dispersion=? where id like ?",
      req.Name,
      req.Suicides,
      req.MysticDispersion,
      pid)
  if err != nil {
    return nil, err
  }

  return &resp, nil
}


func (s *mcocServer) getPlayerIdFromSession(sessionId string) (int, error) {
	var playerId int
  var created time.Time

	err := s.db.QueryRow("Select player_id, created from sessions where id like ?", sessionId).Scan(&playerId, &created)
  if err != nil {
    return 0, err
  }

  fmt.Printf("%v\n%v\n%v\n", created, time.Now(), created.Add(2* time.Hour))
  if created.Add(48 * time.Hour).Before(time.Now()) {
    return 0, fmt.Errorf("Session Expired")
  }

	return playerId, err
}

func (s *mcocServer) getPlayerIdFromChampId(champId int) (int, error) {
	var playerId int

	err := s.db.QueryRow("Select player from champ where id like ?", champId).Scan(&playerId)
  if err != nil {
    return 0, err
  }

	return playerId, err
}

func (s *mcocServer) getPlayerIdFromEmail(email string) (int, error) {
	var playerId int
	err := s.db.QueryRow("Select id from players where email like ?", email).Scan(&playerId)
	return playerId, err
}

func (s *mcocServer) getPlayerId(player string) (int, error) {
	var playerId int
	err := s.db.QueryRow("Select id from players where name like ?", player).Scan(&playerId)
	return playerId, err
}

func (s *mcocServer) getAlliance(playerId int) (int, error) {
	var allianceId int
	err := s.db.QueryRow("Select alliance from players where id = ?", playerId).Scan(&allianceId)
	return allianceId, err
}

func (s *mcocServer) getAccessLevelForPlayer(playerId int) (int, error) {
	var access int
	err := s.db.QueryRow("Select access_level from players where id = ?", playerId).Scan(&access)
	return access, err
}

func (s *mcocServer) getPlayerName(id int) (string, error) {
	var playerName string
	err := s.db.QueryRow("Select name from players where id = ?", id).Scan(&playerName)
	return playerName, err
}

func (s *mcocServer) authorizeRequest(sessionId string, playerId int, allianceId int, mutatePlayer bool, mutateAlliance bool) bool {
  fmt.Printf("trying to authorize from session %v\n", sessionId)
  pid, err := s.getPlayerIdFromSession(sessionId)
  if err != nil {
    fmt.Printf("Unable to get playerid from session %s: %v\n", sessionId, err)
    return false;
  }

  access, err := s.getAccessLevelForPlayer(pid)
  if err != nil {
    fmt.Printf("Unable to get access level for player %v: %v\n", pid, err)
    return false
  }

  userAlliance, err := s.getAlliance(pid)
  if err != nil {
    fmt.Printf("Unable to get alliance id for player %v: %v\n", pid, err)
    return false
  }

  var targetAlliance int
  if playerId != 0 {
    targetAlliance, err = s.getAlliance(playerId)
    if err != nil {
      fmt.Printf("Unable to get alliance id for player %v: %v\n", pid, err)
      return false
    }
  }

  if mutatePlayer && !(playerId != pid || (access > 0 && userAlliance == targetAlliance)) {
    fmt.Printf("player %v is unable to mutate player %v\n", pid, playerId)
    return false
  }

  if mutateAlliance && (userAlliance != allianceId || access == 0) {
    fmt.Printf("Player %v in alliance %v is unable to mutate alliance %v\n", pid, userAlliance, allianceId)
    return false
  }

  return true
}

func (s *mcocServer) addChamp(playerId int, champ HeroVal, stars int, rank int, sig int, node int, update bool) error {

  var err error
	if update {
		_, err = s.db.Exec("update champ set stars=?, herorank=?, signature=?, locked=?, deleted=0 where player=? and heroval=?",
			stars, rank, sig, node, playerId, champ)
	} else {
    fmt.Printf("Adding %v to %v\n", champ, playerId)
		_, err = s.db.Exec("Insert into champ (player, heroval, stars, herorank, signature, locked) Values(?, ?, ?, ?, ?, ?)",
			playerId, champ, stars, rank, sig, node)
	}
	if err != nil {
		return err
	}

	return nil
}

type Player struct {
	Id       int
	Suicides bool
	MD       int
	Name     string
}

func ToInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	return int(i)
}

func (s *mcocServer) UpdateChamp(ctx context.Context, req *pb.AddChampRequest) (*pb.AddChampResponse, error) {
	log.Printf("Add Champ Called\n")
	champValue := NameToValue(req.Champ.ChampName) // TODO this code currently panics on invalid name
  if champValue == Empty {
    return nil, fmt.Errorf("No champ named %v", req.Champ.ChampName)
  }
  playerId := int(req.GetId())
  if playerId == 0 {
    playerId, _ = s.getPlayerId(req.GetPlayer())
    fmt.Printf("Player %v is %v\n", req.GetPlayer(), playerId)
  }
	err := s.addChamp(playerId, champValue, int(req.Champ.Stars), int(req.Champ.Rank), int(req.Champ.Sig), int(req.Champ.LockedNode), true)
	return &pb.AddChampResponse{}, err
}

func (s *mcocServer) AddChamp(ctx context.Context, req *pb.AddChampRequest) (*pb.AddChampResponse, error) {
	fmt.Printf("Add Champ Called")

  playerId := int(req.GetId())
  if playerId == 0 {
    playerId, _ = s.getPlayerId(req.GetPlayer())
    fmt.Printf("Player %v is %v\n", req.GetPlayer(), playerId)
  }

  if !s.authorizeRequest(ctx.Value("session").(string), playerId, 0, true, false) {
    return nil, fmt.Errorf("Access Denied")
  }
	champValue := NameToValue(req.Champ.ChampName) // TODO this code currently panics on invalid name
  if champValue == Empty {
    return nil, fmt.Errorf("No champ named %v", req.Champ.ChampName)
  }
	err := s.addChamp(playerId, champValue, int(req.Champ.Stars), int(req.Champ.Rank), int(req.Champ.Sig), int(req.Champ.LockedNode), false)
	return &pb.AddChampResponse{}, err
}

func (s *mcocServer) DelChamp(ctx context.Context, req *pb.DelChampRequest) (*pb.DelChampResponse, error) {
	fmt.Printf("Del Champ Called")
  playerId, err := s.getPlayerIdFromChampId(int(req.Id))
  if err != nil {
    return nil, err
  }

  if !s.authorizeRequest(ctx.Value("session").(string), playerId, 0, true, false) {
    return nil, fmt.Errorf("Access Denied")
  }
	_, err = s.db.Exec("update champ set deleted = 1 where id = ?", req.Id)
	return &pb.DelChampResponse{}, err
}

func (s *mcocServer) LockChamp(ctx context.Context, req *pb.LockChampRequest) (*pb.LockChampResponse, error) {
	log.Printf("Lock champ called")

	return &pb.LockChampResponse{}, nil
}

func (s *mcocServer) ListPlayers(ctx context.Context, req *pb.ListPlayersRequest) (*pb.ListPlayersResponse, error) {
	log.Printf("List players called")
	ret := pb.ListPlayersResponse{}


	rows, err := s.db.Query("select name, id, bg FROM players where alliance = ?", req.Alliance)
	if err != nil {
		return &ret, err
	}

  type response struct {
    name string
    id int
    bg int
  }


	for rows.Next() {
    var resp pb.PlayerData

		err := rows.Scan(&resp.Name, &resp.Id, &resp.Bg)
		if err != nil {
			return &ret, err
		}

		ret.Players = append(ret.Players, &resp)
	}

	return &ret, nil
}

func (s *mcocServer) SavePlayers(ctx context.Context, req *pb.SavePlayersRequest) (*pb.SavePlayersResponse, error) {
  log.Printf("Save players called\n")

  if !s.authorizeRequest(ctx.Value("session").(string), 0, int(req.Alliance), false, true) {
    return nil, fmt.Errorf("Access Denied")
  }
  for _, p := range req.Players {
    if p.Id < 0 {
      fmt.Printf("%v is new\n", p.Name)
      // This is a new player, insert instead of updating
      _, err := s.db.Exec(`Insert into players
                            (name, alliance, bg, suicides, mystic_dispersion)
                           Values (?, ?, ?, ?, ?)`, p.Name, req.Alliance, p.Bg, 0, 0)
      if err != nil {
        return nil, err
      }
    } else {
      _, err := s.db.Exec("update players set bg=? where id=?", p.Bg, p.Id)
      if err != nil {
        return nil, err
      }
    }
  }

  return &pb.SavePlayersResponse{}, nil
}

func (s *mcocServer) ListChamps(ctx context.Context, req *pb.ListChampsRequest) (*pb.ListChampsResponse, error) {
	log.Printf("List champs called")
	ret := pb.ListChampsResponse{}

  var playerId int
  if req.GetId() != 0 {
    playerId = int(req.GetId())
    name, err := s.getPlayerName(int(req.GetId()))
    if err != nil {
      return &ret, err
    }
    ret.Player = name
  } else {
    id, err := s.getPlayerId(req.GetPlayer())
    if err != nil {
      return &ret, err
    }
    playerId = id
    ret.Player = req.GetPlayer()
  }

  if !s.authorizeRequest(ctx.Value("session").(string), playerId, 0, false, false) {
    return nil, fmt.Errorf("Access Denied")
  }


	rows, err := s.db.Query("select id, heroval, stars, herorank, signature, locked from champ where player = ? and deleted != 1", playerId)
	if err != nil {
		return &ret, err
	}

	for rows.Next() {
		var c pb.Champ
		var hv HeroVal
		err := rows.Scan(&c.Id, &hv, &c.Stars, &c.Rank, &c.Sig, &c.LockedNode)
		if err != nil {
			return &ret, err
		}
		c.ChampName = hv.String()
		ret.Champs = append(ret.Champs, &c)
	}
	return &ret, nil
}

func (s *mcocServer) GetWarDefense(ctx context.Context, req *pb.GetWarDefenseRequest) (*pb.GetWarDefenseResponse, error) {
	var ret pb.GetWarDefenseResponse

  if !s.authorizeRequest(ctx.Value("session").(string), 0, int(req.Alliance), false, true) {
    return nil, fmt.Errorf("Access Denied")
  }
  fmt.Printf("calling war defense\n")
  fmt.Printf("acquiring\n")
  defenseLock.Lock()
  defer defenseLock.Unlock()
	result, err := war.BestWarDefense(int(req.Alliance), int(req.Bg))
  if err != nil {
    return nil, err
  }
	for _, pd := range result {
		assignment := pb.Assignment{
			Player: pd.Player,
		}
		for _, c := range pd.Defenders.Champs {
			var pbc pb.Champ
			pbc.ChampName = c.Champ.String()
			pbc.Stars = c.Stars
			pbc.Rank = c.Level
			pbc.Sig = c.Sig
			pbc.LockedNode = int32(c.LockedNode)
      pbc.AssignedNode = c.AssignedNode
			assignment.Champs = append(assignment.Champs, &pbc)
		}
		ret.Assignments = append(ret.Assignments, &assignment)
	}
	return &ret, nil
}

func (s *mcocServer) GetAllChamps(ctx context.Context, req *pb.GetAllChampsRequest) (*pb.GetAllChampsResponse, error) {
  resp := make([]string, MaxHeroVal)
  for h := Abomination; h < MaxHeroVal; h++ {
    resp[h] = fmt.Sprintf("%v", h)
  }
  fmt.Printf("%v\n", resp)
  return &pb.GetAllChampsResponse{Names: resp}, nil
}

func (s *mcocServer) GetNodePreferences(ctx context.Context, req *pb.GetNodePreferencesRequest) (*pb.GetNodePreferencesResponse, error) {
  fmt.Printf("GetNodePreferences")

  if !s.authorizeRequest(ctx.Value("session").(string), 0, int(req.AllianceId), false, false) {
    return nil, fmt.Errorf("Access Denied")
  }

  ret := pb.GetNodePreferencesResponse{}
	rows, err := s.db.Query("select node, champ1, champ2, champ3, champ4, champ5 from nodes where alliance = ?", req.AllianceId)
	if err != nil {
		return &ret, err
	}

  type Node struct {
    node int
    ids []int
  }
	for rows.Next() {
    n := Node{ids: make([]int, 5)}
		err := rows.Scan(&n.node,
                     &n.ids[0],
                     &n.ids[1],
                     &n.ids[2],
                     &n.ids[3],
                     &n.ids[4])
		if err != nil {
			return &ret, err
    }
    ret.Nodes = append(ret.Nodes, &pb.Node{
      NodeId: int32(n.node),
      ChampName: []string{
        fmt.Sprintf("%v", HeroVal(n.ids[0])),
        fmt.Sprintf("%v", HeroVal(n.ids[1])),
        fmt.Sprintf("%v", HeroVal(n.ids[2])),
        fmt.Sprintf("%v", HeroVal(n.ids[3])),
        fmt.Sprintf("%v", HeroVal(n.ids[4])),
      }})
  }

  return &ret, nil
}

func (s *mcocServer) SetNodePreferences(ctx context.Context, req *pb.SetNodePreferencesRequest) (*pb.SetNodePreferencesResponse, error) {
  fmt.Printf("SetNodePreferences")
  ret := pb.SetNodePreferencesResponse{}

  if !s.authorizeRequest(ctx.Value("session").(string), 0, int(req.AllianceId), false, true) {
    return nil, fmt.Errorf("Access Denied")
  }

  for _, n := range req.Nodes {
    fmt.Printf("setting node %v to %v %v %v", n.NodeId, n.ChampName[0], n.ChampName[1], n.ChampName[2])
    _, err := s.db.Exec("Update nodes set champ1=?, champ2=?, champ3=?, champ4=?, champ5=? where alliance = ? and node=?", 
    NameToValue(n.ChampName[0]),
    NameToValue(n.ChampName[1]), 
    NameToValue(n.ChampName[2]),
    NameToValue(n.ChampName[3]),
    NameToValue(n.ChampName[4]), 
    req.AllianceId, n.NodeId)
    if err != nil {
      return &ret, err
    }
  }

  return &ret, nil
}

