# 生成代码

## user.proto

```bash
goctl rpc protoc api/proto/user.proto --go_out=./service/userService/pb --go-grpc_out=./service/userService/pb --zrpc_out=./service/userService --style goZero
```

## friend.proto

```bash
goctl rpc protoc api/proto/friend.proto --go_out=./service/friendService/pb --go-grpc_out=./service/friendService/pb --zrpc_out=./service/friendService  --style goZero
```

## message.proto

```bash
goctl rpc protoc api/proto/message.proto --go_out=./service/messageService/pb --go-grpc_out=./service/messageService/pb --zrpc_out=./service/messageService  --style goZero
```

## groupMessage.proto

```bash
goctl rpc protoc api/proto/groupMessage.proto  --go_out=./service/groupMessageService/pb  --go-grpc_out=./service/groupMessageService/pb  --zrpc_out=./service/groupMessageService  --style goZero
```
