package auth

import (
	"fmt"
	"github.com/KowalskiPiotr98/ludivault/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/twitch"
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
		"twitch": func(callbackUrl string) (bool, goth.Provider) {
			clientId := utils.GetOptionalConfig("SSO_TWITCH_CLIENT_ID", "")
			clientSecret := utils.GetOptionalConfig("SSO_TWITCH_CLIENT_SECRET", "")

			if clientId == "" || clientSecret == "" {
				return false, nil
			}

			return true, twitch.New(clientId, clientSecret, callbackUrl)
		},
		"github": func(callbackUrl string) (bool, goth.Provider) {
			clientId := utils.GetOptionalConfig("SSO_GITHUB_CLIENT_ID", "")
			clientSecret := utils.GetOptionalConfig("SSO_GITHUB_CLIENT_SECRET", "")

			if clientId == "" || clientSecret == "" {
				return false, nil
			}

			return true, github.New(clientId, clientSecret, callbackUrl)
		},
		"discord": func(callbackUrl string) (bool, goth.Provider) {
			clientId := utils.GetOptionalConfig("SSO_DISCORD_CLIENT_ID", "")
			clientSecret := utils.GetOptionalConfig("SSO_DISCORD_CLIENT_SECRET", "")

			if clientId == "" || clientSecret == "" {
				return false, nil
			}

			return true, discord.New(clientId, clientSecret, callbackUrl)
		},
		"google": func(callbackUrl string) (bool, goth.Provider) {
			clientId := utils.GetOptionalConfig("SSO_GOOGLE_CLIENT_ID", "")
			clientSecret := utils.GetOptionalConfig("SSO_GOOGLE_CLIENT_SECRET", "")

			if clientId == "" || clientSecret == "" {
				return false, nil
			}

			return true, google.New(clientId, clientSecret, callbackUrl)
		},
		"steam": func(callbackUrl string) (bool, goth.Provider) {
			clientSecret := utils.GetOptionalConfig("SSO_STEAM_CLIENT_SECRET", "")

			if clientSecret == "" {
				return false, nil
			}

			return true, steam.New(clientSecret, callbackUrl)
		},
		"oidc": func(callbackUrl string) (bool, goth.Provider) {
			clientId := utils.GetOptionalConfig("SSO_OIDC_CLIENT_ID", "")
			clientSecret := utils.GetOptionalConfig("SSO_OIDC_CLIENT_SECRET", "")
			discoveryUrl := utils.GetOptionalConfig("SSO_OIDC_DISCOVERY_URL", "")

			if clientId == "" || clientSecret == "" || discoveryUrl == "" {
				return false, nil
			}

			provider, err := openidConnect.New(clientId, clientSecret, discoveryUrl, callbackUrl)
			if err != nil {
				log.Warnf("Failed to initialize OpenID Connect provider: %v", err)
				return false, nil
			}

			return true, provider
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
