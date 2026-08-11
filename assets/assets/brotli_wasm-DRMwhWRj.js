import { n as e } from "./chunk-C0phanGV.js";
let M, j, N, P, y, oe, v, _, h, m, ae, b, g, x, re, ie;
let __tla = (async ()=>{
    var t, n = e((()=>{
        t = `/assets/brotli_wasm_bg-B6j4qCtc.wasm`;
    })), r, i = e((()=>{
        r = async (e = {}, t)=>{
            let n;
            if (t.startsWith(`data:`)) {
                let r = t.replace(/^data:.*?base64,/, ``), i;
                if (typeof Buffer == `function` && typeof Buffer.from == `function`) i = Buffer.from(r, `base64`);
                else if (typeof atob == `function`) {
                    let e = atob(r);
                    i = new Uint8Array(e.length);
                    for(let t = 0; t < e.length; t++)i[t] = e.charCodeAt(t);
                } else throw Error(`Cannot decode base64-encoded data URL`);
                n = await WebAssembly.instantiate(i, e);
            } else {
                let r = await fetch(t), i = r.headers.get(`Content-Type`) || ``;
                if (`instantiateStreaming` in WebAssembly && i.startsWith(`application/wasm`)) n = await WebAssembly.instantiateStreaming(r, e);
                else {
                    let t = await r.arrayBuffer();
                    n = await WebAssembly.instantiate(t, e);
                }
            }
            return n.instance.exports;
        };
    }));
    function a(e) {
        return S[e];
    }
    function o() {
        return (w === null || w.buffer !== L.buffer) && (w = new Uint8Array(L.buffer)), w;
    }
    function s(e, t) {
        return C.decode(o().subarray(e, e + t));
    }
    function c(e) {
        T === S.length && S.push(S.length + 1);
        let t = T;
        return T = S[t], S[t] = e, t;
    }
    function ee(e, t, n) {
        if (n === void 0) {
            let n = D.encode(e), r = t(n.length);
            return o().subarray(r, r + n.length).set(n), E = n.length, r;
        }
        let r = e.length, i = t(r), a = o(), s = 0;
        for(; s < r; s++){
            let t = e.charCodeAt(s);
            if (t > 127) break;
            a[i + s] = t;
        }
        if (s !== r) {
            s !== 0 && (e = e.slice(s)), i = n(i, r, r = s + e.length * 3);
            let t = o().subarray(i + s, i + r), a = O(e, t);
            s += a.written;
        }
        return E = s, i;
    }
    function l() {
        return (k === null || k.buffer !== L.buffer) && (k = new Int32Array(L.buffer)), k;
    }
    function te(e) {
        e < 36 || (S[e] = T, T = e);
    }
    function u(e) {
        let t = a(e);
        return te(e), t;
    }
    function d(e, t) {
        let n = t(e.length * 1);
        return o().set(e, n / 1), E = e.length, n;
    }
    function ne(e) {
        if (A == 1) throw Error(`out of js stack`);
        return S[--A] = e, A;
    }
    function f(e, t) {
        return o().subarray(e / 1, e / 1 + t);
    }
    re = function(e, t) {
        try {
            let o = Q(-16), s = d(e, X);
            se(o, s, E, ne(t));
            var n = l()[o / 4 + 0], r = l()[o / 4 + 1], i = l()[o / 4 + 2];
            if (l()[o / 4 + 3]) throw u(i);
            var a = f(n, r).slice();
            return $(n, r * 1), a;
        } finally{
            Q(16), S[A++] = void 0;
        }
    };
    ie = function(e) {
        try {
            let a = Q(-16), o = d(e, X);
            R(a, o, E);
            var t = l()[a / 4 + 0], n = l()[a / 4 + 1], r = l()[a / 4 + 2];
            if (l()[a / 4 + 3]) throw u(r);
            var i = f(t, n).slice();
            return $(t, n * 1), i;
        } finally{
            Q(16);
        }
    };
    function p(e) {
        return e == null;
    }
    m = function(e) {
        return a(e) === void 0;
    };
    h = function(e) {
        let t = a(e);
        return typeof t == `object` && !!t;
    };
    g = function(e, t) {
        return c(s(e, t));
    };
    _ = function(e, t) {
        return c(Error(s(e, t)));
    };
    ae = function(e, t) {
        let n = a(t), r = ee(JSON.stringify(n === void 0 ? null : n), X, Z), i = E;
        l()[e / 4 + 1] = i, l()[e / 4 + 0] = r;
    };
    oe = function() {
        return c(Error());
    };
    v = function(e, t) {
        let n = a(t).stack, r = ee(n, X, Z), i = E;
        l()[e / 4 + 1] = i, l()[e / 4 + 0] = r;
    };
    y = function(e, t) {
        try {
            console.error(s(e, t));
        } finally{
            $(e, t);
        }
    };
    b = function(e) {
        u(e);
    };
    x = function(e, t) {
        throw Error(s(e, t));
    };
    let S, C, w, T, E, D, O, k, A, F, I, L, se, R, z, B, V, H, U, W, G, K, q, J, Y, ce, le, ue, de, X, Z, Q, $, fe;
    F = e((async ()=>{
        await fe(), S = Array(32).fill(void 0), S.push(void 0, null, !0, !1), C = new (typeof TextDecoder > `u` ? (0, module.require)(`util`).TextDecoder : TextDecoder)(`utf-8`, {
            ignoreBOM: !0,
            fatal: !0
        }), C.decode(), w = null, T = S.length, E = 0, D = new (typeof TextEncoder > `u` ? (0, module.require)(`util`).TextEncoder : TextEncoder)(`utf-8`), O = typeof D.encodeInto == `function` ? function(e, t) {
            return D.encodeInto(e, t);
        } : function(e, t) {
            let n = D.encode(e);
            return t.set(n), {
                read: e.length,
                written: n.length
            };
        }, k = null, A = 32, j = Object.freeze({
            ResultSuccess: 1,
            1: `ResultSuccess`,
            NeedsMoreInput: 2,
            2: `NeedsMoreInput`,
            NeedsMoreOutput: 3,
            3: `NeedsMoreOutput`
        }), M = class e {
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
                z(e);
            }
            get code() {
                return B(this.ptr) >>> 0;
            }
            set code(e) {
                V(this.ptr, e);
            }
            get buf() {
                try {
                    let r = Q(-16);
                    H(r, this.ptr);
                    var e = l()[r / 4 + 0], t = l()[r / 4 + 1], n = f(e, t).slice();
                    return $(e, t * 1), n;
                } finally{
                    Q(16);
                }
            }
            set buf(e) {
                let t = d(e, X), n = E;
                U(this.ptr, t, n);
            }
            get input_offset() {
                return W(this.ptr) >>> 0;
            }
            set input_offset(e) {
                G(this.ptr, e);
            }
        }, N = class e {
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
                K(e);
            }
            constructor(t){
                let n = q(!p(t), p(t) ? 0 : t);
                return e.__wrap(n);
            }
            compress(e, t) {
                try {
                    let o = Q(-16);
                    var n = p(e) ? 0 : d(e, X), r = E;
                    J(o, this.ptr, n, r, t);
                    var i = l()[o / 4 + 0], a = l()[o / 4 + 1];
                    if (l()[o / 4 + 2]) throw u(a);
                    return M.__wrap(i);
                } finally{
                    Q(16);
                }
            }
            total_out() {
                return Y(this.ptr) >>> 0;
            }
        }, P = class e {
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
                ce(e);
            }
            constructor(){
                let t = le();
                return e.__wrap(t);
            }
            decompress(e, t) {
                try {
                    let i = Q(-16), a = d(e, X), o = E;
                    ue(i, this.ptr, a, o, t);
                    var n = l()[i / 4 + 0], r = l()[i / 4 + 1];
                    if (l()[i / 4 + 2]) throw u(r);
                    return M.__wrap(n);
                } finally{
                    Q(16);
                }
            }
            total_out() {
                return de(this.ptr) >>> 0;
            }
        };
    }));
    fe = e((async ()=>{
        n(), i(), await F(), URL = globalThis.URL, I = await r({
            "./brotli_wasm_bg.js": {
                __wbindgen_is_undefined: m,
                __wbindgen_is_object: h,
                __wbindgen_string_new: g,
                __wbindgen_error_new: _,
                __wbindgen_json_serialize: ae,
                __wbg_new_693216e109162396: oe,
                __wbg_stack_0ddaca5d1abfb52f: v,
                __wbg_error_09919627ac0992f5: y,
                __wbindgen_object_drop_ref: b,
                __wbindgen_throw: x
            }
        }, t), { memory: L, compress: se, decompress: R, __wbg_brotlistreamresult_free: z, __wbg_get_brotlistreamresult_code: B, __wbg_set_brotlistreamresult_code: V, __wbg_get_brotlistreamresult_buf: H, __wbg_set_brotlistreamresult_buf: U, __wbg_get_brotlistreamresult_input_offset: W, __wbg_set_brotlistreamresult_input_offset: G, __wbg_compressstream_free: K, compressstream_new: q, compressstream_compress: J, compressstream_total_out: Y, __wbg_decompressstream_free: ce, decompressstream_new: le, decompressstream_decompress: ue, decompressstream_total_out: de, __wbindgen_malloc: X, __wbindgen_realloc: Z, __wbindgen_add_to_stack_pointer: Q, __wbindgen_free: $ } = I;
    }));
    await e((async ()=>{
        await F();
    }))();
})();
export { M as BrotliStreamResult, j as BrotliStreamResultCode, N as CompressStream, P as DecompressStream, y as __wbg_error_09919627ac0992f5, oe as __wbg_new_693216e109162396, v as __wbg_stack_0ddaca5d1abfb52f, _ as __wbindgen_error_new, h as __wbindgen_is_object, m as __wbindgen_is_undefined, ae as __wbindgen_json_serialize, b as __wbindgen_object_drop_ref, g as __wbindgen_string_new, x as __wbindgen_throw, re as compress, ie as decompress, __tla };
