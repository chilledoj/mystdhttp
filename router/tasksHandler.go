package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// TasksHandler doubles as both an http service and a db
type TasksHandler struct {
	mu sync.RWMutex
	// data map indexes each task by it's ID
	data map[string]Task
}

// initData will populate the map with dummy data
func (th *TasksHandler) initData() {
	for i := 1; i < 20; i++ {
		st := StatusNew
		if i%2 == 0 {
			st = StatusStarted
		} else if i%5 == 0 {
			st = StatusInactive
		} else if i == 7 {
			st = StatusCompleted
		}
		tsk := NewTask(
			fmt.Sprintf("Task %2d", i),
			fmt.Sprintf(`# %04x - Test %d

Task %d details

- Do Something
- Another thing to do`, i*i, i, i),
			st,
		)
		tsk.AssignedTo = fmt.Sprintf("USER%03d", i%4)
		th.data[tsk.ID] = tsk
	}
}

func NewTasksHandler(initTasks bool) *TasksHandler {
	th := &TasksHandler{
		data: make(map[string]Task),
	}
	if initTasks {
		th.initData()
	}
	return th
}

func (th *TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/" {
			th.listTaskHandler(w, r)
			return
		}
		th.readTaskHandler(w, r)
	case http.MethodPost:
		th.createTaskHandler(w, r)
	case http.MethodPut:
		th.updateTaskHandler(w, r)
	case http.MethodDelete:
		th.deleteTaskHandler(w, r)
	default:
		msg := http.StatusText(http.StatusNotFound)
		ErrJSON(w, httpErr{err: fmt.Errorf(msg), msg: msg}, http.StatusNotFound)
	}
}

/* PAGE ROUTES */

func (th *TasksHandler) taskListPage() http.HandlerFunc {
	files := tmplLayout("./web/templates/index.gohtml")
	tmpl := template.Must(template.New("index").Funcs(defaultFuncs).ParseFiles(files...))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.Method != http.MethodGet {
			notFoundPage(w, r)
			return
		}
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "base", nil); err != nil {
			fmt.Printf("ERR: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, &buf)
	}
}

func (th *TasksHandler) taskViewPage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := extractTaskId(r)
		if err != nil {
			notFoundPage(w, r)
			return
		}
		th.mu.RLock()
		defer th.mu.RUnlock()
		task, ok := th.data[taskId]
		if !ok {
			notFoundPage(w, r)
			return
		}

		files := tmplLayout("./web/templates/taskView.gohtml")
		tmpl := template.Must(template.New("index").Funcs(defaultFuncs).ParseFiles(files...))

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "base", map[string]interface{}{
			"Task": task,
		}); err != nil {
			fmt.Printf("ERR: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, &buf)
	}
}

func (th *TasksHandler) taskCreatePage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		files := tmplLayout("./web/templates/taskForm.gohtml")
		tmpl := template.Must(template.New("form").Funcs(defaultFuncs).ParseFiles(files...))

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "base", map[string]interface{}{
			"FormType": "new",
		}); err != nil {
			fmt.Printf("ERR: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, &buf)
	}
}

func (th *TasksHandler) taskEditPage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := extractTaskId(r)
		if err != nil {
			notFoundPage(w, r)
			return
		}
		th.mu.RLock()
		defer th.mu.RUnlock()
		task, ok := th.data[taskId]
		if !ok {
			notFoundPage(w, r)
			return
		}

		files := tmplLayout("./web/templates/taskForm.gohtml")
		tmpl := template.Must(template.New("form").Funcs(defaultFuncs).ParseFiles(files...))

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "base", map[string]interface{}{
			"FormType": "edit",
			"Task":     task,
		}); err != nil {
			fmt.Printf("ERR: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, &buf)
	}
}

/* API ROUTES */

// listTaskHandler: GET /api/tasks -> returns all tasks.
func (th *TasksHandler) listTaskHandler(w http.ResponseWriter, r *http.Request) {
	th.mu.RLock()
	defer th.mu.RUnlock()
	lists := make([]Task, 0)
	for _, t := range th.data {
		if t.Deleted {
			continue
		}
		lists = append(lists, t)
	}
	JSON(w, lists, http.StatusOK)
}

// extractTaskId is a helper function to extract the task id path param.
// This is very basic in that it just provides everything after the /
func extractTaskId(r *http.Request) (string, error) {
	if len(r.URL.Path) < 2 {
		return "", fmt.Errorf("no task id provided")
	}
	// We know our subrouter will only provide the id as the root path item.
	return strings.Split(r.URL.Path[1:], "/")[0], nil
}

