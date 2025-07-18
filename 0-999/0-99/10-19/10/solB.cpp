#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <algorithm>
using namespace std;

#define N 2000
int n, k, p[N][2];

int main()
{
  int m, xm, lm, rm, dm;
  
  scanf("%d%d", &n, &k);
  for(int i = 1; i <= k; i++) {
    p[i][0] = 0;
    p[i][1] = -1;
  }
  int xc = (k+1)/2;
  for(int i = 0; i < n; i++) {
    scanf("%d", &m);
    if(m > k) {
      printf("-1\n");
      continue;
    }
    dm = -1;
    for(int x = 1; x <= k; x++) {
      if(p[x][0] > p[x][1]) {
	int d = abs(x-xc)*m + (m/2)*((m+1)/2);
	if(dm < 0 || d < dm) {
	  dm = d;
	  xm = x;
	  lm = (k-m)/2+1;
	  rm = lm+m-1;
	}
      } else {
	if(p[x][0] > m) {
	  int d = abs(x-xc)*m + (xc-p[x][0])*m + m*(m+1)/2;
	  if(dm < 0 || d < dm) {
	    dm = d;
	    xm = x;
	    lm = p[x][0] - m;
	    rm = p[x][0] - 1;
	  }
	}
	if(p[x][1] <= k - m) {
	  int d = abs(x-xc)*m + (p[x][1]-xc)*m + m*(m+1)/2;
	  if(dm < 0 || d < dm) {
	    dm = d;
	    xm = x;
	    lm = p[x][1] + 1;
	    rm = p[x][1] + m;
	  }
	}
      }
    }
    if(dm < 0) 
      printf("-1\n");
    else {
      printf("%d %d %d\n", xm, lm, rm);
      if(p[xm][0] > p[xm][1]) {
	p[xm][0] = lm;
	p[xm][1] = rm;
      } else if(p[xm][0] == rm+1) {
	p[xm][0] = lm;
      } else if(p[xm][1] == lm-1) {
	p[xm][1] = rm;
      }
    }
  }
  
  return 0;
}