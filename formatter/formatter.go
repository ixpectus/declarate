package formatter

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ixpectus/declarate/converter"
	"github.com/ixpectus/declarate/suite"
	"gopkg.in/yaml.v3"
)

var (
	idsFrom = 300
	idsTo   = 430
	idsUsed = map[int]bool{}
)

type Converter struct {
	sourceDir string
	targetDir string
}

func New(sourceDir, targetDir string) *Converter {
	return &Converter{
		sourceDir: sourceDir,
		targetDir: targetDir,
	}
}

func (c *Converter) Format() error {
	s := suite.New(c.sourceDir, suite.RunConfig{})
	tt, err := s.AllTests(c.sourceDir)
	if err != nil {
		return err
	}

	for _, v := range tt {
		if !strings.Contains(v, ".yaml") {
			continue
		}

		tests := []converter.DeclarateTest{}
		relatedName := strings.ReplaceAll(v, c.sourceDir, "")
		if relatedName == "" {
			relatedName = path.Base(c.sourceDir)
		}
		data, err := readLines(v)
		if err != nil {
			log.Printf("failed to convert %s, %v", v, err)
			continue
		}

		err = yaml.Unmarshal(data, &tests)

		if len(tests) > 0 && tests[0].Definition != nil && len(tests[0].Definition.Tags) > 0 && tests[0].Definition.ID == 0 {
			v := idsFrom
			if idsUsed[v] {
				for i := v; i <= idsTo; i++ {
					if !idsUsed[i] {
						v = i
						idsUsed[v] = true
						break
					}
				}
			} else {
				idsUsed[v] = true
			}

			tests[0].Definition.ID = v
		}
		for k, v := range tests {
			v.ResponseTmpls = strings.ReplaceAll(v.ResponseTmpls, "$matchRegexp(^.+$)", "$any")
			tests[k] = v
			if len(v.Steps) == 1 {
				name := v.Name
				tests[k] = v.Steps[0]
				tests[k].Name = name
			}
			if tests[k].Poll != nil && tests[k].Poll.Interval == 1*time.Second {
				tests[k].Poll.Interval = 0
			}
		}
		if err != nil {
			log.Printf("failed to unmarshall on convert %s, %v", v, err)
			continue
		}
		bb, _ := yaml.Marshal(tests)
		targetFile := c.targetDir + "/" + relatedName
		res := strings.ReplaceAll(string(bb), "|-", "|")

		if err := os.MkdirAll(path.Dir(targetFile), os.ModePerm); err != nil {
			log.Printf("failed to mkdir for file %s, %v", targetFile, err)
		}
		err = os.WriteFile(targetFile, []byte(res), os.ModePerm)
		if err != nil {
			log.Printf("failed to write to file %s, %v", targetFile, err)
			continue
		}
	}

	return nil
}

func readLines(filePath string) ([]byte, error) {
	readFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, strings.TrimRight(fileScanner.Text(), " "))
	}

	readFile.Close()
	res := []byte{}
	for _, line := range fileLines {
		res = append(res, []byte(line)...)
		res = append(res, []byte("\n")...)
	}
	return res, nil
}
