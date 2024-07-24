package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"github.com/spf13/cobra"
)

type result struct {
	linesNumber int64
	bytesNumber int64
	wordsNumber int64
}

type outputConfig struct {
	showLines bool
	showBytes bool
	showWords bool
}

func main() {
	outputCfg := outputConfig{}
	cmd := cobra.Command{
		Use:  "wc",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileName := ""
			var reader io.Reader
			if len(args) > 0 {
				fileName = args[0]
				fileReader, err := os.Open(args[0])
				if err != nil {
					log.Fatalf("openning file: %v", err)
					return
				}
				defer fileReader.Close()
				reader = fileReader
			} else {
				reader = bufio.NewReader(os.Stdin)
			}
			res, err := doTheMagic(reader)
			if err != nil {
				log.Fatalf("error: %v", err)
				return
			}
			printResult(fileName, res, outputCfg)
		},
	}

	cmd.Flags().BoolVarP(&outputCfg.showLines, "lines", "l", false, "")
	cmd.Flags().BoolVarP(&outputCfg.showBytes, "characters", "c", false, "")
	cmd.Flags().BoolVarP(&outputCfg.showWords, "words", "w", false, "")

	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err)
	}
}

func doTheMagic(inputReader io.Reader) (result, error) {
	var (
		lineNum        int64
		bytesNum       int64
		wordsNum       int64
		wordInProgress bool
	)

	bufReader := bufio.NewReader(inputReader)
	for {
		byt, err := bufReader.ReadByte()
		if err != nil {
			if err == io.EOF {
				if wordInProgress {
					wordsNum += 1
				}
				break
			}
			return result{}, err
		}
		bytesNum += 1
		if byt == '\n' {
			lineNum++
		}
		if unicode.IsSpace(rune(byt)) {
			if wordInProgress {
				wordsNum += 1
			}
			wordInProgress = false
		} else {
			wordInProgress = true
		}
	}
	return result{
		linesNumber: lineNum,
		bytesNumber: bytesNum,
		wordsNumber: wordsNum,
	}, nil
}

func printResult(fileName string, res result, cfg outputConfig) {
	outputLine := " "
	if cfg.showLines {
		outputLine += fmt.Sprintf("%7d ", res.linesNumber)
	}
	if cfg.showWords {
		outputLine += fmt.Sprintf("%7d ", res.wordsNumber)
	}
	if cfg.showBytes {
		outputLine += fmt.Sprintf("%7d ", res.bytesNumber)
	}
	outputLine += fileName
	fmt.Println(outputLine)
}
