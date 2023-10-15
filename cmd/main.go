package main

import "github.com/z416352/Crawler/internal/handler"

func main() {
	handler.Crawler_handler()

	select {}
}
