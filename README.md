# Service wrapper

> **This little wrapper designed to allow service being controlled trough Telegram bot.**

## How to use

### Container preparation

1. Download the binary from [Releases](/../../releases) page or build with `CGO_ENABLED=0 go build -o wrapper ./src`
1. Copy the binary to docker container
1. Create `service.yml` description file somewhere in the container. Example:
    ```yml
    service_name: Test service
    lines_to_preserve: 100
    # if true lines_to_preserve works separately for both outputs
    separate_stdout_stderr: true
    command: /bin/sh
    args:
      - -c
      - | # we also can use inline comments and multiline arguments
        echo 1
        sleep 1s
        echo 2 >&2
        echo 3
    ```
1. Change container entrypoint to `["/path/to/wrapper", "/path/to/service.yml"]`
1. Well done :)

### Container launch
Additionally to regular env, pass two variables more — `WRAPPER_TELEGRAM_CHAT` and `WRAPPER_TELEGRAM_TOKEN`. That's all
