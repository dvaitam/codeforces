#include<cstdio>
#include<cstring>
#include<algorithm>
#include<iostream>
#define ll long long
#define maxn 1000010
using namespace std;
inline int getint(){
	char c=getchar();int x=0;bool p=0; 
	while((c<'0'||c>'9')&&c!='-')c=getchar();
	if(c=='-')p=1,c=getchar();
	while(c>='0'&&c<='9')x=x*10+c-'0',c=getchar();
	if(p)x=-x;return x;
}
struct edge{int v,ne;}e[maxn];
int ans=0,tot,la[maxn],n,a[maxn];
ll sum[maxn],maxx=-1e18;
inline void add(int u,int v){e[++tot].ne=la[u];la[u]=tot;e[tot].v=v;} 
inline void dfs(int u,int fa){
	sum[u]=a[u];
	for(int i=la[u];i;i=e[i].ne){
		int v=e[i].v;
		if(v==fa)continue;
		dfs(v,u);
		if(sum[u]+sum[v]>sum[u])
		sum[u]=sum[v]+sum[u];
	}
}
inline void dfs2(int u,int fa){
	sum[u]=a[u];
	for(int i=la[u];i;i=e[i].ne){
		int v=e[i].v;
		if(v==fa)continue;
		dfs2(v,u);
		if(sum[u]+sum[v]>sum[u])sum[u]=sum[v]+sum[u];
	}
	if(sum[u]==maxx)sum[u]=0,++ans;
}
int main(){
	n=getint();
	for(int i=1;i<=n;++i)a[i]=getint();
	for(int i=1;i<n;++i){
		int u=getint(),v=getint();
		add(u,v),add(v,u);
	}
	dfs(1,0);
	for(int i=1;i<=n;++i)maxx=max(sum[i],maxx);
	dfs2(1,0);
	cout<<1ll*ans*maxx<<" "<<ans<<endl;
	return 0;
}