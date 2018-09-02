package builddepends

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"pault.ag/go/debian/control"
)

// Exported ///////////////////////////////////////////////////////////////////

// BuildGraph returns the DOT graph for a distributable build order
func BuildGraph(controls map[string]*control.Control, sorted bool) (string, error) {
	return graph(controls, true, sorted)
}

// DebianDirectories returns all directories with debian/control files
func DebianDirectories(files []os.FileInfo) ([]string, error) {
	var directories []string

	for _, file := range files {
		if file.IsDir() {

			stat, err := os.Lstat(fmt.Sprintf("%s/debian/control", file.Name()))
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return directories, err
			}

			switch mode := stat.Mode(); {
			case mode.IsRegular():
				directories = append(directories, file.Name())
			}

		}
	}

	return directories, nil
}

// DependencyGraph returns the DOT graph for the build dependencies locally available
func DependencyGraph(controls map[string]*control.Control, sorted bool) (string, error) {
	return graph(controls, false, sorted)
}

// ParseControls returns the Control struct for each directory
func ParseControls(dirs []string) (map[string]*control.Control, error) {
	cs := make(map[string]*control.Control)

	for _, dir := range dirs {
		con, err := control.ParseControlFile(fmt.Sprintf("%s/debian/control", dir))
		if err != nil {
			return cs, err
		}
		cs[dir] = con
	}

	return cs, nil
}

// Locals /////////////////////////////////////////////////////////////////////

// binPkgToDirectory returns a lookup map for binary package to directory translation
func binPkgToDirectory(controls map[string]*control.Control) map[string]string {
	lookup := make(map[string]string)

	for dir, con := range controls {
		for _, bin := range con.Binaries {
			lookup[bin.Package] = dir
		}
	}

	return lookup
}

func graph(controls map[string]*control.Control, reverse bool, sorted bool) (string, error) {
	lines := make(map[string]bool)
	lookup := binPkgToDirectory(controls)

	// Process nodes and edges
	for dir, con := range controls {
		storeLine(node(dir), lines)
		for _, dep := range strings.Split(con.Source.BuildDepends.String(), ", ") {
			bareDep := strings.Split(dep, " ")[0]
			buildDir, available := lookup[bareDep]
			if available {
				storeLine(node(buildDir), lines)
				if reverse {
					storeLine(edge(buildDir, dir), lines)
				} else {
					storeLine(edge(dir, buildDir), lines)
				}
			}
		}
	}

	var bob strings.Builder
	bob.WriteString("strict digraph \"builddepends\" {\n")

	if sorted {
		var nodes []string
		var edges []string
		for line := range lines {
			if strings.Contains(line, "->") {
				edges = append(edges, line)
			} else {
				nodes = append(nodes, line)
			}
		}
		sort.Strings(nodes)
		sort.Strings(edges)

		for _, n := range nodes {
			bob.WriteString(n)
		}
		for _, e := range edges {
			bob.WriteString(e)
		}
	} else {
		for line := range lines {
			bob.WriteString(line)
		}
	}

	bob.WriteString("}\n")
	return bob.String(), nil
}

func edge(from, to string) string {
	return fmt.Sprintf("\"%s\" -> \"%s\";\n", from, to)
}

func node(name string) string {
	return fmt.Sprintf("\"%s\";\n", name)
}

func storeLine(line string, lines map[string]bool) {
	_, written := lines[line]
	if !written {
		lines[line] = true
	}
}
