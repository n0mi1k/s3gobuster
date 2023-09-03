# s3buster
A fast S3 bucket suffix brute force tool to identify existing buckets. By default, DNS enumeration is used, which does not hit AWS's infrastructure. HTTP mode is also supported with the -http flag, but it is slower and noisier.

***NOTE**: This tool does not identify buckets in the us-east-1 region*

# Installation
To install, just run the below command.
```
go install github.com/n0mi1k/s3buster@latest
```

## Options
```
Usage:
  s3buster -u [URL] [Other Flags]

Flags:
  -n        Prefix to use (OPTIONAL)
  -w        Set the wordlist e.g shubs-subdomains.txt [REQUIRED]
  -t        Number of concurrent goroutines [Default=30]
  -d        Delimeter to use between words (e.g . or -) (OPTIONAL)
  -o        Filename to output results (OPTIONAL)
  -http     Enable HTTP bruteforce mode
```
**How to Use s3buster Properly:**  
- Use s3wordlist.txt for a list of good suffixes, set a delimeter if you require.  
For e.g, if delimiter "-" is set, here is the test format \<Prefix\>**-**\<Suffix\>  
So if prefix "testbucket" is used and word is "prod", this is the test case "testbucket-prod"  
- If no prefix is set, simple wordlist bruteforce willl be used
- DNS brute force is fast and stealthier but may affect your connection
- HTTP brute force works but is noisy and much slower 

## Demo Run


## Disclaimer
This tool is for educational and testing purposes only. Do not use it to exploit the vulnerability on any system that you do not own or have permission to test. The authors of this script are not responsible for any misuse or damage caused by its use.