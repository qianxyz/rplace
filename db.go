// The database access layer.
//
// The board is stored in Redis as a BITFIELD in row-major order.
// Every pixel is represented by a 4-bit color.
// See the commands at https://redis.io/commands/bitfield.

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

// Get the pixel color at the row and column specified.
func (b *db) getPixel(row, col int) (int, error) {
	// Color is 4-bit int, the bit offset is [index] * 4.
	offset := (row*width + col) * 4

	val, err := b.client.BitField(b.ctx, b.key, "GET", "u4", offset).Result()
	// On OOB access, BITFIELD GET doesn't err and return 0 (which is what
	// we want). If the operation does err then it's the database.
	if err != nil {
		return 0, err
	}

	// Client.Bitfield returns a []int64 because the BITFIELD command may
	// contain several operations (GET/SET/INCRBY). In our case we only
	// issue a single GET so the result is in [0].
	return int(val[0]), nil
}

// Set the pixel color at the row and column specified.
func (b *db) setPixel(row, col int, color int) error {
	// Color is 4-bit int, the bit offset is [index] * 4.
	offset := (row*width + col) * 4

	_, err := b.client.BitField(b.ctx, b.key, "SET", "u4", offset, color).Result()
	// On OOB access, BITFIELD SET zero-pads the value.
	// If the operation does err then it's the database.
	return err
}

// Get the whole board as a bytestring.
func (b *db) getBoard() (string, error) {
	val, err := b.client.Get(b.ctx, b.key).Result()

	// In case the key does not exist, GET will return redis.Nil.
	// This is actually fine since we can pretend it exists
	// and return the empty string as default.
	// This is not a problem for the frontend which also treats OOB as 0.
	//
	// Moreover we don't even need to insert; The two other operations,
	// BITFIELD GET and SET, will not err on missing key.
	if err == redis.Nil {
		val, err = "", nil
	}
	return val, err
}
