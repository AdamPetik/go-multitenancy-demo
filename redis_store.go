package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"encoding/json"
)

var ctx = context.Background()

type redisStore struct {
	db *redis.Client
	ctx context.Context
}

func NewRedisStore(db *redis.Client, ctx context.Context) Store {
	return &redisStore{db, ctx}
}

func (r *redisStore) Get(key string) (Session, error) {
	raw, err := r.db.Get(r.ctx, key).Bytes()
    if err != nil {
       return nil, err
    }
	var s session
    err = json.Unmarshal(raw, &s)
	return &s, err
}

func (r *redisStore) Save(key string, session Session) error {
    p, err := json.Marshal(session)
	s := string(p)
    if err != nil {
       return err
    }
    err = r.db.Set(r.ctx, key, s, 0).Err()
	return err
}