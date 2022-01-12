package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)

var (
	usermsgFieldNames          = builder.RawFieldNames(&Usermsg{})
	usermsgRows                = strings.Join(usermsgFieldNames, ",")
	usermsgRowsExpectAutoSet   = strings.Join(stringx.Remove(usermsgFieldNames, "`uid`", "`create_time`", "`update_time`"), ",")
	usermsgRowsWithPlaceHolder = strings.Join(stringx.Remove(usermsgFieldNames, "`uid`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheUsermsgUidPrefix     = "cache:usermsg:uid:"
	cacheUsermsgAccountPrefix = "cache:usermsg:account:"
)

type (
	UsermsgModel interface {
		Insert(data *Usermsg) (sql.Result, error)
		FindOne(uid int64) (*Usermsg, error)
		FindOneByAccount(account string) (*Usermsg, error)
		Update(data *Usermsg) error
		Delete(uid int64) error
	}

	defaultUsermsgModel struct {
		sqlc.CachedConn
		table string
	}

	Usermsg struct {
		Uid      int64  `db:"uid"`
		Name     string `db:"name"`
		Gender   string `db:"gender"`
		Phone    string `db:"phone"`
		Address  string `db:"address"`
		Email    string `db:"email"`
		Account  string `db:"account"`
		Password string `db:"password"`
	}
)

func NewUsermsgModel(conn sqlx.SqlConn, c cache.CacheConf) UsermsgModel {
	return &defaultUsermsgModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`usermsg`",
	}
}

func (m *defaultUsermsgModel) Insert(data *Usermsg) (sql.Result, error) {
	usermsgAccountKey := fmt.Sprintf("%s%v", cacheUsermsgAccountPrefix, data.Account)
	usermsgUidKey := fmt.Sprintf("%s%v", cacheUsermsgUidPrefix, data.Uid)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, usermsgRowsExpectAutoSet)
		return conn.Exec(query, data.Name, data.Gender, data.Phone, data.Address, data.Email, data.Account, data.Password)
	}, usermsgUidKey, usermsgAccountKey)
	return ret, err
}

func (m *defaultUsermsgModel) FindOne(uid int64) (*Usermsg, error) {
	usermsgUidKey := fmt.Sprintf("%s%v", cacheUsermsgUidPrefix, uid)
	var resp Usermsg
	err := m.QueryRow(&resp, usermsgUidKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", usermsgRows, m.table)
		return conn.QueryRow(v, query, uid)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUsermsgModel) FindOneByAccount(account string) (*Usermsg, error) {
	usermsgAccountKey := fmt.Sprintf("%s%v", cacheUsermsgAccountPrefix, account)
	var resp Usermsg
	err := m.QueryRowIndex(&resp, usermsgAccountKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `account` = ? limit 1", usermsgRows, m.table)
		if err := conn.QueryRow(&resp, query, account); err != nil {
			return nil, err
		}
		return resp.Uid, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUsermsgModel) Update(data *Usermsg) error {
	usermsgUidKey := fmt.Sprintf("%s%v", cacheUsermsgUidPrefix, data.Uid)
	usermsgAccountKey := fmt.Sprintf("%s%v", cacheUsermsgAccountPrefix, data.Account)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `uid` = ?", m.table, usermsgRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Gender, data.Phone, data.Address, data.Email, data.Account, data.Password, data.Uid)
	}, usermsgUidKey, usermsgAccountKey)
	return err
}

func (m *defaultUsermsgModel) Delete(uid int64) error {
	data, err := m.FindOne(uid)
	if err != nil {
		return err
	}

	usermsgUidKey := fmt.Sprintf("%s%v", cacheUsermsgUidPrefix, uid)
	usermsgAccountKey := fmt.Sprintf("%s%v", cacheUsermsgAccountPrefix, data.Account)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `uid` = ?", m.table)
		return conn.Exec(query, uid)
	}, usermsgUidKey, usermsgAccountKey)
	return err
}

func (m *defaultUsermsgModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUsermsgUidPrefix, primary)
}

func (m *defaultUsermsgModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `uid` = ? limit 1", usermsgRows, m.table)
	return conn.QueryRow(v, query, primary)
}
