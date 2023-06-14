package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ow0sh/gotest/config"
	"github.com/pkg/errors"
)

type Row struct {
	base  string
	quote string
	rate  int
}

type PSQLConn struct {
	conn *pgx.Conn
}

func NewConn(conf config.TypePSQLConfig) (*PSQLConn, error) {
	config, err := pgx.ParseConfig(fmt.Sprintf("%v://%v:%v@%v:%v/%v", conf.Dsn, conf.User, conf.Password, conf.Host, conf.Port, conf.Dbname))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}

	return &PSQLConn{conn: conn}, nil
}

func (conn *PSQLConn) CloseConn() error {
	err := conn.conn.Close(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to close the connection")
	}
	return nil
}

func (conn *PSQLConn) InsertInfo(bases map[string]string, quotes map[string]struct{}) error {
	id := 0
	for k := range bases {
		for j := range quotes {
			insertStr := `INSERT INTO prices VALUES ($1, $2, $3, $4)`
			_, err := conn.conn.Exec(context.Background(), insertStr, id, k, j, 0)
			if err != nil {
				return errors.Wrap(err, "failed to insert info")
			}
			id += 1
		}
	}
	return nil
}

func (conn *PSQLConn) SelectInfo() (*Row, error) {
	var base, quote string
	var rate int
	err := conn.conn.QueryRow(context.Background(), `SELECT base, quote, rate FROM prices ORDER BY id DESC LIMIT 1;`).Scan(&base, &quote, &rate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select info")
	}

	return &Row{base: base, quote: quote, rate: rate}, nil
}
