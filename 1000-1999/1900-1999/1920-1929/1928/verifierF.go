package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func genTest() string {
	n := rand.Intn(8) + 1
	m := rand.Intn(8) + 1
	q := rand.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(11)-5)
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d", rand.Intn(11)-5)
		if i+1 < m {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		t := rand.Intn(2) + 1
		if t == 1 {
			l := rand.Intn(n) + 1
			r := l + rand.Intn(n-l+1)
			x := rand.Intn(11) - 5
			fmt.Fprintf(&sb, "1 %d %d %d\n", l, r, x)
		} else {
			l := rand.Intn(m) + 1
			r := l + rand.Intn(m-l+1)
			x := rand.Intn(11) - 5
			fmt.Fprintf(&sb, "2 %d %d %d\n", l, r, x)
		}
	}
	return sb.String()
}

func runCmd(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	rand.Seed(47)
	for i := 0; i < 100; i++ {
		input := []byte(genTest())
		candCmd := exec.Command(os.Args[1])
		candOut, err := runCmd(candCmd, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error: %v\n", err)
			os.Exit(1)
		}
		refCmd := exec.Command("go", "run", "1928F.go")
		refOut, err := runCmd(refCmd, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error: %v\n", err)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Println("WA on test", i+1)
			fmt.Println("input:\n" + string(input))
			fmt.Println("expected:\n" + refOut)
			fmt.Println("got:\n" + candOut)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
