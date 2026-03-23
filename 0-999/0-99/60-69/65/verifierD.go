package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = `package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type State [4]int

var comb [5][5]int

func compress(mp map[State]struct{}) []State {
	states := make([]State, 0, len(mp))
	if len(mp) == 0 {
		return states
	}
	low := State{1 << 30, 1 << 30, 1 << 30, 1 << 30}
	high := State{-1, -1, -1, -1}
	total := -1
	for st := range mp {
		states = append(states, st)
		sum := 0
		for i := 0; i < 4; i++ {
			if st[i] < low[i] {
				low[i] = st[i]
			}
			if st[i] > high[i] {
				high[i] = st[i]
			}
			sum += st[i]
		}
		if total == -1 {
			total = sum
		}
	}
	var vars []int
	baseSum := 0
	for i := 0; i < 4; i++ {
		if high[i]-low[i] > 1 {
			return states
		}
		if high[i]-low[i] == 1 {
			vars = append(vars, i)
		}
		baseSum += low[i]
	}
	k := total - baseSum
	if k < 0 || k > len(vars) {
		return states
	}
	if comb[len(vars)][k] != len(states) {
		return states
	}
	gen := make([]State, 0, len(states))
	for mask := 0; mask < (1 << len(vars)); mask++ {
		if bits.OnesCount(uint(mask)) != k {
			continue
		}
		ns := low
		for j, idx := range vars {
			if (mask>>j)&1 == 1 {
				ns[idx]++
			}
		}
		if _, ok := mp[ns]; !ok {
			return states
		}
		gen = append(gen, ns)
	}
	return gen
}

func houseIndex(ch byte) int {
	switch ch {
	case 'G':
		return 0
	case 'H':
		return 1
	case 'R':
		return 2
	default:
		return 3
	}
}

func main() {
	for i := 0; i <= 4; i++ {
		comb[i][0] = 1
		comb[i][i] = 1
		for j := 1; j < i; j++ {
			comb[i][j] = comb[i-1][j-1] + comb[i-1][j]
		}
	}

	in := bufio.NewReader(os.Stdin)
	var n int
	var s string
	fmt.Fscan(in, &n, &s)

	states := []State{{0, 0, 0, 0}}

	for i := 0; i < n; i++ {
		next := make(map[State]struct{}, len(states)*4)
		if s[i] == '?' {
			for _, st := range states {
				mn := st[0]
				for j := 1; j < 4; j++ {
					if st[j] < mn {
						mn = st[j]
					}
				}
				for j := 0; j < 4; j++ {
					if st[j] == mn {
						ns := st
						ns[j]++
						next[ns] = struct{}{}
					}
				}
			}
		} else {
			idx := houseIndex(s[i])
			for _, st := range states {
				ns := st
				ns[idx]++
				next[ns] = struct{}{}
			}
		}
		states = compress(next)
	}

	possible := [4]bool{}
	for _, st := range states {
		mn := st[0]
		for j := 1; j < 4; j++ {
			if st[j] < mn {
				mn = st[j]
			}
		}
		for j := 0; j < 4; j++ {
			if st[j] == mn {
				possible[j] = true
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	names := [4]string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"}
	for i := 0; i < 4; i++ {
		if possible[i] {
			fmt.Fprintln(out, names[i])
		}
	}
}
`

func buildRef() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "refbuild")
	if err != nil {
		return "", nil, err
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	if err := os.WriteFile(srcPath, []byte(refSource), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }, nil
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	letters := []byte{'G', 'H', 'R', 'S', '?'}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	s := string(b)
	return fmt.Sprintf("%d\n%s\n", n, s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, cleanup, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		exp, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
