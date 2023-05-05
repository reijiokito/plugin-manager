package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	go_pdk "github.com/reijiokito/go-pdk"
	"google.golang.org/protobuf/proto"
	"log"
	"strings"
)

var Connection *nats.Conn
var JetStream nats.JetStreamContext

type Config struct {
	NatsUrl      string `yaml:"natsUrl"`
	NatsUsername string `yaml:"natsUsername"`
	NatsPassword string `yaml:"natsPassword"`
}

func New() interface{} {
	return &Config{
		NatsUrl:      "127.0.0.1",
		NatsUsername: "",
		NatsPassword: "",
	}
}

func (conf Config) Access(pdk *go_pdk.PDK) {

	fmt.Println(conf)

	var err error
	var nats_ []string
	for _, n := range strings.Split(conf.NatsUrl, ",") {
		fmt.Printf("Nats configuration: nats://%s:4222\n", n)
		nats_ = append(nats_, fmt.Sprintf("nats://%s:4222", n))
	}

	if conf.NatsUsername != "" && conf.NatsPassword != "" {
		Connection, err = nats.Connect(strings.Join(nats_, ","), nats.UserInfo(conf.NatsUsername, conf.NatsPassword))
	} else {
		Connection, err = nats.Connect(strings.Join(nats_, ","))
	}

	if err != nil {
		log.Println("Can not connect to NATS:", nats_)
	}

	/*init jetstream*/
	JetStream, err = Connection.JetStream()
	if err != nil {
		log.Println(err)
	}
}

func Publish(args ...interface{}) {
	subject := args[0].(string)
	data := args[1].([]byte)
	Connection.Publish(subject, data)
	fmt.Println(fmt.Sprintf("Publish data: %v - %v", subject, data))
}

func Subscribe(args ...interface{}) {
	subject := args[0].(string)

	Connection.Subscribe(subject, func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	fmt.Println(fmt.Sprintf("Subcribe subject: %v", subject))
}

type eventStream struct {
	sender    string
	receiver  string
	executors map[string]func(m *nats.Msg)
}

func PostEvent(args ...interface{}) {
	subject := args[0].(string)
	data := args[1].(proto.Message)
	if data, err := proto.Marshal(data); err == nil {
		JetStream.Publish(subject, data)

	}
}

type SubjectHandler func(data proto.Message)

var eventStreams map[string]*eventStream = make(map[string]*eventStream)

//
//func RegisterNats[R proto.Message](args ...interface{}) {
//	subject := args[0].(string)
//	handler := args[1].(SubjectHandler[R])
//
//	parts := strings.Split(subject, ".")
//	stream := createOrGetEventStream(parts[0])
//	log.Println(fmt.Sprintf("Events: subject = %s, receiver = %s", subject, stream.receiver))
//	var event R
//	ref := reflect.New(reflect.TypeOf(event).Elem())
//	event = ref.Interface().(R)
//
//	stream.executors[subject] = func(m *nats.Msg) {
//		if err := proto.Unmarshal(m.Data, event); err == nil {
//			handler(event)
//		} else {
//			log.Print("Error in parsing data nats:", err)
//		}
//	}
//}

func createOrGetEventStream(sender string) *eventStream {
	if stream, ok := eventStreams[sender]; ok {
		return stream
	}

	stream := &eventStream{
		sender:    sender,
		receiver:  "manager",
		executors: make(map[string]func(m *nats.Msg)),
	}

	eventStreams[sender] = stream
	return stream
}

func (es *eventStream) start(JetStream nats.JetStreamContext) {
	sub, err := JetStream.PullSubscribe("", es.receiver, nats.BindStream(es.sender))

	if err != nil {
		log.Println("Error in start event stream - sender ", es.sender, "- receiver ", es.receiver, " : ", err.Error())
	}

	go func() {
		for {
			if messages, err := sub.Fetch(1); err == nil {
				if len(messages) == 1 {
					m := messages[0]
					if executor, ok := es.executors[m.Subject]; ok {
						executor(m)
					}
					m.Ack()
				}
			}
		}
	}()
}

func StartEventStream(args ...interface{}) {
	for _, e := range eventStreams {
		e.start(JetStream)
	}
}

func Release(args ...interface{}) {
	Connection.Close()
}

func GetServices() map[string]func(...interface{}) {
	services := make(map[string]func(...interface{}))
	services["Publish"] = Publish
	services["Subscribe"] = Subscribe
	services["Release"] = Release
	services["PostEvent"] = PostEvent
	services["StartEventStream"] = StartEventStream

	return services
}

func GetCallers() map[string]func(...interface{}) interface{} {
	callers := make(map[string]func(...interface{}) interface{})

	return callers
}
