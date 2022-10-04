package null

import (
	"github.com/volatiletech/null/v8"
	"time"
)

func GetString(val null.String) string {
	if !val.Valid {
		return ""
	}
	return val.String
}

func GetInt(val null.Int) int {
	if !val.Valid {
		return 0
	}
	return val.Int
}

func GetBool(val null.Bool) bool {
	if !val.Valid {
		return false
	}
	return val.Bool
}

func GetTime(val null.Time) (res time.Time) {
	if !val.Valid {
		return
	}
	return val.Time
}
