#include <bits/stdc++.h>
#include <cstdio>
#include <cstdlib>
#include <algorithm>
#include <string>
#include <iostream>
#include <cmath>
#include <vector>
#include <queue>
#include <set>
#include <sstream>
#include <map>
#include <cstring>
#include <complex>
#include <numeric>
#include <functional>
#include <stack>
#include <stdarg.h>
//#define NDEBUG
#include <assert.h>
using namespace std;
#ifndef NDEBUG
    #define debug(x) cerr<<#x<<"=\""<<x<<"\""<<" at line#"<<__LINE__<<endl;
    #define hline() cerr<<"-----------------------------------------"<<endl;
	#define Unit_test(fun) class Test_##fun { public: Test_##fun (){ fun ();};}Test_obj_##fun
#else
    #define debug(x)
    #define hline()
	#define Unit_test(fun)
#endif
typedef long long int llint;
typedef unsigned long long int ullint;
#define low(x) ((((x)^((x)-1))&x))
#define two(x)  ((1)<<(x))
#define PB(x) push_back((x))
#define SORT(x) sort(x.begin(),x.end())
// print function
template<typename T1> void print_list(const T1& a){cerr<<a<<endl;}
template<typename T1,typename T2> void print_list(const T1& a,const T2& b){cerr<<a<<" "<<b<<endl;}
template<typename T1,typename T2,typename T3> void print_list(const T1& a,const T2& b,const T3& c){cerr<<a<<" "<<b<<" "<<c<<endl;}
template<typename T1,typename T2,typename T3,typename T4> void print_list(const T1& a,const T2& b,const T3& c,const T4& d){cerr<<a<<" "<<b<<" "<<c<<" "<<d<<endl;}
template<typename T1,typename T2,typename T3,typename T4,typename T5> void print_list(const T1& a,const T2& b,const T3& c,const T4& d,const T5& e){cerr<<a<<" "<<b<<" "<<c<<" "<<d<<" "<<e<<endl;}
template<typename T1,typename T2,typename T3,typename T4,typename T5,typename T6> void print_list(const T1& a,const T2& b,const T3& c,const T4& d,const T5& e,const T6& f){cerr<<a<<" "<<b<<" "<<c<<" "<<d<<" "<<e<<" "<<f<<endl;}
template<typename T1,typename T2,typename T3,typename T4,typename T5,typename T6,typename T7> void print_list(const T1& a,const T2& b,const T3& c,const T4& d,const T5& e,const T6& f,const T7& g){cerr<<a<<" "<<b<<" "<<c<<" "<<d<<" "<<e<<" "<<f<<" "<<g<<endl;}

//************************************************
template<typename A,typename B>
ostream& operator<<(ostream& out,const pair<A,B>& pp)
{
	out<<"("<<pp.first<<","<<pp.second<<")";
	return out;
}
template<typename A,typename B>
istream& operator<<(istream& in,pair<A,B>& pp)
{
	cerr<<"A pair wanted"<<endl;
	in>>pp.first>>pp.second;
	return in;
}

template<typename T>
ostream& operator<<(ostream& out,const vector<T>& vect)
{
	out<<"length = "<<vect.size()<<endl;
	for(int i=0;i<vect.size();i++)out<<vect[i]<<" ";
	out<<endl;
	return out;
}
template<typename T>
istream& operator>>(istream& in,vector<T>& vect)
{
	vect.clear();
	int n;
	cerr<<"A integer of length wanted"<<endl;
	in>>n;
	vect.resize(n);
	cerr<<n<<" elements wanted"<<endl;
	for(int i=0;i<n;i++)in>>vect[i];
	return in;
}

template<typename T>
ostream& operator<<(ostream& out,const vector<vector<T> >& vect)
{
	out<<"row number="<<vect.size()<<endl;
	for(int i=0;i<vect.size();i++)
	{
		out<<"row #"<<i<<":";
		for(int j=0;j<vect[i].size();j++)
			out<<" "<<vect[i][j];
		out<<endl;
	}
	return out;
}
template<typename T>
istream& operator>>(istream& in,vector<vector<T> >& vect)
{
	vect.clear();
	int n,m;
	cerr<<"Two integers wnated"<<endl;
	in>>n>>m;
	vect.resize(n);
	cerr<<"A matrix "<<n<<" X "<<m<<" wanted"<<endl;
	for(int i=0;i<n;i++)
	{
		vect[i].resize(m);
		for(int j=0;j<m;j++)in>>vect[i][j];
	}
	return in;
}

