package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxN = 600005

var (
   a    [maxN]int
   t    [maxN]int
   bArr [maxN]int
)

func add(n, x int) {
   for i := x; i <= n; i += i & -i {
       t[i]++
   }
}

func ask(x int) int {
   s := 0
   for i := x; i > 0; i -= i & -i {
       s += t[i]
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       var k int64
       fmt.Fscan(reader, &n, &k)
       // read permutation and init
       for i := 1; i <= n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           a[x] = i
           t[i] = 0
           bArr[i] = 0
       }
       // initial inversion count
       tot := int64(n*(n-1) / 2)
       sum := tot
       for val := 1; val <= n; val++ {
           sum -= int64(ask(a[val] - 1))
           add(n, a[val])
       }
       k -= sum
       if k%2 != 0 || k/2 < 0 || k/2 > (tot-sum) {
           fmt.Fprintln(writer, "NO")
           continue
       }
       k /= 2
       // reset Fenwick tree
       for i := 1; i <= n; i++ {
           t[i] = 0
       }
       // build result
       for iVal := 1; iVal <= n; iVal++ {
           xCnt := ask(a[iVal] - 1)
           if k > int64(xCnt) {
               k -= int64(xCnt)
               add(n, a[iVal])
               continue
           }
           if k > 0 && k <= int64(xCnt) {
               id := 0
               for j := 1; j <= iVal; j++ {
                   if a[j] < a[iVal] && k > 0 {
                       id = j
                       k--
                   }
               }
               for j := 1; j <= id; j++ {
                   bArr[j] = iVal + 1 - j
               }
               for j := id + 1; j <= iVal-1; j++ {
                   bArr[j] = iVal - j
               }
               bArr[iVal] = iVal - id
               add(n, a[iVal])
               continue
           }
           bArr[iVal] = iVal
           add(n, a[iVal])
       }
       fmt.Fprintln(writer, "YES")
       for i := 1; i <= n; i++ {
           if i > 1 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, bArr[i])
       }
       fmt.Fprintln(writer)
   }
}
