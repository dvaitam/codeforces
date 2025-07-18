#include <bits/stdc++.h>

using namespace std;

typedef long long ll;
typedef unsigned long long ull;
typedef unsigned int uint;

char s[ 1000009 ];
int n;

vector< pair<int, char> > answ;

int main()
{
    //freopen("input.txt", "r", stdin);
    //freopen(".out", "w+", stdout);
    //ios_base::sync_with_stdio( false );

    scanf("%s", s+2 );
    n = strlen( s+2 );
    s[ 1 ] = '0';
    s[ 0 ] = '0';

    for ( int i=n+1; i>=1; i-- )
    {
        if ( s[ i ] == '1' )
        {
            if ( s[ i-1 ] == '0' )
                answ.push_back( make_pair( n+1-i, '+' ) );
            else
            {
                int k = 0;
                while ( s[ i-k ] == '1' )
                    s[ i-k ] = 0, k++;

                s[ i-k ] = '1';
                answ.push_back( make_pair( n+1-i, '-' ) );
            }
        }
    }

    printf("%d\n", (int)answ.size() );
    for ( auto it : answ )
    {
        printf("%c2^%d\n", it.second, it.first );
    }

    return 0;
}