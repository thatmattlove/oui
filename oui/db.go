package oui

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gookit/gcli/v3/progress"
	"github.com/thatmattlove/go-macaddr"
	"github.com/thatmattlove/oui/internal/logger"
	"github.com/thatmattlove/oui/internal/util"
	_ "modernc.org/sqlite"
)

type OUIDB struct {
	Directory  string
	FileName   string
	Connection *sql.DB
	Version    string
}

func getFileName() (fn string, err error) {
	dir, err := os.UserConfigDir()

	if err != nil {
		return
	}
	fn = filepath.Join(dir, "oui", "oui.db")
	return
}

func scaffold() (dbf *os.File, dn string, err error) {
	fn, err := getFileName()
	if err != nil {
		return
	}
	dn = filepath.Dir(fn)

	err = os.RemoveAll(dn)
	if err != nil {
		return
	}
	err = os.MkdirAll(dn, 0755)
	if err != nil {
		return
	}
	defer dbf.Close()
	dbf, err = os.Create(fn)
	if err != nil {
		return
	}
	return
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

func (ouidb *OUIDB) Insert(d *VendorDef) (res sql.Result, err error) {
	s, err := ouidb.Connection.Prepare(fmt.Sprintf("INSERT INTO %s(prefix, length, org, registry) values(?,?,?,?)", ouidb.Version))
	if err != nil {
		return
	}
	res, err = s.Exec(d.Prefix, d.Length, d.Org, d.Registry)
	return
}

func (ouidb *OUIDB) Populate(p *progress.Progress) (records int64, err error) {
	p.AdvanceTo(11)
	err = ouidb.Clear()
	if err != nil {
		return
	}
	for _, def := range CollectAll(p, logger.New()) {
		_, err = ouidb.Insert(def)
		if err != nil {
			return
		}
		records++
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

func (ouidb *OUIDB) Find(search string) (matches []*VendorDef, err error) {
	mac, err := macaddr.ParseMACAddress(search)
	if err != nil {
		return matches, err
	}
	q := fmt.Sprintf("SELECT prefix,length,org,registry FROM %s WHERE prefix LIKE '%s%%'", ouidb.Version, mac.OUI())
	rows, err := ouidb.Connection.Query(q)
	if err != nil {
		return matches, err
	}

	defer rows.Close()

	for rows.Next() {
		var prefix string
		var length int
		var org string
		var reg string
		err = rows.Scan(&prefix, &length, &org, &reg)
		if err != nil {
			return matches, err
		}
		def := &VendorDef{Prefix: prefix, Length: length, Org: org, Registry: reg}
		_, mp, err := macaddr.ParseMACPrefix(def.PrefixString())
		if err != nil {
			return matches, err
		}
		_, failure := mp.Match(search)
		if failure == nil {
			matches = append(matches, def)
		}
	}
	return matches, nil
}

func (ouidb *OUIDB) Close() error {
	return ouidb.Connection.Close()
}

func tableExists(conn *sql.DB, ver string) (bool, error) {
	q, err := conn.Prepare(fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", ver))
	if err != nil {
		return false, err
	}
	rs, err := q.Query()
	if err != nil {
		return false, err
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

func New(version string) (ouidb *OUIDB, err error) {
	fileName, err := getFileName()
	if err != nil {
		return nil, err
	}

	var conn *sql.DB

	if !util.PathExists(fileName) {
		_, _, err = scaffold()
		if err != nil {
			return nil, err
		}
		conn, err = sql.Open("sqlite", fileName)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = sql.Open("sqlite", fileName)
		if err != nil {
			return nil, err
		}
	}

	err = conn.Ping()
	if err != nil {
		return
	}

	exists, err := tableExists(conn, version)
	if err != nil {
		return nil, err
	}
	if !exists {
		q := fmt.Sprintf("CREATE TABLE `%s` ( `id` INTEGER PRIMARY KEY AUTOINCREMENT, `prefix` VARCHAR(32) NOT NULL, `length` INTEGER NOT NULL, `org` VARCHAR(64) NOT NULL, `registry` VARCHAR(32) NOT NULL , UNIQUE(prefix, length, registry) ON CONFLICT REPLACE )", version)
		_, err = conn.Exec(q)
		if err != nil {
			return nil, err
		}
	}

	ouidb = &OUIDB{
		FileName:   fileName,
		Directory:  filepath.Dir(fileName),
		Connection: conn,
		Version:    version,
	}
	return
}
