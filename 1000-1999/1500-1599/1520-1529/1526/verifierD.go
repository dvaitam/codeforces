package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `100
TAAAANTTNNAAATTONTA
TNONATTANOT
NNANAO
OONNTOATATN
NTTNTOANNTTNA
ANNOAOANNTOAONANTO
ANONOAOTOOAOAO
NONNNNTAON
TAATAONOAATO
NAOONT
NNATNOTOOA
TONAT
TOOTNOATATNAA
NOATOAATTONNN
ANOO
NTNNONAONAONAATT
OANATOONOO
TNA
TOOO
AOONOOOTNATOTAA
ATNANAONONA
TOOAATOT
NTANANA
OAN
TOOOOTAN
NNANA
AT
TTAONNOOTATNTOO
OOOT
AOOATOAAOANAAAO
AATNOOONTTTOOAO
OANT
OOONNTTOTAANAOO
OATONOTAOOO
OAOATNOOOTATNNONAATT
AOAATONNTATAO
NAAO
OAOTOOONANTATO
NA
NAONNTTTNTONAOOOON
AAANONAANNANN
NTTAOOOANAOOOAOOAT
OONNNNA
ANOOTOTNTTO
ONON
TAOOOONTNNN
TTNAONATTAATONTANTT
TTTAONNOATNTNNOA
OTONOATNONOAOTN
N
TNONNTTOTNATTOONOTOT
NTNNTTATNOTAOONNNONN
NNAANOOONAOO
TOOO
AAOT
TONATTAOONAANTOAOOTA
NNOT
AAAON
ATNOAANONNNTTN
AOTNAOO
NTNAAOTAA
OOAO
NOAA
TNONNNNAATTNT
NAONOATT
TOTNTNN
TOTOTNNTANOOAOTNTNO
OAOATAAAATOTOATATOA
AAOAAAANTOA
NANN
TOTNATAN
TTANN
NTTT
NATNTNAOT
TTTONAONOONTOOTOONA
ONANTNONTT
OOTATTAN
OOAT
NANATTOTTAOTA
ANN
OTNOTN
ONN
NAATOOAT
NANATOAOATNOTAT
TNNAANTOONAOOAN
TATO
TNANNTNOA
NA
OOTOOOONATAANNOTNOO
OAATATTAATNT
TAOTOO
ATONA
TTNANONO
OO
NANONTOAATTAN
OTNO
ONAONATTOONNA
NNONNTNAANTTTANO
ANOTTAOATO
AAOTNOOA
`

type testCase struct {
	s string
}

func solveCase(tc testCase) string {
	s := tc.s
	n := len(s)
	idx := make([][]int, 4)
	for i := range idx {
		idx[i] = make([]int, 0, n)
	}
	for i := 0; i < n; i++ {
		var v int
		switch s[i] {
		case 'A':
			v = 0
		case 'N':
			v = 1
		case 'T':
			v = 2
		case 'O':
			v = 3
		}
		idx[v] = append(idx[v], i)
	}

	var swaps [4][4]int64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				continue
			}
			ii, jj := 0, 0
			ni, nj := len(idx[i]), len(idx[j])
			for ii < ni && jj < nj {
				if idx[i][ii] < idx[j][jj] {
					swaps[i][j] += int64(jj)
					ii++
				} else {
					jj++
				}
			}
			swaps[i][j] += int64(ni-ii) * int64(jj)
		}
	}

	bestVal := int64(-1)
	order := [4]int{0, 1, 2, 3}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 4; k++ {
				if k == i || k == j {
					continue
				}
				l := 6 - i - j - k
				if l < 0 || l >= 4 || l == i || l == j || l == k {
					continue
				}
				tans := swaps[i][j] + swaps[i][k] + swaps[i][l]
				tans += swaps[j][k] + swaps[j][l]
				tans += swaps[k][l]
				if tans > bestVal {
					bestVal = tans
					order = [4]int{i, j, k, l}
				}
			}
		}
	}
	mp := []byte{'A', 'N', 'T', 'O'}
	var res []byte
	for _, c := range order {
		for range idx[c] {
			res = append(res, mp[c])
		}
	}
	return string(res)
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	if len(fields) != t+1 {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(fields)-1)
	}
	res := make([]testCase, t)
	for i := 0; i < t; i++ {
		res[i] = testCase{s: fields[i+1]}
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(tc.s)
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
