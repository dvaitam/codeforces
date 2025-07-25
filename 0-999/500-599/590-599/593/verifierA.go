package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(words []string) int {
	type word struct {
		mask uint32
		len  int
	}
	arr := make([]word, 0, len(words))
	for _, w := range words {
		var m uint32
		for _, ch := range w {
			m |= 1 << (ch - 'a')
		}
		if bits.OnesCount32(m) <= 2 {
			arr = append(arr, word{m, len(w)})
		}
	}
	best := 0
	for i := 0; i < 26; i++ {
		for j := i; j < 26; j++ {
			mask := uint32(1<<i | 1<<j)
			sum := 0
			for _, w := range arr {
				if w.mask&^mask == 0 {
					sum += w.len
				}
			}
			if sum > best {
				best = sum
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
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
	cases := make([][]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]string, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			arr[j] = scan.Text()
		}
		cases[i] = arr
	}
	expected := make([]int, t)
	for i, arr := range cases {
		expected[i] = solve(arr)
	}
	for i, arr := range cases {
		input := fmt.Sprintf("%d\n", len(arr))
		for _, w := range arr {
			input += w + "\n"
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, _ := strconv.Atoi(gotStr)
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
