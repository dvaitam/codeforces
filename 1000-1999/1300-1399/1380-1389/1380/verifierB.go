package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
SRP
PP
SPRRPRPP
RSPPSRSRPR
R
RPSRPSRSR
PSRPRSRP
RPSSR
SSP
SP
PSSRPPSPS
SRPRSPP
PSS
RPSSRR
PPPSRPRPS
SSPSRRSRRR
SRPSPSPPP
SSRPSSRSS
PRPP
SRSPPPPPRS
SSPPSRRSR
SRRSPRSRR
P
P
PRSR
PRRRPS
SPS
PSPPP
RP
PPRPRPS
RSPRRRPRR
PSS
SRSSSPR
SRPSSPSSP
S
RRRPR
PP
PSP
RSR
RSPRSSSRPR
RRSSPS
PRSP
SPRPS
PRRRPSR
PRPSRP
PSSPSRRSR
RR
SRP
SSPPPP
PR
SPRSSRPRPR
RRPRSSP
SS
SRPP
SSRPP
RP
S
R
RRRRSPR
PR
RSRP
SPSPSPP
RS
RRRPSS
PPPPRR
SPRPRS
SPSPPRSRP
RPRP
PR
SPRPPRPRPS
RPRSS
SRRRRRPRPS
SR
S
P
PPRRSP
SS
RRR
PRSSSP
RRS
P
SSSSRRPPSR
S
PRSP
SPSPSPR
PRPPRSP
RRSPSRSRRP
PSPRS
RP
R
PSSPSSSRR
PSPRSP
SSSSPS
RRSS
RSRPPS
SPRSS
SRRSSSPR
PPR
SRPSPSSPPS
SSR`

// Embedded reference logic from 1380B.go.
func solve(s string) string {
	cntR, cntP, cntS := 0, 0, 0
	for _, ch := range s {
		switch ch {
		case 'R':
			cntR++
		case 'P':
			cntP++
		case 'S':
			cntS++
		}
	}
	var ans strings.Builder
	if cntR == cntP && cntP == cntS {
		for _, ch := range s {
			switch ch {
			case 'R':
				ans.WriteByte('P')
			case 'P':
				ans.WriteByte('S')
			case 'S':
				ans.WriteByte('R')
			}
		}
	} else {
		cmx := 'R'
		maxc := cntR
		if cntP > maxc {
			cmx = 'P'
			maxc = cntP
		}
		if cntS > maxc {
			cmx = 'S'
			maxc = cntS
		}
		var play byte
		switch cmx {
		case 'R':
			play = 'P'
		case 'P':
			play = 'S'
		case 'S':
			play = 'R'
		}
		for i := 0; i < len(s); i++ {
			ans.WriteByte(play)
		}
	}
	return ans.String()
}

func parseTestcases(raw string) ([]string, error) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty test data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count")
	}
	if len(fields)-1 != t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(fields)-1)
	}
	return fields[1:], nil
}

func run(bin, input string) (string, error) {
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

func acceptUniformBest(s, ans string) bool {
	if len(ans) != len(s) || len(ans) == 0 {
		return false
	}
	ch := ans[0]
	if ch != 'R' && ch != 'P' && ch != 'S' {
		return false
	}
	for i := 1; i < len(ans); i++ {
		if ans[i] != ch {
			return false
		}
	}
	cntR, cntP, cntS := 0, 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'R':
			cntR++
		case 'P':
			cntP++
		case 'S':
			cntS++
		default:
			return false
		}
	}
	maxc := cntR
	if cntP > maxc {
		maxc = cntP
	}
	if cntS > maxc {
		maxc = cntS
	}
	allowed := map[byte]bool{}
	if cntR == maxc {
		allowed['P'] = true
	}
	if cntP == maxc {
		allowed['S'] = true
	}
	if cntS == maxc {
		allowed['R'] = true
	}
	return allowed[ch]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse tests:", err)
		os.Exit(1)
	}

	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(tests)))
	expected := make([]string, len(tests))
	for i, s := range tests {
		input.WriteString(s)
		input.WriteByte('\n')
		expected[i] = solve(s)
	}

	got, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "candidate failed:", err)
		os.Exit(1)
	}

	outLines := strings.Fields(got)
	if len(outLines) != len(expected) {
		fmt.Fprintf(os.Stderr, "wrong number of outputs: expected %d got %d\n", len(expected), len(outLines))
		os.Exit(1)
	}
	for i, exp := range expected {
		ans := outLines[i]
		if ans != exp && !acceptUniformBest(tests[i], ans) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, exp, ans)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
