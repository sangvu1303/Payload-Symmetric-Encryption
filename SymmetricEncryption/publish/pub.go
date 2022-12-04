package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func mqttBegin(broker string, user string, pw string) mqtt.Client {
	var opts *mqtt.ClientOptions = new(mqtt.ClientOptions)

	opts = mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf(broker))
	opts.SetUsername(user)
	opts.SetPassword(pw)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func main() {
	fmt.Println("Encryption Program v0.01")

	text := []byte("Zao shang hao zung wo,\nShenzai wai-oh bing chilling.\nWo han she huay bing chilling.")
	key := []byte("passphrasewhichneedstobe32bytes!")

	c, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	asmEn := gcm.Seal(nonce, nonce, text, nil)

	/*Ma hoa sau do moi publish du lieu*/
	mqttClient = mqttBegin("localhost:1883", "Sang", "2000")
	for {
		mqttClient.Publish("test/encrypt", 0, false, asmEn)
		fmt.Print("publish: ")
		fmt.Println(asmEn)
		fmt.Println("==================================================")
		time.Sleep(3 * time.Second)
	}
}          
