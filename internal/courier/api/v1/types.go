// Package apiv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.3 DO NOT EDIT.
package apiv1

import "encoding/json"

// Email defines model for Email.
type Email struct {
	Content EmailContent    `json:"content"`
	From    EmailAddress    `json:"from"`
	ReplyTo *[]EmailAddress `json:"reply_to,omitempty"`
	Subject string          `json:"subject"`
	To      []EmailAddress  `json:"to"`
}

// EmailAddress defines model for EmailAddress.
type EmailAddress string

// EmailContent defines model for EmailContent.
type EmailContent struct {
	Mode    EmailContentMode `json:"mode"`
	Payload json.RawMessage  `json:"payload"`
}

// EmailContentMode defines model for EmailContentMode.
type EmailContentMode string

// EmailTemplate defines model for EmailTemplate.
type EmailTemplate struct {
	Data string `json:"data"`
	Id   string `json:"id"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Code    string                  `json:"code"`
	Data    *map[string]interface{} `json:"data,omitempty"`
	Message string                  `json:"message"`
}

// Html defines model for Html.
type Html struct {
	Html string `json:"html"`
}

// PhoneNumber defines model for PhoneNumber.
type PhoneNumber string

// PhoneNumberObject defines model for PhoneNumberObject.
type PhoneNumberObject struct {
	Number PhoneNumber `json:"number"`
	Region RegionCode  `json:"region"`
}

// RegionCode defines model for RegionCode.
type RegionCode string

// SendStatus defines model for SendStatus.
type SendStatus struct {
	Code       string `json:"code"`
	Identifier string `json:"identifier"`
	Message    string `json:"message"`
}

// Sms defines model for Sms.
type Sms struct {
	Content SmsContent          `json:"content"`
	From    *PhoneNumberObject  `json:"from,omitempty"`
	To      []PhoneNumberObject `json:"to"`
}

// SmsContent defines model for SmsContent.
type SmsContent struct {
	Payload json.RawMessage `json:"payload"`
	Mode    SmsContentMode  `json:"mode"`
}

// SmsContentMode defines model for SmsContentMode.
type SmsContentMode string

// SmsTemplate defines model for SmsTemplate.
type SmsTemplate struct {
	Id     string             `json:"id"`
	Name   *string            `json:"name,omitempty"`
	Params []SmsTemplateParam `json:"params"`
}

// SmsTemplateParam defines model for SmsTemplateParam.
type SmsTemplateParam string

// Text defines model for Text.
type Text struct {
	Text string `json:"text"`
}

// Error defines model for Error.
type Error ErrorResponse

// SendEmailJSONBody defines parameters for SendEmail.
type SendEmailJSONBody Email

// SendSmsJSONBody defines parameters for SendSms.
type SendSmsJSONBody Sms

// SendEmailJSONRequestBody defines body for SendEmail for application/json ContentType.
type SendEmailJSONRequestBody SendEmailJSONBody

// SendSmsJSONRequestBody defines body for SendSms for application/json ContentType.
type SendSmsJSONRequestBody SendSmsJSONBody
