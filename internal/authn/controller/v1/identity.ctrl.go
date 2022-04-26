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

type identityController struct {
	srvc servicev1.Service
}

func newIdentityController(c *controller) *identityController {
	return &identityController{srvc: c.srvc}
}

// (POST /identities).
func (i *identityController) CreateIdentity(c *fiber.Ctx) error {
	var json apiv1.CreateIdentityJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}

	r, err := i.srvc.Identity().CreateIdentity(nil, apiv1.Identity(json))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldIdentity] = *r

	return core.WriteResponse(c, nil, data)
}

// (DELETE /identities/{identity_id}).
func (i *identityController) DeleteIdentity(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.DeleteIdentityJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	err := i.srvc.Identity().DeleteIdentity(nil, apiv1.IdentityId(identityId), json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// Your GET endpoint
// (GET /identities/{identity_id}).
func (i *identityController) GetIdentity(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))

	identityId := c.Params(constant.FldIdentityId)
	r, err := i.srvc.Identity().GetIdentity(nil, apiv1.IdentityId(identityId))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldIdentity] = *r

	return core.WriteResponse(c, nil, data)
}

// (GET /identities/{identity_id}/credentials).
func (i *identityController) GetAllCredentialsOfIdentity(c *fiber.Ctx) error {
	identityId := c.Params(constant.FldIdentityId)
	r, err := i.srvc.Identity().GetAllCredentialsOfIdentity(nil, apiv1.IdentityId(identityId))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldCredentials] = *r

	return core.WriteResponse(c, nil, data)
}

// (DELETE /identities/{identity_id}/credentials/{credential_type}).
func (i *identityController) DeleteCredentialOfIdentity(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.DeleteCredentialOfIdentityJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	credentialType := c.Params(constant.FldCredentialType)
	err := i.srvc.Identity().DeleteCredentialOfIdentity(nil, apiv1.IdentityId(identityId), apiv1.CredentialType(credentialType), json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (PATCH /identities/{identity_id}/credentials/{credential_type}).
func (i *identityController) UpdateCredentialOfIdentity(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.UpdateCredentialOfIdentityJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	credentialType := c.Params(constant.FldCredentialType)
	err := i.srvc.Identity().UpdateCredentialOfIdentity(nil, apiv1.IdentityId(identityId), apiv1.CredentialType(credentialType), json.Config, json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (DELETE /identities/{identity_id}/credentials/{credential_type}/identifiers).
func (i *identityController) UnbindIdentifierToCredentialOfIdentity(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.UnbindIdentifierToCredentialOfIdentityJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	credentialType := c.Params(constant.FldCredentialType)
	err := i.srvc.Identity().UnbindIdentifierToCredentialOfIdentity(nil, apiv1.IdentityId(identityId), apiv1.CredentialType(credentialType), json.DomainId, json.IdentifierType, json.Identifier, json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (PUT /identities/{identity_id}/credentials/{credential_type}/identifiers).
func (i *identityController) BindIdentifierToCredentialOfIdentity(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.BindIdentifierToCredentialOfIdentityJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	credentialType := c.Params(constant.FldCredentialType)
	err := i.srvc.Identity().BindIdentifierToCredentialOfIdentity(nil, apiv1.IdentityId(identityId), apiv1.CredentialType(credentialType), json.DomainId, json.IdentifierType, json.Identifier, json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// (GET /identities/{identity_id}/identifiers).
func (i *identityController) GetAllIdentifiersOfIdentity(c *fiber.Ctx) error {
	identityId := c.Params(constant.FldIdentityId)
	r, err := i.srvc.Identity().GetAllIdentifiersOfIdentity(nil, apiv1.IdentityId(identityId))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldIdentifiers] = *r

	return core.WriteResponse(c, nil, data)
}

// (PUT /identities/{identity_id}/identifiers).
func (i *identityController) CreateIdentityIdentifer(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.CreateIdentityIdentiferJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	err := i.srvc.Identity().CreateIdentityIdentifer(nil, apiv1.IdentityId(identityId), json.IdentifierCredentials, json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}

// Your GET endpoint
// (GET /identities/{identity_id}/state).
func (i *identityController) GetIdentityState(c *fiber.Ctx) error {
	identityId := c.Params(constant.FldIdentityId)
	r, err := i.srvc.Identity().GetIdentityState(nil, apiv1.IdentityId(identityId))
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	data := make(map[string]interface{})
	data[constant.FldState] = r

	return core.WriteResponse(c, nil, data)
}

// (PUT /identities/{identity_id}/state).
func (i *identityController) UpdateIdentityState(c *fiber.Ctx) error {
	log.Debugf("body:%v", string(c.Body()))
	var json apiv1.UpdateIdentityStateJSONRequestBody
	if err := c.BodyParser(&json); err != nil {
		return core.WriteResponse(c, errors.NewErrorC(errors.ErrCdDataBind, err), nil)
	}
	identityId := c.Params(constant.FldIdentityId)
	err := i.srvc.Identity().UpdateIdentityState(nil, apiv1.IdentityId(identityId), json.State, json.Token)
	if err != nil {
		return core.WriteResponse(c, err, nil)
	}

	return core.WriteResponse(c, nil, nil)
}
