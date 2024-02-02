package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/mauricioabreu/mosaic-video/internal/config"
	"github.com/mauricioabreu/mosaic-video/internal/mosaic"
	"github.com/mauricioabreu/mosaic-video/internal/worker"
	"github.com/spf13/cobra"
)

func Simulate() *cobra.Command {
	return &cobra.Command{
		Use:   "simulate",
		Short: "Simulate a mosaic processing using a queue",
		Run: func(cmd *cobra.Command, args []string) {
			if err := godotenv.Load(".env"); err != nil {
				log.Println("Could not load .env file")
			}

			cfg, err := config.NewConfig()
			if err != nil {
				log.Fatal(err)
			}

			redisAddress := cfg.Redis.Host + ":" + cfg.Redis.Port
			client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddress})
			defer client.Close()

			if err := enqueueStartMosaicTask(client, cfg); err != nil {
				log.Println(err)
			}
		},
	}
}

func enqueueStartMosaicTask(client *asynq.Client, cfg *config.Config) error {
	fileData, err := os.ReadFile("./testing/tasks.json")
	if err != nil {
		return err
	}

	var mosaics []mosaic.Mosaic
	if err := json.Unmarshal(fileData, &mosaics); err != nil {
		return err
	}

	// Enqueue a task for each mosaic
	for _, m := range mosaics {
		payload, err := json.Marshal(worker.StartMosaicPayload{Mosaic: m})
		if err != nil {
			return err
		}

		_, err = client.Enqueue(
			asynq.NewTask(worker.TypeStartMosaic, payload),
			asynq.MaxRetry(cfg.MaxRetriesTasks),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
