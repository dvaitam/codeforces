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

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	ns := make([]int, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		ns[i] = n
	}
	out, err := runCandidate(bin, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(strings.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, n := range ns {
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !outScan.Scan() {
				fmt.Printf("missing output for case %d\n", idx+1)
				os.Exit(1)
			}
			v, err := strconv.Atoi(outScan.Text())
			if err != nil {
				fmt.Printf("invalid number in case %d\n", idx+1)
				os.Exit(1)
			}
			arr[i] = v
			if v < 1 || v > 1000 {
				fmt.Printf("case %d: value out of range\n", idx+1)
				os.Exit(1)
			}
		}
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				sum := arr[x] + arr[y]
				for z := 0; z < n; z++ {
					if sum == arr[z] {
						fmt.Printf("case %d failed: %d + %d == %d\n", idx+1, arr[x], arr[y], arr[z])
						os.Exit(1)
					}
				}
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
