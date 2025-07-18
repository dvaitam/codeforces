#include <cstdio>
#include <cstring>
#include <algorithm>
using std::sort;
using std::swap;
typedef long long ll;
const int N = 1e4+5;
const int C = 20;
const int inf = 1e9+7;
#define Reset(x, y) memset(x, y, sizeof(x))
#define Clean(x)  Reset(x, 0)
#define qwq register
#define OvO inline
// const int MAXIN = 1e6;
// char IN[MAXIN],*SS=IN,*TT=IN;
// #define getchar() (SS == TT&&(TT = (SS=IN) + fread(IN,1,MAXIN,stdin), SS==TT) ? EOF:*SS++)
template <typename T> OvO T max(T a, T b) { return a > b ? a : b; }
template <typename T> OvO T min(T a, T b) { return a < b ? a : b; }
template <typename T> OvO void read(T &x) {
    x = 0; T f = 1; char c = getchar();
    while (c<'0' || c>'9') { if (c=='-') f=-1; c = getchar(); }
    while (c>='0' && c<='9') { x = x*10+c-'0'; c = getchar(); }
    x *= f;
}
bool mark1;
int n, a[N];
bool mark2;
int main()
{
    // freopen("testin.txt", "r", stdin);
    // freopen("stdout.txt", "w", stdout);
    // freopen("testout.txt", "w", stdout);
    // printf("Memory:%lfMB\n", (&mark2-&mark1)/1000.0/1000.0);
    read(n);
    for (int i = 1; i <= n; ++i) read(a[i]);
    int day = 0, p = 0, goal = 0;
    while (p < n) {
        ++p; ++day;
        goal = max(goal, a[p]);
        while (p < goal) goal = max(goal, a[++p]);
    }
    printf("%d\n", day);
    return 0;
}