package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/opx-infra/builddepends"
	flag "github.com/spf13/pflag"
)

const version = "0.1.1"

func main() {
	flag.Usage = func() {
		fmt.Printf("bd - Debian Build-Depends Graph Generator\n\nUsage:\n")
		flag.PrintDefaults()
		fmt.Printf("\n")
	}
	reverseGraph := flag.BoolP("build-order", "b", false, "Generate build order instead of dependency order")
	versionMode := flag.BoolP("version", "V", false, "Print version and exit")
	flag.Parse()

	if *versionMode {
		fmt.Printf("bd %s\n", version)
		os.Exit(0)
	}

	var files []os.FileInfo

	if len(flag.Args()) > 0 {
		// processed user-supplied files
		for _, d := range flag.Args() {
			stat, err := os.Lstat(d)
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, stat)
		}
	} else {
		// use all files in current directory
		var err error
		files, err = ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}
	}

	directories, err := builddepends.DebianDirectories(files)
	if err != nil {
		log.Fatal(err)
	}

	controls, err := builddepends.ParseControls(directories)
	if err != nil {
		log.Fatal(err)
	}

	var graph string
	if *reverseGraph {
		graph, err = builddepends.BuildGraph(controls)
	} else {
		graph, err = builddepends.DependencyGraph(controls)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(graph)
}
