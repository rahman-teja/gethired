package chihelper

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ChiURLParamToString(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func ChiURLParamToInt64(r *http.Request, key string, def int64) (res int64) {
	id := ChiURLParamToString(r, key)

	res = ToInt64(id, def)
	return
}

func ToInt64(from string, def int64) (res int64) {
	var err error

	res, err = strconv.ParseInt(from, 10, 64)
	if err != nil {
		res = def
		return
	}

	return
}
