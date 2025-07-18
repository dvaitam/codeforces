/**
 * Educational Codeforces Round 48 (Rated for Div. 2)
 * Problem ,
 * by Abreto<m@abreto.net>.
 **/
#include <cstdio>

using namespace std;

class abio
{
    static const unsigned BUF_SZ = 65536;
    FILE *istream, *ostream;
    char ibuf[BUF_SZ], obuf[BUF_SZ];
    bool reached_eof;
    size_t ip, isz;
    size_t op, osz;
    inline void clear_ibuf(void) { ip = isz = 0u; }
    inline void clear_obuf(void) { op = osz = 0u; }
    inline void clear_buffer(void)
    {
        reached_eof = false;
        clear_ibuf();
        clear_obuf();
    }
    inline size_t read_buffer(void)
    {
        isz = fread(ibuf, sizeof(char), BUF_SZ, istream);
        ip = 0;
        return isz;
    }
    inline size_t write_buffer(void)
    {
        if(osz)
        {
            size_t ret = fwrite(obuf+op, sizeof(char), osz-op, ostream);
            op += ret;
            if(op == osz) clear_obuf();
            return ret;
        }
        return 0;
    }
    inline abio &reach_eof(void)
    {
        reached_eof = true;
        return (*this);
    }
public:
    static const char endl = '\n';
    abio(FILE *input = stdin, FILE *output = stdout)
    {
        this->istream = input;
        this->ostream = output;
        clear_buffer();
    }
    abio(const char *input, const char *output)
    {
        this->istream = fopen(input, "r");
        this->istream = fopen(output, "w+");
        clear_buffer();
    }
    ~abio(void)
    {
        write_buffer();
        fclose(istream);
        fclose(ostream);
    }
    operator bool() const
    {
        return (!reached_eof);
    }
    inline int getchar(void)
    {
        if(0 == isz || isz == ip) read_buffer();
        if(0 == isz || isz == ip) return EOF;
        return ibuf[ip++];
    }
    inline int putchar(int ch)
    {
        if(osz == BUF_SZ) write_buffer();
        if(osz == BUF_SZ) return EOF;
        return (obuf[osz++] = ch);
    }
    abio &read_int(int &x)
    {
        int flag = 0, ch = getchar();
        if(EOF == ch) return (this->reach_eof());
        x = 0;
        while((EOF!=ch)&&(('-'!=ch)&&((ch<'0')||(ch>'9'))))ch=getchar();
        if(EOF==ch) return (this->reach_eof());
        if('-'==ch){flag=1;ch=getchar();}
        //if(EOF==ch) return (this->reach_eof());
        for(;(((ch>='0')&&(ch<='9')));ch=getchar()){x*=10;x+=(ch-'0');}
        //if(EOF==ch)this->reach_eof();
        if(flag)x*=(-1);
        return (*this);
    }
    abio &read_ll(long long int &x)
    {
        int flag = 0, ch = getchar();
        if(EOF == ch) return (this->reach_eof());
        x = 0ll;
        while((EOF!=ch)&&(('-'!=ch)&&((ch<'0')||(ch>'9'))))ch=getchar();
        if(EOF==ch) return (this->reach_eof());
        if('-'==ch){flag=1;ch=getchar();}
        //if(EOF==ch) return (this->reach_eof());
        for(;(((ch>='0')&&(ch<='9')));ch=getchar()){x*=(10ll);x+=(ch-'0');}
        //if(EOF==ch)this->reach_eof();
        if(flag)x*=(-1);
        return (*this);
    }
    /* set interrupt as '\n' to read a whole line. */
    abio &read_s(char *s, const char interrupt = ' ')
    {
        int ch = getchar();
        while((EOF!=ch)&&(ch<'!'||ch>'~'))ch=getchar();
        if(EOF==ch) return (this->reach_eof());
        for(;(EOF!=ch)&&(interrupt!=ch)&&(ch>=' '&&ch<='~');ch=getchar())(*s++)=ch;
        (*s)=0;
        return (*this);
    }
    abio &write_int(int x, char append = 0)
    {
        int d[20],nd=0;
        if(0==x) putchar('0');
        if(x<0){putchar('-');x=-x;}
        while(x){d[nd++]=x%10;x/=10;}
        while(nd--)putchar('0'+d[nd]);
        if(append)putchar(append);
        return (*this);
    }
    abio &write_ll(long long int x, char append = 0)
    {
        int d[20],nd=0;
        if(0==x) putchar('0');
        if(x<0){putchar('-');x=-x;}
        while(x){d[nd++]=x%10;x/=10;}
        while(nd--)putchar('0'+d[nd]);
        if(append)putchar(append);
        return (*this);
    }
    abio &write_s(const char *s, char append = 0)
    {
        while(*s) putchar(*s++);
        if(append) putchar(append);
        return (*this);
    }
    abio &operator>>(char &ch)
    {
        ch = getchar();
        if(EOF==ch) return (this->reach_eof());
        return (*this);
    }
    abio &operator>>(int &x) { return read_int(x); }
    abio &operator>>(long long int &x) { return read_ll(x); }
    abio &operator>>(char *s) { return read_s(s); }
    abio &operator<<(const char ch) { putchar(ch); return (*this); }
    abio &operator<<(const int x) { return write_int(x); }
    abio &operator<<(const long long int x) { return write_ll(x); }
    abio &operator<<(const char *s) { return write_s(s); }
}io;

class acmachine
{
    static const int maxn = 200020;
public:
    int wa(void)
    {
        int n, m;
        int ai;
        int blank;
        io>>n>>m;
        blank = m;
        while(n--)
        {
            int turn = 0;
            io>>ai;
            if(ai <= blank)
            {
                blank -= ai;
            } else {
                ai -= blank;
                int used_page = (ai + m - 1) / m;
                blank = used_page * m - ai;
                turn = used_page;
            }
            if(0 == blank)
            {
                turn++;
                blank = m;
            }
            io<<turn<<' ';
        }
        io<<io.endl;
        return 0;
    }
};

int main(void)
{
    return (new acmachine())->wa();
}