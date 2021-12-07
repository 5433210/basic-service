package servicev1

import (
	"encoding/json"

	apiv1 "wailik.com/internal/authn/api/v1"
	modelv1 "wailik.com/internal/authn/model/v1"
	"wailik.com/internal/pkg/log"
)

type credentialSrvc struct {
	service *service
}

func newCredentialSrvc(s *service) *credentialSrvc {
	return &credentialSrvc{
		service: s,
	}
}

func Mdl2ApiCredentialIdentifiers(credList []modelv1.Credential) *[]apiv1.CredentialIdentifiers {
	r := make([]apiv1.CredentialIdentifiers, len(credList))
	for i, cred := range credList {
		r[i] = *Mdl2ApiCredentialIdentifier(cred)
	}

	return &r
}

func Mdl2ApiIdentifierCredentials(idtfList []modelv1.Identifier) *[]apiv1.IdentifierCredentials {
	r := make([]apiv1.IdentifierCredentials, len(idtfList))
	for i, idtf := range idtfList {
		r[i] = *Mdl2ApiIdentifierCredential(idtf)
	}

	return &r
}

func Mdl2ApiCredentialIdentifier(cred modelv1.Credential) *apiv1.CredentialIdentifiers {
	return &apiv1.CredentialIdentifiers{
		Credential:  *Mdl2ApiCredential(cred),
		Identifiers: *Mdl2ApiIdentifiers(cred.Identifiers),
	}
}

func Mdl2ApiIdentifierCredential(idtf modelv1.Identifier) *apiv1.IdentifierCredentials {
	return &apiv1.IdentifierCredentials{
		Identifier:  Mdl2ApiIdentifier(idtf),
		Credentials: Mdl2ApiCredentials(idtf.Credentials),
	}
}

func Mdl2ApiCredential(cred modelv1.Credential) *apiv1.CredentialObject {
	return &apiv1.CredentialObject{
		Id:         &cred.ID,
		Type:       apiv1.CredentialType(cred.CredentialType),
		IdentityId: (*apiv1.IdentityId)(&cred.IdentityID),
		Config: func() *apiv1.CredentialConfig {
			var config apiv1.CredentialConfig
			err := json.Unmarshal([]byte(cred.Config), &config)
			if err != nil {
				log.Warn("unmashal config error")

				return &config
			}

			return &config
		}(),
	}
}

func Mdl2ApiCredentials(credList []modelv1.Credential) *[]apiv1.CredentialObject {
	r := make([]apiv1.CredentialObject, len(credList))
	for i, cred := range credList {
		r[i] = *Mdl2ApiCredential(cred)
	}

	return &r
}

func Mdl2ApiIdentifier(idtf modelv1.Identifier) *apiv1.IdentifierObject {
	return &apiv1.IdentifierObject{
		Id:             (*apiv1.IdentifierId)(&idtf.ID),
		IdentifierType: apiv1.IdentifierType(idtf.IdentifierType),
		Identifier:     apiv1.Identifier(idtf.Identifier),
		DomainId:       apiv1.DomainId(idtf.DomainID),
	}
}

func Mdl2ApiIdentifiers(idtfList []modelv1.Identifier) *[]apiv1.IdentifierObject {
	r := make([]apiv1.IdentifierObject, len(idtfList))
	for i, idtf := range idtfList {
		r[i] = *Mdl2ApiIdentifier(idtf)
	}

	return &r
}
