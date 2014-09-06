package main

import (
   "log"
   "database/sql"
)

// Create database
func CreateDB(db *sql.DB) bool {
   sqlStmt := `
      create table networks(id integer not null primary key, essid text not null unique, inscope boolean not null default 0, first_seen datetime not null default (datetime(current_timestamp)),
           last_seen datetime not null default (datetime(current_timestamp)), hidden boolean not null default 0,
           encryption text not null default('') check(encryption in ('', 'Open', 'WEP', 'WPA-Personal', 'WPA2-Personal', 'WPA-Enterprise', 'WPA2-Enterprise')),
           handshake boolean not null default 0, key text not null default(''), eap text not null default('') check(eap in ('', 'PEAP', 'TLS')));
      delete from networks;
      insert into networks(essid) values ('');
      create table bssids(id integer not null primary key, bssid text not null unique, network integer not null, first_seen datetime not null default (datetime(current_timestamp)),
           last_seen datetime not null default (datetime(current_timestamp)), channel integer not null, power integer not null, foreign key(network) references networks(id));
      delete from bssids;
      create table clients(id integer not null primary key, mac text not null unique, name text not null default(''), first_seen datetime not null default (datetime(current_timestamp)),
           last_seen datetime not null default (datetime(current_timestamp)), power integer not null, network integer not null, foreign key(network) references networks(id));
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
      return(false)
   }

   return(true)
}

/*
   DB Transaction

   tx, err := db.Begin()
   if err != nil {
      log.Fatal(err)
   }
   stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
   if err != nil {
      log.Fatal(err)
   }
   defer stmt.Close()
   for i := 0; i < 100; i++ {
      _, err = stmt.Exec(i, fmt.Sprintf("index%03d", i))
      if err != nil {
         log.Fatal(err)
      }
   }
   tx.Commit()

   rows, err := db.Query("select id, name from foo")
   if err != nil {
      log.Fatal(err)
   }
   defer rows.Close()
   for rows.Next() {
      var id int
      var name string
      rows.Scan(&id, &name)
      fmt.Println(id, name)
   }
   rows.Close()
*/
