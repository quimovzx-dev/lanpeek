package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/quimovzx-dev/lanpeek/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := cmd.NewRoot().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
