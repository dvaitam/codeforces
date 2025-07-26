package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type teleport struct{ a, b int }

func solveCase(n, m int, tele []teleport) string {
	maxReach := 0
	for _, t := range tele {
		if t.a <= maxReach && t.b > maxReach {
			maxReach = t.b
		}
	}
	if maxReach >= m {
		return "YES"
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	m := rng.Intn(100) + 1
	tele := make([]teleport, n)
	for i := range tele {
		a := rng.Intn(m + 1)
		b := a + rng.Intn(m-a+1)
		tele[i] = teleport{a, b}
	}
	sort.Slice(tele, func(i, j int) bool { return tele[i].a < tele[j].a })
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, t := range tele {
		fmt.Fprintf(&sb, "%d %d\n", t.a, t.b)
	}
	return sb.String(), solveCase(n, m, tele)
}

func runCase(bin, input, expected string) error {
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
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	got := strings.ToUpper(fields[0])
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
