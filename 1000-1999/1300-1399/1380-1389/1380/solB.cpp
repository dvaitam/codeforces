#include<bits/stdc++.h>

            using namespace std;

                        typedef   long long ll;

                        typedef   vector<ll>vl;

                        typedef   vector<int>vi;

                        typedef   pair<ll,ll>pl;

                        #define   F first

                        #define   S second

                        #define   pb push_back

                        #define   rp(i,a,b) for(ll i=a;i<b;i++)

                        #define   yes cout<<"YES"<<"\n"

                        #define   no cout<<"NO"<<"\n"

                        #define   en cout<<endl

                        #define   print(a,n) for(ll i=0; i<n;i++){cout<<a[i]<<" "; }  

                        #define   read(a,n) for(ll i=0; i<n;i++){cin>>a[i]; } 

                        #define   copy(a,b) for(ll i=0;i<a.size();i++){b[i]=a[i]; }

                        // const ll INF=1e18+1;   

                        // ll dx[] = { 1 , 1, -1 , -1 };

                        // ll dy[] = {1 , -1 , -1 , 1 };

                        // ll ddx[]={1,-1,0,0,1,1,-1,-1};

                        // ll ddy[]={0,0,1,-1,1,-1,-1,1};

                        // ll kx[] = { -2, -1, 1, 2, -2, -1, 1, 2 };

                        // ll ky[] = { -1, -2, -2, -1, 1, 2, 2, 1 }; 

            // int gcd(long long a,long long b)

            // {

            //     for (;;)

            //     {

            //         if (a == 0) return b;

            //         b %= a;

            //         if (b == 0) return a;

            //         a %= b;

            //     }

            // }

            // long long  lcm(long long a, long b)

            // {

            //     long long temp = gcd(a, b);

            

            //     return temp ? (a / temp * b) : 0;

            // }

            /*void reverse(long long a[],long long l,long long r){

                int b[r];

                for(long long i=l,j=r;i<=r;i++,j--)

                b[i]=a[j];

                for(long long i=l;i<=r;i++){

                    a[i]=b[i];

                }

            };*/ 

            // bool isPrime(long long int n)

            // {

            //     // Corner cases

            //     if (n <= 1)

            //         return false;

            //     if (n <= 3)

            //         return true;

            //     // This is checqed so that we can sqip

            //     // middle five numbers in below loop

            //     if (n % 2 == 0 || n % 3 == 0)

            //         return false;

            

            //     for (long long int  i = 5; i * i <= n; i = i + 6)

            //         if (n % i == 0 || n % (i + 2) == 0)

            //             return false;

            

            //     return true;

            // }

            //    ll power(ll a, ll b)

            // {

            //     ll res = 1;

            //     while (b)

            //     {

            //         if (b & 1)

            //             res = (res * a);

            //         a = (a * a);

            //         b >>= 1;

            //     }

            //     return res;

            // }

            // bool areBracketsBalanced(string expr)

            // { 

            //     stack<char> s;

            //     char x;

            

            //     for (int i = 0; i < expr.length(); i++)

            //     {

            //         if (expr[i] == '(')

            //         {

            //             s.push(expr[i]);

            //             continue;

            //         } 

            //         if (s.empty())

            //             return false;

            

            //         switch (expr[i]) {

            //         case ')':

            //             x = s.top();

            //             s.pop();

            //         }

            //     }

            //      return (s.empty());

            // }

            // bool isPalindrome(string str)

            // {

            //     int n = str.length();

            //     if (n == 0)

            //         return true; 

            //     return isPalRec(str, 0, n - 1);

            // }

            

        // ll Comb(ll n, ll k) {

        //     ll res = 1;

        //     for (ll i = n - k + 1; i <= n; ++i)

        //         res *= i;

        //     for (ll i = 2; i <= k; ++i)

        //         res /= i;

        //     return res;

        // }

        // bool cmp(vector<int> a,vector<int>b)

        // {

        //     return a[0]<b[0];

        // }

    //     long long binpow(long long a, long long b) {

    //     long long res = 1;

    //     while (b > 0) {

    //         if (b & 1)

    //             res = res * a;

    //         a = a * a;

    //         b >>= 1;

    //     }

    //     return res;

    // }

    // bool cmp(pair<int,int>a,pair<int,int>b){

    //     return a.second+a.first<b.second+b.first;

    // }

        int find(int p[],int q[],int prev,int i,int n,int ans){

            if(i==n-1){

                return 1;

            }

            if(min(p[i],q[i])-1!=prev &&min(p[i],q[i])-1!=prev){

            ans+=find(p,q,min(p[i],q[i])-1,i+1,n,ans);

            }



            ans+=1+find(p,q,min(p[i],q[i])-1,i+1,n,ans);

            return ans;

        }

            void solve(){

            ll n,x,y,ans=0;

            cin>>n;

            int p[n];

            int q[n];

            rp(i,0,n){

                cin>>p[i];

            }

             rp(i,0,n){

                cin>>q[i];

            }

            cout<<find(p,q,min(p[0],q[0])-1,0,n,ans)<<endl;

            }

                

            void solver(){

               string s,ans;

               cin>>s;

               map<char,int>mp;

               rp(i,0,s.size()){

                mp[s[i]]++;

               }

               if(mp['S']==mp['R'] &&mp['S']==mp['P']){

                for(int i=0;i<s.size();i++){

                    if(s[i]=='R')

                    ans+='P';

                    else if(s[i]=='P')

                    ans+='S';

                    else ans+='R';

                }

                cout<<ans<<endl;

                return;

               }

               else{

                int mx=0;

                char cmx;

                for(auto x:mp){

                if(x.second>mx){

                    mx=x.second;

                    cmx=x.first;

                }

                }

                for(int i=0;i<s.size();i++){

                  if(cmx=='P'){

                    ans+='S';

                  }

                  else if(cmx=='R')

                  ans+='P';

                  else ans+='R';

                }

                cout<<ans<<endl;

               }





            }

    

        int main(){

                

                // freopen("input.txt", "r", stdin);

                // freopen("output.txt", "w", stdout);

                ios_base::sync_with_stdio(false);

                cin.tie(0);

                int t=1;

                cin>>t;

                while(t--)

                solver();

                return 0;

            }