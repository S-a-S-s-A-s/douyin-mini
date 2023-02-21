package controller

import (
	"douyin-mini/db"
	"douyin-mini/db/dao"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentListResponse struct {
	Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment CommentResponse `json:"comment,omitempty"`
}

// CommentResponse 评论信息的响应结构体
type CommentResponse struct {
	ID         int64         `json:"id,omitempty"`
	Content    string        `json:"content,omitempty"`
	CreateDate string        `json:"create_date,omitempty"`
	User       UserResponses `json:"user,omitempty"`
}

// UserResponse 用户信息的响应结构体
type UserResponses struct {
	ID            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	// IsFollow      bool   `json:"is_follow,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	token := c.Query("token")
	videoid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType := c.Query("action_type")
	userName := c.Query("name")

	// getUserID, _ := c.Get("user_id")

	// if user, exist := usersLoginInfo[token]; exist {
	// 	if actionType == "1" {
	// 		text := c.Query("comment_text")
	// 		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
	// 			Comment: Comment{
	// 				Id:         1,
	// 				User:       user,
	// 				Content:    text,
	// 				CreateDate: "05-01",
	// 			}})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, Response{StatusCode: 0})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }

	//非合法操作类型
	if userid, exist := db.Token[token]; exist {
		if actionType != "1" && actionType != "2" {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Unsupported actionType",
			})
			c.Abort()
			return
		}
		//合法操作
		//1:发布评论，2:删除评论
		if actionType == "1" { // 发布评论
			text := c.Query("comment_text")
			PostComment(c, userid, userName, text, int64(videoid))
		} else if actionType == "2" { //删除评论
			commentIDStr := c.Query("comment_id")
			commentID, _ := strconv.ParseInt(commentIDStr, 10, 64)
			DeleteComment(c, int64(videoid), int64(commentID))
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// PostComment 发布评论
func PostComment(c *gin.Context, userID int64, userName string, text string, videoID int64) {

	newComment := db.Comment{
		VideoId:    videoID,
		UserId:     userID,
		UserName:   userName,
		Content:    text,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	//发布评论并改变评论数量，获取video作者信息
	err1 := db.DB.Transaction(func(db *gorm.DB) error {
		if err := dao.PostComment(newComment); err != nil {
			return err
		}
		if err := dao.AddCommentCount(videoID); err != nil {
			return err
		}
		return nil
	})

	var getUser db.User
	err2 := dao.GetAUserByID(userID, &getUser)
	_, err3 := dao.GetVideoAuthorID(videoID)

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Failed to post comment",
		})
		// c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "fail to post a comment"})
		c.Abort()
		return
	} else {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "post the comment successfully",
			},
			Comment: CommentResponse{
				ID:         newComment.ID,
				Content:    newComment.Content,
				CreateDate: newComment.CreateTime,
				User: UserResponses{
					ID:            getUser.ID,
					Name:          getUser.Username,
					FollowCount:   getUser.FollowCount,
					FollowerCount: getUser.FollowerCount,
				},
			},
		})
	}
}

// 删除评论
func DeleteComment(c *gin.Context, videoID int64, commentID int64) {

	//删除评论并改变评论数量，获取video作者信息
	err := db.DB.Transaction(func(db *gorm.DB) error {
		if err := dao.DeleteComment(commentID); err != nil {
			return err
		}
		if err := dao.ReduceCommentCount(videoID); err != nil {
			return err
		}
		return nil
	})
	//响应处理
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "fail to delete the comment"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "delete the comment successfully"})
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	// c.JSON(http.StatusOK, CommentListResponse{
	// 	Response:    Response{StatusCode: 0},
	// 	CommentList: DemoComments,
	// })

	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)

	//获取指定videoid的评论表
	commentList, err := dao.GetCommentList(int64(videoID))

	//评论表不存在
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Failed to get commentList",
		})
		c.Abort()
		return
	}

	//评论表存在
	var responseCommentList []CommentResponse
	for i := 0; i < len(commentList); i++ {
		// getUser, err1 := service.GetUser(commentList[i].UserID)
		var getUser db.User
		err1 := dao.GetAUserByID(commentList[i].UserId, &getUser)

		if err1 != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Failed to get commentList.",
			})
			c.Abort()
			return
		}
		responseComment := CommentResponse{
			ID:         commentList[i].ID,
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreateTime, // mm-dd
			User: UserResponses{
				ID:            getUser.ID,
				Name:          getUser.Username,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
			},
		}
		responseCommentList = append(responseCommentList, responseComment)

	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "Got the comment list.",
		},
		CommentList: responseCommentList,
	})
}
