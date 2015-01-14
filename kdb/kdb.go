package kdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Create database
func CreateDB(db *sql.DB, dbfile string) {
	// Create the schema, if necessary
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		fmt.Println("create")
		sqlStmt := `
      create table networks(essid text not null primary key, inscope boolean not null default 0, cloaked boolean not null default 0,
      	   first_seen datetime not null default (datetime(current_timestamp)), last_seen datetime not null default (datetime(current_timestamp)),
           maxrate integer not null default 0, encryption text not null default('') check(encryption in ('', 'Open', 'WEP', 'WPA-Personal', 'WPA2-Personal', 'WPA-Enterprise', 'WPA2-Enterprise')));
      delete from networks;
      create table bssids(bssid text not null primary key, aptype text not null default '', network text not null default '', manuf text not null default '', channel integer not null default 0,
      	   first_seen datetime not null default (datetime(current_timestamp)), last_seen datetime not null default (datetime(current_timestamp)),
	   atype integer not null default 0, rangeip text not null default '', netmaskip text not null default '', gatewayip text not null default '',
	   minpower integer not null default -99, maxpower integer not null default 0);
      delete from bssids;
      create table clients(mac text not null primary key, first_seen datetime not null default (datetime(current_timestamp)),
           last_seen datetime not null default (datetime(current_timestamp)), power integer not null default -99, minpower integer not null default -99,
	   maxpower integer not null default 0, packets integer not null default 0);
      delete from clients;
      create table probes(id integer not null primary key, client integer not null, essid text not null, foreign key(client) references clients(id));
      delete from probes;
      create table credentials(id integer not null primary key, network integer not null, client integer not null, challenge text not null, response text not null,
           foreign key(network) references networks(id), foreign key(client) references clients(id));
      delete from credentials;
   `
		_, err := db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
		}
	}
}

func InsertNetwork(db *sql.DB, essid string, cloaked int, first string, last string, rate int, enc string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	fInt, _ := strconv.Atoi(first)
	lInt, _ := strconv.Atoi(last)

	stmt, err := tx.Prepare("insert or replace into networks(essid, cloaked, first_seen, last_seen, maxrate, encryption) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(essid, cloaked, time.Unix(int64(fInt), 0), time.Unix(int64(lInt), 0), rate, enc)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func InsertAP(db *sql.DB, bssid string, network string, channel int, first string, last string, min int, max int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	fInt, _ := strconv.Atoi(first)
	lInt, _ := strconv.Atoi(last)

	stmt, err := tx.Prepare("insert or replace into bssids(bssid, network, channel, first_seen, last_seen, minpower, maxpower) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(bssid, network, channel, time.Unix(int64(fInt), 0), time.Unix(int64(lInt), 0), min, max)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func InsertClient(db *sql.DB, mac string, first string, last string, power int, min int, max int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	fInt, _ := strconv.Atoi(first)
	lInt, _ := strconv.Atoi(last)

	stmt, err := tx.Prepare("insert or replace into clients(mac, first_seen, last_seen, power, minpower, maxpower) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(mac, time.Unix(int64(fInt), 0), time.Unix(int64(lInt), 0), power, min, max)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
