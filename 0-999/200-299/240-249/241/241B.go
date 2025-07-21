package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const MOD = 1000000007

// group defines a set of pairs: if cross is false, pairs are within A; if cross is true, pairs are combinations between A and B.
type group struct { A, B []int; cross bool }

// solveAllLower computes the sum of XOR lower bits (0..bit) for all pairs in groups.
func solveAllLower(groups []group, bit int) int64 {
   var sum int64
   for b := bit; b >= 0; b-- {
       bitw := int64(1) << uint(b)
       var next []group
       for _, g := range groups {
           if !g.cross {
               var z, o []int
               for _, v := range g.A {
                   if v>>uint(b)&1 == 1 {
                       o = append(o, v)
                   } else {
                       z = append(z, v)
                   }
               }
               if len(z) > 0 && len(o) > 0 {
                   sum = (sum + int64(len(z))*int64(len(o))%MOD*bitw) % MOD
                   next = append(next, group{A: z, B: o, cross: true})
               }
               if len(z) > 1 {
                   next = append(next, group{A: z, cross: false})
               }
               if len(o) > 1 {
                   next = append(next, group{A: o, cross: false})
               }
           } else {
               var a0, a1, b0, b1 []int
               for _, v := range g.A {
                   if v>>uint(b)&1 == 1 {
                       a1 = append(a1, v)
                   } else {
                       a0 = append(a0, v)
                   }
               }
               for _, v := range g.B {
                   if v>>uint(b)&1 == 1 {
                       b1 = append(b1, v)
                   } else {
                       b0 = append(b0, v)
                   }
               }
               if len(a0) > 0 && len(b1) > 0 {
                   sum = (sum + int64(len(a0))*int64(len(b1))%MOD*bitw) % MOD
                   next = append(next, group{A: a0, B: b1, cross: true})
               }
               if len(a1) > 0 && len(b0) > 0 {
                   sum = (sum + int64(len(a1))*int64(len(b0))%MOD*bitw) % MOD
                   next = append(next, group{A: a1, B: b0, cross: true})
               }
               if len(a0) > 0 && len(b0) > 0 {
                   next = append(next, group{A: a0, B: b0, cross: true})
               }
               if len(a1) > 0 && len(b1) > 0 {
                   next = append(next, group{A: a1, B: b1, cross: true})
               }
           }
       }
       groups = next
       if len(groups) == 0 {
           break
       }
   }
   return sum
}

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   // read n and m
   scanner.Scan()
   n, _ := strconv.Atoi(scanner.Text())
   scanner.Scan()
   m, _ := strconv.ParseInt(scanner.Text(), 10, 64)
   // read attractiveness values
   a := make([]int, n)
   for i := 0; i < n; i++ {
       scanner.Scan()
       ai, _ := strconv.Atoi(scanner.Text())
       a[i] = ai
   }
   // initial group
   groups := []group{{A: a, cross: false}}
   var ans int64
   // greedy by bits from high to low
   for b := 30; b >= 0 && m > 0; b-- {
       bitw := int64(1) << uint(b)
       var highCnt int64
       var highGroups, lowGroups []group
       // partition current groups
       for _, g := range groups {
           if !g.cross {
               var z, o []int
               for _, v := range g.A {
                   if v>>uint(b)&1 == 1 {
                       o = append(o, v)
                   } else {
                       z = append(z, v)
                   }
               }
               if len(z) > 0 && len(o) > 0 {
                   highCnt += int64(len(z)) * int64(len(o))
                   highGroups = append(highGroups, group{A: z, B: o, cross: true})
               }
               if len(z) > 1 {
                   lowGroups = append(lowGroups, group{A: z, cross: false})
               }
               if len(o) > 1 {
                   lowGroups = append(lowGroups, group{A: o, cross: false})
               }
           } else {
               var a0, a1, b0, b1 []int
               for _, v := range g.A {
                   if v>>uint(b)&1 == 1 {
                       a1 = append(a1, v)
                   } else {
                       a0 = append(a0, v)
                   }
               }
               for _, v := range g.B {
                   if v>>uint(b)&1 == 1 {
                       b1 = append(b1, v)
                   } else {
                       b0 = append(b0, v)
                   }
               }
               if len(a0) > 0 && len(b1) > 0 {
                   highCnt += int64(len(a0)) * int64(len(b1))
                   highGroups = append(highGroups, group{A: a0, B: b1, cross: true})
               }
               if len(a1) > 0 && len(b0) > 0 {
                   highCnt += int64(len(a1)) * int64(len(b0))
                   highGroups = append(highGroups, group{A: a1, B: b0, cross: true})
               }
               if len(a0) > 0 && len(b0) > 0 {
                   lowGroups = append(lowGroups, group{A: a0, B: b0, cross: true})
               }
               if len(a1) > 0 && len(b1) > 0 {
                   lowGroups = append(lowGroups, group{A: a1, B: b1, cross: true})
               }
           }
       }
       if highCnt >= m {
           ans = (ans + m*bitw) % MOD
           groups = highGroups
       } else {
           ans = (ans + highCnt*bitw) % MOD
           if highCnt > 0 {
               ans = (ans + solveAllLower(highGroups, b-1)) % MOD
           }
           m -= highCnt
           groups = lowGroups
       }
   }
   fmt.Println(ans)
}
