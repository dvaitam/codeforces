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

func expected(a []int) string {
	n := len(a)
	if a[n-1] != 0 {
		return "NO"
	}
	flag := -1
	for i := 0; i < n-1; i++ {
		if a[i] == 0 {
			flag = i
			break
		}
	}
	if flag == n-2 {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	if flag >= 0 {
		if flag > 0 {
			sb.WriteString("(")
			for i := 0; i < flag-1; i++ {
				fmt.Fprintf(&sb, "%d->", a[i])
			}
			fmt.Fprintf(&sb, "%d)->", a[flag-1])
		}
		sb.WriteString("(0)->(")
		for i := flag + 1; i < n-1; i++ {
			fmt.Fprintf(&sb, "%d->", a[i])
		}
		fmt.Fprintf(&sb, "%d)", a[n-2])
		sb.WriteString("->0")
	} else {
		for i := 0; i < n-1; i++ {
			fmt.Fprintf(&sb, "%d->", a[i])
		}
		sb.WriteString("0")
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(2)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteString("\n")
	input := sb.String()
	exp := expected(a)
	return input, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
