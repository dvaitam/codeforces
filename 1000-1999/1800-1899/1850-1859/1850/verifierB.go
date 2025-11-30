package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesB.txt.
const testcasesB = "NyA0OSAyNyAzIDE3IDMzIDMyIDI2IDIwIDMxIDIzIDM4IDE0IDMzIDkKNSA5IDQ5IDcgNDAgMTcgMzUgNDYgMzkgMTAgMjAKMiA0NyA1IDQ0IDIyCjggMzYgNyAyMyAyOCAyMSA0MCA0MSAxNCAzNiAzMSAyOSAzNCAxNyA0IDM2IDEKMiA0NyAyNiA0NiA0MwoxIDQwIDMyCjYgMTYgNDcgMjEgNDYgNSAxMyAzNyAxNSAxNiAxMCAzNSAyOQoyIDYgMjEgMzMgMzIKMiAyMCAzNiAxOSA0NgoyIDM2IDIyIDM1IDE0CjEwIDM2IDM4IDE5IDI5IDYgMzkgMjUgMjEgMzcgMTYgMTkgMTIgMTMgMTIgMyA0MCA0MyAxNyAzMSA1CjIgNDQgNDkgOSAxMAoxIDYgNDUKOSA0NCAyNiA0NiAzNCAxOCAzNCAxNiAxNCA0NCAzOCAyNyAzOCAxOCAyOSAzMiA0MyA0MiA0NQo2IDYgMjEgNDAgOCAzMiAzOCA0MSAyMiAxMyAxNiAyIDQ3CjUgOCA0NiAxNSAyNCAxMSAyMiAyOCA0IDcgMTAKNCAzIDM3IDQxIDM1IDM5IDQ0IDUgMgoyIDQxIDEzIDM5IDM3CjIgMjYgNiAyNCA4CjEgMzkgMgo0IDEyIDQ2IDggMzEgMTQgNDcgNCA0NAoxIDM1IDI4CjEwIDcgMTcgNSAxNSA1IDQyIDIwIDIzIDI4IDEyIDQgMzMgMzAgMyAzOSA3IDQ1IDI2IDEzIDE3CjYgNDcgMzEgMzcgMTEgNDUgNDQgMTQgNTAgNCA0NCAxMSAxMQo2IDM0IDE3IDggMzkgMjkgNDMgMTIgMSAzMSA0NCAyNyAzNwo5IDIwIDQyIDIzIDI1IDQzIDE3IDEwIDM2IDQ1IDEgMzAgNDggNiAyMiA0OCAzIDM1IDE4CjMgMTYgNDkgMzEgMjMgNDAgMTkKNiAzOCA0MSA0MCA5IDQ2IDIwIDI1IDQ4IDI3IDQyIDYgMQoxMCAxMyA0NSAyMiAxMSAxNiAxNSA0MSAyOSAyNSA0NiA0NCAzNyAyNyAzIDI2IDQ1IDM3IDI3IDUwIDQzCjEgMTEgMjkKMiAxNyA0NSAxMSAyOQo5IDMyIDM2IDM5IDQ5IDEgMyAzMiAyMSAyMCAzMCA0IDI3IDEzIDM2IDQxIDYgNDcgOQoxIDI2IDQ0CjcgMjEgMSAxNCAxIDQ2IDQ5IDEgNDQgMzQgNDAgNyAxMyA4IDM5CjQgMjAgMTggNDUgMTIgNyAzMSAyNiA0MQoyIDIgMTggMjkgOAo1IDkgNDIgMzQgNDIgNDIgMjMgOCAxMCAxOCAyCjEgMyAxNAo1IDM2IDIxIDI0IDM3IDMgNDggNDUgMzkgNDIgMzIKOCA0MSAyOCAyNCAzNSAxMiAxNCAyNSAzOCAxOSAxIDkgMTAgMTggMjIgMjIgMjQKMiAyMiA1MCA0MCAzCjEgMTggMTEKMyAzOCAxOSAyNCAyNiAzNiA5CjUgOCAzMSA0NyAxNiA0IDIwIDEyIDM0IDQ3IDUKNSAyNiAyMiAyMCAyNyA3IDcgMzYgMzEgMzEgMjIKNiA4IDMxIDggNDUgMzIgMjggMyAyMCAyMiA0OCA0NCAxMAozIDQxIDM3IDI1IDQxIDYgNQoyIDEzIDQ4IDE1IDQKNyAxIDcgMjYgMzYgMzQgMTkgMjkgMzIgMzggNDYgNDQgMTQgMjggNgo2IDE1IDE3IDM4IDUwIDExIDI4IDEzIDIzIDggNSA0NSAyCjkgMjkgNDkgNDQgMTMgOCAzMiAyNiAxNyAxNCA0MiAzIDE0IDQwIDEwIDcgMTMgMzAgMjUKNiAzNSAxMCA3IDM5IDMyIDEwIDM3IDI2IDQxIDQ0IDI4IDM0CjggNDQgMjEgMzIgMzIgNDEgNDMgMTMgMzUgNDAgMTUgMSAyMiA0NiA0OCAyMSAyMQoxIDM0IDEwCjUgMzkgMTAgMjUgMzggMTkgNDYgNDYgMzEgNSA2CjkgMyA1IDE1IDkgMyAyMCAxIDQ5IDI5IDIyIDExIDEwIDQyIDMwIDI0IDMzIDI1IDM0CjkgMyAzNyA2IDQ0IDM0IDQ5IDM5IDUgNDggMjggNDkgMTQgMTkgMzUgMzkgMjcgMzEgMjUKMTAgMzggMTUgMiA0MyAxIDQ4IDEyIDIwIDMzIDM3IDE3IDIyIDUgMzIgMTcgMjAgNTAgMjcgMjUgMjUKMSAxMSA0MgozIDE2IDE5IDQ3IDIyIDQgMwo4IDI3IDEwIDMyIDM5IDQ2IDYgNDQgNDUgMTAgMjMgMjcgMyA0MCAzMCAyNSAzMAoxIDcgMzEKMyAyIDMgMzkgNDAgOSA0MQo2IDcgNDUgMzYgNDIgMjMgMTMgMjUgNTAgNTAgMzIgOCA0CjEwIDQ1IDMwIDQwIDQxIDIyIDQyIDggNDQgNDYgNDAgMTkgOSAyNSAxOSA0OCA0NCA4IDM0IDEzIDMKNyAyOSAyNCA0OSAxMyAzMCAyMyA0MSA1IDMgMyAzMiAxNyAyIDM0CjEwIDM3IDE0IDE1IDYgNTAgNDEgNTAgMzMgNDUgMzQgMjcgMzMgMjAgOCAxMCAyOCAzNyAyOCA2IDcKNyA1IDcgMjcgNTAgMTAgNDcgMiAyOSAyOCA0NCAyNyAyIDMyIDIxCjUgNiAyMyA1IDggMjMgNDUgMiAyMyAyMyAxMgoxIDE1IDI0CjIgMzkgMTAgMTQgMQo0IDQzIDQ0IDQ3IDggNDggMSAxOSAyNAoxIDM5IDE1CjMgMTIgMzAgOCAzMSAyMyA0Ngo1IDkgMiAxNCAyNCAyMiAzMSAxOSAxOSAzNiA0MQo2IDEyIDM4IDYgNyAzNSAzOCAyMCAxMSAyNSAxMCA5IDE1CjYgMzMgMTYgMTYgNDkgMTIgMTkgMjQgMjcgNDMgMyA5IDM5CjEgMjYgNQoyIDkgMjcgMjAgMzYKNyA0OCAxMCAzOCAyOCAyMCA0MSAyMyA2IDE2IDI5IDQxIDI0IDQxIDM0CjEgMjUgMjcKMSAyNyA0Nwo2IDI5IDE0IDI0IDE5IDMxIDYgMTIgNyAxOCA4IDM2IDM5CjMgNDUgMjkgMjYgMTIgNTAgMjcKNyAxMiAxNiAzMCAyMiAzNCAxMCAyMyAzMCA0MSA0MSA2IDMxIDQ5IDE0CjUgMSA0NSAyOSA0MCAzMCAxIDE0IDIwIDggNTAKNSAzNSAzOSAxMCAyOCA0NiA0OSAzMSA2IDQ0IDMyCjQgMzUgNDkgMjYgMTggNDEgMiA4IDE4CjEgMSAxNwo3IDM0IDM4IDQ2IDI2IDI5IDcgNDggMTcgMjMgMTkgNDkgNDQgMTMgMzkKMiAzIDUgMTcgMjAKOSAyMiA4IDM0IDE2IDQ5IDExIDUgMjcgMTkgMTkgMzQgOSAzNyAzNCA0MSAxNCAzNSA3CjcgNDEgMzUgMjYgNDggNTAgMTggMTkgMjkgMjQgMzcgNDEgOSAxMSA4CjIgMjUgMjYgMzggMzAKMyAzNiA0MyAyMCAyMyA0MSAzMQo3IDE0IDMxIDMyIDQ1IDMzIDIxIDMyIDQyIDQgMjkgMjAgMTAgNDggMzIKMSA0MCAxNAoxIDIzIDMxCjcgMSAzNCA1IDQ0IDYgNDQgNDggNDMgMjYgMSAyNCAzIDggNDAKMSAxOCA0MQo1IDQ3IDE1IDEwIDQ5IDM3IDE5IDEzIDcgMjggMzAKNiAyNSAxMSAyMiAyNyA0MiA0NCAyOCAxMCAyOSA0NiAxMCAzNAo2IDkgMTQgMTIgMjkgMjMgMjUgMjggMzIgMjUgNDcgMTUgMTMKOCAxNCAzOCA0NiA0IDI1IDMgMTUgNDEgNiAxMiAyNCA0IDQ4IDQxIDQ0IDEyCjQgNDAgMjAgNDAgNiA0NiAzMyA0OSAxOQo2IDI3IDMwIDQgNDEgNDUgMzQgNDMgNDIgMzYgNDggMjggMzgKOCAzMiAxNyA0NiAzMSAxNCAyMiAxOCAzIDMgNCAxMSAyMyAxIDE5IDQyIDEKMyA1IDI4IDQ0IDE1IDM5IDI2CjkgMTUgMzAgMTMgMjIgMzkgNyAzOSA2IDIxIDIxIDM1IDMwIDIxIDE3IDIgMzQgMyAxMwo2IDYgMTQgMzQgMjMgMTMgMTMgMTcgNDQgNDcgNDggMjAgMjAKOSAyNSAxNyAzMSAyMyA0NiAxNiAzIDIwIDM2IDUgMSAzMCAzMiA0NyAyOSA0IDI3IDMyCjggMjkgOCA2IDYgMTYgNyA0OSAxMCAyNyAxNCAyOSA0MCA1IDI4IDM2IDQ5CjcgMyAxMiAxNiAzMiAxNSA5IDE4IDIzIDIxIDI4IDcgMzYgMTkgNDAKOSAxMyA0NiAxOSA1MCAyOSAzMyAzOSAzMCAzNSA0MSAxNyAxOCAxNSAyIDggNDAgNDYgNwozIDQ3IDI3IDE2IDE0IDE5IDQ4CjEgNDggMzUKOSAyOCA0IDggMjUgNDIgMTggOCA0OCAzNyAyMyAxNSA0NCA0NiA0NiAzNSA0MyAxOSAxNQo0IDUgMzQgMjAgNDQgMjEgMTUgMjQgNDEKOCAxOSAzOCAxMSA5IDEgMzYgMzMgMjEgMjQgMzggNDEgMiA5IDI2IDEwIDEyCjkgNSA5IDQ5IDE0IDUwIDMyIDM3IDUwIDQ1IDE0IDE2IDQ3IDkgMTUgNDkgMjUgMjMgMzkK"

type testCase struct {
	n     int
	pairs [][2]int
}

// Embedded solver logic from 1850B.go.
func solve(tc testCase) string {
	bestIdx := 1
	bestQual := -1
	for i, p := range tc.pairs {
		a, b := p[0], p[1]
		if a <= 10 && b > bestQual {
			bestQual = b
			bestIdx = i + 1
		}
	}
	return strconv.Itoa(bestIdx)
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesB)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	cases := []testCase{}
	for sc.Scan() {
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %v", err)
		}
		pairs := make([][2]int, n)
		for i := 0; i < n; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing a[%d]", len(cases)+1, i)
			}
			a, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d a[%d]: %v", len(cases)+1, i, err)
			}
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing b[%d]", len(cases)+1, i)
			}
			b, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d b[%d]: %v", len(cases)+1, i, err)
			}
			pairs[i] = [2]int{a, b}
		}
		cases = append(cases, testCase{n: n, pairs: pairs})
	}
	if err := sc.Err(); err != nil {
		return nil, err
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
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
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
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", tc.n)
		for _, p := range tc.pairs {
			fmt.Fprintf(&input, "%d %d\n", p[0], p[1])
		}

		want := solve(tc)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
