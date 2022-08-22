# terraform-provider-nginx

**DRAFT**

A terraform provider to manage nginx configuration. Exposes an API for performing the following:

  * Enumerating a list of configuration files
  * Enabling and disabling configuration files
  * Testing nginx configuration
  * Reloading nginx configuration

This repository is currently in development and is not yet ready for use.


## Server API

The server task provides a REST API for creating, removing, limking and unlinking
configurations. The schema for the API is as follows:

| Method | Path Pattern | Body                                      | Description |
| ------ | ------------ | ----------------------------------------- | ----------- |
| GET    | /            | No body                                   | Returns the list of available configurations |
| GET    | /:name       | No body                                   | Returns a configuration |
| POST   | /:name       | `{ "enabled" : <bool>, "body" : <text> }` | Creates a new configuration |
| DELETE | /:name       | No body                                   | Removes a configuration |
| PATCH  | /:name       |`{ "enabled" : <bool> }`                   | Enables or disables a configuration |

