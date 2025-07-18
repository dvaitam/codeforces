#include<iostream> 
#include<cstdio>
#include<cmath>
#include<cstdlib>
#include<cstring>
#include<algorithm>
#include<queue>
#include<vector>
using namespace std;
#define ll long long
#define N 300010
char getc(){char c=getchar();while ((c<'A'||c>'Z')&&(c<'a'||c>'z')&&(c<'0'||c>'9')) c=getchar();return c;}
int gcd(int n,int m){return m==0?n:gcd(m,n%m);}
int read()
{
	int x=0,f=1;char c=getchar();
	while (c<'0'||c>'9') {if (c=='-') f=-1;c=getchar();}
	while (c>='0'&&c<='9') x=(x<<1)+(x<<3)+(c^48),c=getchar();
	return x*f;
}
int n,m;
struct data
{
	int x,y; 
	bool operator <(const data&a) const
	{
		return y<a.y;
	}
}a[N];
priority_queue<int,vector<int>,greater<int> > q;
signed main()
{
#ifndef ONLINE_JUDGE
	freopen("c.in","r",stdin);
	freopen("c.out","w",stdout);
#endif
	n=read(),m=read();
	for (int i=1;i<=n;i++) a[i].x=read(),a[i].y=read();
	sort(a+1,a+n+1);
	ll tot=0,ans=0;
	for (int i=n;i>=1;i--)
	{
		q.push(a[i].x);tot+=a[i].x;
		if (n-i+1>m) tot-=q.top(),q.pop();
		ans=max(ans,tot*a[i].y);
	}
	cout<<ans;
	return 0;
	//NOTICE LONG LONG!!!!!
}