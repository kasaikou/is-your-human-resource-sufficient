package main

import (
	"context"

	"github.com/kasaikou/is-your-human-resource-sufficient/controllers/console"
	"github.com/kasaikou/is-your-human-resource-sufficient/services/maps/seq_map"
)

func main() {
	console.NewGameContext(context.Background(), seq_map.Default)
}
