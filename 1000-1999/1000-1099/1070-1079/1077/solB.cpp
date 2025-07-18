#include <bits/stdc++.h>
using namespace std;
#define ll long long int
#define llu unsigned long long int
void pv(vector<ll> arr)
{
    ll size = arr.size();
    for (ll i = 0; i < size; i++)
    {
        printf("%lld ", arr[i]);
    }
}
void gv(vector<ll> &arr)
{
    ll size = arr.size();
    for (ll i = 0; i < size; i++)
    {
        scanf("%lld", arr[i]);
    }
}
int main(void)
{
    ll n;
    cin >> n;
    vector<ll> arr(n);
    for (ll i = 0; i < n; i++)
    {
        cin >> arr[i];
    }
    ll c = 0;
    for (ll i = 0; i < n; i++)
    {
        if (arr[i] == 0)
        {
            if (i == 0 || i == n - 1)
                continue;
            else if (arr[i + 1] == 1 && arr[i - 1] == 1)
            {
                if (i + 1 == n - 1)
                {
                    arr[i + 1] = 0;
                    c++;
                }
                else
                {
                    arr[i + 1] = 0;
                    c++;
                }
            }
        }
    }
    cout << c;
    return 0;
}