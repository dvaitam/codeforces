package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

type DSU struct {
	parent []int
}

func newDSU(n int) *DSU {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &DSU{p}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int) bool {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return false
	}
	d.parent[y] = x
	return true
}

func solveF(n, m int, vectors [][]int) (string, string) {
	d := newDSU(m + 1) // extra node 0
	usedIdx := []int{}
	for i, vec := range vectors {
		if len(vec) == 1 {
			if d.union(0, vec[0]) {
				usedIdx = append(usedIdx, i+1)
			}
		} else {
			if d.union(vec[0], vec[1]) {
				usedIdx = append(usedIdx, i+1)
			}
		}
	}
	rank := len(usedIdx)
	pow := int64(1)
	for i := 0; i < rank; i++ {
		pow = (pow * 2) % mod
	}
	idxStr := strings.TrimSpace(strings.Join(func() []string {
		s := make([]string, len(usedIdx))
		for i, v := range usedIdx {
			s[i] = fmt.Sprint(v)
		}
		return s
	}(), " "))
	return fmt.Sprintf("%d %d", pow, rank), idxStr
}

func genCases() []string {
	rand.Seed(6)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		m := rand.Intn(4) + 1
		n := rand.Intn(5) + 1
		vectors := make([][]int, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				vectors[j] = []int{rand.Intn(m) + 1}
			} else {
				a := rand.Intn(m) + 1
				b := rand.Intn(m) + 1
				for b == a {
					b = rand.Intn(m) + 1
				}
				vectors[j] = []int{a, b}
			}
		}
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, v := range vectors {
			sb.WriteString(fmt.Sprint(len(v)))
			for _, x := range v {
				sb.WriteByte(' ')
				sb.WriteString(fmt.Sprint(x))
			}
			sb.WriteByte('\n')
		}
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n, m int
		fmt.Sscan(lines[0], &n, &m)
		vectors := make([][]int, n)
		for j := 0; j < n; j++ {
			parts := strings.Fields(lines[j+1])
			k := 0
			fmt.Sscan(parts[0], &k)
			vec := make([]int, k)
			for t := 0; t < k; t++ {
				fmt.Sscan(parts[t+1], &vec[t])
			}
			vectors[j] = vec
		}
		want1, want2 := solveF(n, m, vectors)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		gs := strings.SplitN(got, "\n", 2)
		if len(gs) < 1 {
			fmt.Fprintf(os.Stderr, "Wrong format on case %d\n", i+1)
			os.Exit(1)
		}
		if gs[0] != want1 || (len(gs) > 1 && strings.TrimSpace(gs[1]) != want2) {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected:\n%s\n%s\nGot:\n%s\n", i+1, tc, want1, want2, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
