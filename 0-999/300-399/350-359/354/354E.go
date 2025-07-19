package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

const (
   PERMS  = 28
   DIGITS = 19
)

var p [PERMS][7]int
var bucket [10][]int
var pow10 [DIGITS]int64
var a [6]int64

func next(i int) int {
   if i == 0 {
       return 4
   } else if i == 4 {
       return 7
   }
   return 8
}

func search(n int64, dig int) bool {
   if n == 0 {
       return true
   }
   if dig >= DIGITS {
       return false
   }
   d := int((n / pow10[dig]) % 10)
   for _, idx := range bucket[d] {
       sumDigit := int64(p[idx][6])
       if n < sumDigit*pow10[dig] {
           continue
       }
       for j := 0; j < 6; j++ {
           a[j] += int64(p[idx][j]) * pow10[dig]
       }
       if search(n - sumDigit*pow10[dig], dig+1) {
           return true
       }
       for j := 0; j < 6; j++ {
           a[j] -= int64(p[idx][j]) * pow10[dig]
       }
   }
   return false
}

func main() {
   // Prepare permutations
   pid := 0
   for i := 0; i <= 7; i = next(i) {
       for j := i; j <= 7; j = next(j) {
           for k := j; k <= 7; k = next(k) {
               for l := k; l <= 7; l = next(l) {
                   for m := l; m <= 7; m = next(m) {
                       for n2 := m; n2 <= 7; n2 = next(n2) {
                           sum := i + j + k + l + m + n2
                           p[pid][0], p[pid][1], p[pid][2] = i, j, k
                           p[pid][3], p[pid][4], p[pid][5] = l, m, n2
                           p[pid][6] = sum
                           bucket[sum%10] = append(bucket[sum%10], pid)
                           pid++
                       }
                   }
               }
           }
       }
   }
   // Prepare powers of 10
   pow10[0] = 1
   for i := 1; i < DIGITS; i++ {
       pow10[i] = pow10[i-1] * 10
   }

   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var nVal int64
       fmt.Fscan(reader, &nVal)
       for i := 0; i < 6; i++ {
           a[i] = 0
       }
       if search(nVal, 0) {
           for i := 0; i < 6; i++ {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               writer.WriteString(strconv.FormatInt(a[i], 10))
           }
           writer.WriteByte('\n')
       } else {
           writer.WriteString("-1\n")
       }
   }
}
