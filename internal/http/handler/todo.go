package handler

import (
	"errors"
	netHttp "net/http"
	"strconv"
	"strings"
	"todolist/internal/adapter/delivery/http"
	"todolist/internal/domain/shared"
	"todolist/internal/domain/todo/entity"
	"todolist/internal/domain/todo/valueobject"
	"todolist/internal/dto"
	ucTodo "todolist/internal/usecase/todo"
)

// TodoHandler handles todo-related HTTP requests
type TodoHandler struct {
	createTodoUC    ucTodo.CreateTodoUseCase
	updateTodoUC    ucTodo.UpdateTodoUseCase
	completeTodoUC  ucTodo.CompleteTodoUseCase
	deleteTodoUC    ucTodo.DeleteTodoUseCase
	getTodoUC       ucTodo.GetTodoUseCase
	listTodosUC     ucTodo.ListTodosUseCase
	getStatisticsUC ucTodo.GetStatisticsUseCase
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(
	createTodoUC ucTodo.CreateTodoUseCase,
	updateTodoUC ucTodo.UpdateTodoUseCase,
	completeTodoUC ucTodo.CompleteTodoUseCase,
	deleteTodoUC ucTodo.DeleteTodoUseCase,
	getTodoUC ucTodo.GetTodoUseCase,
	listTodosUC ucTodo.ListTodosUseCase,
	getStatisticsUC ucTodo.GetStatisticsUseCase,
) *TodoHandler {
	return &TodoHandler{
		createTodoUC:    createTodoUC,
		updateTodoUC:    updateTodoUC,
		completeTodoUC:  completeTodoUC,
		deleteTodoUC:    deleteTodoUC,
		getTodoUC:       getTodoUC,
		listTodosUC:     listTodosUC,
		getStatisticsUC: getStatisticsUC,
	}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body dto.CreateTodoRequest true "Todo data"
// @Success 201 {object} dto.Response{data=dto.TodoResponse}
// @Failure 400 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/todos [post]
func (h *TodoHandler) CreateTodo(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))

		ctx.Abort()
		return
	}

	var input dto.CreateTodoRequest

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", parseError(err)))

		ctx.Abort()
		return
	}

	// Create todo
	todo, err := h.createTodoUC.Execute(ctx.Context(), userID, input)
	if err != nil {
		ctx.JSON(netHttp.StatusInternalServerError,
			dto.ErrorResponse("CREATE_FAILED", "Failed to create todo", nil))

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusCreated, dto.SuccessResponse(todo, "Todo created successfully"))
}

// ListTodos godoc
// @Summary List todos
// @Description List todos with filters and pagination
// @Tags todos
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param status query []string false "Filter by status" Enums(pending,in_progress,completed,cancelled)
// @Param priority query []string false "Filter by priority" Enums(low,medium,high,critical)
// @Param tags query []string false "Filter by tags"
// @Param search query string false "Search in title and description"
// @Param is_overdue query bool false "Filter overdue todos"
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.TodoResponse}
// @Failure 400 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/todos [get]
func (h *TodoHandler) ListTodos(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))

		ctx.Abort()
		return
	}

	// Parse query parameters
	var queryParams dto.QueryParams
	if err := ctx.BindQuery(&queryParams); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_QUERY", "Invalid query parameters", parseError(err)))

		ctx.Abort()
		return
	}
	queryParams.SetDefaults()

	// Build filters
	filters := valueobject.TodoFilterCriteria{
		UserID:     userID,
		SearchTerm: queryParams.Search,
	}

	// Parse status filter
	if statusStr := ctx.GetQuery("status"); statusStr != "" {
		filters.Status = strings.Split(statusStr, ",")
	}

	// Parse priority filter
	if priorityStr := ctx.GetQuery("priority"); priorityStr != "" {
		filters.Priority = strings.Split(priorityStr, ",")
	}

	// Parse tags filter
	if tagsStr := ctx.GetQuery("tags"); tagsStr != "" {
		filters.Tags = strings.Split(tagsStr, ",")
	}

	// Parse is_overdue filter
	if overdueStr := ctx.GetQuery("is_overdue"); overdueStr != "" {
		isOverdue, _ := strconv.ParseBool(overdueStr)
		filters.IsOverdue = &isOverdue
	}

	// Build query options
	options := shared.QueryOptions{
		Limit:     queryParams.PageSize,
		Offset:    queryParams.GetOffset(),
		OrderBy:   queryParams.OrderBy,
		OrderDesc: queryParams.OrderDir == "desc",
	}

	// Default ordering
	if options.OrderBy == "" {
		options.OrderBy = "created_at"
		options.OrderDesc = true
	}

	// List todos
	result, err := h.listTodosUC.Execute(ctx.Context(), userID, filters, options)
	if err != nil {
		ctx.JSON(netHttp.StatusInternalServerError,
			dto.ErrorResponse("LIST_FAILED", "Failed to list todos", nil))

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.PaginatedSuccessResponse(
		result.Todos,
		queryParams.Page,
		queryParams.PageSize,
		result.TotalCount,
	))
}

