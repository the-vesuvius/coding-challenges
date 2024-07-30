package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:  "diff",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			//fileReader1, err := os.Open(args[0])
			//if err != nil {
			//	log.Fatalf("openning file %s: %v", args[0], err)
			//	return
			//}
			//defer fileReader1.Close()
			//fileReader2, err := os.Open(args[1])
			//if err != nil {
			//	log.Fatalf("openning file %s: %v", args[1], err)
			//	return
			//}
			//defer fileReader1.Close()
			fmt.Println(lcs("abc", "ac"))
		},
	}
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err)
	}
}

func lcs(str1, str2 string) string {
	m := len(str1)
	n := len(str2)
	result := make([][]string, m+1)
	for i := range result {
		result[i] = make([]string, n+1)
	}
	for i := 0; i < m+1; i++ {
		for j := 0; j < n+1; j++ {
			if i == 0 || j == 0 {
				result[i][j] = ""
			} else if str1[i-1] == str2[j-1] {
				result[i][j] = result[i-1][j-1] + string(str1[i-1])
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
