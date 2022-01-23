package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gookit/gcli/v3/progress"
	"github.com/thatmattlove/go-macaddr"
	_ "modernc.org/sqlite"
)

type OUIDB struct {
	Directory  string
	FileName   string
	File       *os.File
	Connection *sql.DB
	Version    string
}

func NewOUIDB() (ouidb *OUIDB, err error) {
	ouidb = &OUIDB{}
	dir, fileName, err := ouidb.getFileName()
	if err != nil {
		return
	}
	ouidb.FileName = fileName
	ouidb.Directory = dir
	file, err := ouidb.getFile(fileName, dir)
	if err != nil {
		return
	}
	ouidb.File = file
	conn, err := ouidb.getConnection()
	if err != nil {
		return
	}
	ouidb.Connection = conn
	_, err = ouidb.getVersion()
	if err != nil {
		err = ouidb.createTable()
		return
	}
	ver, err := ouidb.getVersion()
	if err != nil {
		return
	}
	ouidb.Version = ver
	return
}

func (ouidb *OUIDB) getFileName() (d, f string, err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	d = fmt.Sprintf("%s/oui", dir)
	f = d + "/oui.db"
	return
}

func (ouidb *OUIDB) getFile(dn, fn string) (f *os.File, err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	if !pathExists(dn) {
		err = os.Mkdir(dn, 0755)
		if err != nil {
			return
		}
	}
	if !pathExists(fn) {
		f, err = os.Create(fn)
		if err != nil {
			return
		}
	}
	f, err = os.Open(fn)
	return
}

func (ouidb *OUIDB) getConnection() (conn *sql.DB, err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	conn, err = sql.Open("sqlite", ouidb.FileName)
	if err != nil {
		return
	}
	return
}

func (ouidb *OUIDB) getVersion() (v string, err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	q := fmt.Sprintf("SELECT name FROM sqlite_schema WHERE type='table' AND name LIKE '%s'", _tableVersion)
	conn, err := ouidb.getConnection()
	if err != nil {
		return
	}
	res, err := conn.Query(q)
	if err != nil {
		return
	}
	for res.Next() {
		res.Scan(&v)
	}
	if v == "" {
		err = fmt.Errorf(_updateMsg, _tableVersion)
		return
	}
	return
}

func (ouidb *OUIDB) createTable() (err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	q := fmt.Sprintf("CREATE TABLE `%s` ( `id` INTEGER PRIMARY KEY AUTOINCREMENT, `prefix` VARCHAR(32) NOT NULL, `length` INTEGER NOT NULL, `org` VARCHAR(64) NOT NULL, UNIQUE(prefix, length) ON CONFLICT REPLACE )", _tableVersion)
	_, err = ouidb.Connection.Exec(q)
	return err
}

func (ouidb *OUIDB) Delete() (err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	_, f, err := ouidb.getFileName()
	if err != nil {
		return
	}
	if pathExists(f) {
		err = os.Remove(f)
		if err != nil {
			return
		}
	}
	return
}

func (ouidb *OUIDB) Insert(d VendorDef) (res sql.Result, err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	s, err := ouidb.Connection.Prepare(fmt.Sprintf("INSERT INTO %s(prefix, length, org) values(?,?,?)", _tableVersion))
	if err != nil {
		return
	}
	res, err = s.Exec(d.Prefix, d.Length, d.Org)
	return
}

func Find(search string) (matches chan VendorDef) {
	matches = make(chan VendorDef)
	ouidb, err := NewOUIDB()
	MaybePanic(err)

	go func() {
		mac, err := macaddr.ParseMACAddress(search)
		MaybePanic(err)
		q := fmt.Sprintf("SELECT prefix,length,org FROM %s WHERE prefix LIKE '%s%%'", _tableVersion, mac.OUI())
		rows, err := ouidb.Connection.Query(q)
		MaybePanic(err)

		defer rows.Close()
		defer ouidb.Connection.Close()

		for rows.Next() {
			var prefix string
			var length int
			var org string
			err := rows.Scan(&prefix, &length, &org)
			MaybePanic(err)
			def := VendorDef{Prefix: prefix, Length: length, Org: org}
			_, mp, err := macaddr.ParseMACPrefix(def.PrefixString())
			MaybePanic(err)
			_, failure := mp.Match(search)
			if failure == nil {
				matches <- def
			}
		}
		close(matches)
	}()
	return
}

func (ouidb *OUIDB) Populate(p *progress.Progress) (records int, err error) {
	if ouidb == nil {
		err = fmt.Errorf("OUIDB is not initialized")
		return
	}
	p.Start()
	err = ouidb.Delete()
	p.AdvanceTo(10)
	if err != nil {
		return
	}
	err = ouidb.createTable()
	if err != nil {
		return
	}
	p.AdvanceTo(20)

	f, n, err := DownloadFile(ouidb.Directory, p)
	if err != nil {
		return
	}
	total := n / 48
	p.AdvanceTo(50)

	defer ouidb.Connection.Close()

	for def := range Collect(f, total, p) {
		_, err = ouidb.Insert(def)
		if err != nil {
			return
		}
		records++
	}
	return
}
