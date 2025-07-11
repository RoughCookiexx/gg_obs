package main

import (
	"log"
	"os"
	"obs/obscontrol"
)

func main() {
	client := obscontrol.NewOBSClient(os.Getenv("OBS_PASSWORD"))
	err := client.Connect(os.Getenv("OBS_CLIENT_ADDRESS"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	log.Println("Trying to switch to 'Be right back'")
	err = client.SwitchScene("pee pee or poo poo")
	if err != nil {
		log.Println("failed to switch scene:", err)
	}
}
