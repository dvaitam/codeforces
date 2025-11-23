package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	const numTestCases = 100
	inputBuf := new(bytes.Buffer)
	fmt.Fprintln(inputBuf, numTestCases)

	type TestCase struct {
		N    int
		L, R int64
		A    []int
	}

	testCases := make([]TestCase, numTestCases)

	for i := 0; i < numTestCases; i++ {
		n := rand.Intn(8) + 1 // Small N for brute force
		l := int64(rand.Intn(100) + 1)
		r := l + int64(rand.Intn(100))
		
		testCases[i].N = n
		testCases[i].L = l
		testCases[i].R = r
		testCases[i].A = make([]int, n)
		for j := 0; j < n; j++ {
			testCases[i].A[j] = rand.Intn(200) + 1
		}

		fmt.Fprintf(inputBuf, "%d %d %d\n", n, l, r)
		for j, val := range testCases[i].A {
			if j > 0 {
				fmt.Fprint(inputBuf, " ")
			}
			fmt.Fprint(inputBuf, val)
		}
		fmt.Fprintln(inputBuf)
	}

	var solutionBin string
	var cleanup func()

	if len(os.Args) > 1 {
		solutionBin = os.Args[1]
		if !strings.Contains(solutionBin, "/") {
			solutionBin = "./" + solutionBin
		}
		cleanup = func() {}
	} else {
		solutionSrc := "2157D.go"
		solutionBin = "./solution_d_temp"
		if _, err := os.Stat(solutionSrc); os.IsNotExist(err) {
			fmt.Printf("Error: %s not found.\n", solutionSrc)
			os.Exit(1)
		}
		cmdBuild := exec.Command("go", "build", "-o", solutionBin, solutionSrc)
		cmdBuild.Stderr = os.Stderr
		cmdBuild.Stdout = os.Stdout
		if err := cmdBuild.Run(); err != nil {
			fmt.Printf("Failed to build: %v\n", err)
			os.Exit(1)
		}
		cleanup = func() { os.Remove(solutionBin) }
	}
	defer cleanup()

	cmdRun := exec.Command(solutionBin)
	cmdRun.Stdin = inputBuf
	var outBuf bytes.Buffer
	cmdRun.Stdout = &outBuf
	cmdRun.Stderr = os.Stderr

	if err := cmdRun.Run(); err != nil {
		fmt.Printf("Failed to run solution: %v\n", err)
		os.Exit(1)
	}

	outputLines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
	var validLines []string
	for _, l := range outputLines {
		if strings.TrimSpace(l) != "" {
			validLines = append(validLines, strings.TrimSpace(l))
		}
	}

	if len(validLines) != numTestCases {
		fmt.Printf("Expected %d output lines, got %d\n", numTestCases, len(validLines))
		os.Exit(1)
	}

	passed := 0
	failed := 0
	for i, line := range validLines {
		got, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Printf("Test case %d FAILED: Invalid output '%s'\n", i+1, line)
			failed++
			continue
		}

		expected := solveBrute(testCases[i].N, testCases[i].L, testCases[i].R, testCases[i].A)
		
		if got != expected {
			fmt.Printf("Test case %d FAILED\n", i+1)
			fmt.Printf("Input: N=%d, L=%d, R=%d, A=%v\n", testCases[i].N, testCases[i].L, testCases[i].R, testCases[i].A)
			fmt.Printf("Expected: %d, Got: %d\n", expected, got)
			failed++
			if failed >= 5 {
				break
			}
		} else {
			passed++
		}
	}

	fmt.Printf("Passed %d/%d test cases.\n", passed, numTestCases)
	if failed > 0 {
		os.Exit(1)
	}
}

func solveBrute(n int, l, r int64, a []int) int64 {
	// 3^n brute force
	var maxVal int64 = 0
	
	// We represent choices as an array of -1, 0, 1
	choices := make([]int, n)
	var rec func(idx int)
	rec = func(idx int) {
		if idx == n {
			var sl, sr int64
			for k := 0; k < n; k++ {
				if choices[k] == 1 { // Claim >= a[k]
					sl += l - int64(a[k])
					sr += r - int64(a[k])
				} else if choices[k] == -1 { // Claim <= a[k]
					sl += int64(a[k]) - l
					sr += int64(a[k]) - r
				}
			}
			minS := sl
			if sr < minS {
				minS = sr
			}
			if minS > maxVal {
				maxVal = minS
			}
			return
		}
		
		// Option 0
		choices[idx] = 0
		rec(idx + 1)
		
		// Option 1
		choices[idx] = 1
		rec(idx + 1)
		
		// Option -1
		choices[idx] = -1
		rec(idx + 1)
	}
	
	rec(0)
	return maxVal
}
