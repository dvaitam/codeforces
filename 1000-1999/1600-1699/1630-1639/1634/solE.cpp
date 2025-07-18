#include<bits/stdc++.h>
#define ll long long
using namespace std;
struct pai{
	int x,xh;
};
bool operator <(const pai &x,const pai &y){return x.x<y.x;}
int m,n[100005];vector<pai> a[100005],b[100005];
struct tpi{
	int x,xh1,xh2;
};
bool operator <(const tpi &x,const tpi &y){return x.x<y.x;}
tpi ta[200005];int tot;
vector<int> ans[100005];
int lw[100005],bi[400005][3],bj[200005],bs=1;
void add(int u,int v,int w){bi[++bs][0]=lw[u],bi[lw[u]=bs][1]=v,bi[bs][2]=w;}
void pri(int i,int j,int x,int y){ans[i][x]=0,ans[j][y]=1;}
void dfs(int w,int st)
{
	for(;;)
	{
		int o_o=lw[w];
		if(!bj[o_o>>1])
		{
			int u=w,v=bi[o_o][1],w1=bi[o_o^1][2],w2=bi[o_o][2];bj[o_o>>1]=1;
			pri(u,v,w1,w2);
			if(v!=st)dfs(v,st);
			return;
		}
		lw[w]=bi[o_o][0];
	}
}
int main()
{
	scanf("%d",&m);
	for(int i=1;i<=m;i++)
	{
		scanf("%d",n+i);
		for(int j=0;j<n[i];j++)
		{
			int x;scanf("%d",&x);
			a[i].push_back((pai){x,j}),ans[i].push_back(0);
		}
		sort(a[i].begin(),a[i].end());
		for(int j=0;j<n[i];j++)
		if(j+1<n[i]&&a[i][j].x==a[i][j+1].x)
		{
			ans[i][a[i][j].xh]=0,ans[i][a[i][j+1].xh]=1;
			++j;
		}
		else b[i].push_back(a[i][j]),ta[++tot]=(tpi){a[i][j].x,i,a[i][j].xh};
	}
	sort(ta+1,ta+tot+1);
	for(int i=1;i<=tot;i+=2)
	{
		if(ta[i].x!=ta[i+1].x){printf("NO");return 0;}
		int u=ta[i].xh1,v=ta[i+1].xh1,w1=ta[i].xh2,w2=ta[i+1].xh2;
		add(u,v,w2),add(v,u,w1);
	}
	for(int i=1;i<=m;i++)
	{
		while(lw[i])
		{
			int o_o=lw[i];
			if(!bj[o_o>>1])
			{
				int u=i,v=bi[o_o][1],w1=bi[o_o^1][2],w2=bi[o_o][2];bj[o_o>>1]=1;
				pri(u,v,w1,w2),dfs(v,i);
			}
			lw[i]=bi[o_o][0];
		}
	}
	printf("YES\n");
	for(int i=1;i<=m;i++)
	{
		for(int j=0;j<n[i];j++)putchar(ans[i][j]==0?'L':'R');
		putchar('\n');
	}
}