package mock

// DB holds a mocked in-memory database
type DB struct {
	Memory map[string]string
}

// Put inserts data into the database
func (db *DB) Put(secretId string) error {
	db.Memory[secretId] = secretId
	return nil
}

// Get retrieves data from the database
func (db *DB) Get(secretId string) (string, error) {
	return db.Memory[secretId], nil
}

// Has checks the existance in the database
func (db *DB) Has(secretId string) (bool, error) {
	_, ok := db.Memory[secretId]
	return ok, nil
}

// Remove deletes an entry from the database
func (db *DB) Remove(secretId string) error {
	delete(db.Memory, secretId)
	return nil
}
