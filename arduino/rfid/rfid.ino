void setup() {
  // Start the serial communication
  Serial.begin(9600);
  
  // Wait for the serial communication to start
  while (!Serial);
}

void loop() {
  // Sample wid and hash
  String wid = "000000000000000000773030376D6171";
  String hash = "55f8b969f2a7c33cfb87edaa2d1afafd";
  
  // Send the wid and hash over serial
  Serial.print(wid);
  Serial.print(",");
  Serial.println(hash);
  
  // Wait for 3 seconds before sending again
  delay(3000);
}
