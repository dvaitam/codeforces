package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type node struct {
	next [26]int32
	cnt  int32
}

const testcasesRaw = `5
a
baaac
cbab
b
c

3
abbca
ccbaa
a

3
baa
ccb
cab

4
a
a
cab
acbaa

5
bcaaa
ba
bb
bc
acbaa

1
acb

2
cba
cccc

3
bc
c
cabb

5
bcc
ccbb
aaa
bcab
ccbc

2
a
ccaa

4
bcaab
bbba
ca
b

2
c
cccbc

3
ab
c
ccbc

1
cb

4
cccb
a
b
cbbc

2
bbc
c

4
cba
cbbc
ccbbc
bbc

2
abb
baca

3
ba
baaba
a

5
a
ac
cbcba
bbb
aabaa

1
cbcb

3
acbb
aacc
abc

3
ababa
a
abb

1
a

5
a
abcac
ac
aaaca
caaab

4
a
cbcc
ca
bc

4
aacca
aacbb
bcccc
bacc

1
cb

3
bcb
a
a

3
ca
ab
cabb

2
a
cba

1
bacbb

2
c
bb

1
bcb

1
aba

4
babac
cab
acba
bb

5
a
accc
b
cca
aac

2
aba
caacc

2
ccbc
bbbaa

4
aba
caca
cbbba
cb

4
cccc
ba
ca
caa

2
abb
cbbcc

4
babc
bbc
bcc
b

5
abbcc
caab
abcaa
cbcb
bcbcc

4
acacb
bb
a
a

5
cabb
cabb
bccc
ab
bbccc

1
aca

2
b
aa

4
acca
cb
acbab
abb

1
aac

1
cbc

3
ac
bbca
b

3
bcbb
ba
caba

2
bbbb
aab

2
aba
c

5
bcbaa
abb
a
bbbb
a

3
ac
ba
c

5
caa
cabb
caca
a
ca

4
abacc
a
bac
cbbc

2
baac
a

5
ccb
b
bbbaa
aba
aa

4
abaa
aaa
cbcbb
cbccb

3
c
abacc
b

4
aa
cccc
cccbb
cba

3
accc
c
aaac

1
aa

3
caa
a
bb

5
cbbab
c
ca
caba
ac

2
bbb
cbb

3
cbb
bbcbb
caa

4
bbc
ab
ac
ab

5
acacb
b
ba
a
accac

1
aacac

3
aab
abcc
ca

5
bcbaa
acac
ac
bcb
bc

3
aabaa
bb
bb

5
b
b
bbabc
b
aac

3
acca
cc
bccbc

1
bbba

4
b
a
bbc
bb

1
babab

2
aabcb
cabba

5
caa
a
a
ccaa
bbbc

3
ccab
bacc
accc

4
acb
aac
aab
abbc

5
ba
a
baac
acaab
ccca

2
cc
abbac

1
aaca

5
ab
cb
ccaab
acbcb
acbac

4
abb
abc
ba
c

5
abab
bbaac
bcaca
a
bb

5
ca
aaabb
aba
cbc
cb

5
ca
bcca
bb
ca
aa

1
bbc

3
ac
baba
babab

4
cc
ccba
caacb
acca

5
a
ba
aac
cacbb
b

2
c
cb

5
baa
babca
bbb
b
b

3
ba
aa
acbc`

func expected(strs []string) string {
	nodes := make([]node, 1) // root
	totalLen := 0
	for _, s := range strs {
		cur := 0
		for i := 0; i < len(s); i++ {
			idx := s[i] - 'a'
			child := nodes[cur].next[idx]
			if child == 0 {
				nodes = append(nodes, node{})
				nodes[cur].next[idx] = int32(len(nodes) - 1)
				child = nodes[cur].next[idx]
			}
			cur = int(child)
			nodes[cur].cnt++
		}
		totalLen += len(s)
	}

	var sumLCP int64
	for _, s := range strs {
		cur := 0
		for i := len(s) - 1; i >= 0; i-- {
			idx := s[i] - 'a'
			child := nodes[cur].next[idx]
			if child == 0 {
				break
			}
			cur = int(child)
			sumLCP += int64(nodes[cur].cnt)
		}
	}
	n := int64(len(strs))
	total := int64(totalLen)
	res := 2*n*total - 2*sumLCP
	return fmt.Sprintf("%d", res)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func loadCases() ([]string, []string) {
	blocks := strings.Split(strings.TrimSpace(testcasesRaw), "\n\n")
	var inputs []string
	var expects []string
	for idx, blk := range blocks {
		fields := strings.Fields(blk)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+1 {
			fmt.Printf("invalid testcase block %d\n", idx+1)
			os.Exit(1)
		}
		strs := fields[1:]
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for _, s := range strs {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		inputs = append(inputs, sb.String())
		expects = append(expects, expected(strs))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expects[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expects[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
