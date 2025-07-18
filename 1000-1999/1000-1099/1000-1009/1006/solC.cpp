#include <bits/stdc++.h>
using namespace std;
#define REP(i, a, b) for (int i = (a), i##_end_ = (b); i <= i##_end_; ++i)
#define RPD(i, b, a) for (int i = (b), i##_end_ = (a); i >= i##_end_; --i)
#define pii pair<int, int>
#define PB push_back
#define SZ(x) (int)((x).size())

typedef long long LL;
const int oo = 0x3f3f3f3f;
const int MAXN = 300010;

inline int read()
{
    char c = getchar(); int x = 0, f = 1;
    while(c < '0' || c > '9'){if(c == '-') f = -1; c = getchar();}
    while(c >= '0' && c <= '9'){x = x * 10 + (c - '0'); c = getchar();}
    return x * f;
}

int a[MAXN];

int main()
{
    int n;
    while(cin >> n){
        REP(i, 1, n) a[i] = read();
        LL suml = a[1], sumr = a[n], ans = 0;
        int pl = 1, pr = n;
        while(pl < pr){
            if(suml < sumr) suml += a[++ pl];
            else if(sumr < suml) sumr += a[-- pr];
            else{
                ans = suml;
                suml += a[++ pl];
                sumr += a[-- pr];
            }
        }
        //cout << pl << ' ' << pr << endl;
        cout << ans << endl;
    }
    return 0;
}