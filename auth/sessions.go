package auth

import (
	"github.com/KowalskiPiotr98/ludivault/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//todo: storing sessions in cookie only is bad as it prevents any form of server-side tracking of sessions
// but it'll have to do for now

var (
	UserSessionName = "ludivault-session"
	SetupSession    = func(baseDomain string) sessions.Store {
		key := utils.GetRequiredConfig("SESSION_KEY")
		maxAge := 86400 * 30
		store := sessions.NewCookieStore([]byte(key))
		store.MaxAge(maxAge)
		store.Options.Path = "/"
		store.Options.HttpOnly = true
		store.Options.Secure = gin.Mode() != gin.DebugMode
		store.Options.Domain = baseDomain
		store.Options.SameSite = http.SameSiteLaxMode
		return store
	}

	authStore sessions.Store
)

// InitSessionStore initialises the store for user sessions.
//
// This function must be called before any other function from this package.
func InitSessionStore(baseDomain string) {
	authStore = SetupSession(baseDomain)
	gothic.Store = authStore
}

// StoreUserInSession adds the user data to session store.
func StoreUserInSession(c *gin.Context, userId uuid.UUID) error {
	ensureSessionStoreInit()

	session, err := authStore.Get(c.Request, UserSessionName)
	if err != nil {
		log.Warnf("Failed to get session: %v", err)
		return err
	}

	session.Values["userId"] = userId

	if err = session.Save(c.Request, c.Writer); err != nil {
		log.Warnf("Failed to save session: %v", err)
		return err
	}

	return nil
}

// RetrieveUserFromSession attempts to get user data from session store.
func RetrieveUserFromSession(c *gin.Context) (uuid.UUID, error) {
	ensureSessionStoreInit()

	session, err := authStore.Get(c.Request, UserSessionName)
	if err != nil {
		log.Warnf("Failed to get session: %v", err)
		return uuid.Nil, err
	}

	id, ok := session.Values["userId"].(uuid.UUID)
	if !ok {
		return uuid.Nil, UserIdNotStored
	}
	return id, nil
}

func RemoveUserSession(c *gin.Context) error {
	ensureSessionStoreInit()

	session, err := authStore.Get(c.Request, UserSessionName)
	if err != nil {
		log.Warnf("Failed to get session: %v", err)
		return err
	}

	// remove what's left of the session
	session.Options.MaxAge = -1
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Warnf("Failed to save session: %v", err)
		return err
	}

	return nil
}

func ensureSessionStoreInit() {
	if authStore == nil {
		log.Panic("Session store is not initialized")
	}
}
