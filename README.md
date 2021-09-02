# cobra2snooty
[![CI](https://github.com/mongodb-labs/cobra2snooty/actions/workflows/pr.yml/badge.svg)](https://github.com/mongodb-labs/cobra2snooty/actions/workflows/pr.yml)

## Generate Snooty docs for the entire command tree

This program can actually generate docs for the `mongocli` command in the MongoDB CLI project

```go
package main

import (
	"log"
	"os"

	"github.com/mongodb/mongocli/internal/cli/root"
	"github.com/mongodb-labs/cobra2snooty"
)

func main() {
	var profile string
	const docsPermissions = 0766
	if err := os.MkdirAll("./docs/command", docsPermissions); err != nil {
		log.Fatal(err)
	}

	mongocli := root.Builder(&profile, []string{})

	if err := cobra2snooty.GenSnootyTree(mongocli, "./docs/command"); err != nil {
		log.Fatal(err)
	}
}
```

This will generate a whole series of files, one for each command in the tree, in the directory specified (in this case "./docs/command")


## License

`cobra2snooty` is released under the Apache 2.0 license. See [LICENSE](LICENSE)
