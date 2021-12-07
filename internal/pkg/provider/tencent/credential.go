package tencent

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func NewCredential() common.CredentialIface {
	credential := common.NewCredential(
		"AKIDnQnwqTqwmLddNq2xb4GCrGPHsCF1WkhN",
		"QZS4PTIsEK9i01lFSO3SaR1f3A0Zsuf9",
	)

	return credential
}
