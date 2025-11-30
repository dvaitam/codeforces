package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `100
150 37
737 983
165 457
722 519
695 437
558 853
226 1000
646 817
712 529
462 229
537 665
32 405
692 590
823 329
676 647
437 61
756 306
129 992
218 897
49 314
73 880
79 318
940 962
306 762
163 427
579 259
134 9
575 900
871 39
605 840
223 986
923 584
472 176
848 889
891 998
799 721
638 522
39 388
206 356
102 211
588 691
919 444
606 199
505 107
961 682
400 304
517 512
18 334
627 893
412 922
289 19
161 206
879 336
831 577
802 139
348 440
219 273
691 99
858 389
955 561
353 937
904 858
704 548
497 787
546 241
67 743
42 87
137 174
171 933
552 219
275 778
341 615
519 862
262 377
347 349
117 299
241 889
967 619
799 978
733 909
501 139
594 565
789 107
329 41
417 75
390 887
808 151
849 129
350 118
630 602
801 949
388 79
585 564
230 580
84 976
274 374
913 303
578 548
948 118
469 919`

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func parseTestcases() ([][2]int, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	if len(fields) != 1+2*t {
		return nil, fmt.Errorf("expected %d pairs, found %d numbers", t, len(fields)-1)
	}
	res := make([][2]int, t)
	for i := 0; i < t; i++ {
		a, err := strconv.Atoi(fields[1+2*i])
		if err != nil {
			return nil, fmt.Errorf("parse a%d: %v", i+1, err)
		}
		b, err := strconv.Atoi(fields[1+2*i+1])
		if err != nil {
			return nil, fmt.Errorf("parse b%d: %v", i+1, err)
		}
		res[i] = [2]int{a, b}
	}
	return res, nil
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

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d %d\n", tc[0], tc[1])
	}

	var expected strings.Builder
	for i, tc := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(strconv.Itoa(gcd(tc[0], tc[1])))
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("verifier failed\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
