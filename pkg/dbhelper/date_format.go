package dbhelper

import "time"

func ToSqlFormat(dt string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", dt)

	return t
}

func ToSqlFormatFromEpoch(dt int64) time.Time {
	return time.Unix(0, dt*int64(time.Millisecond))
}
