package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Point struct {
	x, y int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const numTestCases = 100
	inputBuf := new(bytes.Buffer)
	fmt.Fprintln(inputBuf, numTestCases)

	type TestCase struct {
		N    int
		X, Y int
		S    string
	}

	testCases := make([]TestCase, numTestCases)

	for i := 0; i < numTestCases; i++ {
		n := rand.Intn(10) + 1 // Small N for simulation
		testCases[i].N = n
		// Random X, Y in range -25 to 25 to be safe with expansion
		testCases[i].X = rand.Intn(50) - 25
		testCases[i].Y = rand.Intn(50) - 25
		
		var sb strings.Builder
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('4')
			} else {
				sb.WriteByte('8')
			}
		}
		testCases[i].S = sb.String()

		fmt.Fprintf(inputBuf, "%d %d %d\n", testCases[i].N, testCases[i].X, testCases[i].Y)
		fmt.Fprintln(inputBuf, testCases[i].S)
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
		solutionSrc := "2157B.go"
		solutionBin = "./solution_b_temp"
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
		got := strings.TrimSpace(line)
		expected := solveSim(testCases[i].N, testCases[i].X, testCases[i].Y, testCases[i].S)
		
		// Case insensitive check
		if strings.ToUpper(got) != expected {
			fmt.Printf("Test case %d FAILED\n", i+1)
			fmt.Printf("Input: %d %d %d, S=%s\n", testCases[i].N, testCases[i].X, testCases[i].Y, testCases[i].S)
			fmt.Printf("Expected: %s, Got: %s\n", expected, got)
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

func solveSim(n int, tx, ty int, s string) string {
	reached := make(map[Point]bool)
	reached[Point{0, 0}] = true

	for _, op := range s {
		next := make(map[Point]bool)
		// Copy existing (cumulative)
		for p := range reached {
			next[p] = true
		}
		
		for p := range reached {
			// Neighbors
			dx := []int{0, 0, 1, -1}
			dy := []int{1, -1, 0, 0}
			for k := 0; k < 4; k++ {
				next[Point{p.x + dx[k], p.y + dy[k]}] = true
			}
			
			if op == '8' {
				// Diagonals
				dx8 := []int{1, 1, -1, -1}
				dy8 := []int{1, -1, 1, -1}
				for k := 0; k < 4; k++ {
					next[Point{p.x + dx8[k], p.y + dy8[k]}] = true
				}
			}
		}
		reached = next
	}

	if reached[Point{tx, ty}] {
		return "YES"
	}
	return "NO"
}
