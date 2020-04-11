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
  // "strconv"
  // "encoding/json"
  _ "github.com/go-sql-driver/mysql"
)

func GetGlobalConfiguration() (datastruct.Configuration, error) {
  ret_struct := datastruct.Configuration{}
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return ret_struct, err
  }

  data, err := db.Query("SELECT * FROM config")

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

    getThis := false
    getIndex := ""
		for _, col := range values {

      // Get index name
      if getThis == false {
        getIndex = string(col)
        getThis = !getThis
        continue
      }

      // log.Println("Value for "+getIndex+" is "+string(col))

      switch getIndex {
      case "audioformat":
        ret_struct.AudioFormat = string(col)
        break
      case "audioquality":
        ret_struct.AudioQuality = string(col)
        break
      case "upload_hour":
        ret_struct.UploadHour = string(col)
        break
      default:
        break
      }
      getThis = !getThis
		}

  }

  ret_struct.Directories, err = getDirectories()
  if err != nil {
    return ret_struct, err
  }

  ret_struct.ExternalUrl, err = getExternalUrl()
  if err != nil {
    return ret_struct, err
  }


  return ret_struct, nil
}


func getDirectories() (datastruct.DirConfiguration, error) {
  ret_struct := datastruct.DirConfiguration{}
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return ret_struct, err
  }

  data, err := db.Query("SELECT * FROM dirs")

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

    getThis := false
    getIndex := ""


		for i, col := range values {
      if columns[i] == "dirname" {
        getThis = false
      }

      // Get index name
      if getThis == false {
        getIndex = string(col)
        getThis = !getThis
        continue
      }

      switch getIndex {
      case "podcast_upload":
        ret_struct.PodcastUpload = string(col)
        break
      case "podcast_download":
        ret_struct.PodcastDownload = string(col)
        break
      case "radiocore_dir":
        ret_struct.Radiocore = string(col)
        break
      case "radiocloud_dir":
        ret_struct.Radiocloud = string(col)
        break
      case "repeat_dir":
        ret_struct.Repeat = string(col)
        break
      }
    }
  }

  return ret_struct, nil
}

func getExternalUrl() (datastruct.ExternalUrlConfiguration, error) {
  ret_struct := datastruct.ExternalUrlConfiguration{}
  err, db := DBConnect()
  defer db.Close()
  if err != nil {
    return ret_struct, err
  }

  data, err := db.Query("SELECT * FROM dirs")

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

    getThis := false
    getIndex := ""
    for i, col := range values {
      if columns[i] == "dirname" {
        getThis = false
      }

      // Get index name
      if getThis == false {
        getIndex = string(col)
        getThis = !getThis
        continue
      }

      switch getIndex {
      case "arrosa_xmlrpc":
        ret_struct.ArrosaXmlrpc = string(col)
        break
      case "external_upload":
        ret_struct.External = string(col)
        break
      case "arrosa_upload":
        ret_struct.Arrosa = string(col)
        break
      }
    }
  }

  return ret_struct, nil
}
