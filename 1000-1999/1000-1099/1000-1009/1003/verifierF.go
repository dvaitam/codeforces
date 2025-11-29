package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	words []string
}

const testcasesFData = `100
5
bcc cab c a b
8
b aca c b ab a c baa
1
a
4
a ba cca c
4
ba ba b b
5
cac cba bb cb b
8
aab c ba bbb aba aca a bc
6
bba cbb b a cc ab
5
cab bba c c bba
1
c
7
a cbb cb a abb a cac
7
a ba b ab acb cc bb
5
ba b acb aa bc
6
b cca bc cbb cc bc
3
abc bc bb
6
bc cb ac bb aca b
8
ca ccb c ccb cc bca a bcc
3
c c aba
3
a b b
4
abb bc aca c
6
b bb c bac bab cbc
5
ab cc bba c cb
4
c b aba cc
3
cc bb b
5
cb ac bcb bb b
5
ba aa cba aa ab
5
b a ba bbb bcc
2
a bc
6
bcc b bc caa b bca
4
b cc ccc bb
7
cab ab c cb bba b cbb
3
cc a aba
2
abc ab
8
a a cac b baa aba c a
2
bca ac
2
bca b
2
cc bb
6
a bac ccb cca ac b
5
aab cb a acb c
1
a
8
b bb c c baa ca aac aa
7
ab cc bb c cb bbc c
4
aba a a a
1
bc
6
b bc bb aa c a
6
ba bab b aab b ba
1
bba
2
bbb aca
1
c
1
aca
7
bcb b ba cba a aac bcb
1
ab
7
a c cba aab b c c
1
ba
1
bb
5
aaa c cba aab a
3
aca abc bb
7
cac b bac a ba bb a
2
ac b
4
ab aac b aa
7
bc b aca ac aca bba ca
6
b a bc ba c caa
2
a bb
1
c
7
b b aa cac b cac abb
7
b ba cb b c c bc
1
a
7
aaa c b bbb a ab cba
6
ac c bb a b aa
3
ac ab cb
2
b aac
3
ab acb bcb
7
cc cac b ca cab aa b
7
ab aaa c bc ccc b aaa
6
bba bba cac bbb c ccb
4
a bb aaa baa
2
ccc cbb
4
bc abb ac ab
1
a
5
ccc b b acc ca
1
ac
2
ab bac
8
a c a b b acb bac cca
6
a a ba cc a cb
7
a bac b bc b cba b
6
abb b cc c b ba
7
a ab caa ac ba baa ba
2
a aa
5
bbc a a ab abc
8
aaa cb ccc caa b bbb a ab
4
ab b ab b
6
ac aab a bab a acb
6
bab a ca a ab c
6
aa cc aba c bb bc
8
c b c c aa cc a a
7
c cab c a bba bbb c
2
cc cc
7
ba cac ac c bac cb ab
6
b c acb a aab b
6
bb c b ac ba acb
`

func solve(words []string) int {
	// Embedded solution logic from 1003F.go.
	n := len(words)
	a := make([]int, n+1)
	lens := make([]int, n+1)
	myMap := make(map[string]int)
	for i := 1; i <= n; i++ {
		s := words[i-1]
		if myMap[s] == 0 {
			myMap[s] = i
			a[i] = i
			lens[i] = len(s)
		} else {
			a[i] = myMap[s]
		}
	}
	sum := n - 1
	for i := 1; i <= n; i++ {
		sum += lens[a[i]]
	}
	ans := 0
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			cnt := 1
			x := -1
			for k := i; k <= j; k++ {
				x += lens[a[k]]
			}
			poi := j + 1
			length := j - i + 1
			for poi+length-1 <= n {
				ok := true
				for k := 0; k < length; k++ {
					if a[i+k] != a[poi+k] {
						ok = false
						break
					}
				}
				if ok {
					cnt++
					poi += length
				} else {
					poi++
				}
			}
			if cnt > 1 {
				if v := x * cnt; v > ans {
					ans = v
				}
			}
		}
	}
	return sum - ans
}

func parseTestCases(data string) ([]testCase, error) {
	tokens := strings.Fields(data)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no embedded testcases found")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d missing length", i+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d has invalid length: %w", i+1, err)
		}
		idx++
		if idx+n > len(tokens) {
			return nil, fmt.Errorf("test %d missing %d words", i+1, idx+n-len(tokens))
		}
		words := make([]string, n)
		copy(words, tokens[idx:idx+n])
		idx += n
		cases = append(cases, testCase{words: words})
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("embedded data has %d extra tokens", len(tokens)-idx)
	}
	return cases, nil
}

func runCase(bin string, words []string) error {
	n := len(words)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, w := range words {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(w)
	}
	sb.WriteByte('\n')
	input := sb.String()
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	want := solve(words)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestCases(testcasesFData)
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.words); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
