package main

import (
	"context"
	"fmt"

	"github.com/benjamin-maynard/kubernetes-cloud-mysql-backup/internal/config"
)

func main() {

	ctx := context.Background()

	_, err := config.NewConfigFromEnvironment(ctx)
	if err != nil {
		fmt.Println(err)
	}

}
