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

   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       var s string
       fmt.Fscan(reader, &s)
       // pos: last occurrence, st: first occurrence
       pos := [10]int{}
       st := [10]int{}
       for i := 0; i < 10; i++ {
           pos[i] = -1
           st[i] = n
       }
       vis := make([]bool, n)
       for i := 0; i < n; i++ {
           d := int(s[i] - '0')
           pos[d] = i
           if st[d] == n {
               st[d] = i
           }
           vis[i] = true
       }
       // Greedy assign to color 1 (vis=false)
       i := 0
       for now := 0; now <= 9; now++ {
           var j int
           for j = i; j <= pos[now]; j++ {
               if int(s[j]-'0') == now {
                   vis[j] = false
               }
           }
           if st[now] < i {
               break
           }
           i = j
       }
       // Check that color 2 (vis=true) is non-decreasing
       last := 0
       ok := true
       for idx := 0; idx < n; idx++ {
           if vis[idx] {
               d := int(s[idx] - '0')
               if d >= last {
                   last = d
               } else {
                   ok = false
                   break
               }
           }
       }
       if !ok {
           writer.WriteString("-\n")
       } else {
           // build output
           res := make([]byte, n)
           for idx := 0; idx < n; idx++ {
               if !vis[idx] {
                   res[idx] = '1'
               } else {
                   res[idx] = '2'
               }
           }
           writer.Write(res)
           writer.WriteByte('\n')
       }
   }
}
