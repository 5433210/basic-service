package errors

var (
	ErrCdCommon  = ErrorCode{Code: "Z0000", Message: "common error", HttpStatus: StatusInternalServerError}
	ErrCdUnknown = ErrorCode{Code: "Z9999", Message: "unknown error", HttpStatus: StatusInternalServerError}
	ErrCdSuccess = ErrorCode{Code: "00000", Message: "success", HttpStatus: StatusOK}
)

var (
	ErrCdInterSys = ErrorCode{
		Code:       "B0000",
		Message:    "internal system error",
		HttpStatus: StatusInternalServerError,
	}
	ErrCdDataBind          = ErrorCode{Code: "A0001", Message: "data bind error", HttpStatus: StatusBadRequest}
	ErrCdPathParse         = ErrorCode{Code: "A0002", Message: "path parse error", HttpStatus: StatusBadRequest}
	ErrCdDataExist         = ErrorCode{Code: "A0003", Message: "data already exist", HttpStatus: StatusBadRequest}
	ErrCdDataNotExist      = ErrorCode{Code: "A0003", Message: "data do not exist", HttpStatus: StatusBadRequest}
	ErrCdOptionIsNull      = ErrorCode{Code: "A0005", Message: "option is null", HttpStatus: StatusInternalServerError}
	ErrCdNeedInit          = ErrorCode{Code: "A0006", Message: "need to be initialized", HttpStatus: StatusInternalServerError}
	ErrCdKeyIsNull         = ErrorCode{Code: "A0007", Message: "key is null", HttpStatus: StatusInternalServerError}
	ErrCdFuncIsNull        = ErrorCode{Code: "A0008", Message: "func is null", HttpStatus: StatusInternalServerError}
	ErrCdArgIsNull         = ErrorCode{Code: "A0009", Message: "argument is null", HttpStatus: StatusInternalServerError}
	ErrCdDataNotFound      = ErrorCode{Code: "A0010", Message: "data not found", HttpStatus: StatusInternalServerError}
	ErrCdInvalidPermission = ErrorCode{Code: "A0011", Message: "invalid permission", HttpStatus: StatusBadRequest}
	ErrCdInvalidRole       = ErrorCode{Code: "A0012", Message: "invalid role", HttpStatus: StatusBadRequest}
	ErrCdInvalidDeny       = ErrorCode{Code: "A0013", Message: "invalid deny", HttpStatus: StatusBadRequest}
	ErrCdInvalidGroup      = ErrorCode{Code: "A0014", Message: "invalid group", HttpStatus: StatusBadRequest}
	ErrCdInvalidSubject    = ErrorCode{Code: "A0015", Message: "invalid subject", HttpStatus: StatusBadRequest}
	ErrCdInvalidDomain     = ErrorCode{Code: "A0016", Message: "invalid domain", HttpStatus: StatusBadRequest}

	ErrCdNoCredential        = ErrorCode{Code: "A0101", Message: "no credential data", HttpStatus: StatusBadRequest}
	ErrCdNoIdentifier        = ErrorCode{Code: "A0102", Message: "no identifier data", HttpStatus: StatusBadRequest}
	ErrCdInvalidCredCfg      = ErrorCode{Code: "A0102", Message: "invalid credential config", HttpStatus: StatusBadRequest}
	ErrCdInvalidIdtf         = ErrorCode{Code: "A0103", Message: "invalid identifier", HttpStatus: StatusBadRequest}
	ErrCdInvalidCred         = ErrorCode{Code: "A0103", Message: "invalid credential", HttpStatus: StatusBadRequest}
	ErrCdInvalidAuthenToken  = ErrorCode{Code: "A0104", Message: "invalid authenticated token", HttpStatus: StatusBadRequest}
	ErrCdInvalidVerifyCode   = ErrorCode{Code: "A0104", Message: "invalid verify code", HttpStatus: StatusBadRequest}
	ErrCdUnsupportedIdtfType = ErrorCode{Code: "A0105", Message: "unsupported identifier type", HttpStatus: StatusBadRequest}
	ErrCdSaveCacheError      = ErrorCode{Code: "A0106", Message: "saving in cache unsuccessful", HttpStatus: StatusInternalServerError}
	ErrCdPickNodeError       = ErrorCode{Code: "A0107", Message: "picking node unsuccessful", HttpStatus: StatusInternalServerError}

	ErrCdInvalidCacheType = ErrorCode{Code: "A0201", Message: "invalid cache type", HttpStatus: StatusInternalServerError}
)
