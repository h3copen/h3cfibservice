package main

import (
	"fmt"
	"log"
	"time"
	"strings"
    "os"
	"github.com/h3copen/h3cfibservice/gen-go/ipprefix"
	"github.com/h3copen/h3cfibservice/gen-go/platform"
	"github.com/h3copen/comwaresdk/tproto/t_openr"
)

const defaultTimeout uint = 20
const localIfName string = "GE0_0_1"

type FibHandler struct {
	timeout    uint //timeout value, del route if counter > timeout
	counter    uint //timeout counter
	aliveSince int64
	offline    bool
	status     platform.ServiceStatus
	ticker     *time.Ticker
}

func NewFibHandler(timeout ...uint) *FibHandler {
	var s platform.ServiceStatus = platform.ServiceStatus_ALIVE
	keepAliveTicker := time.NewTicker(time.Second)
	handler := &FibHandler{timeout: defaultTimeout,
		aliveSince: time.Now().Unix(),
		status:     s,
		ticker:     keepAliveTicker}
	go handler.keepAlive()
	return handler
}

func (fh *FibHandler) keepAlive() {
	for _ = range fh.ticker.C {
		if fh.offline {
			continue
		}
		fh.counter++
		if fh.counter > fh.timeout {
			//del route if exist
			log.Println("keep alive timer expired, cleaning route")
			fh.offline = true
		}
	}
}

func ipv6Convert(preaddr []byte)(ipv6 string){
    for i := 0; i<16; i++ {
        ipv6 = ipv6 + fmt.Sprintf("%02x", preaddr[i])
        if (i%2 == 1)&&(i < 15){
            ipv6 = ipv6 + ":"
        }
    }
    return ipv6
}

func writeRoutesTxt(data string){
    file, err := os.Create("./routes.txt")
    if err != nil {
        fmt.Println(err)
    }
    file.WriteString(data)
    file.Close()
}
  // Parameters:
  //  - ClientId
  //  - Route
func (fh *FibHandler)  AddUnicastRoute( clientId int16, route *ipprefix.UnicastRoute) (err error) {
	fmt.Printf("AddUnicastRoute: \n")

    var mTRouteMsg t_openr.TRouteMsg
    var typeIp t_openr.TAddrType
    var preIp string = ""
    var nextIp string = ""

    routeLen := len(route.Dest.PrefixAddress.Addr)   
    // preIfName := route.Dest.PrefixAddress.IfName 
    prefixLength := route.Dest.PrefixLength    

    if routeLen == 4 {
        typeIp = t_openr.TAddrType_T_V4
        preIp = strings.Replace(strings.Trim(fmt.Sprint((*(*(route.Dest)).PrefixAddress).Addr), "[]"), " ", ".", -1)
    }else{
        typeIp = t_openr.TAddrType_T_V6
        preIp = ipv6Convert(route.Dest.PrefixAddress.Addr)
    }    
    
    data := preIp

    mTRouteMsg.Route = make([]*t_openr.TUnicstRoute, 1)
    mTRouteMsg.Route[0] = new(t_openr.TUnicstRoute)
    mTRouteMsg.Route[0].PrefixAddress = new(t_openr.TAddress)
    mTRouteMsg.Route[0].Path = make([]*t_openr.TRoutePath, len(route.Nexthops))

    for i := 0; i < len(route.Nexthops); i++ {
        if routeLen == 4 {
            nextIp = strings.Replace(strings.Trim(fmt.Sprint((*(route.Nexthops[i])).Addr), "[]"), " ", ".", -1)
        }else{
            nextIp = ipv6Convert(route.Nexthops[i].Addr)
        } 
           
        mTRouteMsg.Route[0].Path[i] = new(t_openr.TRoutePath)
        mTRouteMsg.Route[0].Path[i].NexthopAddress = new(t_openr.TAddress)

        mTRouteMsg.IndexOfRouteMsg = uint64 (i)
        mTRouteMsg.EnOperType = 2
        mTRouteMsg.VrfName = "0"
        mTRouteMsg.Route[0].PrefixAddress.Type = typeIp
        mTRouteMsg.Route[0].PrefixAddress.Address = preIp
        mTRouteMsg.Route[0].PrefixLen = uint32 (prefixLength)
        mTRouteMsg.Route[0].Preference = 0
        mTRouteMsg.Route[0].Path[i].LocalIfName = localIfName
        mTRouteMsg.Route[0].Path[i].NexthopAddress.Type = typeIp
        mTRouteMsg.Route[0].Path[i].NexthopAddress.Address = nextIp
        mTRouteMsg.Route[0].Path[i].Cost = 0

        data = data + "\n" + nextIp

    }

    if(isWrite){
        writeRoutesTxt(data)
    }
    if(isGrpc){
        err = SendRoute(&mTRouteMsg)
        return err
    }
    return nil
}

  // Parameters:
  //  - ClientId
  //  - Prefix
