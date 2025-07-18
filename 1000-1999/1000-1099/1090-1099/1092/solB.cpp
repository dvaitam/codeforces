#include <bits/stdc++.h>
using namespace std;

int main ()
{
    int a[105], n, ans = 0;
    cin  >> n;
    for (int i = 0; i < n; ++i)
    {
        cin >> a[i];
    }
    sort(a, a+n);
    for (int i = 0; i < n; ++i)
    {
        ans += a[i + 1] - a[i];
        ++i;
    }
    cout << ans << endl;
}