package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/codeboris/pocket_telegram_bot/pkg/config"
	"github.com/codeboris/pocket_telegram_bot/pkg/repository"
	"github.com/codeboris/pocket_telegram_bot/pkg/repository/boltdb"
	"github.com/codeboris/pocket_telegram_bot/pkg/server"
	"github.com/codeboris/pocket_telegram_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	botApi.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	dbApi, err := initBolt(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tRepository := boltdb.NewTokenRepository(dbApi)

	telegramBot := telegram.NewBot(botApi, pocketClient, tRepository, cfg.AuthServerURL, cfg.Messages)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tRepository, cfg.TelegramBotURL)

	go func() {
		if err := authorizationServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initBolt(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
