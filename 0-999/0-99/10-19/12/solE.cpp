#include<cstdio>

#include<cstring>

#include<algorithm>

#include<cctype>

#include<ctime>

#include<cstdlib>

#include<string>

#include<queue>

#include<cmath>

#define mk(a,b) make_pair(a,b)

#define PDD pair<double,double>

#define rep(x,a,b) for (int x=a;x<=int(b);x++)

#define drp(x,a,b) for (int x=a;x>=int(b);x--)

#define cross(x,a) for (int x=hd[a];~x;x=nx[x])

#define ll long long

#define oo 2147483647

#define fr first

#define sc second

#define PII pair<int,int>

#define pb push_back

using namespace std;

inline int rd(){

	int x=0,ch=getchar(),f=1;

	while (!isdigit(ch)&&(ch!='-')&&(ch!=EOF)) ch=getchar();

	if (ch=='-'){f=-1;ch=getchar();}

	while (isdigit(ch)){x=(x<<1)+(x<<3)+ch-'0';ch=getchar();}

	return x*f;

}

inline void rt(ll x){

	if (x<0) putchar('-'),x=-x;

	if (x>=10) rt(x/10),putchar(x%10+'0');

		else putchar(x+'0');

}

const int maxn=1005;

int n,a[maxn][maxn];

int main(){

	n=rd();

	rep(i,2,n-1){

		int now=a[i-1][1];

		rep(j,1,i-1) now=now%(n-1)+1,a[i][j]=now;

	}

	a[n][1]=n-1;

	rep(j,2,n-1) a[n][j]=(a[n][j-1]%(n-1)+1)%(n-1)+1;

	rep(i,1,n){

		rep(j,1,i) rt(a[i][j]),putchar(' ');

		rep(j,i+1,n) rt(a[j][i]),putchar(' ');

		putchar('\n');

	}

}