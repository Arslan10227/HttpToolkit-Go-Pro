const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/shellApi-Cu6Vku3R.js","assets/shellApi-CIY0dWSz.js"])))=>i.map(i=>d[i]);
import { n as e } from "./chunk-CMxvf4Kt.js";
import { o as t, r as n } from "./authConfig-D2sit_fT.js";
let m, h, x, o, p, s, S;
let __tla = (async ()=>{
    let r, i, a;
    s = e((()=>{
        r = `modulepreload`, i = function(e) {
            return `/` + e;
        }, a = {}, o = function(e, t, n) {
            let o = Promise.resolve();
            if (t && t.length > 0) {
                let e = document.getElementsByTagName(`link`), s = document.querySelector(`meta[property=csp-nonce]`), c = s?.nonce || s?.getAttribute(`nonce`);
                function l(e) {
                    return Promise.all(e.map((e)=>Promise.resolve(e).then((e)=>({
                                status: `fulfilled`,
                                value: e
                            }), (e)=>({
                                status: `rejected`,
                                reason: e
                            }))));
                }
                o = l(t.map((t)=>{
                    if (t = i(t, n), t in a) return;
                    a[t] = !0;
                    let o = t.endsWith(`.css`), s = o ? `[rel="stylesheet"]` : ``;
                    if (n) for(let n = e.length - 1; n >= 0; n--){
                        let r = e[n];
                        if (r.href === t && (!o || r.rel === `stylesheet`)) return;
                    }
                    else if (document.querySelector(`link[href="${t}"]${s}`)) return;
                    let l = document.createElement(`link`);
                    if (l.rel = o ? `stylesheet` : r, o || (l.as = `script`), l.crossOrigin = ``, l.href = t, c && l.setAttribute(`nonce`, c), document.head.appendChild(l), o) return new Promise((e, n)=>{
                        l.addEventListener(`load`, e), l.addEventListener(`error`, ()=>n(Error(`Unable to preload CSS for ${t}`)));
                    });
                }));
            }
            function s(e) {
                let t = new Event(`vite:preloadError`, {
                    cancelable: !0
                });
                if (t.payload = e, window.dispatchEvent(t), !t.defaultPrevented) throw e;
            }
            return o.then((t)=>{
                for (let e of t || [])e.status === `rejected` && s(e.reason);
                return e().catch(s);
            });
        };
    }));
    s();
    var c = null;
    function l() {
        try {
            return t.getState().verboseLogging === !0;
        } catch  {
            return !1;
        }
    }
    function u(e, t) {
        return `[${e}] ${t}`;
    }
    async function d(e, t, n, r) {
        f(e, t, n, r);
    }
    function f(e, t, n, r) {
        try {
            let i = window.go?.main?.ShellApp?.FrontendLog;
            if (i) {
                let a = r === void 0 ? `` : ` ${typeof r == `string` ? r : JSON.stringify(r)}`, o = `${u(t, n)}${a}`;
                i(e === `error` ? `error` : `info`, o);
            }
        } catch  {}
    }
    typeof window < `u` && (window.addEventListener(`error`, (e)=>{
        let t = `${e.message} at ${e.filename}:${e.lineno}:${e.colno}`;
        h.error(`window`, `Unhandled error: ${t}`, e.error);
    }), window.addEventListener(`unhandledrejection`, (e)=>{
        let t = e.reason instanceof Error ? e.reason.message : String(e.reason);
        h.error(`window`, `Unhandled promise rejection: ${t}`, e.reason);
    }));
    p = async function() {
        if (!n()) return null;
        if (c) return c;
        try {
            let { getShellApiAsync: e } = await o(async ()=>{
                let { getShellApiAsync: e } = await import(`./shellApi-Cu6Vku3R.js`).then(async (m)=>{
                    await m.__tla;
                    return m;
                });
                return {
                    getShellApiAsync: e
                };
            }, __vite__mapDeps([0,1]));
            return c = await (await e()).getLogPath?.() ?? null, c;
        } catch  {
            return null;
        }
    };
    m = async function() {
        if (n()) try {
            let { getShellApiAsync: e } = await o(async ()=>{
                let { getShellApiAsync: e } = await import(`./shellApi-Cu6Vku3R.js`).then(async (m)=>{
                    await m.__tla;
                    return m;
                });
                return {
                    getShellApiAsync: e
                };
            }, __vite__mapDeps([0,1]));
            await (await e()).openLogFile?.();
        } catch  {}
    };
    let g, _;
    h = {
        verbose (e, t, n) {
            l() && (n === void 0 ? console.log(u(e, t)) : console.log(u(e, t), n), d(`verbose`, e, t, n));
        },
        info (e, t, n) {
            n === void 0 ? console.info(u(e, t)) : console.info(u(e, t), n), d(`info`, e, t, n);
        },
        warn (e, t, n) {
            n === void 0 ? console.warn(u(e, t)) : console.warn(u(e, t), n), d(`warn`, e, t, n);
        },
        error (e, t, n) {
            n === void 0 ? console.error(u(e, t)) : console.error(u(e, t), n), d(`error`, e, t, n);
        }
    };
    g = new Set;
    _ = 200;
    function v(e) {
        let t = 0, n = e.slice(0, 120);
        for(let e = 0; e < n.length; e++)t = (t << 5) - t + n.charCodeAt(e), t |= 0;
        return `m${Math.abs(t)}`;
    }
    function y(e) {
        return e.replace(/https?:\/\/[^\s]+/gi, `[url]`).replace(/[a-f0-9]{32,}/gi, `[hash]`).slice(0, 120);
    }
    function b(e, t) {
        return `${t?.domain ?? `general`}:${t?.code ?? ``}:${v(e)}`;
    }
    x = function(e, t = {}) {
        let n = y(e instanceof Error ? e.message : typeof e == `string` ? e : String(e)), r = b(n, t);
        if (!g.has(r)) {
            if (g.add(r), g.size > _) {
                let e = g.values().next().value;
                e && g.delete(e);
            }
            h.warn(t.domain ?? `app`, n, {
                code: t.code,
                fatal: t.fatal
            });
        }
    };
    S = function() {
        typeof window > `u` || (window.addEventListener(`error`, (e)=>{
            x(e.error ?? e.message, {
                domain: `window`,
                code: `uncaught`
            });
        }), window.addEventListener(`unhandledrejection`, (e)=>{
            x(e.reason, {
                domain: `window`,
                code: `unhandled_rejection`
            });
        }));
    };
})();
export { m as a, h as i, x as n, o, p as r, s, S as t, __tla };
