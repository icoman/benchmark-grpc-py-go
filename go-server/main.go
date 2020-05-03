package main

import (
	"context"
	"flag"
	"fmt"
	"hash/crc32"
	"net"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "../srvapi"
)

var host = flag.String("host", "localhost:8080", "address of gRPC server")
var maxRecvMsgSize = flag.Int("mrsz", 10*1024, "MaxRecvMsgSize in kBytes")
var maxSendMsgSize = flag.Int("mssz", 1*1024, "MaxSendMsgSize in kBytes")
var useSSL = flag.Bool("secured", false, "Use SSL")

var crt = "server.crt"
var key = "server.key"

func init() {
	flag.Parse()
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "02-Jan-2006 15:04:05",
	})

}

func main() {

	logrus.Infof("go RPC Server starts on %s", *host)
	logrus.Infof("\tMaxRecvMsgSize = %.2f Mb", float32(*maxRecvMsgSize)/1024)
	logrus.Infof("\tMaxSendMsgSize = %.2f Mb", float32(*maxSendMsgSize)/1024)

	defer func() {
		logrus.Warn("defer EXIT")
	}()

	lis, err := net.Listen("tcp", *host)
	if err != nil {
		logrus.Errorf("could not listen to %s: %v", *host, err)
		return
	}

	var s *grpc.Server
	// https://grpc.io/docs/guides/auth/
	if *useSSL {
		// https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html
		// Create the TLS credentials
		creds, err := credentials.NewServerTLSFromFile(crt, key)
		if err != nil {
			logrus.Errorf("could not load TLS keys: %s", err)
			return
		}
		s = grpc.NewServer(
			grpc.Creds(creds),
			grpc.MaxRecvMsgSize(*maxRecvMsgSize*1024),
			grpc.MaxSendMsgSize(*maxSendMsgSize*1024),
		)
		logrus.Info("\tStart secured (SSL)")
	} else {
		s = grpc.NewServer()
		logrus.Info("\tStart insecured (no SSL)")
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		logrus.Error("Ctrl-C detected, exiting...")
		s.Stop()
	}()

	pb.RegisterMyRPC1Server(s, server{})
	err = s.Serve(lis)
	if err != nil {
		logrus.Errorf("could not server: %v", err)
		return
	}

	logrus.Info("go RPC Server stops")
}

type server struct{}

var crc32Table = crc32.MakeTable(crc32.IEEE)

func (server) Test1(ctx context.Context, bench *pb.Benchmark) (*pb.Response, error) {
	logrus.Infof("Test: %s (%.2f KBytes)", bench.Name, float32(len(bench.Buffer))/1024)
	crc := crc32.Checksum(bench.Buffer, crc32Table)
	msg := fmt.Sprintf("Test '%s' ok.", bench.Name)
	return &pb.Response{Text: msg, Crc32: crc}, nil
}
