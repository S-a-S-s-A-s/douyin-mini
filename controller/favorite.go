package controller

import (
	"douyin-mini/db"
	"douyin-mini/db/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid

func FavoriteAction(c *gin.Context) {
	videoid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	token := c.Query("token")
	actiontype := c.Query("action_type")
	if userid, exist := db.Token[token]; exist {
		if actiontype == "1" {
			dao.LikeVideo(userid, videoid)
		} else {
			dao.NoLikeVideo(userid, videoid)
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userid, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")

	if myid, exist := db.Token[token]; exist {
		videos := dao.FindAllLikeVideo(userid, myid)
		c.JSON(http.StatusOK, VideoListResponse{
			StatusCode: 0,
			VideoList:  videos,
		})
	} else {
		c.JSON(http.StatusOK, VideoListResponse{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
			VideoList:  nil,
		})
	}
}
