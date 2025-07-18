#include<bits/stdc++.h>

using namespace std;

 

/*

 

*/

 

void solve(int testCase)

{

    long long n,m,k;

    cin >> n >> m >> k;

    long long stBig = 1;

    while(k--)

    {

        long long temp = m;

        long long big = (n%m);

        long long stSmall = stBig;

        while(temp>0LL)

        {

            if(big>0LL)

            {

                cout << (n+m-1)/m << " ";

                long long eles = (n+m-1)/m;

                while(eles>0LL)

                {

                    cout << stBig << " ";

                    ++stBig;

                    if(stBig>n)

                    {

                        stBig = 1LL;

                    }

                    stSmall = stBig;

                    --eles;

                }

                cout << "\n";

                --big;

            }

            else

            {

                cout << n/m << " ";

                long long eles = n/m;

                while(eles>0LL)

                {

                    cout << stSmall << " ";

                    ++stSmall;

                    if(stSmall>n)

                        stSmall = 1LL;

                    --eles;

                }

                cout << "\n";

            }

            --temp; 

        }

    }

    cout << "\n";

}   

 

/*

 

*/

 

int main()

{

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    int t; cin >> t;

    for(int testCase=1;testCase<=t;++testCase)

    {

        solve(testCase);

    }

}