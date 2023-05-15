package adrepo

import (
	"context"
	"fmt"
	"homework9/internal/ads"
	"homework9/internal/app"
	"sync"
	"time"
)

func New() app.AdRepository {
	return &adRepo{make(map[int64]ads.Ad, 0), 0, sync.Mutex{}}
}

type adRepo struct {
	ads   map[int64]ads.Ad
	idx   int64
	mutex sync.Mutex
}

func (r *adRepo) CreateAd(ctx context.Context, Title string, Text string, UserID int64) (ads.Ad, error) {
	r.mutex.Lock()
	newAd := ads.Ad{ID: r.idx, Title: Title, Text: Text, AuthorID: UserID, DateCreate: time.Now().UTC(), DateUpdate: time.Now().UTC()}
	r.ads[r.idx] = newAd
	r.idx++
	r.mutex.Unlock()
	return newAd, nil
}
func (r *adRepo) ChangeAdStatus(ctx context.Context, adID int64, Published bool) (ads.Ad, error) {
	ad, ok := r.ads[adID]
	if !ok {
		return ads.Ad{}, fmt.Errorf("can not find ad")
	}
	ad.Published = Published
	ad.DateUpdate = time.Now().UTC()
	r.ads[adID] = ad
	return ad, nil
}
func (r *adRepo) UpdateAd(ctx context.Context, adID int64, Title string, Text string) (ads.Ad, error) {
	ad, ok := r.ads[adID]
	if !ok {
		return ads.Ad{}, fmt.Errorf("can not find ad")
	}
	ad.Text = Text
	ad.Title = Title
	ad.DateUpdate = time.Now().UTC()
	r.ads[adID] = ad
	return ad, nil
}

func (r *adRepo) GetAd(ctx context.Context, index int64) (ads.Ad, error) {
	ad, ok := r.ads[index]
	if !ok {
		return ads.Ad{}, fmt.Errorf("can not find ad")
	}
	return ad, nil
}

func (r *adRepo) GetAdByTitle(ctx context.Context, Title string) (ads.Ad, error) {
	for i := range r.ads {
		if r.ads[i].Title == Title {
			return r.ads[i], nil
		}
	}
	return ads.Ad{}, fmt.Errorf("ad not found")
}

func (r *adRepo) GetAdsByUserID(ctx context.Context, ID int64) []ads.Ad {
	ads := make([]ads.Ad, 0)
	for _, ad := range ads {
		if ad.AuthorID == ID && ad.Published {
			ads = append(ads, ad)
		}
	}
	return ads
}

func (r *adRepo) GetAds(ctx context.Context) ([]ads.Ad, error) {
	ads := make([]ads.Ad, 0)
	for _, ad := range r.ads {
		if ad.Published {
			ads = append(ads, ad)
		}
	}
	return ads, nil
}

func (r *adRepo) GetAdsByTime(ctx context.Context, Time time.Time) []ads.Ad {
	ads := make([]ads.Ad, 0)
	for _, ad := range ads {
		if ad.DateCreate == Time && ad.Published {
			ads = append(ads, ad)
		}
	}
	return ads
}

func (r *adRepo) DeleteAd(ctx context.Context, adID int64) error {
	_, ok := r.ads[adID]
	if !ok {
		return fmt.Errorf("can not find ad")
	}
	delete(r.ads, adID)
	return nil
}
