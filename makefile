.PHONY: all userService friendService messageService group-messageService gateway

all:  userService friendService messageService group-messageService gateway

userService:
	go build -o bin/userService.exe ./service/userService/user.go

friendService:
	go build -o bin/friendService.exe ./service/friendService/friend.go

messageService:
	go build -o bin/messageService.exe ./service/messageService/message.go

group-messageService:
	go build -o bin/groupMessageService.exe ./service/groupMessageService/groupmessage.go

gateway:
	go build -o bin/gateway.exe ./gateway/main.go

