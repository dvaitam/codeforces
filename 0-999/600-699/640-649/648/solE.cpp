#include<bits/stdc++.h>
#define ll long long
#define ull unsigned long long
#define rep(i, j, k) for(int i=j;i<=k;++i)
#define rep1(i, j, k) for(int i=k;i>=j;--i)
#define pii pair<int, int> 
#define pb push_back
#define ft first
#define sd second


using namespace std;

const int N=1e6+5, M=1e3+5, K=1e4+5, inf=1e9;
int n, k, a[N], c[K], tot, _10[12];
pii fro[M], b[N];
bool tun[12][M];

int read(){
	int s=0;
	char x=getchar();
	while(!isdigit(x)) x=getchar();
	while(isdigit(x)){
		s=(s<<1)+(s<<3)+x-'0';
		x=getchar();
	}
	return s;
}

void dfs(int u){
	if(!u) return;
	dfs(fro[u].ft);
	printf("%d", a[fro[u].sd]);
}

int dis[M];
bool vis[M];

vector<int> vec[K];
bool dij(){
	rep(i, 1, k) dis[i]=inf;
	vec[0].pb(0);
	int i=0, j=0, cnt=1;
	while(cnt){
		while(vec[i].size()<=j) ++i, j=0;
		--cnt;
		int u=vec[i][j];
		if(u==k) return 1;
		if(vis[u]){
			++j;
			continue;
		}
		vis[u]=1;
		rep(l, 1, tot){
			int v=(u*_10[b[c[l]].ft]+b[c[l]].sd)%k;
			if(!v) v=k;
			if(vis[v]) continue;
			if(dis[u]+b[c[l]].ft<dis[v]){
				dis[v]=dis[u]+b[c[l]].ft;
				fro[v]={u, c[l]};
				vec[dis[v]].pb(v);
				++cnt;
			}
		}
		++j;
	}
	return 0;
}

int main(){
	
	n=read(), k=read();
	
	_10[0]=1;
	rep(i, 1, 9) _10[i]=_10[i-1]*10%k;
	
	rep(i, 1, n){
		a[i]=read();
		int x=a[i], len=0;
		while(x) x/=10, ++len;
		if(!a[i]) len=1;
		b[i]={len, a[i]%k};
		if(!tun[len][a[i]%k]) tun[len][a[i]%k]=1, c[++tot]=i;
	}
	
	if(!dij()) puts("NO");
	else{
		puts("YES");
		dfs(k);
	}
	
	return 0;
}