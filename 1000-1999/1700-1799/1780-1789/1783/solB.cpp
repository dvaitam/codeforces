#include<bits/stdc++.h>

using namespace std;

 

#define ll long long

#define in(a,n) for(int i=0;i<n;i++)cin>>a[i];

#define fr(i,x,y) for(int i= x; i<y;i++)

#define bfr(i,x,y) for(int i= x; i>y;i--)

#define sortvec(vec) sort(vec.begin(),vec.end())

#define gsortvec(vec) sort(vec.rbegin(),vec.rend())

#define gsortarr(arr) sort(arr,arr+sizeof(arr)/sizeof(arr[0]),greater<int>())

#define sortarr(arr) sort(arr,arr+sizeof(arr)/sizeof(arr[0]))

#define v(y) vector<ll>y

#define vp(y) vector<pair<ll,ll>>y

#define pb(x) push_back(x)

#define pbm(x,y) push_back(make_pair(x,y))

#define mp(x) map<int ,int>x 

 

//int dp[100002];

//ll mod = pow(10,9)+7



//int solve(){



//}



//LOOK FOR MODULO 



int main(){

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);cout.tie(NULL);

    int t;

    cin>>t;

    while(t--){

        int n;

        cin>>n;

        int a[n][n];

        int s = 1,l = n*n;



        fr(i,0,n){

            if(i%2== 0){

                fr(j,0,n){            

                    if(j%2 == 0)a[i][j] = s++;

                    else a[i][j] = l--;

                }

            }

            if(i%2== 1){

                bfr(j,n-1,-1){

                    if(j%2 == 0)a[i][j] = l--;

                    else a[i][j] = s++;

                }

            }

        }

        fr(i,0,n){

            fr(j,0,n){

                cout<<a[i][j]<<" ";

            }

            cout<<endl;

        }

        cout<<endl;

        



    }

    return 0;





//Tips : for most of the cases

// 1) Greedy 

// 2) After try Binary Search 

// 2) After try dp 

//Tips : Binary Representaion/ Logic use 

// 1) when ever constraints are given in power of 2 then think about Traversinh through all 32 bits 

// it will reduce TLE and Logic is usuall =y implemented using >> or << (shits) 

// 2) operations of GCC likw bit_ctz are there to compare are count , etc bits



}