// readTaskHandler: GET /api/tasks/:taskId -> returns a single task by taskId.
func (th *TasksHandler) readTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId, err := extractTaskId(r)
	if err != nil {
		ErrJSON(w, httpErr{err: err, msg: "not found"}, http.StatusNotFound)
		return
	}
	th.mu.RLock()
	defer th.mu.RUnlock()
	task, ok := th.data[taskId]
	if !ok {
		ErrJSON(w, httpErr{err: err, msg: "not found"}, http.StatusNotFound)
		return
	}
	JSON(w, task, http.StatusOK)
}

// deleteTaskHandler: DELETE /api/tasks/:taskId -> deletes a task by taskId.
// No content is returned
func (th *TasksHandler) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId, err := extractTaskId(r)
	if err != nil {
		ErrJSON(w, httpErr{err: err, msg: "not found"}, http.StatusNotFound)
		return
	}
	th.mu.Lock()
	defer th.mu.Unlock()
	v, ok := th.data[taskId]
	if !ok {
		ErrJSON(w, httpErr{err: err, msg: "not found"}, http.StatusNotFound)
		return
	}
	// soft delete
	v.Deleted = true
	th.data[taskId] = v

	w.WriteHeader(http.StatusNoContent)
}

// JSON API

type taskRequest struct {
	Title   string `json:"title"`
	Details string `json:"details"`
	Status  string `json:"status,omitempty"`
}

// createTaskHandler: POST /api/tasks/:taskId -> creates a task.
func (th *TasksHandler) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	tr := taskRequest{}
	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		ErrJSON(w, httpErr{err: err, msg: "bad request"}, http.StatusBadRequest)
		return
	}

	task := NewTask(tr.Title, tr.Details, StatusNew)

	th.mu.Lock()
	th.data[task.ID] = task
	th.mu.Unlock()

	w.Header().Set("Location", fmt.Sprintf("/api/task/%s", task.ID))
	w.Header().Set("X-ViewLocation", fmt.Sprintf("/view/%s", task.ID))
	JSON(w, task, http.StatusCreated)
}

// updateTaskHandler: PUT /api/tasks/:taskId -> updates a task by taskId.
func (th *TasksHandler) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId, err := extractTaskId(r)
	if err != nil {
		ErrJSON(w, httpErr{err: err, msg: "not found"}, http.StatusNotFound)
		return
	}

	tr := taskRequest{}
	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		ErrJSON(w, httpErr{err: err, msg: "bad request"}, http.StatusBadRequest)
		return
	}
	th.mu.Lock()
	defer th.mu.Unlock()
	task, ok := th.data[taskId]
	if !ok {
		ErrJSON(w, httpErr{err: fmt.Errorf("invalid task"), msg: "not found"}, http.StatusNotFound)
		return
	}

	nw := time.Now()
	task.Title = tr.Title
	task.Details = tr.Details
	if tr.Status != "" && tr.Status != task.Status {

		task.Comments = append(task.Comments, Comment{
			Body:     fmt.Sprintf("Status changed from %s to %s", task.Status, tr.Status),
			Username: "ADMIN",
			Tsp:      nw,
		})
		task.Status = tr.Status
	}
	task.UpdatedAt = nw

	th.data[task.ID] = task

	w.Header().Set("Location", fmt.Sprintf("/api/task/%s", task.ID))
	w.Header().Set("X-ViewLocation", fmt.Sprintf("/view/%s", task.ID))
	JSON(w, task, http.StatusOK)

}

type commentRequest struct {
	TaskID string `json:"taskId"`
	Body   string `json:"body"`
}

// addComment: POST /api/comments -> adds a comment to a task.
func (th *TasksHandler) addComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := http.StatusText(http.StatusNotFound)
		ErrJSON(w, httpErr{err: fmt.Errorf(msg), msg: msg}, http.StatusNotFound)
		return
	}
	cr := commentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		ErrJSON(w, httpErr{err: err, msg: "bad request"}, http.StatusBadRequest)
		return
	}
	th.mu.Lock()
	defer th.mu.Unlock()
	task, ok := th.data[cr.TaskID]
	if !ok {
		ErrJSON(w, httpErr{err: fmt.Errorf("invalid task"), msg: "not found"}, http.StatusNotFound)
		return
	}

	task.Comments = append(task.Comments, Comment{
		Body:     cr.Body,
		Username: "ADMIN",
		Tsp:      time.Now(),
	})

	th.data[cr.TaskID] = task
	w.Header().Set("Location", fmt.Sprintf("/api/task/%s", task.ID))
	w.Header().Set("X-ViewLocation", fmt.Sprintf("/view/%s", task.ID))
	JSON(w, task, http.StatusOK)
}
