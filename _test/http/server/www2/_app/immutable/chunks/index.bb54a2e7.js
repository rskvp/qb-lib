function b(){}function G(t,e){for(const n in e)t[n]=e[n];return t}function J(t){return!!t&&(typeof t=="object"||typeof t=="function")&&typeof t.then=="function"}function L(t){return t()}function T(){return Object.create(null)}function x(t){t.forEach(L)}function S(t){return typeof t=="function"}function _t(t,e){return t!=t?e==e:t!==e||t&&typeof t=="object"||typeof t=="function"}let $;function dt(t,e){return $||($=document.createElement("a")),$.href=e,t===$.href}function K(t){return Object.keys(t).length===0}function P(t,...e){if(t==null)return b;const n=t.subscribe(...e);return n.unsubscribe?()=>n.unsubscribe():n}function ht(t){let e;return P(t,n=>e=n)(),e}function mt(t,e,n){t.$$.on_destroy.push(P(e,n))}function pt(t,e,n,r){if(t){const c=z(t,e,n,r);return t[0](c)}}function z(t,e,n,r){return t[1]&&r?G(n.ctx.slice(),t[1](r(e))):n.ctx}function yt(t,e,n,r){if(t[2]&&r){const c=t[2](r(n));if(e.dirty===void 0)return c;if(typeof c=="object"){const s=[],i=Math.max(e.dirty.length,c.length);for(let a=0;a<i;a+=1)s[a]=e.dirty[a]|c[a];return s}return e.dirty|c}return e.dirty}function bt(t,e,n,r,c,s){if(c){const i=z(e,n,r,s);t.p(i,c)}}function gt(t){if(t.ctx.length>32){const e=[],n=t.ctx.length/32;for(let r=0;r<n;r++)e[r]=-1;return e}return-1}function xt(t){const e={};for(const n in t)n[0]!=="$"&&(e[n]=t[n]);return e}function $t(t){const e={};for(const n in t)e[n]=!0;return e}function wt(t){return t&&S(t.destroy)?t.destroy:b}let k=!1;function Q(){k=!0}function R(){k=!1}function U(t,e,n,r){for(;t<e;){const c=t+(e-t>>1);n(c)<=r?t=c+1:e=c}return t}function V(t){if(t.hydrate_init)return;t.hydrate_init=!0;let e=t.childNodes;if(t.nodeName==="HEAD"){const l=[];for(let u=0;u<e.length;u++){const f=e[u];f.claim_order!==void 0&&l.push(f)}e=l}const n=new Int32Array(e.length+1),r=new Int32Array(e.length);n[0]=-1;let c=0;for(let l=0;l<e.length;l++){const u=e[l].claim_order,f=(c>0&&e[n[c]].claim_order<=u?c+1:U(1,c,d=>e[n[d]].claim_order,u))-1;r[l]=n[f]+1;const o=f+1;n[o]=l,c=Math.max(o,c)}const s=[],i=[];let a=e.length-1;for(let l=n[c]+1;l!=0;l=r[l-1]){for(s.push(e[l-1]);a>=l;a--)i.push(e[a]);a--}for(;a>=0;a--)i.push(e[a]);s.reverse(),i.sort((l,u)=>l.claim_order-u.claim_order);for(let l=0,u=0;l<i.length;l++){for(;u<s.length&&i[l].claim_order>=s[u].claim_order;)u++;const f=u<s.length?s[u]:null;t.insertBefore(i[l],f)}}function X(t,e){if(k){for(V(t),(t.actual_end_child===void 0||t.actual_end_child!==null&&t.actual_end_child.parentNode!==t)&&(t.actual_end_child=t.firstChild);t.actual_end_child!==null&&t.actual_end_child.claim_order===void 0;)t.actual_end_child=t.actual_end_child.nextSibling;e!==t.actual_end_child?(e.claim_order!==void 0||e.parentNode!==t)&&t.insertBefore(e,t.actual_end_child):t.actual_end_child=e.nextSibling}else(e.parentNode!==t||e.nextSibling!==null)&&t.appendChild(e)}function kt(t,e,n){k&&!n?X(t,e):(e.parentNode!==t||e.nextSibling!=n)&&t.insertBefore(e,n||null)}function Y(t){t.parentNode&&t.parentNode.removeChild(t)}function Z(t){return document.createElement(t)}function tt(t){return document.createElementNS("http://www.w3.org/2000/svg",t)}function j(t){return document.createTextNode(t)}function vt(){return j(" ")}function Et(){return j("")}function Nt(t,e,n,r){return t.addEventListener(e,n,r),()=>t.removeEventListener(e,n,r)}function St(t,e,n){n==null?t.removeAttribute(e):t.getAttribute(e)!==n&&t.setAttribute(e,n)}function et(t){return Array.from(t.childNodes)}function nt(t){t.claim_info===void 0&&(t.claim_info={last_index:0,total_claimed:0})}function D(t,e,n,r,c=!1){nt(t);const s=(()=>{for(let i=t.claim_info.last_index;i<t.length;i++){const a=t[i];if(e(a)){const l=n(a);return l===void 0?t.splice(i,1):t[i]=l,c||(t.claim_info.last_index=i),a}}for(let i=t.claim_info.last_index-1;i>=0;i--){const a=t[i];if(e(a)){const l=n(a);return l===void 0?t.splice(i,1):t[i]=l,c?l===void 0&&t.claim_info.last_index--:t.claim_info.last_index=i,a}}return r()})();return s.claim_order=t.claim_info.total_claimed,t.claim_info.total_claimed+=1,s}function F(t,e,n,r){return D(t,c=>c.nodeName===e,c=>{const s=[];for(let i=0;i<c.attributes.length;i++){const a=c.attributes[i];n[a.name]||s.push(a.name)}s.forEach(i=>c.removeAttribute(i))},()=>r(e))}function jt(t,e,n){return F(t,e,n,Z)}function At(t,e,n){return F(t,e,n,tt)}function rt(t,e){return D(t,n=>n.nodeType===3,n=>{const r=""+e;if(n.data.startsWith(r)){if(n.data.length!==r.length)return n.splitText(r.length)}else n.data=r},()=>j(e),!0)}function Ct(t){return rt(t," ")}function Mt(t,e){e=""+e,t.wholeText!==e&&(t.data=e)}function Ot(t,e,n,r){n===null?t.style.removeProperty(e):t.style.setProperty(e,n,r?"important":"")}function Tt(t,e){return new t(e)}let g;function _(t){g=t}function A(){if(!g)throw new Error("Function called outside component initialization");return g}function qt(t){A().$$.on_mount.push(t)}function Bt(t){A().$$.after_update.push(t)}function Lt(t,e){const n=t.$$.callbacks[e.type];n&&n.slice().forEach(r=>r.call(this,e))}const p=[],q=[];let y=[];const B=[],H=Promise.resolve();let E=!1;function I(){E||(E=!0,H.then(C))}function Pt(){return I(),H}function N(t){y.push(t)}const v=new Set;let m=0;function C(){if(m!==0)return;const t=g;do{try{for(;m<p.length;){const e=p[m];m++,_(e),ct(e.$$)}}catch(e){throw p.length=0,m=0,e}for(_(null),p.length=0,m=0;q.length;)q.pop()();for(let e=0;e<y.length;e+=1){const n=y[e];v.has(n)||(v.add(n),n())}y.length=0}while(p.length);for(;B.length;)B.pop()();E=!1,v.clear(),_(t)}function ct(t){if(t.fragment!==null){t.update(),x(t.before_update);const e=t.dirty;t.dirty=[-1],t.fragment&&t.fragment.p(t.ctx,e),t.after_update.forEach(N)}}function it(t){const e=[],n=[];y.forEach(r=>t.indexOf(r)===-1?e.push(r):n.push(r)),n.forEach(r=>r()),y=e}const w=new Set;let h;function lt(){h={r:0,c:[],p:h}}function ut(){h.r||x(h.c),h=h.p}function W(t,e){t&&t.i&&(w.delete(t),t.i(e))}function st(t,e,n,r){if(t&&t.o){if(w.has(t))return;w.add(t),h.c.push(()=>{w.delete(t),r&&(n&&t.d(1),r())}),t.o(e)}else r&&r()}function zt(t,e){const n=e.token={};function r(c,s,i,a){if(e.token!==n)return;e.resolved=a;let l=e.ctx;i!==void 0&&(l=l.slice(),l[i]=a);const u=c&&(e.current=c)(l);let f=!1;e.block&&(e.blocks?e.blocks.forEach((o,d)=>{d!==s&&o&&(lt(),st(o,1,1,()=>{e.blocks[d]===o&&(e.blocks[d]=null)}),ut())}):e.block.d(1),u.c(),W(u,1),u.m(e.mount(),e.anchor),f=!0),e.block=u,e.blocks&&(e.blocks[s]=u),f&&C()}if(J(t)){const c=A();if(t.then(s=>{_(c),r(e.then,1,e.value,s),_(null)},s=>{if(_(c),r(e.catch,2,e.error,s),_(null),!e.hasCatch)throw s}),e.current!==e.pending)return r(e.pending,0),!0}else{if(e.current!==e.then)return r(e.then,1,e.value,t),!0;e.resolved=t}}function Dt(t,e,n){const r=e.slice(),{resolved:c}=t;t.current===t.then&&(r[t.value]=c),t.current===t.catch&&(r[t.error]=c),t.block.p(r,n)}function Ft(t){t&&t.c()}function Ht(t,e){t&&t.l(e)}function at(t,e,n,r){const{fragment:c,after_update:s}=t.$$;c&&c.m(e,n),r||N(()=>{const i=t.$$.on_mount.map(L).filter(S);t.$$.on_destroy?t.$$.on_destroy.push(...i):x(i),t.$$.on_mount=[]}),s.forEach(N)}function ot(t,e){const n=t.$$;n.fragment!==null&&(it(n.after_update),x(n.on_destroy),n.fragment&&n.fragment.d(e),n.on_destroy=n.fragment=null,n.ctx=[])}function ft(t,e){t.$$.dirty[0]===-1&&(p.push(t),I(),t.$$.dirty.fill(0)),t.$$.dirty[e/31|0]|=1<<e%31}function It(t,e,n,r,c,s,i,a=[-1]){const l=g;_(t);const u=t.$$={fragment:null,ctx:[],props:s,update:b,not_equal:c,bound:T(),on_mount:[],on_destroy:[],on_disconnect:[],before_update:[],after_update:[],context:new Map(e.context||(l?l.$$.context:[])),callbacks:T(),dirty:a,skip_bound:!1,root:e.target||l.$$.root};i&&i(u.root);let f=!1;if(u.ctx=n?n(t,e.props||{},(o,d,...M)=>{const O=M.length?M[0]:d;return u.ctx&&c(u.ctx[o],u.ctx[o]=O)&&(!u.skip_bound&&u.bound[o]&&u.bound[o](O),f&&ft(t,o)),d}):[],u.update(),f=!0,x(u.before_update),u.fragment=r?r(u.ctx):!1,e.target){if(e.hydrate){Q();const o=et(e.target);u.fragment&&u.fragment.l(o),o.forEach(Y)}else u.fragment&&u.fragment.c();e.intro&&W(t.$$.fragment),at(t,e.target,e.anchor,e.customElement),R(),C()}_(l)}class Wt{$destroy(){ot(this,1),this.$destroy=b}$on(e,n){if(!S(n))return b;const r=this.$$.callbacks[e]||(this.$$.callbacks[e]=[]);return r.push(n),()=>{const c=r.indexOf(n);c!==-1&&r.splice(c,1)}}$set(e){this.$$set&&!K(e)&&(this.$$.skip_bound=!0,this.$$set(e),this.$$.skip_bound=!1)}}export{at as A,ot as B,ht as C,pt as D,X as E,bt as F,gt as G,yt as H,$t as I,G as J,xt as K,Nt as L,Lt as M,b as N,x as O,tt as P,At as Q,dt as R,Wt as S,wt as T,S as U,mt as V,zt as W,Dt as X,vt as a,kt as b,Ct as c,st as d,Et as e,ut as f,W as g,Y as h,It as i,Bt as j,Z as k,jt as l,et as m,St as n,qt as o,Ot as p,j as q,rt as r,_t as s,Pt as t,Mt as u,lt as v,q as w,Tt as x,Ft as y,Ht as z};