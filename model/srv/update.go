package srv

import (
	"time"
)

func (srv *Srv) Update() (err error) {
	c := srv.GetC()
	defer c.Database.Session.Close()

	srv.Ut = time.Now()
	err = c.UpdateId(srv.ID, srv)
	if err != nil {
		return
	}

	return
}
