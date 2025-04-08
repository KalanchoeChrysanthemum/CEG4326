import os
import lmdb
from siphashc import siphash


MEGABYTE = 1_048_576
DATABASE_NAME = '../../tags.db'
DATABASE_SIZE = 10 * MEGABYTE
KEY_PATH = '.sip.key'

def generate_key(path=KEY_PATH) -> bytes:
    key = os.urandom(16)
    with open(path, 'wb') as f:
        f.write(key)
    os.chmod(path, 0o600)
    return key

def load_key(path=KEY_PATH) -> bytes:
    if not os.path.exists(path):
        generate_key()
    with open(path, 'rb') as f:
        return f.read()

def encode_data(text: str, key: bytes) -> bytes:
    return siphash(key, text.encode()).to_bytes(8, byteorder='big')

def query_db(db, rfid: str, key: bytes) -> bytes | None:
    h_key = encode_data(rfid, key)
    with db.begin() as d:
        return d.get(h_key)

def write_users(db, key: bytes) -> None:
    users = {
        'John': 'Doe'
    }
    for rfid, led in users.items():
        db.put(encode_data(rfid, key), encode_data(led, key))

def is_valid(actual: str, expected: bytes, key: bytes) -> bool:
    return encode_data(actual, key) == expected

def is_db_init(db) -> bool:
    with db.begin() as d:
        return d.get(b'__init__') == b'true'

def init_db(db, key: bytes) -> bytes:
    print('[INFO] Initializing database...')
    with db.begin(write=True) as d:
        write_users(d, key)
        d.put(b'__init__', b'true')
    print('[INFO] Initialization complete.')

def main() -> None:
    key = load_key()

    with lmdb.open(DATABASE_NAME, map_size=DATABASE_SIZE) as db:
        if not is_db_init(db):
            init_db(db, key)

    with lmdb.open(DATABASE_NAME, map_size=DATABASE_SIZE) as db:
        rfid = 'John'
        led = 'Doe'

        res = query_db(db, rfid, key)
        print(f'[DEBUG] Hex Dump Of RFID Hash: {rfid} -> {encode_data(rfid, key).hex()}')
        print(f'[DEBUG] Hex Dump Of LED Matrix Hash: {led} -> {encode_data(led, key).hex()}')

        if res:
            print(f'[DEBUG] Hex Dump Of Queried RFID: {res.hex()}')
            print('[PASS] Valid User') if is_valid(led, res, key) else print('[FAILED] Invalid LED Matrix')
        else:
            print('[FAILED] Invalid RFID')


if __name__ == '__main__':
    main()
