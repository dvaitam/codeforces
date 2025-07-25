package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

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

func expected(n, s int, a []int) string {
	costChief := 0
	if a[s-1] != 0 {
		costChief = 1
	}
	r := make([]int, n+2)
	for i := 1; i <= n; i++ {
		if i == s {
			continue
		}
		ai := a[i-1]
		if ai < 0 {
			ai = 0
		}
		if ai > n {
			ai = n
		}
		r[ai]++
	}
	r0 := r[0]
	total := n - 1
	big := total - r0
	bestP := total + 5
	z := 0
	for D := 1; D <= n; D++ {
		if D <= n {
			big -= r[D]
		}
		if D < len(r) {
			if r[D] == 0 {
				z++
			}
		} else {
			z++
		}
		A := r0 + big
		P := z
		if A > P {
			P = A
		}
		if P < bestP {
			bestP = P
		}
	}
	ans := costChief + bestP
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	r := rand.New(rand.NewSource(5))
	for tc := 1; tc <= 100; tc++ {
		n := r.Intn(10) + 1
		sID := r.Intn(n) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = r.Intn(n)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, sID)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", arr[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := expected(n, sID, arr)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", tc, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", tc, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
