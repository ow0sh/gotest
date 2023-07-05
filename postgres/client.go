package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ow0sh/gotest/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

func (conn *PSQLConn) InsertInfo(log *logrus.Logger, base string, quote string, rate float64) error {
	id := 0
	selectStr := `SELECT id FROM prices ORDER BY id DESC LIMIT 1;`
	err := conn.conn.QueryRow(context.Background(), selectStr).Scan(&id)
	if err != nil {
		log.Error("failed to get last id, creating new one")
	}

	id++
	insertStr := "INSERT INTO prices VALUES ($1, $2, $3, $4);"
	_, err = conn.conn.Exec(context.Background(), insertStr, id, base, quote, rate)
	if err != nil {
		log.Error("failed to insert into db")
		return errors.Wrap(err, "failed to insert into db")
	}

	log.Info("Inserted into DB successfully")
	return nil
}

func (conn *PSQLConn) SelectInfo() (*Row, error) {
	var base, quote string
	var rate int
	selectStr := `SELECT base, quote, rate FROM prices ORDER BY id DESC LIMIT 1;`
	err := conn.conn.QueryRow(context.Background(), selectStr).Scan(&base, &quote, &rate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select info")
	}

	return &Row{base: base, quote: quote, rate: rate}, nil
}

func (conn *PSQLConn) UpdateInfo(log *logrus.Logger, base string, rate float64) error {
	updateStr := "UPDATE prices SET rate = $1 WHERE base_asset = $2;"
	_, err := conn.conn.Exec(context.Background(), updateStr, rate, base)
	if err != nil {
		log.Error("failed to update info")
		return errors.Wrap(err, "failed to update info")
	}
	log.Info("Updated rates successfully")
	return nil
}

func (conn *PSQLConn) Exist(base string) bool {
	exitstStr := "SELECT id FROM prices WHERE base_asset = $1;"
	var tmp int
	conn.conn.QueryRow(context.Background(), exitstStr, base).Scan(&tmp)

	if tmp == 0 {
		return false
	}
	return true
}
