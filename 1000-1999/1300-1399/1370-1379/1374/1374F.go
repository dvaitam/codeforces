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
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       b := make([]int, n)
       copy(b, a)
       sort.Ints(b)
       dup := false
       for i := 1; i < n; i++ {
           if b[i] == b[i-1] {
               dup = true
               break
           }
       }
       var ops []int
       // sort with triple rotations
       for i := 0; i <= n-3; i++ {
           mi := i
           for j := i; j < n; j++ {
               if a[j] < a[mi] {
                   mi = j
               }
           }
           for mi != i {
               if mi-2 >= i {
                   // rotate at mi-2
                   x := mi - 2
                   ops = append(ops, x+1)
                   t := a[x+2]
                   a[x+2] = a[x+1]
                   a[x+1] = a[x]
                   a[x] = t
                   mi -= 2
               } else {
                   // rotate at mi-1
                   x := mi - 1
                   ops = append(ops, x+1)
                   t := a[x+2]
                   a[x+2] = a[x+1]
                   a[x+1] = a[x]
                   a[x] = t
                   mi++
               }
           }
       }
       // fix last two if needed
       if n >= 2 && a[n-2] > a[n-1] {
           if !dup {
               fmt.Fprintln(writer, -1)
               continue
           }
           p := n - 3
           for p >= 0 && a[p+1] > a[p+2] {
               x := p
               ops = append(ops, x+1)
               t := a[x+2]
               a[x+2] = a[x+1]
               a[x+1] = a[x]
               a[x] = t
               p--
           }
       }
       // output
       fmt.Fprintln(writer, len(ops))
       for i, v := range ops {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprintln(writer)
   }
}
