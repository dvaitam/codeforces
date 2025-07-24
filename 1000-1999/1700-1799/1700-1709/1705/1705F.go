package main

import (
    "bufio"
    "fmt"
    "os"
)

const (
    rows = 27
    cols = 25
)

// ask sends a query string to the grader and reads the integer response.
func ask(q []byte, in *bufio.Reader, out *bufio.Writer) int {
    fmt.Fprintln(out, string(q))
    out.Flush()
    var res int
    fmt.Fscan(in, &res)
    return res
}

func askSingle(pos int, n int, in *bufio.Reader, out *bufio.Writer, base int) int {
    q := make([]byte, n)
    for i := range q {
        q[i] = 'F'
    }
    q[pos] = 'T'
    r := ask(q, in, out)
    return (r - base + 1) / 2
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }

    // mapping of questions to pairs (a,b)
    type pair struct{ a, b int }
    ab := make([]pair, n)
    countsA := make([]int, rows)
    countsB := make([]int, cols)
    idx := [rows][cols][]int{}

    skipZeroZero := false
    for i := 0; i < n; i++ {
        a := (i + 1) % rows
        b := (i + 1) % cols
        if a == 0 && b == 0 && !skipZeroZero {
            b = 1 // keep pair (0,0) empty for baseline
            skipZeroZero = true
        }
        ab[i] = pair{a, b}
        idx[a][b] = append(idx[a][b], i)
        countsA[a]++
        countsB[b]++
    }

    // baseline query (all 'F')
    q0 := make([]byte, n)
    for i := range q0 {
        q0[i] = 'F'
    }
    base := ask(q0, in, out)

    // queries for rows
    sumA := make([]int, rows)
    for a := 0; a < rows; a++ {
        q := make([]byte, n)
        for i := 0; i < n; i++ {
            if ab[i].a == a {
                q[i] = 'T'
            } else {
                q[i] = 'F'
            }
        }
        r := ask(q, in, out)
        sumA[a] = (r - base + countsA[a]) / 2
    }

    // queries for columns
    sumB := make([]int, cols)
    for b := 0; b < cols; b++ {
        q := make([]byte, n)
        for i := 0; i < n; i++ {
            if ab[i].b == b {
                q[i] = 'T'
            } else {
                q[i] = 'F'
            }
        }
        r := ask(q, in, out)
        sumB[b] = (r - base + countsB[b]) / 2
    }

    // elimination of pairs using row/column sums
    pairSum := [rows][cols]int{}
    known := [rows][cols]bool{}

    remA := make([]int, rows)
    remB := make([]int, cols)
    leftA := make([]map[int]struct{}, rows)
    leftB := make([]map[int]struct{}, cols)
    for a := 0; a < rows; a++ {
        leftA[a] = map[int]struct{}{}
        for b := 0; b < cols; b++ {
            if len(idx[a][b]) > 0 {
                leftA[a][b] = struct{}{}
                remA[a]++
                if leftB[b] == nil {
                    leftB[b] = map[int]struct{}{}
                }
                leftB[b][a] = struct{}{}
                remB[b]++
            }
        }
    }

    queueA := []int{}
    queueB := []int{}
    for a := 0; a < rows; a++ {
        if remA[a] == 1 {
            queueA = append(queueA, a)
        }
    }
    for b := 0; b < cols; b++ {
        if remB[b] == 1 {
            queueB = append(queueB, b)
        }
    }

    for len(queueA) > 0 || len(queueB) > 0 {
        if len(queueA) > 0 {
            a := queueA[0]
            queueA = queueA[1:]
            if remA[a] != 1 {
                continue
            }
            var b int
            for bb := range leftA[a] {
                b = bb
            }
            pairSum[a][b] = sumA[a]
            known[a][b] = true
            sumA[a] -= pairSum[a][b]
            sumB[b] -= pairSum[a][b]
            remA[a]--
            remB[b]--
            delete(leftA[a], b)
            delete(leftB[b], a)
            if remA[a] == 1 {
                queueA = append(queueA, a)
            }
            if remB[b] == 1 {
                queueB = append(queueB, b)
            }
        } else {
            b := queueB[0]
            queueB = queueB[1:]
            if remB[b] != 1 {
                continue
            }
            var a int
            for aa := range leftB[b] {
                a = aa
            }
            pairSum[a][b] = sumB[b]
            known[a][b] = true
            sumA[a] -= pairSum[a][b]
            sumB[b] -= pairSum[a][b]
            remA[a]--
            remB[b]--
            delete(leftA[a], b)
            delete(leftB[b], a)
            if remA[a] == 1 {
                queueA = append(queueA, a)
            }
            if remB[b] == 1 {
                queueB = append(queueB, b)
            }
        }
    }

    ans := make([]byte, n)
    for i := range ans {
        ans[i] = 'F'
    }

    // resolve remaining pairs
    for a := 0; a < rows; a++ {
        for b := 0; b < cols; b++ {
            if len(idx[a][b]) == 0 {
                continue
            }
            if !known[a][b] {
                // query one element if pair was not deduced
                if len(idx[a][b]) >= 1 {
                    i := idx[a][b][0]
                    val := askSingle(i, n, in, out, base)
                    ans[i] = mapBool(val)
                    if len(idx[a][b]) == 2 {
                        j := idx[a][b][1]
                        pairVal := val
                        ans[j] = mapBool(pairVal)
                    }
                }
                continue
            }
            val := pairSum[a][b]
            if len(idx[a][b]) == 1 {
                ans[idx[a][b][0]] = mapBool(val)
            } else if len(idx[a][b]) == 2 {
                i := idx[a][b][0]
                val1 := askSingle(i, n, in, out, base)
                ans[i] = mapBool(val1)
                ans[idx[a][b][1]] = mapBool(val - val1)
            }
        }
    }

    fmt.Fprintln(out, "!", string(ans))
    out.Flush()
}

func mapBool(v int) byte {
    if v != 0 {
        return 'T'
    }
    return 'F'
}

func remove(m *map[int]struct{}, key int) {
    delete(*m, key)
}

