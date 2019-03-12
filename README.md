# H3cfibservice
OpenR FIB service by H3C, fibservice receives the route sent by OpenR and sends it to the H3C device.

## Build 
Get into directory `/fibservice/fibhandler`. The command go build compile source file and generate the `fibhandler`. Usage of `./fibhandler -h` to view the meaning of parameters.As follows： 
  -ac string
    	addr to comware (default "192.168.18.102")
  -addr string
    	Address to listen to (default ":60100")
  -buffered
    	Use buffered transport
  -framed
    	Use framed transport (default true)
  -gc
    	open grpc connect to comware
  -p string
    	Specify the protocol (binary, compact, json, simplejson) (default "binary")
  -pc uint
    	grpc port to comware (default 50051)
  -pwc string
    	password to comware (default "123456")
  -tls
    	Use TLS secure transport
  -uc string
    	username to comware (default "2")
  -wr
    	write routes to txt 
## Run
Fibservice run in container, generated by dockerfile. An openr container corresponds to a fibservice container. Openr and fibservice Shared the same network. For more information, please refer to [`h3copenr/build/test.h`](https://github.com/h3copen/h3copenr/blob/master/build/test.sh).
 
