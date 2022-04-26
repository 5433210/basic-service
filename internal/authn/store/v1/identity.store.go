package storev1

import (
	apiv1 "wailik.com/internal/authn/api/v1"
	modelv1 "wailik.com/internal/authn/model/v1"
	"wailik.com/internal/pkg/constant"
	dbstore "wailik.com/internal/pkg/store/db"
)

type identityStore struct {
	store Store
}

func newIdentityStore(store Store) *identityStore {
	s := identityStore{}
	s.store = store

	return &s
}

func (s *identityStore) Create(transaction *dbstore.Transaction, identity apiv1.Identity) (modelv1.Identity, error) {
	model := modelv1.Identity{
		DomainID: string(*identity.DomainId),
	}

	result := s.store.DB().Transaction(transaction).Handler.Create(&model)

	return model, result.Error
}

func (s *identityStore) Delete(transaction *dbstore.Transaction, identityId apiv1.IdentityId) error {
	model := modelv1.Identity{}
	model.ID = string(identityId)

	result := s.store.DB().Transaction(transaction).Handler.Delete(&model)

	return result.Error
}

func (s *identityStore) Retrieve(transaction *dbstore.Transaction, identityId apiv1.IdentityId) (modelv1.Identity, error) {
	model := modelv1.Identity{}

	result := s.store.DB().Transaction(transaction).Handler.Preload(constant.ColCredsIdtfs).Find(&model, identityId)

	return model, result.Error
}

func (s *identityStore) UpdateState(transaction *dbstore.Transaction, identityId apiv1.IdentityId, state apiv1.State) error {
	model := modelv1.Identity{}

	result := s.store.DB().Transaction(transaction).Handler.Where(&model).Updates(modelv1.Identity{
		Stat: string(state),
	})

	return result.Error
}

func (s *identityStore) RetrieveState(transaction *dbstore.Transaction, identityId apiv1.IdentityId) (apiv1.State, error) {
	model := modelv1.Identity{}

	result := s.store.DB().Transaction(transaction).Handler.Find(&model, identityId)

	return apiv1.State(model.Stat), result.Error
}
