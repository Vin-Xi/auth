package main

import (
	"context"
	"fmt"
	"os"

	internal "github.com/Vin-Xi/auth/internal/database"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	fmt.Print(databaseUrl)
	_, err := internal.InitDB(ctx, databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "db initialization failed: %v", err)
		os.Exit(1)
	} else {
		fmt.Println("Connection is successful!")
	}

}
