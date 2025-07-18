#include <stdio.h>

#include <memory.h>

#define min(a , b) ((a) < (b) ? (a) : (b))

const int MAXBUF = 1<<22; char B[MAXBUF] , *S = B , *T = B;

#define getch() (S == T && (T = (S = B) + fread(B,1,MAXBUF,stdin) , S == T ) ? EOF : *S ++)

inline int F() { register int aa , bb , ch;

	while(ch = getch() , (ch<'0'||ch>'9') && ch != '-'); ch == '-' ? aa=bb=0 : (aa=ch-'0',bb=1);

	while(ch = getch() , ch>='0'&&ch<='9') aa = aa*10 + ch-'0'; return bb ? aa : -aa;

}

const int Maxn = 400010;

const int Maxm = 800010;

#define R register

int n , m , to[Maxm] , next[Maxm] , g[Maxn] , ecnt = 1 , pos[Maxn] , low[Maxn] , dfn = 0 , ans[Maxn][2] , q[Maxn] , sz;

bool used[Maxm] , vis[Maxn] , inq[Maxn];

inline void ins(register int a , register int b) {

	to[++ecnt] = b; next[ecnt] = g[a]; g[a] = ecnt;

	to[++ecnt] = a; next[ecnt] = g[b]; g[b] = ecnt;

}

void dfs1(int x , int fa) {

	pos[x] = low[x] = ++dfn;

	for(int i=g[x]; i; i=next[i]) {

		if(to[i] == fa) continue;

		if(!pos[to[i]]) {

			dfs1(to[i] , x);

			if(low[to[i]] > pos[x])

				used[i] = used[i^1] = 1;

			low[x] = min(low[x] , low[to[i]]);

		}

		else low[x] = min(low[x] , pos[to[i]]);

	}

}

void dfs2(int x) {

	vis[x] = 1; ++sz;

	for(int i=g[x]; i; i=next[i]) {

		if(used[i]) continue;

		ans[i>>1][0] = x , ans[i>>1][1] = to[i];

		if(!vis[to[i]]) dfs2(to[i]);

	}

}

void bfs(const int&x) {

	R int h , t;

	inq[q[h = t = 1] = x] = 1;

	while(h <= t) {

		R int now = q[h++];

		for(R int i=g[now]; i; i=next[i]) {

			if(inq[to[i]]) continue;

			if(used[i]) ans[i>>1][0] = to[i] , ans[i>>1][1] = now;

			inq[to[i]] = 1;

			q[++t] = to[i];

		}

	}

}

const int buff = 1<<23; char O[buff] , *oo = O;

#define putch(x) (*oo++ = x)

inline void P(int&x) {

	if(!x) *oo++ = '0';

	else {

		static int st[21]; register int top = 0;

		while(x) st[++top] = x%10,x/=10;

		while(top) *oo++ = st[top--] + '0';

	}

}

int main() {

	n = F() , m = F();

	for(R int i=1 , __u = m; i<=__u; ++i) ins(F() , F());

	dfs1(1,0);

	R int mx = 0 , mxl = 1;

	for(R int i=1 , __u = n; i<=__u; ++i) {

		if(vis[i]) continue;

		sz = 0;

		dfs2(i);

		if(sz > mx) mx = sz , mxl = i;

	}

	bfs(mxl);

	P(mx); putch('\n');

	for(R int i=1 , __u = m; i<=__u; ++i)

		P(ans[i][0]) , putch(' ') , P(ans[i][1]) , putch('\n');

	fwrite(O,1,oo-O,stdout);

}