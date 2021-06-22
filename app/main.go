package main

import (
	_ "errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"

	"database/sql"
	"fmt"
	"forums/config"
	"forums/internal/forum"
	forumDelivery "forums/internal/forum/delivery"
	forumRepo "forums/internal/forum/repository"
	forumUcase "forums/internal/forum/usecase"
	"forums/internal/thread"
	threadDelivery "forums/internal/thread/delivery"
	threadRepo "forums/internal/thread/repository"
	threadUcase "forums/internal/thread/usecase"
	"forums/internal/user"
	userDelivery "forums/internal/user/delivery/http"
	userRepo "forums/internal/user/repository"
	userUcase "forums/internal/user/usecase"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

type Router struct {
	user   user.UserHandler
	forum  forum.ForumHandler
	thread thread.ThreadHandler
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("paramValues = ", mux.Vars(r), "\nPath = ", r.RequestURI, "\nBody", r.Body)
	return
}

func newRouter(r Router) *mux.Router {
	router := mux.NewRouter()
	router.Use(LogMiddleware)

	_user := router.PathPrefix("/api/user").Subrouter()
	_user.HandleFunc("/{nickname}/create", r.user.Create).Methods(http.MethodPost)
	_user.HandleFunc("/{nickname}/profile", r.user.GetUserData).Methods(http.MethodGet)
	_user.HandleFunc("/{nickname}/profile", r.user.UpdateUserData).Methods(http.MethodPost)

	_forum := router.PathPrefix("/api/forum").Subrouter()
	_forum.HandleFunc("/create", r.forum.CreateForum).Methods(http.MethodPost)
	_forum.HandleFunc("/{slug}/details", r.forum.GetForumBySlug).Methods(http.MethodGet)
	_forum.HandleFunc("/{slug}/create", r.forum.CreateThread).Methods(http.MethodPost)
	_forum.HandleFunc("/{slug}/users", handler).Methods(http.MethodGet)
	_forum.HandleFunc("/{slug}/threads", r.forum.GetThreadsInForum).Methods(http.MethodGet)

	_post := router.PathPrefix("/api/post").Subrouter()
	_post.HandleFunc("/{id}/details", handler).Methods(http.MethodGet)
	_post.HandleFunc("/{id}/details", handler).Methods(http.MethodPost)

	_service := router.PathPrefix("/api/service").Subrouter()
	_service.HandleFunc("/clear", handler).Methods(http.MethodPost)
	_service.HandleFunc("/status", handler).Methods(http.MethodGet)

	_thread := router.PathPrefix("/api/thread").Subrouter()
	_thread.HandleFunc("/{slugOrId}/create", r.thread.AddPosts).Methods(http.MethodPost)
	_thread.HandleFunc("/{slugOrId}/details", r.thread.GetThread).Methods(http.MethodGet)
	_thread.HandleFunc("/{slugOrId}/details", handler).Methods(http.MethodPost)
	_thread.HandleFunc("/{slugOrId}/posts", handler).Methods(http.MethodGet)
	_thread.HandleFunc("/{slugOrId}/vote", r.thread.Vote).Methods(http.MethodPost)

	return router
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}\n uri=${uri}\n status=${status}\n time=${time_rfc3339_nano}\n\n",
	}))

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open(config.PostgresDB, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userRepo_ := userRepo.NewUserRepo(db)
	userUcase_ := userUcase.NewUserUsecase(userRepo_)
	userHandler_ := userDelivery.NewUserHandler(userUcase_)

	forumRepo_ := forumRepo.NewForumRepo(db, userRepo_)
	forumUcase_ := forumUcase.NewForumUsecase(forumRepo_)
	forumHandler_ := forumDelivery.NewForumHandler(forumUcase_)

	threadRepo_ := threadRepo.NewThreadRepo(db, forumRepo_)
	threadUcase_ := threadUcase.NewThreadUsecase(threadRepo_)
	threadHandler_ := threadDelivery.NewThreadHandler(threadUcase_)

	routes := Router{
		user:   userHandler_,
		forum:  forumHandler_,
		thread: threadHandler_,
	}

	router := newRouter(routes)

	server := &http.Server{
		Handler: router,
		Addr:    ":5000",
	}

	fmt.Println("start server at ", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("cant start this, error: ", err)
	}
}
