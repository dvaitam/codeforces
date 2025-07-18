#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <algorithm>
#include <map>
#include <set>
#include <vector>
#include <map>
#include <string>
#include <math.h>
#include <stdlib.h>
using namespace std;

#define forlr(i,l,r) for(i=l;i<r;i++)
#define forrl(i,r,l) for(i=r-1;i>=l;i--)
#define for0r(i,r) forlr(i,0,r)

#define forltr(i,l,r) for(i=l;i<=r;i++)
#define for1tr(i,r) forltr(i,1,r)
#define forrtl(i,r,l) for(i=r;i>=l;i--)

#define forcnt(n) for(int cnt=n; cnt; cnt--)

#define X(x) (x).first
#define Y(x) (x).second

#define mymemset(x, val) memset(x, val, sizeof(x))
#define sz(x) ((int)x.size())

typedef int LD;

int i,j,n,m;

int main() {
  scanf("%d %d", &n, &m);
  if(n<4) {
    printf("-1\n");
    return 0;
  }
  if(n==4) {
    if(m!=3) {
      printf("-1\n");
      return 0;
    }
  }
  if(m==2) {
    printf("%d\n", n-1);
    for0r(i, n-1) {
      printf("%d %d\n", i+1, i+2);
    }
    return 0;
  }
  if(m==3) {
    printf("%d\n", 3+(n-4)*2);
    for1tr(i, 3) {
      printf("%d %d\n", i, i+1);
    }
    for(i=5;i<=n;i++) {
      printf("%d %d\n", i, 1);
      printf("%d %d\n", i, 2);
    }
    return 0;
  }
  printf("-1\n");
  return 0;
}