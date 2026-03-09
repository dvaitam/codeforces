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

func solveB(a1, a2, a3, a4 int) string {
	type config struct {
		start      byte
		seg4, seg7 int
	}
	var cfgs []config
	switch {
	case a3 == a4:
		k := a3
		// try start='4' first (smaller number); fall back to start='7'
		cfgs = []config{
			{'4', k + 1, k},
			{'7', k, k + 1},
		}
	case a3 == a4+1:
		k := a4
		cfgs = []config{{'4', k + 1, k + 1}}
	case a4 == a3+1:
		k := a3
		cfgs = []config{{'7', k + 1, k + 1}}
	default:
		return "-1"
	}
	for _, cfg := range cfgs {
		if a1 < cfg.seg4 || a2 < cfg.seg7 {
			continue
		}
		size4First := a1 - (cfg.seg4 - 1)
		size7Last := a2 - (cfg.seg7 - 1)
		var b strings.Builder
		b.Grow(a1 + a2)
		idx4, idx7 := 0, 0
		current := cfg.start
		for i := 0; i < cfg.seg4+cfg.seg7; i++ {
			if current == '4' {
				idx4++
				cnt := 1
				if idx4 == 1 {
					cnt = size4First
				}
				b.WriteString(strings.Repeat("4", cnt))
				current = '7'
			} else {
				idx7++
				cnt := 1
				if idx7 == cfg.seg7 {
					cnt = size7Last
				}
				b.WriteString(strings.Repeat("7", cnt))
				current = '4'
			}
		}
		return b.String()
	}
	return "-1"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	// problem requires 1 <= a1,a2,a3,a4
	a1 := rng.Intn(5) + 1
	a2 := rng.Intn(5) + 1
	a3 := rng.Intn(5) + 1
	a4 := rng.Intn(5) + 1
	input := fmt.Sprintf("%d %d %d %d\n", a1, a2, a3, a4)
	return input, solveB(a1, a2, a3, a4)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
