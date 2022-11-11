package httphandler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rahman-teja/gethired/internal/model"
	"github.com/rahman-teja/gethired/internal/repository"
	"github.com/rahman-teja/gethired/internal/usecase"
	"github.com/rahman-teja/gethired/pkg/chihelper"
	"github.com/rahman-teja/gethired/pkg/responser"
	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
)

type ToDoHttpHandler struct {
	Command usecase.ToDoUsecaseCommand
	Query   usecase.ToDoUsecaseQuery
	Decoder rdecoder.Decoder
}

func NewToDoHttpHandler(prop HTTPHandlerProperty) http.Handler {
	props := usecase.ToDoUsecaseProps{
		Repository: repository.NewToDoRepositoryMysql(prop.DB),
	}

	handler := ToDoHttpHandler{
		Command: usecase.NewToDoCommandUsecase(props),
		Query:   usecase.NewToDoQueryUsecase(props),
		Decoder: prop.DefaultDecoder,
	}

	r := chi.NewRouter()

	r.Get("/", handler.GetHandler)
	r.Get("/{id}", handler.GetOneHandler)

	r.Post("/", handler.CreateHandler)
	r.Patch("/{id}", handler.UpdateHandler)
	r.Delete("/{id}", handler.DeleteHandler)

	return r
}

func (a ToDoHttpHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var payload model.ToDo

	err = rdecoder.DecodeRest(r, a.Decoder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"",
			nil,
		)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	err = payload.Validate()
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "payload.Validate",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	if payload.ActivityGroupID == nil {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"activity_group_id cannot be null",
			"",
			nil,
		)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	if payload.Title == nil {
		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"title cannot be null",
			"",
			nil,
		)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	act, err := a.Command.Create(ctx, payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "Command.Create",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccessCreate(act).Write(w, a.Decoder)
}

func (a ToDoHttpHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var payload model.ToDo

	id := chihelper.ChiURLParamToInt64(r, "id", 0)
	if id < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ToDoHttpHandler",
			nil,
		)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	err = rdecoder.DecodeRest(r, a.Decoder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ToDoHttpHandler",
			nil,
		)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	logrus.Info(fmt.Sprintf("data %+v", payload))

	err = payload.Validate()
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "payload.Validate",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	act, err := a.Command.Update(ctx, id, payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "Command.Create",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(act).Write(w, a.Decoder)
}

func (a ToDoHttpHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	id := chihelper.ChiURLParamToInt64(r, "id", 0)
	if id < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ToDoHttpHandler",
			nil,
		)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	err = a.Command.Delete(ctx, id)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "Command.Create",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(map[string]string{}).Write(w, a.Decoder)
}

func (a ToDoHttpHandler) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	id := chihelper.ChiURLParamToInt64(r, "id", 0)
	if id < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ToDoHttpHandler",
			nil,
		)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	act, err := a.Query.GetOne(ctx, id)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "Command.Create",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(act).Write(w, a.Decoder)
}

func (a ToDoHttpHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	actId := chihelper.QueryStringToPointerInt64(r, "activity_group_id", 0)

	filter := repository.ToDoQueryRepositoryFilter{
		ActivityId: actId,
	}

	acts, err := a.Query.Get(ctx, filter)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "Command.Create",
				"at":     "ToDoHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(acts).Write(w, a.Decoder)
}
