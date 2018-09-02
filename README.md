# builddepends

*Debian Build-Depends graph generator*

Generates build dependency graphs for Debian control files. Includes a cli program for scanning directories and outputting the DOT graph.

## The library, `builddepends`

### Installation

```bash
$ go get -u github.com/opx-infra/builddepends
```

### Usage

```go
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/opx-infra/builddepends"
)

func main() {
	files, _ := ioutil.ReadDir(".")
	dirs, _ := builddepends.DebianDirectories(files)
	controls, _ := builddepends.ParseControls(dirs)
	graph, _ := builddepends.DependencyGraph(controls)
	fmt.Print(graph)
}
```

## The program, `bd`

### Installation

```bash
$ curl -sLo bd https://github.com/opx-infra/builddepends/releases/download/v0.1.0/bd-linux-amd64
$ chmod +x bd
```

### Usage

With one or more directories present, run `bd`.

**Bonus!** Use the flag `--build-order` for easy distributed building.

```bash
$ for r in opx-nas-acl opx-nas-daemon opx-alarm opx-logging opx-common-utils; do
    git clone "https://github.com/open-switch/$r"
  done

$ bd
strict digraph "builddepends" {
"opx-alarm";
"opx-common-utils";
"opx-logging";
"opx-nas-acl";
"opx-nas-daemon";
"opx-common-utils" -> "opx-logging";
"opx-nas-acl" -> "opx-common-utils";
"opx-nas-acl" -> "opx-logging";
"opx-nas-daemon" -> "opx-common-utils";
"opx-nas-daemon" -> "opx-logging";
"opx-nas-daemon" -> "opx-nas-acl";
}
```

## Why not [controlgraph](https://github.com/opx-infra/controlgraph)?

Python takes too long to start. Using the directories from above, this Go rewrite runs ~50x faster. The program name is 6x shorter to celebrate this.

```bash
$ time controlgraph --graph >/dev/null

real	0m0.483s
user	0m0.364s
sys 	0m0.102s

$ time bd >/dev/null

real	0m0.009s
user	0m0.004s
sys 	0m0.004s
```
