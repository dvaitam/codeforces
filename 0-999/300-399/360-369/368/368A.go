package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, d int
   if _, err := fmt.Fscan(reader, &n, &d); err != nil {
       return
   }
   costs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &costs[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   sort.Ints(costs)
   // Guests use hooks with minimal available cost
   used := m
   if used > n {
       used = n
   }
   profit := 0
   for i := 0; i < used; i++ {
       profit += costs[i]
   }
   // remaining guests incur fine
   if m > n {
       profit -= (m - n) * d
   }
   fmt.Println(profit)
}
