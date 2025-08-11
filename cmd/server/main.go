package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	hellopb "grpc_tutorial/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// hello_grpc.pb.goに書いてあるサービスのInterfaceを実装
type helloServer struct {
	// サービスの前方互換性を保つためのおまじない
	hellopb.UnimplementedGreetingServiceServer
}

// Helloサービスのメソッド定義
func (s *helloServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Body: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

// HelloServerStreamサービスの定義
func (s *helloServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Body: fmt.Sprintf("[%d] Hello %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

// コンストラクタ
func NewHelloServer() *helloServer {
	return &helloServer{}
}

func main() {
	const port = 8080
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// grpcサーバーの作成
	serv := grpc.NewServer()

	// サービスの登録
	hellopb.RegisterGreetingServiceServer(serv, NewHelloServer())

	reflection.Register(serv)

	// 非同期でクライアントの接続待ち
	go func() {
		log.Printf("listen at port %d", port)
		serv.Serve(listner)
	}()

	// メインスレッドではsigint待ちをするのみ
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("server stop...")
	serv.GracefulStop()
}
