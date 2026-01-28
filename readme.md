# Skeleton Project

## Docker Setup Kafka

### Docker compose file
name it "docker-compose.yaml"
```yaml
version: '2'
services:
  zookeeper:
    image: wurstmeister/zookeeper:3.4.6
    container_name: zookeeper
    ports:
      - "2181:2181"
    environment:
      - ALLOW_ANONYMOUS_LOGIN=yes
  kafka:
    image: wurstmeister/kafka:2.13-2.7.2
    container_name: "kafka-01"
    hostname: kafka-01
    ports:
      - "29092:29092"
    environment:
      - KAFKA_BROKER_ID=1
      - KAFKA_LISTENERS=INSIDE://:9092,OUTSIDE://:29092
      - KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT
      - KAFKA_ADVERTISED_LISTENERS=INSIDE://kafka-01:9092,OUTSIDE://localhost:29092
      - KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
      - KAFKA_INTER_BROKER_LISTENER_NAME=INSIDE
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_CREATE_TOPICS=demand.trips:1:1:compact,demand.jobs:1:1:compact

  kowl:
    image: docker.redpanda.com/vectorized/console:v2.1.1
    container_name: "kowl"
    hostname: kowl
    ports:
      - "8080:8080"
    environment:
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - zookeeper
      - kafka

```


### Install librd kafka

https://github.com/confluentinc/librdkafka


### docker compose run command

``` docker-compose -f docker-compose.yaml up```

 ### run the skeleton/service

 initially build the skeleton/service binary

 ```go build -tags=dynamic -o skeleton main.go```
 
### Use Following run config on VSCode 
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Golang Skeleton",
            "asRoot": true,
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/main.go",
            "envFile": "${workspaceFolder}/configs.env",
            "buildFlags": "-tags=dynamic"
        }
    ]
}
```

find the vscode config named "Golang Skeleton" and run your skeleton/service

if you encounter panics related to topics, find ip of your system using `ifconfig` and connect to red panda by going to the url
`<your ip>:8080` and do the configuration as follows
![img](readme-data/image.png) then restart the service