package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesC.txt.
const testcasesC = "OTYgNCAzMzkxIDkzOTUKNSAzIDc5MyA4OTUKMSA0IDU0MjQgNjc5NAo0NDQ4Mzc5MyA0IDM4NzggNjg3OAoxNzY2IDQgMTQ1MSA3Njk1CjI0NzkzODY5OTcyIDEgOSA5CjYzNDA0MDc3MDczIDQgMzEyMCA1MjIwCjQ4NzAwMjMzIDEgNSA5CjE1MjIzIDQgNDkxNCA4OTI1Cjg0MzY4MjUyMjE3IDMgMTk5IDQ5OQo1NjA3MTQgMSA0IDcKMzIyNzEgNCAwMzEyIDM4MjYKMTkwMjY4ODI5NSAzIDc3NyA3OTgKMzkgMSA2IDkKMzk5ODU1OTE0IDIgMDcgMTcKNjM3IDMgNDc2IDQ4NgowNjcwNTQ5NiA0IDEwODYgMTk5Ngo5NjM1Mjc1MjIxNTEgNCAwODI2IDA4ODgKNjQ0Mjg4MDE2IDEgNSA4CjcgNCAxNjY3IDg4ODgKMTQ5NjIgMyAwNjEgMDY5Cjc2MTk3IDIgMjcgNTkKODQ2NjE3NyAzIDM3OSA3ODkKNjAwNzE0MzY3IDMgMzg0IDY5NQo4IDEgNyA4CjQyMjQ3NTk0NCAzIDM4OCAzOTkKNDMyMjkzODE2MiAyIDE2IDE5CjE5MiAyIDg3IDk5CjI5ODg1NyA0IDgxODcgOTk4Nwo5NTU3NCA0IDUyMTcgNzc1NwowIDMgNjI1IDc0Nwo5IDIgOTcgOTkKMDU4MjUzNSA0IDQxOTMgODk5OAozOTU0IDIgNjMgNzgKMzA4ODkgNCAxNzI2IDQ4NTgKODM1MTAgMSA1IDYKNTIwODIyNjQ4NjYgMyA1NjIgNjY4CjE3NjE3OTA1NzkgNCA5NDQ3IDk5ODcKMDcyOTU4NzMxMjcyIDMgMzAxIDk2OQo0MjMgMSA0IDkKNzA0NTEgNCAxNjM3IDE5NDgKNDk2Mzc1MTM4NDYyIDQgNzc2NiA4ODk3CjE2IDQgNTE4NiA4Mzk2Cjk5MTQ4ODYgMiAzNSA5Ngo3MzMzMTYgMiA5NSA5NwoyOTMxNTIgNCAwNTcxIDg3ODcKMzc4IDIgNzIgODQKMzEzMyAyIDc0IDg5CjYgMiAxNiA1OQo1NzY3Mzc1NjU3ODUgNCA0OTI1IDg5NjYKODA3NTUgMiA3NSA3NQo3MzA5NTg4NDc1OTUgMSA0IDcKMyAxIDMgOQo4MTQwIDEgMiA1CjcxOTIxMDk0IDMgNDczIDQ4NAowMTA5MTY2MjMyIDIgOTAgOTMKNTk2NjA3MjY0IDQgNDYzMyA0NjU1CjM4MzIxNTcyNzE1IDMgNTc3IDY4Nwo3MDQzOTY0MzEgMiAwMCAxNAozMzcgMyA4OTIgODkyCjQ5MDc5OTkwIDQgMDM5NiAwNDk2CjQ3OSAzIDU0MiA4NjUKNjY4MzU0MzM5NCAyIDI4IDI4CjIwMDYwMDEzOSAxIDMgNAo0NDcxNjU0IDEgNSA5CjE3NjcwNjc0NzA5IDIgODQgODkKNTY5MzU3NTM3IDEgNiA5CjggNCA1NjA3IDc2NTgKMzkxNyAyIDMxIDU3CjYzOTY3IDQgMzExNiA5NDc3CjA5MTI3Mjc5MDk0MCAzIDkwMCA5NzQKOTkwOSAyIDAxIDE2Cjc2ODYyOSAxIDEgNQoyNzI0NDI5NyAyIDg5IDg5CjU1OCAyIDQ3IDQ5CjI1MzQyMjgxMjAgNCA4Mzk4IDk0OTgKNCA0IDQ5NjMgNjk4OAo5MDY0ODc4MCAzIDY3MCA3NzUKMiA0IDI2NjQgODg3NAo4NzY3NiA0IDE1MTAgNTY4MQo3IDMgNjM3IDY5OAoxNTM5OCAyIDgyIDk5CjIzODEgNCA5NDkzIDk3OTQKOTQgMyAyNjYgMjY2CjcwODEwNjU2Mzc0MCAzIDQ2OSA0NjkKNjc1OSAzIDU1NiA2ODcKMTU1NTQ1MTc1MzIgMSAwIDYKODQzIDQgMDMzMCAxODY4CjA0IDQgOTk1MCA5OTc0CjcgMiA0NSA4NQozMjA5NTIxMDYgMyAxMDYgNDA4CjU0MTQ0MCAyIDcyIDk4CjU2OTg4NCAyIDY4IDk4CjM3MjA4IDQgMjMwNyA2Mzk3Cjc1NTM0IDMgMDE0IDYxOAo1NTY0ODM4MDkgMiA2NiA5OAozMjkwOTc5NTQgMyAzMzcgNTU3CjYyMTg5IDIgMDcgNzcKNDA1MjU5MzEzNyAyIDY4IDg5CjA3MzUgNCA2MTMxIDgxNTEKOTY0NzMwOTg0NDYwIDQgNzcyNSA4Nzg1CjMwIDMgNDUzIDY3NAo3OTMzIDQgOTQ5NiA5NTk3Cjk0IDIgNTYgNTcKNzc4IDMgMTA0IDQxNAo5MTA2OTQwNSAxIDEgNgo2NTAxOTM0IDIgNDYgNzgKNjYzMTA1NyAyIDYwIDY4CjA0ODAxMTE0IDEgMSA5CjcgMiAxMiA2NQo3OCAyIDI0IDg3CjIgMiAwNiAxNwo2NDQ5OTM0MzU2IDQgNzkzNiA4OTY3CjI2NjUyMTExNTEwIDQgODA0MyA5OTc5Cjk5ODkzNjI3IDMgNTI4IDgyOAo4NTI3MTQzIDEgNCA2CjQyMjAxMzAzODQgMSA1IDcKNzk0MDkzNDkyOCAxIDUgOAo5Nzk5IDIgMzAgNDQKODYzIDMgNTM3IDg1Nwo="

