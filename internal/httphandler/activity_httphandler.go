package httphandler

import (
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

type ActivityHttpHandler struct {
	Command usecase.ActivityUsecaseCommand
	Query   usecase.ActivityUsecaseQuery
	Decoder rdecoder.Decoder
}

func NewActivityHttpHandler(prop HTTPHandlerProperty) http.Handler {
	props := usecase.ActivityUsecaseProps{
		Repository: repository.NewActivityRepositoryMysql(prop.DB),
	}

	handler := ActivityHttpHandler{
		Command: usecase.NewActivityCommandUsecase(props),
		Query:   usecase.NewActivityQueryUsecase(props),
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

func (a ActivityHttpHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var payload model.Activity

	err = rdecoder.DecodeRest(r, a.Decoder, &payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "try to decode data",
				"at":     "ActivityHttpHandler",
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
				"at":     "ActivityHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	act, err := a.Command.Create(ctx, payload)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "Command.Create",
				"at":     "ActivityHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccessCreate(act).Write(w, a.Decoder)
}

func (a ActivityHttpHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var payload model.Activity

	id := chihelper.ChiURLParamToInt64(r, "id", 0)
	if id < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ActivityHttpHandler",
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
				"at":     "ActivityHttpHandler",
			}).
			Error(err)

		err = rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"Invalid body",
			"ActivityHttpHandler",
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
				"at":     "ActivityHttpHandler",
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
				"at":     "ActivityHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(act).Write(w, a.Decoder)
}

func (a ActivityHttpHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	id := chihelper.ChiURLParamToInt64(r, "id", 0)
	if id < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ActivityHttpHandler",
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
				"at":     "ActivityHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(map[string]string{}).Write(w, a.Decoder)
}

func (a ActivityHttpHandler) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	id := chihelper.ChiURLParamToInt64(r, "id", 0)
	if id < 1 {
		err = rapperror.ErrUnprocessableEntity(
			rapperror.AppErrorCodeBadRequest,
			"Invalid id",
			"ActivityHttpHandler",
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
				"at":     "ActivityHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(act).Write(w, a.Decoder)
}

func (a ActivityHttpHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	acts, err := a.Query.Get(ctx)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"action": "Command.Create",
				"at":     "ActivityHttpHandler",
			}).
			Error(err)

		responser.NewResponserErr(err).Write(w, a.Decoder)
		return
	}

	responser.NewResponseSuccess(acts).Write(w, a.Decoder)
}
