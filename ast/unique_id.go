package ast

import (
	"fmt"
)

var uniqueCounter int = 0

func NextUniqueId() Identifier {
	uniqueCounter++

	return Identifier{fmt.Sprintf("`unique_%d`", uniqueCounter)}
}
