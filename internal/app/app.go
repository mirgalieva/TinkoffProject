package app

import (
	"context"
	"fmt"
	"github.com/mirgalieva/valid"
	"github.com/pkg/errors"
	"homework9/internal/ads"
	"homework9/internal/users"
)

var ErrWrongUser = errors.New("user has no rights")
var ErrValidationFail = errors.New("ad is not valid")

type App interface {
	CreateAd(ctx context.Context, Title string, Text string, UserID int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, UserID int64, Published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, UserID int64, Title string, Text string) (ads.Ad, error)
	CreateUser(ctx context.Context, Nickname string, Email string) (users.User, error)
	DeleteUser(ctx context.Context, ID int64) error
	GetUser(ctx context.Context, ID int64) (users.User, error)
	GetAd(ctx context.Context, index int64) (ads.Ad, error)
	GetAdByTitle(ctx context.Context, Title string) (ads.Ad, error)
	GetUsers(ctx context.Context) map[int64]users.User
	GetAds(ctx context.Context) ([]ads.Ad, error)
	GetAdsPrams(ctx context.Context, param map[string]interface{}) ([]ads.Ad, error)
	DeleteAd(ctx context.Context, adID int64, userID int64) error
}

type AdRepository interface {
	CreateAd(ctx context.Context, Title string, Text string, UserID int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, Published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, Title string, Text string) (ads.Ad, error)
	GetAd(ctx context.Context, index int64) (ads.Ad, error)
	GetAdByTitle(ctx context.Context, Title string) (ads.Ad, error)
	GetAds(ctx context.Context) ([]ads.Ad, error)
	DeleteAd(ctx context.Context, adID int64) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, Nickname string, Email string) (users.User, error)
	DeleteUser(ctx context.Context, ID int64) error
	GetUser(ctx context.Context, ID int64) (users.User, error)
	GetUsers(ctx context.Context) map[int64]users.User
}

func NewApp(adRepo AdRepository, userRepo UserRepository) App {
	return &app{adRepo: adRepo, userRepo: userRepo}
}

type app struct {
	adRepo   AdRepository
	userRepo UserRepository
}

func (a *app) DeleteAd(ctx context.Context, adID int64, userID int64) error {
	ad, err := a.GetAd(ctx, adID)
	if err != nil {
		return err
	}
	if ad.AuthorID != userID {
		return ErrWrongUser
	}
	err = a.adRepo.DeleteAd(ctx, adID)
	if err != nil {
		return err
	}
	return nil
}

func (a *app) GetAdsPrams(ctx context.Context, param map[string]any) ([]ads.Ad, error) {
	ans, err := a.adRepo.GetAds(ctx)
	if err != nil {
		return nil, err
	}
	var adArr = make([]ads.Ad, 0)
LP:
	for _, an := range ans {
		for k, v := range param {
			if k == "published" && an.Published != v {
				continue LP
			}
			if k == "author_id" && an.AuthorID != v {
				continue LP
			}
			if k == "date_create" && an.DateCreate != v {
				continue LP
			}
		}
		adArr = append(adArr, an)
	}
	return adArr, nil
}

func (a *app) GetAds(ctx context.Context) ([]ads.Ad, error) {
	Ads, err := a.adRepo.GetAds(ctx)
	if err != nil {
		return make([]ads.Ad, 0), err
	}
	return Ads, nil
}

func (a *app) GetUsers(ctx context.Context) map[int64]users.User {
	return a.userRepo.GetUsers(ctx)
}

type ValidTitleAndText struct {
	Title string `validate:"min:1,max:100"`
	Text  string `validate:"min:1,max:500"`
}

type ValidNicknameAndEmail struct {
	Nickname string `validate:"min:1,max:100"`
	Email    string `validate:"min:1,max:100"`
}

func (a *app) CreateAd(ctx context.Context, Title string, Text string, UserID int64) (ads.Ad, error) {
	valid := ValidTitleAndText{Title, Text}
	err := homework.Validate(valid)
	if err != nil {
		return ads.Ad{}, ErrValidationFail
	}
	ad, err := a.adRepo.CreateAd(ctx, Title, Text, UserID)
	if err != nil {
		return ad, err
	}
	return ad, nil
}
func (a *app) ChangeAdStatus(ctx context.Context, adID int64, UserID int64, Published bool) (ads.Ad, error) {
	ad, err := a.adRepo.GetAd(ctx, adID)
	if err != nil {
		return ads.Ad{}, fmt.Errorf("invalid adId")
	}

	if ad.AuthorID != UserID {
		return ads.Ad{}, ErrWrongUser
	}

	updatedAd, err := a.adRepo.ChangeAdStatus(ctx, adID, Published)
	if err != nil {
		return ads.Ad{}, err
	}
	return updatedAd, nil
}

func (a *app) UpdateAd(ctx context.Context, adID int64, UserID int64, Title string, Text string) (ads.Ad, error) {
	ad, err := a.adRepo.GetAd(ctx, adID)
	if err != nil {
		return ads.Ad{}, fmt.Errorf("invalid adId")
	}
	if ad.AuthorID != UserID {
		return ads.Ad{}, ErrWrongUser
	}
	valid := ValidTitleAndText{Title, Text}
	err = homework.Validate(valid)
	if err != nil {
		return ads.Ad{}, ErrValidationFail
	}
	updatedAd, err := a.adRepo.UpdateAd(ctx, adID, Title, Text)
	if err != nil {
		return ads.Ad{}, err
	}
	return updatedAd, nil
}

func (a *app) CreateUser(ctx context.Context, Nickname string, Email string) (users.User, error) {
	err := homework.Validate(ValidNicknameAndEmail{Nickname: Nickname, Email: Email})
	if err != nil {
		return users.User{}, ErrValidationFail
	}
	user, err := a.userRepo.CreateUser(ctx, Nickname, Email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *app) DeleteUser(ctx context.Context, ID int64) error {
	err := a.userRepo.DeleteUser(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *app) GetUser(ctx context.Context, ID int64) (users.User, error) {
	user, err := a.userRepo.GetUser(ctx, ID)
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (a *app) GetAd(ctx context.Context, ID int64) (ads.Ad, error) {
	ad, err := a.adRepo.GetAd(ctx, ID)
	if err != nil {
		return ads.Ad{}, err
	}
	return ad, nil
}

func (a *app) GetAdByTitle(ctx context.Context, Title string) (ads.Ad, error) {
	err := homework.Validate(ValidTitleAndText{Title, "1"})
	if err != nil {
		return ads.Ad{}, ErrValidationFail
	}
	ad, err := a.adRepo.GetAdByTitle(ctx, Title)
	if err != nil {
		return ads.Ad{}, err
	}
	return ad, nil
}
