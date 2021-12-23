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
	"sync"
	"time"
)

var waitGroup sync.WaitGroup

func main() {

	var output, searchKey string
	flag.StringVar(&output, "o", "", "Write output to a file")
	flag.Parse()

	args := flag.Args()
	if len(args) != 0 {
		searchKey = args[0]
	}

	switch len(args) {
	case 0:
		fmt.Println("Please provide proper arguments")
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

	inputs := readInputFromConsole()
	matchedLines := findMatches(searchKey, inputs)
	return matchedLines
}

func searchInFile(searchKey, fileName string) []string {

	inputs := readInputFromFile(fileName)

	matchedLines := findMatches(searchKey, inputs)
	return matchedLines
}

func performRecursiveMatching(searchKey, directory string) {

	now := time.Now()
	filePaths := returnFilePaths(directory)
	afterFilePath := time.Since(now)
	waitGroup.Add(len(filePaths))
	results := make(chan map[string][]string)
	findMultipleFileMatches(searchKey, filePaths, results)

	for result := range results {
		for key, value := range result {
			prefix := key + " : "
			printOnConsole(value, prefix)
		}
	}
	fmt.Println(afterFilePath)
	fmt.Println(time.Since(now))
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

func findMultipleFileMatches(searchKey string, filePaths []string, results chan map[string][]string) {

	for _, path := range filePaths {
		go paralleSearchInFile(searchKey, path, results)
	}

	// Launch a goroutine to monitor when all the work is done.
	go func() {
		waitGroup.Wait()
		close(results)
	}()

}

func paralleSearchInFile(searchKey, path string, results chan map[string][]string) {
	matchMap := make(map[string][]string)
	matchMap[path] = searchInFile(searchKey, path)
	results <- matchMap
	waitGroup.Done()
}

func readInputFromConsole() []string {

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
	return inputs
}

func readInputFromFile(fileName string) []string {

	file, err := os.Open(fileName)
	handleErr(err)
	defer file.Close()

	inputs := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inputs = append(inputs, scanner.Text())
	}
	return inputs
}
