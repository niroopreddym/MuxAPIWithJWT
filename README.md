# muxapiwithjwt

API endpoints backed up by JWT Tokens as auth and orchestrated using docker-compose

## How to run

Run the docker compose file on your machine it will start a API server port 9293

Before running the compose file please check  DbConnectionString = "mongodb://192.168.99.100:27017"
in this file [a link](https://github.com/niroopreddym/muxapiwithjwt/blob/develop/internal/services/employeeservice.go)
and modify it according if you are running native docker app then it would be DbConnectionString = "mongodb://192.168.99.100:27017" else if you are using docker toolbox like me then you dont need to change the connection string
