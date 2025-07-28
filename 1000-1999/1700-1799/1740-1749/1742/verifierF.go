package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type operation struct {
	d int
	k int
	x string
}

type testCase struct {
	ops []operation
}

func randomString() string {
	length := rand.Intn(3) + 1
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = byte('a' + rand.Intn(3))
	}
	return string(bytes)
}

func generateCases() []testCase {
	rand.Seed(6)
	cases := make([]testCase, 100)
	for i := range cases {
		q := rand.Intn(5) + 1
		ops := make([]operation, q)
		for j := range ops {
			ops[j] = operation{
				d: 1 + rand.Intn(2),
				k: rand.Intn(5) + 1,
				x: randomString(),
			}
		}
		cases[i] = testCase{ops: ops}
	}
	return cases
}

func expected(tc testCase) string {
	cntS := make([]int64, 26)
	cntT := make([]int64, 26)
	cntS[0], cntT[0] = 1, 1
	lenS, lenT := int64(1), int64(1)
	var out strings.Builder
	for _, op := range tc.ops {
		freq := make([]int64, 26)
		for i := 0; i < len(op.x); i++ {
			freq[op.x[i]-'a']++
		}
		if op.d == 1 {
			for i := 0; i < 26; i++ {
				if freq[i] > 0 {
					cntS[i] += int64(op.k) * freq[i]
				}
			}
			lenS += int64(op.k) * int64(len(op.x))
		} else {
			for i := 0; i < 26; i++ {
				if freq[i] > 0 {
					cntT[i] += int64(op.k) * freq[i]
				}
			}
			lenT += int64(op.k) * int64(len(op.x))
		}
		hasTBig := false
		for i := 1; i < 26; i++ {
			if cntT[i] > 0 {
				hasTBig = true
				break
			}
		}
		if hasTBig {
			out.WriteString("YES\n")
			continue
		}
		hasSBig := false
		for i := 1; i < 26; i++ {
			if cntS[i] > 0 {
				hasSBig = true
				break
			}
		}
		if hasSBig {
			out.WriteString("NO\n")
			continue
		}
		if lenS < lenT {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
	}
	return strings.TrimSuffix(out.String(), "\n")
}

func buildIO(cases []testCase) (string, string) {
	var inBuilder strings.Builder
	var outBuilder strings.Builder
	fmt.Fprintf(&inBuilder, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&inBuilder, "%d\n", len(tc.ops))
		for _, op := range tc.ops {
			fmt.Fprintf(&inBuilder, "%d %d %s\n", op.d, op.k, op.x)
		}
		outBuilder.WriteString(expected(tc))
		outBuilder.WriteByte('\n')
	}
	return inBuilder.String(), outBuilder.String()
}

func run(binary, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(binary, ".go") {
		cmd = exec.Command("go", "run", binary)
	} else {
		cmd = exec.Command(binary)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func normalizeTokens(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return strings.Fields(s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	cases := generateCases()
	input, expectedOutput := buildIO(cases)
	actualOutput, err := run(binary, input)
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		os.Exit(1)
	}
	if strings.Join(normalizeTokens(actualOutput), " ") != strings.Join(normalizeTokens(expectedOutput), " ") {
		fmt.Println("Wrong answer")
		fmt.Println("Expected:")
		fmt.Println(expectedOutput)
		fmt.Println("Got:")
		fmt.Println(actualOutput)
		os.Exit(1)
	}
	fmt.Println("All test cases passed!")
}
