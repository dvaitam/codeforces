#include<iostream>
#include<cstdio>
using namespace std;
const int N=1002019;
int n,tot,m,vis[N],cnt[2],ans[2][N],dep[N],v[N];
struct Edge{int to,next;}e[N<<1];
void add(int x,int y){
	e[++tot].to=y;
	e[tot].next=v[x]; v[x]=tot;
}
int read(){
	int x=0,f=1;char ch=getchar();
	while(ch<'0' || ch>'9'){if(ch=='-')f=-1;ch=getchar();}
	while(ch>='0' && ch<='9'){x=x*10+ch-'0';ch=getchar();}
	return x*f;
}
void dfs(int x){
	for(int p=v[x];p;p=e[p].next){
		int kp=e[p].to;
		if(!vis[kp]){
			vis[kp]=1;
			dep[kp]=dep[x]+1;
			dfs(kp);
		}
	}
}
int main()
{
	int T=read();
	while(T--){
		n=read(); m=read();
		for(int i=1;i<=m;i++){
			int u=read(),v=read();
			add(u,v); add(v,u);
		}dfs(1);
		for(int i=1;i<=n;i++){
			if(dep[i]&1)ans[0][++cnt[0]]=i;
			if(!(dep[i]&1))ans[1][++cnt[1]]=i;
		}
		bool pr=cnt[1]<cnt[0];
		printf("%d\n",cnt[pr]);
		for(int i=1;i<=cnt[pr];i++)
			printf("%d ",ans[pr][i]);
		putchar('\n');
		for(int i=1;i<=n;i++)vis[i]=dep[i]=v[i]=0;
		cnt[0]=cnt[1]=tot=0;
	}
	return 0;
}