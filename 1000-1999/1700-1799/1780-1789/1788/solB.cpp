#include <bits/stdc++.h>



using namespace std;



#define int long long



const int N = 1010;



int ans[10] = {0, 1, 2, 10, 18, 100, 180, 1000, 1800, 10000};



signed main() {

    ios::sync_with_stdio(0), cin.tie(0), cout.tie(0);

    int t, n;

    cin >> t;

    while(t --){

        cin >> n;

        if(n % 2 == 0) cout << n / 2 << ' ' << n / 2 << '\n';

        else if(n % 10 != 9) cout << n / 2 << ' ' << n - n /2 << '\n';

        else{

            int p = n / 2, q = n - p;

            int cnt = 0;

            while(n % 10 == 9) n /= 10, cnt ++;

            int val = n % 10;

            if(val == 0) val = 9, cnt --;

            if(val % 2 == 0) cout << p + ans[cnt - 1] * 5 << ' ' << q - ans[cnt - 1] * 5 << '\n';

            else cout << p + ans[cnt] * 5 << ' ' << q - ans[cnt] * 5 << '\n';

        }

    }

    return 0;

}