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

type Test struct{ n int64 }

func (t Test) Input() string {
	return fmt.Sprintf("1\n%d\n", t.n)
}

func expected(t Test) string {
	n := t.n
	days := 0
	for (1<<days)-1 < int(n) {
		days++
	}
	days--
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", days))
	rem := n - 1
	last := int64(1)
	for i := days; i >= 1; i-- {
		low := last
		high := last * 2
		ans := high
		for low <= high {
			mid := (low + high) / 2
			mn := int64(i) * mid
			mx := int64((1<<i)-1) * mid
			if mn > rem {
				high = mid - 1
			} else if mx < rem {
				low = mid + 1
			} else {
				ans = mid
				break
			}
		}
		sb.WriteString(fmt.Sprintf("%d", ans-last))
		if i > 1 {
			sb.WriteString(" ")
		}
		rem -= ans
		last = ans
	}
	if days >= 1 {
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func runProg(bin, input string) (string, error) {
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

func genTest(rng *rand.Rand) Test {
	n := rng.Int63n(1_000_000_000-2) + 2
	return Test{n}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
