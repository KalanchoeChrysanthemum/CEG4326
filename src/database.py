import lmdb
import hashlib as hl

MEGABYTE = 1_048_576
DATABASE_NAME = 'tags.db'
DATABASE_SIZE = 10 * MEGABYTE

"""
Hash provided string input and convert to bytes.
"""
def encode_data(text: str) -> bytes:
    return hl.sha256(text.encode()).digest()

"""
Poll the database for the provided RFID.
RFID must be a string.
Returns the corresponding (hashed) value if present,
None if not.
"""
def query_db(db, rfid: str) -> str:
    key = encode_data(rfid)
    with db.begin() as d:
        res = d.get(key)
        return res if res is not None else None

"""
Write test users to database.
"""
def initialize_db(db) -> None:
    key = encode_data('John')
    val = encode_data('Doe')
    with db.begin(write=True) as d:
        d.put(key, val)

"""
Compares the hash of the provided string (actual),
with the expected bytes
"""
def is_valid(actual: str, expected: bytes) -> bool:
    return encode_data(actual) == expected

def main() -> None:
    with lmdb.open(DATABASE_NAME, map_size=DATABASE_SIZE) as db:
        initialize_db(db)
        
        # Scan RFID & LED matrix here
        rfid = 'John' # Placeholder
        led = 'Doe' # Placeholder

        res = query_db(db, rfid)
        
        if res is not None:
            print('[PASS] Valid User') if is_valid(led, res) else print('[FAILED] Invalid LED Matrix')
        else:
            print('[FAILED] Invalid RFID')
            


if __name__ == '__main__':
    main()
