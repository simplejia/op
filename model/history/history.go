package history

import (
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/op/mongo"
)

type HistoryDetail struct {
	M  map[string]string
	Ct time.Time
}

type History struct {
	ID            bson.ObjectId    `json:"id" bson:"_id,omitempty"`
	Uid           int64            `json:"uid" bson:"uid"`
	SrvID         int64            `json:"srv_id,omitempty" bson:"srv_id,omitempty"`
	SrvActionPath string           `json:"srv_action_path,omitempty" bson:"srv_action_path,omitempty"`
	Details       []*HistoryDetail `json:"details,omitempty" bson:"details,omitempty"`
	Ut            time.Time        `json:"ut" bson:"ut"`
}

func NewHistory() *History {
	return &History{}
}

func (history *History) Db() (db string) {
	return "op"
}

func (history *History) Table() (table string) {
	return "history"
}

func (history *History) GetC() (c *mgo.Collection) {
	db, table := history.Db(), history.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	c = sessionCopy.DB(db).C(table)
	return
}
