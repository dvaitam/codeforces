#include <bits/stdc++.h>
#pragma GCC optimize("unroll-loops")
using namespace std;
mt19937 eng(chrono::steady_clock::now().time_since_epoch().count());
int rnd(int l, int r) { return uniform_int_distribution<int>(l, r)(eng); }
namespace Fread { const int SZ = 1 << 10; char buf[SZ], *S, *T;
	inline char getchar() { if (S == T) { T = (S = buf) + fread(buf, 1, SZ, stdin);
			if (S == T) return '\n'; } return *S++; } } 
namespace Fwrite { const int SZ = 1 << 10; char buf[SZ], *S = buf, *T = buf + SZ;
	inline void flush() { fwrite(buf, 1, S - buf, stdout); S = buf; }
	inline void putchar(char c) { *S++ = c; if (S == T) flush(); }
	struct NTR { ~ NTR() { flush(); }} ztr; } 
#define getchar Fread :: getchar
#define putchar Fwrite :: putchar
namespace Fast {
	struct Reader { template<typename T> Reader& operator >> (T& x) { 
			char c = getchar(); while (!isdigit(c)) c = getchar(); x = 0; 
			while (isdigit(c)) x = x * 10 + (c - '0'), c = getchar(); return *this; } } cin;
	struct Writer { template<typename T> Writer& operator << (T x) {
			static int sta[45]; int top = 0; while (x) { sta[++top] = x % 10; x /= 10; }
			while (top) { putchar(sta[top] + '0'); --top; } return *this; }
		Writer& operator << (char c) { putchar(c); return *this; } } cout; } 
#define cin Fast :: cin
#define cout Fast :: cout
const int MX=5001, MN=2001;
int T, u[MX], v[MX], pa[MN], val[MN];
vector<int> g[MN];
bool vis[MN];
inline int find(int x) { return pa[x] == x ? x : pa[x] = find(pa[x]); }
inline void merge(int x, int y) { x = find(x), y = find(y); if (x != y) pa[x] = y; }
struct Val { const bool operator () (const int &a, const int &b) const { return val[a] > val[b]; } };
void solve() {
	int N1, N2, M; cin>>N1>>N2>>M; int N = N1 + N2;
	for (int i = 1; i <= N; ++i) pa[i] = i;
	for (int i = 1; i <= M; ++i) { cin>>u[i]>>v[i]; if (find(u[i]) != find(v[i])) 
		g[u[i]].push_back(v[i]), g[v[i]].push_back(u[i]), merge(u[i], v[i]); }
	bool flag = false; if (N1 > N2) flag = true, swap(N1, N2);
	while (true) {
		for (int i = 1; i <= N; ++i) val[i] = rnd(1, 3*N);
		priority_queue<int, vector<int>, Val> Q;
		Q.push(rnd(1, N));
		for (int i = 1; i <= N; ++i) vis[i] = false;
		int need = N1;
		for (int u; !Q.empty(); ) {
			vis[u = Q.top()] = true; Q.pop();
			if ((--need) == 0) break;
			for (int v : g[u]) if (!vis[v]) Q.push(v);
		}
		for (int i = 1; i <= N; ++i) pa[i] = i;
		for (int i = 1; i <= M; ++i) if (!vis[u[i]] && !vis[v[i]]) merge(u[i], v[i]);
		int IV = -1; bool flag2 = true;
		for (int i = 1; flag2 && i <= N; ++i) if (!vis[i]) 
			{if (IV == -1) IV = find(i); else if (IV != find(i)) flag2 = false; }
		if (flag2) { vector<int> vec; vec.reserve(N2);	
			for (int i = 1; i <= N; ++i) 
				if(vis[i]==flag) vec.push_back(i);
				else cout<<i<<' ';			
			cout<<'\n'; for (int x : vec) cout<<x<<' '; 
			cout<<'\n'; break;
		}
	}
	if(T) for (int i = 1; i <= N; ++i) g[i].clear();
}
int main() {
	cin>>T;
	while (T--) solve();
}