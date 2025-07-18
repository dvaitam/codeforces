#include <bits/stdc++.h>
     
    using namespace std;
     
    typedef long long int ll;
    typedef long long int lli;
    typedef double db;
    typedef vector<ll> vi; 
    typedef vector<vi> vvi; 
    typedef pair<ll,ll> ii; 
     
    #define F first
    #define S second
    #define pb push_back
    #define mp make_pair
    #define sz(a) ll((a).size()) 
    #define all(c) (c).begin(),(c).end() 
    #define tr(c,i) for(auto i = (c).begin(); i != (c).end(); i++) 
    #define present(c,x) ((c).find(x) != (c).end())             // for sets,maps etc O(log(n));
    #define cpresent(c,x) (find(all(c),x) != (c).end())         // for vector and others O(n);
    #define rep(i,a,b) for( long long int i = a; i < b; i++)
    #define _USE_MATH_DEFINES
    #define M 1e9+7
    
    int main() {
        ios_base::sync_with_stdio(0);cin.tie(0);
        ll n,a;cin>>n>>a;
        ll v1=1,v2=2;
        ll v_c=3,v_s=2;
        double center=double(360)/double(n);
        double diff=abs( center/double(2) -double(a) );
        center+=double(360)/double(n);
        double min_diff=diff;
        ll ver=3;
        while(center<=180){
            double diff1=abs( center/double(2) -double(a) );
            if(diff1<min_diff){
                v_c=ver+1;
                v_s=ver;
                min_diff=diff1;
            }
            double diff2=abs( double(180)-center/double(2) -double(a) );
            if(diff2<min_diff){
                v_c=ver-1;
                v_s=ver;
                min_diff=diff2;
            }
            ver++;
            center+=double(360)/double(n);
        }
        cout<<1<<" "<<v_c<<" "<<v_s<<endl;
        return 0;
    }