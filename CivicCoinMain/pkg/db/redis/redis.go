package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Client es el cliente Redis que se compartirá en todo el proyecto.
var Client *redis.Client
var Ctx context.Context

// Init inicializa el cliente Redis.
func Init(addr, password string, db int) error {
	Ctx = context.Background()
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Verificar la conexión
	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

// GetClient devuelve el cliente Redis.
func GetClient() *redis.Client {
	return Client
}
