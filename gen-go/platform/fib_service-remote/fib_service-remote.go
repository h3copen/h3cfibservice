// Autogenerated by Thrift Compiler (facebook)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
// @generated

package main

import (
        "flag"
        "fmt"
        "math"
        "net"
        "net/url"
        "os"
        "strconv"
        "strings"
        thrift "github.com/facebook/fbthrift-go"
        "github.com/h3copen/h3cfibservice/gen-go/platform"
)

func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  void addUnicastRoute(i16 clientId, UnicastRoute route)")
  fmt.Fprintln(os.Stderr, "  void deleteUnicastRoute(i16 clientId, IpPrefix prefix)")
  fmt.Fprintln(os.Stderr, "  void addUnicastRoutes(i16 clientId,  routes)")
  fmt.Fprintln(os.Stderr, "  void deleteUnicastRoutes(i16 clientId,  prefixes)")
  fmt.Fprintln(os.Stderr, "  void syncFib(i16 clientId,  routes)")
  fmt.Fprintln(os.Stderr, "  i64 periodicKeepAlive(i16 clientId)")
  fmt.Fprintln(os.Stderr, "  i64 aliveSince()")
  fmt.Fprintln(os.Stderr, "  ServiceStatus getStatus()")
  fmt.Fprintln(os.Stderr, "   getCounters()")
  fmt.Fprintln(os.Stderr, "   getRouteTableByClient(i16 clientId)")
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
  var trans thrift.Transport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host")
  flag.IntVar(&port, "p", 60100, "Specify port")
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
    trans, err = thrift.NewHTTPPostClient(parsedUrl.String())
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewSocket(thrift.SocketAddr(net.JoinHostPort(host, portStr)))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.ProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  client := platform.NewFibServiceClientFactory(trans, protocolFactory)
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "addUnicastRoute":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "AddUnicastRoute requires 2 args")
      flag.Usage()
    }
    tmp0, err108 := (strconv.Atoi(flag.Arg(1)))
    if err108 != nil {
      Usage()
      return
    }
    argvalue0 := byte(tmp0)
    value0 := argvalue0
    arg109 := flag.Arg(2)
    mbTrans110 := thrift.NewMemoryBufferLen(len(arg109))
    defer mbTrans110.Close()
    _, err111 := mbTrans110.WriteString(arg109)
    if err111 != nil {
      Usage()
      return
    }
    factory112 := thrift.NewSimpleJSONProtocolFactory()
    jsProt113 := factory112.GetProtocol(mbTrans110)
    argvalue1 := platform.NewUnicastRoute()
    err114 := argvalue1.Read(jsProt113)
    if err114 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.AddUnicastRoute(value0, value1))
    fmt.Print("\n")
    break
  case "deleteUnicastRoute":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "DeleteUnicastRoute requires 2 args")
      flag.Usage()
    }
    tmp0, err115 := (strconv.Atoi(flag.Arg(1)))
    if err115 != nil {
      Usage()
      return
    }
    argvalue0 := byte(tmp0)
    value0 := argvalue0
    arg116 := flag.Arg(2)
    mbTrans117 := thrift.NewMemoryBufferLen(len(arg116))
    defer mbTrans117.Close()
    _, err118 := mbTrans117.WriteString(arg116)
    if err118 != nil {
      Usage()
      return
    }
    factory119 := thrift.NewSimpleJSONProtocolFactory()
    jsProt120 := factory119.GetProtocol(mbTrans117)
    argvalue1 := platform.NewIpPrefix()
    err121 := argvalue1.Read(jsProt120)
    if err121 != nil {
      Usage()
      return
    }
    value1 := argvalue1
    fmt.Print(client.DeleteUnicastRoute(value0, value1))
    fmt.Print("\n")
    break
  case "addUnicastRoutes":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "AddUnicastRoutes requires 2 args")
      flag.Usage()
    }
    tmp0, err122 := (strconv.Atoi(flag.Arg(1)))
    if err122 != nil {
      Usage()
      return
    }
    argvalue0 := byte(tmp0)
    value0 := argvalue0
    arg123 := flag.Arg(2)
    mbTrans124 := thrift.NewMemoryBufferLen(len(arg123))
    defer mbTrans124.Close()
    _, err125 := mbTrans124.WriteString(arg123)
    if err125 != nil { 
      Usage()
      return
    }
    factory126 := thrift.NewSimpleJSONProtocolFactory()
    jsProt127 := factory126.GetProtocol(mbTrans124)
    containerStruct1 := platform.NewFibServiceAddUnicastRoutesArgs()
    err128 := containerStruct1.ReadField2(jsProt127)
    if err128 != nil {
      Usage()
      return
    }
    argvalue1 := containerStruct1.Routes
    value1 := argvalue1
    fmt.Print(client.AddUnicastRoutes(value0, value1))
    fmt.Print("\n")
    break
  case "deleteUnicastRoutes":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "DeleteUnicastRoutes requires 2 args")
      flag.Usage()
    }
    tmp0, err129 := (strconv.Atoi(flag.Arg(1)))
    if err129 != nil {
      Usage()
      return
    }
    argvalue0 := byte(tmp0)
    value0 := argvalue0
    arg130 := flag.Arg(2)
    mbTrans131 := thrift.NewMemoryBufferLen(len(arg130))
    defer mbTrans131.Close()
    _, err132 := mbTrans131.WriteString(arg130)
    if err132 != nil { 
      Usage()
      return
    }
    factory133 := thrift.NewSimpleJSONProtocolFactory()
    jsProt134 := factory133.GetProtocol(mbTrans131)
    containerStruct1 := platform.NewFibServiceDeleteUnicastRoutesArgs()
    err135 := containerStruct1.ReadField2(jsProt134)
    if err135 != nil {
      Usage()
      return
    }
    argvalue1 := containerStruct1.Prefixes
    value1 := argvalue1
    fmt.Print(client.DeleteUnicastRoutes(value0, value1))
    fmt.Print("\n")
    break
  case "syncFib":
    if flag.NArg() - 1 != 2 {
      fmt.Fprintln(os.Stderr, "SyncFib requires 2 args")
      flag.Usage()
    }
    tmp0, err136 := (strconv.Atoi(flag.Arg(1)))
    if err136 != nil {
      Usage()
      return
    }
    argvalue0 := byte(tmp0)
    value0 := argvalue0
    arg137 := flag.Arg(2)
    mbTrans138 := thrift.NewMemoryBufferLen(len(arg137))
    defer mbTrans138.Close()
    _, err139 := mbTrans138.WriteString(arg137)
    if err139 != nil { 
      Usage()
      return
    }
    factory140 := thrift.NewSimpleJSONProtocolFactory()
    jsProt141 := factory140.GetProtocol(mbTrans138)
    containerStruct1 := platform.NewFibServiceSyncFibArgs()
    err142 := containerStruct1.ReadField2(jsProt141)
    if err142 != nil {
      Usage()
      return
    }
    argvalue1 := containerStruct1.Routes
    value1 := argvalue1
    fmt.Print(client.SyncFib(value0, value1))
    fmt.Print("\n")
    break
  case "periodicKeepAlive":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "PeriodicKeepAlive requires 1 args")
      flag.Usage()
    }
    tmp0, err143 := (strconv.Atoi(flag.Arg(1)))
    if err143 != nil {
      Usage()
      return
    }
    argvalue0 := byte(tmp0)
    value0 := argvalue0
    fmt.Print(client.PeriodicKeepAlive(value0))
    fmt.Print("\n")
    break
  case "aliveSince":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "AliveSince requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.AliveSince())
    fmt.Print("\n")
    break
  case "getStatus":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetStatus requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetStatus())
    fmt.Print("\n")
    break
  case "getCounters":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetCounters requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetCounters())
    fmt.Print("\n")
    break
  case "getRouteTableByClient":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "GetRouteTableByClient requires 1 args")
      flag.Usage()
    }
    tmp0, err144 := (strconv.Atoi(flag.Arg(1)))
    if err144 != nil {
      Usage()
      return
    }
    argvalue0 := byte(tmp0)
    value0 := argvalue0
    fmt.Print(client.GetRouteTableByClient(value0))
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
