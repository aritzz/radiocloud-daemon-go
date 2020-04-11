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
  "log"
  "math/rand"
  "net/url"
  "strings"
  "strconv"
  "time"
  "fmt"
  "os"
  "io"
)


func WriteLog(logtype, text string) {
  log.Println(logtype+" >> "+text)
}

func WriteHello() {
  log.Println("······························")
  log.Println("   RadioCloud daemon v1.2")
  log.Println("······························")
  log.Println()
}

func GetRandString(length int) (string) {
  var b strings.Builder
  rand.Seed(time.Now().UnixNano())
  chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

  for i := 0; i < length; i++ {
    b.WriteRune(chars[rand.Intn(len(chars))])
  }
  return b.String()
}

func copy(src, dst string) (error) {
  sourceFileStat, err := os.Stat(src)
  if err != nil {
          return err
  }

  if !sourceFileStat.Mode().IsRegular() {
          return fmt.Errorf("%s is not a regular file", src)
  }

  source, err := os.Open(src)
  if err != nil {
          return err
  }
  defer source.Close()

  destination, err := os.Create(dst)
  if err != nil {
          return err
  }
  defer destination.Close()
  _, err = io.Copy(destination, source)
  return err
}

func checkUrl(str string) (bool) {
    u, err := url.Parse(str)
    return err == nil && u.Scheme != "" && u.Host != ""
}

func getSize(file string) (string) {
  filesize, err :=  os.Stat(file);
  if err != nil {
    return "0"
  }

  return strconv.Itoa(int(filesize.Size()))
}

func getEnv(env, fallback string) (string) {
  data := os.Getenv(env)

  if len(data) == 0 {
    return fallback
  }

  return data
}
