package main

import (
   "bufio"
   "fmt"
   "os"
)

const MX = 1000000

var (
   used [MX]bool
   cur  int64
)

func dfs(start int, limit int, a, b int) {
   // iterative DFS
   stack := []int{start}
   used[start] = true
   cur++
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       v := u + a
       if v <= limit && !used[v] {
           used[v] = true
           cur++
           stack = append(stack, v)
       }
       v = u - b
       if v >= 0 && !used[v] {
           used[v] = true
           cur++
           stack = append(stack, v)
       }
   }
}

func gcd(x, y int) int {
   for y != 0 {
       x, y = y, x%y
   }
   return x
}

// sum computes the total reachable count sum for i from 0 to n
func sum(n int64, d int64) int64 {
   // sum_{i=0..n} (floor(i/d) + 1)
   res := n + 1
   prev := n
   for prev%d != 0 {
       prev--
   }
   prev--
   res += (n - prev) * (n / d)
   prev /= d
   res += d * (prev*(prev+1)/2)
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int64
   var a, b int
   if _, err := fmt.Fscan(in, &m, &a, &b); err != nil {
       return
   }
   // initialize
   used[0] = true
   cur = 1
   var ans int64 = 1
   // process up to MX-1 or m
   lim := MX
   if m+1 < int64(MX) {
       lim = int(m + 1)
   }
   for i := 1; i < lim; i++ {
       if i >= a && used[i-a] {
           dfs(i, i, a, b)
       }
       ans += cur
   }
   // if m exceeds precomputed range, use formula
   if m >= MX {
       d := int64(gcd(a, b))
       ans += sum(m, d) - sum(int64(MX-1), d)
   }
   // output result
   fmt.Println(ans)
}
