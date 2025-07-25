package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

var mapping = map[rune]string{
	'0': "",
	'1': "",
	'2': "2",
	'3': "3",
	'4': "322",
	'5': "5",
	'6': "53",
	'7': "7",
	'8': "7222",
	'9': "7332",
}

func solveDigits(s string) string {
	var res []byte
	for _, ch := range s {
		if rep, ok := mapping[ch]; ok {
			for i := 0; i < len(rep); i++ {
				res = append(res, rep[i])
			}
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i] > res[j] })
	return string(res)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		digits := scan.Text()
		expected[i] = solveDigits(digits)
		_ = n // n not needed
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		if outScan.Text() != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], outScan.Text())
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
