#include<bits/stdc++.h>
using namespace std;
namespace whatever{
	long long readull(){
		char ch=getchar();
		while(!isdigit(ch))
			ch=getchar();
		long long value=ch-'0';
		ch=getchar();
		while(isdigit(ch)){
			value=value*10+ch-'0';
			ch=getchar();
		}
		return value;
	}
	void writeull(long long n){
		if(n<10)
			putchar(n+'0');
		else{
			writeull(n/10);
			putchar(n%10+'0');
		}
	}
	void run(){
		int n=readull();
		long long t=readull();
		static long long a[200001];
		for(int i=1; i<=n; ++i)
			a[i]=readull();
		static int x[200001];
		for(int i=1; i<=n; ++i)
			x[i]=readull();
		static int min_possible_cnt[200002];
		static long long min_possible[200001];
		static long long max_possible[200001];
		for(int i=1; i<=n; ++i){
			min_possible[i]=a[i]+t;
			max_possible[i]=3000000000000000000ll;
		}
		bool ok=true;
		for(int i=1; i<=n; ++i){
			if(x[i]<i){
				ok=false;
				break;
			}
			++min_possible_cnt[i];
			--min_possible_cnt[x[i]];
			if(x[i]!=n)
				max_possible[x[i]]=min(max_possible[x[i]], a[x[i]+1]+t-1);
		}
		if(!ok){
			puts("No");
			return;
		}
		int sum=0;
		for(int i=1; i<=n; ++i){
			sum+=min_possible_cnt[i];
			if(sum)
				min_possible[i]=max(min_possible[i], a[i+1]+t);
		}
		static long long b[200001];
		long long prev=0;
		for(int i=1; i<=n; ++i){
			long long cur=max(prev+1, min_possible[i]);
			if(cur>max_possible[i]){
				ok=false;
				break;
			}
			b[i]=cur;
			prev=cur;
		}
		if(!ok){
			puts("No");
			return;
		}
		puts("Yes");
		for(int i=1; i<=n; ++i){
			writeull(b[i]);
			putchar(' ');
		}
		putchar('\n');
	}
}
int main(){
	whatever::run();
}