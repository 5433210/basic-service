package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func NewCredential() common.CredentialIface {
	credential := common.NewCredential(
		"", "",
	)

	return credential
}
