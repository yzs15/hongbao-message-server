
https://docs.qq.com/doc/DR2JWSHROdFVIa1pi

```text
Usage of msd:
  -kend value
        kubernetes' http endpoint
  -net
        in internet environment
  -spb
        in superbahn environment
  -thing
        run as one Thing
  -wang
        run as Wang
  -wend string
        Wang's czmq endpoint (default "tcp://127.0.0.1:5553")
  -tend value
        thing's czmq endpoint
  -ws string
        web socket service address (default "0.0.0.0:5554")
  -zmq string
        zmq service address (default "tcp://0.0.0.0:5553")
  -log string
        log service address (default "0.0.0.0:5552")
```

## Run
Run as wang in internet
```bash
go run cmd/msd/msd.go -wang -net
```

Run as one thing in superbahn
```bash
go run cmd/msd/msd.go -thing -spb
```

# Log Server

可通过指定多个 -f 参数，来指定多个 log 文件

```text
Usage of ./bin/logserverd:
  -addr string
        log service address (default "0.0.0.0:5552")
  -f value
        the path to log file
```

# 测试脚本
在script里面