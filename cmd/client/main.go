package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	hellopb "myrpc/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	scanner *bufio.Scanner
	client  hellopb.GreetingServiceClient
)

func main() {
	fmt.Println("start gRPC Client.")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	// grpc.Dial は非推奨になり、 grpc.NewClient を使用
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//　非同期でコネクション確立を行った方が効率がいいので、 grpc.WithBlock() は非推奨
		// grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer conn.Close()
	// ここでコネクション確立を待つ
	conn.Connect()

	client = hellopb.NewGreetingServiceClient(conn)
	for {
		fmt.Println("1: send Request")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: exit")
		fmt.Print("please enter > ")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()
		case "2":
			HelloServerStream()
		case "3":
			fmt.Println("bye")
			goto M
		}
	}
M:
}

func Hello() {
	fmt.Println("please enter your name >")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream() {
	fmt.Println("please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received.")
			break
		}
		if err != nil {
			fmt.Println("error: ", err)
			break
		}
		fmt.Println(res.GetMessage())
	}
}
