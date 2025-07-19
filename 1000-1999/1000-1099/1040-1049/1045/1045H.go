package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007
const Nmax = 200005

var jc, inv [Nmax]int

func add(x, y int) int {
   x += y
   if x >= mod {
       x -= mod
   }
   return x
}

func C(x, y int) int {
   if x < y || y < 0 {
       return 0
   }
   return int((int64(jc[x]) * int64(inv[y]) % mod) * int64(inv[x-y]) % mod)
}

func query(l, r, x int) int {
   if l < 0 {
       l = 0
   }
   if l > r {
       return 0
   }
   v := C(r+1, x+1) - C(l, x+1)
   if v < 0 {
       v += mod
   }
   return v
}

func calc(a []int, aN, c00, c01, c10, c11 int) int {
   an := len(a)
   if aN == 0 {
       return 0
   }
   if c00+c11 > aN {
       return 0
   }
   if c10 == 0 {
       if c11 < aN {
           return 1
       }
       if an == 1 {
           return 1
       }
       return 0
   }
   if c00+c11 < aN {
       return int(int64(C(c00-1, c10-1)) * int64(C(c11-1, c01-1)) % mod)
   }
   ans := 0
   for i := 1; i <= an; i++ {
       ai := a[i-1]
       if i&1 == 1 {
           if c01 == 1 {
               if c11 > ai {
                   return ans
               }
               return add(ans, 1)
           }
           t := int(int64(query(c11-ai, c11-2, c01-2)) * int64(C(c00-1, c10-1)) % mod)
           ans = add(ans, t)
           c11 -= ai
           c01--
           if c11 < c01 {
               return ans
           }
       } else {
           if c10 == 1 {
               if c00 < ai {
                   return ans
               }
               return add(ans, 1)
           }
           t := int(int64(query(0, c00-ai-2, c10-2)) * int64(C(c11-1, c01-1)) % mod)
           ans = add(ans, t)
           c00 -= ai
           c10--
           if c00 < c10 {
               return ans
           }
       }
   }
   return ans
}

func modPow(a, b int) int {
   res := 1
   for b > 0 {
       if b&1 == 1 {
           res = int(int64(res) * int64(a) % mod)
       }
       a = int(int64(a) * int64(a) % mod)
       b >>= 1
   }
   return res
}

func main() {
   // precompute factorials and inverses
   jc[0] = 1
   for i := 1; i < Nmax; i++ {
       jc[i] = int(int64(jc[i-1]) * int64(i) % mod)
   }
   inv[Nmax-1] = modPow(jc[Nmax-1], mod-2)
   for i := Nmax - 1; i > 0; i-- {
       inv[i-1] = int(int64(inv[i]) * int64(i) % mod)
   }

   reader := bufio.NewReader(os.Stdin)
   var s1, s2 string
   fmt.Fscan(reader, &s1)
   n := len(s1)
   // adjust first string
   la := 0
   for i := n - 1; i >= 0; i-- {
       if s1[i] == '1' {
           la = i
           break
       }
   }
   b1 := []byte(s1)
   if la == 0 {
       n--
       b1 = make([]byte, n)
       for i := range b1 {
           b1[i] = '1'
       }
   } else {
       b1[la] = '0'
       for i := la + 1; i < n; i++ {
           b1[i] = '1'
       }
   }
   // build segments a
   a := make([]int, 0, n)
   prev := 0
   for i := 1; i < n; i++ {
       if b1[i] != b1[i-1] {
           a = append(a, i-prev)
           prev = i
       }
   }
   a = append(a, n-prev)
   aN := n

   // second string
   fmt.Fscan(reader, &s2)
   n2 := len(s2)
   b2 := []byte(s2)
   bSegments := make([]int, 0, n2)
   prev = 0
   for i := 1; i < n2; i++ {
       if b2[i] != b2[i-1] {
           bSegments = append(bSegments, i-prev)
           prev = i
       }
   }
   bSegments = append(bSegments, n2-prev)
   bN := n2

   var c00, c01, c10, c11 int
   fmt.Fscan(reader, &c00, &c01, &c10, &c11)
   // transform counts
   c11 += c01 + 1
   c00 += c10
   c01++
   // save originals
   oc00, oc01, oc10, oc11 := c00, c01, c10, c11
   // validate
   if (c10 != c01 && c10 != c01-1) || c11 < c01 || c00 < c10 {
       fmt.Println(0)
       return
   }
   // first calc
   ans1 := calc(a, aN, c00, c01, c10, c11)
   // second calc
   ans2 := calc(bSegments, bN, oc00, oc01, oc10, oc11)
   ans := ans2 - ans1
   if ans < 0 {
       ans += mod
   }
   fmt.Println(ans)
}
