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
ABBAAABBBABBABAABAB
ABBBBAABAAA
BAABABBBAAAA
BBABBAAAAAAABBBBBBAB
AB
BABAABBBABBBBAABBBAA
A
BABBBBBBB
BBBABBBBB
ABBBBBABABAA
A
BBAAABABBABA
BBAA
BAAAAABBAAAAABBBAAAB
BBBBABBBBAABBBBB
BAAAAAABBAAABBA
ABABBA
ABAAA
BABAAA
ABAABABBBABAAB
AAAABAAAA
AAABAAABABABAAAB
AAAA
BBBAAABABBAABB
AABB
ABAAABABABBABA
ABAAB
BBBBABAABAAAAAAABBBB
BABABBB
ABBBAABAABA
ABB
AABBB
BBABAAAAAAAAABBAAAA
ABBAABAB
BBABA
AB
BAABABBABAABBBA
AABBAABBABABAABBBB
BBAAABBAAB
A
BAABABBBAAAABABABAB
BAABAABBAABBB
BBBABABABAAA
BB
AABBAABBBBBAABBBBB
AABAABAA
AAABBBAAB
AAAAAB
BBBBBAAAAABABB
BAAABABAABBBBBBBBBA
ABBBAABABBBBAB
BAAABBAABBAABA
AABABBBBABBBA
BAAAAAABABBBBAAB
AAAAB
ABBBABBAAAABAB
A
AAA
BBABAAABABAAABABA
BBABAAA
BBBABAABAB
ABBBBBBBBABA
BAAAAABBAABAABBABBAB
ABABBABAAABAAABABB
BBAAAAA
ABABBBAAAA
BAAAABBBBBABBBBBAAA
BABBAAAAAABAAB
AABABB
AABAAABBAABBBABAABBA
BBABABB
BBBBBBAABB
ABA
AAA
BABAAAAABAAAAAABB
ABA
ABBABBBABBAAABA
ABBAA
BBAABABBABBA
B
AAAAAABAB
AB
BAAABBAAA
BABBBBAA
BAAAB
B
BABAABAAAAAA
ABABBAABABBBAABB
B
ABAAA
AAA
ABAAABA
BABA
ABBABA
BABAAAB
ABAB
ABBAAAAAAABBBABABA
BBABBBABB
BBAABB
ABBABAABAA`

func solve(s string) int {
	bpos := []int{}
	for i := 0; i < len(s); i++ {
		if s[i] == 'B' {
			bpos = append(bpos, i)
		}
	}
	if len(bpos) == 0 {
		return 0
	}
	k := len(bpos)
	segments := make([]int, k+1)
	segments[0] = bpos[0]
	for i := 0; i < k-1; i++ {
		segments[i+1] = bpos[i+1] - bpos[i] - 1
	}
	segments[k] = len(s) - 1 - bpos[k-1]
	edges := make([]int, 0, 2*k)
	for i := 0; i < k; i++ {
		edges = append(edges, segments[i])
		edges = append(edges, segments[i+1])
	}
	if len(edges) == 0 {
		return 0
	}
	dpPrev2 := 0
	dpPrev1 := edges[0]
	if dpPrev1 < 0 {
		dpPrev1 = 0
	}
	for i := 2; i <= len(edges); i++ {
		val := dpPrev1
		if dpPrev2+edges[i-1] > val {
			val = dpPrev2 + edges[i-1]
		}
		dpPrev2, dpPrev1 = dpPrev1, val
	}
	return dpPrev1
}

func expected(s string) string {
	return fmt.Sprintf("%d\n", solve(s))
}

func loadCases() ([]string, []string) {
	lines := strings.Fields(testcasesRaw)
	if len(lines) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(lines[0])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	if len(lines)-1 < t {
		fmt.Println("insufficient testcases embedded")
		os.Exit(1)
	}
	var inputs []string
	var expects []string
	for i := 0; i < t; i++ {
		s := lines[i+1]
		inputs = append(inputs, fmt.Sprintf("1\n%s\n", s))
		expects = append(expects, expected(s))
	}
	return inputs, expects
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		if err := runCase(exe, input, expects[idx]); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
