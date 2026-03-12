package store

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

const (
	configDirName  = "ziti-mcp-server"
	configFileName = "config.json"
	currentVersion = 2
)

var profileNameRE = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// storeData is the on-disk format (v2).
type storeData struct {
	Version       int                          `json:"version"`
	ActiveProfile string                       `json:"active_profile"`
	Profiles      map[string]map[string]string `json:"profiles"`
}

// Store provides typed access to a JSON credential file at
// ~/.config/ziti-mcp-server/config.json with 0600 permissions.
type Store struct {
	mu   sync.RWMutex
	path string
	data storeData
}

// ClearResult describes the outcome of deleting a single key.
type ClearResult struct {
	Key     string
	Success bool
	Error   error
}

// DefaultPath returns the default config file path.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", configDirName, configFileName)
}

// New creates a Store backed by the given file path.
// If path is empty, DefaultPath() is used.
func New(path string) *Store {
	if path == "" {
		path = DefaultPath()
	}
	return &Store{
		path: path,
		data: storeData{
			Version:  currentVersion,
			Profiles: make(map[string]map[string]string),
		},
	}
}

// Load reads the config file from disk. If the file doesn't exist, the store starts empty.
// v1 flat-map files are automatically migrated to v2 format.
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		s.data = storeData{
			Version:  currentVersion,
			Profiles: make(map[string]map[string]string),
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	// Probe for version key to distinguish v1 from v2.
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	if _, hasVersion := probe["version"]; hasVersion {
		// v2 format
		var d storeData
		if err := json.Unmarshal(raw, &d); err != nil {
			return fmt.Errorf("parsing v2 config: %w", err)
		}
		if d.Profiles == nil {
			d.Profiles = make(map[string]map[string]string)
		}
		s.data = d
	} else {
		// v1 flat map — migrate
		var flat map[string]string
		if err := json.Unmarshal(raw, &flat); err != nil {
			return fmt.Errorf("parsing v1 config: %w", err)
		}
		slog.Info("migrating config from v1 to v2 format")
		s.data = storeData{
			Version:       currentVersion,
			ActiveProfile: "default",
			Profiles: map[string]map[string]string{
				"default": flat,
			},
		}
		if err := s.save(); err != nil {
			return fmt.Errorf("saving migrated config: %w", err)
		}
	}

	return nil
}

// save writes the current data to disk with 0600 permissions.
func (s *Store) save() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	if err := os.WriteFile(s.path, raw, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

// activeData returns the map for the active profile, or nil if none.
func (s *Store) activeData() map[string]string {
	if s.data.ActiveProfile == "" {
		return nil
	}
	return s.data.Profiles[s.data.ActiveProfile]
}

// ensureActiveData returns the map for the active profile, creating it if needed.
func (s *Store) ensureActiveData() map[string]string {
	if s.data.ActiveProfile == "" {
		return nil
	}
	m := s.data.Profiles[s.data.ActiveProfile]
	if m == nil {
		m = make(map[string]string)
		s.data.Profiles[s.data.ActiveProfile] = m
	}
	return m
}

// --- Core accessors (delegate to active profile) ---

// Get retrieves a value by key from the active profile. Returns "" if not found.
func (s *Store) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.activeData()
	if m == nil {
		return ""
	}
	return m[key]
}

// Set stores a key-value pair in the active profile and persists to disk.
func (s *Store) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.ensureActiveData()
	if m == nil {
		return fmt.Errorf("no active profile")
	}
	m[key] = value
	return s.save()
}

// Delete removes a key from the active profile and persists to disk.
func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.activeData()
	if m == nil {
		return nil
	}
	delete(m, key)
	return s.save()
}

// ClearAll removes all known keys from the active profile and persists.
func (s *Store) ClearAll() []ClearResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.activeData()
	var results []ClearResult
	for _, key := range AllKeys {
		existed := m != nil && m[key] != ""
		if m != nil {
			delete(m, key)
		}
		results = append(results, ClearResult{Key: key, Success: existed})
	}

	if err := s.save(); err != nil {
		slog.Error("failed to save after clear", "error", err)
		for i := range results {
			results[i].Success = false
			results[i].Error = err
		}
	}

	return results
}

// Has returns true if a key exists and is non-empty in the active profile.
func (s *Store) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.activeData()
	if m == nil {
		return false
	}
	return m[key] != ""
}

// --- Profile management ---

// ActiveProfile returns the name of the active profile.
func (s *Store) ActiveProfile() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.ActiveProfile
}

// SetActiveProfile switches the active profile. Returns error if it doesn't exist.
func (s *Store) SetActiveProfile(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data.Profiles[name]; !ok {
		return fmt.Errorf("profile %q does not exist", name)
	}
	s.data.ActiveProfile = name
	return s.save()
}

// ProfileNames returns sorted list of all profile names.
func (s *Store) ProfileNames() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.data.Profiles))
	for name := range s.data.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// CreateProfile creates a new empty profile. Errors if name is invalid or already exists.
