# RoboViewer

#run from root directory

#docker-compose
docker-compose up -d

#robot gateway server
go run robot_gateway/server/server.go

#robot gateway client
go run robot_gateway/server/client.go

