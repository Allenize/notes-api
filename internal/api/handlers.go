package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Allenize/notes-api/internal/auth"
	"github.com/Allenize/notes-api/internal/store"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type Handlers struct {
	store  *store.Store
	secret []byte
}

func New(s *store.Store, secret []byte) *Handlers {
	return &Handlers{store: s, secret: secret}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type tokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (h *Handlers) issuePair(userID string) (tokenPair, error) {
	access, err := auth.IssueToken(userID, "access", h.secret, accessTokenTTL)
	if err != nil {
		return tokenPair{}, err
	}
	refresh, err := auth.IssueToken(userID, "refresh", h.secret, refreshTokenTTL)
	if err != nil {
		return tokenPair{}, err
	}
	return tokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(accessTokenTTL.Seconds()),
	}, nil
}

type signupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 {
		writeErr(w, http.StatusBadRequest, "username must be at least 3 characters")
		return
	}
	if len(req.Password) < 8 {
		writeErr(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not process password")
		return
	}

	user, err := h.store.CreateUser(req.Username, hash)
	if err == store.ErrConflict {
		writeErr(w, http.StatusConflict, "username already taken")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not create user")
		return
	}

	pair, err := h.issuePair(user.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not issue tokens")
		return
	}
	writeJSON(w, http.StatusCreated, pair)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := h.store.GetUserByUsername(strings.TrimSpace(req.Username))
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		writeErr(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	pair, err := h.issuePair(user.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not issue tokens")
		return
	}
	writeJSON(w, http.StatusOK, pair)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	userID, tokenType, err := auth.VerifyToken(req.RefreshToken, h.secret)
	if err != nil || tokenType != "refresh" {
		writeErr(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	pair, err := h.issuePair(userID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not issue tokens")
		return
	}
	writeJSON(w, http.StatusOK, pair)
}

type taskRequest struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		writeErr(w, http.StatusBadRequest, "title is required")
		return
	}

	task, err := h.store.CreateTask(userID, req.Title, req.Notes)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not create task")
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

type taskListResponse struct {
	Tasks  []*store.Task `json:"tasks"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	tasks := h.store.ListTasks(userID)

	q := r.URL.Query()

	if doneStr := q.Get("done"); doneStr != "" {
		want := doneStr == "true"
		filtered := tasks[:0:0]
		for _, t := range tasks {
			if t.Done == want {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	if search := strings.TrimSpace(q.Get("q")); search != "" {
		search = strings.ToLower(search)
		filtered := tasks[:0:0]
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), search) || strings.Contains(strings.ToLower(t.Notes), search) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	total := len(tasks)

	limit := 20
	if l, err := strconv.Atoi(q.Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(q.Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	end := offset + limit
	if offset > total {
		offset = total
	}
	if end > total {
		end = total
	}
	page := tasks[offset:end]

	writeJSON(w, http.StatusOK, taskListResponse{
		Tasks:  page,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	taskID := r.PathValue("id")

	task, err := h.store.GetTask(userID, taskID)
	if err == store.ErrNotFound {
		writeErr(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not fetch task")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

type updateTaskRequest struct {
	Title *string `json:"title"`
	Notes *string `json:"notes"`
	Done  *bool   `json:"done"`
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	taskID := r.PathValue("id")

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	task, err := h.store.UpdateTask(userID, taskID, req.Title, req.Notes, req.Done)
	if err == store.ErrNotFound {
		writeErr(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not update task")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	taskID := r.PathValue("id")

	err := h.store.DeleteTask(userID, taskID)
	if err == store.ErrNotFound {
		writeErr(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not delete task")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
