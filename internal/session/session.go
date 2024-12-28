package session

import (
	"context"
	"errors"
	"example.com/go-web-base/internal/application"
	"github.com/google/uuid"
	"sync"
	"time"
)

// 604800 seconds = 7 days
var sessionDuration = time.Duration(604800)

type cache struct {
	mutex          sync.Mutex
	cachedSessions map[string]UserSession
}

func (c *cache) Add(session UserSession) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cachedSessions[session.Id] = session
}

func (c *cache) Delete(session UserSession) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cachedSessions, session.Id)
}

func (c *cache) Get(sessionId string) (UserSession, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, found := c.cachedSessions[sessionId]
	if found {
		return v, found
	}

	return UserSession{}, false
}

var sessionCache cache

func init() {
	sessionCache = cache{
		mutex:          sync.Mutex{},
		cachedSessions: make(map[string]UserSession),
	}
}

type UserSession struct {
	Id        string
	UserId    string
	ExpiresAt time.Time
}

func (s UserSession) writeToDB(ctx context.Context, app application.Application) {
	_, err := app.DB.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)", s.Id, s.UserId, s.ExpiresAt)
	if err != nil {
		app.Log.Error(ctx, "unable to create session record: "+err.Error())
	}
}

func (s UserSession) removeFromDB(ctx context.Context, app application.Application) {
	_, err := app.DB.Exec("DELETE FROM sessions WHERE id = $1", s.Id)
	if err != nil {
		app.Log.Error(ctx, "unable to remove user session: "+err.Error())
	}
}

// Refresh will simply call GetById, Invalidate & Create
// After Refresh is called the old session is no longer valid
func (s UserSession) Refresh(ctx context.Context, app application.Application) (UserSession, error) {
	userSession, err := GetById(ctx, app, s.Id)
	if err != nil {
		app.Log.Error(ctx, "unable to get user from sessionId: "+err.Error())
		return UserSession{}, errors.New("unable to get user from sessionId: " + err.Error())
	}

	err = s.Invalidate(ctx, app)
	if err != nil {
		app.Log.Error(ctx, err.Error())
		return UserSession{}, err
	}

	newSession, err := Create(ctx, app, userSession.UserId)
	if err != nil {
		app.Log.Error(ctx, err.Error())
		return UserSession{}, err
	}

	return newSession, nil
}

func (s UserSession) Invalidate(ctx context.Context, app application.Application) error {
	sessionCache.Delete(s)

	go s.removeFromDB(ctx, app)

	return nil
}

func Create(ctx context.Context, app application.Application, userId string) (UserSession, error) {
	sessionId := uuid.New().String()
	sessionExpiration := time.Now().Add(time.Second * sessionDuration)

	session := UserSession{
		Id:        sessionId,
		UserId:    userId,
		ExpiresAt: sessionExpiration,
	}

	sessionCache.Add(session)

	go session.writeToDB(ctx, app)

	return session, nil
}

func GetById(ctx context.Context, app application.Application, sessionId string) (UserSession, error) {
	v, found := sessionCache.Get(sessionId)
	if found && v.ExpiresAt.After(time.Now()) {
		return v, nil
	} else if found {
		_ = v.Invalidate(ctx, app)
	}

	row := app.DB.QueryRow("SELECT id, user_id, expires_at FROM sessions WHERE id = $1 AND expires_at > NOW()", sessionId)

	var userSession UserSession
	err := row.Scan(&userSession.Id, &userSession.UserId, &userSession.ExpiresAt)
	if err != nil {
		app.Log.Error(ctx, "no user found for given sessionId: "+err.Error())
		return UserSession{}, err
	}

	sessionCache.Add(userSession)

	return userSession, nil
}

func PurgeOldSessionsFromDB(ctx context.Context, app application.Application) {
	_, err := app.DB.Exec("DELETE FROM sessions WHERE expires_at < NOW()")
	if err != nil {
		app.Log.Error(ctx, "unable to purge old user sessions: "+err.Error())
		return
	}

	app.Log.Info(ctx, "successfully purged old session from database")
}
