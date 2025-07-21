package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // build log table
   logs := make([]int, n+2)
   logs[1] = 0
   for i := 2; i <= n; i++ {
       logs[i] = logs[i/2] + 1
   }
   kmax := logs[n] + 1
   // build sparse table for gcd
   st := make([][]int, kmax)
   st[0] = make([]int, n+1)
   for i := 1; i <= n; i++ {
       st[0][i] = a[i]
   }
   for k := 1; k < kmax; k++ {
       step := 1 << (k - 1)
       st[k] = make([]int, n+1)
       for i := 1; i+ (1<<k) -1 <= n; i++ {
           st[k][i] = gcd(st[k-1][i], st[k-1][i+step])
       }
   }
   // query function
   query := func(l, r int) int {
       length := r - l + 1
       k := logs[length]
       j := r - (1 << k) + 1
       return gcd(st[k][l], st[k][j])
   }
   var maxd int
   lvals := make([]int, 0, 16)
   // for each position, find largest segment where gcd equals a[j]
   for j := 1; j <= n; j++ {
       aj := a[j]
       // find left boundary
       low, high := 1, j
       for low < high {
           mid := (low + high) >> 1
           if query(mid, j) == aj {
               high = mid
           } else {
               low = mid + 1
           }
       }
       lj := low
       // find right boundary (first false then -1)
       low, high = j, n+1
       for low < high {
           mid := (low + high) >> 1
           if mid <= n && query(j, mid) == aj {
               low = mid + 1
           } else {
               high = mid
           }
       }
       rj := low - 1
       d := rj - lj
       if d > maxd {
           maxd = d
           lvals = lvals[:0]
           lvals = append(lvals, lj)
       } else if d == maxd {
           lvals = append(lvals, lj)
       }
   }
   // sort and dedup lvals
   sort.Ints(lvals)
   uniq := make([]int, 0, len(lvals))
   for i, v := range lvals {
       if i == 0 || v != lvals[i-1] {
           uniq = append(uniq, v)
       }
   }
   // output
   fmt.Fprintln(out, len(uniq), maxd)
   for i, v := range uniq {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
