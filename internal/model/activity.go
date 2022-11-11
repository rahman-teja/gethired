package model

import "gitlab.com/rteja-library3/rapperror"

type Activity struct {
	Email *string `json:"email"`
	Title *string `json:"title"`
}

func (t Activity) Validate() error {
	if t.Title == nil {
		return rapperror.ErrBadRequest(
			rapperror.AppErrorCodeBadRequest,
			"title cannot be null",
			"",
			nil,
		)
	}

	return nil
}
