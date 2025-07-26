package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type queryD struct {
	u uint64
	s string
}

func solveQuery(n uint64, q queryD) uint64 {
	exp := bits.Len64(n+1) - 1
	H := exp - 1
	path := make([]uint8, 0, H)
	l, r := uint64(1), n
	h := H
	u := q.u
	for {
		mid := (l + r) >> 1
		if u == mid {
			break
		} else if u < mid {
			path = append(path, 0)
			r = mid - 1
		} else {
			path = append(path, 1)
			l = mid + 1
		}
		h--
	}
	curH := h
	for _, c := range q.s {
		switch c {
		case 'L':
			if curH > 0 {
				path = append(path, 0)
				curH--
			}
		case 'R':
			if curH > 0 {
				path = append(path, 1)
				curH--
			}
		case 'U':
			if len(path) > 0 {
				path = path[:len(path)-1]
				curH++
			}
		}
	}
	l, r = 1, n
	for _, b := range path {
		mid := (l + r) >> 1
		if b == 0 {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return (l + r) >> 1
}

func genCaseD(rng *rand.Rand) (string, []uint64) {
	exp := rng.Intn(5) + 2 // exponent 2..6 => n up to 63
	n := uint64(1<<exp) - 1
	q := rng.Intn(5) + 1
	queries := make([]queryD, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < q; i++ {
		u := uint64(rng.Intn(int(n)) + 1)
		l := rng.Intn(15) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			switch rng.Intn(3) {
			case 0:
				b[j] = 'L'
			case 1:
				b[j] = 'R'
			default:
				b[j] = 'U'
			}
		}
		s := string(b)
		queries[i] = queryD{u, s}
		sb.WriteString(fmt.Sprintf("%d\n%s\n", u, s))
	}
	answers := make([]uint64, q)
	for i, qu := range queries {
		answers[i] = solveQuery(n, qu)
	}
	return sb.String(), answers
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCaseD(bin string, in string, exp []uint64) error {
	out, err := runCandidate(bin, []byte(in))
	if err != nil {
		return err
	}
	tokens := strings.Fields(out)
	if len(tokens) != len(exp) {
		return fmt.Errorf("expected %d numbers, got %d", len(exp), len(tokens))
	}
	for i, tok := range tokens {
		val, err := strconv.ParseUint(tok, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %q", tok)
		}
		if val != exp[i] {
			return fmt.Errorf("at line %d expected %d got %d", i+1, exp[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseD(rng)
		if err := runCaseD(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
