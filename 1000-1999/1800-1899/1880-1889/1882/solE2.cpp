// LUOGU_RID: 156905377
#include<bits/stdc++.h>
using namespace std;
const int N=2510;
int n,m,a[N],b[N],f1[N],f2[N],p[N];
bool vis[N];
void work(int x,int tp)
{
	if(tp==1) memset(vis,0,sizeof(vis));
	vis[x]=1;
	for(int y=p[x];y!=x;y=p[y]) vis[y]=1;
}
int calc(int n,int *a,int w)
{
	for(int i=0;i<=n;i++) p[i]=(a[i]+w)%(n+1);
	int ans=n;
	work(0,1);
	for(int i=1;i<=n;i++)
		if(p[i]==i) ans--;
		else if(!vis[i]) work(i,0),ans++;
	return ans;
}
vector<int> v,v1,v2;
void dfs(int x,int y)
{
	if(x==y) return;
	dfs(p[x],y);
	v.push_back(x-p[x]);
}
void solve(int n,int *a,int w)
{
	v.clear();
	for(int i=0;i<=n;i++) p[i]=(a[i]+w)%(n+1);
	int cur=0;
	work(0,1);
	for(int i=1;i<=n;i++)
		if(p[i]!=i&&!vis[i])
		{
			v.push_back(i-cur);
			swap(p[cur],p[i]),cur=i;
			work(cur,1);
		}
	dfs(p[cur],cur);
}
int F(int n,int x){return x>0?x:n+1+x;}
int main()
{
	scanf("%d%d",&n,&m);
	for(int i=1;i<=n;i++) scanf("%d",&a[i]);
	for(int i=1;i<=m;i++) scanf("%d",&b[i]);
	for(int i=0;i<=n;i++) f1[i]=calc(n,a,i);
	for(int i=0;i<=m;i++) f2[i]=calc(m,b,i);
	int ans=1e9;
	for(int i=0;i<=n;i++)
		for(int j=0;j<=m;j++)
			if((f1[i]^f2[j])%2==0)
				ans=min(ans,max(f1[i],f2[j]));
	for(int i=0;i<=n;i++)
		for(int j=0;j<=m;j++)
			if((f1[i]^f2[j])%2==0&&max(f1[i],f2[j])==ans)
			{
				solve(n,a,i),v1=v;
				solve(m,b,j),v2=v;
				while(v1.size()<v2.size()) v1.push_back(1),v1.push_back(n);
				while(v2.size()<v1.size()) v2.push_back(1),v2.push_back(m);
				printf("%d\n",(int)v1.size());
				for(int k=0;k<(int)v1.size();k++)
					printf("%d %d\n",F(n,v1[k]),F(m,v2[k]));
				return 0;
			}
	return puts("-1"),0;
}