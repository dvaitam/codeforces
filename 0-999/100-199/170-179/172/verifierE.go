package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `<c><c><c/><a><d></d></a></c></c>|1|c c a d
<d/>|2|d|d
<b></b>|3|a|a b|b
<d/>|3|d|a|d
<d></d>|1|d
<b></b>|3|b|b|b
<c></c>|1|c a
<c><c></c><a><b/><c></c></a></c>|2|c a b|c a
<b></b>|2|d|a
<d/>|3|d|d|d c
<a></a>|2|a|a
<c><a><c></c><d><d></d></d></a></c>|3|d c|c a d d|c
<a><c/></a>|2|a|b
<d></d>|1|b
<b></b>|3|b|b|b
<b></b>|3|b|c d|b
<c></c>|1|b d
<c></c>|2|a a|c
<a/>|1|b
<d><b></b><b></b></d>|3|d b|d|a a
<d><b></b></d>|2|b a|b
<a><c><d></d></c></a>|3|d|c c|a
<a></a>|3|a|a|c
<c></c>|1|c b
<a><b><b><b/><d></d></b></b></a>|2|c a|d
<a><c></c></a>|2|a c|d a
<c></c>|3|c|b|c
<d></d>|1|d
<b><b></b></b>|1|d
<d><c></c></d>|2|d c|d
<a></a>|3|d b|a|a
<b></b>|1|c
<b/>|1|b
<c><d><c></c></d></c>|2|d|d
<c><a><a><a/></a></a></c>|3|d d|b a|c
<d/>|3|d|c d|d
<c><b><c><a/><b></b></c></b><b/></c>|1|c b c a
<a></a>|1|c
<a></a>|2|c|a
<a><c/></a>|2|a|a c
<d><d></d></d>|2|c|d
<d><d></d></d>|2|b d|d d
<c></c>|1|c
<a/>|2|a|a
<b></b>|2|c|a
<d><d></d></d>|3|d|d|c
<a></a>|1|d b
<d><b><c></c></b></d>|1|b
<c><a><c></c></a></c>|1|c a
<b></b>|1|d a
<c></c>|2|c|d c
<d><c></c><b></b></d>|3|d b|d c|d c
<c/>|2|c|b a
<a></a>|1|a
<c><b><a><b></b><c></c></a><c></c></b><d><b><a/></b><b></b></d></c>|2|d c|c c
<c><c><d><a/></d></c><d><d><b></b><b/></d></d></c>|1|a
<c><b></b><d></d></c>|2|a|a
<b/>|2|b|b
<c><b><d></d><b></b></b><b></b></c>|1|c c
<b><c></c><d/></b>|2|d a|b c
<d/>|3|d|d|a b
<b><c><a><d/></a></c><d><a></a></d></b>|2|b c|b c a d
<a></a>|2|b b|a
<b><b><c></c><a/></b></b>|2|a|c c
<c><a><d></d></a></c>|3|d|b c|c
<b><b></b></b>|2|c|d
<a></a>|1|a a
<b></b>|2|b|a
<d></d>|3|d|d|d
<b><a></a><a><c><b></b><a></a></c><c></c></a></b>|2|b c|b
<d/>|2|a b|a a
<a></a>|2|a|a
<a><b></b></a>|3|a b|b a|a
<b><b></b><a/></b>|2|b b|d
<c/>|1|b d
<b><b><c><b/><d/></c></b></b>|3|b|a c|b b c d
<b><a><b/></a></b>|2|d a|b c
<a><d><a><a/><b></b></a></d></a>|1|a d
<b></b>|3|c|a|c
<b><a></a><d><b/></d></b>|3|b d|b|b a
<c><a></a></c>|2|b a|d d
<c><d/></c>|3|b a|b c|c a
<a><d></d><b></b></a>|1|c
<c></c>|2|c|c
<c></c>|3|a|d c|c
<b></b>|2|b c|b
<d></d>|3|d|c d|d
<b><c><d/><c><a></a><a/></c></c><d></d></b>|1|b c
<b></b>|1|c a
<a><d></d><d></d></a>|1|a
<c><b><c><c></c></c></b><d><a></a><b/></d></c>|3|d d|c d a|b b
<c><d></d><d></d></c>|2|c d|c
<c><c></c></c>|3|a d|a|a b
<a></a>|1|a
<a></a>|2|c c|a
<d><c></c></d>|3|a c|d|d c
<c><b><d><a/></d><b><c></c></b></b><c><c/></c></c>|2|c b b|c b b
<c></c>|1|c
<c><a></a></c>|2|a|c a
<a><a><b><a></a></b></a><d><a/></d></a>|2|a c|a a b`

