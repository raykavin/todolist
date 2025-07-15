package handler

import (
	"errors"
	netHttp "net/http"
	"todolist/internal/adapter/delivery/http"
	"todolist/internal/domain/person/valueobject"
	"todolist/internal/domain/shared"
	"todolist/internal/dto"
	ucPerson "todolist/internal/usecase/person"
)

// PersonHandler handles person-related HTTP requests
type PersonHandler struct {
	createPersonUC ucPerson.CreatePersonUseCase
	updatePersonUC ucPerson.UpdatePersonUseCase
	getPersonUC    ucPerson.GetPersonUseCase
}

// NewPersonHandler creates a new person handler
func NewPersonHandler(
	createPersonUC ucPerson.CreatePersonUseCase,
	updatePersonUC ucPerson.UpdatePersonUseCase,
	getPersonUC ucPerson.GetPersonUseCase,
) *PersonHandler {
	return &PersonHandler{
		createPersonUC: createPersonUC,
		updatePersonUC: updatePersonUC,
		getPersonUC:    getPersonUC,
	}
}

// Register registers person routes
// func (h *PersonHandler) Register(router *gin.RouterGroup) {
// 	people := router.Group("/people")
// 	{
// 		// Public routes
// 		people.POST("/", h.CreatePerson)

// 		// Protected routes
// 		people.Use(authMiddleware)
// 		people.GET("/:id", h.GetPerson)
// 		people.PUT("/:id", h.UpdatePerson)
// 	}
// }

// CreatePerson godoc
// @Summary Create a new person
// @Description Create a new person record
// @Tags people
// @Accept json
// @Produce json
// @Param person body dto.CreatePersonRequest true "Person data"
// @Success 201 {object} dto.Response{data=dto.PersonResponse}
// @Failure 400 {object} dto.Response
// @Failure 409 {object} dto.Response
// @Router /api/v1/people [post]
func (h *PersonHandler) CreatePerson(ctx http.RequestContext) {
	var input dto.CreatePersonRequest

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", parseError(err)))

		ctx.Abort()
		return
	}

	// Create person
	person, err := h.createPersonUC.Execute(ctx.Context(), input)
	if err != nil {
		if errors.Is(err, valueobject.ErrInvalidEmail) {
			ctx.JSON(netHttp.StatusConflict,
				dto.ErrorResponse("EMAIL_EXISTS", "Email already exists", nil))
		} else {
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("CREATE_FAILED", "Failed to create person", nil))
		}

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusCreated, dto.SuccessResponse(person, "Person created successfully"))
}

// GetPerson godoc
// @Summary Get person by ID
// @Description Get person details by ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} dto.Response{data=dto.PersonResponse}
// @Failure 404 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/people/{id} [get]
func (h *PersonHandler) GetPerson(ctx http.RequestContext) {
	personID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid ID parameter", parseError(err)))

		ctx.Abort()
		return
	}

	person, err := h.getPersonUC.Execute(ctx.Context(), personID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			ctx.JSON(netHttp.StatusNotFound,
				dto.ErrorResponse("NOT_FOUND", "Person not found", nil))
		} else {
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("GET_FAILED", "Failed to get person", nil))
		}

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(person, ""))
}

// UpdatePerson godoc
// @Summary Update person
// @Description Update person details
// @Tags people
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Param person body dto.UpdatePersonRequest true "Person data"
// @Success 200 {object} dto.Response{data=dto.PersonResponse}
// @Failure 400 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/people/{id} [put]
func (h *PersonHandler) UpdatePerson(ctx http.RequestContext) {
	personID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid ID parameter", parseError(err)))

		ctx.Abort()
		return
	}

	// userRole := middleware.GetUserRole(ctx)

	// // Only admin or the person themselves can update
	// if userRole != "admin" {
	// 	// In a real app, check if the person belongs to the authenticated user
	// 	c.JSON(http.StatusForbidden, httpDto.ErrorResponse("FORBIDDEN", "Access denied", nil))
	// 	return
	// }

	var input dto.UpdatePersonRequest

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", parseError(err)))

		ctx.Abort()
		return
	}

	// Update person
	person, err := h.updatePersonUC.Execute(ctx.Context(), personID, input)
	if err != nil {
		if err == shared.ErrNotFound {
			ctx.JSON(netHttp.StatusNotFound,
				dto.ErrorResponse("NOT_FOUND", "Person not found", nil))
		} else {
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("UPDATE_FAILED", "Failed to update person", nil))
		}

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK,
		dto.SuccessResponse(person, "Person updated successfully"))
}
