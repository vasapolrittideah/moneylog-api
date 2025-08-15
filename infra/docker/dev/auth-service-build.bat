set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o build/auth_service ./services/auth_service/cmd/main.go
