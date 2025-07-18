#include <bits/stdc++.h> // DON'T WATCH THE CLOCK; DO WHAT IT DOES.  KEEP GOING.
using namespace std;
#define gc getchar_unlocked
#define endl '\n'
#define no cout << "NO" << endl
#define yes cout << "YES" << endl
#define deb(x) cout << #x << "  =  " << x << endl
#define deb2(x, y) cout << #x << "  =  " << x << "  ,  " << #y << "  =  " << y << endl
#define deb3(x, y, z) cout << #x << "  =  " << x << "  ,  " << #y << "  =  " << y << "  ,  " << #z << "  =  " << z << endl
#define deb4(x, y, z, w) cout << #x << "  =  " << x << "  ,  " << #y << "  =  " << y << "  ,  " << #z << "  =  " << z << "  ,  " << #w << "  =  " << w << endl
#define F first
#define S second
#define all(x) x.begin(), x.end()
#define allr(a) a.rbegin(), a.rend()
#define coutd cout << fixed << setprecision(15) // coutd<<res<<endl;
#define pqb priority_queue<int>
#define pqs priority_queue<int, vector<int>, greater<int>>
#define In ios_base::sync_with_stdio(false);
#define The cin.tie(NULL);
#define Name_of_ALLAH cout.tie(NULL);
#define ll long long
#define ull unsigned long long
#define pl pair<ll, ll>
//======================= Global_Variable
const ll N = 400050;
const ll INF = 1e18 + 10;
const ull mod = 1e9 + 7;
//======================= user define function

void Sayem()
{
    ll n = 0, m = 0, a = 0, b = 0, c = 0, i = 0, j = 0, k = 0, l = 0, r = 0, ans = 0, temp = 0, cnt = 0, sum = 0;
    cin >> n >> m;
    vector<ll> aa(n + 5, -1000), bb, cc, dd;
    // map<ll, ll> mp;
    // set<ll> st;
    // pqb pq;
    for (int i = 1; i <= n; i++)
    {
        if ((i * (i + 1)) / 2 <= m)
        {
            temp = i;
        }
    }
    // deb(temp);
    cnt = m;
    cnt -= (temp * (temp + 1)) / 2;
    // deb(cnt);
    for (int i = 0; i < temp; i++)
    {
        aa[i] = i + 2;
    }
    if (cnt == 0)
    {
        for (int i = 0; i < n; i++)
            cout << aa[i] << " ";
        cout << endl;
        return;
    }
    ll rem = 0;
    for (int i = temp - 1; i >= cnt - 1; i--)
        rem += aa[i];
    aa[temp] = (ll)-1 * (rem - 1);
    // bb.push_back(rem);
    // deb3(rem, temp, cnt);

    for (int i = 0; i < n; i++)
        cout << aa[i] << " ";
    cout << endl;

    return;
}

int32_t main()
{
    In The Name_of_ALLAH;
    int _TB = 1;
    cin >> _TB;
    while (_TB--)
    {
        Sayem();
    }
    return 0;
}