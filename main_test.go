package main

import (
	"context"
	"testing"
)

func Test_processItems(t *testing.T) {
	ctx := context.TODO()
	processItems(ctx, 10)
}

func Test_readConfig(t *testing.T) {
	var cfg Config
	readFile(&cfg)
	println(cfg.NoOfItems)
}
