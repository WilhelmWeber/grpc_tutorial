package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	hellopb "grpc_tutorial/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Hello(client hellopb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("please enter your name")

	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	// helloサービスの呼び出し
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(res.GetBody())
	}
}

func HelloServerStream(client hellopb.GreetingServiceClient, scanner *bufio.Scanner) {
	fmt.Println("please enter your name")

	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	// HelloServerStreamの呼び出し
	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// streamからeofになるまで取り出す
	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(res.GetBody())
	}
}

func main() {
	fmt.Println("client program start...")

	scanner := bufio.NewScanner(os.Stdin)
	const target = "dns:///localhost:8080"

	// grpcサーバーとの接続確立
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// grpcクライアントインスタンスの生成
	client := hellopb.NewGreetingServiceClient(conn)

OUTER_:
	for {
		fmt.Println("1: Hello")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello(client, scanner)
		case "2":
			HelloServerStream(client, scanner)
		case "3":
			fmt.Println("Good Bye...")
			break OUTER_
		}
	}
}
