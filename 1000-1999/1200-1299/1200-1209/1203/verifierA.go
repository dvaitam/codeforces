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

func solve(arr []int) string {
	n := len(arr)
	idx := -1
	for i, v := range arr {
		if v == 1 {
			idx = i
			break
		}
	}
	ok1 := true
	for i := 0; i < n; i++ {
		if arr[(idx+i)%n] != i+1 {
			ok1 = false
			break
		}
	}
	ok2 := true
	for i := 0; i < n; i++ {
		if arr[(idx-i+n)%n] != i+1 {
			ok2 = false
			break
		}
	}
	if ok1 || ok2 {
		return "YES"
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	arr := make([]string, n)
	for i, v := range perm {
		arr[i] = fmt.Sprintf("%d", v)
	}
	input := fmt.Sprintf("1\n%d\n%s\n", n, strings.Join(arr, " "))
	expected := solve(perm)
	return input, expected
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	gotLines := strings.Split(got, "\n")
	if len(gotLines) > 0 {
		got = strings.TrimSpace(gotLines[0])
	}
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
