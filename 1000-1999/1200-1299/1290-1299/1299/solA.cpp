#include<bits/stdc++.h>

using namespace std;

typedef long long ll;



#define v(li)  vector<ll>

#define vp(li)  vector<pair<ll, ll>>

#define m(li)  map<ll, ll>

#define s(li)  set<ll>

#define yes    cout<<"YES"<<endl

#define no     cout<<"NO"<<endl

#define ss     second

#define ff     first

#define f(i, n)   for (ll i = 0; i < n; i++)

#define nl        cout << "\n"

#define vecsort(v)  sort(v.begin(), v.end())

#define pb          push_back

#define mk          make_pair



void fastio(){

    std::ios::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

}



vector<ll> freq(ll a[],ll n){

    ll c=1;

    v(li) v;

    for(ll i=0;i<n-1;i++){

        if(a[i]==a[i+1])

        c++;

        else{

            v.push_back(c);

            c=1;

        }

    }

    v.push_back(c);

    return v;

}



vector<pair<ll,ll>> freqp(ll a[],ll n){

    ll c=1;

    vp(li) v;

    for(ll i=0;i<n-1;i++){

        if(a[i]==a[i+1])

        c++;

        else{

            v.push_back(make_pair(c,a[i]));

            c=1;

        }

    }

    v.push_back(make_pair(c,a[n-1]));

    return v;

}



ll power(ll a,ll b){

    if(b==0)

    return 1;

    ll res=power(a,b/2);

    if(b%2==0)

    return res*res;

    else

    return res*res*a;

}



ll powerm(ll a,ll b,ll M){

    if(b==0)

    return 1;

    ll res=powerm(a,b/2,M);

    if(b%2==0)

    return ((res%M)*(res%M))%M;

    else

    return (((res%M)*(res%M)%M)*(a%M))%M;

    

}







int main(){

    fastio();

    

        ll n;

        cin>>n;

        ll ma=0;

        ll a[n];

        f(i,n){

            cin>>a[i];

            if(a[i]>ma){

                ma=a[i];

            }

        }

        //cout<<ma<<endl;

        ll ind=0;

        for(ll i=0;i<64;i++){

            if(power(2,i)>ma){

                ind=i;

                break;

            }

        }

        //cout<<ind<<endl;

        ll p[n],s[n];

        ll xo=power(2,ind)-1;

        ll fi;

        ll st;

        p[0]=xo^a[0];

        s[n-1]=xo^a[n-1];

        for(ll i=1;i<n;i++)

        p[i]=p[i-1]&(xo^a[i]);

        for(ll i=n-2;i>=0;i--)

        s[i]=s[i+1]&(xo^a[i]);

        f(i,n){

            if(i==0){

                fi=s[i+1]&a[i];

                st=a[i];

            }

            else if(i==n-1){

                ll x=a[i]&p[i-1];

                if(x>fi){

                    fi=x;

                    st=a[i];

                }

            }

            else{

                ll x=a[i]&p[i-1]&s[i+1];

                if(x>fi){

                    fi=x;

                    st=a[i];

                }

            }

        }

        int flag=1;

        cout<<st<<" ";

        f(i,n){

            if(a[i]==st){

                if(flag)

                flag=0;

                else

                cout<<a[i]<<" ";

            }

            else

            cout<<a[i]<<" ";

        }

        nl;

    }