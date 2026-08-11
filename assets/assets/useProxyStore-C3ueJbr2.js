const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/useNotificationStore-NSh_n0Ie.js","assets/useNotificationStore-80beXspH.js","assets/react-DjOMyCbT.js","assets/chunk-CMxvf4Kt.js","assets/dist-CFICGU_H.js","assets/errorHandler-D84ZyjpV.js","assets/errorHandler-CyIT4NLx.js","assets/errorReporting-C7itArmc.js","assets/authConfig-D2sit_fT.js","assets/proxyAdminClient-D-l-Dsud.js","assets/proxyAdminClient-FRHE1_Zo.js","assets/stream-browserify-TNxEQb9t.js","assets/url-v4rlrrTJ.js","assets/__vite-browser-external-CM3Eme-n.js","assets/dist-B3kALMnJ.js","assets/shellApi-CIY0dWSz.js","assets/har-BulMv3VR.js","assets/bodyRegistry-DeYe45A6.js","assets/helpers-C32VkB97.js","assets/json-schema-draft-06-Ck4H_zqx.js","assets/useAuthStore-CEO3CpNF.js","assets/useAuthStore-BMJSyawZ.js"])))=>i.map(i=>d[i]);
import { o as e } from "./chunk-CMxvf4Kt.js";
import { n as t, t as n } from "./react-DjOMyCbT.js";
import { i as r, n as i, o as a, s as o, __tla as __tla_0 } from "./errorReporting-C7itArmc.js";
import { o as s, s as c } from "./authConfig-D2sit_fT.js";
import { t as l } from "./useNotificationStore-80beXspH.js";
import { t as u } from "./errorHandler-CyIT4NLx.js";
import { i as d, r as f } from "./shellApi-CIY0dWSz.js";
import { C as ee, S as te, _ as ne, a as re, b as ie, c as ae, i as oe, l as se, n as p, o as ce, s as le, t as ue, u as de, x as m, y as fe } from "./proxyAdminClient-FRHE1_Zo.js";
import { n as h, r as pe, t as me } from "./dist-B3kALMnJ.js";
import { t as he } from "./useAuthStore-BMJSyawZ.js";
import { a as ge, c as _e, f as ve, i as ye, l as be, m as xe, n as Se, o as Ce, s as we, t as Te, u as Ee, __tla as __tla_1 } from "./bodyRegistry-DeYe45A6.js";
import { c as De, i as Oe, l as ke, n as Ae, o as je, s as Me } from "./bodyHarExport-_jFEZ_en.js";
let y, Er, Wt, Mr, Cr, ii, Nn, gt, Ct, Mn, x, S, fn, ui, V, li, z, Re, dn, jr, ai, O, Gt, _t, _, Ue, Fe, oi, gi, Wr, Xr, Rr, g, zr, yr, Ar, _i, Tr, Ir, Zr, mr, Z, hi, Vr, He, q, dr, $, Lr, ci, Qr, Br, Ur, Vt;
let __tla = Promise.all([
    (()=>{
        try {
            return __tla_0;
        } catch  {}
    })(),
    (()=>{
        try {
            return __tla_1;
        } catch  {}
    })()
]).then(async ()=>{
    g = function(e) {
        return !!e && `items` in e;
    };
    function Ne(e) {
        return g(e) && `isRoot` in e && e.isRoot === !0;
    }
    _ = function(e) {
        return {
            id: `root`,
            title: `HTTP Toolkit Rules`,
            isRoot: !0,
            items: e
        };
    };
    function Pe(e, t) {
        let n = 0;
        for(; e[n] !== void 0 && t[n] !== void 0;){
            let r = t[n] - e[n];
            if (r !== 0) return r;
            n += 1;
        }
        return e[n] === void 0 ? t[n] === void 0 ? 0 : 1 : -1;
    }
    function v(e, t) {
        return y(e, t.slice(0, -1));
    }
    y = function(e, t) {
        return t.reduce((e, n, r)=>{
            if (!g(e)) throw Error(`Invalid path ${t} at step #${r}`);
            return e.items[n];
        }, e);
    };
    Fe = function(e, t, n = []) {
        let r = (e)=>typeof t == `function` ? t(e) : Object.entries(t).every(([t, n])=>e[t] === n);
        if (r(e)) return n;
        for(let i = 0; i < e.items.length; i++){
            let a = e.items[i];
            if (g(a)) {
                let e = Fe(a, t, n.concat(i));
                if (e) return e;
            } else if (r(a)) return n.concat(i);
        }
    };
    function Ie(e, t) {
        let n = v(e, t), r = t[t.length - 1];
        n.items.splice(r, 1), n.items.length === 0 && !Ne(n) && Ie(e, t.slice(0, -1));
    }
    function Le(e) {
        return g(e) ? {
            ...e,
            items: e.items.map((e)=>Le(e)),
            collapsed: !0,
            id: crypto.randomUUID()
        } : {
            ...e,
            matchers: [
                ...e.matchers
            ],
            steps: [
                ...e.steps
            ],
            id: crypto.randomUUID()
        };
    }
    Re = function(e, t, n, r) {
        let i = e[n], a = e[r], o = Pe(i, a), s = v(t, a), c = a[a.length - 1], l = s.items.length > c ? y(t, a) : void 0, u = i.slice(0, -1), d = a.slice(0, -1);
        if (l && g(l) && !l.collapsed && o > 0) return {
            sourcePath: i,
            targetPath: a.concat(0)
        };
        if (l === void 0 && o > 0) {
            let e = d.slice(0, -1), t = JSON.stringify(u) === JSON.stringify(e) ? 0 : 1;
            return {
                sourcePath: i,
                targetPath: a.slice(0, -2).concat(a[a.length - 2] + t)
            };
        }
        if (!ze(u, d)) {
            let e = a[a.length - 1];
            return {
                sourcePath: i,
                targetPath: o < 0 ? d.concat(e) : d.concat(e + 1)
            };
        }
        return {
            sourcePath: i,
            targetPath: a
        };
    };
    function ze(e, t) {
        return e.length === t.length && e.every((e, n)=>e === t[n]);
    }
    function b(e) {
        return JSON.parse(JSON.stringify(e));
    }
    function Be(e) {
        let t = [], n = (e)=>{
            for (let r of e.items)g(r) ? n(r) : t.push(r);
        };
        return n(e), t;
    }
    var Ve = e(t(), 1);
    He = function() {
        return p() && Number.isFinite(ue());
    };
    x = function() {
        if (p()) {
            let e = ue();
            return $.getState().httpProxyPort !== e && $.setState({
                httpProxyPort: e
            }), e;
        }
        let { httpProxyPort: e, isInitialized: t, serverOnline: n } = $.getState();
        return t && n && e ? e : e || 8e3;
    };
    Ue = function() {
        return (0, Ve.useMemo)(()=>x(), [
            $((e)=>e.httpProxyPort),
            $((e)=>e.streamDisconnected),
            $((e)=>e.isInitialized)
        ]);
    };
    function We(e) {
        let t = e.trim().toLowerCase();
        if (!t) return !1;
        let n = t.includes(`:`) ? t.split(`:`)[0] : t;
        return n === `127.0.0.1` || n === `localhost` || n === `::1`;
    }
    S = function(e) {
        return (e ?? []).filter((e)=>e.trim() && !We(e));
    };
    function Ge(e) {
        let t = S(e);
        return t.length > 0 ? t : void 0;
    }
    pe();
    var C = ee();
    function Ke(e) {
        if (e) return {
            proxyUrl: e.proxyUrl,
            noProxy: e.noProxy ? e.noProxy.split(`,`).map((e)=>e.trim()).filter(Boolean) : void 0
        };
    }
    function qe(e) {
        return `docker-tunnel-proxy-${e}`;
    }
    function Je(e, t, n, r) {
        if (e !== `direct`) {
            if (e === `system`) {
                if (!r?.proxyUrl) return;
                try {
                    let e = new URL(r.proxyUrl).hostname;
                    if (e === `localhost` || e.startsWith(`127.0.0`)) return;
                } catch  {
                    return;
                }
                return r;
            }
            if (t) return {
                proxyUrl: `${e}://${t}`,
                noProxy: n.length ? n.join(`,`) : void 0
            };
        }
    }
    function Ye(e) {
        let t = $.getState(), n = Je(e.upstreamProxyType, e.upstreamProxyHost, e.upstreamNoProxyHosts, t.systemProxyConfig), r = qe(x()), i = Ke(n), a = i;
        if (t.ruleParameterKeys.includes(r)) {
            let e = {
                [C.MOCKTTP_PARAM_REF]: r
            };
            a = i ? [
                e,
                i
            ] : e;
        }
        let o = {};
        Object.entries(e.clientCertificateHostMap).forEach(([e, t])=>{
            o[e] = {
                pfx: h.from(t.pfx, `base64`),
                passphrase: t.passphrase
            };
        });
        let s = Ge(t.dnsServers);
        return {
            ignoreHostHttpsErrors: [
                ...e.whitelistedCertificateHosts
            ],
            clientCertificateHostMap: o,
            proxyConfig: a,
            lookupOptions: s ? {
                servers: s
            } : void 0,
            simulateConnectionErrors: !0
        };
    }
    var w = (e)=>e.id.startsWith(`default-`) || e.id === `rule-passthrough-all` || e.id === `rule-passthrough-ws-all`, T = [
        {
            id: `default-group`,
            title: `Default Rules`,
            collapsed: !0,
            items: [
                {
                    id: `default-amiusing`,
                    title: `Diagnostic check: amiusing.httptoolkit.tech`,
                    activated: !0,
                    matchers: [
                        {
                            type: `method`,
                            method: `GET`
                        },
                        {
                            type: `host`,
                            host: `amiusing.httptoolkit.tech`
                        },
                        {
                            type: `path`,
                            path: `/`
                        }
                    ],
                    steps: [
                        {
                            type: `fixed-response`,
                            statusCode: 200,
                            contentType: `text/html`,
                            headers: {
                                "content-type": `text/html`,
                                "cache-control": `no-store`,
                                "httptoolkit-active": `true`
                            },
                            body: `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
  <link rel="icon" type="image/png" href="/favicon.png" />
  <script type="application/json" id="amiusing">
    { "amiusing": true }
  <\/script>
  <title>
    Are you using HTTP Toolkit? Yes!
  </title>
  <style>
    html {
      height: 100%;
    }

    body {
      min-height: 100%;
      box-sizing: border-box;
      margin: 0;
      padding: 8px;

      background-color: #fafafa;
      color: #1e2028;

      font-family: "DM Sans", Arial, sans-serif;
      letter-spacing: -0.5px;
      line-height: 1.3;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
    }

    .content {
      max-width: 600px;
    }

    h1 {
      font-size: 48px;
      letter-spacing: -2px;
    }

    p {
      font-size: 24px;
    }

    .logo {
      display: block;
      margin: 40px auto;
      height: 200px;
      width: 200px;
    }

    @media (prefers-color-scheme: dark) {
      body {
        background-color: #32343B;
        color: #ffffff;
      }
    }

    @media not (prefers-color-scheme: dark) {
      body {
        background-color: #fafafa;
        color: #1e2028;
      }
    }
  </style>
</head>
<body>
  <div class="content">
    <h1>You're being intercepted by HTTP Toolkit</h1>
    <p>
      This response came from HTTP Toolkit, which is currently intercepting this connection.
    </p>
    <p>
      All requests made by this browser will be recorded by HTTP Toolkit.
      Take a look at the 'View' tab there now to see the request & response
      that brought you this page, or start browsing elsewhere to collect more data.
    </p>
  </div>
  <link href="https://fonts.cdnfonts.com/css/dm-sans" rel="stylesheet">
</body>
</html>`
                        }
                    ]
                },
                {
                    id: `default-certificate`,
                    title: `CA Certificate Download: /certificate`,
                    activated: !0,
                    matchers: [
                        {
                            type: `method`,
                            method: `GET`
                        },
                        {
                            type: `host`,
                            host: `amiusing.httptoolkit.tech`
                        },
                        {
                            type: `path`,
                            path: `/certificate`
                        }
                    ],
                    steps: [
                        {
                            type: `from-file`,
                            filePath: ``,
                            contentType: `application/x-x509-ca-cert`,
                            headers: {
                                "content-type": `application/x-x509-ca-cert`
                            }
                        }
                    ]
                },
                {
                    id: `default-android-certificate`,
                    title: `Android App Config & Cert`,
                    activated: !0,
                    matchers: [
                        {
                            type: `method`,
                            method: `GET`
                        },
                        {
                            type: `host`,
                            host: `android.httptoolkit.tech`
                        },
                        {
                            type: `path`,
                            path: `/config`
                        }
                    ],
                    steps: [
                        {
                            type: `fixed-response`,
                            statusCode: 200,
                            contentType: `application/json`,
                            body: JSON.stringify({
                                certificate: ``
                            })
                        }
                    ]
                },
                {
                    id: `rule-passthrough-all`,
                    title: `Passthrough everything else cleanly`,
                    activated: !0,
                    matchers: [
                        {
                            type: `wildcard`
                        }
                    ],
                    steps: [
                        {
                            type: `passthrough`
                        }
                    ]
                },
                {
                    id: `rule-passthrough-ws-all`,
                    title: `Passthrough WebSocket traffic cleanly`,
                    activated: !0,
                    matchers: [
                        {
                            type: `websocket`
                        }
                    ],
                    steps: [
                        {
                            type: `ws-passthrough`
                        }
                    ]
                }
            ]
        }
    ], E = `default-group`, Xe = [
        `rule-passthrough-all`,
        `rule-passthrough-ws-all`
    ];
    function Ze(e) {
        if (!e.activated || w(e)) return !1;
        let t = e.matchers ?? [];
        if (t.length !== 1 || t[0].type !== `wildcard`) return !1;
        let n = e.steps ?? [];
        if (n.length !== 1) return !1;
        let r = n[0];
        if (r.type !== `fixed-response`) return !1;
        let i = typeof r.body == `string` ? r.body.trim() : ``;
        return i === `{}` || i === ``;
    }
    function Qe(e) {
        let t = 0, n = (e)=>e.map((e)=>g(e) ? {
                    ...e,
                    items: n(e.items)
                } : Ze(e) ? (t += 1, {
                    ...e,
                    activated: !1
                }) : e), r = n(e);
        return t > 0 && console.warn(`[rules] Deactivated ${t} global wildcard mock rule(s) that returned "{}" for all traffic. Re-enable on the Rules page only if intentional.`), r;
    }
    function $e(e, t = new Set) {
        for (let n of e)g(n) ? $e(n.items, t) : t.add(n.id);
        return t;
    }
    function D(e) {
        let t = (e)=>e.map((e)=>{
                if (g(e)) return {
                    ...e,
                    items: t(e.items)
                };
                if (w(e)) {
                    let t = T.find((e)=>g(e) && e.id === E)?.items.find((t)=>!(`items` in t) && t.id === e.id);
                    if (t) {
                        let n = e.steps?.[0]?.filePath || ``, r = JSON.parse(JSON.stringify(t.steps));
                        return r[0] && n && e.id === `default-certificate` && (r[0].filePath = n), {
                            ...e,
                            title: t.title,
                            matchers: JSON.parse(JSON.stringify(t.matchers)),
                            steps: r
                        };
                    }
                }
                return e;
            }), n = et(tt(t(e)));
        n = Qe(n);
        let r = $e(n), i = Xe.some((e)=>!r.has(e)), a = n.some((e)=>g(e) && e.id === E);
        if (!i && a) return n;
        let o = T.find((e)=>g(e) && e.id === E);
        return o ? a ? n.map((e)=>{
            if (!g(e) || e.id !== E) return e;
            let t = [
                ...e.items
            ];
            for (let e of o.items)!g(e) && !r.has(e.id) && t.push(JSON.parse(JSON.stringify(e)));
            return {
                ...e,
                items: t
            };
        }) : [
            ...n,
            JSON.parse(JSON.stringify(o))
        ] : n;
    }
    function et(e) {
        let t = [], n = [], r = null, i = (e)=>{
            let n = [];
            for (let r of e.items)g(r) ? n.push(i(r)) : w(r) ? n.push(r) : t.push(r);
            return {
                ...e,
                collapsed: e.collapsed !== !1,
                items: n
            };
        };
        for (let a of e)g(a) && a.id === E ? r = i(a) : g(a) ? n.push(a) : w(a) || t.push(a);
        let a = [
            ...t,
            ...n
        ];
        return r && a.push(r), a;
    }
    function tt(e) {
        return e.map((e)=>g(e) && e.id === E ? {
                ...e,
                collapsed: e.collapsed !== !1
            } : g(e) ? {
                ...e,
                items: tt(e.items)
            } : e);
    }
    var nt = new Set([
        `content-encoding`,
        `transfer-encoding`,
        `date`,
        `expires`
    ]);
    function rt(e, t) {
        if (!e) return;
        let n = {};
        for (let [t, r] of Object.entries(e)){
            let e = t.toLowerCase();
            nt.has(e) || (n[e] = r);
        }
        return t !== void 0 && t >= 0 && (n[`content-length`] = String(t)), Object.keys(n).length > 0 ? n : void 0;
    }
    pe();
    function it(e) {
        switch(e.type){
            case `wildcard`:
            case `default-wildcard`:
                return new C.matchers.WildcardMatcher;
            case `websocket`:
            case `ws-wildcard`:
            case `default-ws-wildcard`:
                return null;
            case `method`:
                {
                    let t = String(e.method ?? `GET`).toUpperCase(), n = C.Method[t] === void 0 ? C.Method.GET : C.Method[t];
                    return new C.matchers.MethodMatcher(n);
                }
            case `path`:
            case `simple-path`:
                return new C.matchers.FlexiblePathMatcher(String(e.path ?? ``));
            case `regex-path`:
                {
                    let t = String(e.regex ?? e.regexSource ?? ``);
                    return t ? new C.matchers.RegexPathMatcher(new RegExp(t)) : null;
                }
            case `regex-url`:
                {
                    let t = String(e.regex ?? e.regexSource ?? ``);
                    return t ? new C.matchers.RegexUrlMatcher(new RegExp(t)) : null;
                }
            case `host`:
                return new C.matchers.HostMatcher(String(e.host ?? e.hostname ?? ``));
            case `hostname`:
                return new C.matchers.HostnameMatcher(String(e.hostname ?? e.host ?? ``));
            case `header`:
                return e.headers && typeof e.headers == `object` ? new C.matchers.HeaderMatcher(e.headers) : e.key ? new C.matchers.HeaderMatcher({
                    [String(e.key)]: String(e.value ?? ``)
                }) : null;
            case `query`:
                return e.query && typeof e.query == `object` ? new C.matchers.QueryMatcher(e.query) : e.name ? new C.matchers.QueryMatcher({
                    [String(e.name)]: String(e.value ?? ``)
                }) : null;
            case `raw-body-includes`:
                return new C.matchers.RawBodyIncludesMatcher(String(e.content ?? e.value ?? ``));
            case `json-body-matching`:
                return new C.matchers.JsonBodyMatcher(e.value ?? e.body);
            default:
                return null;
        }
    }
    function at(e, t) {
        switch(e.type){
            case `fixed-response`:
                {
                    let t = e.body, n = rt(e.headers ?? (e.contentType ? {
                        "content-type": String(e.contentType)
                    } : void 0), t == null ? void 0 : me.byteLength(t, `utf8`));
                    return new C.requestSteps.FixedResponseStep(Number(e.statusCode ?? 200), void 0, t, n);
                }
            case `from-file`:
                return e.filePath ? new C.requestSteps.FileStep(Number(e.statusCode ?? 200), void 0, String(e.filePath), rt(e.headers)) : null;
            case `passthrough`:
                return new C.requestSteps.PassThroughStep({
                    ...t,
                    ...e.connectionOptions
                });
            case `forward-to-host`:
                {
                    if (!e.host) return null;
                    let n = e.protocol;
                    return new C.requestSteps.PassThroughStep({
                        ...t,
                        transformRequest: {
                            replaceHost: {
                                targetHost: String(e.host),
                                updateHostHeader: e.updateHostHeader !== !1
                            },
                            ...n ? {
                                setProtocol: n
                            } : {}
                        }
                    });
                }
            case `req-res-transformer`:
                return new C.requestSteps.PassThroughStep({
                    ...t,
                    transformRequest: e.transformRequest,
                    transformResponse: e.transformResponse
                });
            case `breakpoint`:
                return new C.requestSteps.PassThroughStep(t);
            case `timeout`:
                return new C.requestSteps.TimeoutStep;
            case `delay`:
                return new C.requestSteps.DelayStep(Number(e.delayMs ?? e.delay ?? 0));
            case `close-connection`:
                return new C.requestSteps.CloseConnectionStep;
            case `reset-connection`:
                return new C.requestSteps.ResetConnectionStep;
        }
        return null;
    }
    function ot(e, t) {
        switch(e.type){
            case `ws-passthrough`:
                return new C.webSocketSteps.PassThroughWebSocketStep({
                    ...t,
                    ...e.connectionOptions
                });
            case `ws-echo`:
                return new C.webSocketSteps.EchoWebSocketStep;
            default:
                return null;
        }
    }
    function st(e, t) {
        if (!e.activated || !e.matchers?.length) return null;
        let n = e.matchers.filter((e)=>e.type !== `websocket`).map((e)=>it(e)).filter((e)=>e !== null);
        if (n.length === 0) return null;
        let r = (e.steps ?? []).map((e)=>at(e, t)).filter((e)=>e !== null);
        return r.length === 0 ? null : {
            id: e.id,
            matchers: n,
            steps: r,
            completionChecker: new C.completionCheckers.Always
        };
    }
    function ct(e, t) {
        if (!e.activated || !e.matchers?.length || !e.matchers.some((e)=>e.type === `websocket`)) return null;
        let n = e.matchers.filter((e)=>e.type !== `websocket`).map((e)=>it(e)).filter((e)=>e !== null);
        n.length === 0 && n.push(new C.matchers.WildcardMatcher);
        let r = (e.steps ?? []).map((e)=>ot(e, t)).filter((e)=>e !== null);
        return r.length === 0 ? null : {
            id: e.id,
            matchers: n,
            steps: r,
            completionChecker: new C.completionCheckers.Always
        };
    }
    function lt(e) {
        return {
            httpRules: [
                {
                    id: `rule-passthrough-fallback`,
                    matchers: [
                        new C.matchers.WildcardMatcher
                    ],
                    steps: [
                        new C.requestSteps.PassThroughStep(e)
                    ],
                    completionChecker: new C.completionCheckers.Always
                }
            ],
            wsRules: [
                {
                    id: `rule-passthrough-ws-fallback`,
                    matchers: [
                        new C.matchers.WildcardMatcher
                    ],
                    steps: [
                        new C.webSocketSteps.PassThroughWebSocketStep(e)
                    ],
                    completionChecker: new C.completionCheckers.Always
                }
            ]
        };
    }
    function ut(e, t, n) {
        let r = e.find((e)=>`items` in e && e.id === `default-group`);
        if (!r?.items) return;
        let i = r.items.find((e)=>!(`items` in e) && e.id === `default-certificate`);
        i?.steps?.[0] && t ? i.steps[0].filePath = t : i && (i.activated = !1);
        let a = r.items.find((e)=>!(`items` in e) && e.id === `default-android-certificate`);
        a?.steps?.[0] && n ? a.steps[0].body = JSON.stringify({
            certificate: n
        }) : a && (a.activated = !1);
    }
    function dt(e, t, n, r) {
        let i = D(e), a = JSON.parse(JSON.stringify(i));
        ut(a, n, r);
        let o = Be(_(a)), s = Ye(t), c = [], l = [], u = [];
        for (let e of o)if (e.matchers.some((e)=>e.type === `websocket`)) {
            let t = ct(e, s);
            t ? l.push(t) : u.push(e.id);
        } else {
            let t = st(e, s);
            t ? c.push(t) : u.push(e.id);
        }
        return c.length === 0 && c.push(...lt(s).httpRules), l.length === 0 && l.push(...lt(s).wsRules), {
            httpRules: c,
            wsRules: l,
            skippedRuleIds: u
        };
    }
    o();
    var ft = null, pt = !1, mt = null;
    async function ht(e, t, n) {
        let { certPath: r, certContent: i, serverOnline: o } = $.getState();
        if (!o || !p()) return;
        let { httpRules: s, wsRules: c, skippedRuleIds: l } = dt(e, t, r, i), u, d;
        try {
            await re(s);
        } catch (e) {
            u = e, console.error(`Failed to sync HTTP rules to proxy`, e);
        }
        try {
            await le(c);
        } catch (e) {
            d = e, console.error(`Failed to sync WebSocket rules to proxy`, e);
        }
        if (n?.silent) return;
        let { useNotificationStore: f } = await a(async ()=>{
            let { useNotificationStore: e } = await import(`./useNotificationStore-NSh_n0Ie.js`).then(async (m)=>{
                await m.__tla;
                return m;
            });
            return {
                useNotificationStore: e
            };
        }, __vite__mapDeps([0,1,2,3,4]));
        u || d ? f.getState().addNotification(`error`, `Failed to sync some rules to proxy`) : l.length > 0 ? f.getState().addNotification(`warning`, `Saved rules; skipped ${l.length} invalid rule(s)`) : f.getState().addNotification(`success`, `Rules saved to proxy`);
    }
    gt = function(e, t, n) {
        if (ft) {
            pt = !0, mt = {
                draftRules: e,
                state: t,
                options: n
            };
            return;
        }
        ft = ht(e, t, n).catch((e)=>{
            n?.silent || a(async ()=>{
                let { errorHandler: e } = await import(`./errorHandler-D84ZyjpV.js`).then(async (m)=>{
                    await m.__tla;
                    return m;
                });
                return {
                    errorHandler: e
                };
            }, __vite__mapDeps([5,6,7,3,8,2,4,1])).then(({ errorHandler: t })=>{
                t(e, `Failed to sync rules with server`);
            });
        }).finally(()=>{
            if (ft = null, pt && mt) {
                pt = !1;
                let e = mt;
                mt = null, gt(e.draftRules, e.state, e.options);
            }
        });
    };
    _t = async function() {
        ft && await ft;
    };
    function vt() {
        return null;
    }
    o();
    function yt(e = !1) {
        let t = he.getState().user;
        t && a(async ()=>{
            let { scheduleSync: e } = await import(`./userConfigSync-BIjNXwGj.js`).then(async (m)=>{
                await m.__tla;
                return m;
            });
            return {
                scheduleSync: e
            };
        }, []).then(({ scheduleSync: n })=>n(t, e ? {
                force: !0
            } : void 0));
    }
    O = n()(c((e, t)=>({
            favorites: [],
            isFavorite: (e)=>t().favorites.some((t)=>t.id === e),
            addFavorite: (n)=>{
                let r = t().favorites.find((e)=>e.id === n.id), i = structuredClone({
                    ...n,
                    savedAt: Date.now()
                });
                r ? e((e)=>({
                        favorites: e.favorites.map((e)=>e.id === n.id ? i : e)
                    })) : (e((e)=>({
                        favorites: [
                            i,
                            ...e.favorites
                        ]
                    })), q({
                    type: `success`,
                    domain: `favorites`,
                    message: `Added to favorites`
                })), yt();
            },
            removeFavorite: (t)=>{
                e((e)=>({
                        favorites: e.favorites.filter((e)=>e.id !== t)
                    })), q({
                    type: `info`,
                    domain: `favorites`,
                    message: `Removed from favorites`
                }), a(async ()=>{
                    let { recordSyncDelete: e } = await import(`./userConfigSync-BIjNXwGj.js`).then(async (m)=>{
                        await m.__tla;
                        return m;
                    });
                    return {
                        recordSyncDelete: e
                    };
                }, []).then(({ recordSyncDelete: e })=>e(`favorites`, t)), yt(!0);
            },
            clearFavorites: ()=>{
                let t = O.getState().favorites.map((e)=>e.id);
                e({
                    favorites: []
                }), q({
                    type: `info`,
                    domain: `favorites`,
                    message: `Cleared all favorites`
                }), a(async ()=>{
                    let { recordSyncDeletes: e } = await import(`./userConfigSync-BIjNXwGj.js`).then(async (m)=>{
                        await m.__tla;
                        return m;
                    });
                    return {
                        recordSyncDeletes: e
                    };
                }, []).then(({ recordSyncDeletes: e })=>e(`favorites`, t)), yt(!0);
            }
        }), {
        name: `httptoolkit-favorites`
    }));
    function bt(e) {
        let t = new TextEncoder().encode(e), n = ``;
        for(let e = 0; e < t.length; e++)n += String.fromCharCode(t[e]);
        return btoa(n);
    }
    function xt(e) {
        let t = {};
        return e ? Array.isArray(e) ? (e.forEach((e)=>{
            if (Array.isArray(e)) {
                let [n, r] = e;
                n && (t[n.toLowerCase()] = r ?? ``);
            } else e && typeof e == `object` && e.name && (t[e.name.toLowerCase()] = e.value ?? ``);
        }), t) : (Object.entries(e).forEach(([e, n])=>{
            e && (t[e.toLowerCase()] = String(n));
        }), t) : t;
    }
    var St = new Set([
        `connection`,
        `keep-alive`,
        `proxy-authenticate`,
        `proxy-authorization`,
        `te`,
        `trailers`,
        `transfer-encoding`,
        `upgrade`,
        `host`,
        `content-length`,
        `content-encoding`
    ]);
    Ct = function(e) {
        return e.type === `http` && !!e.method && !!e.url;
    };
    function wt(e) {
        let t = Object.entries(e.requestHeaders || {}).filter(([e])=>!St.has(e.toLowerCase())).map(([e, t])=>({
                key: e,
                value: String(t),
                enabled: !0
            })), n = e.contentType || e.requestHeaders?.[`content-type`] || e.requestHeaders?.[`Content-Type`] || `application/json`, r = e.requestBody ?? ``;
        if (!r) {
            let t = ye(e.id)?.request?.decodedData;
            t && (r = t.toString(`utf8`));
        }
        let i = t.find((e)=>e.key.toLowerCase() === `content-type`);
        return i ? n = i.value.split(`;`)[0]?.trim() || n : n && t.push({
            key: `Content-Type`,
            value: n,
            enabled: !0
        }), {
            method: e.method || `GET`,
            url: e.url || ``,
            headers: t,
            body: r,
            contentType: n.split(`;`)[0]?.trim() || `application/json`
        };
    }
    var Tt = {
        response: 3,
        "websocket-accepted": 2,
        abort: 1
    }, k = {};
    function Et() {
        Object.keys(k).forEach((e)=>delete k[e]);
    }
    function Dt(e) {
        delete k[e];
    }
    function Ot(e, t, n) {
        let r = Tt[t] ?? 0, i = k[e] ?? [];
        i.push({
            type: t,
            event: n,
            priority: r
        }), i.sort((e, t)=>t.priority - e.priority), k[e] = i;
    }
    function kt(e) {
        let t = k[e];
        return t?.length ? (delete k[e], [
            ...t
        ].sort((e, t)=>t.priority - e.priority)) : [];
    }
    var At = new Set([
        `set-cookie`
    ]);
    function jt(e) {
        return At.has(e.toLowerCase()) ? `
` : `, `;
    }
    function Mt(e, t) {
        if (Array.isArray(e)) {
            let n = t ? jt(t) : `, `;
            return e.map((e)=>Mt(e, t)).join(n);
        }
        if (e && typeof e == `object`) {
            let t = e;
            return `value` in t ? Mt(t.value) : JSON.stringify(e);
        }
        return String(e ?? ``);
    }
    function A(e, t, n) {
        let r = t.toLowerCase(), i = Mt(n, r);
        if (i) if (e[r]) {
            let t = jt(r);
            e[r] = `${e[r]}${t}${i}`;
        } else e[r] = i;
    }
    function Nt(e) {
        let t = e.split(`,`).map((e)=>e.trim()).filter(Boolean);
        if (t.length <= 1) return e;
        let n = [
            ...t
        ];
        for(; n.length % 2 == 0;){
            let e = n.length / 2;
            if (n.slice(0, e).join(`,`) !== n.slice(e).join(`,`)) break;
            n = n.slice(0, e);
        }
        let r = [];
        for (let e of n)(r.length === 0 || r[r.length - 1] !== e) && r.push(e);
        return r.join(`, `);
    }
    function Pt(e, t) {
        let n = e.toLowerCase();
        return [
            `content-encoding`,
            `accept-encoding`,
            `transfer-encoding`
        ].includes(n) ? Nt(t) : t;
    }
    function Ft(e, t, n) {
        let r = e.trim(), i = t.trim();
        if (!r) return i;
        if (!i || r === i) return r;
        if (At.has(n.toLowerCase())) return `${r}\n${i}`;
        let a = (e)=>e.split(`,`).map((e)=>e.trim()).filter(Boolean), o = new Set, s = [];
        for (let e of [
            ...a(r),
            ...a(i)
        ]){
            let t = e.toLowerCase();
            o.has(t) || (o.add(t), s.push(e));
        }
        return s.join(`, `);
    }
    function j(e) {
        if (!e) return {};
        let t = {};
        if (typeof e == `string`) {
            let n = e.trim();
            if (!n) return {};
            try {
                return j(JSON.parse(n));
            } catch  {
                for (let e of n.split(/\r?\n/)){
                    let n = e.indexOf(`:`);
                    n > 0 && A(t, e.slice(0, n).trim(), e.slice(n + 1).trim());
                }
                return t;
            }
        }
        if (e && typeof e == `object` && !Array.isArray(e)) {
            let t = e;
            if (t.headers !== void 0) return j(t.headers);
            if (t.rawHeaders !== void 0) return j(t.rawHeaders);
        }
        if (Array.isArray(e)) {
            if (e.length > 1 && e.every((e)=>typeof e == `string`)) {
                for(let n = 0; n < e.length - 1; n += 2)A(t, String(e[n] ?? ``), e[n + 1]);
                return t;
            }
            for (let n of e)if (Array.isArray(n) && n.length >= 2) A(t, String(n[0] ?? ``), n[1]);
            else if (typeof n == `string` && n.includes(`:`)) {
                let [e, ...r] = n.split(`:`);
                A(t, e, r.join(`:`).trim());
            } else if (n && typeof n == `object`) {
                let e = n;
                typeof e.name == `string` ? A(t, e.name, e.value) : typeof e.key == `string` && A(t, e.key, e.value);
            }
            return t;
        }
        if (typeof e == `object`) {
            let n = Object.entries(e);
            if (n.length > 1 && n.every(([e, t])=>/^\d+$/.test(e) && typeof t == `string`)) return j(n.sort((e, t)=>Number(e[0]) - Number(t[0])).map(([, e])=>String(e)));
            n.forEach(([e, n])=>A(t, e, n));
            for (let [e, n] of Object.entries(t))t[e] = Pt(e, n);
            return t;
        }
        return t;
    }
    function M(e, t) {
        let n = j(e), r = j(t);
        if (Object.keys(r).length === 0) return n;
        if (Object.keys(n).length === 0) {
            let e = {
                ...r
            };
            for (let [t, n] of Object.entries(e))e[t] = Pt(t, n);
            return e;
        }
        let i = {
            ...n
        };
        for (let [e, t] of Object.entries(r)){
            let n = e.toLowerCase(), r = i[n] ?? i[e];
            if (!r) {
                i[n] = t;
                continue;
            }
            i[n] = Ft(r, t, n);
        }
        for (let [e, t] of Object.entries(i))i[e] = Pt(e, t);
        return i;
    }
    function It(e) {
        let t = String(e.url ?? ``).trim();
        if (t) return t;
        let n = e.destination, r = String(n?.hostname ?? e.hostname ?? ``).trim();
        if (!r) return t;
        let i = n?.port ?? (e.port == null ? void 0 : Number(e.port)), a = String(e.path ?? `/`), o = a.startsWith(`/`) ? a : `/${a}`, s = String(e.protocol ?? `https`).replace(/:$/, ``) || `https`, c = s === `http` ? 80 : 443;
        return `${s}://${r}${i != null && !Number.isNaN(i) && i !== c ? `:${i}` : ``}${o}`;
    }
    var Lt = 8e3;
    function Rt(e) {
        return e.length <= 8e3 ? e : e.slice(0, Lt);
    }
    function N(e, t) {
        return Rt([
            t,
            ...e
        ]);
    }
    function P(e, t, n) {
        let r = e.findIndex((e)=>e.id === t);
        if (r === -1) return null;
        let i = [
            ...e
        ];
        return i[r] = n(i[r]), i;
    }
    let zt, Bt, Ht;
    zt = [
        `vercel.app`,
        `vercel-dns.com`,
        `googleapis.com`,
        `www.googleapis.com`,
        `www.gstatic.com`,
        `gstatic.com`,
        `accounts.google.com`
    ];
    Bt = [
        `tauri.localhost`
    ];
    Vt = [
        `localhost`,
        `127.0.0.1`,
        `::1`,
        `0.0.0.0`
    ];
    Ht = new Set(Vt);
    function Ut() {
        let e = [
            d(),
            f(),
            5173
        ];
        return p() && e.push(ue()), [
            ...new Set(e)
        ];
    }
    Wt = function() {
        return [
            ...Vt
        ];
    };
    Gt = function() {
        return [
            ...Bt,
            ...zt.map((e)=>`*.${e}`)
        ];
    };
    function Kt(e) {
        let t = new Set, n = [];
        for (let r of [
            ...Wt(),
            ...e ?? []
        ]){
            let e = r.trim(), i = e.toLowerCase();
            !i || t.has(i) || (t.add(i), n.push(e));
        }
        return n;
    }
    function qt(e) {
        let t = e.trim().toLowerCase();
        return t.startsWith(`[`) && t.endsWith(`]`) ? t.slice(1, -1) : t;
    }
    function Jt(e) {
        let t = qt(e);
        return t ? Bt.some((e)=>t === e.toLowerCase()) ? !0 : zt.some((e)=>t === e || t.endsWith(`.${e}`)) : !1;
    }
    function Yt(e, t) {
        let n = qt(e);
        return !Ht.has(n) || t === void 0 ? !1 : Ut().includes(t);
    }
    function Xt(e) {
        let t = e.trim();
        if (!t) return null;
        try {
            let e = t.includes(`://`) ? t : `https://${t}`, n = new URL(e), r = n.port ? Number(n.port) : void 0;
            return {
                hostname: n.hostname,
                port: r
            };
        } catch  {
            let e = t.split(`/`)[0] ?? t, n = e.lastIndexOf(`:`);
            if (n > 0 && !e.includes(`]`)) {
                let t = Number(e.slice(n + 1));
                if (!Number.isNaN(t)) return {
                    hostname: e.slice(0, n),
                    port: t
                };
            }
            return {
                hostname: e
            };
        }
    }
    function Zt(e) {
        try {
            let t = s.getState();
            if (!t.bypassTelemetryEnabled) return !1;
            let n = e.trim().toLowerCase();
            return t.bypassTelemetryDomains.some((e)=>{
                let t = e.trim().toLowerCase();
                return t.startsWith(`*.`) && (t = t.slice(2)), n === t || n.endsWith(`.${t}`);
            });
        } catch  {
            return !1;
        }
    }
    function F(e) {
        let t = Xt(e);
        return t ? !!(Jt(t.hostname) || Yt(t.hostname, t.port) || Zt(t.hostname)) : !1;
    }
    function Qt(e, t) {
        return !!(Jt(e) || Yt(e, t) || Zt(e));
    }
    var $t = [
        `request-initiated`,
        `request`,
        `response`,
        `abort`,
        `tls-client-error`,
        `tls-passthrough-opened`,
        `tls-passthrough-closed`,
        `client-error`,
        `raw-passthrough-opened`,
        `raw-passthrough-closed`,
        `rule-event`,
        `websocket-request`,
        `websocket-accepted`,
        `websocket-message-received`,
        `websocket-message-sent`,
        `websocket-close`
    ];
    function en(e, t, n, r, i) {
        e((e)=>{
            let a = P(e.events, t, (e)=>Me(e, n, r, i));
            return a ? {
                events: a
            } : (Ot(t, `response`, n), e);
        });
    }
    function tn(e, t, n, r) {
        (async ()=>{
            try {
                let i = Se(n)?.byteLength ?? 0, a = await Ce(n, r);
                Ee(t, a);
                let o = a.encodedByteLength || i;
                o !== i && e((e)=>{
                    let n = P(e.events, t, (e)=>({
                            ...e,
                            responseBodySize: o,
                            size: ke(o)
                        }));
                    return n ? {
                        events: n
                    } : e;
                });
            } catch (e) {
                console.error(`[EventStream] Failed to ingest response body`, e);
            }
        })();
    }
    function nn(e, t) {
        let n = kt(e);
        if (!n.length) return;
        let r = t.setState, i = t.getState;
        for (let t of n){
            let n = t.event;
            if (t.type === `response`) {
                let t = i().events.find((t)=>t.id === e);
                if (F(String(n.url || t?.url || ``))) {
                    _e(e), r((t)=>({
                            events: t.events.filter((t)=>t.id !== e)
                        }));
                    continue;
                }
                let a = M(n.headers, n.rawHeaders);
                en(r, e, n, a, Se(n)?.byteLength ?? 0), tn(r, e, n, a);
            } else t.type === `abort` ? r((t)=>{
                let r = P(t.events, e, (e)=>je(e, n));
                return r ? {
                    events: r
                } : t;
            }) : t.type === `websocket-accepted` && r((t)=>{
                let r = P(t.events, e, (e)=>({
                        ...e,
                        responseHeaders: M(n.headers, n.rawHeaders),
                        isCompleted: !0
                    }));
                return r ? {
                    events: r
                } : t;
            });
        }
    }
    function rn(e) {
        let t = e.getState, n = e.setState, r = (e)=>{
            _e(e), Dt(e), n((t)=>({
                    events: t.events.filter((t)=>t.id !== e)
                }));
        }, i = (e)=>{
            n((t)=>t.events.some((t)=>t.id === e.id) ? t : {
                    events: N(t.events, e),
                    droppedEventCount: t.events.length >= 8e3 ? t.droppedEventCount + 1 : t.droppedEventCount
                });
        }, a = (n)=>{
            if (t().isPaused) return;
            let r = It(n);
            if (F(r)) return;
            let a = String(n.id), o = M(n.headers, n.rawHeaders), s = {
                id: a,
                method: String(n.method || `GET`),
                url: r,
                type: `http`,
                timestamp: n.timingEvents?.startTime || Date.now(),
                duration: 0,
                size: `...`,
                requestHeaders: o,
                responseHeaders: {},
                requestBodySize: 0,
                isCompleted: !1
            };
            s.category = De(s), ge(a, o), i(s), nn(a, e);
        }, o = (i)=>{
            if (t().isPaused) return;
            let a = It(i), o = String(i.id);
            if (F(a)) {
                r(o);
                return;
            }
            let s = M(i.headers, i.rawHeaders);
            (async ()=>{
                let t = await we(o, i, s), r = t.encodedByteLength;
                t.waitForDecoding().catch(()=>{}), n((e)=>{
                    let t = e.events.find((e)=>e.id === o), n = t ?? {
                        id: o,
                        method: String(i.method || `GET`),
                        url: a,
                        type: `http`,
                        timestamp: i.timingEvents?.startTime || Date.now(),
                        duration: 0,
                        size: `...`,
                        requestHeaders: s,
                        responseHeaders: {},
                        isCompleted: !1
                    }, c = {
                        ...n,
                        method: String(i.method || n.method),
                        url: a || n.url,
                        requestHeaders: Object.keys(s).length ? s : n.requestHeaders,
                        requestBodySize: r,
                        httpVersion: i.httpVersion ? String(i.httpVersion) : n.httpVersion
                    };
                    if (c.category = De(c), t) {
                        let t = P(e.events, o, ()=>c);
                        return t ? {
                            events: t
                        } : e;
                    }
                    let l = e.events.length >= 8e3 ? e.droppedEventCount + 1 : e.droppedEventCount;
                    return {
                        events: N(e.events, c),
                        droppedEventCount: l
                    };
                }), nn(o, e);
            })();
        }, s = (e)=>{
            if (t().isPaused) return;
            let i = String(e.id), a = t().events.find((e)=>e.id === i);
            if (F(It(e) || a?.url || ``)) {
                a && r(i), Dt(i);
                return;
            }
            let o = M(e.headers, e.rawHeaders), s = Se(e)?.byteLength ?? 0;
            if (!a) {
                Ot(i, `response`, e);
                return;
            }
            en(n, i, e, o, s), tn(n, i, e, o);
        }, c = (e)=>{
            if (t().isPaused) return;
            let i = String(e.id), a = t().events.find((e)=>e.id === i);
            if (a && F(a.url)) {
                r(i);
                return;
            }
            n((t)=>{
                let n = P(t.events, i, (t)=>je(t, e));
                return n ? {
                    events: n
                } : (Ot(i, `abort`, e), t);
            });
        }, l = (e)=>{
            if (t().isPaused) return;
            let n = String(e.hostname || ``);
            Qt(n) || i({
                id: String(e.id || `tls-${Date.now()}`),
                method: `CONNECT`,
                url: n || `unknown`,
                type: `tls-failure`,
                category: `aborted`,
                timestamp: e.timingEvents?.startTime || Date.now(),
                duration: 0,
                size: `—`,
                requestHeaders: {},
                responseHeaders: {},
                isCompleted: !0,
                isAborted: !0
            });
        }, u = (e)=>{
            if (t().isPaused) return;
            let n = String(e.hostname || e.host || ``);
            n && Qt(n) || i({
                id: String(e.id || `tunnel-${Date.now()}`),
                method: `CONNECT`,
                url: n || String(e.destination || `TLS tunnel`),
                type: `tunnel`,
                category: `tunnel`,
                timestamp: Date.now(),
                duration: 0,
                size: `—`,
                requestHeaders: {},
                responseHeaders: {},
                isCompleted: !1
            });
        }, d = (e)=>{
            let t = String(e.id);
            n((n)=>{
                let r = P(n.events, t, (t)=>({
                        ...t,
                        isCompleted: !0,
                        duration: e.duration || t.duration
                    }));
                return r ? {
                    events: r
                } : n;
            });
        }, f = (e)=>{
            if (t().isPaused) return;
            let n = String(e.id || `client-error-${Date.now()}`), r = String(e.error?.message || e.message || `Client error`);
            i({
                id: n,
                method: `ERR`,
                url: r,
                type: `http`,
                category: `aborted`,
                timestamp: Date.now(),
                duration: 0,
                size: `—`,
                requestHeaders: {},
                responseHeaders: {},
                isCompleted: !1,
                isAborted: !0,
                statusText: r
            });
        }, ee = (e)=>{
            if (t().isPaused) return;
            let n = String(e.hostname || e.host || e.destination || ``);
            n && Qt(n) || i({
                id: String(e.id || `raw-${Date.now()}`),
                method: `TCP`,
                url: String(e.destination || e.host || `Raw passthrough`),
                type: `tunnel`,
                category: `tunnel`,
                timestamp: Date.now(),
                duration: 0,
                size: `—`,
                requestHeaders: {},
                responseHeaders: {},
                isCompleted: !1
            });
        }, te = (e)=>{
            let t = String(e.id);
            n((e)=>{
                let n = P(e.events, t, (e)=>({
                        ...e,
                        isCompleted: !0
                    }));
                return n ? {
                    events: n
                } : e;
            });
        }, ne = (e)=>{
            let t = String(e.requestId || e.id || ``);
            if (!t) return;
            let r = String(e.ruleId || ``);
            n((e)=>{
                let n = P(e.events, t, (e)=>({
                        ...e,
                        matchedRuleId: r || e.matchedRuleId
                    }));
                return n ? {
                    events: n
                } : e;
            });
        }, re = (e)=>{
            if (t().isPaused) return;
            let n = String(e.url || ``);
            if (F(n)) return;
            let r = String(e.id), a = M(e.headers, e.rawHeaders);
            i({
                id: r,
                method: `GET`,
                url: n,
                statusCode: 101,
                statusText: `Switching Protocols`,
                type: `websocket`,
                category: `websocket`,
                timestamp: e.timingEvents?.startTime || Date.now(),
                duration: 0,
                size: `0 msgs`,
                requestHeaders: a,
                responseHeaders: {},
                isCompleted: !1
            });
        }, ie = (e)=>{
            if (t().isPaused) return;
            let i = String(e.id), a = t().events.find((e)=>e.id === i);
            if (a && F(a.url)) {
                r(i);
                return;
            }
            n((t)=>{
                let n = P(t.events, i, (t)=>({
                        ...t,
                        responseHeaders: M(e.headers, e.rawHeaders),
                        isCompleted: !0
                    }));
                return n || Ot(i, `websocket-accepted`, e), n ? {
                    events: n
                } : t;
            });
        }, ae = (e, n)=>{
            if (t().isPaused) return;
            let r = String(e.streamId), i = t().events.find((e)=>e.id === r);
            i && F(i.url) || t().addWebSocketMessage(r, {
                direction: n,
                timestamp: e.timingEvents?.startTime || Date.now(),
                content: String(e.content ?? ``),
                type: e.isBinary ? `binary` : `text`
            });
        };
        return {
            "request-initiated": a,
            request: o,
            response: s,
            abort: c,
            "tls-client-error": l,
            "tls-passthrough-opened": u,
            "tls-passthrough-closed": d,
            "client-error": f,
            "raw-passthrough-opened": ee,
            "raw-passthrough-closed": te,
            "rule-event": ne,
            "websocket-request": re,
            "websocket-accepted": ie,
            "websocket-message-received": (e)=>ae(e, `received`),
            "websocket-message-sent": (e)=>ae(e, `sent`),
            "websocket-close": (e)=>{
                if (t().isPaused) return;
                let r = String(e.streamId || e.id);
                n((t)=>{
                    let n = P(t.events, r, (t)=>({
                            ...t,
                            isCompleted: !0,
                            statusText: String(e.closeCode ?? `closed`)
                        }));
                    return n ? {
                        events: n
                    } : t;
                });
            }
        };
    }
    o();
    var an = [
        `peer-connected`,
        `peer-disconnected`,
        `external-peer-attached`,
        `data-channel-opened`,
        `data-channel-message-sent`,
        `data-channel-message-received`,
        `data-channel-closed`,
        `media-track-opened`,
        `media-track-stats`,
        `media-track-closed`
    ], on = 0;
    function sn(e, t, n) {
        on = t, a(async ()=>{
            let { onMockRTCEvent: e } = await import(`./proxyAdminClient-D-l-Dsud.js`).then(async (m)=>{
                await m.__tla;
                return m;
            });
            return {
                onMockRTCEvent: e
            };
        }, __vite__mapDeps([9,10,3,11,4,12,13,14,8,2,15])).then(({ onMockRTCEvent: r })=>{
            if (t !== n()) return;
            let i = {
                "peer-connected": (t)=>{
                    let n = e;
                    if (n.getState().isPaused) return;
                    let r = String(t.id || `rtc-${Date.now()}`);
                    n.setState((e)=>e.events.some((e)=>e.id === r) ? e : {
                            events: N(e.events, {
                                id: r,
                                method: `RTC`,
                                url: String(t.sessionId || `WebRTC peer`),
                                type: `webrtc`,
                                category: `rtc-data`,
                                timestamp: Date.now(),
                                duration: 0,
                                size: `—`,
                                requestHeaders: {},
                                responseHeaders: {},
                                isCompleted: !1
                            }),
                            droppedEventCount: e.events.length >= 8e3 ? e.droppedEventCount + 1 : e.droppedEventCount
                        });
                },
                "peer-disconnected": (t)=>{
                    let n = String(t.id);
                    e.setState((e)=>{
                        let t = e.events.map((e)=>e.id === n ? {
                                ...e,
                                isCompleted: !0,
                                size: `closed`
                            } : e);
                        return e.events.some((e)=>e.id === n) ? {
                            events: t
                        } : e;
                    });
                },
                "external-peer-attached": ()=>{},
                "data-channel-opened": ()=>{},
                "data-channel-message-sent": ()=>{},
                "data-channel-message-received": ()=>{},
                "data-channel-closed": ()=>{},
                "media-track-opened": ()=>{},
                "media-track-stats": ()=>{},
                "media-track-closed": ()=>{}
            };
            for (let e of an){
                let a = i[e];
                a && r(e, (r)=>{
                    if (!(t !== n() || on !== t)) try {
                        a(r);
                    } catch (t) {
                        console.error(`[EventStream] RTC error handling ${e}`, t);
                    }
                }).catch((t)=>{
                    console.warn(`[EventStream] RTC subscribe skipped for ${e}`, t);
                });
            }
        });
    }
    function cn() {
        on += 1;
    }
    var I = 0;
    async function ln(e) {
        I += 1;
        let t = I;
        Et();
        let n = rn(e), r = await Promise.allSettled($t.map(async (e)=>{
            let r = n[e];
            r && await oe(e, (n)=>{
                if (t === I) try {
                    r(n);
                } catch (t) {
                    console.error(`[EventStream] Error handling ${e}`, t);
                }
            });
        })), i = r.filter((e)=>e.status === `rejected`);
        for (let e of i){
            let t = e.reason instanceof Error ? e.reason.message : String(e.reason ?? `unknown`);
            console.warn(`[EventStream] Subscription failed:`, t);
        }
        if (r.length - i.length === 0) throw Error(`Could not subscribe to mockttp traffic events (admin stream unavailable)`);
        sn(e, t, ()=>I);
    }
    function un() {
        I += 1, Et(), cn();
    }
    o();
    dn = n((e, t)=>({
            sessions: [],
            activeSessionId: void 0,
            openHar: async (t, n)=>{
                try {
                    let { parseHar: r } = await a(async ()=>{
                        let { parseHar: e } = await import(`./har-BulMv3VR.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        });
                        return {
                            parseHar: e
                        };
                    }, __vite__mapDeps([16,3,17,7,8,2,4,11,14,18,13,19])), { events: i } = await r(t), o = i.filter((e)=>e.isCompleted).length, c = i.length - o, u = `har-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
                    e((e)=>({
                            sessions: [
                                ...e.sessions,
                                {
                                    id: u,
                                    label: n.label,
                                    sourcePath: n.sourcePath,
                                    events: Rt(i),
                                    importedAt: Date.now()
                                }
                            ],
                            activeSessionId: u
                        })), s.getState().setActiveView(`har`);
                    let d = c > 0 ? `Opened ${n.label}: ${i.length} entries (${o} complete, ${c} incomplete)` : `Opened ${n.label}: ${i.length} entries`;
                    l.getState().addNotification(`success`, d);
                } catch (e) {
                    u(e, `Failed to parse HAR file`);
                }
            },
            closeSession: (t)=>{
                e((e)=>{
                    let n = e.sessions.filter((e)=>e.id !== t), r = e.activeSessionId;
                    return r === t && (r = n[n.length - 1]?.id), n.length === 0 && s.getState().activeView === `har` && s.getState().setActiveView(`view`), {
                        sessions: n,
                        activeSessionId: r
                    };
                });
            },
            setActiveSession: (n)=>{
                t().sessions.some((e)=>e.id === n) && e({
                    activeSessionId: n
                });
            },
            selectEvent: (t, n)=>{
                e((e)=>({
                        sessions: e.sessions.map((e)=>e.id === t ? {
                                ...e,
                                activeEventId: n
                            } : e)
                    }));
            }
        }));
    fn = async function(e) {
        let t = await e.text(), n = JSON.parse(t);
        await dn.getState().openHar(n, {
            label: e.name
        });
    };
    var pn = `httptoolkit-pro`, mn = 1, L = `traffic-snapshot`, hn = `har-sessions`, gn = 1440 * 60 * 1e3, _n = 256 * 1024, vn = 2e3;
    function yn() {
        return new Promise((e, t)=>{
            let n = indexedDB.open(pn, mn);
            n.onerror = ()=>t(n.error), n.onsuccess = ()=>e(n.result), n.onupgradeneeded = ()=>{
                let e = n.result;
                e.objectStoreNames.contains(L) || e.createObjectStore(L), e.objectStoreNames.contains(hn) || e.createObjectStore(hn);
            };
        });
    }
    function bn(e) {
        return e.map((e)=>{
            let t = {
                ...e
            };
            return t.requestBody && t.requestBody.length > _n && (t.requestBody = void 0), t.responseBody && t.responseBody.length > _n && (t.responseBody = void 0), t;
        });
    }
    async function xn(e, t) {
        let n = await yn();
        return new Promise((r, i)=>{
            let a = n.transaction(e, `readonly`).objectStore(e).get(t);
            a.onsuccess = ()=>r(a.result), a.onerror = ()=>i(a.error);
        });
    }
    async function Sn(e, t, n) {
        let r = await yn();
        return new Promise((i, a)=>{
            let o = r.transaction(e, `readwrite`).objectStore(e).put(n, t);
            o.onsuccess = ()=>i(), o.onerror = ()=>a(o.error);
        });
    }
    async function Cn(e, t) {
        let n = await yn();
        return new Promise((r, i)=>{
            let a = n.transaction(e, `readwrite`).objectStore(e).delete(t);
            a.onsuccess = ()=>r(), a.onerror = ()=>i(a.error);
        });
    }
    async function wn(e) {
        s.getState().persistTrafficOnRefresh && e.events.length !== 0 && await Sn(L, `latest`, {
            ...e,
            events: bn(e.events),
            savedAt: Date.now()
        });
    }
    async function Tn() {
        if (!s.getState().persistTrafficOnRefresh) return null;
        let e = await xn(L, `latest`);
        return e ? Date.now() - e.savedAt > gn ? (await En(), null) : e : null;
    }
    async function En() {
        await Cn(L, `latest`).catch(()=>{});
    }
    async function Dn(e) {
        s.getState().persistTrafficOnRefresh && await Sn(hn, `latest`, {
            sessions: e,
            savedAt: Date.now()
        });
    }
    async function On() {
        if (!s.getState().persistTrafficOnRefresh) return null;
        let e = await xn(hn, `latest`);
        return e ? Date.now() - e.savedAt > gn ? (await Cn(hn, `latest`).catch(()=>{}), null) : e.sessions : null;
    }
    var R = null;
    function kn() {
        let { events: e, websocketMessages: t, droppedEventCount: n } = z.getState();
        return {
            events: e,
            websocketMessages: t,
            droppedEventCount: n
        };
    }
    function An() {
        s.getState().persistTrafficOnRefresh && (R && clearTimeout(R), R = setTimeout(()=>{
            R = null, wn(kn()).catch(()=>{});
        }, vn));
    }
    function jn() {
        s.getState().persistTrafficOnRefresh && wn(kn()).catch(()=>{});
    }
    Mn = function() {
        let e = z.subscribe(()=>{
            An();
        }), t = dn.subscribe((e)=>{
            Dn(e.sessions).catch(()=>{});
        }), n = ()=>jn();
        return window.addEventListener(`beforeunload`, n), ()=>{
            e(), t(), window.removeEventListener(`beforeunload`, n), R && clearTimeout(R);
        };
    };
    Nn = async function() {
        let e = await Tn();
        e && e.events.length > 0 && z.setState({
            events: e.events,
            websocketMessages: e.websocketMessages,
            droppedEventCount: e.droppedEventCount
        });
        let t = await On();
        t && t.length > 0 && dn.setState({
            sessions: t,
            activeSessionId: t[t.length - 1]?.id
        });
    };
    o();
    let Pn;
    Pn = null;
    z = n((e, t)=>({
            events: [],
            websocketMessages: {},
            droppedEventCount: 0,
            isPaused: !1,
            activeSearch: ``,
            streamConnected: !1,
            reconnecting: !1,
            reconnectAttempt: 0,
            connectToEventStream: ()=>{
                if (!p()) {
                    e({
                        streamConnected: !1
                    });
                    return;
                }
                t().streamConnected || Pn || (Pn = ln(z).then(()=>{
                    x(), $.getState().setStreamDisconnected(!1), e({
                        streamConnected: !0,
                        reconnecting: !1,
                        reconnectAttempt: 0
                    }), console.log(`[EventStream] PluggableAdmin subscriptions active`);
                }).catch((t)=>{
                    u(t, `Failed to subscribe to mockttp events`), e({
                        streamConnected: !1,
                        reconnecting: !1
                    });
                }).finally(()=>{
                    Pn = null;
                }));
            },
            disconnectStream: ()=>{
                un(), Pn = null, e({
                    streamConnected: !1,
                    reconnecting: !1,
                    reconnectAttempt: 0
                });
            },
            togglePause: ()=>e((e)=>({
                        isPaused: !e.isPaused
                    })),
            clearEvents: ()=>{
                Te(), En(), e({
                    events: [],
                    websocketMessages: {},
                    droppedEventCount: 0
                });
            },
            addEvent: (t)=>e((e)=>e.isPaused ? e : {
                        events: N(e.events, t),
                        droppedEventCount: e.events.length >= 8e3 ? e.droppedEventCount + 1 : e.droppedEventCount
                    }),
            prependSentEvent: (t)=>e((e)=>({
                        events: N(e.events, t),
                        droppedEventCount: e.events.length >= 8e3 ? e.droppedEventCount + 1 : e.droppedEventCount
                    })),
            updateEvent: (t, n)=>e((e)=>{
                    let r = P(e.events, t, (e)=>({
                            ...e,
                            ...n
                        }));
                    return r ? {
                        events: r
                    } : e;
                }),
            togglePinEvent: (t)=>e((e)=>{
                    let n = P(e.events, t, (e)=>({
                            ...e,
                            pinned: !e.pinned
                        }));
                    return n ? {
                        events: n
                    } : e;
                }),
            addWebSocketMessage: (t, n)=>e((e)=>{
                    let r = [
                        ...e.websocketMessages[t] || [],
                        n
                    ];
                    r.length > 200 && (r = r.slice(-200));
                    let i = P(e.events, t, (e)=>({
                            ...e,
                            size: `${r.length} msgs`
                        }));
                    return {
                        websocketMessages: {
                            ...e.websocketMessages,
                            [t]: r
                        },
                        ...i ? {
                            events: i
                        } : {}
                    };
                }),
            importHar: (t)=>{
                a(async ()=>{
                    let { parseHar: e } = await import(`./har-BulMv3VR.js`).then(async (m)=>{
                        await m.__tla;
                        return m;
                    });
                    return {
                        parseHar: e
                    };
                }, __vite__mapDeps([16,3,17,7,8,2,4,11,14,18,13,19])).then(({ parseHar: e })=>e(t)).then(({ events: t })=>{
                    let n = t.filter((e)=>e.isCompleted).length, r = t.length - n;
                    e({
                        events: Rt(t),
                        droppedEventCount: 0,
                        websocketMessages: {}
                    });
                    let i = r > 0 ? `Imported ${t.length} entries (${n} complete, ${r} incomplete — no response saved in HAR)` : `Imported ${t.length} entries`;
                    l.getState().addNotification(`success`, i);
                }).catch((e)=>{
                    u(e, `Failed to parse HAR file`);
                });
            },
            exportHar: async (e = `httptoolkit-export`)=>{
                let { events: n } = t();
                try {
                    let t = await Ae(n), r = t._meta?.skipped ?? 0, i = t.log.entries.length;
                    Oe(t, `${e}.har`);
                    let a = r > 0 ? `Exported ${i} completed exchanges (${r} pending/aborted skipped)` : `Exported ${i} exchanges`;
                    l.getState().addNotification(`success`, a);
                } catch (e) {
                    u(e, `Failed to export HAR file`);
                }
            }
        }));
    pe();
    function Fn(e) {
        let t = {};
        for (let n of e.headers)!n.enabled || !n.key.trim() || (t[n.key.toLowerCase()] = n.value);
        let n = t[`content-type`] || e.contentType;
        return n && !t[`content-type`] && (t[`content-type`] = n), t;
    }
    function In(e, t, n) {
        let r = Fn(t), i = {
            id: e,
            method: t.method,
            url: n,
            type: `http`,
            timestamp: Date.now(),
            duration: 0,
            size: `...`,
            requestHeaders: r,
            responseHeaders: {},
            isCompleted: !1,
            contentType: (r[`content-type`] || t.contentType || ``).split(`;`)[0]
        };
        i.category = De(i), ![
            `GET`,
            `HEAD`,
            `OPTIONS`
        ].includes(t.method.toUpperCase()) && t.body && be(e, new ve({
            body: h.from(t.body, `utf8`)
        }, r)), z.getState().prependSentEvent(i);
    }
    function Ln(e, t) {
        let n = {
            ...t.headers
        };
        Ee(e, new ve({
            body: t.body
        }, n));
        let r = n[`content-type`]?.split(`;`)[0] || ``, i = {
            statusCode: t.statusCode,
            statusText: t.statusMessage,
            responseHeaders: n,
            responseBodySize: t.body.byteLength,
            duration: t.duration,
            size: ke(t.body.byteLength),
            isCompleted: !0,
            isAborted: !1,
            contentType: r
        };
        i.category = De({
            ...z.getState().events.find((t)=>t.id === e),
            ...i
        }), z.getState().updateEvent(e, i);
    }
    function Rn(e, t) {
        z.getState().updateEvent(e, {
            isAborted: !0,
            isCompleted: !1,
            statusText: t,
            category: `aborted`
        });
    }
    function zn(e, t, n, r) {
        z.getState().events.find((t)=>t.id === e) || In(e, t, n), Ln(e, {
            statusCode: r.statusCode,
            statusMessage: r.statusMessage,
            headers: r.headers,
            body: h.from(r.bodyBase64, `base64`),
            duration: r.duration
        });
    }
    pe(), o();
    var Bn = {
        method: `GET`,
        url: ``,
        headers: [
            {
                key: `Accept`,
                value: `*/*`,
                enabled: !0
            }
        ],
        body: ``,
        contentType: `application/json`
    };
    function Vn(e) {
        return {
            id: `tab-${Date.now()}-${Math.random().toString(36).slice(2)}`,
            title: `New Request`,
            request: {
                ...Bn,
                headers: [
                    {
                        key: `Accept`,
                        value: `*/*`,
                        enabled: !0
                    }
                ]
            },
            isLoading: !1,
            ...e
        };
    }
    var B = new Map;
    function Hn(e) {
        let t = atob(e), n = h.alloc(t.length);
        for(let e = 0; e < t.length; e++)n[e] = t.charCodeAt(e);
        return n;
    }
    async function Un(e, t) {
        let n = e.getReader(), r = new TextDecoder(`utf-8`), i = ``, a = 200, o = `OK`, s = {}, c = [];
        for(;;){
            let { done: e, value: t } = await n.read();
            if (e) break;
            i += r.decode(t, {
                stream: !0
            });
            let l = i.split(`
`);
            i = l.pop() || ``;
            for (let e of l){
                if (!e.trim()) continue;
                let t = JSON.parse(e);
                if (t.type === `error`) throw Error(t.error?.message || `Send request failed on server`);
                t.type === `response-head` ? (a = t.statusCode ?? a, o = t.statusMessage || o, t.headers && (s = xt(t.headers))) : t.type === `response-body-part` && t.rawBody && c.push(Hn(t.rawBody));
            }
        }
        if (i.trim()) {
            let e = JSON.parse(i);
            if (e.type === `error`) throw Error(e.error?.message || `Send request failed on server`);
            e.type === `response-body-part` && e.rawBody && c.push(Hn(e.rawBody));
        }
        let l = h.concat(c);
        return {
            statusCode: a,
            statusMessage: o,
            headers: s,
            bodyBase64: l.toString(`base64`),
            contentType: s[`content-type`]?.split(`;`)[0] || ``,
            duration: Date.now() - t,
            bodySize: l.byteLength
        };
    }
    V = n()(c((e, t)=>{
        let n = Vn();
        return {
            tabs: [
                n
            ],
            activeTabId: n.id,
            history: [],
            addTab: ()=>{
                let t = Vn();
                e((e)=>({
                        tabs: [
                            ...e.tabs,
                            t
                        ],
                        activeTabId: t.id
                    }));
            },
            closeTab: (t)=>{
                e((e)=>{
                    if (e.tabs.length <= 1) {
                        let e = Vn();
                        return {
                            tabs: [
                                e
                            ],
                            activeTabId: e.id
                        };
                    }
                    let n = e.tabs.filter((e)=>e.id !== t);
                    return {
                        tabs: n,
                        activeTabId: e.activeTabId === t ? n[Math.max(0, e.tabs.findIndex((e)=>e.id === t) - 1)]?.id || n[0].id : e.activeTabId
                    };
                }), B.get(t)?.abort(), B.delete(t);
            },
            setActiveTab: (t)=>e({
                    activeTabId: t
                }),
            updateRequest: (t, n)=>e((e)=>({
                        tabs: e.tabs.map((e)=>e.id === t ? {
                                ...e,
                                request: {
                                    ...e.request,
                                    ...n
                                },
                                title: n.url ? Wn(n.url) : e.title
                            } : e)
                    })),
            addHeader: (t)=>e((e)=>({
                        tabs: e.tabs.map((e)=>e.id === t ? {
                                ...e,
                                request: {
                                    ...e.request,
                                    headers: [
                                        ...e.request.headers,
                                        {
                                            key: ``,
                                            value: ``,
                                            enabled: !0
                                        }
                                    ]
                                }
                            } : e)
                    })),
            removeHeader: (t, n)=>e((e)=>({
                        tabs: e.tabs.map((e)=>e.id === t ? {
                                ...e,
                                request: {
                                    ...e.request,
                                    headers: e.request.headers.filter((e, t)=>t !== n)
                                }
                            } : e)
                    })),
            updateHeader: (t, n, r, i)=>e((e)=>({
                        tabs: e.tabs.map((e)=>e.id === t ? {
                                ...e,
                                request: {
                                    ...e.request,
                                    headers: e.request.headers.map((e, t)=>t === n ? {
                                            ...e,
                                            key: r,
                                            value: i
                                        } : e)
                                }
                            } : e)
                    })),
            toggleHeader: (t, n)=>e((e)=>({
                        tabs: e.tabs.map((e)=>e.id === t ? {
                                ...e,
                                request: {
                                    ...e.request,
                                    headers: e.request.headers.map((e, t)=>t === n ? {
                                            ...e,
                                            enabled: !e.enabled
                                        } : e)
                                }
                            } : e)
                    })),
            sendRequest: async (n, r)=>{
                let i = t().tabs.find((e)=>e.id === n);
                if (!i) return;
                let a = new AbortController;
                B.set(n, a), e((e)=>({
                        tabs: e.tabs.map((e)=>e.id === n ? {
                                ...e,
                                isLoading: !0,
                                sentEventId: void 0,
                                response: void 0,
                                error: void 0
                            } : e)
                    }));
                let o = Date.now(), s = crypto.randomUUID();
                try {
                    let { url: t } = i.request;
                    t && !t.match(/^https?:\/\//) && (t = `https://${t}`);
                    let r = {};
                    i.request.headers.filter((e)=>e.enabled && e.key.trim()).forEach((e)=>{
                        r[e.key] = e.value;
                    });
                    let c = ![
                        `GET`,
                        `HEAD`,
                        `OPTIONS`
                    ].includes(i.request.method.toUpperCase());
                    c && i.request.body && !r[`content-type`] && !r[`Content-Type`] && (r[`Content-Type`] = i.request.contentType);
                    try {
                        let e = new URL(t);
                        !r.host && !r.Host && (r.Host = e.host);
                    } catch  {}
                    if (c && i.request.body) {
                        let e = new TextEncoder().encode(i.request.body);
                        !r[`content-length`] && !r[`Content-Length`] && (r[`Content-Length`] = String(e.length));
                    }
                    if (!$.getState().serverOnline) throw Error(`HttpToolkit server is offline. Start httptoolkit-server (npm start in httptoolkit-server) to send requests.`);
                    In(s, i.request, t), e((e)=>({
                            tabs: e.tabs.map((e)=>e.id === n ? {
                                    ...e,
                                    sentEventId: s
                                } : e)
                        }));
                    let l = Ye(Z.getState()), u = c && i.request.body ? bt(i.request.body) : void 0, d = await te(`POST`, `/client/send`, {
                        body: {
                            request: {
                                method: i.request.method,
                                url: t,
                                headers: Object.entries(r).map(([e, t])=>[
                                        e,
                                        t
                                    ]),
                                rawBody: u
                            },
                            options: l
                        },
                        signal: a.signal
                    });
                    if (!d.ok || !d.body) throw Error(`Server send failed: HTTP ${d.status}`);
                    let f = await Un(d.body, o);
                    Ln(s, {
                        statusCode: f.statusCode,
                        statusMessage: f.statusMessage,
                        headers: f.headers,
                        body: h.from(f.bodyBase64, `base64`),
                        duration: f.duration
                    }), e((e)=>{
                        let t = [
                            {
                                id: `hist-${Date.now()}-${Math.random().toString(36).slice(2)}`,
                                timestamp: Date.now(),
                                title: i.title,
                                request: JSON.parse(JSON.stringify(i.request)),
                                response: {
                                    ...f,
                                    exchangeId: s
                                }
                            },
                            ...e.history
                        ].slice(0, 50);
                        return {
                            tabs: e.tabs.map((e)=>e.id === n ? {
                                    ...e,
                                    isLoading: !1,
                                    sentEventId: s,
                                    response: {
                                        ...f,
                                        exchangeId: s
                                    }
                                } : e),
                            history: t
                        };
                    });
                } catch (t) {
                    if (t.name === `AbortError`) Rn(s, `Request cancelled`), e((e)=>({
                            tabs: e.tabs.map((e)=>e.id === n ? {
                                    ...e,
                                    isLoading: !1,
                                    sentEventId: s,
                                    error: `Request cancelled`
                                } : e)
                        }));
                    else {
                        let r = t instanceof de ? t.message : t?.message || `Request failed`;
                        Rn(s, r), e((e)=>({
                                tabs: e.tabs.map((e)=>e.id === n ? {
                                        ...e,
                                        isLoading: !1,
                                        sentEventId: s,
                                        error: r
                                    } : e)
                            }));
                    }
                } finally{
                    B.delete(n);
                }
            },
            cancelRequest: (t)=>{
                B.get(t)?.abort(), B.delete(t), e((e)=>({
                        tabs: e.tabs.map((e)=>e.id === t ? {
                                ...e,
                                isLoading: !1
                            } : e)
                    }));
            },
            openEventInSend: (t, n)=>{
                let r = wt(t), i = Vn({
                    title: Wn(r.url),
                    request: r
                });
                e((e)=>{
                    if (n?.replaceActiveTab) {
                        let t = e.activeTabId;
                        return {
                            tabs: e.tabs.map((e)=>e.id === t ? {
                                    ...i,
                                    id: t
                                } : e),
                            activeTabId: t
                        };
                    }
                    return {
                        tabs: [
                            ...e.tabs,
                            i
                        ],
                        activeTabId: i.id
                    };
                }), q({
                    type: `success`,
                    domain: `send`,
                    message: `Request loaded in Send — edit and send when ready`,
                    dedupKey: `send:loaded`
                });
            },
            loadRequestFromEvent: (e)=>{
                t().openEventInSend(e);
            },
            importCurl: (t, n)=>{
                try {
                    let r = Gn(n);
                    e((e)=>({
                            tabs: e.tabs.map((e)=>e.id === t ? {
                                    ...e,
                                    title: Wn(r.url),
                                    request: {
                                        ...e.request,
                                        method: r.method,
                                        url: r.url,
                                        headers: r.headers,
                                        body: r.body || ``
                                    }
                                } : e)
                        }));
                } catch  {
                    u(`Failed to parse cURL command`, `cURL parsing error`);
                }
            },
            clearHistory: ()=>{
                let t = V.getState().history.map((e)=>e.id);
                Promise.all([
                    a(()=>import(`./userConfigSync-BIjNXwGj.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        }), []),
                    a(()=>import(`./useAuthStore-CEO3CpNF.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        }), __vite__mapDeps([20,21,2,3,4,8,15,6,7,1]))
                ]).then(([{ recordSyncDeletes: e, triggerSyncAfterDelete: n }, { useAuthStore: r }])=>{
                    e(`sendHistory`, t), n(()=>r.getState().user);
                }), e({
                    history: []
                });
            },
            deleteHistoryEntry: (t)=>{
                Promise.all([
                    a(()=>import(`./userConfigSync-BIjNXwGj.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        }), []),
                    a(()=>import(`./useAuthStore-CEO3CpNF.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        }), __vite__mapDeps([20,21,2,3,4,8,15,6,7,1]))
                ]).then(([{ recordSyncDelete: e, triggerSyncAfterDelete: n }, { useAuthStore: r }])=>{
                    e(`sendHistory`, t), n(()=>r.getState().user);
                }), e((e)=>({
                        history: e.history.filter((e)=>e.id !== t)
                    }));
            },
            loadFromHistory: (t)=>{
                let n = `hist-restore-${t.id}-${Date.now()}`, r = t.request.url;
                r && !r.match(/^https?:\/\//) && (r = `https://${r}`), zn(n, t.request, r, t.response);
                let i = Vn({
                    title: t.title,
                    request: JSON.parse(JSON.stringify(t.request)),
                    sentEventId: n,
                    response: {
                        ...t.response,
                        exchangeId: n
                    }
                });
                e((e)=>({
                        tabs: [
                            ...e.tabs,
                            i
                        ],
                        activeTabId: i.id
                    }));
            }
        };
    }, {
        name: `httptoolkit-send-state`,
        partialize: (e)=>({
                tabs: e.tabs.map((e)=>({
                        ...e,
                        isLoading: !1,
                        error: void 0,
                        sentEventId: e.sentEventId
                    })),
                activeTabId: e.activeTabId,
                history: e.history
            })
    }));
    function Wn(e) {
        try {
            let t = new URL(e.match(/^https?:\/\//) ? e : `https://${e}`), n = t.pathname === `/` ? `` : t.pathname;
            return `${t.hostname}${n}`.slice(0, 30) || `New Request`;
        } catch  {
            return e.slice(0, 30) || `New Request`;
        }
    }
    function Gn(e) {
        let t = `GET`, n = ``, r = [], i = ``, a = e.match(/curl\s+(?:'([^']+)'|"([^"]+)"|(\S+))/);
        a && (n = a[1] || a[2] || a[3] || ``);
        let o = e.match(/-X\s+(\w+)/);
        o && (t = o[1].toUpperCase());
        let s = e.matchAll(/-H\s+(?:'([^']+)'|"([^"]+)")/g);
        for (let e of s){
            let t = e[1] || e[2] || ``, n = t.indexOf(`:`);
            n > 0 && r.push({
                key: t.slice(0, n).trim(),
                value: t.slice(n + 1).trim(),
                enabled: !0
            });
        }
        let c = e.match(/(?:--data|-d)\s+(?:'([^']*)'|"([^"]*)")/);
        return c && (i = c[1] || c[2] || ``, t === `GET` && (t = `POST`)), {
            method: t,
            url: n,
            headers: r,
            body: i
        };
    }
    var Kn = {
        favorites: [],
        userRules: [],
        sendHistory: []
    }, qn = 500;
    function Jn(e, t) {
        if (!t || e.includes(t)) return e;
        let n = [
            ...e,
            t
        ];
        return n.length <= qn ? n : n.slice(n.length - qn);
    }
    function Yn(e, t, n) {
        return {
            ...e,
            [t]: Jn(e[t], n)
        };
    }
    function Xn(e, t) {
        if (!t) return e;
        let n = (e, t)=>{
            if (!t?.length) return e;
            let n = new Set(e), r = [
                ...e
            ];
            for (let e of t)n.has(e) || (n.add(e), r.push(e));
            return r.length > qn ? r.slice(r.length - qn) : r;
        };
        return {
            favorites: n(e.favorites, t.favorites),
            userRules: n(e.userRules, t.userRules),
            sendHistory: n(e.sendHistory, t.sendHistory)
        };
    }
    function H(e, t, n) {
        let r = new Set(n[t]);
        return e.filter((e)=>!r.has(e.id));
    }
    var Zn = `default-group`;
    function Qn(e) {
        return JSON.parse(JSON.stringify(e));
    }
    function $n(e) {
        let t = [], n = (e)=>{
            let t = [];
            for (let r of e.items)if (g(r)) {
                let e = n(r);
                e && t.push(e);
            } else w(r) || t.push(r);
            return t.length ? {
                ...e,
                items: t
            } : null;
        };
        for (let r of e)if (g(r)) {
            if (r.id === Zn) continue;
            let e = n(r);
            e && t.push(e);
        } else w(r) || t.push(r);
        return t;
    }
    function er(e) {
        let t = e.find((e)=>g(e) && e.id === Zn);
        if (t && g(t)) return Qn(t);
        let n = T.find((e)=>g(e) && e.id === Zn);
        return n && g(n) ? Qn(n) : null;
    }
    function tr(e, t, n) {
        let r = new Map, i = (e)=>{
            for (let t of e)n.userRules.includes(t.id) || r.set(t.id, t);
        };
        return i(H(t, `userRules`, n)), i(H(e, `userRules`, n)), [
            ...r.values()
        ];
    }
    function nr(e, t, n = Kn) {
        let r = D(e), i = tr($n(r), t, n), a = er(r), o = [
            ...i
        ];
        return a && o.push(a), D(o);
    }
    function rr(e) {
        return Array.isArray(e.userRules) ? e.userRules : Array.isArray(e.draftRules) ? $n(e.draftRules) : [];
    }
    let ir, ar, or, sr, cr, lr, ur, fr;
    ir = 1024 * 1024;
    ar = `httptoolkit-sync-tombstones`;
    or = 64 * 1024;
    sr = 50;
    cr = 200;
    lr = 1500;
    ur = 900 * 1024;
    dr = 5;
    fr = 45457;
    function pr() {
        return `http://127.0.0.1:${d?.() ?? fr}`;
    }
    mr = function() {
        return Math.max(5, s.getState().cloudSyncIntervalHours ?? 5) * 60 * 60 * 1e3;
    };
    var U = !1, hr = null, W = {
        lastSyncAt: null,
        lastSyncError: null
    };
    function gr(e) {
        return `httptoolkit-last-sync-${e}`;
    }
    function _r(e) {
        try {
            let t = localStorage.getItem(gr(e));
            if (!t) return {};
            let n = JSON.parse(t);
            return {
                lastSyncAt: typeof n.lastSyncAt == `number` ? n.lastSyncAt : null,
                lastSyncError: n.lastSyncError ?? null
            };
        } catch  {
            return {};
        }
    }
    function vr(e, t) {
        try {
            localStorage.setItem(gr(e), JSON.stringify({
                lastSyncAt: t.lastSyncAt,
                lastSyncError: t.lastSyncError
            }));
        } catch  {}
    }
    yr = function() {
        return {
            ...W
        };
    };
    function br(e) {
        W = {
            lastSyncAt: Date.now(),
            lastSyncError: null
        }, vr(e, W);
    }
    function xr(e, t) {
        W = {
            ...W,
            lastSyncError: t
        }, vr(e, W);
    }
    let Sr;
    Sr = (e)=>{
        let t = (e || ``).toLowerCase();
        return t.startsWith(`image/`) || t.startsWith(`video/`) || t.startsWith(`audio/`);
    };
    Cr = (e)=>Sr(e.contentType) || (e.requestBody?.length || 0) > ir || (e.responseBody?.length || 0) > ir;
    function wr(e) {
        if (!Cr(e)) return e;
        let { requestBody: t, responseBody: n, ...r } = e;
        return r;
    }
    Tr = function(e) {
        if (e == null) return e;
        if (Array.isArray(e)) return e.map((e)=>Tr(e));
        if (typeof e == `object`) {
            let t = {};
            for (let [n, r] of Object.entries(e))r !== void 0 && (t[n] = Tr(r));
            return t;
        }
        return e;
    };
    Er = function(e) {
        return e.slice(0, sr).map((e)=>{
            let t = e.request.body, n = e.response.bodyBase64;
            return {
                ...e,
                request: {
                    ...e.request,
                    body: t.length > or ? `` : t
                },
                response: {
                    ...e.response,
                    bodyBase64: n.length > or ? `` : n
                }
            };
        });
    };
    function Dr(e) {
        let { updatedAt: t, ...n } = e;
        return JSON.stringify(n);
    }
    var G = Or();
    function Or() {
        try {
            let e = localStorage.getItem(ar);
            return e ? Xn(Kn, JSON.parse(e)) : {
                ...Kn
            };
        } catch  {
            return {
                ...Kn
            };
        }
    }
    function kr() {
        try {
            localStorage.setItem(ar, JSON.stringify(G));
        } catch  {}
    }
    Ar = function() {
        return G;
    };
    jr = function(e, t) {
        t && (G = Yn(G, e, t), kr());
    };
    Mr = function(e, t) {
        for (let n of t)jr(e, n);
    };
    function Nr() {
        let e = `httptoolkit-device-id`;
        try {
            let t = localStorage.getItem(e);
            return t || (t = crypto.randomUUID(), localStorage.setItem(e, t)), t;
        } catch  {
            return `unknown`;
        }
    }
    function Pr() {
        let e = Z.getState(), t = $.getState(), n = S(t.dnsServers);
        return {
            whitelistedCertificateHosts: e.whitelistedCertificateHosts,
            upstreamProxyType: e.upstreamProxyType,
            upstreamProxyHost: e.upstreamProxyHost ?? ``,
            upstreamNoProxyHosts: e.upstreamNoProxyHosts,
            http2Enabled: t.http2Enabled,
            ...n.length ? {
                dnsServers: n
            } : {},
            ...t.portConfig ? {
                portConfig: t.portConfig
            } : {}
        };
    }
    function Fr() {
        let e = s.getState();
        return {
            theme: e.theme,
            detailsCardOrder: e.detailsCardOrder,
            detailsCardExpanded: e.detailsCardExpanded,
            detailsCardDisabled: e.detailsCardDisabled,
            detailsCardMode: e.detailsCardMode,
            exportSnippetFormat: e.exportSnippetFormat,
            preferredShell: e.preferredShell,
            activeFilter: e.activeFilter,
            customFilters: e.customFilters
        };
    }
    Ir = function(e) {
        return {
            rules: `users/${e}/rules/draft.json`,
            favorites: `users/${e}/favorites.json`,
            sendHistory: `users/${e}/send/history.json`
        };
    };
    Lr = function(e) {
        let t = Ir(e);
        return {
            version: 3,
            updatedAt: Date.now(),
            deviceId: Nr(),
            storagePaths: t,
            connection: Pr(),
            ui: Fr()
        };
    };
    Rr = function() {
        let e = H(O.getState().favorites.slice(0, cr).map(wr), `favorites`, G), t = H($n(Z.getState().draftRules).filter((e)=>typeof e == `object` && !!e && `id` in e), `userRules`, G), n = H(Er(V.getState().history), `sendHistory`, G);
        return {
            version: 5,
            updatedAt: Date.now(),
            deviceId: Nr(),
            favorites: e,
            rules: {
                userRules: t
            },
            sendHistory: n,
            connection: Pr(),
            ui: Fr(),
            deleted: G
        };
    };
    zr = function() {
        return Rr();
    };
    Br = async function(e, t) {
        let n = {
            version: 4,
            updatedAt: e.updatedAt,
            deviceId: e.deviceId ?? `unknown`,
            connection: e.connection,
            ui: e.ui
        }, r = await t(e.storagePaths.rules);
        r && typeof r == `object` && r && (`userRules` in r && Array.isArray(r.userRules) ? n.rules = {
            userRules: r.userRules
        } : `draftRules` in r && Array.isArray(r.draftRules) && (n.rules = {
            userRules: $n(r.draftRules)
        }));
        let i = await t(e.storagePaths.favorites);
        Array.isArray(i) && (n.favorites = i);
        let a = await t(e.storagePaths.sendHistory);
        return Array.isArray(a) && (n.sendHistory = a), n;
    };
    Vr = function() {
        let e = Rr();
        return {
            version: 2,
            updatedAt: e.updatedAt,
            deviceId: e.deviceId,
            favorites: e.favorites,
            rules: {
                draftRules: e.rules.userRules
            },
            connection: e.connection,
            ui: e.ui
        };
    };
    function Hr(e) {
        let t = {};
        Array.isArray(e.favorites) && (t.favorites = e.favorites), Array.isArray(e.rules) && (t.rules = {
            draftRules: e.rules
        });
        let n = e.settings;
        return n && (t.ui = {
            theme: n.theme || `automatic`,
            detailsCardOrder: n.detailsCardOrder || [],
            detailsCardMode: n.detailsCardMode || {},
            exportSnippetFormat: n.exportSnippetFormat || `curl`,
            preferredShell: n.preferredShell || `bash`,
            activeFilter: ``,
            customFilters: {}
        }), t;
    }
    Ur = function(e, t, n = G) {
        let r = H(t, `favorites`, n), i = H(e, `favorites`, n), a = new Map;
        for (let e of r)a.set(e.id, e);
        for (let e of i){
            let t = a.get(e.id);
            (!t || e.savedAt > t.savedAt) && a.set(e.id, e);
        }
        return [
            ...a.values()
        ].sort((e, t)=>t.savedAt - e.savedAt);
    };
    Wr = function(e, t, n = G) {
        let r = H(t, `sendHistory`, n), i = H(e, `sendHistory`, n), a = new Map;
        for (let e of r)a.set(e.id, e);
        for (let e of i){
            let t = a.get(e.id);
            (!t || e.timestamp > t.timestamp) && a.set(e.id, e);
        }
        return [
            ...a.values()
        ].sort((e, t)=>t.timestamp - e.timestamp).slice(0, sr);
    };
    function Gr(e) {
        Z.setState({
            whitelistedCertificateHosts: e.whitelistedCertificateHosts ?? [],
            upstreamProxyType: e.upstreamProxyType,
            upstreamProxyHost: e.upstreamProxyHost ?? ``,
            upstreamNoProxyHosts: e.upstreamNoProxyHosts ?? []
        }), (e.http2Enabled !== void 0 || e.dnsServers || e.portConfig) && $.setState({
            ...e.http2Enabled === void 0 ? {} : {
                http2Enabled: e.http2Enabled
            },
            ...e.dnsServers ? {
                dnsServers: S(e.dnsServers)
            } : {},
            ...e.portConfig ? {
                portConfig: e.portConfig
            } : {}
        });
    }
    function Kr(e) {
        let t = s.getState();
        e.theme && t.setTheme(e.theme), e.detailsCardOrder?.length && t.setDetailsCardOrder(e.detailsCardOrder), e.detailsCardExpanded && s.setState({
            detailsCardExpanded: e.detailsCardExpanded
        }), e.detailsCardDisabled && t.setDetailsCardDisabled(e.detailsCardDisabled), e.detailsCardMode && Object.entries(e.detailsCardMode).forEach(([e, n])=>t.setDetailsCardMode(e, n)), e.exportSnippetFormat && t.setExportSnippetFormat(e.exportSnippetFormat), e.preferredShell && t.setPreferredShell(e.preferredShell), e.activeFilter !== void 0 && t.setActiveFilter(e.activeFilter), e.customFilters && s.setState({
            customFilters: e.customFilters
        });
    }
    async function qr(e) {
        return null;
    }
    function Jr(e, t) {
        if (`deleted` in e && e.deleted && (G = Xn(G, e.deleted), kr()), e.connection && Gr(e.connection), e.ui && Kr(e.ui), e.favorites) {
            let n = O.getState().favorites, r = H(e.favorites, `favorites`, G), i = t === `replace` ? r : Ur(n, r, G);
            O.setState({
                favorites: i
            });
        }
        if (e.rules) {
            let n = rr(e.rules), r = Z.getState().draftRules, i = nr(t === `replace` ? [] : r, n, G);
            Z.setState({
                draftRules: i,
                rules: i
            });
        }
        let n = `sendHistory` in e && Array.isArray(e.sendHistory) ? e.sendHistory : void 0;
        if (n) {
            let e = V.getState().history, r = H(n, `sendHistory`, G), i = t === `replace` ? r : Wr(e, r, G);
            V.setState({
                history: i
            });
        }
    }
    async function Yr() {
        let { activeSessionId: e } = $.getState();
        if (e) try {
            Z.getState().saveRules({
                silent: !0
            });
        } catch (e) {
            r.warn(`sync`, `saveRules after restore failed`, e);
        }
    }
    Xr = async function(e, t = {}) {
        let n = t.mode ?? `merge`;
        if (e.version === 5 || e.version === 4) Jr(e, n);
        else if (e.version === 3 && e.storagePaths) Jr(await Br(e, qr), n);
        else {
            let t;
            t = e.version === 2 ? e : Hr(e), Jr(t, n);
        }
        await Yr();
    };
    Zr = async function(e, t) {
        if (U || !e?.email) return;
        let n = Tr(Rr()), a = Dr(n);
        if (a === hr && t?.reason !== `manual`) {
            r.verbose(`supabase`, `Skipping backup — snapshot unchanged`);
            return;
        }
        let o = JSON.stringify(n);
        o.length > ur && r.warn(`supabase`, `Backup payload is large (${o.length} bytes)`), r.verbose(`supabase`, `Writing settings backup`, {
            bytes: o.length,
            favorites: n.favorites.length,
            rules: n.rules.userRules.length,
            sendHistory: n.sendHistory.length
        });
        try {
            let i = new Date().toISOString(), o = await fetch(`${pr()}/config/backup`, {
                method: `POST`,
                headers: {
                    "Content-Type": `application/json`
                },
                body: JSON.stringify(n)
            });
            if (!o.ok) {
                let e = await o.text().catch(()=>`backup request failed`);
                throw Error(e);
            }
            hr = a, ti = new Date(i).getTime(), br(e.id), r.verbose(`sync`, `Settings backup written`), t?.reason === `manual` && mi(e, {
                type: `success`,
                title: `Cloud sync complete`,
                message: `Your settings were synced to the cloud.`,
                source: `sync`
            });
        } catch (t) {
            let n = t instanceof Error ? t.message : String(t);
            throw xr(e.id, n), i(t, {
                domain: `supabase`,
                code: `backup_failed`
            }), mi(e, {
                type: `error`,
                title: `Cloud sync failed`,
                message: n.slice(0, 500),
                source: `sync`
            }), q({
                type: `error`,
                domain: `sync`,
                message: `Cloud sync failed: ${n}`,
                dedupKey: `supabase-sync-error`
            }), t;
        }
    };
    Qr = async function(e) {
        if (!e?.email) return {
            hadCloudData: !1
        };
        W = {
            ...W,
            ..._r(e.id)
        }, U = !0;
        let t = !1;
        try {
            r.verbose(`sync`, `Restoring user data`);
            let e = await fetch(`${pr()}/config/restore`);
            if (!e.ok) {
                if (e.status === 404) return {
                    hadCloudData: !1
                };
                let t = await e.text().catch(()=>`restore request failed`);
                throw Error(t);
            }
            let n = await e.json();
            n?.payload && typeof n.payload == `object` && (t = !0, await Xr(n.payload, {
                mode: `merge`
            }));
        } finally{
            U = !1;
            try {
                await Zr(e, {
                    reason: `change`
                });
            } catch (e) {
                r.warn(`supabase`, `Post-restore backup failed`, e);
            }
        }
        return {
            hadCloudData: t
        };
    };
    var K = null, $r = null, ei = null, ti = 0;
    function ni() {
        K &&= (clearTimeout(K), null);
    }
    function ri() {
        K = null;
        let e = $r;
        $r = null, e && Zr(e, {
            reason: `change`
        }).catch((e)=>{
            r.warn(`sync`, `sync failed`, e);
        });
    }
    ii = function(e, t) {
        if (!U) {
            if ($r = e, t?.force) {
                ni(), Zr(e, {
                    reason: `manual`
                }).catch((e)=>{
                    r.warn(`sync`, `sync failed`, e);
                });
                return;
            }
            K && clearTimeout(K), K = setTimeout(ri, lr);
        }
    };
    ai = ii;
    oi = function(e) {
        let t = ()=>{
            let t = e();
            t && !U && ii(t);
        }, n = (e)=>JSON.stringify({
                draftRules: e.draftRules,
                whitelistedCertificateHosts: e.whitelistedCertificateHosts,
                upstreamNoProxyHosts: e.upstreamNoProxyHosts,
                upstreamProxyType: e.upstreamProxyType,
                upstreamProxyHost: e.upstreamProxyHost
            }), r = (e)=>JSON.stringify({
                http2Enabled: e.http2Enabled,
                dnsServers: e.dnsServers,
                portConfig: e.portConfig
            }), i = (e)=>JSON.stringify(si(e)), a = (e)=>JSON.stringify(e.favorites.map((e)=>e.id)), o = (e)=>JSON.stringify(e.history.map((e)=>e.id)), c = n(Z.getState()), l = r($.getState()), u = i(s.getState()), d = a(O.getState()), f = o(V.getState());
        Z.subscribe((e)=>{
            let r = n(e);
            r !== c && (c = r, t());
        }), $.subscribe((e)=>{
            let n = r(e);
            n !== l && (l = n, t());
        }), s.subscribe((e)=>{
            let n = i(e);
            n !== u && (u = n, t());
        }), O.subscribe((e)=>{
            let n = a(e);
            n !== d && (d = n, t());
        }), V.subscribe((e)=>{
            let n = o(e);
            n !== f && (f = n, t());
        });
    };
    function si(e) {
        return {
            theme: e.theme,
            detailsCardOrder: e.detailsCardOrder,
            detailsCardExpanded: e.detailsCardExpanded,
            detailsCardDisabled: e.detailsCardDisabled,
            detailsCardMode: e.detailsCardMode,
            exportSnippetFormat: e.exportSnippetFormat,
            preferredShell: e.preferredShell,
            activeFilter: e.activeFilter,
            customFilters: e.customFilters
        };
    }
    ci = function(e) {
        ei?.(), ei = null;
        let t = e(), n = vt();
        if (!t || !n) return;
        let i = n.channel(`user_settings:${t.id}`).on(`postgres_changes`, {
            event: `UPDATE`,
            schema: `public`,
            table: `user_settings`,
            filter: `user_id=eq.${t.id}`
        }, (e)=>{
            let t = e.new, n = t.updated_at ? new Date(t.updated_at).getTime() : 0;
            n <= ti || n <= (W.lastSyncAt ?? 0) || !t.payload || typeof t.payload != `object` || (ti = n, U = !0, Xr(t.payload, {
                mode: `merge`
            }).catch((e)=>r.warn(`sync`, `realtime pull failed`, e)).finally(()=>{
                U = !1;
            }));
        }).subscribe();
        ei = ()=>{
            n.removeChannel(i);
        };
    };
    li = function() {
        ei?.(), ei = null;
    };
    ui = function(e) {
        let t = e();
        t && ii(t, {
            force: !0
        });
    };
    var di = 50, fi = 30;
    function pi(e) {
        return {
            id: e.id,
            type: e.type ?? `info`,
            title: e.title,
            message: e.message,
            createdAt: new Date(e.created_at).getTime(),
            read: e.read,
            source: e.source ?? void 0
        };
    }
    async function mi(e, t) {
        let n = vt();
        if (!n) return null;
        try {
            let { data: r, error: i } = await n.from(`notification_inbox`).insert({
                user_id: e.id,
                type: t.type,
                title: t.title.slice(0, 200),
                message: t.message.slice(0, 2e3),
                read: t.read ?? !1,
                source: t.source ?? null
            }).select(`id`).single();
            if (i) throw i;
            let { data: a } = await n.from(`notification_inbox`).select(`id`).eq(`user_id`, e.id).order(`created_at`, {
                ascending: !1
            });
            if (a && a.length > di) {
                let e = a.slice(di).map((e)=>e.id);
                await n.from(`notification_inbox`).delete().in(`id`, e);
            }
            return r?.id ?? null;
        } catch (e) {
            return r.warn(`supabase`, `pushInboxItem failed`, e), null;
        }
    }
    hi = function(e, t) {
        let n = vt();
        if (!n) return null;
        n.from(`notification_inbox`).select(`*`).eq(`user_id`, e.id).order(`created_at`, {
            ascending: !1
        }).limit(fi).then(({ data: e, error: n })=>{
            if (n) {
                r.warn(`supabase`, `inbox initial load error`, n), t([]);
                return;
            }
            t((e ?? []).map(pi));
        });
        let i = n.channel(`inbox:${e.id}`).on(`postgres_changes`, {
            event: `*`,
            schema: `public`,
            table: `notification_inbox`,
            filter: `user_id=eq.${e.id}`
        }, ()=>{
            n.from(`notification_inbox`).select(`*`).eq(`user_id`, e.id).order(`created_at`, {
                ascending: !1
            }).limit(fi).then(({ data: e })=>t((e ?? []).map(pi)));
        }).subscribe();
        return ()=>{
            n.removeChannel(i);
        };
    };
    gi = async function(e, t) {
        let n = vt();
        if (!n || t.length === 0) return;
        let { error: r } = await n.from(`notification_inbox`).update({
            read: !0
        }).eq(`user_id`, e.id).in(`id`, t);
        if (r) throw r;
    };
    _i = async function(e, t) {
        let n = vt();
        if (!n || t.length === 0) return;
        let { error: r } = await n.from(`notification_inbox`).delete().eq(`user_id`, e.id).in(`id`, t);
        if (r) throw r;
    };
    var vi = new Map, yi = 2e3;
    function bi(e) {
        let t = Date.now(), n = vi.get(e);
        if (n !== void 0 && t - n < yi) return !0;
        if (vi.set(e, t), vi.size > 100) {
            let e = t - yi;
            for (let [t, n] of vi)n < e && vi.delete(t);
        }
        return !1;
    }
    q = function(e) {
        let t = e.domain ?? `general`;
        if (r.verbose(t, e.message), e.silent || bi(e.dedupKey ?? `${t}:${e.type}:${e.message}`)) return;
        let n = l.getState().addNotification(e.type, e.message, e.durationMs);
        if (t === `sync` && (e.type === `error` || e.type === `warning`)) {
            let t = he.getState().user;
            t && mi(t, {
                type: e.type,
                title: e.type === `error` ? `Sync error` : `Sync warning`,
                message: e.message.slice(0, 500),
                source: `sync`
            });
        }
        return n;
    };
    o();
    var xi = 50, J = [], Y = -1;
    function X(e) {
        Y < J.length - 1 && J.splice(Y + 1), J.push(JSON.parse(JSON.stringify(e))), J.length > xi ? J.shift() : Y++;
    }
    Z = n()(c((e, t)=>({
            rules: T,
            draftRules: T,
            canUndo: !1,
            canRedo: !1,
            upstreamProxyType: `system`,
            upstreamProxyHost: void 0,
            upstreamNoProxyHosts: Kt([]),
            whitelistedCertificateHosts: [
                `localhost`
            ],
            clientCertificateHostMap: {},
            undo: ()=>{
                Y > 0 && (Y--, e({
                    draftRules: JSON.parse(JSON.stringify(J[Y])),
                    canUndo: Y > 0,
                    canRedo: Y < J.length - 1
                }));
            },
            redo: ()=>{
                Y < J.length - 1 && (Y++, e({
                    draftRules: JSON.parse(JSON.stringify(J[Y])),
                    canUndo: Y > 0,
                    canRedo: Y < J.length - 1
                }));
            },
            setUpstreamProxy: (t, n, r = [])=>e({
                    upstreamProxyType: t,
                    upstreamProxyHost: n,
                    upstreamNoProxyHosts: r
                }),
            addCertificateHost: (t)=>e((e)=>({
                        whitelistedCertificateHosts: [
                            ...e.whitelistedCertificateHosts.filter((e)=>e !== t),
                            t
                        ]
                    })),
            removeCertificateHost: (t)=>e((e)=>({
                        whitelistedCertificateHosts: e.whitelistedCertificateHosts.filter((e)=>e !== t)
                    })),
            addUpstreamNoProxyHost: (t)=>e((e)=>({
                        upstreamNoProxyHosts: Kt([
                            ...e.upstreamNoProxyHosts,
                            t.trim()
                        ])
                    })),
            removeUpstreamNoProxyHost: (t)=>e((e)=>({
                        upstreamNoProxyHosts: Kt(e.upstreamNoProxyHosts.filter((e)=>e !== t))
                    })),
            addClientCertificate: (t, n)=>e((e)=>({
                        clientCertificateHostMap: {
                            ...e.clientCertificateHostMap,
                            [t]: n
                        }
                    })),
            removeClientCertificate: (t)=>e((e)=>{
                    let n = {
                        ...e.clientCertificateHostMap
                    };
                    return delete n[t], {
                        clientCertificateHostMap: n
                    };
                }),
            setRules: (t)=>e({
                    rules: t
                }),
            setDraftRules: (n)=>{
                X(t().draftRules), e({
                    draftRules: n,
                    canUndo: Y > 0,
                    canRedo: Y < J.length - 1
                });
            },
            saveRules: (t)=>e((e)=>{
                    X(e.draftRules);
                    let n = D(e.draftRules);
                    try {
                        let { serverOnline: r } = $.getState();
                        r && p() && gt(n, e, t);
                    } catch (e) {
                        t?.silent || u(e, `Failed to sync rules with server`);
                    }
                    return {
                        rules: n,
                        draftRules: n
                    };
                }),
            resetRuleDrafts: ()=>e((e)=>(X(e.draftRules), {
                        draftRules: [
                            ...e.rules
                        ],
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    })),
            resetRulesToDefault: ()=>e((e)=>(X(e.draftRules), {
                        rules: T,
                        draftRules: T,
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    })),
            addUserRule: (t)=>e((e)=>(X(e.draftRules), q({
                        type: `success`,
                        domain: `rules`,
                        message: `Rule added (save changes to apply to proxy)`,
                        dedupKey: `rules:added`
                    }), {
                        draftRules: [
                            t,
                            ...e.draftRules
                        ],
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    })),
            addRule: (t, n)=>e((e)=>{
                    if (X(e.draftRules), !n) return {
                        draftRules: [
                            t,
                            ...e.draftRules
                        ],
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                    let r = (e)=>e.map((e)=>`items` in e ? e.id === n ? {
                                ...e,
                                items: [
                                    t,
                                    ...e.items
                                ]
                            } : {
                                ...e,
                                items: r(e.items)
                            } : e);
                    return {
                        draftRules: r(e.draftRules),
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                }),
            deleteRule: (t)=>{
                Promise.all([
                    a(()=>import(`./userConfigSync-BIjNXwGj.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        }), []),
                    a(()=>import(`./useAuthStore-CEO3CpNF.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        }), __vite__mapDeps([20,21,2,3,4,8,15,6,7,1]))
                ]).then(([{ recordSyncDelete: e, triggerSyncAfterDelete: n }, { useAuthStore: r }])=>{
                    e(`userRules`, t), n(()=>r.getState().user);
                }), e((e)=>{
                    X(e.draftRules);
                    let n = (e)=>e.filter((e)=>e.id !== t).map((e)=>`items` in e ? {
                                ...e,
                                items: n(e.items)
                            } : e);
                    return {
                        draftRules: n(e.draftRules),
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                });
            },
            updateRule: (t, n)=>e((e)=>{
                    X(e.draftRules);
                    let r = (e)=>e.map((e)=>e.id === t ? {
                                ...e,
                                ...n
                            } : `items` in e ? {
                                ...e,
                                items: r(e.items)
                            } : e);
                    return {
                        draftRules: r(e.draftRules),
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                }),
            cloneRule: (t)=>e((e)=>{
                    X(e.draftRules);
                    let n = (e)=>{
                        let r = [];
                        return e.forEach((e)=>{
                            if (r.push(e), e.id === t) {
                                let t = JSON.parse(JSON.stringify(e));
                                t.id = `rule-${Date.now()}-${Math.random().toString(36).slice(2)}`, t.title = `${t.title} (Copy)`, r.push(t);
                            } else `items` in e && (r[r.length - 1] = {
                                ...e,
                                items: n(e.items)
                            });
                        }), r;
                    };
                    return {
                        draftRules: n(e.draftRules),
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                }),
            importRules: (t)=>e((e)=>(X(e.draftRules), {
                        draftRules: [
                            ...e.draftRules,
                            ...t
                        ],
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    })),
            moveDraftRule: (t, n)=>e((e)=>{
                    X(e.draftRules);
                    let r = b(e.draftRules), i = _(r), a = v(i, t), o = v(i, n), s = t[t.length - 1], c = n[n.length - 1], [l] = a.items.splice(s, 1);
                    if (o.items.splice(c, 0, l), a.items.length === 0 && !Ne(a)) {
                        let e = Fe(i, {
                            id: a.id
                        });
                        e && Ie(i, e);
                    }
                    return {
                        draftRules: r,
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                }),
            combineDraftRulesAsGroup: (t, n)=>e((e)=>{
                    X(e.draftRules);
                    let r = b(e.draftRules), i = _(r), a = y(i, t), o = v(i, n), s = n[n.length - 1], c = o.items[s];
                    o.items[s] = {
                        id: crypto.randomUUID(),
                        title: `New group`,
                        collapsed: !0,
                        items: [
                            c,
                            a
                        ]
                    };
                    let l = v(i, t), u = t[t.length - 1];
                    if (l.items.splice(u, 1), l.items.length === 0 && !Ne(l)) {
                        let e = Fe(i, {
                            id: l.id
                        });
                        e && Ie(i, e);
                    }
                    return {
                        draftRules: r,
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                }),
            deleteDraftItemAtPath: (t)=>{
                e((e)=>{
                    X(e.draftRules);
                    let n = b(e.draftRules);
                    return Ie(_(n), t), {
                        draftRules: n,
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                }), q({
                    type: `info`,
                    domain: `rules`,
                    message: `Rule removed from draft`
                });
            },
            updateGroupTitle: (t, n)=>e((e)=>{
                    let r = (e)=>e.map((e)=>g(e) ? e.id === t ? {
                                ...e,
                                title: n
                            } : {
                                ...e,
                                items: r(e.items)
                            } : e);
                    return {
                        draftRules: r(e.draftRules)
                    };
                }),
            toggleGroupCollapsed: (t)=>e((e)=>{
                    let n = (e)=>e.map((e)=>g(e) ? e.id === t ? {
                                ...e,
                                collapsed: !e.collapsed
                            } : {
                                ...e,
                                items: n(e.items)
                            } : e);
                    return {
                        draftRules: n(e.draftRules)
                    };
                }),
            saveItemAtPath: (t)=>{
                e((e)=>{
                    let n = _(b(e.rules)), r = y(_(b(e.draftRules)), t), i = v(n, t), a = t[t.length - 1];
                    return i.items[a] = JSON.parse(JSON.stringify(r)), {
                        rules: n.items
                    };
                }), q({
                    type: `success`,
                    domain: `rules`,
                    message: `Rule saved`
                });
            },
            resetItemAtPath: (t)=>e((e)=>{
                    let n = b(e.draftRules), r = y(_(b(e.rules)), t), i = v(_(n), t), a = t[t.length - 1];
                    return i.items[a] = JSON.parse(JSON.stringify(r)), {
                        draftRules: n
                    };
                }),
            updateRuleAtPath: (t, n)=>e((e)=>{
                    X(e.draftRules);
                    let r = b(e.draftRules), i = _(r), a = y(i, t);
                    if (!g(a)) {
                        let e = v(i, t), r = t[t.length - 1];
                        e.items[r] = {
                            ...a,
                            ...n
                        };
                    }
                    return {
                        draftRules: r,
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                }),
            cloneItemAtPath: (t)=>e((e)=>{
                    X(e.draftRules);
                    let n = b(e.draftRules), r = _(n), i = Le(y(r, t)), a = v(r, t), o = t[t.length - 1];
                    return a.items.splice(o + 1, 0, i), {
                        draftRules: n,
                        canUndo: Y > 0,
                        canRedo: Y < J.length - 1
                    };
                })
        }), {
        name: `httptoolkit-rules-state`,
        partialize: (e)=>({
                rules: e.rules,
                draftRules: e.draftRules,
                upstreamProxyType: e.upstreamProxyType,
                upstreamProxyHost: e.upstreamProxyHost,
                upstreamNoProxyHosts: e.upstreamNoProxyHosts,
                whitelistedCertificateHosts: e.whitelistedCertificateHosts,
                clientCertificateHostMap: e.clientCertificateHostMap
            }),
        merge: (e, t)=>{
            let n = {
                ...t,
                ...e
            }, r = D(n.draftRules ?? T), i = D(n.rules ?? r);
            return {
                ...n,
                draftRules: r,
                rules: i,
                upstreamNoProxyHosts: Kt(n.upstreamNoProxyHosts)
            };
        },
        onRehydrateStorage: ()=>(e)=>{
                e && (e.draftRules = D(e.draftRules), e.rules = D(e.rules), setTimeout(()=>{
                    let { serverOnline: e } = $.getState();
                    e && p() && Z.getState().saveRules({
                        silent: !0
                    });
                }, 0));
            }
    }));
    o();
    let Q;
    Q = null;
    $ = n((e, t)=>({
            isInitialized: !1,
            serverOnline: !1,
            serverVersion: `1.0.0`,
            certPath: ``,
            certContent: ``,
            certFingerprint: ``,
            certFiles: void 0,
            externalNetworkAddresses: [],
            systemProxyConfig: void 0,
            dnsServers: [],
            ruleParameterKeys: [],
            toolPaths: void 0,
            httpProxyPort: 8e3,
            activeSessionId: ``,
            streamDisconnected: !1,
            http2Enabled: `fallback`,
            portConfig: void 0,
            isSystemCertInstalled: null,
            certInstallMessage: null,
            javaVersions: [],
            certJustInstalled: !1,
            certStatusAlreadyNotified: !1,
            setHttp2Enabled: (t)=>e({
                    http2Enabled: t
                }),
            setPortConfig: (t)=>e({
                    portConfig: t
                }),
            setDnsServers: (t)=>e({
                    dnsServers: t
                }),
            setStreamDisconnected: (t)=>e({
                    streamDisconnected: t
                }),
            reinitializeSession: async ()=>{
                let { useEventsStore: n } = await a(async ()=>{
                    let { useEventsStore: e } = await import(`./useEventsStore-B6HXM7p5.js`).then(async (m)=>{
                        await m.__tla;
                        return m;
                    });
                    return {
                        useEventsStore: e
                    };
                }, []);
                n.getState().disconnectStream(), e({
                    isInitialized: !1,
                    streamDisconnected: !0
                }), await se(), await t().initialize(), t().serverOnline && n.getState().connectToEventStream();
            },
            initialize: async ()=>{
                ce((n)=>{
                    if (t().setStreamDisconnected(n), !n && p()) {
                        let n = ue();
                        n !== t().httpProxyPort && e({
                            httpProxyPort: n
                        });
                    }
                    a(async ()=>{
                        let { useEventsStore: e } = await import(`./useEventsStore-B6HXM7p5.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        });
                        return {
                            useEventsStore: e
                        };
                    }, []).then(({ useEventsStore: e })=>{
                        n ? e.setState((e)=>({
                                streamConnected: !1,
                                reconnecting: !0,
                                reconnectAttempt: e.reconnectAttempt + 1
                            })) : (e.setState({
                            reconnecting: !1
                        }), e.getState().connectToEventStream());
                    });
                });
                let { isInitialized: n, serverOnline: r } = t();
                if (n && r) {
                    if (p()) {
                        let n = ue();
                        n !== t().httpProxyPort && e({
                            httpProxyPort: n
                        });
                    }
                    return;
                }
                if (Q) return Q;
                Q = (async ()=>{
                    let n = 8e3, { http2Enabled: r, portConfig: i, serverOnline: a } = t();
                    !a && p() && await se();
                    try {
                        let { httpProxyPort: e } = await ae({
                            http2Enabled: r,
                            portConfig: i
                        });
                        n = e;
                    } catch (t) {
                        u(t, `Failed to start Mockttp via PluggableAdmin`), e({
                            isInitialized: !0,
                            serverOnline: !1,
                            serverVersion: `Offline`,
                            certPath: ``,
                            certContent: ``,
                            certFingerprint: ``,
                            certFiles: void 0,
                            httpProxyPort: n,
                            activeSessionId: ``,
                            streamDisconnected: !0
                        });
                        return;
                    }
                    try {
                        let { version: r } = await m(`GET`, `/version`), { config: i } = await m(`GET`, `/config`, {
                            query: {
                                proxyPort: n
                            }
                        });
                        ie(i.authToken);
                        let { networkInterfaces: a } = await m(`GET`, `/config/network-interfaces`), o = [];
                        Object.entries(a || {}).forEach(([e, t])=>{
                            e !== `docker0` && !e.startsWith(`br-`) && !e.startsWith(`veth`) && Array.isArray(t) && t.forEach((e)=>{
                                e.family === `IPv4` && !e.internal && e.address && o.push(e.address);
                            });
                        }), e({
                            isInitialized: !0,
                            serverOnline: !0,
                            serverVersion: r,
                            certPath: i.certificatePath,
                            certContent: i.certificateContent || ``,
                            certFingerprint: i.certificateFingerprint || ``,
                            certFiles: i.certificateFiles,
                            systemProxyConfig: i.systemProxy,
                            dnsServers: S(i.dnsServers || []),
                            ruleParameterKeys: i.ruleParameterKeys || [],
                            toolPaths: i.toolPaths,
                            externalNetworkAddresses: o,
                            httpProxyPort: n,
                            activeSessionId: `pluggable-admin`,
                            streamDisconnected: !1
                        }), Z.getState().saveRules({
                            silent: !0
                        }), t().checkSystemCertStatus(), xe();
                    } catch (t) {
                        u(t, `Failed to initialize proxy store`), p() && await se().catch(()=>void 0), e({
                            isInitialized: !0,
                            serverOnline: !1,
                            serverVersion: `Offline`,
                            certPath: ``,
                            certContent: ``,
                            certFingerprint: ``,
                            certFiles: void 0,
                            httpProxyPort: n,
                            activeSessionId: ``,
                            streamDisconnected: !0
                        });
                    }
                })();
                try {
                    await Q;
                } finally{
                    Q = null;
                }
            },
            refreshConfig: async ()=>{
                try {
                    let { config: t } = await m(`GET`, `/config`, {
                        query: {
                            proxyPort: x()
                        }
                    });
                    e({
                        certPath: t.certificatePath,
                        certContent: t.certificateContent || ``,
                        certFingerprint: t.certificateFingerprint || ``,
                        certFiles: t.certificateFiles,
                        systemProxyConfig: t.systemProxy,
                        dnsServers: S(t.dnsServers || []),
                        ruleParameterKeys: t.ruleParameterKeys || [],
                        toolPaths: t.toolPaths
                    });
                } catch (e) {
                    u(e, `Failed to refresh proxy config`);
                }
            },
            checkSystemCertStatus: async (n = !1)=>{
                try {
                    let r = await ne(), i = t().certJustInstalled, o = t().certStatusAlreadyNotified;
                    if (e({
                        isSystemCertInstalled: r.installed,
                        certInstallMessage: null
                    }), r.installed) {
                        let { useNotificationStore: t } = await a(async ()=>{
                            let { useNotificationStore: e } = await import(`./useNotificationStore-NSh_n0Ie.js`).then(async (m)=>{
                                await m.__tla;
                                return m;
                            });
                            return {
                                useNotificationStore: e
                            };
                        }, __vite__mapDeps([0,1,2,3,4]));
                        i ? (t.getState().addNotification(`success`, `HttpToolkit Pro CA Certificate has been successfully installed now! Please restart your browser to capture traffic.`, 8e3), e({
                            certJustInstalled: !1,
                            certStatusAlreadyNotified: !0
                        })) : (n || !o) && (t.getState().addNotification(`success`, `HttpToolkit Pro CA Certificate is already installed and trusted.`, 5e3), e({
                            certStatusAlreadyNotified: !0
                        }));
                    }
                } catch (t) {
                    u(t, `Failed to check certificate status`), e({
                        isSystemCertInstalled: null
                    });
                }
            },
            installSystemCert: async ()=>{
                try {
                    let n = x();
                    e({
                        certJustInstalled: !0
                    });
                    let r = await fe(n);
                    await t().checkSystemCertStatus();
                    let i = r.message || (r.installed ? `CA installed — restart your browser, then capture traffic again.` : `Install finished but CA trust was not detected. Use manual install or run as Administrator.`);
                    e({
                        certInstallMessage: i
                    });
                    let { useNotificationStore: o } = await a(async ()=>{
                        let { useNotificationStore: e } = await import(`./useNotificationStore-NSh_n0Ie.js`).then(async (m)=>{
                            await m.__tla;
                            return m;
                        });
                        return {
                            useNotificationStore: e
                        };
                    }, __vite__mapDeps([0,1,2,3,4]));
                    r.installed || o.getState().addNotification(`warning`, i, 8e3);
                } catch (t) {
                    e({
                        certInstallMessage: t instanceof Error ? t.message : `Failed to install certificate`
                    }), u(t, `Failed to install certificate`);
                }
            },
            loadJavaVersions: async ()=>{
                try {
                    let n = (await m(`GET`, `/java/versions`)).versions || [];
                    e({
                        javaVersions: await Promise.all(n.map(async (e)=>({
                                ...e,
                                certInstalled: await t().checkJavaCertStatus(e.javaPath)
                            })))
                    });
                } catch (t) {
                    u(t, `Failed to load Java versions`), e({
                        javaVersions: []
                    });
                }
            },
            checkJavaCertStatus: async (e)=>{
                try {
                    return !!(await m(`GET`, `/java/certificate/status`, {
                        query: {
                            javaPath: e
                        }
                    })).installed;
                } catch  {
                    return !1;
                }
            },
            installJavaCert: async (n)=>{
                try {
                    await m(`POST`, `/java/certificate/install`, {
                        query: {
                            proxyPort: x()
                        },
                        body: {
                            javaPath: n
                        }
                    });
                    let r = await t().checkJavaCertStatus(n);
                    return e((e)=>({
                            javaVersions: e.javaVersions.map((e)=>e.javaPath === n ? {
                                    ...e,
                                    certInstalled: r
                                } : e)
                        })), await t().loadJavaVersions(), {
                        success: !0
                    };
                } catch (e) {
                    return u(e, `Failed to install certificate for Java: ${n}`), {
                        success: !1,
                        error: e instanceof Error ? e.message : String(e)
                    };
                }
            }
        }));
});
export { y as $, Er as A, Wt as B, Mr as C, Cr as D, ii as E, Nn as F, gt as G, Ct as H, Mn as I, x as J, S as K, fn as L, ui as M, V as N, li as O, z as P, Re as Q, dn as R, jr as S, ai as T, O as U, Gt as V, _t as W, _ as X, Ue as Y, Fe as Z, oi as _, gi as a, Wr as b, Xr as c, Rr as d, g as et, zr as f, yr as g, Ar as h, _i as i, Tr as j, Ir as k, Zr as l, mr as m, Z as n, hi as o, Vr as p, He as q, q as r, dr as s, $ as t, Lr as u, ci as v, Qr as w, Br as x, Ur as y, Vt as z, __tla };
