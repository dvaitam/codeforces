#include<cstdio>  
#include<iostream>  
#include<algorithm>  
#include<cstdlib>  
#include<cstring>
#include<string>
#include<climits>
#include<vector>
#include<cmath>
#include<map>
#define LL long long
#define Vdata vector <data>

using namespace std;

inline char nc(){
  static char buf[1000000],*p1=buf,*p2=buf;
  if (p1==p2) { p2=(p1=buf)+fread(buf,1,100000,stdin); if (p1==p2) return EOF; }
  return *p1++;
}

inline void read(int &x){
  char c=nc(),b=1;
  for (;!(c>='0' && c<='9');c=nc()) if (c=='-') b=-1;
  for (x=0;c>='0' && c<='9';x=x*10+c-'0',c=nc()); x*=b;
}

inline void read(LL &x){
  char c=nc(),b=1;
  for (;!(c>='0' && c<='9');c=nc()) if (c=='-') b=-1;
  for (x=0;c>='0' && c<='9';x=x*10+c-'0',c=nc()); x*=b;
}

inline void read(char &x){
  for (x=nc();!(x>='A' && x<='Z');x=nc());
}

int wt,ss[19];
inline void print(int x){
	if (x<0) x=-x,putchar('-'); 
	if (!x) putchar(48); else {for (wt=0;x;ss[++wt]=x%10,x/=10);for (;wt;putchar(ss[wt]+48),wt--);}
}
inline void print(LL x){
	if (x<0) x=-x,putchar('-');
	if (!x) putchar(48); else {for (wt=0;x;ss[++wt]=x%10,x/=10);for (;wt;putchar(ss[wt]+48),wt--);}
}

int n,m,f[100010],g[100010],h[100010];
map<int,int> p;

int main()
{
	read(n);
	for (int i=1;i<=n;i++)
		read(f[i]);
	int x,y;
	for (int i=1;i<=n;i++)
		if (f[i]!=f[f[i]]) {puts("-1");return 0;}
	int s=0;
	for (int i=1;i<=n;i++)
		if (p[f[i]]==0) p[f[i]]=++s,h[s]=f[i],g[i]=s;
		else g[i]=p[f[i]];
	print(s),putchar('\n');
	for (int i=1;i<n;i++)
		print(g[i]),putchar(' ');
	print(g[n]),putchar('\n');
	for (int i=1;i<s;i++)
		print(h[i]),putchar(' ');
	print(h[s]),putchar('\n');
	return 0;
}