package projectHandler

import (
	"log"
	"net/http"

	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/projectEntity"
	"test-va/internals/service/projectService"

	"github.com/gin-gonic/gin"
)

type projectHandler struct {
	srv projectService.ProjectService
}

func NewProjectHandler(srv projectService.ProjectService) *projectHandler {
	return &projectHandler{srv: srv}
}

func (p *projectHandler) CreateProject(c *gin.Context){
	var req projectEntity.CreateProjectReq
	value := c.GetString("userId")
	log.Println("userId is: ", value)
	if value == "" {

		c.AbortWithStatusJSON(http.StatusUnauthorized, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "you are not allowed to access this resource", nil, nil))
		return
	}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}
	req.UserId = value
	log.Println("create project req",req)

	project, errRes := p.srv.PersistProject(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error creating Project", errRes, nil))
		return
	}

	c.JSON(http.StatusOK,
		ResponseEntity.BuildSuccessResponse(http.StatusOK, "Created Project Successfully", project, nil))

}

func (p *projectHandler) GetAllUsersProjects(c *gin.Context){
	userId := c.GetString("userId")
	log.Println(userId)
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	projects, errRes := p.srv.GetListOfUsersProjects(userId)
	if projects == nil{
		message := "user with id " + userId + " has no project"
		c.AbortWithStatusJSON(http.StatusOK,
			ResponseEntity.BuildSuccessResponse(http.StatusNoContent, message, projects, nil))
		return
	}
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failure To Find all users project", errRes, nil))
		return
	}

	c.JSON(http.StatusOK,
		ResponseEntity.BuildSuccessResponse(http.StatusOK, "Users projects returned successfully", projects, nil))
}


