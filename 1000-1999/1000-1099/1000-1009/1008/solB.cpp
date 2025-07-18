#include <iostream>
#include <cstdio>
#include <cstring>
#include <algorithm>
#include <queue>
#include <cmath>
#include <map>
#define ll long long
#define out(a) printf("%d ",a)
#define writeln printf("\n")
#define N 100050
using namespace std;
int n;
int now,a,b;
bool flag;
int read()
{
	int s=0,t=1; char c;
	while (c<'0'||c>'9'){if (c=='-') t=-1; c=getchar();}
	while (c>='0'&&c<='9'){s=s*10+c-'0'; c=getchar();}
	return s*t;
}
ll readl()
{
	ll s=0,t=1; char c;
	while (c<'0'||c>'9'){if (c=='-') t=-1; c=getchar();}
	while (c>='0'&&c<='9'){s=s*10+c-'0'; c=getchar();}
	return s*t;
}
int main()
{
	n=read(); flag=true;
	for (int i=1;i<=n;i++){
	  a=read(); b=read();
	  if (i==1) now=max(a,b);
	  else {
	    if (now>=a&&now>=b) now=max(a,b);
	    else if (now>=a&&now<b) now=a;
	    else if (now>=b&&now<a) now=b;
	    else if (now<a&&now<b) {
	      flag=false;
	      break;
	  }
    }
  }
    if (flag) puts("YES");
    else puts("NO");
	return 0;
}