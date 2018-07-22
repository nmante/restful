package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	restclient "github.com/nmante/restful/restclient"
)

const (
	postsUrl = "https://jsonplaceholder.typicode.com/posts"
)

var (
	postsHeaders http.Header
	restClient   = restclient.New()
)

// Post represents a post from jsonplaceholder
type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// PostsResponse is a list of posts
type PostsResponse struct {
	Total   int    `json:"total"`
	Results []Post `json:"results"`
}

// PostHandler returns a post
func GetPostHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	url := fmt.Sprintf("%s/%s", postsUrl, params.ByName("id"))
	response, err := restClient.Get(url, r, nil)

	if err != nil {
		restClient.WriteErrorResponse(w, "server_error", "Server error occured")
		return
	}

	var post Post
	restClient.WriteJSONResponse(w, response, post)
}

// PostsHandler returns a list of posts
func GetPostsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	url := postsUrl

	if qs := r.URL.RawQuery; len(qs) > 0 {
		url = fmt.Sprintf("%s?%s", postsUrl, r.URL.RawQuery)
	}

	response, err := restClient.Get(url, r, r.Header)

	if err != nil {
		restClient.WriteErrorResponse(w, "server_error", "Server error occured")
		return
	}

	var posts []Post
	restClient.WriteJSONResponse(w, response, posts)
}

// CreatePostHandler creates a post
func CreatePostHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	response, err := restClient.Post(postsUrl, r, r.Header)

	if err != nil {
		restClient.WriteErrorResponse(w, "server_error", "Server error occured")
		return
	}

	var post Post
	restClient.WriteJSONResponse(w, response, post)
}

// UpdatePostHandler creates a post
func UpdatePostHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	url := fmt.Sprintf("%s/%s", postsUrl, params.ByName("id"))
	response, err := restClient.Put(url, r, r.Header)

	if err != nil {
		restClient.WriteErrorResponse(w, "server_error", "Server error occured")
		return
	}

	var post Post
	restClient.WriteJSONResponse(w, response, post)
}

// UpdatePostHandler creates a post
func PatchPostHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	url := fmt.Sprintf("%s/%s", postsUrl, params.ByName("id"))
	response, err := restClient.Patch(url, r, r.Header)

	if err != nil {
		restClient.WriteErrorResponse(w, "server_error", "Server error occured")
		return
	}

	var post Post
	restClient.WriteJSONResponse(w, response, post)
}

// DeletePostHandler creates a post
func DeletePostHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	url := fmt.Sprintf("%s/%s", postsUrl, params.ByName("id"))
	response, err := restClient.Delete(url, r, r.Header)

	if err != nil {
		restClient.WriteErrorResponse(w, "server_error", "Server error occured")
		return
	}

	var post Post
	restClient.WriteJSONResponse(w, response, post)
}
