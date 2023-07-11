package example

import future.keywords.if

default allow := false

allow if {
    input.method == "GET"
    input.user == "bob"
}
