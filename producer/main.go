package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

// Sarama configuration options
var (
	brokers                      = ""
	topics                       = ""
	clientid                     = ""
	conect_to_Instaclutr_cluster = false
	saslusername                 = ""
	saslpassword                 = ""
)

func init() {
	flag.StringVar(&brokers, "brokers", "", "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&topics, "topics", "", "Kafka topics to be consumed, as a comma separated list")
	flag.StringVar(&clientid, "clientid", "", "ClientId for the Producer")
	flag.BoolVar(&conect_to_Instaclutr_cluster, "connecttoinstaclustr", false, "Connection to Instaclustr")
	flag.StringVar(&saslusername, "saslusername", "", "UserName to Instaclustr")
	flag.StringVar(&saslpassword, "saslpassword", "", "Password to Instaclustr")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topics) == 0 {
		panic("no topics given to be consumed, please set the -topics flag")
	}

	if conect_to_Instaclutr_cluster {
		if saslpassword == "" || saslusername == "" {
			panic("Provide the username and Password for the Instaclustr Kafka Cluster")

		}

	}

}

func main() {
	// create producer
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	// read command line input
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter msg: ")
		msg, _ := reader.ReadString('\n')

		// publish without goroutene
		publish(msg, producer)

	}
}

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout

	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// Producer configuration

	// parse Kafka cluster version
	version, err := sarama.ParseKafkaVersion("2.4.0")
	if err != nil {
		panic(err)
	}

	// init config, enable errors and notifications
	config := sarama.NewConfig()
	config.Version = version
	config.Metadata.Full = true
	config.ClientID = "test-producer"
	config.Producer.Return.Successes = true

	// Kafka SASL configuration
	if conect_to_Instaclutr_cluster {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = saslusername
		config.Net.SASL.Password = saslpassword
		config.Net.SASL.Handshake = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	}
	// init producer
	broker := []string{brokers}
	producer, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		panic(err)
	}

	return producer, err

}

func publish(message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topics,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
}
