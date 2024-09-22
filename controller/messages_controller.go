package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/yaninyzwitty/messaging-service/helpers"
	"github.com/yaninyzwitty/messaging-service/models"
	"github.com/yaninyzwitty/messaging-service/service"
)

type MessageController struct {
	service service.MessagesService
}

func NewMessageController(service service.MessagesService) *MessageController {
	return &MessageController{service: service}
}

func (c *MessageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	var ctx = r.Context()
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if message.Body == "" {
		http.Error(w, "Message body cannot be empty", http.StatusBadRequest)
		return
	}
	if message.ConversationID.String() == "" || message.SenderId.String() == "" {
		http.Error(w, "Conversation ID and Sender ID are required", http.StatusBadRequest)
		return
	}

	// initialize the defaults
	message.ID = gocql.TimeUUID()
	message.CreatedAt = time.Now()
	message.IsSoftDeleted = false
	createdMessage, err := c.service.CreateMessage(ctx, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = helpers.NewResponseToJson(w, http.StatusCreated, createdMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *MessageController) GetMessages(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	messages, err := c.service.GetMessages(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = helpers.NewResponseToJson(w, http.StatusOK, messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *MessageController) GetMessagesByPagingState(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	// Read query parameters for page size and paging state
	pageSizeInStr := r.URL.Query().Get("pageSize")
	pagingStateInStr := r.URL.Query().Get("paging_state")

	pageSize, err := strconv.Atoi(pageSizeInStr)
	if err != nil || pageSize <= 0 {
		// use default page size
		pageSize = 10
	}
	pagingState := []byte(pagingStateInStr) //we are getting the previous paging state

	// we try to fetch paginated messages
	messages, newPagingState, err := c.service.GetMessagesByPagingState(ctx, pageSize, pagingState)
	if err != nil {
		http.Error(w, "error getting the paginated messages"+err.Error(), http.StatusInternalServerError)
	}
	// Return messages and next page token (paging state)
	response := models.MessageWithPagingState{
		Messages:      messages,
		NextPageToken: string(newPagingState),
	}

	err = helpers.NewResponseToJson(w, http.StatusOK, response)
	if err != nil {
		http.Error(w, "error decoding the response"+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *MessageController) GetMessage(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	idStr := r.PathValue("id")
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Failed to parse the id into gocql uuid format", http.StatusBadRequest)
		return
	}
	message, err := c.service.GetMessage(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get the message: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = helpers.NewResponseToJson(w, http.StatusOK, message)
	if err != nil {
		http.Error(w, "Failed to fully decode the message"+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *MessageController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	var ctx = r.Context()
	var idStr = r.PathValue("id")

	// Parse UUID from the path
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Failed to parse the id into gocql UUID format", http.StatusBadRequest)
		return
	}

	// Parse the request body into the message object
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request payload"+err.Error(), http.StatusBadRequest)
		return
	}

	// Set the updated timestamp
	message.UpdatedAt = time.Now()

	// Ensure the ID in the URL is the one being updated
	message.ID = id

	// Call the service to update the message
	updatedMessage, err := c.service.UpdateMessage(ctx, id, message)
	if err != nil {
		http.Error(w, "failed to update message"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response back
	err = helpers.NewResponseToJson(w, http.StatusOK, updatedMessage)
	if err != nil {
		http.Error(w, "Error marshaling the response", http.StatusInternalServerError)
		return
	}
}

func (c *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	idStr := r.PathValue("id")
	id, err := gocql.ParseUUID(idStr)
	if err != nil {
		http.Error(w, "Failed to parse the id into gocql uuid format", http.StatusBadRequest)
		return
	}

	err = c.service.DeleteMessage(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = helpers.NewResponseToJson(w, http.StatusOK, "Message deleted successfully")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
