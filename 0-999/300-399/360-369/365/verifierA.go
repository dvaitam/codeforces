package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func countKGood(n, k int, nums []string) int {
	ans := 0
	for _, s := range nums {
		found := make([]bool, k+1)
		need := k + 1
		for i := 0; i < len(s); i++ {
			d := int(s[i] - '0')
			if d <= k && !found[d] {
				found[d] = true
				need--
				if need == 0 {
					break
				}
			}
		}
		if need == 0 {
			ans++
		}
	}
	return ans
}

func runCase(bin string, n, k int, nums []string) error {
	var in bytes.Buffer
	fmt.Fprintf(&in, "%d %d\n", n, k)
	for _, s := range nums {
		fmt.Fprintln(&in, s)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", countKGood(n, k, nums))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "missing n on case %d\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "missing k on case %d\n", i+1)
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		nums := make([]string, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "missing value on case %d line %d\n", i+1, j+1)
				os.Exit(1)
			}
			nums[j] = scan.Text()
		}
		if err := runCase(bin, n, k, nums); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
