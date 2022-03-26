package util

import (
	"encoding/json"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
	"io"
	"net/http"
)

func GetTokenData(ctx *krogo.Context) (*model.User, error) {
	tData := model.User{}

	token, err := ctx.Request().Cookie("auth")
	if err != nil {
		ctx.Logger.Error(err)
		return nil, errors.Unauthenticated{}
	}

	client := http.DefaultClient

	req, _ := http.NewRequest(http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token.Value)

	res, err := client.Do(req)
	if err != nil {
		ctx.Logger.Error(err)
		return nil, errors.InternalServerErr{}
	}

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
