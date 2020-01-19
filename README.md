# Mongo Change Stream POC

This repo shows how to watch changes in a collection with the ChangeStream API.

It consists of a mongo database with two nodes in replicaset, a watcher container that prints all documents inserted into the "logs" collection, and a generator container that generate random messages and make bulk inserts into the same "logs" collection. The generator container can be started with multiple replicas.


## Requirements

- Docker and Docker compose installed

## Usage

### Setup Mongo replica set

1. First you need to start mongo in replicaset. Start the nodes: 

`docker-compose up -d mongo0 mongo1`

2. Connect to one nodo and set the configuration:

`docker exec -it mongo0 mongo`

3. Once inside Mongo shell console, set the config and initiate the replicaset: 

```
config={"_id":"rs0","members":[{"_id":0,"host":"mongo0:27017"},{"_id":1,"host":"mongo1:27017"}]}
rs.initiate()
```

### Start the watcher

Start the ChangeStream watcher container:

`docker-compose up watcher`

### Start the generators

You can start multiple generators with the scale command:

`docker scale generator=10`

Set scale to 0 if you want to stop inserting messages into the logs collection


## References

* [How to simply set up Mongo's replica set locally with Docker](https://37yonub.ru/articles/mongo-replica-set-docker-localhost)
* [go - Watch for MongoDB Change Streams - Stack Overflow](https://stackoverflow.com/questions/49151104/watch-for-mongodb-change-streams)
* [Using the official MongoDB Go driver](https://vkt.sh/go-mongodb-driver-cookbook/)


