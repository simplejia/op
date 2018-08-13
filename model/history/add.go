package history

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (history *History) Add() (err error) {
	c := history.GetC()
	defer c.Database.Session.Close()

	up := bson.M{
		"$push": bson.M{
			"details": bson.M{
				"$each":     history.Details,
				"$slice":    60,
				"$position": 0,
			},
		},
		"$set": bson.M{
			"ut": time.Now(),
		},
	}

	q := bson.M{
		"uid":             history.Uid,
		"srv_id":          history.SrvID,
		"srv_action_path": history.SrvActionPath,
	}
	_, err = c.Upsert(q, up)
	if err != nil {
		return
	}

	return
}
