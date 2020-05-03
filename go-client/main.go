package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"hash/crc32"
	rnd "math/rand"
	"time"

	pb "../srvapi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var host = flag.String("host", "localhost:8080", "address of gRPC server")
var numSize = flag.Int("n", 10, "Number of tests")
var batchSize = flag.Int("bs", 10, "Batch size in KBytes (max 4 megs)")
var maxRecvMsgSize = flag.Int("mrsz", 10*1024, "MaxRecvMsgSize in kBytes")
var maxSendMsgSize = flag.Int("mssz", 10*1024, "MaxSendMsgSize in kBytes")

var useSSL = flag.Bool("secured", false, "Use SSL")

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

	var conn *grpc.ClientConn
	var err error
	if *useSSL {
		logrus.Infof("go RPC Client starts secured (SSL) for %s", *host)
		config := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err = grpc.Dial(*host, grpc.WithTransportCredentials(credentials.NewTLS(config)),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(*maxRecvMsgSize*1024),
				grpc.MaxCallSendMsgSize(*maxSendMsgSize*1024)))

		if err != nil {
			logrus.Errorf("could not connect to %s: %v", *host, err)
			return
		}
	} else {
		logrus.Infof("go RPC Client starts insecured (no SSL) for %s", *host)
		conn, err = grpc.Dial(*host, grpc.WithInsecure(),
			// this option is not working ?
			grpc.WithMaxMsgSize(6144*1024),
			// these options are not working ?
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(*maxRecvMsgSize*1024),
				grpc.MaxCallSendMsgSize(*maxSendMsgSize*1024)))

		if err != nil {
			logrus.Errorf("could not connect to %s: %v", *host, err)
			return
		}
	}
	defer conn.Close()

	client := pb.NewMyRPC1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50)*time.Second)
	// or use:
	//clientDeadline := time.Now().Add(time.Duration(3) * time.Second)
	//ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	rnd.NewSource(time.Now().UnixNano())
	crc32Table := crc32.MakeTable(crc32.IEEE)
	for i := 0; i < *numSize; i++ {
		t1 := time.Now().UnixNano()
		bench := &pb.Benchmark{}
		bench.Name = fmt.Sprintf("Test go %d", i)
		bench.Buffer = make([]byte, *batchSize*1024)
		rand.Read(bench.Buffer)
		crc := crc32.Checksum(bench.Buffer, crc32Table)
		res, err := client.Test1(ctx, bench)
		if err != nil {
			logrus.Errorf("RPC error: %v", err)
			break
		}
		t2 := time.Now().UnixNano()
		elapsed := float32(t2-t1) / 1e9
		if crc == res.Crc32 {
			logrus.Infof("Response: %v in %.3f sec, valid: %v.", res.Text, elapsed, crc == res.Crc32)
		} else {
			logrus.Errorf("Response: %v in %.3f sec, valid: %v.", res.Text, elapsed, crc == res.Crc32)
		}
	}
	logrus.Info("go RPC Client stops")
}
