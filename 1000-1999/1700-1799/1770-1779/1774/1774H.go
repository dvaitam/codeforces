package main

import (
   "bufio"
   "fmt"
   "os"
)

var seq = make([]int, 200005)
var ans = make([]int, 200005)

func f(x, y, k int) int {
   return (x-1)*k + y
}

// build sequence in seq[1..nLocal] based on kLocal
func buildSeq(nLocal, kLocal int) {
   if kLocal%2 == 0 {
       // even kLocal: simple interleave
       for i := 1; i <= nLocal/2; i++ {
           seq[2*i-1] = i
           seq[2*i] = nLocal + 1 - i
       }
       return
   }
   m := nLocal / kLocal
   cur := 3 * m
   // for blocks 4..kLocal
   for i := 4; i <= kLocal; i++ {
       if i%2 == 1 {
           for j := m; j >= 1; j-- {
               cur++
               seq[f(j, i, kLocal)] = cur
           }
       } else {
           for j := 1; j <= m; j++ {
               cur++
               seq[f(j, i, kLocal)] = cur
           }
       }
   }
   half := (m + 1) >> 1
   mid := (3*m + 3) >> 1
   for i := 1; i <= half; i++ {
       seq[f(i, 1, kLocal)] = 2*i - 1
       seq[f(i, 2, kLocal)] = mid - i
       seq[f(i, 3, kLocal)] = 3*m - i + 1
   }
   start := (m + 3) >> 1
   for i := start; i <= m; i++ {
       delta := i - start
       seq[f(i, 1, kLocal)] = mid + delta
       seq[f(i, 2, kLocal)] = (m<<1) + 1 + delta
       seq[f(i, 3, kLocal)] = m - (m & 1) - (delta << 1)
   }
}

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       // ensure ans indexed up to n
       if n%k == 0 {
           buildSeq(n, k)
           for i := 1; i <= n; i++ {
               ans[i] = seq[i]
           }
       } else {
           q := n / k
           r := n % k
           if r == 1 {
               cur := 0
               delta := (q<<1) + 1
               // positions mod k == 1
               for i := 1; i <= n; i += k {
                   cur++
                   ans[i] = cur
               }
               // positions mod k == 0 (1-indexed)
               for i := n - k + 1; i >= 2; i -= k {
                   cur++
                   ans[i] = cur
               }
               // middle blocks
               buildSeq(q*(k-2), k-2)
               cur = 0
               for i := 3; i <= n; i += k {
                   for j := i; j <= i+k-3; j++ {
                       cur++
                       ans[j] = seq[cur] + delta
                   }
               }
           } else if k-r == 1 {
               if q == 1 {
                   cur := 0
                   // place one
                   ans[k] = n
                   buildSeq(n-1, k-1)
                   for i := 1; i < k; i++ {
                       cur++
                       ans[i] = seq[cur]
                   }
                   for i := k + 1; i <= n; i++ {
                       cur++
                       ans[i] = seq[cur]
                   }
               } else {
                   cur := n + 1
                   delta := q + 1
                   for i := k; i <= n; i += k {
                       cur--
                       ans[i] = cur
                   }
                   cur = 0
                   for i := 1; i <= n; i += k {
                       cur++
                       ans[i] = cur
                   }
                   buildSeq((q+1)*(r-1), r-1)
                   cur = 0
                   for i := 2; i <= n; i += k {
                       for j := i; j <= i+r-2; j++ {
                           cur++
                           ans[j] = seq[cur] + delta
                       }
                   }
               }
           } else {
               cur := 0
               delta := (q+1)*r
               buildSeq((q+1)*r, r)
               for i := 1; i <= n; i += k {
                   for j := i; j <= i+r-1; j++ {
                       cur++
                       ans[j] = seq[cur]
                   }
               }
               buildSeq(q*(k-r), k-r)
               cur = 0
               for i := r + 1; i <= n; i += k {
                   for j := i; j <= i+(k-r)-1; j++ {
                       cur++
                       ans[j] = seq[cur] + delta
                   }
               }
           }
       }
       // compute minimum sum of window size k and output
       var sum, best int64
       for i := 1; i <= k; i++ {
           sum += int64(ans[i])
       }
       best = sum
       for i := k + 1; i <= n; i++ {
           sum += int64(ans[i]) - int64(ans[i-k])
           if tmp := sum; tmp < best {
               best = tmp
           }
       }
       fmt.Fprintln(writer, best)
       // print permutation
       for i := 1; i <= n; i++ {
           writer.WriteString(fmt.Sprintf("%d", ans[i]))
           if i < n {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}
