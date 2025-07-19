package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n      int
   a1, a2 [][]uint64
   s1, s2 []int
   words  int
   writer *bufio.Writer
)

func R(x int) int { return x ^ 1 }

func SIM() {
   writer.WriteString("SIMILAR\n")
   writer.Flush()
   os.Exit(0)
}

func NONSIM(s []int) {
   // print complement assignment
   for i := 0; i < n; i++ {
       val := s[i]
       if val < 0 {
           val = 0
       }
       writer.WriteByte('0' + byte(val^1))
       if i+1 < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
   writer.Flush()
   os.Exit(0)
}

func readf(m int, a [][]uint64) {
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       bx := 0
       if x < 0 {
           bx = 1
           x = -x
       }
       x--
       xi := x*2 + bx
       by := 0
       if y < 0 {
           by = 1
           y = -y
       }
       y--
       yi := y*2 + by
       // implications: ¬xi -> yi, ¬yi -> xi
       from := R(xi)
       a[from][yi/64] |= 1 << (uint(yi) % 64)
       from2 := R(yi)
       a[from2][xi/64] |= 1 << (uint(xi) % 64)
   }
   // self reachability
   lim := n * 2
   for i := 0; i < lim; i++ {
       a[i][i/64] |= 1 << (uint(i) % 64)
   }
}

var reader *bufio.Reader

func scan() {
   var m1, m2 int
   fmt.Fscan(reader, &n, &m1, &m2)
   words = (2*n + 63) / 64
   size := 2 * n
   a1 = make([][]uint64, size)
   a2 = make([][]uint64, size)
   for i := 0; i < size; i++ {
       a1[i] = make([]uint64, words)
       a2[i] = make([]uint64, words)
   }
   s1 = make([]int, n)
   s2 = make([]int, n)
   for i := range s1 {
       s1[i] = -1
       s2[i] = -1
   }
   readf(m1, a1)
   readf(m2, a2)
}

func closure(a [][]uint64) {
   lim := n * 2
   for k := 0; k < lim; k++ {
       wk := k / 64
       bk := uint(k % 64)
       mask := uint64(1) << bk
       for i := 0; i < lim; i++ {
           if a[i][wk]&mask != 0 {
               // a[i] |= a[k]
               for w := 0; w < words; w++ {
                   a[i][w] |= a[k][w]
               }
           }
       }
   }
}

func take(a [][]uint64, s []int, x int) {
   idx := x / 2
   if s[idx] != -1 {
       return
   }
   s[idx] = x & 1
   lim := n * 2
   // propagate
   for i := 0; i < lim; i++ {
       if s[i/2] == -1 {
           if (a[x][i/64]>>(uint(i)%64))&1 != 0 {
               take(a, s, i)
           }
       }
   }
   // clear row x and x^1
   x1 := x ^ 1
   for w := 0; w < words; w++ {
       a[x][w] = 0
       a[x1][w] = 0
   }
   // clear columns x and x^1
   for i := 0; i < lim; i++ {
       // column x
       w := x / 64
       b := uint(x % 64)
       a[i][w] &^= 1 << b
       // column x1
       w1 := x1 / 64
       b1 := uint(x1 % 64)
       a[i][w1] &^= 1 << b1
   }
}

func reduct(a [][]uint64, s []int) bool {
   lim := n * 2
   for i := 0; i < lim; i += 2 {
       // unsat check
       // if a[i][i^1] and a[i^1][i]
       if ((a[i][(i^1)/64]>>(uint((i^1)%64)))&1) != 0 &&
           ((a[i^1][i/64]>>(uint(i%64)))&1) != 0 {
           return false
       }
   }
   for i := 0; i < lim; i += 2 {
       if ((a[i][(i^1)/64]>>(uint((i^1)%64)))&1) != 0 {
           take(a, s, i^1)
       } else if ((a[i^1][i/64]>>(uint(i%64)))&1) != 0 {
           take(a, s, i)
       }
   }
   return true
}

func solve_full(a [][]uint64, s []int) {
   lim := n * 2
   for i := 0; i < lim; i++ {
       if s[i/2] == -1 {
           take(a, s, i)
       }
   }
}

func solve() {
   scan()
   closure(a1)
   closure(a2)
   res1 := reduct(a1, s1)
   res2 := reduct(a2, s2)
   if !res1 && !res2 {
       SIM()
   } else if res1 != res2 {
       if res1 {
           solve_full(a1, s1)
           NONSIM(s1)
       } else {
           solve_full(a2, s2)
           NONSIM(s2)
       }
   }
   // check direct differences
   for i := 0; i < n; i++ {
       if s1[i] != -1 && s2[i] != -1 && s1[i] != s2[i] {
           solve_full(a1, s1)
           NONSIM(s1)
       }
   }
   // one assigned, one unassigned
   for i := 0; i < n; i++ {
       if s1[i] != s2[i] {
           if s1[i] == -1 {
               take(a1, s1, i*2+ (1 - s2[i]))
               solve_full(a1, s1)
               NONSIM(s1)
           } else {
               take(a2, s2, i*2+ (1 - s1[i]))
               solve_full(a2, s2)
               NONSIM(s2)
           }
       }
   }
   // find differing implication
   lim := n * 2
   for i := 0; i < lim; i++ {
       for j := 0; j < lim; j++ {
           if s1[i/2] == -1 && s1[j/2] == -1 {
               bi := (a1[i][j/64] >> (uint(j) % 64)) & 1
               bj := (a2[i][j/64] >> (uint(j) % 64)) & 1
               if bi != bj {
                   if bi == 1 {
                       take(a2, s2, i)
                       if s2[j/2] == -1 {
                           take(a2, s2, j^1)
                       }
                       solve_full(a2, s2)
                       NONSIM(s2)
                   } else {
                       take(a1, s1, i)
                       if s1[j/2] == -1 {
                           take(a1, s1, j^1)
                       }
                       solve_full(a1, s1)
                       NONSIM(s1)
                   }
               }
           }
       }
   }
   SIM()
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   solve()
}
