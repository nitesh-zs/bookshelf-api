package main

import (
	"github.com/krogertechnology/krogo/cmd/krogo/migration"
	dbmigration "github.com/krogertechnology/krogo/cmd/krogo/migration/dbMigration"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/handler/auth"
	bHandler "github.com/nitesh-zs/bookshelf-api/handler/book"

	// "github.com/nitesh-zs/bookshelf-api/middleware"
	"github.com/nitesh-zs/bookshelf-api/migrations"
	bSvc "github.com/nitesh-zs/bookshelf-api/service/book"
	uSvc "github.com/nitesh-zs/bookshelf-api/service/user"
	"github.com/nitesh-zs/bookshelf-api/store/book"
	"github.com/nitesh-zs/bookshelf-api/store/user"
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

	//nolint:gocritic //will remove the code later upon finalization of client side auth flow
	// conf := &oauth2.Config{
	//	ClientID:     k.Config.Get("CLIENT_ID"),
	//	ClientSecret: k.Config.Get("CLIENT_SECRET"),
	//	RedirectURL:  k.Config.Get("REDIRECT_URL"),
	//	Scopes: []string{
	//		"email",
	//		"openid",
	//		"profile",
	//	},
	//	Endpoint: google.Endpoint,
	// }

	// set auth middleware
	// k.Server.UseMiddleware(middleware.Login(conf), middleware.Redirect(conf),
	// middleware.ValidateToken(conf), middleware.Logout)

	uStore := user.New()
	userSvc := uSvc.New(uStore)
	authHandler := auth.New(userSvc)

	bookStore := book.New()
	bookSvc := bSvc.New(bookStore)
	bookHandler := bHandler.New(bookSvc)

	k.GET("/hello", func(c *krogo.Context) (interface{}, error) {
		return "hello", nil
	})

	k.GET("/register", authHandler.Register)
	k.GET("/book", bookHandler.Get)
	k.GET("/book/{id}", bookHandler.GetByID)
	k.GET("/list/{param}", bookHandler.GetFilters)
	k.Start()
}