func (s *Store) CreateProfile(name string) error {
	if !profileNameRE.MatchString(name) {
		return fmt.Errorf("invalid profile name %q: must match [a-zA-Z0-9_-]+", name)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data.Profiles[name]; ok {
		return fmt.Errorf("profile %q already exists", name)
	}
	s.data.Profiles[name] = make(map[string]string)
	return s.save()
}

// DeleteProfile removes a profile. Errors if it's the active profile or the last one.
func (s *Store) DeleteProfile(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data.Profiles[name]; !ok {
		return fmt.Errorf("profile %q does not exist", name)
	}
	if name == s.data.ActiveProfile {
		return fmt.Errorf("cannot delete the active profile %q", name)
	}
	if len(s.data.Profiles) <= 1 {
		return fmt.Errorf("cannot delete the last profile")
	}
	delete(s.data.Profiles, name)
	return s.save()
}

// HasProfile returns true if the named profile exists.
func (s *Store) HasProfile(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data.Profiles[name]
	return ok
}

// ProfileData returns a copy of the data for a named profile (for display).
func (s *Store) ProfileData(name string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.data.Profiles[name]
	if !ok {
		return nil
	}
	cp := make(map[string]string, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// ClearProfile clears all known keys from a specific profile.
func (s *Store) ClearProfile(name string) []ClearResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.data.Profiles[name]
	var results []ClearResult
	for _, key := range AllKeys {
		existed := m != nil && m[key] != ""
		if m != nil {
			delete(m, key)
		}
		results = append(results, ClearResult{Key: key, Success: existed})
	}

	if err := s.save(); err != nil {
		slog.Error("failed to save after clear", "error", err)
		for i := range results {
			results[i].Success = false
			results[i].Error = err
		}
	}

	return results
}

// GetForProfile reads a key from a specific profile.
func (s *Store) GetForProfile(profile, key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.data.Profiles[profile]
	if m == nil {
		return ""
	}
	return m[key]
}

// SetForProfile writes a key to a specific profile and persists.
func (s *Store) SetForProfile(profile, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.Profiles[profile]
	if m == nil {
		return fmt.Errorf("profile %q does not exist", profile)
	}
	m[key] = value
	return s.save()
}

// EnsureProfile creates the profile if it doesn't exist, then sets it as active.
func (s *Store) EnsureProfile(name string) error {
	if !profileNameRE.MatchString(name) {
		return fmt.Errorf("invalid profile name %q: must match [a-zA-Z0-9_-]+", name)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data.Profiles[name]; !ok {
		s.data.Profiles[name] = make(map[string]string)
	}
	s.data.ActiveProfile = name
	return s.save()
}

// --- Typed convenience getters/setters ---

func (s *Store) Token() string            { return s.Get(KeyToken) }
func (s *Store) SetToken(v string) error   { return s.Set(KeyToken, v) }

func (s *Store) ControllerHost() string          { return s.Get(KeyControllerHost) }
func (s *Store) SetControllerHost(v string) error { return s.Set(KeyControllerHost, v) }

func (s *Store) Domain() string            { return s.Get(KeyDomain) }
func (s *Store) SetDomain(v string) error   { return s.Set(KeyDomain, v) }

func (s *Store) RefreshToken() string            { return s.Get(KeyRefreshToken) }
func (s *Store) SetRefreshToken(v string) error   { return s.Set(KeyRefreshToken, v) }

func (s *Store) TokenExpiresAt() int64 {
	v := s.Get(KeyTokenExpiresAt)
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseInt(v, 10, 64)
	return n
}

func (s *Store) SetTokenExpiresAt(ms int64) error {
	return s.Set(KeyTokenExpiresAt, strconv.FormatInt(ms, 10))
}

func (s *Store) IdentityCert() string            { return s.Get(KeyIdentityCert) }
func (s *Store) SetIdentityCert(v string) error   { return s.Set(KeyIdentityCert, v) }

func (s *Store) IdentityKey() string            { return s.Get(KeyIdentityKey) }
func (s *Store) SetIdentityKey(v string) error   { return s.Set(KeyIdentityKey, v) }

func (s *Store) IdentityCA() string            { return s.Get(KeyIdentityCA) }
func (s *Store) SetIdentityCA(v string) error   { return s.Set(KeyIdentityCA, v) }

func (s *Store) UpdbUsername() string            { return s.Get(KeyUpdbUsername) }
func (s *Store) SetUpdbUsername(v string) error   { return s.Set(KeyUpdbUsername, v) }

func (s *Store) UpdbPassword() string            { return s.Get(KeyUpdbPassword) }
func (s *Store) SetUpdbPassword(v string) error   { return s.Set(KeyUpdbPassword, v) }

func (s *Store) ControllerCA() string            { return s.Get(KeyControllerCA) }
func (s *Store) SetControllerCA(v string) error   { return s.Set(KeyControllerCA, v) }

func (s *Store) IDPClientID() string            { return s.Get(KeyIDPClientID) }
func (s *Store) SetIDPClientID(v string) error   { return s.Set(KeyIDPClientID, v) }
