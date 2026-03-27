package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const embeddedSolver611E = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func getCaps(S int, A, B, C int) []int {
	var caps []int
	for sub := 1; sub <= 7; sub++ {
		if (S&sub) == sub {
			sum := 0
			if (sub&1) != 0 {
				sum += A
			}
			if (sub&2) != 0 {
				sum += B
			}
			if (sub&4) != 0 {
				sum += C
			}
			caps = append(caps, sum)
		}
	}
	sort.Slice(caps, func(i, j int) bool {
		return caps[i] > caps[j]
	})
	return caps
}

func compareCaps(c1, c2 []int) int {
	for i := 0; i < len(c1) && i < len(c2); i++ {
		if c1[i] != c2[i] {
			if c1[i] < c2[i] {
				return -1
			}
			return 1
		}
	}
	if len(c1) < len(c2) {
		return -1
	} else if len(c1) > len(c2) {
		return 1
	}
	return 0
}

type StateRank struct {
	S    int
	caps []int
}

type Transition struct {
	S   int
	R   int
	cap int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	var a, b, c int
	fmt.Fscan(reader, &a, &b, &c)

	criminals := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &criminals[i])
	}

	m := []int{a, b, c}
	sort.Ints(m)
	A, B, C := m[0], m[1], m[2]

	sort.Slice(criminals, func(i, j int) bool {
		return criminals[i] > criminals[j]
	})

	if criminals[0] > A+B+C {
		fmt.Println("-1")
		return
	}

	ranks := make([]int, 8)
	var states []StateRank
	for S := 0; S < 8; S++ {
		states = append(states, StateRank{S, getCaps(S, A, B, C)})
	}
	sort.Slice(states, func(i, j int) bool {
		return compareCaps(states[i].caps, states[j].caps) < 0
	})
	for i, sr := range states {
		ranks[sr.S] = i
	}

	var transitions []Transition
	for S := 1; S < 8; S++ {
		for R := 0; R < S; R++ {
			if (S&R) == R {
				cap := 0
				diff := S ^ R
				if (diff&1) != 0 {
					cap += A
				}
				if (diff&2) != 0 {
					cap += B
				}
				if (diff&4) != 0 {
					cap += C
				}
				transitions = append(transitions, Transition{S, R, cap})
			}
		}
	}

	sort.Slice(transitions, func(i, j int) bool {
		t1, t2 := transitions[i], transitions[j]
		rS1 := ranks[t1.S]
		rS2 := ranks[t2.S]
		if rS1 != rS2 {
			return rS1 < rS2
		}
		rR1 := ranks[t1.R]
		rR2 := ranks[t2.R]
		if rR1 != rR2 {
			return rR1 > rR2
		}
		return t1.cap < t2.cap
	})

	check := func(H int) bool {
		c_counts := make([]int, 8)
		c_counts[7] = H
		for _, X := range criminals {
			found := false
			for _, t := range transitions {
				if c_counts[t.S] > 0 && t.cap >= X {
					c_counts[t.S]--
					c_counts[t.R]++
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}

	low := 1
	high := n
	ans := -1
	for low <= high {
		mid := (low + high) / 2
		if check(mid) {
			ans = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	fmt.Println(ans)
}
`

func buildEmbeddedSolver() (string, func(), error) {
	tmpSrc, err := os.CreateTemp("", "solver611E-*.go")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpSrc.WriteString(embeddedSolver611E); err != nil {
		tmpSrc.Close()
		os.Remove(tmpSrc.Name())
		return "", nil, err
	}
	tmpSrc.Close()

	tmpBin, err := os.CreateTemp("", "solver611E-bin-*")
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

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func runCandidate(bin string, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, cleanup, err := buildEmbeddedSolver()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build embedded solver: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(5))
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		strengths := [3]int{rng.Intn(50) + 1, rng.Intn(50) + 1, rng.Intn(50) + 1}
		criminals := make([]int, n)
		for i := 0; i < n; i++ {
			criminals[i] = rng.Intn(50) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		fmt.Fprintf(&sb, "%d %d %d\n", strengths[0], strengths[1], strengths[2])
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", criminals[i])
		}
		sb.WriteByte('\n')
		input := sb.String()

		want, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", t+1, err)
			os.Exit(1)
		}
		wantVal, _ := strconv.Atoi(want)

		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != wantVal {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", t+1, wantVal, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
