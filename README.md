### Clone the repo

```bash
git clone github.com/mc-solo/friendy-backend.git

cd friendy
```

### start the psql container

```docker
docker-composer up -d
docker ps
```

### download deps

```go
go mod download
go run main.go
```
