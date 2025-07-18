#include<cstdio>
#include<iostream>
#define maxn 200005
#define pii pair<int,int>
using namespace std;
struct Edge{
	int v,next;
}edge[maxn*2];
int head[maxn],cnt;
void add(int u,int v){
	edge[++cnt]=(Edge){v,head[u]};head[u]=cnt;
}
int ind[maxn];
pii f[maxn][2];
void dfs(int u,int fa){
	f[u][0]=make_pair(0,1);
	f[u][1]=make_pair(1,ind[u]);
	for(int i=head[u];i;i=edge[i].next){
		int v=edge[i].v;
		if(v==fa) continue;
		dfs(v,u);
		f[u][1].first+=f[v][0].first;
		f[u][1].second+=f[v][0].second;
		if(f[v][0].first>f[v][1].first){
			f[u][0].first+=f[v][0].first;
			f[u][0].second+=f[v][0].second;
		}
		else if(f[v][0].first==f[v][1].first){
			f[u][0].first+=f[v][0].first;
			f[u][0].second+=min(f[v][0].second,f[v][1].second);
		}
		else if(f[v][0].first<f[v][1].first){
			f[u][0].first+=f[v][1].first;
			f[u][0].second+=f[v][1].second;
		}
	}
}
int col[maxn];
void work(int u,int fa,int flag){
	if(flag){
		col[u]=ind[u];
	}
	else col[u]=1;
	for(int i=head[u];i;i=edge[i].next){
		int v=edge[i].v;
		if(v==fa) continue;
		if(flag) work(v,u,0);
		else{
			if(f[v][0].first>f[v][1].first){
				work(v,u,0);
			}
			else if(f[v][0].first==f[v][1].first){
				if(f[v][0].second<f[v][1].second)	
					work(v,u,0);
					else work(v,u,1);
			}
			else if(f[v][0].first<f[v][1].first){
				work(v,u,1);
			}	
		}
	}
}
int main(){
	int n;
	scanf("%d",&n);
	for(int i=1;i<n;i++){
		int u,v;
		scanf("%d%d",&u,&v);
		add(u,v);add(v,u);
		ind[u]++;ind[v]++;
	}
	if(n==2){
		puts("2 2");
		puts("1 1");
		return 0;
	}
	dfs(1,0);
	if(f[1][0].first>f[1][1].first){
		printf("%d %d\n",f[1][0].first,f[1][0].second);
		work(1,0,0);
	}
	else if(f[1][0].first==f[1][1].first){
		printf("%d %d\n",f[1][0].first,min(f[1][0].second,f[1][1].second));
		if(f[1][0].second<f[1][1].second)
			work(1,0,0);
		else work(1,0,1);
	}
	else if(f[1][0].first<f[1][1].first){
		printf("%d %d\n",f[1][1].first,f[1][1].second);
		work(1,0,1);
	}
	for(int i=1;i<=n;i++){
		printf("%d ",col[i]);
	}
	return 0;
}