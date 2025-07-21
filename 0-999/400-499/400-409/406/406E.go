package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   var m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   aFs := make([]int64, 0, m)
   bFs := make([]int64, 0, m)
   for i := 0; i < m; i++ {
       var s, f int64
       fmt.Fscan(in, &s, &f)
       if s == 0 {
           aFs = append(aFs, f)
       } else {
           bFs = append(bFs, f)
       }
   }
   a := int64(len(aFs))
   b := int64(len(bFs))
   // sort f-values
   sort.Slice(aFs, func(i, j int) bool { return aFs[i] < aFs[j] })
   sort.Slice(bFs, func(i, j int) bool { return bFs[i] < bFs[j] })
   // compute ranges
   var minA, maxA, rangeA int64
   if a > 0 {
       minA = aFs[0]
       maxA = aFs[len(aFs)-1]
       rangeA = maxA - minA
   }
   var minB, maxB, rangeB int64
   if b > 0 {
       minB = bFs[0]
       maxB = bFs[len(bFs)-1]
       rangeB = maxB - minB
   }
   // compute scores and counts for each pattern
   var S_AAA, S_BBB, S_AAB, S_ABB int64
   var cnt_AAA, cnt_BBB, cnt_AAB, cnt_ABB int64
   // homogeneous AAA
   if a >= 3 {
       S_AAA = 2 * rangeA
       cnt_AAA = countHomogeneous(a, minA, maxA, aFs)
   }
   // homogeneous BBB
   if b >= 3 {
       S_BBB = 2 * rangeB
       cnt_BBB = countHomogeneous(b, minB, maxB, bFs)
   }
   // mixed AAB
   if a >= 2 && b >= 1 {
       // precompute c2A, cminA, cmaxA
       c2A := comb2(a)
       // occurrences of minA and maxA
       var cminA, cmaxA int64
       for _, v := range aFs {
           if v == minA {
               cminA++
           }
       }
       for i := len(aFs) - 1; i >= 0 && aFs[i] == maxA; i-- {
           cmaxA++
       }
       // find best S_AAB and count
       // case Bs inside [minA, maxA] give S = 2*n
       insideCnt := int64(0)
       for _, f := range bFs {
           if f >= minA && f <= maxA {
               insideCnt++
           }
       }
       if insideCnt > 0 {
           S_AAB = 2 * n
           // count pairs covering inside Bs
           for _, f := range bFs {
               if f < minA || f > maxA {
                   continue
               }
               less := int64(sort.Search(len(aFs), func(i int) bool { return aFs[i] >= f }))
               greater := a - int64(sort.Search(len(aFs), func(i int) bool { return aFs[i] > f }))
               cnt_AAB += c2A - comb2(less) - comb2(greater)
           }
       } else {
           // Bs all on one side
           // find max fB < minA and count
           var best1, best2 int64 = -1 << 60, -1 << 60
           var cnt1, cnt2 int64
           for _, f := range bFs {
               if f < minA {
                   val := 2*n + 2*(f-minA)
                   if val > best1 {
                       best1 = val; cnt1 = 1
                   } else if val == best1 {
                       cnt1++
                   }
               } else if f > maxA {
                   val := 2*n + 2*(maxA-f)
                   if val > best2 {
                       best2 = val; cnt2 = 1
                   } else if val == best2 {
                       cnt2++
                   }
               }
           }
           // choose best among
           if best1 >= best2 {
               S_AAB = best1
               // pairs with fA1 = minA
               cntPairs := c2A - comb2(a-cminA)
               cnt_AAB = cnt1 * cntPairs
           }
           if best2 >= best1 {
               S_AAB = best2
               cntPairs := c2A - comb2(a-cmaxA)
               cnt_AAB += cnt2 * cntPairs
           }
       }
   }
   // mixed ABB
   if b >= 2 && a >= 1 {
       c2B := comb2(b)
       var cminB, cmaxB int64
       for _, v := range bFs {
           if v == minB {
               cminB++
           }
       }
       for i := len(bFs) - 1; i >= 0 && bFs[i] == maxB; i-- {
           cmaxB++
       }
       insideCnt := int64(0)
       for _, f := range aFs {
           if f >= minB && f <= maxB {
               insideCnt++
           }
       }
       if insideCnt > 0 {
           S_ABB = 2 * n
           for _, f := range aFs {
               if f < minB || f > maxB {
                   continue
               }
               less := int64(sort.Search(len(bFs), func(i int) bool { return bFs[i] >= f }))
               greater := b - int64(sort.Search(len(bFs), func(i int) bool { return bFs[i] > f }))
               cnt_ABB += c2B - comb2(less) - comb2(greater)
           }
       } else {
           var best1, best2 int64 = -1 << 60, -1 << 60
           var cnt1, cnt2 int64
           for _, f := range aFs {
               if f < minB {
                   val := 2*n + 2*(f-minB)
                   if val > best1 {
                       best1 = val; cnt1 = 1
                   } else if val == best1 {
                       cnt1++
                   }
               } else if f > maxB {
                   val := 2*n + 2*(maxB-f)
                   if val > best2 {
                       best2 = val; cnt2 = 1
                   } else if val == best2 {
                       cnt2++
                   }
               }
           }
           if best1 >= best2 {
               S_ABB = best1
               cntPairs := c2B - comb2(b-cminB)
               cnt_ABB = cnt1 * cntPairs
           }
           if best2 >= best1 {
               S_ABB = best2
               cntPairs := c2B - comb2(b-cmaxB)
               cnt_ABB += cnt2 * cntPairs
           }
       }
   }
   // find global max
   S_max := S_AAA
   if S_BBB > S_max {
       S_max = S_BBB
   }
   if S_AAB > S_max {
       S_max = S_AAB
   }
   if S_ABB > S_max {
       S_max = S_ABB
   }
   // sum counts
   var result int64
   if S_AAA == S_max {
       result += cnt_AAA
   }
   if S_BBB == S_max {
       result += cnt_BBB
   }
   if S_AAB == S_max {
       result += cnt_AAB
   }
   if S_ABB == S_max {
       result += cnt_ABB
   }
   fmt.Println(result)
}

// comb2 returns C(x,2)
func comb2(x int64) int64 {
   if x < 2 {
       return 0
   }
   return x * (x - 1) / 2
}

// comb3 returns C(x,3)
func comb3(x int64) int64 {
   if x < 3 {
       return 0
   }
   return x * (x - 1) * (x - 2) / 6
}

// countHomogeneous counts triples in fs slice achieving max range = maxVal-minVal
func countHomogeneous(count int64, minVal, maxVal int64, fs []int64) int64 {
   if minVal == maxVal {
       // all equal
       return comb3(count)
   }
   // count occurrences
   var cmin, cmax int64
   // fs is sorted
   // count min
   for _, v := range fs {
       if v == minVal {
           cmin++
       }
   }
   for i := len(fs) - 1; i >= 0; i-- {
       if fs[i] == maxVal {
           cmax++
       }
   }
   mid := count - cmin - cmax
   // inclusion-exclusion: C(count,3) - C(count-cmin,3) - C(count-cmax,3) + C(mid,3)
   return comb3(count) - comb3(count-cmin) - comb3(count-cmax) + comb3(mid)
}

func max(x, y int64) int64 {
   if x > y {
       return x
   }
   return y
}
