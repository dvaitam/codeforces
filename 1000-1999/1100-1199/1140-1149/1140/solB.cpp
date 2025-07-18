#include <iostream>
#include <bits/stdc++.h>
using namespace std;

int main()
{
    long long t;
    cin>>t;
    while(t--)
    {
        long long  n;
        string s;
        cin >> n >> s;
        long long res1 = 0;

        for(int i=0;i<n;i++)
        {

            if(s[i] != '<'){
                    break;
            }
            else
                {
                    res1++;
                }
        }
        long long res2 = 0;
        for(int i=n-1;i>=0;i--)
        {

            if(s[i] != '>'){
                    break;
            }
            else{
                res2++;
            }
        }
        if(res1 < res2)cout << res1 << endl;
        else cout << res2 << endl;
    }

    return 0;
}