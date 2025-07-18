#include <bits/stdc++.h>
using namespace std;
const int N=5005;
int t,n,m,rd[N];
long double f[N][N],dp[N],mp[N];
queue<int>q;
struct nod{
	struct node{int next,to;}a[200005];
	int h[N],cnt;
	void add(int x,int y){a[++cnt].to=y,a[cnt].next=h[x],h[x]=cnt;}
}a,b;
int main()
{
	ios::sync_with_stdio(0),cin.tie(0),cout.tie(0);
	for(int i=1;i<=5000;i++)
	{
		f[i][1]=1;
		for(int j=2;j<=i;j++) f[i][j]=(j-2)*f[i-2][j-2]+(i-j)*f[i-2][j-1];
		for(int j=1;j<=i;j++) f[i][j]/=i;
	}
	cin>>t;
	while(t--)
	{
		cin>>n>>m,a.cnt=b.cnt=0,dp[n]=1;
		for(int x,y,i=1;i<=m;i++) cin>>x>>y,a.add(y,x),b.add(x,y),rd[x]++;
		for(int i=1;i<=n;i++) if(!rd[i]) q.push(i);
		while(!q.empty())
		{
			int x=q.front(),now=0;q.pop();
			for(int i=b.h[x];i;i=b.a[i].next) mp[++now]=dp[b.a[i].to];
			sort(mp+1,mp+now+1);
			for(int i=1;i<=now;i++) dp[x]+=f[now][now-i+1]*mp[i];
			for(int i=a.h[x];i;i=a.a[i].next)
			{
				int k=a.a[i].to;
				if(!(--rd[k])) q.push(k);
			}
		}
		cout<<fixed<<setprecision(20)<<dp[1]<<'\n';
		for(int i=1;i<=n;i++) dp[i]=a.h[i]=b.h[i]=0;
	}
	return 0;
}