func (fh *FibHandler)  DeleteUnicastRoute( clientId int16, prefix *ipprefix.IpPrefix) (err error) {
	fmt.Printf("DeleteUnicastRoute: \n")

    routeLen := len(prefix.PrefixAddress.Addr)
    
    prefixLength := prefix.PrefixLength    
    
    var mTRouteMsg t_openr.TRouteMsg
    var typeIp t_openr.TAddrType
    var preIp string = ""

    if routeLen == 4{
        typeIp = t_openr.TAddrType_T_V4
        preIp = strings.Replace(strings.Trim(fmt.Sprint((*(prefix.PrefixAddress)).Addr), "[]"), " ", ".", -1)
    }else{
        typeIp = t_openr.TAddrType_T_V6
        preIp = ipv6Convert(prefix.PrefixAddress.Addr)
    }

    mTRouteMsg.Route = make([]*t_openr.TUnicstRoute, 1)
    mTRouteMsg.Route[0] = new(t_openr.TUnicstRoute)
    mTRouteMsg.Route[0].PrefixAddress = new(t_openr.TAddress)
    mTRouteMsg.Route[0].Path = make([]*t_openr.TRoutePath, 1)
    mTRouteMsg.Route[0].Path[0] = new(t_openr.TRoutePath)
    mTRouteMsg.Route[0].Path[0].NexthopAddress = new(t_openr.TAddress)

    mTRouteMsg.IndexOfRouteMsg = 1
    mTRouteMsg.EnOperType = 3
    mTRouteMsg.VrfName = "0"
    mTRouteMsg.Route[0].PrefixAddress.Type = typeIp
    mTRouteMsg.Route[0].PrefixAddress.Address = preIp
    mTRouteMsg.Route[0].PrefixLen = uint32 (prefixLength)
    mTRouteMsg.Route[0].Preference = 0
    mTRouteMsg.Route[0].Path[0].LocalIfName = localIfName
    // mTRouteMsg.Route[0].Path[0].NexthopAddress.Type = typeIp
    // mTRouteMsg.Route[0].Path[0].NexthopAddress.Address = ""
    mTRouteMsg.Route[0].Path[0].Cost = 0

    data := preIp
    if(isWrite){
        writeRoutesTxt(data)
    }

    if(isGrpc){
        err = SendRoute(&mTRouteMsg)
        return err
    }
    return nil

}

  // Parameters:
  //  - ClientId
  //  - Routes
func (fh *FibHandler)  AddUnicastRoutes( clientId int16, routes []*ipprefix.UnicastRoute) (err error) {
	fmt.Printf("AddUnicastRoutes: \n")

    numRoutes := len (routes)
    var mTRouteMsg t_openr.TRouteMsg
    var typeIp t_openr.TAddrType
    var data string = ""
    var preIp string = ""
    var nextIp string = ""

    mTRouteMsg.Route = make([]*t_openr.TUnicstRoute, len(routes))
    for i := 0; i < numRoutes; i++ {
        routeLen := len(routes[i].Dest.PrefixAddress.Addr)

        mTRouteMsg.Route[i] = new(t_openr.TUnicstRoute)
        mTRouteMsg.Route[i].PrefixAddress = new(t_openr.TAddress)
        mTRouteMsg.Route[i].Path = make([]*t_openr.TRoutePath, len(routes[i].Nexthops))    
        prefixLength := routes[i].Dest.PrefixLength   


        data = data + preIp +"\n"

        if routeLen == 4{
            typeIp = t_openr.TAddrType_T_V4
            preIp = strings.Replace(strings.Trim(fmt.Sprint((*(*(routes[i].Dest)).PrefixAddress).Addr), "[]"), " ", ".", -1)
        }else{
            typeIp = t_openr.TAddrType_T_V6
            preIp = ipv6Convert(routes[i].Dest.PrefixAddress.Addr)
        }    

        for j := 0; j < len(routes[i].Nexthops); j++ {

            if routeLen == 4 {
                nextIp = strings.Replace(strings.Trim(fmt.Sprint((*(routes[i].Nexthops[j])).Addr), "[]"), " ", ".", -1) 
            }else{
                nextIp = ipv6Convert(routes[i].Nexthops[j].Addr)
            } 
            
            mTRouteMsg.Route[i].Path[j] = new(t_openr.TRoutePath)
            mTRouteMsg.Route[i].Path[j].NexthopAddress = new(t_openr.TAddress)

            mTRouteMsg.IndexOfRouteMsg = uint64 (i)
            mTRouteMsg.EnOperType = 2
            mTRouteMsg.VrfName = "0"
            mTRouteMsg.Route[i].PrefixAddress.Type = typeIp
            mTRouteMsg.Route[i].PrefixAddress.Address = preIp
            mTRouteMsg.Route[i].PrefixLen = uint32 (prefixLength)
            mTRouteMsg.Route[i].Preference = 0
            mTRouteMsg.Route[i].Path[j].LocalIfName = localIfName
            mTRouteMsg.Route[i].Path[j].NexthopAddress.Type = typeIp
            mTRouteMsg.Route[i].Path[j].NexthopAddress.Address = nextIp
            mTRouteMsg.Route[i].Path[j].Cost = 0

            if j == (len(routes[i].Nexthops)-1){
                data = data + nextIp + "\n\n"
            }else{
                data = data + nextIp + "\n"
            }
        }  
    }

    if(isWrite){
        writeRoutesTxt(data)
    }

    if(isGrpc){
        err = SendRoute(&mTRouteMsg)
        return err
    }
    return nil
}

  // Parameters:
  //  - ClientId
  //  - Prefixes
