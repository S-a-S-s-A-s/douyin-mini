package controller

import (
	"douyin-mini/db"
	"douyin-mini/db/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"time"
)

type DouyinFeedRequest struct {
	LatestTime int64  `query:"latest_time"`
	Token      string `query:"token"`
}

type DouyinFeedResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMsg  string     ` json:"status_msg"`
	VideoList  []db.Video `json:"video_list"`
	NextTime   int64      ` json:"next_time"`
}

func Feed(c *gin.Context) {
	t, _ := strconv.ParseInt(c.Query("latest_time"), 10, 32)
	var req = DouyinFeedRequest{LatestTime: t, Token: c.Query("token")}
	videos := dao.GetVideos(time.Unix(req.LatestTime, 0), "")
	if len(videos) == 0 {
		c.JSON(http.StatusOK, DouyinFeedResponse{
			StatusCode: -1,
			StatusMsg:  "no new videos",
			VideoList:  nil,
			NextTime:   t,
		})
	}
	c.JSON(http.StatusOK, DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videos,
		NextTime:   videos[len(videos)-1].CreateTime.Unix(),
	})
}
