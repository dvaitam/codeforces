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
   var tt int
   fmt.Fscan(reader, &tt)
   for tt > 0 {
       tt--
       solve(reader, writer)
   }
}

// solve processes one test case
func solve(r *bufio.Reader, w *bufio.Writer) {
   var n, m int
   var x int64
   fmt.Fscan(r, &n, &m, &x)
   a := make([]pair, n)
   for i := 0; i < n; i++ {
       var v int64
       fmt.Fscan(r, &v)
       a[i] = pair{v, i}
   }
   sort.Slice(a, func(i, j int) bool {
       return a[i].value < a[j].value
   })
   b := make([]int, n)
   fmt.Fprintln(w, "YES")
   for i := 0; i < n; i++ {
       idx := a[i].index
       b[idx] = i%m + 1
   }
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Fprint(w, " ")
       }
       fmt.Fprint(w, b[i])
   }
   fmt.Fprintln(w)
}

// pair holds value and original index
type pair struct {
   value int64
   index int
}
