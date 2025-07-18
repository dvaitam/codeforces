#include <bits/stdc++.h>
using namespace std;

int n;
int a[1005];
int d[1005];
int fact[1005], inv[1005];
const int MOD = 998244353;
inline int lgput(int x, int p){
    int aux = x, ans = 1;
    for(int i = 1; i <= p ; i = i * 2){
        if(i & p) ans = (1LL * ans * aux) % MOD;
        aux = (1LL * aux * aux) % MOD;
    }
    return ans;
}
inline int comb(int n, int k){
    int ans = fact[n];
    ans = (1LL * ans * inv[k]) % MOD;
    ans = (1LL * ans * inv[n - k]) % MOD;
    return ans;
}
int main(){
    scanf("%d", &n);
    for(int i = 1; i <= n ; ++i)
        scanf("%d", &a[i]);

    fact[0] = 1; inv[0] = 1;
    for(int i = 1; i <= n ; ++i){
        fact[i] = (1LL * fact[i - 1] * i) % MOD;
        inv[i] = lgput(fact[i], MOD - 2);
    }

    int Sol = 0;
    for(int i = n; i >= 1 ; --i){
        if(a[i] <= 0) continue ;
        int ram = n - i;
        if(a[i] > ram) continue ;
        d[i] = comb(ram, a[i]);
        for(int j = i + a[i] + 1, nr = a[i]; j <= n ; ++j, ++nr){
            int aux = comb(nr, a[i]);
            aux = (1LL * aux * d[j]) % MOD;
            d[i] += aux;
            if(d[i] >= MOD) d[i] -= MOD;
        }
        Sol += d[i];
        if(Sol >= MOD) Sol -= MOD;
    }
    printf("%d", Sol);
    return 0;
}