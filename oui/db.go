package oui

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gookit/gcli/v3/progress"
	"github.com/thatmattlove/go-macaddr"
	"github.com/thatmattlove/oui/v2/internal/util"
)

type OUIDB struct {
	Connection  *sql.DB
	Version     string
	Progress    *progress.Progress
	Logger      *LoggerType
	useLogging  bool
	useProgress bool
	dialect     int
}

func tableExists(dialect int, conn *sql.DB, ver string) (bool, error) {
	q, err := conn.Prepare(fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", ver))
	if err != nil {
		return false, err
	}
	rs, err := q.Query()
	if err != nil {
		return false, nil
	}
	defer rs.Close()

	var table string
	rs.Next()
	rs.Scan(&table)
	if table != "" {
		return true, nil
	}
	return false, nil
}

func (ouidb *OUIDB) Clear() (err error) {
	err = ouidb.Connection.Ping()
	if err != nil {
		return
	}
	query, err := ouidb.Connection.Prepare(fmt.Sprintf("DELETE FROM %s", ouidb.Version))
	if err != nil {
		return
	}
	_, err = query.Exec()
	return
}

func (ouidb *OUIDB) Delete() error {
	err := ouidb.Connection.Ping()
	if err != nil {
		return err
	}
	query, err := ouidb.Connection.Prepare(fmt.Sprintf("DROP TABLE %s", ouidb.Version))
	if err != nil {
		return err
	}
	_, err = query.Exec()
	return err
}

func (ouidb *OUIDB) Insert(d *VendorDef) (res sql.Result, err error) {
	var q string
	switch ouidb.dialect {
	case dialectSqlite:
		q = fmt.Sprintf("INSERT INTO %s(prefix, length, org, registry) values(?,?,?,?)", ouidb.Version)
	case dialectPsql:
		q = fmt.Sprintf("INSERT INTO %s(prefix, length, org, registry) values($1,$2,$3,$4) ON CONFLICT (prefix, length, registry) DO UPDATE SET prefix = excluded.prefix, length = excluded.length, registry = excluded.registry", ouidb.Version)
	}
	s, err := ouidb.Connection.Prepare(q)
	if err != nil {
		return
	}
	res, err = s.Exec(d.Prefix, d.Length, d.Org, d.Registry)
	return
}

func (ouidb *OUIDB) BulkInsert(defs []*VendorDef) (int64, error) {

	var statement string
	var tmpl string
	var maxRecords int
	switch ouidb.dialect {
	case dialectSqlite:
		tmpl = "(?,?,?,?)"
		statement = "INSERT INTO %s(prefix, length, org, registry) values%s"
		maxRecords = maxVarsSqlite
	case dialectPsql:
		tmpl = "($%d,$%d,$%d,$%d)"
		statement = "INSERT INTO %s(prefix, length, org, registry) values%s ON CONFLICT (prefix, length, registry) DO UPDATE SET prefix = excluded.prefix, length = excluded.length, registry = excluded.registry"
		maxRecords = maxVarsPsql
	}

	splitDefs := util.SplitSlice(defs, maxRecords/4)
	inserted := int64(0)

	for _, split := range splitDefs {
		placeholders := make([]string, 0, len(split))
		args := make([]interface{}, 0, len(split)*4)
		i := 0
		for _, def := range split {
			def := def
			var placeholder string
			switch ouidb.dialect {
			case dialectSqlite:
				placeholder = tmpl
			case dialectPsql:
				placeholder = fmt.Sprintf(tmpl, i*4+1, i*4+2, i*4+3, i*4+4)
			}
			placeholders = append(placeholders, placeholder)
			args = append(args, def.Prefix, def.Length, def.Org, def.Registry)
			i++
		}
		q := fmt.Sprintf(statement, ouidb.Version, strings.Join(placeholders, ","))
		res, err := ouidb.Connection.Exec(q, args...)
		if err != nil {
			return inserted, err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return inserted, err
		}
		inserted += rows
		if err != nil {
			return inserted, err
		}
	}
	return inserted, nil
}

func (ouidb *OUIDB) Populate() (records int64, err error) {
	err = ouidb.Clear()
	if err != nil {
		return
	}
	var p *progress.Progress = nil
	var l LoggerType = nil
	if ouidb.useLogging {
		l = *ouidb.Logger
	}
	if ouidb.useProgress {
		p = ouidb.Progress
	}
	defs, err := CollectAll(p, l)
	if err != nil {
		return
	}
	records, err = ouidb.BulkInsert(defs)
	if err != nil {
		return
	}
	return
}

func (ouidb *OUIDB) Count() (count int64, err error) {
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s", ouidb.Version)
	rows, err := ouidb.Connection.Query(q)
	if err != nil {
		return
	}
	var countS string
	for rows.Next() {
		err = rows.Scan(&countS)
		if err != nil {
			return
		}
	}
	count, err = strconv.ParseInt(countS, 10, 64)
	return
}

func (ouidb *OUIDB) Find(search string) ([]*VendorDef, error) {
	mac, err := macaddr.ParseMACAddress(search)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("SELECT prefix,length,org,registry FROM %s WHERE prefix LIKE '%s%%'", ouidb.Version, mac.OUI())
	rows, err := ouidb.Connection.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	errs := make([]error, 0)
	matches := make([]*VendorDef, 0)

	for rows.Next() {
		var prefix string
		var length int
		var org string
		var reg string
		err = rows.Scan(&prefix, &length, &org, &reg)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		def := &VendorDef{Prefix: prefix, Length: length, Org: org, Registry: reg}
		_, mp, err := macaddr.ParseMACPrefix(def.PrefixString())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		_, failure := mp.Match(search)
		if failure == nil {
			matches = append(matches, def)
		}
	}
	err = errors.Join(errs...)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (ouidb *OUIDB) Close() error {
	return ouidb.Connection.Close()
}

func New(opts ...Option) (*OUIDB, error) {
	options := getOptions(opts...)
	err := options.Connection.Ping()
	if err != nil {
		return nil, err
	}
	if options.dialect == dialectSqlite {
		exists, err := tableExists(options.dialect, options.Connection, options.Version)
		if err != nil {
			return nil, err
		}
		if !exists {
			q := fmt.Sprintf("CREATE TABLE `%s` ( `id` INTEGER PRIMARY KEY AUTOINCREMENT, `prefix` VARCHAR(32) NOT NULL, `length` INTEGER NOT NULL, `org` VARCHAR(256) NOT NULL, `registry` VARCHAR(32) NOT NULL , UNIQUE(prefix, length, registry) ON CONFLICT REPLACE )", options.Version)
			_, err := options.Connection.Exec(q)
			if err != nil {
				return nil, err
			}
		}
	} else if options.dialect == dialectPsql {
		q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v ( id SERIAL PRIMARY KEY, prefix VARCHAR(32) NOT NULL, length INT NOT NULL, org VARCHAR(256) NOT NULL, registry VARCHAR(32) NOT NULL, UNIQUE(prefix, length, registry) )", options.Version)
		options.Connection.SetMaxOpenConns(int(options.MaxConnections))
		_, err := options.Connection.Exec(q)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unknown SQL dialect")
	}

	indexQ := fmt.Sprintf("CREATE INDEX IF NOT EXISTS prefix_idx ON %v (prefix)", options.Version)
	_, err = options.Connection.Exec(indexQ)
	if err != nil {
		return nil, err
	}

	ouidb := &OUIDB{
		Connection:  options.Connection,
		Version:     options.Version,
		Logger:      options.Logger,
		Progress:    options.Progress,
		useLogging:  options.Logger != nil,
		useProgress: options.Progress != nil,
		dialect:     options.dialect,
	}
	return ouidb, nil
}
