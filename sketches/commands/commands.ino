// This file is part of uno-r4-wifi-fwuploader-plugin.
//
// Copyright (c) 2023 Arduino LLC.  All right reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// To genereate the binaries run:
// arduino-cli compile -e --profile unor4wifi

#include <Modem.h>

void reboot() {
  std::string res = "";
  modem.write(std::string(PROMPT("+RESTARTBOOTLOADER")), res, CMD("+RESTARTBOOTLOADER"));
}

void version() {
  static char fw_version[12];
  std::string res = "";
  if(modem.write(std::string(PROMPT(_FWVERSION)), res, CMD_READ(_FWVERSION))) {
      memset(fw_version,0x00,12);
      memcpy(fw_version, res.c_str(), res.size() < 12 ? res.size() : 11);
      Serial.println(fw_version);
   } else {
       Serial.println("0.0.0");
   }
}

void setup() {
  Serial.begin(9600);
  modem.begin();
}

void loop() {
  while (true) {
    if (Serial.available()) {
      char choice = Serial.read();
      switch (choice) {
        case 'r':
            reboot();
            break;
        case 'v':
            version();
            break;
        default:
          continue;
      }
    }
  }
}
