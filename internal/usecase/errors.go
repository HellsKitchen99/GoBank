package usecase

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")
var ErrAmountToTransfer = errors.New("not enough balance")
var ErrSameUser = errors.New("you cant't transfer yourself")
