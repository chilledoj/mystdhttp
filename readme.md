# My Go Standard Library Website

A basic website built with just the go standard library. This provides an example of how the standard libray api is used.

## Running

Use the `makefile` and run `make run` which will pre-populate the in-memory db with data. A second run command is available which will leave the in-memory db empty - `make run-empty`.

The binary takes 3 flags:

- **-host** _string_ Host http address to listen on
- **-init-tasks** _bool_ Set to false to not prepopulate the in-memory db (default true). N.B. As the default is true, the syntax for setting to false is to use `-init-tasks=false`.
- **-port** _string_ Port number for http listener (default "8000")
