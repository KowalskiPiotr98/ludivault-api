package auth

import (
	"fmt"
	"github.com/KowalskiPiotr98/ludivault/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gitea"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	providers []string

	setups = map[string]func(callbackUrl string) (bool, goth.Provider){
		"gitea": func(callbackUrl string) (bool, goth.Provider) {
			clientId := utils.GetOptionalConfig("SSO_GITEA_CLIENT_ID", "")
			clientSecret := utils.GetOptionalConfig("SSO_GITEA_CLIENT_SECRET", "")
			url := utils.GetOptionalConfig("SSO_GITEA_URL", "")

			if clientId == "" || clientSecret == "" || url == "" {
				return false, nil
			}
			url = strings.TrimRight(url, "/")

			return true, gitea.NewCustomisedURL(clientId, clientSecret, callbackUrl, fmt.Sprintf("%s/login/oauth/authorize", url), fmt.Sprintf("%s/login/oauth/access_token", url), fmt.Sprintf("%s/api/v1/user", url))
		},
	}
)

func SetupProviders(baseUrl string) error {
	if areProvidersSet() {
		// prevent duplicate provider setup
		return ProvidersAlreadyInitialised
	}

	baseUrl = strings.TrimRight(baseUrl, "/")
	callbackUrl := fmt.Sprintf("%s/api/v1/auth/callback?provider=%%s", baseUrl)
	log.Debugf("Setting provider callback url to: %s", callbackUrl)

	enabledProviders := make([]goth.Provider, 0)

	for providerName, setup := range setups {
		ok, provider := setup(fmt.Sprintf(callbackUrl, providerName))
		if ok {
			providers = append(providers, providerName)
			enabledProviders = append(enabledProviders, provider)
			log.Debugf("Registered login provider %s", providerName)
		}
	}

	goth.UseProviders(enabledProviders...)

	if !areProvidersSet() {
		return NoProvidersSet
	}

	return nil
}

func GetEnabledProviders() []string {
	return providers
}

func areProvidersSet() bool {
	return len(providers) > 0
}
