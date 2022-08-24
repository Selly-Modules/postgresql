package postgresql

import "github.com/volatiletech/null/v8"

func NewValidNullString(val string) null.String {
	return null.String{
		String: val,
		Valid:  true,
	}
}
