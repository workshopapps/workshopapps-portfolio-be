package taskService

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/notificationEntity"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/entity/vaEntity"
	"test-va/internals/service/loggerService"
	"test-va/internals/service/notificationService"
	"test-va/internals/service/reminderService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
	PersistTask(req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *ResponseEntity.ServiceError)
	GetPendingTasks(userId string) ([]*taskEntity.GetPendingTasksRes, *ResponseEntity.ServiceError)
	SearchTask(req *taskEntity.SearchTitleParams) ([]*taskEntity.SearchTaskRes, *ResponseEntity.ServiceError)
	GetListOfExpiredTasks() ([]*taskEntity.GetAllExpiredRes, *ResponseEntity.ServiceError)
	GetListOfPendingTasks() ([]*taskEntity.GetAllPendingRes, *ResponseEntity.ServiceError)
	DeleteTaskByID(taskId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	GetAllTask(userId string) ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError)
	GetTaskByID(taskId string) (*taskEntity.GetTasksByIdRes, *ResponseEntity.ServiceError)
	DeleteAllTask(userId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	UpdateTaskStatusByID(taskId string, req *taskEntity.UpdateTaskStatus) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	EditTaskByID(taskId string, req *taskEntity.EditTaskReq) (*taskEntity.EditTaskRes, *ResponseEntity.ServiceError)

	GetVADetails(userId string) (string, *ResponseEntity.ServiceError)
	AssignTaskToVA(req *taskEntity.AssignReq) *ResponseEntity.ServiceError
	GetTaskAssignedToVA(vaId string) ([]*vaEntity.VATask, *ResponseEntity.ServiceError)
	GetAllTaskForVA() ([]*vaEntity.VATaskAll, *ResponseEntity.ServiceError)

	//comments
	PersistComment(req *taskEntity.CreateCommentReq) (*taskEntity.CreateCommentRes, *ResponseEntity.ServiceError)
	GetAllComments(taskId string) ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError)
	DeleteCommentByID(commentId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	GetComments() ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError)
}

type taskSrv struct {
	repo          taskRepo.TaskRepository
	timeSrv       timeSrv.TimeService
	validationSrv validationService.ValidationSrv
	logger        loggerService.LogSrv
	remindSrv     reminderService.ReminderSrv
	nSrv          notificationService.NotificationSrv
}

func NewTaskSrv(repo taskRepo.TaskRepository, timeSrv timeSrv.TimeService,
	srv validationService.ValidationSrv, logSrv loggerService.LogSrv,
	reminderSrv reminderService.ReminderSrv,
	notificationSrv notificationService.NotificationSrv) TaskService {
	return &taskSrv{repo: repo, timeSrv: timeSrv, validationSrv: srv,
		logger: logSrv, remindSrv: reminderSrv, nSrv: notificationSrv}
}

// Get Tasks Assigned To VA godoc
// @Summary	Get all tasks assigned to a VA
// @Description	Tasks assigned to VA route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	vaId	path	string	true	"VA Id"
// @Success	200  {object}  []vaEntity.VATask
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user/assigned-tasks/{vaId} [get]

// Get Tasks Assigned To VA godoc
// @Summary	Get all tasks assigned to a VA
// @Description	Tasks assigned to VA route
// @Tags	VA - Tasks
// @Accept	json
// @Produce	json
// @Success	200  {object}  []vaEntity.VATask
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/all/va [get]
func (t *taskSrv) GetTaskAssignedToVA(vaId string) ([]*vaEntity.VATask, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	va, err := t.repo.GetAllTaskAssignedToVA(ctx, vaId)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return va, nil
}

// Get All Tasks For VA godoc
// @Summary	Get all tasks for a VA
// @Description	All tasks for VA route
// @Tags	VA - Tasks
// @Accept	json
// @Produce	json
// @Success	200  {object}  []vaEntity.VATaskAll
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/all [get]
func (t *taskSrv) GetAllTaskForVA() ([]*vaEntity.VATaskAll, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	va, err := t.repo.GetAllTaskForVA(ctx)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return va, nil
}

