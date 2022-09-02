
var "listener" {
    type = string
    description = "Which port server listens on"
    default = ":80"
}

httpserver "main" {
    listen = var.listener
}

