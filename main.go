package main

import (
	"fmt"
	"lab2/client"
	"lab2/server"
	"log"
)

func main() {
	go server.RunServer()
	fmt.Println("Какой сценарий запустить")
	var n int
	fmt.Scan(&n)
	switch n {
	case 1:
		client.Scenario1()
	case 2:
		client.Scenario2()
	default:
		log.Fatalln("Указан неизвестный номер сценария")
	}

}
