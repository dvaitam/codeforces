#include <bits/stdc++.h>
using namespace std;

#define fast_io   ios_base::sync_with_stdio(false); cin.tie(nullptr);
#define nl        "\n"

using ll = long long ;

void solve()
{
    int n , a , b ;
    cin >> n >> a >> b ;
    if (abs(a - b) > 1)
    {
        cout << -1 << nl ;
        return ;
    }
    else if (a + b > n - 2)
    {
        cout << -1 << nl ;
        return ;
    }

    vector<int> permut(n) ;
    iota(permut.begin(), permut.end() , 1) ;

    int temp = min(a, b) ;

    if (a > b)
    {
        int i = 1 ;
        for (i = 1 ; temp > 0 ; i += 2)
        {
            swap(permut[i], permut[i + 1]) ;
            temp-- ;
        }
        for (int k = i; k < n - 1 ; k++)
        {
            swap(permut[k], permut[k + 1]) ;
        }
    }
    else if (a < b)
    {
        int i = n - 2 ;
        for (i = n - 2 ; temp > 0 ; i -= 2)
        {
            swap(permut[i], permut[i - 1]) ;
            temp-- ;
        }
        for (int k = i; k > 0 ; k--)
        {
            swap(permut[k], permut[k - 1]) ;
        }
    }
    else if (a == b)
    {
        int i = 1 ;
        for (i = 1 ; temp > 0 ; i += 2)
        {
            swap(permut[i], permut[i + 1]) ;
            temp-- ;
        }
    }
   
    for (int i = 0 ; i < n ; i++)
    {
        cout << permut[i] << " " ;
    }
    cout << nl ;
}


int32_t main()
{
    fast_io
    int tc;
    cin >> tc;
    while (tc--) solve();
}