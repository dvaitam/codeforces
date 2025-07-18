/** بِسْمِ اللَّهِ الرَّحْمَنِ الرَّحِيم**/



#include<iostream>

#include<cstdio>

#include<sstream>

#include<cstdlib>

#include<cmath>

#include<algorithm>

#include<set>

#include<map>

#include<string>

#include<cstring>

#include<vector>

#include<stack>

#include<queue>

#include<list>

#include<fstream>

#include<numeric>

#include<iterator>



using namespace std;



#define MXN

#define MXE

#define MXQ

#define MOD

#define INF

#define lli long long int

#define lu unsigned long

#define llu unsigned long long int

#define PI acos(0.0)

#define pb push_back

#define ppb pop_back

#define pfs printf("*")

#define nl printf("\n")

#define pf1 printf("%d",x)

#define pf2 printf("%d %d\n",x,y)

#define pf3 printf("%d %d %d\n",x,y,z)

#define max3(a,b,c) max(a,max(b,c))

#define min3(a,b,c) min(a,min(b,c))

#define sf1(a) scanf("%d",&a)

#define sf2(a,b) scanf("%d %d",&a,&b)

#define sf1ll(a) scanf("%lld",&a)

#define sf2ll(a,b) scanf("%lld %lld",&a,&b)

#define takell(a) scanf("%I64d", &a)

#define rep1(i,n) for(int i=0;i<n;i++)

#define rep2(i,n) for(int i=1;i<=n;i++)

#define rep3(n,i) for(int i=n;i>=0;i--)

#define rep(a,n) for(int i=a;i<=n;i++)



/***************************** MAIN PROGRAM STARTS HARE ***************************/



int main()

{

    char arr[110];

    bool ok=true;

    string str,str1;

    int t,a,b,c=0;

    cin>>str>>str1;

    int l,i;

    l = str.size();



    for(i=0;i<l;i++)

    {

        if(str[i]!='z')

            break;

    }

    if(i==l)

    {

        printf("No such string\n");

        return 0;

    }

    for(i=l-1;i>=0;i--)

    {

        if(str[i]=='z')

        {

            str[i]='a';

        }

        else

        {

             str[i]++;

             break;

        }

    }

    if(str<str1)

        cout<<str<<endl;

    else

      printf("No such string\n");



    return 0;

}