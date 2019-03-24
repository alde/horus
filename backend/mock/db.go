package mock

import "errors"

// DB holds a mocked in-memory database
type DB struct {
	Memory map[string]map[string]string
}

// Get mock
func (d DB) Get(repository, key string) (string, error) {
	if r, ok := d.Memory[repository]; ok {
		if s, ok := r[key]; ok {
			return s, nil
		}
	}

	return "", errors.New("no such secret")
}

// Put mock
func (d DB) Put(repository, key string, secret []byte) error {
	if d.Memory[repository] == nil {
		d.Memory[repository] = make(map[string]string)
	}
	d.Memory[repository][key] = string(secret)
	return nil
}

// Remove mock
func (d DB) Remove(repository, key string) error {
	var r map[string]string
	var ok bool
	if r, ok = d.Memory[repository]; !ok {
		return nil
	}
	delete(r, key)

	return nil
}

// Has mock
func (d DB) Has(repository string, key string) bool {
	if r, ok := d.Memory[repository]; ok {
		if _, ok := r[key]; ok {
			return true
		}
	}
	return false
}

// List mock
func (d DB) List(repository string) []string {
	secrets := []string{}

	if r, ok := d.Memory[repository]; ok {
		for k := range r {
			secrets = append(secrets, k)
		}
	}

	return secrets
}
