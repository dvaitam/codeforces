package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type contestInfo struct {
	ID       string
	Path     string
	Problems []string
}

var contests map[string]*contestInfo
var db *sql.DB

var indexTmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>ggg</title>
    <style>
      :root {
        --bg: #0b0f17;
        --panel: #111827;
        --text: #e5e7eb;
        --muted: #9ca3af;
        --accent: #60a5fa;
        --border: #1f2937;
        --good: #10b981;
        --bad: #ef4444;
      }
      html, body { height: 100%; }
      body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Helvetica, Arial, sans-serif; background: var(--bg); color: var(--text); }
      a { color: var(--accent); text-decoration: none; }
      a:hover { text-decoration: underline; }
      .container { max-width: 960px; margin: 0 auto; padding: 24px; }
      .nav { display:flex; gap:16px; align-items:center; padding:16px 0; border-bottom:1px solid var(--border); margin-bottom:24px; }
      .brand { font-weight: 600; color: var(--text); }
      .grid { display:grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap:16px; }
      .card { background: var(--panel); border:1px solid var(--border); border-radius:10px; padding:16px; }
      .muted { color: var(--muted); }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="nav">
        <div class="brand">Codeforces Helper</div>
        <a href="/leaderboard">Leaderboard</a>
        <a href="/model">Models</a>
        <a href="/submission">Submissions</a>
      </div>
      <h1>Contests</h1>
      <div class="grid">
        {{range .}}
        <div class="card">
          <div style="font-size:22px; font-weight:600; margin-bottom:8px;"><a href="/contest/{{.ID}}">#{{.ID}}</a></div>
          <div class="muted">Problems: {{len .Problems}}</div>
        </div>
        {{end}}
      </div>
    </div>
  </body>
</html>`))

var contestTmpl = template.Must(template.New("contest").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Contest {{.ID}}</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Helvetica, Arial, sans-serif; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      .container{ max-width: 960px; margin: 0 auto; padding: 24px; }
      .nav { display:flex; gap:16px; align-items:center; padding:16px 0; border-bottom:1px solid var(--border); margin-bottom:24px; }
      .brand { font-weight:600; }
      .grid { display:grid; grid-template-columns: repeat(auto-fill, minmax(140px,1fr)); gap:12px; }
      .chip { background:var(--panel); border:1px solid var(--border); border-radius:999px; padding:10px 14px; text-align:center; }
      .panel { background:var(--panel); border:1px solid var(--border); border-radius:10px; padding:16px; }
      input, textarea { background:#0f172a; color:var(--text); border:1px solid var(--border); border-radius:8px; padding:10px; width:100%; box-sizing:border-box; }
      label { display:block; margin-top:8px; margin-bottom:6px; color:var(--muted); }
      .row { display:flex; gap:12px; }
      .row > div { flex:1; }
      .btn { background:#1f2937; color:var(--text); border:1px solid var(--border); padding:10px 14px; border-radius:8px; cursor:pointer; }
      .btn:hover { background:#243042; }
      details { margin-top:24px; }
      summary { cursor:pointer; color:var(--muted); }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="nav">
        <a class="brand" href="/">← Contests</a>
        <a href="/leaderboard">Leaderboard</a>
        <a href="/model">Models</a>
        <a href="/submission">Submissions</a>
      </div>
      <h1>Contest {{.ID}}</h1>
      <h2>Problems</h2>
      <div class="grid">
        {{range .Problems}}
          <a class="chip" href="/contest/{{$.ID}}/problem/{{.}}">Problem {{.}}</a>
        {{end}}
      </div>
      <details>
        <summary>Add Problem</summary>
        <div class="panel" style="margin-top:12px;">
          <form action="/addproblem" method="post">
            <input type="hidden" name="contest" value="{{.ID}}">
            <div class="row">
              <div>
                <label>Letter</label>
                <input name="letter" placeholder="A">
              </div>
              <div>
                <label>Admin Key</label>
                <input type="password" name="adminkey" placeholder="••••••">
              </div>
            </div>
            <label style="margin-top:12px;">Problem statement (plain text)</label>
            <textarea name="statement" rows="10" placeholder="Paste the problem text here..."></textarea>
            <div style="margin-top:12px;">
              <button class="btn" type="submit">Add Problem</button>
            </div>
          </form>
        </div>
      </details>
    </div>
  </body>
</html>`))

