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

type Test struct {
	rows []string
}

func (t Test) Input() string {
	return strings.Join(t.rows, "\n") + "\n"
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "926D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref failed: %v: %s", err, out)
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genRow(rng *rand.Rand) string {
	seats := []byte("..-..-..")
	for i := 0; i < 8; i++ {
		if seats[i] == '-' {
			continue
		}
		if rng.Intn(2) == 0 {
			seats[i] = '*'
		} else {
			seats[i] = '.'
		}
	}
	return string(seats)
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		rows := make([]string, 6)
		hasVacant := false
		for j := 0; j < 6; j++ {
			row := genRow(rng)
			if strings.Contains(row, ".") {
				hasVacant = true
			}
			rows[j] = row
		}
		if !hasVacant {
			// ensure at least one vacant seat
			rows[0] = strings.Replace(rows[0], "*", ".", 1)
		}
		tests = append(tests, Test{rows})
	}
	tests = append(tests, Test{rows: []string{"..-..-..", "..-..-..", "..-..-..", "..-..-..", "..-..-..", "..-..-.."}})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
