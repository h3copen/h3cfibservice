package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"io"
	"log"
	"time"
	thrift "github.com/facebook/fbthrift-go"
	"github.com/h3copen/h3cfibservice/gen-go/platform"
	sdk "github.com/h3copen/comwaresdk/sdk"
	"github.com/h3copen/comwaresdk/tproto/t_openr"
)

var (
	addr       string
	protocol   string
	useTLS     bool
	isFramed   bool
	isBuffered bool
)


var(
	address      string = "192.168.18.102"
	port         uint = 50051
	username     string = "2"
	password     string = "123456"
	grpcSession  *sdk.GrpcSession
	isWrite      bool
	isGrpc       bool
)

func SendRoute(mTRouteMsg *t_openr.TRouteMsg) (err error){

	grpcSession, err = sdk.NewClient(address, port, username , password)
	if err != nil {
		log.Println("Failed to open session.")
	}
	defer grpcSession.Close()

	ctx_with_token, cancel := sdk.CtxWithToken(grpcSession.Token, time.Second)
	defer cancel()
	src := t_openr.NewTAgentOperClient(grpcSession.Conn)

	stream, err := src.SyncRoutes(ctx_with_token)

	stream.Send(mTRouteMsg)
    mTRouteMsgRsp,err := stream.Recv()
    if err == io.EOF {
        log.Printf("stream recv EOF: %v", err)
            return err
    }
    if err != nil {
        log.Printf("stream recv end: %v\n", err)
        return err
    }

    fmt.Printf("SyncRoutes ErrorStatus:%v\n", 
        mTRouteMsgRsp.ErrorStatus.Status)
    return nil
}

func init() {
	flag.StringVar(&addr, "addr", ":60100", "Address to listen to")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS secure transport")
	flag.StringVar(&protocol, "p", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	flag.BoolVar(&isFramed, "framed", true, "Use framed transport")
	flag.BoolVar(&isBuffered, "buffered", false, "Use buffered transport")

	flag.StringVar(&address, "ac", "192.168.18.102", "addr to comware")
	flag.UintVar(&port, "pc", 50051, "grpc port to comware")
	flag.StringVar(&username, "uc", "2", "username to comware")
	flag.StringVar(&password, "pwc", "123456", "password to comware")
	flag.BoolVar(&isWrite, "wr", false, "write routes to txt")
	flag.BoolVar(&isGrpc, "gc", false, "open grpc connect to comware")
}

func main() {
	flag.Parse()

	var err error
	var protocolFactory thrift.ProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		return
	}

	fmt.Printf("protocol: %v,", protocol)

	var transportFactory thrift.TransportFactory
	if isBuffered {
		fmt.Printf("buffered,")
		transportFactory = thrift.NewBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTransportFactory()
	}

	if isFramed {
		fmt.Printf("framed,")
		transportFactory = thrift.NewFramedTransportFactory(transportFactory)
	}

	err = runServer(transportFactory, protocolFactory, addr, useTLS)
	if err != nil {
		fmt.Println("Failed to run fib handler:", err)
	}

	return
}

func runServer(transportFactory thrift.TransportFactory, protocolFactory thrift.ProtocolFactory, addr string, secure bool) error {
	var transport thrift.ServerTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewServerSocket(addr)
	}

	if err != nil {
		return err
	}
	fmt.Printf("%T\n", transport)
	var handler = NewFibHandler()
	processor := platform.NewFibServiceProcessor(handler)
	server := thrift.NewSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the h3c fib handler on ", addr)
	return server.Serve()
}
