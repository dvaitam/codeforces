#include<bits/stdc++.h>
using namespace std;
#define ll long long
const int inf=1e9+7;ll id,p3[35];
int n,tt,ans,a[35],b[35],c[35],o[35],mx[2000005],f[2000005];
map<pair<int,int>,int>mp;
void dfs(int x,int l,int d1,int d2,int t,int s,int o)
{
	if(x==l)
	{
		if(!o)
		{
			int &p=mp[make_pair(d1,d2)];if(!p)p=++tt;
			if(t>mx[p]){mx[p]=t;f[p]=s;}return;
		}
		int p=mp[make_pair(-d1,-d2)];if(!p)return;
		if(t+mx[p]>ans){ans=t+mx[p];id=1ll*p3[n-n/2]*f[p]+s;}
		return;
	}
	dfs(x+1,l,d1+a[x]-b[x],d2+b[x],t+a[x],s*3+2,o);
	dfs(x+1,l,d1+a[x],d2-c[x],t+a[x],s*3+1,o);
	dfs(x+1,l,d1-b[x],d2+b[x]-c[x],t,s*3,o);
}
int main()
{
	scanf("%d",&n);p3[0]=1;ans=-inf;
	for(int i=1;i<2000005;i++)mx[i]=-inf;
	for(int i=1;i<=25;i++)p3[i]=p3[i-1]*3;
	for(int i=0;i<n;i++)scanf("%d%d%d",&a[i],&b[i],&c[i]);
	if(n==1){if(!a[0]&&!b[0])puts("LM");else if(!a[0]&&!c[0])puts("LW");else if(!b[0]&&!c[0])puts("MW");else puts("Impossible");return 0;}
	dfs(0,n/2,0,0,0,0,0);dfs(n/2,n,0,0,0,0,1);
	if(ans==-inf){puts("Impossible");return 0;}
	for(int i=n;i>=1;i--)o[i]=id%3,id/=3;
	for(int i=1;i<=n;i++)puts((o[i]==0)?"MW":(o[i]==1?"LW":"LM"));
	return 0;
}