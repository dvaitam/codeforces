#include<iostream> 
#include<cstdio>
#include<cmath>
#include<cstdlib>
#include<cstring>
#include<algorithm>
using namespace std;
int read()
{
	int x=0,f=1;char c=getchar();
	while (c<'0'||c>'9') {if (c=='-') f=-1;c=getchar();}
	while (c>='0'&&c<='9') x=(x<<1)+(x<<3)+(c^48),c=getchar();
	return x*f;
}
#define N 200010
int n,m,a[N],ans=0;
int main()
{
	n=read(),m=read();
	for (int i=1;i<=n;i++) a[i]=read();
	sort(a+1,a+n+1);reverse(a+1,a+n+1);
	for (int i=1;a[i]!=a[n];)
	{
		int t=i;long long tot=0;
		while (t<n&&tot+1ll*(a[t]-a[t+1])*t<=m) tot+=(a[t]-a[t+1])*t,t++;
		ans++;
		i=t;a[i]-=(m-tot)/t;
	}
	cout<<ans;
	return 0;
}