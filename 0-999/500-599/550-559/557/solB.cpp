#include<cstdio>
#include<cstdlib>
#include<algorithm>
using namespace std;
typedef long long LL;
const int MAXN = 200100;
const int INF = 0x7777777f;
inline int read() {
    int x=0,f=1;char c=getchar();
    while(c<'0'||c>'9'){if(c=='-')f=-1;c=getchar();}
    while(c>='0'&&c<='9'){x=x*10+c-'0';c=getchar();}
    return x*f;
}/*====================Template====================*/
int n, m, a[MAXN];

int main(){
    n = read(); m = read();
    for (int i = 1; i <= 2*n; i++) a[i] = read();
    sort(a+1, a+1+2*n);
    
    double mid = double(a[n+1]);
    double a1 = double(a[1]);
    double mm = double(m);
    if (a1*2 <= mid) printf("%.7lf", min(mm, a1*n*3));
    else printf("%.7lf", min(mm, mid*3/2*n));
    
    return 0;
}