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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase() string {
	n := rand.Intn(5) + 1
	k := rand.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		v := rand.Intn(7) - 3
		if v == 0 {
			v = 1
		}
		sb.WriteString(fmt.Sprintf("%d ", v))
	}
	sb.WriteByte('\n')
	q := rand.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		typ := rand.Intn(2) + 1
		if typ == 1 {
			idx := rand.Intn(n) + 1
			val := rand.Intn(7) - 3
			if val == 0 {
				val = 1
			}
			sb.WriteString(fmt.Sprintf("1 %d %d\n", idx, val))
		} else {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := generateCase()
		exp, err := runBinary("1340F.go", input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed running official solution:", err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
