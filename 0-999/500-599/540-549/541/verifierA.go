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

const reference = "541A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		vids := make([][2]int, n)
		for j := 0; j < n; j++ {
			l := rng.Intn(10)
			r := l + rng.Intn(5)
			vids[j] = [2]int{l, r}
		}
		chans := make([]struct{ a, b, c int }, m)
		for j := 0; j < m; j++ {
			a := rng.Intn(10)
			b := a + rng.Intn(5)
			c := rng.Intn(5) + 1
			chans[j] = struct{ a, b, c int }{a, b, c}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", vids[j][0], vids[j][1]))
		}
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", chans[j].a, chans[j].b, chans[j].c))
		}
		input := sb.String()
		exp, err := run(reference, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
