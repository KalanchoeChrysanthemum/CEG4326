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
        if res is not None:
            return res
        else:
            return None

"""
Write test users to database.
"""
def initialize_db(db):
    key = encode_data('John')
    val = encode_data('Doe')
    with db.begin(write=True) as d:
        d.put(key, val)

"""
Compares the hash of the provided string (actual),
with the expected bytes
"""
def compare(actual: str, expected: bytes) -> bool:
    return encode_data(actual) == expected

def main():
    with lmdb.open(DATABASE_NAME, map_size=DATABASE_SIZE) as db:
        initialize_db(db)
        
        # Scan RFID & LED matrix here
        rfid = 'John' # Placeholder
        led = 'Doe' # Placeholder

        res = query_db(db, rfid)
        
        if res is None:
            print('Invalid RFID')
        else:
            if (compare(led, res)):
                print('Success')
            else:
                print('Invalid LED matrix')


if __name__ == '__main__':
    main()
