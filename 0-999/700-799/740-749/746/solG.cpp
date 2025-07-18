#define me AcrossTheSky&HalfSummer11  
#include <cstdio>  
#include <cmath>  
#include <ctime>  
#include <string>  
#include <cstring>  
#include <cstdlib>  
#include <iostream>  
#include <algorithm>  

#include <map>
#include <set> 
#include <stack>  
#include <queue>  
#include <vector>  
 
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
  
#define maxn 200005
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
int t[maxn], node[maxn], num[maxn], lowp[maxn], cap[maxn];
int sum, Max;
int main(){
	int n, T, k; read(n); read(T); read(k);
	FORP(i, 1, T) read(t[i]);
	FORP(i, 1, T) cap[i] = t[i] - 1, lowp[i] = max(0, t[i] - t[i + 1]), sum += lowp[i], Max += cap[i];
	Max -= cap[T];
	sum -= lowp[T];
	
	if (sum > k - t[T] || Max < k - t[T]) {puts("-1"); return 0;}
	sum = k - t[T] - sum;
	for (int i = 1; i <= T && sum ; i ++)
		num[i] = lowp[i] + min(sum, cap[i] - lowp[i]), sum -= min(sum, cap[i] - lowp[i]);
	node[1] = 1;
	int cnt = 1, l = 1, r = 1, nod = 1; 
	write(n,'\n');
	for (int i = 1; i <= T; i++){
		int temp = cnt, templ = l;
		FOR(nod, temp + 1, temp + t[i]) {
			write(nod, ' '); write(node[l],'\n'); l++; if (l > r) l = templ;
			node[++cnt] = nod;
		}
		r = cnt, l = cnt - (t[i] - num[i]) + 1;
	}
}