#include <bits/stdc++.h>
#define endl '\n'
using namespace std;
typedef unsigned long long ull;
typedef long long ll;
typedef pair<int, int> ii;
int n, ans;

int main()
{
    ios_base::sync_with_stdio(false);
    cin.tie(0);
    int q;
    cin >> q;
    while (q--)
    {
        cin >> n; string s; cin >> s;
        if (n == 2) 
        {
            if (s[1] > s[0]) 
            {
                cout << "YES" << endl;
                cout << 2 << endl << s[0] << " " << s[1] << endl;
            }
            else cout << "NO" << endl;
        }
        else 
        {
            cout << "YES" << endl;
            cout << 2 << endl;
            cout << s[0] << " ";
            for (int i = 1; i < n; i++) cout << s[i];
            cout << endl;
        }
    }
    
    return 0;
}