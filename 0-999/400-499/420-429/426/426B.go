package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // compute maximum power of two dividing n
   maxK := 0
   for t := n; t%2 == 0; t /= 2 {
       maxK++
   }
   // try smallest x = n/(2^k)
   for k := maxK; k >= 0; k-- {
       div := 1 << k
       if n%div != 0 {
           continue
       }
       x := n / div
       if check(a, n, x) {
           fmt.Println(x)
           return
       }
   }
}

// check if initial rows x can generate a by mirrorings
func check(a [][]int, total, x int) bool {
   n := total
   for i := x; i < n; i++ {
       idx := mapIndex(i, n, x)
       // compare row i with row idx
       for j := range a[i] {
           if a[i][j] != a[idx][j] {
               return false
           }
       }
   }
   return true
}

// mapIndex maps row index i in final matrix of length total to initial index < x
func mapIndex(i, total, x int) int {
   cur := i
   length := total
   for length > x {
       half := length / 2
       if cur >= half {
           cur = length - cur - 1
       }
       length = half
   }
   return cur
}
