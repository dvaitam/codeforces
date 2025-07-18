/*+lmake
 * STD = c++14
 * DEFINE += MDEBUG
 */
// Create your own template by modifying this file!
#include <bits/stdc++.h>
using namespace std;
using LL=long long;
using ULL=unsigned long long;
#ifndef ONLINE_JUDGE
     #define debug(args...)            {dbg,args; cerr<<endl;}
#else
    #define debug(args...)              // Just strip off all debug tokens
#endif

struct debugger
{
    template<typename T> debugger& operator , (const T& v)
    {    
        cerr<<v<<" ";    
        return *this;    
    }
} dbg;

template <typename T>
inline void chmax(T &a, T b)
{
    a = std::max(a, b);
}

template <typename T>
inline void chmin(T &a, T b)
{
    a = std::min(a, b);
}

template <size_t _I_Buffer_Size = 1 << 23, size_t _O_Buffer_Size = 1 << 23>
struct IO_Tp
{
    char _I_Buffer[_I_Buffer_Size];
    char *_I_pos;
    const char *_I_end;

    char _O_Buffer[_O_Buffer_Size];
    char *_O_pos;
    const char *_O_end;

    IO_Tp()
        : _I_pos(_I_Buffer)
        , _O_pos(_O_Buffer)
        , _I_end(_I_Buffer + _I_Buffer_Size)
        , _O_end(_O_Buffer + _O_Buffer_Size)
    {
    }

    ~IO_Tp() { fwrite(_O_Buffer, 1, _O_pos - _O_Buffer, stdout); }

    inline char getchar()
    {
        char res = *_I_pos;
        nextchar();
        return res;
    }

    inline bool is_digit(const char ch) { return '0' <= ch && ch <= '9'; }

    inline void nextchar()
    {
        ++_I_pos;
        if (_I_pos == _I_end) {
            fread(_I_Buffer, 1, _I_Buffer_Size, stdin);
            _I_pos = _I_Buffer;
        }
    }

    template <typename Int>
    inline IO_Tp &operator>>(Int &res)
    {
        res = 0;
        int k = 1;
        while (!is_digit(*_I_pos)) {
            if (*_I_pos == '-')
                k = -1;
            nextchar();
        }
        do {
            (res *= 10) += (*_I_pos) - '0';
            nextchar();
        } while (is_digit(*_I_pos));
        res *= k;
        return *this;
    }

    inline IO_Tp &operator>>(char &res)
    {
        do {
            res = *_I_pos;
            nextchar();
        } while (res == ' ' || res == '\0' || res == '\t' || res == '\n' || res == '\r');
        return *this;
    }

    inline void putchar(char x)
    {
        if (_O_pos == _O_end) {
            fwrite(_O_Buffer, 1, _O_pos - _O_Buffer, stdout);
            _O_pos = _O_Buffer;
        }
        *_O_pos++ = x;
    }
    template <typename Int>
    inline IO_Tp &operator<<(Int n)
    {
        if (n < 0) {
            putchar('-');
            n = -n;
        }
        static char _buf[20];
        char *_pos(_buf);
        do
            *_pos++ = '0' + n % 10;
        while (n /= 10);
        while (_pos != _buf)
            putchar(*--_pos);
        return *this;
    }

    inline IO_Tp &operator<<(char ch)
    {
        putchar(ch);
        return *this;
    }

    inline IO_Tp &operator<<(const char *s)
    {
        while (*s != 0) {
            putchar(*s);
            ++s;
        }
        return *this;
    }
};
IO_Tp<> IO;
const int MAXN=100000;
struct E
{
	int u,v,c,id;
}ee[MAXN+10];
int n,m;
int tid[MAXN+10];
vector<int> e[MAXN+10];
int in[MAXN+10];
bool check(int mid)
{
	for(int i=1; i<=n; ++i) {
		e[i].clear();
	}
	memset(in,0,sizeof(in));
	for(int i=mid; i<=m; ++i) {
		e[ee[i].u].push_back(ee[i].v);
		in[ee[i].v]++;
	}
	queue<int> q;
	for(int i=1; i<=n; ++i) {
		if (in[i]==0) q.push(i);
	}
	int id_cnt=0;
	while(!q.empty()) {
		int now=q.front(); q.pop();
		tid[now]=++id_cnt;
		for(auto i:e[now]) {
			in[i]--;
			if (in[i]==0) {
				q.push(i);
			}
		}
	}
	if (id_cnt==n) return true;
	else return false;
}
int main() 
{
#ifdef MDEBUG
	freopen("in.txt","r",stdin);
#endif
	IO>>n>>m;
	for(int i=1; i<=m; ++i) {
		IO>>ee[i].u>>ee[i].v>>ee[i].c;
		ee[i].id=i;
	}
	sort(ee+1,ee+1+m,[](const E& a,const E& b){return a.c<b.c;});
	int l=1,r=m,ans=m+1;
	while(l<=r) {
		int mid=(l+r)/2;
		if (check(mid)) {
			ans=mid;
			r=mid-1;
		} else {
			l=mid+1;
		}
	}
	check(ans);
	IO<<ee[ans-1].c<<' ';
	vector<int> al;
	for(int i=1; i<ans; ++i) {
		if (tid[ee[i].u]>tid[ee[i].v]) {
			al.push_back(ee[i].id);
		}
	}
	sort(al.begin(),al.end());
	IO<<al.size()<<'\n';
	for(auto i:al) {
		IO<<i<<' ';
	}
	IO<<'\n';
	return 0;
}