//**************************************************
// Tool functions
template<typename T>
inline void updateMax(T& a,const T& b){a=max(a,b);}
template<typename T>
inline void updateMin(T& a,const T& b){a=min(a,b);}
template<typename T>
inline vector<T> erase(vector<T> table,int ind)
{
	assert(ind<table.size());
	table.erase(table.begin()+ind);
	return table;
}
template<typename T>
vector<T> unique(vector<T> table)
{
	SORT(table);
	return vector<T>(table.begin(),unique(table.begin(),table.end()));
}

//Be carefull for cut into strings!!!!
template<class T>
vector<T> parse(const string& ss,const char* cut=" ")
{
	vector<T> re;
	for(int j=0;j<ss.size();j++)
	{
		string s;
		while(j<ss.size()&&NULL==strchr(cut,ss[j]))
			s+=ss[j++];
		if(!s.empty())
		{
			T tmp;
			istringstream is(s);
			is>>tmp;
			re.push_back(tmp);
		}
	}
	return re;
}
//****************************************8
/************ bitwise functions  ************/
int countBit(int n)
{
	int re=0;
	while(n)re++,n^=low(n);
	return re;
}
// Most significant bit
int MSB(int n)
{
	if(n==0)return 0;
	while(low(n)!=n)n^=low(n);
	return n;
}
void initHash(int n,int hash[])
{
	for(int i=0;i<n;i++)
		hash[two(i)]=i;
}
void initBcnt(int n,int bcnt[])
{
	bcnt[0]=0;
	for(int i=1;i<two(n);i++)
		bcnt[i]=bcnt[i^low(i)]+1;
}


//************************************************
//**********Number functions*****************
template<typename T>
T __gcd(T n,T m,T& a,T& b)
{
	T a1=0,b1=1;
	a=1,b=0;
	// a*n+b*m=n;  a1*n+b1*m=m;
	while(m)
	{
		T c=n/m;
		T r=n-m*c;
		T t;
		t=a;a=a1;a1=t-c*a1;
		t=b;b=b1;b1=t-c*b1;
		n=m;m=r;
	}
	return n;
}

//*******************************************
struct Node
{
	int r,c,d;
	int v;
	Node(int x,int y,int z,int v1):r(x),c(y),d(z),v(v1){};
	bool operator<(const Node& other)const
	{
		return v>other.v;
	};
};
// A small heap
priority_queue<Node> heap;

///*******************************************
const int dir[][2]={{-1,0},{0,1},{1,0},{0,-1},{1,1},{1,-1},{-1,1},{-1,-1}};
const char dname[]="NWSE";
//**********************************************

//Linear time kmp algorithm. It returns all matched positions
vector<int> KMPmatch(const vector<int>& text,const vector<int>& pattern)
{
	const int n=text.size();
	const int m=pattern.size();
	assert(m);
	vector<int> fail(m);
	fail[0]=-1;
	for(int i=1,k=-1;i<pattern.size();i++)
	{
		while(k!=-1&&pattern[k+1]!=pattern[i])k=fail[k];
		if(pattern[k+1]==pattern[i])k++;
		fail[i]=k;
	}
	vector<int> ret;
	for(int i=0,k=-1;i<n;i++)
	{
		while(k!=-1&&text[i]!=pattern[k+1])k=fail[k];
		if(pattern[k+1]==text[i])k++;
		if(k==m-1)ret.push_back(i-m+1),k=fail[k];
	}
	return ret;
}


//extend[i] = the largest k such that pattern.substr(i,k)=pattern.substr(0,k)
void extended_kmp(int pattern[],int n,int extend[])
{
	extend[0]=1;
	// reach=Max{ extend[j]+j-1 | 0< j <i }
	// idx = arg_max{ extend[j]+j-1 | 0< j <i}
	int reach=0,idx=0;
	for(int i=1;i<n;i++)
	{
		int p=0;
		if(reach>=i)p=min(reach+1-i,extend[i-idx]);
		while(p+i<n&&pattern[p]==pattern[p+i])p++;
		extend[i]=p;
		if(p+i-1>reach)reach=p+i-1,idx=i;
	}
	extend[0]=n;
}
// text[i]= Max{ k | pattern.substr(0,k)=text.substr(i,k) }
// match is supposed to be filled by extended_kmp(pattern,n,match)
void longest_common_prefix(int pattern[],int n,int text[],int m,int extend[]=NULL)
{
	bool allo=false;
	if(extend==NULL)
	{
		allo=true;
		extend=new int[n];
		extended_kmp(pattern,n,extend);
	}
	int idx=-1,reach=-1;
	for(int i=0;i<m;i++)
	{
		int p=i;
		if(reach>=i)p=min(reach+1,extend[i-idx]+i);
		while(p-i<n&&pattern[p-i]==text[p])p++;
		text[i]=p-i;
		if(reach<p-1)reach=p-1,idx=i;
	}
	if(allo)delete [] extend;
}

