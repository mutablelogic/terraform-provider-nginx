# terraform-provider-nginx

**DRAFT**

A terraform provider to manage nginx configuration. Exposes an API for performing the following:

  * Enumerating a list of configuration files
  * Enabling and disabling configuration files
  * Testing nginx configuration
  * Reloading nginx configuration

There are two elements to the provider:

  * An API gateway which manages the nginx configuration files and server, and listens for requests from the terraform provider.
    This server can be run in a docker container, more details are below.
  * A terraform provider which exposes the API gateway as a resource.

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

