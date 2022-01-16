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
	lendFieldNames          = builder.RawFieldNames(&Lend{})
	lendRows                = strings.Join(lendFieldNames, ",")
	lendRowsExpectAutoSet   = strings.Join(stringx.Remove(lendFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	lendRowsWithPlaceHolder = strings.Join(stringx.Remove(lendFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheLendIdPrefix        = "cache:lend:id:"
	cacheLendBookidUidPrefix = "cache:lend:bookid:uid:"
)

type (
	LendModel interface {
		Insert(data *Lend) (sql.Result, error)
		FindOne(id int64) (*Lend, error)
		FindOneByBookidUid(bookid string, uid int64) (*Lend, error)
		Update(data *Lend) error
		Delete(id int64) error
		SearchAllMsgByUid(uid int64) ([]string, error)
		SearchHowManyBook(uid int64) (int64, error)
		SearchWhoLookBook(bookid string) ([]*int64, error)
	}

	defaultLendModel struct {
		sqlc.CachedConn
		table string
	}

	Lend struct {
		Bookid string `db:"bookid"`
		Uid    int64  `db:"uid"`
		Id     int64  `db:"id"`
	}
)

func NewLendModel(conn sqlx.SqlConn, c cache.CacheConf) LendModel {
	return &defaultLendModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`lend`",
	}
}
func (m *defaultLendModel) SearchWhoLookBook(bookid string) ([]*int64, error) {
	var nums []*int64
	query := fmt.Sprintf("select `uid` from %s where `bookid` = ?", m.table)
	err := m.QueryRowsNoCache(&nums, query, bookid)

	switch err {
	case nil:
		return nums, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLendModel) SearchHowManyBook(uid int64) (int64, error) {
	var num int64
	query := fmt.Sprintf("select count(uid) from %s group by `uid` having uid = ?", m.table)
	err := m.QueryRowNoCache(&num, query, uid)
	switch err {
	case nil:
		return num, nil
	case sqlc.ErrNotFound:
		return -1, ErrNotFound
	default:
		return -1, err
	}
}

func (m *defaultLendModel) SearchAllMsgByUid(uid int64) ([]string, error) {
	var bookids []string

	query := fmt.Sprintf("select %s from %s where uid = ?", "bookid", m.table)
	err := m.QueryRowsNoCache(&bookids, query, uid)

	switch err {
	case nil:
		return bookids, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLendModel) Insert(data *Lend) (sql.Result, error) {
	lendIdKey := fmt.Sprintf("%s%v", cacheLendIdPrefix, data.Id)
	lendBookidUidKey := fmt.Sprintf("%s%v:%v", cacheLendBookidUidPrefix, data.Bookid, data.Uid)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, lendRowsExpectAutoSet)
		return conn.Exec(query, data.Bookid, data.Uid)
	}, lendIdKey, lendBookidUidKey)
	return ret, err
}

func (m *defaultLendModel) FindOne(id int64) (*Lend, error) {
	lendIdKey := fmt.Sprintf("%s%v", cacheLendIdPrefix, id)
	var resp Lend
	err := m.QueryRow(&resp, lendIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", lendRows, m.table)
		return conn.QueryRow(v, query, id)
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

func (m *defaultLendModel) FindOneByBookidUid(bookid string, uid int64) (*Lend, error) {
	lendBookidUidKey := fmt.Sprintf("%s%v:%v", cacheLendBookidUidPrefix, bookid, uid)
	var resp Lend
	err := m.QueryRowIndex(&resp, lendBookidUidKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `bookid` = ? and `uid` = ? limit 1", lendRows, m.table)
		if err := conn.QueryRow(&resp, query, bookid, uid); err != nil {
			return nil, err
		}
		return resp.Id, nil
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

func (m *defaultLendModel) Update(data *Lend) error {
	lendIdKey := fmt.Sprintf("%s%v", cacheLendIdPrefix, data.Id)
	lendBookidUidKey := fmt.Sprintf("%s%v:%v", cacheLendBookidUidPrefix, data.Bookid, data.Uid)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, lendRowsWithPlaceHolder)
		return conn.Exec(query, data.Bookid, data.Uid, data.Id)
	}, lendIdKey, lendBookidUidKey)
	return err
}

func (m *defaultLendModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	lendIdKey := fmt.Sprintf("%s%v", cacheLendIdPrefix, id)
	lendBookidUidKey := fmt.Sprintf("%s%v:%v", cacheLendBookidUidPrefix, data.Bookid, data.Uid)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, lendIdKey, lendBookidUidKey)
	return err
}

func (m *defaultLendModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheLendIdPrefix, primary)
}

func (m *defaultLendModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", lendRows, m.table)
	return conn.QueryRow(v, query, primary)
}
