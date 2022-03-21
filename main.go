package main

import "github.com/krogertechnology/krogo/pkg/krogo"

func main() {
	k := krogo.New()

	k.Server.ValidateHeaders = false

	// enabling /swagger endpoint for Swagger UI
	k.EnableSwaggerUI()

	k.Start()
}