// Assign task to VA godoc
// @Summary	A user assign task to a VA
// @Description	Assign task to VA route
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	taskId	path	string	true	"Task Id"
// @Param	request	body taskEntity.AssignReq	true	"Task Details"
// @Success	200  {string}  string    "ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/assign/{taskId} [post]
func (t *taskSrv) AssignTaskToVA(req *taskEntity.AssignReq) *ResponseEntity.ServiceError {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	vaID, serviceError := t.GetVADetails(req.UserId)
	if serviceError != nil {
		log.Println(" error here", serviceError)
		return serviceError
	}

	err := t.repo.AssignTaskToVa(ctx, vaID, req.TaskId)
	if err != nil {
		log.Println(" error here 2", err)
		return ResponseEntity.NewInternalServiceError(err)
	}

	// t.nSrv.SendNotificationToVA(req.UserId, "Task Assigned", fmt.Sprintf("%s Just Assigned a Task to You", req.UserId), data)

	return nil
}

func (t *taskSrv) GetVADetails(userId string) (string, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	vaId, err := t.repo.GetVADetails(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ResponseEntity.NewInternalServiceError("No VA assigned yet")
		}
		return "", ResponseEntity.NewInternalServiceError(err)
	}

	return vaId, nil
}

// Get Pending Tasks godoc
// @Summary	Get pending tasks of a particular user
// @Description	Get all pending task
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	userId	path	string	true	"User Id"
// @Success	200  {object}  []taskEntity.GetPendingTasksRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task/pending/{userId} [get]
func (t *taskSrv) GetPendingTasks(userId string) ([]*taskEntity.GetPendingTasksRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	tasks, err := t.repo.GetPendingTasks(userId, ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return tasks, nil
}

// Create Task godoc
// @Summary	Create a task
// @Description	Create task
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	request	body	taskEntity.CreateTaskReq	true "Tasks Details"
// @Success	200  {object}  taskEntity.CreateTaskRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task [post]
func (t *taskSrv) PersistTask(req *taskEntity.CreateTaskReq) (*taskEntity.CreateTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	// implement validation for struct

	err := t.validationSrv.Validate(req)
	if err != nil {
		// log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}

	//set time
	req.CreatedAt = t.timeSrv.CurrentTimeString()
	req.UpdatedAt = t.timeSrv.CurrentTimeString()
	//set id
	req.TaskId = uuid.New().String()
	req.Status = "PENDING"

	//set start time and endtime
	if req.StartTime == "" {
		req.StartTime = t.timeSrv.CurrentTimeString()
	}

	if req.EndTime == "" {
		req.EndTime = t.timeSrv.CalcEndTimeString()
	}

	//check if timeDueDate and StartDate is valid
	req.EndTime, err = t.timeSrv.CheckFor339Format(req.EndTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad Start-Time Input", err)
	}

	req.StartTime, err = t.timeSrv.CheckFor339Format(req.StartTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad End-Time Input", err)
	}

	var schedule time.Time
	if req.ScheduledDate != "" {
		schedule, err = t.timeSrv.Parse(req.ScheduledDate)
		if err != nil {
			return nil, ResponseEntity.NewCustomServiceError("Error when parsing scheduled date", err)
		}

		// log.Println(err)
		ok := t.timeSrv.ScheduleTimeAfter(schedule)
		if err != nil || !ok {
			return nil, ResponseEntity.NewCustomServiceError("Invalid schedule date, schedule time cannot be in the past", err)
		}

		req.ScheduledDate = schedule.Format(time.RFC3339)
		req.EndTime = t.timeSrv.CalcScheduleEndTimeString(schedule)
	}

	// create a reminder
	switch req.Repeat {
	case "never":
		err = t.remindSrv.SetReminder(req)
		if err != nil {
			// log.Println(err)
			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	case "daily":
		err = t.remindSrv.SetDailyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Daily Input")
		}
	case "weekly":
		err = t.remindSrv.SetWeeklyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Weekly Input")
		}
	case "bi-weekly":
		err = t.remindSrv.SetBiWeeklyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Bi Weekly Input")
		}
	case "monthly":
		err = t.remindSrv.SetMonthlyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Monthly Input")
		}
	case "yearly":
		err = t.remindSrv.SetYearlyReminder(req)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Yearly Input")
		}
	default:
		req.Repeat = "never"
		// err = t.remindSrv.SetReminder(req)

		// if err != nil {
		// 	log.Println("From setting reminder",err)
		// 	return nil, ResponseEntity.NewInternalServiceError(err)
		// }

	}

	// find features {Look for a better way to handle this!!!}
	var features taskEntity.TaskFeatures
	if req.Assigned != "" {
		features.IsAssigned = true
	}
	if req.ScheduledDate != "" {
		features.IsScheduled = true
	}

	data := taskEntity.CreateTaskRes{
		TaskId:        req.TaskId,
		Title:         req.Title,
		Description:   req.Description,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Notify:        req.Notify,
		VAOption:      req.VAOption,
		Repeat:        req.Repeat,
		TaskFeatures:  features,
		CreatedAt:     req.CreatedAt,
		ProjectId:     req.ProjectId,
		ScheduledDate: req.ScheduledDate,
	}

	tokens, vaId, username, err := t.nSrv.GetUserVaToken(req.UserId)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(vaId)

	body := []notificationEntity.NotificationBody{
		{
			Content: fmt.Sprintf("%s Just Created a Task", username),
			Color:   notificationEntity.CreatedColor,
			Time:    time.Now().UTC().String(),
		},
	}

	if vaId != "" {
		err = t.nSrv.CreateNotification(vaId, "Task Created", time.Now().String(), fmt.Sprintf("%s just created a new task", username), notificationEntity.CreatedColor, req.TaskId)
		if err != nil {
			fmt.Println("Error Uploading Notification to DB", err)
		}
	}
	if len(tokens) > 0 {
		err := t.nSrv.SendBatchNotifications(tokens, "Task Created", body, data)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("User Has Not VA or VA Has Not Registered For Notifications")
	}

	switch req.Assigned {
	case "assigned":
		err = t.repo.PersistAndAssign(ctx, req)
		if err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), `"virtual_Assistant_id": converting NULL to string is unsupported`) {
				return nil, ResponseEntity.NewInternalServiceError("YOU DON'T HAVE A VA. GET YA MONEY UP. BROKE BOY.")
			}

			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	default:
		// insert into db
		err = t.repo.Persist(ctx, req)
		if err != nil {
			// log.Println(err)
			return nil, ResponseEntity.NewInternalServiceError(err)
		}

	}

	return &data, nil
}

