// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"tera"
)

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  string Get(string table, string key)")
	fmt.Fprintln(os.Stderr, "  Status Put(string table, string key, string value)")
	fmt.Fprintln(os.Stderr, "   BatchGet(string table,  keys)")
	fmt.Fprintln(os.Stderr, "   BatchPut(string table,  kvs)")
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
	var parsedUrl url.URL
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
		parsedUrl, err := url.Parse(urlString)
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
	client := tera.NewProxyClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "Get":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "Get requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := argvalue0
		argvalue1 := flag.Arg(2)
		value1 := argvalue1
		fmt.Print(client.Get(value0, value1))
		fmt.Print("\n")
		break
	case "Put":
		if flag.NArg()-1 != 3 {
			fmt.Fprintln(os.Stderr, "Put requires 3 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := argvalue0
		argvalue1 := flag.Arg(2)
		value1 := argvalue1
		argvalue2 := flag.Arg(3)
		value2 := argvalue2
		fmt.Print(client.Put(value0, value1, value2))
		fmt.Print("\n")
		break
	case "BatchGet":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "BatchGet requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := argvalue0
		arg20 := flag.Arg(2)
		mbTrans21 := thrift.NewTMemoryBufferLen(len(arg20))
		defer mbTrans21.Close()
		_, err22 := mbTrans21.WriteString(arg20)
		if err22 != nil {
			Usage()
			return
		}
		factory23 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt24 := factory23.GetProtocol(mbTrans21)
		containerStruct1 := tera.NewProxyBatchGetArgs()
		err25 := containerStruct1.ReadField2(jsProt24)
		if err25 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Keys
		value1 := argvalue1
		fmt.Print(client.BatchGet(value0, value1))
		fmt.Print("\n")
		break
	case "BatchPut":
		if flag.NArg()-1 != 2 {
			fmt.Fprintln(os.Stderr, "BatchPut requires 2 args")
			flag.Usage()
		}
		argvalue0 := flag.Arg(1)
		value0 := argvalue0
		arg27 := flag.Arg(2)
		mbTrans28 := thrift.NewTMemoryBufferLen(len(arg27))
		defer mbTrans28.Close()
		_, err29 := mbTrans28.WriteString(arg27)
		if err29 != nil {
			Usage()
			return
		}
		factory30 := thrift.NewTSimpleJSONProtocolFactory()
		jsProt31 := factory30.GetProtocol(mbTrans28)
		containerStruct1 := tera.NewProxyBatchPutArgs()
		err32 := containerStruct1.ReadField2(jsProt31)
		if err32 != nil {
			Usage()
			return
		}
		argvalue1 := containerStruct1.Kvs
		value1 := argvalue1
		fmt.Print(client.BatchPut(value0, value1))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}