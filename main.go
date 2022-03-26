package main

import (
	"github.com/krogertechnology/krogo/cmd/krogo/migration"
	dbmigration "github.com/krogertechnology/krogo/cmd/krogo/migration/dbMigration"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/handler/auth"
	"github.com/nitesh-zs/bookshelf-api/middleware"
	"github.com/nitesh-zs/bookshelf-api/migrations"
	uSvc "github.com/nitesh-zs/bookshelf-api/service/user"
	"github.com/nitesh-zs/bookshelf-api/store/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	k := krogo.New()

	k.Server.ValidateHeaders = false

	// enabling /swagger endpoint for Swagger UI
	k.EnableSwaggerUI()

	err := migration.Migrate("bookshelf-api", dbmigration.NewGorm(k.GORM()), migrations.All(), "UP", k.Logger)
	if err != nil {
		return
	}

	conf := &oauth2.Config{
		ClientID:     k.Config.Get("CLIENT_ID"),
		ClientSecret: k.Config.Get("CLIENT_SECRET"),
		RedirectURL:  k.Config.Get("REDIRECT_URL"),
		Scopes: []string{
			"email",
			"openid",
			"profile",
		},
		Endpoint: google.Endpoint,
	}

	// set auth middleware
	k.Server.UseMiddleware(middleware.Login(conf), middleware.Redirect(k.Logger, conf),
		middleware.ValidateToken(k.Logger, conf), middleware.Logout)

	uStore := user.New()
	uSvc := uSvc.New(uStore)
	authHandler := auth.New(uSvc)

	k.GET("/hello", func(c *krogo.Context) (interface{}, error) {
		return "hello", nil
	})

	k.GET("/register", authHandler.Register)

	k.Start()
}
