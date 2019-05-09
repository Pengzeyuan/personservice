package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "personservice/tutorial"

	"google.golang.org/grpc"
)

// 个人信息服务端
type personServer struct {
	persons sync.Map
}

// AddPerson 添加一个个人信息
func (s *personServer) AddPerson(ctx context.Context, person *pb.Person) (*pb.Result, error) {
	s.persons.LoadOrStore(person.Name, person)
	return &pb.Result{
		Success: true,
	}, nil
}

// AddPersons 添加多个个人信息
func (s *personServer) AddPersons(stream pb.Manage_AddPersonsServer) error {
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Result{
				Success: true,
			})
		}

		if err != nil {
			return err
		}

		s.persons.LoadOrStore(person.Name, person)
	}
}

// GetPersonsLimit 获取限定数目的个人信息
func (s *personServer) GetPersonsLimit(limitNum *pb.ReqNum, stream pb.Manage_GetPersonsLimitServer) error {
	var err error
	var i int32
	s.persons.Range(func(key, value interface{}) bool {
		person, ok := value.(*pb.Person)
		if !ok {
			return false
		}
		err = stream.Send(person)
		if err != nil {
			return false
		}
		i++
		if i >= (limitNum.Num) {
			return false
		}
		return true
	})
	return err
}

// GetPersons 获取给定名字的所有个人信息
func (s *personServer) GetPersons(stream pb.Manage_GetPersonsServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		value, ok := s.persons.Load(in.Name)
		if !ok {
			continue
		}
		person, ok := value.(*pb.Person)
		if !ok {
			continue
		}
		err = stream.Send(person)
		if err != nil {
			return err
		}
	}
}

func newServer() *personServer {
	s := &personServer{}
	return s
}

func main() {
	address := "localhost:50001"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterManageServer(grpcServer, newServer())
	fmt.Println("Server listening on:", address)
	grpcServer.Serve(lis)
}
