import { n as e } from "./chunk-CMxvf4Kt.js";
import { r as t, t as n } from "./dist-B3kALMnJ.js";
let P, N, F, I, ae, S, ie, b, v, _, x, oe, y, C, ne, re;
let __tla = (async ()=>{
    var r, i = e((()=>{
        r = `/assets/brotli_wasm_bg-B6j4qCtc.wasm`;
    })), a, o = e((()=>{
        t(), a = async (e = {}, t)=>{
            let r;
            if (t.startsWith(`data:`)) {
                let i = t.replace(/^data:.*?base64,/, ``), a;
                if (typeof n == `function` && typeof n.from == `function`) a = n.from(i, `base64`);
                else if (typeof atob == `function`) {
                    let e = atob(i);
                    a = new Uint8Array(e.length);
                    for(let t = 0; t < e.length; t++)a[t] = e.charCodeAt(t);
                } else throw Error(`Cannot decode base64-encoded data URL`);
                r = await WebAssembly.instantiate(a, e);
            } else {
                let n = await fetch(t), i = n.headers.get(`Content-Type`) || ``;
                if (`instantiateStreaming` in WebAssembly && i.startsWith(`application/wasm`)) r = await WebAssembly.instantiateStreaming(n, e);
                else {
                    let t = await n.arrayBuffer();
                    r = await WebAssembly.instantiate(t, e);
                }
            }
            return r.instance.exports;
        };
    }));
    function s(e) {
        return w[e];
    }
    function c() {
        return (E === null || E.buffer !== z.buffer) && (E = new Uint8Array(z.buffer)), E;
    }
    function l(e, t) {
        return T.decode(c().subarray(e, e + t));
    }
    function u(e) {
        D === w.length && w.push(w.length + 1);
        let t = D;
        return D = w[t], w[t] = e, t;
    }
    function d(e, t, n) {
        if (n === void 0) {
            let n = k.encode(e), r = t(n.length);
            return c().subarray(r, r + n.length).set(n), O = n.length, r;
        }
        let r = e.length, i = t(r), a = c(), o = 0;
        for(; o < r; o++){
            let t = e.charCodeAt(o);
            if (t > 127) break;
            a[i + o] = t;
        }
        if (o !== r) {
            o !== 0 && (e = e.slice(o)), i = n(i, r, r = o + e.length * 3);
            let t = c().subarray(i + o, i + r), a = A(e, t);
            o += a.written;
        }
        return O = o, i;
    }
    function f() {
        return (j === null || j.buffer !== z.buffer) && (j = new Int32Array(z.buffer)), j;
    }
    function ee(e) {
        e < 36 || (w[e] = D, D = e);
    }
    function p(e) {
        let t = s(e);
        return ee(e), t;
    }
    function m(e, t) {
        let n = t(e.length * 1);
        return c().set(e, n / 1), O = e.length, n;
    }
    function te(e) {
        if (M == 1) throw Error(`out of js stack`);
        return w[--M] = e, M;
    }
    function h(e, t) {
        return c().subarray(e / 1, e / 1 + t);
    }
    ne = function(e, t) {
        try {
            let o = Q(-16), s = m(e, X);
            se(o, s, O, te(t));
            var n = f()[o / 4 + 0], r = f()[o / 4 + 1], i = f()[o / 4 + 2];
            if (f()[o / 4 + 3]) throw p(i);
            var a = h(n, r).slice();
            return $(n, r * 1), a;
        } finally{
            Q(16), w[M++] = void 0;
        }
    };
    re = function(e) {
        try {
            let a = Q(-16), o = m(e, X);
            B(a, o, O);
            var t = f()[a / 4 + 0], n = f()[a / 4 + 1], r = f()[a / 4 + 2];
            if (f()[a / 4 + 3]) throw p(r);
            var i = h(t, n).slice();
            return $(t, n * 1), i;
        } finally{
            Q(16);
        }
    };
    function g(e) {
        return e == null;
    }
    _ = function(e) {
        return s(e) === void 0;
    };
    v = function(e) {
        let t = s(e);
        return typeof t == `object` && !!t;
    };
    y = function(e, t) {
        return u(l(e, t));
    };
    b = function(e, t) {
        return u(Error(l(e, t)));
    };
    x = function(e, t) {
        let n = s(t), r = d(JSON.stringify(n === void 0 ? null : n), X, Z), i = O;
        f()[e / 4 + 1] = i, f()[e / 4 + 0] = r;
    };
    S = function() {
        return u(Error());
    };
    ie = function(e, t) {
        let n = s(t).stack, r = d(n, X, Z), i = O;
        f()[e / 4 + 1] = i, f()[e / 4 + 0] = r;
    };
    ae = function(e, t) {
        try {
            console.error(l(e, t));
        } finally{
            $(e, t);
        }
    };
    oe = function(e) {
        p(e);
    };
    C = function(e, t) {
        throw Error(l(e, t));
    };
    let w, T, E, D, O, k, A, j, M, L, R, z, se, B, V, H, U, W, G, K, q, J, Y, ce, le, ue, de, fe, pe, X, Z, Q, $, me;
    L = e((async ()=>{
        await me(), w = Array(32).fill(void 0), w.push(void 0, null, !0, !1), T = new (typeof TextDecoder > `u` ? (0, module.require)(`util`).TextDecoder : TextDecoder)(`utf-8`, {
            ignoreBOM: !0,
            fatal: !0
        }), T.decode(), E = null, D = w.length, O = 0, k = new (typeof TextEncoder > `u` ? (0, module.require)(`util`).TextEncoder : TextEncoder)(`utf-8`), A = typeof k.encodeInto == `function` ? function(e, t) {
            return k.encodeInto(e, t);
        } : function(e, t) {
            let n = k.encode(e);
            return t.set(n), {
                read: e.length,
                written: n.length
            };
        }, j = null, M = 32, N = Object.freeze({
            ResultSuccess: 1,
            1: `ResultSuccess`,
            NeedsMoreInput: 2,
            2: `NeedsMoreInput`,
            NeedsMoreOutput: 3,
            3: `NeedsMoreOutput`
        }), P = class e {
            static __wrap(t) {
                let n = Object.create(e.prototype);
                return n.ptr = t, n;
            }
            __destroy_into_raw() {
                let e = this.ptr;
                return this.ptr = 0, e;
            }
            free() {
                let e = this.__destroy_into_raw();
                V(e);
            }
            get code() {
                return H(this.ptr) >>> 0;
            }
            set code(e) {
                U(this.ptr, e);
            }
            get buf() {
                try {
                    let r = Q(-16);
                    W(r, this.ptr);
                    var e = f()[r / 4 + 0], t = f()[r / 4 + 1], n = h(e, t).slice();
                    return $(e, t * 1), n;
                } finally{
                    Q(16);
                }
            }
            set buf(e) {
                let t = m(e, X), n = O;
                G(this.ptr, t, n);
            }
            get input_offset() {
                return K(this.ptr) >>> 0;
            }
            set input_offset(e) {
                q(this.ptr, e);
            }
        }, F = class e {
            static __wrap(t) {
                let n = Object.create(e.prototype);
                return n.ptr = t, n;
            }
            __destroy_into_raw() {
                let e = this.ptr;
                return this.ptr = 0, e;
            }
            free() {
                let e = this.__destroy_into_raw();
                J(e);
            }
            constructor(t){
                let n = Y(!g(t), g(t) ? 0 : t);
                return e.__wrap(n);
            }
            compress(e, t) {
                try {
                    let o = Q(-16);
                    var n = g(e) ? 0 : m(e, X), r = O;
                    ce(o, this.ptr, n, r, t);
                    var i = f()[o / 4 + 0], a = f()[o / 4 + 1];
                    if (f()[o / 4 + 2]) throw p(a);
                    return P.__wrap(i);
                } finally{
                    Q(16);
                }
            }
            total_out() {
                return le(this.ptr) >>> 0;
            }
        }, I = class e {
            static __wrap(t) {
                let n = Object.create(e.prototype);
                return n.ptr = t, n;
            }
            __destroy_into_raw() {
                let e = this.ptr;
                return this.ptr = 0, e;
            }
            free() {
                let e = this.__destroy_into_raw();
                ue(e);
            }
            constructor(){
                let t = de();
                return e.__wrap(t);
            }
            decompress(e, t) {
                try {
                    let i = Q(-16), a = m(e, X), o = O;
                    fe(i, this.ptr, a, o, t);
                    var n = f()[i / 4 + 0], r = f()[i / 4 + 1];
                    if (f()[i / 4 + 2]) throw p(r);
                    return P.__wrap(n);
                } finally{
                    Q(16);
                }
            }
            total_out() {
                return pe(this.ptr) >>> 0;
            }
        };
    }));
    me = e((async ()=>{
        i(), o(), await L(), URL = globalThis.URL, R = await a({
            "./brotli_wasm_bg.js": {
                __wbindgen_is_undefined: _,
                __wbindgen_is_object: v,
                __wbindgen_string_new: y,
                __wbindgen_error_new: b,
                __wbindgen_json_serialize: x,
                __wbg_new_693216e109162396: S,
                __wbg_stack_0ddaca5d1abfb52f: ie,
                __wbg_error_09919627ac0992f5: ae,
                __wbindgen_object_drop_ref: oe,
                __wbindgen_throw: C
            }
        }, r), { memory: z, compress: se, decompress: B, __wbg_brotlistreamresult_free: V, __wbg_get_brotlistreamresult_code: H, __wbg_set_brotlistreamresult_code: U, __wbg_get_brotlistreamresult_buf: W, __wbg_set_brotlistreamresult_buf: G, __wbg_get_brotlistreamresult_input_offset: K, __wbg_set_brotlistreamresult_input_offset: q, __wbg_compressstream_free: J, compressstream_new: Y, compressstream_compress: ce, compressstream_total_out: le, __wbg_decompressstream_free: ue, decompressstream_new: de, decompressstream_decompress: fe, decompressstream_total_out: pe, __wbindgen_malloc: X, __wbindgen_realloc: Z, __wbindgen_add_to_stack_pointer: Q, __wbindgen_free: $ } = R;
    }));
    await e((async ()=>{
        await L();
    }))();
})();
export { P as BrotliStreamResult, N as BrotliStreamResultCode, F as CompressStream, I as DecompressStream, ae as __wbg_error_09919627ac0992f5, S as __wbg_new_693216e109162396, ie as __wbg_stack_0ddaca5d1abfb52f, b as __wbindgen_error_new, v as __wbindgen_is_object, _ as __wbindgen_is_undefined, x as __wbindgen_json_serialize, oe as __wbindgen_object_drop_ref, y as __wbindgen_string_new, C as __wbindgen_throw, ne as compress, re as decompress, __tla };
