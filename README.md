Builds and runs a go module and watches files in a given directory to rebuild and run when a go file in the directory is modified. Plan to extend it eventually to be more configurable.

Build it with "go build ."

Run it with "go-hot-reloader [go-project-dir] [go-proj-executable-name] [poll-wait-time-ms]"
By default the poll-wait-time-ms is 2000 and can be omitted.

Stop it with Ctrl+C (SIGTERM)
