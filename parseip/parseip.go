package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ipAddressGrammer *regexp.Regexp

func initParser() {
	partIP := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	grammar := partIP + "\\." + partIP + "\\." + partIP + "\\." + partIP
	ipAddressGrammer = regexp.MustCompile(grammar)
}

func findIPAddress(input string) string {
	if ipAddressGrammer == nil {
		initParser()
	}
	return ipAddressGrammer.FindString(input)
}

func parseFile(filename string, ipAddresses map[string]int) {

	f, err := os.Open(filename)

	if err != nil {
		fmt.Printf("error opening file %s\n", err)
		os.Exit(-1)
	}

	defer f.Close()

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')

		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading file %s", err)
			}
			break
		}

		ip := findIPAddress(line)

		if net.ParseIP(ip).To4() == nil {
			continue
		}

		_, ok := ipAddresses[ip]

		if ok {
			ipAddresses[ip] = ipAddresses[ip] + 1
		} else {
			ipAddresses[ip] = 1
		}
	}

}
func printAddresses(addresses map[string]int, minCount int) {

	for key, v := range addresses {
		if v < minCount {
			continue
		} else if v > 10 {
			names, _ := net.LookupAddr(key)
			fmt.Printf("%d %s\n", v, strings.Join(append(names, key), "|"))
		} else {
			fmt.Printf("%d %s\n", v, key)
		}
	}
}
func main() {

	minusC := flag.Int("c", 10, "Min # occurrences before reporting address ")
	flag.Parse()
	flags := flag.Args()

	if len(flags) == 0 {
		fmt.Printf("Usage: %s [OPTIONS] logFile...\n OPTIONS\n -c N\t only display ip if occurs N times", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	ipAddresses := make(map[string]int)

	for _, filename := range flag.Args() {
		parseFile(filename, ipAddresses)
	}

	printAddresses(ipAddresses, *minusC)

}
