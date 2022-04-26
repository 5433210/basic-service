package storev1

import (
	"encoding/json"

	apiv1 "wailik.com/internal/authn/api/v1"
	modelv1 "wailik.com/internal/authn/model/v1"
	"wailik.com/internal/pkg/constant"
	dbstore "wailik.com/internal/pkg/store/db"
)

type credentialStore struct {
	store Store
}

func newCredentialStore(store Store) *credentialStore {
	s := credentialStore{}
	s.store = store

	return &s
}

func (s *credentialStore) Create(transaction *dbstore.Transaction, credential apiv1.CredentialObject) (modelv1.Credential, error) {
	config, _ := json.Marshal(credential.Config)
	model := modelv1.Credential{
		IdentityID:     string(*credential.IdentityId),
		CredentialType: string(credential.Type),
		Config:         string(config),
	}

	result := s.store.DB().Transaction(transaction).Handler.Create(&model)

	return model, result.Error
}

func (s *credentialStore) Exist(transaction *dbstore.Transaction, credential apiv1.CredentialObject) *modelv1.Credential {
	model := modelv1.Credential{
		IdentityID:     string(*credential.IdentityId),
		CredentialType: string(credential.Type),
	}

	result := s.store.DB().Transaction(transaction).Handler.Find(&model)

	if result.Error == nil {
		return &model
	}

	return nil
}

func (s *credentialStore) DeleteByIdentityId(transaction *dbstore.Transaction, identityId apiv1.IdentityId) error {
	model := modelv1.Credential{}
	model.IdentityID = string(identityId)
	result := s.store.DB().Transaction(transaction).Handler.Delete(&model)

	return result.Error
}

func (s *credentialStore) DeleteByIdentityIdAddType(transaction *dbstore.Transaction, identityId apiv1.IdentityId, typ apiv1.CredentialType) error {
	model := modelv1.Credential{}
	model.IdentityID = string(identityId)
	model.CredentialType = string(typ)
	result := s.store.DB().Transaction(transaction).Handler.Delete(&model)

	return result.Error
}

func (s *credentialStore) Update(transaction *dbstore.Transaction, identityId apiv1.IdentityId, typ apiv1.CredentialType, config apiv1.CredentialConfig) error {
	model := modelv1.Credential{}
	model.IdentityID = string(identityId)
	model.CredentialType = string(typ)
	marshaled, err := json.Marshal(config)
	if err != nil {
		return err
	}

	result := s.store.DB().Transaction(transaction).Handler.Where(&model).Updates(modelv1.Credential{
		Config: string(marshaled),
	})

	return result.Error
}

func (s *credentialStore) RetrieveByIdentityId(transaction *dbstore.Transaction, identityId apiv1.IdentityId) (*[]modelv1.Credential, error) {
	models := []modelv1.Credential{}

	result := s.store.DB().Transaction(transaction).Handler.Preload(constant.ColIdentifiers).Where(&modelv1.Credential{
		IdentityID: string(identityId),
	}).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return &models, nil
}

func (s *credentialStore) UnbindIdentifier(transaction *dbstore.Transaction, identityId apiv1.IdentityId, typ apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) error {
	model := modelv1.Identifier{}
	model.IdentitiyID = string(identityId)
	model.Identifier = string(identifier)
	model.IdentifierType = string(identifierType)

	result := s.store.DB().Transaction(transaction).Handler.Find(&model)
	if result.Error != nil {
		return result.Error
	}

	model2 := modelv1.Credential{}
	model2.IdentityID = string(identityId)
	model2.CredentialType = string(typ)

	result = s.store.DB().Transaction(transaction).Handler.Find(&model2)
	if result.Error != nil {
		return result.Error
	}

	model3 := modelv1.IdentifierCredential{}
	model3.CredentialID = model2.ID
	model3.IdentifierID = model.ID

	result = s.store.DB().Transaction(transaction).Handler.Delete(&model3)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *credentialStore) BindIdentifier(transaction *dbstore.Transaction, identityId apiv1.IdentityId, typ apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) error {
	model := modelv1.Identifier{}
	model.IdentitiyID = string(identityId)
	model.Identifier = string(identifier)
	model.IdentifierType = string(identifierType)

	result := s.store.DB().Transaction(transaction).Handler.Find(&model)
	if result.Error != nil {
		return result.Error
	}

	model2 := modelv1.Credential{}
	model2.IdentityID = string(identityId)
	model2.CredentialType = string(typ)

	result = s.store.DB().Transaction(transaction).Handler.Find(&model2)
	if result.Error != nil {
		return result.Error
	}

	model3 := modelv1.IdentifierCredential{}
	model3.CredentialID = model2.ID
	model3.IdentifierID = model.ID

	result = s.store.DB().Transaction(transaction).Handler.Create(&model3)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *credentialStore) RetrieveConfigByTypeAndIdentifier(transaction *dbstore.Transaction, identityId apiv1.IdentityId, typ apiv1.CredentialType, domainId apiv1.DomainId, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) (*apiv1.CredentialConfig, error) {
	model := modelv1.Credential{}
	model.IdentityID = string(identityId)
	model.CredentialType = string(typ)

	result := s.store.DB().Transaction(transaction).Handler.Preload("Identifiers", "domain_id=? and identifier_type=? and identifier=?", domainId, identifierType, identifier).Find(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	config := apiv1.CredentialConfig{}
	err := json.Unmarshal([]byte(model.Config), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
