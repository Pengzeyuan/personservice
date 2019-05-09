package main

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "personservice/tutorial"
	"time"

	"google.golang.org/grpc"
)

const (
	rpcTimeOut = 10
)

// addPerson 用于添加个人信息
func addPerson(client pb.ManageClient, person *pb.Person) bool {
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeOut*time.Second)
	defer cancel()
	res, err := client.AddPerson(ctx, person)
	if err != nil {
		log.Printf("client.AddPerson failed, error: %v\n", err)
		return false
	}
	return res.Success

}

// addPersons 用来添加多个个人信息
func addPersons(client pb.ManageClient, persons []*pb.Person) bool {
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeOut*time.Second)
	defer cancel()
	stream, err := client.AddPersons(ctx)
	if err != nil {
		log.Printf("client.AddPersons failed, error: %v\n", err)
		return false
	}
	for _, person := range persons {
		if err := stream.Send(person); err != nil {
			log.Printf("stream.Send failed, error: %v\n", err)
			return false
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("stream.CloseAndRecv failed, error: %v\n", err)
		return false
	}
	return res.Success
}

// getPersonsLimit 用来获取指定数目的个人信息
func getPersonsLimit(client pb.ManageClient, limitNum int32) ([]*pb.Person, error) {
	var persons []*pb.Person
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeOut*time.Second)
	defer cancel()
	num := pb.ReqNum{
		Num: limitNum,
	}
	stream, err := client.GetPersonsLimit(ctx, &num)
	if err != nil {
		log.Printf("client.GetPersonsLimit failed, error: %v\n", err)
		return persons, err
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("stream.Recv failed, error: %v\n", err)
			return persons, err
		}
		persons = append(persons, person)
	}

	return persons, nil
}

// getPersons 用来获取指定名字的所有个人信息
func getPersons(client pb.ManageClient, personNames []string) ([]*pb.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeOut*time.Second)
	defer cancel()
	stream, err := client.GetPersons(ctx)
	if err != nil {
		log.Printf("client.GetPersons failed, error: %v\n", err)
		return nil, err
	}
	waitc := make(chan struct{})
	// 发送个人名字信息
	go func() {
		for _, personName := range personNames {
			name := pb.ReqName{
				Name: personName,
			}
			if err := stream.Send(&name); err != nil {
				log.Printf("stream.Send failed, error: %v\n", err)
				break
			}
		}
		err := stream.CloseSend()
		if err != nil {
			log.Printf("stream.CloseSend failed, error: %v\n", err)
		}
		close(waitc)
	}()
	// 获取对应的所有个人信息
	var persons []*pb.Person
	var in *pb.Person
	for {
		in, err = stream.Recv()
		if err != nil {
			break
		}
		persons = append(persons, in)
	}

	<-waitc
	// 检查读取结果, err应该不会为nil
	if err == io.EOF || err == nil {
		return persons, nil
	}
	log.Fatalf("stream.Recv failed, error: %v\n", err)
	return persons, err
}

func makePerson(name string, id int32, email string) pb.Person {
	return pb.Person{
		Name:  name,
		Id:    id,
		Email: email,
	}
}

func printPersons(persons []*pb.Person) {
	for _, person := range persons {
		fmt.Printf("%+v\n", person)
	}
	fmt.Println("")
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:50001", opts...)
	if err != nil {
		log.Fatalf("grpc.Dial failed, error: %v\n", err)
	}
	defer conn.Close()
	client := pb.NewManageClient(conn)

	person := makePerson("Tom", 1, "tom@gmail.com")

	suc := addPerson(client, &person)
	if !suc {
		log.Fatalf("addPerson failed.\n")
	}

	person = makePerson("Lilly", 2, "lilly@gmail.com")
	person2 := makePerson("Jim", 3, "jim@gmail.com")

	persons := []*pb.Person{&person, &person2}
	suc = addPersons(client, persons)
	if !suc {
		log.Fatalf("addPersons failed.\n")
	}

	resPersons, err := getPersonsLimit(client, 5)
	if err != nil {
		log.Fatalf("getPersonsLimit failed, error: %v\n", err)
	}
	fmt.Println("getPersonsLimit output:")
	printPersons(resPersons)

	var personNames []string
	for _, person := range persons {
		personNames = append(personNames, person.GetName())
	}
	resPersons, err = getPersons(client, personNames)
	if err != nil {
		log.Fatalf("getPersons failed, error: %v\n", err)
	}
	fmt.Println("getPersons output:")
	printPersons(resPersons)
}
