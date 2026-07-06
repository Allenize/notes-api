package store

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "data.json"))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	return s
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func TestCreateAndGetUser(t *testing.T) {
	s := newTestStore(t)

	u, err := s.CreateUser("alice", "hashedpw")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	got, err := s.GetUserByUsername("alice")
	if err != nil {
		t.Fatalf("GetUserByUsername returned error: %v", err)
	}
	if got.ID != u.ID {
		t.Fatalf("expected user ID %s, got %s", u.ID, got.ID)
	}
}

func TestCreateUserConflict(t *testing.T) {
	s := newTestStore(t)

	if _, err := s.CreateUser("bob", "hash1"); err != nil {
		t.Fatalf("first CreateUser returned error: %v", err)
	}
	if _, err := s.CreateUser("bob", "hash2"); err != ErrConflict {
		t.Fatalf("expected ErrConflict for duplicate username, got %v", err)
	}
}

func TestTaskCRUD(t *testing.T) {
	s := newTestStore(t)
	u, _ := s.CreateUser("carol", "hash")

	task, err := s.CreateTask(u.ID, "Buy milk", "2%", "Errands", []string{"grocery"})
	if err != nil {
		t.Fatalf("CreateTask returned error: %v", err)
	}
	if task.Category != "Errands" {
		t.Fatalf("expected category Errands, got %q", task.Category)
	}

	tasks := s.ListTasks(u.ID)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	newTitle := "Buy oat milk"
	updated, err := s.UpdateTask(u.ID, task.ID, TaskUpdate{Title: strPtr(newTitle)})
	if err != nil {
		t.Fatalf("UpdateTask returned error: %v", err)
	}
	if updated.Title != newTitle {
		t.Fatalf("expected title %q, got %q", newTitle, updated.Title)
	}

	if err := s.DeleteTask(u.ID, task.ID); err != nil {
		t.Fatalf("DeleteTask returned error: %v", err)
	}

	if _, err := s.GetTask(u.ID, task.ID); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestTaskIsolatedByUser(t *testing.T) {
	s := newTestStore(t)
	u1, _ := s.CreateUser("dave", "hash")
	u2, _ := s.CreateUser("erin", "hash")

	task, _ := s.CreateTask(u1.ID, "Dave's task", "", "", nil)

	if _, err := s.GetTask(u2.ID, task.ID); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound when a different user accesses the task, got %v", err)
	}
}

func TestStorePersistsAcrossReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")

	s1, _ := New(path)
	u, _ := s1.CreateUser("frank", "hash")
	s1.CreateTask(u.ID, "Persisted task", "", "", nil)

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected data file to exist: %v", err)
	}

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reloading store returned error: %v", err)
	}

	got, err := s2.GetUserByUsername("frank")
	if err != nil {
		t.Fatalf("expected user to persist across reload: %v", err)
	}

	tasks := s2.ListTasks(got.ID)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 persisted task, got %d", len(tasks))
	}
}

func TestTaskFavoritePinArchive(t *testing.T) {
	s := newTestStore(t)
	u, _ := s.CreateUser("gina", "hash")
	task, _ := s.CreateTask(u.ID, "Idea", "", "", nil)

	updated, err := s.UpdateTask(u.ID, task.ID, TaskUpdate{Favorite: boolPtr(true), Pinned: boolPtr(true)})
	if err != nil {
		t.Fatalf("UpdateTask returned error: %v", err)
	}
	if !updated.Favorite || !updated.Pinned {
		t.Fatalf("expected favorite and pinned to be true, got favorite=%v pinned=%v", updated.Favorite, updated.Pinned)
	}

	archived, err := s.UpdateTask(u.ID, task.ID, TaskUpdate{Archived: boolPtr(true)})
	if err != nil {
		t.Fatalf("UpdateTask (archive) returned error: %v", err)
	}
	if !archived.Archived {
		t.Fatal("expected archived to be true")
	}
}

func TestTaskTrashAndRestore(t *testing.T) {
	s := newTestStore(t)
	u, _ := s.CreateUser("hank", "hash")
	task, _ := s.CreateTask(u.ID, "Old note", "", "", nil)

	trashed, err := s.TrashTask(u.ID, task.ID)
	if err != nil {
		t.Fatalf("TrashTask returned error: %v", err)
	}
	if !trashed.Trashed || trashed.TrashedAt == nil {
		t.Fatal("expected task to be marked trashed with a timestamp")
	}

	restored, err := s.RestoreTask(u.ID, task.ID)
	if err != nil {
		t.Fatalf("RestoreTask returned error: %v", err)
	}
	if restored.Trashed || restored.TrashedAt != nil {
		t.Fatal("expected task to be restored (not trashed)")
	}
}
