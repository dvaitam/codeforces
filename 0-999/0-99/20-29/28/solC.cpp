#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <cmath>
#include <cstring>
#include <string>
#include <cctype>
#include <algorithm>
#define maxn 55

using namespace std;
typedef long long LL;

double ans,pre,f[maxn][maxn],p[maxn],s;
int n,m,a[maxn];
LL c[maxn][maxn];

void init()
 {
  scanf("%d%d",&n,&m);
  for (int i=1;i<=m;i++)scanf("%d",&a[i]);
 }
 
void work()
 {
  int i,j,k,kk;
  p[0]=1;
  for (i=1;i<=n;i++)p[i]=p[i-1]*1.0/m;
  c[0][0]=1;
  for (i=1;i<=n;i++)
   {
    c[i][0]=1;
    for (j=1;j<=n;j++)c[i][j]=c[i-1][j-1]+c[i-1][j];
   }
  for (kk=1;kk<=n;kk++)
   {
    for (i=0;i<=m;i++)
     for (j=0;j<=n;j++)f[i][j]=0;
    f[0][0]=1;
    for (i=0;i<m;i++)
     for (j=0;j<=n;j++)
      for (k=0;k<=a[i+1]*kk && j+k<=n;k++)f[i+1][j+k]+=f[i][j]*p[k]*c[n-j][k];
    ans+=(f[m][n]-pre)*kk;
    pre=f[m][n];
   }
  printf("%.20lf\n",ans);
 }
 
int main()
 {
  init();
  work();
  return 0;
 }