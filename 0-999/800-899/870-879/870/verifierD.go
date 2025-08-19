package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "sync"
)

func solveD(n int, a, b []int) (int, []int) {
	count := 0
	var ans []int
	for pos0 := 0; pos0 < n; pos0++ {
		perm := make([]int, n)
		used := make([]bool, n)
		ok := true
		for i := 0; i < n; i++ {
			v := a[i] ^ pos0
			if v < 0 || v >= n || used[v] {
				ok = false
				break
			}
			perm[i] = v
			used[v] = true
		}
		if !ok {
			continue
		}
		p0 := perm[0]
		inv := make([]int, n)
		for j := 0; j < n; j++ {
			inv[j] = b[j] ^ p0
		}
		for i := 0; i < n && ok; i++ {
			if inv[perm[i]] != i {
				ok = false
			}
		}
		if ok {
			count++
			if ans == nil {
				ans = perm
			}
		}
	}
	return count, ans
}

// Run candidate interactively: respond to lines like "? i j" with a[i]^b[j].
func runCaseInteractive(exe string, n int, a, b []int, expect string) error {
    cmd := exec.Command(exe)
    stdin, err := cmd.StdinPipe()
    if err != nil { return err }
    stdout, err := cmd.StdoutPipe()
    if err != nil { return err }
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    if err := cmd.Start(); err != nil { return err }

    // Feed initial n
    if _, err := fmt.Fprintf(stdin, "%d\n", n); err != nil {
        return fmt.Errorf("failed to write n: %v", err)
    }

    gotBuf := &bytes.Buffer{}
    scanner := bufio.NewScanner(stdout)
    scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        for scanner.Scan() {
            ln := strings.TrimSpace(scanner.Text())
            if strings.HasPrefix(ln, "?") {
                // parse indices
                parts := strings.Fields(ln)
                if len(parts) == 3 {
                    i, _ := strconv.Atoi(parts[1])
                    j, _ := strconv.Atoi(parts[2])
                    // guard indices
                    if i >= 0 && i < n && j >= 0 && j < n {
                        ans := a[i] ^ b[j]
                        fmt.Fprintf(stdin, "%d\n", ans)
                    } else {
                        // invalid query; respond 0 to avoid deadlock
                        fmt.Fprintln(stdin, 0)
                    }
                } else {
                    // malformed; respond 0
                    fmt.Fprintln(stdin, 0)
                }
                continue
            }
            if ln == "" {
                continue
            }
            // record non-query output
            if gotBuf.Len() > 0 { gotBuf.WriteByte('\n') }
            gotBuf.WriteString(ln)
        }
    }()

    wg.Wait()
    if err := cmd.Wait(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    got := strings.TrimSpace(gotBuf.String())
    exp := strings.TrimSpace(expect)
    if got != exp {
        return fmt.Errorf("expected %q got %q", exp, got)
    }
    return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			a[i], _ = strconv.Atoi(scan.Text())
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			b[i], _ = strconv.Atoi(scan.Text())
		}
        cnt, perm := solveD(n, a, b)
        expectBuilder := &strings.Builder{}
        fmt.Fprintln(expectBuilder, "!")
        fmt.Fprintln(expectBuilder, cnt)
        for i := 0; i < n; i++ {
            if i > 0 {
                expectBuilder.WriteByte(' ')
            }
            fmt.Fprintf(expectBuilder, "%d", perm[i])
        }
        expectBuilder.WriteByte('\n')
        if err := runCaseInteractive(exe, n, a, b, expectBuilder.String()); err != nil {
            fmt.Printf("case %d failed: %v\n", caseNum, err)
            os.Exit(1)
        }
	}
	fmt.Println("All tests passed")
}