func (fh *FibHandler)  DeleteUnicastRoutes( clientId int16, prefixes []*ipprefix.IpPrefix) (err error) {
	fmt.Printf("DeleteUnicastRoutes: \n")

    numRoutes := len(prefixes)
    var mTRouteMsg t_openr.TRouteMsg
    var data string = ""
    var preIp string = ""

    mTRouteMsg.Route = make([]*t_openr.TUnicstRoute, len(prefixes))
    for i := 0; i < numRoutes; i++ {
        routeLen := len(prefixes[i].PrefixAddress.Addr)
         
        prefixLength := prefixes[i].PrefixLength

        mTRouteMsg.Route[i] = new(t_openr.TUnicstRoute)
        mTRouteMsg.Route[i].PrefixAddress = new(t_openr.TAddress)
        mTRouteMsg.Route[i].Path = make([]*t_openr.TRoutePath, 1)
        mTRouteMsg.Route[i].Path[0] = new(t_openr.TRoutePath)
        mTRouteMsg.Route[i].Path[0].NexthopAddress = new(t_openr.TAddress)

        var typeIp t_openr.TAddrType
        if routeLen == 4{
            typeIp = t_openr.TAddrType_T_V4
            preIp = strings.Replace(strings.Trim(fmt.Sprint((*(prefixes[i].PrefixAddress)).Addr), "[]"), " ", ".", -1)
        }else{
            typeIp = t_openr.TAddrType_T_V6
            preIp = ipv6Convert(prefixes[i].PrefixAddress.Addr)
        }   

        mTRouteMsg.IndexOfRouteMsg = uint64 (i)
        mTRouteMsg.EnOperType = 3
        mTRouteMsg.VrfName = "0"
        mTRouteMsg.Route[i].PrefixAddress.Type = typeIp
        mTRouteMsg.Route[i].PrefixAddress.Address = preIp
        mTRouteMsg.Route[i].PrefixLen = uint32 (prefixLength)
        mTRouteMsg.Route[i].Preference = 0
        mTRouteMsg.Route[i].Path[0].LocalIfName = localIfName
        mTRouteMsg.Route[i].Path[0].Cost = 0

        data = data + "\n\n" + preIp
    }

    if(isWrite){
        writeRoutesTxt(data)
    }

    if(isGrpc){
        err = SendRoute(&mTRouteMsg)
        return err
    }
    return nil
}

// Parameters:
//  - ClientId
//  - Routes
func (fh *FibHandler) SyncFib(clientId int16, routes []*ipprefix.UnicastRoute) (err error) {
	fmt.Printf("SyncFib: \n")

    err = fh.AddUnicastRoutes(clientId, routes)

    return err
}

// DEPRECATED ... Use `aliveSince` API instead
// openr should periodically call this to let Fib know that it is alive
//
// Parameters:
//  - ClientId
func (fh *FibHandler) PeriodicKeepAlive(clientId int16) (r int64, err error) {
	fmt.Printf("PeriodicKeepAlive: %v, timeout: %v\n", clientId, fh.counter)
	fh.online()
	return 0, nil
}

// Returns the unix time that the service has been running since
func (fh *FibHandler) AliveSince() (r int64, err error) {
	fmt.Printf("AliveSince: %v\n", fh.aliveSince)
	fh.online()
	return fh.aliveSince, nil
}

// Get the status of this service
func (fh *FibHandler) GetStatus() (r platform.ServiceStatus, err error) {
	fmt.Printf("GetStatus: status: %v, aliveSince: %v\n", string(fh.status), fh.aliveSince)
	fh.online()
	return fh.status, nil
}

// Get number of routes
func (fh *FibHandler) GetCounters() (r map[string]int64, err error) {
	fmt.Printf("GetCounters was called\n")

	var counters = make(map[string]int64)
	counters["BGP"] = 10
	return counters, nil
}

// Parameters:
//  - ClientId
func (fh *FibHandler) GetRouteTableByClient(clientId int16) (r []*ipprefix.UnicastRoute, err error) {
	fmt.Printf("GetRouteTableByClient: %v\n", clientId)
	var rt []*ipprefix.UnicastRoute

	return rt, nil
}

func (fh *FibHandler) online() {
	if fh.offline {
		fh.counter = 0
		fh.offline = false
	}
}
