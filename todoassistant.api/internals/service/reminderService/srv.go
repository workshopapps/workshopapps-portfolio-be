package reminderService

import (
	"errors"
	"fmt"
	"log"
	"test-va/internals/Repository/reminderRepo"
	"test-va/internals/entity/notificationEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/service/notificationService"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
)

type ReminderSrv interface {
	SetReminder(data *taskEntity.CreateTaskReq) error
	SetReminderEvery30Min()
	SetReminderEvery5Min()
	SetDailyReminder(data *taskEntity.CreateTaskReq) error
	SetWeeklyReminder(data *taskEntity.CreateTaskReq) error
	SetBiWeeklyReminder(data *taskEntity.CreateTaskReq) error
	SetMonthlyReminder(data *taskEntity.CreateTaskReq) error
	SetYearlyReminder(data *taskEntity.CreateTaskReq) error
	ScheduleNotificationEverySixHours()
	ScheduleNotificationDaily()
}

type reminderSrv struct {
	cron *gocron.Scheduler
	// conn *sql.DB
	repo reminderRepo.ReminderRepository
	nSrv notificationService.NotificationSrv
}

func (r *reminderSrv) SetBiWeeklyReminder(data *taskEntity.CreateTaskReq) error {
	r.cron.RemoveByTag(data.TaskId)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}

	r.cron.Every(2).Weeks().Tag(data.TaskId).StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		// r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 14).Format(time.RFC3339)
		data.Status = "PENDING"
		log.Println(data)

		err = r.repo.CreateNewTask(data)
		if err != nil {
			return err
		}
		return nil
	})

	r.cron.StartAsync()
	log.Println("created new eventService.")
	return nil
}

func (r *reminderSrv) SetYearlyReminder(data *taskEntity.CreateTaskReq) error {
	r.cron.RemoveByTag(data.TaskId)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}
	r.cron.Every(12).Months().Tag(data.TaskId).StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		// r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(1, 0, 0).Format(time.RFC3339)
		data.Status = "PENDING"
		log.Println(data)

		err = r.repo.CreateNewTask(data)
		if err != nil {
			return err
		}
		return nil
	})

	r.cron.StartAsync()
	log.Println("created new eventService.")
	return nil
}

func (r *reminderSrv) SetMonthlyReminder(data *taskEntity.CreateTaskReq) error {
	r.cron.RemoveByTag(data.TaskId)
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}

	r.cron.Every(1).Months().Tag(data.TaskId).StartAt(dDate).Do(func() error {
		log.Println("setting status to expired")
		// r.repo.SetTaskToExpired(data.TaskId)
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 1).Format(time.RFC3339)
		data.Status = "PENDING"
		log.Println(data)

		err = r.repo.CreateNewTask(data)
		if err != nil {
			return err
		}
		r.cron.StartAsync()
		return nil
	})
	log.Println("created new eventService.")
	return nil
}

func (r *reminderSrv) SetWeeklyReminder(data *taskEntity.CreateTaskReq) error {
	r.cron.RemoveByTag(data.TaskId)
	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}

	r.cron.Every(1).Week().Tag(data.TaskId).StartAt(dDate).Do(func() error {
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 7).Format(time.RFC3339)
		data.Status = "PENDING"
		data.TaskId = uuid.New().String()

		err = r.repo.CreateNewTask(data)
		if err != nil {
			return err
		}

		return nil
	})

	log.Println("created new eventService.")
	r.cron.StartAsync()
	return nil
}

