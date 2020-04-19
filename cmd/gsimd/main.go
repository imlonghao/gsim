package main

import (
	"github.com/imlonghao/gsim/db"
	"github.com/imlonghao/gsim/github"
	"github.com/imlonghao/gsim/log"
	"github.com/imlonghao/gsim/sentry"
	"github.com/imlonghao/gsim/types"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
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
		log.Logger.Info("task finished", zap.Uint("id", task.ID))
		task.NextScanTime = time.Now().Add(time.Duration(task.Interval) * time.Second)
		err := db.DB.Table("tasks").Save(&task).Error
		if err != nil {
			sentry.SENTRY.CaptureException(err)
		}
		taskIsRunning[task.ID] = false
	}()
	results, err := github.Search(task.Rule)
	if err != nil {
		sentry.SENTRY.CaptureException(err)
		log.Logger.Error("task failed", zap.Uint("id", task.ID), zap.Error(err))
		return
	}
	whitelists, err := db.GetWhitelists()
	if err != nil {
		sentry.SENTRY.CaptureException(err)
		return
	}
	for _, result := range results {
		matchWhitelist := isMatchWhitelist(result, whitelists)
		existed := db.IfResultExisted(result.ID)
		if matchWhitelist {
			if existed {
				db.DB.Table("results").Delete(&result)
			}
			log.Logger.Info("whitelist matched",
				zap.Uint("id", task.ID),
				zap.String("username", result.Username),
				zap.String("repo", result.Repo),
				zap.String("path", result.Path))
			continue
		}
		result.TaskID = task.ID
		result.Task = task
		if existed {
			log.Logger.Info("result added",
				zap.Uint("id", task.ID),
				zap.String("username", result.Username),
				zap.String("repo", result.Repo),
				zap.String("path", result.Path))
			db.DB.Table("results").Create(&result)
		}
	}
}

func main() {
	taskIsRunning = make(map[uint]bool)
	c := cron.New()
	c.AddFunc("* * * * *", func() {
		log.Logger.Info("Cron job starting", zap.String("module", "updater"))
		tasks, err := db.GetTasks()
		if err != nil {
			sentry.SENTRY.CaptureException(err)
			return
		}
		for _, task := range tasks {
			if taskIsRunning[task.ID] {
				log.Logger.Warn("task is running", zap.Uint("id", task.ID))
				continue
			}
			log.Logger.Info("task is starting", zap.Uint("id", task.ID))
			taskIsRunning[task.ID] = true
			go worker(task)
		}
	})
	c.AddFunc("* * * * *", func() {
		log.Logger.Info("Cron job starting", zap.String("module", "result cleaner"))
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
