package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/kushsharma/go-kafutil/config"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func initReader(conf config.App) *cobra.Command {
	thisCmd := &cobra.Command{
		Use: "read",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fds := &descriptorpb.FileDescriptorSet{}
			descSetBytes, err := ioutil.ReadFile(conf.DescriptorSetPath)
			if err != nil {
				panic(err)
			}
			if err := proto.Unmarshal(descSetBytes, fds); err != nil {
				panic(err)
			}

			descriptorFiles, err := protodesc.NewFiles(fds)
			if err != nil {
				panic(err)
			}
			fmt.Printf("reading messages from topic %s, decoding with proto %s...\n", conf.Topic, conf.Schema)
			if err := readMessageFromKafka(descriptorFiles, conf); err != nil {
				panic(err)
			}
			return nil
		},
	}
	return thisCmd
}

// decodeMessageValue deserialize message using descriptor
func decodeMessageValue(desc protoreflect.Descriptor, raw []byte) string {

	sampleReplaceDescriptorMessage, ok := desc.(protoreflect.MessageDescriptor)
	if !ok {
		log.Fatal("Unable to assert into MessageDescriptor")
	}

	sampleReplaceMessage := dynamicpb.NewMessage(sampleReplaceDescriptorMessage)
	if err := proto.Unmarshal(raw, sampleReplaceMessage); err != nil {
		log.Fatal("failed to parse error")
		panic(err)
	}

	return sampleReplaceMessage.String()
}

func readMessageFromKafka(files *protoregistry.Files, conf config.App) error {
	sampleReplaceDescriptor, err := files.FindDescriptorByName(protoreflect.FullName(conf.Schema))
	if err != nil {
		return err
	}

	// make a new reader that consumes
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{conf.Host},
		Topic:     conf.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	// to set offset
	//r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), decodeMessageValue(sampleReplaceDescriptor, m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	return nil
}

// readDummyFromKafka - DEPRECATE
// sample function that demonstrate reading from kafka
func readDummyFromKafka() error {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "kwrite-test-A",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	//r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	return nil
}
