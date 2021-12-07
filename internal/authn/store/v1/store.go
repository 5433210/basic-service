package storev1

import (
	"wailik.com/internal/pkg/log"
	dbstore "wailik.com/internal/pkg/store/db"
)

type Store interface {
	Identity() *identityStore
	Identifier() *identifierStore
	Credential() *credentialStore
	DB() dbstore.Store
}

type store struct {
	dbstore.Store
}

type Options struct {
	DSN string
}

var _ Store = &store{}

func (s *store) Identity() *identityStore     { return newIdentityStore(s) }
func (s *store) Identifier() *identifierStore { return newIdentifierStore(s) }
func (s *store) Credential() *credentialStore { return newCredentialStore(s) }
func (s *store) DB() dbstore.Store            { return s.Store }

func New(opts *Options) (Store, error) {
	log.Info("new store...")
	s, err := dbstore.New(&dbstore.Options{
		DSN: opts.DSN,
	})
	if err != nil {
		return nil, err
	}

	if err = s.Open(); err != nil {
		return nil, err
	}

	log.Info("store created")

	return &store{
		Store: s,
	}, nil
}
