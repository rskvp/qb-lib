function g(){}function W(t,n){for(const e in n)t[e]=n[e];return t}function D(t){return t()}function P(){return Object.create(null)}function $(t){t.forEach(D)}function C(t){return typeof t=="function"}function at(t,n){return t!=t?n==n:t!==n||t&&typeof t=="object"||typeof t=="function"}let v;function ft(t,n){return v||(v=document.createElement("a")),v.href=n,t===v.href}function G(t){return Object.keys(t).length===0}function L(t,...n){if(t==null)return g;const e=t.subscribe(...n);return e.unsubscribe?()=>e.unsubscribe():e}function _t(t){let n;return L(t,e=>n=e)(),n}function dt(t,n,e){t.$$.on_destroy.push(L(n,e))}function ht(t,n,e,r){if(t){const i=q(t,n,e,r);return t[0](i)}}function q(t,n,e,r){return t[1]&&r?W(e.ctx.slice(),t[1](r(n))):e.ctx}function mt(t,n,e,r){if(t[2]&&r){const i=t[2](r(e));if(n.dirty===void 0)return i;if(typeof i=="object"){const l=[],c=Math.max(n.dirty.length,i.length);for(let o=0;o<c;o+=1)l[o]=n.dirty[o]|i[o];return l}return n.dirty|i}return n.dirty}function pt(t,n,e,r,i,l){if(i){const c=q(n,e,r,l);t.p(c,i)}}function yt(t){if(t.ctx.length>32){const n=[],e=t.ctx.length/32;for(let r=0;r<e;r++)n[r]=-1;return n}return-1}function gt(t){const n={};for(const e in t)e[0]!=="$"&&(n[e]=t[e]);return n}function xt(t,n){const e={};n=new Set(n);for(const r in t)!n.has(r)&&r[0]!=="$"&&(e[r]=t[r]);return e}function $t(t){const n={};for(const e in t)n[e]=!0;return n}function bt(t){return t&&C(t.destroy)?t.destroy:g}let E=!1;function J(){E=!0}function K(){E=!1}function Q(t,n,e,r){for(;t<n;){const i=t+(n-t>>1);e(i)<=r?t=i+1:n=i}return t}function R(t){if(t.hydrate_init)return;t.hydrate_init=!0;let n=t.childNodes;if(t.nodeName==="HEAD"){const s=[];for(let u=0;u<n.length;u++){const f=n[u];f.claim_order!==void 0&&s.push(f)}n=s}const e=new Int32Array(n.length+1),r=new Int32Array(n.length);e[0]=-1;let i=0;for(let s=0;s<n.length;s++){const u=n[s].claim_order,f=(i>0&&n[e[i]].claim_order<=u?i+1:Q(1,i,b=>n[e[b]].claim_order,u))-1;r[s]=e[f]+1;const a=f+1;e[a]=s,i=Math.max(a,i)}const l=[],c=[];let o=n.length-1;for(let s=e[i]+1;s!=0;s=r[s-1]){for(l.push(n[s-1]);o>=s;o--)c.push(n[o]);o--}for(;o>=0;o--)c.push(n[o]);l.reverse(),c.sort((s,u)=>s.claim_order-u.claim_order);for(let s=0,u=0;s<c.length;s++){for(;u<l.length&&c[s].claim_order>=l[u].claim_order;)u++;const f=u<l.length?l[u]:null;t.insertBefore(c[s],f)}}function U(t,n){if(E){for(R(t),(t.actual_end_child===void 0||t.actual_end_child!==null&&t.actual_end_child.parentNode!==t)&&(t.actual_end_child=t.firstChild);t.actual_end_child!==null&&t.actual_end_child.claim_order===void 0;)t.actual_end_child=t.actual_end_child.nextSibling;n!==t.actual_end_child?(n.claim_order!==void 0||n.parentNode!==t)&&t.insertBefore(n,t.actual_end_child):t.actual_end_child=n.nextSibling}else(n.parentNode!==t||n.nextSibling!==null)&&t.appendChild(n)}function vt(t,n,e){E&&!e?U(t,n):(n.parentNode!==t||n.nextSibling!=e)&&t.insertBefore(n,e||null)}function V(t){t.parentNode&&t.parentNode.removeChild(t)}function X(t){return document.createElement(t)}function Y(t){return document.createElementNS("http://www.w3.org/2000/svg",t)}function j(t){return document.createTextNode(t)}function wt(){return j(" ")}function Et(){return j("")}function kt(t,n,e,r){return t.addEventListener(n,e,r),()=>t.removeEventListener(n,e,r)}function Z(t,n,e){e==null?t.removeAttribute(n):t.getAttribute(n)!==e&&t.setAttribute(n,e)}function Nt(t,n){const e=Object.getOwnPropertyDescriptors(t.__proto__);for(const r in n)n[r]==null?t.removeAttribute(r):r==="style"?t.style.cssText=n[r]:r==="__value"?t.value=t[r]=n[r]:e[r]&&e[r].set?t[r]=n[r]:Z(t,r,n[r])}function tt(t){return Array.from(t.childNodes)}function nt(t){t.claim_info===void 0&&(t.claim_info={last_index:0,total_claimed:0})}function B(t,n,e,r,i=!1){nt(t);const l=(()=>{for(let c=t.claim_info.last_index;c<t.length;c++){const o=t[c];if(n(o)){const s=e(o);return s===void 0?t.splice(c,1):t[c]=s,i||(t.claim_info.last_index=c),o}}for(let c=t.claim_info.last_index-1;c>=0;c--){const o=t[c];if(n(o)){const s=e(o);return s===void 0?t.splice(c,1):t[c]=s,i?s===void 0&&t.claim_info.last_index--:t.claim_info.last_index=c,o}}return r()})();return l.claim_order=t.claim_info.total_claimed,t.claim_info.total_claimed+=1,l}function z(t,n,e,r){return B(t,i=>i.nodeName===n,i=>{const l=[];for(let c=0;c<i.attributes.length;c++){const o=i.attributes[c];e[o.name]||l.push(o.name)}l.forEach(c=>i.removeAttribute(c))},()=>r(n))}function St(t,n,e){return z(t,n,e,X)}function At(t,n,e){return z(t,n,e,Y)}function et(t,n){return B(t,e=>e.nodeType===3,e=>{const r=""+n;if(e.data.startsWith(r)){if(e.data.length!==r.length)return e.splitText(r.length)}else e.data=r},()=>j(n),!0)}function Ct(t){return et(t," ")}function jt(t,n){n=""+n,t.wholeText!==n&&(t.data=n)}function Ot(t,n){t.value=n??""}function Mt(t,n,e,r){e===null?t.style.removeProperty(n):t.style.setProperty(n,e,r?"important":"")}function Pt(t,n,e){t.classList[e?"add":"remove"](n)}function rt(t,n,{bubbles:e=!1,cancelable:r=!1}={}){const i=document.createEvent("CustomEvent");return i.initCustomEvent(t,e,r,n),i}function Tt(t,n){return new t(n)}let x;function y(t){x=t}function p(){if(!x)throw new Error("Function called outside component initialization");return x}function Dt(t){p().$$.on_mount.push(t)}function Lt(t){p().$$.after_update.push(t)}function qt(t){p().$$.on_destroy.push(t)}function Bt(){const t=p();return(n,e,{cancelable:r=!1}={})=>{const i=t.$$.callbacks[n];if(i){const l=rt(n,e,{cancelable:r});return i.slice().forEach(c=>{c.call(t,l)}),!l.defaultPrevented}return!0}}function zt(t,n){return p().$$.context.set(t,n),n}function Ft(t){return p().$$.context.get(t)}function Ht(t,n){const e=t.$$.callbacks[n.type];e&&e.slice().forEach(r=>r.call(this,n))}const h=[],T=[];let m=[];const N=[],F=Promise.resolve();let S=!1;function H(){S||(S=!0,F.then(I))}function It(){return H(),F}function A(t){m.push(t)}function Wt(t){N.push(t)}const k=new Set;let d=0;function I(){if(d!==0)return;const t=x;do{try{for(;d<h.length;){const n=h[d];d++,y(n),it(n.$$)}}catch(n){throw h.length=0,d=0,n}for(y(null),h.length=0,d=0;T.length;)T.pop()();for(let n=0;n<m.length;n+=1){const e=m[n];k.has(e)||(k.add(e),e())}m.length=0}while(h.length);for(;N.length;)N.pop()();S=!1,k.clear(),y(t)}function it(t){if(t.fragment!==null){t.update(),$(t.before_update);const n=t.dirty;t.dirty=[-1],t.fragment&&t.fragment.p(t.ctx,n),t.after_update.forEach(A)}}function ct(t){const n=[],e=[];m.forEach(r=>t.indexOf(r)===-1?n.push(r):e.push(r)),e.forEach(r=>r()),m=n}const w=new Set;let _;function Gt(){_={r:0,c:[],p:_}}function Jt(){_.r||$(_.c),_=_.p}function st(t,n){t&&t.i&&(w.delete(t),t.i(n))}function Kt(t,n,e,r){if(t&&t.o){if(w.has(t))return;w.add(t),_.c.push(()=>{w.delete(t),r&&(e&&t.d(1),r())}),t.o(n)}else r&&r()}function Qt(t,n){const e={},r={},i={$$scope:1};let l=t.length;for(;l--;){const c=t[l],o=n[l];if(o){for(const s in c)s in o||(r[s]=1);for(const s in o)i[s]||(e[s]=o[s],i[s]=1);t[l]=o}else for(const s in c)i[s]=1}for(const c in r)c in e||(e[c]=void 0);return e}function Rt(t,n,e){const r=t.$$.props[n];r!==void 0&&(t.$$.bound[r]=e,e(t.$$.ctx[r]))}function Ut(t){t&&t.c()}function Vt(t,n){t&&t.l(n)}function ot(t,n,e,r){const{fragment:i,after_update:l}=t.$$;i&&i.m(n,e),r||A(()=>{const c=t.$$.on_mount.map(D).filter(C);t.$$.on_destroy?t.$$.on_destroy.push(...c):$(c),t.$$.on_mount=[]}),l.forEach(A)}function ut(t,n){const e=t.$$;e.fragment!==null&&(ct(e.after_update),$(e.on_destroy),e.fragment&&e.fragment.d(n),e.on_destroy=e.fragment=null,e.ctx=[])}function lt(t,n){t.$$.dirty[0]===-1&&(h.push(t),H(),t.$$.dirty.fill(0)),t.$$.dirty[n/31|0]|=1<<n%31}function Xt(t,n,e,r,i,l,c,o=[-1]){const s=x;y(t);const u=t.$$={fragment:null,ctx:[],props:l,update:g,not_equal:i,bound:P(),on_mount:[],on_destroy:[],on_disconnect:[],before_update:[],after_update:[],context:new Map(n.context||(s?s.$$.context:[])),callbacks:P(),dirty:o,skip_bound:!1,root:n.target||s.$$.root};c&&c(u.root);let f=!1;if(u.ctx=e?e(t,n.props||{},(a,b,...O)=>{const M=O.length?O[0]:b;return u.ctx&&i(u.ctx[a],u.ctx[a]=M)&&(!u.skip_bound&&u.bound[a]&&u.bound[a](M),f&&lt(t,a)),b}):[],u.update(),f=!0,$(u.before_update),u.fragment=r?r(u.ctx):!1,n.target){if(n.hydrate){J();const a=tt(n.target);u.fragment&&u.fragment.l(a),a.forEach(V)}else u.fragment&&u.fragment.c();n.intro&&st(t.$$.fragment),ot(t,n.target,n.anchor,n.customElement),K(),I()}y(s)}class Yt{$destroy(){ut(this,1),this.$destroy=g}$on(n,e){if(!C(e))return g;const r=this.$$.callbacks[n]||(this.$$.callbacks[n]=[]);return r.push(e),()=>{const i=r.indexOf(e);i!==-1&&r.splice(i,1)}}$set(n){this.$$set&&!G(n)&&(this.$$.skip_bound=!0,this.$$set(n),this.$$.skip_bound=!1)}}export{Pt as $,ot as A,ut as B,_t as C,ht as D,U as E,pt as F,yt as G,mt as H,$t as I,W as J,gt as K,kt as L,Ht as M,g as N,$ as O,Y as P,At as Q,ft as R,Yt as S,bt as T,C as U,qt as V,Rt as W,Wt as X,dt as Y,Ot as Z,Nt as _,wt as a,Qt as a0,xt as a1,Bt as a2,Ft as a3,zt as a4,vt as b,Ct as c,Kt as d,Et as e,Jt as f,st as g,V as h,Xt as i,Lt as j,X as k,St as l,tt as m,Z as n,Dt as o,Mt as p,j as q,et as r,at as s,It as t,jt as u,Gt as v,T as w,Tt as x,Ut as y,Vt as z};
