package vaService

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"test-va/internals/Repository/vaRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/vaEntity"
	"test-va/internals/service/cryptoService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"

	"github.com/google/uuid"
)

type VAService interface {
	SignUp(req *vaEntity.CreateVAReq) (*vaEntity.CreateVARes, *ResponseEntity.ServiceError)
	Login(req *vaEntity.LoginReq) (*vaEntity.FindByEmailRes, *ResponseEntity.ServiceError)
	GetVA(id string) (*vaEntity.FindByIdRes, *ResponseEntity.ServiceError)
	FindByEmail(email string) (*vaEntity.FindByEmailRes, *ResponseEntity.ServiceError)
	UpdateVA(req *vaEntity.EditVaReq, id string) (*vaEntity.EditVARes, *ResponseEntity.ServiceError)
	ChangePassword(req *vaEntity.ChangeVAPassword) *ResponseEntity.ServiceError
	DeleteVA(id string) *ResponseEntity.ServiceError
	GetAllUserToVa(vaId string) ([]*vaEntity.VAStruct, *ResponseEntity.ServiceError)
}

type vaSrv struct {
	repo      vaRepo.VARepo
	validator validationService.ValidationSrv
	timeSrv   timeSrv.TimeService
	cryptoSrv cryptoService.CryptoSrv
}

// Get All Users Assigned To VA godoc
// @Summary	Get all users in the system assigned to a particular VA
// @Description	Get all users assigned to VA route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	vaId	path	string	true	"VA Id"
// @Success	200  {object}  []vaEntity.VAStruct
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/user/{vaId} [get]
func (v *vaSrv) GetAllUserToVa(vaId string) ([]*vaEntity.VAStruct, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()
	va, err := v.repo.GetUserAssignedToVa(ctx, vaId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ResponseEntity.NewInternalServiceError("No User Found Pls")
		}
		return nil, ResponseEntity.NewInternalServiceError("Error Getting User")
	}

	return va, nil
}

// Login Virtual Assistant godoc
// @Summary	Provide email and password to be logged in
// @Description	Login as a va
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	request	body	userEntity.LoginReq	true "Login Details"
// @Success	200  {object}  vaEntity.FindByIdRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/va/login [post]
func (v *vaSrv) Login(req *vaEntity.LoginReq) (*vaEntity.FindByEmailRes, *ResponseEntity.ServiceError) {
	// validate request first
	err := v.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(fmt.Sprintf("Bad Request: %v", err))
	}

	//find the user with email
	user, errRes := v.FindByEmail(req.Email)
	if errRes != nil {
		log.Println("err")
		return nil, ResponseEntity.NewValidatingError("Email Not Found")
	}

	//compare passwords
	err = v.cryptoSrv.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError("Passwords Don't Match")
	}

	return user, nil
}

func (v *vaSrv) FindByEmail(email string) (*vaEntity.FindByEmailRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	user, err := v.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(fmt.Sprintf("Error Finding User: %v", err))
	}
	return user, nil
}

// Update VA godoc
// @Summary	Update a virtual assistant profile
// @Description	Update va route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	vaId	path	string	true	"Virtual Assistant Id"
// @Param	request	body	vaEntity.EditVaReq	true "Update VA Details"
// @Success	200  {object}  vaEntity.EditVARes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/va/{vaId} [post]
func (v *vaSrv) UpdateVA(req *vaEntity.EditVaReq, id string) (*vaEntity.EditVARes, *ResponseEntity.ServiceError) {
	// validate request first
	err := v.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(fmt.Sprintf("Bad Request: %v", err))
	}

	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err = v.repo.UpdateUser(ctx, req, id)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(fmt.Sprintf("Error Updating User: %v", err))
	}
	data := vaEntity.EditVARes{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		ProfilePicture: req.ProfilePicture,
	}
	return &data, nil
}

// Change VA Password godoc
// @Summary	Change a va password
// @Description	Change va password route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Success	200  {string}  string    "ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/va/change-password [post]
func (v *vaSrv) ChangePassword(req *vaEntity.ChangeVAPassword) *ResponseEntity.ServiceError {
	// validate request first
	err := v.validator.Validate(req)
	if err != nil {
		return ResponseEntity.NewValidatingError(fmt.Sprintf("Bad Request: %v", err))
	}

	// hash new password
	pass, err := v.cryptoSrv.HashPassword(req.NewPassword)
	if err != nil {
		return ResponseEntity.NewInternalServiceError(err)
	}
	req.NewPassword = pass

	// create context
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	// change password
	errRes := v.repo.ChangePassword(ctx, req)
	if errRes != nil {
		return ResponseEntity.NewInternalServiceError(fmt.Sprintf("Failed to Change Password: %v", err))
	}

	return nil
}

// Delete VA godoc
// @Summary	Delete a va from the database
// @Description	Delete va route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	vaId	path	string	true	"VA Id"
// @Success	200  {string}  string    "ok"
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/va/{vaId} [delete]
func (v *vaSrv) DeleteVA(id string) *ResponseEntity.ServiceError {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := v.repo.DeleteUser(ctx, id)
	if err != nil {
		return ResponseEntity.NewInternalServiceError(fmt.Sprintf("Failed to delete User: %v", err))
	}
	return nil
}

// Register VA godoc
// @Summary	Register a virtual assistant
// @Description	Register va route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	request	body	vaEntity.CreateVAReq	true "VA Details"
// @Success	200  {object}  vaEntity.CreateVARes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Security ApiKeyAuth
// @Router	/va/signup [post]
func (v *vaSrv) SignUp(req *vaEntity.CreateVAReq) (*vaEntity.CreateVARes, *ResponseEntity.ServiceError) {
	err := v.validator.Validate(req)
	if err != nil {
		return nil, ResponseEntity.NewValidatingError(fmt.Sprintf("Bad Request: %v", err))
	}

	//find the user with email
	user, errRes := v.FindByEmail(req.Email)
	if errRes == nil {
		log.Println(user)
		log.Println(errRes)
		return nil, ResponseEntity.NewValidatingError("User Already Exists")
	}

	//compare passwords
	pass, err := v.cryptoSrv.HashPassword(req.Password)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError("Passwords Don't Match")
	}

	req.CreatedAt = v.timeSrv.CurrentTime().Format(time.RFC3339)
	req.Password = pass
	req.VaId = uuid.New().String()

	// save user to repo
	// create context
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err = v.repo.Persist(ctx, req)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(fmt.Sprintf("Error creating User: %v", err))
	}

	// return user
	data := vaEntity.CreateVARes{
		VaId:           req.VaId,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		ProfilePicture: req.ProfilePicture,
		AccountType:    req.AccountType,
	}

	return &data, nil
}

// Get Specific VA godoc
// @Summary	Get a particular VA by the Id
// @Description	Get VA route
// @Tags	VA
// @Accept	json
// @Produce	json
// @Param	vaId	path	string	true	"Virtual Assistant Id"
// @Success	200  {object}  vaEntity.FindByIdRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/va/{vaId} [get]
func (v *vaSrv) GetVA(id string) (*vaEntity.FindByIdRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	user, err := v.repo.FindById(ctx, id)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(fmt.Sprintf("Error Finding User: %v", err))
	}
	return user, nil
}

func NewVaService(repo vaRepo.VARepo, validator validationService.ValidationSrv,
	timeSrv timeSrv.TimeService, cryptoSrv cryptoService.CryptoSrv) VAService {
	return &vaSrv{repo: repo, validator: validator, timeSrv: timeSrv, cryptoSrv: cryptoSrv}
}
