package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct { l, r int }

var (
   t, n int
   a    [1005]int
   b    [1005]int
   ans  []pair
   rdr  = bufio.NewReader(os.Stdin)
   wrt  = bufio.NewWriter(os.Stdout)
)

func Print() {
   fmt.Fprintln(wrt, len(ans))
   for _, p := range ans {
       fmt.Fprintln(wrt, p.l, p.r)
   }
}

func Move(l, r, x int) {
   length := r - l + 1
   if x < 0 {
       x += length
   }
   for i := l; i <= r; i++ {
       idx := (i + x - l) % length + l
       b[idx] = a[i]
   }
   for i := l; i <= r; i++ {
       a[i] = b[i]
   }
}

func main() {
   defer wrt.Flush()
   fmt.Fscan(rdr, &t)
   for ; t > 0; t-- {
       ans = ans[:0]
       fmt.Fscan(rdr, &n)
       for i := 1; i <= n; i++ {
           fmt.Fscan(rdr, &a[i])
       }
       flag := true
       for i := 1; i <= n; i++ {
           if a[i] != i {
               flag = false
               break
           }
       }
       if flag {
           fmt.Fprintln(wrt, n)
           fmt.Fprintln(wrt, 0)
           continue
       }
       flag = true
       for i := 1; i < n; i++ {
           diff := a[i+1] - a[i]
           if diff != 1 && diff != 1-n {
               flag = false
               break
           }
       }
       if flag {
           fmt.Fprintln(wrt, n-1)
           p := 0
           for i := 1; i <= n; i++ {
               if a[i] == 1 {
                   p = i
                   break
               }
           }
           fmt.Fprintln(wrt, p-1)
           for i := 1; i < p; i++ {
               fmt.Fprintln(wrt, 2, 1)
           }
           continue
       }
       inv := 0
       for i := 1; i < n; i++ {
           for j := i + 1; j <= n; j++ {
               if a[i] > a[j] {
                   inv++
               }
           }
       }
       if inv&1 == 1 && (n-1)&1 == 1 {
           fmt.Fprintln(wrt, n-3)
           posn := 0
           for i := 1; i <= n; i++ {
               if a[i] == n {
                   posn = i
                   break
               }
           }
           if posn == 1 {
               ans = append(ans, pair{2, 1})
               Move(1, n-2, -1)
               posn = n - 2
           } else if posn == 2 {
               ans = append(ans, pair{1, 2})
               Move(1, n-2, 1)
               posn = 3
           }
           for i := posn; i < n; i++ {
               ans = append(ans, pair{3, 4})
           }
           Move(3, n, n-posn)
           n--
           if a[1] == 1 {
               ans = append(ans, pair{1, 2})
               tmp := a[n-1]
               for i := n - 1; i >= 2; i-- {
                   a[i] = a[i-1]
               }
               a[1] = tmp
           } else if a[n] == 1 {
               ans = append(ans, pair{2, 3})
               for i := n; i >= 3; i-- {
                   a[i] = a[i-1]
               }
               a[2] = 1
           }
           for i := 2; i <= n-2; i++ {
               pos1, posi := 0, 0
               for j := 1; j <= n; j++ {
                   if a[j] == 1 {
                       pos1 = j
                   }
                   if a[j] == i {
                       posi = j
                   }
               }
               if posi == 1 {
                   if pos1 == 2 {
                       ans = append(ans, pair{1, 2})
                       tmp := a[n-1]
                       for j := n - 1; j >= 2; j-- {
                           a[j] = a[j-1]
                       }
                       a[1] = tmp
                       pos1++; posi++
                   } else {
                       ans = append(ans, pair{2, 1})
                       for j := 1; j <= n-2; j++ {
                           a[j] = a[j+1]
                       }
                       a[n-1] = i
                       pos1--; posi = n - 1
                   }
               }
               if posi != n {
                   if posi < pos1 {
                       for j := 2; j <= posi; j++ {
                           ans = append(ans, pair{3, 2})
                       }
                       Move(2, n, -(posi-1))
                       pos1 -= posi - 1
                       posi = n
                   } else {
                       for j := posi; j < n; j++ {
                           ans = append(ans, pair{2, 3})
                       }
                       Move(2, n, n-posi)
                       pos1 += n - posi
                       posi = n
                   }
               }
               for j := pos1 + i - 2; j < n-1; j++ {
                   ans = append(ans, pair{1, 2})
               }
               Move(1, n-1, n-1-(pos1+i-2))
               pos1 += n-1 - (pos1 + i - 2)
               ans = append(ans, pair{3, 2})
               Move(2, n, -1)
           }
           ans = append(ans, pair{2, 1})
           Move(1, n-1, -1)
           if a[n-1] == n-1 {
               Print()
               continue
           } else {
               for i := 1; i <= (n-3)/2; i++ {
                   ans = append(ans, pair{1, 2})
                   ans = append(ans, pair{2, 3})
                   ans = append(ans, pair{2, 3})
                   ans = append(ans, pair{2, 1})
                   ans = append(ans, pair{3, 2})
                   ans = append(ans, pair{2, 1})
                   ans = append(ans, pair{2, 3})
                   ans = append(ans, pair{2, 3})
               }
               ans = append(ans, pair{2, 3})
               Print()
           }
       } else {
           fmt.Fprintln(wrt, n-2)
           if a[1] == 1 {
               ans = append(ans, pair{1, 2})
               tmp := a[n-1]
               for i := n - 1; i >= 2; i-- {
                   a[i] = a[i-1]
               }
               a[1] = tmp
           } else if a[n] == 1 {
               ans = append(ans, pair{2, 3})
               for i := n; i >= 3; i-- {
                   a[i] = a[i-1]
               }
               a[2] = 1
           }
           for i := 2; i <= n-2; i++ {
               pos1, posi := 0, 0
               for j := 1; j <= n; j++ {
                   if a[j] == 1 {
                       pos1 = j
                   }
                   if a[j] == i {
                       posi = j
                   }
               }
               if posi == 1 {
                   if pos1 == 2 {
                       ans = append(ans, pair{1, 2})
                       tmp := a[n-1]
                       for j := n - 1; j >= 2; j-- {
                           a[j] = a[j-1]
                       }
                       a[1] = tmp
                       pos1++; posi++
                   } else {
                       ans = append(ans, pair{2, 1})
                       for j := 1; j <= n-2; j++ {
                           a[j] = a[j+1]
                       }
                       a[n-1] = i
                       pos1--; posi = n - 1
                   }
               }
               if posi != n {
                   if posi < pos1 {
                       for j := 2; j <= posi; j++ {
                           ans = append(ans, pair{3, 2})
                       }
                       Move(2, n, -(posi-1))
                       pos1 -= posi - 1
                       posi = n
                   } else {
                       for j := posi; j < n; j++ {
                           ans = append(ans, pair{2, 3})
                       }
                       Move(2, n, n-posi)
                       pos1 += n - posi
                       posi = n
                   }
               }
               for j := pos1 + i - 2; j < n-1; j++ {
                   ans = append(ans, pair{1, 2})
               }
               Move(1, n-1, n-1-(pos1+i-2))
               pos1 += n-1 - (pos1 + i - 2)
               ans = append(ans, pair{3, 2})
               Move(2, n, -1)
           }
           ans = append(ans, pair{2, 1})
           Move(1, n-1, -1)
           if a[n-1] == n-1 {
               Print()
               continue
           } else {
               for i := 1; i <= (n-3)/2; i++ {
                   ans = append(ans, pair{1, 2})
                   ans = append(ans, pair{2, 3})
                   ans = append(ans, pair{2, 3})
                   ans = append(ans, pair{2, 1})
                   ans = append(ans, pair{3, 2})
                   ans = append(ans, pair{2, 1})
                   ans = append(ans, pair{2, 3})
                   ans = append(ans, pair{2, 3})
               }
               ans = append(ans, pair{2, 3})
               Print()
           }
       }
   }
}
