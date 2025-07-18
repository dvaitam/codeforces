#include <bits/stdc++.h>
#pragma comment(linker, "/stack:200000000")
#pragma GCC optimize("Ofast")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")
#include <bits/stdc++.h>
using namespace std;
typedef long long LL;
#define lop(i,a,b) for(register int i = (a); i <= (b); ++i)
#define dlop(i,a,b) for(register int i = (a); i >= (b); --i)
//#define getchar() (*p1++)
//char buf[50<<20], *p1=buf;
inline int read(){
	register int c = getchar(), x = 0, f = 1;
	while(!isdigit(c)) {if (c == '-') f = -1; c = getchar();}
	while(isdigit(c)) x = (x<<3)+(x<<1)+(c^48), c = getchar();
	return x * f;
}
const int MAXN = 2e5+7;
int dg[MAXN];

int n, pre, flag;
int main(void){
	n = read();
	for(int i = 1; i < n; ++i) {
		int x = read(), y = read();
		++dg[x], ++dg[y];
	}
	for(int i = 1; i <= n; ++i) {
		if (dg[i] > 2) {
			if (flag) {
				puts("No");
				return 0;
			}
			flag = 1;
		}
	}
	for(int i = 1; i <= n; ++i) 
		if (dg[i] > dg[pre]) pre = i;
	printf("Yes\n%d\n", dg[pre]);
	for(int i = 1; i <= n; ++i) {
		if (pre == i) continue;
		if (dg[i] == 1) printf("%d %d\n", pre, i);
	}
	return 0; 
}