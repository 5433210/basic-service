package servicev1

import (
	"fmt"
	"time"

	"github.com/rs/xid"

	apiv1 "wailik.com/internal/authn/api/v1"
	"wailik.com/internal/pkg/cache"
	"wailik.com/internal/pkg/constant"
	"wailik.com/internal/pkg/errors"
	"wailik.com/internal/pkg/log"
	"wailik.com/internal/pkg/util"
)

func CheckIndentifierVerifyCode(s *service, cache cache.Cache, identifier apiv1.Identifier, code apiv1.IdentifierVerifyToken) bool {
	if _, ok := cache.Get(fmt.Sprintf("%s:%v:%v", "v_token", identifier, code)); ok {
		cache.Del(code)

		return true
	}

	return false
}

func GenerateIdentifierVerifyCode(s *service, cache cache.Cache, identifierType apiv1.IdentifierType, identifier apiv1.Identifier) error {
	code := fmt.Sprintf("%d", util.RandomInt(999999, 100000))
	log.Debugf("identifier verify code:%v", code)

	switch identifierType {
	case constant.IdentifierTypeEmail:
		err := SendEmail(s, string(identifier), code)
		if err != nil {
			return err
		}
	default:
		return errors.NewErrorC(errors.ErrCdUnsupportedIdtfType, nil)
	}

	ret := cache.Set(fmt.Sprintf("%s:%v:%v", "v_token", identifier, code), "", 60*time.Second)
	if !ret {
		return errors.NewErrorC(errors.ErrCdSaveCacheError, nil)
	}

	return nil
}

func CheckAuthenticatedToken(s *service, cache cache.Cache, code apiv1.AuthenticatedToken) bool {
	if _, ok := cache.Get(code); ok {
		cache.Del(code)

		return true
	}

	return false
}

func GenerateAuthenticatedToken(s *service, cache cache.Cache) apiv1.AuthenticatedToken {
	code := xid.New().String()
	if !cache.Set(code, "", 60*time.Second) {
		return ""
	}

	log.Debugf("authenticated code:%v", code)

	return apiv1.AuthenticatedToken(code)
}

func IsVerifiable(identifierType apiv1.IdentifierType) bool {
	switch identifierType {
	case constant.IdentifierTypeEmail:
		return true
	case constant.IdentifierTypeMobile:
		return true
	}

	return false
}
