package emailsrepo

import "errors"

var DatabaseError = errors.New("Couldn't reach the database")
var NotFoundError = errors.New("Not found")
