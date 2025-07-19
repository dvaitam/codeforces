package main

import (
   "bufio"
   "fmt"
   "os"
)

const md = 998244353

var (
   k    int
   A, h int64
   mid  int
   w    []int64
   a, b []int64
   mp1, mp2 map[int64]int
   flagFound bool
)

func modPow(a, b int64) int64 {
   res := int64(1)
   for b > 0 {
       if b&1 == 1 {
           res = res * a % md
       }
       a = a * a % md
       b >>= 1
   }
   return res
}

func dfs1(t, lc int, s, num int64, nn int) {
   // find next empty position
   for a[lc] != 0 {
       lc++
       if lc > mid {
           lc = 1
           s = s/2 + 1
       }
   }
   if s == 2 {
       // store two scenarios
       r1 := (num + int64(lc)*w[1]) % md
       mp1[r1] = nn + (1 << mid)
       r2 := (num + int64(lc)*w[2]) % md
       mp2[r2] = nn + (1 << mid)
       return
   }
   cnt := lc
   a[cnt] = s
   pp := s
   kk := (num + int64(cnt)*w[int(s)]) % md
   // find next
   for a[cnt] != 0 {
       cnt++
   }
   cnt++
   if cnt > mid {
       cnt = 1
       pp = pp/2 + 1
   }
   dfs1(t+1, cnt, pp, kk, nn | (1 << (t-1)))
   // backtrack
   cnt = lc
   a[cnt] = 0
   cnt++
   for a[cnt] != 0 {
       cnt++
   }
   q := cnt
   a[cnt] = s
   num = (num + int64(cnt)*w[int(s)]) % md
   cnt++
   if cnt > mid {
       cnt = 1
       s = s/2 + 1
   }
   dfs1(t+1, cnt, s, num, nn)
   a[q] = 0
}

func work(x int) {
   kk := x
   cnt := 1
   cc := mid + 1
   for i := 0; i < mid-1; i++ {
       // find next empty in first half (with wrap)
       for b[cnt] != 0 {
           cnt++
           if cnt > mid {
               cnt = 1
               cc = cc/2 + 1
           }
       }
       if kk&(1<<i) != 0 {
           // assign and move past the assigned slot
           b[cnt] = int64(cc)
           // find next empty slot (no wrap)
           for b[cnt] != 0 {
               cnt++
           }
           cnt++
           if cnt > mid {
               cnt = 1
               cc = cc/2 + 1
           }
       } else {
           // skip current and assign to next empty (no wrap)
           cnt++
           for b[cnt] != 0 {
               cnt++
           }
           b[cnt] = int64(cc)
       }
   }
   // copy second half
   for i := mid + 1; i <= mid*2; i++ {
       b[i] = a[i]
   }
}

func dfs2(t, lc int, s, num int64) {
   if flagFound {
       return
   }
   for a[lc] != 0 {
       lc++
       if lc > mid*2 {
           lc = mid + 1
           s = s/2 + 1
       }
   }
   if s == 2 {
       // try matchings
       r1 := (h - (num + int64(lc)*w[1])%md + md) % md
       if v, ok := mp2[r1]; ok && v != 0 {
           flagFound = true
           work(v)
           // fill remaining b slots
           bo := false
           for i := 1; i <= mid*2; i++ {
               if b[i] == 0 {
                   if !bo {
                       bo = true
                       b[i] = 2
                   } else {
                       b[i] = 1
                   }
               }
           }
           return
       }
       r2 := (h - (num + int64(lc)*w[2])%md + md) % md
       if v, ok := mp1[r2]; ok && v != 0 {
           flagFound = true
           work(v)
           bo := false
           for i := 1; i <= mid*2; i++ {
               if b[i] == 0 {
                   if !bo {
                       bo = true
                       b[i] = 1
                   } else {
                       b[i] = 2
                   }
               }
           }
           return
       }
       return
   }
   cnt := lc
   a[cnt] = s
   pp := s
   kk := (num + int64(cnt)*w[int(s)]) % md
   for a[cnt] != 0 {
       cnt++
   }
   cnt++
   if cnt > mid*2 {
       cnt = mid + 1
       pp = pp/2 + 1
   }
   dfs2(t+1, cnt, pp, kk)
   if flagFound {
       return
   }
   cnt = lc
   a[cnt] = 0
   cnt++
   for a[cnt] != 0 {
       cnt++
   }
   q := cnt
   a[cnt] = s
   num = (num + int64(cnt)*w[int(s)]) % md
   cnt++
   if cnt > mid*2 {
       cnt = mid + 1
       s = s/2 + 1
   }
   dfs2(t+1, cnt, s, num)
   a[q] = 0
}

func main() {
   in := bufio.NewReader(os.Stdin)
   _, _ = fmt.Fscan(in, &k, &A, &h)
   // number of leaves per half
   if k < 1 {
       fmt.Println(-1)
       return
   }
   mid = 1 << (k - 1)
   // precompute powers
   w = make([]int64, mid*2+3)
   w[0] = 1
   for i := 1; i <= mid*2; i++ {
       w[i] = w[i-1] * A % md
   }
   // special case k==1
   if k == 1 {
       if (w[1] + 2*w[2])%md == h {
           fmt.Println("1 2")
       } else if (2*w[1] + w[2])%md == h {
           fmt.Println("2 1")
       } else {
           fmt.Println(-1)
       }
       return
   }
   a = make([]int64, mid*2+3)
   b = make([]int64, mid*2+3)
   mp1 = make(map[int64]int)
   mp2 = make(map[int64]int)
   // build first half
   dfs1(1, 1, int64(mid+1), 0, 0)
   // search second half
   dfs2(1, mid+1, int64(mid+1), 0)
   if !flagFound {
       fmt.Println(-1)
       return
   }
   // output result
   out := bufio.NewWriter(os.Stdout)
   for i := 1; i <= mid*2; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       out.WriteString(fmt.Sprint(b[i]))
   }
   out.WriteByte('\n')
   out.Flush()
}
