import lmdb

MEGABYTE = 1_048_576
DATABASE_NAME = 'test.db'
DATABASE_SIZE = 10 * MEGABYTE

def db_init():
    env = lmdb.open(DATABASE_NAME, map_size=DATABASE_SIZE)
    _ = env

def main():
    db_init()
    print('Creating test database...')

if __name__ == '__main__':
    main()
