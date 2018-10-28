void setup() {
  Serial.begin(19200);
}

void loop() {
  if (Serial.available()) {
    const uint8_t incoming_byte = Serial.read();
    Serial.write(incoming_byte + 1);
  }
}
