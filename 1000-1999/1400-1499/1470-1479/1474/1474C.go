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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       N := n * 2
       a := make([]int, N)
       for i := 0; i < N; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       // Original counts
       orig := make(map[int]int, N)
       for _, v := range a {
           orig[v]++
       }
       found := false
       var result [][2]int
       var initialSum int
       // Try pairing largest with each
       for i := 0; i < N-1; i++ {
           x := a[N-1] + a[i]
           b := make(map[int]int, len(orig))
           for k, v := range orig {
               b[k] = v
           }
           tf := x
           idx := N - 1
           tmp := make([][2]int, 0, n)
           ok := true
           for k := 0; k < n; k++ {
               // find current largest available
               for idx >= 0 && b[a[idx]] == 0 {
                   idx--
               }
               if idx < 0 {
                   ok = false
                   break
               }
               v := a[idx]
               b[v]--
               need := tf - v
               if b[need] == 0 {
                   ok = false
                   break
               }
               b[need]--
               tmp = append(tmp, [2]int{need, v})
               tf = v
           }
           if ok {
               found = true
               result = tmp
               initialSum = x
               break
           }
       }
       if !found {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       fmt.Fprintln(writer, initialSum)
       for _, p := range result {
           fmt.Fprintln(writer, p[0], p[1])
       }
   }
}
