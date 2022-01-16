package model

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	sqlxo "github.com/jmoiron/sqlx"
	"github.com/tal-tech/go-zero/core/stores/builder"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
)

var (
	bookmsgFieldNames          = builder.RawFieldNames(&Bookmsg{})
	bookmsgRows                = strings.Join(bookmsgFieldNames, ",")
	bookmsgRowsExpectAutoSet   = strings.Join(stringx.Remove(bookmsgFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	bookmsgRowsWithPlaceHolder = strings.Join(stringx.Remove(bookmsgFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheBookmsgIdPrefix     = "cache:bookmsg:id:"
	cacheBookmsgBookIdPrefix = "cache:bookmsg:bookId:"
)

type (
	BookmsgModel interface {
		Insert(data *Bookmsg) (sql.Result, error)
		FindOne(id int64) (*Bookmsg, error)
		FindOneByBookId(bookId string) (*Bookmsg, error)
		Update(data *Bookmsg) error
		Delete(id int64) error
		SearchByName(name string) ([]*Bookmsg, error)
		FindAllMsg([]string) ([]*Bookmsg, error)
	}

	defaultBookmsgModel struct {
		sqlc.CachedConn
		table string
	}

	Bookmsg struct {
		Id     int64  `db:"id"`
		BookId string `db:"book_id"`
		Name   string `db:"name"`
		Count  int64  `db:"count"`
	}
)

func NewBookmsgModel(conn sqlx.SqlConn, c cache.CacheConf) BookmsgModel {
	return &defaultBookmsgModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`bookmsg`",
	}
}

func (m *defaultBookmsgModel) FindAllMsg(bookids []string) ([]*Bookmsg, error) {
	var bookmsgs []*Bookmsg

	query, args, err := sqlxo.In(fmt.Sprintf("select %s from %s where `book_id` in (?)", bookmsgRows, m.table), bookids)
	if err != nil {
		log.Println("FindAllMsg err", err)
		return nil, err
	}
	log.Println(query)
	log.Println(args)
	// 他要的是空接口切片  args要加上...
	err = m.QueryRowsNoCache(&bookmsgs, query, args...)
	switch err {
	case nil:
		return bookmsgs, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultBookmsgModel) SearchByName(name string) ([]*Bookmsg, error) {
	var bookmsgs []*Bookmsg

	query := fmt.Sprintf("select %s from %s where `name` = ?", bookmsgRows, m.table)
	err := m.QueryRowsNoCache(&bookmsgs, query, name)

	switch err {
	case nil:
		return bookmsgs, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultBookmsgModel) Insert(data *Bookmsg) (sql.Result, error) {
	bookmsgIdKey := fmt.Sprintf("%s%v", cacheBookmsgIdPrefix, data.Id)
	bookmsgBookIdKey := fmt.Sprintf("%s%v", cacheBookmsgBookIdPrefix, data.BookId)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, bookmsgRowsExpectAutoSet)
		return conn.Exec(query, data.BookId, data.Name, data.Count)
	}, bookmsgBookIdKey, bookmsgIdKey)
	return ret, err
}

func (m *defaultBookmsgModel) FindOne(id int64) (*Bookmsg, error) {
	bookmsgIdKey := fmt.Sprintf("%s%v", cacheBookmsgIdPrefix, id)
	var resp Bookmsg
	err := m.QueryRow(&resp, bookmsgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", bookmsgRows, m.table)
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

func (m *defaultBookmsgModel) FindOneByBookId(bookId string) (*Bookmsg, error) {
	bookmsgBookIdKey := fmt.Sprintf("%s%v", cacheBookmsgBookIdPrefix, bookId)
	var resp Bookmsg
	err := m.QueryRowIndex(&resp, bookmsgBookIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `book_id` = ? limit 1", bookmsgRows, m.table)
		if err := conn.QueryRow(&resp, query, bookId); err != nil {
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

func (m *defaultBookmsgModel) Update(data *Bookmsg) error {
	bookmsgIdKey := fmt.Sprintf("%s%v", cacheBookmsgIdPrefix, data.Id)
	bookmsgBookIdKey := fmt.Sprintf("%s%v", cacheBookmsgBookIdPrefix, data.BookId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, bookmsgRowsWithPlaceHolder)
		return conn.Exec(query, data.BookId, data.Name, data.Count, data.Id)
	}, bookmsgBookIdKey, bookmsgIdKey)
	return err
}

func (m *defaultBookmsgModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	bookmsgIdKey := fmt.Sprintf("%s%v", cacheBookmsgIdPrefix, id)
	bookmsgBookIdKey := fmt.Sprintf("%s%v", cacheBookmsgBookIdPrefix, data.BookId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, bookmsgIdKey, bookmsgBookIdKey)
	return err
}

func (m *defaultBookmsgModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheBookmsgIdPrefix, primary)
}

func (m *defaultBookmsgModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", bookmsgRows, m.table)
	return conn.QueryRow(v, query, primary)
}
