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

package datastruct


type Configuration struct {
  AudioFormat string
  AudioQuality string
  UploadHour string
  Directories DirConfiguration
  ExternalUrl ExternalUrlConfiguration
}

type DirConfiguration struct {
  PodcastUpload string
  PodcastDownload string
  Radiocore string
  Radiocloud string
  Repeat string
}

type ExternalUrlConfiguration struct {
  ArrosaXmlrpc string
  Arrosa string
  External string
}
