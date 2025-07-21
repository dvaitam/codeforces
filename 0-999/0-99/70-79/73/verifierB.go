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

type Racer struct {
	name string
	a    int
}

func rankRange(racers []Racer, b []int, vasyaName string) (int, int) {
	n := len(racers)
	var vasyaA int
	others := make([]Racer, 0, n-1)
	for _, r := range racers {
		if r.name == vasyaName {
			vasyaA = r.a
		} else {
			others = append(others, r)
		}
	}
	total := n
	B := make([]int, 0, total)
	for i := 0; i < len(b); i++ {
		B = append(B, b[i])
	}
	for i := 0; i < total-len(b); i++ {
		B = append(B, 0)
	}
	sort.Slice(B, func(i, j int) bool { return B[i] > B[j] })
	bestB := append([]int(nil), B[1:]...)
	threshBest := vasyaA + B[0]
	worstB := append([]int(nil), B[:len(B)-1]...)
	threshWorst := vasyaA + B[len(B)-1]
	sort.Slice(others, func(i, j int) bool {
		if others[i].a != others[j].a {
			return others[i].a < others[j].a
		}
		return others[i].name < others[j].name
	})
	threatsBest := 0
	for i, r := range others {
		sum := r.a + bestB[i]
		if sum > threshBest || (sum == threshBest && r.name < vasyaName) {
			threatsBest++
		}
	}
	sort.Slice(others, func(i, j int) bool {
		if others[i].a != others[j].a {
			return others[i].a > others[j].a
		}
		return others[i].name > others[j].name
	})
	threatsWorst := 0
	for i, r := range others {
		sum := r.a + worstB[i]
		if sum > threshWorst || (sum == threshWorst && r.name < vasyaName) {
			threatsWorst++
		}
	}
	return 1 + threatsBest, 1 + threatsWorst
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 2
		racers := make([]Racer, n)
		for j := 0; j < n; j++ {
			racers[j] = Racer{
				name: fmt.Sprintf("p%d", j),
				a:    rng.Intn(30),
			}
		}
		m := rng.Intn(n) + 1
		b := make([]int, m)
		for j := 0; j < m; j++ {
			b[j] = rng.Intn(20)
		}
		vasyaIndex := rng.Intn(n)
		vasyaName := racers[vasyaIndex].name
		best, worst := rankRange(racers, b, vasyaName)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, r := range racers {
			sb.WriteString(fmt.Sprintf("%s %d\n", r.name, r.a))
		}
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(vasyaName)
		sb.WriteByte('\n')
		input := sb.String()
		expected := fmt.Sprintf("%d %d", best, worst)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
