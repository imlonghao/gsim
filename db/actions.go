package db

import (
	"github.com/imlonghao/gsim/types"
	"time"
)

func GetTasks() []types.Task {
	var tasks []types.Task
	DB.Where("next_scan_time < ?", time.Now()).Find(&tasks)
	return tasks
}

func GetWhitelists() []types.Whitelist {
	var whitelists []types.Whitelist
	DB.Table("whitelists").Find(&whitelists)
	return whitelists
}

func IfResultExisted(id string) bool {
	return DB.First(&types.Result{}, "id = ?", id).RecordNotFound()
}
