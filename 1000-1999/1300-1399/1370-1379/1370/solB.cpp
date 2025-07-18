#include <bits/stdc++.h>

using namespace std;

#define ll long long

#define mod 1000000007

int main(){

    ios_base::sync_with_stdio(false);   cout.tie(nullptr);  cin.tie(nullptr);

    int T=1; cin>>T;

    while(T--) {

       ll n,x; cin>>n;

       vector<ll>even,odd;

       for(int i=0;i<2*n;i++){

           cin>>x;

           if(x%2==0) even.push_back(i+1);

           else odd.push_back(i+1);

       }

      n--;

       if(even.size()) {

           for (int i = 0; i < even.size() - 1; i += 2) {

               if (n == 0)break;

               cout << even[i] << ' ' << even[i + 1] << '\n';

               n--;

           }

       }

       if(odd.size()) {

           for (int i = 0; i < odd.size() - 1; i += 2) {

               if (n == 0)break;

               cout << odd[i] << ' ' << odd[i + 1] << '\n';

               n--;

           }

       }



    }

    return 0;

}