package main

import . "github.com/c6h12o6/mcoc/globals"
import "github.com/c6h12o6/mcoc/war"

import "fmt"
import "os"
import "strconv"
import "log"
import "net"
import "context"

import (
	"database/sql"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"

  //"net/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

  //pbtools "github.com/protocolbuffers/protobuf-go"

	//"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/examples/data"

	"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"

  "github.com/gin-gonic/gin"
  "github.com/itsjamie/gin-cors"

  "time"

	pb "github.com/c6h12o6/mcoc/proto"
)

var sqlPassword = os.Getenv("CLOUD_SQL_PASSWORD")

type mcocServer struct {
	pb.UnimplementedMcocServiceServer
	db *sql.DB
}

var mcoc *mcocServer

func (s *mcocServer) getPlayerId(player string) (int, error) {
	var playerId int
	err := s.db.QueryRow("Select id from players where name like ?", player).Scan(&playerId)
	return playerId, err
}

func (s *mcocServer) getPlayerName(id int) (string, error) {
	var playerName string
	err := s.db.QueryRow("Select name from players where id = ?", id).Scan(&playerName)
	return playerName, err
}

func (s *mcocServer) addChamp(playerId int, champ HeroVal, stars int, rank int, sig int, node int, update bool) error {

  var err error
	if update {
		_, err = s.db.Exec("update champ set stars=?, herorank=?, signature=? locked=?, deleted=0 where player=? and heroval=?",
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
	champValue := NameToValue(req.Champ.ChampName) // TODO this code currently panics on invalid name
  if champValue == Empty {
    return nil, fmt.Errorf("No champ named %v", req.Champ.ChampName)
  }
  playerId := int(req.GetId())
  if playerId == 0 {
    playerId, _ = s.getPlayerId(req.GetPlayer())
    fmt.Printf("Player %v is %v\n", req.GetPlayer(), playerId)
  }
	err := s.addChamp(playerId, champValue, int(req.Champ.Stars), int(req.Champ.Rank), int(req.Champ.Sig), int(req.Champ.LockedNode), false)
	return &pb.AddChampResponse{}, err
}

func (s *mcocServer) DelChamp(ctx context.Context, req *pb.DelChampRequest) (*pb.DelChampResponse, error) {
	fmt.Printf("Del Champ Called")
	_, err := s.db.Exec("update champ set deleted = 1 where id = ?", req.Id)
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
  fmt.Printf("calling war defense")
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

func newServer() *mcocServer {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@cloudsql(homeproject:us-east1:champdb)/champdb", sqlPassword))
	if err != nil {
		log.Fatalf("error creating db connection: %v\n", err)
	}
	s := &mcocServer{db: db}
	return s
}

func (s *mcocServer) GetNodePreferences(ctx context.Context, req *pb.GetNodePreferencesRequest) (*pb.GetNodePreferencesResponse, error) {
  fmt.Printf("GetNodePreferences")
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

type getNodePreferencesPost struct {
  AllianceId int32 `json:"alliance_id"`
}

type Node struct {
  NodeId int32 `json:"nodeId"`
  ChampName []string `json:"champName"`
}

type setNodePreferencesPost struct {
  AllianceId int32 `json:"alliance_id"`
  Nodes []Node `json:"nodes"`

}
type listChampsPost struct {
  Name string `json:"name"`
  PlayerId int32 `json:"player_id"`
}

type listPlayersPost struct {
  AllianceID int32 `json:"alliance_id"`
}

type addChampPost struct {
  PlayerId int32 `json:"player_id"`
  ChampName string `json:"champ_name"`
  Stars int32 `json:"stars"`
  Rank int32 `json:"rank"`
  Sigs int32 `json:"sigs"`
  Node int32 `json:"node"`
}

type delChampPost struct {
  Id int32 `json:"id"`
}

type getWarDefensePost struct {
  AllianceId int32 `json:"alliance_id"`
  Bg int32 `json:"bg"`
}

func postGetNodePreferences(c *gin.Context) {
  var req getNodePreferencesPost

  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("hiii\n")
    return
  }

  pbReq := pb.GetNodePreferencesRequest{AllianceId: req.AllianceId}
  resp, err := mcoc.GetNodePreferences(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    return
  }

  m := jsonpb.Marshaler{}
  result, err := m.MarshalToString(resp)
  if err != nil {
    fmt.Printf("booo: %v\n", err)
    return
  }

  fmt.Printf("%+v\n", result)
  fmt.Printf("request: %+v\n", req)
  c.Writer.Header().Set("Content-Type", "application/json")
  c.Writer.Header().Set("Access-Control-Allow-Origin", "*") 
  c.String(200, result)
}

func postSetNodePreferences(c *gin.Context) {
  var req setNodePreferencesPost

  fmt.Printf("set node preferences\n")
  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("hiii %v\n", err)
    return
  }

  fmt.Printf("bound json %v\n", req)
  pbReq := pb.SetNodePreferencesRequest{AllianceId: req.AllianceId}
  pbReq.Nodes = make([]*pb.Node, len(req.Nodes))
  for idx, n := range req.Nodes {
    pbReq.Nodes[idx] = &pb.Node{NodeId: n.NodeId, ChampName: n.ChampName}
  }
  fmt.Printf("pb %v\n", pbReq)
  _, err := mcoc.SetNodePreferences(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    return
  }

  c.String(200, "noice")
}


func postListPlayers(c *gin.Context) {
  var req listPlayersPost

  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("hiii\n")
    return
  }

  pbReq := pb.ListPlayersRequest{Alliance: req.AllianceID}
  resp, err := mcoc.ListPlayers(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    return
  }

  m := jsonpb.Marshaler{}
  result, err := m.MarshalToString(resp)
  if err != nil {
    fmt.Printf("booo: %v\n", err)
    return
  }

  fmt.Printf("%+v\n", result)
  //c.IndentedJson(pbtools.Format(resp))

  fmt.Printf("request: %+v\n", req)
  c.Writer.Header().Set("Content-Type", "application/json")
  c.Writer.Header().Set("Access-Control-Allow-Origin", "*") 
  c.String(200, result)

}

type playerInfo struct {
  Name string `json:"name"`
  Id string `json:"id"`
  Bg int32 `json:"bg"`
}

type savePlayersPost struct {
  AllianceID int32 `json:"alliance_id"`
  Players []playerInfo `json:"players"`
}

func postSavePlayers(c *gin.Context) {
  var req savePlayersPost
  fmt.Printf("Save players post\n")

  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("hiii\n")
    return
  }

  pbReq := pb.SavePlayersRequest{
      Alliance: req.AllianceID,
      Players: make([]*pb.PlayerData, len(req.Players)),
  }

  for idx, p := range req.Players {
    id, err := strconv.Atoi(p.Id)
    if err != nil {
      c.String(400, fmt.Sprintf("bad id: %v\n", int32(id)))
      return
    }
    pbReq.Players[idx] = &pb.PlayerData{Name: p.Name, Id: int32(id), Bg: p.Bg}
  }
    
  _, err := mcoc.SavePlayers(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    return
  }

  c.String(200, "noice")
}

