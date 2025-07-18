#include<bits/stdc++.h>

using namespace std;



#define lli long long int

#define ll long long

#define no cout<<"NO \n";

#define yes cout<<"YES \n";

#define test lli trtyuio ; cin>>trtyuio; while(trtyuio--)

int main(){

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    test

    {

        lli n,l,r;

        cin>>n>>l>>r;

        lli a[n+1],boo1=1;

        for(int i = 1 ; i <= n ; i++)

        {

            a[i] = ((l - 1) / i + 1) * i;

            if(a[i]>r)

            {boo1=0;

            break;}

        }



        if(boo1)

        {

            yes

            for (int i = 1; i <=n; i++)

            {

                cout<<a[i]<<" ";

            }

            cout<<"\n";

        }

        else

        no

    }



    return 0;

}