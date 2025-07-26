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

func expected(perms [][]int) string {
	M := len(perms)
	if M == 0 {
		return "0\n"
	}
	N := len(perms[0])
	a := make([][]int, M)
	for i := 0; i < M; i++ {
		a[i] = make([]int, N+2)
		copy(a[i][1:], perms[i])
	}
	father := make([]int, N+1)
	cnt := make([]int, N+1)
	for j := 1; j <= N; j++ {
		father[a[0][j]] = a[0][j+1]
	}
	for i := 1; i < M; i++ {
		for j := 1; j <= N; j++ {
			if father[a[i][j]] != a[i][j+1] {
				father[a[i][j]] = 0
			}
		}
	}
	for i := 1; i <= N; i++ {
		if father[i] == 0 {
			father[i] = i
		}
	}
	var getfather func(int) int
	getfather = func(x int) int {
		if father[x] != x {
			father[x] = getfather(father[x])
		}
		return father[x]
	}
	for i := 1; i <= N; i++ {
		root := getfather(i)
		if root != 0 {
			cnt[root]++
		}
	}
	ans := int64(N)
	for i := 1; i <= N; i++ {
		if cnt[i] > 1 {
			ans += int64(cnt[i]) * int64(cnt[i]-1) / 2
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func runCase(bin, input, want string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %q got %q", want, got)
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

	tests := [][][]int{{{1, 2, 3}}}
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		m := rng.Intn(3) + 1
		casep := make([][]int, m)
		for j := 0; j < m; j++ {
			perm := rand.Perm(n)
			for k := range perm {
				perm[k]++
			}
			casep[j] = perm
		}
		tests = append(tests, casep)
	}

	for idx, tc := range tests {
		m := len(tc)
		n := len(tc[0])
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < m; i++ {
			for j, v := range tc[i] {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			sb.WriteByte('\n')
		}
		want := expected(tc)
		if err := runCase(bin, sb.String(), strings.TrimSpace(want)); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
