/*=================================
# Created time: 2017-02-18 16:58
# Filename: d.cpp
# Description: 
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
#include <map> 
#include <stack>  
#include <queue>  
#include <vector> 
#include <bitset> 
 
#define pb push_back 
#define lb lower_bound 
#define ub upper_bound 
#define sqr(x) (x)*(x) 
#define lowbit(x) (x)&(-x)  
#define ls(a,b) (((a)+(b)) << 1)  
#define rs(a,b) (((a)+(b)) >> 1)  
#define fin(a) freopen(a,"r",stdin); 
#define fout(a) freopen(a,"w",stdout); 
#define Abs(x) ((x) > 0 ? (x) : (-(x)))  
#define FOR(i,a,b) for((i)=(a);(i)<=(b);(i)++)  
#define FORP(i,a,b) for(int i=(a);i<=(b);i++)  
#define FORM(i,a,b) for(int i=(a);i>=(b);i--)  
 
#define maxn 1000005 
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
int n, m, k;
struct P{
	int d, id;
	bool  operator < (const P&x)const{
		return d < x.d;
	}
}p[maxn];
int a[maxn], ans[maxn], cnt;
int main(){
	//fin("d.in");
	read(n); read(m); read(k);
	FORP(i, 1, n) read(a[i]);
	sort(a + 1, a + 1 + n);
	FORP(i, 1, m) {read(p[i].d); p[i].id = i;}//read(p[i].id); }
	int head = 1, head2 = 1, d = 0; //q.insert(Node{INF, INF});
	sort(p + 1, p + 1 + m);
	/*while (head <= tail || q.size() > 1){
		if (d < a[head]) {
			Node x = q.top(); int tot = 0;
			while (x.v <= a[head] && tot < k && q.s	ize() > k) 
				tot ++, ans[++cnt] = x.id, q.pop();
		}
		else {
			int tot = 0;
			while (a[head] >= d && head <= tail && tot <= k) head++;
			if (a[head + 1] == a[head]) 
		}
	}*/
	FORP(d, 0, n + m / k + 1){
		while (p[head2].d < d && head2 <= m) head2++;
		int num = 0;
		if (head <= n && a[head] < d) {puts("-1"); return 0;}
		while (num < k && (head2 <= m || head <= n)){
			P x = p[head2]; if (head2 > m) x.d = INF;
			if (x.d < a[head] || head > n) ans[++cnt] = x.id, head2++;
			else head++;
			num++;
		}
		if (head2 > m && head > n) break;
	}
	write(cnt,'\n');
	sort(ans + 1, ans + 1 + cnt);
	FORP(i, 1, cnt) write(ans[i],' ');
}