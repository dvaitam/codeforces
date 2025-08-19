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

// decode applies the inverse of the Right-Left Cipher.
// Given t (the encrypted string), it reconstructs the original s.
func decode(t string) string {
    s := ""
    for len(t) > 0 {
        n := len(t)
        if n%2 == 1 {
            s = string(t[0]) + s
            t = t[1:]
        } else {
            s = string(t[n-1]) + s
            t = t[:n-1]
        }
    }
    return s
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
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

type Case struct {
	t string
}

func genCases() []Case {
	r := rand.New(rand.NewSource(1085))
	cases := make([]Case, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := range cases {
		n := r.Intn(50) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteRune(letters[r.Intn(len(letters))])
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    cand := os.Args[1]
    rand.Seed(time.Now().UnixNano())
    cases := genCases()
    for i, c := range cases {
        input := c.t + "\n"
        want := decode(c.t)
        got, err := run(cand, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\n", i+1, err)
            os.Exit(1)
        }
		if want != got {
			fmt.Printf("case %d failed\ninput: %s\nexpected: %s\ngot: %s\n", i+1, c.t, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
