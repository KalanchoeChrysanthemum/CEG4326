#include <Crypto.h>
#include <SHA256.h>
#include <MFRC522.h>
#include <MFRC522.h>

#define RST_PIN         9         
#define SS_PIN          10  

MFRC522 mfrc522(SS_PIN, RST_PIN);
MFRC522::MIFARE_Key keyA;
MFRC522::MIFARE_Key keyB;

// Key for the HMAC hash

void generate_pwm_signal() {
    // Pin 3 has 228KHz signal
    pinMode(3, OUTPUT);

    // Modify timer2 to generate desired signal
    TCCR2A = 0;
    TCCR2B = 0;
    TCNT2 = 0;

    TCCR2A = (1 << WGM20);
    TCCR2B = (1 << WGM22);

    TCCR2A |= (1 << COM2B1);

    TCCR2B |= (1 << CS20);

    OCR2A = 34;
    OCR2B = 14;

    // Fix pin 11's output signal
    TCCR2A &= ~( (1 << COM2A1) | (1 << COM2A0) );
}

String scan_Rfid(byte blockAddress) {

  byte len = 18;
  byte dataBlock[18];
  memset(dataBlock, 0x00, len);
  
  byte status = mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_B, 7, &keyB, &(mfrc522.uid));

  if(status != MFRC522::STATUS_OK){
      Serial.println("Authentication failed");
      return;
  }

MFRC522::StatusCode st;

st = mfrc522.MIFARE_Read(blockAddress, dataBlock, &len);
if(st != MFRC522::STATUS_OK){
  Serial.println("Reading Error");
  Serial.println(mfrc522.GetStatusCodeName(st));
  return;
}
String readData;
for(int i=0; i<16; i++) {
  if (dataBlock[i] < 0x10) {
      readData += "0";
  } 
    readData += String(dataBlock[i], HEX);
}
return readData;
}

void setup() {
    generate_pwm_signal(); // Pin 3 outputs 228KHz signal
    Serial.begin(9600);


    SPI.begin();        // Init SPI bus
    mfrc522.PCD_Init(); // Init MFRC522 card
mfrc522.PCD_Init();

for(byte i=0; i<6; i++)
  keyA.keyByte[i] = 0xff;

for(byte i=0; i<6; i++)
  keyB.keyByte[i] = 0xff;
   
}

void loop() {
   if(!mfrc522.PICC_IsNewCardPresent())
    return;
  if(!mfrc522.PICC_ReadCardSerial())
    return;

   byte block4 = 4;
   byte block5 = 5;
   String wid, rfid;

    rfid = scan_Rfid(block4);
    wid = scan_Rfid(block5);

      Serial.print(wid);
      Serial.print(",");
      Serial.println(rfid);
   


    mfrc522.PICC_HaltA();
  mfrc522.PCD_StopCrypto1();
}

