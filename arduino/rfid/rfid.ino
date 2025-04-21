#include <Arduino.h>

#define OUTPUT_PIN 9
#define FREQ 228000
#define CHANNEL 0
#define RESOLUTION 1
#define DUTY_CYCLE 1

void setup() {
  ledcSetup(0, FREQ, RESOLUTION);
  ledcAttachPin(OUTPUT_PIN, 0);
  ledcWrite(0, DUTY_CYCLE);
}

void loop() {

}
