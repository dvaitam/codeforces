#include <bits/stdc++.h>
#include <fstream>
#define ll long long
using namespace std;
long long ncr (int n, int r) ;
long long factorial(long long x ) ;
ll power (ll x, ll y ) ;
bool palindromic( ll h , ll m ) ;
void sexy_code (){
    ios_base::sync_with_stdio(0) ;
    cin.tie(0); cin.tie(0);
}


int main()
{
    sexy_code() ;
    ll n , t ;
    cin >> n >> t ;
    vector<ll> buses ;
    for ( int i = 0 ; i < n ; i++){
        ll d , s ;
        cin >> d >> s ;
        while(d < t ){
            d+= s ;
        }
        buses.push_back(d) ;
    }
    //for ( int i = 0 ; i < buses.size() ;i ++){
      //  cout << buses[i] << endl ;
    //}
    ll mini = 1e18;
    ll number ;
    for ( int i = 0 ; i < buses.size() ; i++){
        if ( buses[i] <= mini ){
            mini = buses[i] ;
            number = i + 1 ;
        }
    }
    cout << number << endl ;


    /*ll n , arr[1005] = {} ;
    cin >> n ;
    vector< pair <ll , ll > > fuck ;
    for ( int i = 0 ; i < n ; i++){
        cin >> arr[i] ;
    }
    for ( int i = 0 ; i < n-1 ; i += 2){
        fuck.push_back( make_pair(arr[i],arr[i+1])) ;
    }
    for ( int i = 0 ; i < fuck.size() ; i++){
        //cout << fuck[i].first << "\t" << fuck[i].second << endl ;
    }
    bool intersect = false ;
    for ( int i = 0 ; i < fuck.size() ; i++){
            for ( int j = 0 ; j < fuck.size() ; j++){
                if ( i == j ){
                    continue;
                }
                if ( (fuck[i].first < fuck[j].first && fuck[i].second < fuck[j].second && fuck[i].first > fuck[j].second )
                    || ( fuck[i].first > fuck[j].first && fuck[i].second > fuck[j].second && fuck[i].first < fuck[j].second  )){
                    intersect = true;
                    cout << "x" << endl ;
                }
            }
    }
    if ( intersect){
        cout << "yes" ;
    }
    else {
        cout << "no" ;
    }*/


    return 0 ;
}






bool palindromic( ll h , ll m ){
    ll new_num = 0 , mod ;
    while ( h != 0){
        mod = h % 10 ;
        new_num = new_num * 10 + mod ;
        h /= 10 ;
    }
    if ( new_num == m){
        return true ;
    }
    else {
        return false ;
    }

}












ll power (ll x, ll y )
{
    ll h = 1 ;
    for (int i = 0 ; i < y ; i++ )
    {
        h = h * x ;
    }
    return h ;
}


long long factorial(long long x )
{
    long long res = 1 ;
    for (int i = 1 ; i <= x ; i++)
    {
        res = res * i ;
    }


    return res ;
}

long long ncr (int n, int r)
{
    long long ans = 1 ;
    int rfact = 2 ;
    for ( int i = n-r+1 ; i<= n ; i++)
    {
        ans *= i ;
        if (ans % rfact == 0 && rfact <=r)
        {
            ans /= rfact ;
            rfact ++ ;
        }
        if ( ans >= 1000000000 + 9)
        {
            ans = ans % (1000000009) ;
        }
    }
    return ans ;

}