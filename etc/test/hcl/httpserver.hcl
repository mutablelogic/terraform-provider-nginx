
var "listener" {
    type = "string"
    description = "Which port server listens on"
    default = ":80"
}

var "test" {
    type = "string"
}

router "main-router" {}

httpserver "main-server" {
    router = router.main-router
}

/*
var "test" {
    type = "string"
    description = "Which port server listens on"
    default = ":80"
}

router "test-router" {}

*/
