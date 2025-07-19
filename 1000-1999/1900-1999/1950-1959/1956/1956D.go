package main

import (
    "bufio"
    "fmt"
    "os"
)

type pair struct { first, second int }

var operations []pair

func permuteRec(n, offset int) {
    if n == 1 {
        return
    }
    solveRec(n-1, offset)
    if n > 2 {
        operations = append(operations, pair{offset + 1, offset + n - 2})
    }
    permuteRec(n-1, offset+1)
}

func solveRec(n, offset int) {
    permuteRec(n, offset)
    operations = append(operations, pair{offset, offset + n - 1})
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    arr := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &arr[i])
    }
    // brute-force masks of length n with values 0,1
    mask := make([]int, n)
    var maxMask []int
    var maxSum int
    for mask[0] < 2 {
        sum := 0
        streak := 0
        for i := 0; i < n; i++ {
            if mask[i] == 0 {
                streak++
            } else {
                sum += streak * streak
                sum += arr[i]
                streak = 0
            }
        }
        sum += streak * streak
        if sum > maxSum {
            maxSum = sum
            maxMask = append([]int(nil), mask...)
        }
        // increment mask
        idx := n - 1
        mask[idx]++
        for idx > 0 && mask[idx] == 2 {
            mask[idx] = 0
            idx--
            mask[idx]++
        }
    }

    // build operations
    operations = operations[:0]
    fmt.Fprint(writer, maxSum, " ")
    prev := -1
    for i := 0; i < n; i++ {
        if maxMask[i] == 0 && arr[i] != 0 {
            operations = append(operations, pair{i, i})
        }
        if maxMask[i] == 0 && prev == -1 {
            prev = i
        } else if maxMask[i] == 1 {
            if prev != -1 {
                solveRec(i-prev, prev)
            }
            prev = -1
        }
    }
    if prev != -1 {
        solveRec(n-prev, prev)
    }

    fmt.Fprintln(writer, len(operations))
    for _, op := range operations {
        // output 1-based indices
        fmt.Fprintln(writer, op.first+1, op.second+1)
    }
}
