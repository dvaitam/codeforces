#include <bits/stdc++.h>
#define TIME ios_base::sync_with_stdio(0)
using namespace std;
    int n;
    char z = 'a';
    int main()
    {
            TIME;
        cin >> n;
        cout << 'a';
        n=n-1;
        for( int i = 1; i <= n/2; i ++ )
        {
            if( i % 2 == 1 )
            {
                z = 'b';
                cout << "bb";
            }
            else
            {
                z = 'a';
                cout << "aa";
            }
        }
        if( n%2==1)
        {
            if(z=='a' )
                cout << 'b';
            else
                cout <<'a';
        }
    }