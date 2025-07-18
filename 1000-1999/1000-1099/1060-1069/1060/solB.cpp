#include <bits/stdc++.h>
#define LL long long
#define MEM(x,y) memset(x,y,sizeof(x))
#define MOD(x) ((x)%1000000007)

using namespace std;
LL n;
int main()
{
    ios_base::sync_with_stdio(false);
    cin.tie(NULL); cout.tie(NULL);
    cin >> n;
    LL tmp = n, a = 0;
    int ans = 0;
    while(tmp >= 10)
    {
        tmp /= 10;
        a = 10 * a + 9;
        ans += 9;
    }
    LL b = n - a;
    //cout << a << " " << b << endl;
    while(b)
    {
        ans += (b % 10);
        b /= 10;
    }
    cout << ans;
    return 0;
}