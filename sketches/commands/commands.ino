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
