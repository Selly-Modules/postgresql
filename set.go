package postgresql

import (
	"github.com/volatiletech/null/v8"
	"time"
)

func SetString(val string) null.String {
	return null.String{
		String: val,
		Valid:  true,
	}
}

func SetInt(val int) null.Int {
	return null.Int{
		Int:   val,
		Valid: true,
	}
}

func SetBool(val bool) null.Bool {
	return null.Bool{
		Bool:  val,
		Valid: true,
	}
}

func SetTime(val time.Time) null.Time {
	return null.Time{
		Time:  val,
		Valid: true,
	}
}

func SetJSON(val []byte) null.JSON {
	return null.JSON{
		JSON:  val,
		Valid: true,
	}
}
