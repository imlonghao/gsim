package main

import (
	"fmt"
	"github.com/imlonghao/gsim/db"
	"github.com/imlonghao/gsim/github"
	"github.com/imlonghao/gsim/types"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

var taskIsRunning map[uint]bool

func isMatchWhitelist(result types.Result, whitelists []types.Whitelist) bool {
	for _, whitelist := range whitelists {
		switch whitelist.Type {
		case 1:
			if result.Username == whitelist.Content {
				return true
			}
		case 2:
			if strings.Contains(strings.ToLower(result.Repo), whitelist.Content) {
				return true
			}
		case 3:
			if strings.Contains(strings.ToLower(result.Code), whitelist.Content) {
				return true
			}
		case 4:
			if strings.Contains(strings.ToLower(result.Path), whitelist.Content) {
				return true
			}
		}
	}
	return false
}

func worker(task types.Task) {
	defer func() {
		fmt.Printf("Task %d is finishing\n", task.ID)
		task.NextScanTime = time.Now().Add(time.Duration(task.Interval) * time.Second)
		db.DB.Table("tasks").Save(&task)
		taskIsRunning[task.ID] = false
	}()
	results, err := github.Search(task.Rule)
	if err != nil {
		fmt.Printf("Task %d fail, %v\n", task.ID, err)
		return
	}
	whitelists := db.GetWhitelists()
	for _, result := range results {
		matchWhitelist := isMatchWhitelist(result, whitelists)
		existed := db.IfResultExisted(result.ID)
		if matchWhitelist {
			if existed {
				db.DB.Table("results").Delete(&result)
			}
			fmt.Printf("Task %d with %s/%s %s matching the whitelist\n", task.ID, result.Username, result.Repo, result.Path)
			continue
		}
		result.TaskID = task.ID
		result.Task = task
		if existed {
			fmt.Printf("Task %d with %s/%s %s adding to db\n", task.ID, result.Username, result.Repo, result.Path)
			db.DB.Table("results").Create(&result)
		}
	}
}

func main() {
	taskIsRunning = make(map[uint]bool)
	c := cron.New()
	c.AddFunc("* * * * *", func() {
		fmt.Printf("Running cron job\n")
		tasks := db.GetTasks()
		for _, task := range tasks {
			if taskIsRunning[task.ID] {
				fmt.Printf("Task %d is running, skipping\n", task.ID)
				continue
			}
			fmt.Printf("Task %d is starting\n", task.ID)
			taskIsRunning[task.ID] = true
			go worker(task)
		}
	})
	c.AddFunc("* * * * *", func() {
		fmt.Printf("Running cron job results cleaner\n")
		var results []types.Result
		var whitelists []types.Whitelist
		db.DB.Table("results").Where("status = 0").Find(&results)
		db.DB.Table("whitelists").Find(&whitelists)
		for _, result := range results {
			if isMatchWhitelist(result, whitelists) {
				db.DB.Table("results").Delete(&result)
			}
		}
	})
	c.Start()
	select {}
}
