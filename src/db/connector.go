/*
    RadioCloud Daemon 1.2 - Part of RadioCloud automation system

    Copyright (C) 2020 - Aritz Olea Zubikarai <aritzolea@gmail.com>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package db

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)


// Variables
var user string
var password string
var host string
var port string
var dbname string

// Makes a connection to database
func DBConnect() (error, *sql.DB) {
  dbstat, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+dbname+"?charset=utf8")

  if err != nil {
    return err, nil
  }

  return nil, dbstat
}

// Stores info in global vars
func SaveInfo(luser, lpassword, lhost, lport, ldbname string) {
  user = luser
  password = lpassword
  host = lhost
  port = lport
  dbname = ldbname
}

// Connection tester: makes a connection and it closes
func DBConnectTest() (error) {
  if err, dbstat := DBConnect(); err != nil {
    return err
  } else {
    dbstat.Close()
  }
  return nil
}
