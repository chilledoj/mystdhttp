package router

import (
	"crypto/rand"
	"fmt"
	"time"
)

type Comment struct {
	Body     string    `json:"body"`
	Username string    `json:"username"`
	Tsp      time.Time `json:"tsp"`
}

const (
	StatusInactive  = "inactive"
	StatusNew       = "new"
	StatusStarted   = "started"
	StatusCompleted = "completed"
	StatusReopened  = "reopened"
)

type Task struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Details    string    `json:"details"`
	Status     string    `json:"status"`
	Deleted    bool      `json:"-"`
	CreatedAt  time.Time `json:"crtAt"`
	UpdatedAt  time.Time `json:"updAt"`
	Creator    string    `json:"crtUsr"`
	Updater    string    `json:"updUsr"`
	AssignedTo string    `json:"assignedTo"`
	Comments   []Comment `json:"comments,omitempty"`
}

// AllowedStatuses: This provides the possible statuses the Task can be moved into.
func (t Task) AllowedStatuses() []string {
	if t.Status == StatusInactive {
		return []string{StatusInactive, StatusReopened}
	} else if t.Status == StatusReopened {
		return []string{StatusReopened, StatusStarted, StatusCompleted, StatusInactive}
	} else if t.Status == StatusCompleted {
		return []string{StatusCompleted, StatusInactive, StatusReopened}
	} else if t.Status == StatusStarted {
		return []string{StatusStarted, StatusCompleted, StatusInactive}
	}
	return []string{StatusNew, StatusStarted, StatusCompleted, StatusInactive}
}

// Fake ID generator
func genId() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func NewTask(title, details, status string) Task {
	st := "new"
	if status != "" {
		st = status
	}
	return Task{
		ID:         genId(),
		Title:      title,
		Details:    details,
		Status:     st,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Creator:    "ADMIN",
		Updater:    "ADMIN",
		AssignedTo: "ADMIN",
	}
}
