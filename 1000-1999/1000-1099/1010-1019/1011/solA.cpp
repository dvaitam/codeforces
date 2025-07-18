///...................SUBHASHIS MOLLICK...................///
///.....DEPARTMENT OF COMPUTER SCIENCE AND ENGINEERING....///
///.............ISLAMIC UNIVERSITY,BANGLADESH.............///
///....................SESSION-(14-15)....................///
#include<bits/stdc++.h>
using namespace std;
#define sf(a) scanf("%lld",&a)
#define sf2(a,b) scanf("%lld %lld",&a,&b)
#define sf3(a,b,c) scanf("%lld %lld %lld",&a,&b,&c)
#define pf(a) printf("%lld",a)
#define pf2(a,b) printf("%lld %lld",a,b)
#define pf3(a,b,c) printf("%lld %lld %lld",a,b,c)
#define nl printf("\n")
#define   timesave              ios_base::sync_with_stdio(false); cin.tie(0); cout.tie(0);
#define ll long long
#define pb push_back
#define MPI map<int,int>mp;
#define fr(i,n) for(i=0;i<n;i++)
#define fr1(i,n) for(i=1;i<=n;i++)
#define frl(i,a,b) for(i=a;i<=b;i++)
/*primes in range 1 - n
1 - 100(1e2) -> 25 pimes
1 - 1000(1e3) -> 168 primes
1 - 10000(1e4) -> 1229 primes
1 - 100000(1e5) -> 9592 primes
1 - 1000000(1e6) -> 78498 primes
1 - 10000000(1e7) -> 664579 primes
large primes ->
104729 1299709 15485863 179424673 2147483647 32416190071 112272535095293 48112959837082048697
*/
//freopen("Input.txt","r",stdin);
//freopen("Output.txt","w",stdout);
//const int fx[]={+1,-1,+0,+0};
//const int fy[]={+0,+0,+1,-1};
//const int fx[]={+0,+0,+1,-1,-1,+1,-1,+1};   // Kings Move
//const int fy[]={-1,+1,+0,+0,+1,+1,-1,-1};  // Kings Move
//const int fx[]={-2, -2, -1, -1,  1,  1,  2,  2};  // Knights Move
//const int fy[]={-1,  1, -2,  2, -2,  2, -1,  1}; // Knights Move
long n,k;
main()
{
    timesave;
    while(cin>>n>>k)
    {
        string s1;
        cin>>s1;
        if(k>13||k>n)
        {
            cout<<-1<<endl;
        }
        else
        {
            long i,j,x,cnt,ans=0,mn=100000000,f=0;
            sort(s1.begin(),s1.end());
            for(i=0; i<n; i++)
            {
                cnt=1,x=0,ans=s1[i]-96;
                char ch=s1[i];
                for(j=i; j<n; j++)
                {
                    x=s1[j]-ch;
                    //cout<<x<<"yes";
                    if(x>1)
                    {
                        ch=s1[j];
                        ans+=(s1[j]-96);
                        cnt++;
                    }
                    if(cnt==k)
                    {
                        f=1;
                        mn=min(ans,mn);
                        break;
                    }
                }
            }
            if(f==0)
                cout<<-1<<endl;
            else
                cout<<mn<<endl;
        }

    }
}