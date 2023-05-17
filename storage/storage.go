package storage

type Storage struct {
	store map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		store: make(map[string]string),
	}
}

func (s *Storage) SaveLongURL(hashedURL, longURL string) {
	s.store[hashedURL] = longURL
}

func (s *Storage) GetLongURL(hashedURL string) (string, bool) {
	longURL := s.store[hashedURL]

	if longURL == "" {
		return "", false
	}

	return longURL, true
}
