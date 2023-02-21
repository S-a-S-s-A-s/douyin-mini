package db

import (
	"time"
)

func CreateData() {
	user := []*User{
		{Username: "sbbbbbbb", Password: "123456", FollowCount: 123, FollowerCount: 123456},
		{Username: "whale", Password: "123456", FollowCount: 123, FollowerCount: 123456},
	}
	if err := DB.Create(user).Error; err != nil {
		println(err)
	}

	videos := []*Video{
		{AuthorId: 1, AuthorName: "海贼王", CreateTime: time.Now(), PlayURL: "https://644-27-2.vod.tv.itc.cn/sohu/v1/Tmw30EIsWBdHge04g8sRDEeDM8cFP2wXPLEFD8vmPGXUyYbSoO27fSx.mp4?k=7DtxpY&p=j9lvzSwUqmkiqSoBoSri0S1AqSPCopkUhRYAtUxIgYeiwm12ZD6Sotxcyp0Gvm1mRDE&r=TUldziJCtpCmhWB3tSCGhWlvsmCUqpxWtWaizY&q=OpCUhW7IWhodRDvswmfCyY2sWh1Hfhd45G6tRhXtWGo2ZDvtRhWsWYb4wm4cWJvXY&nid=644", CoverURL: "https://photocdn.tv.sohu.com/img/20230208/frag_item_1675830939360_1.jpg", FavoriteCount: 1203, CommentCount: 1234, Title: "hhhh"},
		{AuthorId: 1, AuthorName: "海贼王", CreateTime: time.Now(), PlayURL: "https://644-41-1.vod.tv.itc.cn/sohu/v1/TmwA0EIsWYcLgTyO86cMe66BgebU8Efio6PIkI1bhXUyYbSoO27fSx.mp4?k=2IyCIZ&p=j9lvzSwUqmkiqSoBoSri0S1AqSPCopkUhRYAtUxIgYeiwm12ZD6Sotxcyp0Gvm1mRDE&r=TUldziJCtpCmhWB3tSCGhWlvsmCUqpxWtWaizY&q=OpCAhW7IWhodRDvswmfCyY2sWh1Hfhd45G6tRhXtWGo2ZDvtRhWsWYb4vm4cWJ1tY&nid=644", CoverURL: "https://photocdn.tv.sohu.com/img/20230208/frag_item_1675830939431_2.jpg", FavoriteCount: 123, CommentCount: 1234, Title: "hhhh"},
		{AuthorId: 1, AuthorName: "海贼王", CreateTime: time.Now(), PlayURL: "https://644-36-1.vod.tv.itc.cn/sohu/v1/TmP3qKIsWJNHDmo3POd2PLC7oMPGuT02fLE6yJy7PJXUyYbSoO27fSx.mp4?k=YPc9JY&p=j9lvzSwUqmkiqSoBoSri0S1AqSPCopkUhRYAtUxIgYeiwm12ZD6Sotxcyp0Gvm1mRDE&r=TUldziJCtpCmhWB3tSCGhWlvsmCUqpxWtWaizY&q=OpCGoEOyzSwWsSCAoKOL4HrIWh6s5G6XfFXsWBAHfBNS0F2OfBAOWh14fBoURDvsWZ&nid=644", CoverURL: "https://photocdn.tv.sohu.com/img/20230208/frag_item_1675842825344_5.jpg", FavoriteCount: 123, CommentCount: 1234, Title: "eeeee"},
	}
	if err := DB.Create(videos).Error; err != nil {
		println(err)
	}

	comments := []*Comment{
		{UserId: 1, UserName: "海贼王", VideoId: 1, Content: "test测试", CreateTime: time.Now().Format("2006-01-02 15:04:05")},
	}
	if err := DB.Create(comments).Error; err != nil {
		println(err)
	}
}
