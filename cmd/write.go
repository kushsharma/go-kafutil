package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/kushsharma/go-kafutil/config"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func initWriter(conf config.App) *cobra.Command {
	thisCmd := &cobra.Command{
		Use: "write",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fds := &descriptorpb.FileDescriptorSet{}
			descSetBytes, err := ioutil.ReadFile(conf.DescriptorSetPath)
			if err != nil {
				panic(err)
			}
			if err := proto.Unmarshal(descSetBytes, fds); err != nil {
				panic(err)
			}

			files, _ := protodesc.NewFiles(fds)
			sampleReplaceDescriptor, err := files.FindDescriptorByName(protoreflect.FullName(conf.Schema))
			if err != nil {
				panic(err)
			}
			sampleReplaceDescriptorMessage, ok := sampleReplaceDescriptor.(protoreflect.MessageDescriptor)
			if !ok {
				log.Fatal("Unable to assert into MessageDescriptor")
			}

			sampleReplaceMessage := dynamicpb.NewMessage(sampleReplaceDescriptorMessage)
			sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("hakai"), protoreflect.ValueOfString("beeru"))
			sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("rasengan"), protoreflect.ValueOfString("naru"))
			sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("over"), protoreflect.ValueOfInt64(int64(rand.Intn(10000))))
			messageInBytes, err := proto.Marshal(sampleReplaceMessage)
			if err != nil {
				panic(err)
			}

			// if err := writeMessageToKafka(conf, []byte("key-3"), nil); err != nil {
			// 	panic(err)
			// }

			if err := writeMessageToKafka(conf, []byte("key-2"), messageInBytes); err != nil {
				panic(err)
			}

			return nil
		},
	}
	return thisCmd
}

func writeMessageToKafka(conf config.App, key, msg []byte) error {

	// make a writer that produces to topic-A, using the least-bytes distribution
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{conf.Host},
		Topic:    conf.Topic,
		Balancer: &kafka.LeastBytes{},
	})

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   key,
			Value: msg,
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	return nil
}

func writeDummyToKafka() error {

	// make a writer that produces to topic-A, using the least-bytes distribution
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "kwrite-test-A",
		Balancer: &kafka.LeastBytes{},
	})

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	return nil
}
