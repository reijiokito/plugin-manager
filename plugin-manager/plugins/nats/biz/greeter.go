package biz

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"log"
	"time"
)

// Here is a real implementation of Greeter
type GreeterHello struct {
	Logger hclog.Logger
}

func (g *GreeterHello) Calculate(a, b int32) int32 {
	g.Logger.Debug("message from GreeterHello.Calculate")
	count := 0
	for i := 0; i < 10; i++ {
		log.Println(fmt.Sprintf("Hello: %d", count))
		count++
		time.Sleep(time.Second * 1)
	}
	return a + b
}

func (g *GreeterHello) Greet(name string) string {
	g.Logger.Debug("message from GreeterHello.Greet")

	return fmt.Sprintf("Hello, %s!", name)
}