// Search task by name
// Search task godoc
// @Summary	Search task by title
// @Description	Search task route
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	q    query     string  false  "name search by q"
// @Success	200  {object}	[]taskEntity.SearchTaskRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task/search [get]
func (t *taskSrv) SearchTask(title *taskEntity.SearchTitleParams) ([]*taskEntity.SearchTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.validationSrv.Validate(title)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError(err)
	}
	tasks, err := t.repo.SearchTasks(title, ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return tasks, nil
}

// Get Task godoc
// @Summary	Get a single task
// @Description	Get a particular task
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	taskId	path	string	true	"Task Id"
// @Success	200  {object}  taskEntity.GetTasksByIdRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task/{taskId} [get]
func (t *taskSrv) GetTaskByID(taskId string) (*taskEntity.GetTasksByIdRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetTaskByID(ctx, taskId)

	if task == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	log.Println("From getByID", task)
	return task, nil
}

// Get Expired Tasks godoc
// @Summary	Get all expired tasks
// @Description	Get all expired task
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Success	200  {object}  []taskEntity.GetAllExpiredRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task/expired [get]
func (t *taskSrv) GetListOfExpiredTasks() ([]*taskEntity.GetAllExpiredRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetListOfExpiredTasks(ctx)

	if task == nil {
		log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError("No Task")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

// Get Pending Tasks godoc
// @Summary	Get list of pending tasks
// @Description	Get all pending task
// @Tags	VA - Tasks
// @Accept	json
// @Produce	json
// @Success	200  {object}  []taskEntity.GetAllPendingRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/all/pendingtasks [get]
func (t *taskSrv) GetListOfPendingTasks() ([]*taskEntity.GetAllPendingRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetListOfPendingTasks(ctx)

	if task == nil {
		log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError("No Task")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

// Get all task service
// Get All Tasks godoc
// @Summary	Get all tasks created by a user
// @Description	Get all tasks
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Success	200  {object}  []taskEntity.GetAllTaskRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task [get]

// Get all task service
// Get All Tasks godoc
// @Summary	Get all tasks created by a user
// @Description	Get all tasks
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	userId	path	string	true	"User Id"
// @Success	200  {object}  []taskEntity.GetAllTaskRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/user/task/{userId} [get]
func (t *taskSrv) GetAllTask(userId string) ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError) {
	log.Println("inside Fn")
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	task, err := t.repo.GetAllTasks(ctx, userId)

	if task == nil {
		log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return task, nil

}

// get all tasks of user assigned to va
// func (t *taskSrv) GetAllVaTask(userId string) ([]*taskEntity.GetAllTaskRes, *ResponseEntity.ServiceError) {
// 	log.Println("inside Fn")
// 	// create context of 1 minute
// 	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
// 	defer cancelFunc()

// 	task, err := t.repo.GetAllTasks(ctx, userId)

// 	if task == nil {
// 		log.Println("no rows returned")
// 		return nil, ResponseEntity.NewInternalServiceError(err)
// 	}
// 	if err != nil {
// 		log.Println(err)
// 		return nil, ResponseEntity.NewInternalServiceError(err)
// 	}
// 	return task, nil

// }

// Delete task by Id
// Delete task godoc
// @Summary	Delete task by Id
// @Description	Delete task by taskId
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	taskId	path	string	true	"Task Id"
// @Success	200  {object}  ResponseEntity.ResponseMessage
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task/{taskId} [delete]
func (t *taskSrv) DeleteTaskByID(taskId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteTaskByID(ctx, taskId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "Deleted successfully", nil, nil), nil
}

// Delete All Task
// Summary	Delete all users task
// Description	Delete all task route
// Tags	Tasks
// Accept	json
// Produce	json
// Success	200  {object}  ResponseEntity.ResponseMessage
// Failure	400  {object}  ResponseEntity.ServiceError
// Failure	404  {object}  ResponseEntity.ServiceError
// Failure	500  {object}  ResponseEntity.ServiceError
// Security ApiKeyAuth
// Router	/task [delete]
func (t *taskSrv) DeleteAllTask(userId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteAllTask(ctx, userId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "deleted user successfully", nil, nil), nil

}

// Update task status
// Update task status godoc
// @Summary	Update the status of a task
// @Description	Update task status route
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	taskId	path	string	true	"Task Id"
// @Param	request	body	taskEntity.UpdateTaskStatus	true	"Update task status"
// @Success	200  {object}  ResponseEntity.ResponseMessage
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task/{taskId}/status [post]
func (t *taskSrv) UpdateTaskStatusByID(taskId string, req *taskEntity.UpdateTaskStatus) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	//validating the struct
	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err = t.repo.UpdateTaskStatusByID(ctx, taskId, req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "Updated status successfully", nil, nil), nil

}

// Update task by Id
// Update task status godoc
// @Summary	Update the status of a task
// @Description	Update task status route
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	taskId	path	string	true	"Task Id"
// @Success	200  {object}  ResponseEntity.ResponseMessage
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/task/{taskId} [put]
func (t *taskSrv) EditTaskByID(taskId string, req *taskEntity.EditTaskReq) (*taskEntity.EditTaskRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	//validating the struct
	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}

	// Get task by ID
	task, err := t.repo.GetTaskByID(ctx, taskId)
	if task == nil {
		return nil, ResponseEntity.NewInternalServiceError("No task with that ID")
	}

	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	req1 := t.updateTask(req, task)
	req1.UpdatedAt = t.timeSrv.CurrentTimeString()

	//check if timeDueDate and StartDate is valid
	req1.EndTime, err = t.timeSrv.CheckFor339Format(req1.EndTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad End Time Input", err)
	}

	req1.StartTime, err = t.timeSrv.CheckFor339Format(req1.StartTime)
	if err != nil {
		return nil, ResponseEntity.NewCustomServiceError("Bad Start Time Input", err)
	}

	var schedule time.Time
	if req.ScheduledDate != "" {
		schedule, err = time.Parse(time.RFC3339, req1.ScheduledDate)

		ok := t.timeSrv.ScheduleTimeAfter(schedule)
		if err != nil || !ok {
			return nil, ResponseEntity.NewCustomServiceError("Invalid schedule date, schedule time cannot be in the past", err)
		}

		req1.ScheduledDate = schedule.Format(time.RFC3339)
		req1.EndTime = t.timeSrv.CalcScheduleEndTimeString(schedule)
	}

	switch req1.Repeat {
	case "never":
		err = t.remindSrv.SetReminder(req1)
		if err != nil {
			log.Println(err)
			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	case "daily":
		err = t.remindSrv.SetDailyReminder(req1)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Daily Input")
		}
	case "weekly":
		err = t.remindSrv.SetWeeklyReminder(req1)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Weekly Input")
		}
	case "bi-weekly":
		err = t.remindSrv.SetBiWeeklyReminder(req1)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Bi Weekly Input")
		}
	case "monthly":
		err = t.remindSrv.SetMonthlyReminder(req1)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Monthly Input")
		}
	case "yearly":
		err = t.remindSrv.SetYearlyReminder(req1)
		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Yearly Input")
		}
	default:
		req.Repeat = "never"
		// err = t.remindSrv.SetReminder(req)

		// if err != nil {
		// 	log.Println("From setting reminder",err)
		// 	return nil, ResponseEntity.NewInternalServiceError(err)
		// }

	}

	var features taskEntity.TaskFeatures
	if req.Assigned != "" {
		features.IsAssigned = true
	}
	if req.ScheduledDate != "" {
		features.IsScheduled = true
	}

	// Returning Data and Edit Data
	data := taskEntity.EditTaskReq{
		Title:         req1.Title,
		Description:   req1.Description,
		StartTime:     req1.StartTime,
		EndTime:       req1.EndTime,
		Notify:        req1.Notify,
		Repeat:        req1.Repeat,
		ProjectId:     req1.ProjectId,
		ScheduledDate: req1.ScheduledDate,
		Files:         req1.Files,
		Assigned:      req1.Assigned,
		UpdatedAt:     req1.UpdatedAt,
		Status:        req1.Status,
	}

	tokens, vaId, username, err := t.nSrv.GetUserVaToken(req1.UserId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vaId)

	if vaId == "" {
		log.Println(err)
		if strings.Contains(err.Error(), `"virtual_Assistant_id": converting NULL to string is unsupported`) {
			return nil, ResponseEntity.NewInternalServiceError("YOU DON'T HAVE A VA. GET YA MONEY UP. BROKE BOY.")
		}

		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	body := []notificationEntity.NotificationBody{
		{
			Content: fmt.Sprintf("%s Just Updated a Task", username),
			Color:   notificationEntity.CreatedColor,
			Time:    t.timeSrv.CurrentTimeString(),
		},
	}

	if len(tokens) > 0 {
		err := t.nSrv.SendBatchNotifications(tokens, "Task Updated", body, data)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("User Has Not VA or VA Has Not Registered For Notifications")
	}

	//Update Task
	err = t.repo.EditTaskById(ctx, taskId, &data)
	if err != nil {
		log.Println(err, "error updating data")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	// updateAt := t.timeSrv.CurrentTime().Format(time.RFC3339)
	// ndate := &taskEntity.CreateTaskReq{
	// 	TaskId:      taskId,
	// 	UserId:      task.UserId,
	// 	Title:       data.Title,
	// 	Description: data.Description,
	// 	Repeat:      data.Repeat,
	// 	StartTime:   data.StartTime,
	// 	EndTime:     data.EndTime,
	// 	VAOption:    data.VAOption,
	// 	Status:      data.Status,
	// 	UpdatedAt:   updateAt,
	// 	CreatedAt:   task.CreatedAt,
	// }

	// // delete former task
	// _, err2 := t.DeleteTaskByID(taskId)
	// if err2 != nil {
	// 	log.Println(err2)
	// 	return nil, err2
	// }
	// log.Println("Deleted task line 381")

	// create new task

	// //check if timeDueDate and StartDate is valid
	// err = t.timeSrv.CheckFor339Format(ndate.EndTime)
	// if err != nil {
	// 	return nil, ResponseEntity.NewCustomServiceError("Bad Time Input", err)
	// }

	// err = t.timeSrv.CheckFor339Format(ndate.StartTime)
	// if err != nil {
	// 	return nil, ResponseEntity.NewCustomServiceError("Bad Time Input", err)
	// }

	// create a reminder
	// switch req.Repeat {
	// case "never":
	// 	err = t.remindSrv.SetReminder(ndate)

	// 	if err != nil {
	// 		log.Println(err)
	// 		return nil, ResponseEntity.NewInternalServiceError(err)
	// 	}
	// case "daily":
	// 	err = t.remindSrv.SetDailyReminder(ndate)
	// 	if err != nil {
	// 		return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Daily Input")
	// 	}
	// case "weekly":
	// 	err = t.remindSrv.SetWeeklyReminder(ndate)
	// 	if err != nil {
	// 		return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Weekly Input")
	// 	}
	// case "bi-weekly":
	// 	err = t.remindSrv.SetBiWeeklyReminder(ndate)
	// 	if err != nil {
	// 		return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Bi Weekly Input")
	// 	}
	// case "monthly":
	// 	err = t.remindSrv.SetMonthlyReminder(ndate)
	// 	if err != nil {
	// 		return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Monthly Input")
	// 	}
	// case "yearly":
	// 	err = t.remindSrv.SetYearlyReminder(ndate)
	// 	if err != nil {
	// 		return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Yearly Input")
	// 	}
	// default:
	// 	return nil, ResponseEntity.NewInternalServiceError("Bad Recurrent Input(check enum data)")
	// }

	// insert into db
	// err = t.repo.Persist(ctx, ndate)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, ResponseEntity.NewInternalServiceError(err)
	// }
	log.Println("update complete")

	return &taskEntity.EditTaskRes{
		Title:        data.Title,
		Description:  data.Description,
		Repeat:       data.Repeat,
		StartTime:    data.StartTime,
		EndTime:      data.EndTime,
		Status:       data.Status,
		UpdatedAt:    data.UpdatedAt,
		TaskFeatures: features,
	}, nil
}

// Create a comment
// Create Comment godoc
// @Summary	Create comment for a task
// @Description	Create comment route
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	request	body	taskEntity.CreateCommentReq	true "Create Comment"
// @Success	200  {object}  taskEntity.CreateCommentRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/comment [post]
func (t *taskSrv) PersistComment(req *taskEntity.CreateCommentReq) (*taskEntity.CreateCommentRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	// implement validation for struct

	err := t.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}

	//set time
	req.CreatedAt = t.timeSrv.CurrentTimeString() // Format(time.RFC3339)

	// insert into db
	err = t.repo.PersistComment(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := taskEntity.CreateCommentRes{
		TaskId:  req.TaskId,
		Comment: req.Comment,
	}

	return &data, nil

}

// Get all comments for a task godoc
// @Summary	Get all comments by both user and VA on a task
// @Description	Get comments for a task route
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	taskId	path	string	true	"Task Id"
// @Success	200  {object}  []taskEntity.GetCommentRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/comment/{taskId} [get]
func (t *taskSrv) GetAllComments(taskId string) ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	comments, err := t.repo.GetAllComments(ctx, taskId)

	if comments == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return comments, nil

}

