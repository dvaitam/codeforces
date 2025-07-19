package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n   int
   a   []int
   d   []int
   pArr []int
   writer *bufio.Writer
)

func max(x, y int) int {
   if x > y {
       return x
   }
   return y
}

func even(m int) int {
   if m <= 0 {
       return 0
   }
   if d[m] != 0 {
       return d[m]
   }
   if m == 2 {
       d[m] = max(a[0], a[1])
       pArr[m] = 0
       return d[m]
   }
   smallest := even(m-2) + max(a[m-2], a[m-1])
   k := m - 2
   sum := 0
   for i := m-2; i >= 0; i -= 2 {
       cur := even(i) + sum + max(a[i], a[m-1])
       sum += max(a[i], a[i-1])
       if cur < smallest {
           smallest = cur
           k = i
       }
   }
   sum = 0
   for i := m-3; i >= 0; i -= 2 {
       cur := odd(i) + sum + max(a[i], a[m-1])
       sum += max(a[i], a[i+1])
       if cur < smallest {
           smallest = cur
           k = i
       }
   }
   d[m] = smallest
   pArr[m] = k
   return smallest
}

func odd(m int) int {
   if m == 1 {
       // single customer
       return max(a[0], a[2])
   }
   if d[m] != 0 {
       return d[m]
   }
   smallest := even(m-1) + max(a[m-1], a[m+1])
   k := m - 1
   sum := 0
   for i := m-1; i >= 0; i -= 2 {
       cur := even(i) + sum + max(a[i], a[m+1])
       sum += max(a[i], a[i-1])
       if cur < smallest {
           smallest = cur
           k = i
       }
   }
   sum = 0
   for i := m-2; i >= 0; i -= 2 {
       cur := odd(i) + sum + max(a[i], a[m+1])
       sum += max(a[i], a[i+1])
       if cur < smallest {
           smallest = cur
           k = i
       }
   }
   d[m] = smallest
   pArr[m] = k
   return smallest
}

func printEven(m int) {
   if m <= 0 {
       return
   }
   k := pArr[m]
   if k%2 != 0 {
       printOdd(k)
       for i := k + 2; i+1 < m; i += 2 {
           fmt.Fprintf(writer, "%d %d\n", i+1, i+2)
       }
   } else {
       printEven(k)
       for i := k + 1; i+1 < m; i += 2 {
           fmt.Fprintf(writer, "%d %d\n", i+1, i+2)
       }
   }
   // final pair
   fmt.Fprintf(writer, "%d %d\n", k+1, m)
}

func printOdd(m int) {
   if m <= 0 {
       return
   }
   k := pArr[m]
   if k%2 != 0 {
       printOdd(k)
       for i := k + 2; i+1 <= m-1; i += 2 {
           fmt.Fprintf(writer, "%d %d\n", i+1, i+2)
       }
   } else {
       printEven(k)
       for i := k + 1; i+1 <= m-1; i += 2 {
           fmt.Fprintf(writer, "%d %d\n", i+1, i+2)
       }
   }
   // last single or pair
   if a[m+1] != 0 {
       fmt.Fprintf(writer, "%d %d\n", k+1, m+2)
   } else {
       fmt.Fprintf(writer, "%d\n", k+1)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   a = make([]int, n+3)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   d = make([]int, n+3)
   pArr = make([]int, n+3)
   if n%2 != 0 {
       res := odd(n)
       fmt.Fprintln(writer, res)
       printOdd(n)
   } else {
       res := even(n)
       fmt.Fprintln(writer, res)
       printEven(n)
   }
}
