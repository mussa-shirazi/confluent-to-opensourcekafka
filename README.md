# Confluent-To-Opensourcekafka

This repo is for Instruction to migrate from confluent Kafka to Apache Kafka

## Clone Migration Repo 

Clone the migration repo, we have provided a test setup to get you started.

Run the following command to clone the repo
```
git clone https://github.com/mussa-shirazi/confluent-to-opensourcekafka.git

```


## Setup test Conflunt Cluster 

Before performing the migration on your production cluster. We recommend you setup a test confluent environment that can be used for migration testing. The test environment should mimic the production environment as much as possible. Ingest realistic data into the cluster and try migrating data to Instaclustr Managed cluster. We have provided a docker compose which deploy one node Kafka cluster with Zookeeper. As an exercise we are going to Create two topics following the steps in next page to run Kafka cluster.

### Create Docker Network 

Create Docker network exists

```
docker network create migration_workshop || true

```
```
docker-compose -f confluent/docker-compose.yml up -d

```


Check both Zookeerper and Kafka Containers are up . If not run the above command again and this sould bring up both Kafka and Zookeeper containers

```

docker ps

CONTAINER ID   IMAGE                       COMMAND                  CREATED          STATUS          PORTS                                                 NAMES
e438dff09a3a   confluentinc/cp-kafka       "/etc/confluent/dock…"   36 seconds ago   Up 36 seconds   9092/tcp, 0.0.0.0:9094->9094/tcp, :::9094->9094/tcp   confluent_kafka_1
b6cf81cddf41   confluentinc/cp-zookeeper   "/etc/confluent/dock…"   37 seconds ago   Up 36 seconds   2181/tcp, 2888/tcp, 3888/tcp                          confluent_zookeeper_1

```

## Produce data to test Confluent Cluster

1. For the test environment we have built test producer which can produce data to both test and Instaclustr cluster
1. Producer is written in Golang you always further change suits your need. 
1. For this exercise we will produce some message to demonstrate the procedure to migrating data from confluent to Instaclustr Apache Kafka
1. Use the producer located in “producer/producer_app” folder and run the following command to produce messages to the confluent test cluster . For us we are running this exercise on Macos , also there are executables for other operating systems.
1. You can skip this step if you already have some producer app to test the procedure 


```
./producer/producer_app/producer_macOS -connecttoinstaclustr=false -brokers="localhost:9094" -topics="app_topic"

```

Open a New Terminal window


## Consumer data from test Confluent Cluster

1. For the test environment we have also built test consumer which can use to test consume data from to both test and Instaclustr cluster
1. Consumer is written in Golang and you always change to suit your need. Code can be found in can found on the /consumer/. 
1. Use the consumer located in “consumer/consumer_app” folder and run the following command to consumer messages from the confluent test cluster . For us we are running this exercise on Macos , also there are executables for other operating systems.
1. You can skip this step if you already have some consumer app to test the consumer 

```
./consumer/consumer_app/consumer_macOS -connecttoinstaclustr=false  -brokers="localhost:9094" -topics="app_topic" -group="consumer-app-group"

```

## Setup Instaclustr Cluster

Follow the steps to create the Kafka cluster on the "https://console2.instaclustr.com/"

Note the Public IP and username and password once the cluster is created. 


## Deploy Kafka Management UI 

1. Update the “docker-compose” file with following section
1. Update “bootstrap.server” with the Instaclustr Public IP bootstrap addresses
1. Update Instaclustr “username” and ”password”

```

InstaClustr-Cluster:
 properties:
   bootstrap.servers: "x.x.x.x:9092"
   security.protocol: SASL_PLAINTEXT
   sasl.mechanism: SCRAM-SHA-256
   sasl.jaas.config: org.apache.kafka.common.security.scram.ScramLoginModule required username="username" password="password";

```

### Bring up the AKHQ container 

```
docker-compose -f akhq/docker-compose.yml up -d

```

### Access AKHQ 


Open browser and access using AKHQ https://localhost:9080


## Deploy MirrorMaker2

1. In this step we will deploy mirrormaker2 which will copy topics and data from confluent cluster to New Instaclustr Kafka cluster.
1. Cd into mirrormaker folder and edit mm2.properties file and update the following section from on the migration repo 

### Update Bootstrap address 

```
target.bootstrap.servers = x.x.x.x:9092

```
### Update Password for the Instaclustr on mirrormaker.properties

```
target.security.protocol = SASL_PLAINTEXT
target.sasl.mechanism = SCRAM-SHA-256
target.sasl.jaas.config = org.apache.kafka.common.security.scram.ScramLoginModule required \
  username="Instaclustrusername" \
  password="Insaclustrpassword";

```

### Bring the mirrormker2 up 

```
docker-compose -f mirrormaker/docker-compose.yml up -d

```

## Move Consumers to New Instaclustr Kafka cluster 

1. In this step we will move consumers from the confluent cluster to new Instaclustr Kafka cluster
1. As a part of the mirror maker conigs we can replicate topic data, offsets and Consumer group information 
1. That mean when we move our consumers to the new Kafka cluster, the consumer group only consume the new messages in Instaclustr, which were not consumed in confluent cluster. 
1. Producers will continue to produce message on confluent cluster which  will be replicated to Instaclustr Kafka cluster 
1. Close consumer connecting to the confluent cluster CTRL+C . And run the following command with information of new Instaclustr details 

use the username and password of Instaclustr and bootstrap address of the Instaclustr

```
./consumer/consumer_app/consumer_macOS -connecttoinstaclustr=true -saslusername="ickafka" -saslpassword="9f239ea6831d8cb701212f8fe830fbbc6a3799e538e431ab4ac94bae88d040e0" -brokers="54.191.182.150:9092" -topics="app_topic" -group="consumer-app-group"

```

### Produce messages to Confluent Cluster

Produce messages to Confluent Cluster and Expect to see the new messages consume by the consumer connecting to Instaclustr 

```
./producer_macOS  -connecttoinstaclustr=false -brokers="localhost:9094" -topics="app_topic”

```

This test proves that we have successfully migrated our consumers. You could perform the similar step without losing the single message   


### Move Producers to Instaclustr Kafka

1. Now that we have successfully move our consumers. Next step is to move our producers
1. Close the existing producer CTRL+C
1. Connect your Producer to Instaclustr Kafka cluster using the following command 

```

./producer/producer_app/producer_macOS -connecttoinstaclustr=true -saslusername="ickafka" -saslpassword="9f239ea6831d8cb701212f8fe830fbbc6a3799e538e431ab4ac94bae88d040e0" -brokers="54.191.182.150:9092" -topics="app_topic" 

```

### Shut the confluent Kafka cluster

Once you have succesully migrated the Producers and consumer you can shutdown Confluent Cluster ,x


```
docker-compose -f confluent/docker-compose.yml down

```

### Shutdown Mirrormaker2 

Also we wont need mirror maker which we can also shutdown. 


```
docker-compose -f confluent/docker-compose.yml down

```
