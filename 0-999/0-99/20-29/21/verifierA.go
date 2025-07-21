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

func isValidToken(s string, minLen, maxLen int) bool {
	n := len(s)
	if n < minLen || n > maxLen {
		return false
	}
	for i := 0; i < n; i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			continue
		}
		return false
	}
	return true
}

func isValidJabberID(s string) bool {
	if strings.Count(s, "@") != 1 {
		return false
	}
	parts := strings.SplitN(s, "@", 2)
	user := parts[0]
	rest := parts[1]
	if !isValidToken(user, 1, 16) {
		return false
	}
	if strings.Count(rest, "/") > 1 {
		return false
	}
	var host, res string
	if idx := strings.Index(rest, "/"); idx >= 0 {
		host = rest[:idx]
		res = rest[idx+1:]
		if !isValidToken(res, 1, 16) {
			return false
		}
	} else {
		host = rest
	}
	if len(host) < 1 || len(host) > 32 {
		return false
	}
	labels := strings.Split(host, ".")
	for _, lbl := range labels {
		if !isValidToken(lbl, 1, 16) {
			return false
		}
	}
	return true
}

func randomToken(rng *rand.Rand, minLen, maxLen int) string {
	l := rng.Intn(maxLen-minLen+1) + minLen
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		typ := rng.Intn(3)
		switch typ {
		case 0:
			b[i] = byte('a' + rng.Intn(26))
		case 1:
			b[i] = byte('A' + rng.Intn(26))
		default:
			b[i] = byte('0' + rng.Intn(10))
		}
	}
	return string(b)
}

func generateValid(rng *rand.Rand) string {
	user := randomToken(rng, 1, 16)
	// host labels
	var hostParts []string
	for {
		part := randomToken(rng, 1, 16)
		if len(strings.Join(append(hostParts, part), ".")) > 32 {
			continue
		}
		hostParts = append(hostParts, part)
		if len(strings.Join(hostParts, ".")) >= rng.Intn(32-1)+1 || rng.Float64() < 0.3 {
			break
		}
	}
	host := strings.Join(hostParts, ".")
	res := ""
	if rng.Float64() < 0.5 {
		res = randomToken(rng, 1, 16)
	}
	if res != "" {
		return user + "@" + host + "/" + res
	}
	return user + "@" + host
}

func mutateInvalid(rng *rand.Rand, s string) string {
	choice := rng.Intn(6)
	switch choice {
	case 0:
		// remove @
		return strings.ReplaceAll(s, "@", "")
	case 1:
		// add extra @
		return s + "@" + randomToken(rng, 1, 5)
	case 2:
		// invalid char in username
		return "#" + s
	case 3:
		// host too long
		return strings.Split(s, "@")[0] + "@" + strings.Repeat("a", 33)
	case 4:
		// invalid resource char
		return s + "/" + "!" + randomToken(rng, 1, 3)
	default:
		return s + "//" + randomToken(rng, 1, 3)
	}
}

func generateCase(rng *rand.Rand) (string, string) {
	valid := rng.Float64() < 0.5
	id := generateValid(rng)
	if !valid {
		id = mutateInvalid(rng, id)
	}
	expected := "NO"
	if isValidJabberID(id) {
		expected = "YES"
	}
	return id + "\n", expected
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
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
