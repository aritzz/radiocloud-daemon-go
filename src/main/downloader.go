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
  "time"
  "os"
  "io"
  "net/http"
  // "log"
  "github.com/aritzz/radiocloud-daemon/src/db"
  "github.com/nandosousafr/podfeed"
  "path/filepath"
  "errors"
)

func PodcastDownloader(signalstop *bool) {
  WriteLog("Downloader", "Starting downloader")

  /* Downloader loop */
  for !*signalstop {
    time.Sleep(1 * time.Minute)

    dat, err := db.DownloadGetPending()
    if err != nil {
      continue
    }

    for _, downElement := range dat {
      currentTime := time.Now()
      if ((downElement.DownloadDay == int(currentTime.Weekday())) && (downElement.DownloadHour == currentTime.Hour())) || downElement.Force  {
        urls, filename, err := getLatestPodcast(downElement.Url, downElement.DownloadAll, downElement.LastFile)
        if err != nil {
          WriteLog("Downloader", "Error parsing RSS Feed "+downElement.Name)
          continue
        }

        if downElement.LastFile == filename {
          continue
        }

        WriteLog("Downloader", "Starting download for "+downElement.Name)
        for _, url := range urls {
          if err = downloadPodcast(url, globalConfig.Directories.Radiocore+"/"+downElement.File, downElement.DownloadAll); err != nil {
            WriteLog("Downloader", "Error downloading Podcast for "+downElement.Name)
            continue
          }
        }


        WriteLog("Downloader", "Downloaded podcast for "+downElement.Name)

        // Update last file downloaded
        db.UpdateLastFileDownload(downElement.Id, filename)

      }
    }
  }
  // Downloader loop  end

  WriteLog("Downloader", "Stopping downloader")
}

func downloadPodcast(url, filename string, all int) (error) {
  thisfile := filename
  if all == 1 {
    thisfile = filename+"/"+filepath.Base(url)
    os.MkdirAll(filename, os.ModePerm)
  }

  out, err := os.Create(thisfile)
  if err != nil {
    return err
  }

  resp, err := http.Get(url)
  if err != nil {
    return err
  }

  defer out.Close()
  defer resp.Body.Close()

  _, err = io.Copy(out, resp.Body)
  if err != nil {
    return err
  }

  return nil
}


func getLatestPodcast(url string, all int, latest string) ([]string, string, error) {
  podcast, err := podfeed.Fetch(url)
  ret_string := []string{}
  if err != nil {
    return ret_string, "", err
  }

  if len(podcast.Items) == 0 {
    return ret_string, "", errors.New("Invalid item detected")
  }

  ret_string = append(ret_string, podcast.Items[0].Enclosure.Url)

  if all == 1 {
    for _, el := range podcast.Items {
      if filepath.Base(el.Enclosure.Url) == latest {
        break
      }

      ret_string = append(ret_string, el.Enclosure.Url)
    }
  }

  return ret_string, filepath.Base(podcast.Items[0].Enclosure.Url), nil
}
