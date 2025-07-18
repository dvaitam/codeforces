#include <iostream>
#include <cmath>
#include <cstring>
#include <algorithm>
#include <queue>
#include <cstdio>
#include <stack>
#include <map>
#include <string>
#include <set>
#include <unordered_set>
#include <iomanip>
#include <bitset>
#pragma GCC optimize(2)
#define eps 1e-5
#define mod 998244353
#define pi acos(-1)
#define MAXN 100005
#define ee 2.71828182845904523536
using namespace std;
bool Finish_read;
template<class T>inline void read(T &x) { Finish_read = 0; x = 0; int f = 1; char ch = getchar(); while (!isdigit(ch)) { if (ch == '-')f = -1; if (ch == EOF)return; ch = getchar(); }while (isdigit(ch))x = x * 10 + ch - '0', ch = getchar(); x *= f; Finish_read = 1; }
template<class T>inline void print(T x) { if (x / 10 != 0)print(x / 10); putchar(x % 10 + '0'); }
template<class T>inline void writeln(T x) { if (x<0)putchar('-'); x = abs(x); print(x); putchar('\n'); }
template<class T>inline void write(T x) { if (x<0)putchar('-'); x = abs(x); print(x); putchar(' '); }

int a[1005];
int b[1005];
int main()
{
    int n,s;
    scanf("%d%d",&n,&s);
    for (int i=1;i<=n;i++)
        scanf("%d",&a[i]);
    for (int i=1;i<=n;i++)
        scanf("%d",&b[i]);
    if (a[1]==0)
    {
        printf("NO\n");
        return 0;
    }
    if (a[s]==1)
    {
        printf("YES\n");
        return 0;
    }
    if (b[s]==1)
    {
        for (int i=s;i<=n;i++)
        {
            if (a[i]&&b[i])
            {
                printf("YES\n");
                return 0;
            }
        }
    }
    printf("NO\n");
    return 0;
}