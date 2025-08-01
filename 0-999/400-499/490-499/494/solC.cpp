#include <cstdio>
#include <algorithm>
#define N 100005
#define M 5005
#define FI(a, b, c) for(int a = (b); a <= (c); a++)
#define FD(a, b, c) for(int a = (b); a >= (c); a--)
#define fe(a, b, c) for(int a = (b); a; a = c[a])
using namespace std;

int n, m, a[N], f[N], v[N + N], o[N + N], p, b[M], lim[M], stk[M], ptr;
double dp[M][M], nd[M], ans;
struct query{
	int l, r;
	double p;
	bool operator < (const query &T) const{
		return l != T.l ? l < T.l : r > T.r;
	}
} q[M];

void bd(int a, int b){
	v[++p] = b; o[p] = f[a]; f[a] = p;
}

void dfs(int x){
	int r = q[x].r, mx = 0;
	dp[x][0] = 1;
	
	fe(i, f[x], o){
		FI(j, q[v[i]].r + 1, r) b[x] = max(b[x], a[j]);
		r = q[v[i]].l - 1;
	}
	FI(i, q[x].l, r) b[x] = max(b[x], a[i]);
	
	fe(i, f[x], o){
		dfs(v[i]);
		
		int rb = max(b[x], b[v[i]]);
		int rl = max(b[x] + lim[x] - rb, b[v[i]] + lim[v[i]] - rb);
		FI(j, 0, rl) nd[j] = 0;
		
		FI(j, 0, lim[v[i]]) FI(k, 0, lim[x]){
			int ind = max(b[v[i]] + j, b[x] + k) - rb;
			nd[ind] += dp[v[i]][j] * dp[x][k];
		}
		
		lim[x] = rl;
		b[x] = rb;
		FI(j, 0, rl) dp[x][j] = nd[j];
	}
	
	lim[x]++;
	dp[x][lim[x]] = 0;
	FD(i, lim[x], 1) dp[x][i] = dp[x][i] * (1 - q[x].p) + dp[x][i - 1] * q[x].p;
	dp[x][0] *= 1 - q[x].p;
}

int main(){
	scanf("%d %d", &n, &m);
	FI(i, 1, n) scanf("%d", &a[i]);
	FI(i, 1, m) scanf("%d %d %lf", &q[i].l, &q[i].r, &q[i].p);
	q[0].l = 1; q[0].r = n; q[0].p = 0;
	sort(q, q + m + 1);
	
	stk[ptr++] = 0;
	FI(i, 1, m){
		while(q[i].r > q[stk[ptr - 1]].r) ptr--;
		bd(stk[ptr - 1], i);
		while(i < m && q[i + 1].l == q[i].l && q[i + 1].r == q[i].r){
			bd(i, i + 1);
			i++;
		}
		stk[ptr++] = i;
	}
	
	dfs(0);
	
	FI(i, 0, lim[0]) ans += dp[0][i] * (b[0] + i);
	printf("%.6lf\n", ans);
	scanf("\n");
}