package adrepo

import (
	"homework9/internal/ads"
	"homework9/internal/app"
	"homework9/internal/users"
	"strings"
	"sync"
)

const (
	published    string = "published"
	userId       string = "user_id"
	dateCreating string = "date_creating"
	dateFormat   string = "2006-01-02"
)

type repositoryMap struct {
	dictAds        map[int64]ads.Ad
	dictUsers      map[int64]users.User
	dictAdsByTitle map[string][]ads.Ad

	counterAds   int64
	counterUsers int64

	mu sync.Mutex
}

func (repo *repositoryMap) GetAdById(id int64) (ads.Ad, error) {
	ad, ok := repo.dictAds[id]
	if !ok {
		return ad, app.IncorrectAdId
	}
	return ad, nil
}

func (repo *repositoryMap) AddAd(ad *ads.Ad) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.dictAds[ad.ID] = *ad
	repo.dictAdsByTitle[ad.Title] = append(repo.dictAdsByTitle[ad.Title], *ad)
	repo.counterAds++
}

func New() app.Repository {
	return &repositoryMap{dictAds: make(map[int64]ads.Ad), dictUsers: make(map[int64]users.User), dictAdsByTitle: make(map[string][]ads.Ad), counterAds: 0, counterUsers: 0}
}

func (repo *repositoryMap) GetAdsPrimaryKey() int64 {
	return repo.counterAds
}

func (repo *repositoryMap) GetUsersPrimaryKey() int64 {
	return repo.counterUsers
}

func (repo *repositoryMap) ChangeAd(ad *ads.Ad) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, ok := repo.dictAds[ad.ID]
	if ok {
		repo.dictAds[ad.ID] = *ad
		for idx := range repo.dictAdsByTitle[ad.Title] {
			if repo.dictAdsByTitle[ad.Title][idx].ID == ad.ID {
				repo.dictAdsByTitle[ad.Title][idx] = *ad
				break
			}
		}
	}
	return ok
}

func (repo *repositoryMap) GetAdsByTitle(pattern string) []ads.Ad {
	var listByTitle []ads.Ad
	for title, list := range repo.dictAdsByTitle {
		if strings.HasPrefix(title, pattern) {
			listByTitle = append(listByTitle, list...)
		}
	}

	return listByTitle
}

func (repo *repositoryMap) GetAds(filters map[string]any) []ads.Ad {
	var list []ads.Ad
	var selectedAds = repo.dictAds
	if len(filters) == 0 {
		selectedAds = SelectByPublished(selectedAds, true)
	} else {
		if filter, ok := filters[published]; ok {
			selectedAds = SelectByPublished(selectedAds, filter)
		}

		if filter, ok := filters[userId]; ok {
			selectedAds = SelectByUserId(selectedAds, filter)
		}

		if filter, ok := filters[dateCreating]; ok {
			selectedAds = SelectByDateCreating(selectedAds, filter)
		}
	}

	for _, val := range selectedAds {
		list = append(list, val)
	}
	return list
}

func SelectByPublished(dict map[int64]ads.Ad, published any) map[int64]ads.Ad {
	repoWithFilter := make(map[int64]ads.Ad)
	for id := range dict {
		if dict[id].Published == published {
			repoWithFilter[id] = dict[id]
		}
	}
	return repoWithFilter
}

func SelectByUserId(dict map[int64]ads.Ad, userId any) map[int64]ads.Ad {
	repoWithFilter := make(map[int64]ads.Ad)
	for id := range dict {
		if dict[id].AuthorID == userId {
			repoWithFilter[id] = dict[id]
		}
	}
	return repoWithFilter
}

func SelectByDateCreating(dict map[int64]ads.Ad, dateCreating any) map[int64]ads.Ad {
	repoWithFilter := make(map[int64]ads.Ad)
	for id := range dict {
		elem := dict[id].DateCreating.Format(dateFormat)
		filter := dateCreating.(string)[:10]
		if elem == filter {
			repoWithFilter[id] = dict[id]
		}
	}
	return repoWithFilter
}

func (repo *repositoryMap) GetUserById(id int64) (users.User, error) {
	user, ok := repo.dictUsers[id]
	if !ok {
		return user, app.IncorrectUserId
	}
	return user, nil
}

func (repo *repositoryMap) AddUser(user *users.User) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.dictUsers[user.ID] = *user
	repo.counterUsers++
}

func (repo *repositoryMap) ChangeUser(user *users.User) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, ok := repo.dictUsers[user.ID]
	if ok {
		repo.dictUsers[user.ID] = *user
	}
	return ok
}

func (repo *repositoryMap) DeleteUser(userId int64) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.dictUsers, userId)
}

func (repo *repositoryMap) DeleteAd(adId int64) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	title := repo.dictAds[adId].Title

	for idx, ad := range repo.dictAdsByTitle[title] {
		if ad.ID == adId {
			repo.dictAdsByTitle[title][idx] = repo.dictAdsByTitle[title][len(repo.dictAdsByTitle[title])-1]
			repo.dictAdsByTitle[title] = repo.dictAdsByTitle[title][:len(repo.dictAdsByTitle[title])-1]
			if len(repo.dictAdsByTitle[title]) == 0 {
				delete(repo.dictAdsByTitle, title)
			}
			break
		}
	}

	delete(repo.dictAds, adId)
}
