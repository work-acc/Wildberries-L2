package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/net/html"
)

// Реализовать утилиту wget с возможностью скачивать сайты целиком.

type wgetParameters struct {
	URL       string
	OutputDir string
}

func main() {
	parameters := parceCmdArgs()

	err := download(parameters.URL, parameters.OutputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Download error: %v\n", err)
		os.Exit(1)
	}
}

func parceCmdArgs() wgetParameters {
	parameters := wgetParameters{}

	flag.StringVar(&parameters.URL, "url", "", "The URL of the download site")
	flag.StringVar(&parameters.OutputDir, "output", ".", "The directory for saving files")

	flag.Parse()

	return parameters
}

func download(urlStr string, outputDir string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	err = downloadPage(urlStr, parsedURL, outputDir)
	if err != nil {
		return err
	}

	fmt.Printf("The site is uploaded to the directory: %s\n", outputDir)
	return nil
}

func downloadPage(urlStr string, baseURL *url.URL, outputDir string) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the page could not be downloaded, the status: %d", resp.StatusCode)
	}

	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					linkURL, err := baseURL.Parse(attr.Val)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error in URL analysis: %v\n", err)
						continue
					}

					err = downloadResource(linkURL, outputDir)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error downloading the resource: %v\n", err)
					}
				}
			}
		}
	}
}

func downloadResource(url *url.URL, outputDir string) error {
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download the resource %s, status: %d", url.String(), resp.StatusCode)
	}

	filePath := path.Join(outputDir, url.Host, url.Path)
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("A resource has been downloaded: %s\n", filePath)
	return nil
}
