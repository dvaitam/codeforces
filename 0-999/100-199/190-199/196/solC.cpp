#include <iostream>
#include <cstdio>
#include <memory.h>
#include <cstring>
#include <cmath>

#include <vector>
#include <deque>
#include <queue>
#include <stack>
#include <algorithm>

#define ABS(a) ((a)<0?-(a):(a))
#define SIGN(a) ((a)<0?-1:((a)>0?1:0))
#define SQR(a) ((a)*(a))
#define MAX(a,b) (((a)>(b))?(a):(b))
#define MIN(a,b) (((a)<(b))?(a):(b))

#define REP(i, n) for(int i=0; i<(n); ++i)
#define FOR(i, a, b) for(int i=(a); i<(b); ++i)

#define in ({int x;scanf("%d", &x);x;})

#define PI   (3.1415926)
#define INF  (2147483647)
#define INF2 (1073741823)
#define EPS  (0.000001)
#define y1 stupid_cmath

typedef long long LL;

using namespace std;

vector<int> g[2000];
int N, mm[2000][4], L=0, px, py, u[2000], mc[2000];

int dfs(int v){
    u[v]=1;
    int r=1;
    for(int i=0;i<g[v].size();++i)
        if(!u[g[v][i]])
            r+=dfs(g[v][i]);
    mc[v]=r;
    return r;
}

void sw(int a, int b){
    int t;
    t=mm[a][0]; mm[a][0]=mm[b][0]; mm[b][0]=t;
    t=mm[a][1]; mm[a][1]=mm[b][1]; mm[b][1]=t;
    t=mm[a][2]; mm[a][2]=mm[b][2]; mm[b][2]=t;
    t=mm[a][3]; mm[a][3]=mm[b][3]; mm[b][3]=t;
}

LL vect(LL ax, LL ay, LL bx, LL by){
    return ax*by-ay*bx;
}

void qs(int _l, int _r){
	int i=_l, j=_r;
	int _x=mm[(i+j)>>1][0], _y=mm[(i+j)>>1][1];

	do{
		while(vect(mm[i][0]-px, mm[i][1]-py, _x-px, _y-py)<0) ++i;
		while(vect(mm[j][0]-px, mm[j][1]-py, _x-px, _y-py)>0) --j;
		if(i<=j){
		    sw(i, j);
			++i;
			--j;
		}
	}while(i<=j);

	if(_l<j) qs(_l, j);
	if(i<_r) qs(i, _r);
}

void rec(int a, int v){
    mm[a][2]=v;
    u[v]=1;
    sw(L, a);
    px=mm[L][0];
    py=mm[L][1];
    qs(L+1, L+mc[v]-1);
    L++;
    for(int i=0;i<g[v].size();++i){
        if(!u[g[v][i]])
            rec(L, g[v][i]);
    }
}

int main(){
//	freopen("input.txt","r",stdin); freopen("output.txt","w",stdout);
	cin>>N;
	int i, a, b, m=0;
	for(i=0;i<N-1;++i){
        cin>>a>>b;
        g[a].push_back(b);
        g[b].push_back(a);
	}
	for(i=0;i<N;++i){
        cin>>mm[i][0]>>mm[i][1];
        mm[i][2]=0;
        mm[i][3]=i;
        if(mm[i][1] > mm[m][1]) m=i;
    }
    memset(u, 0, sizeof(u));
    memset(mc, 0, sizeof(mc));
    dfs(1);
    memset(u, 0, sizeof(u));
	rec(m, 1);
	int r[2000];
	for(i=0;i<N;++i) r[mm[i][3]]=mm[i][2];
	for(i=0;i<N;++i) cout<<r[i]<<" ";

	return 0;
}