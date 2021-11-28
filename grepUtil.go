package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	var output string
	flag.StringVar(&output, "o", "", "Write output to a file")
	flag.Parse()

	fmt.Println(output)
	args := flag.Args()
	searchKey := args[0]

	switch len(args) {
	//Search for a string from standard input
	case 1:
		matchedLines := searchFromConsole(searchKey)
		printOnConsole(matchedLines, "")
	case 2:
		file, err := os.Stat(args[1])
		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
		}
		//Search for a string recursively in a directory
		if file.IsDir() {
			searchKey := args[0]
			directory := args[1]
			performRecursiveMatching(searchKey, directory)
		} else { //Search for a string in a file
			fileName := args[1]
			matchedLines := searchInFile(searchKey, fileName)
			printOnConsole(matchedLines, "")
		}
	//Write the output to a file
	case 4:
		if output != "" {
			inputFileName := args[1]
			outputFileName := output
			matchedLines := searchInFile(searchKey, inputFileName)
			writeToFile(matchedLines, outputFileName)

		} else {
			fmt.Println("Wrong parameter provided")
		}

	default:
		fmt.Println("Wrong parameters provided")

	}

}

func searchFromConsole(searchKey string) []string {

	inputs := []string{}
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		inputs = append(inputs, input)
		if err != io.EOF {
			handleErr(err)
		} else {
			inputs = inputs[:len(inputs)-1]
			break
		}
	}
	matchedLines := findMatches(searchKey, inputs)
	return matchedLines
}

func searchInFile(searchKey, fileName string) []string {

	file, err := os.Open(fileName)
	handleErr(err)
	defer file.Close()

	inputs := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputs = append(inputs, scanner.Text())
	}

	matchedLines := findMatches(searchKey, inputs)
	return matchedLines
}

func performRecursiveMatching(searchKey, directory string) {

	filePaths := returnFilePaths(directory)

	matchesMap := findMultipleFileMatches(searchKey, filePaths)

	for key, value := range matchesMap {
		prefix := key + " : "
		printOnConsole(value, prefix)
	}
}

func findMatches(searchKey string, inputs []string) []string {

	matchedLines := []string{}
	searchPattern, err := regexp.Compile("(?i)" + searchKey)
	handleErr(err)

	for _, value := range inputs {
		if searchPattern.MatchString(value) {
			matchedLines = append(matchedLines, value)
		}
	}
	return matchedLines
}

func printOnConsole(matchedLines []string, prefix string) {
	for _, value := range matchedLines {
		fmt.Println(prefix + value)
	}
}

func writeToFile(matchedLines []string, fileName string) {
	file, err := os.Create(fileName)
	handleErr(err)
	for _, value := range matchedLines {
		_, error := io.WriteString(file, value+"\n")
		handleErr(error)
	}

}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func returnFilePaths(directory string) []string {
	filePaths := []string{}
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filePaths = append(filePaths, path)
			}
			return nil
		})
	handleErr(err)
	return filePaths
}

func findMultipleFileMatches(searchKey string, filePaths []string) map[string][]string {

	matchMap := make(map[string][]string)
	for _, path := range filePaths {
		matchMap[path] = searchInFile(searchKey, path)
	}
	return matchMap
}
