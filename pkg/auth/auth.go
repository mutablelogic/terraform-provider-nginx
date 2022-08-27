// auth package manages the authentication tokens
package auth

import (
	"context"
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

	event "github.com/mutablelogic/terraform-provider-nginx/pkg/event"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Path  string
	File  string
	Delta time.Duration
}

type auth struct {
	sync.RWMutex

	delta    time.Duration
	path     string
	tokens   map[string]*token
	modified bool
	ch       chan *event.Event
}

type token struct {
	Token string    `json:"token"`
	Time  time.Time `json:"atime"`
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultFile                 = "auth,json"
	defaultLength               = 32
	defaultDelta                = time.Second * 30
	defaultEventChannelCapacity = 1000
)

const (
	AdminToken = "admin"
)

var (
	reValidName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]+$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New() (*auth, error) {
	this := new(auth)
	this.ch = make(chan *event.Event, defaultEventChannelCapacity)
	this.delta = defaultDelta

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
		this.path = fn
	}

	// Read the file if it exists
	if tokens, err := fileRead(this.path); err != nil {
		return nil, err
	} else {
		this.tokens = tokens
	}

	// If the admin token does not exist, then create it
	if _, ok := this.tokens[AdminToken]; !ok {
		// Create a new token
		this.tokens[AdminToken] = newToken(defaultLength)
	}

	// Write tokens to disk
	if err := fileWrite(this.path, this.tokens); err != nil {
		return nil, err
	}

	// Set delta
	if c.Delta != 0 {
		this.delta = c.Delta
	}

	// Return success
	return this, nil
}

// Run will write the authorization tokens back to disk if they have been modified
func (c *auth) Run(ctx context.Context) error {
	ticker := time.NewTicker(c.delta)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(c.ch)
			_, err := c.writeIfModified()
			return err
		case <-ticker.C:
			if written, err := c.writeIfModified(); err != nil {
				event.NewError(err).Emit(c.ch)
			} else if written {
				event.NewEvent(nil, "Written tokens to disk").Emit(c.ch)
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *auth) String() string {
	str := "<auth"
	str += fmt.Sprintf(" path=%q", c.path)
	for k, v := range c.tokens {
		str += fmt.Sprintf(" %v=%v", k, v)
	}
	str += fmt.Sprint(" delta=", c.delta)
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

// Return event channel
func (c *auth) C() <-chan *event.Event {
	return c.ch
}

// Return true if a token associated with the name already exists
func (c *auth) Exists(name string) bool {
	c.RLock()
	defer c.RUnlock()

	_, ok := c.tokens[name]
	return ok
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
	if name == AdminToken {
		return "", ErrBadParameter.Withf("token is reserved: %q", name)
	}

	// Create a new token
	c.tokens[name] = newToken(defaultLength)

	// Set modified flag
	c.setModified(true)

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
	if name == AdminToken {
		// Rotate the token
		c.tokens[name] = newToken(defaultLength)
	} else {
		// Delete the token
		delete(c.tokens, name)
	}

	// Set modified flag
	c.setModified(true)

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

// setModified sets a new modified value, and returns true if changed
func (c *auth) setModified(modified bool) bool {
	if modified != c.modified {
		c.modified = modified
		return true
	} else {
		return false
	}
}

// write the tokens to disk if modified
func (c *auth) writeIfModified() (bool, error) {
	modified := c.setModified(false)
	if modified {
		if err := fileWrite(c.path, c.tokens); err != nil {
			return modified, err
		}
	}

	// Return success
	return modified, nil
}

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
