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

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       b[i] = a[i]
   }
   sort.Ints(b)
   perm := make([]int, n)
   for i := 0; i < n; i++ {
       // find rank of a[i] in sorted array b
       perm[i] = sort.SearchInts(b, a[i])
   }
   visited := make([]bool, n)
   cycles := make([][]int, 0, n)
   for i := 0; i < n; i++ {
       if !visited[i] {
           cur := []int{}
           x := i
           for !visited[x] {
               visited[x] = true
               cur = append(cur, x+1)
               x = perm[x]
           }
           cycles = append(cycles, cur)
       }
   }
   fmt.Fprintln(writer, len(cycles))
   for _, cycle := range cycles {
       fmt.Fprint(writer, len(cycle))
       for _, idx := range cycle {
           fmt.Fprint(writer, " ", idx)
       }
       fmt.Fprintln(writer)
   }
}
