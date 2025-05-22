package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	repoURL  = "http://mps.gov.com/service/rest/v1/components"
	repoName = "harbor"
	user     = "admin"
	password = "admin"
)

type Item struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Response struct {
	Items             []Item `json:"items"`
	ContinuationToken string `json:"continuationToken"`
}

var imageVersions = make(map[string]string)

func main() {
	fetchImages()

	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/query", queryHandler)

	fmt.Println("镜像管理服务启动，监听端口8080")
	http.ListenAndServe(":8080", nil)
}

func fetchImages() {
	continuationToken := ""
	for {
		var url string
		if continuationToken == "" {
			url = fmt.Sprintf("%s?repository=%s", repoURL, repoName)
		} else {
			url = fmt.Sprintf("%s?continuationToken=%s&repository=%s", repoURL, continuationToken, repoName)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.SetBasicAuth(user, password)
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var response Response
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}

		for _, item := range response.Items {
			imageName := item.Name
			version := item.Version

			if existing, ok := imageVersions[imageName]; ok {
				imageVersions[imageName] = existing + "," + version
			} else {
				imageVersions[imageName] = version
			}
		}

		continuationToken = response.ContinuationToken

		if continuationToken == "" {
			break
		}
	}
}

// 列出所有镜像及其版本
func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imageVersions)
}

// 模糊查询镜像
func queryHandler(w http.ResponseWriter, r *http.Request) {
	queryImage := r.URL.Query().Get("image")
	if queryImage == "" {
		http.Error(w, "参数 image 不能为空", http.StatusBadRequest)
		return
	}

	result := make(map[string]string)
	for imageName, versions := range imageVersions {
		if strings.Contains(imageName, queryImage) {
			result[imageName] = versions
		}
	}

	if len(result) == 0 {
		http.Error(w, "镜像不存在或无版本信息", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
