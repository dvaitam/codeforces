package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   s := make([]int, n+1)
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   for i := 1; i <= n; i++ {
       p[s[i]] = i
   }

   ans := make([][2]int, 0, n*5)

   doSwap := func(x, y int) {
       ans = append(ans, [2]int{x, y})
       s[x], s[y] = s[y], s[x]
       p[s[x]] = x
       p[s[y]] = y
   }

   var swapPos func(x, y int)
   swapPos = func(x, y int) {
       if x == y {
           return
       }
       if 2*abs(x-y) >= n {
           doSwap(x, y)
           return
       }
       if x < y {
           x, y = y, x
       }
       if 2*abs(x-n) >= n {
           doSwap(x, n)
           doSwap(y, n)
           doSwap(x, n)
           return
       }
       if 2*abs(y-1) >= n {
           doSwap(y, 1)
           doSwap(x, 1)
           doSwap(y, 1)
           return
       }
       doSwap(1, x)
       doSwap(1, n)
       doSwap(y, n)
       doSwap(1, n)
       doSwap(1, x)
   }

   for i := 1; i <= n; i++ {
       if s[i] != i {
           swapPos(i, p[i])
       }
   }

   writer.WriteString(strconv.Itoa(len(ans)))
   writer.WriteByte('\n')
   for _, pr := range ans {
       writer.WriteString(strconv.Itoa(pr[0]))
       writer.WriteByte(' ')
       writer.WriteString(strconv.Itoa(pr[1]))
       writer.WriteByte('\n')
   }
}
