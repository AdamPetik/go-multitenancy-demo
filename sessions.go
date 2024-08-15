package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Session interface {
	GetID() string
	Get(name string) interface{}
	Set(name string, data interface{})
	Save()
	setStore(store Store)
}

const defaultKey = "session_data"
const cookieName = "sessionid"

type session struct {
	ID    string
	store Store
	data  map[string]interface{}
}

func NewSession(store Store) Session {
	id := uuid.New()
	return &session{id.String(), store, make(map[string]interface{})}
}

func (s *session) GetID() string {
	return s.ID
}

func (s *session) Get(name string) interface{} {
	return s.data[name]
}

func (s *session) Save() {
	s.store.Save(s.ID, s)
}

func (s *session) setStore(store Store) {
	s.store = store;
}


func (s *session) Set(name string, data interface{}) {
	s.data[name] = data
}

func GetSession(c *gin.Context) Session {
	s, exists := c.Get(defaultKey)
	if !exists {
		panic("session does not exists, have you forgotten to add middleware?")
	}
	return s.(Session)
}

func getDomain(c *gin.Context) string {
	return strings.Split(c.Request.Host, ":")[0]
}

func SessionMiddleware(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionCookie, err := c.Cookie(cookieName)

		if err != nil {
			log.Println(err.Error())
		}

		var s Session
		domain := getDomain(c)

		if err != nil {
			log.Println("No session in cookies, creating a new one")
			s = NewSession(store)
		} else {
			s, err = store.Get(sessionCookie)

			if err != nil {
				// TODO it is not correct to handle it like that
				log.Println("No session in store, creating a new one")
				s = NewSession(store)
			}
		}

		s.setStore(store)
		c.Set(defaultKey, s)
		c.SetCookie(cookieName, s.GetID(), 3600, "/", domain, false, false)

		c.Next()

		s.Save()
		log.Println(domain)
	}
}
