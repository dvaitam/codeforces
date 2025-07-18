#include<cstdio>
#include<cstring>
#include<cstdlib>
#include<cmath>
#include<iostream>
#include<algorithm>
#include<queue>
#define maxn 200010

using namespace std;

struct yts
{
	int id,w,w1;
}b[maxn];

priority_queue<pair<int,int> > Q;

int n,m,k;
int d[maxn],d1[maxn],dp[maxn],q[maxn],dp1[maxn],D[maxn];
int head[maxn],to[maxn],nxt[maxn];
int head1[maxn],to1[maxn],nxt1[maxn];
int num,num1;
int a[maxn];
bool flag[maxn];

void addedge(int x,int y)
{
	d[y]++;
	num++;to[num]=y;nxt[num]=head[x];head[x]=num;
}

void addedge1(int x,int y)
{
	d1[y]++;
	num1++;to1[num1]=y;nxt1[num1]=head1[x];head1[x]=num1;
}

bool cmp(yts x,yts y) {return x.w1>y.w1;}

int main()
{
	//freopen("1.in","r",stdin);
	scanf("%d%d%d",&n,&m,&k);
	if (n<k) {puts("-1");return 0;}
	for (int i=1;i<=n;i++) {scanf("%d",&a[i]);flag[a[i]]=1;dp[i]=a[i];}
	for (int i=1;i<=m;i++)
	{
		int x,y;
		scanf("%d%d",&x,&y);
		addedge(y,x);addedge1(x,y);
	}
	int l=0,r=0;
	for (int i=1;i<=n;i++) if (!d[i]) q[++r]=i;
	for (int i=1;i<=n;i++) D[i]=d[i];
	for (int i=1;i<=n;i++) if (!dp[i]) dp[i]=1;
	while (l<r)
	{
		int x=q[++l];
		for (int p=head[x];p;p=nxt[p])
		{
			dp[to[p]]=max(dp[to[p]],dp[x]+1);
			d[to[p]]--;
			if (!d[to[p]]) q[++r]=to[p];
		}
	}
	if (r!=n) {puts("-1");return 0;}
	for (int i=1;i<=n;i++) if (a[i] && dp[i]!=a[i]) {puts("-1");return 0;}
	for (int i=1;i<=n;i++) if (dp[i]>k) {puts("-1");return 0;}
	l=r=0;
	for (int i=1;i<=n;i++)
	{
		if (!a[i]) dp1[i]=k; else dp1[i]=a[i];
		if (!d1[i]) q[++r]=i;
	}
	while (l<r)
	{
		int x=q[++l];
		for (int p=head1[x];p;p=nxt1[p])
		{
			dp1[to1[p]]=min(dp1[to1[p]],dp1[x]-1);
			d1[to1[p]]--;
			if (!d1[to1[p]]) q[++r]=to1[p];
		}
	}	
	r=0;
	for (int i=1;i<=n;i++) if (!a[i]) b[++r]=(yts){i,dp[i],dp1[i]};
	sort(b+1,b+r+1,cmp);
	int now=1;
	for (int i=k;i>=1;i--)
	  if (!flag[i])
	  {
	  	while (now<=r && b[now].w1>=i) Q.push(make_pair(b[now].w,b[now].id)),now++;
	  	int id=0;
	  	while (!Q.empty())
	  	{
	  		id=Q.top().second;
	  		if (dp[id]>i) Q.pop();
	  		else break;
	  	}
	  	if (Q.empty()) {puts("-1");return 0;}
	  	a[id]=i;Q.pop();
	  }
	
	l=r=0;
	for (int i=1;i<=n;i++) dp[i]=a[i],d[i]=D[i];
	for (int i=1;i<=n;i++) if (!d[i]) q[++r]=i;
	for (int i=1;i<=n;i++) if (!dp[i]) dp[i]=1;
	while (l<r)
	{
		int x=q[++l];
		for (int p=head[x];p;p=nxt[p])
		{
			dp[to[p]]=max(dp[to[p]],dp[x]+1);
			d[to[p]]--;
			if (!d[to[p]]) q[++r]=to[p];
		}
	}
	for (int i=1;i<n;i++) printf("%d ",dp[i]);printf("%d\n",dp[n]);
	return 0;
}