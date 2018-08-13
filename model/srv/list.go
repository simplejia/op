package srv

func (srv *Srv) List(offset, limit int) (srvs []*Srv, err error) {
	c := srv.GetC()
	defer c.Database.Session.Close()

	err = c.Find(nil).Sort("_id").Skip(offset).Limit(limit).All(&srvs)
	if err != nil {
		return
	}

	return
}
