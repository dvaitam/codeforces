package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // k-th largest is at index n-k
   fmt.Println(a[n-k])
}
