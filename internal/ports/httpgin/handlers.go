package httpgin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"homework9/internal/app"
	"net/http"
	"strconv"
	"time"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
		}
		ad, err := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			if errors.Is(err, app.ErrWrongUser) {
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
				return
			}
			if errors.Is(err, app.ErrValidationFail) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		ad, err := a.ChangeAdStatus(c, adID, reqBody.UserID, reqBody.Published)
		if err != nil {
			if errors.Is(err, app.ErrWrongUser) {
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
				return
			}
			if errors.Is(err, app.ErrValidationFail) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if err != nil {
			if errors.Is(err, app.ErrWrongUser) {
				c.JSON(http.StatusForbidden, AdErrorResponse(err))
				return
			}
			if errors.Is(err, app.ErrValidationFail) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
		}
		user, err := a.CreateUser(c, reqBody.Nickname, reqBody.Email)
		if err != nil {
			if errors.Is(err, app.ErrValidationFail) {
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(&user))
	}
}

func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody getUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		err = a.DeleteUser(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessDelete())
	}
}

func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody getUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		user, err := a.GetUser(c, userID)
		if err != nil {
			if errors.Is(err, app.ErrValidationFail) {
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(&user))
	}
}

func getAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody getAdID
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}
		ad, err := a.GetAd(c, adID)
		if err != nil {
			if errors.Is(err, app.ErrValidationFail) {
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func getAdByTitle(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody getTitle
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		title := c.Param("title")
		ad, err := a.GetAdByTitle(c, title)
		if err != nil {
			if errors.Is(err, app.ErrValidationFail) {
				c.JSON(http.StatusBadRequest, UserErrorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, UserErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func getAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, err := a.GetAds(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdsSuccessResponse(ad))
	}
}
func getAdsFilter(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody paramsAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		for k, v := range reqBody.Params {
			switch k {
			case "author_id":
				val, ok := v.(float64)
				if !ok {
					c.JSON(http.StatusBadRequest, AdErrorResponse(fmt.Errorf("bad parameters")))
					return
				}
				reqBody.Params[k] = int64(val)
			case "date_create":
				val, ok := v.(time.Time)
				if !ok {
					c.JSON(http.StatusBadRequest, AdErrorResponse(fmt.Errorf("bad parameters")))
					return
				}
				reqBody.Params[k] = val.UTC()
			case "published":
				val, ok := v.(bool)
				if !ok {
					c.JSON(http.StatusBadRequest, AdErrorResponse(fmt.Errorf("bad parameters")))
					return
				}
				reqBody.Params[k] = val
			default:
				c.JSON(http.StatusBadRequest, AdErrorResponse(fmt.Errorf("bad parameters")))
				return
			}
		}
		ad, err := a.GetAdsPrams(c, reqBody.Params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdsSuccessResponse(ad))
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
		}
		_, err := a.GetUser(c, int64(reqBody.UserId))
		if err != nil {
			c.JSON(http.StatusUnauthorized, AdErrorResponse(err))
			return
		}
		adIDs := c.Param("ad_id")
		if adIDs == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid ad_id"})
			return
		}
		adID, err := strconv.Atoi(adIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		err = a.DeleteAd(c, int64(adID), reqBody.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessDelete())
	}
}
