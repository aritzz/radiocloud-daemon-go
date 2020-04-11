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
  "github.com/aritzz/radiocloud-daemon/src/transcoder"
  "log"
  // "encoding/json"
  "mime/multipart"
	"net/http"
	"os"
  "time"
  "errors"
  "io"
	"path/filepath"
  "bytes"
  // "github.com/aritzz/radiocloud-daemon/src/transcoder"
)

const TEMPDIR = "tmp"

func PodcastUploader(signalstop *bool) {
  WriteLog("Uploader", "Starting uploader")
  os.MkdirAll(TEMPDIR, os.ModePerm)

  for !*signalstop {
    time.Sleep(1 * time.Minute)
    dat, err := db.UploadGetPending()
    if err != nil {
      return
    }

    for _, uploadElement := range dat {
      if *signalstop {
        break
      }
      WriteLog("Uploader", "Starting upload pipeline for "+uploadElement.Progname)
      err = executeUploadPipeline(uploadElement)
      if err != nil {
        WriteLog("Uploader", "Error in upload pipeline for "+uploadElement.Progname)
        db.UpdateUploadError(uploadElement.Id)
      } else {
        WriteLog("Uploader", "Upload finished for "+uploadElement.Progname)
      }
    }
  }

  // transcoder.ConvertToMP3()

  WriteLog("Uploader", "Stopping uploader")
}


func executeUploadPipeline(podcast datastruct.UploaderPending) (error) {
  var err error
  var link string

  // -1. Set status
  db.UpdateUploading(podcast.Id)

  // 0. Create temporary transcoded file
  inputFile := globalConfig.Directories.Radiocloud+globalConfig.Directories.PodcastUpload+"/"+podcast.Tempfile
  outputFile := TEMPDIR+"/"+GetRandString(5)+"-"+podcast.Username+"."+globalConfig.AudioFormat

  WriteLog("Transcoder", "Transcoding started for "+podcast.Progname)
  start := time.Now()
  err = transcoder.CreateTranscodedFile(inputFile, outputFile, globalConfig.AudioFormat, globalConfig.AudioQuality)
  if err != nil {
    return err
  }
  elapsed := time.Since(start)
  WriteLog("Transcoder", "Transcoding finished for "+podcast.Progname+" in "+elapsed.String())

  // 1. Add to repeat
  if (podcast.Repeat == 1) && (podcast.RepeatAdded == 0) {
    err = doAddToRepeat(outputFile, globalConfig.Directories.Repeat+"/"+podcast.Destfile)
    if err != nil {
      goto END
    }

    err = db.UpdateRepeat(podcast.Id)
    if err != nil {
      goto END
    }

    WriteLog("Transcoder", "Repeat file added for "+podcast.Progname)
  }


  // 2. Upload to external
  if (podcast.Podcast == 1) && (podcast.Uploaded == 0) {
    WriteLog("Uploader", "Uploading to external server for "+podcast.Progname)
    link, err = doExternalUpload(podcast, globalConfig.ExternalUrl.External, outputFile, globalConfig.AudioFormat)
    // Error checking
    if err != nil {
      goto END
    }

    // Uploaded to website
    err = db.UpdateUploaded(podcast.Id)
    if err != nil {
      goto END
    }
    WriteLog("Uploader", "Uploading finished for "+podcast.Progname)
    podcast.Uploaded = 1
  }



  // 3. Upload to Arrosa only if is uploaded already to external
  if podcast.Podcast == 1 && podcast.Uploaded == 1 && podcast.Arrosa == 1 {
    WriteLog("Uploader", "Uploading to Arrosa for "+podcast.Progname)
    err = doArrosaUpload(podcast, globalConfig.ExternalUrl.Arrosa, globalConfig.Directories.Radiocloud+"/"+podcast.Image, globalConfig.AudioFormat, link)
    if err != nil {
      goto END
    }
  }

  END:

  os.Remove(outputFile)

  return err
}

func doAddToRepeat(tmpfile, destfile string) (error) {
  err := copy(tmpfile, destfile)
  return err
}

// Thank you, https://matt.aimonetti.net/posts/2013-07-golang-multipart-file-upload-example/
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}


func doExternalUpload(podcast datastruct.UploaderPending, url, filepath, extension string) (string, error) {
  var retUrl string

	extraParams := map[string]string{
		"date":  podcast.Date.Format("2006-01-02"),
		"type": extension,
    "title": podcast.Title,
    "text": podcast.Text,
    "name": podcast.Progname,
    "user": podcast.Username,
	}
	request, err := newfileUploadRequest(url, extraParams, "file", filepath)
	if err != nil {
		return retUrl, err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return retUrl, err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
    if err != nil {
			return retUrl, err
		}
    resp.Body.Close()
    return body.String(), nil
	}

  return retUrl, nil
}

func doArrosaUpload(podcast datastruct.UploaderPending, url, filepath, extension, link string) (error) {

  if len(link) == 0 || !checkUrl(link) || len(podcast.ArrosaUser) == 0 {
    // return errors.New("Invalid link")
    return nil
  }

	extraParams := map[string]string{
		"arrosa_user": podcast.ArrosaUser,
    "arrosa_pass": podcast.ArrosaPass,
    "arrosa_category": podcast.ArrosaCategory,
    "xmlrpc_url": globalConfig.ExternalUrl.ArrosaXmlrpc,
    "title": podcast.Title,
    "progname": podcast.Progname,
    "text": podcast.Text,
    "name": podcast.Username,
    "filesize": getSize(filepath),
    "filetype": extension,
    "link": link,
	}
	request, err := newfileUploadRequest(url, extraParams, "file", filepath)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
    if err != nil {
			return err
		}
    resp.Body.Close()
		/*log.Println(resp.StatusCode)
		log.Println(resp.Header)
		log.Println(body)*/
    if resp.StatusCode != 200 {
      return errors.New("Unknown status code")
    }
	}

  return nil
}
