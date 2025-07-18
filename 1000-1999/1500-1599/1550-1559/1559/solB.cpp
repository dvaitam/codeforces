#include <bits/stdc++.h>

using namespace std;

typedef long long ll;

typedef unsigned long long ull;

typedef long double ld;

typedef vector<int> vi;

typedef vector<ll> vll;

typedef vector<ld> vld;

#define forn(i, n) for(int i = 0; i < int(n); ++i)





int main() {

    ios::sync_with_stdio(0);

    cin.tie(0), cout.tie(0);

    ll  h, m, n, k, k0, t = 1, l, r, w, x, y, z, t1, t2, t3;

    cin >> t;

    string s;

    k=0;

    l=0;

    h=0;

    while (t--) {

        cin>>n;

        cin>>s;

        l=-1;

        forn(i, n)

        if(s[i]!='?')

            l=i;

        if(l==-1){

            s[0]='B';

            l=0;

        }

        r=l-1;

        while(r>=0){

            if(s[r]=='?'){

                if(s[r+1]=='B')

                    s[r]='R';

                else s[r]='B';

            }

            --r;

        }

        r=l+1;

        while(r<n){

            if(s[r]=='?'){

                if(s[r-1]=='B')

                    s[r]='R';

                else s[r]='B';

            }

            ++r;

        }

        cout<<s<<'\n';

    }



}