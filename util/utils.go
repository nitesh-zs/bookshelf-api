package util

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
)

func GetTokenData(ctx *krogo.Context) (*model.User, error) {
	tData := model.User{}

	token, err := ctx.Request().Cookie("auth")
	if err != nil {
		ctx.Logger.Error(err)
		return nil, errors.Unauthenticated{}
	}

	client := http.DefaultClient

	req, _ := http.NewRequest(http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token.Value)

	res, err := client.Do(req)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, errors.InternalServerErr{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, errors.InternalServerErr{}
	}

	err = json.Unmarshal(body, &tData)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, errors.InternalServerErr{}
	}

	return &tData, nil
}

// Pagination filters page query parameter and returns page instance
func Pagination(ctx *krogo.Context) (model.Page, error) {
	page := model.Page{}

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

// FilterList returns a list of available filters for querying books
func FilterList() []string {
	return []string{"genre", "author", "year", "language"}
}
