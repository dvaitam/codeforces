/*
* @Author: XiaoBanni
* @Email:  xjc5069@gmail.com
* @Date:   2019-05-15 23:41:37
* @Last Modified by:   Y
* @Last Modified time: 2019-05-15 23:41:37
*/

#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef pair<int,int> pint;
typedef pair<ll,ll> pll;
typedef priority_queue<int> maxpque;
typedef priority_queue<int,vector<int>,greater<int> > minpque;
const int MAXN = 1e5+5;
const int INF=INT_MAX/2;
const ll LLINF=9223372036854775807/2;
#define in(a) scanf("%d",&a)
#define inll(a) scanf("%lld",&a)
#define out(a) printf("%d",a)
#define outll(a) printf("%lld",a)
#define outln(a) printf("%d\n",a)
#define outllln(a) printf("%lld\n",a)
#define fin(i,st,n) for(int i=st;i<n;i++)
#define fin2(i,st,n) for(int i=st;i<=n;i++)
#define mk make_pair
#define maxpque(type) priority_queue<type>
#define minpque(type) priority_queue<type,vector<type>,greater<type> > //pay attention that no ',' in type
#define gcd __gcd
#define IOS ios::sync_with_stdio(false);cin.tie(0); cout.tie(0); 

char res[2*MAXN];
char s[2*MAXN];

int main(){
	int n;
	in(n);scanf("%s",&s);
	int p=1;
	int q=1;
	for(int i=0;i<n;i++){
		if(s[i]=='('){
			if(p){
				res[i]='0';
			}
			else res[i]='1';
			p^=1;
		}
		else{
			if(q){
				res[i]='0';
			}
			else res[i]='1';
			q^=1;
		}
	}
	printf("%s\n",res);
	return 0;
}