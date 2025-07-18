package main

import (
    "bufio"
    "fmt"
    "os"
)

type Pair struct{
    file int
    idx int
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, m int
    if _, err := fmt.Fscan(in, &n, &m); err != nil {
        return
    }
    a := make([][]int, m)
    totalPos := make([]bool, n+1)
    for i := 0; i < m; i++ {
        var ni int
        fmt.Fscan(in, &ni)
        a[i] = make([]int, ni)
        for j := 0; j < ni; j++ {
            fmt.Fscan(in, &a[i][j])
            totalPos[a[i][j]] = true
        }
    }
    // find free cluster
    freePos := -1
    for i := 1; i <= n; i++ {
        if !totalPos[i] {
            freePos = i
            break
        }
    }
    if freePos == -1 {
        // no free pos (should not happen due to constraints)
        freePos = 1
    }
    // mapping from position to (file, idx)
    posTo := make([]Pair, n+1)
    for i := 1; i <= n; i++ {
        posTo[i] = Pair{-1, -1}
    }
    for i := 0; i < m; i++ {
        for j := 0; j < len(a[i]); j++ {
            posTo[a[i][j]] = Pair{i, j}
        }
    }

    ops := make([][2]int, 0)
    dest := 1
    for i := 0; i < m; i++ {
        for j := 0; j < len(a[i]); j++ {
            curPos := a[i][j]
            if curPos == dest {
                dest++
                continue
            }
            // if dest is occupied by some cluster (not free)
            if dest != freePos {
                p := posTo[dest]
                if p.file != -1 {
                    // move occupant at dest to freePos
                    ops = append(ops, [2]int{dest, freePos})
                    // update mapping for occupant
                    posTo[freePos] = p
                    a[p.file][p.idx] = freePos
                    posTo[dest] = Pair{-1, -1}
                }
            }
            // move our current cluster to dest
            ops = append(ops, [2]int{curPos, dest})
            // update mapping
            posTo[dest] = Pair{i, j}
            a[i][j] = dest
            posTo[curPos] = Pair{-1, -1}
            // update freePos to curPos (old location now free)
            freePos = curPos
            dest++
        }
    }
    // Print operations
    fmt.Fprintln(out, len(ops))
    for _, op := range ops {
        fmt.Fprintf(out, "%d %d\n", op[0], op[1])
    }
}

