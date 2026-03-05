package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesB = `4 ae a cdba cdba
4 ha c de mj
1 pnkdfh
2 psjgwnq avrmps
4 lepha leph g ep
4 ebompclq clqkdijfh cl k
3 ehitoa hito l
5 h a hj edi jfg
1 dpi
2 egc f
5 rqaijlmn hkprqaijlmn mnc gb edfhkprqaij
5 gfajm dp bih rlkdpo lkdpocstqngfajm
3 acrdhnmbplfeqj crdhnmbplfeqj dhnmbplfeqjkig
1 c
1 nkam
3 fns vqeipmywuhxodfn eipmyw
3 gdeiclhjf deicl deic
1 kf
1 ghkafde
1 cab
3 a a a
2 b ac
2 bda d
4 ig ig ig dejhbfa
5 jid jidfg dfga ec ji
2 ceb h
1 adhfecg
2 qhbpj mi
4 esranci iglm kfesra gl
3 tas rle d
4 g bdfgi g d
3 a xopybwvnfkigulqcejrzs kigulqc
3 a bc dbcae
4 a a a a
4 l bfiajneg l ne
1 ckgjb
3 j efohsqg rmpnilbe
4 mav jpncb wkti dqgox
1 cjgdhm
4 d ef fdba ba
2 g gabecfd
4 a a a a
4 a a a a
1 a
1 b
4 cbd da fc g
3 cdhbegram be h
5 ead bfc bf a ea
4 kdbgfcjeiahl fcjeia l fcjeiah
2 xsmitkg kgrqhn
5 c bd d bd bd
2 ogusnciprhbjfmtke tkedv
1 tpdexacqi
1 bac
4 e bc de bcad
1 jgebhdifk
4 mfvie tjgaolbn gao jga
4 e cb dcbiehf ehfg
3 aedc a aedcghf
4 dhec hecf idhecf cf
5 rbahtg fijwdlpmsoqurbahtgevc jwdlpmsoqurbahtgev cn c
4 ihbafdce afdce ce ihbafdce
2 db db
4 bjchi iuef l pk
2 km hne
2 gptacihjfdnmvuskqbero kqbero
1 b
1 kchobd
2 f h
3 a cba ba
5 a a a a a
2 ace e
3 agbdch i dche
3 rlfomnpgusc e u
1 e
4 sd xz lvqe unptirl
3 e h cd
3 jbdgck k ljbd
1 fmdalceh
3 jq cflbtpnsrkyghaidwu dwu
4 f e d fcae
5 busd xrcmbus rcm sd sd
2 hkz g
4 qi afmcjbdlgqi n cjbdlgq
3 de ebc deb
1 hd
4 ghsy ci yn nxdue
4 fbqolgk lg ap djapifbqolgkme
4 ekhanmbfgloijc lo d cd
1 ij
1 gfmjnc
2 kjh dek
5 jcb hoaps i bhoap de
2 bfghc ca
2 c i
5 f b a edc af
3 bdkoiecq rah oiecqpr
4 krblpashqnjt jtfg gc gc
3 kerognmtj tj jqi
3 hft ftckgopqiw wvaldbnmeu`


type testCase struct {
	n     int
	frags []string
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesB), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n", idx+1)
		}
		if len(fields)-1 < n {
			return nil, fmt.Errorf("line %d: expected %d fragments got %d", idx+1, n, len(fields)-1)
		}
		frags := make([]string, n)
		copy(frags, fields[1:1+n])
		cases = append(cases, testCase{n: n, frags: frags})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

// validate checks that out is a valid minimum-length answer for the given fragments.
func validate(frags []string, out string) string {
	// Count distinct letters across all fragments (minimum required length).
	letterSet := make(map[rune]bool)
	for _, f := range frags {
		for _, c := range f {
			letterSet[c] = true
		}
	}
	minLen := len(letterSet)

	if len(out) != minLen {
		return fmt.Sprintf("length %d, want %d", len(out), minLen)
	}

	// All characters in output must be distinct.
	seen := make(map[rune]bool)
	for _, c := range out {
		if seen[c] {
			return fmt.Sprintf("duplicate character %c", c)
		}
		seen[c] = true
	}

	// Every fragment must be a substring of output.
	for _, f := range frags {
		if !strings.Contains(out, f) {
			return fmt.Sprintf("fragment %q not found in %q", f, out)
		}
	}

	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, s := range tc.frags {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if msg := validate(tc.frags, got); msg != "" {
			fmt.Fprintf(os.Stderr, "case %d failed: %s\nfragments: %v\ngot: %s\n", idx+1, msg, tc.frags, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
