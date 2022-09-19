// tokenauth-gw plugin is a gateway for token authentication. It provides:
//
//  - Middleware for token authentication, which checks for a "Authorization: Token <value>"
//    http header and validates against the tokenauth plugin;
//  - HTTP handlers for creating and revoking tokens;
//  - Methods which can return a token name from a token value.
//
package main