type testCase struct {
	s string
	m int
	l string
	r string
}

// Embedded solver logic from 1845C.go.
func solve(tc testCase) string {
	s := tc.s
	m := tc.m
	l := tc.l
	r := tc.r

	n := len(s)
	nextPos := make([][10]int, n+1)
	for d := 0; d < 10; d++ {
		nextPos[n][d] = n
	}
	for i := n - 1; i >= 0; i-- {
		for d := 0; d < 10; d++ {
			nextPos[i][d] = nextPos[i+1][d]
		}
		digit := int(s[i] - '0')
		nextPos[i][digit] = i
	}

	cur := map[int]struct{}{-1: {}}
	found := false
	for i := 0; i < m && !found; i++ {
		nextSet := make(map[int]struct{})
		lo := int(l[i] - '0')
		hi := int(r[i] - '0')
		for pos := range cur {
			start := pos + 1
			if start > n {
				start = n
			}
			for d := lo; d <= hi; d++ {
				nxt := nextPos[start][d]
				if nxt == n {
					found = true
					break
				}
				nextSet[nxt] = struct{}{}
			}
			if found {
				break
			}
		}
		cur = nextSet
	}
	if found {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesC)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 4 {
			return nil, fmt.Errorf("bad line %d: %q", idx+1, line)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d m: %v", idx+1, err)
		}
		cases = append(cases, testCase{s: parts[0], m: m, l: parts[2], r: parts[3]})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		input := fmt.Sprintf("1\n%s\n%d\n%s\n%s\n", tc.s, tc.m, tc.l, tc.r)
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
