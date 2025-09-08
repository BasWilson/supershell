package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type Record struct {
	Nickname   string `json:"nickname"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	AuthMethod string `json:"auth_method"`
	KeyPath    string `json:"key_path,omitempty"`
	Password   string `json:"password,omitempty"`
}

type Store struct {
	path string
	mu   sync.RWMutex
	data map[string]Record
}

func New() (*Store, error) {
	configDir, err := os.UserConfigDir()
	if err != nil { return nil, err }
	appDir := filepath.Join(configDir, "supershell")
	if err := os.MkdirAll(appDir, 0o700); err != nil { return nil, err }
	p := filepath.Join(appDir, "connections.json")
	s := &Store{ path: p, data: map[string]Record{} }
	if err := s.load(); err != nil {
		if errors.Is(err, fs.ErrNotExist) { return s, nil }
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil { return err }
	defer f.Close()
	dec := json.NewDecoder(f)
	return dec.Decode(&s.data)
}

func (s *Store) save() error {
	tmp := s.path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil { return err }
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s.data); err != nil { f.Close(); return err }
	if err := f.Close(); err != nil { return err }
	return os.Rename(tmp, s.path)
}

func (s *Store) Add(r Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[r.Nickname]; exists {
		return fmt.Errorf("nickname already exists: %s", r.Nickname)
	}
	s.data[r.Nickname] = r
	return s.save()
}

func (s *Store) Update(nickname string, mutate func(*Record)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.data[nickname]
	if !ok { return fmt.Errorf("not found: %s", nickname) }
	mutate(&r)
	s.data[nickname] = r
	return s.save()
}

func (s *Store) Delete(nickname string) error {
	s.mu.Lock(); defer s.mu.Unlock()
	if _, ok := s.data[nickname]; !ok { return fmt.Errorf("not found: %s", nickname) }
	delete(s.data, nickname)
	return s.save()
}

func (s *Store) Get(nickname string) (Record, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	r, ok := s.data[nickname]
	if !ok { return Record{}, fmt.Errorf("not found: %s", nickname) }
	return r, nil
}

func (s *Store) List() []Record {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]Record, 0, len(s.data))
	for _, r := range s.data { out = append(out, r) }
	sort.Slice(out, func(i, j int) bool { return out[i].Nickname < out[j].Nickname })
	return out
}
