package main

type Store interface {
	Get(key string) (Session, error)
	Save(key string, session Session) error
}






