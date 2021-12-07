package controllerv1

import (
	"github.com/gofiber/fiber/v2"

	apiv1 "wailik.com/internal/authn/api/v1"
	servicev1 "wailik.com/internal/authn/service/v1"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/core"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
)

type identifierController struct {
	srvc servicev1.Service
}

func newIdentifierController(c *controller) *identifierController {
	return &identifierController{srvc: c.srvc}
}

// Your GET endpoint
// (GET /identifiers/{identifier_id}).
func (i *identifierController) GetIdentifier(c *fiber.Ctx) error {
	identifierId := c.Params(constant.FldCredentialType)
	r, err := i.srvc.Identifier().GetIdentifier(nil, apiv1.IdentifierId(identifierId))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldIdentifier] = *r

	return core.WriteResponse(c, nil, data)
}

// (PATCH /identifiers/{identifier_id}).
func (i *identifierController) ChangeIdentifier(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.ChangeIdentifierJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identifierId := c.Params(constant.FldCredentialType)
	err := i.srvc.Identifier().ChangeIdentifier(nil, apiv1.IdentifierId(identifierId), json.Idetifier, json.IdentifierVerifyToken, json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (DELETE /identifiers/{identifier_id}).
func (i *identifierController) DeleteIdentifier(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.DeleteIdentifierJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	identifierId := c.Params(constant.FldCredentialType)
	err := i.srvc.Identifier().DeleteIdentifier(nil, apiv1.IdentityId(identityId), apiv1.IdentifierId(identifierId), json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// Your GET endpoint
// (GET /identifiers/id).
func (i *identifierController) GetIdentifierId(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.GetIdentifierIdJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r, err := i.srvc.Identifier().GetIdentifierId(nil, json.DomainId, json.IdentifierType, json.Identifier)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldIdentifierId] = r

	return core.WriteResponse(c, nil, data)
}

// (PUT /identifiers/verification).
func (i *identifierController) RequestIdentifierVerifyCode(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.RequestIdentifierVerifyTokenJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	err := i.srvc.Identifier().RequestIdentifierVerifyCode(nil, json.DomainId, json.IdentifierType, json.Identifier)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (POST /identifiers/authentication).
func (i *identifierController) CheckCredentialForAuthentication(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.CheckCredentialForAuthenticationJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r, err := i.srvc.Identifier().CheckCredentialForAuthentication(nil, json.CredentialType, json.DomainId, json.IdentifierType, json.Identifier, json.CredentialConfig)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}
	data := make(map[string]interface{})
	data[constant.FldToken] = r

	return core.WriteResponse(c, nil, data)
}

// (PUT /identifiers/authentication).
func (i *identifierController) GenerateCredentialForAuthentication(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.GenerateCredentialForAuthenticationJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	err := i.srvc.Identifier().GenerateCredentialForAuthentication(nil, json.CredentialType, json.DomainId, json.IdentifierType, json.Identifier)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}
