/*=================================
# Created time: 2019-01-11 23:07
# Filename: c.cpp
# Description: 
=================================*/
#include <bits/stdc++.h> 
 
#define fi first 
#define se second 
#define mp make_pair 
#define pb push_back 
#define lb lower_bound 
#define ub upper_bound 
#define sqr(x) (x)*(x) 
#define lowbit(x) (x)&(-x) 
#define pii pair<int, int> 
#define ls(a,b) (((a)+(b)) << 1) 
#define rs(a,b) (((a)+(b)) >> 1) 
#define fin(a) freopen(a,"r",stdin); 
#define fout(a) freopen(a,"w",stdout); 
#define FORP(i,a,b) for(int i=(a);i<=(b);i++) 
#define FORM(i,a,b) for(int i=(a);i>=(b);i--) 
 
#define maxn 100005 
#define maxm 100005 
#define MOD 1000000007 
#define INF 1070000000 
using namespace std; 
typedef long long ll; 
typedef unsigned long long ull; 
 
template<class T> inline   
void read(T& num){   
    num = 0; bool f = 1;char ch = getchar();   
    while(ch < '0' || ch > '9') { if(ch == '-') f = 0;ch = getchar();}   
    while(ch >= '0' && ch <= '9') {num = num * 10 + ch - '0';ch = getchar();}   
    num = f ? num: -num;   
}  
template<class T> inline 
void write(T x,char ch){ 
	int s[100]; 
	if (x==0) {putchar('0'); putchar(ch); return;} 
	if (x<0) {putchar('-'); x=-x;} int num=0; 
	while (x){ s[num++]=(x%10); x=x/10;} 
	FORM(i,num-1,0) putchar(s[i]+'0'); putchar(ch); 
} 
/*==================split line==================*/ 
const double pi = acos(-1); 
const double eps = 1e-8; 
struct seg
{
	int id, x, y;
	bool operator <(const seg &t) const{
		if (x == t.x) return y < t.y;
		else return x < t.x;
	}
}p[maxn];
int out[maxn];
int main(){ 
	//fin("c.in");
	int cas; read(cas);
	while (cas--){
		int n; read(n);
		int ans = 1;
		FORP(i, 1, n){
			int x, y; read(x); read(y);
			p[i].x = x, p[i].y = y; p[i].id = i;
		}
		sort(p + 1, p + 1 + n);
		int r = p[1].y;
		FORP(i, 2, n){
			if (p[i].x <= r) r = max(r, p[i].y);
			else {ans++; r = max(r, p[i].y);}
		}
		if (ans < 2) puts("-1");
		else {
			ans = 1; out[p[1].id] = 1;
			r = p[1].y;
			FORP(i, 2, n){
				if (p[i].x <= r) {out[p[i].id] = ans; r = max(r, p[i].y);}
				else {ans++; if (ans > 2) ans = 2; out[p[i].id] = ans; r = max(r, p[i].y);}
			}
			FORP(i, 1, n) write(out[i],' ');
			puts("");
		}
	}
}