// GetTodo godoc
// @Summary Get todo by ID
// @Description Get todo details by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.Response{data=dto.TodoResponse}
// @Failure 404 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetTodo(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))

		ctx.Abort()
		return
	}

	todoID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid ID parameter", parseError(err)))

		ctx.Abort()
		return
	}

	todo, err := h.getTodoUC.Execute(ctx.Context(), userID, todoID)
	if err != nil {
		if err == shared.ErrNotFound || err.Error() == "unauthorized to access this todo" {
			ctx.JSON(netHttp.StatusNotFound,
				dto.ErrorResponse("NOT_FOUND", "Todo not found", nil))
		} else {
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("GET_FAILED", "Failed to get todo", nil))

		}

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(todo, ""))
}

// UpdateTodo godoc
// @Summary Update todo
// @Description Update todo details
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body dto.UpdateTodoRequest true "Todo data"
// @Success 200 {object} dto.Response{data=dto.TodoResponse}
// @Failure 400 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/todos/{id} [put]
func (h *TodoHandler) UpdateTodo(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))

		ctx.Abort()
		return
	}

	todoID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid ID parameter", parseError(err)))

		ctx.Abort()
		return
	}

	var input dto.UpdateTodoRequest

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", parseError(err)))

		ctx.Abort()
		return
	}

	// Update todo
	todo, err := h.updateTodoUC.Execute(ctx.Context(), userID, todoID, input)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrNotFound):
		case errors.Is(err, entity.ErrUnauthorizedTodoAccess):
			ctx.JSON(netHttp.StatusNotFound,
				dto.ErrorResponse("NOT_FOUND", "Todo not found", nil))
		case errors.Is(err, entity.ErrInvalidStatusTransition):
			ctx.JSON(netHttp.StatusBadRequest,
				dto.ErrorResponse("INVALID_TRANSITION", "Invalid status transition", nil))
		default:
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("UPDATE_FAILED", "Failed to update todo", nil))
		}

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(todo, "Todo updated successfully"))
}

// CompleteTodo godoc
// @Summary Complete todo
// @Description Mark todo as completed
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.Response{data=dto.TodoResponse}
// @Failure 400 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/todos/{id}/complete [post]
func (h *TodoHandler) CompleteTodo(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))

		ctx.Abort()
		return
	}

	todoID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid ID parameter", parseError(err)))

		ctx.Abort()
		return
	}

	todo, err := h.completeTodoUC.Execute(ctx.Context(), userID, todoID)
	if err != nil {
		switch {
		case errors.Is(err, shared.ErrNotFound):
		case errors.Is(err, entity.ErrUnauthorizedTodoAccess):
			ctx.JSON(netHttp.StatusNotFound,
				dto.ErrorResponse("NOT_FOUND", "Todo not found", nil))
		case errors.Is(err, entity.ErrTodoAlreadyCompleted):
			ctx.JSON(netHttp.StatusBadRequest,
				dto.ErrorResponse("ALREADY_COMPLETED", "Todo is already completed", nil))
		default:
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("COMPLETE_FAILED", "Failed to complete todo", nil))
		}

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(todo, "Todo completed successfully"))
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description Delete a todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))

		ctx.Abort()
		return
	}

	todoID, err := getIDParam(ctx)
	if err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid ID parameter", parseError(err)))

		ctx.Abort()
		return
	}

	err = h.deleteTodoUC.Execute(ctx.Context(), userID, todoID)
	if err != nil {
		if err == shared.ErrNotFound || err.Error() == "unauthorized to access this todo" {
			ctx.JSON(netHttp.StatusNotFound,
				dto.ErrorResponse("NOT_FOUND", "Todo not found", nil))
		} else {
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("DELETE_FAILED", "Failed to delete todo", nil))
		}

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(nil, "Todo deleted successfully"))
}

// GetStatistics godoc
// @Summary Get todo statistics
// @Description Get todo statistics for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response{data=valueobject.TodoStatistics}
// @Security BearerAuth
// @Router /api/v1/todos/statistics [get]
func (h *TodoHandler) GetStatistics(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))

		ctx.Abort()
		return
	}

	stats, err := h.getStatisticsUC.Execute(ctx.Context(), userID)
	if err != nil {
		ctx.JSON(netHttp.StatusInternalServerError,
			dto.ErrorResponse("STATS_FAILED", "Failed to get statistics", nil))

		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(stats, ""))
}

// // Register registers todo routes
// func (h *TodoHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
// 	todos := router.Group("/todos")
// 	todos.Use(authMiddleware)
// 	{
// 		todos.POST("/", h.CreateTodo)
// 		todos.GET("/", h.ListTodos)
// 		todos.GET("/statistics", h.GetStatistics)
// 		todos.GET("/:id", h.GetTodo)
// 		todos.PUT("/:id", h.UpdateTodo)
// 		todos.POST("/:id/complete", h.CompleteTodo)
// 		todos.DELETE("/:id", h.DeleteTodo)
// 	}
// }
