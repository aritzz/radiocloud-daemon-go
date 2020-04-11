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

import(
  // "log"
  "database/sql"
  "github.com/aritzz/radiocloud-daemon/src/datastruct"
  "strconv"
  "time"
	// "fmt"
	// "io"
  // "encoding/json"
  _ "github.com/go-sql-driver/mysql"
)

func DownloadGetPending() ([]datastruct.DownloaderPending, error) {
  ret_struct := []datastruct.DownloaderPending{}
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return ret_struct, err
  }

  data, err := db.Query("SELECT podcast_download.*, blocks.vars as file, blocks.desc as izen FROM podcast_download INNER JOIN blocks ON podcast_download.blockid=blocks.id")

  // Error control
  if err != nil {
    return ret_struct, err
  }

  // Get columns
  columns, err := data.Columns()
	if err != nil {
		return ret_struct, err
	}

  // Make arrays for storage
  values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

  // Get data
  for data.Next() {
    err = data.Scan(scanArgs...)
    if err != nil {
      return ret_struct, err
    }
    tmpUploadElement := datastruct.DownloaderPending{}
		for i, col := range values {

      switch columns[i] {
      case "id":
        tmpUploadElement.Id = string(col)
        break
      case "url":
        tmpUploadElement.Url = string(col)
        break
      case "dday":
        tmpUploadElement.DownloadDay,_ = strconv.Atoi(string(col))
        break
      case "dhour":
        tmpUploadElement.DownloadHour,_ = strconv.Atoi(string(col))
        break
      case "last_update":
        tmpUploadElement.Force = false
        if string(col) == "FORCE" {
          tmpUploadElement.LastUpdate = time.Now()
          tmpUploadElement.Force = true
        } else {
          tmpUploadElement.LastUpdate, _ = time.Parse("2006-01-02 15:04:05", string(col))
        }
        break
      case "file":
        tmpUploadElement.File = string(col)
        break
      case "izen":
        tmpUploadElement.Name = string(col)
        break
      case "download_all":
        tmpUploadElement.DownloadAll,_ = strconv.Atoi(string(col))
        break
      case "last_file":
        tmpUploadElement.LastFile = string(col)
        break
      default:
        break
      }
		}

    ret_struct = append(ret_struct, tmpUploadElement)
  }


  return ret_struct, nil
}

// Updates last file downloaded
func UpdateLastFileDownload(id, name string) (error) {
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return err
  }

  _, err = db.Query("UPDATE podcast_download SET last_file='"+name+"', last_update=NOW() WHERE id="+id)

  // Error control
  if err != nil {
    return err
  }

  return nil
}

// Updates uploaded status to 1
/*func UpdateUploaded(podcastid string) (error) {
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return err
  }

  _, err = db.Query("UPDATE podcast_upload SET uploaded=1 WHERE id="+podcastid)

  // Error control
  if err != nil {
    return err
  }

  return nil
}*/
