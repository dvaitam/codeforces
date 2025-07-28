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

type node struct {
	val  int
	need bool
}

func solveD(n int, a, b []int, razors []int) string {
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return "NO"
		}
	}
	cnt := make(map[int]int)
	for _, x := range razors {
		cnt[x]++
	}
	need := make(map[int]int)
	stack := make([]node, 0)
	for i := 0; i < n; i++ {
		bi := b[i]
		diff := a[i] > b[i]
		for len(stack) > 0 && stack[len(stack)-1].val < bi {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if top.need {
				need[top.val]++
			}
		}
		if len(stack) > 0 && stack[len(stack)-1].val == bi {
			if diff {
				stack[len(stack)-1].need = true
			}
		} else {
			stack = append(stack, node{val: bi, need: diff})
		}
	}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top.need {
			need[top.val]++
		}
	}
	for x, v := range need {
		if cnt[x] < v {
			return "NO"
		}
	}
	return "YES"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 3
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10) + 1
		delta := rng.Intn(5)
		if rng.Intn(2) == 0 {
			if a[i] > delta {
				b[i] = a[i] - delta
			} else {
				b[i] = a[i]
			}
		} else {
			b[i] = a[i] + delta
		}
	}
	m := rng.Intn(n) + 1
	razors := make([]int, m)
	for i := 0; i < m; i++ {
		razors[i] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i, v := range razors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	input := sb.String()
	expected := solveD(n, a, b, razors)
	return input, expected
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
