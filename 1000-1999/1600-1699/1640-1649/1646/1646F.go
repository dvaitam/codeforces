package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(in, &a[i][j])
           a[i][j]--
       }
   }
   c := make([][]int, n)
   for i := 0; i < n; i++ {
       c[i] = make([]int, n)
       for j := 0; j < n; j++ {
           c[i][a[i][j]]++
       }
   }
   var s [][]int
   // build initial permutations resolving duplicates
   for {
       found := false
       var i0, j0 int
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               if c[i][j] > 1 {
                   found = true
                   i0, j0 = i, j
                   break
               }
           }
           if found {
               break
           }
       }
       if !found {
           break
       }
       cur := make([]int, n)
       cur[i0] = j0
       c[i0][j0]--
       j := j0
       p := (i0 + 1) % n
       for p != i0 {
           c[p][j]++
           var q int
           for q = 0; q < n; q++ {
               if c[p][q] > 1 {
                   break
               }
           }
           cur[p] = q
           c[p][q]--
           j = q
           p = (p + 1) % n
       }
       c[i0][j]++
       s = append(s, cur)
   }
   // append cyclic shift permutations
   for i := 1; i < n; i++ {
       for j := i; j >= 1; j-- {
           cur := make([]int, n)
           for p := 0; p < n; p++ {
               cur[p] = (p + j) % n
           }
           s = append(s, cur)
       }
   }
   // output
   fmt.Fprintln(out, len(s))
   var sb strings.Builder
   for _, cur := range s {
       sb.Reset()
       for idx, v := range cur {
           sb.WriteString(strconv.Itoa(v + 1))
           if idx+1 < n {
               sb.WriteByte(' ')
           }
       }
       fmt.Fprintln(out, sb.String())
   }
}
