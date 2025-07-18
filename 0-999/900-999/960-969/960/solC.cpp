#include <bits/stdc++.h>
/*
  _   _     _____    __  __    ____      ______    ______      ___     ______
 /'\_/`\  /\  __`\  /\ \/\ \  /\  _`\   /\__  _\  /\  _  \    /'___\  /\  _  \
/\      \ \ \ \/\ \ \ \ \ \ \ \ \,\L\_  \/_/\ \/  \ \ \L\ \  /\ \__/  \ \ \L\ \
\ \ \__\ \ \ \ \ \ \ \ \ \ \ \ \/_\__ \    \ \ \   \ \  __ \ \ \ ,__\  \ \  __ \
 \ \ \_/\ \ \ \ \_\ \ \ \ \_/ \  /\ \L\ \   \ \ \   \ \ \/\ \ \ \ \_/   \ \ \/\ \
  \ \_\\ \_\ \ \_____\ \ `\___/  \ `\____\   \ \_\   \ \_\ \_\ \ \_\     \ \_\ \_\
   \/_/ \/_/  \/_____/  `\/__/    \/_____/    \/_/    \/_/\/_/  \/_/      \/_/\/_/
  __ __
 _\ \\ \__
/\__  _  _\
\/_L\ \\ \L_
  /\_   _  _\
  \/_/\_\\_\/
     \/_//_/
*/
#define io ios_base::sync_with_stdio(0), cin.tie(0), cout.tie(0);
#define ll long long
#define ld long double
#define EPS (1e-15)
#define mod 1000000007
#define pii pair<ll, ll>
#define piii pair<pii, ll>
#define SIZE(x) (int) x.size()
#define to__string(x) static_cast<std::ostringstream&>((std::ostringstream() << std::dec << x)).str()
#define lp(i, a, n) for (int(i) = (a); (i) <= (int)(n); ++(i))
#define lpd(i, n, a) for (int(i) = (n); (i) >= (a); --(i))
#define bitcount(n) __builtin_popcountll(n)
#define readi(x) scanf("%d", &x);
#define readl(x) scanf("%lld", &x);
#define reads(x) scanf("%s", &x);
#define readc(x) scanf("%c", &x);
#define readd(x) scanf("%lf", &x);

using namespace std;
const int N = 1000006;
const int K = (int)log2(N) + 2;
const int SQRTN = (int) sqrt(N) + 7;
const ll infi = 1ll << 60;
ll x, d, n, last;
vector<ll> ans;

int main() {
//	freopen("input.txt", "r", stdin);
//	freopen("output.txt", "w", stdout);
    io;
    cin >> x >> d;
    last = 1;
    if (x <= 1e4) {
        cout << x << endl;
        lp (i, 1, x)
            cout << last << " ", last += (d + 2);
        cout << endl;
        return 0;
    }
    while(x) {
        ll len = 1;
        while ((1 << len) - 1 <= x)
            ++len;
        --len;
        lp (i, 1, len)
            ans.push_back(last);
        last += (d + 2);
        x -= (1 << len) - 1;
    }
    cout << SIZE(ans) << endl;
    for (auto& itr : ans)
        cout << itr << " ";
    cout << endl;
    return 0;
}