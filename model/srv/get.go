package srv

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (srv *Srv) Get() (srvRet *Srv, err error) {
	c := srv.GetC()
	defer c.Database.Session.Close()

	err = c.FindId(srv.ID).One(&srvRet)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}
		err = nil
		return
	}

	return
}

func (srv *Srv) GetByName() (srvRet *Srv, err error) {
	c := srv.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"m_name": srv.MName,
		"s_name": srv.SName,
	}
	err = c.Find(q).One(&srvRet)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}
		err = nil
		return
	}

	return
}
