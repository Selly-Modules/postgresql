package postgresql

import (
	"github.com/Masterminds/squirrel"
)

// GetStmBuilder ...
func GetStmBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
