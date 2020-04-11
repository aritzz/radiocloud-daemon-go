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

func UploadGetPending() ([]datastruct.UploaderPending, error) {
  ret_struct := []datastruct.UploaderPending{}
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return ret_struct, err
  }

  data, err := db.Query("SELECT podcast_upload.*, users.username as username, users.image as image, programs.arrosa_user as arrosa_user, programs.arrosa_pass as arrosa_pass, programs.arrosa_category as arrosa_category, programs.name as progname, blocks.vars as destfile FROM ((podcast_upload INNER JOIN users ON users.id=podcast_upload.userid) INNER JOIN programs ON programs.id=users.programid) INNER JOIN blocks ON programs.blockid=blocks.id WHERE podcast_upload.is_trash=0 AND users.enabled=1 AND podcast_upload.uploaded=0 AND DATE(podcast_upload.date) <= CURDATE() ORDER BY podcast_upload.id ASC limit 5")

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
    tmpUploadElement := datastruct.UploaderPending{}
		for i, col := range values {

      switch columns[i] {
      case "id":
        tmpUploadElement.Id = string(col)
        break
      case "tmpname":
        tmpUploadElement.Tempfile = string(col)
        break
      case "progname":
        tmpUploadElement.Progname = string(col)
        break
      case "repeat_added":
        tmpUploadElement.RepeatAdded, _ = strconv.Atoi(string(col))
        break
      case "add_repeat":
        tmpUploadElement.Repeat, _ = strconv.Atoi(string(col))
        break
      case "destfile":
        tmpUploadElement.Destfile = string(col)
        break
      case "uploaded":
        tmpUploadElement.Uploaded, _ = strconv.Atoi(string(col))
        break
      case "add_podcast":
        tmpUploadElement.Podcast, _ = strconv.Atoi(string(col))
        break
      case "add_arrosa":
        tmpUploadElement.Arrosa, _ = strconv.Atoi(string(col))
        break
      case "username":
        tmpUploadElement.Username = string(col)
        break
      case "image":
        tmpUploadElement.Image = string(col)
        break
      case "arrosa_user":
        tmpUploadElement.ArrosaUser = string(col)
        break
      case "arrosa_pass":
        tmpUploadElement.ArrosaPass = string(col)
        break
      case "arrosa_category":
        tmpUploadElement.ArrosaCategory = string(col)
        break
      case "date":
        date, _ := time.Parse("2006-01-02", string(col))
        tmpUploadElement.Date = date
        break
      case "title":
        tmpUploadElement.Title = string(col)
        break
      case "text":
        tmpUploadElement.Text = string(col)
        break
      default:
        break
      }
		}

    ret_struct = append(ret_struct, tmpUploadElement)
  }


  return ret_struct, nil
}

// Updates repeat status to 1
func UpdateRepeat(podcastid string) (error) {
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return err
  }

  _, err = db.Query("UPDATE podcast_upload SET repeat_added=1 WHERE id="+podcastid)

  // Error control
  if err != nil {
    return err
  }

  return nil
}

// Updates uploaded status to 1
func UpdateUploaded(podcastid string) (error) {
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return err
  }

  _, err = db.Query("UPDATE podcast_upload SET uploaded=1, uploaded_date=NOW() WHERE id="+podcastid)

  // Error control
  if err != nil {
    return err
  }

  return nil
}

// Updates uploaded status to 2
func UpdateUploading(podcastid string) (error) {
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return err
  }

  _, err = db.Query("UPDATE podcast_upload SET uploaded=2 WHERE id="+podcastid)

  // Error control
  if err != nil {
    return err
  }

  return nil
}

func UpdateUploadError(podcastid string) (error) {
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return err
  }

  _, err = db.Query("UPDATE podcast_upload SET uploaded=0 WHERE id="+podcastid)

  // Error control
  if err != nil {
    return err
  }

  return nil
}
