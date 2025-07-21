package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func manacher(s []byte) []int {
   n := len(s)
   d := make([]int, n)
   l, r := 0, -1
   for i := 0; i < n; i++ {
       var k int
       if i > r {
           k = 1
       } else {
           k = min(d[l+r-i], r-i) + 1
       }
       for i-k >= 0 && i+k < n && s[i-k] == s[i+k] {
           k++
       }
       d[i] = k - 1
       if i+k-1 > r {
           l = i - k + 1
           r = i + k - 1
       }
   }
   return d
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t string
   fmt.Fscan(reader, &t)
   s := []byte(t)
   n := len(s)

   // palindromic radius
   pal := manacher(s)

   // prev and next occurrence
   const K = 26
   prev := make([]int, n)
   last := make([]int, K)
   for i := range last {
       last[i] = -1
   }
   for i := 0; i < n; i++ {
       c := int(s[i] - 'a')
       prev[i] = last[c]
       last[c] = i
   }
   next := make([]int, n)
   for i := range last {
       last[i] = n
   }
   for i := n - 1; i >= 0; i-- {
       c := int(s[i] - 'a')
       next[i] = last[c]
       last[c] = i
   }
   uniq := make([]int, n)
   for i := 0; i < n; i++ {
       ldist := n
       if prev[i] >= 0 {
           ldist = i - prev[i]
       }
       rdist := n
       if next[i] < n {
           rdist = next[i] - i
       }
       tmp := min(ldist, rdist) - 1
       if tmp < 0 {
           tmp = 0
       }
       uniq[i] = tmp
   }

   // levels
   maxH := 0
   for h := 1; ; h++ {
       if (1<<h)-1 > n {
           break
       }
       maxH = h
   }
   d := make([]int, maxH+1)
   off := make([]int, maxH+1)
   for h := 1; h <= maxH; h++ {
       if h == 1 {
           d[h] = 0
           off[h] = 0
       } else {
           d[h] = (1 << (h - 1)) - 1
           off[h] = 1 << (h - 2)
       }
   }

   // dp arrays
   good := make([][]bool, n)
   dp := make([][]bool, n)
   for i := 0; i < n; i++ {
       good[i] = make([]bool, maxH+1)
       dp[i] = make([]bool, maxH+1)
       good[i][1] = true
       dp[i][1] = true
   }
   for h := 2; h <= maxH; h++ {
       dh := d[h]
       oh := off[h]
       for k := 0; k < n; k++ {
           if k-oh >= 0 && k+oh < n && pal[k] >= dh && good[k][h-1] && good[k-oh][h-1] {
               good[k][h] = true
               if uniq[k] >= dh {
                   dp[k][h] = true
               }
           }
       }
   }

   // precompute sq lens
   lenh := make([]int64, maxH+1)
   sqh := make([]int64, maxH+1)
   for h := 1; h <= maxH; h++ {
       l := int64((1<<h) - 1)
       lenh[h] = l
       sqh[h] = l * l
   }

   var origSum int64
   bestDelta := int64(0)
   for k := 0; k < n; k++ {
       Hk := 0
       Gk := 0
       for h := 1; h <= maxH; h++ {
           if dp[k][h] {
               origSum += sqh[h]
               Hk = h
           }
           if good[k][h] {
               Gk = h
           }
       }
       if Gk > Hk {
           var delta int64
           for h := Hk + 1; h <= Gk; h++ {
               delta += sqh[h]
           }
           if delta > bestDelta {
               bestDelta = delta
           }
       }
   }
   fmt.Println(origSum + bestDelta)
}
