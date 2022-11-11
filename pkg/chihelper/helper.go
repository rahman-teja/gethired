package chihelper

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rahman-teja/rhelper"
)

func ChiURLParamToString(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func ChiURLParamToInt64(r *http.Request, key string, def int64) (res int64) {
	id := ChiURLParamToString(r, key)

	res = ToInt64(id, def)
	return
}

func QueryStringToPointerInt64(r *http.Request, key string, def int64) *int64 {
	s := rhelper.QueryString(r, key)
	if s == "" {
		return nil
	}

	val := ToInt64(s, def)
	return &val
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
