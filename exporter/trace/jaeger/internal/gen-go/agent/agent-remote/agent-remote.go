// Autogenerated by Thrift Compiler (0.11.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"

	"go.opencensus.io/exporter/trace/jaeger/internal/gen-go/agent"
	"go.opencensus.io/exporter/trace/jaeger/internal/gen-go/jaeger"
	"go.opencensus.io/exporter/trace/jaeger/internal/gen-go/zipkincore"
)

var _ = jaeger.GoUnusedProtection__
var _ = zipkincore.GoUnusedProtection__

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  void emitZipkinBatch( spans)")
	fmt.Fprintln(os.Stderr, "  void emitBatch(Batch batch)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
	var parsedUrl *url.URL
	var trans thrift.TTransport
	_ = strconv.Atoi
	_ = math.Abs
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost", "Specify host and port")
	flag.IntVar(&port, "p", 9090, "Specify port")
	flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
	flag.StringVar(&urlString, "u", "", "Specify the url")
	flag.BoolVar(&framed, "framed", false, "Use framed transport")
	flag.BoolVar(&useHttp, "http", false, "Use http")
	flag.Parse()

	if len(urlString) > 0 {
		var err error
		parsedUrl, err = url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
	}

	cmd := flag.Arg(0)
	var err error
	if useHttp {
		trans, err = thrift.NewTHttpClient(parsedUrl.String())
	} else {
		portStr := fmt.Sprint(port)
		if strings.Contains(host, ":") {
			host, portStr, err = net.SplitHostPort(host)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error with host:", err)
				os.Exit(1)
			}
		}
		trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}
		if framed {
			trans = thrift.NewTFramedTransport(trans)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating transport", err)
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
		break
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
		break
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
		Usage()
		os.Exit(1)
	}
	iprot := protocolFactory.GetProtocol(trans)
	oprot := protocolFactory.GetProtocol(trans)
	client := agent.NewAgentClient(thrift.NewTStandardClient(iprot, oprot))
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "emitZipkinBatch":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "EmitZipkinBatch requires 1 args")
			flag.Usage()
		}
		arg5 := flag.Arg(1)
		mbTrans6 := thrift.NewTMemoryBufferLen(len(arg5))
		defer mbTrans6.Close()
		_, err7 := mbTrans6.WriteString(arg5)
		if err7 != nil {
			Usage()
			return
		}
		factory8 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt9 := factory8.GetProtocol(mbTrans6)
		containerStruct0 := agent.NewAgentEmitZipkinBatchArgs()
		err10 := containerStruct0.ReadField1(jsProt9)
		if err10 != nil {
			Usage()
			return
		}
		argvalue0 := containerStruct0.Spans
		value0 := argvalue0
		fmt.Print(client.EmitZipkinBatch(context.Background(), value0))
		fmt.Print("\n")
		break
	case "emitBatch":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "EmitBatch requires 1 args")
			flag.Usage()
		}
		arg11 := flag.Arg(1)
		mbTrans12 := thrift.NewTMemoryBufferLen(len(arg11))
		defer mbTrans12.Close()
		_, err13 := mbTrans12.WriteString(arg11)
		if err13 != nil {
			Usage()
			return
		}
		factory14 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt15 := factory14.GetProtocol(mbTrans12)
		argvalue0 := jaeger.NewBatch()
		err16 := argvalue0.Read(jsProt15)
		if err16 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.EmitBatch(context.Background(), value0))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}
