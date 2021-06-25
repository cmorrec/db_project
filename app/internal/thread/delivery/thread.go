package delivery

import (
	"encoding/json"
	"net/http"

	forumModel "github.com/forums/app/internal/forum"
	threadModel "github.com/forums/app/internal/thread"
	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/response"
	"github.com/gorilla/mux"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type Handler struct {
	threadRepo threadModel.ThreadRepo
	userRepo   userModel.UserRepo
	forumRepo  forumModel.ForumRepo
}

func NewThreadHandler(threadRepo threadModel.ThreadRepo, userRepo userModel.UserRepo,
	forumRepo forumModel.ForumRepo) threadModel.ThreadHandler {
	return &Handler{
		threadRepo: threadRepo,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
	}
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slug := vars["slug"]
	newThread := new(models.Thread)
	err := json.NewDecoder(r.Body).Decode(&newThread)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	user, err := h.userRepo.GetUserByName(ctx, newThread.Author)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	forum, err := h.forumRepo.GetForumBySlug(ctx, slug)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if user == nil {
		message := models.Message{
			Message: "Can't find user with id #" + newThread.Author + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}
	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + slug + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	oldThread, err := h.threadRepo.GetThreadBySlugOrId(ctx, newThread.Slug)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if newThread.Slug != "" && oldThread != nil {
		response.New(http.StatusConflict, oldThread).SendSuccess(w)
		return
	}

	newThread.Forum = forum.Slug
	id, err := h.threadRepo.CreateThread(ctx, newThread)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	newThread.Id = id
	response.New(http.StatusCreated, newThread).SendSuccess(w)
}

func (h *Handler) GetDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]

	thread, err := h.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if thread == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slugOrId + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	response.New(http.StatusOK, thread).SendSuccess(w)
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	newThread := new(models.Thread)
	err := json.NewDecoder(r.Body).Decode(&newThread)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	threadOld, err := h.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if threadOld == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slugOrId + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	if newThread.Title != "" {
		threadOld.Title = newThread.Title
	}

	if newThread.Message != "" {
		threadOld.Message = newThread.Message
	}

	err = h.threadRepo.UpdateThreadBySlug(ctx, threadOld)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	response.New(http.StatusOK, threadOld).SendSuccess(w)
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	threadPosts := new(models.ThreadPosts)

	threadPosts.SlugOrId = vars["slug_or_id"]
	threadPosts.Limit = r.URL.Query().Get("limit")
	threadPosts.Since = r.URL.Query().Get("since")
	threadPosts.Sort = r.URL.Query().Get("sort")

	desc := r.URL.Query().Get("desc")
	if desc == "false" || desc == "" {
		threadPosts.Desc = false
	} else {
		threadPosts.Desc = true
	}

	thread, err := h.threadRepo.GetThreadBySlugOrId(ctx, threadPosts.SlugOrId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if thread == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + threadPosts.SlugOrId + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	threadPosts.ThreadId = thread.Id
	posts, err := h.threadRepo.GetPosts(ctx, threadPosts)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	response.New(http.StatusOK, posts).SendSuccess(w)
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	vote := new(models.Vote)
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer r.Body.Close()

	thread, err := h.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if thread == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slugOrId + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	vote.Thread = thread.Id
	err = h.threadRepo.AddVote(ctx, vote)
	if err != nil {
		if pqErr, ok := err.(pgx.PgError); ok {
			switch pqErr.Code {
			case pgerrcode.ForeignKeyViolation:
				message := models.Message{
					Message: "Can't find user with id #" + vote.User + "\n",
				}
				response.New(http.StatusNotFound, message).SendSuccess(w)
				return

			case pgerrcode.UniqueViolation:
				err = h.threadRepo.UpdateVote(ctx, vote)
				if err != nil {
					w.WriteHeader(500)
					return
				}

				thread, err = h.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
				if err != nil {
					w.WriteHeader(500)
					return
				}

			default:
				w.WriteHeader(500)
				return
			}
		}
	} else {
		thread.Votes += vote.Voice
	}

	response.New(http.StatusOK, thread).SendSuccess(w)
}