func (r *reminderSrv) SetDailyReminder(data *taskEntity.CreateTaskReq) error {
	r.cron.RemoveByTag(data.TaskId)
	dDate, err := time.Parse(time.RFC3339, data.EndTime)
	if err != nil {
		return err
	}

	if dDate.Before(time.Now().Local()) {
		return errors.New("invalid Time, try again")
	}

	r.cron.Every(1).Day().Tag(data.TaskId).StartAt(dDate).Do(func() error {
		endDate, err := time.Parse(time.RFC3339, data.EndTime)
		if err != nil {
			return err
		}

		data.StartTime = data.EndTime
		data.UpdatedAt = data.EndTime
		data.EndTime = endDate.AddDate(0, 0, 1).Format(time.RFC3339)
		data.CreatedAt = time.Now().Local().Format(time.RFC3339)
		data.Status = "PENDING"
		data.TaskId = uuid.New().String()

		// log.Println(data)

		err = r.repo.CreateNewTask(data)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("created new eventService.")

		//// send notification
		//task, err := r.noSrv.GetTaskFromUser(data.UserId)
		//if err != nil {
		//	return err
		//}
		//
		//fmt.Println("notification sent out")
		//r.noSrv.SendNotification(task.DeviceId,
		//	"Your Notification is about to expire",
		//	"your Task is due in 5 miutes",
		//	[]string{task.TaskId},
		//)

		return nil
	})
	r.cron.StartAsync()
	return nil
}

func (r *reminderSrv) SetReminder(data *taskEntity.CreateTaskReq) error {
	r.cron.RemoveByTag(data.TaskId)
	dueDate := data.EndTime
	taskId := data.TaskId

	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return err
	}

	r.cron.Every(1).Tag(data.TaskId).StartAt(dDate).Do(func() {
		//Send VA Notifications to Firebase
		vaTokens, vaId, username, err := r.nSrv.GetUserVaToken(data.UserId)
		if err != nil {
			fmt.Println("Error Getting VA Tokens", err)
		}
		if vaId != "" && username != "" {
			//Upload the Notifications to DB
			err := r.nSrv.CreateNotification(vaId, "Expired Task", time.Now().String(), fmt.Sprintf("%s has an expired task", username), notificationEntity.ExpiredColor, taskId)
			if err != nil {
				fmt.Println("Error Uploading Notification to DB", err)
			}
		}
		if len(vaTokens) < 1 {
			fmt.Println("User Has No VA, Or VA Has Not Registered For Notifications")
		}

		//Send User Notification With Firebase
		userTokens, username, err := r.nSrv.GetUserToken(data.UserId)
		if err != nil {
			fmt.Println("Error Getting User Tokens", err)
		}
		if username != "" {
			err := r.nSrv.CreateNotification(data.UserId, "Expired Task", time.Now().String(), fmt.Sprintf("%s has an expired task", username), notificationEntity.ExpiredColor, taskId)
			if err != nil {
				fmt.Println("Error Uploading Notification to DB", err)
			}
		}
		if err != nil {
			fmt.Println("Error Uploading Notification to DB", err)
		}
		if len(userTokens) < 1 {
			fmt.Println("User Has Not Registered For Notifications")
		}

		body := []notificationEntity.NotificationBody{
			{
				Content: "This Task Has Expired",
				Color:   notificationEntity.ExpiredColor,
				Time:    time.Now().Local().Format(time.RFC3339),
			},
		}

		allTokens := append(userTokens, vaTokens...)
		if len(allTokens) > 0 {
			err = r.nSrv.SendBatchNotifications(allTokens, "Expired", body, []interface{}{data})
			if err != nil {
				fmt.Println("Error Sending Notifications", err)
			}
		}
	})

	r.cron.LimitRunsTo(1)
	r.cron.StartAsync()
	return nil
}

func (r *reminderSrv) SetReminderEvery5Min() {
	tasks, err := r.repo.GetAllUsersPendingTasks()
	if err != nil {
		log.Println(err)
		return
	}

	for _, task := range tasks {
		//check if time until is <=30 minutes or 5 minutes
		yes := checkIfTimeElapsed5Minutes(task.EndTime)

		if yes {
			fmt.Println("notification sent out")
			r.nSrv.SendNotification(task.DeviceId,
				"Your Notification is about to expire",
				"your Task is due in 5 minutes",
				task.TaskId,
			)
			continue
		}
	}
}

