package srv

func (srv *Srv) Delete() (err error) {
	c := srv.GetC()
	defer c.Database.Session.Close()

	err = c.RemoveId(srv.ID)
	if err != nil {
		return
	}

	return
}
