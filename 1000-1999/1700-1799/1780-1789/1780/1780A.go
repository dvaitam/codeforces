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

   var tc int
   if _, err := fmt.Fscan(reader, &tc); err != nil {
       return
   }
   for i := 0; i < tc; i++ {
       solve(reader, writer)
   }
}

func solve(r *bufio.Reader, w *bufio.Writer) {
   var n int
   fmt.Fscan(r, &n)
   odd := make([]int, 0, n)
   even := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       var a int
       fmt.Fscan(r, &a)
       if a&1 == 1 {
           odd = append(odd, i)
       } else {
           even = append(even, i)
       }
   }
   ok := false
   ans := make([]int, 3)
   if len(odd) >= 3 {
       ok = true
       ans[0], ans[1], ans[2] = odd[0], odd[1], odd[2]
   } else if len(odd) >= 1 && len(even) >= 2 {
       ok = true
       ans[0] = odd[0]
       ans[1], ans[2] = even[0], even[1]
   }
   if ok {
       fmt.Fprintln(w, "YES")
       fmt.Fprintf(w, "%d %d %d\n", ans[0], ans[1], ans[2])
   } else {
       fmt.Fprintln(w, "NO")
   }
}
