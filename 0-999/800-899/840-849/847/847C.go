package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n64, k int64
    if _, err := fmt.Fscan(in, &n64, &k); err != nil {
        return
    }
    sum := n64 * (n64 - 1) / 2
    if sum < k {
        fmt.Println("Impossible")
        return
    }
    var builder strings.Builder
    // k == 0: print n pairs of ()
    if k == 0 {
        for i := int64(0); i < n64; i++ {
            builder.WriteString("()")
        }
        fmt.Println(builder.String())
        return
    }
    origN := int(n64)
    used := 0
    remSum := sum
    i := n64 - 1
    // build prefix of simple pairs until adding next exceeds k
    for remSum >= k {
        remSum -= i
        if remSum < k {
            // need to insert a block
            d := k - remSum
            rem := origN - used
            // add opening parentheses, insert () at position d
            for j := 0; int64(j) < int64(rem); j++ {
                builder.WriteByte('(')
                if int64(j) == d-1 {
                    builder.WriteString("()")
                    rem--
                }
            }
            // close remaining opens
            for j := 0; j < rem; j++ {
                builder.WriteByte(')')
            }
            break
        } else {
            builder.WriteString("()")
            used++
        }
        i--
    }
    fmt.Println(builder.String())
}
