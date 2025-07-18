#include <cstdio>
#include <queue>
#include <algorithm>
//#define INF 0x3FFFFFF
#define MN 1005
using namespace std;
const int fx[4][2]={{1,0},{0,1},{0,-1},{-1,0}};
int n,ans[MN][2],pin,hr[MN],nex[MN],to[MN],d[MN];

inline void ins(int x,int y) {++pin; nex[pin]=hr[x]; to[pin]=y; hr[x]=pin;}
inline int abs(int x) {return x<0?-x:x;}
inline int read()
{
	int n=0,f=1; char c=getchar();
	while (c<'0' || c>'9') {if(c=='-')f=-1; c=getchar();}
	while (c>='0' && c<='9') {n=n*10+c-'0'; c=getchar();}
	return n*f;
}

//inline void dal(char ch,int x,int y)
//{
//	if (ch=='#'||ch=='*'||ch=='&') dis[x][1]=min(dis[x][1],y);
//	else if (ch>='a' && ch<='z') dis[x][2]=min(dis[x][2],y);
//	else if (ch>='0' && ch<='9') dis[x][3]=min(dis[x][3],y);
//}
//bool pd()
//{
//	for (int i=1;i<=n;++i) if (a[i]) return false;
//	return true;
//}
//bool cmp(int x,int y) {return b[x]<b[y] || b[x]==b[y] && di[x]<0 && di[y]>0;}

inline void dfs(int x,int y,int fa,int lf)
{
	int i,j;
	for (i=hr[x],j=0;i;i=nex[i])
	{
		if (to[i]==fa) continue;
		if (lf==j) ++j;
		ans[to[i]][0]=ans[x][0]+fx[j][0]*y;
		ans[to[i]][1]=ans[x][1]+fx[j][1]*y;
		dfs(to[i],y>>1,x,3-j); ++j;
	}
}

int main()
{
//	freopen("test.in","r",stdin);
	register int i,x,y;
	n=read();
	for (i=1;i<n;++i)
	{
		x=read(); y=read(); ++d[x]; ++d[y];
		ins(x,y); ins(y,x);
	}
	for (i=1;i<=n;++i) if (d[i]>4) return 0*printf("NO");
	ans[1][0]=ans[1][1]=0; dfs(1,1<<30,0,-1);
	printf("YES\n");
	for (i=1;i<=n;++i) printf("%d %d\n",ans[i][0],ans[i][1]);
//	dfs(1,0); if (u) return 0*printf("NO");
}