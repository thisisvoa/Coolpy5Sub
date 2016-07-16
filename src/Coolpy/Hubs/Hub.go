package Hubs

import (
	"github.com/garyburd/redigo/redis"
	"Coolpy/Incr"
	"encoding/json"
	"strconv"
)

type Hub struct {
	Id        int64
	Ukey      string
	Title     string `validate:"required"`
	About     string
	Tabs      string
	Public    bool
	Local     string `validate:"required"`
	Latitude  float64 `validate:"gte=-90,lte=90"`
	Longitude float64 `validate:"gte=-180,lte=180"`
}

var rds redis.Conn

func Connect(addr string, pwd string) {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	_, err = c.Do("AUTH", pwd)
	if err != nil {
		panic(err)
	}
	rds = c
	rds.Do("SELECT", "2")
}

func HubCreate(hub *Hub) error {
	v, err := Incr.HubInrc()
	if err != nil {
		return err
	}
	hub.Id = v
	hub.Public = false
	json, err := json.Marshal(hub)
	if err != nil {
		return err
	}
	key := hub.Ukey + ":" + strconv.FormatInt(hub.Id, 10)
	_, err = rds.Do("SET", key, json)
	if err != nil {
		return err
	}
	return nil
}

func HubStartWith(k string) ([]*Hub, error) {
	data, err := redis.Strings(rds.Do("KEYSSTART", k))
	if err != nil {
		return nil, err
	}
	var ndata []*Hub
	for _, v := range data {
		o, _ := redis.String(rds.Do("GET", v))
		h := &Hub{}
		json.Unmarshal([]byte(o), &h)
		h.Ukey = ""
		ndata = append(ndata, h)
	}
	return ndata, nil
}