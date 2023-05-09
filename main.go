package main

import (
	"encoding/json"
	"fmt"
	"log" // @todo switch to nicer logger
	"net/http"
	"time"
)
import "github.com/gorilla/mux"
import "github.com/eduncan911/podcast"

type CloudCasts struct {
	Data []struct {
		Key           string        `json:"key"`
		URL           string        `json:"url"`
		Name          string        `json:"name"`
		Tags          []interface{} `json:"tags"`
		CreatedTime   time.Time     `json:"created_time"`
		UpdatedTime   time.Time     `json:"updated_time"`
		PlayCount     int           `json:"play_count"`
		FavoriteCount int           `json:"favorite_count"`
		CommentCount  int           `json:"comment_count"`
		ListenerCount int           `json:"listener_count"`
		RepostCount   int           `json:"repost_count"`
		Pictures      struct {
			Small         string `json:"small"`
			Thumbnail     string `json:"thumbnail"`
			MediumMobile  string `json:"medium_mobile"`
			Medium        string `json:"medium"`
			Large         string `json:"large"`
			Three20Wx320H string `json:"320wx320h"`
			ExtraLarge    string `json:"extra_large"`
			Six40Wx640H   string `json:"640wx640h"`
			Seven68Wx768H string `json:"768wx768h"`
			One024Wx1024H string `json:"1024wx1024h"`
		} `json:"pictures"`
		Slug string `json:"slug"`
		User struct {
			Key      string `json:"key"`
			URL      string `json:"url"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Pictures struct {
				Small         string `json:"small"`
				Thumbnail     string `json:"thumbnail"`
				MediumMobile  string `json:"medium_mobile"`
				Medium        string `json:"medium"`
				Large         string `json:"large"`
				Three20Wx320H string `json:"320wx320h"`
				ExtraLarge    string `json:"extra_large"`
				Six40Wx640H   string `json:"640wx640h"`
			} `json:"pictures"`
		} `json:"user"`
		AudioLength int `json:"audio_length"`
	} `json:"data"`
	Paging struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"paging"`
	Name string `json:"name"`
}

// UserHandler returns an RSS feed for a MixCloud username
func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	url := fmt.Sprintf("https://api.mixcloud.com/%s/cloudcasts/", username)
	log.Println(url)

	// @todo cache this response
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	CloudCastsData := CloudCasts{}

	//_, body := io.ReadAll(resp.Body)
	//log.Println(body)

	err = json.NewDecoder(resp.Body).Decode(&CloudCastsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	p := podcast.New(
		CloudCastsData.Name,
		fmt.Sprintf("https://www.mixcloud.com/%s", username),
		CloudCastsData.Name,
		&CloudCastsData.Data[0].CreatedTime,
		&CloudCastsData.Data[0].UpdatedTime,
	)

	// @todo loop through CloudCastsData.Data, adding items to the podcast
	for i := range CloudCastsData.Data {
		post := CloudCastsData.Data[i]

		item := podcast.Item{}

		item.GUID = post.Key
		item.Link = post.URL
		item.Title = post.Name
		item.Description = post.Name
		item.PubDate = &post.CreatedTime
		item.IImage = &podcast.IImage{
			HREF: post.Pictures.ExtraLarge,
		}
		item.Author = &podcast.Author{Name: username}
		item.AddDuration(int64(post.AudioLength))
		_, err := p.AddItem(item)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/rss+xml")
	w.Write([]byte(p.String()))
}

func main() {
	port := "8080" // @todo take this from an environment variable

	r := mux.NewRouter()
	// @todo serve a homepage
	// @todo serve more than just an RSS feed
	r.HandleFunc("/u/{username}", UserHandler)

	log.Printf("Listening on http://127.0.0.1:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
