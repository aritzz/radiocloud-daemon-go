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

package main

import(
	"github.com/aritzz/radiocloud-daemon/src/db"
	"github.com/aritzz/radiocloud-daemon/src/datastruct"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var globalConfig datastruct.Configuration

func main() {
	var err error
	mainStatusFinish := false
	signalrecv := make(chan os.Signal, 1)
  signal.Notify(signalrecv, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,	syscall.SIGQUIT)

	// Signal handler
	go func() {
		for {
			<- signalrecv
			mainStatusFinish = true
		}
	}()

	// Say hello
	WriteHello()
	WriteLog("Main", "Starting daemon routines")

	// Save databse info
	db.SaveInfo(getEnv("DBUSER", "root"), getEnv("DBPASS", "radixu"), getEnv("DBHOST", "localhost"), getEnv("DBPORT", "3306"), getEnv("DBNAME", "radiocloud12"))

	// Connect to database
	if err := db.DBConnectTest(); err != nil {
		WriteLog("Main", "Error connecting to database")
		return
	}

	// Connected to database - proceed
	WriteLog("Main", "Connection established to database")

	// Get global configuration
	globalConfig, err = db.GetGlobalConfiguration()
	if err != nil  {
		WriteLog("Main", "Error getting configuration from database: " +err.Error())
		return
	}

	WriteLog("Main", "Global configuration loaded")

	// Create downloader process
	go PodcastDownloader(&mainStatusFinish)

	// Create uploader process
	go PodcastUploader(&mainStatusFinish)


	for !mainStatusFinish {
		time.Sleep(2 * time.Second)
	}

	WriteLog("Main", "Terminating main program")
	time.Sleep(2 * time.Second)

}
