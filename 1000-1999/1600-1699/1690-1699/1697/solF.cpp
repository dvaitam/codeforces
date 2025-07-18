#include <bits/stdc++.h>
using namespace std;

#define Int register int
#define MAXN 500005

template <typename T> inline void read (T &t){t = 0;char c = getchar();int f = 1;while (c < '0' || c > '9'){if (c == '-') f = -f;c = getchar();}while (c >= '0' && c <= '9'){t = (t << 3) + (t << 1) + c - '0';c = getchar();} t *= f;}
template <typename T,typename ... Args> inline void read (T &t,Args&... args){read (t);read (args...);}
template <typename T> inline void write (T x){if (x < 0){x = -x;putchar ('-');}if (x > 9) write (x / 10);putchar (x % 10 + '0');}
template <typename T> inline void chkmax (T &a,T b){a = max (a,b);}
template <typename T> inline void chkmin (T &a,T b){a = min (a,b);}

int T,n,m,K,Siz;

vector <int> g[MAXN];
int getind (int i,int j,int k){
	return (j - 2) * n + i + k * Siz;
}
void link (int u,int v){
	g[u].push_back (v);
}

bool vis[MAXN];stack <int> S;
int ind,cnt,dfn[MAXN],low[MAXN],bel[MAXN];

void tarjan (int u){
	vis[u] = 1,S.push (u),dfn[u] = low[u] = ++ ind;
	for (Int v : g[u]){
		if (!dfn[v]) tarjan (v),chkmin (low[u],low[v]);
		else if (vis[v]) chkmin (low[u],dfn[v]);
	}
	if (dfn[u] == low[u]){
		++ cnt;
		while (1){
			int now = S.top();S.pop ();
			bel[now] = cnt,vis[now] = 0;
			if (now == u) break;
		}
	}
}

void linkit (int u,int v){
	link (u,v),link (v > Siz ? v - Siz : v + Siz,u > Siz ? u - Siz : u + Siz);
}

int val[MAXN];

signed main(){
	read (T);
	while (T --> 0){
		read (n,m,K),ind = cnt = 0,Siz = n * (K - 1);
		for (Int u = 1;u <= 2 * n * K;++ u) dfn[u] = low[u] = 0,g[u].clear ();
		for (Int u = 1;u < n;++ u)
			for (Int i = 2;i <= K;++ i) linkit (getind (u,i,0),getind (u + 1,i,0));
		for (Int u = 1;u <= n;++ u)
			for (Int i = 3;i <= K;++ i) linkit (getind (u,i,0),getind (u,i - 1,0));
		for (Int t = 1;t <= m;++ t){
			int opt,i,j,v;read (opt);
			if (opt == 1){
				read (i,v);
				if (v == K) link (getind (i,v,0),getind (i,v,1));
				else if (v == 1) link (getind (i,2,1),getind (i,2,0));
				else linkit (getind (i,v,0),getind (i,v + 1,0));
			}
			else if (opt == 2){
				read (i,j,v);if (i > j) swap (i,j);
				if (v <= K) link (getind (j,v,0),getind (j,v,1));
				for (Int p = 2;p <= v && p <= K;++ p)
					if (v - p + 1 <= K) linkit (getind (i,p,0),getind (j,max (p,v - p + 1),1));
			}
			else{
				read (i,j,v);if (i > j) swap (i,j);
				for (Int p = 2;p <= (v + 1 >> 1);++ p)
					if (v - p + 1 <= K) linkit (getind (i,p,1),getind (j,v - p + 1,0));
					else link (getind (i,p,1),getind (i,p,0));;
			}
		}
		for (Int u = 1;u <= 2 * Siz;++ u) if (!dfn[u]) tarjan (u);
		int lst;
		for (Int u = 1;u <= Siz;++ u) if (bel[u] == bel[u + Siz]){
			puts ("-1");
			goto there;
		}
		for (Int x = 1;x <= n;++ x){
			lst = 1;
			for (Int i = 2;i <= K;++ i) lst += bel[getind (x,i,0)] < bel[getind (x,i,1)];
			write (lst),putchar (' ');
		}
		putchar ('\n');
		there:;
	}
	return 0;
}
/*
错误的，偏激的，极右翼的，非马恩主义的，女权的，失败的，人民日报的，乐的
三措并举激活经济“大动脉”
英勇无畏的普京大帝！
深圳每套房可省5块钱！
失业中年人“表演”上班 ×
灵活就业中年人 √
考虑到疫情的不确定性，NOI科学委员会决定暂停NOI 2022网络同步赛一年。
欲华润百家，奈何贫贱不能移
比工资，心胸越比越狭窄。谈奉献，格局越谈越广阔。
像石榴籽一样紧紧抱在一起！
律师黄码无法出庭被撤诉
媒体:中国文艺需要不拘一格降人才 
-----------------------------------------------------------------
*/