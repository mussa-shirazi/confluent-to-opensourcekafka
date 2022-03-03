# Consumergroup example

This example shows you how to use the Sarama consumer group consumer. The example simply starts consuming the given Kafka topics and logs the consumed messages.

```bash
$ go run main.go -brokers="localhost:9094" -topics="messages" -group="consumer-group"
```

Working local


./consumer_app -connecttoinstaclustr=false  -brokers="localhost:9094" -topics="messages" -group="consumer-group"


to instacluster

./consumer_app -connecttoinstaclustr=true -saslusername="ickafka" -saslpassword="ESF5XlG&m1HushPv" -brokers="35.155.222.66:9092" -topics="messages" -group="consumer-group"
