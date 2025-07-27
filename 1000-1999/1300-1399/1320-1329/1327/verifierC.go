package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expectedPath(n, m int) string {
	var sb strings.Builder
	for i := 0; i < n-1; i++ {
		sb.WriteByte('U')
	}
	for i := 0; i < m-1; i++ {
		sb.WriteByte('L')
	}
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			for j := 0; j < m-1; j++ {
				sb.WriteByte('R')
			}
		} else {
			for j := 0; j < m-1; j++ {
				sb.WriteByte('L')
			}
		}
		sb.WriteByte('D')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(3)
	t := 100
	ns := make([]int, t)
	ms := make([]int, t)
	ks := make([]int, t)
	starts := make([][][2]int, t)
	ends := make([][][2]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		k := rand.Intn(5)
		ns[i], ms[i], ks[i] = n, m, k
		st := make([][2]int, k)
		ed := make([][2]int, k)
		for j := 0; j < k; j++ {
			st[j] = [2]int{rand.Intn(n) + 1, rand.Intn(m) + 1}
		}
		for j := 0; j < k; j++ {
			ed[j] = [2]int{rand.Intn(n) + 1, rand.Intn(m) + 1}
		}
		starts[i] = st
		ends[i] = ed
	}

	for idx := 0; idx < t; idx++ {
		n, m, k := ns[idx], ms[idx], ks[idx]
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for j := 0; j < k; j++ {
			input.WriteString(fmt.Sprintf("%d %d\n", starts[idx][j][0], starts[idx][j][1]))
		}
		for j := 0; j < k; j++ {
			input.WriteString(fmt.Sprintf("%d %d\n", ends[idx][j][0], ends[idx][j][1]))
		}
		in := input.String()
		wantPath := expectedPath(n, m)
		want := fmt.Sprintf("%d\n%s\n", len(wantPath), wantPath)

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Runtime error: %v\n%s", err, out.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		want = strings.TrimSpace(want)
		if got != want {
			fmt.Printf("Wrong answer on test %d\nExpected:\n%s\nGot:\n%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
