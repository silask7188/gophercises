package main

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Package struct {
	Name      string
	RelPath   string
	LineCount int
}

type Exercise struct {
	DirName     string
	Title       string
	Goal        string
	SummaryFile string
	Packages    []Package
	TotalLines  int
	HasCode     bool
}

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting root directory: %v\n", err)
		os.Exit(1)
	}

	exercises := scanWorkspace(rootDir)
	readmePath := filepath.Join(rootDir, "README.md")
	generateReadme(readmePath, exercises)
}

func scanWorkspace(rootDir string) []Exercise {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Printf("error reading root dir: %v\n", err)
		return nil
	}

	var exercises []Exercise

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "scripts" {
			continue
		}

		ex := parseExercise(rootDir, name)
		exercises = append(exercises, ex)
	}

	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].DirName < exercises[j].DirName
	})

	return exercises
}

func countLines(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // fix: check error im too lazy to fix ita
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines
}

func getGoPackageName(filePath string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
	if err != nil || node == nil || node.Name == nil {
		return ""
	}
	return node.Name.Name
}

func parseExercise(rootDir, dirName string) Exercise {
	dirPath := filepath.Join(rootDir, dirName)
	ex := Exercise{
		DirName: dirName,
		Title:   strings.ReplaceAll(dirName, "-", " "),
	}

	pkgMap := make(map[string]*Package)

	_ = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() != dirName && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(rootDir, path)
		relDir := filepath.Dir(relPath)

		if strings.HasSuffix(d.Name(), ".md") {
			if filepath.Dir(path) == dirPath && ex.SummaryFile == "" {
				ex.SummaryFile = d.Name()
				title, goal := parseMarkdownGoal(path)
				if title != "" {
					ex.Title = title
				}
				if goal != "" {
					ex.Goal = goal
				}
			}
		} else if strings.HasSuffix(d.Name(), ".go") {
			pkgName := getGoPackageName(path)
			if pkgName == "" {
				return nil
			}
			lines := countLines(path)
			ex.TotalLines += lines
			ex.HasCode = true

			key := relDir
			if p, exists := pkgMap[key]; exists {
				p.LineCount += lines
			} else {
				pkgMap[key] = &Package{
					Name:      pkgName,
					RelPath:   relDir,
					LineCount: lines,
				}
			}
		}
		return nil
	})

	if ex.Goal == "" {
		ex.Goal = "go exercise solution"
	}

	var packages []Package
	for _, p := range pkgMap {
		packages = append(packages, *p)
	}

	sort.Slice(packages, func(i, j int) bool {
		if packages[i].Name == "main" && packages[j].Name != "main" {
			return true
		}
		if packages[i].Name != "main" && packages[j].Name == "main" {
			return false
		}
		if packages[i].Name != packages[j].Name {
			return packages[i].Name < packages[j].Name
		}
		return packages[i].RelPath < packages[j].RelPath
	})

	ex.Packages = packages

	return ex
}

func parseMarkdownGoal(filePath string) (string, string) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var title string
	var goalLines []string
	inGoal := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "# ") && title == "" {
			title = strings.TrimPrefix(trimmed, "# ")
			continue
		}

		if strings.HasPrefix(strings.ToLower(trimmed), "## goal") {
			inGoal = true
			continue
		}

		if inGoal {
			if strings.HasPrefix(trimmed, "## ") {
				inGoal = false
				continue
			}
			if trimmed != "" {
				goalLines = append(goalLines, line)
			}
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	goal := strings.Join(goalLines, "\n")
	return title, goal
}

func generateReadme(readmePath string, exercises []Exercise) {
	var sb strings.Builder
	totalRepoLines := 0

	for _, ex := range exercises {
		totalRepoLines += ex.TotalLines
	}

	sb.WriteString("# gophercises & go projects\n\n")
	sb.WriteString("solutions and notes for go coding exercises\n\n")
	sb.WriteString(fmt.Sprintf("total lines of code: ``%d``\n\n", totalRepoLines))
	sb.WriteString("## solutions\n\n")

	for _, ex := range exercises {
		if ex.TotalLines > 0 {
			sb.WriteString(fmt.Sprintf("### [%s](./%s) — ``%d lines``\n", ex.Title, ex.DirName, ex.TotalLines))
		} else {
			sb.WriteString(fmt.Sprintf("### [%s](./%s)\n", ex.Title, ex.DirName))
		}
		sb.WriteString(fmt.Sprintf("%s\n\n", ex.Goal))

		if ex.SummaryFile != "" {
			sb.WriteString(fmt.Sprintf("- **summary:** [%s/%s](./%s/%s)\n", ex.DirName, ex.SummaryFile, ex.DirName, ex.SummaryFile))
		}
		for _, pkg := range ex.Packages {
			sb.WriteString(fmt.Sprintf("- **package:** [%s](./%s) (``%d lines``)\n", pkg.Name, pkg.RelPath, pkg.LineCount))
		}
		sb.WriteString("\n")
	}

	err := os.WriteFile(readmePath, []byte(sb.String()), 0644)
	if err != nil {
		fmt.Printf("error writing readme: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("readme successfully updated")
}
