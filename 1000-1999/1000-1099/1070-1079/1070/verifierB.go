package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	left, right int
	origWhite   bool
	origBlack   bool
	gotBlack    bool
}

var nodes []Node

func insertTrie(ip uint32, length int, ruleType int) {
	curr := 0
	for i := 0; i < length; i++ {
		bit := (ip >> (31 - i)) & 1
		if bit == 0 {
			if nodes[curr].left == 0 {
				nodes = append(nodes, Node{})
				nodes[curr].left = len(nodes) - 1
			}
			curr = nodes[curr].left
		} else {
			if nodes[curr].right == 0 {
				nodes = append(nodes, Node{})
				nodes[curr].right = len(nodes) - 1
			}
			curr = nodes[curr].right
		}
	}
	if ruleType == 1 {
		nodes[curr].origWhite = true
	} else if ruleType == 2 {
		nodes[curr].origBlack = true
	} else if ruleType == 3 {
		nodes[curr].gotBlack = true
	}
}

func dfs(u int, pWhite, pBlack, pGot bool) error {
	cWhite := pWhite || nodes[u].origWhite
	cBlack := pBlack || nodes[u].origBlack
	cGot := pGot || nodes[u].gotBlack

	if nodes[u].left == 0 {
		if cWhite && cGot {
			return fmt.Errorf("whitelisted IP blocked")
		}
		if cBlack && !cGot {
			return fmt.Errorf("blacklisted IP not blocked")
		}
	} else {
		if err := dfs(nodes[u].left, cWhite, cBlack, cGot); err != nil {
			return err
		}
	}

	if nodes[u].right == 0 {
		if cWhite && cGot {
			return fmt.Errorf("whitelisted IP blocked")
		}
		if cBlack && !cGot {
			return fmt.Errorf("blacklisted IP not blocked")
		}
	} else {
		if err := dfs(nodes[u].right, cWhite, cBlack, cGot); err != nil {
			return err
		}
	}

	return nil
}

func parseIP(s string) (uint32, int, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, "/")
	if len(parts) > 2 {
		return 0, 0, fmt.Errorf("invalid format: %s", s)
	}
	ipParts := strings.Split(parts[0], ".")
	if len(ipParts) != 4 {
		return 0, 0, fmt.Errorf("invalid IP: %s", s)
	}
	var ip uint32
	for i := 0; i < 4; i++ {
		val, err := strconv.Atoi(ipParts[i])
		if err != nil || val < 0 || val > 255 {
			return 0, 0, fmt.Errorf("invalid IP octet: %s", s)
		}
		ip = (ip << 8) | uint32(val)
	}
	length := 32
	if len(parts) == 2 {
		l, err := strconv.Atoi(parts[1])
		if err != nil || l < 0 || l > 32 {
			return 0, 0, fmt.Errorf("invalid length: %s", s)
		}
		length = l
	}
	return ip, length, nil
}

func verify(input string, want string, got string) error {
	wantStr := strings.TrimSpace(want)
	gotStr := strings.TrimSpace(got)

	wantLines := strings.Split(wantStr, "\n")
	if len(wantLines) > 0 && wantLines[0] == "-1" {
		if gotStr != "-1" {
			return fmt.Errorf("expected -1, got something else")
		}
		return nil
	}

	if gotStr == "-1" {
		return fmt.Errorf("expected a solution, got -1")
	}

	gotLines := strings.Split(gotStr, "\n")
	if len(gotLines) == 0 {
		return fmt.Errorf("got empty output")
	}

	wantLen, err := strconv.Atoi(strings.TrimSpace(wantLines[0]))
	if err != nil {
		return fmt.Errorf("failed to parse want length: %v", err)
	}

	gotLen, err := strconv.Atoi(strings.TrimSpace(gotLines[0]))
	if err != nil {
		return fmt.Errorf("failed to parse got length: %v", err)
	}

	if gotLen > wantLen {
		return fmt.Errorf("expected at most %d rules, got %d rules", wantLen, gotLen)
	}

	var gotRules []string
	for i := 1; i < len(gotLines); i++ {
		line := strings.TrimSpace(gotLines[i])
		if line != "" {
			gotRules = append(gotRules, line)
		}
	}
	if len(gotRules) != gotLen {
		return fmt.Errorf("length mismatch: gotLen %d, but found %d rules", gotLen, len(gotRules))
	}

	nodes = []Node{{}}

	inLines := strings.Split(strings.TrimSpace(input), "\n")
	for i := 1; i < len(inLines); i++ {
		line := strings.TrimSpace(inLines[i])
		if line == "" {
			continue
		}
		isWhite := line[0] == '+'
		ip, length, err := parseIP(line[1:])
		if err != nil {
			return fmt.Errorf("invalid input line %s: %v", line, err)
		}
		if isWhite {
			insertTrie(ip, length, 1)
		} else {
			insertTrie(ip, length, 2)
		}
	}

	for _, line := range gotRules {
		ip, length, err := parseIP(line)
		if err != nil {
			return fmt.Errorf("invalid output line %s: %v", line, err)
		}
		insertTrie(ip, length, 3)
	}

	return dfs(0, false, false, false)
}

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1070B.go")
	bin := filepath.Join(os.TempDir(), "1070B_ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runProg(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genSubnet(rng *rand.Rand) string {
	a := rng.Intn(256)
	b := rng.Intn(256)
	c := rng.Intn(256)
	d := rng.Intn(256)
	p := rng.Intn(33)
	ip := uint32(a)<<24 | uint32(b)<<16 | uint32(c)<<8 | uint32(d)
	if p < 32 {
		ip &= ^uint32((1 << (32 - p)) - 1)
	}
	a = int(ip>>24) & 255
	b = int(ip>>16) & 255
	c = int(ip>>8) & 255
	d = int(ip) & 255
	if p == 32 {
		return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
	}
	return fmt.Sprintf("%d.%d.%d.%d/%d", a, b, c, d, p)
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sign := '+'
		if i == 0 || rng.Intn(2) == 0 {
			sign = '-'
		}
		sb.WriteByte(byte(sign))
		sb.WriteString(genSubnet(rng))
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 1000; i++ {
		input := genCase(rng)
		want, err := runProg(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if err := verify(string(input), want, got); err != nil {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\nerror:\n%v\n", i, string(input), want, got, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
