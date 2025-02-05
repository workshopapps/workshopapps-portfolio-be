package callService

import (
	"context"
	"log"
	"test-va/internals/Repository/callRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/callEntity"
	"test-va/internals/service/loggerService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"
)

type CallService interface {
	GetCalls() ([]*callEntity.CallRes, *ResponseEntity.ResponseMessage)
}

type callSrv struct {
	repo          callRepo.CallRepository
	timeSrv       timeSrv.TimeService
	validationSrv validationService.ValidationSrv
	logger        loggerService.LogSrv
}

// Get calls godoc
// @Summary	Get all your calls
// @Description	Get call route
// @Tags	Calls
// @Accept	json
// @Produce	json
// @Success	200  {object}	[]callEntity.CallRes
// @Failure	400  {object}  ResponseEntity.ResponseMessage
// @Failure	404  {object}  ResponseEntity.ResponseMessage
// @Failure	500  {object}  ResponseEntity.ResponseMessage
// @Router	/calls [get]
func (t callSrv) GetCalls() ([]*callEntity.CallRes, *ResponseEntity.ResponseMessage) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	calls, err := t.repo.GetCalls(ctx)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewCustomError(500, "Internal Server Error")
	}
	return calls, nil
}

func NewCallSrv(repo callRepo.CallRepository, timeSrv timeSrv.TimeService, srv validationService.ValidationSrv, logSrv loggerService.LogSrv) CallService {
	return &callSrv{repo: repo, timeSrv: timeSrv, validationSrv: srv, logger: logSrv}
}
