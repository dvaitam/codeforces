package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a int
   if _, err := fmt.Fscan(reader, &n, &a); err != nil {
       return
   }
   xs := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i])
   }
   if n <= 1 {
       fmt.Println(0)
       return
   }
   sort.Ints(xs)
   mn := xs[0]
   mx := xs[n-1]
   secondMn := xs[1]
   secondMx := xs[n-2]
   // option 1: skip the smallest
   skipLeft := (mx - secondMn) + min(abs(a-secondMn), abs(a-mx))
   // option 2: skip the largest
   skipRight := (secondMx - mn) + min(abs(a-mn), abs(a-secondMx))
   ans := skipLeft
   if skipRight < ans {
       ans = skipRight
   }
   fmt.Println(ans)
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
