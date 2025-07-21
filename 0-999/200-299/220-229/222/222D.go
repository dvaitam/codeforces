package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x int
   if _, err := fmt.Fscan(reader, &n, &x); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   sort.Ints(a)
   sort.Ints(b)
   // Best place is always 1 (Vasya could have the highest possible sum)
   best := 1
   // Compute worst place: maximum number of participants that could have sum >= x
   i, j := 0, n-1
   worst := 0
   for i < n && j >= 0 {
       if a[i]+b[j] >= x {
           worst++
           i++
           j--
       } else {
           i++
       }
   }
   // Vasya could be placed at worst = worst (others above) if ties break against him
   fmt.Printf("%d %d\n", best, worst)
}
