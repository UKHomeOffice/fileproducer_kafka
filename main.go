package main

import (
	"fmt"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Broker struct {
	topic     string
	partition int
	hostname  string
}

func newBroker(hostname string, topic string, partition int) *Broker {
	return &Broker{
		topic:     topic,
		partition: partition,
		hostname:  hostname,
	}
}

func isFileValid(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return false
	}
	mode := fi.Mode()
	if mode.IsRegular() {
		return true
	}
	return false
}

func produceFile(broker kafka.Client, topic string, partition int, data []byte) {
	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{Value: data}
	if _, err := producer.Produce(topic, int32(partition), msg); err != nil {
		log.Fatalf("cannot produce message to %s:%d: %s", topic, partition, err)
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	file := os.Getenv("FILE")
	partition, err1 := strconv.Atoi(os.Getenv("PARTITION"))
	check(err1)
	topic := os.Getenv("TOPIC")
	kafkaAddrs := strings.Split(os.Getenv("BROKERS"), ",")
	name := os.Getenv("CLIENT_NAME")
	broker, err2 := kafka.Dial(kafkaAddrs, kafka.NewBrokerConf(name))
	check(err2)
	defer broker.Close()

	if file != "" {
		f, err3 := os.Open(file)
		check(err3)
		data, err4 := ioutil.ReadAll(f)
		check(err4)
		produceFile(broker, topic, partition, data)

	} else {
		fmt.Println("File not defined")
	}

}
