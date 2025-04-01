from Crypto.Cipher import AES
import os
import binascii

SECRET_KEY = b'16byteslongkey!!'

def encrypt(nonce):
    cipher = AES.new(SECRET_KEY, AES.MODE_ECB)
    return cipher.encrypt(nonce.ljust(16, b'\x00'))

def verify(nonce, encrypted_response):
    expected = encrypt(nonce)
    return expected == encrypted_response

nonce = os.urandom(8)

encrypted_nonce = encrypt(nonce)

print(f"Nonce: {binascii.hexlify(nonce).decode()}")
print(f"Encrypted Nonce: {binascii.hexlify(encrypted_nonce).decode()}")

if verify(nonce, encrypted_nonce):
    print("\nTest passed")
else:
    print("\nTest failed")
