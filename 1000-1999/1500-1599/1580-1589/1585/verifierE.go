package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type freqGroup struct {
	freq   int
	values []int
	count  int
}

type QuerySpec struct {
	k             int
	filteredTotal int
	valueFreq     map[int]int
	groups        []freqGroup
}

type Case struct {
	input    string
	specs    []QuerySpec
	ansCount int
}

func runBinary(bin, input string) (string, error) {
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

type query struct {
	v, l, k int
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		var sb strings.Builder
		n := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(n) + 1
		}
		parent := make([]int, n+1)
		for v := 2; v <= n; v++ {
			parent[v] = rng.Intn(v-1) + 1
		}
		qs := make([]query, q)
		fmt.Fprintf(&sb, "1\n%d %d\n", n, q)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(a[j]))
		}
		sb.WriteByte('\n')
		for vtx := 2; vtx <= n; vtx++ {
			if vtx > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(parent[vtx]))
		}
		sb.WriteByte('\n')
		for j := 0; j < q; j++ {
			qs[j].v = rng.Intn(n) + 1
			qs[j].l = rng.Intn(n) + 1
			qs[j].k = rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", qs[j].v, qs[j].l, qs[j].k)
		}
		specs := buildSpecs(n, q, a, parent, qs)
		cases[i] = Case{input: sb.String(), specs: specs, ansCount: q}
	}
	return cases
}

func buildSpecs(n, q int, values []int, parent []int, queries []query) []QuerySpec {
	specs := make([]QuerySpec, q)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = values[i-1]
	}
	for idx, qq := range queries {
		freq := make(map[int]int)
		v := qq.v
		for v > 0 {
			freq[a[v]]++
			v = parent[v]
		}
		valueFreq := make(map[int]int)
		groupMap := make(map[int][]int)
		total := 0
		for val, cnt := range freq {
			if cnt >= qq.l {
				valueFreq[val] = cnt
				groupMap[cnt] = append(groupMap[cnt], val)
				total++
			}
		}
		var freqs []int
		for f := range groupMap {
			freqs = append(freqs, f)
		}
		sort.Ints(freqs)
		groups := make([]freqGroup, 0, len(freqs))
		for _, f := range freqs {
			groups = append(groups, freqGroup{
				freq:   f,
				values: append([]int(nil), groupMap[f]...),
				count:  len(groupMap[f]),
			})
		}
		specs[idx] = QuerySpec{
			k:             qq.k,
			filteredTotal: total,
			valueFreq:     valueFreq,
			groups:        groups,
		}
	}
	return specs
}

func parseAnswers(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	ans := make([]int, expected)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		ans[i] = v
	}
	return ans, nil
}

func answerValid(spec QuerySpec, ans int) bool {
	if spec.filteredTotal < spec.k {
		return ans == -1
	}
	if ans == -1 {
		return false
	}
	freq, ok := spec.valueFreq[ans]
	if !ok {
		return false
	}
	before := 0
	for _, g := range spec.groups {
		if g.freq < freq {
			before += g.count
			continue
		}
		if g.freq == freq {
			if before < spec.k && spec.k <= before+g.count {
				for _, v := range g.values {
					if v == ans {
						return true
					}
				}
			}
			return false
		}
		break
	}
	return false
}

func runCase(bin string, c Case) error {
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	answers, err := parseAnswers(got, c.ansCount)
	if err != nil {
		return err
	}
	for i, spec := range c.specs {
		if !answerValid(spec, answers[i]) {
			return fmt.Errorf("query %d invalid answer %d", i+1, answers[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
