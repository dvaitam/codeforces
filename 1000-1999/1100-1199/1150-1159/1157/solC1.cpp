#include<map>
#include<cmath>
#include<stack>
#include<queue>
#include<cstdio>
#include<vector>
#include<cstring>
#include<iostream>
#include<algorithm>
using namespace std;
typedef long long ll;
template<class Read>void in(Read &x){
    x=0;
    int f=0;
    char ch=getchar();
    while(ch<'0'||ch>'9'){
        f|=(ch=='-');
        ch=getchar();
    }
    while(ch>='0'&&ch<='9'){
        x=(x<<1)+(x<<3)+(ch^48);
        ch=getchar();
    }
    x=f?-x:x;
    return;
}
int n,a[200005],l,r,now;
string s;
int main(){
	in(n);
	l=1;
	r=n;
	for(int i=1;i<=n;i++)in(a[i]);
	while(l<=r){
		if(a[l]<=now&&a[r]<=now)break;
		if(a[l]<a[r]){
			if(now<a[l]){
				s+='L';
				now=a[l];
				l++;
			}
			else{
				s+='R';
				now=a[r];
				r--;
			}
		}
		else{
			if(now<a[r]){
				s+='R';
				now=a[r];
				r--;
			}
			else{
				s+='L';
				now=a[l];
				l++;
			}
		}
	}
	cout<<s.size()<<endl;
	cout<<s<<endl;
	return 0;
}