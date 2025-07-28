package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expectedQuery(pref []int, pos map[int][]int, l, r int) string {
	total := pref[r] ^ pref[l-1]
	if total == 0 {
		return "YES"
	}
	list1 := pos[pref[r]]
	i := sort.Search(len(list1), func(i int) bool { return list1[i] >= l })
	if i == len(list1) || list1[i] >= r {
		return "NO"
	}
	j1 := list1[i]
	list2 := pos[pref[l-1]]
	j := sort.Search(len(list2), func(i int) bool { return list2[i] > j1 })
	if j == len(list2) || list2[j] >= r {
		return "NO"
	}
	return "YES"
}

func solve(n, q int, a []int, queries [][2]int) []string {
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] ^ a[i-1]
	}
	pos := make(map[int][]int, n+1)
	for i := 0; i <= n; i++ {
		v := pref[i]
		pos[v] = append(pos[v], i)
	}
	ans := make([]string, len(queries))
	for idx, qr := range queries {
		l, r := qr[0], qr[1]
		ans[idx] = expectedQuery(pref, pos, l, r)
	}
	return ans
}

func genCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(8) + 2
	q := rng.Intn(5) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(4)
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n-1) + 1
		r := rng.Intn(n-l+1) + l
		if r == l {
			r++
			if r > n {
				l--
			}
		}
		queries[i] = [2]int{l, r}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		fmt.Fprintf(&sb, "%d %d\n", queries[i][0], queries[i][1])
	}
	expect := solve(n, q, a, queries)
	return sb.String(), expect
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, expArr := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != len(expArr) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(expArr), len(gotLines), inp)
			os.Exit(1)
		}
		for j, exp := range expArr {
			if strings.TrimSpace(gotLines[j]) != exp {
				fmt.Fprintf(os.Stderr, "case %d failed at query %d\nexpected: %s\ngot: %s\ninput:\n%s", i+1, j+1, exp, gotLines[j], inp)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
