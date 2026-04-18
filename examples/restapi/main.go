package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rluders/httpsuite/v3"
	"github.com/rluders/httpsuite/validation/playground"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin member viewer"`
}

type GetUserRequest struct {
	ID int `json:"id"`
}

func (r *GetUserRequest) SetParam(fieldName, value string) error {
	if fieldName != "id" {
		return nil
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	r.ID = id
	return nil
}

type UserStore struct {
	mu     sync.RWMutex
	nextID int
	users  map[int]User
}

func NewUserStore() *UserStore {
	return &UserStore{
		nextID: 3,
		users: map[int]User{
			1: {
				ID:        1,
				Name:      "Ada Lovelace",
				Email:     "ada@example.com",
				Role:      "admin",
				CreatedAt: time.Now().Add(-72 * time.Hour),
			},
			2: {
				ID:        2,
				Name:      "Grace Hopper",
				Email:     "grace@example.com",
				Role:      "member",
				CreatedAt: time.Now().Add(-48 * time.Hour),
			},
		},
	}
}

func (s *UserStore) Create(req *CreateUserRequest) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:        s.nextID,
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: time.Now().UTC(),
	}
	s.users[user.ID] = user
	s.nextID++
	return user
}

func (s *UserStore) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	return user, ok
}

func (s *UserStore) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	return users
}

func main() {
	store := NewUserStore()
	httpsuite.SetValidator(playground.NewWithValidator(nil, &httpsuite.ProblemConfig{
		BaseURL: "http://localhost:8080",
	}))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		page := readPositiveInt(r, "page", 1)
		pageSize := readPositiveInt(r, "page_size", 2)

		users := store.List()
		start, end := clampPageWindow(page, pageSize, len(users))

		httpsuite.Reply().
			Meta(httpsuite.NewPageMeta(page, pageSize, len(users))).
			OK(w, users[start:end])
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		req, err := httpsuite.ParseRequest[*GetUserRequest](w, r, chi.URLParam, nil, "id")
		if err != nil {
			return
		}

		user, ok := store.Get(req.ID)
		if !ok {
			problem := httpsuite.ProblemNotFound(fmt.Sprintf("User %d does not exist", req.ID)).
				Title("User Not Found").
				Instance(r.URL.Path).
				Extension("request_id", middleware.GetReqID(r.Context())).
				Build()
			httpsuite.ProblemResponse(w, problem)
			return
		}

		httpsuite.OK(w, user)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		req, err := httpsuite.ParseRequest[*CreateUserRequest](w, r, nilParamExtractor, &httpsuite.ParseOptions{
			MaxBodyBytes: 2 << 10,
		})
		if err != nil {
			return
		}

		user := store.Create(req)
		httpsuite.Reply().
			Created(w, user, fmt.Sprintf("/users/%d", user.ID))
	})

	r.Get("/feed", func(w http.ResponseWriter, r *http.Request) {
		cursor := strings.TrimSpace(r.URL.Query().Get("cursor"))
		items := []map[string]any{
			{"event": "user.created", "resource_id": 1},
			{"event": "user.updated", "resource_id": 2},
		}

		meta := httpsuite.NewCursorMeta("next-page-token", cursor, true, cursor != "")
		httpsuite.Reply().
			Meta(meta).
			OK(w, items)
	})

	log.Println("Starting REST API example on :8080")
	log.Println("POST /users")
	log.Println("GET  /users?page=1&page_size=2")
	log.Println("GET  /users/{id}")
	log.Println("GET  /feed?cursor=next-page-token")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func readPositiveInt(r *http.Request, key string, fallback int) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func clampPageWindow(page, pageSize, total int) (int, int) {
	if total <= 0 {
		return 0, 0
	}
	if page <= 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = total
	}

	totalPages := 1 + (total-1)/pageSize
	if page > totalPages {
		return total, total
	}

	start := (page - 1) * pageSize
	if start > total {
		start = total
	}

	end := total
	if remaining := total - start; pageSize < remaining {
		end = start + pageSize
	}

	return start, end
}

func nilParamExtractor(*http.Request, string) string {
	return ""
}