func postAddChamp(c *gin.Context) {
  var req addChampPost
  fmt.Printf("Add champ post\n")

  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("Error Binding Json: %v\n", err)
    return
  }

  pbReq := pb.AddChampRequest{
    Champ: &pb.Champ {
      ChampName: req.ChampName,
      Stars: req.Stars,
      Rank: req.Rank,
      Sig: req.Sigs,
      LockedNode: req.Node,
    },
    Identifier: &pb.AddChampRequest_Id{req.PlayerId},
  }

  fmt.Printf("%+v\n", pbReq)
  _, err := mcoc.AddChamp(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    _, err := mcoc.UpdateChamp(context.TODO(), &pbReq)
    if err != nil {
      fmt.Printf("lame: %v\n", err)
      return
    }
  }

  c.String(200, "noice")
}

func postDelChamp(c *gin.Context) {
  var req delChampPost
  fmt.Printf("del champ post\n")

  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("Error Binding Json: %v\n", err)
    return
  }

  pbReq := pb.DelChampRequest{
    Id: req.Id,
  }

  fmt.Printf("%+v\n", pbReq)
  _, err := mcoc.DelChamp(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    return
  }

  c.String(200, "noice")
}

func getAllChamps(c *gin.Context) {
  pbReq := pb.GetAllChampsRequest{}
  resp, err := mcoc.GetAllChamps(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("bad: %v\n", err)
    return
  }
  m := jsonpb.Marshaler{}
  result, err := m.MarshalToString(resp)
  if err != nil {
    fmt.Printf("booo: %v\n", err)
    return
  }

  fmt.Printf("%+v\n", result)
  c.Writer.Header().Set("Content-Type", "application/json")
  c.String(200, result)
}

