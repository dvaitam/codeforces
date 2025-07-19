package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func ncount(fs1, ls1, fs2, ls2, num1, num2, k, x int) int {
   // simulate sequence counts until k
   if num2 > x {
       return -1
   }
   // sequence 1: (fs1, ls1, count1=num1), sequence 2: (fs2, ls2, count2=num2)
   count1, count2 := num1, num2
   oldFs, oldLs := fs1, ls1
   newFs, newLs := fs2, ls2
   for it := 2; it < k; it++ {
       // next sequence = old + new
       nextCount := count1 + count2
       if oldLs == 0 && newFs == 2 {
           nextCount++
       }
       // rotate
       count1, count2 = count2, nextCount
       oldFs, oldLs, newFs, newLs = newFs, newLs, oldFs, newLs
       if count2 > x {
           return -1
       }
   }
   return count2
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var k, x, n, m int
   if _, err := fmt.Fscan(in, &k, &x, &n, &m); err != nil {
       return
   }
   // try all configurations
   for fs1 := 0; fs1 < 3; fs1++ {
       for ls1 := 0; ls1 < 3; ls1++ {
           for fs2 := 0; fs2 < 3; fs2++ {
               for ls2 := 0; ls2 < 3; ls2++ {
                   for num1 := 0; num1 <= n/2; num1++ {
                       for num2 := 0; num2 <= m/2; num2++ {
                           // validity checks
                           ok := true
                           if n == 1 && fs1 != ls1 {
                               ok = false
                           }
                           if m == 1 && fs2 != ls2 {
                               ok = false
                           }
                           if n > 1 && n%2 == 0 && num1 == n/2 && (fs1 != 0 || ls1 != 2) {
                               ok = false
                           }
                           if m > 1 && m%2 == 0 && num2 == m/2 && (fs2 != 0 || ls2 != 2) {
                               ok = false
                           }
                           if n > 1 && n%2 == 1 && num1 == n/2 && (fs1 != 0 && ls1 != 2) {
                               ok = false
                           }
                           if m > 1 && m%2 == 1 && num2 == m/2 && (fs2 != 0 && ls2 != 2) {
                               ok = false
                           }
                           if n == 2 && fs1 == 0 && ls1 == 2 && num1 != 1 {
                               ok = false
                           }
                           if m == 2 && fs2 == 0 && ls2 == 2 && num2 != 1 {
                               ok = false
                           }
                           if !ok {
                               continue
                           }
                           if ncount(fs1, ls1, fs2, ls2, num1, num2, k, x) == x {
                               // build and output strings
                               // first string
                               rem1, rem2 := num1, num2
                               var sb1 strings.Builder
                               switch fs1 {
                               case 0:
                                   sb1.WriteByte('A')
                               case 1:
                                   sb1.WriteByte('X')
                               case 2:
                                   sb1.WriteByte('C')
                               }
                               for i := 1; i < n-1; i++ {
                                   if rem1 > 0 {
                                       prev := sb1.String()[i-1]
                                       if prev == 'A' {
                                           sb1.WriteByte('C')
                                           rem1--
                                       } else {
                                           sb1.WriteByte('A')
                                       }
                                   } else {
                                       sb1.WriteByte('X')
                                   }
                               }
                               if n > 1 {
                                   switch ls1 {
                                   case 0:
                                       sb1.WriteByte('A')
                                   case 1:
                                       sb1.WriteByte('X')
                                   case 2:
                                       sb1.WriteByte('C')
                                   }
                               }
                               fmt.Println(sb1.String())
                               // second string
                               var sb2 strings.Builder
                               switch fs2 {
                               case 0:
                                   sb2.WriteByte('A')
                               case 1:
                                   sb2.WriteByte('X')
                               case 2:
                                   sb2.WriteByte('C')
                               }
                               for i := 1; i < m-1; i++ {
                                   if rem2 > 0 {
                                       prev := sb2.String()[i-1]
                                       if prev == 'A' {
                                           sb2.WriteByte('C')
                                           rem2--
                                       } else {
                                           sb2.WriteByte('A')
                                       }
                                   } else {
                                       sb2.WriteByte('X')
                                   }
                               }
                               if m > 1 {
                                   switch ls2 {
                                   case 0:
                                       sb2.WriteByte('A')
                                   case 1:
                                       sb2.WriteByte('X')
                                   case 2:
                                       sb2.WriteByte('C')
                                   }
                               }
                               fmt.Println(sb2.String())
                               return
                           }
                       }
                   }
               }
           }
       }
   }
   fmt.Println("Happy new year!")
}
