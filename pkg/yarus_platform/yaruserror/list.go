package yaruserror

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ErrNotSet is error for case when some thing is not set
var ErrNotSet error = errors.New("Not set.")

// ErrNotFound is error for case when entity not found
var ErrNotFound error = errors.New("Not found.")

// ErrAlreadyExists is error for case when some value is already exists
var ErrAlreadyExists error = errors.New("Already exists.")

// ErrEmptyParams is error for case when params are empty
var ErrEmptyParams error = errors.New("Empty params.")

// ErrBadRequest is error for case when bad request
var ErrBadParams error = errors.New("Bad params.")

// ErrBadRequest is error for case when bad request
var ErrBadRequest error = errors.New("Bad request.")

// ErrInternal is error for case when smth went wrong
var ErrInternal error = errors.New("Internal error.")

var ErrTokenHasExpired error = errors.New("Token has expired.")

type ErrList struct {
	list     map[int]interface{}
	listJSON []byte
	message  string
}

func NewErrList(message string, list map[int]interface{}) *ErrList {
	if list == nil || len(list) == 0 {
		return nil
	}

	listJSON, err := json.Marshal(list)
	if err != nil {
		listJSON = make([]byte, 0)
	}

	if message == "" {
		message = "Error list: "
	}

	return &ErrList{
		message:  message,
		list:     list,
		listJSON: listJSON,
	}
}

func (e *ErrList) Error() string {
	return e.message + string(e.listJSON)
}

func (e *ErrList) List() map[int]interface{} {
	return e.list
}

type ErrAlreadyExistsList struct {
	*ErrList
}

func NewErrAlreadyExistsList(message string, list map[int]interface{}) *ErrAlreadyExistsList {
	if message == "" {
		message = "Already exists: "
	}
	return &ErrAlreadyExistsList{NewErrList(message, list)}
}

type ErrNotFoundList struct {
	*ErrList
}

func NewErrNotFoundList(message string, list map[int]interface{}) *ErrNotFoundList {
	if message == "" {
		message = "Not found: "
	}
	return &ErrNotFoundList{NewErrList(message, list)}
}
