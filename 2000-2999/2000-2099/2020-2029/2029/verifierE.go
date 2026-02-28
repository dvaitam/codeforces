package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const maxVal = 400000

var spf [maxVal + 1]int
var isPrimeArr [maxVal + 1]bool

func sieve() {
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			spf[i] = i
			isPrimeArr[i] = true
			if i <= maxVal/i {
				for j := i * i; j <= maxVal; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			spf[i] = i
		}
	}
}

// refSolve is the reference solution for one test case.
func refSolve(a []int) int {
	primeCount := 0
	primeVal := -1
	minVal := maxVal + 1
	for _, v := range a {
		if v < minVal {
			minVal = v
		}
		if isPrimeArr[v] {
			primeCount++
			primeVal = v
		}
	}
	if primeCount == 0 {
		return 2
	}
	if primeCount > 1 {
		return -1
	}
	p := primeVal
	if p != minVal {
		return -1
	}
	twice := 2 * p
	for _, y := range a {
		if y == p {
			continue
		}
		if isPrimeArr[y] {
			return -1
		}
		if y < twice {
			return -1
		}
		if y == twice {
			continue
		}
		if y == twice+1 {
			return -1
		}
		s := spf[y]
		if y-s < twice {
			return -1
		}
	}
	return p
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func buildTests(rng *rand.Rand) ([][]int, string) {
	var cases [][]int
	var sb strings.Builder

	add := func(a []int) {
		cases = append(cases, a)
	}

	// Examples from problem statement.
	add([]int{8, 9, 10})
	add([]int{2, 3, 4, 5})
	add([]int{7, 15})
	add([]int{453, 6, 8, 25, 100000})

	// Edge cases.
	add([]int{2})
	add([]int{3})
	add([]int{4})
	add([]int{400000})
	add([]int{2, 4})
	add([]int{3, 6})
	add([]int{3, 9})
	add([]int{5, 10})
	add([]int{5, 11})
	add([]int{3, 5})       // two odd primes
	add([]int{7, 14, 21})
	add([]int{2, 3})       // prime + even

	// Random small cases.
	for len(cases) < 80 {
		n := rng.Intn(5) + 1
		used := make(map[int]bool)
		a := make([]int, 0, n)
		for len(a) < n {
			v := rng.Intn(30) + 2
			if !used[v] {
				used[v] = true
				a = append(a, v)
			}
		}
		add(a)
	}

	// Random medium cases.
	for len(cases) < 100 {
		n := rng.Intn(20) + 2
		used := make(map[int]bool)
		a := make([]int, 0, n)
		for len(a) < n {
			v := rng.Intn(1000) + 2
			if !used[v] {
				used[v] = true
				a = append(a, v)
			}
		}
		add(a)
	}

	// Format input.
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, a := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", len(a)))
		parts := make([]string, len(a))
		for i, v := range a {
			parts[i] = fmt.Sprintf("%d", v)
		}
		sb.WriteString(strings.Join(parts, " "))
		sb.WriteByte('\n')
	}
	return cases, sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	sieve()

	rng := rand.New(rand.NewSource(42))
	cases, input := buildTests(rng)
	inputData := []byte(input)

	// Run reference solution.
	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "2029E.go")

	refOut, err := runProgram(refPath, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refFields := strings.Fields(refOut)
	if len(refFields) < len(cases) {
		fmt.Fprintf(os.Stderr, "reference output too short: expected %d, got %d\n", len(cases), len(refFields))
		os.Exit(1)
	}

	// Run candidate.
	candOut, err := runProgram(target, inputData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candFields := strings.Fields(candOut)
	if len(candFields) < len(cases) {
		fmt.Fprintf(os.Stderr, "candidate output too short: expected %d, got %d\n", len(cases), len(candFields))
		os.Exit(1)
	}

	for i := range cases {
		if refFields[i] != candFields[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s got %s\n", i+1, refFields[i], candFields[i])
			// Print the failing input.
			parts := make([]string, len(cases[i]))
			for j, v := range cases[i] {
				parts[j] = fmt.Sprintf("%d", v)
			}
			fmt.Fprintf(os.Stderr, "input: n=%d a=[%s]\n", len(cases[i]), strings.Join(parts, " "))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(cases))
}
