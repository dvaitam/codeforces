package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveD(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)
	times := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &times[i])
	}
	sortSlice(times)
	var sum int64
	count := 0
	for _, t := range times {
		if sum <= t {
			count++
			sum += t
		}
	}
	return strconv.Itoa(count)
}

func sortSlice(a []int64) {
	if len(a) < 2 {
		return
	}
	quickSort(a, 0, len(a)-1)
}

func quickSort(a []int64, l, r int) {
	if l >= r {
		return
	}
	p := a[(l+r)/2]
	i, j := l, r
	for i <= j {
		for a[i] < p {
			i++
		}
		for a[j] > p {
			j--
		}
		if i <= j {
			a[i], a[j] = a[j], a[i]
			i++
			j--
		}
	}
	if l < j {
		quickSort(a, l, j)
	}
	if i < r {
		quickSort(a, i, r)
	}
}

func genTests() []string {
	rand.Seed(4)
	tests := make([]string, 0, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintln(n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d ", rand.Int63n(20)+1))
		}
		tests = append(tests, strings.TrimSpace(sb.String()))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierD <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		expected := solveD(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
