package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type db struct {
	ctx    context.Context
	client *redis.Client
	key    string
}

var board = &db{
	ctx: context.Background(),
	client: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}),
	key: "place",
}

func (b *db) getPixel(row, col int) (int, error) {
	offset := int64(row*width + col)
	val, err := b.client.GetBit(b.ctx, b.key, offset).Result()
	// GETBIT returns 0 on missing key, err is genuine
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func (b *db) togglePixel(row, col int) (int, error) {
	offset := row*width + col
	val, err := b.client.BitField(
		b.ctx, b.key, "INCRBY", "u1", offset, 1,
	).Result()
	// BITFIELD create on missing key, err is genuine
	if err != nil {
		return 0, err
	}
	return int(val[0]), nil
}

func (b *db) getBoard() (string, error) {
	val, err := b.client.Get(b.ctx, b.key).Result()
	// GET treats missing key as error
	if err == redis.Nil {
		val, err = "", nil
	}
	return val, err
}
