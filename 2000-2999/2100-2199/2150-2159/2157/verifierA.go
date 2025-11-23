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
		N int
		A []int
	}

	testCases := make([]TestCase, numTestCases)

	for i := 0; i < numTestCases; i++ {
		n := rand.Intn(100) + 1 // 1 to 100
		testCases[i].N = n
		testCases[i].A = make([]int, n)
		for j := 0; j < n; j++ {
			testCases[i].A[j] = rand.Intn(n + 1) // 0 to n
		}

		fmt.Fprintln(inputBuf, n)
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
		// Compile the solution
		solutionSrc := "2157A.go"
		solutionBin = "./solution_a_temp"

		// Check if source exists
		if _, err := os.Stat(solutionSrc); os.IsNotExist(err) {
			fmt.Printf("Error: %s not found in current directory.\n", solutionSrc)
			os.Exit(1)
		}

		cmdBuild := exec.Command("go", "build", "-o", solutionBin, solutionSrc)
		cmdBuild.Stderr = os.Stderr
		cmdBuild.Stdout = os.Stdout
		if err := cmdBuild.Run(); err != nil {
			fmt.Printf("Failed to build solution: %v\n", err)
			os.Exit(1)
		}
		cleanup = func() { os.Remove(solutionBin) }
	}
	defer cleanup()

	// Run the solution
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
	// Filter out empty lines if any
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
		got, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Test case %d: Invalid output '%s': %v\n", i+1, line, err)
			failed++
			continue
		}

		expected := solve(testCases[i].N, testCases[i].A)
		if got != expected {
			fmt.Printf("Test case %d FAILED\n", i+1)
			fmt.Printf("Input N: %d\n", testCases[i].N)
			fmt.Printf("Input A: %v\n", testCases[i].A)
			fmt.Printf("Expected: %d, Got: %d\n", expected, got)
			failed++
			if failed >= 5 {
				fmt.Println("Too many failures, stopping.")
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

func solve(n int, a []int) int {
	cnt := make(map[int]int)
	for _, x := range a {
		cnt[x]++
	}

	kept := 0
	for x, c := range cnt {
		if x == 0 {
			continue
		}
		if c >= x {
			kept += x
		}
	}
	return n - kept
}
