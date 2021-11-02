package main

import . "github.com/c6h12o6/mcoc/globals"
import "github.com/c6h12o6/mcoc/oauth"

import "fmt"
import "strconv"
import "log"
import "net"
import "context"

import (
	"database/sql"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"

  "net/http"
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

  "sync"

	pb "github.com/c6h12o6/mcoc/proto"
)

var mcoc *mcocServer

func newServer() *mcocServer {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@cloudsql(homeproject:us-east1:champdb)/champdb?parseTime=true", sqlPassword))
	if err != nil {
		log.Fatalf("error creating db connection: %v\n", err)
	}
	s := &mcocServer{db: db}
	return s
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
  PlayerId int32 `json:"player_id"`
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

  sess := c.GetHeader("SESSION-ID")

  pbReq := pb.GetNodePreferencesRequest{AllianceId: req.AllianceId}
  resp, err := mcoc.GetNodePreferences(context.WithValue(context.Background(), "session", sess), &pbReq)
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

  sess := c.GetHeader("SESSION-ID")

  fmt.Printf("bound json %v\n", req)
  pbReq := pb.SetNodePreferencesRequest{AllianceId: req.AllianceId}
  pbReq.Nodes = make([]*pb.Node, len(req.Nodes))
  for idx, n := range req.Nodes {
    pbReq.Nodes[idx] = &pb.Node{NodeId: n.NodeId, ChampName: n.ChampName}
  }
  fmt.Printf("pb %v\n", pbReq)
  _, err := mcoc.SetNodePreferences(context.WithValue(context.Background(), "session", sess), &pbReq)
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

  sess := c.GetHeader("SESSION-ID")
  pbReq := pb.ListPlayersRequest{Alliance: req.AllianceID}
  resp, err := mcoc.ListPlayers(context.WithValue(context.Background(), "session", sess), &pbReq)
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
    
  sess := c.GetHeader("SESSION-ID")

  _, err := mcoc.SavePlayers(context.WithValue(context.Background(), "session", sess), &pbReq)
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
  sess := c.GetHeader("SESSION-ID")

  _, err := mcoc.AddChamp(context.WithValue(context.Background(), "session", sess), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    _, err := mcoc.UpdateChamp(context.WithValue(context.Background(), "session", sess), &pbReq)
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
  sess := c.GetHeader("SESSION-ID")

  _, err := mcoc.DelChamp(context.WithValue(context.Background(), "session", sess), &pbReq)
  if err != nil {
    fmt.Printf("lame: %v\n", err)
    return
  }

  c.String(200, "noice")
}

func getAllChamps(c *gin.Context) {
  pbReq := pb.GetAllChampsRequest{}

  sess := c.GetHeader("SESSION-ID")
  resp, err := mcoc.GetAllChamps(context.WithValue(context.Background(), "session", sess), &pbReq)
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

  sess := c.GetHeader("SESSION-ID")

  resp, err := mcoc.GetWarDefense(context.WithValue(context.Background(), "session", sess), &pbReq)
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

  sess := c.GetHeader("SESSION-ID")

  fmt.Printf("%+v\n", pbReq)
  resp, err := mcoc.ListChamps(context.WithValue(context.Background(), "session", sess), &pbReq)
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

type googleAuthPost struct {
  TokenId string `json:"token"`
}

type googleAuthPostResponse struct {
  Email string `json:"email"`
  SessionId string `json:"session_id"`
}

func postGoogleAuth(c *gin.Context) {
  var req googleAuthPost
  c.Writer.Header().Set("Content-Type", "application/json")
  fmt.Printf("Login thingy\n")
  if err := c.BindJSON(&req); err != nil {
    c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
    fmt.Printf("Error Binding Json: %v\n", err)
    return
  }
  fmt.Printf("found token: %v\n", req.TokenId)
  result, err := oauth.VerifyIdToken(req.TokenId)
  if err != nil {
    c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
    fmt.Printf("Error in oauth: %v\n", err)
    return
  }
  fmt.Printf("Result: %+v\n", result)
  fmt.Printf("Email: %v\n", result.Email)
  playerId, sessionId, err := mcoc.createSession(result.Email)
  if err != nil {
    fmt.Printf("error creating session: %v\n", err)
    c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
    return
  }

  allianceId, err := mcoc.getAlliance(playerId)
  if err != nil {
    c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
    fmt.Printf("Error getting alliance %v\n", err)
    return
  }

  c.JSON(200, gin.H{"sessionId": sessionId, "playerId": playerId, "allianceId": allianceId})
}

type createAlliancePost struct {
  Name string `json:"name"`
  PlayerId int32 `json:"player_id"`
}
func postCreateAlliance(c *gin.Context) {
  var req createAlliancePost
  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("Error Binding Json: %v\n", err)
    c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
    return
  }

  pbReq := pb.CreateAllianceRequest{Name: req.Name, PlayerId: req.PlayerId}
  fmt.Printf("Create alliance %+v\n", pbReq)

  sess := c.GetHeader("SESSION-ID")
  resp, err := mcoc.CreateAlliance(context.WithValue(context.Background(), "session", sess), &pbReq)
  if err != nil {
    fmt.Printf("bad: %v\n", err)
    c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
    return
  }
  m := jsonpb.Marshaler{}
  result, err := m.MarshalToString(resp)
  if err != nil {
    fmt.Printf("booo: %v\n", err)
    c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
    return
  }

  fmt.Printf("%+v\n", result)
  c.Writer.Header().Set("Content-Type", "application/json")
  c.String(200, result)
}

type allianceInfoPost struct {
  AllianceId int32 `json:"alliance_id"`
}

func postAllianceInfo(c *gin.Context) {
  var req allianceInfoPost
  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("Error Binding Json: %v\n", err)
    return
  }

  fmt.Printf("alliance info request for %v\n", req.AllianceId)

  pbReq := pb.GetAllianceInfoRequest{AllianceId: req.AllianceId}

  sess := c.GetHeader("SESSION-ID")
  resp, err := mcoc.GetAllianceInfo(context.WithValue(context.Background(), "session", sess), &pbReq)
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
type playerInfoPost struct {
  SessionId string `json:"session_id"`
}

func postPlayerInfo(c *gin.Context) {
  var req playerInfoPost
  if err := c.BindJSON(&req); err != nil {
    fmt.Printf("Error Binding Json: %v\n", err)
    return
  }

  fmt.Printf("player info request for %v\n", req.SessionId)

  pbReq := pb.GetPlayerInfoRequest{SessionId: req.SessionId}
  sess := c.GetHeader("SESSION-ID")
  resp, err := mcoc.GetPlayerInfo(context.WithValue(context.Background(), "session", sess), &pbReq)
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

//func (s *mcocServer) ListChamps(ctx context.Context, req *pb.ListChampsRequest) (*pb.ListChampsResponse, error) {
  
func main() {
	lis, err := net.Listen("tcp", "localhost:4567")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

  defenseLock = sync.Mutex{}

  r := gin.Default()
  r.Use(cors.Middleware(cors.Config{
    Origins:        "*",
    Methods:        "GET, PUT, POST, DELETE",
    RequestHeaders: "Origin, Authorization, Content-Type, Session-id",
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
  r.POST("/CreateAlliance", postCreateAlliance)
  r.POST("/AllianceInfo", postAllianceInfo)
  r.POST("/PlayerInfo", postPlayerInfo)
  r.POST("/api/v1/auth/google", postGoogleAuth)
  go r.Run("0.0.0.0:8080")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
  mcoc = newServer()
	pb.RegisterMcocServiceServer(grpcServer, mcoc)
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}
