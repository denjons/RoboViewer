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

robot progress processor server:
```
go run robot_progress_processor/server/server.go
```

robot gateway client:
```
go run robot_gateway/client/client.go
```

robot service server:
```
go run robot_service/server/server.go
```
