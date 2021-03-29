package alienVault

import "os"

type AVLabAPI struct {
	token string
}

// create new
func New(token string) *AVLabAPI {

	os.Setenv("X_OTX_API_KEY", token)

	return &AVLabAPI{
		token: token,
	}
}
