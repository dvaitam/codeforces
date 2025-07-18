//私は猫です
#include <bits/stdc++.h>
using namespace std;
#define int long long
int read(){
    int xx = 0, f = 1; char ch = getchar();
    for (;!isdigit(ch); ch = getchar())
        f = (ch == '-' ? -1 : 1);
    for (; isdigit(ch); ch = getchar())
        xx = (xx << 3) + (xx << 1) + ch - '0';
    return xx * f;
}
const int N = 2555;
void solve(int n, int *a, int *op){
    op[0] = 0;
    for (int i = 1; i <= n; ++i)
        for (int j = i; j <= n; ++j)if (a[j] == i){
            if (j != i)
                op[++op[0]] = i, op[++op[0]] = j - i, op[++op[0]] = n - j + 1;
            swap(a[i], a[j]);
            break;
        }
}
int n, m, a[N], b[N], opa[10100], opb[10100];
signed main(){
    n = read(), m = read();
    for (int i = 1; i <= n; ++i)a[i] = read();
    for (int j = 1; j <= m; ++j)b[j] = read();
    solve(n, a, opa), solve(m, b, opb);
    // cout<<opa[0]<<" "<<opb[0]<<endl;
    if ((opa[0] + opb[0]) & 1){
        if ((n & 1) || (m & 1)){
            if (n & 1)for (int i = 1; i <= n; ++i)opa[++opa[0]] = n;
            else for (int i = 1; i <= m; ++i)opb[++opb[0]] = m;
        }
        else {
            printf("-1\n"); return 0;
        }
    }
    while(opa[0] < opb[0])opa[++opa[0]] = 1, opa[++opa[0]] = n;
    while(opb[0] < opa[0])opb[++opb[0]] = 1, opb[++opb[0]] = m;
    printf("%lld\n", opa[0]);
    for (int i = 1; i <= opa[0]; ++i)printf("%lld %lld\n", opa[i], opb[i]);
    return 0;
}