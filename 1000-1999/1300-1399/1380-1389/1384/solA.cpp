#include<bits/stdc++.h>

using namespace std;

#define ll long long



int main()

{

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);





    int t;

    cin>>t;

    while(t--)

    {

        int n;

        cin>>n;

        int ar[n];

        string sam;

        for(int i=0;i<=100;i++) sam.push_back('a');

        for(int i=0; i<n; i++)

        {

            cin>>ar[i];

        }



        for(int i=0; i<=n; i++)

        {

           if(i==0) cout<<sam<<endl;

           else{

            if(sam[ar[i-1]]=='a') sam[ar[i-1]]++;

            else sam[ar[i-1]]--;

            cout<<sam<<endl;

           }

        }



    }

    return 0;

}