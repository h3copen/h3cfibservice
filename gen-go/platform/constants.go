// Autogenerated by Thrift Compiler (facebook)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
// @generated

package platform

import (
	"bytes"
	"context"
	"sync"
	"fmt"
	thrift "github.com/facebook/fbthrift-go"
	ipprefix0 "github.com/h3copen/h3cfibservice/gen-go/ipprefix"

)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = sync.Mutex{}
var _ = bytes.Equal
var _ = context.Background

var _ = ipprefix0.GoUnusedProtection__
var ClientIdtoProtocolId map[int16]int16
var ClientIdtoPriority map[int16]int16

func init() {
ClientIdtoProtocolId = map[int16]int16{
  786: 99,
  0: 253,
  64: 64,
}

ClientIdtoPriority = map[int16]int16{
  786: 10,
  0: 20,
  64: 11,
}

}

