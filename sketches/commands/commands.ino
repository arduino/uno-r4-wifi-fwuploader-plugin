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
  modem.begin();
  std::string res = "";
  modem.write(std::string(PROMPT("+RESET=1")), res, CMD("+RESET=1"));
}

void version() {
}

void setup() {
  Serial.begin(9600);
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
