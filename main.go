package main

import (
		"fmt"
		"bufio"
		"log"
		"os"
		"net"
		"net/http"
		"crypto/tls"
		"strings"
		"time"
		"strconv"
		"flag"
)

func readWordlist(filename string) []string {
	logger := log.New(os.Stderr, "", 0)
    file, err := os.Open(filename)
    if err != nil {
        logger.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
	const maxCapacity int = 20000000
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	words := make([]string, 0)
    for scanner.Scan() {
		words = append(words, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        logger.Fatal("[!] Error opening file")
    }
	return words
}


func writeToFile(filename, content string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(content + "\n")
    if err != nil {
        return err
    }
    return nil
}


func enumS3dns(word string, baseString string, delim string, filename string) {
	cname, err:= net.LookupCNAME(baseString + delim + word + ".s3.amazonaws.com")

	if err != nil {
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
			return
		} else {
			fmt.Println("Network/Timeout issue:", word)
		}
		return
	}

	if !(strings.Contains(cname, "s3-1-w.amazonaws.com.") || (strings.Contains(cname, "s3-w.us-east-1.amazonaws.com."))) {
		bucket := baseString + delim + word
		fmt.Println(bucket)
		if filename == "" {
			writeToFile(filename, bucket)
		}
	} 
}


func enumS3http(word string, baseString string, delim string, filename string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout: 5 * time.Second,
	}

	var url = "https://" + baseString + delim + word + ".s3.amazonaws.com"

	resp, err := client.Get(url)
	
	if err != nil {
		return
	}

	if resp.StatusCode == 200 || resp.StatusCode == 403 {
		bucket := baseString + delim + word
		fmt.Println(bucket)
		if filename == "" {
			writeToFile(filename, bucket)
		}
	}
}


func main() {	
	var wordFile string
	var baseName string
	var delimiter string
	var outputFile string
	var threadCount int
	var httpMode bool

	flag.StringVar(&wordFile, "w", "", "Wordlist for brute-force")
	flag.StringVar(&baseName, "n", "", "Prefix to use (optional)")
	flag.IntVar(&threadCount, "t", 30, "Number of threads to use")
	flag.StringVar(&delimiter, "d", "", "Delimeter to use between words")
	flag.StringVar(&outputFile, "o", "", "File to output results")
	flag.BoolVar(&httpMode, "http", false, "Use HTTP mode or not")
	flag.Parse()

	if wordFile == "" {
		fmt.Println("[!] Some arguments are missing")
		os.Exit(1)
	}

	fmt.Println("[+] Starting S3brute with " + strconv.Itoa(threadCount) + " threads")
	fmt.Println("[*] Using wordlist: " + wordFile)

	if delimiter != "" {
		fmt.Println("[*] Using delimiter: " + delimiter)
	}

	if outputFile != "" {
		fmt.Println("[*] Outputing results to: " + outputFile)
	}

	wordlist := readWordlist(wordFile)

	fmt.Println("[*] Loaded " + strconv.Itoa(len(wordlist)) + " words")

	if httpMode {
		fmt.Println("[*] Discovering S3 with HTTP")
	} else {
		fmt.Println("[*] Discovering S3 with DNS")
	}

	semaphore := make(chan struct{}, threadCount)

	for _, word := range wordlist {
		semaphore <- struct{}{}

		go func(word string) {
			if httpMode {
				enumS3http(word, baseName, delimiter, outputFile)
			} else {
				enumS3dns(word, baseName, delimiter, outputFile)
			}
			<-semaphore
		}(word)
	}

	for i := 0; i < cap(semaphore); i++ {
		semaphore <- struct{}{}
	}
}