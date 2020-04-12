# RoboViewer

## How to run
Run the following from the project root.

docker-compose:
```
docker-compose up -d
```

robot gateway server:
```
go run robot_gateway/server/server.go
```

robot service server:
```
go run robot_service/server/server.go
```

RoboViewer client:
```
go run client/client.go
```
