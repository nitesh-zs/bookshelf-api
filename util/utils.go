package util

import (
	"golang.org/x/net/context"
	"google.golang.org/api/idtoken"
	"strconv"
	"strings"

	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
)

func GetTokenData(ctx *krogo.Context) (*model.User, error) {
	tData := &model.User{}

	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		return nil, errors.Unauthenticated{}
	}

	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	payload, _ := idtoken.Validate(context.Background(), token, "")

	tData.Name = payload.Claims["name"].(string)
	tData.Email = payload.Claims["email"].(string)

	return tData, nil
}

// Pagination filters page query parameter and returns page instance
func Pagination(ctx *krogo.Context) (*model.Page, error) {
	page := &model.Page{}

	size := ctx.Param("size")
	if size == "" {
		page.Size = model.DefaultPageSize
	} else {
		size, err := strconv.Atoi(size)
		if err != nil {
			return page, errors.InvalidParam{Param: []string{"size"}}
		}

		page.Size = size
	}

	p := ctx.Param("page")
	if p == "" {
		page.Offset = model.DefaultPageOffset
	} else {
		p, err := strconv.Atoi(p)
		if err != nil {
			return page, errors.InvalidParam{Param: []string{"offset"}}
		}

		page.Offset = (p - 1) * page.Size
	}

	return page, nil
}