func postWarDefense(c *gin.Context) {
  pbReq := pb.GetWarDefenseRequest{}
  req := getWarDefensePost{}

  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("Error Binding Json: %v\n", err)
    return
  }

  pbReq.Alliance = req.AllianceId
  pbReq.Bg = req.Bg

  resp, err := mcoc.GetWarDefense(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("bad: %v\n", err)
    return
  }
  m := jsonpb.Marshaler{}
  result, err := m.MarshalToString(resp)
  if err != nil {
    fmt.Printf("booo: %v\n", err)
    return
  }

  fmt.Printf("%+v\n", result)
  c.Writer.Header().Set("Content-Type", "application/json")
  c.String(200, result)
}

func postListChamps(c *gin.Context) {
  var req listChampsPost
  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("hiii %v\n", err)
    return
  }
  fmt.Printf("req: %+v\n", req)

  var pbReq pb.ListChampsRequest
  if req.PlayerId != 0 {
    pbReq = pb.ListChampsRequest{Identifier: &pb.ListChampsRequest_Id{req.PlayerId}}
  } else {
    pbReq = pb.ListChampsRequest{Identifier: &pb.ListChampsRequest_Player{req.Name}}
  }

  fmt.Printf("%+v\n", pbReq)
  resp, err := mcoc.ListChamps(context.TODO(), &pbReq)
  if err != nil {
    fmt.Printf("that failed: %v\n", err)
    return
  }

  m := jsonpb.Marshaler{}
  result, err := m.MarshalToString(resp)
  if err != nil {
    fmt.Printf("booo: %v\n", err)
    return
  }

  fmt.Printf("%+v\n", result)

  fmt.Printf("request: %+v\n", req)
  c.String(200, result)
}

//func (s *mcocServer) ListChamps(ctx context.Context, req *pb.ListChampsRequest) (*pb.ListChampsResponse, error) {
  
func main() {
	lis, err := net.Listen("tcp", "localhost:4567")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

  r := gin.Default()
  r.Use(cors.Middleware(cors.Config{
    Origins:        "*",
      Methods:        "GET, PUT, POST, DELETE",
        RequestHeaders: "Origin, Authorization, Content-Type",
          ExposedHeaders: "",
            MaxAge: 50 * time.Second,
              Credentials: true,
                ValidateHeaders: false,
  }))
  r.POST("/ListPlayers", postListPlayers)
  r.POST("/ListChamps", postListChamps)
  r.POST("/SavePlayers", postSavePlayers)
  r.GET("/AllChamps", getAllChamps)
  r.POST("/AddChamp", postAddChamp)
  r.POST("/DelChamp", postDelChamp)
  r.POST("/GetNodePreferences", postGetNodePreferences)
  r.POST("/SetNodePreferences", postSetNodePreferences)
  r.POST("/GetWarDefense", postWarDefense)
  go r.Run("localhost:8080")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
  mcoc = newServer()
	pb.RegisterMcocServiceServer(grpcServer, mcoc)
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}
