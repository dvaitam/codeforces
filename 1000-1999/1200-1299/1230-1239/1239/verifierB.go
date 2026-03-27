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

const embeddedSolver1239B = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)

	if n%2 != 0 {
		fmt.Println(0)
		fmt.Println("1 1")
		return
	}

	open := 0
	close := 0
	for i := 0; i < n; i++ {
		if s[i] == '(' {
			open++
		} else {
			close++
		}
	}
	if open != close {
		fmt.Println(0)
		fmt.Println("1 1")
		return
	}

	cur := 0
	minPref := 0
	minIdx := 0
	for i := 0; i < n; i++ {
		if s[i] == '(' {
			cur++
		} else {
			cur--
		}
		if cur < minPref {
			minPref = cur
			minIdx = i
		}
	}

	shift := (minIdx + 1) % n
	t := make([]byte, n)
	for i := 0; i < n; i++ {
		t[i] = s[(shift+i)%n]
	}

	S := make([]int, n+1)
	S[0] = 0
	for i := 0; i < n; i++ {
		if t[i] == '(' {
			S[i+1] = S[i] + 1
		} else {
			S[i+1] = S[i] - 1
		}
	}

	c0 := 0
	for i := 1; i <= n; i++ {
		if S[i] == 0 {
			c0++
		}
	}

	maxBeauty := c0
	bestL := 1
	bestR := 1

	prevZero := 0
	for i := 1; i <= n; i++ {
		if S[i] == 0 {
			L := prevZero + 1
			R := i
			cnt1 := 0
			for j := L; j <= R; j++ {
				if S[j] == 1 {
					cnt1++
				}
			}
			if cnt1 > maxBeauty {
				maxBeauty = cnt1
				bestL = L
				bestR = R
			}
			prevZero = i
		}
	}

	prevLeq1 := 0
	for i := 1; i <= n; i++ {
		if S[i] <= 1 {
			if prevLeq1 > 0 && S[prevLeq1] == 1 && S[i] == 1 && i > prevLeq1+1 {
				L := prevLeq1 + 1
				R := i
				cnt2 := 0
				for j := L; j <= R; j++ {
					if S[j] == 2 {
						cnt2++
					}
				}
				if c0+cnt2 > maxBeauty {
					maxBeauty = c0 + cnt2
					bestL = L
					bestR = R
				}
			}
			prevLeq1 = i
		}
	}

	origL := (shift + bestL - 1) % n + 1
	origR := (shift + bestR - 1) % n + 1

	fmt.Println(maxBeauty)
	fmt.Printf("%d %d\n", origL, origR)
}
`

func buildEmbeddedSolver() (string, func(), error) {
	tmpSrc, err := os.CreateTemp("", "solver1239B-*.go")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpSrc.WriteString(embeddedSolver1239B); err != nil {
		tmpSrc.Close()
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpSrc.Close()

	tmpBin, err := os.CreateTemp("", "solver1239B-bin-*")
	if err != nil {
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpBin.Close()

	cmd := exec.Command("go", "build", "-o", tmpBin.Name(), tmpSrc.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpSrc.Name())
		os.Remove(tmpBin.Name())
		return "", nil, fmt.Errorf("build solver: %v\n%s", err, out)
	}
	os.Remove(tmpSrc.Name())
	return tmpBin.Name(), func() { os.Remove(tmpBin.Name()) }, nil
}

func run(bin string, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	var b strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b.WriteByte('(')
		} else {
			b.WriteByte(')')
		}
	}
	return fmt.Sprintf("%d\n%s\n", n, b.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, cleanup, err := buildEmbeddedSolver()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build embedded solver: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)

		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
