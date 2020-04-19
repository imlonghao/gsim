package db

import (
	"github.com/imlonghao/gsim/types"
	"time"
)

func GetTasks() ([]types.Task, error) {
	var tasks []types.Task
	err := DB.Where("next_scan_time < ?", time.Now()).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetWhitelists() ([]types.Whitelist, error) {
	var whitelists []types.Whitelist
	err := DB.Table("whitelists").Find(&whitelists).Error
	if err != nil {
		return nil, err
	}
	return whitelists, nil
}

func IfResultExisted(id string) bool {
	return DB.First(&types.Result{}, "id = ?", id).RecordNotFound()
}
