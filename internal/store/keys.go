package store

// Key constants matching the TypeScript keychain item keys.
const (
	KeyToken            = "token"
	KeyControllerHost   = "ziti_controller_host"
	KeyDomain           = "domain"
	KeyRefreshToken     = "refresh_token"
	KeyTokenExpiresAt   = "token_expires_at"
	KeyIdentityCert     = "identity_cert"
	KeyIdentityKey      = "identity_key"
	KeyIdentityCA       = "identity_ca"
	KeyUpdbUsername      = "updb_username"
	KeyUpdbPassword      = "updb_password"
	KeyControllerCA     = "controller_ca"
	KeyIDPClientID      = "idp_client_id"
)

// AllKeys is the list of all credential store keys, used for clear operations.
var AllKeys = []string{
	KeyToken,
	KeyControllerHost,
	KeyDomain,
	KeyRefreshToken,
	KeyTokenExpiresAt,
	KeyIdentityCert,
	KeyIdentityKey,
	KeyIdentityCA,
	KeyUpdbUsername,
	KeyUpdbPassword,
	KeyControllerCA,
	KeyIDPClientID,
}
