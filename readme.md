
## Project setup
```bash
$ go mod init <project_name>
```

## RUN IN PRODUCTION
```bash
    docker build -t control-panel .
    docker run -p 8080:8080 control-panel
```

## RUN LOCALLY
```bash
    docker-compose up
```

### Format Project
```bash
    go fmt ./...
```

### Compile and run the project
```bash
# run
$ go run cmd/main.go
```
### RUN LOCALLY
```bash
    npm run dev
```

### RUN DOCKER FILE  WITH THE ENVIRONMENT VARIABLES
```bash
    docker run --env-file .env myapp
```
### RUN DOCKER COMPOSE
```bash
    docker-compose up
```

### RUN DOCKER COMPOSE <SPECIFIC SERVICE>
```bash
    docker-compose up <service-name>
```

### RUN DOCKER COMPOSE <SPECIFIC SERVICE> IN DETACHED MODE
```bash
    docker-compose up -d <service-name>
```

### STOP DOCKER COMPOSE
```bash
    docker-compose down
```

### STOP DOCKER COMPOSE <SPECIFIC SERVICE>
```bash
    docker-compose down <service-name> || docker-compose stop <service-name>
```

### REBUILD A CONTAINER
```bash
    docker-compose up --build 
```

### Remove containers and volumes
```bash
    docker-compose down -v
```

### DOCKER REACTING WITH MONGODB

```bash
    
    # docker-compose exec control-panel-app env | grep MONGO_URI
    docker-compose exec control-panel-app \
    mongosh "mongodb://root:example@mongodb:27017" --eval "db.runCommand({ping:1})"
    
    # verify environment variables
    docker-compose exec control-panel-app env | grep MONGO_URI
    
    # checks logs
    docker logs mongodb_container | tail -50
    
    # manually testing connection in docker and opening the mongo shell
    docker exec -it mongodb_container mongosh -u admin -p password --authenticationDatabase admin
    
    # list of dbs
    show dbs
    
    # select db
    use flowCx
    
    # list of collections
    show collections

    # display all documents in the users collection in a readable format
    db.users.find().pretty()
    
    # count collection document
    db.your_collection_name.countDocuments()
    
    # find a certain doc based on inputs
    db.users.find({}, { name: 1, email: 1, _id: 0 }).pretty()
    
    #
```


### Remove specific volume manually
```bash
    docker volume rm <volume_name>
```

### Shows logs in real time
```bash
    docker-compose logs -f
```

### Test ALl File
```bash
    go test ./...
```

### Test Files in a module
```bash
    cd module
    go test 
```

### Test Coverage Percentage
```bash
    go test -cover -coverprofile=coverage.out
```

### Test coverage in a specific package
```bash
    go test -cover 
```


```bash
    java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
```