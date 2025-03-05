package main

import (
	"log"

	"YANDEX_GOLANG/internal/agent"
)

func main() {
	agent := agent.NewAgent()
	log.Println("Agent is Starting...")
	agent.Start()
}
