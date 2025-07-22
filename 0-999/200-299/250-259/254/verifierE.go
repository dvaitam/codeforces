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

type friend struct{ l, r, f int }

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
	return out.String(), nil
}

func expectedRating(input string) (int, error) {
	out, err := run("254E.go", input)
	if err != nil {
		return 0, err
	}
	var rating int
	if _, err := fmt.Fscan(strings.NewReader(out), &rating); err != nil {
		return 0, err
	}
	return rating, nil
}

func simulate(n, v int, a []int, friends []friend, schedule [][]int) (int, error) {
	carry := 0
	rating := 0
	for day := 0; day < n; day++ {
		avail := a[day] + carry
		used := make(map[int]bool)
		for _, idx := range schedule[day] {
			if idx < 1 || idx > len(friends) {
				return 0, fmt.Errorf("bad friend index")
			}
			if used[idx] {
				return 0, fmt.Errorf("duplicate friend")
			}
			used[idx] = true
			f := friends[idx-1]
			if day+1 < f.l || day+1 > f.r {
				return 0, fmt.Errorf("feeding outside interval")
			}
			avail -= f.f
			rating++
		}
		avail -= v
		if avail < 0 {
			return 0, fmt.Errorf("not enough food")
		}
		if avail > v {
			avail = v
		}
		carry = avail
	}
	return rating, nil
}

func generateCase(rng *rand.Rand) (string, int, []int, []friend) {
	n := rng.Intn(4) + 1
	v := rng.Intn(5) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = v + rng.Intn(5)
	}
	m := rng.Intn(3) + 1
	friends := make([]friend, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		f := rng.Intn(v) + 1
		friends[i] = friend{l, r, f}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, v))
	for i, x := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for _, fr := range friends {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", fr.l, fr.r, fr.f))
	}
	return sb.String(), v, a, friends
}

func parseSchedule(out string, n int) (int, [][]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return 0, nil, fmt.Errorf("no output")
	}
	rating, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, nil, fmt.Errorf("bad rating")
	}
	pos := 1
	schedule := make([][]int, n)
	for i := 0; i < n; i++ {
		if pos >= len(tokens) {
			return 0, nil, fmt.Errorf("missing day %d", i+1)
		}
		k, err := strconv.Atoi(tokens[pos])
		if err != nil {
			return 0, nil, fmt.Errorf("bad k")
		}
		pos++
		if pos+k > len(tokens) {
			return 0, nil, fmt.Errorf("missing indices")
		}
		day := make([]int, k)
		for j := 0; j < k; j++ {
			val, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return 0, nil, fmt.Errorf("bad index")
			}
			day[j] = val
			pos++
		}
		schedule[i] = day
	}
	if pos != len(tokens) {
		return 0, nil, fmt.Errorf("extra tokens")
	}
	return rating, schedule, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, v, a, friends := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		rating, schedule, err := parseSchedule(out, len(a))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
		simRating, err := simulate(len(a), v, a, friends, schedule)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid schedule: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
		if simRating != rating {
			fmt.Fprintf(os.Stderr, "case %d rating mismatch: reported %d actual %d\ninput:\n%soutput:\n%s", i+1, rating, simRating, input, out)
			os.Exit(1)
		}
		exp, err := expectedRating(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d internal solver error: %v", i+1, err)
			os.Exit(1)
		}
		if rating != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected rating %d got %d\ninput:\n%soutput:\n%s", i+1, exp, rating, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
