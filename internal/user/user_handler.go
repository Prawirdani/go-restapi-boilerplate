package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (h *UserHandler) Routes(r chi.Router) {
	r.Get("/users", h.List)
	r.Get("/users/{userId}", h.ById)
	r.Post("/users", h.Create)
}

//	@Summary		List Users
//	@Description	Returning list of users
//	@Produce		json
//	@Tags			Users
//	@Success		200		{object}	httputil.Response{data=[]user.User}
//	@Failure		default	{object}	httputil.ErrorResponse	"400 & 500 status, error field can be string or object"
//	@Router			/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.FindAll(r.Context())
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJson(w, http.StatusOK, users)
}

//	@Summary		User by id
//	@Description	Returning an user object with given id
//	@Produce		json
//	@Param			userId	path	string	true	"user_id (ULID)"
//	@Tags			Users
//	@Success		200		{object}	httputil.Response{data=user.User}
//	@Failure		default	{object}	httputil.ErrorResponse	"400 & 500 status, error field can be string or object"
//	@Router			/users/{userId} [get]
func (h *UserHandler) ById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")


	user, err := h.userService.FindById(r.Context(), userId)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJson(w, http.StatusOK, user)
}

//	@Summary		Create User
//	@Description	Create New User
//	@Accept			json
//	@Param			User	body	user.User	true	"please exclude non-required fields before firing the request"
//	@Produce		json
//	@Tags			Users
//	@Success		201		{object}	httputil.Response
//	@Failure		default	{object}	httputil.ErrorResponse	"400 & 500 status, error field can be string or object"
//	@Router			/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody User
	if err := httputil.BindJson(r, &reqBody); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.userService.Save(r.Context(), reqBody); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJson(w, 201, "Success Create User")
}