// Get all comments godoc
// @Summary	Get all comments by both user and VA on a task
// @Description	Create task
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Success	200  {object}  []taskEntity.GetCommentRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/comment/all [get]
func (t *taskSrv) GetComments() ([]*taskEntity.GetCommentRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	comments, err := t.repo.GetComments(ctx)

	if comments == nil {
		log.Println("no rows returned")
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return comments, nil

}

// Delete Comment By Id godoc
// @Summary	Delete a particular comment using it's id
// @Description	Delete comment route
// @Tags	Tasks
// @Accept	json
// @Produce	json
// @Param	commentId	path	string	true	"Comment Id"
// @Success	200  {object}  ResponseEntity.ResponseMessage
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/comment/{commentId} [delete]
func (t *taskSrv) DeleteCommentByID(commentId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := t.repo.DeleteCommentByID(ctx, commentId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "Deleted successfully", nil, nil), nil
}

// Auxillary function
func (t *taskSrv) updateTask(req *taskEntity.EditTaskReq, task *taskEntity.GetTasksByIdRes) *taskEntity.CreateTaskReq {
	log.Println(task)

	if req.Title != "" && req.Title != task.Title {
		task.Title = req.Title
	}

	if req.Description != "" && req.Description != task.Description {
		task.Description = req.Description
	}

	if !req.Notify && req.Notify != task.Notify {
		task.Notify = req.Notify
	}

	if req.Status != "" && req.Status != task.Status {
		task.Status = req.Status
	}

	if req.Assigned != "" && req.Assigned != task.Assigned {
		task.Assigned = req.Assigned
	}

	if req.Repeat != "" && req.Repeat != task.Repeat {
		task.Repeat = req.Repeat
	}

	if req.ProjectId != "" && req.ProjectId != task.ProjectId {
		task.ProjectId = req.ProjectId
	}

	if req.Assigned != "" && req.Assigned != task.Assigned {
		task.Assigned = req.Assigned
	}

	if req.StartTime != "" {
		task.StartTime = req.StartTime
	}

	if req.EndTime != "" {
		task.EndTime = req.EndTime
	}

	if req.ScheduledDate != "" && req.ScheduledDate != task.ScheduledDate {
		task.ScheduledDate = req.ScheduledDate
	}

	log.Println(task)
	return &taskEntity.CreateTaskReq{
		TaskId:        task.TaskId,
		UserId:        task.UserId,
		Title:         task.Title,
		Description:   task.Description,
		Notify:        task.Notify,
		Status:        task.Status,
		Repeat:        task.Repeat,
		Assigned:      task.Assigned,
		Files:         task.Files,
		VAOption:      task.VAOption,
		StartTime:     task.StartTime,
		EndTime:       task.EndTime,
		ProjectId:     task.ProjectId,
		UpdatedAt:     task.UpdatedAt,
		ScheduledDate: task.ScheduledDate,
	}
}
