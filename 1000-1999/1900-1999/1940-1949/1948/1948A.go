package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t, n int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        if _, err := fmt.Fscan(reader, &n); err != nil {
            return
        }
        if n%2 == 1 {
            writer.WriteString("NO\n")
        } else {
            writer.WriteString("YES\n")
            halves := n / 2
            for i := 1; i <= halves; i++ {
                if i%2 == 1 {
                    writer.WriteString("AA")
                } else {
                    writer.WriteString("BB")
                }
            }
            writer.WriteByte('\n')
        }
    }
}
