package main

import (
	"GoBank/internal/app"
	"fmt"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Printf("error while trying to start app - %v\n", err)
	}
}
