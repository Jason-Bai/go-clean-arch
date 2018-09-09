package user

import (
	. "github.com/Jason-Bai/go-clean-arch/handler"
	"github.com/Jason-Bai/go-clean-arch/model"
	"github.com/Jason-Bai/go-clean-arch/pkg/auth"
	"github.com/Jason-Bai/go-clean-arch/pkg/errno"
	"github.com/Jason-Bai/go-clean-arch/pkg/token"

	"github.com/gin-gonic/gin"
)

// @Summary Login generates the authentication token
// @Produce  json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MjgwMTY5MjIsImlkIjowLCJuYmYiOjE1MjgwMTY5MjIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LjxrK9DuAwAzUD8-9v43NzWBN7HXsSLfebw92DKd1JQ"}}"
// @Router /login [post]
func Login(c *gin.Context) {
	var u model.UserModel

	if err := c.Bind(&u); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	d, err := model.GetUser(u.Username)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	if err := auth.Compare(d.Password, u.Password); err != nil {
		SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	t, err := token.Sign(c, token.Context{ID: d.ID, Username: d.Username}, "")
	if err != nil {
		SendResponse(c, errno.ErrToken, nil)
		return
	}

	SendResponse(c, nil, model.Token{Token: t})
}
