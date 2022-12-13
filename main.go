package main

import (
	"fmt"
	"log"
	"os"

	"github.com/FahsaRwatit/wechatbot-chatgpt/config"
	"github.com/FahsaRwatit/wechatbot-chatgpt/openai"
	"github.com/FahsaRwatit/wechatbot-chatgpt/session"
)

func main() {
	configPath, err := os.UserConfigDir()
	fmt.Println(configPath)
	persistentConfig, err := config.LoadOrCreatePersistentConfig()
	if err != nil {
		log.Fatalf("Couldn't load config: %v", err)
	}
	if persistentConfig.OpenAISession == "" {
		token, err := session.GetSession()
		if err != nil {
			log.Fatalf("Couldn't get OpenAI session: %v", err)
		}
		fmt.Println(token)
		if err = persistentConfig.SetSessionToken(token); err != nil {
			log.Fatalf("Couldn't save OpenAI session: %v", err)
		}
	}
	chatGPT := openai.Init(persistentConfig)
	log.Println("Started ChatGPT")

	fmt.Println(chatGPT)
}
