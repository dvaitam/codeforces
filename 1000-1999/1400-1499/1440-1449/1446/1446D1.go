package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   maxVal := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxVal {
           maxVal = a[i]
       }
   }
   // sliding window helper
   counts := make([]int, maxVal+1)
   freqCount := make([]int, n+2)
   var check = func(L int) bool {
       if L < 2 {
           return false
       }
       // reset counts and freqCount
       for i := 1; i <= maxVal; i++ {
           counts[i] = 0
       }
       for i := 0; i <= L; i++ {
           freqCount[i] = 0
       }
       fmax := 0
       // init first window
       for i := 0; i < L; i++ {
           v := a[i]
           old := counts[v]
           if old > 0 {
               freqCount[old]--
           }
           counts[v] = old + 1
           freqCount[old+1]++
           if old+1 > fmax {
               fmax = old + 1
           }
       }
       if freqCount[fmax] >= 2 {
           return true
       }
       // slide
       for i := L; i < n; i++ {
           // add a[i]
           v := a[i]
           old := counts[v]
           if old > 0 {
               freqCount[old]--
           }
           counts[v] = old + 1
           freqCount[old+1]++
           if old+1 > fmax {
               fmax = old + 1
           }
           // remove a[i-L]
           u := a[i-L]
           old2 := counts[u]
           freqCount[old2]--
           counts[u] = old2 - 1
           if old2-1 > 0 {
               freqCount[old2-1]++
           }
           if old2 == fmax && freqCount[old2] == 0 {
               // decrease fmax
               for fmax > 0 && freqCount[fmax] == 0 {
                   fmax--
               }
           }
           if freqCount[fmax] >= 2 {
               return true
           }
       }
       return false
   }
   // binary search
   lo, hi := 0, n+1
   for lo+1 < hi {
       mid := (lo + hi) >> 1
       if check(mid) {
           lo = mid
       } else {
           hi = mid
       }
   }
   fmt.Println(lo)
}
