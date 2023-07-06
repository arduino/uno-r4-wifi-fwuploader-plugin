// To genereate the binaries run:
// arduino-cli compile ./commands.ino --fqbn arduino:renesas_uno:unor4wifi -e

#include <Modem.h>

char waitResponse() {
  while (true) {
    if (Serial.available()) {
      char choice = Serial.read();
      switch (choice) {
        case 'r':
          return 'r';
        case 'v':
          return 'v';
        default:
          continue;
      }
    }
  }
}

void reboot() {
  std::string res = "";
  modem.write(std::string(PROMPT("+RESET=1")), res, CMD("+RESET=1"));
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
  char command = waitResponse();

  switch (command) {
    case 'r':
      reboot();
      break;
    case 'v':
      version();
      break;
  }
}

void loop() {

}
