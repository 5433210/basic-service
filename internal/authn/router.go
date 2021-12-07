package authn

import (
	"github.com/gofiber/fiber/v2"

	controllerv1 "wailik.com/internal/authn/controller/v1"
	servicev1 "wailik.com/internal/authn/service/v1"
)

func route(app *fiber.App, srvc servicev1.Service) {
	c := controllerv1.NewController(srvc)

	v1 := app.Group("v1")
	{
		identityGroup(v1, c)
		identifierGroup(v1, c)
	}
}

func identityGroup(v1 fiber.Router, c controllerv1.Controller) {
	v1Identity := v1.Group("identities")
	{
		v1Identity.Post("", c.Identities().CreateIdentity)
		v1Identity.Delete(":identity_id", c.Identities().DeleteIdentity)
		v1Identity.Get(":identity_id", c.Identities().GetIdentity)
		v1Identity.Get(":identity_id/credentials", c.Identities().GetAllCredentialsOfIdentity)
		v1Identity.Delete(":identity_id/credentials/:credential_type", c.Identities().DeleteCredentialOfIdentity)
		v1Identity.Patch(":identity_id/credentials/:credential_type", c.Identities().UpdateCredentialOfIdentity)
		v1Identity.Delete(":identity_id/credentials/:credential_type/identifiers", c.Identities().UnbindIdentifierToCredentialOfIdentity)
		v1Identity.Put(":identity_id/credentials/:credential_type/identifiers", c.Identities().BindIdentifierToCredentialOfIdentity)
		v1Identity.Get(":identity_id/identifiers", c.Identities().GetAllIdentifiersOfIdentity)
		v1Identity.Put(":identity_id/identifiers", c.Identities().CreateIdentityIdentifer)
		v1Identity.Get(":identity_id/state", c.Identities().GetIdentityState)
		v1Identity.Put(":identity_id/state", c.Identities().UpdateIdentityState)
	}
}

func identifierGroup(v1 fiber.Router, c controllerv1.Controller) {
	v1Identifier := v1.Group("identifiers")
	{
		v1Identifier.Delete(":identifier_id", c.Identifiers().DeleteIdentifier)
		v1Identifier.Get(":identifier_id", c.Identifiers().GetIdentifier)
		v1Identifier.Patch(":identifier_id", c.Identifiers().ChangeIdentifier)
		v1Identifier.Get("id", c.Identifiers().GetIdentifierId)
		v1Identifier.Put("verification", c.Identifiers().RequestIdentifierVerifyCode)
		v1Identifier.Post("authentication", c.Identifiers().CheckCredentialForAuthentication)
		v1Identifier.Put("authentication", c.Identifiers().GenerateCredentialForAuthentication)
	}
}
