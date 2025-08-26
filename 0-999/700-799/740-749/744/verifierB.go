package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

type Matrix [][]int

func generateMatrix(rng *rand.Rand) Matrix {
	n := rng.Intn(3) + 2 // 2..4
	m := make(Matrix, n)
	for i := 0; i < n; i++ {
		m[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				m[i][j] = 0
			} else {
				m[i][j] = rng.Intn(10)
			}
		}
	}
	return m
}

func expectedRowMins(mat Matrix) []int {
	n := len(mat)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		best := math.MaxInt32
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if mat[i][j] < best {
				best = mat[i][j]
			}
		}
		res[i] = best
	}
	return res
}

func runCase(bin string, mat Matrix) error {
    n := len(mat)
    var stderr bytes.Buffer
    cmd := exec.Command(bin)
    stdin, _ := cmd.StdinPipe()
    stdout, _ := cmd.StdoutPipe()
    cmd.Stderr = &stderr

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start: %v", err)
    }

    // Send only n as initial input (interactive protocol)
    fmt.Fprintf(stdin, "%d\n", n)

    reader := bufio.NewReader(stdout)

    // Helper to read next integer token from candidate output
    readInt := func() (int, error) {
        var x int
        _, err := fmt.Fscan(reader, &x)
        return x, err
    }

    for {
        t, err := readInt()
        if err != nil {
            cmd.Wait()
            return fmt.Errorf("read error: %v stderr:%s", err, stderr.String())
        }
        if t == -1 {
            // Final answers follow
            ans := make([]int, n)
            for i := 0; i < n; i++ {
                v, err := readInt()
                if err != nil {
                    cmd.Wait()
                    return fmt.Errorf("failed to read final answers: %v stderr:%s", err, stderr.String())
                }
                ans[i] = v
            }
            expect := expectedRowMins(mat)
            for i := 0; i < n; i++ {
                if ans[i] != expect[i] {
                    cmd.Wait()
                    return fmt.Errorf("row %d expected %d got %d", i+1, expect[i], ans[i])
                }
            }
            break
        }
        // Read t indices (1-based)
        idx := make([]int, t)
        for i := 0; i < t; i++ {
            v, err := readInt()
            if err != nil {
                cmd.Wait()
                return fmt.Errorf("failed to read indices: %v stderr:%s", err, stderr.String())
            }
            if v < 1 || v > n {
                cmd.Wait()
                return fmt.Errorf("index out of range: %d", v)
            }
            idx[i] = v - 1
        }
        // Compute and send responses: for each k, min over j in idx of mat[k][j]
        var sb strings.Builder
        for k := 0; k < n; k++ {
            best := math.MaxInt32
            for _, j := range idx {
                if mat[k][j] < best {
                    best = mat[k][j]
                }
            }
            if k > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(strconv.Itoa(best))
        }
        sb.WriteByte('\n')
        if _, err := stdin.Write([]byte(sb.String())); err != nil {
            cmd.Wait()
            return fmt.Errorf("failed to write response: %v", err)
        }
    }

    if err := cmd.Wait(); err != nil {
        return fmt.Errorf("runtime error: %v stderr:%s", err, stderr.String())
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
	for t := 1; t <= 100; t++ {
		mat := generateMatrix(rng)
		if err := runCase(bin, mat); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
