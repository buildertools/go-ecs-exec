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

The pairs are specified in the form `ENVVAR=FILEPATH`. The list of pairs is space delimited and terminated by a `--` token. Everything after the `--` token will be executed using the `execve(2)` system call. 

The environment variables that are written to files are unset in the environment (for safety). Note that in the following example your secret will not be included in the environment for the process:

```shell
NOT_SECRET=myEnemysPassword \
SOME_SECRET=myPassword \
ecs-exec \
  SOME_SECRET=/run/secrets/some_secret \
  -- \
  env
```

### Base64 Encoded Values

If you're trying to pack a bunch of crap into environment variables because systems like ECS don't support writing files into your environment then you might want to base64 encode those larger or structured values and stuff them into an environment variable. In that case `ecs-exec` can automatically decode those values prior to writing them to files. Simply add a `b64,` prefix to the filename:

```shell
MY_CONFIG=eyJ1c2VybmFtZSI6ImF3c2Zvb2wiLCJzZWNyZXRfa2V5Ijoic3VwZXJzZWNyZXRnYXJiYWdlIn0K \
ecs-exec \
  MY_CONFIG=b64,/run/secrets/creds \
  -- \
  cat /run/secrets/creds
```

## Licensing

See the [LICENSE](LICENSE) file.

## Building it

It is pure Go. Run `go build`. Or maybe some variation like:

    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w -extldflags "-static"'

You can also use the `Makefile` included with the repository.
