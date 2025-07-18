/*=================================

# Created time: 2016-11-07 15:26

# Filename: f.cpp

# Description: Learn from "wo_de_xiao_hao"

=================================*/

#define me AcrossTheSky&HalfSummer11  

#include <cstdio>  

#include <cmath>  

#include <ctime>  

#include <string>  

#include <cstring>  

#include <cstdlib>  

#include <iostream>  

#include <algorithm>  

  

#include <set> 

#include <stack>  

#include <queue>  

#include <vector>  

#include <map>



#define pb push_back 

#define lb lower_bound 

#define sqr(x) (x)*(x) 

#define lowbit(x) (x)&(-x)  

#define Abs(x) ((x) > 0 ? (x) : (-(x)))  

#define FOR(i,a,b) for((i)=(a);(i)<=(b);(i)++)  

#define FORP(i,a,b) for(int i=(a);i<=(b);i++)  

#define FORM(i,a,b) for(int i=(a);i>=(b);i--)  

#define ls(a,b) (((a)+(b)) << 1)  

#define rs(a,b) (((a)+(b)) >> 1)  

#define getlc(a) ch[(a)][0]  

#define getrc(a) ch[(a)][1]  

  

#define maxn 100005 

#define maxm 100005 

#define INF 1070000000  

using namespace std;  

typedef long long ll;  

typedef unsigned long long ull;  

  

template<class T> inline  

void read(T& num){  

    num = 0; bool f = true;char ch = getchar();  

    while(ch < '0' || ch > '9') { if(ch == '-') f = false;ch = getchar();}  

    while(ch >= '0' && ch <= '9') {num = num * 10 + ch - '0';ch = getchar();}  

    num = f ? num: -num;  

} 

int out[100]; 

template<class T> inline 

void write(T x,char ch){ 

	if (x==0) {putchar('0'); putchar(ch); return;} 

	if (x<0) {putchar('-'); x=-x;} 

	int num=0; 

	while (x){ out[num++]=(x%10); x=x/10;} 

	FORM(i,num-1,0) putchar(out[i]+'0'); putchar(ch); 

} 

/*==================split line==================*/ 

#define Rep(i,n) for(int i=1;i<=n;++i)



const int N = 1024;

int id[N][N], cx[N][N], cy[N][N], dx[N], dy[N], ans[N], nx, ny;



void dfs(int x[][N], int y[][N], int u, int v, int cu, int cv) { //还可以这么传参...学到了

	//点u, 点v,cu是点v没有的最小颜色,cv是点u没有的最小颜色, 那么我们要把它uv之间染上cv,就要在两边乱搞...

	int w = y[v][cv];

	if(!x[w][cu]) {

		x[w][cu] = v, y[v][cu] = w;

		x[w][cv] = 0;

	}

	else dfs(y, x, v, w, cv, cu); //交换过来搞... 好强大...

	x[u][cv] = v; y[v][cv] = u;

}



int main()

{

	int m;

	read(nx); read(ny); read(m);

	int u, v;

	Rep(i, m) {

		read(u); read(v);

		id[u][v] = i;

		++dx[u]; ++dy[v];

	}

	int d = max(*max_element(dx + 1, dx + 1 + nx), *max_element(dy + 1, dy + 1 + ny));

	write(d,'\n');

	FORP(u, 1, nx) 

		FORP(v, 1, ny) if(id[u][v]) {

			int c = 1;

			while(cx[u][c] || cy[v][c]) ++c; //选一个两边都没有的颜色.

			if(c <= d) {

				cx[u][c] = v; cy[v][c] = u;

				continue;

			}

			int cu = 1, cv = 1; //各选一个没有的颜色

			while(cy[v][cu]) ++cu;

			while(cx[u][cv]) ++cv;

			dfs(cx, cy, u, v, cu, cv);

	}

	FORP(u, 1, nx) FORP(c, 1, d) if(cx[u][c]) ans[id[u][cx[u][c]]] = c;

	FORP(i, 1, m) write(ans[i],' ');

	puts("");

	return 0;

}