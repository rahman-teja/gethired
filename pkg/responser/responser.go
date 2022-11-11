package responser

import (
	"net/http"

	"gitlab.com/rteja-library3/rapperror"
	"gitlab.com/rteja-library3/rdecoder"
)

type Responser struct {
	HttpCode int         `json:"-"`
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}

func NewResponserErr(err error) Responser {
	appErr, ok := err.(*rapperror.AppError)
	if !ok {
		return NewResponser(http.StatusInternalServerError, "Internal Server", err.Error(), nil)
	}

	status := "Internal Server"
	switch appErr.Status {
	case http.StatusNotFound:
		status = "Not Found"
	case http.StatusConflict:
		status = "Conflict"
	case http.StatusBadRequest:
		status = "Bad Request"
	}

	return NewResponser(appErr.Status, status, appErr.Message, nil)
}

func NewResponseSuccess(data interface{}) Responser {
	return NewResponser(http.StatusOK, "Success", "Success", data)
}

func NewResponseSuccessCreate(data interface{}) Responser {
	return NewResponser(http.StatusCreated, "Success", "Success", data)
}

func NewResponser(code int, status, message string, data interface{}) Responser {
	return Responser{
		HttpCode: code,
		Status:   status,
		Message:  message,
		Data:     data,
	}
}

func (r Responser) Write(w http.ResponseWriter, enc rdecoder.Encode) (err error) {
	w.Header().Set("Content-type", enc.ContentType())
	w.WriteHeader(r.HttpCode)

	err = enc.Encode(w, r)
	return
}
