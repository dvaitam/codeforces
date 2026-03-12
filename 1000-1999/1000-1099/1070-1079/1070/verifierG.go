package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func buildRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		_, cur, _, _ := runtime.Caller(0)
		dir := filepath.Dir(cur)
		src = filepath.Join(dir, "1070G.go")
	}
	bin := filepath.Join(os.TempDir(), "1070G_ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runProg(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) ([]byte, int, int, []int, []int, []int) {
	n := rng.Intn(5) + 1
	m := rng.Intn(n) + 1
	// Generate board first
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(11) - 5
	}
	// Pick m distinct positions for heroes and ensure they are blank
	perm := rng.Perm(n)
	starts := make([]int, m)
	hps := make([]int, m)
	for i := 0; i < m; i++ {
		a[perm[i]] = 0 // ensure hero positions are blank
		starts[i] = perm[i] + 1
		hps[i] = rng.Intn(5) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", starts[i], hps[i]))
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	return []byte(sb.String()), n, m, starts, hps, a
}

// simulate checks if a given rally point and hero order is valid.
func simulate(n, m int, starts, hps []int, a []int, rally int, order []int) bool {
	board := make([]int, n+1)
	for i := 1; i <= n; i++ {
		board[i] = a[i-1]
	}
	for _, idx := range order {
		hp := hps[idx]
		pos := starts[idx]
		if pos < rally {
			for p := pos + 1; p <= rally; p++ {
				v := board[p]
				if v < 0 {
					if -v > hp {
						return false
					}
					hp += v
					board[p] = 0
				} else {
					hp += v
					board[p] = 0
				}
			}
		} else if pos > rally {
			for p := pos - 1; p >= rally; p-- {
				v := board[p]
				if v < 0 {
					if -v > hp {
						return false
					}
					hp += v
					board[p] = 0
				} else {
					hp += v
					board[p] = 0
				}
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, n, m, starts, hps, a := genCase(rng)

		// Run reference to know if a solution exists
		refOut, err := runProg(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}

		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}

		refFields := strings.Fields(refOut)
		gotFields := strings.Fields(got)

		refIsNeg1 := len(refFields) > 0 && refFields[0] == "-1"

		if len(gotFields) == 0 {
			fmt.Printf("wrong answer on test %d: empty output\ninput:\n%s", i, string(input))
			os.Exit(1)
		}

		// If candidate says -1, check reference agrees
		if gotFields[0] == "-1" {
			if !refIsNeg1 {
				fmt.Printf("wrong answer on test %d: candidate says -1 but solution exists\ninput:\n%sexpected:\n%s\n", i, string(input), refOut)
				os.Exit(1)
			}
			continue
		}

		// Candidate gives an answer - validate it structurally regardless of reference

		if len(gotFields) != m+1 {
			fmt.Printf("wrong answer on test %d: expected %d tokens, got %d\ninput:\n%sgot:\n%s\n", i, m+1, len(gotFields), string(input), got)
			os.Exit(1)
		}

		rally, err := strconv.Atoi(gotFields[0])
		if err != nil || rally < 1 || rally > n {
			fmt.Printf("wrong answer on test %d: invalid rally point %s\ninput:\n%s", i, gotFields[0], string(input))
			os.Exit(1)
		}

		order := make([]int, m)
		seen := make(map[int]bool)
		for j := 0; j < m; j++ {
			v, err := strconv.Atoi(gotFields[j+1])
			if err != nil || v < 1 || v > m {
				fmt.Printf("wrong answer on test %d: invalid hero index %s\ninput:\n%s", i, gotFields[j+1], string(input))
				os.Exit(1)
			}
			if seen[v] {
				fmt.Printf("wrong answer on test %d: duplicate hero %d\ninput:\n%s", i, v, string(input))
				os.Exit(1)
			}
			seen[v] = true
			order[j] = v - 1
		}

		if !simulate(n, m, starts, hps, a, rally, order) {
			fmt.Printf("wrong answer on test %d: heroes don't all survive\ninput:\n%sgot:\n%s\n", i, string(input), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
