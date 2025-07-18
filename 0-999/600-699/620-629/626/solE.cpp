#include <cstdio>
#include <algorithm>
#include <vector>
#include <cstring>
#include <iostream>

#define rep(i, l, r) for(int i=l; i<=r; i++)
#define down(i, l, r) for(int i=l; i>=r; i--)
#define travel(x) for(edge *p=fir[x]; p; p=p->n)
#define clr(x,c) memset(x, c, sizeof(x))
#define pb push_back
#define lowbit(x) (x&-x)
#define l(x) Left[x]
#define r(x) Right[x]
#define ll long long

#define maxn 200009

using namespace std;
inline int read()
{
	int x=0, f=1; char ch=getchar();
	while (ch<'0' || '9'<ch) {if (ch=='-') f=-1; ch=getchar();}
	while ('0'<=ch && ch<='9') x=x*10+ch-'0', ch=getchar();
	return x*f;
}

int n, k[maxn], g1, g2, g3, now, tot; ll all;
double ans;

int main()
{
	n=read(); rep(i, 1, n) k[i]=read(); sort(k+1, k+1+n);
	
	ans=0; g1=1, g2=1, g3=0;
	now=0; all=k[1], tot=1; rep(i, 2, n)
	{
		all-=k[i-now-1], all+=k[i];
		
		while (i-now-1>0 && n+1-now-1>i && 1LL*(k[i-now-1]+k[n+1-now-1])*tot>=2*all)
			now++, tot+=2, all+=k[i-now]+k[n+1-now];
		while (now && 1LL*(k[i-now]+k[n+1-now])*tot<2*all)
			all-=k[i-now]+k[n+1-now], tot-=2, now--;
		
		if (ans<1.0*all/tot-k[i])
			ans=1.0*all/tot-k[i], g1=1, g2=i, g3=now;
		
		if (n+1-now==i+1) all-=k[i-now]+k[n+1-now], tot-=2, now--;
	}
	
	now=0; all=k[1]+k[2], tot=2; rep(i, 2, n-1)
	{
		all-=k[i-now-1], all+=k[i+1];
		
		while (i-now-1>0 && n+1-now-1>i+1 && 1LL*(k[i-now-1]+k[n+1-now-1])*tot>=2*all)
			now++, tot+=2, all+=k[i-now]+k[n+1-now];
		while (now && 1LL*(k[i-now]+k[n+1-now])*tot<2*all)
			all-=k[i-now]+k[n+1-now], tot-=2, now--;
		
		if (ans<1.0*all/tot-1.0*(k[i]+k[i+1])/2)
			ans=1.0*all/tot-1.0*(k[i]+k[i+1])/2, g1=2, g2=i, g3=now;
		
		if (n+1-now==i+2) all-=k[i-now]+k[n+1-now], tot-=2, now--;
	}
	
	if (g1==1)
	{
		printf("%d\n", g3*2+1);
		down(i, g3, 1) printf("%d ", k[g2-i]);
		printf("%d", k[g2]);
		down(i, g3, 1) printf(" %d", k[n+1-i]);
		puts("");
	}
	else
	{
		printf("%d\n", g3*2+2);
		down(i, g3, 1) printf("%d ", k[g2-i]);
		printf("%d %d", k[g2], k[g2+1]);
		down(i, g3, 1) printf(" %d", k[n+1-i]);
		puts("");
	}
	return 0;
}