package srv

func (srv *Srv) Count() (total int, err error) {
	c := srv.GetC()
	defer c.Database.Session.Close()

	total, err = c.Find(nil).Count()
	if err != nil {
		return
	}

	return
}
