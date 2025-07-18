#include <bits/stdc++.h>

using namespace std;

#define mx 300005
#define int long long
#define pii pair <int, int>
#define piii pair <int, pii>
#define f first
#define s second

int ara[mx];
int32_t main(){
    //freopen("in.txt", "r", stdin);
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    cout.tie(0);
    int n, m, k, p;
    while(cin >> n ){
        int a, b;
        a = b = 0;
        for(int i=1;i<=n;i++){
            int x;
            cin >> x;
            a += x;
        }
        for(int i=1;i<=n;i++){
            int x;
            cin >> x;
            b += x;
        }
        if(a>=b) cout << "Yes" << endl;
        else cout << "No" << endl;
    }
    return 0;
}