package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n     int
	k     int
	a     int
	shots []int
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(50) + 5
	a := rng.Intn(n/2) + 1
	maxShips := (n + 1) / (a + 1)
	k := rng.Intn(maxShips) + 1
	m := rng.Intn(50) + 1
	shots := make([]int, m)
	for i := 0; i < m; i++ {
		shots[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n%d\n", n, k, a, m))
	for i, s := range shots {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(s))
	}
	sb.WriteByte('\n')
	return Test{n: n, k: k, a: a, shots: shots, input: sb.String()}
}

func calc(len, a int) int {
	if len < 0 {
		return 0
	}
	return (len + 1) / (a + 1)
}

// solver replicates 567D.go
func solve(t Test) string {
	// use ordered set as []int slice for simplicity
	cuts := []int{0, t.n + 1}
	total := calc(t.n, t.a)
	for i, x := range t.shots {
		// find position to insert
		idx := 0
		for idx < len(cuts) && cuts[idx] < x {
			idx++
		}
		prev := cuts[idx-1]
		next := cuts[idx]
		total -= calc(next-prev-1, t.a)
		total += calc(x-prev-1, t.a)
		total += calc(next-x-1, t.a)
		// insert x maintaining order if not already present
		inserted := false
		if cuts[idx-1] != x && cuts[idx] != x {
			cuts = append(cuts, 0)
			copy(cuts[idx+1:], cuts[idx:])
			cuts[idx] = x
			inserted = true
		}
		if total < t.k {
			return strconv.Itoa(i + 1)
		}
		if !inserted {
			// duplicate shot, but still after result check
		}
	}
	return "-1"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
