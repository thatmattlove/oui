package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gookit/gcli/v3/progress"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thatmattlove/go-macaddr"
)

func DBFileName() (d, f string, err error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	d = fmt.Sprintf("%s/oui", dir)
	f = d + "/oui.db"
	return
}

func GetOrCreateDBFile() (f *os.File, err error) {
	dn, fn, err := DBFileName()
	if err != nil {
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

func GetDB() (db *sql.DB) {
	f, err := GetOrCreateDBFile()
	MaybePanic(err)
	db, err = sql.Open("sqlite3", f.Name())
	MaybePanic(err)
	return
}

func TableVersion(db *sql.DB) (string, error) {
	q := fmt.Sprintf("SELECT name FROM sqlite_schema WHERE type='table' AND name LIKE '%s'", _tableVersion)
	res, err := db.Query(q)
	MaybePanic(err)
	var e string
	for res.Next() {
		res.Scan(&e)
	}
	if e == "" {
		return e, fmt.Errorf(_updateMsg, _tableVersion)
	}
	return e, nil
}

func CreateDB(db *sql.DB) error {
	c := fmt.Sprintf("CREATE TABLE `%s` ( `id` INTEGER PRIMARY KEY AUTOINCREMENT, `prefix` VARCHAR(32) NOT NULL, `length` INTEGER NOT NULL, `org` VARCHAR(64) NOT NULL, UNIQUE(prefix, length) ON CONFLICT REPLACE )", _tableVersion)
	_, err := db.Exec(c)
	return err
}

func DeleteDB() (err error) {
	_, f, err := DBFileName()
	if err != nil {
		return
	}
	if pathExists(f) {
		err = os.Remove(f)
		if err != nil {
			return
		}
	}
	return nil
}

func UpdateTable(p *progress.Progress) (int, error) {
	p.Start()
	err := DeleteDB()
	p.AdvanceTo(10)
	if err != nil {
		return 0, err
	}
	db := GetDB()
	err = CreateDB(db)
	if err != nil {
		return 0, err
	}
	p.AdvanceTo(20)
	return Populate(p)
}

func CheckTable(db *sql.DB) (string, error) {
	ver, err := TableVersion(db)
	if err != nil {
		return ver, err
	}
	if ver != _tableVersion {
		return ver, fmt.Errorf(_updateMsg, _tableVersion)
	}
	return ver, nil
}

func Insert(db *sql.DB, d VendorDef) (res sql.Result, err error) {
	s, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(prefix, length, org) values(?,?,?)", _tableVersion))
	MaybePanic(err)
	res, err = s.Exec(d.Prefix, d.Length, d.Org)
	return
}

func SelectAll() (chnl chan VendorDef) {
	chnl = make(chan VendorDef)
	go func() {
		db := GetDB()
		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", _tableVersion))
		MaybePanic(err)

		defer rows.Close()

		for rows.Next() {
			var i int
			var p string
			var l int
			var o string
			err := rows.Scan(&i, &p, &l, &o)
			MaybePanic(err)
			chnl <- VendorDef{Prefix: p, Length: l, Org: o}
		}
		close(chnl)
		db.Close()
	}()
	return
}

func Find(search string) (matches chan VendorDef) {
	matches = make(chan VendorDef)

	go func() {
		db := GetDB()
		mac, err := macaddr.ParseMACAddress(search)
		MaybePanic(err)
		q := fmt.Sprintf("SELECT prefix,length,org FROM %s WHERE prefix LIKE '%s%%'", _tableVersion, mac.OUI())
		rows, err := db.Query(q)
		MaybePanic(err)

		defer rows.Close()

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
		db.Close()
	}()
	return
}
