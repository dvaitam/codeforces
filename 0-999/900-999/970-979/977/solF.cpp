#pragma comment(linker, "/STACK:1024000000,1024000000")
#include<algorithm>
#include<iostream>
#include<cstring>
#include<fstream>
#include<bitset>
#include<cstdio>
#include<string>
#include<vector>
#include<cmath>
#include<queue>
#include<stack>
#include<map>
#include<set>
#define INF 0X3F3F3F3F
#define N 1000005
#define M 200005
#define LL long long
#define ULL unsigned long long
#define FF(i, a, b) for(int i = a; i <= b; ++i)
#define RR(i, a, b) for(int i = a; i >= b; --i)
#define FJ(i, a, b) for(int i = a; i < b; ++i)
#define SC(x) scanf("%d", &x)
#define SCC(x, y) scanf("%d%d", &x, &y)
#define SCCC(x, y, z) scanf("%d%d%d", &x, &y, &z)
#define SS(x) scanf("%s", x)
#define PR(x) printf("%d\n", x)
#define PII pair<int, int>
#define PLL pair<unsigned long long, unsigned long long>
#define MP make_pair
#define PB push_back
#define IN freopen("in.txt", "r", stdin)
#define OUT freopen("out.txt", "w", stdout)
using namespace std;
inline void II(int &n){char ch = getchar(); n = 0; bool t = 0;
for(; ch < '0'; ch = getchar()) if(ch == '-') t = 1;
for(; ch >= '0'; n = n * 10 + ch - '0', ch = getchar()); if(t) n =- n;}
inline void OO(int a){if(a < 0) {putchar('-'); a = -a;}
if(a >= 10) OO(a / 10); putchar(a % 10 + '0');}
int n, a[N], hs[N], f[N];
int pre[N], pos[N];
int main(){
   // IN;
    SC(n);
    FF(i, 1, n){
        SC(a[i]);
        hs[(i-1)*3+1] = a[i];
        hs[(i-1)*3+2] = a[i]-1;
        hs[(i-1)*3+3] = a[i]+1;
    }
    int ans = 0;
    sort(hs + 1, hs + 1 + n*3);
    int m = unique(hs + 1, hs + 1 + n*3) - hs - 1;
    FF(i, 1, n) a[i] = lower_bound(hs + 1, hs + 1 + m, a[i]) - hs;
    int anspos;
    FF(i, 1, n){
        if(f[a[i] - 1] + 1 > f[a[i]]){
            f[a[i]] = f[a[i] - 1] + 1;
            pre[i] = pos[a[i] - 1];
            pos[a[i]] = i;
        }
        if(f[a[i]] > ans){
            ans = f[a[i]];
            anspos = i;
        }
    }
    cout << ans << endl;
    vector<int> vec;
    int p = anspos;
    while(p){
        vec.PB(p);
        p = pre[p];
    }
    reverse(vec.begin(), vec.end());
    int sz = vec.size() - 1;
    for(int i = 0; i <= sz; ++i) printf("%d%c", vec[i], i == sz ? '\n' : ' ');
    return 0;
}