package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func produceFile(broker kafka.Client, topic string, partition int, data []byte, key []byte) {
	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{
		Value: data,
		Key:   key}
	if _, err := producer.Produce(topic, int32(partition), msg); err != nil {
		log.Fatalf("cannot produce message to %s:%d: %s", topic, partition, err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getAllSubdirectories(parentPath string) (paths []string, err error) {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	}
	err = filepath.Walk(parentPath, walkFn)
	return paths, err
}

func processFile(file string) {

	partition, err1 := strconv.Atoi(os.Getenv("PARTITION"))
	if err1 != nil {
		partition = 0
	}

	topic := os.Getenv("TOPIC")
	kafkaAddrs := strings.Split(os.Getenv("BROKERS"), ",")
	name := os.Getenv("CLIENT_NAME")
	broker, err2 := kafka.Dial(kafkaAddrs, kafka.NewBrokerConf(name))
	check(err2)
	defer broker.Close()
	if file != "" {
		f, err3 := os.Open(file)
		if err3 == nil {
			data, err4 := ioutil.ReadAll(f)
			check(err4)
			produceFile(broker, topic, partition, data, []byte(file))
		}

	} else {
		fmt.Println("File not defined")
	}
}

func main() {
	sourcedir := os.Getenv("SOURCEDIR")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != fsnotify.Chmod {
					log.Println("Pushing Files: " + event.Name)
					processFile(event.Name)
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	if sourcedir != "" {
		paths, err := getAllSubdirectories(sourcedir)
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range paths {
			watcher.Add(path)
		}
	} else {
		fmt.Println("Please, enter a valid SOURCEDIR to inspect")
	}

	// Daemon mode
	for {
		time.Sleep(1000 * time.Millisecond)
	}

}
