/*
 * This is my code,
 * my code is amazing...
 */
//Template v2.0
//iostream is too mainstream
#include<iostream>
#include<string>
#include<algorithm>
#include<queue>
#include<map>
#include<set>
#include<unordered_map>
#include<unordered_set>
#include<vector>
#include<iomanip>
//clibraries
#include<cstring>
#include<cmath>
#include<cstdlib>
#include<cstdio>
#include<ctime>
//defines
#define ll long long
#define lld long double
#define pll pair<ll,ll>
#define pld pair<lld,lld>
#define vll vector<ll>
#define vvll vector<vll>
#define INF 1000000000000000047
const char en='\n';
#define debug(x){cerr<<x<<en;}
#define prime 47
#define lprime 1000000000000000009
#define lldmin LDBL_MIN
#define MP make_pair
#define PB push_back
using namespace std;




int main(){
	ios::sync_with_stdio(false);
int n;
cin>>n;
if(n>3){

    if(n%2==0){
    cout<<"YES"<<endl;
        for(int i=n; i>4; i-=2){
            cout<<i<<" - "<<i-1<<" = 1"<<en;
        }
        for(int i=n; i>4; i-=2)
            cout<<"1 * 1 = 1"<<en;

        cout<<"2 * 3 = 6"<<en;
        cout<<"6 * 4 = 24"<<en;
        cout<<"24 * 1 = 24"<<en;


    }
    else if(n%2==1){
        cout<<"YES"<<en;
        cout<<"5 - 1 = 4"<<en;
        cout<<"4 - 2 = 2"<<en;
        cout<<"2 * 3 = 6"<<en;
        cout<<"6 * 4 = 24"<<en;
        for(int i=n; i>5; i-=2)
            cout<<i<<" - "<<i-1<<" = 1"<<en<<"24 * 1 = 24"<<en;

        
    }


/*
    else if(n>=9){
    cout<<"YES"<<endl;
        cout<<"5 - 4 = 1"<<en;
        cout<<"7 - 6 = 1"<<en;
        cout<<"9 - 8 = 1"<<en;
        cout<<"1 + 1 = 2"<<en;
        cout<<"1 + 2 = 3"<<en;
        cout<<"3 + 1 = 4"<<en;
        cout<<"2 * 3 = 6"<<en;
        cout<<"6 * 4 = 24"<<en;
        for(int i=11; i<=n; i+=2){
            cout<<i<<" - "<<i-1<<" = 1"<<en;
            cout<<"24 * 1 = 24"<<en;

        }



    }
    else cout<<"NO"<<en;

*/

}
else cout<<"NO"<<en;



}