#include <bits/stdc++.h>
using namespace std;

signed main()
{
    string s; cin >> s; s=s+'_';

    for(int i=1; i<s.length(); ++i)
    {
        if(s[i-1] != 'a' && s[i-1] != 'e' && s[i-1] != 'i' && s[i-1] != 'o' && s[i-1] != 'u' && s[i-1] != 'n')
        {
            if(s[i] != 'a' && s[i] != 'e' && s[i] != 'i' && s[i] != 'o' && s[i] != 'u')
            {
                cout << "NO" << endl;
                return 0;
            }
        }
    }
    cout << "YES" << endl;
    return 0;
}