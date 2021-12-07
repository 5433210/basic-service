package servicev1

import (
	apiv1 "wailik.com/internal/authn/api/v1"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
	dbstore "wailik.com/internal/pkg/store/db"
)

type identitySrvc struct {
	service *service
}

func newIdentitySrvc(s *service) *identitySrvc {
	return &identitySrvc{
		service: s,
	}
}

// (POST /identities).
func (s *identitySrvc) CreateIdentity(transaction *dbstore.Transaction, identity apiv1.Identity) (*apiv1.Identity, error) {
	log.Info("create identity...")
	tx := s.service.store.DB().BeginTxn()
	defer s.service.store.DB().RollbackTxn(tx)

	identityModel, err := s.service.store.Identity().Create(&tx, identity)
	if err != nil {
		return nil, err
	}
	credentials := identity.Credentials
	for i := range *credentials {
		credential := &((*credentials)[i].Credential)
		credential.IdentityId = (*apiv1.IdentityId)(&identityModel.ID)
		credentialModel, err := s.service.store.Credential().Create(&tx, *credential)
		if err != nil {
			return nil, err
		}
		credential.Id = &credentialModel.ID
		identifiers := &((*credentials)[i].Identifiers)
		for i = range *identifiers {
			identifier := &((*identifiers)[i])
			if IsVerifiable(identifier.IdentifierType) {
				if !CheckIndentifierVerifyCode(s.service, s.service.cache, identifier.Identifier, *identifier.IdentifierVerifiyToken) {
					return nil, errors.NewErrorC(errors.ErrCdInvalidVerifyCode, nil)
				}
			}
			identifier.IdentityId = (*apiv1.IdentityId)(&identityModel.ID)
			identifierModel, err := s.service.store.Identifier().Create(&tx, *identifier)
			if err != nil {
				return nil, err
			}
			identifier.Id = (*apiv1.IdentifierId)(&identifierModel.ID)
			err = s.service.store.Identifier().BindCredential(&tx, apiv1.IdentifierId(identifierModel.ID), credentialModel.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	s.service.store.DB().CommitTxn(tx)
	identity.Id = (*apiv1.IdentityId)(&identityModel.ID)
	identity.State = (*apiv1.State)(&identityModel.Stat)
	log.Debugf("identity created:%+v", identity)

	return &identity, nil
}

// (DELETE /identities/{identity_id}).
func (s *identitySrvc) DeleteIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId, token apiv1.AuthenticatedToken) error {
	log.Info("delete identity...")
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}
	tx := s.service.store.DB().BeginTxn()
	defer s.service.store.DB().RollbackTxn(tx)
	// delete related entities by foreign key
	err := s.service.store.Identity().Delete(&tx, identityId)
	if err != nil {
		return err
	}
	s.service.store.DB().CommitTxn(tx)
	log.Debugf("identity deleted:%+v", identityId)

	return nil
}

// Your GET endpoint
// (GET /identities/{identity_id}).
func (s *identitySrvc) GetIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId) (*apiv1.Identity, error) {
	mIdtt, err := s.service.store.Identity().Retrieve(nil, identityId)
	if err != nil {
		return nil, err
	}
	aIdtt := apiv1.Identity{
		Id:          (*apiv1.IdentityId)(&mIdtt.ID),
		State:       (*apiv1.State)(&mIdtt.Stat),
		DomainId:    (*apiv1.DomainId)(&mIdtt.DomainID),
		Credentials: Mdl2ApiCredentialIdentifiers(mIdtt.Credentials),
	}

	return &aIdtt, nil
}

// (GET /identities/{identity_id}/credentials).
func (s *identitySrvc) GetAllCredentialsOfIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId) (*[]apiv1.CredentialIdentifiers, error) {
	mCreds, err := s.service.store.Credential().RetrieveByIdentityId(nil, identityId)
	if err != nil {
		return nil, err
	}

	return Mdl2ApiCredentialIdentifiers(*mCreds), nil
}

