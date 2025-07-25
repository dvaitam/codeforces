package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type event struct {
	orig string
	cur  string
}

func solveB(reqs [][2]string) string {
	events := make([]event, 0, len(reqs))
	curIndex := make(map[string]int)
	for _, q := range reqs {
		a, b := q[0], q[1]
		if idx, ok := curIndex[a]; ok {
			events[idx].cur = b
			delete(curIndex, a)
			curIndex[b] = idx
		} else {
			events = append(events, event{orig: a, cur: b})
			curIndex[b] = len(events) - 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(events)))
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("%s %s\n", e.orig, e.cur))
	}
	return sb.String()
}

func genHandle(rng *rand.Rand, used map[string]bool) string {
	for {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for i := 0; i < l; i++ {
			b[i] = byte('a' + rng.Intn(26))
		}
		s := string(b)
		if !used[s] {
			used[s] = true
			return s
		}
	}
}

func genCase(rng *rand.Rand) (string, string) {
	q := rng.Intn(10) + 1
	used := make(map[string]bool)
	active := make([]string, 0)
	// initialize with one user
	first := genHandle(rng, used)
	active = append(active, first)
	reqs := make([][2]string, 0, q)
	for i := 0; i < q; i++ {
		old := active[rng.Intn(len(active))]
		newh := genHandle(rng, used)
		// update active list
		for idx, v := range active {
			if v == old {
				active[idx] = newh
			}
		}
		reqs = append(reqs, [2]string{old, newh})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for _, r := range reqs {
		sb.WriteString(fmt.Sprintf("%s %s\n", r[0], r[1]))
	}
	expect := solveB(reqs)
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := out.String()
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
