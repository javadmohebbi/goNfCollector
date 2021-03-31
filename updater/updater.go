package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

// writeCounter -  Can be used to count the number of bytes which is written
type writeCounter struct {
	URL   string
	Total uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.printProgress()
	return n, nil
}

// printProgress prints the progress of a file write
func (wc writeCounter) printProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading (%v)... %s complete", wc.URL, humanize.Bytes(wc.Total))
}

// DownloadDatabase - Download the lastest version of IP2Location
func DownloadDatabase(dbFileName string, path string, tmpPath string, filetoDownlad string, needUnzip bool) {
	fmt.Println("Download Started:", filetoDownlad)

	fileURL := filetoDownlad
	err := downloadFile(fileURL, tmpPath, path, dbFileName, needUnzip)
	if err != nil {
		fmt.Println("Download error", err)
		os.Exit(1)
	}

	fmt.Println("Download Finished")
}

// downloadFile - Download File and rename it to the one that we want
func downloadFile(url string, tmpPath string, filepath string, dbFileName string, needUnzip bool) error {

	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	out, err := os.Create(tmpPath + dbFileName + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &writeCounter{URL: url}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println()

	if needUnzip {
		files, err := unzipFile(tmpPath+dbFileName+".tmp", filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		err = os.Remove(tmpPath + dbFileName + ".tmp")

		fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	} else {
		return copyFile(tmpPath+dbFileName+".tmp", filepath+"/"+dbFileName)
	}

	return nil

}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		fmt.Println("os stat error: ", err)
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func unzipFile(src string, dest string) ([]string, error) {
	fmt.Println("UZIPING...!")
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
