package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// referenceSolutionSource embeds the original 391B solution for traceability.
const referenceSolutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   ans := 1
   // For each starting position p0 as bottom of pile
   for p0 := 0; p0 < n; p0++ {
       // Memoization for state (p_tip, width)
       memo := make(map[int]int)
       // dfs returns max height from state (p_tip, w)
       var dfs func(pTip, w int) int
       dfs = func(pTip, w int) int {
           key := pTip*(n+1) + w
           if v, ok := memo[key]; ok {
               return v
           }
           best := 1
           // f is fold index: between f and f+1, left width = f+1
           // require pTip <= f and new pos pNew < w
           // f <= floor((w + pTip - 2)/2)
           maxF := (w + pTip - 2) / 2
           for f := pTip; f <= maxF; f++ {
               // new position from reflection of pTip over f
               pNew := 2*f - pTip + 1
               if pNew < 0 || pNew >= w {
                   continue
               }
               if s[pNew] != s[p0] {
                   continue
               }
               h := 1 + dfs(pNew, f+1)
               if h > best {
                   best = h
               }
           }
           memo[key] = best
           return best
       }
       h := dfs(p0, n)
       if h > ans {
           ans = h
       }
   }
   fmt.Println(ans)
}
`

var _ = referenceSolutionSource

const rawTestcases = `M 1
Y 1
N 1
B 1
I 1
Q 1
P 1
M 1
Z 1
J 1
P 1
L 1
S 1
G 1
Q 1
EJ 1
EY 1
DT 1
ZI 1
RW 1
ZT 1
EJ 1
DX 1
CV 1
KP 1
RD 1
LN 1
KT 1
UG 1
RP 1
OQI 1
BZR 1
ACX 1
MWZ 1
VUA 1
TPK 1
HXK 1
WCG 1
SHH 2
ZEZ 1
ROC 1
CKQ 1
PDJ 1
RJW 1
DRK 1
RGZT 1
RSJO 1
CTZM 1
KSHJ 1
FGFB 1
TVIP 1
CCVY 2
EEBC 2
WRVM 1
WQIQ 1
ZHGV 1
SNSI 1
OPVU 1
WZLC 1
KTDP 1
SUKGH 1
AXIDW 1
HLZFK 1
NBDZE 1
WHBSU 1
RTVCA 1
DUGTS 1
DMCLD 1
BTAGF 1
WDPGX 1
ZBVAR 1
NTDIC 1
HCUJL 1
NFBQO 1
BTDWM 1
GILXPS 1
FWVGYB 1
ZVFFKQ 2
IDTOVF 1
APVNSQ 1
JULMVI 1
ERWAOX 1
CKXBRI 1
EHYPLT 1
JVLSUT 1
EWJMXN 1
UCATGW 1
KFHHUO 2
MWVSNB 1
MWSNYV 1
WBFOCIW 1
FOQPRTY 1
ABPKJOB 2
ZZNGRUC 2
XEAMVNK 1
AGAWYAV 2
QTDGDTU 1
GJIWFDP 1
MUCAIOZ 1
ZDIEUQU 1`

type testcase struct {
	input    string
	expected int
	hasExp   bool
}

func parseTestcases(data string) ([]testcase, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	out := make([]testcase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		tc := testcase{input: fields[0]}
		if len(fields) > 1 {
			v, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("parse expected for %q: %w", line, err)
			}
			tc.expected = v
			tc.hasExp = true
		}
		out = append(out, tc)
	}
	return out, nil
}

func solveCase(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	ans := 1
	for p0 := 0; p0 < n; p0++ {
		memo := make(map[int]int)
		var dfs func(pTip, w int) int
		dfs = func(pTip, w int) int {
			key := pTip*(n+1) + w
			if v, ok := memo[key]; ok {
				return v
			}
			best := 1
			maxF := (w + pTip - 2) / 2
			for f := pTip; f <= maxF; f++ {
				pNew := 2*f - pTip + 1
				if pNew < 0 || pNew >= w {
					continue
				}
				if s[pNew] != s[p0] {
					continue
				}
				h := 1 + dfs(pNew, f+1)
				if h > best {
					best = h
				}
			}
			memo[key] = best
			return best
		}
		h := dfs(p0, n)
		if h > ans {
			ans = h
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(rawTestcases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solveCase(tc.input)
		if tc.hasExp && expect != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d embedded expected %d but oracle got %d\n", idx+1, tc.expected, expect)
			os.Exit(1)
		}

		got, err := run(bin, tc.input+"\n")
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != strconv.Itoa(expect) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
