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
	"google.golang.org/grpc/status"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
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
		grpc.WithUnaryInterceptor(myUnaryClientInterceptor1),
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
		fmt.Println("3: HelloClientStream")
		fmt.Println("4: HelloBiStreams")
		fmt.Println("5: exit")
		fmt.Print("please enter > ")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()
		case "2":
			HelloServerStream()
		case "3":
			HelloClientStream()
		case "4":
			HelloBiStreams()
		case "5":
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
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("code: %v\n", stat.Code())
			fmt.Printf("message: %v\n", stat.Message())
			fmt.Printf("details: %v\n", stat.Details())
		} else {
			fmt.Println("error: ", err)
		}
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

func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	sendCount := 5
	fmt.Printf("Please enter %d names\n", sendCount)
	for i := 0; i < sendCount; i++ {
		scanner.Scan()
		name := scanner.Text()

		req := &hellopb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			fmt.Println("error: ", err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("error: ", err)
		return
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloBiStreams() {
	stream, err := client.HelloBiStreams(context.Background())
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	sendNum := 5
	fmt.Printf("Please enter %d names\n", sendNum)

	var sendEnd, recvEnd bool
	sendCount := 0
	for !(sendEnd && recvEnd) {
		if !sendEnd {
			scanner.Scan()
			name := scanner.Text()

			sendCount++
			req := &hellopb.HelloRequest{
				Name: name,
			}
			if err := stream.Send(req); err != nil {
				fmt.Println("error: ", err)
				sendEnd = true
			}

			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println("error: ", err)
				}
			}
		}

		if !recvEnd {
			res, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				recvEnd = true
				continue
			}
			if err != nil {
				fmt.Println("error: ", err)
			}
			fmt.Println(res.GetMessage())
		}
	}
}
