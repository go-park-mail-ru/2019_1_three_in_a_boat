package game

import "errors"

var ErrWrongState = errors.New("the method can not be called with this state")
