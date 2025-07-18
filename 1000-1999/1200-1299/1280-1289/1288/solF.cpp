// LUOGU_RID: 92279710
#include <bits/stdc++.h>
using namespace std;

#define Int register int
#define inf 0x3f3f3f3f
#define int long long
#define MAXM 50005
#define MAXN 5005

template <typename T> void read (T &x){char c = getchar ();x = 0;int f = 1;while (c < '0' || c > '9') f = (c == '-' ? -1 : 1),c = getchar ();while (c >= '0' && c <= '9') x = x * 10 + c - '0',c = getchar ();x *= f;}
template <typename T,typename ... Args> void read (T &x,Args& ... args){read (x),read (args...);}
template <typename T> void write (T x){if (x < 0) x = -x,putchar ('-');if (x > 9) write (x / 10);putchar (x % 10 + '0');}
template <typename T> void chkmax (T &a,T b){a = max (a,b);}

struct edge{
	int v,flw,cost,nxt;
}e[MAXM << 1];

int toop = 1,head[MAXN];
int dis[MAXN],cur[MAXN];bool vis[MAXN];
bool spfa (int s,int t){
	memcpy (cur,head,sizeof (head)),memset (dis,0x3f,sizeof (dis)),memset (vis,0,sizeof (vis));queue <int> q;dis[s] = 0,vis[s] = 1,q.push (s);
	int st0 = dis[t];
	while (!q.empty()){
		int u = q.front();q.pop(),vis[u] = 0;
		for (Int i = head[u];i;i = e[i].nxt){
			int v = e[i].v,flw = e[i].flw,cst = e[i].cost;
			if (flw && dis[v] > dis[u] + cst){
				dis[v] = dis[u] + cst;
				if (!vis[v]) vis[v] = 1,q.push (v); 
			}
		}
	}
	return dis[t] != st0;
}

int dfs (int u,int flow,int t,int &ans){
	if (u == t) return ans += flow * dis[t],flow;
	int rst = flow,tmp = flow;vis[u] = 1;
	for (Int i = cur[u];i && rst;i = e[i].nxt){
		cur[u] = i;
		int v = e[i].v,flw = e[i].flw,cst = e[i].cost;
		if (flw && dis[v] == dis[u] + cst && !vis[v]){
			int had = dfs (v,min (rst,flw),t,ans);
			if (!had) dis[v] = -inf;
			e[i].flw -= had,e[i ^ 1].flw += had,rst -= had;
			if (!rst) break;
		}
	}
	return tmp - rst;
}

#define pii pair<int,int>
#define se second
#define fi first
pii Maxflow (int s,int t){
	int ans = 0,flw = 0;
	while (spfa (s,t)) flw += dfs (s,inf,t,ans);
	return {flw,ans};
}

int deg[MAXN];
void addit (int u,int v,int l,int r,int c){
	deg[u] -= l,deg[v] += l;
	e[++ toop] = edge{v,r - l,c,head[u]},head[u] = toop;
	e[++ toop] = edge{u,0,-c,head[v]},head[v] = toop;
}

int n1,n2,m,R,B;char str[MAXN];

signed main(){
	read (n1,n2,m,R,B);
	int s = n1 + n2 + 1,t = s + 1,
		S = t + 1,T = S + 1;
	scanf ("%s",str + 1);
	for (Int i = 1;i <= n1;++ i) 
		if (str[i] == 'R') addit (s,i,1,inf,0);
		else if (str[i] == 'B') addit (i,t,1,inf,0);
		else addit (s,i,0,inf,0),addit (i,t,0,inf,0);
	scanf ("%s",str + 1);
	for (Int i = 1;i <= n2;++ i) 
		if (str[i] == 'R') addit (n1 + i,t,1,inf,0);
		else if (str[i] == 'B') addit (s,n1 + i,1,inf,0);
		else addit (s,n1 + i,0,inf,0),addit (n1 + i,t,0,inf,0);
	int sht = toop;
	for (Int i = 1,u,v;i <= m;++ i) read (u,v),addit (u,v + n1,0,1,R),addit (v + n1,u,0,1,B);
	int sum = 0;
	addit (t,s,0,inf,0);
	for (Int i = 1;i <= t;++ i){
		if (deg[i] > 0) sum += deg[i],addit (S,i,0,deg[i],0);
		if (deg[i] < 0) addit (i,T,0,-deg[i],0);
	}
	pii it = Maxflow (S,T);
	if (it.fi < sum) return puts ("-1") & 0;
	write (it.se),putchar ('\n');
	for (Int i = sht + 1;i <= sht + 4 * m;i += 4){
		if (!e[i].flw) putchar ('R');
		else if (!e[i + 2].flw) putchar ('B');
		else putchar ('U');
	}
	putchar ('\n');
	return 0;
}