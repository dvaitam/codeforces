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
#define N 500010
int n,a[N],mn=1000000001;
long long ans;
int main()
{
	n=read();
	for (int i=1;i<=n;i++) a[i]=read(),ans+=abs(a[i]);
	if (n==1) {cout<<a[1];return 0;}
	bool flag=1;
	for (int i=1;i<=n;i++) if (a[i]<0) {flag=0;break;}
	if (flag)
	{
		for (int i=1;i<=n;i++)
		mn=min(mn,a[i]);
		ans-=mn<<1;
	}
	else
	{
		flag=1;
		for (int i=1;i<=n;i++) if (a[i]>0) {flag=0;break;}
		if (flag)
		{
			for (int i=1;i<=n;i++)
			mn=min(mn,abs(a[i]));
			ans-=mn<<1;
		}
	}
	cout<<ans;
	return 0;
}