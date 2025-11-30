package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const mod int64 = 998244353
const testcasesB64 = "MTAwCjQgNgoyIDEgMQozIDEgNAo0IDEgMgo0IDEwCjIgMSA4CjMgMiAxCjQgMSA1CjQgNgoyIDEgMQozIDIgNQo0IDEgMwo1IDMKMiAxIDMKMyAyIDIKNCAzIDIKNSAyIDMKMyAxMAoyIDEgMgozIDEgOAo0IDcKMiAxIDMKMyAyIDIKNCAzIDUKNSA5CjIgMSA4CjMgMSA0CjQgMyA1CjUgMiAxCjIgMQoyIDEgMQozIDEwCjIgMSA3CjMgMSA2CjQgOQoyIDEgMQozIDEgNQo0IDIgMgo0IDUKMiAxIDEKMyAxIDMKNCAyIDQKMiAzCjIgMSAyCjQgOAoyIDEgMQozIDEgOAo0IDEgOAo0IDgKMiAxIDEKMyAyIDYKNCAzIDEKMyAxMAoyIDEgNgozIDEgMwo0IDEKMiAxIDEKMyAxIDEKNCAxIDEKNCA2CjIgMSAzCjMgMiAxCjQgMiAyCjQgOAoyIDEgNAozIDIgMwo0IDEgNgoyIDQKMiAxIDQKMiA4CjIgMSA0CjQgOAoyIDEgNgozIDEgMgo0IDMgMwozIDEKMiAxIDEKMyAxIDEKNCAxCjIgMSAxCjMgMiAxCjQgMSAxCjIgMTAKMiAxIDQKNCA0CjIgMSAzCjMgMiAzCjQgMiAxCjQgMwoyIDEgMgozIDEgMQo0IDMgMwoyIDQKMiAxIDMKNCAzCjIgMSAzCjMgMSAxCjQgMSAyCjQgOQoyIDEgMQozIDEgMQo0IDIgMQoyIDcKMiAxIDEKNSA3CjIgMSA2CjMgMiA1CjQgMiAzCjUgMyA3CjMgOQoyIDEgOAozIDEgMgo1IDIKMiAxIDIKMyAyIDIKNCAxIDEKNSAyIDIKNCA3CjIgMSA2CjMgMiA3CjQgMiA0CjUgMQoyIDEgMQozIDEgMQo0IDIgMQo1IDIgMQo0IDEwCjIgMSA5CjMgMSA3CjQgMyAyCjQgOQoyIDEgNQozIDIgNwo0IDMgNwozIDEwCjIgMSA5CjMgMSA4CjUgMgoyIDEgMQozIDIgMgo0IDMgMgo1IDIgMQoyIDEwCjIgMSAyCjUgNgoyIDEgNAozIDEgMwo0IDIgMQo1IDEgNgo1IDYKMiAxIDYKMyAxIDEKNCAxIDIKNSAxIDIKMyA4CjIgMSA3CjMgMiA0CjMgNgoyIDEgNQozIDIgNAo0IDIKMiAxIDIKMyAxIDIKNCAxIDIKMiAxCjIgMSAxCjQgMQoyIDEgMQozIDIgMQo0IDMgMQo1IDcKMiAxIDUKMyAxIDYKNCAxIDQKNSAyIDIKNSA2CjIgMSAyCjMgMiA2CjQgMSAyCjUgMSAxCjIgMwoyIDEgMwo1IDMKMiAxIDEKMyAyIDIKNCAzIDMKNSA0IDMKMyA4CjIgMSA4CjMgMiA2CjMgMwoyIDEgMwozIDEgMwo1IDQKMiAxIDQKMyAyIDEKNCAxIDQKNSAzIDEKNSAyCjIgMSAyCjMgMiAxCjQgMyAxCjUgMiAyCjIgOQoyIDEgMQo0IDkKMiAxIDgKMyAyIDkKNCAyIDYKNSA2CjIgMSAzCjMgMSAxCjQgMyA1CjUgMyA2CjQgMQoyIDEgMQozIDEgMQo0IDMgMQo0IDgKMiAxIDMKMyAxIDUKNCAyIDMKNCA1CjIgMSAxCjMgMSA1CjQgMSA1CjMgOAoyIDEgMwozIDIgNAo1IDMKMiAxIDIKMyAxIDEKNCAzIDEKNSA0IDIKNSA0CjIgMSAyCjMgMSAzCjQgMSAzCjUgMSAxCjMgOAoyIDEgOAozIDIgOAozIDcKMiAxIDQKMyAyIDIKNSAxCjIgMSAxCjMgMiAxCjQgMyAxCjUgMSAxCjUgMTAKMiAxIDYKMyAyIDUKNCAxIDUKNSA0IDgKMyAxMAoyIDEgOQozIDIgMwo0IDkKMiAxIDEKMyAxIDEKNCAyIDIKMiA4CjIgMSA0CjQgNQoyIDEgMQozIDEgNQo0IDEgMQozIDEwCjIgMSAyCjMgMSAxMAozIDcKMiAxIDEKMyAxIDMKNCA4CjIgMSA4CjMgMiA3CjQgMyA2CjQgMwoyIDEgMgozIDEgMgo0IDIgMgoyIDUKMiAxIDIKNSAyCjIgMSAxCjMgMiAxCjQgMyAyCjUgMSAxCjUgNAoyIDEgNAozIDIgMQo0IDEgMQo1IDEgMwo1IDEKMiAxIDEKMyAxIDEKNCAzIDEKNSAzIDEKNCA3CjIgMSA1CjMgMiA0CjQgMiA0CjIgOQoyIDEgMQo1IDgKMiAxIDYKMyAyIDUKNCAyIDEKNSAxIDYKNCAyCjIgMSAyCjMgMSAxCjQgMyAxCjUgNAoyIDEgNAozIDEgMwo0IDIgMQo1IDQgMgozIDEKMiAxIDEKMyAxIDEKNSAyCjIgMSAyCjMgMiAxCjQgMSAyCjUgMSAyCjIgNwoyIDEgMgo1IDkKMiAxIDQKMyAxIDEKNCAzIDMKNSA0IDEKNSAxMAoyIDEgNAozIDEgMQo0IDEgMQo1IDMgMgo1IDEwCjIgMSAxMAozIDEgMgo0IDEgOAo1IDIgMgozIDcKMiAxIDcKMyAyIDMKMiAxMAoyIDEgNgozIDEKMiAxIDEKMyAxIDEKNCAzCjIgMSAxCjMgMiAzCjQgMSAyCjMgMTAKMiAxIDMKMyAxIDEwCjUgMwoyIDEgMgozIDEgMQo0IDEgMQo1IDEgMwo0IDcKMiAxIDYKMyAxIDIKNCAxIDYKNCA5CjIgMSA3CjMgMSAzCjQgMyAxCjMgOAoyIDEgNQozIDIgNwo="

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	for i := range p {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{p, sz}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func powMod(a, e int64) int64 {
	a %= mod
	if a < 0 {
		a += mod
	}
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func expected(n int, S int64, edges [][3]int64) string {
	freq := make(map[int64]int64)
	for _, e := range edges {
		freq[e[2]]++
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i][2] < edges[j][2] })
	dsu := NewDSU(n)
	counts := make(map[int64]int64)
	for _, e := range edges {
		u := dsu.find(int(e[0]) - 1)
		v := dsu.find(int(e[1]) - 1)
		if u != v {
			counts[e[2]] += int64(dsu.size[u]) * int64(dsu.size[v])
			dsu.union(u, v)
		}
	}
	ans := int64(1)
	zero := false
	for w, c := range counts {
		exp := c - freq[w]
		if exp <= 0 {
			continue
		}
		diff := S - w
		if diff <= 0 {
			zero = true
			break
		}
		ans = ans * powMod(diff%mod, exp) % mod
	}
	if zero {
		return "0"
	}
	return fmt.Sprintf("%d", ans%mod)
}

