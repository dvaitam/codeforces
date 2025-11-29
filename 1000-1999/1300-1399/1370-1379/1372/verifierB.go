package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `100
140893
596855
888600
841237
800877
66174
267461
123648
519503
797928
471327
495187
683246
398057
827038
220155
98420
511556
29726
936712
876365
408746
453791
636946
799310
804425
2210
729635
467024
279269
756591
840777
239876
619871
991190
107194
945217
332851
32077
23408
26683
681100
567714
9654
984771
924042
399723
719832
227122
442623
761113
30453
553261
232462
800800
459160
984789
519898
579717
244408
362495
242083
709729
229410
797913
481931
998502
303860
971514
22535
436398
878266
960780
583486
966986
673496
104859
194938
659926
758792
901721
310789
126764
779247
348858
939080
756533
745740
525128
981931
442613
532382
870357
954400
702868
199073
318106
297964
616124
925348`

type testCase struct {
	n int64
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func referenceSolve(n int64) (int64, int64) {
	var sp int64
	if n%2 == 0 {
		sp = 2
	} else {
		for i := int64(3); i*i <= n; i += 2 {
			if n%i == 0 {
				sp = i
				break
			}
		}
	}
	if sp == 0 {
		sp = n
	}
	b := n / sp
	a := n - b
	return a, b
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesB))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("empty testcases")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing n for case %d", i+1)
		}
		n, err := strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse n case %d: %w", i+1, err)
		}
		cases = append(cases, testCase{n: n})
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("1\n%d\n", tc.n)
		wantA, wantB := referenceSolve(tc.n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		parts := strings.Fields(got)
		if len(parts) != 2 {
			fmt.Printf("case %d failed: expected two numbers got %q\n", idx+1, got)
			os.Exit(1)
		}
		a, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			fmt.Printf("case %d failed: bad first number %q\n", idx+1, parts[0])
			os.Exit(1)
		}
		b, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			fmt.Printf("case %d failed: bad second number %q\n", idx+1, parts[1])
			os.Exit(1)
		}
		if a != wantA || b != wantB {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%d %d\ngot:\n%d %d\n", idx+1, input, wantA, wantB, a, b)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
