package main

import "Crawler/internal/handler"

func main() {
	handler.Crawler_handler()

	select {}
}
