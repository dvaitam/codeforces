package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

// readInt reads the next integer from standard input
func readInt() int {
   var c byte
   var x int
   // skip non-digits
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c >= '0' && c <= '9' {
           x = int(c - '0')
           break
       }
   }
   // read remaining digits
   for {
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       c = b
       if c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
   }
   return x
}

func main() {
   defer writer.Flush()
   t := readInt()
   for tc := 0; tc < t; tc++ {
       n := readInt()
       used := make([]bool, n+1)
       unmatched := 0
       for i := 1; i <= n; i++ {
           m := readInt()
           matched := false
           for j := 0; j < m; j++ {
               p := readInt()
               if !matched && p >= 1 && p <= n && !used[p] {
                   used[p] = true
                   matched = true
               }
           }
           if !matched {
               unmatched = i
           }
       }
       if unmatched == 0 {
           writer.WriteString("OPTIMAL\n")
       } else {
           writer.WriteString("IMPROVE\n")
           var prince int
           for j := 1; j <= n; j++ {
               if !used[j] {
                   prince = j
                   break
               }
           }
           writer.WriteString(fmt.Sprintf("%d %d\n", unmatched, prince))
       }
   }
}
