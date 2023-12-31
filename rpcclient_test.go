package xrpc

import (
	"context"
	"fmt"
	pb "github.com/mix-go/xrpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestNewGrpcClient(t *testing.T) {
	conn, err := NewGrpcClient("127.0.0.1:50000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewAppMessagesClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), CallTimeout)
	resp, err := client.Send(ctx, &pb.SendRequest{
		Text: "123456789",
	})
	fmt.Println(resp, err)
}

func TestNewGrpcTLSClient(t *testing.T) {
	dir, _ := os.Getwd()
	tlsConf, err := LoadClientTLSConfig(dir+"/certificates/ca.pem", dir+"/certificates/client.pem", dir+"/certificates/client.key")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := NewGrpcClient("127.0.0.1:50000", grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewAppMessagesClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), CallTimeout)
	resp, err := client.Send(ctx, &pb.SendRequest{
		Text: "123456789",
	})
	fmt.Println(resp, err)
}

func TestNewGatewayClient(t *testing.T) {
	client := &http.Client{}
	resp, err := client.Post("http://127.0.0.1:50001/v1/send_message", "application/json", strings.NewReader(`{"order_number":"123456789"}`))
	if err != nil {
		log.Fatal(err)
	}
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b), err)
}

func TestNewGatewayTLSClient(t *testing.T) {
	dir, _ := os.Getwd()
	tlsConf, err := LoadClientTLSConfig(dir+"/certificates/ca.pem", dir+"/certificates/client.pem", dir+"/certificates/client.key")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConf,
		},
	}
	defer client.CloseIdleConnections()
	resp, err := client.Post("https://127.0.0.1:50001/v1/send_message", "application/json", strings.NewReader(`{"order_number":"123456789"}`))
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b), err)
}
