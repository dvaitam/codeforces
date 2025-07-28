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

func runCandidate(bin, input string) (string, error) {
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(hC, dC, hM, dM, k, w, a int64) string {
	for i := int64(0); i <= k; i++ {
		attack := dC + i*w
		health := hC + (k-i)*a
		turnsHero := (hM + attack - 1) / attack
		turnsMonster := (health + dM - 1) / dM
		if turnsHero <= turnsMonster {
			return "YES"
		}
	}
	return "NO"
}

func generateCase(r *rand.Rand) (string, string) {
	hC := int64(r.Intn(100) + 1)
	dC := int64(r.Intn(100) + 1)
	hM := int64(r.Intn(100) + 1)
	dM := int64(r.Intn(100) + 1)
	k := int64(r.Intn(10))
	w := int64(r.Intn(10) + 1)
	a := int64(r.Intn(10) + 1)
	expect := solveC(hC, dC, hM, dM, k, w, a)
	input := fmt.Sprintf("1\n%d %d\n%d %d\n%d %d %d\n", hC, dC, hM, dM, k, w, a)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