type occ struct {
	qid       int
	positions []int
}

type testCase struct {
	doc     string
	queries []string
}

func runProg(bin, input string) (string, error) {
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

func parseQueries(raw []string) ([][]string, []int, map[string]map[int][]int) {
	m := len(raw)
	queries := make([][]string, m)
	lengths := make([]int, m)
	tagMap := make(map[string]map[int][]int)
	for i, line := range raw {
		parts := strings.Fields(strings.TrimSpace(line))
		for _, token := range parts {
			queries[i] = append(queries[i], token)
			if tagMap[token] == nil {
				tagMap[token] = make(map[int][]int)
			}
			tagMap[token][i] = append(tagMap[token][i], len(queries[i]))
		}
		lengths[i] = len(queries[i])
	}
	return queries, lengths, tagMap
}

func solve(doc string, queryLines []string) []int {
	_, lengths, tagMap := parseQueries(queryLines)

	occMap := make(map[string][]occ)
	for tag, m1 := range tagMap {
		tmp := make([]occ, 0, len(m1))
		for qi, ps := range m1 {
			for l, r := 0, len(ps)-1; l < r; l, r = l+1, r-1 {
				ps[l], ps[r] = ps[r], ps[l]
			}
			tmp = append(tmp, occ{qid: qi, positions: ps})
		}
		occMap[tag] = tmp
	}

	m := len(queryLines)
	dp := make([][]bool, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]bool, lengths[i]+1)
		dp[i][0] = true
	}
	ans := make([]int, m)

	var changeStack []int
	type change struct{ qid, pos int }
	var changes []change

	n := len(doc)
	for i := 0; i < n; {
		if doc[i] != '<' {
			i++
			continue
		}
		j := i + 1
		isClose := false
		if j < n && doc[j] == '/' {
			isClose = true
			j++
		}
		k := j
		for k < n && doc[k] != '>' {
			k++
		}
		content := doc[j:k]
		selfClose := false
		if !isClose && len(content) > 0 && content[len(content)-1] == '/' {
			selfClose = true
			content = content[:len(content)-1]
		}
		tag := content
		if isClose {
			last := changeStack[len(changeStack)-1]
			changeStack = changeStack[:len(changeStack)-1]
			for len(changes) > last {
				c := changes[len(changes)-1]
				dp[c.qid][c.pos] = false
				changes = changes[:len(changes)-1]
			}
		} else {
			prev := len(changes)
			changeStack = append(changeStack, prev)
			if occs, ok := occMap[tag]; ok {
				for _, oc := range occs {
					q := oc.qid
					for _, posi := range oc.positions {
						if !dp[q][posi] && dp[q][posi-1] {
							dp[q][posi] = true
							changes = append(changes, change{qid: q, pos: posi})
							if posi == lengths[q] {
								ans[q]++
							}
						}
					}
				}
			}
			if selfClose {
				last := changeStack[len(changeStack)-1]
				changeStack = changeStack[:len(changeStack)-1]
				for len(changes) > last {
					c := changes[len(changes)-1]
					dp[c.qid][c.pos] = false
					changes = changes[:len(changes)-1]
				}
			}
		}
		i = k + 1
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		doc := parts[0]
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d parse m: %v", idx+1, err)
		}
		if len(parts) != m+2 {
			return nil, fmt.Errorf("line %d expected %d queries got %d", idx+1, m, len(parts)-2)
		}
		queries := make([]string, m)
		for i := 0; i < m; i++ {
			queries[i] = strings.TrimSpace(parts[2+i])
		}
		cases = append(cases, testCase{doc: doc, queries: queries})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		expected := solve(tc.doc, tc.queries)

		var input strings.Builder
		input.WriteString(tc.doc)
		input.WriteByte('\n')
		fmt.Fprintf(&input, "%d\n", len(tc.queries))
		for _, q := range tc.queries {
			input.WriteString(q)
			input.WriteByte('\n')
		}

		got, err := runProg(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(outLines) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d wrong number of lines\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < len(expected); i++ {
			if strings.TrimSpace(outLines[i]) != strconv.Itoa(expected[i]) {
				fmt.Fprintf(os.Stderr, "case %d query %d failed: expected %d got %s\n", idx+1, i+1, expected[i], strings.TrimSpace(outLines[i]))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
