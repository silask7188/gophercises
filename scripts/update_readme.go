package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CodeFile struct {
	Name      string
	LineCount int
}

type Exercise struct {
	DirName     string
	Title       string
	Goal        string
	SummaryFile string
	CodeFiles   []CodeFile
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

		dirPath := filepath.Join(rootDir, name)
		ex := parseExercise(dirPath, name)
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

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines
}

func parseExercise(dirPath, dirName string) Exercise {
	ex := Exercise{
		DirName: dirName,
		Title:   strings.ReplaceAll(dirName, "-", " "),
	}

	entries, err := os.ReadDir(dirPath)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				ex.SummaryFile = e.Name()
				title, goal := parseMarkdownGoal(filepath.Join(dirPath, e.Name()))
				if title != "" {
					ex.Title = title
				}
				if goal != "" {
					ex.Goal = goal
				}
			} else if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
				filePath := filepath.Join(dirPath, e.Name())
				lines := countLines(filePath)
				ex.CodeFiles = append(ex.CodeFiles, CodeFile{Name: e.Name(), LineCount: lines})
				ex.TotalLines += lines
				ex.HasCode = true
			}
		}
	}

	if ex.Goal == "" {
		ex.Goal = "go exercise solution"
	}

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
		for _, codeFile := range ex.CodeFiles {
			sb.WriteString(fmt.Sprintf("- **file:** [%s/%s](./%s/%s) (``%d lines``)\n", ex.DirName, codeFile.Name, ex.DirName, codeFile.Name, codeFile.LineCount))
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
