package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+2)
   pred := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > 0 {
           pred[a[i]] = i
       }
   }

   taken := make([]bool, n+2)
   avail := make(map[int]struct{})
   freeUnfixed := make(map[int]struct{})
   // initialize availability: all seats initially available
   for i := 1; i <= n; i++ {
       avail[i] = struct{}{}
       if a[i] == 0 {
           freeUnfixed[i] = struct{}{}
       }
   }

   ans := 1
   for t := 1; t <= n; t++ {
       if pos := pred[t]; pos > 0 {
           // predetermined seat
           if _, ok := avail[pos]; !ok {
               fmt.Fprintln(writer, 0)
               return
           }
           remove(pos, n, taken, avail, freeUnfixed)
       } else {
           k := len(freeUnfixed)
           if k == 0 {
               fmt.Fprintln(writer, 0)
               return
           }
           ans = int(int64(ans) * int64(k) % MOD)
           // pick arbitrary seat
           var pos int
           for idx := range freeUnfixed {
               pos = idx
               break
           }
           remove(pos, n, taken, avail, freeUnfixed)
       }
   }
   fmt.Fprintln(writer, ans)
}

// remove marks seat i taken, updates availability and freeUnfixed
func remove(i, n int, taken []bool, avail, freeUnfixed map[int]struct{}) {
   taken[i] = true
   delete(avail, i)
   if _, ok := freeUnfixed[i]; ok {
       delete(freeUnfixed, i)
   }
   // update neighbors
   for _, j := range []int{i - 1, i + 1} {
       if j < 1 || j > n || taken[j] {
           continue
       }
       // check if j now has both neighbors taken
       leftTaken := j-1 >= 1 && taken[j-1]
       rightTaken := j+1 <= n && taken[j+1]
       if leftTaken && rightTaken {
           // becomes unavailable
           if _, ok := avail[j]; ok {
               delete(avail, j)
           }
           if _, ok := freeUnfixed[j]; ok {
               delete(freeUnfixed, j)
           }
       }
   }
}
