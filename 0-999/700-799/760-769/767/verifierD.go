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

const INF = 1070000000

type testCase struct {
	n, m, k int
	a       []int
	tasks   []int
}

func solve(n, m, k int, a []int, tasks []int) string {
	aa := append([]int(nil), a...)
	sort.Ints(aa)
	type task struct{ d, id int }
	ts := make([]task, m)
	for i := 0; i < m; i++ {
		ts[i] = task{tasks[i], i + 1}
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].d < ts[j].d })
	headA, headT := 0, 0
	ans := []int{}
	for d := 0; ; d++ {
		for headT < m && ts[headT].d < d {
			headT++
		}
		if headA < n && aa[headA] < d {
			return "-1"
		}
		cnt := 0
		for cnt < k && (headT < m || headA < n) {
			td := INF
			if headT < m {
				td = ts[headT].d
			}
			if headA >= n || td < aa[headA] {
				ans = append(ans, ts[headT].id)
				headT++
			} else {
				headA++
			}
			cnt++
		}
		if headT >= m && headA >= n {
			break
		}
	}
	sort.Ints(ans)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	if len(ans) > 0 {
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func (tc testCase) input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.tasks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	k := rng.Intn(3) + 1
	a := make([]int, n)
	tasks := make([]int, m)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10)
	}
	for i := 0; i < m; i++ {
		tasks[i] = rng.Intn(10)
	}
	return testCase{n: n, m: m, k: k, a: a, tasks: tasks}
}

func runProgram(bin, input string) (string, error) {
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

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc.n, tc.m, tc.k, append([]int(nil), tc.a...), append([]int(nil), tc.tasks...))
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{randomCase(rng)}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