func (r *reminderSrv) SetReminderEvery30Min() {
	tasks, err := r.repo.GetAllUsersPendingTasks()
	if err != nil {
		log.Println(err)
		return
	}

	for _, task := range tasks {
		//check if time until is <=30 minutes and greater than 5 minutes
		yes := checkIfTimeElapsed30Minutes(task.EndTime)

		if yes {
			fmt.Println("notification sent out")
			r.nSrv.SendNotification(task.DeviceId,
				"Your Notification is about to expire",
				"your Task is due in 30 miutes",
				task.TaskId,
			)
			continue
		}
	}
}

func checkIfTimeElapsed30Minutes(endTime string) bool {
	dueTime, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		log.Println(err)
		return false
	}

	if dueTime.Before(time.Now()) {
		return false
	}

	minutes := time.Until(dueTime).Minutes()
	if minutes <= float64(30) && minutes > float64(5) {
		return true
	} else {
		return false
	}
}

func checkIfTimeElapsed5Minutes(due string) bool {
	dueTime, err := time.Parse(time.RFC3339, due)
	if err != nil {
		log.Println(err)
		return false
	}

	if dueTime.Before(time.Now()) {
		return false
	}

	minutes := time.Until(dueTime).Minutes()
	if minutes <= float64(5) {
		return true
	} else {
		return false
	}
}

// Everyday By 12:00am you get Notifications For All Tasks That are Due That Day
func (r *reminderSrv) ScheduleNotificationDaily() {
	fmt.Println("Daily Notifications Setup")
	r.cron.Every(1).Days().At("00:00").Do(func() {
		tasks, err := r.nSrv.GetTasksToExpireToday()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Daily")

		if len(tasks) < 1 {
			fmt.Println("No Notifications to Send Just Yet", tasks)
			return
		}

		for k, v := range tasks {
			body := []notificationEntity.NotificationBody{
				{
					Content: fmt.Sprintf("You Have %v tasks due today", len(v)),
					Color:   notificationEntity.DueColor,
					Time:    time.Now().String(),
				},
			}

			err = r.nSrv.SendNotification(k, "Due Today", body, v[0])
			if err != nil {
				fmt.Println(err)
			}
		}
	})
}

// Your Pending Tasks are Checked On Six Hour Intervals to Get Tasks That Are Just About To Expire
func (r *reminderSrv) ScheduleNotificationEverySixHours() {
	fmt.Println("Six Hour Notifications Setup")
	r.cron.Every(6).Hours().Do(func() {
		tasks, err := r.nSrv.GetTasksToExpireInAFewHours()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Six Hourly")

		if len(tasks) < 1 {
			fmt.Println("No Notifications to Send Just Yet", tasks)
			return
		}

		for k, v := range tasks {
			body := []notificationEntity.NotificationBody{
				{
					Content: fmt.Sprintf("You Have %v tasks due in a few hours", len(v)),
					Color:   notificationEntity.DueColor,
					Time:    time.Now().String(),
				},
			}

			err = r.nSrv.SendNotification(k, "Due Shortly", body, v[0])
			if err != nil {
				fmt.Println(err)
			}
		}
	})
}

// func getExpiredTasks(conn *sql.DB) ([]notificationEntity.GetExpiredTasksWithDeviceId, error) {
// 	stmt := fmt.Sprint(`
// 		SELECT task_id, Tasks.user_id, title ,description, end_time, device_id, va_id
// 		FROM Tasks
// 		INNER JOIN Notification_Tokens ON Tasks.user_id = Notification_Tokens.user_id
// 		WHERE Tasks.status = 'EXPIRED';
// 	`)

// 	var tasks []notificationEntity.GetExpiredTasksWithDeviceId

// 	query, err := conn.Query(stmt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for query.Next() {
// 		var task notificationEntity.GetExpiredTasksWithDeviceId
// 		var deviceId string
// 		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &deviceId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tasks = append(tasks, task)
// 	}
// 	return tasks, nil
// }

func NewReminderSrv(s *gocron.Scheduler, reminderRepo reminderRepo.ReminderRepository, nSrv notificationService.NotificationSrv) ReminderSrv {
	return &reminderSrv{cron: s, repo: reminderRepo, nSrv: nSrv}
}
