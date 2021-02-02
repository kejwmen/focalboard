package auth

import (
	"time"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/services/config"
	"github.com/mattermost/focalboard/server/services/store"
	"github.com/pkg/errors"
)

// Auth authenticates sessions
type Auth struct {
	config *config.Configuration
	store  store.Store
}

// New returns a new Auth
func New(config *config.Configuration, store store.Store) *Auth {
	return &Auth{config: config, store: store}
}

// GetSession Get a user active session and refresh the session if is needed
func (a *Auth) GetSession(token string) (*model.Session, error) {
	if len(token) < 1 {
		return nil, errors.New("no session token")
	}

	session, err := a.store.GetSession(token, a.config.SessionExpireTime)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get the session for the token")
	}
	if session.UpdateAt < (time.Now().Unix() - a.config.SessionRefreshTime) {
		a.store.RefreshSession(session)
	}
	return session, nil
}