//  ©pandeypranavdeval

#include<bits/stdc++.h>

using namespace std;

typedef long long LL;

#define fr(i,n) for (int i=0; i<n; i++)

#define fr1(i,n) for (int i=1; i<n; i++)

#define pb push_back

#define ppb pop_back

#define V vector<int>

#define mp make_pair

#define endl '\n';

int main()

{

ios_base::sync_with_stdio(false);

cin.tie(NULL);

    int t;

    cin>>t;

    while(t--){

        int n;

        cin>>n;

        vector<LL>v;

        fr(i,n){

            LL x;

            cin>>x;

            v.pb(x);

        }

        LL gc1=v[0],gc2=v[1];

        fr(i,n){

            if(i&1){

                gc2=__gcd(gc2,v[i]);

            }

            else{

                gc1=__gcd(gc1,v[i]);

            }

        }

        if(gc1==gc2){

            cout<<"0\n";

        }

        else{

            bool flag=false;

            if(gc1>1){

                flag=true;

                for(int i=1;i<n;i+=2){

                    if(v[i]%gc1==0){

                        flag=false;

                        break;

                    }

                }

                if(flag){

                    cout<<gc1<<"\n";

                }

            }

            if(flag==false){

                flag=true;

                for(int i=0;i<n;i+=2){

                    if(v[i]%gc2==0){

                        flag=false;

                        break;

                    }

                }

                if(flag){

                    cout<<gc2<<"\n";

                }

                else{

                    cout<<"0\n";

                }

            }

        }

    }

return 0;

}