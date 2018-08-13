package srv

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

func (srv *Srv) Add() (err error) {
	c := srv.GetC()
	defer c.Database.Session.Close()

	if srv.ID <= 0 {
		var _srv *Srv
		err = c.Find(nil).Sort("-_id").Limit(1).One(&_srv)
		if err != nil {
			if err != mgo.ErrNotFound {
				return
			}
			srv.ID = 100000
			err = nil
		} else {
			srv.ID = _srv.ID + 1
		}
	}

	srv.Ct = time.Now()
	err = c.Insert(srv)
	if err != nil {
		return
	}

	return
}