func runCase(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputs, exps := loadCases()
	for idx, input := range inputs {
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exps[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}

func loadCases() ([]string, []string) {
	data, err := base64.StdEncoding.DecodeString(testcasesB64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode embedded testcases: %v\n", err)
		os.Exit(1)
	}
	tokens := strings.Fields(string(bytes.TrimSpace(data)))
	if len(tokens) == 0 {
		fmt.Fprintln(os.Stderr, "no testcases found")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid test count header\n")
		os.Exit(1)
	}
	pos := 1
	var inputs []string
	var exps []string
	for caseNum := 1; caseNum <= t; caseNum++ {
		if pos+1 >= len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing header\n", caseNum)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		S, errS := strconv.ParseInt(tokens[pos+1], 10, 64)
		if errN != nil || errS != nil {
			fmt.Fprintf(os.Stderr, "invalid header on case %d\n", caseNum)
			os.Exit(1)
		}
		pos += 2
		edges := make([][3]int64, n-1)
		if pos+3*(n-1) > len(tokens) {
			fmt.Fprintf(os.Stderr, "case %d missing edges\n", caseNum)
			os.Exit(1)
		}
		for i := 0; i < n-1; i++ {
			u, errU := strconv.ParseInt(tokens[pos+3*i], 10, 64)
			v, errV := strconv.ParseInt(tokens[pos+3*i+1], 10, 64)
			w, errW := strconv.ParseInt(tokens[pos+3*i+2], 10, 64)
			if errU != nil || errV != nil || errW != nil {
				fmt.Fprintf(os.Stderr, "invalid edge on case %d\n", caseNum)
				os.Exit(1)
			}
			edges[i] = [3]int64{u, v, w}
		}
		pos += 3 * (n - 1)
		want := expected(n, S, edges)
		lines := []string{"1", fmt.Sprintf("%d %d", n, S)}
		for _, e := range edges {
			lines = append(lines, fmt.Sprintf("%d %d %d", e[0], e[1], e[2]))
		}
		inputs = append(inputs, strings.Join(lines, "\n")+"\n")
		exps = append(exps, want)
	}
	return inputs, exps
}
