package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type item struct {
	length int64
	p      int64
}

type byOrder []item

func (a byOrder) Len() int      { return len(a) }
func (a byOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byOrder) Less(i, j int) bool {
	ai, aj := a[i], a[j]
	left := ai.p * ai.length * (100 - aj.p)
	right := aj.p * aj.length * (100 - ai.p)
	return left > right
}

func solveCase(arr []item) string {
	sort.Sort(byOrder(arr))
	var S, t float64
	for _, it := range arr {
		S += 10000.0*float64(it.length) + t*(100.0-float64(it.p))
		t += float64(it.p) * float64(it.length)
	}
	result := S / 10000.0
	return fmt.Sprintf("%.10f", result)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	arr := make([]item, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		length := int64(rng.Intn(986) + 15)
		p := int64(rng.Intn(101))
		arr[i] = item{length: length, p: p}
		sb.WriteString(fmt.Sprintf("%d %d\n", length, p))
	}
	expected := solveCase(arr)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
