package history

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (history *History) GetByUidAndSrv() (historyRet *History, err error) {
	c := history.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"uid":             history.Uid,
		"srv_id":          history.SrvID,
		"srv_action_path": history.SrvActionPath,
	}
	err = c.Find(q).One(&historyRet)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}

	return
}
