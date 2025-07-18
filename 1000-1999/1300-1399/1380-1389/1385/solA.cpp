/*

        author :- sairaj

*/

#pragma GCC optimize("O1")

#include<bits/stdc++.h>

using namespace std;

#define all(a) a.begin(),a.end()

#define sortv(a) sort(all(a))

#define sortvg(a) sort(all(a),greater<>());

#define int long long

#define endl "\n"

#define SPEED {ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL);}

template<class T>

void display(T a)

{

    for(auto& it:a)

    {

        cout<<it<<" ";

        cout<<endl;

    }

}

void jawab()

{

    int a,b,c;

    cin>>a>>b>>c;

    if(a==b&&b==c)

        cout<<"YES"<<endl<<a<<" "<<b<<" "<<c;

    else if(a==b || b==c || c==a)

    {

        int maxi = max(max(a,b),c);

        if(a==b)

        {

            if(a==maxi)

            {

                cout<<"YES"<<endl;

                cout<<maxi<<" "<<1<<" "<<c;

            }

            else

            {

                cout<<"NO";

            }

        }

        if(b==c)

        {

            if(b==maxi)

            {

                cout<<"YES"<<endl;

                cout<<a<<" "<<maxi<<" "<<1;

            }

            else

            {

                cout<<"NO";

            }

        }

        if(a==c)

        {

            if(a==maxi)

            {

                cout<<"YES"<<endl;

                cout<<1<<" "<<b<<" "<<maxi;

            }

            else

            {

                cout<<"NO";

            }

        }

    }

    else

        cout<<"NO";

    cout<<endl;

}

int32_t main()

{

    SPEED;

    int t=1;

    cin>>t;

    while(t--)

    {

        jawab();

    }

    return 0;

}