package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type JdkData struct {
	Version         string
	DownloadPath    string
	UnzipFolderName string
}

var (
	java21 = JdkData{
		Version:         "21",
		DownloadPath:    "https://corretto.aws/downloads/latest/amazon-corretto-21-aarch64-macos-jdk.tar.gz",
		UnzipFolderName: "amazon-corretto-21.jdk",
	}

	java17 = JdkData{
		Version:         "17",
		DownloadPath:    "https://corretto.aws/downloads/latest/amazon-corretto-17-aarch64-macos-jdk.tar.gz",
		UnzipFolderName: "amazon-corretto-17.jdk",
	}

	java11 = JdkData{
		Version:         "11",
		DownloadPath:    "https://corretto.aws/downloads/latest/amazon-corretto-11-aarch64-macos-jdk.tar.gz",
		UnzipFolderName: "amazon-corretto-11.jdk",
	}

	java8 = JdkData{
		Version:         "8",
		DownloadPath:    "https://corretto.aws/downloads/latest/amazon-corretto-8-aarch64-macos-jdk.tar.gz",
		UnzipFolderName: "amazon-corretto-8.jdk",
	}
)

var jdkDataList = []JdkData{java21, java17, java11, java8}
var jdkVersionMap = map[string]JdkData{
	"21": java21,
	"17": java17,
	"11": java11,
	"8":  java8,
}

func main() {
	jdkData, err := ShowJdkList(jdkDataList)
	if err != nil {
		panic(err)
	}

	if !CheckIfJdkExists(jdkData) {
		err = DownloadJdk(jdkData)
		if err != nil {
			panic(err)
		}
	}
	err = SetSymLinkToBinary(jdkData)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println()
	color.Blue("Please add the following to your bash profile, ignore if you have already done it")
	color.Green("# Set java home to sym link which will be swapped out by the jdk manager")
	color.Green("export JAVA_HOME=jdk")
}

func ShowJdkList(jdkDataList []JdkData) (JdkData, error) {
	var versions []string
	for _, data := range jdkDataList {
		versions = append(versions, data.Version)
	}

	prompt := promptui.Select{
		Label: "Select Version",
		Items: versions,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return JdkData{}, err
	}

	return jdkVersionMap[result], nil
}

func DownloadFile(url, outputFile string) error {
	// Create or truncate the output file
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Perform the HTTP request
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}

	// Copy the contents of the response body to the output file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func UnzipFile(zipFile, outputDir string) error {
	// Open the tar.gz file
	file, err := os.Open(zipFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Extract files from the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Reached the end of the archive
		}
		if err != nil {
			return err
		}

		// Create the full path for the extracted file
		targetPath := filepath.Join(outputDir, header.Name)

		// Create directories if needed
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(targetPath, 0755)
			if err != nil {
				return err
			}
			continue
		}

		// Create the file
		targetFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		// Copy the file content from the tar archive to the new file
		_, err = io.Copy(targetFile, tarReader)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckIfJdkExists(data JdkData) bool {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error while fetching wd:", err)
		return false
	}
	outputDir := filepath.Join("jdk-dump", data.Version)
	javaBinaryPath := filepath.Join(wd, outputDir, data.UnzipFolderName, "Contents/Home/bin/java")
	if _, err := os.Stat(javaBinaryPath); os.IsNotExist(err) {
		return false
	} else {
		color.Yellow("JDK is already present at %s", javaBinaryPath)
		return true
	}
}

func DownloadJdk(data JdkData) error {
	fmt.Println("Downloading Java " + data.Version)
	outputFile := "jdk-temp.tar.gz"
	err := DownloadFile(data.DownloadPath,
		outputFile)

	if err != nil {
		fmt.Println("Error downloading file:", err)
		return err
	}

	fmt.Println("File downloaded successfully:", outputFile)
	outputDir := filepath.Join("jdk-dump", data.Version)
	err = UnzipFile(outputFile, outputDir)
	if err != nil {
		fmt.Println("Error unzipping file:", err)
		return err
	}

	fmt.Println("File unzipped successfully.")

	// Delete the tar.gz file
	err = os.Remove(outputFile)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return err
	}

	fmt.Println("File deleted:", outputFile)
	return nil
}

func SetSymLinkToBinary(data JdkData) error {
	fmt.Println()
	fmt.Println()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error while fetching wd:", err)
		return err
	}

	outputDir := filepath.Join("jdk-dump", data.Version)
	javaBinaryPath := filepath.Join(wd, outputDir, data.UnzipFolderName, "Contents/Home/bin/java")
	// Set execute permission on the Java binary
	err = os.Chmod(javaBinaryPath, 0755)
	if err != nil {
		fmt.Println("Error setting execute permission:", err)
		return err
	}

	color.Blue("Execute permission set successfully on:", javaBinaryPath)

	var targetPath = filepath.Join(wd, outputDir, data.UnzipFolderName, "Contents/Home")

	if err := os.Remove("jdk"); err != nil && !os.IsNotExist(err) {
		fmt.Println("Error removing existing symlink:", err)
		return err
	}

	err = os.Symlink(targetPath, "jdk")
	if err != nil {
		fmt.Println("Error creating symbolic link:", err)
		return err
	}

	color.Blue("Symbolic link created: %s -> %s\n", "jdk", targetPath)
	return nil
}
