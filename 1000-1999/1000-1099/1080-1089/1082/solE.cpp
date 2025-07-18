#include<iostream>
#include<cstring>
#include<cstdio>
#define N 500010
using namespace std;
int f[N],l[N],a[N],s[N];
int read()
{
	int x=0;
	char c=getchar();
	while(c<'0'||c>'9') c=getchar();
	while(c>='0'&&c<='9') x=x*10+c-'0',c=getchar();
	return x;
}
int n,c,nw=0,ans=0;
int main()
{
	n = read(),c = read();
	for (int i=1;i<=n;i++)
	{
		a[i] = read();
		s[i]=s[i-1]+(int)(a[i]==c);
	}
	for(int i=1;i<=n;i++)
	{
		int cmx = max(f[l[a[i]]]+1,s[i-1]+1);
		f[i]=max(cmx,s[l[a[i]]-1]+1+(l[a[i]]!=0));
		l[a[i]]=i;
	}
	for(int i=n;i;i--)
	{
		ans=max(ans,f[i]+nw);
		nw+=(a[i]==c);
	}
	cout<<ans;
}