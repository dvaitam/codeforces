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

func expectedC(b []int) string {
	n := len(b)
	used := make([]bool, 2*n+1)
	for _, v := range b {
		if v >= 1 && v <= 2*n {
			used[v] = true
		}
	}
	rem := make([]int, 0, 2*n)
	for v := 1; v <= 2*n; v++ {
		if !used[v] {
			rem = append(rem, v)
		}
	}
	a := make([]int, 2*n)
	for i := 0; i < n; i++ {
		a[2*i] = b[i]
		idx := -1
		for j, rv := range rem {
			if rv > b[i] {
				idx = j
				break
			}
		}
		if idx == -1 {
			return "-1"
		}
		a[2*i+1] = rem[idx]
		rem = append(rem[:idx], rem[idx+1:]...)
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		used := make(map[int]bool)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			for {
				v := rand.Intn(2*n) + 1
				if !used[v] {
					used[v] = true
					b[j] = v
					break
				}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("1\n%d\n", n))
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := expectedC(b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
