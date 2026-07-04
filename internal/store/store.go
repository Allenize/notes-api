package store

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("already exists")

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type Task struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Notes     string    `json:"notes"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type snapshot struct {
	Users map[string]*User `json:"users"`
	Tasks map[string]*Task `json:"tasks"`
}

type Store struct {
	mu       sync.Mutex
	path     string
	users    map[string]*User
	tasks    map[string]*Task
	usersIdx map[string]string
}

func New(path string) (*Store, error) {
	s := &Store{
		path:     path,
		users:    make(map[string]*User),
		tasks:    make(map[string]*Task),
		usersIdx: make(map[string]string),
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func newID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var snap snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return err
	}
	if snap.Users != nil {
		s.users = snap.Users
	}
	if snap.Tasks != nil {
		s.tasks = snap.Tasks
	}
	for id, u := range s.users {
		s.usersIdx[u.Username] = id
	}
	return nil
}

func (s *Store) saveLocked() error {
	snap := snapshot{Users: s.users, Tasks: s.tasks}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

func (s *Store) CreateUser(username, passwordHash string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersIdx[username]; exists {
		return nil, ErrConflict
	}

	u := &User{
		ID:           newID(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	s.users[u.ID] = u
	s.usersIdx[username] = u.ID

	if err := s.saveLocked(); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, ok := s.usersIdx[username]
	if !ok {
		return nil, ErrNotFound
	}
	return s.users[id], nil
}

func (s *Store) CreateTask(userID, title, notes string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	t := &Task{
		ID:        newID(),
		UserID:    userID,
		Title:     title,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.tasks[t.ID] = t

	if err := s.saveLocked(); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Store) ListTasks(userID string) []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]*Task, 0)
	for _, t := range s.tasks {
		if t.UserID == userID {
			result = append(result, t)
		}
	}
	return result
}

func (s *Store) GetTask(userID, taskID string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[taskID]
	if !ok || t.UserID != userID {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *Store) UpdateTask(userID, taskID string, title, notes *string, done *bool) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[taskID]
	if !ok || t.UserID != userID {
		return nil, ErrNotFound
	}

	if title != nil {
		t.Title = *title
	}
	if notes != nil {
		t.Notes = *notes
	}
	if done != nil {
		t.Done = *done
	}
	t.UpdatedAt = time.Now()

	if err := s.saveLocked(); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Store) DeleteTask(userID, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[taskID]
	if !ok || t.UserID != userID {
		return ErrNotFound
	}
	delete(s.tasks, taskID)

	return s.saveLocked()
}
