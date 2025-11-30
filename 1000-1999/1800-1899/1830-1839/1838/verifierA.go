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

// Base64-encoded contents of testcasesA.txt.
const testcasesA = "MTAwCjkgNDcgMyAtNDUgLTE3IDE1IDEyIDEgNTAgLTEyCjEwIC01IDI0IC0yMyAxNCAtMzMgLTE0IC0zMyA0NiAtMzggMjkKNyAxOCA0MCAyNyAtMzIgLTExIC0zOCA0Mwo0IDM3IC04IDEwIDIxCjQgLTUgNSAtMTAgMjgKNiAyMCAxMSA2IDE2IC0xNyAtNDMKMyAtMzkgNDIgMQozIDI4IDEzIC04CjYgNDMgLTkgNDAgLTQyIC0yNiAyMgo2IC0yMCAtMzIgMTkgNyAtMzkgLTQwCjggMTUgMTIgLTM3IC0xMiAyMCAtMTMgNDAgLTM1CjggMTkgLTI0IDI3IDIwIDI1IC0xNCA2IC0zOQo5IC0xMCAyMyAtMjAgLTEzIC0yNyAtMjYgLTI3IC00NiAyOAo3IDEwIC00MiAtMzkgMzYgNDYgLTM0IC0zMQozIC00MCAzOSAxOQo5IDQwIDE3IC0xNSAxNiAtMjAgLTIzIDM2IDI1IDMKNyA3IDEzIDM0IDMyIDM5IC01IC00MAo4IDI4IC0zNiAxMiAyNSAzMCAtOCAtMjYgLTE5CjMgNDMgLTE2IC0zNgo2IC0zIC0yOSAtOCA0IC00MyAtMzgKNSAzOSAtMjIgLTQ1IDIzIDMxCjQgLTQ3IC0zNSAzMSAtMjYKNCAwIC0zOSAtMyAtMzYKMyAyNyAtNDggLTI2CjUgNDEgLTM1IDExIC0yNCA0MwozIDM2IC00OCAxOQo5IDI5IC0zOCAtMTcgLTQyIC0yMiAtNDEgMzIgLTEyIC02CjkgLTI3IC00MyAxNCA5IC00NSAyNiAtMzggMzkgMAo2IC0xNyAtNSA0MyAxMCAyMiAtMjkKNiA0OCAtNDMgNTAgMzYgLTMwIC0zMAo4IDE3IC0xOCAtMzUgMjYgNiAzNSAtMjggLTQ5CjEwIDM3IDIgMjIgMTUgLTExIDMzIC01IC0xIDM0IC0xOAo1IDIxIDM4IC00OSA4IDQ0CjQgLTggNDQgLTQ1IDE5CjcgLTMzIC0yMCA0NyAxMSAtNSAyOCAtMTQKOCAyNSAzMSAyOSAtMzQgNDEgLTExIC0xIDQ1CjkgMzMgLTQwIC01MCAyNiAtMjYgMzkgLTggLTMwIC0yMAo2IDMxIDcgLTIgNDAgMzYgMjIKOSAtNDYgMSAzOSAyMiAzIDQ4IDM0IDQwIC00NQo1IDcgLTQyIC0xNyAzOSAtMzAKMTAgMTcgMTIgMjEgMjcgNDYgLTUwIC00NiAxMyAtOSAtMTEKMTAgLTQ0IDMgLTI2IDIwIDMxIC00MCA0MiAtMzQgLTQ5IDEKOSAtMTAgLTUwIC0yMyAtNDkgNDEgNDYgLTUwIDM2IDE3CjQgLTI2IC0zNSAyNyAzMwo2IC0xMiAtMTUgMzggLTI3IC0zOCAxMAo5IDMwIC00MCAtNDggLTE1IDcgLTM2IC0xOCAtMzMgMzMKOCAtMzYgLTMxIC0xNSAtNDggLTQ1IC00NSAtMjQgMzcKNyAyMSAtMTAgLTQgMjIgLTQ1IDQ1IDM5CjEwIDQxIDMyIDggMzEgNSAtMyAxOCAtMjggLTI0IC0yCjcgLTQ5IC0zMyAtMzEgLTE2IC04IC03IC0zCjQgLTcgNDkgMjkgLTQ2CjMgLTE2IC0zMCAtMzEKNyAtNCAwIDIwIC0zNCAtMTMgLTM2IDExCjYgLTQ0IC0xMSAtMjggMTYgNDMgLTQxCjcgMSAtOCAtMTIgMyAtMzcgLTM4IDIxCjEwIDEwIC03IC03IC0zNSAxMSAtMzYgMzkgMTMgNCAtNDYKNyAtOCA0NCAzNyAtMzEgLTI5IDMwIDIyCjkgMzEgLTM5IC00MiAtNDAgLTI1IDQ1IC0yMiAtNDMgLTEKMyAtMzggMCAyMQo3IDcgMTIgNTAgMjQgNDEgMzYgLTIzCjkgLTQwIC0zIC0yMiAtMTcgMjQgNDkgLTI5IDUgLTI2CjggLTM2IC00MiAzOSAtNDcgMTcgNyA0NiAzNgo2IC0zNSAxMyAwIC0xOCAtMjQgMzIKMyAtMjMgMjkgLTMyCjQgLTI1IDggLTIgLTQKNSAtMzcgMjYgMTIgLTMyIDIyCjkgMzEgMzcgNCAxNiAxMyAzNiAtOSAxMyAxMwo2IDE5IDI4IC0yMiAtNDkgLTcgNDAKOCAtOSAtNDYgMTcgLTMyIC0xOCAyNyA1MCAtMzEKOSAyNCAtMTMgNDEgNDAgMTAgLTQyIC00MCAxNiAtNDUKNCAtMjIgLTM0IC00NSAtMTIKMyA0NyA3IC04CjUgLTMxIDMzIDggLTMgMTQKOSAxNyAxNCAtNDYgMjMgLTM5IDM2IDE2IDQ3IDI2CjQgNDUgNCA0NiAtMjQKNyAxOCAyNiAzIDExIC0xIDI3IDI1CjYgLTQ4IDM0IC01MCA0NCAtMjcgLTEyCjcgLTggLTQyIDEzIC0xNyAtMTIgNDggMgo5IC0xIC00MyAtMzAgMzIgLTM0IC0yMCAtMTQgNDMgLTgKMyAtNDYgMTEgMwo1IDEyIDI3IDQxIC00MCAzNgo1IC01IDIgLTQ2IDI4IDkKOSA4IC00NCAtMzggMTAgNDkgLTMxIC00OCAtNDYgMjYKNSAzMCAtOSAtMzcgMzkgMjAKOCAtMjYgLTEgNTAgNDkgNDkgMTIgLTM2IC00MwoxMCAyOCAzMCAtNyAzMyAtMzUgMzcgNDEgMjkgLTEzIDUwCjUgLTEgLTEzIDQ1IDM3IC0zNQo2IC00NiA1MCAwIDYgLTMgNDYKNiA4IC01IDMwIC00MSAtNDUgLTQ1CjEwIC0xOCAtNDcgMTYgMzUgMjIgMjMgLTIzIC0yMSAtMzkgNDkKOSAxNCAtMTEgLTM2IC0zMiA0IDIyIDQgLTQwIC0zNwo5IC00MiAtMzggMyA0OSAtMzEgNDMgLTQ3IDcgNQo5IC00NyAxMyAtOSA0MiAtMTggLTQwIC01IC00MSAtMzUKOCAzOCAtNDcgLTYgLTYgLTI4IC00OSAtMjEgLTQKNCAyNiAtMzIgLTI0IC01MAo2IDM0IDM2IDQzIC0zNSA0NSAtNTAKNyAtMyAzOCAtNDcgMjcgLTIxIC0zMiAtMjcKMTAgLTM2IDExIC02IDQwIC0xNyAtMzQgLTQ3IC0yNCAtNCAtOAoxMCAtMTMgLTEzIDIwIDMxIC05IC0yNyAyNSAtNDAgLTM3IDE4CjcgLTMwIC0yIC0zMiAtMzQgLTIyIC0xMCAxNQo="

type testCase struct {
	arr []int
}

// Embedded solver logic from 1838A.go.
func solve(tc testCase) string {
	lastNeg := 0
	hasNeg := false
	maxVal := 0
	for i, x := range tc.arr {
		if x < 0 {
			hasNeg = true
			lastNeg = x
		}
		if i == 0 || x > maxVal {
			maxVal = x
		}
	}
	if hasNeg {
		return strconv.Itoa(lastNeg)
	}
	return strconv.Itoa(maxVal)
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesA)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d n: %v", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing value %d", i+1, j)
			}
			arr[j], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d value %d: %v", i+1, j, err)
			}
		}
		cases = append(cases, testCase{arr: arr})
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
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
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
		fmt.Fprintf(&input, "1\n%d\n", len(tc.arr))
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		want := solve(tc)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
