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

func buildRef() (string, error) {
	refBin := "./1762G_ref"
	if err := exec.Command("go", "build", "-o", refBin, "1762G.go").Run(); err != nil {
		return "", err
	}
	return refBin, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

type testG struct {
	n   int
	arr []int
}

func genTestG(rng *rand.Rand) testG {
	n := rng.Intn(8) + 3
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	return testG{n: n, arr: arr}
}

func parseOutput(out string, tc testG, expectYes bool) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	first := strings.ToUpper(fields[0])
	if first == "NO" {
		if expectYes {
			return fmt.Errorf("expected YES got NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("first token should be YES or NO")
	}
	if !expectYes {
		return fmt.Errorf("expected NO got YES")
	}
	if len(fields) != 1+tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(fields)-1)
	}
	perm := make([]int, tc.n)
	used := make([]bool, tc.n+1)
	for i := 0; i < tc.n; i++ {
		v, err := strconv.Atoi(fields[1+i])
		if err != nil || v < 1 || v > tc.n || used[v] {
			return fmt.Errorf("invalid permutation")
		}
		used[v] = true
		perm[i] = v
	}
	for i := 1; i < tc.n; i++ {
		if tc.arr[perm[i-1]-1] == tc.arr[perm[i]-1] {
			return fmt.Errorf("adjacent equal values")
		}
	}
	for i := 2; i < tc.n; i++ {
		if perm[i-2] >= perm[i] {
			return fmt.Errorf("order condition failed")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	userBin := os.Args[1]
	refBin, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		tc := genTestG(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		refOut, err := run(refBin, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expectYes := strings.HasPrefix(strings.ToUpper(strings.TrimSpace(refOut)), "YES")
		userOut, err := run(userBin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", t+1, err, userOut)
			os.Exit(1)
		}
		if err := parseOutput(userOut, tc, expectYes); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", t+1, err, input, userOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
