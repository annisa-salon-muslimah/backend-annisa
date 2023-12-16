package main

import (
	_ "annisa-salon/docs"
	"annisa-salon/handler"
)

func main() {
	// @title Sweager Service API
	// @description Sweager service API in Go using Gin framework
	// @host backend-annisa-production.up.railway.app
	// @securitydefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	handler.StartApp()
}
