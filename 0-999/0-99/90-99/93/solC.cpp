#include <algorithm>
#include <cstring>
#include <cstdlib>
#include <cstdio>
#define For(i,x,y) for (i=x;i<=y;i++)
using namespace std;
const int b[4]={1,2,4,8};
struct ww {
	int a,b,c,d;
} a[100];
int i,j,k,n,m,an;
inline void print() {
	printf("%d\n",an-1);
	int i;
	For(i,2,an) {
		printf("lea e%cx, [",i+'a'-1);
		if (a[i].b) printf("e%cx + ",a[i].b+'a'-1);
		if (a[i].d) printf("%d*",b[a[i].d]);
		printf("e%cx]\n",a[i].c+'a'-1);
	}
	exit(0);
}
void dfs(int x) {
	if (x==an) {
		if (a[x].a==n) print();
		return;
	}
	int i,j,k;
	For(i,0,x)For(j,1,x)For(k,0,3) {
		int A=a[i].a+b[k]*a[j].a;
		if (A<=a[x].a||A>n) continue;
		a[x+1]=(ww){A,i,j,k};
		dfs(x+1);
	}
}
int main() {
	scanf("%d",&n);
	for (an=1;;an++) {
		a[1].a=1;
		dfs(1);
	}
	return 0;
}