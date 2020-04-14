package github

import (
	"github.com/imlonghao/gsim/db"
	"github.com/imlonghao/gsim/types"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

func randomToken() string {
	var token types.Token
	// rand() in MySQL
	// random() in PostgreSQL / Sqlite3
	db.DB.Order(gorm.Expr("rand()")).Limit(1).Find(&token)
	return token.Token
}

func httpGet(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
