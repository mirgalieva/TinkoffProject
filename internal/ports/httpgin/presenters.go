package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/ads"
	"homework9/internal/users"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type getTitle struct {
	Title string `json:"title"`
}

type adResponse struct {
	ID        int64  `json:"ad_id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type userResponse struct {
	ID       int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type getUserRequest struct {
	UserId int64 `json:"user_id"`
}
type deleteAdRequest struct {
	UserId int64 `json:"user_id"`
	AdId   int64 `json:"ad_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}
type getAdID struct {
	ID int64 `json:"id"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
		},
		"error": nil,
	}
}

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
		"error": nil,
	}
}

func UserSuccessDelete() *gin.H {
	return &gin.H{
		"data":  nil,
		"error": nil,
	}
}
func AdSuccessDelete() *gin.H {
	return &gin.H{
		"data":  nil,
		"error": nil,
	}
}
func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
func UserErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
func AdsSuccessResponse(ads []ads.Ad) *gin.H {
	ans := make([]adResponse, len(ads))
	for i, v := range ads {
		ans[i] = adResponse{
			ID:        v.ID,
			Title:     v.Title,
			Text:      v.Text,
			AuthorID:  v.AuthorID,
			Published: v.Published,
		}
	}
	return &gin.H{
		"data":  ans,
		"error": nil,
	}
}

type paramsAdRequest struct {
	Params map[string]any `json:"params"`
}
