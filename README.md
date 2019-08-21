# go-ecs-exec 

Writes the content of an environment variables to a named files and execute another program.

## Using it

The program takes environment variable name and filename pairs as well as a full command line command to execute. For example:

```shell
TESTVAR=testValue1 \
SECOND_VAR=testValue2 \
ecs-exec \
  TESTVAR=/run/secrets/testvar \
  SECOND_VAR=/run/secrets/secondvar \
  -- \
  cat /run/secrets/testvar
```

The pairs are specified in the form `ENVVAR=FILEPATH`. The list of pairs is space delimited and terminated by a `--` token. Everything after the `--` token will be executed using the `execve(2)` system call. The environment variables that are written to files are unset in the environment (for safety).

## Licensing

See the [LICENSE](LICENSE) file.

## Building it

It is pure Go. Run `go build`. Or maybe some variation like:

    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w -extldflags "-static"'

You can also use the `Makefile` included with the repository.
