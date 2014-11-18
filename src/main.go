package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// TODO: Split this off into some sort of model package
type Revision string
type RevisionContents struct {
	Revision Revision
	Paths []string
}

type Alias struct {
	Branch string
	Revision Revision
}

func (alias Alias) PrintVerbose() {
	fmt.Printf("Branch: \"%s\",\tRevision: \"%s\"\n", alias.Branch, string(alias.Revision))
	for _, path := range alias.Revision.Load().Paths {
		fmt.Printf("\tPath: \"%s\"\n", path)
	}
}

// TODO: Create a package that wraps all of the calls to git commands
func (revision Revision) Load() *RevisionContents {
	out, err := exec.Command("git", "ls-tree", "-r", string(revision)).Output()
	if err != nil {
		log.Fatal(err)
	}
	listOutput := strings.Trim(string(out), "\n ")
	lines := strings.Split(listOutput, "\n")
	paths := make([]string, len(lines))
	for index, line := range lines {
		line = strings.Replace(lines[index], "\t", " ", -1)
		lineParts := strings.Split(line, " ")
		paths[index] = lineParts[len(lineParts) - 1]
	}
	return &RevisionContents{revision, paths}
}

func parseBranchListLine(line string) Alias {
	line = strings.Trim(line, "* ")
	splitLine := strings.Split(line, " ")
	masterName := splitLine[0]
	for _, lineComponent := range splitLine[1:] {
		if len(lineComponent) == 40 {
			revisionHash := lineComponent
			return Alias{masterName, Revision(revisionHash)}
		}
	}
	return Alias{Branch: masterName}
}

func ListBranches() []Alias {
	out, err := exec.Command("git", "branch", "-av", "--list", "--abbrev=40").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.Trim(string(out), " \n"), "\n")
	aliases := make([]Alias, len(lines))
	index := 0
	for _, line := range lines {
		if line != "" {
			aliases[index] = parseBranchListLine(line)
			index += 1
		}
	}
	return aliases
}

// TODO: Serve a webpage instead of printing to stdout
func main() {
	for _, alias := range ListBranches() {
		alias.PrintVerbose()
	}
}