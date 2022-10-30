package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// like this: MASK: [rule_type] [rule]

const (
	typeMatch = iota
	typeRegexp
)

func main() {
	// search repo all comment with MASK: line and mask with rule in after MASK:
	args := os.Args[1:]
	if len(args) == 0 {
		panic("please input repo path")
	}
	repoPath := args[0]
	// walk repo, exclude .dir
	walkDirSearchReplace(repoPath)
}

func walkDirSearchReplace(dir string) {
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		fmt.Printf("walk dir: %s\n", path)
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		// if isn't go file, skip
		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}
		// search MASK: line
		input, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// if line start with MASK: then replace with rule
		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			cline := line
			if strings.Contains(line, "// MASK:") && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/*") {
				// get last index with MASK:

				ruleIndex := strings.LastIndex(line, "// MASK:")
				codeStr, ruleStr := line[:ruleIndex], line[ruleIndex+7:]
				if !strings.HasSuffix(ruleStr, "]") {
					continue
				}
				ruleType := ruleStr[strings.Index(ruleStr, "[")+1 : strings.Index(ruleStr, "]")]
				rule := ruleStr[strings.LastIndex(ruleStr, "[")+1 : strings.LastIndex(ruleStr, "]")]
				fmt.Println("ruleType:", ruleType, "rule:", rule)
				// replace with rule
				switch ruleType {
				case "match":
					// replace with rule
					codeStr = strings.ReplaceAll(codeStr, rule, "******")
				case "regexp":
					re, _ := regexp.Compile(rule)
					codeStr = re.ReplaceAllString(codeStr, "******")
				}
				line = codeStr + "// MASK_DONE"
				fmt.Printf("replace result:\n \tsource: %s\n \tdest: %s\n", cline, line)
			}
			if cline != line {
				lines[i] = line
			}
		}
		output := strings.Join(lines, "\n")
		err = os.WriteFile(path, []byte(output), 0644)
		// replace with rule
		return nil
	})
}
