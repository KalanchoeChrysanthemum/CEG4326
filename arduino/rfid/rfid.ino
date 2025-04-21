#include <SPI.h>
#include <MFRC522.h>

#define SS_PIN 10
#define RST_PIN 7

MFRC522 rfid(SS_PIN, RST_PIN);

MFRC522::MIFARE_Key key;

void setup() {
Serial.begin(9600);
SPI.begin();
rfid.PCD_Init();
}

void loop() {
  // put your main code here, to run repeatedly:
  if ( ! rfid.PICC_IsNewCardPresent())
  return;
}
