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

type query struct {
	s byte
	x int
}

func expected(n, q int, arr []int, ops []query) []int {
	for _, op := range ops {
		if op.s == '>' {
			for i := 0; i < n; i++ {
				if arr[i] > op.x {
					arr[i] = -arr[i]
				}
			}
		} else {
			for i := 0; i < n; i++ {
				if arr[i] < op.x {
					arr[i] = -arr[i]
				}
			}
		}
	}
	return arr
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(11) - 5
	}
	ops := make([]query, q)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			ops[i].s = '>'
		} else {
			ops[i].s = '<'
		}
		ops[i].x = rng.Intn(11) - 5
		fmt.Fprintf(&sb, "%c %d\n", ops[i].s, ops[i].x)
	}
	expArr := expected(n, q, append([]int(nil), arr...), ops)
	var sbExp strings.Builder
	for i, v := range expArr {
		if i > 0 {
			sbExp.WriteByte(' ')
		}
		sbExp.WriteString(fmt.Sprint(v))
	}
	return sb.String(), sbExp.String()
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
