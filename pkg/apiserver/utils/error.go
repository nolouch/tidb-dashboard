// Copyright 2021 PingCAP, Inc. Licensed under Apache-2.0.

package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
)

var (
	ErrNS    = errorx.NewNamespace("error.api")
	ErrOther = ErrNS.NewType("other")
)

var ErrUnauthorized = ErrNS.NewType("unauthorized")

func MakeUnauthorizedError(c *gin.Context) {
	_ = c.Error(ErrUnauthorized.New("Sign in is required"))
	c.Status(http.StatusUnauthorized)
}

var ErrInsufficientPrivilege = ErrNS.NewType("insufficient_privilege")

func MakeInsufficientPrivilegeError(c *gin.Context) {
	_ = c.Error(ErrInsufficientPrivilege.New("Insufficient privilege"))
	c.Status(http.StatusForbidden)
}

var ErrInvalidRequest = ErrNS.NewType("invalid_request")

func MakeInvalidRequestErrorWithMessage(c *gin.Context, message string, args ...interface{}) {
	_ = c.Error(ErrInvalidRequest.New(message, args...))
	c.Status(http.StatusBadRequest)
}

func MakeInvalidRequestErrorFromError(c *gin.Context, err error) {
	_ = c.Error(ErrInvalidRequest.WrapWithNoMessage(err))
	c.Status(http.StatusBadRequest)
}

var ErrExpNotEnabled = ErrNS.NewType("experimental_feature_not_enabled")

var ErrFeatureNotSupported = ErrNS.NewType("feature_not_supported")

type APIError struct {
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	FullText string `json:"full_text"`
}

func NewAPIError(err error) *APIError {
	innerErr := errorx.Cast(err)
	if innerErr == nil {
		innerErr = ErrOther.WrapWithNoMessage(err)
	}
	return &APIError{
		Error:    true,
		Message:  innerErr.Error(),
		Code:     errorx.GetTypeName(innerErr),
		FullText: fmt.Sprintf("%+v", innerErr),
	}
}

// MWHandleErrors creates a middleware that turns (last) error in the context into an APIError json response.
// In handlers, `c.Error(err)` can be used to attach the error to the context.
// When error is attached in the context:
// - The handler can optionally assign the HTTP status code.
// - The handler must not self-generate a response body.
func MWHandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		statusCode := c.Writer.Status()
		if statusCode == http.StatusOK {
			statusCode = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(statusCode, NewAPIError(err.Err))
	}
}
