package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded source for the reference solution (was 1258D.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd computes the greatest common divisor of x and y using the Euclidean algorithm.
func gcd(x, y int64) int64 {
   if x < 0 {
       x = -x
   }
   if y < 0 {
       y = -y
   }
   for y != 0 {
       x, y = y, x%y
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b int64
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   result := gcd(a, b)
   fmt.Println(result)
}
`

const testcasesRaw = `925 594
307 579
741 787
437 333
779 101
786 836
762 497
427 946
822 710
158 120
910 869
450 706
400 646
697 301
592 957
157 797
133 549
985 824
423 965
794 174
569 933
206 20
850 727
759 198
749 434
467 782
324 97
885 978
621 161
439 990
991 468
424 693
127 287
247 572
521 277
527 534
640 815
16 33
92 449
372 162
461 381
532 29
610 322
781 420
163 977
212 186
756 723
83 996
871 433
227 636
446 912
11 659
39 212
252 458
41 625
900 744
926 438
266 25
434 787
152 66
130 522
159 737
976 778
340 486
210 450
18 805
97 803
213 818
293 296
371 683
552 703
961 734
969 531
966 693
810 156
885 762
923 858
106 522
940 798
241 334
804 489
904 450
485 333
103 844
371 36
837 748
158 739
850 159
749 919
756 462
304 593
780 378
425 313
895 771
418 728
578 39
152 113
19 338
843 155
636 138`

var _ = solutionSource

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type testCase struct {
	a int64
	b int64
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line %d", idx+1)
		}
		a, _ := strconv.ParseInt(parts[0], 10, 64)
		b, _ := strconv.ParseInt(parts[1], 10, 64)
		tests = append(tests, testCase{a: a, b: b})
	}
	return tests, nil
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
		want := fmt.Sprintf("%d", gcd(tc.a, tc.b))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
