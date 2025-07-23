package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
      return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
      fmt.Fscan(reader, &a[i])
   }

   lmax := make([]int, n)
   rmax := make([]int, n)
   lmin := make([]int, n)
   rmin := make([]int, n)
   stack := make([]int, 0, n)

   // previous greater (strict)
   for i := 0; i < n; i++ {
      for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
         stack = stack[:len(stack)-1]
      }
      if len(stack) == 0 {
         lmax[i] = -1
      } else {
         lmax[i] = stack[len(stack)-1]
      }
      stack = append(stack, i)
   }
   // next greater or equal
   stack = stack[:0]
   for i := n - 1; i >= 0; i-- {
      for len(stack) > 0 && a[stack[len(stack)-1]] < a[i] {
         stack = stack[:len(stack)-1]
      }
      if len(stack) == 0 {
         rmax[i] = n
      } else {
         rmax[i] = stack[len(stack)-1]
      }
      stack = append(stack, i)
   }

   // previous less (strict)
   stack = stack[:0]
   for i := 0; i < n; i++ {
      for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
         stack = stack[:len(stack)-1]
      }
      if len(stack) == 0 {
         lmin[i] = -1
      } else {
         lmin[i] = stack[len(stack)-1]
      }
      stack = append(stack, i)
   }
   // next less or equal
   stack = stack[:0]
   for i := n - 1; i >= 0; i-- {
      for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
         stack = stack[:len(stack)-1]
      }
      if len(stack) == 0 {
         rmin[i] = n
      } else {
         rmin[i] = stack[len(stack)-1]
      }
      stack = append(stack, i)
   }

   var result int64
   for i := 0; i < n; i++ {
      leftMax := i - lmax[i]
      rightMax := rmax[i] - i
      leftMin := i - lmin[i]
      rightMin := rmin[i] - i
      contrib := int64(leftMax)*int64(rightMax) - int64(leftMin)*int64(rightMin)
      result += contrib * a[i]
   }
   fmt.Fprintln(writer, result)
}
