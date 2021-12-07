package storev1

import (
	apiv1 "wailik.com/internal/authn/api/v1"
	modelv1 "wailik.com/internal/authn/model/v1"
	dbstore "wailik.com/internal/pkg/store/db"
)

type identifierStore struct {
	store Store
}

func newIdentifierStore(store Store) *identifierStore {
	s := identifierStore{}
	s.store = store

	return &s
}

func (s *identifierStore) Create(transaction *dbstore.Transaction, identifier apiv1.IdentifierObject) (modelv1.Identifier, error) {
	model := modelv1.Identifier{
		IdentifierType: string(identifier.IdentifierType),
		Identifier:     string(identifier.Identifier),
		DomainID:       string(identifier.DomainId),
		IdentitiyID:    string(*identifier.IdentityId),
	}

	result := s.store.DB().Transaction(transaction).Handler.Create(&model)

	return model, result.Error
}

func (s *identifierStore) BindCredential(transaction *dbstore.Transaction, identifierId apiv1.IdentifierId, credentialId string) error {
	model := modelv1.IdentifierCredential{
		IdentifierID: string(identifierId),
		CredentialID: credentialId,
	}

	result := s.store.DB().Transaction(transaction).Handler.Create(&model)

	return result.Error
}

func (s *identifierStore) Delete(transaction *dbstore.Transaction, identifierId apiv1.IdentifierId) error {
	model := modelv1.Identifier{}
	model.IdentitiyID = string(identifierId)

	result := s.store.DB().Transaction(transaction).Handler.Delete(&model)

	return result.Error
}

func (s *identifierStore) DeleteByIdentityId(transaction *dbstore.Transaction, identityId apiv1.IdentityId) error {
	model := modelv1.Identifier{}
	model.IdentitiyID = string(identityId)
	result := s.store.DB().Transaction(transaction).Handler.Delete(&model)

	return result.Error
}

func (s *identifierStore) DeleteByIdentifier(transaction *dbstore.Transaction, identifier apiv1.Identifier, identifierType apiv1.IdentifierType, domainId apiv1.DomainId) error {
	model := modelv1.Identifier{}
	model.DomainID = string(domainId)
	model.Identifier = string(identifier)
	model.IdentifierType = string(identifierType)
	result := s.store.DB().Transaction(transaction).Handler.Delete(&model)

	return result.Error
}

func (s *identifierStore) DeleteByIdentityIdAndCredentialType(transaction *dbstore.Transaction, identityId apiv1.IdentityId, typ apiv1.CredentialType) error {
	model := modelv1.Identifier{}
	model.IdentitiyID = string(identityId)
	model.Credentials = make([]modelv1.Credential, 1)

	result := s.store.DB().Transaction(transaction).Handler.Delete(&model)

	return result.Error
}

func (s *identifierStore) RetrieveByIdentityId(transaction *dbstore.Transaction, identityId apiv1.IdentityId) ([]modelv1.Identifier, error) {
	models := []modelv1.Identifier{}
	result := s.store.DB().Transaction(transaction).Handler.Preload("Credentials").Where("identity_id=?", identityId).Find(&models)

	return models, result.Error
}

func (s *identifierStore) Retrieve(transaction *dbstore.Transaction, identifierId apiv1.IdentifierId) (modelv1.Identifier, error) {
	model := modelv1.Identifier{}
	model.ID = string(identifierId)

	result := s.store.DB().Transaction(transaction).Handler.Preload("Credentials").Find(&model)

	return model, result.Error
}

func (s *identifierStore) UpdateIdentifier(transaction *dbstore.Transaction, identifierId apiv1.IdentifierId, identifier apiv1.Identifier) error {
	model := modelv1.Identifier{}
	model.ID = string(identifierId)

	result := s.store.DB().Transaction(transaction).Handler.Where(&model).Updates(modelv1.Identifier{
		Identifier: string(identifier),
	})

	return result.Error
}

func (s *identifierStore) RetrieveIdByUniquekey(transaction *dbstore.Transaction, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) (apiv1.IdentifierId, error) {
	model := modelv1.Identifier{}
	model.DomainID = string(domainId)
	model.IdentifierType = string(identifierType)
	model.Identifier = string(identifier)

	result := s.store.DB().Transaction(transaction).Handler.Find(&model)

	return apiv1.IdentifierId(model.ID), result.Error
}

func (s *identifierStore) RetrieveByUniquekey(transaction *dbstore.Transaction, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) (modelv1.Identifier, error) {
	model := modelv1.Identifier{}
	model.DomainID = string(domainId)
	model.IdentifierType = string(identifierType)
	model.Identifier = string(identifier)

	result := s.store.DB().Transaction(transaction).Handler.Find(&model)

	return model, result.Error
}

func (s *identifierStore) RetrieveByUniquekeyAndCredType(transaction *dbstore.Transaction, credentialType apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) (modelv1.Identifier, error) {
	model := modelv1.Identifier{}
	model.DomainID = string(domainId)
	model.IdentifierType = string(identifierType)
	model.Identifier = string(identifier)

	result := s.store.DB().Transaction(transaction).Handler.Preload("Credentials", "credential_type=?", credentialType).Find(&model)

	return model, result.Error
}
