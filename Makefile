build:
	go build -v ./cmd/tgclient/run.go && run.exe messages


.DEFAULT_GOAL:=build
 