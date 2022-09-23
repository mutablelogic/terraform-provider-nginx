
var "listener" {
    type = "string"
    description = "Which port server listens on"
    default = ":80"
}

httpserver "server" {
    listen = var.listener
}

tokenauth-gateway "tokenauth" {
    router = httpserver.server
    prefix = "/api/tokenauth/v1"
}

nginx-gateway "nginx" {
    router = httpserver.server
    prefix = "/api/nginx/v1"
    middleware = [ "tokenauth" ]
}
