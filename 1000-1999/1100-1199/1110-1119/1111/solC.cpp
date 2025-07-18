#include<iostream>
#include<cstdio>
#include<cmath>
#include<string>
#include<cstring>
#include<cstdlib>
#include<algorithm>
#include<vector>
#include<queue>
#include<stack>
#include<set>
#include<map>
using namespace std;
#define isNum(a) (a>='0'&&a<='9')
#define SP putchar(' ')
#define EL putchar('\n')
#define inf 2147483647
#define N 100005
#define File(a) freopen(a".in", "r", stdin), freopen(a".out", "w", stdout)
template<class T>inline void read(T&);
template<class T>inline void write(const T&);
typedef long long int ll;
ll work(int, int, int, int);
int p[N];
int A, B;
int main () {
    int n, k;
    read(n), read(k), read(A), read(B);
    for (int i=1; i<=k; ++i) {
        read(p[i]);
    }
    sort(p+1, p+k+1);
    write(work(n, 1, k, 1)), EL;
    return 0;
}
template<class T>void read(T &Re) {
    T k=0;
    char ch=getchar();
    int flag=1;
    while (!isNum(ch)) {
        if (ch=='-') {
            flag=-1;
        }
        ch=getchar();
    }
    while (isNum(ch)) {
        k=(k<<1)+(k<<3)+ch-'0';
        ch=getchar();
    }
    Re=flag*k;
}
template<class T>void write(const T& Wr) {
    if (Wr<0) {
        putchar('-');
        write(-Wr);
    }else {
        if (Wr<10) {
            putchar(Wr+'0');
        }else {
            write(Wr/10);
            putchar((Wr%10)+'0');
        }
    }
}
ll work(int l, int beg, int end, int st) {
    if (beg>end) {
        // cout<<l<<' '<<st<<endl;
        return A;
    }
    if (l==0) {
        // cout<<st<<endl;
        return (end-beg+1)*1ll*B;
    }
    int m=st+(1<<(l-1));
    int k=lower_bound(p+beg, p+end+1, m)-p;
    // cout<<l<<' '<<beg<<' '<<end<<' '<<st<<' '<<m<<' '<<k<<endl;
    return min((end-beg+1)*1ll*(1<<l)*B, work(l-1, beg, k-1, st)+work(l-1, k, end, m));
    // cout<<mini<<endl;
    // return mini;
}