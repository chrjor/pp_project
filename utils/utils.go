package utils

import (
	"bufio"
	"os"
)

// Function adapted from HW1 Problem 2
func ReadFile(filePath string) []string {

	inFile, _ := os.Open(filePath)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var input []string

	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, line)
	}
	return input
}
