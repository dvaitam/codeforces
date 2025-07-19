package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   readInt := func() int {
       var x int
       var sign = 1
       b, err := reader.ReadByte()
       for err == nil && (b < '0' || b > '9') && b != '-' {
           b, err = reader.ReadByte()
       }
       if err != nil {
           return 0
       }
       if b == '-' {
           sign = -1
           b, _ = reader.ReadByte()
       }
       for err == nil && b >= '0' && b <= '9' {
           x = x*10 + int(b - '0')
           b, err = reader.ReadByte()
       }
       return x * sign
   }

   t := readInt()
   for t > 0 {
       t--
       n := readInt()
       a := make([]int, n)
       for i := 0; i < n; i++ {
           a[i] = readInt()
       }
       // reverse a
       for i := 0; i < n/2; i++ {
           a[i], a[n-1-i] = a[n-1-i], a[i]
       }
       // linked list via nxt
       nxt := make([]int, n+2)
       prev := n + 1
       for _, v := range a {
           nxt[prev] = v
           prev = v
       }
       nxt[prev] = 0
       ans := make([]int, n+1)
       vs := make([]bool, n+2)
       col := 0
       // process rounds
       for nxt[n+1] != 0 {
           col++
           // collect current nodes
           s := make([]int, 0, n)
           for i := n + 1; nxt[i] != 0; i = nxt[i] {
               x := nxt[i]
               vs[x] = false
               s = append(s, x)
           }
           sort.Ints(s)
           // process
           i := n + 1
           for nxt[i] != 0 {
               x := nxt[i]
               if !vs[x] {
                   // find index in s
                   idx := sort.SearchInts(s, x)
                   // mark neighbors
                   if idx > 0 {
                       vs[s[idx-1]] = true
                   }
                   if idx+1 < len(s) {
                       vs[s[idx+1]] = true
                   }
                   ans[x] = col
                   // remove x from linked list
                   nxt[i] = nxt[x]
                   // remove x from s
                   s = append(s[:idx], s[idx+1:]...)
               } else {
                   // skip, remove from s
                   i = x
                   // remove x from s
                   idx := sort.SearchInts(s, x)
                   if idx < len(s) && s[idx] == x {
                       s = append(s[:idx], s[idx+1:]...)
                   }
               }
           }
       }
       // output
       for i := 1; i <= n; i++ {
           fmt.Fprintf(writer, "%d ", ans[i])
       }
       fmt.Fprint(writer, "\n")
   }
}