// extend[i] = Max{ k | str.substr(i-k,2*k+1) is a palindrome}
void palindrom_extend(int str[],int n,int extend[])
{
	// reach = Max{extend[j]+j |  j<i }
	int idx=0,reach=0;
	extend[0]=0;
	for(int i=1;i<n;i++)
	{
		int len=0;
		if(reach>i)len=min(reach-i,extend[2*idx-i]);
		while(i-len-1>=0&&i+len+1<n&&str[i-len-1]==str[i+len+1])len++;
		extend[i]=len;
		if(reach<extend[i]+i)
		{
			reach=extend[i]+i;
			idx=i;
		}
	}
}
char str[100001],rstr[100001];
int buf[100000],extend[100000],fail[100000];
int match[100000],pos[100000];
// Code review!!!!!!!!!!!!!!!!!!!!!!!!
int main()
{
	fgets(str, sizeof(str), stdin);
	const int n=strlen(str);
//	int n=1000;
//	scanf("%d",&n);
//	srand(time(NULL));
	for(int i=0;i<n;i++)
	{
		//buf[i]=str[i]=(rand()%2)+'a';
		buf[i]=str[i];
		rstr[i]=str[i];
	}
//	str[n]=0;
	//printf("%s\n",str);
	reverse(rstr,rstr+n);
	palindrom_extend(buf,n,extend);
//	for(int i=0;i<n;i++)cerr<<extend[i]<<" ";cerr<<endl;
	fail[0]=-1;
	for(int i=1,p=-1;i<n;i++)
	{
		while(p!=-1&&rstr[p+1]!=rstr[i])p=fail[p];
		if(rstr[p+1]==rstr[i])p++;
		fail[i]=p;
	}
	for(int i=0,p=-1;i<n;i++)
	{
		while(p!=-1&&rstr[p+1]!=str[i])p=fail[p];
		if(rstr[p+1]==str[i])p++;
		while(p+1+i+1>n)p=fail[p];
		pos[i]=i;match[i]=p+1;
	//	cerr<<match[i]<<" ";
		if(i&&match[i]<match[i-1])
		{
			match[i]=match[i-1];
			pos[i]=pos[i-1];
		}
	}
//	cerr<<endl;
	int x1=0,l1=0,x2=0,l2=0,x3=0,l3=0;
	int best=0;
	for(int i=0;i<n;i++)
	{
//		debug(i);
		int ll=i-extend[i]-1;
		int rr=i+extend[i]+1;
		int p=0,pl=0;
		int s=0,sl=0;
		if(ll>=0&&rr<n)
		{
			int len=min(match[ll],n-rr);
			p=pos[ll]-len+1;
			assert(p>=0);
			s=n-len;
			pl=sl=len;
		}
		if(pl+sl+2*extend[i]+1>best)
		{
			best=pl+sl+2*extend[i]+1;
/*			debug(i);
			debug(best);
			debug(pl);
			debug(p);
			debug(extend[i]);*/
			x1=p,l1=pl;
			x2=ll+1,l2=2*extend[i]+1;
			x3=s,l3=sl;
		}
	}
	assert(l2);
	int ans=0;
	if(l1)ans++;
	if(l2)ans++;
	if(l3)ans++;
	cout<<ans<<endl;
	if(l1)cout<<x1+1<<" "<<l1<<endl;
	cout<<x2+1<<" "<<l2<<endl;
	if(l3)cout<<x3+1<<" "<<l3<<endl;
	return 0;
	assert(l1==l3);
	for(int i=0;i<l2;i++)assert(str[i+x2]==str[x2+l2-i-1]);
	for(int i=0;i<l1;i++)assert(str[x1+i]==str[x3+l1-1-i]);
	return 0;
}