// (DELETE /identities/{identity_id}/credentials/{credential_type}).
func (s *identitySrvc) DeleteCredentialOfIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId, credentialType apiv1.CredentialType, token apiv1.AuthenticatedToken) error {
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}

	return s.service.store.Credential().DeleteByIdentityIdAddType(nil, identityId, credentialType)
}

// (PATCH /identities/{identity_id}/credentials/{credential_type}).
func (s *identitySrvc) UpdateCredentialOfIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId, credentialType apiv1.CredentialType, config apiv1.CredentialConfig, token apiv1.AuthenticatedToken) error {
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}

	return s.service.store.Credential().Update(nil, identityId, credentialType, config)
}

// (DELETE /identities/{identity_id}/credentials/{credential_type}/identifiers).
func (s *identitySrvc) UnbindIdentifierToCredentialOfIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId, credentialType apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier, token apiv1.AuthenticatedToken) error {
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}

	return s.service.store.Credential().UnbindIdentifier(nil, identityId, credentialType, domainId, identifierType, identifier)
}

// (PUT /identities/{identity_id}/credentials/{credential_type}/identifiers).
func (s *identitySrvc) BindIdentifierToCredentialOfIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId, credentialType apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier, token apiv1.AuthenticatedToken) error {
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}

	return s.service.store.Credential().BindIdentifier(nil, identityId, credentialType, domainId, identifierType, identifier)
}

// (GET /identities/{identity_id}/identifiers).
func (s *identitySrvc) GetAllIdentifiersOfIdentity(transaction *dbstore.Transaction, identityId apiv1.IdentityId) (*[]apiv1.IdentifierCredentials, error) {
	models, err := s.service.store.Identifier().RetrieveByIdentityId(nil, identityId)
	if err != nil {
		return nil, err
	}

	return Mdl2ApiIdentifierCredentials(models), nil
}

// (PUT /identities/{identity_id}/identifiers).
func (s *identitySrvc) CreateIdentityIdentifer(transaction *dbstore.Transaction, identityId apiv1.IdentityId, identifier apiv1.IdentifierCredentials, token apiv1.AuthenticatedToken) error {
	log.Info("create identifier...")
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}
	tx := s.service.store.DB().BeginTxn()
	defer s.service.store.DB().RollbackTxn(tx)

	if identifier.Identifier == nil {
		return errors.NewErrorC(errors.ErrCdInvalidIdtf, nil)
	}

	if IsVerifiable(identifier.Identifier.IdentifierType) {
		if !CheckIndentifierVerifyCode(s.service, s.service.cache, identifier.Identifier.Identifier, *identifier.Identifier.IdentifierVerifiyToken) {
			return errors.NewErrorC(errors.ErrCdInvalidVerifyCode, nil)
		}
	}

	mIdtf, err := s.service.store.Identifier().Create(&tx, *identifier.Identifier)
	if err != nil {
		return err
	}
	if identifier.Credentials == nil {
		return errors.NewErrorC(errors.ErrCdInvalidCred, nil)
	}
	for _, aCred := range *identifier.Credentials {
		mCred := s.service.store.Credential().Exist(&tx, aCred)
		if mCred == nil {
			*mCred, err = s.service.store.Credential().Create(&tx, aCred)
			if err != nil {
				return err
			}
		}

		err = s.service.store.Identifier().BindCredential(&tx, apiv1.IdentifierId(mIdtf.ID), mCred.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Your GET endpoint
// (GET /identities/{identity_id}/state).
func (s *identitySrvc) GetIdentityState(transaction *dbstore.Transaction, identityId apiv1.IdentityId) (apiv1.State, error) {
	return s.service.store.Identity().RetrieveState(nil, identityId)
}

// (PUT /identities/{identity_id}/state).
func (s *identitySrvc) UpdateIdentityState(transaction *dbstore.Transaction, identityId apiv1.IdentityId, state apiv1.State, token apiv1.AuthenticatedToken) error {
	if !CheckAuthenticatedToken(s.service, s.service.cache, token) {
		return errors.NewErrorC(errors.ErrCdInvalidAuthenToken, nil)
	}

	return s.service.store.Identity().UpdateState(nil, identityId, state)
}
