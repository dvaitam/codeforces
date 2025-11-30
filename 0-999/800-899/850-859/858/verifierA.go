package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (from testcasesA.txt) to keep the verifier standalone.
const rawTestcases = `906691060 6
813847340 6
43469774 4
548977049 7
434794719 4
511742082 5
626401696 3
541903390 2
302621085 2
811538591 1
663968656 4
976832603 8
757172937 2
333018423 1
783650879 1
965120264 5
506959382 8
108127102 5
466188457 5
655934895 3
593267778 7
475338373 8
279701489 0
864392047 8
983541587 0
100149904 6
762628806 0
657019496 7
889126175 5
261897307 5
755530427 1
205156724 3
256211902 2
862407392 8
481003666 1
86378037 5
939617817 8
525367681 1
323676027 8
312556225 1
587810203 5
874527132 8
218186328 8
630949022 4
477803330 1
640258140 6
340426464 3
311738932 2
203357386 2
35403864 4
511671264 1
96448169 2
941425019 2
991472819 0
904584779 1
964406046 8
733900803 6
899650944 8
295959888 8
871479690 3
912128613 3
961040770 6
622442781 4
483788449 7
708933077 5
88447325 5
657970850 1
522315486 5
907395137 3
260957516 0
785430573 4
125771989 3
399496591 2
357057976 6
876080017 0
108026500 2
918395672 3
48569709 8
646575639 1
28665467 1
681825965 3
651050913 1
420057910 1
397434736 1
39075653 0
208940442 2
771068319 1
514573219 3
780776981 0
729443502 0
584368328 6
666364151 1
897543643 4
75166547 3
77279012 4
376125287 6
193614874 0
540775580 7
42282563 1
750892481 6
`

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(n int64, k int) int64 {
	ten := int64(1)
	for i := 0; i < k; i++ {
		ten *= 10
	}
	g := gcd(n, ten)
	return n * (ten / g)
}

type testCase struct {
	n int64
	k int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcases, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields, got %d", idx+1, len(fields))
		}
		n, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %w", idx+1, err)
		}
		cases = append(cases, testCase{n: n, k: k})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n", tc.n, tc.k)
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := strconv.FormatInt(solve(tc.n, tc.k), 10)
		input := buildInput(tc)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
