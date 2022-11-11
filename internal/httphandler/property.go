package httphandler

import (
	"database/sql"

	"gitlab.com/rteja-library3/rdecoder"
)

type HTTPHandlerProperty struct {
	DB             *sql.DB
	DefaultDecoder rdecoder.Decoder
}
