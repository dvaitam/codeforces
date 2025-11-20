package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// The original solveB function, kept for reference in error messages
func solveB(n, k int) string {
	var sb strings.Builder
	sb.WriteString("YES\n")
	blank := strings.Repeat(".", n)
	if k%2 == 0 {
	half := k / 2
		row := "." + strings.Repeat("#", half) + strings.Repeat(".", n-1-half)
		sb.WriteString(blank + "\n")
		sb.WriteString(row + "\n")
		sb.WriteString(row + "\n")
		sb.WriteString(blank + "\n")
	} else if k >= n-2 {
		sb.WriteString(blank + "\n")
		row2 := "." + strings.Repeat("#", n-2) + "."
		rem := k - (n - 2)
		half := rem / 2
		midDots := (n - 2) - rem
		row3 := "." + strings.Repeat("#", half) + strings.Repeat(".", midDots) + strings.Repeat("#", half) + "."
		sb.WriteString(row2 + "\n")
		sb.WriteString(row3 + "\n")
		sb.WriteString(blank + "\n")
	} else {
		sb.WriteString(blank + "\n")
		left := (n - k) / 2
		row2 := strings.Repeat(".", left) + strings.Repeat("#", k) + strings.Repeat(".", left)
		sb.WriteString(row2 + "\n")
		sb.WriteString(blank + "\n")
		sb.WriteString(blank + "\n")
	}
	return sb.String()
}

func solveBFromInput(input string) string {
    var n, k int
    fmt.Sscanf(input, "%d %d", &n, &k)
    return solveB(n,k)
}

type Test struct {
	input string
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(10)*2 + 3 // odd between 3 and 21
		k := rand.Intn(2*(n-2) + 1)
		if n == 13 && k == 17 { // Make sure the failing test case is included
			tests = append(tests, Test{input: "13 17\n"})
			continue
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		tests = append(tests, Test{input: input})
	}
	// Add the failing test case if not added
	found := false
	for _, t := range tests {
		if t.input == "13 17\n" {
			found = true
			break
		}
	}
	if !found {
		tests = append(tests, Test{input: "13 17\n"})
	}

	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.String(), nil
}

func checkSolution(input string, output string) error {
    var n, k int
    fmt.Sscanf(input, "%d %d", &n, &k)

    lines := strings.Split(strings.TrimSpace(output), "\n")
    if len(lines) == 0 || lines[0] != "YES" {
        return fmt.Errorf("output should start with YES")
    }

    if len(lines) != 5 { // YES + 4 rows
        return fmt.Errorf("expected 5 lines in output, got %d", len(lines))
    }

    grid := lines[1:]
    
    hotelCount := 0
    for r := 0; r < 4; r++ {
        if len(grid[r]) != n {
            return fmt.Errorf("row %d has wrong length %d, expected %d", r+1, len(grid[r]), n)
        }
        for c := 0; c < n; c++ {
            if grid[r][c] == '#' {
                hotelCount++
                if r == 0 || r == 3 {
                    return fmt.Errorf("hotel at forbidden row %d", r+1)
                }
                if grid[r][n-1-c] != '#' {
                    return fmt.Errorf("asymmetry at row %d, col %d", r+1, c+1)
                }
            }
        }
    }

    if hotelCount != k {
        return fmt.Errorf("wrong number of hotels, expected %d, got %d", k, hotelCount)
    }

    return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(binary, t.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		
        err = checkSolution(t.input, got)
		if err != nil {
			expected := solveBFromInput(t.input)
			fmt.Printf("Test %d failed. Input: %q\nError: %v\n\nExpected (one of possible solutions):\n%s\nGot:\n%s\n", i+1, t.input, err, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
