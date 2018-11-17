package srv

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/simplejia/op/mongo"
)

type FieldKind int

const (
	FieldKindNone FieldKind = iota
	FieldKindString
	FieldKindInteger
	FieldKindFloat
	FieldKindBool
	FieldKindMap
	FieldKindArray
	FieldKindFile
)

type ActionKind int

const (
	ActionKindNone ActionKind = iota
	ActionKindList
	ActionKindUpdate
	ActionKindDelete
	ActionKindCustomer
	ActionKindTransparent
)

type FieldSource int

const (
	FieldSourceNone FieldSource = iota
	FieldSourceUser
	FieldSourceArray
	FieldSourceUrl
)

type SrvAction struct {
	Path string     `json:"path" bson:"path"`
	Desc string     `json:"desc" bson:"desc"`
	Kind ActionKind `json:"kind" bson:"kind"`
}

func (srvAction *SrvAction) Regular() (ok bool) {
	if srvAction == nil {
		return
	}

	srvAction.Path = strings.TrimSpace(srvAction.Path)
	if srvAction.Path == "" {
		return
	}

	if srvAction.Kind == ActionKindNone {
		return
	}

	ok = true
	return
}

type SrvField struct {
	Name     string      `json:"name" bson:"name"`
	Required bool        `json:"required" bson:"required"`
	Kind     FieldKind   `json:"kind" bson:"kind"`
	Source   FieldSource `json:"source" bson:"source"`
	Param    string      `json:"param,omitempty" bson:"param,omitempty"`
}

func (srvField *SrvField) Regular() (ok bool) {
	if srvField == nil {
		return
	}

	srvField.Name = strings.TrimSpace(srvField.Name)
	if srvField.Name == "" {
		return
	}

	if srvField.Source == FieldSourceNone {
		return
	}

	ok = true
	return
}

type SrvActionField struct {
	Action *SrvAction  `json:"action" bson:"action"`
	Fields []*SrvField `json:"fields" bson:"fields"`
}

func (srvActionField *SrvActionField) Regular() (ok bool) {
	if srvActionField == nil {
		return
	}

	if !srvActionField.Action.Regular() {
		return
	}

	for _, field := range srvActionField.Fields {
		if !field.Regular() {
			return
		}
	}

	ok = true
	return
}

type Srv struct {
	ID           int64             `json:"id" bson:"_id"`
	MName        string            `json:"m_name" bson:"m_name"`
	SName        string            `json:"s_name" bson:"s_name"`
	Desc         string            `json:"desc" bisn:"desc"`
	Addr         string            `json:"addr" bson:"addr"`
	ActionFields []*SrvActionField `json:"action_fields" bson:"action_fields"`
	Ct           time.Time         `json:"ct" bson:"ct"`
	Ut           time.Time         `json:"ut" bson:"ut"`
}

func NewSrv() *Srv {
	return &Srv{}
}

func (srv *Srv) Db() (db string) {
	return "op"
}

func (srv *Srv) Table() (table string) {
	return "srv"
}

func (srv *Srv) GetC() (c *mgo.Collection) {
	db, table := srv.Db(), srv.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	c = sessionCopy.DB(db).C(table)
	return
}

func (srv *Srv) Regular() (ok bool) {
	if srv == nil {
		return
	}

	srv.MName = strings.TrimSpace(srv.MName)
	if srv.MName == "" {
		return
	}

	srv.SName = strings.TrimSpace(srv.SName)
	if srv.SName == "" {
		return
	}

	srv.Addr = strings.TrimSpace(srv.Addr)
	if srv.Addr == "" {
		return
	}

	for _, actionField := range srv.ActionFields {
		if !actionField.Regular() {
			return
		}
	}

	ok = true
	return
}

func (srv *Srv) ParseFromRequest(r *http.Request) (err error) {
	if srv == nil {
		err = errors.New("srv empty")
		return
	}

	srv.ID, _ = strconv.ParseInt(r.PostFormValue("id"), 10, 64)
	srv.MName = r.PostFormValue("m_name")
	srv.SName = r.PostFormValue("s_name")
	srv.Desc = r.PostFormValue("desc")
	srv.Addr = r.PostFormValue("addr")

	actionKinds := map[ActionKind]int{}
	actionPaths := map[string]int{}

	maxActionNum := 50

	for actionPos := 0; actionPos < maxActionNum; actionPos++ {
		actionPosStr := strconv.Itoa(actionPos)

		actionPath := r.PostFormValue("action_path_" + actionPosStr)
		if actionPath == "" {
			continue
		}

		if actionPaths[actionPath]++; actionPaths[actionPath] > 1 {
			err = errors.New("action path err, multi exist")
			return
		}

		actionKind, _ := strconv.Atoi(r.PostFormValue("action_kind_" + actionPosStr))
		_actionKind := ActionKind(actionKind)

		if _actionKind != ActionKindCustomer && _actionKind != ActionKindTransparent {
			if actionKinds[_actionKind]++; actionKinds[_actionKind] > 1 {
				err = errors.New("action kind err, multi exist, but customer or transparent kind is unlimit")
				return
			}
		}

		srvActionField := &SrvActionField{}
		actionDesc := r.PostFormValue("action_desc_" + actionPosStr)
		srvActionField.Action = &SrvAction{
			Path: actionPath,
			Desc: actionDesc,
			Kind: _actionKind,
		}

		fieldNames := r.PostForm["field_name_"+actionPosStr]
		fieldKinds := r.PostForm["field_kind_"+actionPosStr]
		fieldRequires := r.PostForm["field_require_"+actionPosStr]
		fieldSources := r.PostForm["field_source_"+actionPosStr]
		fieldParams := r.PostForm["field_param_"+actionPosStr]
		for fieldPos, fieldName := range fieldNames {
			if fieldName == "" {
				continue
			}

			fieldRequire, _ := strconv.ParseBool(fieldRequires[fieldPos])
			fieldKind, _ := strconv.Atoi(fieldKinds[fieldPos])
			fieldSource, _ := strconv.Atoi(fieldSources[fieldPos])
			fieldParam := fieldParams[fieldPos]
			srvActionField.Fields = append(srvActionField.Fields, &SrvField{
				Name:     fieldName,
				Required: fieldRequire,
				Kind:     FieldKind(fieldKind),
				Source:   FieldSource(fieldSource),
				Param:    fieldParam,
			})
		}

		srv.ActionFields = append(srv.ActionFields, srvActionField)
	}

	if !srv.Regular() {
		err = errors.New("input invalid")
		return
	}

	return
}
