package router

/*

// AddMiddleware adds a middleware handler with a unique key.
func (r *router) AddMiddleware(key string, fn func(http.HandlerFunc) http.HandlerFunc) error {
	// Preconditions
	if !reValidName.MatchString(key) {
		return ErrBadParameter.Withf("AddMiddleWare: %q", key)
	}
	if fn == nil {
		return ErrBadParameter.Withf("AddMiddleWare: %q", key)
	}

	// Check for duplicate entry
	r.RLock()
	_, exists := r.middleware[key]
	r.RUnlock()
	if exists {
		return ErrDuplicateEntry.Withf("AddMiddleWare: %q", key)
	}

	// Set middleware mapping
	r.Lock()
	r.middleware[key] = fn
	r.Unlock()

	// Return success
	return nil
}

// SetMiddleware binds an array of middleware functions to a prefix. The prefix should
// already exist in the router.
func (r *router) SetMiddleware(prefix string, chain ...string) error {
	prefix = normalizePath(prefix, true)
	fmt.Println("SetMiddleware", prefix, chain)

	return nil
}
*/
