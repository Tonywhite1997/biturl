package helper

import (
	"math/big"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node        *snowflake.Node
	once        sync.Once
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// InitializeSnowflake initializes the Snowflake node (call once at startup)
func InitializeSnowflake(nodeNumber int64) error {
	var err error
	once.Do(func() {
		node, err = snowflake.NewNode(nodeNumber)
	})
	return err
}

// GenerateShortCode generates a URL-friendly short code from Snowflake ID
func GenerateShortCode() string {
	if node == nil {
		panic("Snowflake node not initialized")
	}

	id := node.Generate() // get Snowflake ID (int64)
	return base62Encode(id.Int64())
}

// base62Encode converts int64 to base62 string
func base62Encode(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var result []byte
	n := big.NewInt(num)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for n.Cmp(zero) > 0 {
		n.DivMod(n, base, mod)
		result = append([]byte{base62Chars[mod.Int64()]}, result...)
	}
	return string(result)
}
