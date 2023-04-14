# Plugin Manager over RPC

This example builds plugin manager which compose 3 clients and a server communicate through netrpc:
- Init Client
- Export functions for client to call include: Send/get normal data, Publish/Subcribe message though NATs broker
- Plugins manager (not yet)
To build this example:

```sh
# This command build the server main CLI
$ cd server
$ go build -o . .

# This builds the client plugins written in Go
$ cd clients/client1
$ go build . .

$ cd clients/client2
$ go build . .

$ cd clients/client3
$ go build . .
```


