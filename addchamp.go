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

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	//"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/examples/data"

	//"github.com/golang/protobuf/proto"

	pb "github.com/c6h12o6/mcoc/proto"
)

var sqlPassword = os.Getenv("CLOUD_SQL_PASSWORD")

func (s *mcocServer) getPlayerId(player string) (int, error) {
	var playerId int
	err := s.db.QueryRow("Select id from players where name like ?", player).Scan(&playerId)
	return playerId, err
}

func (s *mcocServer) addChamp(player string, champ HeroVal, stars int, rank int, sig int, update bool) error {

	playerId, err := s.getPlayerId(player)
	fmt.Printf("Player %v is %v\n", player, playerId)

	if update {
		_, err = s.db.Exec("update champ set stars=?, herorank=?, signature=? where player=? and heroval=?",
			stars, rank, sig, playerId, champ)
	} else {
		_, err = s.db.Exec("Insert into champ (player, heroval, stars, herorank, signature) Values(?, ?, ?, ?, ?)",
			playerId, champ, stars, rank, sig)
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

type mcocServer struct {
	pb.UnimplementedMcocServiceServer
	db *sql.DB
}

func (s *mcocServer) UpdateChamp(ctx context.Context, req *pb.AddChampRequest) (*pb.AddChampResponse, error) {
	log.Printf("Add Champ Called")
	champValue := NameToValue(req.Champ.ChampName) // TODO this code currently panics on invalid name
	err := s.addChamp(req.Player, champValue, int(req.Champ.Stars), int(req.Champ.Rank), int(req.Champ.Sig), true)
	return &pb.AddChampResponse{}, err
}

func (s *mcocServer) AddChamp(ctx context.Context, req *pb.AddChampRequest) (*pb.AddChampResponse, error) {
	log.Printf("Add Champ Called")
	champValue := NameToValue(req.Champ.ChampName) // TODO this code currently panics on invalid name
	err := s.addChamp(req.Player, champValue, int(req.Champ.Stars), int(req.Champ.Rank), int(req.Champ.Sig), false)
	return &pb.AddChampResponse{}, err
}

func (s *mcocServer) LockChamp(ctx context.Context, req *pb.LockChampRequest) (*pb.LockChampResponse, error) {
	log.Printf("Lock champ called")
	return &pb.LockChampResponse{}, nil
}

func (s *mcocServer) ListChamps(ctx context.Context, req *pb.ListChampsRequest) (*pb.ListChampsResponse, error) {
	log.Printf("List champs called")
	ret := pb.ListChampsResponse{}

	playerId, err := s.getPlayerId(req.Player)
	if err != nil {
		return &ret, err
	}

	rows, err := s.db.Query("select heroval, stars, herorank, signature, locked from champ where player = ?", playerId)
	if err != nil {
		return &ret, err
	}

	for rows.Next() {
		var c pb.Champ
		var hv HeroVal
		err := rows.Scan(&hv, &c.Stars, &c.Rank, &c.Sig, &c.LockedNode)
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
	result := war.BestWarDefense(int(req.Alliance), int(req.Bg))
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

func newServer() *mcocServer {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@cloudsql(homeproject:us-east1:champdb)/champdb", sqlPassword))
	if err != nil {
		log.Fatalf("error creating db connection: %v\n", err)
	}
	s := &mcocServer{db: db}
	return s
}

func main() {
	/*
	   fmt.Printf("%v\n", os.Args)
	   player := os.Args[1]
	   champ := os.Args[2]
	   stars := os.Args[3]
	   rank := os.Args[4]
	   sig := os.Args[5]

	   champValue := NameToValue(champ)
	   err := addChamp(player, champValue,
	     ToInt(stars), ToInt(rank), ToInt(sig))
	   fmt.Printf("%v\n", err)
	*/
	lis, err := net.Listen("tcp", "localhost:4567")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMcocServiceServer(grpcServer, newServer())
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}
