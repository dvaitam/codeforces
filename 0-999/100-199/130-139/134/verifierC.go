package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
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

type Pair struct{ first, second int }

type solver struct {
	first []int
	nxt   []int
	mm    []Pair
	ans   []Pair
}

func newSolver(maxn int) *solver {
	return &solver{
		first: make([]int, maxn),
		nxt:   make([]int, maxn),
		mm:    make([]Pair, maxn),
		ans:   make([]Pair, 0, maxn),
	}
}

func (s *solver) tj(x, y int) {
	s.nxt[x] = s.first[y]
	s.first[y] = x
}

func (s *solver) solve(n, sum int, a []int) string {
	if sum%2 == 1 {
		return "No"
	}
	s.ans = s.ans[:0]
	for i := range s.first {
		s.first[i] = 0
	}
	for i := range s.nxt {
		s.nxt[i] = 0
	}
	for i := 1; i <= n; i++ {
		s.tj(i, a[i-1])
	}
	for i := sum; i >= 1; i-- {
		for s.first[i] != 0 {
			x := s.first[i]
			s.first[i] = s.nxt[x]
			j := i
			for k := 1; k <= i; k++ {
				for j > 0 && s.first[j] == 0 {
					j--
				}
				if j == 0 {
					return "No"
				}
				mid := s.first[j]
				s.mm[k] = Pair{mid, j}
				s.first[j] = s.nxt[mid]
				s.ans = append(s.ans, Pair{x, mid})
			}
			for k := 1; k <= i; k++ {
				p := s.mm[k]
				s.tj(p.first, p.second-1)
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("Yes\n")
	sb.WriteString(fmt.Sprintln(len(s.ans)))
	for _, p := range s.ans {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.first, p.second))
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(r *rand.Rand, slv *solver) (string, string) {
	n := r.Intn(10) + 1
	sum := r.Intn(10*n + 1)
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, sum))
	for i := 0; i < n; i++ {
		if sum > 0 {
			arr[i] = r.Intn(sum + 1)
		} else {
			arr[i] = 0
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	expect := slv.solve(n, sum, arr)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	slv := newSolver(201000)
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r, slv)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
