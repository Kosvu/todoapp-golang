package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/Kosvu/todoapp-golang/internal/core/logger"
	core_http_request "github.com/Kosvu/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/Kosvu/todoapp-golang/internal/core/transport/http/response"
)

func (h *TasksHTTPHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPRespsonseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	if err := h.taskService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)

		return
	}

	responseHandler.NoContentResponse()
}
