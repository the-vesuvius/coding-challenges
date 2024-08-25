package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:  "diff",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fileReader1, err := os.Open(args[0])
			if err != nil {
				log.Fatalf("openning file %s: %v", args[0], err)
				return
			}
			defer fileReader1.Close()
			fileReader2, err := os.Open(args[1])
			if err != nil {
				log.Fatalf("openning file %s: %v", args[1], err)
				return
			}
			defer fileReader1.Close()

			lines1 := readFile(fileReader1)
			lines2 := readFile(fileReader2)

			results := diff(lines1, lines2)
			for _, line := range results {
				fmt.Println(line)
			}
		},
	}
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err)
	}
}

func readFile(reader io.Reader) []string {
	results := make([]string, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}
	return results
}

func diff(lines1, lines2 []string) []string {
	results := []string{}
	linesLcs := lcsStrings(lines1, lines2)
	idx1, idx2, idxLcs := 0, 0, 0
	for ; idxLcs < len(linesLcs); idxLcs++ {
		for ; idx1 < len(lines1) && lines1[idx1] != linesLcs[idxLcs]; idx1++ {
			results = append(results, "< "+lines1[idx1])
		}
		for ; idx2 < len(lines2) && lines1[idx2] != linesLcs[idxLcs]; idx2++ {
			results = append(results, "> "+lines2[idx2])
		}
		idx1++
		idx2++
	}
	for ; idx1 < len(lines1); idx1++ {
		results = append(results, "< "+lines1[idx1])
	}
	for ; idx2 < len(lines2); idx2++ {
		results = append(results, "> "+lines2[idx2])
	}
	return results
}

func lcsStrings(lines1, lines2 []string) []string {
	m := len(lines1)
	n := len(lines2)
	result := make([][][]string, m+1)
	for i := range result {
		result[i] = make([][]string, n+1)
	}
	for i := 0; i < m+1; i++ {
		for j := 0; j < n+1; j++ {
			if i == 0 || j == 0 {
				result[i][j] = []string{}
			} else if lines1[i-1] == lines2[j-1] {
				result[i][j] = []string{}
				result[i][j] = append(result[i][j], result[i-1][j-1]...)
				result[i][j] = append(result[i][j], lines1[i-1])
			} else {
				if len(result[i-1][j]) > len(result[i][j-1]) {
					result[i][j] = result[i-1][j]
				} else {
					result[i][j] = result[i][j-1]
				}
			}
		}
	}

	return result[m][n]
}

//func lcs(str1, str2 string) string {
//	m := len(str1)
//	n := len(str2)
//	result := make([][]string, m+1)
//	for i := range result {
//		result[i] = make([]string, n+1)
//	}
//	for i := 0; i < m+1; i++ {
//		for j := 0; j < n+1; j++ {
//			if i == 0 || j == 0 {
//				result[i][j] = ""
//			} else if str1[i-1] == str2[j-1] {
//				result[i][j] = result[i-1][j-1] + string(str1[i-1])
//			} else {
//				if len(result[i-1][j]) > len(result[i][j-1]) {
//					result[i][j] = result[i-1][j]
//				} else {
//					result[i][j] = result[i][j-1]
//				}
//			}
//		}
//	}
//
//	return result[m][n]
//}
