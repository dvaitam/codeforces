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

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "571C.go")
	bin := filepath.Join(os.TempDir(), "oracle571C.bin")
	
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func verifyAssignment(input, output string) error {
    lines := strings.Fields(input)
    if len(lines) == 0 { return nil }
    nc, _ := strconv.Atoi(lines[0])
    m, _ := strconv.Atoi(lines[1])
    
    cursor := 2
    clauses := make([][]int, nc)
    for i := 0; i < nc; i++ {
        k, _ := strconv.Atoi(lines[cursor])
        cursor++
        for j := 0; j < k; j++ {
            lit, _ := strconv.Atoi(lines[cursor])
            clauses[i] = append(clauses[i], lit)
            cursor++
        }
    }
    
    outFields := strings.Fields(output)
    if len(outFields) == 0 {
        return fmt.Errorf("no output")
    }
    
    if outFields[0] != "YES" {
        return fmt.Errorf("expected YES, got %s", outFields[0])
    }
    if len(outFields) < 2 {
        return fmt.Errorf("missing assignment string")
    }
    assignStr := outFields[1]
    if len(assignStr) != m {
        return fmt.Errorf("assignment length %d != variables %d", len(assignStr), m)
    }
    
    vals := make([]int, m+1)
    for i, ch := range assignStr {
        if ch == '1' {
            vals[i+1] = 1
        } else if ch == '0' {
            vals[i+1] = 0
        } else {
            return fmt.Errorf("invalid char in assignment: %c", ch)
        }
    }
    
    for i, cl := range clauses {
        sat := false
        for _, lit := range cl {
            v := abs(lit)
            val := vals[v]
            if (lit > 0 && val == 1) || (lit < 0 && val == 0) {
                sat = true
                break
            }
        }
        if !sat {
            return fmt.Errorf("clause %d %v not satisfied by %s", i+1, cl, assignStr)
        }
    }
    return nil
}

func runCase(bin, oracle, input string) error {
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle runtime error: %v", err)
	}
	oracleOut := strings.TrimSpace(outO.String())
    oracleVerdict := ""
    if len(oracleOut) > 0 {
        oracleVerdict = strings.Fields(oracleOut)[0]
    }

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	got := strings.TrimSpace(out.String())
    gotFields := strings.Fields(got)
    if len(gotFields) == 0 {
        return fmt.Errorf("no output from solution")
    }
    gotVerdict := gotFields[0]
    
    if oracleVerdict == "NO" {
        if gotVerdict != "NO" {
            return fmt.Errorf("oracle says NO, but solution says %s", gotVerdict)
        }
        return nil
    }
    
    if gotVerdict == "NO" {
        return fmt.Errorf("oracle says YES, but solution says NO")
    }
    
    if err := verifyAssignment(input, got); err != nil {
        return fmt.Errorf("invalid assignment: %v", err)
    }
    
	return nil
}

func genCase(rng *rand.Rand) string {
    for {
        nc := rng.Intn(10) + 1 
        minN := (nc + 1) / 2
        n := rng.Intn(10) + minN 
        
        clauses := make([][]int, nc)
        
        for i := 1; i <= n; i++ {
            occ := rng.Intn(3) 
            if occ > nc { occ = nc }
            
            if occ > 0 {
                cIndices := rng.Perm(nc)[:occ]
                for _, ci := range cIndices {
                    val := i
                    if rng.Intn(2) == 0 { val = -val }
                    clauses[ci] = append(clauses[ci], val)
                }
            }
        }
        
        valid := true
        for _, cl := range clauses {
            if len(cl) == 0 {
                valid = false
                break
            }
        }
        
        if valid {
            var sb strings.Builder
            fmt.Fprintf(&sb, "%d %d\n", nc, n)
            for _, cl := range clauses {
                fmt.Fprintf(&sb, "%d", len(cl))
                for _, v := range cl {
                    fmt.Fprintf(&sb, " %d", v)
                }
                sb.WriteByte('\n')
            }
            return sb.String()
        }
    }
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}