#include<cstdio>
#include<algorithm>
using namespace std;
#define rep(i,n) for (int i=1;i<=n;i++)
#define N 205
int w[N],d[N][N],f[N][N],b[N],c[N],n,v,x,y;
void dfs(int x,int y)
{
	rep(j,n) f[x][j]=w[d[x][j]]+v;
	rep(i,n) if (d[x][i]==1 && y-i){dfs(i,x);
		rep(j,n) f[x][j]+=min(f[i][j]-v,f[i][b[i]]);}
	b[x]=1; rep(j,n) if (f[x][j]<f[x][b[x]]) b[x]=j;
}
void prt(int x,int y,int z)
{c[x]=z; rep(i,n) if (d[x][i]==1 && y-i) prt(i,x,f[i][z]-v<f[i][b[i]]?z:b[i]);}
int main()
{
	scanf("%d%d",&n,&v); rep(i,n-1) scanf("%d",w+i);
	rep(i,n) rep(j,n) d[i][j]=(i!=j)<<22;
	rep(i,n-1) scanf("%d%d",&x,&y),d[x][y]=d[y][x]=1;
	rep(k,n) rep(i,n) rep(j,n) d[i][j]=min(d[i][j],d[i][k]+d[k][j]);
	dfs(1,0),printf("%d\n",f[1][b[1]]),prt(1,0,b[1]);
	rep(i,n) printf("%d%s",c[i],i<n?" ":"\n"); return 0;
}