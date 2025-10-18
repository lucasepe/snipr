# Examples

## Shell

This is a basic snippet with shell command:

```sh
echo "Hello, (anonymous) mdrun!"
```

By adding {name=hello} you can annotate the block and give it a descriptive name:

```sh {name=hello}
echo "Hello, mdrun!"
```


By adding {skip=true} the block wull not be executed:

```sh {name=hello-skipped skip=true}
echo "Hello, skipped!"
```


Blocks can contain multiple lines as well:

```sh {name=multiple-echoes}
echo "1"
echo "2"
echo "3"
```

Use {timeout=2} to set a timeout (in seconds) for the block execution:

```sh {name=hello-timed timeout=2}
sleep 1
echo "Hello..."
```

Since each block is executed as a single script, it works seamlessly with cd, pushd, and similar commands:

```sh {name=pusher}
temp_dir=$(mktemp -d -t "mdrun-XXXXXXX")
pushd $temp_dir
echo "hi!" > hi.txt
pwd
cat hi.txt
popd
pwd
```

With {file=/path/to/file.txt} you can save a block’s output into specified file:

```bash {name=you-to-file file=/tmp/you.txt}
whoami
```


With {export=VAR_NAME} you can save a block’s output into an environment variable:

```bash {name=you-to-env export=YOU}
whoami
```

…and use it in a later block:

```bash {name=hello-you}
echo "Hello, ${YOU}!"
```

## Go

You can also execute Go code snippets:

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Printf("Hello from Go, %s!\n", os.Getenv("YOU"))
}
```

```go {name=test-import}
package main

import (
    "fmt"
    "github.com/lucasepe/x/env"
)

func main() {
    fmt.Printf("Hello from Go again, %s!\n", env.Str("YOU", "nobody"))
}
```