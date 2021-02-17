package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/golang/protobuf/ptypes"
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

			// use this to push oob date for BQ testing
			// evtTsProto, _ := ptypes.TimestampProto(time.Date(9999, 11, 11, 24, 24, 24, 24, time.UTC))
			// eventTimestampValue := evtTsProto.ProtoReflect()

			// default ts
			eventTimestampValue := ptypes.TimestampNow().ProtoReflect()
			_ = eventTimestampValue

			sampleReplaceMessage := dynamicpb.NewMessage(sampleReplaceDescriptorMessage)
			sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("hakai"), protoreflect.ValueOfString("beeru"))
			sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("rasengan"), protoreflect.ValueOfString("naru"))
			sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("over"), protoreflect.ValueOfInt64(int64(rand.Intn(10000))))
			sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("event_timestamp"), protoreflect.ValueOfMessage(eventTimestampValue))

			// to cause invalidProtocolBufferException
			// sampleReplaceMessage.Set(sampleReplaceDescriptorMessage.Fields().ByName("event_timestamp"), protoreflect.ValueOfString("invalid-timestamp"))

			messageInBytes, err := proto.Marshal(sampleReplaceMessage)
			if err != nil {
				panic(err)
			}

			fmt.Printf("writing message to topic %s, encoding with proto %s...\n", conf.Topic, conf.Schema)

			// to test null message value
			// if err := writeMessageToKafka(conf, []byte("key-3"), nil); err != nil {
			// 	panic(err)
			// }

			if err := writeMessageToKafka(conf, []byte(fmt.Sprintf("key-%d", int64(rand.Intn(1000)))), messageInBytes); err != nil {
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

// writeDummyToKafka - DEPRECATE
// sample function that demonstrate writing to kafka
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
