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

package transcoder

import (
    "github.com/xfrr/goffmpeg/transcoder"
    // "fmt"
)


// Creates transcoded file from input, according to
// a format and bitrate specified
func CreateTranscodedFile(input, output, format, bitrate string) (error) {
  trans := new(transcoder.Transcoder)
  err := trans.Initialize(input, output)
  if err != nil {
    return err
  }

  // Set transcoding info and start process
  trans.MediaFile().SetAudioCodec(getCodec(format))
  trans.MediaFile().SetAudioBitRate(bitrate+"k")

	done := trans.Run(false)
	err = <-done

  return err
}


func getCodec(format string) (string) {
  switch (format) {
  case "mp3":
    return "mp3"
  case "ogg":
    return "libvorbis"
  }

  return "mp3"
}
