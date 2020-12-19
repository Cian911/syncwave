package main

import (
	"context"
	"github.com/cian911/raspberry-pi-provisioner/pkg/cli/syncwave"
)

func main() {
	ctx := context.Background()
	syncwave.New().ExecuteContext(ctx)
}
