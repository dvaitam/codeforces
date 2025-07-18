#include <bits/stdc++.h>
#define mn(a, b) (a < b ? a : b)
#define mx(a, b) (a > b ? a : b)
#define f first
#define s second
#define all(v) (v).begin(), (v).end()
#define base 331

using namespace std;

typedef long long ll;
typedef unsigned long long ull;

int const MAXn = 5e3 + 2;

int best[MAXn];
int dp[MAXn][MAXn];
int arr[MAXn];

int main(){

    //freopen(".in","r",stdin);
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int n, k;
    cin>>n>>k;
    for(int i = 1; i <= n; i++){
        cin>>arr[i];
    }
    sort(arr + 1, arr + n + 1);
    for(int i = 1; i <= n; i++){
        for(int j = i; j >= 1; j--){
            if(arr[i] - arr[j] > 5) break;
            best[i] = j;
        }
    }
    for(int i = 1; i <= n; i++){
        int l = best[i];
        dp[i][1] = mx(dp[i - 1][1], i - l + 1);
        for(int j = 1; j < k; j++){
            if(dp[l - 1][j] != 0)
                dp[i][j + 1] = mx(dp[i - 1][j + 1], dp[l - 1][j] + i - l + 1);
        }
    }
    int sol = 0;
    for(int i = 1; i <= k; i++){
        if(dp[n][i] >= k)
            sol = mx(sol, dp[n][i]);
    }
    cout<<sol;
    return 0;
}