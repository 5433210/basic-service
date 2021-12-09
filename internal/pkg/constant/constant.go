package constant

const (
	RegoFilename = "rbac.rego"

	FldSlash = "/"

	FldIndexes      = "indexes"
	FldDomains      = "domains"
	FldDomainId     = "domainId"
	FldName         = "name"
	FldPermissions  = "permissions"
	FldPermissionId = "permissionId"
	FldDenies       = "denies"
	FldDenyId       = "denyId"
	FldRoles        = "roles"
	FldRoldId       = "roleId"
	FldGroups       = "groups"
	FldGroupId      = "groupId"
	FldSubjects     = "subjects"
	FldSubjectId    = "subjectId"
	FldMembers      = "members"
	FldOptions      = "options"
	FldIncluded     = "included"
	FldExcluded     = "excluded"

	RuleNULL                = "data.rbac.rl_null"
	RuleDomains             = "data.rbac.rl_domains"
	RuleGroups              = "data.rbac.rl_groups"
	RuleSubjects            = "data.rbac.rl_subjects"
	RuleDomain              = "data.rbac.rl_domain"
	RuleDeniesCanBeGranted  = "data.rbac.rl_denies_can_be_granted_to_subject"
	RuleDeniesCanBeAccessed = "data.rbac.rl_denies_can_be_accessed_by_subject"
	RuleRolesCanBeGranted   = "data.rbac.rl_roles_can_be_granted_to_subject"
	RuleRolesCanBeAccessed  = "data.rbac.rl_roles_can_be_accessed_by_subject"
	RuleDenyPermissions     = "data.rbac.rl_permissions_of_deny"
	RuleRolePermissions     = "data.rbac.rl_permissions_of_role"
	RuleGroupPermissions    = "data.rbac.rl_permissions_of_group"
	RuleSubjectPermissions  = "data.rbac.rl_permissions_of_subject"

	TypeV1Domain     = "apiv1.Domain"
	TypeV1Permission = "apiv1.Permission"
	TypeV1Role       = "apiv1.Role"
	TypeV1Deny       = "apiv1.Deny"
	TypeV1Group      = "apiv1.Group"
	TypeV1Subject    = "apiv1.Subject"

	DefaultDomainId = "78839721-a274-4a01-a2be-2725903bcf82"

	DsMysql  = "mysql"
	DsSqlite = "sqlite"

	StatActived   = "actived"
	StatUnactived = "unactived"

	FldIdentity       = "identity"
	FldIdentityId     = "identity_id"
	FldCredentials    = "credentials"
	FldCredentialType = "credential_type"
	FldIdentifierId   = "identifier_id"
	FldIdentifier     = "identifier"
	FldIdentifiers    = "identifiers"
	FldState          = "state"
	FldToken          = "token"

	ColCredsIdtfs  = "Credentials.Identifiers"
	ColIdentifiers = "Identifiers"

	IdentifierTypeEmail  = "email"
	IdentifierTypeMobile = "mobile"

	FldStatus = "status"

	FldOk            = "ok"
	FldChallengeMode = "challenge_mode"
	FldChallenge     = "challenge"

	DiscoveryPrifex = "discovery"
	MasterPrifex    = "master"

	ServiceNameCourier = "courier"
	ServiceNameAuthn   = "authn"
	ServiceNameAuthz   = "authz"
	ServiceNameCaptcha = "captcha"
	ServiceNameSched   = "sched"

	CourierTypeEmail = "email"
	CourierTypeSms   = "sms"

	EmailModeTemplate       = "template"
	EmailModeText           = "text"
	EmailModeHtml           = "html"
	EmailOtpTemplateId      = "19768"
	EmailOtpTemplatePattern = "{\"code\":\"%s\"}"
	EmailOtpAddrNoReply     = "noreply@taijik.com"
	EmailOtpSubject         = "动态验证码"

	FldJobId = "job_id"
)
