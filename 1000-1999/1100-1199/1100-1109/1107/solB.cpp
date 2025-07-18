#include<bits/stdc++.h>
#define ll long long
#define ld long double
#define ull unsigned long long
#define m_p make_pair
#define ss second
#define ff first
#define pb push_back
using namespace std;

const int N = 5200;
const ll MOD = 1e9+7;


int main()
{
    ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);
    int n;
    cin >> n;
    for (int i=0; i<n; ++i)
    {
        ll x, k;
        cin >> k >>x;
        cout << x + (9*(k-1)) << '\n';
    }


}