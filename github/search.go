package github

import (
	"encoding/json"
	"fmt"
	"github.com/imlonghao/gsim/types"
	"strings"
)

func Search(keyword string) ([]types.Result, error) {
	token := randomToken()
	resp, err := httpGet(fmt.Sprintf("https://api.github.com/search/code?q=\"%s\"&sort=indexed", keyword), token)
	if err != nil {
		return nil, err
	}
	var searchResult types.GithubSearchCodeResultSchema
	err = json.Unmarshal(resp, &searchResult)
	if err != nil {
		return nil, err
	}
	var results []types.Result
	for _, result := range searchResult.Items {
		commit := strings.Split(result.URL, "=")[1]
		code, err := httpGet(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", result.Repository.FullName, commit, result.Path), "")
		if err != nil {
			code = []byte("")
		}
		results = append(results, types.Result{
			ID:       result.Sha,
			Status:   0,
			Username: result.Repository.Owner.Login,
			Repo:     result.Repository.Name,
			Path:     result.Path,
			Url:      result.HTMLURL,
			Code:     string(code),
		})
	}
	return results, nil
}
