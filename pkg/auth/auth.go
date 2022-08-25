// auth package manages the authentication tokens
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	// Module imports

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Path string
	File string
}

type auth struct {
	sync.RWMutex
	filename string
	tokens   map[string]*token
}

type token struct {
	Token string    `json:"token"`
	Time  time.Time `json:"atime"`
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultFile   = "auth,json"
	defaultLength = 32
	adminName     = "admin"
)

var (
	reValidName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]+$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New() (*auth, error) {
	this := new(auth)

	// Check for path
	if stat, err := os.Stat(c.Path); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, ErrBadParameter.Withf("not a directory: %q", c.Path)
	}

	// Check for filename
	if c.File == "" {
		c.File = defaultFile
	}

	// Set filename
	if fn, err := filepath.Abs(filepath.Join(c.Path, c.File)); err != nil {
		return nil, err
	} else {
		this.filename = fn
	}

	// Read the file if it exists
	if tokens, err := fileRead(this.filename); err != nil {
		return nil, err
	} else {
		this.tokens = tokens
	}

	// If the admin token does not exist, then create it
	if _, ok := this.tokens[adminName]; !ok {
		// Create a new token
		this.tokens[adminName] = newToken(defaultLength)
	}

	// Write tokens to disk
	if err := fileWrite(this.filename, this.tokens); err != nil {
		return nil, err
	}

	// Return success
	return this, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *auth) String() string {
	str := "<auth"
	str += fmt.Sprintf(" filename=%q", c.filename)
	for k, v := range c.tokens {
		str += fmt.Sprintf(" %v=%v", k, v)
	}
	return str + ">"
}

func (t *token) String() string {
	str := "<token"
	str += fmt.Sprintf(" token=%q", t.Token)
	str += fmt.Sprintf("  last_accessed=%q", t.Time.Format(time.RFC3339))
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return true if a token associated with the name already exists
func (c *auth) Exists(name string) bool {
	c.RLock()
	defer c.RUnlock()

	_, ok := c.tokens[name]
	return ok
}

// Return true if a token associated with the name is the admin token
func (c *auth) IsAdmin(name string) bool {
	return name == adminName
}

// Create a new token associated with a name and return it.
func (c *auth) Create(name string) (string, error) {
	c.Lock()
	defer c.Unlock()

	// If the name is invalid, then return an error
	if !reValidName.MatchString(name) {
		return "", ErrBadParameter.Withf("invalid name: %q", name)
	}
	// If the name exists already, then return an error
	if _, ok := c.tokens[name]; ok {
		return "", ErrDuplicateEntry.Withf("token already exists: %q", name)
	}
	// If the name is the admin token, then return an error
	if name == adminName {
		return "", ErrBadParameter.Withf("token is reserved: %q", name)
	}

	// Create a new token
	c.tokens[name] = newToken(defaultLength)

	// Write tokens to disk
	if err := fileWrite(c.filename, c.tokens); err != nil {
		delete(c.tokens, name)
		return "", err
	}

	// Success: return the token value
	return c.tokens[name].Token, nil
}

// Revoke a token associated with a name. For the admin token, it is
// rotated rather than revoked.
func (c *auth) Revoke(name string) error {
	c.Lock()
	defer c.Unlock()

	// If the name does not exist, then return an error
	if _, ok := c.tokens[name]; !ok {
		return ErrNotFound.Withf("token already exists: %q", name)
	}

	// Either delete or rotate the token
	if name == adminName {
		// Rotate the token
		c.tokens[name] = newToken(defaultLength)
	} else {
		// Delete the token
		delete(c.tokens, name)
	}

	// Write tokens to disk
	if err := fileWrite(c.filename, c.tokens); err != nil {
		delete(c.tokens, name)
		return err
	}

	// Return success
	return nil
}

// Return all token names and their last access times
func (c *auth) Enumerate() map[string]time.Time {
	c.RLock()
	defer c.RUnlock()

	var result = make(map[string]time.Time)
	for k, v := range c.tokens {
		result[k] = v.Time
	}

	// Return the result
	return result
}

// Returns the name of the token if a value matches. Updates
// the access time for the token. If token with value not
// found, then return empty string
func (c *auth) Matches(value string) string {
	c.Lock()
	defer c.Unlock()

	for k, v := range c.tokens {
		if v.Token == value {
			v.Time = time.Now()
			// TODO: Update time
			return k
		}
	}

	// Token not found
	return ""
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func fileRead(filename string) (map[string]*token, error) {
	var result = map[string]*token{}

	// If the file doesn't exist, return empty result
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return result, nil
	} else if err != nil {
		return nil, err
	}

	// Open the file
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	// Decode the file
	if err := json.NewDecoder(fh).Decode(&result); err != nil {
		return nil, err
	}

	// Return success
	return result, nil
}

func fileWrite(filename string, tokens map[string]*token) error {
	if tokens == nil {
		return ErrBadParameter.Withf("tokens is nil")
	}

	// Create the file
	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	// Write the tokens
	if err := json.NewEncoder(fh).Encode(tokens); err != nil {
		return err
	}

	// Return success
	return nil
}

func newToken(length int) *token {
	return &token{
		Token: generateToken(length),
		Time:  time.Now(),
	}
}

func generateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