var addProblemTmpl = template.Must(template.New("addproblem").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Add Problem</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      .container{ max-width: 880px; margin: 0 auto; padding: 24px; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      .panel { background:var(--panel); border:1px solid var(--border); border-radius:10px; padding:16px; }
      label { display:block; color:var(--muted); margin: 8px 0 6px; }
      input, textarea { background:#0f172a; color:var(--text); border:1px solid var(--border); border-radius:8px; padding:10px; width:100%; box-sizing:border-box; }
      .row { display:flex; gap:12px; }
      .row > div { flex:1; }
      .btn { background:#1f2937; color:var(--text); border:1px solid var(--border); padding:10px 14px; border-radius:8px; cursor:pointer; }
      .btn:hover { background:#243042; }
    </style>
  </head>
  <body>
    <div class="container">
      <a href="/">← Back</a>
      <h1 style="margin-top:8px;">Add Problem</h1>
      <div class="panel">
        <form action="/addproblem" method="post">
          <div class="row">
            <div>
              <label>Contest ID</label>
              <input name="contest" value="{{.Contest}}" placeholder="1772">
            </div>
            <div>
              <label>Letter</label>
              <input name="letter" value="{{.Letter}}" placeholder="A">
            </div>
          </div>
          <label>Admin Key</label>
          <input type="password" name="adminkey" placeholder="••••••">
          <label>Problem statement (plain text)</label>
          <textarea name="statement" rows="12" placeholder="Paste the problem text here...">{{.Statement}}</textarea>
          <div style="margin-top:12px;">
            <button class="btn" type="submit">Add Problem</button>
          </div>
        </form>
      </div>
    </div>
  </body>
</html>`))

var problemTmpl = template.Must(template.New("problem").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Contest {{.Contest}} Problem {{.Letter}}</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      .container{ max-width: 1040px; margin: 0 auto; padding: 24px; }
      .nav { display:flex; gap:16px; align-items:center; padding:16px 0; border-bottom:1px solid var(--border); margin-bottom:24px; }
      .panel { background:var(--panel); border:1px solid var(--border); border-radius:10px; padding:16px; }
      .cols { display:grid; grid-template-columns: 1fr 1fr; gap:16px; }
      pre { white-space: pre-wrap; word-wrap: break-word; background:#0f172a; border:1px solid var(--border); border-radius:8px; padding:12px; }
      label { display:block; color:var(--muted); margin: 8px 0 6px; }
      select, input, textarea { background:#0f172a; color:var(--text); border:1px solid var(--border); border-radius:8px; padding:10px; width:100%; box-sizing:border-box; }
      textarea { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace; }
      .btn { background:#1f2937; color:var(--text); border:1px solid var(--border); padding:10px 14px; border-radius:8px; cursor:pointer; }
      .btn:hover { background:#243042; }
      .actions { display:flex; gap:8px; align-items:center; }
      .copy { float:right; font-size:12px; color:var(--muted); cursor:pointer; }
      table { width:100%; border-collapse: collapse; background:var(--panel); border:1px solid var(--border); border-radius:10px; overflow:hidden; }
      th, td { padding:10px 12px; border-bottom:1px solid var(--border); text-align:left; }
      th { color: var(--muted); font-weight:600; }
      tr:hover td { background:#0f172a; }
      @media (max-width: 900px) { .cols { grid-template-columns: 1fr; } }
    </style>
    <script>
      function copyTextFallback(text){
        const ta = document.createElement('textarea');
        ta.value = text;
        ta.setAttribute('readonly','');
        ta.style.position = 'absolute';
        ta.style.left = '-9999px';
        document.body.appendChild(ta);
        const selection = document.getSelection();
        const selected = selection && selection.rangeCount > 0 ? selection.getRangeAt(0) : null;
        ta.select();
        let ok = false;
        try { ok = document.execCommand('copy'); } catch (e) { ok = false; }
        document.body.removeChild(ta);
        if (selected) { selection.removeAllRanges(); selection.addRange(selected); }
        return ok;
      }
      function copyPre(id){
        const el = document.getElementById(id);
        if(!el) return;
        const text = el.innerText;
        const done = () => {
          const msg = document.getElementById(id+"-copied");
          if(msg){ msg.style.opacity = 1; setTimeout(()=> msg.style.opacity = 0, 800); }
        };
        if (navigator.clipboard && navigator.clipboard.writeText) {
          navigator.clipboard.writeText(text).then(done).catch(()=>{
            if (copyTextFallback(text)) done();
          });
        } else {
          if (copyTextFallback(text)) done();
        }
      }

      function mapLang(l){
        if(!l) return 'cpp';
        l = (''+l).trim().toLowerCase();
        if(l === 'c++' || l === 'cpp' || l.includes('cpp')) return 'cpp';
        if(l === 'python' || l === 'py' || l === 'python3' || l.includes('python')) return 'python';
        if(l === 'golang' || l === 'go' || l.includes('golang') || l.startsWith('go')) return 'go';
        if(l === 'rust' || l === 'rs' || l.includes('rust')) return 'rust';
        if(l === 'c') return 'c';
        if(l === 'java' || l.includes('java')) return 'java';
        return l;
      }

      async function retryFrom(kind, id, lang){
        try{
          const url = kind === 'evaluation' ? '/evaluation/raw/response/' + id : '/submission/raw/code/' + id;
          const res = await fetch(url);
          if(!res.ok){ alert('Failed to fetch code'); return; }
          const code = await res.text();
          const form = document.querySelector('form[action$="/submit"]');
          if(!form){ alert('Submit form not found'); return; }
          const sel = form.querySelector('select[name="lang"]');
          const ta = form.querySelector('textarea[name="code"]');
          if(sel) sel.value = mapLang(lang);
          if(ta) ta.value = code;
          form.submit();
        }catch(e){
          console.error(e);
          alert('Retry failed: ' + e);
        }
      }
    </script>
  </head>
  <body>
    <div class="container">
      <div class="nav">
        <a href="/contest/{{.Contest}}">← Contest {{.Contest}}</a>
        <a href="https://codeforces.com/contest/{{.Contest}}/problem/{{.Letter}}" target="_blank" rel="noopener">Open on Codeforces ↗</a>
      </div>
      <h1>Problem {{.Letter}}</h1>
      <div class="cols">
        <div class="panel">
          <div class="actions">
            <div style="flex:1; color:var(--muted);">Statement</div>
            <span class="copy" onclick="copyPre('statement')">Copy</span>
            <span id="statement-copied" style="opacity:0; transition:opacity .2s; color:#10b981;">Copied</span>
          </div>
          <pre id="statement">{{.Statement}}</pre>
        </div>
        <div class="panel">
          <form action="/contest/{{.Contest}}/problem/{{.Letter}}/submit" method="post" enctype="multipart/form-data">
            <label>Language</label>
            <select name="lang">
              <option value="c">C</option>
              <option value="cpp">C++</option>
              <option value="java">Java</option>
              <option value="python">Python 3</option>
              <option value="go">Go</option>
              <option value="rust">Rust</option>
            </select>
            <label style="margin-top:8px;">Paste your solution</label>
            <textarea name="code" rows="18" placeholder="// Paste your solution here"></textarea>
            <label style="margin-top:8px;">…or upload a file</label>
            <input type="file" name="file">
            <div style="margin-top:12px;">
              <button class="btn" type="submit">Submit</button>
            </div>
          </form>
        </div>
      </div>
      <div style="height:40px"></div>
      {{if .Submissions}}
      <h2>Submissions</h2>
      <div class="panel">
        <table>
          <tr><th>ID</th><th>Language</th><th>Exit Code</th><th>Timestamp</th><th>Code</th><th>Stdout</th><th>Stderr</th><th>Retry</th></tr>
          {{range .Submissions}}
          <tr>
            <td><a href="/submission/generate/fix/prompt/{{.ID}}">{{.ID}}</a></td>
            <td>{{.Lang}}</td>
            <td>{{.ExitCode}}</td>
            <td>{{.Timestamp}}</td>
            <td><a href="/submission/code/{{.ID}}">View</a></td>
            <td><a href="/submission/stdout/{{.ID}}">View</a></td>
            <td><a href="/submission/stderr/{{.ID}}">View</a></td>
            <td><a href="#" onclick="retryFrom('submission', {{.ID}}, {{printf "%q" .Lang}}); return false;">Retry</a></td>
          </tr>
          {{end}}
        </table>
      </div>
      {{end}}

      {{if .Evals}}
      <h2 style="margin-top:24px;">Evaluations</h2>
      <div class="panel">
        <table>
          <tr><th>Eval ID</th><th>Run ID</th><th>Provider</th><th>Model</th><th>Lang</th><th>Success</th><th>Timestamp</th><th>Prompt</th><th>Response</th><th>Stdout</th><th>Stderr</th><th>Retry</th></tr>
          {{range .Evals}}
          <tr>
            <td><a href="/evaluation/generate/fix/prompt/{{.ID}}">{{.ID}}</a></td>
            <td>{{.RunID}}</td>
            <td>{{.Provider}}</td>
            <td>{{.Model}}</td>
            <td>{{.Lang}}</td>
            <td>{{.Success}}</td>
            <td>{{.Timestamp}}</td>
            <td><a href="/evaluation/prompt/{{.ID}}">View</a></td>
            <td><a href="/evaluation/response/{{.ID}}">View</a></td>
            <td><a href="/evaluation/stdout/{{.ID}}">View</a></td>
            <td><a href="/evaluation/stderr/{{.ID}}">View</a></td>
            <td><a href="#" onclick="retryFrom('evaluation', {{.ID}}, {{printf "%q" .Lang}}); return false;">Retry</a></td>
          </tr>
          {{end}}
        </table>
      </div>
      {{end}}
    </div>
  </body>
</html>`))

var resultTmpl = template.Must(template.New("result").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Result {{.Contest}}{{.Letter}}</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; --good:#10b981; --bad:#ef4444; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      .container{ max-width: 960px; margin: 0 auto; padding: 24px; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      .panel { background:var(--panel); border:1px solid var(--border); border-radius:10px; padding:16px; }
      pre { white-space: pre-wrap; word-wrap: break-word; background:#0f172a; border:1px solid var(--border); border-radius:8px; padding:12px; }
    </style>
  </head>
  <body>
    <div class="container">
      <a href="/contest/{{.Contest}}/problem/{{.Letter}}">← Back to Problem</a>
      <h1 style="margin-top:8px;">Result for {{.Contest}}{{.Letter}}</h1>
      <div class="panel">
        <pre>{{.Output}}</pre>
      </div>
    </div>
  </body>
</html>`))

var textTmpl = template.Must(template.New("text").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>View</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      .container{ max-width: 960px; margin: 0 auto; padding: 24px; }
      pre { white-space: pre-wrap; word-wrap: break-word; background:#0f172a; border:1px solid var(--border); border-radius:8px; padding:12px; }
      .row { display:flex; gap:12px; align-items:center; }
      .btn { background:#1f2937; color:var(--text); border:1px solid var(--border); padding:8px 12px; border-radius:8px; cursor:pointer; }
      .btn:hover { background:#243042; }
    </style>
    <script>
      function copyTextFallback(text){
        const ta = document.createElement('textarea');
        ta.value = text;
        ta.setAttribute('readonly','');
        ta.style.position = 'absolute';
        ta.style.left = '-9999px';
        document.body.appendChild(ta);
        const selection = document.getSelection();
        const selected = selection && selection.rangeCount > 0 ? selection.getRangeAt(0) : null;
        ta.select();
        let ok = false;
        try { ok = document.execCommand('copy'); } catch (e) { ok = false; }
        document.body.removeChild(ta);
        if (selected) { selection.removeAllRanges(); selection.addRange(selected); }
        return ok;
      }
      async function copyAll(){
        const t = document.getElementById("content").innerText;
        try {
          if (navigator.clipboard && navigator.clipboard.writeText) {
            await navigator.clipboard.writeText(t);
          } else if (!copyTextFallback(t)) {
            throw new Error('No clipboard API and fallback failed');
          }
          const btn = document.getElementById("copyBtn");
          const original = btn.innerText;
          btn.innerText = "Copied!";
          setTimeout(() => { btn.innerText = original; }, 2000);
        } catch (err) {
          console.error('Copy failed', err);
          alert('Copy failed. Select all and copy manually.');
        }
      }
    </script>
  </head>
  <body>
    <div class="container">
      <div class="row">
        <a href="/">← Home</a>
        <button id="copyBtn" class="btn" onclick="copyAll()">Copy</button>
      </div>
      <pre id="content">{{.}}</pre>
    </div>
  </body>
</html>`))

var leaderboardTmpl = template.Must(template.New("leaderboard").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Leaderboard</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; --good:#10b981; --bad:#ef4444; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      .container{ max-width: 1200px; margin: 0 auto; padding: 24px; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      table { width:100%; border-collapse: collapse; background:var(--panel); border:1px solid var(--border); border-radius:10px; overflow:hidden; }
      th, td { padding:10px 12px; border-bottom:1px solid var(--border); text-align:left; }
      th { color: var(--muted); font-weight:600; }
      tr:hover td { background:#0f172a; }
      .nav { display:flex; gap:16px; align-items:center; padding:16px 0; border-bottom:1px solid var(--border); margin-bottom:24px; }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="nav">
        <a href="/">← Contests</a>
        <a href="/model">Models</a>
        <a href="/submission">Submissions</a>
      </div>
      <h1>Leaderboard</h1>
      <table>
        <tr><th>Run ID</th><th>Model</th><th>Lang</th><th>Rating</th><th>Timestamp</th></tr>
        {{range .Leaders}}
        <tr><td><a href="/leaderboard?run={{.RunID}}">{{.RunID}}</a></td><td>{{.Model}}</td><td>{{.Lang}}</td><td>{{.Rating}}</td><td>{{.Timestamp}}</td></tr>
        {{end}}
      </table>
      {{if .Evals}}
      <h2 style="margin-top:24px;">Evaluation History for {{.RunID}}</h2>
      <table>
        <tr><th>Eval ID</th><th>Run ID</th><th>Model</th><th>Lang</th><th>Problem</th><th>Rating</th><th>Success</th><th>Timestamp</th><th>Prompt</th><th>Response</th><th>Stdout</th><th>Stderr</th></tr>
        {{range .Evals}}
        <tr>
          <td><a href="/evaluation/generate/fix/prompt/{{.ID}}">{{.ID}}</a></td>
          <td>{{.RunID}}</td>
          <td>{{.Model}}</td>
          <td>{{.Lang}}</td>
          <td><a href="/contest/{{.ContestID}}/problem/{{.IndexName}}">{{.ContestID}}{{.IndexName}}</a> (<a href="https://codeforces.com/contest/{{.ContestID}}/problem/{{.IndexName}}" target="_blank" rel="noopener">CF</a>)</td>
          <td>{{.Rating}}</td>
          <td>{{.Success}}</td>
          <td>{{.Timestamp}}</td>
          <td><a href="/evaluation/prompt/{{.ID}}">View</a></td>
          <td><a href="/evaluation/response/{{.ID}}">View</a></td>
          <td><a href="/evaluation/stdout/{{.ID}}">View</a></td>
          <td><a href="/evaluation/stderr/{{.ID}}">View</a></td>
        </tr>
        {{end}}
      </table>
      {{end}}
    </div>
  </body>
</html>`))

var modelsTmpl = template.Must(template.New("models").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Models</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      .container{ max-width: 960px; margin: 0 auto; padding: 24px; }
      .grid { display:grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap:12px; }
      .card { background:var(--panel); border:1px solid var(--border); border-radius:10px; padding:16px; }
    </style>
  </head>
  <body>
    <div class="container">
      <a href="/">← Contests</a>
      <h1>Models</h1>
      <div class="grid">
        {{range .}}
          <div class="card"><a href="/model?name={{.}}">{{.}}</a></div>
        {{end}}
      </div>
    </div>
  </body>
</html>`))

var modelTmpl = template.Must(template.New("model").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Evaluations for {{.Model}}</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      .container{ max-width: 1200px; margin: 0 auto; padding: 24px; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      table { width:100%; border-collapse: collapse; background:var(--panel); border:1px solid var(--border); border-radius:10px; overflow:hidden; }
      th, td { padding:10px 12px; border-bottom:1px solid var(--border); text-align:left; }
      th { color: var(--muted); font-weight:600; }
      tr:hover td { background:#0f172a; }
    </style>
  </head>
  <body>
    <div class="container">
      <a href="/model">← Models</a>
      <h1>Evaluations for {{.Model}}</h1>
      <table>
        <tr><th>Eval ID</th><th>Run ID</th><th>Problem</th><th>Rating</th><th>Success</th><th>Timestamp</th><th>Prompt</th><th>Response</th><th>Stdout</th><th>Stderr</th></tr>
        {{range .Evals}}
        <tr><td><a href="/evaluation/generate/fix/prompt/{{.ID}}">{{.ID}}</a></td><td>{{.RunID}}</td><td><a href="/contest/{{.ContestID}}/problem/{{.IndexName}}">{{.ContestID}}{{.IndexName}}</a> (<a href="https://codeforces.com/contest/{{.ContestID}}/problem/{{.IndexName}}" target="_blank" rel="noopener">CF</a>)</td><td>{{.Rating}}</td><td>{{.Success}}</td><td>{{.Timestamp}}</td><td><a href="/evaluation/prompt/{{.ID}}">View</a></td><td><a href="/evaluation/response/{{.ID}}">View</a></td><td><a href="/evaluation/stdout/{{.ID}}">View</a></td><td><a href="/evaluation/stderr/{{.ID}}">View</a></td></tr>
        {{end}}
      </table>
    </div>
  </body>
</html>`))

var failedTmpl = template.Must(template.New("failed").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Failed Evaluations for {{.Model}}</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      .container{ max-width: 1200px; margin: 0 auto; padding: 24px; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      table { width:100%; border-collapse: collapse; background:var(--panel); border:1px solid var(--border); border-radius:10px; overflow:hidden; }
      th, td { padding:10px 12px; border-bottom:1px solid var(--border); text-align:left; }
      th { color: var(--muted); font-weight:600; }
      tr:hover td { background:#0f172a; }
      .controls { display:flex; gap:12px; align-items:center; margin:12px 0; }
      .btn { background:#1f2937; color:var(--text); border:1px solid var(--border); padding:8px 12px; border-radius:8px; cursor:pointer; }
      .btn:hover { background:#243042; }
      .muted { color: var(--muted); }
    </style>
  </head>
  <body>
    <div class="container">
      <a href="/model?name={{.Model}}">← Back to {{.Model}}</a>
      <h1>Failed evaluations for {{.Model}}</h1>
      <div class="controls">
        <form method="get" action="/failed/{{.Model}}" style="display:flex; gap:8px; align-items:center;">
          <label class="muted"><input type="checkbox" name="unreviewed" value="1" {{if .Unreviewed}}checked{{end}}> Only non-reviewed</label>
          <button class="btn" type="submit">Apply</button>
        </form>
        <form id="markForm" method="post" action="/evaluation/mark-reviewed" style="margin-left:auto;">
          <input type="hidden" name="redirect" value="/failed/{{.Model}}{{if .Unreviewed}}?unreviewed=1{{end}}">
          <button class="btn" type="submit">Mark Reviewed</button>
        </form>
      </div>
      <table>
        <tr><th><input type="checkbox" id="checkAll" onclick="toggleAll(this)"></th><th>Eval ID</th><th>Run ID</th><th>Problem</th><th>Rating</th><th>Reviewed</th><th>Timestamp</th><th>Prompt</th><th>Response</th><th>Stdout</th><th>Stderr</th></tr>
        {{range .Evals}}
        <tr>
          <td><input form="markForm" type="checkbox" name="ids" value="{{.ID}}"></td>
          <td><a href="/evaluation/generate/fix/prompt/{{.ID}}">{{.ID}}</a></td>
          <td>{{.RunID}}</td>
          <td><a href="/contest/{{.ContestID}}/problem/{{.IndexName}}">{{.ContestID}}{{.IndexName}}</a> (<a href="https://codeforces.com/contest/{{.ContestID}}/problem/{{.IndexName}}" target="_blank" rel="noopener">CF</a>)</td>
          <td>{{.Rating}}</td>
          <td>{{.Reviewed}}</td>
          <td>{{.Timestamp}}</td>
          <td><a href="/evaluation/prompt/{{.ID}}">View</a></td>
          <td><a href="/evaluation/response/{{.ID}}">View</a></td>
          <td><a href="/evaluation/stdout/{{.ID}}">View</a></td>
          <td><a href="/evaluation/stderr/{{.ID}}">View</a></td>
        </tr>
        {{end}}
      </table>
      <script>
        function toggleAll(src){
          const boxes = document.querySelectorAll('input[name="ids"]');
          boxes.forEach(b => b.checked = src.checked);
        }
      </script>
    </div>
  </body>
</html>`))

var submissionsTmpl = template.Must(template.New("submissions").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Submissions</title>
    <style>
      :root { --bg:#0b0f17; --panel:#111827; --text:#e5e7eb; --muted:#9ca3af; --accent:#60a5fa; --border:#1f2937; }
      body { margin:0; background:var(--bg); color:var(--text); font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif; }
      .container{ max-width: 1200px; margin: 0 auto; padding: 24px; }
      a { color: var(--accent); text-decoration:none; } a:hover{ text-decoration:underline; }
      table { width:100%; border-collapse: collapse; background:var(--panel); border:1px solid var(--border); border-radius:10px; overflow:hidden; }
      th, td { padding:10px 12px; border-bottom:1px solid var(--border); text-align:left; }
      th { color: var(--muted); font-weight:600; }
      tr:hover td { background:#0f172a; }
      .nav { display:flex; gap:16px; align-items:center; padding:16px 0; border-bottom:1px solid var(--border); margin-bottom:24px; }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="nav">
        <a href="/">← Contests</a>
        <a href="/leaderboard">Leaderboard</a>
        <a href="/model">Models</a>
      </div>
      <h1>Submissions</h1>
      <table>
        <tr><th>ID</th><th>Problem</th><th>Language</th><th>Exit Code</th><th>Timestamp</th><th>Code</th><th>Stdout</th><th>Stderr</th></tr>
        {{range .}}
        <tr>
          <td><a href="/submission/generate/fix/prompt/{{.ID}}">{{.ID}}</a></td>
          <td><a href="/contest/{{.ContestID}}/problem/{{.Letter}}">{{.ContestID}}{{.Letter}}</a> (<a href="https://codeforces.com/contest/{{.ContestID}}/problem/{{.Letter}}" target="_blank" rel="noopener">CF</a>)</td>
          <td>{{.Lang}}</td><td>{{.ExitCode}}</td><td>{{.Timestamp}}</td>
          <td><a href="/submission/code/{{.ID}}">View</a></td><td><a href="/submission/stdout/{{.ID}}">View</a></td><td><a href="/submission/stderr/{{.ID}}">View</a></td>
        </tr>
        {{end}}
      </table>
    </div>
  </body>
</html>`))

func scanContests(root string) (map[string]*contestInfo, error) {
	result := make(map[string]*contestInfo)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if _, err := strconv.Atoi(base); err == nil {
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil
			}
			var probs []string
			for _, e := range entries {
				name := e.Name()
				if strings.HasPrefix(name, "problem") && strings.HasSuffix(name, ".txt") {
					letter := strings.TrimSuffix(strings.TrimPrefix(name, "problem"), ".txt")
					probs = append(probs, letter)
				}
			}
			if len(probs) > 0 {
				sort.Strings(probs)
				result[base] = &contestInfo{ID: base, Path: path, Problems: probs}
			}
		}
		return nil
	})
	return result, err
}

func findVerifier(dir, letter string) string {
	cand := filepath.Join(dir, "verifier"+letter+".go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	cand = filepath.Join(dir, "verifier.go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	return ""
}

func detectJavaClassName(src []byte) string {
	re := regexp.MustCompile(`(?m)^\s*public\s+(?:class|interface|enum|record)\s+([A-Za-z_][A-Za-z0-9_]*)`)
	if m := re.FindSubmatch(src); m != nil {
		return string(m[1])
	}
	return "Main"
}

func compileSource(srcPath, lang string) (string, string, error) {
	tmpDir, err := os.MkdirTemp("", "submit")
	if err != nil {
		return "", "", err
	}
	exe := filepath.Join(tmpDir, "main")
	var cmd *exec.Cmd
	switch lang {
	case "c":
		cmd = exec.Command("gcc", srcPath, "-O2", "-std=c11", "-o", exe)
	case "cpp", "c++":
		cmd = exec.Command("g++", srcPath, "-O2", "-std=c++17", "-o", exe)
	case "go":
		cmd = exec.Command("go", "build", "-o", exe, srcPath)
	case "rust":
		cmd = exec.Command("rustc", "-O", srcPath, "-o", exe)
	case "java":
		javaDir := filepath.Join(tmpDir, "java")
		if err := os.Mkdir(javaDir, 0755); err != nil {
			return "", "", err
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return "", "", err
		}
		className := detectJavaClassName(data)
		javaSrc := srcPath
		if filepath.Base(srcPath) != className+".java" {
			javaSrc = filepath.Join(tmpDir, className+".java")
			if err := os.WriteFile(javaSrc, data, 0644); err != nil {
				return "", "", err
			}
		}
		cmd = exec.Command("javac", "-d", javaDir, javaSrc)
		exe = filepath.Join(tmpDir, "run-java.sh")
		script := fmt.Sprintf("#!/bin/sh\njava -cp %s %s \"$@\"\n", javaDir, className)
		if err := os.WriteFile(exe, []byte(script), 0755); err != nil {
			return "", "", err
		}
	case "python":
		absSrc, err := filepath.Abs(srcPath)
		if err != nil {
			return "", "", err
		}
		exe = filepath.Join(tmpDir, "run-python.sh")
		script := fmt.Sprintf("#!/bin/sh\npython3 %s \"$@\"\n", absSrc)
		if err := os.WriteFile(exe, []byte(script), 0755); err != nil {
			return "", "", err
		}
		return exe, "", nil
	default:
		return "", "", fmt.Errorf("unknown language")
	}

	if cmd != nil {
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", string(out), err
		}
		if lang == "java" {
			// ensure wrapper is executable
			return exe, string(out), nil
		}
	}
	return exe, "", nil
}

func submitSolution(w http.ResponseWriter, r *http.Request, c *contestInfo, letter string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
    lang := r.FormValue("lang")
    // Normalize language variants from UI or retry helper
    l := strings.ToLower(strings.TrimSpace(lang))
    if l == "c++" || strings.Contains(l, "cpp") {
        l = "cpp"
    } else if l == "py" || l == "python3" || strings.Contains(l, "python") {
        l = "python"
    } else if l == "golang" || l == "go" || strings.Contains(l, "golang") || strings.HasPrefix(l, "go") {
        // Accept variants like "Go (Golang)"
        l = "go"
    } else if l == "rs" || strings.Contains(l, "rust") {
        l = "rust"
    } else if l == "java" || strings.Contains(l, "java") {
        l = "java"
    } else if l == "c" {
        l = "c"
    }
    lang = l
    var data []byte
    file, _, err := r.FormFile("file")
    if err == nil {
        defer file.Close()
        data, _ = io.ReadAll(file)
    } else {
        data = []byte(r.FormValue("code"))
    }
    extMap := map[string]string{
        "c": ".c",
        "cpp": ".cpp",
        "c++": ".cpp",
        "java": ".java",
        "python": ".py",
        "python3": ".py",
        "py": ".py",
        "go": ".go",
        "golang": ".go",
        "rust": ".rs",
        "rs": ".rs",
    }
    ext := extMap[lang]
    if ext == "" {
        // Heuristic fallback: detect language from code contents (helps retries with odd labels)
        code := strings.ToLower(string(data))
        switch {
        case strings.Contains(code, "package main") && strings.Contains(code, "func main"):
            lang, ext = "go", ".go"
        case strings.Contains(code, "#include") || strings.Contains(code, "int main("):
            // Try to disambiguate C++ from C using common C++ markers
            if strings.Contains(code, "std::") || strings.Contains(code, "using namespace std") || strings.Contains(code, "<iostream>") {
                lang, ext = "cpp", ".cpp"
            } else {
                lang, ext = "c", ".c"
            }
        case strings.Contains(code, "using namespace std") || strings.Contains(code, "std::") || strings.Contains(code, "#include <iostream>"):
            lang, ext = "cpp", ".cpp"
        case strings.Contains(code, "public class") && strings.Contains(code, "static void main"):
            lang, ext = "java", ".java"
        case strings.Contains(code, "def main(") || strings.Contains(code, "#!/usr/bin/env python") || strings.Contains(code, "print("):
            lang, ext = "python", ".py"
        case strings.Contains(code, "fn main()") || strings.Contains(code, "use std::"):
            lang, ext = "rust", ".rs"
        }
        if ext == "" {
            http.Error(w, "unknown language", http.StatusBadRequest)
            return
        }
    }
	srcPath := filepath.Join(c.Path, "user"+strings.ToUpper(letter)+ext)
	if err := os.WriteFile(srcPath, data, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exe, compileOut, err := compileSource(srcPath, lang)
	output := bytes.Buffer{}
	exitCode := -1
	stdoutStr := ""
	stderrStr := ""
	if err != nil {
		output.WriteString("Compilation failed:\n")
		output.WriteString(compileOut)
		output.WriteString(err.Error())
		stderrStr = compileOut + err.Error()
	} else {
		verifier := findVerifier(c.Path, letter)
		if verifier != "" {
            ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "go", "run", filepath.Base(verifier), exe)
			cmd.Dir = c.Path
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

			var stdoutBuf, stderrBuf bytes.Buffer
			cmd.Stdout = &stdoutBuf
			cmd.Stderr = &stderrBuf

			if err := cmd.Start(); err != nil {
				output.WriteString("Verifier error: " + err.Error())
			} else {
				errCh := make(chan error, 1)
				go func() {
					errCh <- cmd.Wait()
				}()

				var runErr error
				select {
				case <-ctx.Done():
					if cmd.Process != nil {
						pgid, _ := syscall.Getpgid(cmd.Process.Pid)
						syscall.Kill(-pgid, syscall.SIGKILL)
					}
					runErr = <-errCh
				case runErr = <-errCh:
				}

				stdoutStr = stdoutBuf.String()
				stderrStr = stderrBuf.String()
				output.WriteString(stdoutStr)
				output.WriteString(stderrStr)

                if ctx.Err() == context.DeadlineExceeded {
                    output.WriteString("\nVerifier timed out after 120 seconds")
				} else if runErr != nil {
					if ee, ok := runErr.(*exec.ExitError); ok {
						exitCode = ee.ExitCode()
						output.WriteString(fmt.Sprintf("\nVerifier exited with status %d", ee.ExitCode()))
					} else {
						output.WriteString("\nVerifier error: " + runErr.Error())
					}
				} else {
					exitCode = 0
				}
			}
		} else {
			output.WriteString("Compiled successfully. No verifier available.")
		}
	}
	respStr := output.String()
	if _, dbErr := db.Exec("INSERT INTO submissions (contest_id, problem_letter, lang, code, stdout, stderr, response, exit_code) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", c.ID, letter, lang, string(data), stdoutStr, stderrStr, respStr, exitCode); dbErr != nil {
		fmt.Println("failed to insert submission:", dbErr)
	}
	resultTmpl.Execute(w, map[string]string{
		"Contest": c.ID,
		"Letter":  letter,
		"Output":  respStr,
	})
}

func problemPage(w http.ResponseWriter, r *http.Request, c *contestInfo, letter string) {
    stmtPath := filepath.Join(c.Path, "problem"+letter+".txt")
    data, err := os.ReadFile(stmtPath)
    if err != nil {
        http.Error(w, "problem not found", http.StatusNotFound)
        return
    }
    // Load submissions for this problem
    type Sub struct {
        ID        int
        Lang      string
        ExitCode  int
        Timestamp string
    }
    var subs []Sub
    if db != nil {
        rows, err := db.Query("SELECT id, lang, exit_code, timestamp FROM submissions WHERE contest_id = ? AND problem_letter = ? ORDER BY id DESC", c.ID, letter)
        if err == nil {
            defer rows.Close()
            for rows.Next() {
                var s Sub
                if err = rows.Scan(&s.ID, &s.Lang, &s.ExitCode, &s.Timestamp); err == nil {
                    subs = append(subs, s)
                }
            }
        }
    }

    // Load evaluations for this problem (by joining problems)
    type Eval struct {
        ID        int
        RunID     string
        Provider  string
        Model     string
        Lang      string
        Success   bool
        Timestamp string
    }
    var evals []Eval
    if db != nil {
        rows, err := db.Query(`SELECT e.id, e.run_id, COALESCE(e.provider, ''), e.model, e.lang, e.success, e.timestamp
                                FROM evaluations e
                                JOIN problems p ON e.problem_id = p.id
                                WHERE p.contest_id = ? AND p.index_name = ?
                                ORDER BY e.timestamp DESC`, c.ID, letter)
        if err == nil {
            defer rows.Close()
            for rows.Next() {
                var e Eval
                if err = rows.Scan(&e.ID, &e.RunID, &e.Provider, &e.Model, &e.Lang, &e.Success, &e.Timestamp); err == nil {
                    evals = append(evals, e)
                }
            }
        }
    }

    problemTmpl.Execute(w, map[string]interface{}{
        "Contest":     c.ID,
        "Letter":      letter,
        "Statement":   string(data),
        "Submissions": subs,
        "Evals":       evals,
    })
}

func contestPage(w http.ResponseWriter, r *http.Request, cid string) {
	c := contests[cid]
	if c == nil {
		http.NotFound(w, r)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/contest/"+cid), "/")
	if len(parts) >= 3 && parts[1] == "problem" && parts[2] != "" {
		if len(parts) == 4 && parts[3] == "submit" {
			submitSolution(w, r, c, parts[2])
			return
		}
		if r.Method == http.MethodGet {
			problemPage(w, r, c, parts[2])
			return
		}
	}
	if r.URL.Path != "/contest/"+cid {
		http.NotFound(w, r)
		return
	}
	contestTmpl.Execute(w, c)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ids := make([]string, 0, len(contests))
	for id := range contests {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var list []*contestInfo
	for _, id := range ids {
		list = append(list, contests[id])
	}
	indexTmpl.Execute(w, list)
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/path" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, os.Getenv("PATH"))
}

func contestDir(id string) (string, error) {
	n, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	thousands := (n / 1000) * 1000
	tDir := fmt.Sprintf("%d-%d", thousands, thousands+999)
	n %= 1000
	hundreds := (n / 100) * 100
	hDir := fmt.Sprintf("%d-%d", thousands+hundreds, thousands+hundreds+99)
	n %= 100
	tens := (n / 10) * 10
	teDir := fmt.Sprintf("%d-%d", thousands+hundreds+tens, thousands+hundreds+tens+9)
	return filepath.Join(tDir, hDir, teDir, id), nil
}

func ensureContest(id string) (*contestInfo, error) {
	if c := contests[id]; c != nil {
		return c, nil
	}
	dir, err := contestDir(id)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	c := &contestInfo{ID: id, Path: dir}
	contests[id] = c
	return c, nil
}

func addProblemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		addProblemTmpl.Execute(w, map[string]string{
			"Contest":   r.URL.Query().Get("contest"),
			"Letter":    r.URL.Query().Get("letter"),
			"Statement": "",
		})
		return
	case http.MethodPost:
		contestID := r.FormValue("contest")
		letter := strings.ToUpper(r.FormValue("letter"))
		statement := r.FormValue("statement")
		adminkey := r.FormValue("adminkey")
		if adminkey != os.Getenv("ADMIN_KEY") {
			http.Error(w, "admin key mismatch", http.StatusForbidden)
			return
		}
		if contestID == "" || letter == "" {
			http.Error(w, "missing parameters", http.StatusBadRequest)
			return
		}
		c, err := ensureContest(contestID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stmtPath := filepath.Join(c.Path, "problem"+letter+".txt")
		if _, err := os.Stat(stmtPath); err == nil {
			http.Error(w, "problem already exists", http.StatusBadRequest)
			return
		}
		if err := os.WriteFile(stmtPath, []byte(statement), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		found := false
		for _, l := range c.Problems {
			if l == letter {
				found = true
				break
			}
		}
		if !found {
			c.Problems = append(c.Problems, letter)
			sort.Strings(c.Problems)
		}
		http.Redirect(w, r, "/contest/"+contestID+"/problem/"+letter, http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/leaderboard" {
		http.NotFound(w, r)
		return
	}

	type Leader struct {
		RunID     string
		Model     string
		Lang      string
		Rating    int
		Timestamp string
	}
	var leaders []Leader
	rows, err := db.Query("SELECT run_id, model, lang, rating, timestamp FROM leaderboard ORDER BY rating DESC")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var l Leader
			if err = rows.Scan(&l.RunID, &l.Model, &l.Lang, &l.Rating, &l.Timestamp); err == nil {
				leaders = append(leaders, l)
			}
		}
	}

	type Eval struct {
		ID        int
		RunID     string
		Model     string
		Lang      string
		ProblemID int
		ContestID int
		IndexName string
		Rating    int
		Success   bool
		Timestamp string
	}
	var evals []Eval
	runIDFilter := r.URL.Query().Get("run")
	if runIDFilter != "" {
		rows, err = db.Query(`SELECT e.id, e.run_id, e.model, e.lang, e.problem_id, p.contest_id, p.index_name, COALESCE(p.rating, 0), e.success, e.timestamp
                       FROM evaluations e
                       JOIN problems p ON e.problem_id = p.id
                       WHERE e.run_id = ? ORDER BY e.timestamp DESC`, runIDFilter)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var e Eval
				if err = rows.Scan(&e.ID, &e.RunID, &e.Model, &e.Lang, &e.ProblemID, &e.ContestID, &e.IndexName, &e.Rating, &e.Success, &e.Timestamp); err == nil {
					evals = append(evals, e)
				}
			}
		}
	}

	leaderboardTmpl.Execute(w, map[string]interface{}{
		"Leaders": leaders,
		"Evals":   evals,
		"RunID":   runIDFilter,
	})
}

func modelHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/model" {
		http.NotFound(w, r)
		return
	}
	modelName := r.URL.Query().Get("name")
	if modelName == "" {
		var models []string
		rows, err := db.Query("SELECT DISTINCT model FROM evaluations ORDER BY model")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var m string
				if err = rows.Scan(&m); err == nil {
					models = append(models, m)
				}
			}
		}
		modelsTmpl.Execute(w, models)
		return
	}

	type Eval struct {
		ID        int
		RunID     string
		ProblemID int
		ContestID int
		IndexName string
		Rating    int
		Success   bool
		Timestamp string
	}
	var evals []Eval
	rows, err := db.Query(`SELECT e.id, e.run_id, e.problem_id, p.contest_id, p.index_name, COALESCE(p.rating, 0), e.success, e.timestamp
                               FROM evaluations e
                               JOIN problems p ON e.problem_id = p.id
                               WHERE e.model = ? ORDER BY e.timestamp DESC`, modelName)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Eval
			if err = rows.Scan(&e.ID, &e.RunID, &e.ProblemID, &e.ContestID, &e.IndexName, &e.Rating, &e.Success, &e.Timestamp); err == nil {
				evals = append(evals, e)
			}
		}
	}

	modelTmpl.Execute(w, map[string]interface{}{
		"Model": modelName,
		"Evals": evals,
	})
}

func failedHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/failed/") {
		http.NotFound(w, r)
		return
	}
	modelName := strings.TrimPrefix(r.URL.Path, "/failed/")
	if modelName == "" {
		http.NotFound(w, r)
		return
	}
	onlyUnreviewed := r.URL.Query().Get("unreviewed") == "1"
	type Eval struct {
		ID        int
		RunID     string
		ProblemID int
		ContestID int
		IndexName string
		Rating    int
		Reviewed  int
		Timestamp string
	}
	var evals []Eval
	cond := "WHERE e.model = ? AND e.success = 0"
	var args []interface{}
	args = append(args, modelName)
	if onlyUnreviewed {
		cond += " AND COALESCE(e.reviewied, 0) = 0"
	}
	q := `SELECT e.id, e.run_id, e.problem_id, p.contest_id, p.index_name, COALESCE(p.rating, 0), COALESCE(e.reviewied, 0), e.timestamp
            FROM evaluations e
            JOIN problems p ON e.problem_id = p.id ` + cond + ` ORDER BY e.timestamp DESC`
	rows, err := db.Query(q, args...)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Eval
			if err = rows.Scan(&e.ID, &e.RunID, &e.ProblemID, &e.ContestID, &e.IndexName, &e.Rating, &e.Reviewed, &e.Timestamp); err == nil {
				evals = append(evals, e)
			}
		}
	}
	failedTmpl.Execute(w, map[string]interface{}{
		"Model":      modelName,
		"Evals":      evals,
		"Unreviewed": onlyUnreviewed,
	})
}

func markReviewedHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/evaluation/mark-reviewed" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	ids := r.FormValue("ids")
	// Also support multiple ids values
	idVals := r.Form["ids"]
	if ids != "" && len(idVals) == 0 {
		idVals = []string{ids}
	}
	if len(idVals) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Build IN clause safely
	var placeholders []string
	var args []interface{}
	for _, s := range idVals {
		if s == "" {
			continue
		}
		if _, err := strconv.Atoi(s); err != nil {
			continue
		}
		placeholders = append(placeholders, "?")
		args = append(args, s)
	}
	if len(placeholders) > 0 {
		q := "UPDATE evaluations SET reviewied = 1 WHERE id IN (" + strings.Join(placeholders, ",") + ")"
		if _, err := db.Exec(q, args...); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func failedJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/failed/json" {
		http.NotFound(w, r)
		return
	}
	modelName := r.URL.Query().Get("model")
	if modelName == "" {
		http.Error(w, "missing model query parameter", http.StatusBadRequest)
		return
	}
	langFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("lang")))
	type Eval struct {
		ID        int    `json:"id"`
		RunID     string `json:"run_id"`
		ProblemID int    `json:"problem_id"`
		ContestID int    `json:"contest_id"`
		IndexName string `json:"index_name"`
		Rating    int    `json:"rating"`
		Reviewed  int    `json:"reviewied"`
		Timestamp string `json:"timestamp"`
	}
	var evals []Eval
	cond := "WHERE e.model = ? AND e.success = 0"
	var args []interface{}
	args = append(args, modelName)
	if langFilter != "" {
		cond += " AND e.lang = ?"
		args = append(args, langFilter)
	}
	q := `SELECT e.id, e.run_id, e.problem_id, p.contest_id, p.index_name, COALESCE(p.rating, 0), COALESCE(e.reviewied,0), e.timestamp
            FROM evaluations e
            JOIN problems p ON e.problem_id = p.id ` + cond + ` ORDER BY e.timestamp DESC`
	rows, err := db.Query(q, args...)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Eval
			if err = rows.Scan(&e.ID, &e.RunID, &e.ProblemID, &e.ContestID, &e.IndexName, &e.Rating, &e.Reviewed, &e.Timestamp); err == nil {
				evals = append(evals, e)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"model": modelName,
		"lang":  langFilter,
		"evals": evals,
	})
}

func submissionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/submission" {
		http.NotFound(w, r)
		return
	}
    type Sub struct {
        ID        int
        ContestID string
        Letter    string
        Lang      string
        ExitCode  int
        Timestamp string
    }
    var subs []Sub
    rows, err := db.Query("SELECT id, contest_id, problem_letter, lang, exit_code, timestamp FROM submissions ORDER BY id DESC")
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var s Sub
            if err = rows.Scan(&s.ID, &s.ContestID, &s.Letter, &s.Lang, &s.ExitCode, &s.Timestamp); err == nil {
                subs = append(subs, s)
            }
        }
    }
    submissionsTmpl.Execute(w, subs)
}

func submissionFixPromptHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/submission/generate/fix/prompt/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var contestID, letter, code, stdout, stderr string
	err = db.QueryRow("SELECT contest_id, problem_letter, code, stdout, stderr FROM submissions WHERE id = ?", id).Scan(&contestID, &letter, &code, &stdout, &stderr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	output := strings.TrimSpace(stdout + "\n" + stderr)
	prompt := fmt.Sprintf("For problem %s%s this is a correct solution, but verifier ends with %s can you fix the verifier? %s", contestID, letter, output, code)
	textTmpl.Execute(w, prompt)
}

func submissionContentHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/submission/"), "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	field := parts[0]
	if field != "code" && field != "stdout" && field != "stderr" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var content string
	err = db.QueryRow("SELECT "+field+" FROM submissions WHERE id = ?", id).Scan(&content)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	textTmpl.Execute(w, content)
}

func submissionRawCodeHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/submission/raw/code/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    var code string
    err = db.QueryRow("SELECT code FROM submissions WHERE id = ?", id).Scan(&code)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    _, _ = w.Write([]byte(code))
}

func evaluationFixPromptHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/evaluation/generate/fix/prompt/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var contestID, letter, code, stdout, stderr string
	err = db.QueryRow(`SELECT p.contest_id, p.index_name, e.response, e.stdout, e.stderr FROM evaluations e JOIN problems p ON e.problem_id = p.id WHERE e.id = ?`, id).Scan(&contestID, &letter, &code, &stdout, &stderr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	output := strings.TrimSpace(stdout + "\n" + stderr)
	prompt := fmt.Sprintf("For problem %s%s this is a correct solution, but verifier ends with %s can you fix the verifier? %s", contestID, letter, output, code)
	textTmpl.Execute(w, prompt)
}

func evaluationContentHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/evaluation/"), "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	field := parts[0]
	id, err := strconv.Atoi(parts[1])
	if err != nil || (field != "prompt" && field != "response" && field != "stdout" && field != "stderr") {
		http.NotFound(w, r)
		return
	}
	var content string
	err = db.QueryRow("SELECT "+field+" FROM evaluations WHERE id = ?", id).Scan(&content)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	textTmpl.Execute(w, content)
}

func evaluationRawResponseHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/evaluation/raw/response/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var resp string
	if err := db.QueryRow("SELECT response FROM evaluations WHERE id = ?", id).Scan(&resp); err != nil {
		http.NotFound(w, r)
		return
	}
	// Extract code block only if fenced with go or rust
	code := ""
	if m := regexp.MustCompile(`(?s)\x60\x60\x60go\s*(.*?)\x60\x60\x60`).FindStringSubmatch(resp); len(m) > 1 {
		code = strings.TrimSpace(m[1])
	} else if m := regexp.MustCompile(`(?s)\x60\x60\x60rust\s*(.*?)\x60\x60\x60`).FindStringSubmatch(resp); len(m) > 1 {
		code = strings.TrimSpace(m[1])
	} else {
		code = resp
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(code))
}

func main() {
	var err error
	fmt.Println("PATH:", os.Getenv("PATH"))
	contests, err = scanContests(".")
	if err != nil {
		panic(err)
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "user:pass@tcp(127.0.0.1:3306)/dbname"
	}
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if _, err = db.Exec(`
               CREATE TABLE IF NOT EXISTS submissions (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       contest_id VARCHAR(20),
                       problem_letter VARCHAR(10),
                       lang VARCHAR(20),
                       code TEXT,
                       stdout TEXT,
                       stderr TEXT,
                       response TEXT,
                       exit_code INT,
                       timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
               )
       `); err != nil {
		panic(err)
	}
	// Ensure evaluations table has reviewied column
	if _, err = db.Exec(`ALTER TABLE evaluations ADD COLUMN IF NOT EXISTS reviewied TINYINT DEFAULT 0`); err != nil {
		// Fallback for MySQL versions without IF NOT EXISTS
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
			// Try without IF NOT EXISTS and ignore duplicate error
			if _, err2 := db.Exec(`ALTER TABLE evaluations ADD COLUMN reviewied TINYINT DEFAULT 0`); err2 != nil {
				if !strings.Contains(strings.ToLower(err2.Error()), "duplicate column") {
					fmt.Println("warning: could not ensure reviewied column:", err2)
				}
			}
		}
	}
	// Ensure evaluations table has provider column
	if _, err = db.Exec(`ALTER TABLE evaluations ADD COLUMN IF NOT EXISTS provider VARCHAR(255)`); err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
			if _, err2 := db.Exec(`ALTER TABLE evaluations ADD COLUMN provider VARCHAR(255)`); err2 != nil {
				if !strings.Contains(strings.ToLower(err2.Error()), "duplicate column") {
					fmt.Println("warning: could not ensure provider column:", err2)
				}
			}
		}
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/path", pathHandler)
	http.HandleFunc("/addproblem", addProblemHandler)
	http.HandleFunc("/contest/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/contest/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}
		contestPage(w, r, parts[0])
	})
	http.HandleFunc("/leaderboard", leaderboardHandler)
	http.HandleFunc("/model", modelHandler)
	http.HandleFunc("/failed/", failedHandler)
	http.HandleFunc("/evaluation/mark-reviewed", markReviewedHandler)
	http.HandleFunc("/failed/json", failedJSONHandler)
	http.HandleFunc("/submission", submissionsHandler)
	http.HandleFunc("/submission/raw/code/", submissionRawCodeHandler)
	http.HandleFunc("/submission/generate/fix/prompt/", submissionFixPromptHandler)
	http.HandleFunc("/submission/", submissionContentHandler)
	http.HandleFunc("/evaluation/generate/fix/prompt/", evaluationFixPromptHandler)
	http.HandleFunc("/evaluation/", evaluationContentHandler)
	// raw code response
	http.HandleFunc("/evaluation/raw/response/", evaluationRawResponseHandler)
	http.ListenAndServe(":8081", nil)
}
