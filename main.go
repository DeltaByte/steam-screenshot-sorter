package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func cleanName(name string) string {
	invalidChars := []string{"/", "\\", ":", "<", ">", "\"", "|", "?", "*"}
	cleaned := name
	for _, char := range invalidChars {
		cleaned = strings.ReplaceAll(cleaned, char, "_")
	}
	return cleaned
}

func renameDirs(appList *AppList, files []fs.DirEntry) error {
	idRegex := regexp.MustCompile("^\\[(?P<app>[0-9]+)\\] (.+)$")
	regexIndex := idRegex.SubexpIndex("app")

	for _, file := range files {
		// ignore actual files
		if !file.IsDir() {
			continue
		}

		// ignore non-matching dirs
		if !idRegex.MatchString(file.Name()) {
			continue
		}

		// extract appId from dir name
		rawAppId := idRegex.FindStringSubmatch(file.Name())[regexIndex]
		appId, err := strconv.Atoi(rawAppId)
		if err != nil {
			return err
		}

		// rename dir if needed
		appName := appList.Find(appId)
		dirName := fmt.Sprintf("[%d] %s", appId, cleanName(appName))
		if file.Name() != dirName {
			if err := os.Rename(file.Name(), dirName); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	appList, err := LoadAppList()
	if err != nil {
		fmt.Println("Failed to load applist")
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	// scan files
	files, err := os.ReadDir("./")
	if err != nil {
		fmt.Printf("Failed to scan for files")
		os.Exit(1)
	}
	fileCount := len(files)
	fmt.Printf("Found %d files\n\n\n", fileCount)

	// update dir names to account for games having changed their names (or having previously been UNKNOWN)
	renameDirs(appList, files)

	// match and move files
	imageRegex := regexp.MustCompile(`^\w+.(jpe?g|png|gif|bmp)$`)
	for i, file := range files {
		counter := fmt.Sprintf("(%d/%d)", i+1, fileCount)

		// ignore dir/non-image files
		if file.IsDir() || !imageRegex.MatchString(file.Name()) {
			fmt.Printf("%s File is directory or non-image, skipping.\n	[%s]\n\n", counter, file.Name())
			continue
		}

		// match app ID
		rawAppId := strings.Split(file.Name(), "_")[0]
		appId, err := strconv.Atoi(rawAppId)
		if err != nil {
			fmt.Printf("%s Failed to parse appId from filename.\n	[%s]\n\n", counter, file.Name())
		}
		appName := appList.Find(appId)
		fmt.Printf("%s Matched file as \"%s\"\n	[%s]\n\n", counter, appName, file.Name())

		// check for app dir, create if needed
		dirName := fmt.Sprintf("[%d] %s", appId, cleanName(appName))
		_, err = os.Stat(dirName)
		if os.IsNotExist(err) {
			if err := os.Mkdir(dirName, 0777); err != nil {
				fmt.Printf("%s Failed to create directory: %s\n\n", counter, dirName)
				continue
			}
		}

		// move file
		newName := path.Join(dirName, file.Name())
		if err := os.Rename(file.Name(), newName); err != nil {
			fmt.Printf("%s Failed to move file into \"%s\".\n	[%s]\n\n", counter, dirName, file.Name())
			panic(err)
		}
	}

	fmt.Println("Complete, press enter key to exit...")
	fmt.Scanln()
}
