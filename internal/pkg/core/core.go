// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/gofiber/fiber/v2"

	"wailik.com/internal/pkg/errors"
)

// Response defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type Response struct {
	// Code defines the business error code.
	Code string `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`

	Data interface{} `json:"data,omitempty"`
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *fiber.Ctx, err error, data interface{}) error {
	if err != nil {
		code := errors.Code(err)

		return c.JSON(Response{
			Code:      code.Code,
			Message:   code.Message,
			Reference: code.Ref,
		})
	}

	return c.JSON(Response{
		Code:    errors.ErrCdSuccess.Code,
		Message: errors.ErrCdSuccess.Message,
		Data:    data,
	})
}
