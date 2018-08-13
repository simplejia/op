package history

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (history *History) Update() (err error) {
	c := history.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"uid":             history.Uid,
		"srv_id":          history.SrvID,
		"srv_action_path": history.SrvActionPath,
	}
	err = c.Update(q, history)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}

	return
}
