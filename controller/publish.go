package controller

import (
	"bytes"
	"douyin-mini/db"
	"douyin-mini/db/dao"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"time"
)

type VideoListResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMsg  string     ` json:"status_msg"`
	VideoList  []db.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.Query("title")
	if _, exist := db.Token[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := dao.FindUserId(db.Token[token])
	//文件名
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	//保存文件到本地
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//保存封面地址到本地
	if err := GetSnapshot("./public/"+finalName, finalName, 1); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//插入video
	var video = db.Video{
		AuthorId:   user.ID,
		AuthorName: user.Username,
		CreateTime: time.Now(),
		PlayURL:    "http://192.168.224.66:80/" + finalName,
		CoverURL:   "http://192.168.224.66:80/" + finalName + ".png",
		Title:      title,
	}
	dao.CreateVideo(video)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")
	if _, exist := db.Token[token]; exist {
		var videos = dao.FindAllVideo(id)
		c.JSON(http.StatusOK, VideoListResponse{
			StatusCode: 0,
			VideoList:  videos,
		})
	} else {
		c.JSON(http.StatusOK, VideoListResponse{
			StatusCode: 1,
			StatusMsg:  "token is failed",
			VideoList:  nil,
		})
	}

}

func GetSnapshot(videoPath, snapshotPath string, frameNum int) (err error) {
	snapshotPath = "./public/" + snapshotPath
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	return nil
}
