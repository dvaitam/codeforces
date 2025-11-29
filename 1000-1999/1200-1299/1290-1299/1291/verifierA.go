package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded source for the reference solution (1291A.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       var s string
       fmt.Fscan(reader, &s)
       if n < 2 {
           fmt.Println(-1)
           continue
       }
       var fi, se byte
       ct := 0
       for i := 0; i < n; i++ {
           c := s[i]
           if (c-'0')%2 == 1 {
               if fi == 0 {
                   fi = c
                   ct++
               } else if se == 0 {
                   se = c
                   ct++
               }
           }
           if ct == 2 {
               break
           }
       }
       if ct < 2 {
           fmt.Println(-1)
       } else {
           fmt.Printf("%c%c\n", fi, se)
       }
   }
}
`

const testcasesRaw = `14 70487647593824
6 294892
11 21578156593
19 8784080160975351393
9 387115871
11 94185839894
16 2965934232094711
6 301868
10 9339694775
4 6917
20 63304135256012309891
2 23
20 26151090321730086914
4 4145
15 308709163457923
3 325
18 519720769845642807
4 6084
6 475945
20 34661093523376960696
3 371
10 3787890075
11 80638120665
2 40
2 99
5 41934
10 3176104714
6 951240
3 134
19 6590977658236940224
12 651590042294
13 7824173042814
14 64611877551717
15 145229611133060
5 78847
17 46153492635110873
5 86430
8 31376582
5 82966
18 857738930555082492
14 57118013204075
7 3758688
3 289
4 7348
15 869930024894517
10 5666022345
3 176
6 891256
3 867
3 272
2 19
6 618536
17 20979519426418306
16 6375100740899331
18 968412696116116207
15 707541511505520
9 619230310
11 60932271754
6 135574
11 95291189426
6 335833
7 5560290
14 21264862964513
16 6806606573547121
10 2892762662
9 858257173
11 17970341489
6 771738
14 50140046896714
13 5391014485183
7 2644829
18 481686447592211669
16 3845763778570742
17 19305760811605019
2 54
9 394316756
7 6662728
12 332756676337
8 16031250
7 4949184
13 7708869774735
10 1002504021
15 496837359191558
16 6408035138533444
18 747530481077706777
5 21312
15 479168602373245
12 718498347897
19 5430191263340886016
10 2953843318
11 63574922088
12 690262281237
20 43236599271908957740
9 927785142
14 45460305357345
7 5059802
13 1700300351056
6 476254
7 6660648
19 8091662082982153230`

var _ = solutionSource

type testCase struct {
	n int
	s string
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedA(n int, s string) string {
	var digits []byte
	for i := 0; i < len(s); i++ {
		if (s[i]-'0')%2 == 1 {
			digits = append(digits, s[i])
			if len(digits) == 2 {
				break
			}
		}
	}
	if len(digits) < 2 {
		return "-1"
	}
	return string(digits)
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid line %d: %q", idx+1, line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n on line %d: %v", idx+1, err)
		}
		s := fields[1]
		tests = append(tests, testCase{n: n, s: s})
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		expect := expectedA(tc.n, tc.s)
		input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
