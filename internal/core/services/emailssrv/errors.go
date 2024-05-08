package emailssrvc

import (
	"errors"
)

var NoStockError = errors.New("The company does not offer cars with 3 doors nor purple cars")
