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
func BuildGraph(controls map[string]*control.Control) (string, error) {
	return graph(controls, true)
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
func DependencyGraph(controls map[string]*control.Control) (string, error) {
	return graph(controls, false)
}

// ParseControls returns the Control struct for each directory
func ParseControls(dirs []string) (map[string]*control.Control, error) {
	m := make(map[string]*control.Control)

	for _, d := range dirs {
		c, err := control.ParseControlFile(fmt.Sprintf("%s/debian/control", d))
		if err != nil {
			return m, err
		}
		m[d] = c
	}

	return m, nil
}

// Locals /////////////////////////////////////////////////////////////////////

const (
	digraphStart = "strict digraph \"\" {\n"
	digraphEnd   = "}"
)

// binPkgToDirectory returns a lookup map for binary package to directory translation
func binPkgToDirectory(controls map[string]*control.Control) map[string]string {
	m := make(map[string]string)

	for d, c := range controls {
		for _, bin := range c.Binaries {
			m[bin.Package] = d
		}
	}

	return m
}

func graph(controls map[string]*control.Control, reverse bool) (string, error) {
	lines := make(map[string]bool)
	lookup := binPkgToDirectory(controls)

	for d, c := range controls {
		storeLine(node(d), lines)
		for _, dep := range strings.Split(c.Source.BuildDepends.String(), ", ") {
			bareDep := strings.Split(dep, " ")[0]
			buildDir, available := lookup[bareDep]
			if available {
				storeLine(node(buildDir), lines)
				if reverse {
					storeLine(edge(buildDir, d), lines)
				} else {
					storeLine(edge(d, buildDir), lines)
				}
			}
		}
	}

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

	var b strings.Builder
	b.WriteString(digraphStart)

	for _, node := range nodes {
		b.WriteString(node)
	}

	for _, edges := range edges {
		b.WriteString(edges)
	}

	b.WriteString(digraphEnd)
	return b.String(), nil
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
