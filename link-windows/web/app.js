'use strict';
const $=id=>document.getElementById(id),base=new URL('.',location.href);
let stopped=false,busy=false,lastInstall='',lastState='';
function error(e){$('error').hidden=false;$('error').textContent=String(e.message||e);}
function draw(text){const q=new LumenQR(-1,LumenQRError.M);q.addData(text);q.make();const n=q.getModuleCount(),unit=8,c=$('installQR');c.width=c.height=(n+8)*unit;const ctx=c.getContext('2d');ctx.fillStyle='#fff';ctx.fillRect(0,0,c.width,c.height);ctx.fillStyle='#000';for(let y=0;y<n;y++)for(let x=0;x<n;x++)if(q.isDark(y,x))ctx.fillRect((x+4)*unit,(y+4)*unit,unit,unit);}
function render(s){const ready=s.phase==='ready';$('phase').textContent=ready?'接続の準備完了':s.phase==='starting'?'接続を準備しています':'PCの接続設定を確認してください';$('message').textContent=ready&&s.selected?'【Android・Zepp】下のQRを読み取ってください。':s.message;
 $('setup').hidden=ready;$('retry').disabled=s.phase==='starting';$('consent').hidden=true;if(s.setupURL){const u=new URL(s.setupURL);if(u.protocol==='https:'&&u.hostname==='login.tailscale.com'){$('consent').href=u.href;$('consent').hidden=false;}}
 $('select').hidden=s.selected;$('install').hidden=!ready||!s.selected;
 $('request').textContent='ファイル要求：'+(s.lastRequest||'未到達');$('downloadURL').value=s.downloadURL||'';
 if(ready&&s.selected&&lastInstall!==s.installURI){lastInstall=s.installURI;draw(lastInstall);}
}
async function post(path,body=''){const r=await fetch(new URL(path,base),{method:'POST',headers:{'X-Lumen-Action':'1','Content-Type':'application/octet-stream'},body,cache:'no-store'});if(!r.ok)throw Error(await r.text());}
async function update(){const r=await fetch(new URL('state',base),{cache:'no-store'});if(!r.ok)throw Error('接続できません。PCのツールとTailscaleを確認してください。');const state=await r.text();if(state!==lastState){render(JSON.parse(state));lastState=state;}}
async function poll(){if(stopped)return;try{await update();}catch(e){lastState='';$('phase').textContent='接続が終了、または中断しました';}finally{if(!stopped)setTimeout(poll,1500);}}
$('file').addEventListener('change',async()=>{if(busy)return;const f=$('file').files[0];if(!f)return;busy=true;$('error').hidden=true;try{if(f.size>8*1024*1024)throw Error('ZPKは8MB以内で選んでください。');await post('select',await f.arrayBuffer());await update();}catch(e){error(e);}finally{busy=false;}});
$('retry').addEventListener('click',async()=>{try{await post('retry');await update();}catch(e){error(e);}});
$('stop').addEventListener('click',async()=>{try{await post('stop');stopped=true;$('phase').textContent='終了しました';$('message').textContent='このPC画面を閉じて構いません。';for(const id of ['select','install','setup','finish'])$(id).hidden=true;}catch(e){error(e);}});
poll();
