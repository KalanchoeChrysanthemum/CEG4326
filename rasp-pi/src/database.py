import os
import lmdb
from siphashc import siphash

MEGABYTE = 1_048_576
DATABASE_NAME = '../../tags.db' # Location of database relative to this source file
DATABASE_SIZE = 10 * MEGABYTE # 10MB database for storing a large amount of users
KEY_PATH = '.sip.key' # Location of siphash key

'''
Generates and returns a randomized 128-bit key.

Key is also written to a file.
'''
def generate_key(path=KEY_PATH) -> bytes:
    key = os.urandom(16)
    with open(path, 'wb') as f:
        f.write(key)
    os.chmod(path, 0o600)
    return key

'''
First checks if hash key exists.

If key exists, returns key

Else, creates key and then returns key
'''
def load_key(path=KEY_PATH) -> bytes:
    if not os.path.exists(path):
        generate_key()
    with open(path, 'rb') as f:
        return f.read()

'''
Hashes provided text using the desired string in big endian order
'''
def encode_data(text: str, key: bytes) -> bytes:
    return siphash(key, text.encode()).to_bytes(8, byteorder='big')

'''
Hashes provided key and checks if it exists in the database.

If key exists returns the associated value (LED matrix).

If key does not exist returns None
'''
def query_db(db, rfid: str, key: bytes) -> bytes | None:
    h_key = encode_data(rfid, key)
    with db.begin() as d:
        return d.get(h_key)

'''
Opens provided database and writes hard-coded users into database.

NOTE: This function is for testing purposes
'''
def write_users(db, key: bytes) -> None:
    users = {
        '0D0031EF53': '00011101010100000001'
    }
    for rfid, led in users.items():
        db.put(encode_data(rfid, key), encode_data(led, key))

'''
Checks if the hash of the provided string (LED matrix) matches the expected hash (LED matrix)
'''
def is_valid(actual: str, expected: bytes, key: bytes) -> bool:
    return encode_data(actual, key) == expected

'''
Checks if the database is already built and doesn't need to be created
'''
def is_db_init(db) -> bool:
    with db.begin() as d:
        return d.get(b'__init__') == b'true'

'''
Initializes the database, writing the default users and hashing with the provided key
'''
def init_db(db, key: bytes) -> bytes:
    print('[INFO] Initializing database...')
    with db.begin(write=True) as d:
        write_users(d, key)
        d.put(b'__init__', b'true')
    print('[INFO] Initialization complete.')

def main() -> None:
    key = load_key()

    # Check database exists
    with lmdb.open(DATABASE_NAME, map_size=DATABASE_SIZE) as db:
        if not is_db_init(db):
            init_db(db, key)

    # Fully open the database to the program, set to readonly mode to prevent unwanted writes
    with lmdb.open(DATABASE_NAME, map_size=DATABASE_SIZE, readonly=True) as db:
        # Create dummy user to test functionality
        rfid = '0D0031EF53'
        led = '00011101010100000001'

        # Retrieve user data
        res = query_db(db, rfid, key)
        print(f'[DEBUG] Hex Dump Of RFID Hash: {rfid} -> {encode_data(rfid, key).hex()}')
        print(f'[DEBUG] Hex Dump Of LED Matrix Hash: {led} -> {encode_data(led, key).hex()}')

        # Check if user has valid RFID && valid LED matrix
        if res:
            print(f'[DEBUG] Hex Dump Of Queried RFID: {res.hex()}')
            print('[PASS] Valid User') if is_valid(led, res, key) else print('[FAILED] Invalid LED Matrix')
        else:
            print('[FAILED] Invalid RFID')


if __name__ == '__main__':
    main()
