package main

import (
   "bufio"
   "os"
   "strconv"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   // increase buffer for large inputs
   scanner.Buffer(make([]byte, 1024), 1e9)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   nextInt := func() int {
       if !scanner.Scan() {
           return 0
       }
       val, _ := strconv.Atoi(scanner.Text())
       return val
   }

   t := nextInt()
   for ; t > 0; t-- {
       n := nextInt()
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           a[i] = int64(nextInt())
       }
       // reverse slice
       for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
           a[i], a[j] = a[j], a[i]
       }
       for i := 0; i < n; i++ {
           var v int64
           if 2*i < n {
               v = -a[i]
           } else {
               v = a[i]
           }
           writer.WriteString(strconv.FormatInt(v, 10))
           if i+1 < n {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}
