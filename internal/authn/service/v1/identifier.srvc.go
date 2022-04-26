package servicev1

import (
	"reflect"

	apiv1 "wailik.com/internal/authn/api/v1"
	"wailik.com/internal/pkg/errors"
	dbstore "wailik.com/internal/pkg/store/db"
)

type identifierSrvc struct {
	service *service
}

func newIdentifierSrvc(s *service) *identifierSrvc {
	return &identifierSrvc{
		service: s,
	}
}

// (GET /identifiers/id).
func (s *identifierSrvc) GetIdentifierId(transaction *dbstore.Transaction, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) (apiv1.IdentifierId, error) {
	return s.service.store.Identifier().RetrieveIdByUniquekey(nil, domainId, identifierType, identifier)
}

// (PUT /identifiers/verification).
func (s *identifierSrvc) RequestIdentifierVerifyCode(transaction *dbstore.Transaction, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) error {
	_, err := s.service.store.Identifier().RetrieveIdByUniquekey(nil, domainId, identifierType, identifier)
	if err != nil {
		return err
	}
	// generate token & put into cache
	err = GenerateIdentifierVerifyCode(s.service, s.service.cache, identifierType, identifier)
	if err != nil {
		return err
	}

	return err
}

// Your GET endpoint
// (GET /identifiers/{identifier_id}).
func (s *identifierSrvc) GetIdentifier(transaction *dbstore.Transaction, identifierId apiv1.IdentifierId) (*apiv1.IdentifierCredentials, error) {
	mIdtf, err := s.service.store.Identifier().Retrieve(nil, identifierId)
	if err != nil {
		return nil, err
	}
	aIdtf := Mdl2ApiIdentifierCredential(mIdtf)

	return aIdtf, nil
}

// (PATCH /identifiers/{identifier_id}).
func (s *identifierSrvc) ChangeIdentifier(transaction *dbstore.Transaction, identifierId apiv1.IdentifierId, identifier apiv1.Identifier, vToken *apiv1.IdentifierVerifyToken, token apiv1.AuthenticatedToken) error {
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}

	mIdtf, err := s.service.store.Identifier().Retrieve(nil, identifierId)
	if err != nil {
		return err
	}

	if IsVerifiable(apiv1.IdentifierType(mIdtf.IdentifierType)) {
		if !CheckIndentifierVerifyCode(s.service, s.service.cache, identifier, *vToken) {
			return errors.NewErrorC(errors.ErrCdInvalidVerifyCode, nil)
		}
	}

	return s.service.store.Identifier().UpdateIdentifier(nil, identifierId, identifier)
}

// (DELETE /identifiers/{identifier_id}).
func (s *identifierSrvc) DeleteIdentifier(transaction *dbstore.Transaction, identityId apiv1.IdentityId, identifierId apiv1.IdentifierId, token apiv1.AuthenticatedToken) error {
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}

	return s.service.store.Identifier().Delete(nil, identifierId)
}

// (POST /identifiers/authentication).
func (s *identifierSrvc) CheckCredentialForAuthentication(transaction *dbstore.Transaction, credentialType apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier, config apiv1.CredentialConfig) (*apiv1.AuthenticatedToken, error) {
	mIdtf, err := s.service.store.Identifier().RetrieveByUniquekeyAndCredType(nil, credentialType, domainId, identifierType, identifier)
	if err != nil {
		return nil, err
	}
	if mIdtf.Credentials == nil || len(mIdtf.Credentials) == 0 {
		return nil, errors.NewErrorC(errors.ErrCdDataNotFound, nil)
	}

	origConfig := mIdtf.Credentials[0].Config

	if reflect.DeepEqual(origConfig, config) {
		token := GenerateAuthenticatedToken(s.service, s.service.cache)

		return &token, nil
	}

	return nil, nil
}

// (PUT /identifiers/authentication).
func (s *identifierSrvc) GenerateCredentialForAuthentication(transaction *dbstore.Transaction, credentialType apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) error {
	_, err := s.service.store.Identifier().RetrieveByUniquekeyAndCredType(nil, credentialType, domainId, identifierType, identifier)
	if err != nil {
		return err
	}

	return nil
}
