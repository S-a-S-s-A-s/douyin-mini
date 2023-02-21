package controller

import (
	"douyin-mini/db"
	"douyin-mini/db/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User db.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := dao.FindUserName(username)
	if user.ID != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}
	token := username + password
	db.Token[token] = dao.CreateUser(username, password)

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   db.Token[token],
		Token:    token,
	})

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user := dao.FindUserPassword(username, password)
	if user.ID == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	token := username + password
	db.Token[token] = user.ID
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.ID,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")
	if _, exist := db.Token[token]; exist {
		user := dao.FindUserId(id)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
