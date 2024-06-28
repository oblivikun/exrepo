package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	// Check if the correct number of arguments were provided
	if len(os.Args)!= 3 || os.Args[1]!= "add" && os.Args[1]!= "del" {
		fmt.Println("Usage:./main add <reponame> OR./main del <reponame>")
		return
	}

	action := os.Args[1]
	repoName := os.Args[2]

	switch action {
	case "add":
		repoURL := fmt.Sprintf("https://summer.exherbolinux.org/repositories/%s/index.html", repoName)

		// Fetch HTML content
		resp, err := http.Get(repoURL)
		if err!= nil {
			fmt.Println("Error fetching URL:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err!= nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// Define regex patterns
		syncPattern := regexp.MustCompile(`(?s)<th>sync</th>\s*<td>(.*?)</td>`)
		formatPattern := regexp.MustCompile(`(?s)<th>format</th>\s*<td>(.*?)</td>`)

		// Find matches
		syncMatch := syncPattern.FindStringSubmatch(string(body))
		formatMatch := formatPattern.FindStringSubmatch(string(body))

		// Check if matches were found
		if len(syncMatch) < 2 || len(formatMatch) < 2 {
			fmt.Println("Sync or Format not found.")
			return
		}

		// Extract values
		cleanedSyncValue := cleanHTMLTags(syncMatch[1])
		cleanedFormatValue := cleanHTMLTags(formatMatch[1])

		// Construct the file path
		filePath := fmt.Sprintf("/etc/paludis/repositories/%s.conf", repoName)

		// Create the file if it doesn't exist, or truncate it if it does
		file, err := os.Create(filePath)
		if err!= nil {
			fmt.Printf("Error creating or opening file: %v\n", err)
			return
		}
		defer file.Close()

		// Writing to the file
		_, err = file.WriteString(fmt.Sprintf("format = %s\nlocation = /var/db/paludis/repositories/%s\nsync = %s\n", cleanedFormatValue, repoName, cleanedSyncValue))
		if err!= nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		fmt.Println("Data written to config successfully.")
	case "del":
		// Construct the file path
		filePath := fmt.Sprintf("/etc/paludis/repositories/%s.conf", repoName)

		// Attempt to delete the file
		err := os.Remove(filePath)
		if err!= nil {
			fmt.Printf("Error deleting file: %v\n", err)
		} else {
			fmt.Printf("Deleted file: %s\n", filePath)
		}
	default:
		fmt.Println("Invalid action. Please use 'add' or 'del'.")
	}
}

// Function to remove HTML tags from a string
func cleanHTMLTags(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(input, "")
}
