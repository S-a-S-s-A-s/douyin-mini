package dao

import (
	"douyin-mini/db"

	"gorm.io/gorm"

	"time"
)

func GetVideos(lasttime time.Time, token string) []db.Video {
	userid := db.Token[token]
	var result = make([]db.Video, 0)
	err := db.DB.Model(&db.Video{}).Where("create_time < ? and author_id != ?", lasttime, userid).Order("create_time DESC").Limit(30).Find(&result).Error
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(result); i++ {
		if err = db.DB.Model(&db.User{}).Where("id = ?", result[i].AuthorId).Find(&result[i].Author).Error; err != nil {
			panic(err)
		}

		//查看是否点赞
		var like db.Like
		if err = db.DB.Model(&db.Like{}).Where("video_id = ? and user_id = ?", result[i].ID, userid).Find(&like).Error; err != nil {
			panic(err)
		}
		if like.ID != 0 {
			result[i].IsFavorite = true
		}
	}

	return result
}

func CreateVideo(video db.Video) {
	if err := db.DB.Create(&video).Error; err != nil {
		println(err)
	}
}

func FindUserName(name string) db.User {
	var user db.User
	err := db.DB.Model(&db.User{}).Where("username = ?", name).Find(&user).Error
	if err != nil {
		panic(err)
	}
	return user
}

func CreateUser(name, password string) int64 {
	user := &db.User{Username: name, Password: password}
	if err := db.DB.Create(user).Error; err != nil {
		println(err)
	}
	return user.ID
}

func FindUserPassword(name, password string) db.User {
	var user db.User
	err := db.DB.Model(&db.User{}).Where("username = ? and password = ?", name, password).Find(&user).Error
	if err != nil {
		panic(err)
	}
	return user
}

func FindUserId(id int64) db.User {
	var user db.User
	err := db.DB.Model(&db.User{}).Where("id = ?", id).Find(&user).Error
	if err != nil {
		panic(err)
	}
	return user
}

// 列出主页所有视频
func FindAllVideo(id int64) []db.Video {
	var result = make([]db.Video, 0)
	err := db.DB.Model(&db.Video{}).Where("author_id = ?", id).Find(&result).Error
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(result); i++ {
		if err = db.DB.Model(&db.User{}).Where("id = ?", result[i].AuthorId).Find(&result[i].Author).Error; err != nil {
			panic(err)
		}
		//查看是否点赞
		var like db.Like
		if err = db.DB.Model(&db.Like{}).Where("video_id = ? and user_id = ?", result[i].ID, id).Find(&like).Error; err != nil {
			panic(err)
		}
		if like.ID != 0 {
			result[i].IsFavorite = true
		}
	}
	return result
}

// 点赞
func LikeVideo(userid, videoid int64) {
	like := db.Like{UserId: userid, VideoId: videoid}
	err := db.DB.Create(&like).Error
	if err != nil {
		panic(err)
	}
	err = db.DB.Model(&db.Video{ID: videoid}).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	if err != nil {
		panic(err)
	}
}

// 取消点赞
func NoLikeVideo(userid, videoid int64) {
	err := db.DB.Where("user_id = ? and video_id = ?", userid, videoid).Delete(&db.Like{}).Error
	if err != nil {
		panic(err)
	}
	err = db.DB.Model(&db.Video{ID: videoid}).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
	if err != nil {
		panic(err)
	}

}

// 列出所有喜欢的视频
func FindAllLikeVideo(id, myid int64) []db.Video {
	var results = make([]db.Video, 0)
	var likes = make([]db.Like, 0)
	var ids = make([]int64, 0)

	//找到所有点赞记录
	err := db.DB.Model(&db.Like{}).Where("user_id = ?", id).Find(&likes).Error
	if err != nil {
		panic(err)
	}
	if len(likes) == 0 {
		return nil
	}
	//找到所有视频id
	for i := 0; i < len(likes); i++ {
		ids = append(ids, likes[i].VideoId)
	}

	err = db.DB.Model(&db.Video{}).Find(&results, ids).Error

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(results); i++ {
		//查看我是否点赞
		var like db.Like
		if err = db.DB.Model(&db.Like{}).Where("video_id = ? and user_id = ?", results[i].ID, myid).Find(&like).Error; err != nil {
			panic(err)
		}
		if like.ID != 0 {
			results[i].IsFavorite = true
		}
	}
	return results
}

// GetCommentList 获取指定videoID的评论表
func GetCommentList(videoID int64) ([]db.Comment, error) {
	var commentList []db.Comment
	if err := db.DB.Model(&db.Comment{}).Where("video_id=?", videoID).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

// err := db.DB.Model(&db.Like{}).Where("user_id = ?", id).Find(&likes).Error

// PostComment 发布评论
func PostComment(comment db.Comment) error {
	if err := db.DB.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteComment 删除指定commentID的评论
func DeleteComment(commentID int64) error {
	if err := db.DB.Model(&db.Comment{}).Where("id = ?", commentID).Delete(&db.Comment{}).Error; err != nil {
		return err
	}
	return nil
}

// AddCommentCount add comment_count
func AddCommentCount(videoID int64) error {

	if err := db.DB.Model(&db.Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}
	return nil
}

// ReduceCommentCount reduce comment_count
func ReduceCommentCount(videoID int64) error {
	if err := db.DB.Model(&db.Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return err
	}
	return nil
}

// Get user info by ID
func GetAUserByID(userID int64, user *db.User) error {
	if err := db.DB.Model(&db.User{}).Where("id = ?", userID).Find(user).Error; err != nil {
		return err
	}
	return nil
}

// GetVideoAuthor get video author
func GetVideoAuthorID(videoID int64) (int64, error) {
	var video db.Video
	if err := db.DB.Model(&db.Video{}).Where("id = ?", videoID).Find(&video).Error; err != nil {
		return int64(video.ID), err
	}
	return video.AuthorId, nil
}
