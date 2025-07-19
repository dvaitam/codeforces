package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, 2*n)
   for i := 0; i < 2*n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)

   // dp bitsets: chk[i][j] is bitset of achievable sums using j picks from a[2..i]
   chk := make([][]*big.Int, 2*n)
   for i := 0; i < 2*n; i++ {
       chk[i] = make([]*big.Int, n+1)
       for j := 0; j <= n; j++ {
           chk[i][j] = new(big.Int)
       }
   }
   // base: using first element beyond a[0],a[1]
   chk[1][0].SetBit(chk[1][0], 0, 1)

   for i := 2; i < 2*n; i++ {
       for j := 0; j <= i-1 && j <= n-1; j++ {
           // carry over without taking a[i]
           chk[i][j].Set(chk[i-1][j])
           // try taking a[i]
           if j > 0 {
               shifted := new(big.Int).Lsh(chk[i-1][j-1], uint(a[i]))
               chk[i][j].Or(chk[i][j], shifted)
           }
       }
   }

   // total sum of elements from index 2 onward
   tot := 0
   for i := 2; i < 2*n; i++ {
       tot += a[i]
   }
   // find best achievable sum closest to tot/2
   finalDP := chk[2*n-1][n-1]
   best := -1
   for s := 0; s <= tot; s++ {
       if finalDP.Bit(s) == 1 {
           if best < 0 || abs(tot-2*best) > abs(tot-2*s) {
               best = s
           }
       }
   }

   // reconstruct picks
   sum := best
   gae := n - 1
   v := make([]int, 0, n)
   u := make([]int, 0, n)
   for i := 2*n - 1; i >= 2; i-- {
       if sum >= a[i] && gae > 0 && chk[i-1][gae-1].Bit(sum-a[i]) == 1 {
           v = append(v, a[i])
           sum -= a[i]
           gae--
       } else {
           u = append(u, a[i])
       }
   }
   // include the two smallest
   v = append(v, a[0])
   u = append(u, a[1])

   sort.Ints(v)
   sort.Ints(u)
   // reverse u for descending
   for i, j := 0, len(u)-1; i < j; i, j = i+1, j-1 {
       u[i], u[j] = u[j], u[i]
   }

   // output
   for i, x := range v {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, x)
   }
   fmt.Fprintln(writer)
   for i, x := range u {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, x)
   }
   fmt.Fprintln(writer)
}
