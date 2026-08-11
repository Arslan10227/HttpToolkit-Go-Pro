import { r as e, t } from "./chunk-C0phanGV.js";
import { n, r, t as i } from "./helpers-C-xDqguz.js";
(async ()=>{
    const a = [
        EvalError,
        RangeError,
        ReferenceError,
        SyntaxError,
        TypeError,
        URIError,
        globalThis.DOMException,
        globalThis.AssertionError,
        globalThis.SystemError
    ].filter(Boolean).map((e)=>[
            e.name,
            e
        ]), o = new Map(a), s = [
        {
            property: `name`,
            enumerable: !1
        },
        {
            property: `message`,
            enumerable: !1
        },
        {
            property: `stack`,
            enumerable: !1
        },
        {
            property: `code`,
            enumerable: !0
        },
        {
            property: `cause`,
            enumerable: !1
        }
    ], c = new WeakSet, l = (e)=>{
        c.add(e);
        let t = e.toJSON();
        return c.delete(e), t;
    }, u = (e)=>o.get(e) ?? Error, d = ({ from: e, seen: t, to: n, forceEnumerable: r, maxDepth: i, depth: a, useToJSON: o, serialize: f })=>{
        if (n ||= Array.isArray(e) ? [] : !f && p(e) ? new (u(e.name)) : {}, t.push(e), a >= i) return n;
        if (o && typeof e.toJSON == `function` && !c.has(e)) return l(e);
        let m = (e)=>d({
                from: e,
                seen: [
                    ...t
                ],
                forceEnumerable: r,
                maxDepth: i,
                depth: a,
                useToJSON: o,
                serialize: f
            });
        for (let [r, i] of Object.entries(e)){
            if (i && i instanceof Uint8Array && i.constructor.name === `Buffer`) {
                n[r] = `[object Buffer]`;
                continue;
            }
            if (typeof i == `object` && i && typeof i.pipe == `function`) {
                n[r] = `[object Stream]`;
                continue;
            }
            if (typeof i != `function`) {
                if (!i || typeof i != `object`) {
                    try {
                        n[r] = i;
                    } catch  {}
                    continue;
                }
                if (!t.includes(e[r])) {
                    a++, n[r] = m(e[r]);
                    continue;
                }
                n[r] = `[Circular]`;
            }
        }
        for (let { property: t, enumerable: i } of s)e[t] !== void 0 && e[t] !== null && Object.defineProperty(n, t, {
            value: p(e[t]) ? m(e[t]) : e[t],
            enumerable: r ? !0 : i,
            configurable: !0,
            writable: !0
        });
        return n;
    };
    function f(e, t = {}) {
        let { maxDepth: n = 1 / 0, useToJSON: r = !0 } = t;
        return typeof e == `object` && e ? d({
            from: e,
            seen: [],
            forceEnumerable: !0,
            maxDepth: n,
            depth: 0,
            useToJSON: r,
            serialize: !0
        }) : typeof e == `function` ? `[Function: ${e.name || `anonymous`}]` : e;
    }
    function p(e) {
        return !!e && typeof e == `object` && `name` in e && `message` in e && `stack` in e;
    }
    var m = t(((e)=>{
        Object.defineProperties(e, {
            __esModule: {
                value: !0
            },
            [Symbol.toStringTag]: {
                value: `Module`
            }
        });
        var t = {}, n = {};
        n.byteLength = u, n.toByteArray = f, n.fromByteArray = h;
        for(var r = [], i = [], a = typeof Uint8Array < `u` ? Uint8Array : Array, o = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/`, s = 0, c = o.length; s < c; ++s)r[s] = o[s], i[o.charCodeAt(s)] = s;
        i[45] = 62, i[95] = 63;
        function l(e) {
            var t = e.length;
            if (t % 4 > 0) throw Error(`Invalid string. Length must be a multiple of 4`);
            var n = e.indexOf(`=`);
            n === -1 && (n = t);
            var r = n === t ? 0 : 4 - n % 4;
            return [
                n,
                r
            ];
        }
        function u(e) {
            var t = l(e), n = t[0], r = t[1];
            return (n + r) * 3 / 4 - r;
        }
        function d(e, t, n) {
            return (t + n) * 3 / 4 - n;
        }
        function f(e) {
            var t, n = l(e), r = n[0], o = n[1], s = new a(d(e, r, o)), c = 0, u = o > 0 ? r - 4 : r, f;
            for(f = 0; f < u; f += 4)t = i[e.charCodeAt(f)] << 18 | i[e.charCodeAt(f + 1)] << 12 | i[e.charCodeAt(f + 2)] << 6 | i[e.charCodeAt(f + 3)], s[c++] = t >> 16 & 255, s[c++] = t >> 8 & 255, s[c++] = t & 255;
            return o === 2 && (t = i[e.charCodeAt(f)] << 2 | i[e.charCodeAt(f + 1)] >> 4, s[c++] = t & 255), o === 1 && (t = i[e.charCodeAt(f)] << 10 | i[e.charCodeAt(f + 1)] << 4 | i[e.charCodeAt(f + 2)] >> 2, s[c++] = t >> 8 & 255, s[c++] = t & 255), s;
        }
        function p(e) {
            return r[e >> 18 & 63] + r[e >> 12 & 63] + r[e >> 6 & 63] + r[e & 63];
        }
        function m(e, t, n) {
            for(var r, i = [], a = t; a < n; a += 3)r = (e[a] << 16 & 16711680) + (e[a + 1] << 8 & 65280) + (e[a + 2] & 255), i.push(p(r));
            return i.join(``);
        }
        function h(e) {
            for(var t, n = e.length, i = n % 3, a = [], o = 16383, s = 0, c = n - i; s < c; s += o)a.push(m(e, s, s + o > c ? c : s + o));
            return i === 1 ? (t = e[n - 1], a.push(r[t >> 2] + r[t << 4 & 63] + `==`)) : i === 2 && (t = (e[n - 2] << 8) + e[n - 1], a.push(r[t >> 10] + r[t >> 4 & 63] + r[t << 2 & 63] + `=`)), a.join(``);
        }
        var g = {};
        g.read = function(e, t, n, r, i) {
            var a, o, s = i * 8 - r - 1, c = (1 << s) - 1, l = c >> 1, u = -7, d = n ? i - 1 : 0, f = n ? -1 : 1, p = e[t + d];
            for(d += f, a = p & (1 << -u) - 1, p >>= -u, u += s; u > 0; a = a * 256 + e[t + d], d += f, u -= 8);
            for(o = a & (1 << -u) - 1, a >>= -u, u += r; u > 0; o = o * 256 + e[t + d], d += f, u -= 8);
            if (a === 0) a = 1 - l;
            else if (a === c) return o ? NaN : (p ? -1 : 1) * (1 / 0);
            else o += 2 ** r, a -= l;
            return (p ? -1 : 1) * o * 2 ** (a - r);
        }, g.write = function(e, t, n, r, i, a) {
            var o, s, c, l = a * 8 - i - 1, u = (1 << l) - 1, d = u >> 1, f = i === 23 ? 2 ** -24 - 2 ** -77 : 0, p = r ? 0 : a - 1, m = r ? 1 : -1, h = +(t < 0 || t === 0 && 1 / t < 0);
            for(t = Math.abs(t), isNaN(t) || t === 1 / 0 ? (s = +!!isNaN(t), o = u) : (o = Math.floor(Math.log(t) / Math.LN2), t * (c = 2 ** -o) < 1 && (o--, c *= 2), o + d >= 1 ? t += f / c : t += f * 2 ** (1 - d), t * c >= 2 && (o++, c /= 2), o + d >= u ? (s = 0, o = u) : o + d >= 1 ? (s = (t * c - 1) * 2 ** i, o += d) : (s = t * 2 ** (d - 1) * 2 ** i, o = 0)); i >= 8; e[n + p] = s & 255, p += m, s /= 256, i -= 8);
            for(o = o << i | s, l += i; l > 0; e[n + p] = o & 255, p += m, o /= 256, l -= 8);
            e[n + p - m] |= h * 128;
        }, (function(e) {
            let t = n, r = g, i = typeof Symbol == `function` && typeof Symbol.for == `function` ? Symbol.for(`nodejs.util.inspect.custom`) : null;
            e.Buffer = d, e.SlowBuffer = C, e.INSPECT_MAX_BYTES = 50;
            let a = 2147483647;
            e.kMaxLength = a;
            let { Uint8Array: o, ArrayBuffer: s, SharedArrayBuffer: c } = globalThis;
            d.TYPED_ARRAY_SUPPORT = l(), !d.TYPED_ARRAY_SUPPORT && typeof console < `u` && typeof console.error == `function` && console.error("This browser lacks typed array (Uint8Array) support which is required by `buffer` v5.x. Use `buffer` v4.x if you require old browser support.");
            function l() {
                try {
                    let e = new o(1), t = {
                        foo: function() {
                            return 42;
                        }
                    };
                    return Object.setPrototypeOf(t, o.prototype), Object.setPrototypeOf(e, t), e.foo() === 42;
                } catch  {
                    return !1;
                }
            }
            Object.defineProperty(d.prototype, "parent", {
                enumerable: !0,
                get: function() {
                    if (d.isBuffer(this)) return this.buffer;
                }
            }), Object.defineProperty(d.prototype, "offset", {
                enumerable: !0,
                get: function() {
                    if (d.isBuffer(this)) return this.byteOffset;
                }
            });
            function u(e) {
                if (e > a) throw RangeError(`The value "` + e + `" is invalid for option "size"`);
                let t = new o(e);
                return Object.setPrototypeOf(t, d.prototype), t;
            }
            function d(e, t, n) {
                if (typeof e == `number`) {
                    if (typeof t == `string`) throw TypeError(`The "string" argument must be of type string. Received type number`);
                    return h(e);
                }
                return f(e, t, n);
            }
            d.poolSize = 8192;
            function f(e, t, n) {
                if (typeof e == `string`) return _(e, t);
                if (s.isView(e)) return y(e);
                if (e == null) throw TypeError(`The first argument must be one of type string, Buffer, ArrayBuffer, Array, or Array-like Object. Received type ` + typeof e);
                if (de(e, s) || e && de(e.buffer, s) || c !== void 0 && (de(e, c) || e && de(e.buffer, c))) return b(e, t, n);
                if (typeof e == `number`) throw TypeError(`The "value" argument must not be of type number. Received type number`);
                let r = e.valueOf && e.valueOf();
                if (r != null && r !== e) return d.from(r, t, n);
                let i = x(e);
                if (i) return i;
                if (typeof Symbol < `u` && Symbol.toPrimitive != null && typeof e[Symbol.toPrimitive] == `function`) return d.from(e[Symbol.toPrimitive](`string`), t, n);
                throw TypeError(`The first argument must be one of type string, Buffer, ArrayBuffer, Array, or Array-like Object. Received type ` + typeof e);
            }
            d.from = function(e, t, n) {
                return f(e, t, n);
            }, Object.setPrototypeOf(d.prototype, o.prototype), Object.setPrototypeOf(d, o);
            function p(e) {
                if (typeof e != `number`) throw TypeError(`"size" argument must be of type number`);
                if (e < 0) throw RangeError(`The value "` + e + `" is invalid for option "size"`);
            }
            function m(e, t, n) {
                return p(e), e <= 0 || t === void 0 ? u(e) : typeof n == `string` ? u(e).fill(t, n) : u(e).fill(t);
            }
            d.alloc = function(e, t, n) {
                return m(e, t, n);
            };
            function h(e) {
                return p(e), u(e < 0 ? 0 : S(e) | 0);
            }
            d.allocUnsafe = function(e) {
                return h(e);
            }, d.allocUnsafeSlow = function(e) {
                return h(e);
            };
            function _(e, t) {
                if ((typeof t != `string` || t === ``) && (t = `utf8`), !d.isEncoding(t)) throw TypeError(`Unknown encoding: ` + t);
                let n = w(e, t) | 0, r = u(n), i = r.write(e, t);
                return i !== n && (r = r.slice(0, i)), r;
            }
            function v(e) {
                let t = e.length < 0 ? 0 : S(e.length) | 0, n = u(t);
                for(let r = 0; r < t; r += 1)n[r] = e[r] & 255;
                return n;
            }
            function y(e) {
                if (de(e, o)) {
                    let t = new o(e);
                    return b(t.buffer, t.byteOffset, t.byteLength);
                }
                return v(e);
            }
            function b(e, t, n) {
                if (t < 0 || e.byteLength < t) throw RangeError(`"offset" is outside of buffer bounds`);
                if (e.byteLength < t + (n || 0)) throw RangeError(`"length" is outside of buffer bounds`);
                let r;
                return r = t === void 0 && n === void 0 ? new o(e) : n === void 0 ? new o(e, t) : new o(e, t, n), Object.setPrototypeOf(r, d.prototype), r;
            }
            function x(e) {
                if (d.isBuffer(e)) {
                    let t = S(e.length) | 0, n = u(t);
                    return n.length === 0 || e.copy(n, 0, 0, t), n;
                }
                if (e.length !== void 0) return typeof e.length != `number` || Q(e.length) ? u(0) : v(e);
                if (e.type === `Buffer` && Array.isArray(e.data)) return v(e.data);
            }
            function S(e) {
                if (e >= a) throw RangeError(`Attempt to allocate Buffer larger than maximum size: 0x` + a.toString(16) + ` bytes`);
                return e | 0;
            }
            function C(e) {
                return +e != e && (e = 0), d.alloc(+e);
            }
            d.isBuffer = function(e) {
                return e != null && e._isBuffer === !0 && e !== d.prototype;
            }, d.compare = function(e, t) {
                if (de(e, o) && (e = d.from(e, e.offset, e.byteLength)), de(t, o) && (t = d.from(t, t.offset, t.byteLength)), !d.isBuffer(e) || !d.isBuffer(t)) throw TypeError(`The "buf1", "buf2" arguments must be one of type Buffer or Uint8Array`);
                if (e === t) return 0;
                let n = e.length, r = t.length;
                for(let i = 0, a = Math.min(n, r); i < a; ++i)if (e[i] !== t[i]) {
                    n = e[i], r = t[i];
                    break;
                }
                return n < r ? -1 : +(r < n);
            }, d.isEncoding = function(e) {
                switch(String(e).toLowerCase()){
                    case `hex`:
                    case `utf8`:
                    case `utf-8`:
                    case `ascii`:
                    case `latin1`:
                    case `binary`:
                    case `base64`:
                    case `ucs2`:
                    case `ucs-2`:
                    case `utf16le`:
                    case `utf-16le`:
                        return !0;
                    default:
                        return !1;
                }
            }, d.concat = function(e, t) {
                if (!Array.isArray(e)) throw TypeError(`"list" argument must be an Array of Buffers`);
                if (e.length === 0) return d.alloc(0);
                let n;
                if (t === void 0) for(t = 0, n = 0; n < e.length; ++n)t += e[n].length;
                let r = d.allocUnsafe(t), i = 0;
                for(n = 0; n < e.length; ++n){
                    let t = e[n];
                    if (de(t, o)) i + t.length > r.length ? (d.isBuffer(t) || (t = d.from(t)), t.copy(r, i)) : o.prototype.set.call(r, t, i);
                    else if (d.isBuffer(t)) t.copy(r, i);
                    else throw TypeError(`"list" argument must be an Array of Buffers`);
                    i += t.length;
                }
                return r;
            };
            function w(e, t) {
                if (d.isBuffer(e)) return e.length;
                if (s.isView(e) || de(e, s)) return e.byteLength;
                if (typeof e != `string`) throw TypeError(`The "string" argument must be one of type string, Buffer, or ArrayBuffer. Received type ` + typeof e);
                let n = e.length, r = arguments.length > 2 && arguments[2] === !0;
                if (!r && n === 0) return 0;
                let i = !1;
                for(;;)switch(t){
                    case `ascii`:
                    case `latin1`:
                    case `binary`:
                        return n;
                    case `utf8`:
                    case `utf-8`:
                        return ue(e).length;
                    case `ucs2`:
                    case `ucs-2`:
                    case `utf16le`:
                    case `utf-16le`:
                        return n * 2;
                    case `hex`:
                        return n >>> 1;
                    case `base64`:
                        return X(e).length;
                    default:
                        if (i) return r ? -1 : ue(e).length;
                        t = (`` + t).toLowerCase(), i = !0;
                }
            }
            d.byteLength = w;
            function T(e, t, n) {
                let r = !1;
                if ((t === void 0 || t < 0) && (t = 0), t > this.length || ((n === void 0 || n > this.length) && (n = this.length), n <= 0) || (n >>>= 0, t >>>= 0, n <= t)) return ``;
                for(e ||= `utf8`;;)switch(e){
                    case `hex`:
                        return re(this, t, n);
                    case `utf8`:
                    case `utf-8`:
                        return ee(this, t, n);
                    case `ascii`:
                        return I(this, t, n);
                    case `latin1`:
                    case `binary`:
                        return ne(this, t, n);
                    case `base64`:
                        return P(this, t, n);
                    case `ucs2`:
                    case `ucs-2`:
                    case `utf16le`:
                    case `utf-16le`:
                        return L(this, t, n);
                    default:
                        if (r) throw TypeError(`Unknown encoding: ` + e);
                        e = (e + ``).toLowerCase(), r = !0;
                }
            }
            d.prototype._isBuffer = !0;
            function E(e, t, n) {
                let r = e[t];
                e[t] = e[n], e[n] = r;
            }
            d.prototype.swap16 = function() {
                let e = this.length;
                if (e % 2 != 0) throw RangeError(`Buffer size must be a multiple of 16-bits`);
                for(let t = 0; t < e; t += 2)E(this, t, t + 1);
                return this;
            }, d.prototype.swap32 = function() {
                let e = this.length;
                if (e % 4 != 0) throw RangeError(`Buffer size must be a multiple of 32-bits`);
                for(let t = 0; t < e; t += 4)E(this, t, t + 3), E(this, t + 1, t + 2);
                return this;
            }, d.prototype.swap64 = function() {
                let e = this.length;
                if (e % 8 != 0) throw RangeError(`Buffer size must be a multiple of 64-bits`);
                for(let t = 0; t < e; t += 8)E(this, t, t + 7), E(this, t + 1, t + 6), E(this, t + 2, t + 5), E(this, t + 3, t + 4);
                return this;
            }, d.prototype.toString = function() {
                let e = this.length;
                return e === 0 ? `` : arguments.length === 0 ? ee(this, 0, e) : T.apply(this, arguments);
            }, d.prototype.toLocaleString = d.prototype.toString, d.prototype.equals = function(e) {
                if (!d.isBuffer(e)) throw TypeError(`Argument must be a Buffer`);
                return this === e ? !0 : d.compare(this, e) === 0;
            }, d.prototype.inspect = function() {
                let t = ``, n = e.INSPECT_MAX_BYTES;
                return t = this.toString(`hex`, 0, n).replace(/(.{2})/g, `$1 `).trim(), this.length > n && (t += ` ... `), `<Buffer ` + t + `>`;
            }, i && (d.prototype[i] = d.prototype.inspect), d.prototype.compare = function(e, t, n, r, i) {
                if (de(e, o) && (e = d.from(e, e.offset, e.byteLength)), !d.isBuffer(e)) throw TypeError(`The "target" argument must be one of type Buffer or Uint8Array. Received type ` + typeof e);
                if (t === void 0 && (t = 0), n === void 0 && (n = e ? e.length : 0), r === void 0 && (r = 0), i === void 0 && (i = this.length), t < 0 || n > e.length || r < 0 || i > this.length) throw RangeError(`out of range index`);
                if (r >= i && t >= n) return 0;
                if (r >= i) return -1;
                if (t >= n) return 1;
                if (t >>>= 0, n >>>= 0, r >>>= 0, i >>>= 0, this === e) return 0;
                let a = i - r, s = n - t, c = Math.min(a, s), l = this.slice(r, i), u = e.slice(t, n);
                for(let e = 0; e < c; ++e)if (l[e] !== u[e]) {
                    a = l[e], s = u[e];
                    break;
                }
                return a < s ? -1 : +(s < a);
            };
            function D(e, t, n, r, i) {
                if (e.length === 0) return -1;
                if (typeof n == `string` ? (r = n, n = 0) : n > 2147483647 ? n = 2147483647 : n < -2147483648 && (n = -2147483648), n = +n, Q(n) && (n = i ? 0 : e.length - 1), n < 0 && (n = e.length + n), n >= e.length) {
                    if (i) return -1;
                    n = e.length - 1;
                } else if (n < 0) if (i) n = 0;
                else return -1;
                if (typeof t == `string` && (t = d.from(t, r)), d.isBuffer(t)) return t.length === 0 ? -1 : O(e, t, n, r, i);
                if (typeof t == `number`) return t &= 255, typeof o.prototype.indexOf == `function` ? i ? o.prototype.indexOf.call(e, t, n) : o.prototype.lastIndexOf.call(e, t, n) : O(e, [
                    t
                ], n, r, i);
                throw TypeError(`val must be string, number or Buffer`);
            }
            function O(e, t, n, r, i) {
                let a = 1, o = e.length, s = t.length;
                if (r !== void 0 && (r = String(r).toLowerCase(), r === `ucs2` || r === `ucs-2` || r === `utf16le` || r === `utf-16le`)) {
                    if (e.length < 2 || t.length < 2) return -1;
                    a = 2, o /= 2, s /= 2, n /= 2;
                }
                function c(e, t) {
                    return a === 1 ? e[t] : e.readUInt16BE(t * a);
                }
                let l;
                if (i) {
                    let r = -1;
                    for(l = n; l < o; l++)if (c(e, l) === c(t, r === -1 ? 0 : l - r)) {
                        if (r === -1 && (r = l), l - r + 1 === s) return r * a;
                    } else r !== -1 && (l -= l - r), r = -1;
                } else for(n + s > o && (n = o - s), l = n; l >= 0; l--){
                    let n = !0;
                    for(let r = 0; r < s; r++)if (c(e, l + r) !== c(t, r)) {
                        n = !1;
                        break;
                    }
                    if (n) return l;
                }
                return -1;
            }
            d.prototype.includes = function(e, t, n) {
                return this.indexOf(e, t, n) !== -1;
            }, d.prototype.indexOf = function(e, t, n) {
                return D(this, e, t, n, !0);
            }, d.prototype.lastIndexOf = function(e, t, n) {
                return D(this, e, t, n, !1);
            };
            function k(e, t, n, r) {
                n = Number(n) || 0;
                let i = e.length - n;
                r ? (r = Number(r), r > i && (r = i)) : r = i;
                let a = t.length;
                r > a / 2 && (r = a / 2);
                let o;
                for(o = 0; o < r; ++o){
                    let r = parseInt(t.substr(o * 2, 2), 16);
                    if (Q(r)) return o;
                    e[n + o] = r;
                }
                return o;
            }
            function A(e, t, n, r) {
                return Z(ue(t, e.length - n), e, n, r);
            }
            function j(e, t, n, r) {
                return Z(J(t), e, n, r);
            }
            function M(e, t, n, r) {
                return Z(X(t), e, n, r);
            }
            function N(e, t, n, r) {
                return Z(Y(t, e.length - n), e, n, r);
            }
            d.prototype.write = function(e, t, n, r) {
                if (t === void 0) r = `utf8`, n = this.length, t = 0;
                else if (n === void 0 && typeof t == `string`) r = t, n = this.length, t = 0;
                else if (isFinite(t)) t >>>= 0, isFinite(n) ? (n >>>= 0, r === void 0 && (r = `utf8`)) : (r = n, n = void 0);
                else throw Error(`Buffer.write(string, encoding, offset[, length]) is no longer supported`);
                let i = this.length - t;
                if ((n === void 0 || n > i) && (n = i), e.length > 0 && (n < 0 || t < 0) || t > this.length) throw RangeError(`Attempt to write outside buffer bounds`);
                r ||= `utf8`;
                let a = !1;
                for(;;)switch(r){
                    case `hex`:
                        return k(this, e, t, n);
                    case `utf8`:
                    case `utf-8`:
                        return A(this, e, t, n);
                    case `ascii`:
                    case `latin1`:
                    case `binary`:
                        return j(this, e, t, n);
                    case `base64`:
                        return M(this, e, t, n);
                    case `ucs2`:
                    case `ucs-2`:
                    case `utf16le`:
                    case `utf-16le`:
                        return N(this, e, t, n);
                    default:
                        if (a) throw TypeError(`Unknown encoding: ` + r);
                        r = (`` + r).toLowerCase(), a = !0;
                }
            }, d.prototype.toJSON = function() {
                return {
                    type: `Buffer`,
                    data: Array.prototype.slice.call(this._arr || this, 0)
                };
            };
            function P(e, n, r) {
                return n === 0 && r === e.length ? t.fromByteArray(e) : t.fromByteArray(e.slice(n, r));
            }
            function ee(e, t, n) {
                n = Math.min(e.length, n);
                let r = [], i = t;
                for(; i < n;){
                    let t = e[i], a = null, o = t > 239 ? 4 : t > 223 ? 3 : t > 191 ? 2 : 1;
                    if (i + o <= n) {
                        let n, r, s, c;
                        switch(o){
                            case 1:
                                t < 128 && (a = t);
                                break;
                            case 2:
                                n = e[i + 1], (n & 192) == 128 && (c = (t & 31) << 6 | n & 63, c > 127 && (a = c));
                                break;
                            case 3:
                                n = e[i + 1], r = e[i + 2], (n & 192) == 128 && (r & 192) == 128 && (c = (t & 15) << 12 | (n & 63) << 6 | r & 63, c > 2047 && (c < 55296 || c > 57343) && (a = c));
                                break;
                            case 4:
                                n = e[i + 1], r = e[i + 2], s = e[i + 3], (n & 192) == 128 && (r & 192) == 128 && (s & 192) == 128 && (c = (t & 15) << 18 | (n & 63) << 12 | (r & 63) << 6 | s & 63, c > 65535 && c < 1114112 && (a = c));
                        }
                    }
                    a === null ? (a = 65533, o = 1) : a > 65535 && (a -= 65536, r.push(a >>> 10 & 1023 | 55296), a = 56320 | a & 1023), r.push(a), i += o;
                }
                return te(r);
            }
            let F = 4096;
            function te(e) {
                let t = e.length;
                if (t <= F) return String.fromCharCode.apply(String, e);
                let n = ``, r = 0;
                for(; r < t;)n += String.fromCharCode.apply(String, e.slice(r, r += F));
                return n;
            }
            function I(e, t, n) {
                let r = ``;
                n = Math.min(e.length, n);
                for(let i = t; i < n; ++i)r += String.fromCharCode(e[i] & 127);
                return r;
            }
            function ne(e, t, n) {
                let r = ``;
                n = Math.min(e.length, n);
                for(let i = t; i < n; ++i)r += String.fromCharCode(e[i]);
                return r;
            }
            function re(e, t, n) {
                let r = e.length;
                (!t || t < 0) && (t = 0), (!n || n < 0 || n > r) && (n = r);
                let i = ``;
                for(let r = t; r < n; ++r)i += fe[e[r]];
                return i;
            }
            function L(e, t, n) {
                let r = e.slice(t, n), i = ``;
                for(let e = 0; e < r.length - 1; e += 2)i += String.fromCharCode(r[e] + r[e + 1] * 256);
                return i;
            }
            d.prototype.slice = function(e, t) {
                let n = this.length;
                e = ~~e, t = t === void 0 ? n : ~~t, e < 0 ? (e += n, e < 0 && (e = 0)) : e > n && (e = n), t < 0 ? (t += n, t < 0 && (t = 0)) : t > n && (t = n), t < e && (t = e);
                let r = this.subarray(e, t);
                return Object.setPrototypeOf(r, d.prototype), r;
            };
            function R(e, t, n) {
                if (e % 1 != 0 || e < 0) throw RangeError(`offset is not uint`);
                if (e + t > n) throw RangeError(`Trying to access beyond buffer length`);
            }
            d.prototype.readUintLE = d.prototype.readUIntLE = function(e, t, n) {
                e >>>= 0, t >>>= 0, n || R(e, t, this.length);
                let r = this[e], i = 1, a = 0;
                for(; ++a < t && (i *= 256);)r += this[e + a] * i;
                return r;
            }, d.prototype.readUintBE = d.prototype.readUIntBE = function(e, t, n) {
                e >>>= 0, t >>>= 0, n || R(e, t, this.length);
                let r = this[e + --t], i = 1;
                for(; t > 0 && (i *= 256);)r += this[e + --t] * i;
                return r;
            }, d.prototype.readUint8 = d.prototype.readUInt8 = function(e, t) {
                return e >>>= 0, t || R(e, 1, this.length), this[e];
            }, d.prototype.readUint16LE = d.prototype.readUInt16LE = function(e, t) {
                return e >>>= 0, t || R(e, 2, this.length), this[e] | this[e + 1] << 8;
            }, d.prototype.readUint16BE = d.prototype.readUInt16BE = function(e, t) {
                return e >>>= 0, t || R(e, 2, this.length), this[e] << 8 | this[e + 1];
            }, d.prototype.readUint32LE = d.prototype.readUInt32LE = function(e, t) {
                return e >>>= 0, t || R(e, 4, this.length), (this[e] | this[e + 1] << 8 | this[e + 2] << 16) + this[e + 3] * 16777216;
            }, d.prototype.readUint32BE = d.prototype.readUInt32BE = function(e, t) {
                return e >>>= 0, t || R(e, 4, this.length), this[e] * 16777216 + (this[e + 1] << 16 | this[e + 2] << 8 | this[e + 3]);
            }, d.prototype.readBigUInt64LE = pe(function(e) {
                e >>>= 0, K(e, `offset`);
                let t = this[e], n = this[e + 7];
                (t === void 0 || n === void 0) && ce(e, this.length - 8);
                let r = t + this[++e] * 2 ** 8 + this[++e] * 2 ** 16 + this[++e] * 2 ** 24, i = this[++e] + this[++e] * 2 ** 8 + this[++e] * 2 ** 16 + n * 2 ** 24;
                return BigInt(r) + (BigInt(i) << BigInt(32));
            }), d.prototype.readBigUInt64BE = pe(function(e) {
                e >>>= 0, K(e, `offset`);
                let t = this[e], n = this[e + 7];
                (t === void 0 || n === void 0) && ce(e, this.length - 8);
                let r = t * 2 ** 24 + this[++e] * 2 ** 16 + this[++e] * 2 ** 8 + this[++e], i = this[++e] * 2 ** 24 + this[++e] * 2 ** 16 + this[++e] * 2 ** 8 + n;
                return (BigInt(r) << BigInt(32)) + BigInt(i);
            }), d.prototype.readIntLE = function(e, t, n) {
                e >>>= 0, t >>>= 0, n || R(e, t, this.length);
                let r = this[e], i = 1, a = 0;
                for(; ++a < t && (i *= 256);)r += this[e + a] * i;
                return i *= 128, r >= i && (r -= 2 ** (8 * t)), r;
            }, d.prototype.readIntBE = function(e, t, n) {
                e >>>= 0, t >>>= 0, n || R(e, t, this.length);
                let r = t, i = 1, a = this[e + --r];
                for(; r > 0 && (i *= 256);)a += this[e + --r] * i;
                return i *= 128, a >= i && (a -= 2 ** (8 * t)), a;
            }, d.prototype.readInt8 = function(e, t) {
                return e >>>= 0, t || R(e, 1, this.length), this[e] & 128 ? (255 - this[e] + 1) * -1 : this[e];
            }, d.prototype.readInt16LE = function(e, t) {
                e >>>= 0, t || R(e, 2, this.length);
                let n = this[e] | this[e + 1] << 8;
                return n & 32768 ? n | 4294901760 : n;
            }, d.prototype.readInt16BE = function(e, t) {
                e >>>= 0, t || R(e, 2, this.length);
                let n = this[e + 1] | this[e] << 8;
                return n & 32768 ? n | 4294901760 : n;
            }, d.prototype.readInt32LE = function(e, t) {
                return e >>>= 0, t || R(e, 4, this.length), this[e] | this[e + 1] << 8 | this[e + 2] << 16 | this[e + 3] << 24;
            }, d.prototype.readInt32BE = function(e, t) {
                return e >>>= 0, t || R(e, 4, this.length), this[e] << 24 | this[e + 1] << 16 | this[e + 2] << 8 | this[e + 3];
            }, d.prototype.readBigInt64LE = pe(function(e) {
                e >>>= 0, K(e, `offset`);
                let t = this[e], n = this[e + 7];
                (t === void 0 || n === void 0) && ce(e, this.length - 8);
                let r = this[e + 4] + this[e + 5] * 2 ** 8 + this[e + 6] * 2 ** 16 + (n << 24);
                return (BigInt(r) << BigInt(32)) + BigInt(t + this[++e] * 2 ** 8 + this[++e] * 2 ** 16 + this[++e] * 2 ** 24);
            }), d.prototype.readBigInt64BE = pe(function(e) {
                e >>>= 0, K(e, `offset`);
                let t = this[e], n = this[e + 7];
                (t === void 0 || n === void 0) && ce(e, this.length - 8);
                let r = (t << 24) + this[++e] * 2 ** 16 + this[++e] * 2 ** 8 + this[++e];
                return (BigInt(r) << BigInt(32)) + BigInt(this[++e] * 2 ** 24 + this[++e] * 2 ** 16 + this[++e] * 2 ** 8 + n);
            }), d.prototype.readFloatLE = function(e, t) {
                return e >>>= 0, t || R(e, 4, this.length), r.read(this, e, !0, 23, 4);
            }, d.prototype.readFloatBE = function(e, t) {
                return e >>>= 0, t || R(e, 4, this.length), r.read(this, e, !1, 23, 4);
            }, d.prototype.readDoubleLE = function(e, t) {
                return e >>>= 0, t || R(e, 8, this.length), r.read(this, e, !0, 52, 8);
            }, d.prototype.readDoubleBE = function(e, t) {
                return e >>>= 0, t || R(e, 8, this.length), r.read(this, e, !1, 52, 8);
            };
            function z(e, t, n, r, i, a) {
                if (!d.isBuffer(e)) throw TypeError(`"buffer" argument must be a Buffer instance`);
                if (t > i || t < a) throw RangeError(`"value" argument is out of bounds`);
                if (n + r > e.length) throw RangeError(`Index out of range`);
            }
            d.prototype.writeUintLE = d.prototype.writeUIntLE = function(e, t, n, r) {
                if (e = +e, t >>>= 0, n >>>= 0, !r) {
                    let r = 2 ** (8 * n) - 1;
                    z(this, e, t, n, r, 0);
                }
                let i = 1, a = 0;
                for(this[t] = e & 255; ++a < n && (i *= 256);)this[t + a] = e / i & 255;
                return t + n;
            }, d.prototype.writeUintBE = d.prototype.writeUIntBE = function(e, t, n, r) {
                if (e = +e, t >>>= 0, n >>>= 0, !r) {
                    let r = 2 ** (8 * n) - 1;
                    z(this, e, t, n, r, 0);
                }
                let i = n - 1, a = 1;
                for(this[t + i] = e & 255; --i >= 0 && (a *= 256);)this[t + i] = e / a & 255;
                return t + n;
            }, d.prototype.writeUint8 = d.prototype.writeUInt8 = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 1, 255, 0), this[t] = e & 255, t + 1;
            }, d.prototype.writeUint16LE = d.prototype.writeUInt16LE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 2, 65535, 0), this[t] = e & 255, this[t + 1] = e >>> 8, t + 2;
            }, d.prototype.writeUint16BE = d.prototype.writeUInt16BE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 2, 65535, 0), this[t] = e >>> 8, this[t + 1] = e & 255, t + 2;
            }, d.prototype.writeUint32LE = d.prototype.writeUInt32LE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 4, 4294967295, 0), this[t + 3] = e >>> 24, this[t + 2] = e >>> 16, this[t + 1] = e >>> 8, this[t] = e & 255, t + 4;
            }, d.prototype.writeUint32BE = d.prototype.writeUInt32BE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 4, 4294967295, 0), this[t] = e >>> 24, this[t + 1] = e >>> 16, this[t + 2] = e >>> 8, this[t + 3] = e & 255, t + 4;
            };
            function B(e, t, n, r, i) {
                G(t, r, i, e, n, 7);
                let a = Number(t & BigInt(4294967295));
                e[n++] = a, a >>= 8, e[n++] = a, a >>= 8, e[n++] = a, a >>= 8, e[n++] = a;
                let o = Number(t >> BigInt(32) & BigInt(4294967295));
                return e[n++] = o, o >>= 8, e[n++] = o, o >>= 8, e[n++] = o, o >>= 8, e[n++] = o, n;
            }
            function ie(e, t, n, r, i) {
                G(t, r, i, e, n, 7);
                let a = Number(t & BigInt(4294967295));
                e[n + 7] = a, a >>= 8, e[n + 6] = a, a >>= 8, e[n + 5] = a, a >>= 8, e[n + 4] = a;
                let o = Number(t >> BigInt(32) & BigInt(4294967295));
                return e[n + 3] = o, o >>= 8, e[n + 2] = o, o >>= 8, e[n + 1] = o, o >>= 8, e[n] = o, n + 8;
            }
            d.prototype.writeBigUInt64LE = pe(function(e, t = 0) {
                return B(this, e, t, BigInt(0), BigInt(`0xffffffffffffffff`));
            }), d.prototype.writeBigUInt64BE = pe(function(e, t = 0) {
                return ie(this, e, t, BigInt(0), BigInt(`0xffffffffffffffff`));
            }), d.prototype.writeIntLE = function(e, t, n, r) {
                if (e = +e, t >>>= 0, !r) {
                    let r = 2 ** (8 * n - 1);
                    z(this, e, t, n, r - 1, -r);
                }
                let i = 0, a = 1, o = 0;
                for(this[t] = e & 255; ++i < n && (a *= 256);)e < 0 && o === 0 && this[t + i - 1] !== 0 && (o = 1), this[t + i] = (e / a >> 0) - o & 255;
                return t + n;
            }, d.prototype.writeIntBE = function(e, t, n, r) {
                if (e = +e, t >>>= 0, !r) {
                    let r = 2 ** (8 * n - 1);
                    z(this, e, t, n, r - 1, -r);
                }
                let i = n - 1, a = 1, o = 0;
                for(this[t + i] = e & 255; --i >= 0 && (a *= 256);)e < 0 && o === 0 && this[t + i + 1] !== 0 && (o = 1), this[t + i] = (e / a >> 0) - o & 255;
                return t + n;
            }, d.prototype.writeInt8 = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 1, 127, -128), e < 0 && (e = 255 + e + 1), this[t] = e & 255, t + 1;
            }, d.prototype.writeInt16LE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 2, 32767, -32768), this[t] = e & 255, this[t + 1] = e >>> 8, t + 2;
            }, d.prototype.writeInt16BE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 2, 32767, -32768), this[t] = e >>> 8, this[t + 1] = e & 255, t + 2;
            }, d.prototype.writeInt32LE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 4, 2147483647, -2147483648), this[t] = e & 255, this[t + 1] = e >>> 8, this[t + 2] = e >>> 16, this[t + 3] = e >>> 24, t + 4;
            }, d.prototype.writeInt32BE = function(e, t, n) {
                return e = +e, t >>>= 0, n || z(this, e, t, 4, 2147483647, -2147483648), e < 0 && (e = 4294967295 + e + 1), this[t] = e >>> 24, this[t + 1] = e >>> 16, this[t + 2] = e >>> 8, this[t + 3] = e & 255, t + 4;
            }, d.prototype.writeBigInt64LE = pe(function(e, t = 0) {
                return B(this, e, t, -BigInt(`0x8000000000000000`), BigInt(`0x7fffffffffffffff`));
            }), d.prototype.writeBigInt64BE = pe(function(e, t = 0) {
                return ie(this, e, t, -BigInt(`0x8000000000000000`), BigInt(`0x7fffffffffffffff`));
            });
            function ae(e, t, n, r, i, a) {
                if (n + r > e.length || n < 0) throw RangeError(`Index out of range`);
            }
            function oe(e, t, n, i, a) {
                return t = +t, n >>>= 0, a || ae(e, t, n, 4), r.write(e, t, n, i, 23, 4), n + 4;
            }
            d.prototype.writeFloatLE = function(e, t, n) {
                return oe(this, e, t, !0, n);
            }, d.prototype.writeFloatBE = function(e, t, n) {
                return oe(this, e, t, !1, n);
            };
            function V(e, t, n, i, a) {
                return t = +t, n >>>= 0, a || ae(e, t, n, 8), r.write(e, t, n, i, 52, 8), n + 8;
            }
            d.prototype.writeDoubleLE = function(e, t, n) {
                return V(this, e, t, !0, n);
            }, d.prototype.writeDoubleBE = function(e, t, n) {
                return V(this, e, t, !1, n);
            }, d.prototype.copy = function(e, t, n, r) {
                if (!d.isBuffer(e)) throw TypeError(`argument should be a Buffer`);
                if (n ||= 0, !r && r !== 0 && (r = this.length), t >= e.length && (t = e.length), t ||= 0, r > 0 && r < n && (r = n), r === n || e.length === 0 || this.length === 0) return 0;
                if (t < 0) throw RangeError(`targetStart out of bounds`);
                if (n < 0 || n >= this.length) throw RangeError(`Index out of range`);
                if (r < 0) throw RangeError(`sourceEnd out of bounds`);
                r > this.length && (r = this.length), e.length - t < r - n && (r = e.length - t + n);
                let i = r - n;
                return this === e && typeof o.prototype.copyWithin == `function` ? this.copyWithin(t, n, r) : o.prototype.set.call(e, this.subarray(n, r), t), i;
            }, d.prototype.fill = function(e, t, n, r) {
                if (typeof e == `string`) {
                    if (typeof t == `string` ? (r = t, t = 0, n = this.length) : typeof n == `string` && (r = n, n = this.length), r !== void 0 && typeof r != `string`) throw TypeError(`encoding must be a string`);
                    if (typeof r == `string` && !d.isEncoding(r)) throw TypeError(`Unknown encoding: ` + r);
                    if (e.length === 1) {
                        let t = e.charCodeAt(0);
                        (r === `utf8` && t < 128 || r === `latin1`) && (e = t);
                    }
                } else typeof e == `number` ? e &= 255 : typeof e == `boolean` && (e = Number(e));
                if (t < 0 || this.length < t || this.length < n) throw RangeError(`Out of range index`);
                if (n <= t) return this;
                t >>>= 0, n = n === void 0 ? this.length : n >>> 0, e ||= 0;
                let i;
                if (typeof e == `number`) for(i = t; i < n; ++i)this[i] = e;
                else {
                    let a = d.isBuffer(e) ? e : d.from(e, r), o = a.length;
                    if (o === 0) throw TypeError(`The value "` + e + `" is invalid for argument "value"`);
                    for(i = 0; i < n - t; ++i)this[i + t] = a[i % o];
                }
                return this;
            };
            let H = {};
            function U(e, t, n) {
                H[e] = class extends n {
                    constructor(){
                        super(), Object.defineProperty(this, "message", {
                            value: t.apply(this, arguments),
                            writable: !0,
                            configurable: !0
                        }), this.name = `${this.name} [${e}]`, this.stack, delete this.name;
                    }
                    get code() {
                        return e;
                    }
                    set code(e) {
                        Object.defineProperty(this, "code", {
                            configurable: !0,
                            enumerable: !0,
                            value: e,
                            writable: !0
                        });
                    }
                    toString() {
                        return `${this.name} [${e}]: ${this.message}`;
                    }
                };
            }
            U(`ERR_BUFFER_OUT_OF_BOUNDS`, function(e) {
                return e ? `${e} is outside of buffer bounds` : `Attempt to access memory outside buffer bounds`;
            }, RangeError), U(`ERR_INVALID_ARG_TYPE`, function(e, t) {
                return `The "${e}" argument must be of type number. Received type ${typeof t}`;
            }, TypeError), U(`ERR_OUT_OF_RANGE`, function(e, t, n) {
                let r = `The value of "${e}" is out of range.`, i = n;
                return Number.isInteger(n) && Math.abs(n) > 2 ** 32 ? i = se(String(n)) : typeof n == `bigint` && (i = String(n), (n > BigInt(2) ** BigInt(32) || n < -(BigInt(2) ** BigInt(32))) && (i = se(i)), i += `n`), r += ` It must be ${t}. Received ${i}`, r;
            }, RangeError);
            function se(e) {
                let t = ``, n = e.length, r = +(e[0] === `-`);
                for(; n >= r + 4; n -= 3)t = `_${e.slice(n - 3, n)}${t}`;
                return `${e.slice(0, n)}${t}`;
            }
            function W(e, t, n) {
                K(t, `offset`), (e[t] === void 0 || e[t + n] === void 0) && ce(t, e.length - (n + 1));
            }
            function G(e, t, n, r, i, a) {
                if (e > n || e < t) {
                    let r = typeof t == `bigint` ? `n` : ``, i;
                    throw i = a > 3 ? t === 0 || t === BigInt(0) ? `>= 0${r} and < 2${r} ** ${(a + 1) * 8}${r}` : `>= -(2${r} ** ${(a + 1) * 8 - 1}${r}) and < 2 ** ${(a + 1) * 8 - 1}${r}` : `>= ${t}${r} and <= ${n}${r}`, new H.ERR_OUT_OF_RANGE(`value`, i, e);
                }
                W(r, i, a);
            }
            function K(e, t) {
                if (typeof e != `number`) throw new H.ERR_INVALID_ARG_TYPE(t, `number`, e);
            }
            function ce(e, t, n) {
                throw Math.floor(e) === e ? t < 0 ? new H.ERR_BUFFER_OUT_OF_BOUNDS : new H.ERR_OUT_OF_RANGE(n || `offset`, `>= ${+!!n} and <= ${t}`, e) : (K(e, n), new H.ERR_OUT_OF_RANGE(n || `offset`, `an integer`, e));
            }
            let le = /[^+/0-9A-Za-z-_]/g;
            function q(e) {
                if (e = e.split(`=`)[0], e = e.trim().replace(le, ``), e.length < 2) return ``;
                for(; e.length % 4 != 0;)e += `=`;
                return e;
            }
            function ue(e, t) {
                t ||= 1 / 0;
                let n, r = e.length, i = null, a = [];
                for(let o = 0; o < r; ++o){
                    if (n = e.charCodeAt(o), n > 55295 && n < 57344) {
                        if (!i) {
                            if (n > 56319) {
                                (t -= 3) > -1 && a.push(239, 191, 189);
                                continue;
                            } else if (o + 1 === r) {
                                (t -= 3) > -1 && a.push(239, 191, 189);
                                continue;
                            }
                            i = n;
                            continue;
                        }
                        if (n < 56320) {
                            (t -= 3) > -1 && a.push(239, 191, 189), i = n;
                            continue;
                        }
                        n = (i - 55296 << 10 | n - 56320) + 65536;
                    } else i && (t -= 3) > -1 && a.push(239, 191, 189);
                    if (i = null, n < 128) {
                        if (--t < 0) break;
                        a.push(n);
                    } else if (n < 2048) {
                        if ((t -= 2) < 0) break;
                        a.push(n >> 6 | 192, n & 63 | 128);
                    } else if (n < 65536) {
                        if ((t -= 3) < 0) break;
                        a.push(n >> 12 | 224, n >> 6 & 63 | 128, n & 63 | 128);
                    } else if (n < 1114112) {
                        if ((t -= 4) < 0) break;
                        a.push(n >> 18 | 240, n >> 12 & 63 | 128, n >> 6 & 63 | 128, n & 63 | 128);
                    } else throw Error(`Invalid code point`);
                }
                return a;
            }
            function J(e) {
                let t = [];
                for(let n = 0; n < e.length; ++n)t.push(e.charCodeAt(n) & 255);
                return t;
            }
            function Y(e, t) {
                let n, r, i, a = [];
                for(let o = 0; o < e.length && !((t -= 2) < 0); ++o)n = e.charCodeAt(o), r = n >> 8, i = n % 256, a.push(i), a.push(r);
                return a;
            }
            function X(e) {
                return t.toByteArray(q(e));
            }
            function Z(e, t, n, r) {
                let i;
                for(i = 0; i < r && !(i + n >= t.length || i >= e.length); ++i)t[i + n] = e[i];
                return i;
            }
            function de(e, t) {
                return e instanceof t || e != null && e.constructor != null && e.constructor.name != null && e.constructor.name === t.name;
            }
            function Q(e) {
                return e !== e;
            }
            let fe = (function() {
                let e = `0123456789abcdef`, t = Array(256);
                for(let n = 0; n < 16; ++n){
                    let r = n * 16;
                    for(let i = 0; i < 16; ++i)t[r + i] = e[n] + e[i];
                }
                return t;
            })();
            function pe(e) {
                return typeof BigInt > `u` ? $ : e;
            }
            function $() {
                throw Error(`BigInt not supported`);
            }
        })(t);
        let _ = t.Buffer;
        e.Blob = t.Blob, e.BlobOptions = t.BlobOptions, e.Buffer = t.Buffer, e.File = t.File, e.FileOptions = t.FileOptions, e.INSPECT_MAX_BYTES = t.INSPECT_MAX_BYTES, e.SlowBuffer = t.SlowBuffer, e.TranscodeEncoding = t.TranscodeEncoding, e.atob = t.atob, e.btoa = t.btoa, e.constants = t.constants, e.default = _, e.isAscii = t.isAscii, e.isUtf8 = t.isUtf8, e.kMaxLength = t.kMaxLength, e.kStringMaxLength = t.kStringMaxLength, e.resolveObjectURL = t.resolveObjectURL, e.transcode = t.transcode;
    })), h = t(((e, t)=>{
        var n = typeof Reflect == `object` ? Reflect : null, r = n && typeof n.apply == `function` ? n.apply : function(e, t, n) {
            return Function.prototype.apply.call(e, t, n);
        }, i = n && typeof n.ownKeys == `function` ? n.ownKeys : Object.getOwnPropertySymbols ? function(e) {
            return Object.getOwnPropertyNames(e).concat(Object.getOwnPropertySymbols(e));
        } : function(e) {
            return Object.getOwnPropertyNames(e);
        };
        function a(e) {
            console && console.warn && console.warn(e);
        }
        var o = Number.isNaN || function(e) {
            return e !== e;
        };
        function s() {
            s.init.call(this);
        }
        t.exports = s, t.exports.once = y, s.EventEmitter = s, s.prototype._events = void 0, s.prototype._eventsCount = 0, s.prototype._maxListeners = void 0;
        var c = 10;
        function l(e) {
            if (typeof e != `function`) throw TypeError(`The "listener" argument must be of type Function. Received type ` + typeof e);
        }
        Object.defineProperty(s, "defaultMaxListeners", {
            enumerable: !0,
            get: function() {
                return c;
            },
            set: function(e) {
                if (typeof e != `number` || e < 0 || o(e)) throw RangeError(`The value of "defaultMaxListeners" is out of range. It must be a non-negative number. Received ` + e + `.`);
                c = e;
            }
        }), s.init = function() {
            (this._events === void 0 || this._events === Object.getPrototypeOf(this)._events) && (this._events = Object.create(null), this._eventsCount = 0), this._maxListeners = this._maxListeners || void 0;
        }, s.prototype.setMaxListeners = function(e) {
            if (typeof e != `number` || e < 0 || o(e)) throw RangeError(`The value of "n" is out of range. It must be a non-negative number. Received ` + e + `.`);
            return this._maxListeners = e, this;
        };
        function u(e) {
            return e._maxListeners === void 0 ? s.defaultMaxListeners : e._maxListeners;
        }
        s.prototype.getMaxListeners = function() {
            return u(this);
        }, s.prototype.emit = function(e) {
            for(var t = [], n = 1; n < arguments.length; n++)t.push(arguments[n]);
            var i = e === `error`, a = this._events;
            if (a !== void 0) i &&= a.error === void 0;
            else if (!i) return !1;
            if (i) {
                var o;
                if (t.length > 0 && (o = t[0]), o instanceof Error) throw o;
                var s = Error(`Unhandled error.` + (o ? ` (` + o.message + `)` : ``));
                throw s.context = o, s;
            }
            var c = a[e];
            if (c === void 0) return !1;
            if (typeof c == `function`) r(c, this, t);
            else for(var l = c.length, u = g(c, l), n = 0; n < l; ++n)r(u[n], this, t);
            return !0;
        };
        function d(e, t, n, r) {
            var i, o, s;
            if (l(n), o = e._events, o === void 0 ? (o = e._events = Object.create(null), e._eventsCount = 0) : (o.newListener !== void 0 && (e.emit(`newListener`, t, n.listener ? n.listener : n), o = e._events), s = o[t]), s === void 0) s = o[t] = n, ++e._eventsCount;
            else if (typeof s == `function` ? s = o[t] = r ? [
                n,
                s
            ] : [
                s,
                n
            ] : r ? s.unshift(n) : s.push(n), i = u(e), i > 0 && s.length > i && !s.warned) {
                s.warned = !0;
                var c = Error(`Possible EventEmitter memory leak detected. ` + s.length + ` ` + String(t) + ` listeners added. Use emitter.setMaxListeners() to increase limit`);
                c.name = `MaxListenersExceededWarning`, c.emitter = e, c.type = t, c.count = s.length, a(c);
            }
            return e;
        }
        s.prototype.addListener = function(e, t) {
            return d(this, e, t, !1);
        }, s.prototype.on = s.prototype.addListener, s.prototype.prependListener = function(e, t) {
            return d(this, e, t, !0);
        };
        function f() {
            if (!this.fired) return this.target.removeListener(this.type, this.wrapFn), this.fired = !0, arguments.length === 0 ? this.listener.call(this.target) : this.listener.apply(this.target, arguments);
        }
        function p(e, t, n) {
            var r = {
                fired: !1,
                wrapFn: void 0,
                target: e,
                type: t,
                listener: n
            }, i = f.bind(r);
            return i.listener = n, r.wrapFn = i, i;
        }
        s.prototype.once = function(e, t) {
            return l(t), this.on(e, p(this, e, t)), this;
        }, s.prototype.prependOnceListener = function(e, t) {
            return l(t), this.prependListener(e, p(this, e, t)), this;
        }, s.prototype.removeListener = function(e, t) {
            var n, r, i, a, o;
            if (l(t), r = this._events, r === void 0 || (n = r[e], n === void 0)) return this;
            if (n === t || n.listener === t) --this._eventsCount === 0 ? this._events = Object.create(null) : (delete r[e], r.removeListener && this.emit(`removeListener`, e, n.listener || t));
            else if (typeof n != `function`) {
                for(i = -1, a = n.length - 1; a >= 0; a--)if (n[a] === t || n[a].listener === t) {
                    o = n[a].listener, i = a;
                    break;
                }
                if (i < 0) return this;
                i === 0 ? n.shift() : _(n, i), n.length === 1 && (r[e] = n[0]), r.removeListener !== void 0 && this.emit(`removeListener`, e, o || t);
            }
            return this;
        }, s.prototype.off = s.prototype.removeListener, s.prototype.removeAllListeners = function(e) {
            var t, n = this._events, r;
            if (n === void 0) return this;
            if (n.removeListener === void 0) return arguments.length === 0 ? (this._events = Object.create(null), this._eventsCount = 0) : n[e] !== void 0 && (--this._eventsCount === 0 ? this._events = Object.create(null) : delete n[e]), this;
            if (arguments.length === 0) {
                var i = Object.keys(n), a;
                for(r = 0; r < i.length; ++r)a = i[r], a !== `removeListener` && this.removeAllListeners(a);
                return this.removeAllListeners(`removeListener`), this._events = Object.create(null), this._eventsCount = 0, this;
            }
            if (t = n[e], typeof t == `function`) this.removeListener(e, t);
            else if (t !== void 0) for(r = t.length - 1; r >= 0; r--)this.removeListener(e, t[r]);
            return this;
        };
        function m(e, t, n) {
            var r = e._events;
            if (r === void 0) return [];
            var i = r[t];
            return i === void 0 ? [] : typeof i == `function` ? n ? [
                i.listener || i
            ] : [
                i
            ] : n ? v(i) : g(i, i.length);
        }
        s.prototype.listeners = function(e) {
            return m(this, e, !0);
        }, s.prototype.rawListeners = function(e) {
            return m(this, e, !1);
        }, s.listenerCount = function(e, t) {
            return typeof e.listenerCount == `function` ? e.listenerCount(t) : h.call(e, t);
        }, s.prototype.listenerCount = h;
        function h(e) {
            var t = this._events;
            if (t !== void 0) {
                var n = t[e];
                if (typeof n == `function`) return 1;
                if (n !== void 0) return n.length;
            }
            return 0;
        }
        s.prototype.eventNames = function() {
            return this._eventsCount > 0 ? i(this._events) : [];
        };
        function g(e, t) {
            for(var n = Array(t), r = 0; r < t; ++r)n[r] = e[r];
            return n;
        }
        function _(e, t) {
            for(; t + 1 < e.length; t++)e[t] = e[t + 1];
            e.pop();
        }
        function v(e) {
            for(var t = Array(e.length), n = 0; n < t.length; ++n)t[n] = e[n].listener || e[n];
            return t;
        }
        function y(e, t) {
            return new Promise(function(n, r) {
                function i(n) {
                    e.removeListener(t, a), r(n);
                }
                function a() {
                    typeof e.removeListener == `function` && e.removeListener(`error`, i), n([].slice.call(arguments));
                }
                x(e, t, a, {
                    once: !0
                }), t !== `error` && b(e, i, {
                    once: !0
                });
            });
        }
        function b(e, t, n) {
            typeof e.on == `function` && x(e, `error`, t, n);
        }
        function x(e, t, n, r) {
            if (typeof e.on == `function`) r.once ? e.once(t, n) : e.on(t, n);
            else if (typeof e.addEventListener == `function`) e.addEventListener(t, function i(a) {
                r.once && e.removeEventListener(t, i), n(a);
            });
            else throw TypeError(`The "emitter" argument must be of type EventEmitter. Received type ` + typeof e);
        }
    })), g = t(((e, t)=>{
        typeof Object.create == `function` ? t.exports = function(e, t) {
            t && (e.super_ = t, e.prototype = Object.create(t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }));
        } : t.exports = function(e, t) {
            if (t) {
                e.super_ = t;
                var n = function() {};
                n.prototype = t.prototype, e.prototype = new n, e.prototype.constructor = e;
            }
        };
    })), _ = t(((e, t)=>{
        t.exports = h().EventEmitter;
    })), v = t(((e, t)=>{
        t.exports = function() {
            if (typeof Symbol != `function` || typeof Object.getOwnPropertySymbols != `function`) return !1;
            if (typeof Symbol.iterator == `symbol`) return !0;
            var e = {}, t = Symbol(`test`), n = Object(t);
            if (typeof t == `string` || Object.prototype.toString.call(t) !== `[object Symbol]` || Object.prototype.toString.call(n) !== `[object Symbol]`) return !1;
            var r = 42;
            for(var i in e[t] = r, e)return !1;
            if (typeof Object.keys == `function` && Object.keys(e).length !== 0 || typeof Object.getOwnPropertyNames == `function` && Object.getOwnPropertyNames(e).length !== 0) return !1;
            var a = Object.getOwnPropertySymbols(e);
            if (a.length !== 1 || a[0] !== t || !Object.prototype.propertyIsEnumerable.call(e, t)) return !1;
            if (typeof Object.getOwnPropertyDescriptor == `function`) {
                var o = Object.getOwnPropertyDescriptor(e, t);
                if (o.value !== r || o.enumerable !== !0) return !1;
            }
            return !0;
        };
    })), y = t(((e, t)=>{
        var n = v();
        t.exports = function() {
            return n() && !!Symbol.toStringTag;
        };
    })), b = t(((e, t)=>{
        t.exports = Object;
    })), x = t(((e, t)=>{
        t.exports = Error;
    })), S = t(((e, t)=>{
        t.exports = EvalError;
    })), C = t(((e, t)=>{
        t.exports = RangeError;
    })), w = t(((e, t)=>{
        t.exports = ReferenceError;
    })), T = t(((e, t)=>{
        t.exports = SyntaxError;
    })), E = t(((e, t)=>{
        t.exports = TypeError;
    })), D = t(((e, t)=>{
        t.exports = URIError;
    })), O = t(((e, t)=>{
        t.exports = Math.abs;
    })), k = t(((e, t)=>{
        t.exports = Math.floor;
    })), A = t(((e, t)=>{
        t.exports = Math.max;
    })), j = t(((e, t)=>{
        t.exports = Math.min;
    })), M = t(((e, t)=>{
        t.exports = Math.pow;
    })), N = t(((e, t)=>{
        t.exports = Math.round;
    })), P = t(((e, t)=>{
        t.exports = Number.isNaN || function(e) {
            return e !== e;
        };
    })), ee = t(((e, t)=>{
        var n = P();
        t.exports = function(e) {
            return n(e) || e === 0 ? e : e < 0 ? -1 : 1;
        };
    })), F = t(((e, t)=>{
        t.exports = Object.getOwnPropertyDescriptor;
    })), te = t(((e, t)=>{
        var n = F();
        if (n) try {
            n([], `length`);
        } catch  {
            n = null;
        }
        t.exports = n;
    })), I = t(((e, t)=>{
        var n = Object.defineProperty || !1;
        if (n) try {
            n({}, `a`, {
                value: 1
            });
        } catch  {
            n = !1;
        }
        t.exports = n;
    })), ne = t(((e, t)=>{
        var n = typeof Symbol < `u` && Symbol, r = v();
        t.exports = function() {
            return typeof n != `function` || typeof Symbol != `function` || typeof n(`foo`) != `symbol` || typeof Symbol(`bar`) != `symbol` ? !1 : r();
        };
    })), re = t(((e, t)=>{
        t.exports = typeof Reflect < `u` && Reflect.getPrototypeOf || null;
    })), L = t(((e, t)=>{
        t.exports = b().getPrototypeOf || null;
    })), R = t(((e, t)=>{
        var n = `Function.prototype.bind called on incompatible `, r = Object.prototype.toString, i = Math.max, a = `[object Function]`, o = function(e, t) {
            for(var n = [], r = 0; r < e.length; r += 1)n[r] = e[r];
            for(var i = 0; i < t.length; i += 1)n[i + e.length] = t[i];
            return n;
        }, s = function(e, t) {
            for(var n = [], r = t || 0, i = 0; r < e.length; r += 1, i += 1)n[i] = e[r];
            return n;
        }, c = function(e, t) {
            for(var n = ``, r = 0; r < e.length; r += 1)n += e[r], r + 1 < e.length && (n += t);
            return n;
        };
        t.exports = function(e) {
            var t = this;
            if (typeof t != `function` || r.apply(t) !== a) throw TypeError(n + t);
            for(var l = s(arguments, 1), u, d = function() {
                if (this instanceof u) {
                    var n = t.apply(this, o(l, arguments));
                    return Object(n) === n ? n : this;
                }
                return t.apply(e, o(l, arguments));
            }, f = i(0, t.length - l.length), p = [], m = 0; m < f; m++)p[m] = `$` + m;
            if (u = Function(`binder`, `return function (` + c(p, `,`) + `){ return binder.apply(this,arguments); }`)(d), t.prototype) {
                var h = function() {};
                h.prototype = t.prototype, u.prototype = new h, h.prototype = null;
            }
            return u;
        };
    })), z = t(((e, t)=>{
        var n = R();
        t.exports = Function.prototype.bind || n;
    })), B = t(((e, t)=>{
        t.exports = Function.prototype.call;
    })), ie = t(((e, t)=>{
        t.exports = Function.prototype.apply;
    })), ae = t(((e, t)=>{
        t.exports = typeof Reflect < `u` && Reflect && Reflect.apply;
    })), oe = t(((e, t)=>{
        var n = z(), r = ie(), i = B();
        t.exports = ae() || n.call(i, r);
    })), V = t(((e, t)=>{
        var n = z(), r = E(), i = B(), a = oe();
        t.exports = function(e) {
            if (e.length < 1 || typeof e[0] != `function`) throw new r(`a function is required`);
            return a(n, i, e);
        };
    })), H = t(((e, t)=>{
        var n = V(), r = te(), i;
        try {
            i = [].__proto__ === Array.prototype;
        } catch (e) {
            if (!e || typeof e != `object` || !(`code` in e) || e.code !== `ERR_PROTO_ACCESS`) throw e;
        }
        var a = !!i && r && r(Object.prototype, `__proto__`), o = Object, s = o.getPrototypeOf;
        t.exports = a && typeof a.get == `function` ? n([
            a.get
        ]) : typeof s == `function` ? function(e) {
            return s(e == null ? e : o(e));
        } : !1;
    })), U = t(((e, t)=>{
        var n = re(), r = L(), i = H();
        t.exports = n ? function(e) {
            return n(e);
        } : r ? function(e) {
            if (!e || typeof e != `object` && typeof e != `function`) throw TypeError(`getProto: not an object`);
            return r(e);
        } : i ? function(e) {
            return i(e);
        } : null;
    })), se = t(((e, t)=>{
        var n = Function.prototype.call, r = Object.prototype.hasOwnProperty;
        t.exports = z().call(n, r);
    })), W = t(((e, t)=>{
        var n, r = b(), i = x(), a = S(), o = C(), s = w(), c = T(), l = E(), u = D(), d = O(), f = k(), p = A(), m = j(), h = M(), g = N(), _ = ee(), v = Function, y = function(e) {
            try {
                return v(`"use strict"; return (` + e + `).constructor;`)();
            } catch  {}
        }, P = te(), F = I(), R = function() {
            throw new l;
        }, ae = P ? function() {
            try {
                return arguments.callee, R;
            } catch  {
                try {
                    return P(arguments, `callee`).get;
                } catch  {
                    return R;
                }
            }
        }() : R, oe = ne()(), V = U(), H = L(), W = re(), G = ie(), K = B(), ce = {}, le = typeof Uint8Array > `u` || !V ? n : V(Uint8Array), q = {
            __proto__: null,
            "%AggregateError%": typeof AggregateError > `u` ? n : AggregateError,
            "%Array%": Array,
            "%ArrayBuffer%": typeof ArrayBuffer > `u` ? n : ArrayBuffer,
            "%ArrayIteratorPrototype%": oe && V ? V([][Symbol.iterator]()) : n,
            "%AsyncFromSyncIteratorPrototype%": n,
            "%AsyncFunction%": ce,
            "%AsyncGenerator%": ce,
            "%AsyncGeneratorFunction%": ce,
            "%AsyncIteratorPrototype%": ce,
            "%Atomics%": typeof Atomics > `u` ? n : Atomics,
            "%BigInt%": typeof BigInt > `u` ? n : BigInt,
            "%BigInt64Array%": typeof BigInt64Array > `u` ? n : BigInt64Array,
            "%BigUint64Array%": typeof BigUint64Array > `u` ? n : BigUint64Array,
            "%Boolean%": Boolean,
            "%DataView%": typeof DataView > `u` ? n : DataView,
            "%Date%": Date,
            "%decodeURI%": decodeURI,
            "%decodeURIComponent%": decodeURIComponent,
            "%encodeURI%": encodeURI,
            "%encodeURIComponent%": encodeURIComponent,
            "%Error%": i,
            "%eval%": eval,
            "%EvalError%": a,
            "%Float16Array%": typeof Float16Array > `u` ? n : Float16Array,
            "%Float32Array%": typeof Float32Array > `u` ? n : Float32Array,
            "%Float64Array%": typeof Float64Array > `u` ? n : Float64Array,
            "%FinalizationRegistry%": typeof FinalizationRegistry > `u` ? n : FinalizationRegistry,
            "%Function%": v,
            "%GeneratorFunction%": ce,
            "%Int8Array%": typeof Int8Array > `u` ? n : Int8Array,
            "%Int16Array%": typeof Int16Array > `u` ? n : Int16Array,
            "%Int32Array%": typeof Int32Array > `u` ? n : Int32Array,
            "%isFinite%": isFinite,
            "%isNaN%": isNaN,
            "%IteratorPrototype%": oe && V ? V(V([][Symbol.iterator]())) : n,
            "%JSON%": typeof JSON == `object` ? JSON : n,
            "%Map%": typeof Map > `u` ? n : Map,
            "%MapIteratorPrototype%": typeof Map > `u` || !oe || !V ? n : V(new Map()[Symbol.iterator]()),
            "%Math%": Math,
            "%Number%": Number,
            "%Object%": r,
            "%Object.getOwnPropertyDescriptor%": P,
            "%parseFloat%": parseFloat,
            "%parseInt%": parseInt,
            "%Promise%": typeof Promise > `u` ? n : Promise,
            "%Proxy%": typeof Proxy > `u` ? n : Proxy,
            "%RangeError%": o,
            "%ReferenceError%": s,
            "%Reflect%": typeof Reflect > `u` ? n : Reflect,
            "%RegExp%": RegExp,
            "%Set%": typeof Set > `u` ? n : Set,
            "%SetIteratorPrototype%": typeof Set > `u` || !oe || !V ? n : V(new Set()[Symbol.iterator]()),
            "%SharedArrayBuffer%": typeof SharedArrayBuffer > `u` ? n : SharedArrayBuffer,
            "%String%": String,
            "%StringIteratorPrototype%": oe && V ? V(``[Symbol.iterator]()) : n,
            "%Symbol%": oe ? Symbol : n,
            "%SyntaxError%": c,
            "%ThrowTypeError%": ae,
            "%TypedArray%": le,
            "%TypeError%": l,
            "%Uint8Array%": typeof Uint8Array > `u` ? n : Uint8Array,
            "%Uint8ClampedArray%": typeof Uint8ClampedArray > `u` ? n : Uint8ClampedArray,
            "%Uint16Array%": typeof Uint16Array > `u` ? n : Uint16Array,
            "%Uint32Array%": typeof Uint32Array > `u` ? n : Uint32Array,
            "%URIError%": u,
            "%WeakMap%": typeof WeakMap > `u` ? n : WeakMap,
            "%WeakRef%": typeof WeakRef > `u` ? n : WeakRef,
            "%WeakSet%": typeof WeakSet > `u` ? n : WeakSet,
            "%Function.prototype.call%": K,
            "%Function.prototype.apply%": G,
            "%Object.defineProperty%": F,
            "%Object.getPrototypeOf%": H,
            "%Math.abs%": d,
            "%Math.floor%": f,
            "%Math.max%": p,
            "%Math.min%": m,
            "%Math.pow%": h,
            "%Math.round%": g,
            "%Math.sign%": _,
            "%Reflect.getPrototypeOf%": W
        };
        if (V) try {
            null.error;
        } catch (e) {
            q[`%Error.prototype%`] = V(V(e));
        }
        var ue = function e(t) {
            var n;
            if (t === `%AsyncFunction%`) n = y(`async function () {}`);
            else if (t === `%GeneratorFunction%`) n = y(`function* () {}`);
            else if (t === `%AsyncGeneratorFunction%`) n = y(`async function* () {}`);
            else if (t === `%AsyncGenerator%`) {
                var r = e(`%AsyncGeneratorFunction%`);
                r && (n = r.prototype);
            } else if (t === `%AsyncIteratorPrototype%`) {
                var i = e(`%AsyncGenerator%`);
                i && V && (n = V(i.prototype));
            }
            return q[t] = n, n;
        }, J = {
            __proto__: null,
            "%ArrayBufferPrototype%": [
                `ArrayBuffer`,
                `prototype`
            ],
            "%ArrayPrototype%": [
                `Array`,
                `prototype`
            ],
            "%ArrayProto_entries%": [
                `Array`,
                `prototype`,
                `entries`
            ],
            "%ArrayProto_forEach%": [
                `Array`,
                `prototype`,
                `forEach`
            ],
            "%ArrayProto_keys%": [
                `Array`,
                `prototype`,
                `keys`
            ],
            "%ArrayProto_values%": [
                `Array`,
                `prototype`,
                `values`
            ],
            "%AsyncFunctionPrototype%": [
                `AsyncFunction`,
                `prototype`
            ],
            "%AsyncGenerator%": [
                `AsyncGeneratorFunction`,
                `prototype`
            ],
            "%AsyncGeneratorPrototype%": [
                `AsyncGeneratorFunction`,
                `prototype`,
                `prototype`
            ],
            "%BooleanPrototype%": [
                `Boolean`,
                `prototype`
            ],
            "%DataViewPrototype%": [
                `DataView`,
                `prototype`
            ],
            "%DatePrototype%": [
                `Date`,
                `prototype`
            ],
            "%ErrorPrototype%": [
                `Error`,
                `prototype`
            ],
            "%EvalErrorPrototype%": [
                `EvalError`,
                `prototype`
            ],
            "%Float32ArrayPrototype%": [
                `Float32Array`,
                `prototype`
            ],
            "%Float64ArrayPrototype%": [
                `Float64Array`,
                `prototype`
            ],
            "%FunctionPrototype%": [
                `Function`,
                `prototype`
            ],
            "%Generator%": [
                `GeneratorFunction`,
                `prototype`
            ],
            "%GeneratorPrototype%": [
                `GeneratorFunction`,
                `prototype`,
                `prototype`
            ],
            "%Int8ArrayPrototype%": [
                `Int8Array`,
                `prototype`
            ],
            "%Int16ArrayPrototype%": [
                `Int16Array`,
                `prototype`
            ],
            "%Int32ArrayPrototype%": [
                `Int32Array`,
                `prototype`
            ],
            "%JSONParse%": [
                `JSON`,
                `parse`
            ],
            "%JSONStringify%": [
                `JSON`,
                `stringify`
            ],
            "%MapPrototype%": [
                `Map`,
                `prototype`
            ],
            "%NumberPrototype%": [
                `Number`,
                `prototype`
            ],
            "%ObjectPrototype%": [
                `Object`,
                `prototype`
            ],
            "%ObjProto_toString%": [
                `Object`,
                `prototype`,
                `toString`
            ],
            "%ObjProto_valueOf%": [
                `Object`,
                `prototype`,
                `valueOf`
            ],
            "%PromisePrototype%": [
                `Promise`,
                `prototype`
            ],
            "%PromiseProto_then%": [
                `Promise`,
                `prototype`,
                `then`
            ],
            "%Promise_all%": [
                `Promise`,
                `all`
            ],
            "%Promise_reject%": [
                `Promise`,
                `reject`
            ],
            "%Promise_resolve%": [
                `Promise`,
                `resolve`
            ],
            "%RangeErrorPrototype%": [
                `RangeError`,
                `prototype`
            ],
            "%ReferenceErrorPrototype%": [
                `ReferenceError`,
                `prototype`
            ],
            "%RegExpPrototype%": [
                `RegExp`,
                `prototype`
            ],
            "%SetPrototype%": [
                `Set`,
                `prototype`
            ],
            "%SharedArrayBufferPrototype%": [
                `SharedArrayBuffer`,
                `prototype`
            ],
            "%StringPrototype%": [
                `String`,
                `prototype`
            ],
            "%SymbolPrototype%": [
                `Symbol`,
                `prototype`
            ],
            "%SyntaxErrorPrototype%": [
                `SyntaxError`,
                `prototype`
            ],
            "%TypedArrayPrototype%": [
                `TypedArray`,
                `prototype`
            ],
            "%TypeErrorPrototype%": [
                `TypeError`,
                `prototype`
            ],
            "%Uint8ArrayPrototype%": [
                `Uint8Array`,
                `prototype`
            ],
            "%Uint8ClampedArrayPrototype%": [
                `Uint8ClampedArray`,
                `prototype`
            ],
            "%Uint16ArrayPrototype%": [
                `Uint16Array`,
                `prototype`
            ],
            "%Uint32ArrayPrototype%": [
                `Uint32Array`,
                `prototype`
            ],
            "%URIErrorPrototype%": [
                `URIError`,
                `prototype`
            ],
            "%WeakMapPrototype%": [
                `WeakMap`,
                `prototype`
            ],
            "%WeakSetPrototype%": [
                `WeakSet`,
                `prototype`
            ]
        }, Y = z(), X = se(), Z = Y.call(K, Array.prototype.concat), de = Y.call(G, Array.prototype.splice), Q = Y.call(K, String.prototype.replace), fe = Y.call(K, String.prototype.slice), pe = Y.call(K, RegExp.prototype.exec), $ = /[^%.[\]]+|\[(?:(-?\d+(?:\.\d+)?)|(["'])((?:(?!\2)[^\\]|\\.)*?)\2)\]|(?=(?:\.|\[\])(?:\.|\[\]|%$))/g, me = /\\(\\)?/g, he = function(e) {
            var t = fe(e, 0, 1), n = fe(e, -1);
            if (t === `%` && n !== `%`) throw new c("invalid intrinsic syntax, expected closing `%`");
            if (n === `%` && t !== `%`) throw new c("invalid intrinsic syntax, expected opening `%`");
            var r = [];
            return Q(e, $, function(e, t, n, i) {
                r[r.length] = n ? Q(i, me, `$1`) : t || e;
            }), r;
        }, ge = function(e, t) {
            var n = e, r;
            if (X(J, n) && (r = J[n], n = `%` + r[0] + `%`), X(q, n)) {
                var i = q[n];
                if (i === ce && (i = ue(n)), i === void 0 && !t) throw new l(`intrinsic ` + e + ` exists, but is not available. Please file an issue!`);
                return {
                    alias: r,
                    name: n,
                    value: i
                };
            }
            throw new c(`intrinsic ` + e + ` does not exist!`);
        };
        t.exports = function(e, t) {
            if (typeof e != `string` || e.length === 0) throw new l(`intrinsic name must be a non-empty string`);
            if (arguments.length > 1 && typeof t != `boolean`) throw new l(`"allowMissing" argument must be a boolean`);
            if (pe(/^%?[^%]*%?$/, e) === null) throw new c("`%` may not be present anywhere but at the beginning and end of the intrinsic name");
            var n = he(e), r = n.length > 0 ? n[0] : ``, i = ge(`%` + r + `%`, t), a = i.name, o = i.value, s = !1, u = i.alias;
            u && (r = u[0], de(n, Z([
                0,
                1
            ], u)));
            for(var d = 1, f = !0; d < n.length; d += 1){
                var p = n[d], m = fe(p, 0, 1), h = fe(p, -1);
                if ((m === `"` || m === `'` || m === "`" || h === `"` || h === `'` || h === "`") && m !== h) throw new c(`property names with quotes must have matching quotes`);
                if ((p === `constructor` || !f) && (s = !0), r += `.` + p, a = `%` + r + `%`, X(q, a)) o = q[a];
                else if (o != null) {
                    if (!(p in o)) {
                        if (!t) throw new l(`base intrinsic for ` + e + ` exists, but the property is not available.`);
                        return;
                    }
                    if (P && d + 1 >= n.length) {
                        var g = P(o, p);
                        f = !!g, o = f && `get` in g && !(`originalValue` in g.get) ? g.get : o[p];
                    } else f = X(o, p), o = o[p];
                    f && !s && (q[a] = o);
                }
            }
            return o;
        };
    })), G = t(((e, t)=>{
        var n = W(), r = V(), i = r([
            n(`%String.prototype.indexOf%`)
        ]);
        t.exports = function(e, t) {
            var a = n(e, !!t);
            return typeof a == `function` && i(e, `.prototype.`) > -1 ? r([
                a
            ]) : a;
        };
    })), K = t(((e, t)=>{
        var n = y()(), r = G()(`Object.prototype.toString`), i = function(e) {
            return n && e && typeof e == `object` && Symbol.toStringTag in e ? !1 : r(e) === `[object Arguments]`;
        }, a = function(e) {
            return i(e) ? !0 : typeof e == `object` && !!e && `length` in e && typeof e.length == `number` && e.length >= 0 && r(e) !== `[object Array]` && `callee` in e && r(e.callee) === `[object Function]`;
        }, o = function() {
            return i(arguments);
        }();
        i.isLegacyArguments = a, t.exports = o ? i : a;
    })), ce = t(((e, t)=>{
        var n = G(), r = y()(), i = se(), a = te(), o;
        if (r) {
            var s = n(`RegExp.prototype.exec`), c = {}, l = function() {
                throw c;
            }, u = {
                toString: l,
                valueOf: l
            };
            typeof Symbol.toPrimitive == `symbol` && (u[Symbol.toPrimitive] = l), o = function(e) {
                if (!e || typeof e != `object`) return !1;
                var t = a(e, `lastIndex`);
                if (!(t && i(t, `value`))) return !1;
                try {
                    s(e, u);
                } catch (e) {
                    return e === c;
                }
            };
        } else {
            var d = n(`Object.prototype.toString`), f = `[object RegExp]`;
            o = function(e) {
                return !e || typeof e != `object` && typeof e != `function` ? !1 : d(e) === f;
            };
        }
        t.exports = o;
    })), le = t(((e, t)=>{
        var n = G(), r = ce(), i = n(`RegExp.prototype.exec`), a = E();
        t.exports = function(e) {
            if (!r(e)) throw new a("`regex` must be a RegExp");
            return function(t) {
                return i(e, t) !== null;
            };
        };
    })), q = t(((e, t)=>{
        let n = function*() {}.constructor;
        t.exports = ()=>n;
    })), ue = t(((e, t)=>{
        var n = G(), r = le()(/^\s*(?:function)?\*/), i = y()(), a = U(), o = n(`Object.prototype.toString`), s = n(`Function.prototype.toString`), c = q();
        t.exports = function(e) {
            if (typeof e != `function`) return !1;
            if (r(s(e))) return !0;
            if (!i) return o(e) === `[object GeneratorFunction]`;
            if (!a) return !1;
            var t = c();
            return t && a(e) === t.prototype;
        };
    })), J = t(((e, t)=>{
        var n = Function.prototype.toString, r = typeof Reflect == `object` && Reflect !== null && Reflect.apply, i, a;
        if (typeof r == `function` && typeof Object.defineProperty == `function`) try {
            i = Object.defineProperty({}, "length", {
                get: function() {
                    throw a;
                }
            }), a = {}, r(function() {
                throw 42;
            }, null, i);
        } catch (e) {
            e !== a && (r = null);
        }
        else r = null;
        var o = /^\s*class\b/, s = function(e) {
            try {
                var t = n.call(e);
                return o.test(t);
            } catch  {
                return !1;
            }
        }, c = function(e) {
            try {
                return s(e) ? !1 : (n.call(e), !0);
            } catch  {
                return !1;
            }
        }, l = Object.prototype.toString, u = `[object Object]`, d = `[object Function]`, f = `[object GeneratorFunction]`, p = `[object HTMLAllCollection]`, m = `[object HTML document.all class]`, h = `[object HTMLCollection]`, g = typeof Symbol == `function` && !!Symbol.toStringTag, _ = !(0 in [
            , 
        ]), v = function() {
            return !1;
        };
        if (typeof document == `object`) {
            var y = document.all;
            l.call(y) === l.call(document.all) && (v = function(e) {
                if ((_ || !e) && (e === void 0 || typeof e == `object`)) try {
                    var t = l.call(e);
                    return (t === p || t === m || t === h || t === u) && e(``) == null;
                } catch  {}
                return !1;
            });
        }
        t.exports = r ? function(e) {
            if (v(e)) return !0;
            if (!e || typeof e != `function` && typeof e != `object`) return !1;
            try {
                r(e, null, i);
            } catch (e) {
                if (e !== a) return !1;
            }
            return !s(e) && c(e);
        } : function(e) {
            if (v(e)) return !0;
            if (!e || typeof e != `function` && typeof e != `object`) return !1;
            if (g) return c(e);
            if (s(e)) return !1;
            var t = l.call(e);
            return t !== d && t !== f && !/^\[object HTML/.test(t) ? !1 : c(e);
        };
    })), Y = t(((e, t)=>{
        var n = J(), r = Object.prototype.toString, i = Object.prototype.hasOwnProperty, a = function(e, t, n) {
            for(var r = 0, a = e.length; r < a; r++)i.call(e, r) && (n == null ? t(e[r], r, e) : t.call(n, e[r], r, e));
        }, o = function(e, t, n) {
            for(var r = 0, i = e.length; r < i; r++)n == null ? t(e.charAt(r), r, e) : t.call(n, e.charAt(r), r, e);
        }, s = function(e, t, n) {
            for(var r in e)i.call(e, r) && (n == null ? t(e[r], r, e) : t.call(n, e[r], r, e));
        };
        function c(e) {
            return r.call(e) === `[object Array]`;
        }
        t.exports = function(e, t, r) {
            if (!n(t)) throw TypeError(`iterator must be a function`);
            var i;
            arguments.length >= 3 && (i = r), c(e) ? a(e, t, i) : typeof e == `string` ? o(e, t, i) : s(e, t, i);
        };
    })), X = t(((e, t)=>{
        t.exports = [
            `Float16Array`,
            `Float32Array`,
            `Float64Array`,
            `Int8Array`,
            `Int16Array`,
            `Int32Array`,
            `Uint8Array`,
            `Uint8ClampedArray`,
            `Uint16Array`,
            `Uint32Array`,
            `BigInt64Array`,
            `BigUint64Array`
        ];
    })), Z = t(((e, t)=>{
        var n = X(), r = globalThis;
        t.exports = function() {
            for(var e = [], t = 0; t < n.length; t++)typeof r[n[t]] == `function` && (e[e.length] = n[t]);
            return e;
        };
    })), de = t(((e, t)=>{
        var n = I(), r = T(), i = E(), a = te();
        t.exports = function(e, t, o) {
            if (!e || typeof e != `object` && typeof e != `function`) throw new i("`obj` must be an object or a function`");
            if (typeof t != `string` && typeof t != `symbol`) throw new i("`property` must be a string or a symbol`");
            if (arguments.length > 3 && typeof arguments[3] != `boolean` && arguments[3] !== null) throw new i("`nonEnumerable`, if provided, must be a boolean or null");
            if (arguments.length > 4 && typeof arguments[4] != `boolean` && arguments[4] !== null) throw new i("`nonWritable`, if provided, must be a boolean or null");
            if (arguments.length > 5 && typeof arguments[5] != `boolean` && arguments[5] !== null) throw new i("`nonConfigurable`, if provided, must be a boolean or null");
            if (arguments.length > 6 && typeof arguments[6] != `boolean`) throw new i("`loose`, if provided, must be a boolean");
            var s = arguments.length > 3 ? arguments[3] : null, c = arguments.length > 4 ? arguments[4] : null, l = arguments.length > 5 ? arguments[5] : null, u = arguments.length > 6 ? arguments[6] : !1, d = !!a && a(e, t);
            if (n) n(e, t, {
                configurable: l === null && d ? d.configurable : !l,
                enumerable: s === null && d ? d.enumerable : !s,
                value: o,
                writable: c === null && d ? d.writable : !c
            });
            else if (u || !s && !c && !l) e[t] = o;
            else throw new r(`This environment does not support defining a property as non-configurable, non-writable, or non-enumerable.`);
        };
    })), Q = t(((e, t)=>{
        var n = I(), r = function() {
            return !!n;
        };
        r.hasArrayLengthDefineBug = function() {
            if (!n) return null;
            try {
                return n([], `length`, {
                    value: 1
                }).length !== 1;
            } catch  {
                return !0;
            }
        }, t.exports = r;
    })), fe = t(((e, t)=>{
        var n = W(), r = de(), i = Q()(), a = te(), o = E(), s = n(`%Math.floor%`);
        t.exports = function(e, t) {
            if (typeof e != `function`) throw new o("`fn` is not a function");
            if (typeof t != `number` || t < 0 || t > 4294967295 || s(t) !== t) throw new o("`length` must be a positive 32-bit integer");
            var n = arguments.length > 2 && !!arguments[2], c = !0, l = !0;
            if (`length` in e && a) {
                var u = a(e, `length`);
                u && !u.configurable && (c = !1), u && !u.writable && (l = !1);
            }
            return (c || l || !n) && (i ? r(e, `length`, t, !0, !0) : r(e, `length`, t)), e;
        };
    })), pe = t(((e, t)=>{
        var n = z(), r = ie(), i = oe();
        t.exports = function() {
            return i(n, r, arguments);
        };
    })), $ = t(((e, t)=>{
        var n = fe(), r = I(), i = V(), a = pe();
        t.exports = function(e) {
            var t = i(arguments), r = 1 + e.length - (arguments.length - 1);
            return n(t, r > 0 ? r : 0, !0);
        }, r ? r(t.exports, `apply`, {
            value: a
        }) : t.exports.apply = a;
    })), me = t(((e, t)=>{
        var n = Y(), r = Z(), i = $(), a = G(), o = te(), s = U(), c = a(`Object.prototype.toString`), l = y()(), u = globalThis, d = r(), f = a(`String.prototype.slice`), p = a(`Array.prototype.indexOf`, !0) || function(e, t) {
            for(var n = 0; n < e.length; n += 1)if (e[n] === t) return n;
            return -1;
        }, m = {
            __proto__: null
        };
        l && o && s ? n(d, function(e) {
            var t = new u[e];
            if (Symbol.toStringTag in t && s) {
                var n = s(t), r = o(n, Symbol.toStringTag);
                if (!r && n && (r = o(s(n), Symbol.toStringTag)), r && r.get) {
                    var a = i(r.get);
                    m[`$` + e] = a;
                }
            }
        }) : n(d, function(e) {
            var t = new u[e], n = t.slice || t.set;
            if (n) {
                var r = i(n);
                m[`$` + e] = r;
            }
        });
        var h = function(e) {
            var t = !1;
            return n(m, function(n, r) {
                if (!t) try {
                    `$` + n(e) === r && (t = f(r, 1));
                } catch  {}
            }), t;
        }, g = function(e) {
            var t = !1;
            return n(m, function(n, r) {
                if (!t) try {
                    n(e), t = f(r, 1);
                } catch  {}
            }), t;
        };
        t.exports = function(e) {
            if (!e || typeof e != `object`) return !1;
            if (!l) {
                var t = f(c(e), 8, -1);
                return p(d, t) > -1 ? t : t === `Object` ? g(e) : !1;
            }
            return o ? h(e) : null;
        };
    })), he = t(((e, t)=>{
        var n = me();
        t.exports = function(e) {
            return !!n(e);
        };
    })), ge = t(((e)=>{
        var t = K(), n = ue(), r = me(), i = he();
        function a(e) {
            return e.call.bind(e);
        }
        var o = typeof BigInt < `u`, s = typeof Symbol < `u`, c = a(Object.prototype.toString), l = a(Number.prototype.valueOf), u = a(String.prototype.valueOf), d = a(Boolean.prototype.valueOf);
        if (o) var f = a(BigInt.prototype.valueOf);
        if (s) var p = a(Symbol.prototype.valueOf);
        function m(e, t) {
            if (typeof e != `object`) return !1;
            try {
                return t(e), !0;
            } catch  {
                return !1;
            }
        }
        e.isArgumentsObject = t, e.isGeneratorFunction = n, e.isTypedArray = i;
        function h(e) {
            return typeof Promise < `u` && e instanceof Promise || typeof e == `object` && !!e && typeof e.then == `function` && typeof e.catch == `function`;
        }
        e.isPromise = h;
        function g(e) {
            return typeof ArrayBuffer < `u` && ArrayBuffer.isView ? ArrayBuffer.isView(e) : i(e) || ne(e);
        }
        e.isArrayBufferView = g;
        function _(e) {
            return r(e) === `Uint8Array`;
        }
        e.isUint8Array = _;
        function v(e) {
            return r(e) === `Uint8ClampedArray`;
        }
        e.isUint8ClampedArray = v;
        function y(e) {
            return r(e) === `Uint16Array`;
        }
        e.isUint16Array = y;
        function b(e) {
            return r(e) === `Uint32Array`;
        }
        e.isUint32Array = b;
        function x(e) {
            return r(e) === `Int8Array`;
        }
        e.isInt8Array = x;
        function S(e) {
            return r(e) === `Int16Array`;
        }
        e.isInt16Array = S;
        function C(e) {
            return r(e) === `Int32Array`;
        }
        e.isInt32Array = C;
        function w(e) {
            return r(e) === `Float32Array`;
        }
        e.isFloat32Array = w;
        function T(e) {
            return r(e) === `Float64Array`;
        }
        e.isFloat64Array = T;
        function E(e) {
            return r(e) === `BigInt64Array`;
        }
        e.isBigInt64Array = E;
        function D(e) {
            return r(e) === `BigUint64Array`;
        }
        e.isBigUint64Array = D;
        function O(e) {
            return c(e) === `[object Map]`;
        }
        O.working = typeof Map < `u` && O(new Map);
        function k(e) {
            return typeof Map > `u` ? !1 : O.working ? O(e) : e instanceof Map;
        }
        e.isMap = k;
        function A(e) {
            return c(e) === `[object Set]`;
        }
        A.working = typeof Set < `u` && A(new Set);
        function j(e) {
            return typeof Set > `u` ? !1 : A.working ? A(e) : e instanceof Set;
        }
        e.isSet = j;
        function M(e) {
            return c(e) === `[object WeakMap]`;
        }
        M.working = typeof WeakMap < `u` && M(new WeakMap);
        function N(e) {
            return typeof WeakMap > `u` ? !1 : M.working ? M(e) : e instanceof WeakMap;
        }
        e.isWeakMap = N;
        function P(e) {
            return c(e) === `[object WeakSet]`;
        }
        P.working = typeof WeakSet < `u` && P(new WeakSet);
        function ee(e) {
            return P(e);
        }
        e.isWeakSet = ee;
        function F(e) {
            return c(e) === `[object ArrayBuffer]`;
        }
        F.working = typeof ArrayBuffer < `u` && F(new ArrayBuffer);
        function te(e) {
            return typeof ArrayBuffer > `u` ? !1 : F.working ? F(e) : e instanceof ArrayBuffer;
        }
        e.isArrayBuffer = te;
        function I(e) {
            return c(e) === `[object DataView]`;
        }
        I.working = typeof ArrayBuffer < `u` && typeof DataView < `u` && I(new DataView(new ArrayBuffer(1), 0, 1));
        function ne(e) {
            return typeof DataView > `u` ? !1 : I.working ? I(e) : e instanceof DataView;
        }
        e.isDataView = ne;
        var re = typeof SharedArrayBuffer < `u` ? SharedArrayBuffer : void 0;
        function L(e) {
            return c(e) === `[object SharedArrayBuffer]`;
        }
        function R(e) {
            return re === void 0 ? !1 : (L.working === void 0 && (L.working = L(new re)), L.working ? L(e) : e instanceof re);
        }
        e.isSharedArrayBuffer = R;
        function z(e) {
            return c(e) === `[object AsyncFunction]`;
        }
        e.isAsyncFunction = z;
        function B(e) {
            return c(e) === `[object Map Iterator]`;
        }
        e.isMapIterator = B;
        function ie(e) {
            return c(e) === `[object Set Iterator]`;
        }
        e.isSetIterator = ie;
        function ae(e) {
            return c(e) === `[object Generator]`;
        }
        e.isGeneratorObject = ae;
        function oe(e) {
            return c(e) === `[object WebAssembly.Module]`;
        }
        e.isWebAssemblyCompiledModule = oe;
        function V(e) {
            return m(e, l);
        }
        e.isNumberObject = V;
        function H(e) {
            return m(e, u);
        }
        e.isStringObject = H;
        function U(e) {
            return m(e, d);
        }
        e.isBooleanObject = U;
        function se(e) {
            return o && m(e, f);
        }
        e.isBigIntObject = se;
        function W(e) {
            return s && m(e, p);
        }
        e.isSymbolObject = W;
        function G(e) {
            return V(e) || H(e) || U(e) || se(e) || W(e);
        }
        e.isBoxedPrimitive = G;
        function ce(e) {
            return typeof Uint8Array < `u` && (te(e) || R(e));
        }
        e.isAnyArrayBuffer = ce, [
            `isProxy`,
            `isExternal`,
            `isModuleNamespaceObject`
        ].forEach(function(t) {
            Object.defineProperty(e, t, {
                enumerable: !1,
                value: function() {
                    throw Error(t + ` is not supported in userland`);
                }
            });
        });
    })), _e = t(((e, t)=>{
        t.exports = function(e) {
            return e && typeof e == `object` && typeof e.copy == `function` && typeof e.fill == `function` && typeof e.readUInt8 == `function`;
        };
    })), ve = t(((e)=>{
        var t = Object.getOwnPropertyDescriptors || function(e) {
            for(var t = Object.keys(e), n = {}, r = 0; r < t.length; r++)n[t[r]] = Object.getOwnPropertyDescriptor(e, t[r]);
            return n;
        }, n = /%[sdj%]/g;
        e.format = function(e) {
            if (!S(e)) {
                for(var t = [], r = 0; r < arguments.length; r++)t.push(o(arguments[r]));
                return t.join(` `);
            }
            for(var r = 1, i = arguments, a = i.length, s = String(e).replace(n, function(e) {
                if (e === `%%`) return `%`;
                if (r >= a) return e;
                switch(e){
                    case `%s`:
                        return String(i[r++]);
                    case `%d`:
                        return Number(i[r++]);
                    case `%j`:
                        try {
                            return JSON.stringify(i[r++]);
                        } catch  {
                            return `[Circular]`;
                        }
                    default:
                        return e;
                }
            }), c = i[r]; r < a; c = i[++r])y(c) || !E(c) ? s += ` ` + c : s += ` ` + o(c);
            return s;
        }, e.deprecate = function(t, n) {
            if (typeof process < `u` && process.noDeprecation === !0) return t;
            if (typeof process > `u`) return function() {
                return e.deprecate(t, n).apply(this, arguments);
            };
            var r = !1;
            function i() {
                if (!r) {
                    if (process.throwDeprecation) throw Error(n);
                    process.traceDeprecation ? console.trace(n) : console.error(n), r = !0;
                }
                return t.apply(this, arguments);
            }
            return i;
        };
        var r = {}, i = /^$/;
        if ({}.NODE_DEBUG) {
            var a = {}.NODE_DEBUG;
            a = a.replace(/[|\\{}()[\]^$+?.]/g, `\\$&`).replace(/\*/g, `.*`).replace(/,/g, `$|^`).toUpperCase(), i = RegExp(`^` + a + `$`, `i`);
        }
        e.debuglog = function(t) {
            if (t = t.toUpperCase(), !r[t]) if (i.test(t)) {
                var n = process.pid;
                r[t] = function() {
                    var r = e.format.apply(e, arguments);
                    console.error(`%s %d: %s`, t, n, r);
                };
            } else r[t] = function() {};
            return r[t];
        };
        function o(t, n) {
            var r = {
                seen: [],
                stylize: c
            };
            return arguments.length >= 3 && (r.depth = arguments[2]), arguments.length >= 4 && (r.colors = arguments[3]), v(n) ? r.showHidden = n : n && e._extend(r, n), w(r.showHidden) && (r.showHidden = !1), w(r.depth) && (r.depth = 2), w(r.colors) && (r.colors = !1), w(r.customInspect) && (r.customInspect = !0), r.colors && (r.stylize = s), u(r, t, r.depth);
        }
        e.inspect = o, o.colors = {
            bold: [
                1,
                22
            ],
            italic: [
                3,
                23
            ],
            underline: [
                4,
                24
            ],
            inverse: [
                7,
                27
            ],
            white: [
                37,
                39
            ],
            grey: [
                90,
                39
            ],
            black: [
                30,
                39
            ],
            blue: [
                34,
                39
            ],
            cyan: [
                36,
                39
            ],
            green: [
                32,
                39
            ],
            magenta: [
                35,
                39
            ],
            red: [
                31,
                39
            ],
            yellow: [
                33,
                39
            ]
        }, o.styles = {
            special: `cyan`,
            number: `yellow`,
            boolean: `yellow`,
            undefined: `grey`,
            null: `bold`,
            string: `green`,
            date: `magenta`,
            regexp: `red`
        };
        function s(e, t) {
            var n = o.styles[t];
            return n ? `\x1B[` + o.colors[n][0] + `m` + e + `\x1B[` + o.colors[n][1] + `m` : e;
        }
        function c(e, t) {
            return e;
        }
        function l(e) {
            var t = {};
            return e.forEach(function(e, n) {
                t[e] = !0;
            }), t;
        }
        function u(t, n, r) {
            if (t.customInspect && n && k(n.inspect) && n.inspect !== e.inspect && !(n.constructor && n.constructor.prototype === n)) {
                var i = n.inspect(r, t);
                return S(i) || (i = u(t, i, r)), i;
            }
            var a = d(t, n);
            if (a) return a;
            var o = Object.keys(n), s = l(o);
            if (t.showHidden && (o = Object.getOwnPropertyNames(n)), O(n) && (o.indexOf(`message`) >= 0 || o.indexOf(`description`) >= 0)) return f(n);
            if (o.length === 0) {
                if (k(n)) {
                    var c = n.name ? `: ` + n.name : ``;
                    return t.stylize(`[Function` + c + `]`, `special`);
                }
                if (T(n)) return t.stylize(RegExp.prototype.toString.call(n), `regexp`);
                if (D(n)) return t.stylize(Date.prototype.toString.call(n), `date`);
                if (O(n)) return f(n);
            }
            var g = ``, v = !1, y = [
                `{`,
                `}`
            ];
            if (_(n) && (v = !0, y = [
                `[`,
                `]`
            ]), k(n) && (g = ` [Function` + (n.name ? `: ` + n.name : ``) + `]`), T(n) && (g = ` ` + RegExp.prototype.toString.call(n)), D(n) && (g = ` ` + Date.prototype.toUTCString.call(n)), O(n) && (g = ` ` + f(n)), o.length === 0 && (!v || n.length == 0)) return y[0] + g + y[1];
            if (r < 0) return T(n) ? t.stylize(RegExp.prototype.toString.call(n), `regexp`) : t.stylize(`[Object]`, `special`);
            t.seen.push(n);
            var b = v ? p(t, n, r, s, o) : o.map(function(e) {
                return m(t, n, r, s, e, v);
            });
            return t.seen.pop(), h(b, g, y);
        }
        function d(e, t) {
            if (w(t)) return e.stylize(`undefined`, `undefined`);
            if (S(t)) {
                var n = `'` + JSON.stringify(t).replace(/^"|"$/g, ``).replace(/'/g, `\\'`).replace(/\\"/g, `"`) + `'`;
                return e.stylize(n, `string`);
            }
            if (x(t)) return e.stylize(`` + t, `number`);
            if (v(t)) return e.stylize(`` + t, `boolean`);
            if (y(t)) return e.stylize(`null`, `null`);
        }
        function f(e) {
            return `[` + Error.prototype.toString.call(e) + `]`;
        }
        function p(e, t, n, r, i) {
            for(var a = [], o = 0, s = t.length; o < s; ++o)ee(t, String(o)) ? a.push(m(e, t, n, r, String(o), !0)) : a.push(``);
            return i.forEach(function(i) {
                i.match(/^\d+$/) || a.push(m(e, t, n, r, i, !0));
            }), a;
        }
        function m(e, t, n, r, i, a) {
            var o, s, c = Object.getOwnPropertyDescriptor(t, i) || {
                value: t[i]
            };
            if (c.get ? s = c.set ? e.stylize(`[Getter/Setter]`, `special`) : e.stylize(`[Getter]`, `special`) : c.set && (s = e.stylize(`[Setter]`, `special`)), ee(r, i) || (o = `[` + i + `]`), s || (e.seen.indexOf(c.value) < 0 ? (s = y(n) ? u(e, c.value, null) : u(e, c.value, n - 1), s.indexOf(`
`) > -1 && (s = a ? s.split(`
`).map(function(e) {
                return `  ` + e;
            }).join(`
`).slice(2) : `
` + s.split(`
`).map(function(e) {
                return `   ` + e;
            }).join(`
`))) : s = e.stylize(`[Circular]`, `special`)), w(o)) {
                if (a && i.match(/^\d+$/)) return s;
                o = JSON.stringify(`` + i), o.match(/^"([a-zA-Z_][a-zA-Z_0-9]*)"$/) ? (o = o.slice(1, -1), o = e.stylize(o, `name`)) : (o = o.replace(/'/g, `\\'`).replace(/\\"/g, `"`).replace(/(^"|"$)/g, `'`), o = e.stylize(o, `string`));
            }
            return o + `: ` + s;
        }
        function h(e, t, n) {
            var r = 0;
            return e.reduce(function(e, t) {
                return r++, t.indexOf(`
`) >= 0 && r++, e + t.replace(/\u001b\[\d\d?m/g, ``).length + 1;
            }, 0) > 60 ? n[0] + (t === `` ? `` : t + `
 `) + ` ` + e.join(`,
  `) + ` ` + n[1] : n[0] + t + ` ` + e.join(`, `) + ` ` + n[1];
        }
        e.types = ge();
        function _(e) {
            return Array.isArray(e);
        }
        e.isArray = _;
        function v(e) {
            return typeof e == `boolean`;
        }
        e.isBoolean = v;
        function y(e) {
            return e === null;
        }
        e.isNull = y;
        function b(e) {
            return e == null;
        }
        e.isNullOrUndefined = b;
        function x(e) {
            return typeof e == `number`;
        }
        e.isNumber = x;
        function S(e) {
            return typeof e == `string`;
        }
        e.isString = S;
        function C(e) {
            return typeof e == `symbol`;
        }
        e.isSymbol = C;
        function w(e) {
            return e === void 0;
        }
        e.isUndefined = w;
        function T(e) {
            return E(e) && j(e) === `[object RegExp]`;
        }
        e.isRegExp = T, e.types.isRegExp = T;
        function E(e) {
            return typeof e == `object` && !!e;
        }
        e.isObject = E;
        function D(e) {
            return E(e) && j(e) === `[object Date]`;
        }
        e.isDate = D, e.types.isDate = D;
        function O(e) {
            return E(e) && (j(e) === `[object Error]` || e instanceof Error);
        }
        e.isError = O, e.types.isNativeError = O;
        function k(e) {
            return typeof e == `function`;
        }
        e.isFunction = k;
        function A(e) {
            return e === null || typeof e == `boolean` || typeof e == `number` || typeof e == `string` || typeof e == `symbol` || e === void 0;
        }
        e.isPrimitive = A, e.isBuffer = _e();
        function j(e) {
            return Object.prototype.toString.call(e);
        }
        function M(e) {
            return e < 10 ? `0` + e.toString(10) : e.toString(10);
        }
        var N = [
            `Jan`,
            `Feb`,
            `Mar`,
            `Apr`,
            `May`,
            `Jun`,
            `Jul`,
            `Aug`,
            `Sep`,
            `Oct`,
            `Nov`,
            `Dec`
        ];
        function P() {
            var e = new Date, t = [
                M(e.getHours()),
                M(e.getMinutes()),
                M(e.getSeconds())
            ].join(`:`);
            return [
                e.getDate(),
                N[e.getMonth()],
                t
            ].join(` `);
        }
        e.log = function() {
            console.log(`%s - %s`, P(), e.format.apply(e, arguments));
        }, e.inherits = g(), e._extend = function(e, t) {
            if (!t || !E(t)) return e;
            for(var n = Object.keys(t), r = n.length; r--;)e[n[r]] = t[n[r]];
            return e;
        };
        function ee(e, t) {
            return Object.prototype.hasOwnProperty.call(e, t);
        }
        var F = typeof Symbol < `u` ? Symbol(`util.promisify.custom`) : void 0;
        e.promisify = function(e) {
            if (typeof e != `function`) throw TypeError(`The "original" argument must be of type Function`);
            if (F && e[F]) {
                var n = e[F];
                if (typeof n != `function`) throw TypeError(`The "util.promisify.custom" argument must be of type Function`);
                return Object.defineProperty(n, F, {
                    value: n,
                    enumerable: !1,
                    writable: !1,
                    configurable: !0
                }), n;
            }
            function n() {
                for(var t, n, r = new Promise(function(e, r) {
                    t = e, n = r;
                }), i = [], a = 0; a < arguments.length; a++)i.push(arguments[a]);
                i.push(function(e, r) {
                    e ? n(e) : t(r);
                });
                try {
                    e.apply(this, i);
                } catch (e) {
                    n(e);
                }
                return r;
            }
            return Object.setPrototypeOf(n, Object.getPrototypeOf(e)), F && Object.defineProperty(n, F, {
                value: n,
                enumerable: !1,
                writable: !1,
                configurable: !0
            }), Object.defineProperties(n, t(e));
        }, e.promisify.custom = F;
        function te(e, t) {
            if (!e) {
                var n = Error(`Promise was rejected with a falsy value`);
                n.reason = e, e = n;
            }
            return t(e);
        }
        function I(e) {
            if (typeof e != `function`) throw TypeError(`The "original" argument must be of type Function`);
            function n() {
                for(var t = [], n = 0; n < arguments.length; n++)t.push(arguments[n]);
                var r = t.pop();
                if (typeof r != `function`) throw TypeError(`The last argument must be of type Function`);
                var i = this, a = function() {
                    return r.apply(i, arguments);
                };
                e.apply(this, t).then(function(e) {
                    process.nextTick(a.bind(null, null, e));
                }, function(e) {
                    process.nextTick(te.bind(null, e, a));
                });
            }
            return Object.setPrototypeOf(n, Object.getPrototypeOf(e)), Object.defineProperties(n, t(e)), n;
        }
        e.callbackify = I;
    })), ye = t(((e, t)=>{
        function n(e, t) {
            var n = Object.keys(e);
            if (Object.getOwnPropertySymbols) {
                var r = Object.getOwnPropertySymbols(e);
                t && (r = r.filter(function(t) {
                    return Object.getOwnPropertyDescriptor(e, t).enumerable;
                })), n.push.apply(n, r);
            }
            return n;
        }
        function r(e) {
            for(var t = 1; t < arguments.length; t++){
                var r = arguments[t] == null ? {} : arguments[t];
                t % 2 ? n(Object(r), !0).forEach(function(t) {
                    i(e, t, r[t]);
                }) : Object.getOwnPropertyDescriptors ? Object.defineProperties(e, Object.getOwnPropertyDescriptors(r)) : n(Object(r)).forEach(function(t) {
                    Object.defineProperty(e, t, Object.getOwnPropertyDescriptor(r, t));
                });
            }
            return e;
        }
        function i(e, t, n) {
            return t = c(t), t in e ? Object.defineProperty(e, t, {
                value: n,
                enumerable: !0,
                configurable: !0,
                writable: !0
            }) : e[t] = n, e;
        }
        function a(e, t) {
            if (!(e instanceof t)) throw TypeError(`Cannot call a class as a function`);
        }
        function o(e, t) {
            for(var n = 0; n < t.length; n++){
                var r = t[n];
                r.enumerable = r.enumerable || !1, r.configurable = !0, `value` in r && (r.writable = !0), Object.defineProperty(e, c(r.key), r);
            }
        }
        function s(e, t, n) {
            return t && o(e.prototype, t), n && o(e, n), Object.defineProperty(e, "prototype", {
                writable: !1
            }), e;
        }
        function c(e) {
            var t = l(e, `string`);
            return typeof t == `symbol` ? t : String(t);
        }
        function l(e, t) {
            if (typeof e != `object` || !e) return e;
            var n = e[Symbol.toPrimitive];
            if (n !== void 0) {
                var r = n.call(e, t || `default`);
                if (typeof r != `object`) return r;
                throw TypeError(`@@toPrimitive must return a primitive value.`);
            }
            return (t === `string` ? String : Number)(e);
        }
        var u = m().Buffer, d = ve().inspect, f = d && d.custom || `inspect`;
        function p(e, t, n) {
            u.prototype.copy.call(e, t, n);
        }
        t.exports = function() {
            function e() {
                a(this, e), this.head = null, this.tail = null, this.length = 0;
            }
            return s(e, [
                {
                    key: `push`,
                    value: function(e) {
                        var t = {
                            data: e,
                            next: null
                        };
                        this.length > 0 ? this.tail.next = t : this.head = t, this.tail = t, ++this.length;
                    }
                },
                {
                    key: `unshift`,
                    value: function(e) {
                        var t = {
                            data: e,
                            next: this.head
                        };
                        this.length === 0 && (this.tail = t), this.head = t, ++this.length;
                    }
                },
                {
                    key: `shift`,
                    value: function() {
                        if (this.length !== 0) {
                            var e = this.head.data;
                            return this.length === 1 ? this.head = this.tail = null : this.head = this.head.next, --this.length, e;
                        }
                    }
                },
                {
                    key: `clear`,
                    value: function() {
                        this.head = this.tail = null, this.length = 0;
                    }
                },
                {
                    key: `join`,
                    value: function(e) {
                        if (this.length === 0) return ``;
                        for(var t = this.head, n = `` + t.data; t = t.next;)n += e + t.data;
                        return n;
                    }
                },
                {
                    key: `concat`,
                    value: function(e) {
                        if (this.length === 0) return u.alloc(0);
                        for(var t = u.allocUnsafe(e >>> 0), n = this.head, r = 0; n;)p(n.data, t, r), r += n.data.length, n = n.next;
                        return t;
                    }
                },
                {
                    key: `consume`,
                    value: function(e, t) {
                        var n;
                        return e < this.head.data.length ? (n = this.head.data.slice(0, e), this.head.data = this.head.data.slice(e)) : n = e === this.head.data.length ? this.shift() : t ? this._getString(e) : this._getBuffer(e), n;
                    }
                },
                {
                    key: `first`,
                    value: function() {
                        return this.head.data;
                    }
                },
                {
                    key: `_getString`,
                    value: function(e) {
                        var t = this.head, n = 1, r = t.data;
                        for(e -= r.length; t = t.next;){
                            var i = t.data, a = e > i.length ? i.length : e;
                            if (a === i.length ? r += i : r += i.slice(0, e), e -= a, e === 0) {
                                a === i.length ? (++n, t.next ? this.head = t.next : this.head = this.tail = null) : (this.head = t, t.data = i.slice(a));
                                break;
                            }
                            ++n;
                        }
                        return this.length -= n, r;
                    }
                },
                {
                    key: `_getBuffer`,
                    value: function(e) {
                        var t = u.allocUnsafe(e), n = this.head, r = 1;
                        for(n.data.copy(t), e -= n.data.length; n = n.next;){
                            var i = n.data, a = e > i.length ? i.length : e;
                            if (i.copy(t, t.length - e, 0, a), e -= a, e === 0) {
                                a === i.length ? (++r, n.next ? this.head = n.next : this.head = this.tail = null) : (this.head = n, n.data = i.slice(a));
                                break;
                            }
                            ++r;
                        }
                        return this.length -= r, t;
                    }
                },
                {
                    key: f,
                    value: function(e, t) {
                        return d(this, r(r({}, t), {}, {
                            depth: 0,
                            customInspect: !1
                        }));
                    }
                }
            ]), e;
        }();
    })), be = t(((e, t)=>{
        function n(e, t) {
            var n = this, a = this._readableState && this._readableState.destroyed, s = this._writableState && this._writableState.destroyed;
            return a || s ? (t ? t(e) : e && (this._writableState ? this._writableState.errorEmitted || (this._writableState.errorEmitted = !0, process.nextTick(o, this, e)) : process.nextTick(o, this, e)), this) : (this._readableState && (this._readableState.destroyed = !0), this._writableState && (this._writableState.destroyed = !0), this._destroy(e || null, function(e) {
                !t && e ? n._writableState ? n._writableState.errorEmitted ? process.nextTick(i, n) : (n._writableState.errorEmitted = !0, process.nextTick(r, n, e)) : process.nextTick(r, n, e) : t ? (process.nextTick(i, n), t(e)) : process.nextTick(i, n);
            }), this);
        }
        function r(e, t) {
            o(e, t), i(e);
        }
        function i(e) {
            e._writableState && !e._writableState.emitClose || e._readableState && !e._readableState.emitClose || e.emit(`close`);
        }
        function a() {
            this._readableState && (this._readableState.destroyed = !1, this._readableState.reading = !1, this._readableState.ended = !1, this._readableState.endEmitted = !1), this._writableState && (this._writableState.destroyed = !1, this._writableState.ended = !1, this._writableState.ending = !1, this._writableState.finalCalled = !1, this._writableState.prefinished = !1, this._writableState.finished = !1, this._writableState.errorEmitted = !1);
        }
        function o(e, t) {
            e.emit(`error`, t);
        }
        function s(e, t) {
            var n = e._readableState, r = e._writableState;
            n && n.autoDestroy || r && r.autoDestroy ? e.destroy(t) : e.emit(`error`, t);
        }
        t.exports = {
            destroy: n,
            undestroy: a,
            errorOrDestroy: s
        };
    })), xe = t(((e, t)=>{
        function n(e, t) {
            e.prototype = Object.create(t.prototype), e.prototype.constructor = e, e.__proto__ = t;
        }
        var r = {};
        function i(e, t, i) {
            i ||= Error;
            function a(e, n, r) {
                return typeof t == `string` ? t : t(e, n, r);
            }
            var o = function(e) {
                n(t, e);
                function t(t, n, r) {
                    return e.call(this, a(t, n, r)) || this;
                }
                return t;
            }(i);
            o.prototype.name = i.name, o.prototype.code = e, r[e] = o;
        }
        function a(e, t) {
            if (Array.isArray(e)) {
                var n = e.length;
                return e = e.map(function(e) {
                    return String(e);
                }), n > 2 ? `one of ${t} ${e.slice(0, n - 1).join(`, `)}, or ` + e[n - 1] : n === 2 ? `one of ${t} ${e[0]} or ${e[1]}` : `of ${t} ${e[0]}`;
            } else return `of ${t} ${String(e)}`;
        }
        function o(e, t, n) {
            return e.substr(!n || n < 0 ? 0 : +n, t.length) === t;
        }
        function s(e, t, n) {
            return (n === void 0 || n > e.length) && (n = e.length), e.substring(n - t.length, n) === t;
        }
        function c(e, t, n) {
            return typeof n != `number` && (n = 0), n + t.length > e.length ? !1 : e.indexOf(t, n) !== -1;
        }
        i(`ERR_INVALID_OPT_VALUE`, function(e, t) {
            return `The value "` + t + `" is invalid for option "` + e + `"`;
        }, TypeError), i(`ERR_INVALID_ARG_TYPE`, function(e, t, n) {
            var r;
            typeof t == `string` && o(t, `not `) ? (r = `must not be`, t = t.replace(/^not /, ``)) : r = `must be`;
            var i = s(e, ` argument`) ? `The ${e} ${r} ${a(t, `type`)}` : `The "${e}" ${c(e, `.`) ? `property` : `argument`} ${r} ${a(t, `type`)}`;
            return i += `. Received type ${typeof n}`, i;
        }, TypeError), i(`ERR_STREAM_PUSH_AFTER_EOF`, `stream.push() after EOF`), i(`ERR_METHOD_NOT_IMPLEMENTED`, function(e) {
            return `The ` + e + ` method is not implemented`;
        }), i(`ERR_STREAM_PREMATURE_CLOSE`, `Premature close`), i(`ERR_STREAM_DESTROYED`, function(e) {
            return `Cannot call ` + e + ` after a stream was destroyed`;
        }), i(`ERR_MULTIPLE_CALLBACK`, `Callback called multiple times`), i(`ERR_STREAM_CANNOT_PIPE`, `Cannot pipe, not readable`), i(`ERR_STREAM_WRITE_AFTER_END`, `write after end`), i(`ERR_STREAM_NULL_VALUES`, `May not write null values to stream`, TypeError), i(`ERR_UNKNOWN_ENCODING`, function(e) {
            return `Unknown encoding: ` + e;
        }, TypeError), i(`ERR_STREAM_UNSHIFT_AFTER_END_EVENT`, `stream.unshift() after end event`), t.exports.codes = r;
    })), Se = t(((e, t)=>{
        var n = xe().codes.ERR_INVALID_OPT_VALUE;
        function r(e, t, n) {
            return e.highWaterMark == null ? t ? e[n] : null : e.highWaterMark;
        }
        function i(e, t, i, a) {
            var o = r(t, a, i);
            if (o != null) {
                if (!(isFinite(o) && Math.floor(o) === o) || o < 0) throw new n(a ? i : `highWaterMark`, o);
                return Math.floor(o);
            }
            return e.objectMode ? 16 : 16 * 1024;
        }
        t.exports = {
            getHighWaterMark: i
        };
    })), Ce = t(((e, t)=>{
        t.exports = n;
        function n(e, t) {
            if (r(`noDeprecation`)) return e;
            var n = !1;
            function i() {
                if (!n) {
                    if (r(`throwDeprecation`)) throw Error(t);
                    r(`traceDeprecation`) ? console.trace(t) : console.warn(t), n = !0;
                }
                return e.apply(this, arguments);
            }
            return i;
        }
        function r(e) {
            try {
                if (!globalThis.localStorage) return !1;
            } catch  {
                return !1;
            }
            var t = globalThis.localStorage[e];
            return t == null ? !1 : String(t).toLowerCase() === `true`;
        }
    })), we = t(((e, t)=>{
        t.exports = O;
        function n(e) {
            var t = this;
            this.next = null, this.entry = null, this.finish = function() {
                ie(t, e);
            };
        }
        var r;
        O.WritableState = E;
        var i = {
            deprecate: Ce()
        }, a = _(), o = m().Buffer, s = (typeof globalThis < `u` ? globalThis : typeof window < `u` ? window : typeof self < `u` ? self : {}).Uint8Array || function() {};
        function c(e) {
            return o.from(e);
        }
        function l(e) {
            return o.isBuffer(e) || e instanceof s;
        }
        var u = be(), d = Se().getHighWaterMark, f = xe().codes, p = f.ERR_INVALID_ARG_TYPE, h = f.ERR_METHOD_NOT_IMPLEMENTED, v = f.ERR_MULTIPLE_CALLBACK, y = f.ERR_STREAM_CANNOT_PIPE, b = f.ERR_STREAM_DESTROYED, x = f.ERR_STREAM_NULL_VALUES, S = f.ERR_STREAM_WRITE_AFTER_END, C = f.ERR_UNKNOWN_ENCODING, w = u.errorOrDestroy;
        g()(O, a);
        function T() {}
        function E(e, t, i) {
            r ||= Te(), e ||= {}, typeof i != `boolean` && (i = t instanceof r), this.objectMode = !!e.objectMode, i && (this.objectMode = this.objectMode || !!e.writableObjectMode), this.highWaterMark = d(this, e, `writableHighWaterMark`, i), this.finalCalled = !1, this.needDrain = !1, this.ending = !1, this.ended = !1, this.finished = !1, this.destroyed = !1;
            var a = e.decodeStrings === !1;
            this.decodeStrings = !a, this.defaultEncoding = e.defaultEncoding || `utf8`, this.length = 0, this.writing = !1, this.corked = 0, this.sync = !0, this.bufferProcessing = !1, this.onwrite = function(e) {
                F(t, e);
            }, this.writecb = null, this.writelen = 0, this.bufferedRequest = null, this.lastBufferedRequest = null, this.pendingcb = 0, this.prefinished = !1, this.errorEmitted = !1, this.emitClose = e.emitClose !== !1, this.autoDestroy = !!e.autoDestroy, this.bufferedRequestCount = 0, this.corkedRequestsFree = new n(this);
        }
        E.prototype.getBuffer = function() {
            for(var e = this.bufferedRequest, t = []; e;)t.push(e), e = e.next;
            return t;
        }, (function() {
            try {
                Object.defineProperty(E.prototype, "buffer", {
                    get: i.deprecate(function() {
                        return this.getBuffer();
                    }, `_writableState.buffer is deprecated. Use _writableState.getBuffer instead.`, `DEP0003`)
                });
            } catch  {}
        })();
        var D;
        typeof Symbol == `function` && Symbol.hasInstance && typeof Function.prototype[Symbol.hasInstance] == `function` ? (D = Function.prototype[Symbol.hasInstance], Object.defineProperty(O, Symbol.hasInstance, {
            value: function(e) {
                return D.call(this, e) ? !0 : this === O ? e && e._writableState instanceof E : !1;
            }
        })) : D = function(e) {
            return e instanceof this;
        };
        function O(e) {
            r ||= Te();
            var t = this instanceof r;
            if (!t && !D.call(O, this)) return new O(e);
            this._writableState = new E(e, this, t), this.writable = !0, e && (typeof e.write == `function` && (this._write = e.write), typeof e.writev == `function` && (this._writev = e.writev), typeof e.destroy == `function` && (this._destroy = e.destroy), typeof e.final == `function` && (this._final = e.final)), a.call(this);
        }
        O.prototype.pipe = function() {
            w(this, new y);
        };
        function k(e, t) {
            var n = new S;
            w(e, n), process.nextTick(t, n);
        }
        function A(e, t, n, r) {
            var i;
            return n === null ? i = new x : typeof n != `string` && !t.objectMode && (i = new p(`chunk`, [
                `string`,
                `Buffer`
            ], n)), i ? (w(e, i), process.nextTick(r, i), !1) : !0;
        }
        O.prototype.write = function(e, t, n) {
            var r = this._writableState, i = !1, a = !r.objectMode && l(e);
            return a && !o.isBuffer(e) && (e = c(e)), typeof t == `function` && (n = t, t = null), a ? t = `buffer` : t ||= r.defaultEncoding, typeof n != `function` && (n = T), r.ending ? k(this, n) : (a || A(this, r, e, n)) && (r.pendingcb++, i = M(this, r, a, e, t, n)), i;
        }, O.prototype.cork = function() {
            this._writableState.corked++;
        }, O.prototype.uncork = function() {
            var e = this._writableState;
            e.corked && (e.corked--, !e.writing && !e.corked && !e.bufferProcessing && e.bufferedRequest && ne(this, e));
        }, O.prototype.setDefaultEncoding = function(e) {
            if (typeof e == `string` && (e = e.toLowerCase()), !([
                `hex`,
                `utf8`,
                `utf-8`,
                `ascii`,
                `binary`,
                `base64`,
                `ucs2`,
                `ucs-2`,
                `utf16le`,
                `utf-16le`,
                `raw`
            ].indexOf((e + ``).toLowerCase()) > -1)) throw new C(e);
            return this._writableState.defaultEncoding = e, this;
        }, Object.defineProperty(O.prototype, "writableBuffer", {
            enumerable: !1,
            get: function() {
                return this._writableState && this._writableState.getBuffer();
            }
        });
        function j(e, t, n) {
            return !e.objectMode && e.decodeStrings !== !1 && typeof t == `string` && (t = o.from(t, n)), t;
        }
        Object.defineProperty(O.prototype, "writableHighWaterMark", {
            enumerable: !1,
            get: function() {
                return this._writableState.highWaterMark;
            }
        });
        function M(e, t, n, r, i, a) {
            if (!n) {
                var o = j(t, r, i);
                r !== o && (n = !0, i = `buffer`, r = o);
            }
            var s = t.objectMode ? 1 : r.length;
            t.length += s;
            var c = t.length < t.highWaterMark;
            if (c || (t.needDrain = !0), t.writing || t.corked) {
                var l = t.lastBufferedRequest;
                t.lastBufferedRequest = {
                    chunk: r,
                    encoding: i,
                    isBuf: n,
                    callback: a,
                    next: null
                }, l ? l.next = t.lastBufferedRequest : t.bufferedRequest = t.lastBufferedRequest, t.bufferedRequestCount += 1;
            } else N(e, t, !1, s, r, i, a);
            return c;
        }
        function N(e, t, n, r, i, a, o) {
            t.writelen = r, t.writecb = o, t.writing = !0, t.sync = !0, t.destroyed ? t.onwrite(new b(`write`)) : n ? e._writev(i, t.onwrite) : e._write(i, a, t.onwrite), t.sync = !1;
        }
        function P(e, t, n, r, i) {
            --t.pendingcb, n ? (process.nextTick(i, r), process.nextTick(z, e, t), e._writableState.errorEmitted = !0, w(e, r)) : (i(r), e._writableState.errorEmitted = !0, w(e, r), z(e, t));
        }
        function ee(e) {
            e.writing = !1, e.writecb = null, e.length -= e.writelen, e.writelen = 0;
        }
        function F(e, t) {
            var n = e._writableState, r = n.sync, i = n.writecb;
            if (typeof i != `function`) throw new v;
            if (ee(n), t) P(e, n, r, t, i);
            else {
                var a = re(n) || e.destroyed;
                !a && !n.corked && !n.bufferProcessing && n.bufferedRequest && ne(e, n), r ? process.nextTick(te, e, n, a, i) : te(e, n, a, i);
            }
        }
        function te(e, t, n, r) {
            n || I(e, t), t.pendingcb--, r(), z(e, t);
        }
        function I(e, t) {
            t.length === 0 && t.needDrain && (t.needDrain = !1, e.emit(`drain`));
        }
        function ne(e, t) {
            t.bufferProcessing = !0;
            var r = t.bufferedRequest;
            if (e._writev && r && r.next) {
                var i = t.bufferedRequestCount, a = Array(i), o = t.corkedRequestsFree;
                o.entry = r;
                for(var s = 0, c = !0; r;)a[s] = r, r.isBuf || (c = !1), r = r.next, s += 1;
                a.allBuffers = c, N(e, t, !0, t.length, a, ``, o.finish), t.pendingcb++, t.lastBufferedRequest = null, o.next ? (t.corkedRequestsFree = o.next, o.next = null) : t.corkedRequestsFree = new n(t), t.bufferedRequestCount = 0;
            } else {
                for(; r;){
                    var l = r.chunk, u = r.encoding, d = r.callback;
                    if (N(e, t, !1, t.objectMode ? 1 : l.length, l, u, d), r = r.next, t.bufferedRequestCount--, t.writing) break;
                }
                r === null && (t.lastBufferedRequest = null);
            }
            t.bufferedRequest = r, t.bufferProcessing = !1;
        }
        O.prototype._write = function(e, t, n) {
            n(new h(`_write()`));
        }, O.prototype._writev = null, O.prototype.end = function(e, t, n) {
            var r = this._writableState;
            return typeof e == `function` ? (n = e, e = null, t = null) : typeof t == `function` && (n = t, t = null), e != null && this.write(e, t), r.corked && (r.corked = 1, this.uncork()), r.ending || B(this, r, n), this;
        }, Object.defineProperty(O.prototype, "writableLength", {
            enumerable: !1,
            get: function() {
                return this._writableState.length;
            }
        });
        function re(e) {
            return e.ending && e.length === 0 && e.bufferedRequest === null && !e.finished && !e.writing;
        }
        function L(e, t) {
            e._final(function(n) {
                t.pendingcb--, n && w(e, n), t.prefinished = !0, e.emit(`prefinish`), z(e, t);
            });
        }
        function R(e, t) {
            !t.prefinished && !t.finalCalled && (typeof e._final == `function` && !t.destroyed ? (t.pendingcb++, t.finalCalled = !0, process.nextTick(L, e, t)) : (t.prefinished = !0, e.emit(`prefinish`)));
        }
        function z(e, t) {
            var n = re(t);
            if (n && (R(e, t), t.pendingcb === 0 && (t.finished = !0, e.emit(`finish`), t.autoDestroy))) {
                var r = e._readableState;
                (!r || r.autoDestroy && r.endEmitted) && e.destroy();
            }
            return n;
        }
        function B(e, t, n) {
            t.ending = !0, z(e, t), n && (t.finished ? process.nextTick(n) : e.once(`finish`, n)), t.ended = !0, e.writable = !1;
        }
        function ie(e, t, n) {
            var r = e.entry;
            for(e.entry = null; r;){
                var i = r.callback;
                t.pendingcb--, i(n), r = r.next;
            }
            t.corkedRequestsFree.next = e;
        }
        Object.defineProperty(O.prototype, "destroyed", {
            enumerable: !1,
            get: function() {
                return this._writableState === void 0 ? !1 : this._writableState.destroyed;
            },
            set: function(e) {
                this._writableState && (this._writableState.destroyed = e);
            }
        }), O.prototype.destroy = u.destroy, O.prototype._undestroy = u.undestroy, O.prototype._destroy = function(e, t) {
            t(e);
        };
    })), Te = t(((e, t)=>{
        var n = Object.keys || function(e) {
            var t = [];
            for(var n in e)t.push(n);
            return t;
        };
        t.exports = c;
        var r = je(), i = we();
        g()(c, r);
        for(var a = n(i.prototype), o = 0; o < a.length; o++){
            var s = a[o];
            c.prototype[s] || (c.prototype[s] = i.prototype[s]);
        }
        function c(e) {
            if (!(this instanceof c)) return new c(e);
            r.call(this, e), i.call(this, e), this.allowHalfOpen = !0, e && (e.readable === !1 && (this.readable = !1), e.writable === !1 && (this.writable = !1), e.allowHalfOpen === !1 && (this.allowHalfOpen = !1, this.once(`end`, l)));
        }
        Object.defineProperty(c.prototype, "writableHighWaterMark", {
            enumerable: !1,
            get: function() {
                return this._writableState.highWaterMark;
            }
        }), Object.defineProperty(c.prototype, "writableBuffer", {
            enumerable: !1,
            get: function() {
                return this._writableState && this._writableState.getBuffer();
            }
        }), Object.defineProperty(c.prototype, "writableLength", {
            enumerable: !1,
            get: function() {
                return this._writableState.length;
            }
        });
        function l() {
            this._writableState.ended || process.nextTick(u, this);
        }
        function u(e) {
            e.end();
        }
        Object.defineProperty(c.prototype, "destroyed", {
            enumerable: !1,
            get: function() {
                return this._readableState === void 0 || this._writableState === void 0 ? !1 : this._readableState.destroyed && this._writableState.destroyed;
            },
            set: function(e) {
                this._readableState === void 0 || this._writableState === void 0 || (this._readableState.destroyed = e, this._writableState.destroyed = e);
            }
        });
    })), Ee = t(((e, t)=>{
        var n = m(), r = n.Buffer;
        function i(e, t) {
            for(var n in e)t[n] = e[n];
        }
        r.from && r.alloc && r.allocUnsafe && r.allocUnsafeSlow ? t.exports = n : (i(n, e), e.Buffer = a);
        function a(e, t, n) {
            return r(e, t, n);
        }
        a.prototype = Object.create(r.prototype), i(r, a), a.from = function(e, t, n) {
            if (typeof e == `number`) throw TypeError(`Argument must not be a number`);
            return r(e, t, n);
        }, a.alloc = function(e, t, n) {
            if (typeof e != `number`) throw TypeError(`Argument must be a number`);
            var i = r(e);
            return t === void 0 ? i.fill(0) : typeof n == `string` ? i.fill(t, n) : i.fill(t), i;
        }, a.allocUnsafe = function(e) {
            if (typeof e != `number`) throw TypeError(`Argument must be a number`);
            return r(e);
        }, a.allocUnsafeSlow = function(e) {
            if (typeof e != `number`) throw TypeError(`Argument must be a number`);
            return n.SlowBuffer(e);
        };
    })), De = t(((e)=>{
        var t = Ee().Buffer, n = t.isEncoding || function(e) {
            switch(e = `` + e, e && e.toLowerCase()){
                case `hex`:
                case `utf8`:
                case `utf-8`:
                case `ascii`:
                case `binary`:
                case `base64`:
                case `ucs2`:
                case `ucs-2`:
                case `utf16le`:
                case `utf-16le`:
                case `raw`:
                    return !0;
                default:
                    return !1;
            }
        };
        function r(e) {
            if (!e) return `utf8`;
            for(var t;;)switch(e){
                case `utf8`:
                case `utf-8`:
                    return `utf8`;
                case `ucs2`:
                case `ucs-2`:
                case `utf16le`:
                case `utf-16le`:
                    return `utf16le`;
                case `latin1`:
                case `binary`:
                    return `latin1`;
                case `base64`:
                case `ascii`:
                case `hex`:
                    return e;
                default:
                    if (t) return;
                    e = (`` + e).toLowerCase(), t = !0;
            }
        }
        function i(e) {
            var i = r(e);
            if (typeof i != `string` && (t.isEncoding === n || !n(e))) throw Error(`Unknown encoding: ` + e);
            return i || e;
        }
        e.StringDecoder = a;
        function a(e) {
            this.encoding = i(e);
            var n;
            switch(this.encoding){
                case `utf16le`:
                    this.text = f, this.end = p, n = 4;
                    break;
                case `utf8`:
                    this.fillLast = l, n = 4;
                    break;
                case `base64`:
                    this.text = m, this.end = h, n = 3;
                    break;
                default:
                    this.write = g, this.end = _;
                    return;
            }
            this.lastNeed = 0, this.lastTotal = 0, this.lastChar = t.allocUnsafe(n);
        }
        a.prototype.write = function(e) {
            if (e.length === 0) return ``;
            var t, n;
            if (this.lastNeed) {
                if (t = this.fillLast(e), t === void 0) return ``;
                n = this.lastNeed, this.lastNeed = 0;
            } else n = 0;
            return n < e.length ? t ? t + this.text(e, n) : this.text(e, n) : t || ``;
        }, a.prototype.end = d, a.prototype.text = u, a.prototype.fillLast = function(e) {
            if (this.lastNeed <= e.length) return e.copy(this.lastChar, this.lastTotal - this.lastNeed, 0, this.lastNeed), this.lastChar.toString(this.encoding, 0, this.lastTotal);
            e.copy(this.lastChar, this.lastTotal - this.lastNeed, 0, e.length), this.lastNeed -= e.length;
        };
        function o(e) {
            return e <= 127 ? 0 : e >> 5 == 6 ? 2 : e >> 4 == 14 ? 3 : e >> 3 == 30 ? 4 : e >> 6 == 2 ? -1 : -2;
        }
        function s(e, t, n) {
            var r = t.length - 1;
            if (r < n) return 0;
            var i = o(t[r]);
            return i >= 0 ? (i > 0 && (e.lastNeed = i - 1), i) : --r < n || i === -2 ? 0 : (i = o(t[r]), i >= 0 ? (i > 0 && (e.lastNeed = i - 2), i) : --r < n || i === -2 ? 0 : (i = o(t[r]), i >= 0 ? (i > 0 && (i === 2 ? i = 0 : e.lastNeed = i - 3), i) : 0));
        }
        function c(e, t, n) {
            if ((t[0] & 192) != 128) return e.lastNeed = 0, `�`;
            if (e.lastNeed > 1 && t.length > 1) {
                if ((t[1] & 192) != 128) return e.lastNeed = 1, `�`;
                if (e.lastNeed > 2 && t.length > 2 && (t[2] & 192) != 128) return e.lastNeed = 2, `�`;
            }
        }
        function l(e) {
            var t = this.lastTotal - this.lastNeed, n = c(this, e, t);
            if (n !== void 0) return n;
            if (this.lastNeed <= e.length) return e.copy(this.lastChar, t, 0, this.lastNeed), this.lastChar.toString(this.encoding, 0, this.lastTotal);
            e.copy(this.lastChar, t, 0, e.length), this.lastNeed -= e.length;
        }
        function u(e, t) {
            var n = s(this, e, t);
            if (!this.lastNeed) return e.toString(`utf8`, t);
            this.lastTotal = n;
            var r = e.length - (n - this.lastNeed);
            return e.copy(this.lastChar, 0, r), e.toString(`utf8`, t, r);
        }
        function d(e) {
            var t = e && e.length ? this.write(e) : ``;
            return this.lastNeed ? t + `�` : t;
        }
        function f(e, t) {
            if ((e.length - t) % 2 == 0) {
                var n = e.toString(`utf16le`, t);
                if (n) {
                    var r = n.charCodeAt(n.length - 1);
                    if (r >= 55296 && r <= 56319) return this.lastNeed = 2, this.lastTotal = 4, this.lastChar[0] = e[e.length - 2], this.lastChar[1] = e[e.length - 1], n.slice(0, -1);
                }
                return n;
            }
            return this.lastNeed = 1, this.lastTotal = 2, this.lastChar[0] = e[e.length - 1], e.toString(`utf16le`, t, e.length - 1);
        }
        function p(e) {
            var t = e && e.length ? this.write(e) : ``;
            if (this.lastNeed) {
                var n = this.lastTotal - this.lastNeed;
                return t + this.lastChar.toString(`utf16le`, 0, n);
            }
            return t;
        }
        function m(e, t) {
            var n = (e.length - t) % 3;
            return n === 0 ? e.toString(`base64`, t) : (this.lastNeed = 3 - n, this.lastTotal = 3, n === 1 ? this.lastChar[0] = e[e.length - 1] : (this.lastChar[0] = e[e.length - 2], this.lastChar[1] = e[e.length - 1]), e.toString(`base64`, t, e.length - n));
        }
        function h(e) {
            var t = e && e.length ? this.write(e) : ``;
            return this.lastNeed ? t + this.lastChar.toString(`base64`, 0, 3 - this.lastNeed) : t;
        }
        function g(e) {
            return e.toString(this.encoding);
        }
        function _(e) {
            return e && e.length ? this.write(e) : ``;
        }
    })), Oe = t(((e, t)=>{
        var n = xe().codes.ERR_STREAM_PREMATURE_CLOSE;
        function r(e) {
            var t = !1;
            return function() {
                if (!t) {
                    t = !0;
                    var n = [
                        ...arguments
                    ];
                    e.apply(this, n);
                }
            };
        }
        function i() {}
        function a(e) {
            return e.setHeader && typeof e.abort == `function`;
        }
        function o(e, t, s) {
            if (typeof t == `function`) return o(e, null, t);
            t ||= {}, s = r(s || i);
            var c = t.readable || t.readable !== !1 && e.readable, l = t.writable || t.writable !== !1 && e.writable, u = function() {
                e.writable || f();
            }, d = e._writableState && e._writableState.finished, f = function() {
                l = !1, d = !0, c || s.call(e);
            }, p = e._readableState && e._readableState.endEmitted, m = function() {
                c = !1, p = !0, l || s.call(e);
            }, h = function(t) {
                s.call(e, t);
            }, g = function() {
                var t;
                if (c && !p) return (!e._readableState || !e._readableState.ended) && (t = new n), s.call(e, t);
                if (l && !d) return (!e._writableState || !e._writableState.ended) && (t = new n), s.call(e, t);
            }, _ = function() {
                e.req.on(`finish`, f);
            };
            return a(e) ? (e.on(`complete`, f), e.on(`abort`, g), e.req ? _() : e.on(`request`, _)) : l && !e._writableState && (e.on(`end`, u), e.on(`close`, u)), e.on(`end`, m), e.on(`finish`, f), t.error !== !1 && e.on(`error`, h), e.on(`close`, g), function() {
                e.removeListener(`complete`, f), e.removeListener(`abort`, g), e.removeListener(`request`, _), e.req && e.req.removeListener(`finish`, f), e.removeListener(`end`, u), e.removeListener(`close`, u), e.removeListener(`finish`, f), e.removeListener(`end`, m), e.removeListener(`error`, h), e.removeListener(`close`, g);
            };
        }
        t.exports = o;
    })), ke = t(((e, t)=>{
        var n;
        function r(e, t, n) {
            return t = i(t), t in e ? Object.defineProperty(e, t, {
                value: n,
                enumerable: !0,
                configurable: !0,
                writable: !0
            }) : e[t] = n, e;
        }
        function i(e) {
            var t = a(e, `string`);
            return typeof t == `symbol` ? t : String(t);
        }
        function a(e, t) {
            if (typeof e != `object` || !e) return e;
            var n = e[Symbol.toPrimitive];
            if (n !== void 0) {
                var r = n.call(e, t || `default`);
                if (typeof r != `object`) return r;
                throw TypeError(`@@toPrimitive must return a primitive value.`);
            }
            return (t === `string` ? String : Number)(e);
        }
        var o = Oe(), s = Symbol(`lastResolve`), c = Symbol(`lastReject`), l = Symbol(`error`), u = Symbol(`ended`), d = Symbol(`lastPromise`), f = Symbol(`handlePromise`), p = Symbol(`stream`);
        function m(e, t) {
            return {
                value: e,
                done: t
            };
        }
        function h(e) {
            var t = e[s];
            if (t !== null) {
                var n = e[p].read();
                n !== null && (e[d] = null, e[s] = null, e[c] = null, t(m(n, !1)));
            }
        }
        function g(e) {
            process.nextTick(h, e);
        }
        function _(e, t) {
            return function(n, r) {
                e.then(function() {
                    if (t[u]) {
                        n(m(void 0, !0));
                        return;
                    }
                    t[f](n, r);
                }, r);
            };
        }
        var v = Object.getPrototypeOf(function() {}), y = Object.setPrototypeOf((n = {
            get stream () {
                return this[p];
            },
            next: function() {
                var e = this, t = this[l];
                if (t !== null) return Promise.reject(t);
                if (this[u]) return Promise.resolve(m(void 0, !0));
                if (this[p].destroyed) return new Promise(function(t, n) {
                    process.nextTick(function() {
                        e[l] ? n(e[l]) : t(m(void 0, !0));
                    });
                });
                var n = this[d], r;
                if (n) r = new Promise(_(n, this));
                else {
                    var i = this[p].read();
                    if (i !== null) return Promise.resolve(m(i, !1));
                    r = new Promise(this[f]);
                }
                return this[d] = r, r;
            }
        }, r(n, Symbol.asyncIterator, function() {
            return this;
        }), r(n, `return`, function() {
            var e = this;
            return new Promise(function(t, n) {
                e[p].destroy(null, function(e) {
                    if (e) {
                        n(e);
                        return;
                    }
                    t(m(void 0, !0));
                });
            });
        }), n), v);
        t.exports = function(e) {
            var t, n = Object.create(y, (t = {}, r(t, p, {
                value: e,
                writable: !0
            }), r(t, s, {
                value: null,
                writable: !0
            }), r(t, c, {
                value: null,
                writable: !0
            }), r(t, l, {
                value: null,
                writable: !0
            }), r(t, u, {
                value: e._readableState.endEmitted,
                writable: !0
            }), r(t, f, {
                value: function(e, t) {
                    var r = n[p].read();
                    r ? (n[d] = null, n[s] = null, n[c] = null, e(m(r, !1))) : (n[s] = e, n[c] = t);
                },
                writable: !0
            }), t));
            return n[d] = null, o(e, function(e) {
                if (e && e.code !== `ERR_STREAM_PREMATURE_CLOSE`) {
                    var t = n[c];
                    t !== null && (n[d] = null, n[s] = null, n[c] = null, t(e)), n[l] = e;
                    return;
                }
                var r = n[s];
                r !== null && (n[d] = null, n[s] = null, n[c] = null, r(m(void 0, !0))), n[u] = !0;
            }), e.on(`readable`, g.bind(null, n)), n;
        };
    })), Ae = t(((e, t)=>{
        t.exports = function() {
            throw Error(`Readable.from is not available in the browser`);
        };
    })), je = t(((e, t)=>{
        t.exports = A;
        var n;
        A.ReadableState = k, h().EventEmitter;
        var r = function(e, t) {
            return e.listeners(t).length;
        }, i = _(), a = m().Buffer, o = (typeof globalThis < `u` ? globalThis : typeof window < `u` ? window : typeof self < `u` ? self : {}).Uint8Array || function() {};
        function s(e) {
            return a.from(e);
        }
        function c(e) {
            return a.isBuffer(e) || e instanceof o;
        }
        var l = ve(), u = l && l.debuglog ? l.debuglog(`stream`) : function() {}, d = ye(), f = be(), p = Se().getHighWaterMark, v = xe().codes, y = v.ERR_INVALID_ARG_TYPE, b = v.ERR_STREAM_PUSH_AFTER_EOF, x = v.ERR_METHOD_NOT_IMPLEMENTED, S = v.ERR_STREAM_UNSHIFT_AFTER_END_EVENT, C, w, T;
        g()(A, i);
        var E = f.errorOrDestroy, D = [
            `error`,
            `close`,
            `destroy`,
            `pause`,
            `resume`
        ];
        function O(e, t, n) {
            if (typeof e.prependListener == `function`) return e.prependListener(t, n);
            !e._events || !e._events[t] ? e.on(t, n) : Array.isArray(e._events[t]) ? e._events[t].unshift(n) : e._events[t] = [
                n,
                e._events[t]
            ];
        }
        function k(e, t, r) {
            n ||= Te(), e ||= {}, typeof r != `boolean` && (r = t instanceof n), this.objectMode = !!e.objectMode, r && (this.objectMode = this.objectMode || !!e.readableObjectMode), this.highWaterMark = p(this, e, `readableHighWaterMark`, r), this.buffer = new d, this.length = 0, this.pipes = null, this.pipesCount = 0, this.flowing = null, this.ended = !1, this.endEmitted = !1, this.reading = !1, this.sync = !0, this.needReadable = !1, this.emittedReadable = !1, this.readableListening = !1, this.resumeScheduled = !1, this.paused = !0, this.emitClose = e.emitClose !== !1, this.autoDestroy = !!e.autoDestroy, this.destroyed = !1, this.defaultEncoding = e.defaultEncoding || `utf8`, this.awaitDrain = 0, this.readingMore = !1, this.decoder = null, this.encoding = null, e.encoding && (C ||= De().StringDecoder, this.decoder = new C(e.encoding), this.encoding = e.encoding);
        }
        function A(e) {
            if (n ||= Te(), !(this instanceof A)) return new A(e);
            var t = this instanceof n;
            this._readableState = new k(e, this, t), this.readable = !0, e && (typeof e.read == `function` && (this._read = e.read), typeof e.destroy == `function` && (this._destroy = e.destroy)), i.call(this);
        }
        Object.defineProperty(A.prototype, "destroyed", {
            enumerable: !1,
            get: function() {
                return this._readableState === void 0 ? !1 : this._readableState.destroyed;
            },
            set: function(e) {
                this._readableState && (this._readableState.destroyed = e);
            }
        }), A.prototype.destroy = f.destroy, A.prototype._undestroy = f.undestroy, A.prototype._destroy = function(e, t) {
            t(e);
        }, A.prototype.push = function(e, t) {
            var n = this._readableState, r;
            return n.objectMode ? r = !0 : typeof e == `string` && (t ||= n.defaultEncoding, t !== n.encoding && (e = a.from(e, t), t = ``), r = !0), j(this, e, t, !1, r);
        }, A.prototype.unshift = function(e) {
            return j(this, e, null, !0, !1);
        };
        function j(e, t, n, r, i) {
            u(`readableAddChunk`, t);
            var o = e._readableState;
            if (t === null) o.reading = !1, te(e, o);
            else {
                var c;
                if (i || (c = N(o, t)), c) E(e, c);
                else if (o.objectMode || t && t.length > 0) if (typeof t != `string` && !o.objectMode && Object.getPrototypeOf(t) !== a.prototype && (t = s(t)), r) o.endEmitted ? E(e, new S) : M(e, o, t, !0);
                else if (o.ended) E(e, new b);
                else if (o.destroyed) return !1;
                else o.reading = !1, o.decoder && !n ? (t = o.decoder.write(t), o.objectMode || t.length !== 0 ? M(e, o, t, !1) : re(e, o)) : M(e, o, t, !1);
                else r || (o.reading = !1, re(e, o));
            }
            return !o.ended && (o.length < o.highWaterMark || o.length === 0);
        }
        function M(e, t, n, r) {
            t.flowing && t.length === 0 && !t.sync ? (t.awaitDrain = 0, e.emit(`data`, n)) : (t.length += t.objectMode ? 1 : n.length, r ? t.buffer.unshift(n) : t.buffer.push(n), t.needReadable && I(e)), re(e, t);
        }
        function N(e, t) {
            var n;
            return !c(t) && typeof t != `string` && t !== void 0 && !e.objectMode && (n = new y(`chunk`, [
                `string`,
                `Buffer`,
                `Uint8Array`
            ], t)), n;
        }
        A.prototype.isPaused = function() {
            return this._readableState.flowing === !1;
        }, A.prototype.setEncoding = function(e) {
            C ||= De().StringDecoder;
            var t = new C(e);
            this._readableState.decoder = t, this._readableState.encoding = this._readableState.decoder.encoding;
            for(var n = this._readableState.buffer.head, r = ``; n !== null;)r += t.write(n.data), n = n.next;
            return this._readableState.buffer.clear(), r !== `` && this._readableState.buffer.push(r), this._readableState.length = r.length, this;
        };
        var P = 1073741824;
        function ee(e) {
            return e >= P ? e = P : (e--, e |= e >>> 1, e |= e >>> 2, e |= e >>> 4, e |= e >>> 8, e |= e >>> 16, e++), e;
        }
        function F(e, t) {
            return e <= 0 || t.length === 0 && t.ended ? 0 : t.objectMode ? 1 : e === e ? (e > t.highWaterMark && (t.highWaterMark = ee(e)), e <= t.length ? e : t.ended ? t.length : (t.needReadable = !0, 0)) : t.flowing && t.length ? t.buffer.head.data.length : t.length;
        }
        A.prototype.read = function(e) {
            u(`read`, e), e = parseInt(e, 10);
            var t = this._readableState, n = e;
            if (e !== 0 && (t.emittedReadable = !1), e === 0 && t.needReadable && ((t.highWaterMark === 0 ? t.length > 0 : t.length >= t.highWaterMark) || t.ended)) return u(`read: emitReadable`, t.length, t.ended), t.length === 0 && t.ended ? H(this) : I(this), null;
            if (e = F(e, t), e === 0 && t.ended) return t.length === 0 && H(this), null;
            var r = t.needReadable;
            u(`need readable`, r), (t.length === 0 || t.length - e < t.highWaterMark) && (r = !0, u(`length less than watermark`, r)), t.ended || t.reading ? (r = !1, u(`reading or ended`, r)) : r && (u(`do read`), t.reading = !0, t.sync = !0, t.length === 0 && (t.needReadable = !0), this._read(t.highWaterMark), t.sync = !1, t.reading || (e = F(n, t)));
            var i = e > 0 ? V(e, t) : null;
            return i === null ? (t.needReadable = t.length <= t.highWaterMark, e = 0) : (t.length -= e, t.awaitDrain = 0), t.length === 0 && (t.ended || (t.needReadable = !0), n !== e && t.ended && H(this)), i !== null && this.emit(`data`, i), i;
        };
        function te(e, t) {
            if (u(`onEofChunk`), !t.ended) {
                if (t.decoder) {
                    var n = t.decoder.end();
                    n && n.length && (t.buffer.push(n), t.length += t.objectMode ? 1 : n.length);
                }
                t.ended = !0, t.sync ? I(e) : (t.needReadable = !1, t.emittedReadable || (t.emittedReadable = !0, ne(e)));
            }
        }
        function I(e) {
            var t = e._readableState;
            u(`emitReadable`, t.needReadable, t.emittedReadable), t.needReadable = !1, t.emittedReadable || (u(`emitReadable`, t.flowing), t.emittedReadable = !0, process.nextTick(ne, e));
        }
        function ne(e) {
            var t = e._readableState;
            u(`emitReadable_`, t.destroyed, t.length, t.ended), !t.destroyed && (t.length || t.ended) && (e.emit(`readable`), t.emittedReadable = !1), t.needReadable = !t.flowing && !t.ended && t.length <= t.highWaterMark, oe(e);
        }
        function re(e, t) {
            t.readingMore || (t.readingMore = !0, process.nextTick(L, e, t));
        }
        function L(e, t) {
            for(; !t.reading && !t.ended && (t.length < t.highWaterMark || t.flowing && t.length === 0);){
                var n = t.length;
                if (u(`maybeReadMore read 0`), e.read(0), n === t.length) break;
            }
            t.readingMore = !1;
        }
        A.prototype._read = function(e) {
            E(this, new x(`_read()`));
        }, A.prototype.pipe = function(e, t) {
            var n = this, i = this._readableState;
            switch(i.pipesCount){
                case 0:
                    i.pipes = e;
                    break;
                case 1:
                    i.pipes = [
                        i.pipes,
                        e
                    ];
                    break;
                default:
                    i.pipes.push(e);
                    break;
            }
            i.pipesCount += 1, u(`pipe count=%d opts=%j`, i.pipesCount, t);
            var a = (!t || t.end !== !1) && e !== process.stdout && e !== process.stderr ? s : g;
            i.endEmitted ? process.nextTick(a) : n.once(`end`, a), e.on(`unpipe`, o);
            function o(e, t) {
                u(`onunpipe`), e === n && t && t.hasUnpiped === !1 && (t.hasUnpiped = !0, d());
            }
            function s() {
                u(`onend`), e.end();
            }
            var c = R(n);
            e.on(`drain`, c);
            var l = !1;
            function d() {
                u(`cleanup`), e.removeListener(`close`, m), e.removeListener(`finish`, h), e.removeListener(`drain`, c), e.removeListener(`error`, p), e.removeListener(`unpipe`, o), n.removeListener(`end`, s), n.removeListener(`end`, g), n.removeListener(`data`, f), l = !0, i.awaitDrain && (!e._writableState || e._writableState.needDrain) && c();
            }
            n.on(`data`, f);
            function f(t) {
                u(`ondata`);
                var r = e.write(t);
                u(`dest.write`, r), r === !1 && ((i.pipesCount === 1 && i.pipes === e || i.pipesCount > 1 && se(i.pipes, e) !== -1) && !l && (u(`false write response, pause`, i.awaitDrain), i.awaitDrain++), n.pause());
            }
            function p(t) {
                u(`onerror`, t), g(), e.removeListener(`error`, p), r(e, `error`) === 0 && E(e, t);
            }
            O(e, `error`, p);
            function m() {
                e.removeListener(`finish`, h), g();
            }
            e.once(`close`, m);
            function h() {
                u(`onfinish`), e.removeListener(`close`, m), g();
            }
            e.once(`finish`, h);
            function g() {
                u(`unpipe`), n.unpipe(e);
            }
            return e.emit(`pipe`, n), i.flowing || (u(`pipe resume`), n.resume()), e;
        };
        function R(e) {
            return function() {
                var t = e._readableState;
                u(`pipeOnDrain`, t.awaitDrain), t.awaitDrain && t.awaitDrain--, t.awaitDrain === 0 && r(e, `data`) && (t.flowing = !0, oe(e));
            };
        }
        A.prototype.unpipe = function(e) {
            var t = this._readableState, n = {
                hasUnpiped: !1
            };
            if (t.pipesCount === 0) return this;
            if (t.pipesCount === 1) return e && e !== t.pipes ? this : (e ||= t.pipes, t.pipes = null, t.pipesCount = 0, t.flowing = !1, e && e.emit(`unpipe`, this, n), this);
            if (!e) {
                var r = t.pipes, i = t.pipesCount;
                t.pipes = null, t.pipesCount = 0, t.flowing = !1;
                for(var a = 0; a < i; a++)r[a].emit(`unpipe`, this, {
                    hasUnpiped: !1
                });
                return this;
            }
            var o = se(t.pipes, e);
            return o === -1 ? this : (t.pipes.splice(o, 1), --t.pipesCount, t.pipesCount === 1 && (t.pipes = t.pipes[0]), e.emit(`unpipe`, this, n), this);
        }, A.prototype.on = function(e, t) {
            var n = i.prototype.on.call(this, e, t), r = this._readableState;
            return e === `data` ? (r.readableListening = this.listenerCount(`readable`) > 0, r.flowing !== !1 && this.resume()) : e === `readable` && !r.endEmitted && !r.readableListening && (r.readableListening = r.needReadable = !0, r.flowing = !1, r.emittedReadable = !1, u(`on readable`, r.length, r.reading), r.length ? I(this) : r.reading || process.nextTick(B, this)), n;
        }, A.prototype.addListener = A.prototype.on, A.prototype.removeListener = function(e, t) {
            var n = i.prototype.removeListener.call(this, e, t);
            return e === `readable` && process.nextTick(z, this), n;
        }, A.prototype.removeAllListeners = function(e) {
            var t = i.prototype.removeAllListeners.apply(this, arguments);
            return (e === `readable` || e === void 0) && process.nextTick(z, this), t;
        };
        function z(e) {
            var t = e._readableState;
            t.readableListening = e.listenerCount(`readable`) > 0, t.resumeScheduled && !t.paused ? t.flowing = !0 : e.listenerCount(`data`) > 0 && e.resume();
        }
        function B(e) {
            u(`readable nexttick read 0`), e.read(0);
        }
        A.prototype.resume = function() {
            var e = this._readableState;
            return e.flowing || (u(`resume`), e.flowing = !e.readableListening, ie(this, e)), e.paused = !1, this;
        };
        function ie(e, t) {
            t.resumeScheduled || (t.resumeScheduled = !0, process.nextTick(ae, e, t));
        }
        function ae(e, t) {
            u(`resume`, t.reading), t.reading || e.read(0), t.resumeScheduled = !1, e.emit(`resume`), oe(e), t.flowing && !t.reading && e.read(0);
        }
        A.prototype.pause = function() {
            return u(`call pause flowing=%j`, this._readableState.flowing), this._readableState.flowing !== !1 && (u(`pause`), this._readableState.flowing = !1, this.emit(`pause`)), this._readableState.paused = !0, this;
        };
        function oe(e) {
            var t = e._readableState;
            for(u(`flow`, t.flowing); t.flowing && e.read() !== null;);
        }
        A.prototype.wrap = function(e) {
            var t = this, n = this._readableState, r = !1;
            for(var i in e.on(`end`, function() {
                if (u(`wrapped end`), n.decoder && !n.ended) {
                    var e = n.decoder.end();
                    e && e.length && t.push(e);
                }
                t.push(null);
            }), e.on(`data`, function(i) {
                u(`wrapped data`), n.decoder && (i = n.decoder.write(i)), !(n.objectMode && i == null) && (!n.objectMode && (!i || !i.length) || t.push(i) || (r = !0, e.pause()));
            }), e)this[i] === void 0 && typeof e[i] == `function` && (this[i] = function(t) {
                return function() {
                    return e[t].apply(e, arguments);
                };
            }(i));
            for(var a = 0; a < D.length; a++)e.on(D[a], this.emit.bind(this, D[a]));
            return this._read = function(t) {
                u(`wrapped _read`, t), r && (r = !1, e.resume());
            }, this;
        }, typeof Symbol == `function` && (A.prototype[Symbol.asyncIterator] = function() {
            return w === void 0 && (w = ke()), w(this);
        }), Object.defineProperty(A.prototype, "readableHighWaterMark", {
            enumerable: !1,
            get: function() {
                return this._readableState.highWaterMark;
            }
        }), Object.defineProperty(A.prototype, "readableBuffer", {
            enumerable: !1,
            get: function() {
                return this._readableState && this._readableState.buffer;
            }
        }), Object.defineProperty(A.prototype, "readableFlowing", {
            enumerable: !1,
            get: function() {
                return this._readableState.flowing;
            },
            set: function(e) {
                this._readableState && (this._readableState.flowing = e);
            }
        }), A._fromList = V, Object.defineProperty(A.prototype, "readableLength", {
            enumerable: !1,
            get: function() {
                return this._readableState.length;
            }
        });
        function V(e, t) {
            if (t.length === 0) return null;
            var n;
            return t.objectMode ? n = t.buffer.shift() : !e || e >= t.length ? (n = t.decoder ? t.buffer.join(``) : t.buffer.length === 1 ? t.buffer.first() : t.buffer.concat(t.length), t.buffer.clear()) : n = t.buffer.consume(e, t.decoder), n;
        }
        function H(e) {
            var t = e._readableState;
            u(`endReadable`, t.endEmitted), t.endEmitted || (t.ended = !0, process.nextTick(U, t, e));
        }
        function U(e, t) {
            if (u(`endReadableNT`, e.endEmitted, e.length), !e.endEmitted && e.length === 0 && (e.endEmitted = !0, t.readable = !1, t.emit(`end`), e.autoDestroy)) {
                var n = t._writableState;
                (!n || n.autoDestroy && n.finished) && t.destroy();
            }
        }
        typeof Symbol == `function` && (A.from = function(e, t) {
            return T === void 0 && (T = Ae()), T(A, e, t);
        });
        function se(e, t) {
            for(var n = 0, r = e.length; n < r; n++)if (e[n] === t) return n;
            return -1;
        }
    })), Me = t(((e, t)=>{
        t.exports = l;
        var n = xe().codes, r = n.ERR_METHOD_NOT_IMPLEMENTED, i = n.ERR_MULTIPLE_CALLBACK, a = n.ERR_TRANSFORM_ALREADY_TRANSFORMING, o = n.ERR_TRANSFORM_WITH_LENGTH_0, s = Te();
        g()(l, s);
        function c(e, t) {
            var n = this._transformState;
            n.transforming = !1;
            var r = n.writecb;
            if (r === null) return this.emit(`error`, new i);
            n.writechunk = null, n.writecb = null, t != null && this.push(t), r(e);
            var a = this._readableState;
            a.reading = !1, (a.needReadable || a.length < a.highWaterMark) && this._read(a.highWaterMark);
        }
        function l(e) {
            if (!(this instanceof l)) return new l(e);
            s.call(this, e), this._transformState = {
                afterTransform: c.bind(this),
                needTransform: !1,
                transforming: !1,
                writecb: null,
                writechunk: null,
                writeencoding: null
            }, this._readableState.needReadable = !0, this._readableState.sync = !1, e && (typeof e.transform == `function` && (this._transform = e.transform), typeof e.flush == `function` && (this._flush = e.flush)), this.on(`prefinish`, u);
        }
        function u() {
            var e = this;
            typeof this._flush == `function` && !this._readableState.destroyed ? this._flush(function(t, n) {
                d(e, t, n);
            }) : d(this, null, null);
        }
        l.prototype.push = function(e, t) {
            return this._transformState.needTransform = !1, s.prototype.push.call(this, e, t);
        }, l.prototype._transform = function(e, t, n) {
            n(new r(`_transform()`));
        }, l.prototype._write = function(e, t, n) {
            var r = this._transformState;
            if (r.writecb = n, r.writechunk = e, r.writeencoding = t, !r.transforming) {
                var i = this._readableState;
                (r.needTransform || i.needReadable || i.length < i.highWaterMark) && this._read(i.highWaterMark);
            }
        }, l.prototype._read = function(e) {
            var t = this._transformState;
            t.writechunk !== null && !t.transforming ? (t.transforming = !0, this._transform(t.writechunk, t.writeencoding, t.afterTransform)) : t.needTransform = !0;
        }, l.prototype._destroy = function(e, t) {
            s.prototype._destroy.call(this, e, function(e) {
                t(e);
            });
        };
        function d(e, t, n) {
            if (t) return e.emit(`error`, t);
            if (n != null && e.push(n), e._writableState.length) throw new o;
            if (e._transformState.transforming) throw new a;
            return e.push(null);
        }
    })), Ne = t(((e, t)=>{
        t.exports = r;
        var n = Me();
        g()(r, n);
        function r(e) {
            if (!(this instanceof r)) return new r(e);
            n.call(this, e);
        }
        r.prototype._transform = function(e, t, n) {
            n(null, e);
        };
    })), Pe = t(((e, t)=>{
        var n;
        function r(e) {
            var t = !1;
            return function() {
                t || (t = !0, e.apply(void 0, arguments));
            };
        }
        var i = xe().codes, a = i.ERR_MISSING_ARGS, o = i.ERR_STREAM_DESTROYED;
        function s(e) {
            if (e) throw e;
        }
        function c(e) {
            return e.setHeader && typeof e.abort == `function`;
        }
        function l(e, t, i, a) {
            a = r(a);
            var s = !1;
            e.on(`close`, function() {
                s = !0;
            }), n === void 0 && (n = Oe()), n(e, {
                readable: t,
                writable: i
            }, function(e) {
                if (e) return a(e);
                s = !0, a();
            });
            var l = !1;
            return function(t) {
                if (!s && !l) {
                    if (l = !0, c(e)) return e.abort();
                    if (typeof e.destroy == `function`) return e.destroy();
                    a(t || new o(`pipe`));
                }
            };
        }
        function u(e) {
            e();
        }
        function d(e, t) {
            return e.pipe(t);
        }
        function f(e) {
            return !e.length || typeof e[e.length - 1] != `function` ? s : e.pop();
        }
        function p() {
            var e = [
                ...arguments
            ], t = f(e);
            if (Array.isArray(e[0]) && (e = e[0]), e.length < 2) throw new a(`streams`);
            var n, r = e.map(function(i, a) {
                var o = a < e.length - 1;
                return l(i, o, a > 0, function(e) {
                    n ||= e, e && r.forEach(u), !o && (r.forEach(u), t(n));
                });
            });
            return e.reduce(d);
        }
        t.exports = p;
    })), Fe = t(((e, t)=>{
        t.exports = r;
        var n = h().EventEmitter;
        g()(r, n), r.Readable = je(), r.Writable = we(), r.Duplex = Te(), r.Transform = Me(), r.PassThrough = Ne(), r.finished = Oe(), r.pipeline = Pe(), r.Stream = r;
        function r() {
            n.call(this);
        }
        r.prototype.pipe = function(e, t) {
            var r = this;
            function i(t) {
                e.writable && !1 === e.write(t) && r.pause && r.pause();
            }
            r.on(`data`, i);
            function a() {
                r.readable && r.resume && r.resume();
            }
            e.on(`drain`, a), !e._isStdio && (!t || t.end !== !1) && (r.on(`end`, s), r.on(`close`, c));
            var o = !1;
            function s() {
                o || (o = !0, e.end());
            }
            function c() {
                o || (o = !0, typeof e.destroy == `function` && e.destroy());
            }
            function l(e) {
                if (u(), n.listenerCount(this, `error`) === 0) throw e;
            }
            r.on(`error`, l), e.on(`error`, l);
            function u() {
                r.removeListener(`data`, i), e.removeListener(`drain`, a), r.removeListener(`end`, s), r.removeListener(`close`, c), r.removeListener(`error`, l), e.removeListener(`error`, l), r.removeListener(`end`, u), r.removeListener(`close`, u), e.removeListener(`close`, u);
            }
            return r.on(`end`, u), r.on(`close`, u), e.on(`close`, u), e.emit(`pipe`, r), e;
        };
    })), Ie = t(((e, t)=>{
        function n(e) {
            "@babel/helpers - typeof";
            return n = typeof Symbol == `function` && typeof Symbol.iterator == `symbol` ? function(e) {
                return typeof e;
            } : function(e) {
                return e && typeof Symbol == `function` && e.constructor === Symbol && e !== Symbol.prototype ? `symbol` : typeof e;
            }, n(e);
        }
        function r(e, t) {
            for(var n = 0; n < t.length; n++){
                var r = t[n];
                r.enumerable = r.enumerable || !1, r.configurable = !0, `value` in r && (r.writable = !0), Object.defineProperty(e, a(r.key), r);
            }
        }
        function i(e, t, n) {
            return t && r(e.prototype, t), n && r(e, n), Object.defineProperty(e, "prototype", {
                writable: !1
            }), e;
        }
        function a(e) {
            var t = o(e, `string`);
            return n(t) === `symbol` ? t : String(t);
        }
        function o(e, t) {
            if (n(e) !== `object` || e === null) return e;
            var r = e[Symbol.toPrimitive];
            if (r !== void 0) {
                var i = r.call(e, t || `default`);
                if (n(i) !== `object`) return i;
                throw TypeError(`@@toPrimitive must return a primitive value.`);
            }
            return (t === `string` ? String : Number)(e);
        }
        function s(e, t) {
            if (!(e instanceof t)) throw TypeError(`Cannot call a class as a function`);
        }
        function c(e, t) {
            if (typeof t != `function` && t !== null) throw TypeError(`Super expression must either be null or a function`);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    writable: !0,
                    configurable: !0
                }
            }), Object.defineProperty(e, "prototype", {
                writable: !1
            }), t && l(e, t);
        }
        function l(e, t) {
            return l = Object.setPrototypeOf ? Object.setPrototypeOf.bind() : function(e, t) {
                return e.__proto__ = t, e;
            }, l(e, t);
        }
        function u(e) {
            var t = p();
            return function() {
                var n = m(e), r;
                if (t) {
                    var i = m(this).constructor;
                    r = Reflect.construct(n, arguments, i);
                } else r = n.apply(this, arguments);
                return d(this, r);
            };
        }
        function d(e, t) {
            if (t && (n(t) === `object` || typeof t == `function`)) return t;
            if (t !== void 0) throw TypeError(`Derived constructors may only return object or undefined`);
            return f(e);
        }
        function f(e) {
            if (e === void 0) throw ReferenceError(`this hasn't been initialised - super() hasn't been called`);
            return e;
        }
        function p() {
            if (typeof Reflect > `u` || !Reflect.construct || Reflect.construct.sham) return !1;
            if (typeof Proxy == `function`) return !0;
            try {
                return Boolean.prototype.valueOf.call(Reflect.construct(Boolean, [], function() {})), !0;
            } catch  {
                return !1;
            }
        }
        function m(e) {
            return m = Object.setPrototypeOf ? Object.getPrototypeOf.bind() : function(e) {
                return e.__proto__ || Object.getPrototypeOf(e);
            }, m(e);
        }
        var h = {}, g, _;
        function v(e, t, n) {
            n ||= Error;
            function r(e, n, r) {
                return typeof t == `string` ? t : t(e, n, r);
            }
            h[e] = function(t) {
                c(a, t);
                var n = u(a);
                function a(t, i, o) {
                    var c;
                    return s(this, a), c = n.call(this, r(t, i, o)), c.code = e, c;
                }
                return i(a);
            }(n);
        }
        function y(e, t) {
            if (Array.isArray(e)) {
                var n = e.length;
                return e = e.map(function(e) {
                    return String(e);
                }), n > 2 ? `one of ${t} ${e.slice(0, n - 1).join(`, `)}, or ` + e[n - 1] : n === 2 ? `one of ${t} ${e[0]} or ${e[1]}` : `of ${t} ${e[0]}`;
            } else return `of ${t} ${String(e)}`;
        }
        function b(e, t, n) {
            return e.substr(!n || n < 0 ? 0 : +n, t.length) === t;
        }
        function x(e, t, n) {
            return (n === void 0 || n > e.length) && (n = e.length), e.substring(n - t.length, n) === t;
        }
        function S(e, t, n) {
            return typeof n != `number` && (n = 0), n + t.length > e.length ? !1 : e.indexOf(t, n) !== -1;
        }
        v(`ERR_AMBIGUOUS_ARGUMENT`, `The "%s" argument is ambiguous. %s`, TypeError), v(`ERR_INVALID_ARG_TYPE`, function(e, t, r) {
            g === void 0 && (g = et()), g(typeof e == `string`, `'name' must be a string`);
            var i;
            typeof t == `string` && b(t, `not `) ? (i = `must not be`, t = t.replace(/^not /, ``)) : i = `must be`;
            var a = x(e, ` argument`) ? `The ${e} ${i} ${y(t, `type`)}` : `The "${e}" ${S(e, `.`) ? `property` : `argument`} ${i} ${y(t, `type`)}`;
            return a += `. Received type ${n(r)}`, a;
        }, TypeError), v(`ERR_INVALID_ARG_VALUE`, function(e, t) {
            var n = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : `is invalid`;
            _ === void 0 && (_ = ve());
            var r = _.inspect(t);
            return r.length > 128 && (r = `${r.slice(0, 128)}...`), `The argument '${e}' ${n}. Received ${r}`;
        }, TypeError, RangeError), v(`ERR_INVALID_RETURN_VALUE`, function(e, t, r) {
            return `Expected ${e} to be returned from the "${t}" function but got ${r && r.constructor && r.constructor.name ? `instance of ${r.constructor.name}` : `type ${n(r)}`}.`;
        }, TypeError), v(`ERR_MISSING_ARGS`, function() {
            var e = [
                ...arguments
            ];
            g === void 0 && (g = et()), g(e.length > 0, `At least one arg needs to be specified`);
            var t = `The `, n = e.length;
            switch(e = e.map(function(e) {
                return `"${e}"`;
            }), n){
                case 1:
                    t += `${e[0]} argument`;
                    break;
                case 2:
                    t += `${e[0]} and ${e[1]} arguments`;
                    break;
                default:
                    t += e.slice(0, n - 1).join(`, `), t += `, and ${e[n - 1]} arguments`;
                    break;
            }
            return `${t} must be specified`;
        }, TypeError), t.exports.codes = h;
    })), Le = t(((e, t)=>{
        function n(e, t) {
            var n = Object.keys(e);
            if (Object.getOwnPropertySymbols) {
                var r = Object.getOwnPropertySymbols(e);
                t && (r = r.filter(function(t) {
                    return Object.getOwnPropertyDescriptor(e, t).enumerable;
                })), n.push.apply(n, r);
            }
            return n;
        }
        function r(e) {
            for(var t = 1; t < arguments.length; t++){
                var r = arguments[t] == null ? {} : arguments[t];
                t % 2 ? n(Object(r), !0).forEach(function(t) {
                    i(e, t, r[t]);
                }) : Object.getOwnPropertyDescriptors ? Object.defineProperties(e, Object.getOwnPropertyDescriptors(r)) : n(Object(r)).forEach(function(t) {
                    Object.defineProperty(e, t, Object.getOwnPropertyDescriptor(r, t));
                });
            }
            return e;
        }
        function i(e, t, n) {
            return t = c(t), t in e ? Object.defineProperty(e, t, {
                value: n,
                enumerable: !0,
                configurable: !0,
                writable: !0
            }) : e[t] = n, e;
        }
        function a(e, t) {
            if (!(e instanceof t)) throw TypeError(`Cannot call a class as a function`);
        }
        function o(e, t) {
            for(var n = 0; n < t.length; n++){
                var r = t[n];
                r.enumerable = r.enumerable || !1, r.configurable = !0, `value` in r && (r.writable = !0), Object.defineProperty(e, c(r.key), r);
            }
        }
        function s(e, t, n) {
            return t && o(e.prototype, t), n && o(e, n), Object.defineProperty(e, "prototype", {
                writable: !1
            }), e;
        }
        function c(e) {
            var t = l(e, `string`);
            return b(t) === `symbol` ? t : String(t);
        }
        function l(e, t) {
            if (b(e) !== `object` || e === null) return e;
            var n = e[Symbol.toPrimitive];
            if (n !== void 0) {
                var r = n.call(e, t || `default`);
                if (b(r) !== `object`) return r;
                throw TypeError(`@@toPrimitive must return a primitive value.`);
            }
            return (t === `string` ? String : Number)(e);
        }
        function u(e, t) {
            if (typeof t != `function` && t !== null) throw TypeError(`Super expression must either be null or a function`);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    writable: !0,
                    configurable: !0
                }
            }), Object.defineProperty(e, "prototype", {
                writable: !1
            }), t && v(e, t);
        }
        function d(e) {
            var t = g();
            return function() {
                var n = y(e), r;
                if (t) {
                    var i = y(this).constructor;
                    r = Reflect.construct(n, arguments, i);
                } else r = n.apply(this, arguments);
                return f(this, r);
            };
        }
        function f(e, t) {
            if (t && (b(t) === `object` || typeof t == `function`)) return t;
            if (t !== void 0) throw TypeError(`Derived constructors may only return object or undefined`);
            return p(e);
        }
        function p(e) {
            if (e === void 0) throw ReferenceError(`this hasn't been initialised - super() hasn't been called`);
            return e;
        }
        function m(e) {
            var t = typeof Map == `function` ? new Map : void 0;
            return m = function(e) {
                if (e === null || !_(e)) return e;
                if (typeof e != `function`) throw TypeError(`Super expression must either be null or a function`);
                if (t !== void 0) {
                    if (t.has(e)) return t.get(e);
                    t.set(e, n);
                }
                function n() {
                    return h(e, arguments, y(this).constructor);
                }
                return n.prototype = Object.create(e.prototype, {
                    constructor: {
                        value: n,
                        enumerable: !1,
                        writable: !0,
                        configurable: !0
                    }
                }), v(n, e);
            }, m(e);
        }
        function h(e, t, n) {
            return h = g() ? Reflect.construct.bind() : function(e, t, n) {
                var r = [
                    null
                ];
                r.push.apply(r, t);
                var i = new (Function.bind.apply(e, r));
                return n && v(i, n.prototype), i;
            }, h.apply(null, arguments);
        }
        function g() {
            if (typeof Reflect > `u` || !Reflect.construct || Reflect.construct.sham) return !1;
            if (typeof Proxy == `function`) return !0;
            try {
                return Boolean.prototype.valueOf.call(Reflect.construct(Boolean, [], function() {})), !0;
            } catch  {
                return !1;
            }
        }
        function _(e) {
            return Function.toString.call(e).indexOf(`[native code]`) !== -1;
        }
        function v(e, t) {
            return v = Object.setPrototypeOf ? Object.setPrototypeOf.bind() : function(e, t) {
                return e.__proto__ = t, e;
            }, v(e, t);
        }
        function y(e) {
            return y = Object.setPrototypeOf ? Object.getPrototypeOf.bind() : function(e) {
                return e.__proto__ || Object.getPrototypeOf(e);
            }, y(e);
        }
        function b(e) {
            "@babel/helpers - typeof";
            return b = typeof Symbol == `function` && typeof Symbol.iterator == `symbol` ? function(e) {
                return typeof e;
            } : function(e) {
                return e && typeof Symbol == `function` && e.constructor === Symbol && e !== Symbol.prototype ? `symbol` : typeof e;
            }, b(e);
        }
        var x = ve().inspect, S = Ie().codes.ERR_INVALID_ARG_TYPE;
        function C(e, t, n) {
            return (n === void 0 || n > e.length) && (n = e.length), e.substring(n - t.length, n) === t;
        }
        function w(e, t) {
            if (t = Math.floor(t), e.length == 0 || t == 0) return ``;
            var n = e.length * t;
            for(t = Math.floor(Math.log(t) / Math.log(2)); t;)e += e, t--;
            return e += e.substring(0, n - e.length), e;
        }
        var T = ``, E = ``, D = ``, O = ``, k = {
            deepStrictEqual: `Expected values to be strictly deep-equal:`,
            strictEqual: `Expected values to be strictly equal:`,
            strictEqualObject: `Expected "actual" to be reference-equal to "expected":`,
            deepEqual: `Expected values to be loosely deep-equal:`,
            equal: `Expected values to be loosely equal:`,
            notDeepStrictEqual: `Expected "actual" not to be strictly deep-equal to:`,
            notStrictEqual: `Expected "actual" to be strictly unequal to:`,
            notStrictEqualObject: `Expected "actual" not to be reference-equal to "expected":`,
            notDeepEqual: `Expected "actual" not to be loosely deep-equal to:`,
            notEqual: `Expected "actual" to be loosely unequal to:`,
            notIdentical: `Values identical but not reference-equal:`
        }, A = 10;
        function j(e) {
            var t = Object.keys(e), n = Object.create(Object.getPrototypeOf(e));
            return t.forEach(function(t) {
                n[t] = e[t];
            }), Object.defineProperty(n, "message", {
                value: e.message
            }), n;
        }
        function M(e) {
            return x(e, {
                compact: !1,
                customInspect: !1,
                depth: 1e3,
                maxArrayLength: 1 / 0,
                showHidden: !1,
                breakLength: 1 / 0,
                showProxy: !1,
                sorted: !0,
                getters: !0
            });
        }
        function N(e, t, n) {
            var r = ``, i = ``, a = 0, o = ``, s = !1, c = M(e), l = c.split(`
`), u = M(t).split(`
`), d = 0, f = ``;
            if (n === `strictEqual` && b(e) === `object` && b(t) === `object` && e !== null && t !== null && (n = `strictEqualObject`), l.length === 1 && u.length === 1 && l[0] !== u[0]) {
                var p = l[0].length + u[0].length;
                if (p <= A) {
                    if ((b(e) !== `object` || e === null) && (b(t) !== `object` || t === null) && (e !== 0 || t !== 0)) return `${k[n]}

${l[0]} !== ${u[0]}
`;
                } else if (n !== `strictEqualObject` && p < (process.stderr && process.stderr.isTTY ? process.stderr.columns : 80)) {
                    for(; l[0][d] === u[0][d];)d++;
                    d > 2 && (f = `
  ${w(` `, d)}^`, d = 0);
                }
            }
            for(var m = l[l.length - 1], h = u[u.length - 1]; m === h && (d++ < 2 ? o = `
  ${m}${o}` : r = m, l.pop(), u.pop(), !(l.length === 0 || u.length === 0));)m = l[l.length - 1], h = u[u.length - 1];
            var g = Math.max(l.length, u.length);
            if (g === 0) {
                var _ = c.split(`
`);
                if (_.length > 30) for(_[26] = `${T}...${O}`; _.length > 27;)_.pop();
                return `${k.notIdentical}

${_.join(`
`)}
`;
            }
            d > 3 && (o = `
${T}...${O}${o}`, s = !0), r !== `` && (o = `
  ${r}${o}`, r = ``);
            var v = 0, y = k[n] + `
${E}+ actual${O} ${D}- expected${O}`, x = ` ${T}...${O} Lines skipped`;
            for(d = 0; d < g; d++){
                var S = d - a;
                if (l.length < d + 1) S > 1 && d > 2 && (S > 4 ? (i += `
${T}...${O}`, s = !0) : S > 3 && (i += `
  ${u[d - 2]}`, v++), i += `
  ${u[d - 1]}`, v++), a = d, r += `
${D}-${O} ${u[d]}`, v++;
                else if (u.length < d + 1) S > 1 && d > 2 && (S > 4 ? (i += `
${T}...${O}`, s = !0) : S > 3 && (i += `
  ${l[d - 2]}`, v++), i += `
  ${l[d - 1]}`, v++), a = d, i += `
${E}+${O} ${l[d]}`, v++;
                else {
                    var j = u[d], N = l[d], P = N !== j && (!C(N, `,`) || N.slice(0, -1) !== j);
                    P && C(j, `,`) && j.slice(0, -1) === N && (P = !1, N += `,`), P ? (S > 1 && d > 2 && (S > 4 ? (i += `
${T}...${O}`, s = !0) : S > 3 && (i += `
  ${l[d - 2]}`, v++), i += `
  ${l[d - 1]}`, v++), a = d, i += `
${E}+${O} ${N}`, r += `
${D}-${O} ${j}`, v += 2) : (i += r, r = ``, (S === 1 || d === 0) && (i += `
  ${N}`, v++));
                }
                if (v > 20 && d < g - 2) return `${y}${x}
${i}
${T}...${O}${r}
${T}...${O}`;
            }
            return `${y}${s ? x : ``}
${i}${r}${o}${f}`;
        }
        t.exports = function(e, t) {
            u(i, e);
            var n = d(i);
            function i(e) {
                var t;
                if (a(this, i), b(e) !== `object` || e === null) throw new S(`options`, `Object`, e);
                var r = e.message, o = e.operator, s = e.stackStartFn, c = e.actual, l = e.expected, u = Error.stackTraceLimit;
                if (Error.stackTraceLimit = 0, r != null) t = n.call(this, String(r));
                else if (process.stderr && process.stderr.isTTY && (process.stderr && process.stderr.getColorDepth && process.stderr.getColorDepth() !== 1 ? (T = `\x1B[34m`, E = `\x1B[32m`, O = `\x1B[39m`, D = `\x1B[31m`) : (T = ``, E = ``, O = ``, D = ``)), b(c) === `object` && c !== null && b(l) === `object` && l !== null && `stack` in c && c instanceof Error && `stack` in l && l instanceof Error && (c = j(c), l = j(l)), o === `deepStrictEqual` || o === `strictEqual`) t = n.call(this, N(c, l, o));
                else if (o === `notDeepStrictEqual` || o === `notStrictEqual`) {
                    var d = k[o], m = M(c).split(`
`);
                    if (o === `notStrictEqual` && b(c) === `object` && c !== null && (d = k.notStrictEqualObject), m.length > 30) for(m[26] = `${T}...${O}`; m.length > 27;)m.pop();
                    t = m.length === 1 ? n.call(this, `${d} ${m[0]}`) : n.call(this, `${d}

${m.join(`
`)}
`);
                } else {
                    var h = M(c), g = ``, _ = k[o];
                    o === `notDeepEqual` || o === `notEqual` ? (h = `${k[o]}

${h}`, h.length > 1024 && (h = `${h.slice(0, 1021)}...`)) : (g = `${M(l)}`, h.length > 512 && (h = `${h.slice(0, 509)}...`), g.length > 512 && (g = `${g.slice(0, 509)}...`), o === `deepEqual` || o === `equal` ? h = `${_}

${h}

should equal

` : g = ` ${o} ${g}`), t = n.call(this, `${h}${g}`);
                }
                return Error.stackTraceLimit = u, t.generatedMessage = !r, Object.defineProperty(p(t), "name", {
                    value: `AssertionError [ERR_ASSERTION]`,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }), t.code = `ERR_ASSERTION`, t.actual = c, t.expected = l, t.operator = o, Error.captureStackTrace && Error.captureStackTrace(p(t), s), t.stack, t.name = `AssertionError`, f(t);
            }
            return s(i, [
                {
                    key: `toString`,
                    value: function() {
                        return `${this.name} [${this.code}]: ${this.message}`;
                    }
                },
                {
                    key: t,
                    value: function(e, t) {
                        return x(this, r(r({}, t), {}, {
                            customInspect: !1,
                            depth: 0
                        }));
                    }
                }
            ]), i;
        }(m(Error), x.custom);
    })), Re = t(((e, t)=>{
        var n = Object.prototype.toString;
        t.exports = function(e) {
            var t = n.call(e), r = t === `[object Arguments]`;
            return r ||= t !== `[object Array]` && typeof e == `object` && !!e && typeof e.length == `number` && e.length >= 0 && n.call(e.callee) === `[object Function]`, r;
        };
    })), ze = t(((e, t)=>{
        var n;
        if (!Object.keys) {
            var r = Object.prototype.hasOwnProperty, i = Object.prototype.toString, a = Re(), o = Object.prototype.propertyIsEnumerable, s = !o.call({
                toString: null
            }, `toString`), c = o.call(function() {}, `prototype`), l = [
                `toString`,
                `toLocaleString`,
                `valueOf`,
                `hasOwnProperty`,
                `isPrototypeOf`,
                `propertyIsEnumerable`,
                `constructor`
            ], u = function(e) {
                var t = e.constructor;
                return t && t.prototype === e;
            }, d = {
                $applicationCache: !0,
                $console: !0,
                $external: !0,
                $frame: !0,
                $frameElement: !0,
                $frames: !0,
                $innerHeight: !0,
                $innerWidth: !0,
                $onmozfullscreenchange: !0,
                $onmozfullscreenerror: !0,
                $outerHeight: !0,
                $outerWidth: !0,
                $pageXOffset: !0,
                $pageYOffset: !0,
                $parent: !0,
                $scrollLeft: !0,
                $scrollTop: !0,
                $scrollX: !0,
                $scrollY: !0,
                $self: !0,
                $webkitIndexedDB: !0,
                $webkitStorageInfo: !0,
                $window: !0
            }, f = function() {
                if (typeof window > `u`) return !1;
                for(var e in window)try {
                    if (!d[`$` + e] && r.call(window, e) && window[e] !== null && typeof window[e] == `object`) try {
                        u(window[e]);
                    } catch  {
                        return !0;
                    }
                } catch  {
                    return !0;
                }
                return !1;
            }(), p = function(e) {
                if (typeof window > `u` || !f) return u(e);
                try {
                    return u(e);
                } catch  {
                    return !1;
                }
            };
            n = function(e) {
                var t = typeof e == `object` && !!e, n = i.call(e) === `[object Function]`, o = a(e), u = t && i.call(e) === `[object String]`, d = [];
                if (!t && !n && !o) throw TypeError(`Object.keys called on a non-object`);
                var f = c && n;
                if (u && e.length > 0 && !r.call(e, 0)) for(var m = 0; m < e.length; ++m)d.push(String(m));
                if (o && e.length > 0) for(var h = 0; h < e.length; ++h)d.push(String(h));
                else for(var g in e)!(f && g === `prototype`) && r.call(e, g) && d.push(String(g));
                if (s) for(var _ = p(e), v = 0; v < l.length; ++v)!(_ && l[v] === `constructor`) && r.call(e, l[v]) && d.push(l[v]);
                return d;
            };
        }
        t.exports = n;
    })), Be = t(((e, t)=>{
        var n = Array.prototype.slice, r = Re(), i = Object.keys, a = i ? function(e) {
            return i(e);
        } : ze(), o = Object.keys;
        a.shim = function() {
            return Object.keys ? function() {
                var e = Object.keys(arguments);
                return e && e.length === arguments.length;
            }(1, 2) || (Object.keys = function(e) {
                return r(e) ? o(n.call(e)) : o(e);
            }) : Object.keys = a, Object.keys || a;
        }, t.exports = a;
    })), Ve = t(((e, t)=>{
        var n = Be(), r = v()(), i = G(), a = b(), o = i(`Array.prototype.push`), s = i(`Object.prototype.propertyIsEnumerable`), c = r ? a.getOwnPropertySymbols : null;
        t.exports = function(e, t) {
            if (e == null) throw TypeError(`target must be an object`);
            var i = a(e);
            if (arguments.length === 1) return i;
            for(var l = 1; l < arguments.length; ++l){
                var u = a(arguments[l]), d = n(u), f = r && (a.getOwnPropertySymbols || c);
                if (f) for(var p = f(u), m = 0; m < p.length; ++m){
                    var h = p[m];
                    s(u, h) && o(d, h);
                }
                for(var g = 0; g < d.length; ++g){
                    var _ = d[g];
                    s(u, _) && (i[_] = u[_]);
                }
            }
            return i;
        };
    })), He = t(((e, t)=>{
        var n = Ve(), r = function() {
            if (!Object.assign) return !1;
            for(var e = `abcdefghijklmnopqrst`, t = e.split(``), n = {}, r = 0; r < t.length; ++r)n[t[r]] = t[r];
            var i = Object.assign({}, n), a = ``;
            for(var o in i)a += o;
            return e !== a;
        }, i = function() {
            if (!Object.assign || !Object.preventExtensions) return !1;
            var e = Object.preventExtensions({
                1: 2
            });
            try {
                Object.assign(e, `xy`);
            } catch  {
                return e[1] === `y`;
            }
            return !1;
        };
        t.exports = function() {
            return !Object.assign || r() || i() ? n : Object.assign;
        };
    })), Ue = t(((e, t)=>{
        var n = function(e) {
            return e !== e;
        };
        t.exports = function(e, t) {
            return e === 0 && t === 0 ? 1 / e == 1 / t : !!(e === t || n(e) && n(t));
        };
    })), We = t(((e, t)=>{
        var n = Ue();
        t.exports = function() {
            return typeof Object.is == `function` ? Object.is : n;
        };
    })), Ge = t(((e, t)=>{
        var n = W(), r = $(), i = r(n(`String.prototype.indexOf`));
        t.exports = function(e, t) {
            var a = n(e, !!t);
            return typeof a == `function` && i(e, `.prototype.`) > -1 ? r(a) : a;
        };
    })), Ke = t(((e, t)=>{
        var n = Be(), r = typeof Symbol == `function` && typeof Symbol(`foo`) == `symbol`, i = Object.prototype.toString, a = Array.prototype.concat, o = de(), s = function(e) {
            return typeof e == `function` && i.call(e) === `[object Function]`;
        }, c = Q()(), l = function(e, t, n, r) {
            if (t in e) {
                if (r === !0) {
                    if (e[t] === n) return;
                } else if (!s(r) || !r()) return;
            }
            c ? o(e, t, n, !0) : o(e, t, n);
        }, u = function(e, t) {
            var i = arguments.length > 2 ? arguments[2] : {}, o = n(t);
            r && (o = a.call(o, Object.getOwnPropertySymbols(t)));
            for(var s = 0; s < o.length; s += 1)l(e, o[s], t[o[s]], i[o[s]]);
        };
        u.supportsDescriptors = !!c, t.exports = u;
    })), qe = t(((e, t)=>{
        var n = We(), r = Ke();
        t.exports = function() {
            var e = n();
            return r(Object, {
                is: e
            }, {
                is: function() {
                    return Object.is !== e;
                }
            }), e;
        };
    })), Je = t(((e, t)=>{
        var n = Ke(), r = $(), i = Ue(), a = We(), o = qe(), s = r(a(), Object);
        n(s, {
            getPolyfill: a,
            implementation: i,
            shim: o
        }), t.exports = s;
    })), Ye = t(((e, t)=>{
        t.exports = function(e) {
            return e !== e;
        };
    })), Xe = t(((e, t)=>{
        var n = Ye();
        t.exports = function() {
            return Number.isNaN && !Number.isNaN(`a`) ? Number.isNaN : n;
        };
    })), Ze = t(((e, t)=>{
        var n = Ke(), r = Xe();
        t.exports = function() {
            var e = r();
            return n(Number, {
                isNaN: e
            }, {
                isNaN: function() {
                    return Number.isNaN !== e;
                }
            }), e;
        };
    })), Qe = t(((e, t)=>{
        var n = $(), r = Ke(), i = Ye(), a = Xe(), o = Ze(), s = n(a(), Number);
        r(s, {
            getPolyfill: a,
            implementation: i,
            shim: o
        }), t.exports = s;
    })), $e = t(((e, t)=>{
        function n(e, t) {
            return s(e) || o(e, t) || i(e, t) || r();
        }
        function r() {
            throw TypeError(`Invalid attempt to destructure non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`);
        }
        function i(e, t) {
            if (e) {
                if (typeof e == `string`) return a(e, t);
                var n = Object.prototype.toString.call(e).slice(8, -1);
                if (n === `Object` && e.constructor && (n = e.constructor.name), n === `Map` || n === `Set`) return Array.from(e);
                if (n === `Arguments` || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return a(e, t);
            }
        }
        function a(e, t) {
            (t == null || t > e.length) && (t = e.length);
            for(var n = 0, r = Array(t); n < t; n++)r[n] = e[n];
            return r;
        }
        function o(e, t) {
            var n = e == null ? null : typeof Symbol < `u` && e[Symbol.iterator] || e[`@@iterator`];
            if (n != null) {
                var r, i, a, o, s = [], c = !0, l = !1;
                try {
                    if (a = (n = n.call(e)).next, t === 0) {
                        if (Object(n) !== n) return;
                        c = !1;
                    } else for(; !(c = (r = a.call(n)).done) && (s.push(r.value), s.length !== t); c = !0);
                } catch (e) {
                    l = !0, i = e;
                } finally{
                    try {
                        if (!c && n.return != null && (o = n.return(), Object(o) !== o)) return;
                    } finally{
                        if (l) throw i;
                    }
                }
                return s;
            }
        }
        function s(e) {
            if (Array.isArray(e)) return e;
        }
        function c(e) {
            "@babel/helpers - typeof";
            return c = typeof Symbol == `function` && typeof Symbol.iterator == `symbol` ? function(e) {
                return typeof e;
            } : function(e) {
                return e && typeof Symbol == `function` && e.constructor === Symbol && e !== Symbol.prototype ? `symbol` : typeof e;
            }, c(e);
        }
        var l = /a/g.flags !== void 0, u = function(e) {
            var t = [];
            return e.forEach(function(e) {
                return t.push(e);
            }), t;
        }, d = function(e) {
            var t = [];
            return e.forEach(function(e, n) {
                return t.push([
                    n,
                    e
                ]);
            }), t;
        }, f = Object.is ? Object.is : Je(), p = Object.getOwnPropertySymbols ? Object.getOwnPropertySymbols : function() {
            return [];
        }, m = Number.isNaN ? Number.isNaN : Qe();
        function h(e) {
            return e.call.bind(e);
        }
        var g = h(Object.prototype.hasOwnProperty), _ = h(Object.prototype.propertyIsEnumerable), v = h(Object.prototype.toString), y = ve().types, b = y.isAnyArrayBuffer, x = y.isArrayBufferView, S = y.isDate, C = y.isMap, w = y.isRegExp, T = y.isSet, E = y.isNativeError, D = y.isBoxedPrimitive, O = y.isNumberObject, k = y.isStringObject, A = y.isBooleanObject, j = y.isBigIntObject, M = y.isSymbolObject, N = y.isFloat32Array, P = y.isFloat64Array;
        function ee(e) {
            if (e.length === 0 || e.length > 10) return !0;
            for(var t = 0; t < e.length; t++){
                var n = e.charCodeAt(t);
                if (n < 48 || n > 57) return !0;
            }
            return e.length === 10 && e >= 2 ** 32;
        }
        function F(e) {
            return Object.keys(e).filter(ee).concat(p(e).filter(Object.prototype.propertyIsEnumerable.bind(e)));
        }
        function te(e, t) {
            if (e === t) return 0;
            for(var n = e.length, r = t.length, i = 0, a = Math.min(n, r); i < a; ++i)if (e[i] !== t[i]) {
                n = e[i], r = t[i];
                break;
            }
            return n < r ? -1 : +(r < n);
        }
        var I = void 0, ne = !0, re = !1, L = 0, R = 1, z = 2, B = 3;
        function ie(e, t) {
            return l ? e.source === t.source && e.flags === t.flags : RegExp.prototype.toString.call(e) === RegExp.prototype.toString.call(t);
        }
        function ae(e, t) {
            if (e.byteLength !== t.byteLength) return !1;
            for(var n = 0; n < e.byteLength; n++)if (e[n] !== t[n]) return !1;
            return !0;
        }
        function oe(e, t) {
            return e.byteLength === t.byteLength ? te(new Uint8Array(e.buffer, e.byteOffset, e.byteLength), new Uint8Array(t.buffer, t.byteOffset, t.byteLength)) === 0 : !1;
        }
        function V(e, t) {
            return e.byteLength === t.byteLength && te(new Uint8Array(e), new Uint8Array(t)) === 0;
        }
        function H(e, t) {
            return O(e) ? O(t) && f(Number.prototype.valueOf.call(e), Number.prototype.valueOf.call(t)) : k(e) ? k(t) && String.prototype.valueOf.call(e) === String.prototype.valueOf.call(t) : A(e) ? A(t) && Boolean.prototype.valueOf.call(e) === Boolean.prototype.valueOf.call(t) : j(e) ? j(t) && BigInt.prototype.valueOf.call(e) === BigInt.prototype.valueOf.call(t) : M(t) && Symbol.prototype.valueOf.call(e) === Symbol.prototype.valueOf.call(t);
        }
        function U(e, t, n, r) {
            if (e === t) return e === 0 && n ? f(e, t) : !0;
            if (n) {
                if (c(e) !== `object`) return typeof e == `number` && m(e) && m(t);
                if (c(t) !== `object` || e === null || t === null || Object.getPrototypeOf(e) !== Object.getPrototypeOf(t)) return !1;
            } else {
                if (e === null || c(e) !== `object`) return t === null || c(t) !== `object` ? e == t : !1;
                if (t === null || c(t) !== `object`) return !1;
            }
            var i = v(e);
            if (i !== v(t)) return !1;
            if (Array.isArray(e)) {
                if (e.length !== t.length) return !1;
                var a = F(e, I), o = F(t, I);
                return a.length === o.length ? W(e, t, n, r, R, a) : !1;
            }
            if (i === `[object Object]` && (!C(e) && C(t) || !T(e) && T(t))) return !1;
            if (S(e)) {
                if (!S(t) || Date.prototype.getTime.call(e) !== Date.prototype.getTime.call(t)) return !1;
            } else if (w(e)) {
                if (!w(t) || !ie(e, t)) return !1;
            } else if (E(e) || e instanceof Error) {
                if (e.message !== t.message || e.name !== t.name) return !1;
            } else if (x(e)) {
                if (!n && (N(e) || P(e))) {
                    if (!ae(e, t)) return !1;
                } else if (!oe(e, t)) return !1;
                var s = F(e, I), l = F(t, I);
                return s.length === l.length ? W(e, t, n, r, L, s) : !1;
            } else if (T(e)) return !T(t) || e.size !== t.size ? !1 : W(e, t, n, r, z);
            else if (C(e)) return !C(t) || e.size !== t.size ? !1 : W(e, t, n, r, B);
            else if (b(e)) {
                if (!V(e, t)) return !1;
            } else if (D(e) && !H(e, t)) return !1;
            return W(e, t, n, r, L);
        }
        function se(e, t) {
            return t.filter(function(t) {
                return _(e, t);
            });
        }
        function W(e, t, n, r, i, a) {
            if (arguments.length === 5) {
                a = Object.keys(e);
                var o = Object.keys(t);
                if (a.length !== o.length) return !1;
            }
            for(var s = 0; s < a.length; s++)if (!g(t, a[s])) return !1;
            if (n && arguments.length === 5) {
                var c = p(e);
                if (c.length !== 0) {
                    var l = 0;
                    for(s = 0; s < c.length; s++){
                        var u = c[s];
                        if (_(e, u)) {
                            if (!_(t, u)) return !1;
                            a.push(u), l++;
                        } else if (_(t, u)) return !1;
                    }
                    var d = p(t);
                    if (c.length !== d.length && se(t, d).length !== l) return !1;
                } else {
                    var f = p(t);
                    if (f.length !== 0 && se(t, f).length !== 0) return !1;
                }
            }
            if (a.length === 0 && (i === L || i === R && e.length === 0 || e.size === 0)) return !0;
            if (r === void 0) r = {
                val1: new Map,
                val2: new Map,
                position: 0
            };
            else {
                var m = r.val1.get(e);
                if (m !== void 0) {
                    var h = r.val2.get(t);
                    if (h !== void 0) return m === h;
                }
                r.position++;
            }
            r.val1.set(e, r.position), r.val2.set(t, r.position);
            var v = Y(e, t, n, a, r, i);
            return r.val1.delete(e), r.val2.delete(t), v;
        }
        function G(e, t, n, r) {
            for(var i = u(e), a = 0; a < i.length; a++){
                var o = i[a];
                if (U(t, o, n, r)) return e.delete(o), !0;
            }
            return !1;
        }
        function K(e) {
            switch(c(e)){
                case `undefined`:
                    return null;
                case `object`:
                    return;
                case `symbol`:
                    return !1;
                case `string`:
                    e = +e;
                case `number`:
                    if (m(e)) return !1;
            }
            return !0;
        }
        function ce(e, t, n) {
            var r = K(n);
            return r ?? (t.has(r) && !e.has(r));
        }
        function le(e, t, n, r, i) {
            var a = K(n);
            if (a != null) return a;
            var o = t.get(a);
            return o === void 0 && !t.has(a) || !U(r, o, !1, i) ? !1 : !e.has(a) && U(r, o, !1, i);
        }
        function q(e, t, n, r) {
            for(var i = null, a = u(e), o = 0; o < a.length; o++){
                var s = a[o];
                if (c(s) === `object` && s !== null) i === null && (i = new Set), i.add(s);
                else if (!t.has(s)) {
                    if (n || !ce(e, t, s)) return !1;
                    i === null && (i = new Set), i.add(s);
                }
            }
            if (i !== null) {
                for(var l = u(t), d = 0; d < l.length; d++){
                    var f = l[d];
                    if (c(f) === `object` && f !== null) {
                        if (!G(i, f, n, r)) return !1;
                    } else if (!n && !e.has(f) && !G(i, f, n, r)) return !1;
                }
                return i.size === 0;
            }
            return !0;
        }
        function ue(e, t, n, r, i, a) {
            for(var o = u(e), s = 0; s < o.length; s++){
                var c = o[s];
                if (U(n, c, i, a) && U(r, t.get(c), i, a)) return e.delete(c), !0;
            }
            return !1;
        }
        function J(e, t, r, i) {
            for(var a = null, o = d(e), s = 0; s < o.length; s++){
                var l = n(o[s], 2), u = l[0], f = l[1];
                if (c(u) === `object` && u !== null) a === null && (a = new Set), a.add(u);
                else {
                    var p = t.get(u);
                    if (p === void 0 && !t.has(u) || !U(f, p, r, i)) {
                        if (r || !le(e, t, u, f, i)) return !1;
                        a === null && (a = new Set), a.add(u);
                    }
                }
            }
            if (a !== null) {
                for(var m = d(t), h = 0; h < m.length; h++){
                    var g = n(m[h], 2), _ = g[0], v = g[1];
                    if (c(_) === `object` && _ !== null) {
                        if (!ue(a, e, _, v, r, i)) return !1;
                    } else if (!r && (!e.has(_) || !U(e.get(_), v, !1, i)) && !ue(a, e, _, v, !1, i)) return !1;
                }
                return a.size === 0;
            }
            return !0;
        }
        function Y(e, t, n, r, i, a) {
            var o = 0;
            if (a === z) {
                if (!q(e, t, n, i)) return !1;
            } else if (a === B) {
                if (!J(e, t, n, i)) return !1;
            } else if (a === R) for(; o < e.length; o++)if (g(e, o)) {
                if (!g(t, o) || !U(e[o], t[o], n, i)) return !1;
            } else if (g(t, o)) return !1;
            else {
                for(var s = Object.keys(e); o < s.length; o++){
                    var c = s[o];
                    if (!g(t, c) || !U(e[c], t[c], n, i)) return !1;
                }
                return s.length === Object.keys(t).length;
            }
            for(o = 0; o < r.length; o++){
                var l = r[o];
                if (!U(e[l], t[l], n, i)) return !1;
            }
            return !0;
        }
        function X(e, t) {
            return U(e, t, re);
        }
        function Z(e, t) {
            return U(e, t, ne);
        }
        t.exports = {
            isDeepEqual: X,
            isDeepStrictEqual: Z
        };
    })), et = t(((e, t)=>{
        function n(e) {
            "@babel/helpers - typeof";
            return n = typeof Symbol == `function` && typeof Symbol.iterator == `symbol` ? function(e) {
                return typeof e;
            } : function(e) {
                return e && typeof Symbol == `function` && e.constructor === Symbol && e !== Symbol.prototype ? `symbol` : typeof e;
            }, n(e);
        }
        function r(e, t) {
            for(var n = 0; n < t.length; n++){
                var r = t[n];
                r.enumerable = r.enumerable || !1, r.configurable = !0, `value` in r && (r.writable = !0), Object.defineProperty(e, a(r.key), r);
            }
        }
        function i(e, t, n) {
            return t && r(e.prototype, t), n && r(e, n), Object.defineProperty(e, "prototype", {
                writable: !1
            }), e;
        }
        function a(e) {
            var t = o(e, `string`);
            return n(t) === `symbol` ? t : String(t);
        }
        function o(e, t) {
            if (n(e) !== `object` || e === null) return e;
            var r = e[Symbol.toPrimitive];
            if (r !== void 0) {
                var i = r.call(e, t || `default`);
                if (n(i) !== `object`) return i;
                throw TypeError(`@@toPrimitive must return a primitive value.`);
            }
            return (t === `string` ? String : Number)(e);
        }
        function s(e, t) {
            if (!(e instanceof t)) throw TypeError(`Cannot call a class as a function`);
        }
        var c = Ie().codes, l = c.ERR_AMBIGUOUS_ARGUMENT, u = c.ERR_INVALID_ARG_TYPE, d = c.ERR_INVALID_ARG_VALUE, f = c.ERR_INVALID_RETURN_VALUE, p = c.ERR_MISSING_ARGS, m = Le(), h = ve().inspect, g = ve().types, _ = g.isPromise, v = g.isRegExp, y = He()(), b = We()(), x = Ge()(`RegExp.prototype.test`), S, C;
        function w() {
            var e = $e();
            S = e.isDeepEqual, C = e.isDeepStrictEqual;
        }
        var T = !1, E = t.exports = j, D = {};
        function O(e) {
            throw e.message instanceof Error ? e.message : new m(e);
        }
        function k(e, t, n, r, i) {
            var a = arguments.length, o;
            if (a === 0 ? o = `Failed` : a === 1 ? (n = e, e = void 0) : (T === !1 && (T = !0, (process.emitWarning ? process.emitWarning : console.warn.bind(console))(`assert.fail() with more than one argument is deprecated. Please use assert.strictEqual() instead or only pass a message.`, `DeprecationWarning`, `DEP0094`)), a === 2 && (r = `!=`)), n instanceof Error) throw n;
            var s = {
                actual: e,
                expected: t,
                operator: r === void 0 ? `fail` : r,
                stackStartFn: i || k
            };
            n !== void 0 && (s.message = n);
            var c = new m(s);
            throw o && (c.message = o, c.generatedMessage = !0), c;
        }
        E.fail = k, E.AssertionError = m;
        function A(e, t, n, r) {
            if (!n) {
                var i = !1;
                if (t === 0) i = !0, r = "No value argument passed to `assert.ok()`";
                else if (r instanceof Error) throw r;
                var a = new m({
                    actual: n,
                    expected: !0,
                    message: r,
                    operator: `==`,
                    stackStartFn: e
                });
                throw a.generatedMessage = i, a;
            }
        }
        function j() {
            var e = [
                ...arguments
            ];
            A.apply(void 0, [
                j,
                e.length
            ].concat(e));
        }
        E.ok = j, E.equal = function e(t, n, r) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            t != n && O({
                actual: t,
                expected: n,
                message: r,
                operator: `==`,
                stackStartFn: e
            });
        }, E.notEqual = function e(t, n, r) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            t == n && O({
                actual: t,
                expected: n,
                message: r,
                operator: `!=`,
                stackStartFn: e
            });
        }, E.deepEqual = function e(t, n, r) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            S === void 0 && w(), S(t, n) || O({
                actual: t,
                expected: n,
                message: r,
                operator: `deepEqual`,
                stackStartFn: e
            });
        }, E.notDeepEqual = function e(t, n, r) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            S === void 0 && w(), S(t, n) && O({
                actual: t,
                expected: n,
                message: r,
                operator: `notDeepEqual`,
                stackStartFn: e
            });
        }, E.deepStrictEqual = function e(t, n, r) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            S === void 0 && w(), C(t, n) || O({
                actual: t,
                expected: n,
                message: r,
                operator: `deepStrictEqual`,
                stackStartFn: e
            });
        }, E.notDeepStrictEqual = M;
        function M(e, t, n) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            S === void 0 && w(), C(e, t) && O({
                actual: e,
                expected: t,
                message: n,
                operator: `notDeepStrictEqual`,
                stackStartFn: M
            });
        }
        E.strictEqual = function e(t, n, r) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            b(t, n) || O({
                actual: t,
                expected: n,
                message: r,
                operator: `strictEqual`,
                stackStartFn: e
            });
        }, E.notStrictEqual = function e(t, n, r) {
            if (arguments.length < 2) throw new p(`actual`, `expected`);
            b(t, n) && O({
                actual: t,
                expected: n,
                message: r,
                operator: `notStrictEqual`,
                stackStartFn: e
            });
        };
        var N = i(function e(t, n, r) {
            var i = this;
            s(this, e), n.forEach(function(e) {
                e in t && (r !== void 0 && typeof r[e] == `string` && v(t[e]) && x(t[e], r[e]) ? i[e] = r[e] : i[e] = t[e]);
            });
        });
        function P(e, t, n, r, i, a) {
            if (!(n in e) || !C(e[n], t[n])) {
                if (!r) {
                    var o = new m({
                        actual: new N(e, i),
                        expected: new N(t, i, e),
                        operator: `deepStrictEqual`,
                        stackStartFn: a
                    });
                    throw o.actual = e, o.expected = t, o.operator = a.name, o;
                }
                O({
                    actual: e,
                    expected: t,
                    message: r,
                    operator: a.name,
                    stackStartFn: a
                });
            }
        }
        function ee(e, t, r, i) {
            if (typeof t != `function`) {
                if (v(t)) return x(t, e);
                if (arguments.length === 2) throw new u(`expected`, [
                    `Function`,
                    `RegExp`
                ], t);
                if (n(e) !== `object` || e === null) {
                    var a = new m({
                        actual: e,
                        expected: t,
                        message: r,
                        operator: `deepStrictEqual`,
                        stackStartFn: i
                    });
                    throw a.operator = i.name, a;
                }
                var o = Object.keys(t);
                if (t instanceof Error) o.push(`name`, `message`);
                else if (o.length === 0) throw new d(`error`, t, `may not be an empty object`);
                return S === void 0 && w(), o.forEach(function(n) {
                    typeof e[n] == `string` && v(t[n]) && x(t[n], e[n]) || P(e, t, n, r, o, i);
                }), !0;
            }
            return t.prototype !== void 0 && e instanceof t ? !0 : Error.isPrototypeOf(t) ? !1 : t.call({}, e) === !0;
        }
        function F(e) {
            if (typeof e != `function`) throw new u(`fn`, `Function`, e);
            try {
                e();
            } catch (e) {
                return e;
            }
            return D;
        }
        function te(e) {
            return _(e) || e !== null && n(e) === `object` && typeof e.then == `function` && typeof e.catch == `function`;
        }
        function I(e) {
            return Promise.resolve().then(function() {
                var t;
                if (typeof e == `function`) {
                    if (t = e(), !te(t)) throw new f(`instance of Promise`, `promiseFn`, t);
                } else if (te(e)) t = e;
                else throw new u(`promiseFn`, [
                    `Function`,
                    `Promise`
                ], e);
                return Promise.resolve().then(function() {
                    return t;
                }).then(function() {
                    return D;
                }).catch(function(e) {
                    return e;
                });
            });
        }
        function ne(e, t, r, i) {
            if (typeof r == `string`) {
                if (arguments.length === 4) throw new u(`error`, [
                    `Object`,
                    `Error`,
                    `Function`,
                    `RegExp`
                ], r);
                if (n(t) === `object` && t !== null) {
                    if (t.message === r) throw new l(`error/message`, `The error message "${t.message}" is identical to the message.`);
                } else if (t === r) throw new l(`error/message`, `The error "${t}" is identical to the message.`);
                i = r, r = void 0;
            } else if (r != null && n(r) !== `object` && typeof r != `function`) throw new u(`error`, [
                `Object`,
                `Error`,
                `Function`,
                `RegExp`
            ], r);
            if (t === D) {
                var a = ``;
                r && r.name && (a += ` (${r.name})`), a += i ? `: ${i}` : `.`;
                var o = e.name === `rejects` ? `rejection` : `exception`;
                O({
                    actual: void 0,
                    expected: r,
                    operator: e.name,
                    message: `Missing expected ${o}${a}`,
                    stackStartFn: e
                });
            }
            if (r && !ee(t, r, i, e)) throw t;
        }
        function re(e, t, n, r) {
            if (t !== D) {
                if (typeof n == `string` && (r = n, n = void 0), !n || ee(t, n)) {
                    var i = r ? `: ${r}` : `.`, a = e.name === `doesNotReject` ? `rejection` : `exception`;
                    O({
                        actual: t,
                        expected: n,
                        operator: e.name,
                        message: `Got unwanted ${a}${i}
Actual message: "${t && t.message}"`,
                        stackStartFn: e
                    });
                }
                throw t;
            }
        }
        E.throws = function e(t) {
            var n = [
                ...arguments
            ].slice(1);
            ne.apply(void 0, [
                e,
                F(t)
            ].concat(n));
        }, E.rejects = function e(t) {
            var n = [
                ...arguments
            ].slice(1);
            return I(t).then(function(t) {
                return ne.apply(void 0, [
                    e,
                    t
                ].concat(n));
            });
        }, E.doesNotThrow = function e(t) {
            var n = [
                ...arguments
            ].slice(1);
            re.apply(void 0, [
                e,
                F(t)
            ].concat(n));
        }, E.doesNotReject = function e(t) {
            var n = [
                ...arguments
            ].slice(1);
            return I(t).then(function(t) {
                return re.apply(void 0, [
                    e,
                    t
                ].concat(n));
            });
        }, E.ifError = function e(t) {
            if (t != null) {
                var r = `ifError got unwanted exception: `;
                n(t) === `object` && typeof t.message == `string` ? t.message.length === 0 && t.constructor ? r += t.constructor.name : r += t.message : r += h(t);
                var i = new m({
                    actual: t,
                    expected: null,
                    operator: `ifError`,
                    message: r,
                    stackStartFn: e
                }), a = t.stack;
                if (typeof a == `string`) {
                    var o = a.split(`
`);
                    o.shift();
                    for(var s = i.stack.split(`
`), c = 0; c < o.length; c++){
                        var l = s.indexOf(o[c]);
                        if (l !== -1) {
                            s = s.slice(0, l);
                            break;
                        }
                    }
                    i.stack = `${s.join(`
`)}
${o.join(`
`)}`;
                }
                throw i;
            }
        };
        function L(e, t, r, i, a) {
            if (!v(t)) throw new u(`regexp`, `RegExp`, t);
            var o = a === `match`;
            if (typeof e != `string` || x(t, e) !== o) {
                if (r instanceof Error) throw r;
                var s = !r;
                r ||= typeof e == `string` ? (o ? `The input did not match the regular expression ` : `The input was expected to not match the regular expression `) + `${h(t)}. Input:

${h(e)}
` : `The "string" argument must be of type string. Received type ${n(e)} (${h(e)})`;
                var c = new m({
                    actual: e,
                    expected: t,
                    message: r,
                    operator: a,
                    stackStartFn: i
                });
                throw c.generatedMessage = s, c;
            }
        }
        E.match = function e(t, n, r) {
            L(t, n, r, e, `match`);
        }, E.doesNotMatch = function e(t, n, r) {
            L(t, n, r, e, `doesNotMatch`);
        };
        function R() {
            var e = [
                ...arguments
            ];
            A.apply(void 0, [
                R,
                e.length
            ].concat(e));
        }
        E.strict = y(R, E, {
            equal: E.strictEqual,
            deepEqual: E.deepStrictEqual,
            notEqual: E.notStrictEqual,
            notDeepEqual: E.notDeepStrictEqual
        }), E.strict.strict = E.strict;
    })), tt = t(((e, t)=>{
        function n() {
            this.input = null, this.next_in = 0, this.avail_in = 0, this.total_in = 0, this.output = null, this.next_out = 0, this.avail_out = 0, this.total_out = 0, this.msg = ``, this.state = null, this.data_type = 2, this.adler = 0;
        }
        t.exports = n;
    })), nt = t(((e)=>{
        var t = typeof Uint8Array < `u` && typeof Uint16Array < `u` && typeof Int32Array < `u`;
        function n(e, t) {
            return Object.prototype.hasOwnProperty.call(e, t);
        }
        e.assign = function(e) {
            for(var t = Array.prototype.slice.call(arguments, 1); t.length;){
                var r = t.shift();
                if (r) {
                    if (typeof r != `object`) throw TypeError(r + `must be non-object`);
                    for(var i in r)n(r, i) && (e[i] = r[i]);
                }
            }
            return e;
        }, e.shrinkBuf = function(e, t) {
            return e.length === t ? e : e.subarray ? e.subarray(0, t) : (e.length = t, e);
        };
        var r = {
            arraySet: function(e, t, n, r, i) {
                if (t.subarray && e.subarray) {
                    e.set(t.subarray(n, n + r), i);
                    return;
                }
                for(var a = 0; a < r; a++)e[i + a] = t[n + a];
            },
            flattenChunks: function(e) {
                var t, n, r = 0, i, a, o;
                for(t = 0, n = e.length; t < n; t++)r += e[t].length;
                for(o = new Uint8Array(r), i = 0, t = 0, n = e.length; t < n; t++)a = e[t], o.set(a, i), i += a.length;
                return o;
            }
        }, i = {
            arraySet: function(e, t, n, r, i) {
                for(var a = 0; a < r; a++)e[i + a] = t[n + a];
            },
            flattenChunks: function(e) {
                return [].concat.apply([], e);
            }
        };
        e.setTyped = function(t) {
            t ? (e.Buf8 = Uint8Array, e.Buf16 = Uint16Array, e.Buf32 = Int32Array, e.assign(e, r)) : (e.Buf8 = Array, e.Buf16 = Array, e.Buf32 = Array, e.assign(e, i));
        }, e.setTyped(t);
    })), rt = t(((e)=>{
        var t = nt(), n = 4, r = 0, i = 1, a = 2;
        function o(e) {
            for(var t = e.length; --t >= 0;)e[t] = 0;
        }
        var s = 0, c = 1, l = 2, u = 3, d = 258, f = 29, p = 256, m = p + 1 + f, h = 30, g = 19, _ = 2 * m + 1, v = 15, y = 16, b = 7, x = 256, S = 16, C = 17, w = 18, T = [
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            1,
            1,
            1,
            1,
            2,
            2,
            2,
            2,
            3,
            3,
            3,
            3,
            4,
            4,
            4,
            4,
            5,
            5,
            5,
            5,
            0
        ], E = [
            0,
            0,
            0,
            0,
            1,
            1,
            2,
            2,
            3,
            3,
            4,
            4,
            5,
            5,
            6,
            6,
            7,
            7,
            8,
            8,
            9,
            9,
            10,
            10,
            11,
            11,
            12,
            12,
            13,
            13
        ], D = [
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            2,
            3,
            7
        ], O = [
            16,
            17,
            18,
            0,
            8,
            7,
            9,
            6,
            10,
            5,
            11,
            4,
            12,
            3,
            13,
            2,
            14,
            1,
            15
        ], k = 512, A = Array((m + 2) * 2);
        o(A);
        var j = Array(h * 2);
        o(j);
        var M = Array(k);
        o(M);
        var N = Array(d - u + 1);
        o(N);
        var P = Array(f);
        o(P);
        var ee = Array(h);
        o(ee);
        function F(e, t, n, r, i) {
            this.static_tree = e, this.extra_bits = t, this.extra_base = n, this.elems = r, this.max_length = i, this.has_stree = e && e.length;
        }
        var te, I, ne;
        function re(e, t) {
            this.dyn_tree = e, this.max_code = 0, this.stat_desc = t;
        }
        function L(e) {
            return e < 256 ? M[e] : M[256 + (e >>> 7)];
        }
        function R(e, t) {
            e.pending_buf[e.pending++] = t & 255, e.pending_buf[e.pending++] = t >>> 8 & 255;
        }
        function z(e, t, n) {
            e.bi_valid > y - n ? (e.bi_buf |= t << e.bi_valid & 65535, R(e, e.bi_buf), e.bi_buf = t >> y - e.bi_valid, e.bi_valid += n - y) : (e.bi_buf |= t << e.bi_valid & 65535, e.bi_valid += n);
        }
        function B(e, t, n) {
            z(e, n[t * 2], n[t * 2 + 1]);
        }
        function ie(e, t) {
            var n = 0;
            do n |= e & 1, e >>>= 1, n <<= 1;
            while (--t > 0);
            return n >>> 1;
        }
        function ae(e) {
            e.bi_valid === 16 ? (R(e, e.bi_buf), e.bi_buf = 0, e.bi_valid = 0) : e.bi_valid >= 8 && (e.pending_buf[e.pending++] = e.bi_buf & 255, e.bi_buf >>= 8, e.bi_valid -= 8);
        }
        function oe(e, t) {
            var n = t.dyn_tree, r = t.max_code, i = t.stat_desc.static_tree, a = t.stat_desc.has_stree, o = t.stat_desc.extra_bits, s = t.stat_desc.extra_base, c = t.stat_desc.max_length, l, u, d, f, p, m, h = 0;
            for(f = 0; f <= v; f++)e.bl_count[f] = 0;
            for(n[e.heap[e.heap_max] * 2 + 1] = 0, l = e.heap_max + 1; l < _; l++)u = e.heap[l], f = n[n[u * 2 + 1] * 2 + 1] + 1, f > c && (f = c, h++), n[u * 2 + 1] = f, !(u > r) && (e.bl_count[f]++, p = 0, u >= s && (p = o[u - s]), m = n[u * 2], e.opt_len += m * (f + p), a && (e.static_len += m * (i[u * 2 + 1] + p)));
            if (h !== 0) {
                do {
                    for(f = c - 1; e.bl_count[f] === 0;)f--;
                    e.bl_count[f]--, e.bl_count[f + 1] += 2, e.bl_count[c]--, h -= 2;
                }while (h > 0);
                for(f = c; f !== 0; f--)for(u = e.bl_count[f]; u !== 0;)d = e.heap[--l], !(d > r) && (n[d * 2 + 1] !== f && (e.opt_len += (f - n[d * 2 + 1]) * n[d * 2], n[d * 2 + 1] = f), u--);
            }
        }
        function V(e, t, n) {
            var r = Array(v + 1), i = 0, a, o;
            for(a = 1; a <= v; a++)r[a] = i = i + n[a - 1] << 1;
            for(o = 0; o <= t; o++){
                var s = e[o * 2 + 1];
                s !== 0 && (e[o * 2] = ie(r[s]++, s));
            }
        }
        function H() {
            var e, t, n, r, i, a = Array(v + 1);
            for(n = 0, r = 0; r < f - 1; r++)for(P[r] = n, e = 0; e < 1 << T[r]; e++)N[n++] = r;
            for(N[n - 1] = r, i = 0, r = 0; r < 16; r++)for(ee[r] = i, e = 0; e < 1 << E[r]; e++)M[i++] = r;
            for(i >>= 7; r < h; r++)for(ee[r] = i << 7, e = 0; e < 1 << E[r] - 7; e++)M[256 + i++] = r;
            for(t = 0; t <= v; t++)a[t] = 0;
            for(e = 0; e <= 143;)A[e * 2 + 1] = 8, e++, a[8]++;
            for(; e <= 255;)A[e * 2 + 1] = 9, e++, a[9]++;
            for(; e <= 279;)A[e * 2 + 1] = 7, e++, a[7]++;
            for(; e <= 287;)A[e * 2 + 1] = 8, e++, a[8]++;
            for(V(A, m + 1, a), e = 0; e < h; e++)j[e * 2 + 1] = 5, j[e * 2] = ie(e, 5);
            te = new F(A, T, p + 1, m, v), I = new F(j, E, 0, h, v), ne = new F([], D, 0, g, b);
        }
        function U(e) {
            var t;
            for(t = 0; t < m; t++)e.dyn_ltree[t * 2] = 0;
            for(t = 0; t < h; t++)e.dyn_dtree[t * 2] = 0;
            for(t = 0; t < g; t++)e.bl_tree[t * 2] = 0;
            e.dyn_ltree[x * 2] = 1, e.opt_len = e.static_len = 0, e.last_lit = e.matches = 0;
        }
        function se(e) {
            e.bi_valid > 8 ? R(e, e.bi_buf) : e.bi_valid > 0 && (e.pending_buf[e.pending++] = e.bi_buf), e.bi_buf = 0, e.bi_valid = 0;
        }
        function W(e, n, r, i) {
            se(e), i && (R(e, r), R(e, ~r)), t.arraySet(e.pending_buf, e.window, n, r, e.pending), e.pending += r;
        }
        function G(e, t, n, r) {
            var i = t * 2, a = n * 2;
            return e[i] < e[a] || e[i] === e[a] && r[t] <= r[n];
        }
        function K(e, t, n) {
            for(var r = e.heap[n], i = n << 1; i <= e.heap_len && (i < e.heap_len && G(t, e.heap[i + 1], e.heap[i], e.depth) && i++, !G(t, r, e.heap[i], e.depth));)e.heap[n] = e.heap[i], n = i, i <<= 1;
            e.heap[n] = r;
        }
        function ce(e, t, n) {
            var r, i, a = 0, o, s;
            if (e.last_lit !== 0) do r = e.pending_buf[e.d_buf + a * 2] << 8 | e.pending_buf[e.d_buf + a * 2 + 1], i = e.pending_buf[e.l_buf + a], a++, r === 0 ? B(e, i, t) : (o = N[i], B(e, o + p + 1, t), s = T[o], s !== 0 && (i -= P[o], z(e, i, s)), r--, o = L(r), B(e, o, n), s = E[o], s !== 0 && (r -= ee[o], z(e, r, s)));
            while (a < e.last_lit);
            B(e, x, t);
        }
        function le(e, t) {
            var n = t.dyn_tree, r = t.stat_desc.static_tree, i = t.stat_desc.has_stree, a = t.stat_desc.elems, o, s, c = -1, l;
            for(e.heap_len = 0, e.heap_max = _, o = 0; o < a; o++)n[o * 2] === 0 ? n[o * 2 + 1] = 0 : (e.heap[++e.heap_len] = c = o, e.depth[o] = 0);
            for(; e.heap_len < 2;)l = e.heap[++e.heap_len] = c < 2 ? ++c : 0, n[l * 2] = 1, e.depth[l] = 0, e.opt_len--, i && (e.static_len -= r[l * 2 + 1]);
            for(t.max_code = c, o = e.heap_len >> 1; o >= 1; o--)K(e, n, o);
            l = a;
            do o = e.heap[1], e.heap[1] = e.heap[e.heap_len--], K(e, n, 1), s = e.heap[1], e.heap[--e.heap_max] = o, e.heap[--e.heap_max] = s, n[l * 2] = n[o * 2] + n[s * 2], e.depth[l] = (e.depth[o] >= e.depth[s] ? e.depth[o] : e.depth[s]) + 1, n[o * 2 + 1] = n[s * 2 + 1] = l, e.heap[1] = l++, K(e, n, 1);
            while (e.heap_len >= 2);
            e.heap[--e.heap_max] = e.heap[1], oe(e, t), V(n, c, e.bl_count);
        }
        function q(e, t, n) {
            var r, i = -1, a, o = t[1], s = 0, c = 7, l = 4;
            for(o === 0 && (c = 138, l = 3), t[(n + 1) * 2 + 1] = 65535, r = 0; r <= n; r++)a = o, o = t[(r + 1) * 2 + 1], !(++s < c && a === o) && (s < l ? e.bl_tree[a * 2] += s : a === 0 ? s <= 10 ? e.bl_tree[C * 2]++ : e.bl_tree[w * 2]++ : (a !== i && e.bl_tree[a * 2]++, e.bl_tree[S * 2]++), s = 0, i = a, o === 0 ? (c = 138, l = 3) : a === o ? (c = 6, l = 3) : (c = 7, l = 4));
        }
        function ue(e, t, n) {
            var r, i = -1, a, o = t[1], s = 0, c = 7, l = 4;
            for(o === 0 && (c = 138, l = 3), r = 0; r <= n; r++)if (a = o, o = t[(r + 1) * 2 + 1], !(++s < c && a === o)) {
                if (s < l) do B(e, a, e.bl_tree);
                while (--s !== 0);
                else a === 0 ? s <= 10 ? (B(e, C, e.bl_tree), z(e, s - 3, 3)) : (B(e, w, e.bl_tree), z(e, s - 11, 7)) : (a !== i && (B(e, a, e.bl_tree), s--), B(e, S, e.bl_tree), z(e, s - 3, 2));
                s = 0, i = a, o === 0 ? (c = 138, l = 3) : a === o ? (c = 6, l = 3) : (c = 7, l = 4);
            }
        }
        function J(e) {
            var t;
            for(q(e, e.dyn_ltree, e.l_desc.max_code), q(e, e.dyn_dtree, e.d_desc.max_code), le(e, e.bl_desc), t = g - 1; t >= 3 && e.bl_tree[O[t] * 2 + 1] === 0; t--);
            return e.opt_len += 3 * (t + 1) + 5 + 5 + 4, t;
        }
        function Y(e, t, n, r) {
            var i;
            for(z(e, t - 257, 5), z(e, n - 1, 5), z(e, r - 4, 4), i = 0; i < r; i++)z(e, e.bl_tree[O[i] * 2 + 1], 3);
            ue(e, e.dyn_ltree, t - 1), ue(e, e.dyn_dtree, n - 1);
        }
        function X(e) {
            var t = 4093624447, n;
            for(n = 0; n <= 31; n++, t >>>= 1)if (t & 1 && e.dyn_ltree[n * 2] !== 0) return r;
            if (e.dyn_ltree[18] !== 0 || e.dyn_ltree[20] !== 0 || e.dyn_ltree[26] !== 0) return i;
            for(n = 32; n < p; n++)if (e.dyn_ltree[n * 2] !== 0) return i;
            return r;
        }
        var Z = !1;
        function de(e) {
            Z ||= (H(), !0), e.l_desc = new re(e.dyn_ltree, te), e.d_desc = new re(e.dyn_dtree, I), e.bl_desc = new re(e.bl_tree, ne), e.bi_buf = 0, e.bi_valid = 0, U(e);
        }
        function Q(e, t, n, r) {
            z(e, (s << 1) + +!!r, 3), W(e, t, n, !0);
        }
        function fe(e) {
            z(e, c << 1, 3), B(e, x, A), ae(e);
        }
        function pe(e, t, r, i) {
            var o, s, u = 0;
            e.level > 0 ? (e.strm.data_type === a && (e.strm.data_type = X(e)), le(e, e.l_desc), le(e, e.d_desc), u = J(e), o = e.opt_len + 3 + 7 >>> 3, s = e.static_len + 3 + 7 >>> 3, s <= o && (o = s)) : o = s = r + 5, r + 4 <= o && t !== -1 ? Q(e, t, r, i) : e.strategy === n || s === o ? (z(e, (c << 1) + +!!i, 3), ce(e, A, j)) : (z(e, (l << 1) + +!!i, 3), Y(e, e.l_desc.max_code + 1, e.d_desc.max_code + 1, u + 1), ce(e, e.dyn_ltree, e.dyn_dtree)), U(e), i && se(e);
        }
        function $(e, t, n) {
            return e.pending_buf[e.d_buf + e.last_lit * 2] = t >>> 8 & 255, e.pending_buf[e.d_buf + e.last_lit * 2 + 1] = t & 255, e.pending_buf[e.l_buf + e.last_lit] = n & 255, e.last_lit++, t === 0 ? e.dyn_ltree[n * 2]++ : (e.matches++, t--, e.dyn_ltree[(N[n] + p + 1) * 2]++, e.dyn_dtree[L(t) * 2]++), e.last_lit === e.lit_bufsize - 1;
        }
        e._tr_init = de, e._tr_stored_block = Q, e._tr_flush_block = pe, e._tr_tally = $, e._tr_align = fe;
    })), it = t(((e, t)=>{
        function n(e, t, n, r) {
            for(var i = e & 65535 | 0, a = e >>> 16 & 65535 | 0, o = 0; n !== 0;){
                o = n > 2e3 ? 2e3 : n, n -= o;
                do i = i + t[r++] | 0, a = a + i | 0;
                while (--o);
                i %= 65521, a %= 65521;
            }
            return i | a << 16 | 0;
        }
        t.exports = n;
    })), at = t(((e, t)=>{
        function n() {
            for(var e, t = [], n = 0; n < 256; n++){
                e = n;
                for(var r = 0; r < 8; r++)e = e & 1 ? 3988292384 ^ e >>> 1 : e >>> 1;
                t[n] = e;
            }
            return t;
        }
        var r = n();
        function i(e, t, n, i) {
            var a = r, o = i + n;
            e ^= -1;
            for(var s = i; s < o; s++)e = e >>> 8 ^ a[(e ^ t[s]) & 255];
            return e ^ -1;
        }
        t.exports = i;
    })), ot = t(((e, t)=>{
        t.exports = {
            2: `need dictionary`,
            1: `stream end`,
            0: ``,
            "-1": `file error`,
            "-2": `stream error`,
            "-3": `data error`,
            "-4": `insufficient memory`,
            "-5": `buffer error`,
            "-6": `incompatible version`
        };
    })), st = t(((e)=>{
        var t = nt(), n = rt(), r = it(), i = at(), a = ot(), o = 0, s = 1, c = 3, l = 4, u = 5, d = 0, f = 1, p = -2, m = -3, h = -5, g = -1, _ = 1, v = 2, y = 3, b = 4, x = 0, S = 2, C = 8, w = 9, T = 15, E = 8, D = 286, O = 30, k = 19, A = 2 * D + 1, j = 15, M = 3, N = 258, P = N + M + 1, ee = 32, F = 42, te = 69, I = 73, ne = 91, re = 103, L = 113, R = 666, z = 1, B = 2, ie = 3, ae = 4, oe = 3;
        function V(e, t) {
            return e.msg = a[t], t;
        }
        function H(e) {
            return (e << 1) - (e > 4 ? 9 : 0);
        }
        function U(e) {
            for(var t = e.length; --t >= 0;)e[t] = 0;
        }
        function se(e) {
            var n = e.state, r = n.pending;
            r > e.avail_out && (r = e.avail_out), r !== 0 && (t.arraySet(e.output, n.pending_buf, n.pending_out, r, e.next_out), e.next_out += r, n.pending_out += r, e.total_out += r, e.avail_out -= r, n.pending -= r, n.pending === 0 && (n.pending_out = 0));
        }
        function W(e, t) {
            n._tr_flush_block(e, e.block_start >= 0 ? e.block_start : -1, e.strstart - e.block_start, t), e.block_start = e.strstart, se(e.strm);
        }
        function G(e, t) {
            e.pending_buf[e.pending++] = t;
        }
        function K(e, t) {
            e.pending_buf[e.pending++] = t >>> 8 & 255, e.pending_buf[e.pending++] = t & 255;
        }
        function ce(e, n, a, o) {
            var s = e.avail_in;
            return s > o && (s = o), s === 0 ? 0 : (e.avail_in -= s, t.arraySet(n, e.input, e.next_in, s, a), e.state.wrap === 1 ? e.adler = r(e.adler, n, s, a) : e.state.wrap === 2 && (e.adler = i(e.adler, n, s, a)), e.next_in += s, e.total_in += s, s);
        }
        function le(e, t) {
            var n = e.max_chain_length, r = e.strstart, i, a, o = e.prev_length, s = e.nice_match, c = e.strstart > e.w_size - P ? e.strstart - (e.w_size - P) : 0, l = e.window, u = e.w_mask, d = e.prev, f = e.strstart + N, p = l[r + o - 1], m = l[r + o];
            e.prev_length >= e.good_match && (n >>= 2), s > e.lookahead && (s = e.lookahead);
            do {
                if (i = t, l[i + o] !== m || l[i + o - 1] !== p || l[i] !== l[r] || l[++i] !== l[r + 1]) continue;
                r += 2, i++;
                do ;
                while (l[++r] === l[++i] && l[++r] === l[++i] && l[++r] === l[++i] && l[++r] === l[++i] && l[++r] === l[++i] && l[++r] === l[++i] && l[++r] === l[++i] && l[++r] === l[++i] && r < f);
                if (a = N - (f - r), r = f - N, a > o) {
                    if (e.match_start = t, o = a, a >= s) break;
                    p = l[r + o - 1], m = l[r + o];
                }
            }while ((t = d[t & u]) > c && --n !== 0);
            return o <= e.lookahead ? o : e.lookahead;
        }
        function q(e) {
            var n = e.w_size, r, i, a, o, s;
            do {
                if (o = e.window_size - e.lookahead - e.strstart, e.strstart >= n + (n - P)) {
                    t.arraySet(e.window, e.window, n, n, 0), e.match_start -= n, e.strstart -= n, e.block_start -= n, i = e.hash_size, r = i;
                    do a = e.head[--r], e.head[r] = a >= n ? a - n : 0;
                    while (--i);
                    i = n, r = i;
                    do a = e.prev[--r], e.prev[r] = a >= n ? a - n : 0;
                    while (--i);
                    o += n;
                }
                if (e.strm.avail_in === 0) break;
                if (i = ce(e.strm, e.window, e.strstart + e.lookahead, o), e.lookahead += i, e.lookahead + e.insert >= M) for(s = e.strstart - e.insert, e.ins_h = e.window[s], e.ins_h = (e.ins_h << e.hash_shift ^ e.window[s + 1]) & e.hash_mask; e.insert && (e.ins_h = (e.ins_h << e.hash_shift ^ e.window[s + M - 1]) & e.hash_mask, e.prev[s & e.w_mask] = e.head[e.ins_h], e.head[e.ins_h] = s, s++, e.insert--, !(e.lookahead + e.insert < M)););
            }while (e.lookahead < P && e.strm.avail_in !== 0);
        }
        function ue(e, t) {
            var n = 65535;
            for(n > e.pending_buf_size - 5 && (n = e.pending_buf_size - 5);;){
                if (e.lookahead <= 1) {
                    if (q(e), e.lookahead === 0 && t === o) return z;
                    if (e.lookahead === 0) break;
                }
                e.strstart += e.lookahead, e.lookahead = 0;
                var r = e.block_start + n;
                if ((e.strstart === 0 || e.strstart >= r) && (e.lookahead = e.strstart - r, e.strstart = r, W(e, !1), e.strm.avail_out === 0) || e.strstart - e.block_start >= e.w_size - P && (W(e, !1), e.strm.avail_out === 0)) return z;
            }
            return e.insert = 0, t === l ? (W(e, !0), e.strm.avail_out === 0 ? ie : ae) : (e.strstart > e.block_start && (W(e, !1), e.strm.avail_out), z);
        }
        function J(e, t) {
            for(var r, i;;){
                if (e.lookahead < P) {
                    if (q(e), e.lookahead < P && t === o) return z;
                    if (e.lookahead === 0) break;
                }
                if (r = 0, e.lookahead >= M && (e.ins_h = (e.ins_h << e.hash_shift ^ e.window[e.strstart + M - 1]) & e.hash_mask, r = e.prev[e.strstart & e.w_mask] = e.head[e.ins_h], e.head[e.ins_h] = e.strstart), r !== 0 && e.strstart - r <= e.w_size - P && (e.match_length = le(e, r)), e.match_length >= M) if (i = n._tr_tally(e, e.strstart - e.match_start, e.match_length - M), e.lookahead -= e.match_length, e.match_length <= e.max_lazy_match && e.lookahead >= M) {
                    e.match_length--;
                    do e.strstart++, e.ins_h = (e.ins_h << e.hash_shift ^ e.window[e.strstart + M - 1]) & e.hash_mask, r = e.prev[e.strstart & e.w_mask] = e.head[e.ins_h], e.head[e.ins_h] = e.strstart;
                    while (--e.match_length !== 0);
                    e.strstart++;
                } else e.strstart += e.match_length, e.match_length = 0, e.ins_h = e.window[e.strstart], e.ins_h = (e.ins_h << e.hash_shift ^ e.window[e.strstart + 1]) & e.hash_mask;
                else i = n._tr_tally(e, 0, e.window[e.strstart]), e.lookahead--, e.strstart++;
                if (i && (W(e, !1), e.strm.avail_out === 0)) return z;
            }
            return e.insert = e.strstart < M - 1 ? e.strstart : M - 1, t === l ? (W(e, !0), e.strm.avail_out === 0 ? ie : ae) : e.last_lit && (W(e, !1), e.strm.avail_out === 0) ? z : B;
        }
        function Y(e, t) {
            for(var r, i, a;;){
                if (e.lookahead < P) {
                    if (q(e), e.lookahead < P && t === o) return z;
                    if (e.lookahead === 0) break;
                }
                if (r = 0, e.lookahead >= M && (e.ins_h = (e.ins_h << e.hash_shift ^ e.window[e.strstart + M - 1]) & e.hash_mask, r = e.prev[e.strstart & e.w_mask] = e.head[e.ins_h], e.head[e.ins_h] = e.strstart), e.prev_length = e.match_length, e.prev_match = e.match_start, e.match_length = M - 1, r !== 0 && e.prev_length < e.max_lazy_match && e.strstart - r <= e.w_size - P && (e.match_length = le(e, r), e.match_length <= 5 && (e.strategy === _ || e.match_length === M && e.strstart - e.match_start > 4096) && (e.match_length = M - 1)), e.prev_length >= M && e.match_length <= e.prev_length) {
                    a = e.strstart + e.lookahead - M, i = n._tr_tally(e, e.strstart - 1 - e.prev_match, e.prev_length - M), e.lookahead -= e.prev_length - 1, e.prev_length -= 2;
                    do ++e.strstart <= a && (e.ins_h = (e.ins_h << e.hash_shift ^ e.window[e.strstart + M - 1]) & e.hash_mask, r = e.prev[e.strstart & e.w_mask] = e.head[e.ins_h], e.head[e.ins_h] = e.strstart);
                    while (--e.prev_length !== 0);
                    if (e.match_available = 0, e.match_length = M - 1, e.strstart++, i && (W(e, !1), e.strm.avail_out === 0)) return z;
                } else if (e.match_available) {
                    if (i = n._tr_tally(e, 0, e.window[e.strstart - 1]), i && W(e, !1), e.strstart++, e.lookahead--, e.strm.avail_out === 0) return z;
                } else e.match_available = 1, e.strstart++, e.lookahead--;
            }
            return e.match_available &&= (i = n._tr_tally(e, 0, e.window[e.strstart - 1]), 0), e.insert = e.strstart < M - 1 ? e.strstart : M - 1, t === l ? (W(e, !0), e.strm.avail_out === 0 ? ie : ae) : e.last_lit && (W(e, !1), e.strm.avail_out === 0) ? z : B;
        }
        function X(e, t) {
            for(var r, i, a, s, c = e.window;;){
                if (e.lookahead <= N) {
                    if (q(e), e.lookahead <= N && t === o) return z;
                    if (e.lookahead === 0) break;
                }
                if (e.match_length = 0, e.lookahead >= M && e.strstart > 0 && (a = e.strstart - 1, i = c[a], i === c[++a] && i === c[++a] && i === c[++a])) {
                    s = e.strstart + N;
                    do ;
                    while (i === c[++a] && i === c[++a] && i === c[++a] && i === c[++a] && i === c[++a] && i === c[++a] && i === c[++a] && i === c[++a] && a < s);
                    e.match_length = N - (s - a), e.match_length > e.lookahead && (e.match_length = e.lookahead);
                }
                if (e.match_length >= M ? (r = n._tr_tally(e, 1, e.match_length - M), e.lookahead -= e.match_length, e.strstart += e.match_length, e.match_length = 0) : (r = n._tr_tally(e, 0, e.window[e.strstart]), e.lookahead--, e.strstart++), r && (W(e, !1), e.strm.avail_out === 0)) return z;
            }
            return e.insert = 0, t === l ? (W(e, !0), e.strm.avail_out === 0 ? ie : ae) : e.last_lit && (W(e, !1), e.strm.avail_out === 0) ? z : B;
        }
        function Z(e, t) {
            for(var r;;){
                if (e.lookahead === 0 && (q(e), e.lookahead === 0)) {
                    if (t === o) return z;
                    break;
                }
                if (e.match_length = 0, r = n._tr_tally(e, 0, e.window[e.strstart]), e.lookahead--, e.strstart++, r && (W(e, !1), e.strm.avail_out === 0)) return z;
            }
            return e.insert = 0, t === l ? (W(e, !0), e.strm.avail_out === 0 ? ie : ae) : e.last_lit && (W(e, !1), e.strm.avail_out === 0) ? z : B;
        }
        function de(e, t, n, r, i) {
            this.good_length = e, this.max_lazy = t, this.nice_length = n, this.max_chain = r, this.func = i;
        }
        var Q = [
            new de(0, 0, 0, 0, ue),
            new de(4, 4, 8, 4, J),
            new de(4, 5, 16, 8, J),
            new de(4, 6, 32, 32, J),
            new de(4, 4, 16, 16, Y),
            new de(8, 16, 32, 32, Y),
            new de(8, 16, 128, 128, Y),
            new de(8, 32, 128, 256, Y),
            new de(32, 128, 258, 1024, Y),
            new de(32, 258, 258, 4096, Y)
        ];
        function fe(e) {
            e.window_size = 2 * e.w_size, U(e.head), e.max_lazy_match = Q[e.level].max_lazy, e.good_match = Q[e.level].good_length, e.nice_match = Q[e.level].nice_length, e.max_chain_length = Q[e.level].max_chain, e.strstart = 0, e.block_start = 0, e.lookahead = 0, e.insert = 0, e.match_length = e.prev_length = M - 1, e.match_available = 0, e.ins_h = 0;
        }
        function pe() {
            this.strm = null, this.status = 0, this.pending_buf = null, this.pending_buf_size = 0, this.pending_out = 0, this.pending = 0, this.wrap = 0, this.gzhead = null, this.gzindex = 0, this.method = C, this.last_flush = -1, this.w_size = 0, this.w_bits = 0, this.w_mask = 0, this.window = null, this.window_size = 0, this.prev = null, this.head = null, this.ins_h = 0, this.hash_size = 0, this.hash_bits = 0, this.hash_mask = 0, this.hash_shift = 0, this.block_start = 0, this.match_length = 0, this.prev_match = 0, this.match_available = 0, this.strstart = 0, this.match_start = 0, this.lookahead = 0, this.prev_length = 0, this.max_chain_length = 0, this.max_lazy_match = 0, this.level = 0, this.strategy = 0, this.good_match = 0, this.nice_match = 0, this.dyn_ltree = new t.Buf16(A * 2), this.dyn_dtree = new t.Buf16((2 * O + 1) * 2), this.bl_tree = new t.Buf16((2 * k + 1) * 2), U(this.dyn_ltree), U(this.dyn_dtree), U(this.bl_tree), this.l_desc = null, this.d_desc = null, this.bl_desc = null, this.bl_count = new t.Buf16(j + 1), this.heap = new t.Buf16(2 * D + 1), U(this.heap), this.heap_len = 0, this.heap_max = 0, this.depth = new t.Buf16(2 * D + 1), U(this.depth), this.l_buf = 0, this.lit_bufsize = 0, this.last_lit = 0, this.d_buf = 0, this.opt_len = 0, this.static_len = 0, this.matches = 0, this.insert = 0, this.bi_buf = 0, this.bi_valid = 0;
        }
        function $(e) {
            var t;
            return !e || !e.state ? V(e, p) : (e.total_in = e.total_out = 0, e.data_type = S, t = e.state, t.pending = 0, t.pending_out = 0, t.wrap < 0 && (t.wrap = -t.wrap), t.status = t.wrap ? F : L, e.adler = t.wrap === 2 ? 0 : 1, t.last_flush = o, n._tr_init(t), d);
        }
        function me(e) {
            var t = $(e);
            return t === d && fe(e.state), t;
        }
        function he(e, t) {
            return !e || !e.state || e.state.wrap !== 2 ? p : (e.state.gzhead = t, d);
        }
        function ge(e, n, r, i, a, o) {
            if (!e) return p;
            var s = 1;
            if (n === g && (n = 6), i < 0 ? (s = 0, i = -i) : i > 15 && (s = 2, i -= 16), a < 1 || a > w || r !== C || i < 8 || i > 15 || n < 0 || n > 9 || o < 0 || o > b) return V(e, p);
            i === 8 && (i = 9);
            var c = new pe;
            return e.state = c, c.strm = e, c.wrap = s, c.gzhead = null, c.w_bits = i, c.w_size = 1 << c.w_bits, c.w_mask = c.w_size - 1, c.hash_bits = a + 7, c.hash_size = 1 << c.hash_bits, c.hash_mask = c.hash_size - 1, c.hash_shift = ~~((c.hash_bits + M - 1) / M), c.window = new t.Buf8(c.w_size * 2), c.head = new t.Buf16(c.hash_size), c.prev = new t.Buf16(c.w_size), c.lit_bufsize = 1 << a + 6, c.pending_buf_size = c.lit_bufsize * 4, c.pending_buf = new t.Buf8(c.pending_buf_size), c.d_buf = 1 * c.lit_bufsize, c.l_buf = 3 * c.lit_bufsize, c.level = n, c.strategy = o, c.method = r, me(e);
        }
        function _e(e, t) {
            return ge(e, t, C, T, E, x);
        }
        function ve(e, t) {
            var r, a, m, g;
            if (!e || !e.state || t > u || t < 0) return e ? V(e, p) : p;
            if (a = e.state, !e.output || !e.input && e.avail_in !== 0 || a.status === R && t !== l) return V(e, e.avail_out === 0 ? h : p);
            if (a.strm = e, r = a.last_flush, a.last_flush = t, a.status === F) if (a.wrap === 2) e.adler = 0, G(a, 31), G(a, 139), G(a, 8), a.gzhead ? (G(a, +!!a.gzhead.text + (a.gzhead.hcrc ? 2 : 0) + (a.gzhead.extra ? 4 : 0) + (a.gzhead.name ? 8 : 0) + (a.gzhead.comment ? 16 : 0)), G(a, a.gzhead.time & 255), G(a, a.gzhead.time >> 8 & 255), G(a, a.gzhead.time >> 16 & 255), G(a, a.gzhead.time >> 24 & 255), G(a, a.level === 9 ? 2 : a.strategy >= v || a.level < 2 ? 4 : 0), G(a, a.gzhead.os & 255), a.gzhead.extra && a.gzhead.extra.length && (G(a, a.gzhead.extra.length & 255), G(a, a.gzhead.extra.length >> 8 & 255)), a.gzhead.hcrc && (e.adler = i(e.adler, a.pending_buf, a.pending, 0)), a.gzindex = 0, a.status = te) : (G(a, 0), G(a, 0), G(a, 0), G(a, 0), G(a, 0), G(a, a.level === 9 ? 2 : a.strategy >= v || a.level < 2 ? 4 : 0), G(a, oe), a.status = L);
            else {
                var _ = C + (a.w_bits - 8 << 4) << 8, b = -1;
                b = a.strategy >= v || a.level < 2 ? 0 : a.level < 6 ? 1 : a.level === 6 ? 2 : 3, _ |= b << 6, a.strstart !== 0 && (_ |= ee), _ += 31 - _ % 31, a.status = L, K(a, _), a.strstart !== 0 && (K(a, e.adler >>> 16), K(a, e.adler & 65535)), e.adler = 1;
            }
            if (a.status === te) if (a.gzhead.extra) {
                for(m = a.pending; a.gzindex < (a.gzhead.extra.length & 65535) && !(a.pending === a.pending_buf_size && (a.gzhead.hcrc && a.pending > m && (e.adler = i(e.adler, a.pending_buf, a.pending - m, m)), se(e), m = a.pending, a.pending === a.pending_buf_size));)G(a, a.gzhead.extra[a.gzindex] & 255), a.gzindex++;
                a.gzhead.hcrc && a.pending > m && (e.adler = i(e.adler, a.pending_buf, a.pending - m, m)), a.gzindex === a.gzhead.extra.length && (a.gzindex = 0, a.status = I);
            } else a.status = I;
            if (a.status === I) if (a.gzhead.name) {
                m = a.pending;
                do {
                    if (a.pending === a.pending_buf_size && (a.gzhead.hcrc && a.pending > m && (e.adler = i(e.adler, a.pending_buf, a.pending - m, m)), se(e), m = a.pending, a.pending === a.pending_buf_size)) {
                        g = 1;
                        break;
                    }
                    g = a.gzindex < a.gzhead.name.length ? a.gzhead.name.charCodeAt(a.gzindex++) & 255 : 0, G(a, g);
                }while (g !== 0);
                a.gzhead.hcrc && a.pending > m && (e.adler = i(e.adler, a.pending_buf, a.pending - m, m)), g === 0 && (a.gzindex = 0, a.status = ne);
            } else a.status = ne;
            if (a.status === ne) if (a.gzhead.comment) {
                m = a.pending;
                do {
                    if (a.pending === a.pending_buf_size && (a.gzhead.hcrc && a.pending > m && (e.adler = i(e.adler, a.pending_buf, a.pending - m, m)), se(e), m = a.pending, a.pending === a.pending_buf_size)) {
                        g = 1;
                        break;
                    }
                    g = a.gzindex < a.gzhead.comment.length ? a.gzhead.comment.charCodeAt(a.gzindex++) & 255 : 0, G(a, g);
                }while (g !== 0);
                a.gzhead.hcrc && a.pending > m && (e.adler = i(e.adler, a.pending_buf, a.pending - m, m)), g === 0 && (a.status = re);
            } else a.status = re;
            if (a.status === re && (a.gzhead.hcrc ? (a.pending + 2 > a.pending_buf_size && se(e), a.pending + 2 <= a.pending_buf_size && (G(a, e.adler & 255), G(a, e.adler >> 8 & 255), e.adler = 0, a.status = L)) : a.status = L), a.pending !== 0) {
                if (se(e), e.avail_out === 0) return a.last_flush = -1, d;
            } else if (e.avail_in === 0 && H(t) <= H(r) && t !== l) return V(e, h);
            if (a.status === R && e.avail_in !== 0) return V(e, h);
            if (e.avail_in !== 0 || a.lookahead !== 0 || t !== o && a.status !== R) {
                var x = a.strategy === v ? Z(a, t) : a.strategy === y ? X(a, t) : Q[a.level].func(a, t);
                if ((x === ie || x === ae) && (a.status = R), x === z || x === ie) return e.avail_out === 0 && (a.last_flush = -1), d;
                if (x === B && (t === s ? n._tr_align(a) : t !== u && (n._tr_stored_block(a, 0, 0, !1), t === c && (U(a.head), a.lookahead === 0 && (a.strstart = 0, a.block_start = 0, a.insert = 0))), se(e), e.avail_out === 0)) return a.last_flush = -1, d;
            }
            return t === l ? a.wrap <= 0 ? f : (a.wrap === 2 ? (G(a, e.adler & 255), G(a, e.adler >> 8 & 255), G(a, e.adler >> 16 & 255), G(a, e.adler >> 24 & 255), G(a, e.total_in & 255), G(a, e.total_in >> 8 & 255), G(a, e.total_in >> 16 & 255), G(a, e.total_in >> 24 & 255)) : (K(a, e.adler >>> 16), K(a, e.adler & 65535)), se(e), a.wrap > 0 && (a.wrap = -a.wrap), a.pending === 0 ? f : d) : d;
        }
        function ye(e) {
            var t;
            return !e || !e.state ? p : (t = e.state.status, t !== F && t !== te && t !== I && t !== ne && t !== re && t !== L && t !== R ? V(e, p) : (e.state = null, t === L ? V(e, m) : d));
        }
        function be(e, n) {
            var i = n.length, a, o, s, c, l, u, f, m;
            if (!e || !e.state || (a = e.state, c = a.wrap, c === 2 || c === 1 && a.status !== F || a.lookahead)) return p;
            for(c === 1 && (e.adler = r(e.adler, n, i, 0)), a.wrap = 0, i >= a.w_size && (c === 0 && (U(a.head), a.strstart = 0, a.block_start = 0, a.insert = 0), m = new t.Buf8(a.w_size), t.arraySet(m, n, i - a.w_size, a.w_size, 0), n = m, i = a.w_size), l = e.avail_in, u = e.next_in, f = e.input, e.avail_in = i, e.next_in = 0, e.input = n, q(a); a.lookahead >= M;){
                o = a.strstart, s = a.lookahead - (M - 1);
                do a.ins_h = (a.ins_h << a.hash_shift ^ a.window[o + M - 1]) & a.hash_mask, a.prev[o & a.w_mask] = a.head[a.ins_h], a.head[a.ins_h] = o, o++;
                while (--s);
                a.strstart = o, a.lookahead = M - 1, q(a);
            }
            return a.strstart += a.lookahead, a.block_start = a.strstart, a.insert = a.lookahead, a.lookahead = 0, a.match_length = a.prev_length = M - 1, a.match_available = 0, e.next_in = u, e.input = f, e.avail_in = l, a.wrap = c, d;
        }
        e.deflateInit = _e, e.deflateInit2 = ge, e.deflateReset = me, e.deflateResetKeep = $, e.deflateSetHeader = he, e.deflate = ve, e.deflateEnd = ye, e.deflateSetDictionary = be, e.deflateInfo = `pako deflate (from Nodeca project)`;
    })), ct = t(((e, t)=>{
        var n = 30, r = 12;
        t.exports = function(e, t) {
            var i = e.state, a = e.next_in, o, s, c, l, u, d, f, p, m, h, g, _, v, y, b, x, S, C, w, T, E, D = e.input, O;
            o = a + (e.avail_in - 5), s = e.next_out, O = e.output, c = s - (t - e.avail_out), l = s + (e.avail_out - 257), u = i.dmax, d = i.wsize, f = i.whave, p = i.wnext, m = i.window, h = i.hold, g = i.bits, _ = i.lencode, v = i.distcode, y = (1 << i.lenbits) - 1, b = (1 << i.distbits) - 1;
            top: do {
                g < 15 && (h += D[a++] << g, g += 8, h += D[a++] << g, g += 8), x = _[h & y];
                dolen: for(;;){
                    if (S = x >>> 24, h >>>= S, g -= S, S = x >>> 16 & 255, S === 0) O[s++] = x & 65535;
                    else if (S & 16) {
                        C = x & 65535, S &= 15, S && (g < S && (h += D[a++] << g, g += 8), C += h & (1 << S) - 1, h >>>= S, g -= S), g < 15 && (h += D[a++] << g, g += 8, h += D[a++] << g, g += 8), x = v[h & b];
                        dodist: for(;;){
                            if (S = x >>> 24, h >>>= S, g -= S, S = x >>> 16 & 255, S & 16) {
                                if (w = x & 65535, S &= 15, g < S && (h += D[a++] << g, g += 8, g < S && (h += D[a++] << g, g += 8)), w += h & (1 << S) - 1, w > u) {
                                    e.msg = `invalid distance too far back`, i.mode = n;
                                    break top;
                                }
                                if (h >>>= S, g -= S, S = s - c, w > S) {
                                    if (S = w - S, S > f && i.sane) {
                                        e.msg = `invalid distance too far back`, i.mode = n;
                                        break top;
                                    }
                                    if (T = 0, E = m, p === 0) {
                                        if (T += d - S, S < C) {
                                            C -= S;
                                            do O[s++] = m[T++];
                                            while (--S);
                                            T = s - w, E = O;
                                        }
                                    } else if (p < S) {
                                        if (T += d + p - S, S -= p, S < C) {
                                            C -= S;
                                            do O[s++] = m[T++];
                                            while (--S);
                                            if (T = 0, p < C) {
                                                S = p, C -= S;
                                                do O[s++] = m[T++];
                                                while (--S);
                                                T = s - w, E = O;
                                            }
                                        }
                                    } else if (T += p - S, S < C) {
                                        C -= S;
                                        do O[s++] = m[T++];
                                        while (--S);
                                        T = s - w, E = O;
                                    }
                                    for(; C > 2;)O[s++] = E[T++], O[s++] = E[T++], O[s++] = E[T++], C -= 3;
                                    C && (O[s++] = E[T++], C > 1 && (O[s++] = E[T++]));
                                } else {
                                    T = s - w;
                                    do O[s++] = O[T++], O[s++] = O[T++], O[s++] = O[T++], C -= 3;
                                    while (C > 2);
                                    C && (O[s++] = O[T++], C > 1 && (O[s++] = O[T++]));
                                }
                            } else if (S & 64) {
                                e.msg = `invalid distance code`, i.mode = n;
                                break top;
                            } else {
                                x = v[(x & 65535) + (h & (1 << S) - 1)];
                                continue dodist;
                            }
                            break;
                        }
                    } else if (!(S & 64)) {
                        x = _[(x & 65535) + (h & (1 << S) - 1)];
                        continue dolen;
                    } else if (S & 32) {
                        i.mode = r;
                        break top;
                    } else {
                        e.msg = `invalid literal/length code`, i.mode = n;
                        break top;
                    }
                    break;
                }
            }while (a < o && s < l);
            C = g >> 3, a -= C, g -= C << 3, h &= (1 << g) - 1, e.next_in = a, e.next_out = s, e.avail_in = a < o ? 5 + (o - a) : 5 - (a - o), e.avail_out = s < l ? 257 + (l - s) : 257 - (s - l), i.hold = h, i.bits = g;
        };
    })), lt = t(((e, t)=>{
        var n = nt(), r = 15, i = 852, a = 592, o = 0, s = 1, c = 2, l = [
            3,
            4,
            5,
            6,
            7,
            8,
            9,
            10,
            11,
            13,
            15,
            17,
            19,
            23,
            27,
            31,
            35,
            43,
            51,
            59,
            67,
            83,
            99,
            115,
            131,
            163,
            195,
            227,
            258,
            0,
            0
        ], u = [
            16,
            16,
            16,
            16,
            16,
            16,
            16,
            16,
            17,
            17,
            17,
            17,
            18,
            18,
            18,
            18,
            19,
            19,
            19,
            19,
            20,
            20,
            20,
            20,
            21,
            21,
            21,
            21,
            16,
            72,
            78
        ], d = [
            1,
            2,
            3,
            4,
            5,
            7,
            9,
            13,
            17,
            25,
            33,
            49,
            65,
            97,
            129,
            193,
            257,
            385,
            513,
            769,
            1025,
            1537,
            2049,
            3073,
            4097,
            6145,
            8193,
            12289,
            16385,
            24577,
            0,
            0
        ], f = [
            16,
            16,
            16,
            16,
            17,
            17,
            18,
            18,
            19,
            19,
            20,
            20,
            21,
            21,
            22,
            22,
            23,
            23,
            24,
            24,
            25,
            25,
            26,
            26,
            27,
            27,
            28,
            28,
            29,
            29,
            64,
            64
        ];
        t.exports = function(e, t, p, m, h, g, _, v) {
            var y = v.bits, b = 0, x = 0, S = 0, C = 0, w = 0, T = 0, E = 0, D = 0, O = 0, k = 0, A, j, M, N, P, ee = null, F = 0, te, I = new n.Buf16(r + 1), ne = new n.Buf16(r + 1), re = null, L = 0, R, z, B;
            for(b = 0; b <= r; b++)I[b] = 0;
            for(x = 0; x < m; x++)I[t[p + x]]++;
            for(w = y, C = r; C >= 1 && I[C] === 0; C--);
            if (w > C && (w = C), C === 0) return h[g++] = 20971520, h[g++] = 20971520, v.bits = 1, 0;
            for(S = 1; S < C && I[S] === 0; S++);
            for(w < S && (w = S), D = 1, b = 1; b <= r; b++)if (D <<= 1, D -= I[b], D < 0) return -1;
            if (D > 0 && (e === o || C !== 1)) return -1;
            for(ne[1] = 0, b = 1; b < r; b++)ne[b + 1] = ne[b] + I[b];
            for(x = 0; x < m; x++)t[p + x] !== 0 && (_[ne[t[p + x]]++] = x);
            if (e === o ? (ee = re = _, te = 19) : e === s ? (ee = l, F -= 257, re = u, L -= 257, te = 256) : (ee = d, re = f, te = -1), k = 0, x = 0, b = S, P = g, T = w, E = 0, M = -1, O = 1 << w, N = O - 1, e === s && O > i || e === c && O > a) return 1;
            for(;;){
                R = b - E, _[x] < te ? (z = 0, B = _[x]) : _[x] > te ? (z = re[L + _[x]], B = ee[F + _[x]]) : (z = 96, B = 0), A = 1 << b - E, j = 1 << T, S = j;
                do j -= A, h[P + (k >> E) + j] = R << 24 | z << 16 | B | 0;
                while (j !== 0);
                for(A = 1 << b - 1; k & A;)A >>= 1;
                if (A === 0 ? k = 0 : (k &= A - 1, k += A), x++, --I[b] === 0) {
                    if (b === C) break;
                    b = t[p + _[x]];
                }
                if (b > w && (k & N) !== M) {
                    for(E === 0 && (E = w), P += S, T = b - E, D = 1 << T; T + E < C && (D -= I[T + E], !(D <= 0));)T++, D <<= 1;
                    if (O += 1 << T, e === s && O > i || e === c && O > a) return 1;
                    M = k & N, h[M] = w << 24 | T << 16 | P - g | 0;
                }
            }
            return k !== 0 && (h[P + k] = b - E << 24 | 4194304), v.bits = w, 0;
        };
    })), ut = t(((e)=>{
        var t = nt(), n = it(), r = at(), i = ct(), a = lt(), o = 0, s = 1, c = 2, l = 4, u = 5, d = 6, f = 0, p = 1, m = 2, h = -2, g = -3, _ = -4, v = -5, y = 8, b = 1, x = 2, S = 3, C = 4, w = 5, T = 6, E = 7, D = 8, O = 9, k = 10, A = 11, j = 12, M = 13, N = 14, P = 15, ee = 16, F = 17, te = 18, I = 19, ne = 20, re = 21, L = 22, R = 23, z = 24, B = 25, ie = 26, ae = 27, oe = 28, V = 29, H = 30, U = 31, se = 32, W = 852, G = 592, K = 15;
        function ce(e) {
            return (e >>> 24 & 255) + (e >>> 8 & 65280) + ((e & 65280) << 8) + ((e & 255) << 24);
        }
        function le() {
            this.mode = 0, this.last = !1, this.wrap = 0, this.havedict = !1, this.flags = 0, this.dmax = 0, this.check = 0, this.total = 0, this.head = null, this.wbits = 0, this.wsize = 0, this.whave = 0, this.wnext = 0, this.window = null, this.hold = 0, this.bits = 0, this.length = 0, this.offset = 0, this.extra = 0, this.lencode = null, this.distcode = null, this.lenbits = 0, this.distbits = 0, this.ncode = 0, this.nlen = 0, this.ndist = 0, this.have = 0, this.next = null, this.lens = new t.Buf16(320), this.work = new t.Buf16(288), this.lendyn = null, this.distdyn = null, this.sane = 0, this.back = 0, this.was = 0;
        }
        function q(e) {
            var n;
            return !e || !e.state ? h : (n = e.state, e.total_in = e.total_out = n.total = 0, e.msg = ``, n.wrap && (e.adler = n.wrap & 1), n.mode = b, n.last = 0, n.havedict = 0, n.dmax = 32768, n.head = null, n.hold = 0, n.bits = 0, n.lencode = n.lendyn = new t.Buf32(W), n.distcode = n.distdyn = new t.Buf32(G), n.sane = 1, n.back = -1, f);
        }
        function ue(e) {
            var t;
            return !e || !e.state ? h : (t = e.state, t.wsize = 0, t.whave = 0, t.wnext = 0, q(e));
        }
        function J(e, t) {
            var n, r;
            return !e || !e.state || (r = e.state, t < 0 ? (n = 0, t = -t) : (n = (t >> 4) + 1, t < 48 && (t &= 15)), t && (t < 8 || t > 15)) ? h : (r.window !== null && r.wbits !== t && (r.window = null), r.wrap = n, r.wbits = t, ue(e));
        }
        function Y(e, t) {
            var n, r;
            return e ? (r = new le, e.state = r, r.window = null, n = J(e, t), n !== f && (e.state = null), n) : h;
        }
        function X(e) {
            return Y(e, K);
        }
        var Z = !0, de, Q;
        function fe(e) {
            if (Z) {
                var n;
                for(de = new t.Buf32(512), Q = new t.Buf32(32), n = 0; n < 144;)e.lens[n++] = 8;
                for(; n < 256;)e.lens[n++] = 9;
                for(; n < 280;)e.lens[n++] = 7;
                for(; n < 288;)e.lens[n++] = 8;
                for(a(s, e.lens, 0, 288, de, 0, e.work, {
                    bits: 9
                }), n = 0; n < 32;)e.lens[n++] = 5;
                a(c, e.lens, 0, 32, Q, 0, e.work, {
                    bits: 5
                }), Z = !1;
            }
            e.lencode = de, e.lenbits = 9, e.distcode = Q, e.distbits = 5;
        }
        function pe(e, n, r, i) {
            var a, o = e.state;
            return o.window === null && (o.wsize = 1 << o.wbits, o.wnext = 0, o.whave = 0, o.window = new t.Buf8(o.wsize)), i >= o.wsize ? (t.arraySet(o.window, n, r - o.wsize, o.wsize, 0), o.wnext = 0, o.whave = o.wsize) : (a = o.wsize - o.wnext, a > i && (a = i), t.arraySet(o.window, n, r - i, a, o.wnext), i -= a, i ? (t.arraySet(o.window, n, r - i, i, 0), o.wnext = i, o.whave = o.wsize) : (o.wnext += a, o.wnext === o.wsize && (o.wnext = 0), o.whave < o.wsize && (o.whave += a))), 0;
        }
        function $(e, W) {
            var G, K, le, q, ue, J, Y, X, Z, de, Q, $, me, he, ge = 0, _e, ve, ye, be, xe, Se, Ce, we, Te = new t.Buf8(4), Ee, De, Oe = [
                16,
                17,
                18,
                0,
                8,
                7,
                9,
                6,
                10,
                5,
                11,
                4,
                12,
                3,
                13,
                2,
                14,
                1,
                15
            ];
            if (!e || !e.state || !e.output || !e.input && e.avail_in !== 0) return h;
            G = e.state, G.mode === j && (G.mode = M), ue = e.next_out, le = e.output, Y = e.avail_out, q = e.next_in, K = e.input, J = e.avail_in, X = G.hold, Z = G.bits, de = J, Q = Y, we = f;
            inf_leave: for(;;)switch(G.mode){
                case b:
                    if (G.wrap === 0) {
                        G.mode = M;
                        break;
                    }
                    for(; Z < 16;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    if (G.wrap & 2 && X === 35615) {
                        G.check = 0, Te[0] = X & 255, Te[1] = X >>> 8 & 255, G.check = r(G.check, Te, 2, 0), X = 0, Z = 0, G.mode = x;
                        break;
                    }
                    if (G.flags = 0, G.head && (G.head.done = !1), !(G.wrap & 1) || (((X & 255) << 8) + (X >> 8)) % 31) {
                        e.msg = `incorrect header check`, G.mode = H;
                        break;
                    }
                    if ((X & 15) !== y) {
                        e.msg = `unknown compression method`, G.mode = H;
                        break;
                    }
                    if (X >>>= 4, Z -= 4, Ce = (X & 15) + 8, G.wbits === 0) G.wbits = Ce;
                    else if (Ce > G.wbits) {
                        e.msg = `invalid window size`, G.mode = H;
                        break;
                    }
                    G.dmax = 1 << Ce, e.adler = G.check = 1, G.mode = X & 512 ? k : j, X = 0, Z = 0;
                    break;
                case x:
                    for(; Z < 16;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    if (G.flags = X, (G.flags & 255) !== y) {
                        e.msg = `unknown compression method`, G.mode = H;
                        break;
                    }
                    if (G.flags & 57344) {
                        e.msg = `unknown header flags set`, G.mode = H;
                        break;
                    }
                    G.head && (G.head.text = X >> 8 & 1), G.flags & 512 && (Te[0] = X & 255, Te[1] = X >>> 8 & 255, G.check = r(G.check, Te, 2, 0)), X = 0, Z = 0, G.mode = S;
                case S:
                    for(; Z < 32;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    G.head && (G.head.time = X), G.flags & 512 && (Te[0] = X & 255, Te[1] = X >>> 8 & 255, Te[2] = X >>> 16 & 255, Te[3] = X >>> 24 & 255, G.check = r(G.check, Te, 4, 0)), X = 0, Z = 0, G.mode = C;
                case C:
                    for(; Z < 16;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    G.head && (G.head.xflags = X & 255, G.head.os = X >> 8), G.flags & 512 && (Te[0] = X & 255, Te[1] = X >>> 8 & 255, G.check = r(G.check, Te, 2, 0)), X = 0, Z = 0, G.mode = w;
                case w:
                    if (G.flags & 1024) {
                        for(; Z < 16;){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        G.length = X, G.head && (G.head.extra_len = X), G.flags & 512 && (Te[0] = X & 255, Te[1] = X >>> 8 & 255, G.check = r(G.check, Te, 2, 0)), X = 0, Z = 0;
                    } else G.head && (G.head.extra = null);
                    G.mode = T;
                case T:
                    if (G.flags & 1024 && ($ = G.length, $ > J && ($ = J), $ && (G.head && (Ce = G.head.extra_len - G.length, G.head.extra || (G.head.extra = Array(G.head.extra_len)), t.arraySet(G.head.extra, K, q, $, Ce)), G.flags & 512 && (G.check = r(G.check, K, $, q)), J -= $, q += $, G.length -= $), G.length)) break inf_leave;
                    G.length = 0, G.mode = E;
                case E:
                    if (G.flags & 2048) {
                        if (J === 0) break inf_leave;
                        $ = 0;
                        do Ce = K[q + $++], G.head && Ce && G.length < 65536 && (G.head.name += String.fromCharCode(Ce));
                        while (Ce && $ < J);
                        if (G.flags & 512 && (G.check = r(G.check, K, $, q)), J -= $, q += $, Ce) break inf_leave;
                    } else G.head && (G.head.name = null);
                    G.length = 0, G.mode = D;
                case D:
                    if (G.flags & 4096) {
                        if (J === 0) break inf_leave;
                        $ = 0;
                        do Ce = K[q + $++], G.head && Ce && G.length < 65536 && (G.head.comment += String.fromCharCode(Ce));
                        while (Ce && $ < J);
                        if (G.flags & 512 && (G.check = r(G.check, K, $, q)), J -= $, q += $, Ce) break inf_leave;
                    } else G.head && (G.head.comment = null);
                    G.mode = O;
                case O:
                    if (G.flags & 512) {
                        for(; Z < 16;){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        if (X !== (G.check & 65535)) {
                            e.msg = `header crc mismatch`, G.mode = H;
                            break;
                        }
                        X = 0, Z = 0;
                    }
                    G.head && (G.head.hcrc = G.flags >> 9 & 1, G.head.done = !0), e.adler = G.check = 0, G.mode = j;
                    break;
                case k:
                    for(; Z < 32;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    e.adler = G.check = ce(X), X = 0, Z = 0, G.mode = A;
                case A:
                    if (G.havedict === 0) return e.next_out = ue, e.avail_out = Y, e.next_in = q, e.avail_in = J, G.hold = X, G.bits = Z, m;
                    e.adler = G.check = 1, G.mode = j;
                case j:
                    if (W === u || W === d) break inf_leave;
                case M:
                    if (G.last) {
                        X >>>= Z & 7, Z -= Z & 7, G.mode = ae;
                        break;
                    }
                    for(; Z < 3;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    switch(G.last = X & 1, X >>>= 1, --Z, X & 3){
                        case 0:
                            G.mode = N;
                            break;
                        case 1:
                            if (fe(G), G.mode = ne, W === d) {
                                X >>>= 2, Z -= 2;
                                break inf_leave;
                            }
                            break;
                        case 2:
                            G.mode = F;
                            break;
                        case 3:
                            e.msg = `invalid block type`, G.mode = H;
                    }
                    X >>>= 2, Z -= 2;
                    break;
                case N:
                    for(X >>>= Z & 7, Z -= Z & 7; Z < 32;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    if ((X & 65535) != (X >>> 16 ^ 65535)) {
                        e.msg = `invalid stored block lengths`, G.mode = H;
                        break;
                    }
                    if (G.length = X & 65535, X = 0, Z = 0, G.mode = P, W === d) break inf_leave;
                case P:
                    G.mode = ee;
                case ee:
                    if ($ = G.length, $) {
                        if ($ > J && ($ = J), $ > Y && ($ = Y), $ === 0) break inf_leave;
                        t.arraySet(le, K, q, $, ue), J -= $, q += $, Y -= $, ue += $, G.length -= $;
                        break;
                    }
                    G.mode = j;
                    break;
                case F:
                    for(; Z < 14;){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    if (G.nlen = (X & 31) + 257, X >>>= 5, Z -= 5, G.ndist = (X & 31) + 1, X >>>= 5, Z -= 5, G.ncode = (X & 15) + 4, X >>>= 4, Z -= 4, G.nlen > 286 || G.ndist > 30) {
                        e.msg = `too many length or distance symbols`, G.mode = H;
                        break;
                    }
                    G.have = 0, G.mode = te;
                case te:
                    for(; G.have < G.ncode;){
                        for(; Z < 3;){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        G.lens[Oe[G.have++]] = X & 7, X >>>= 3, Z -= 3;
                    }
                    for(; G.have < 19;)G.lens[Oe[G.have++]] = 0;
                    if (G.lencode = G.lendyn, G.lenbits = 7, Ee = {
                        bits: G.lenbits
                    }, we = a(o, G.lens, 0, 19, G.lencode, 0, G.work, Ee), G.lenbits = Ee.bits, we) {
                        e.msg = `invalid code lengths set`, G.mode = H;
                        break;
                    }
                    G.have = 0, G.mode = I;
                case I:
                    for(; G.have < G.nlen + G.ndist;){
                        for(; ge = G.lencode[X & (1 << G.lenbits) - 1], _e = ge >>> 24, ve = ge >>> 16 & 255, ye = ge & 65535, !(_e <= Z);){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        if (ye < 16) X >>>= _e, Z -= _e, G.lens[G.have++] = ye;
                        else {
                            if (ye === 16) {
                                for(De = _e + 2; Z < De;){
                                    if (J === 0) break inf_leave;
                                    J--, X += K[q++] << Z, Z += 8;
                                }
                                if (X >>>= _e, Z -= _e, G.have === 0) {
                                    e.msg = `invalid bit length repeat`, G.mode = H;
                                    break;
                                }
                                Ce = G.lens[G.have - 1], $ = 3 + (X & 3), X >>>= 2, Z -= 2;
                            } else if (ye === 17) {
                                for(De = _e + 3; Z < De;){
                                    if (J === 0) break inf_leave;
                                    J--, X += K[q++] << Z, Z += 8;
                                }
                                X >>>= _e, Z -= _e, Ce = 0, $ = 3 + (X & 7), X >>>= 3, Z -= 3;
                            } else {
                                for(De = _e + 7; Z < De;){
                                    if (J === 0) break inf_leave;
                                    J--, X += K[q++] << Z, Z += 8;
                                }
                                X >>>= _e, Z -= _e, Ce = 0, $ = 11 + (X & 127), X >>>= 7, Z -= 7;
                            }
                            if (G.have + $ > G.nlen + G.ndist) {
                                e.msg = `invalid bit length repeat`, G.mode = H;
                                break;
                            }
                            for(; $--;)G.lens[G.have++] = Ce;
                        }
                    }
                    if (G.mode === H) break;
                    if (G.lens[256] === 0) {
                        e.msg = `invalid code -- missing end-of-block`, G.mode = H;
                        break;
                    }
                    if (G.lenbits = 9, Ee = {
                        bits: G.lenbits
                    }, we = a(s, G.lens, 0, G.nlen, G.lencode, 0, G.work, Ee), G.lenbits = Ee.bits, we) {
                        e.msg = `invalid literal/lengths set`, G.mode = H;
                        break;
                    }
                    if (G.distbits = 6, G.distcode = G.distdyn, Ee = {
                        bits: G.distbits
                    }, we = a(c, G.lens, G.nlen, G.ndist, G.distcode, 0, G.work, Ee), G.distbits = Ee.bits, we) {
                        e.msg = `invalid distances set`, G.mode = H;
                        break;
                    }
                    if (G.mode = ne, W === d) break inf_leave;
                case ne:
                    G.mode = re;
                case re:
                    if (J >= 6 && Y >= 258) {
                        e.next_out = ue, e.avail_out = Y, e.next_in = q, e.avail_in = J, G.hold = X, G.bits = Z, i(e, Q), ue = e.next_out, le = e.output, Y = e.avail_out, q = e.next_in, K = e.input, J = e.avail_in, X = G.hold, Z = G.bits, G.mode === j && (G.back = -1);
                        break;
                    }
                    for(G.back = 0; ge = G.lencode[X & (1 << G.lenbits) - 1], _e = ge >>> 24, ve = ge >>> 16 & 255, ye = ge & 65535, !(_e <= Z);){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    if (ve && !(ve & 240)) {
                        for(be = _e, xe = ve, Se = ye; ge = G.lencode[Se + ((X & (1 << be + xe) - 1) >> be)], _e = ge >>> 24, ve = ge >>> 16 & 255, ye = ge & 65535, !(be + _e <= Z);){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        X >>>= be, Z -= be, G.back += be;
                    }
                    if (X >>>= _e, Z -= _e, G.back += _e, G.length = ye, ve === 0) {
                        G.mode = ie;
                        break;
                    }
                    if (ve & 32) {
                        G.back = -1, G.mode = j;
                        break;
                    }
                    if (ve & 64) {
                        e.msg = `invalid literal/length code`, G.mode = H;
                        break;
                    }
                    G.extra = ve & 15, G.mode = L;
                case L:
                    if (G.extra) {
                        for(De = G.extra; Z < De;){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        G.length += X & (1 << G.extra) - 1, X >>>= G.extra, Z -= G.extra, G.back += G.extra;
                    }
                    G.was = G.length, G.mode = R;
                case R:
                    for(; ge = G.distcode[X & (1 << G.distbits) - 1], _e = ge >>> 24, ve = ge >>> 16 & 255, ye = ge & 65535, !(_e <= Z);){
                        if (J === 0) break inf_leave;
                        J--, X += K[q++] << Z, Z += 8;
                    }
                    if (!(ve & 240)) {
                        for(be = _e, xe = ve, Se = ye; ge = G.distcode[Se + ((X & (1 << be + xe) - 1) >> be)], _e = ge >>> 24, ve = ge >>> 16 & 255, ye = ge & 65535, !(be + _e <= Z);){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        X >>>= be, Z -= be, G.back += be;
                    }
                    if (X >>>= _e, Z -= _e, G.back += _e, ve & 64) {
                        e.msg = `invalid distance code`, G.mode = H;
                        break;
                    }
                    G.offset = ye, G.extra = ve & 15, G.mode = z;
                case z:
                    if (G.extra) {
                        for(De = G.extra; Z < De;){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        G.offset += X & (1 << G.extra) - 1, X >>>= G.extra, Z -= G.extra, G.back += G.extra;
                    }
                    if (G.offset > G.dmax) {
                        e.msg = `invalid distance too far back`, G.mode = H;
                        break;
                    }
                    G.mode = B;
                case B:
                    if (Y === 0) break inf_leave;
                    if ($ = Q - Y, G.offset > $) {
                        if ($ = G.offset - $, $ > G.whave && G.sane) {
                            e.msg = `invalid distance too far back`, G.mode = H;
                            break;
                        }
                        $ > G.wnext ? ($ -= G.wnext, me = G.wsize - $) : me = G.wnext - $, $ > G.length && ($ = G.length), he = G.window;
                    } else he = le, me = ue - G.offset, $ = G.length;
                    $ > Y && ($ = Y), Y -= $, G.length -= $;
                    do le[ue++] = he[me++];
                    while (--$);
                    G.length === 0 && (G.mode = re);
                    break;
                case ie:
                    if (Y === 0) break inf_leave;
                    le[ue++] = G.length, Y--, G.mode = re;
                    break;
                case ae:
                    if (G.wrap) {
                        for(; Z < 32;){
                            if (J === 0) break inf_leave;
                            J--, X |= K[q++] << Z, Z += 8;
                        }
                        if (Q -= Y, e.total_out += Q, G.total += Q, Q && (e.adler = G.check = G.flags ? r(G.check, le, Q, ue - Q) : n(G.check, le, Q, ue - Q)), Q = Y, (G.flags ? X : ce(X)) !== G.check) {
                            e.msg = `incorrect data check`, G.mode = H;
                            break;
                        }
                        X = 0, Z = 0;
                    }
                    G.mode = oe;
                case oe:
                    if (G.wrap && G.flags) {
                        for(; Z < 32;){
                            if (J === 0) break inf_leave;
                            J--, X += K[q++] << Z, Z += 8;
                        }
                        if (X !== (G.total & 4294967295)) {
                            e.msg = `incorrect length check`, G.mode = H;
                            break;
                        }
                        X = 0, Z = 0;
                    }
                    G.mode = V;
                case V:
                    we = p;
                    break inf_leave;
                case H:
                    we = g;
                    break inf_leave;
                case U:
                    return _;
                case se:
                default:
                    return h;
            }
            return e.next_out = ue, e.avail_out = Y, e.next_in = q, e.avail_in = J, G.hold = X, G.bits = Z, (G.wsize || Q !== e.avail_out && G.mode < H && (G.mode < ae || W !== l)) && pe(e, e.output, e.next_out, Q - e.avail_out) ? (G.mode = U, _) : (de -= e.avail_in, Q -= e.avail_out, e.total_in += de, e.total_out += Q, G.total += Q, G.wrap && Q && (e.adler = G.check = G.flags ? r(G.check, le, Q, e.next_out - Q) : n(G.check, le, Q, e.next_out - Q)), e.data_type = G.bits + (G.last ? 64 : 0) + (G.mode === j ? 128 : 0) + (G.mode === ne || G.mode === P ? 256 : 0), (de === 0 && Q === 0 || W === l) && we === f && (we = v), we);
        }
        function me(e) {
            if (!e || !e.state) return h;
            var t = e.state;
            return t.window &&= null, e.state = null, f;
        }
        function he(e, t) {
            var n;
            return !e || !e.state || (n = e.state, !(n.wrap & 2)) ? h : (n.head = t, t.done = !1, f);
        }
        function ge(e, t) {
            var r = t.length, i, a, o;
            return !e || !e.state || (i = e.state, i.wrap !== 0 && i.mode !== A) ? h : i.mode === A && (a = 1, a = n(a, t, r, 0), a !== i.check) ? g : (o = pe(e, t, r, r), o ? (i.mode = U, _) : (i.havedict = 1, f));
        }
        e.inflateReset = ue, e.inflateReset2 = J, e.inflateResetKeep = q, e.inflateInit = X, e.inflateInit2 = Y, e.inflate = $, e.inflateEnd = me, e.inflateGetHeader = he, e.inflateSetDictionary = ge, e.inflateInfo = `pako inflate (from Nodeca project)`;
    })), dt = t(((e, t)=>{
        t.exports = {
            Z_NO_FLUSH: 0,
            Z_PARTIAL_FLUSH: 1,
            Z_SYNC_FLUSH: 2,
            Z_FULL_FLUSH: 3,
            Z_FINISH: 4,
            Z_BLOCK: 5,
            Z_TREES: 6,
            Z_OK: 0,
            Z_STREAM_END: 1,
            Z_NEED_DICT: 2,
            Z_ERRNO: -1,
            Z_STREAM_ERROR: -2,
            Z_DATA_ERROR: -3,
            Z_BUF_ERROR: -5,
            Z_NO_COMPRESSION: 0,
            Z_BEST_SPEED: 1,
            Z_BEST_COMPRESSION: 9,
            Z_DEFAULT_COMPRESSION: -1,
            Z_FILTERED: 1,
            Z_HUFFMAN_ONLY: 2,
            Z_RLE: 3,
            Z_FIXED: 4,
            Z_DEFAULT_STRATEGY: 0,
            Z_BINARY: 0,
            Z_TEXT: 1,
            Z_UNKNOWN: 2,
            Z_DEFLATED: 8
        };
    })), ft = t(((e)=>{
        var t = et(), n = tt(), r = st(), i = ut(), a = dt();
        for(var o in a)e[o] = a[o];
        e.NONE = 0, e.DEFLATE = 1, e.INFLATE = 2, e.GZIP = 3, e.GUNZIP = 4, e.DEFLATERAW = 5, e.INFLATERAW = 6, e.UNZIP = 7;
        var s = 31, c = 139;
        function l(t) {
            if (typeof t != `number` || t < e.DEFLATE || t > e.UNZIP) throw TypeError(`Bad argument`);
            this.dictionary = null, this.err = 0, this.flush = 0, this.init_done = !1, this.level = 0, this.memLevel = 0, this.mode = t, this.strategy = 0, this.windowBits = 0, this.write_in_progress = !1, this.pending_close = !1, this.gzip_id_bytes_read = 0;
        }
        l.prototype.close = function() {
            if (this.write_in_progress) {
                this.pending_close = !0;
                return;
            }
            this.pending_close = !1, t(this.init_done, `close before init`), t(this.mode <= e.UNZIP), this.mode === e.DEFLATE || this.mode === e.GZIP || this.mode === e.DEFLATERAW ? r.deflateEnd(this.strm) : (this.mode === e.INFLATE || this.mode === e.GUNZIP || this.mode === e.INFLATERAW || this.mode === e.UNZIP) && i.inflateEnd(this.strm), this.mode = e.NONE, this.dictionary = null;
        }, l.prototype.write = function(e, t, n, r, i, a, o) {
            return this._write(!0, e, t, n, r, i, a, o);
        }, l.prototype.writeSync = function(e, t, n, r, i, a, o) {
            return this._write(!1, e, t, n, r, i, a, o);
        }, l.prototype._write = function(n, r, i, a, o, s, c, l) {
            if (t.equal(arguments.length, 8), t(this.init_done, `write before init`), t(this.mode !== e.NONE, `already finalized`), t.equal(!1, this.write_in_progress, `write already in progress`), t.equal(!1, this.pending_close, `close is pending`), this.write_in_progress = !0, t.equal(!1, r === void 0, `must provide flush value`), this.write_in_progress = !0, r !== e.Z_NO_FLUSH && r !== e.Z_PARTIAL_FLUSH && r !== e.Z_SYNC_FLUSH && r !== e.Z_FULL_FLUSH && r !== e.Z_FINISH && r !== e.Z_BLOCK) throw Error(`Invalid flush value`);
            if (i ?? (i = Buffer.alloc(0), o = 0, a = 0), this.strm.avail_in = o, this.strm.input = i, this.strm.next_in = a, this.strm.avail_out = l, this.strm.output = s, this.strm.next_out = c, this.flush = r, !n) return this._process(), this._checkError() ? this._afterSync() : void 0;
            var u = this;
            return process.nextTick(function() {
                u._process(), u._after();
            }), this;
        }, l.prototype._afterSync = function() {
            var e = this.strm.avail_out, t = this.strm.avail_in;
            return this.write_in_progress = !1, [
                t,
                e
            ];
        }, l.prototype._process = function() {
            var t = null;
            switch(this.mode){
                case e.DEFLATE:
                case e.GZIP:
                case e.DEFLATERAW:
                    this.err = r.deflate(this.strm, this.flush);
                    break;
                case e.UNZIP:
                    switch(this.strm.avail_in > 0 && (t = this.strm.next_in), this.gzip_id_bytes_read){
                        case 0:
                            if (t === null) break;
                            if (this.strm.input[t] === s) {
                                if (this.gzip_id_bytes_read = 1, t++, this.strm.avail_in === 1) break;
                            } else {
                                this.mode = e.INFLATE;
                                break;
                            }
                        case 1:
                            if (t === null) break;
                            this.strm.input[t] === c ? (this.gzip_id_bytes_read = 2, this.mode = e.GUNZIP) : this.mode = e.INFLATE;
                            break;
                        default:
                            throw Error(`invalid number of gzip magic number bytes read`);
                    }
                case e.INFLATE:
                case e.GUNZIP:
                case e.INFLATERAW:
                    for(this.err = i.inflate(this.strm, this.flush), this.err === e.Z_NEED_DICT && this.dictionary && (this.err = i.inflateSetDictionary(this.strm, this.dictionary), this.err === e.Z_OK ? this.err = i.inflate(this.strm, this.flush) : this.err === e.Z_DATA_ERROR && (this.err = e.Z_NEED_DICT)); this.strm.avail_in > 0 && this.mode === e.GUNZIP && this.err === e.Z_STREAM_END && this.strm.next_in[0] !== 0;)this.reset(), this.err = i.inflate(this.strm, this.flush);
                    break;
                default:
                    throw Error(`Unknown mode ` + this.mode);
            }
        }, l.prototype._checkError = function() {
            switch(this.err){
                case e.Z_OK:
                case e.Z_BUF_ERROR:
                    if (this.strm.avail_out !== 0 && this.flush === e.Z_FINISH) return this._error(`unexpected end of file`), !1;
                    break;
                case e.Z_STREAM_END:
                    break;
                case e.Z_NEED_DICT:
                    return this.dictionary == null ? this._error(`Missing dictionary`) : this._error(`Bad dictionary`), !1;
                default:
                    return this._error(`Zlib error`), !1;
            }
            return !0;
        }, l.prototype._after = function() {
            if (this._checkError()) {
                var e = this.strm.avail_out, t = this.strm.avail_in;
                this.write_in_progress = !1, this.callback(t, e), this.pending_close && this.close();
            }
        }, l.prototype._error = function(e) {
            this.strm.msg && (e = this.strm.msg), this.onerror(e, this.err), this.write_in_progress = !1, this.pending_close && this.close();
        }, l.prototype.init = function(n, r, i, a, o) {
            t(arguments.length === 4 || arguments.length === 5, `init(windowBits, level, memLevel, strategy, [dictionary])`), t(n >= 8 && n <= 15, `invalid windowBits`), t(r >= -1 && r <= 9, `invalid compression level`), t(i >= 1 && i <= 9, `invalid memlevel`), t(a === e.Z_FILTERED || a === e.Z_HUFFMAN_ONLY || a === e.Z_RLE || a === e.Z_FIXED || a === e.Z_DEFAULT_STRATEGY, `invalid strategy`), this._init(r, n, i, a, o), this._setDictionary();
        }, l.prototype.params = function() {
            throw Error(`deflateParams Not supported`);
        }, l.prototype.reset = function() {
            this._reset(), this._setDictionary();
        }, l.prototype._init = function(t, a, o, s, c) {
            switch(this.level = t, this.windowBits = a, this.memLevel = o, this.strategy = s, this.flush = e.Z_NO_FLUSH, this.err = e.Z_OK, (this.mode === e.GZIP || this.mode === e.GUNZIP) && (this.windowBits += 16), this.mode === e.UNZIP && (this.windowBits += 32), (this.mode === e.DEFLATERAW || this.mode === e.INFLATERAW) && (this.windowBits = -1 * this.windowBits), this.strm = new n, this.mode){
                case e.DEFLATE:
                case e.GZIP:
                case e.DEFLATERAW:
                    this.err = r.deflateInit2(this.strm, this.level, e.Z_DEFLATED, this.windowBits, this.memLevel, this.strategy);
                    break;
                case e.INFLATE:
                case e.GUNZIP:
                case e.INFLATERAW:
                case e.UNZIP:
                    this.err = i.inflateInit2(this.strm, this.windowBits);
                    break;
                default:
                    throw Error(`Unknown mode ` + this.mode);
            }
            this.err !== e.Z_OK && this._error(`Init error`), this.dictionary = c, this.write_in_progress = !1, this.init_done = !0;
        }, l.prototype._setDictionary = function() {
            if (this.dictionary != null) {
                switch(this.err = e.Z_OK, this.mode){
                    case e.DEFLATE:
                    case e.DEFLATERAW:
                        this.err = r.deflateSetDictionary(this.strm, this.dictionary);
                        break;
                    default:
                        break;
                }
                this.err !== e.Z_OK && this._error(`Failed to set dictionary`);
            }
        }, l.prototype._reset = function() {
            switch(this.err = e.Z_OK, this.mode){
                case e.DEFLATE:
                case e.DEFLATERAW:
                case e.GZIP:
                    this.err = r.deflateReset(this.strm);
                    break;
                case e.INFLATE:
                case e.INFLATERAW:
                case e.GUNZIP:
                    this.err = i.inflateReset(this.strm);
                    break;
                default:
                    break;
            }
            this.err !== e.Z_OK && this._error(`Failed to reset stream`);
        }, e.Zlib = l;
    })), pt = t(((e)=>{
        var t = m().Buffer, n = Fe().Transform, r = ft(), i = ve(), a = et().ok, o = m().kMaxLength, s = `Cannot create final Buffer. It would be larger than 0x` + o.toString(16) + ` bytes`;
        r.Z_MIN_WINDOWBITS = 8, r.Z_MAX_WINDOWBITS = 15, r.Z_DEFAULT_WINDOWBITS = 15, r.Z_MIN_CHUNK = 64, r.Z_MAX_CHUNK = 1 / 0, r.Z_DEFAULT_CHUNK = 16 * 1024, r.Z_MIN_MEMLEVEL = 1, r.Z_MAX_MEMLEVEL = 9, r.Z_DEFAULT_MEMLEVEL = 8, r.Z_MIN_LEVEL = -1, r.Z_MAX_LEVEL = 9, r.Z_DEFAULT_LEVEL = r.Z_DEFAULT_COMPRESSION;
        for(var c = Object.keys(r), l = 0; l < c.length; l++){
            var u = c[l];
            u.match(/^Z/) && Object.defineProperty(e, u, {
                enumerable: !0,
                value: r[u],
                writable: !1
            });
        }
        for(var d = {
            Z_OK: r.Z_OK,
            Z_STREAM_END: r.Z_STREAM_END,
            Z_NEED_DICT: r.Z_NEED_DICT,
            Z_ERRNO: r.Z_ERRNO,
            Z_STREAM_ERROR: r.Z_STREAM_ERROR,
            Z_DATA_ERROR: r.Z_DATA_ERROR,
            Z_MEM_ERROR: r.Z_MEM_ERROR,
            Z_BUF_ERROR: r.Z_BUF_ERROR,
            Z_VERSION_ERROR: r.Z_VERSION_ERROR
        }, f = Object.keys(d), p = 0; p < f.length; p++){
            var h = f[p];
            d[d[h]] = h;
        }
        Object.defineProperty(e, "codes", {
            enumerable: !0,
            value: Object.freeze(d),
            writable: !1
        }), e.Deflate = v, e.Inflate = y, e.Gzip = b, e.Gunzip = x, e.DeflateRaw = S, e.InflateRaw = C, e.Unzip = w, e.createDeflate = function(e) {
            return new v(e);
        }, e.createInflate = function(e) {
            return new y(e);
        }, e.createDeflateRaw = function(e) {
            return new S(e);
        }, e.createInflateRaw = function(e) {
            return new C(e);
        }, e.createGzip = function(e) {
            return new b(e);
        }, e.createGunzip = function(e) {
            return new x(e);
        }, e.createUnzip = function(e) {
            return new w(e);
        }, e.deflate = function(e, t, n) {
            return typeof t == `function` && (n = t, t = {}), g(new v(t), e, n);
        }, e.deflateSync = function(e, t) {
            return _(new v(t), e);
        }, e.gzip = function(e, t, n) {
            return typeof t == `function` && (n = t, t = {}), g(new b(t), e, n);
        }, e.gzipSync = function(e, t) {
            return _(new b(t), e);
        }, e.deflateRaw = function(e, t, n) {
            return typeof t == `function` && (n = t, t = {}), g(new S(t), e, n);
        }, e.deflateRawSync = function(e, t) {
            return _(new S(t), e);
        }, e.unzip = function(e, t, n) {
            return typeof t == `function` && (n = t, t = {}), g(new w(t), e, n);
        }, e.unzipSync = function(e, t) {
            return _(new w(t), e);
        }, e.inflate = function(e, t, n) {
            return typeof t == `function` && (n = t, t = {}), g(new y(t), e, n);
        }, e.inflateSync = function(e, t) {
            return _(new y(t), e);
        }, e.gunzip = function(e, t, n) {
            return typeof t == `function` && (n = t, t = {}), g(new x(t), e, n);
        }, e.gunzipSync = function(e, t) {
            return _(new x(t), e);
        }, e.inflateRaw = function(e, t, n) {
            return typeof t == `function` && (n = t, t = {}), g(new C(t), e, n);
        }, e.inflateRawSync = function(e, t) {
            return _(new C(t), e);
        };
        function g(e, n, r) {
            var i = [], a = 0;
            e.on(`error`, l), e.on(`end`, u), e.end(n), c();
            function c() {
                for(var t; (t = e.read()) !== null;)i.push(t), a += t.length;
                e.once(`readable`, c);
            }
            function l(t) {
                e.removeListener(`end`, u), e.removeListener(`readable`, c), r(t);
            }
            function u() {
                var n, c = null;
                a >= o ? c = RangeError(s) : n = t.concat(i, a), i = [], e.close(), r(c, n);
            }
        }
        function _(e, n) {
            if (typeof n == `string` && (n = t.from(n)), !t.isBuffer(n)) throw TypeError(`Not a string or buffer`);
            var r = e._finishFlushFlag;
            return e._processChunk(n, r);
        }
        function v(e) {
            if (!(this instanceof v)) return new v(e);
            E.call(this, e, r.DEFLATE);
        }
        function y(e) {
            if (!(this instanceof y)) return new y(e);
            E.call(this, e, r.INFLATE);
        }
        function b(e) {
            if (!(this instanceof b)) return new b(e);
            E.call(this, e, r.GZIP);
        }
        function x(e) {
            if (!(this instanceof x)) return new x(e);
            E.call(this, e, r.GUNZIP);
        }
        function S(e) {
            if (!(this instanceof S)) return new S(e);
            E.call(this, e, r.DEFLATERAW);
        }
        function C(e) {
            if (!(this instanceof C)) return new C(e);
            E.call(this, e, r.INFLATERAW);
        }
        function w(e) {
            if (!(this instanceof w)) return new w(e);
            E.call(this, e, r.UNZIP);
        }
        function T(e) {
            return e === r.Z_NO_FLUSH || e === r.Z_PARTIAL_FLUSH || e === r.Z_SYNC_FLUSH || e === r.Z_FULL_FLUSH || e === r.Z_FINISH || e === r.Z_BLOCK;
        }
        function E(i, a) {
            var o = this;
            if (this._opts = i ||= {}, this._chunkSize = i.chunkSize || e.Z_DEFAULT_CHUNK, n.call(this, i), i.flush && !T(i.flush)) throw Error(`Invalid flush flag: ` + i.flush);
            if (i.finishFlush && !T(i.finishFlush)) throw Error(`Invalid flush flag: ` + i.finishFlush);
            if (this._flushFlag = i.flush || r.Z_NO_FLUSH, this._finishFlushFlag = i.finishFlush === void 0 ? r.Z_FINISH : i.finishFlush, i.chunkSize && (i.chunkSize < e.Z_MIN_CHUNK || i.chunkSize > e.Z_MAX_CHUNK)) throw Error(`Invalid chunk size: ` + i.chunkSize);
            if (i.windowBits && (i.windowBits < e.Z_MIN_WINDOWBITS || i.windowBits > e.Z_MAX_WINDOWBITS)) throw Error(`Invalid windowBits: ` + i.windowBits);
            if (i.level && (i.level < e.Z_MIN_LEVEL || i.level > e.Z_MAX_LEVEL)) throw Error(`Invalid compression level: ` + i.level);
            if (i.memLevel && (i.memLevel < e.Z_MIN_MEMLEVEL || i.memLevel > e.Z_MAX_MEMLEVEL)) throw Error(`Invalid memLevel: ` + i.memLevel);
            if (i.strategy && i.strategy != e.Z_FILTERED && i.strategy != e.Z_HUFFMAN_ONLY && i.strategy != e.Z_RLE && i.strategy != e.Z_FIXED && i.strategy != e.Z_DEFAULT_STRATEGY) throw Error(`Invalid strategy: ` + i.strategy);
            if (i.dictionary && !t.isBuffer(i.dictionary)) throw Error(`Invalid dictionary: it should be a Buffer instance`);
            this._handle = new r.Zlib(a);
            var s = this;
            this._hadError = !1, this._handle.onerror = function(t, n) {
                D(s), s._hadError = !0;
                var r = Error(t);
                r.errno = n, r.code = e.codes[n], s.emit(`error`, r);
            };
            var c = e.Z_DEFAULT_COMPRESSION;
            typeof i.level == `number` && (c = i.level);
            var l = e.Z_DEFAULT_STRATEGY;
            typeof i.strategy == `number` && (l = i.strategy), this._handle.init(i.windowBits || e.Z_DEFAULT_WINDOWBITS, c, i.memLevel || e.Z_DEFAULT_MEMLEVEL, l, i.dictionary), this._buffer = t.allocUnsafe(this._chunkSize), this._offset = 0, this._level = c, this._strategy = l, this.once(`end`, this.close), Object.defineProperty(this, "_closed", {
                get: function() {
                    return !o._handle;
                },
                configurable: !0,
                enumerable: !0
            });
        }
        i.inherits(E, n), E.prototype.params = function(t, n, i) {
            if (t < e.Z_MIN_LEVEL || t > e.Z_MAX_LEVEL) throw RangeError(`Invalid compression level: ` + t);
            if (n != e.Z_FILTERED && n != e.Z_HUFFMAN_ONLY && n != e.Z_RLE && n != e.Z_FIXED && n != e.Z_DEFAULT_STRATEGY) throw TypeError(`Invalid strategy: ` + n);
            if (this._level !== t || this._strategy !== n) {
                var o = this;
                this.flush(r.Z_SYNC_FLUSH, function() {
                    a(o._handle, `zlib binding closed`), o._handle.params(t, n), o._hadError || (o._level = t, o._strategy = n, i && i());
                });
            } else process.nextTick(i);
        }, E.prototype.reset = function() {
            return a(this._handle, `zlib binding closed`), this._handle.reset();
        }, E.prototype._flush = function(e) {
            this._transform(t.alloc(0), ``, e);
        }, E.prototype.flush = function(e, n) {
            var i = this, a = this._writableState;
            (typeof e == `function` || e === void 0 && !n) && (n = e, e = r.Z_FULL_FLUSH), a.ended ? n && process.nextTick(n) : a.ending ? n && this.once(`end`, n) : a.needDrain ? n && this.once(`drain`, function() {
                return i.flush(e, n);
            }) : (this._flushFlag = e, this.write(t.alloc(0), ``, n));
        }, E.prototype.close = function(e) {
            D(this, e), process.nextTick(O, this);
        };
        function D(e, t) {
            t && process.nextTick(t), e._handle &&= (e._handle.close(), null);
        }
        function O(e) {
            e.emit(`close`);
        }
        E.prototype._transform = function(e, n, i) {
            var a, o = this._writableState, s = (o.ending || o.ended) && (!e || o.length === e.length);
            if (e !== null && !t.isBuffer(e)) return i(Error(`invalid input`));
            if (!this._handle) return i(Error(`zlib binding closed`));
            s ? a = this._finishFlushFlag : (a = this._flushFlag, e.length >= o.length && (this._flushFlag = this._opts.flush || r.Z_NO_FLUSH)), this._processChunk(e, a, i);
        }, E.prototype._processChunk = function(e, n, r) {
            var i = e && e.length, c = this._chunkSize - this._offset, l = 0, u = this, d = typeof r == `function`;
            if (!d) {
                var f = [], p = 0, m;
                this.on(`error`, function(e) {
                    m = e;
                }), a(this._handle, `zlib binding closed`);
                do var h = this._handle.writeSync(n, e, l, i, this._buffer, this._offset, c);
                while (!this._hadError && v(h[0], h[1]));
                if (this._hadError) throw m;
                if (p >= o) throw D(this), RangeError(s);
                var g = t.concat(f, p);
                return D(this), g;
            }
            a(this._handle, `zlib binding closed`);
            var _ = this._handle.write(n, e, l, i, this._buffer, this._offset, c);
            _.buffer = e, _.callback = v;
            function v(o, s) {
                if (this && (this.buffer = null, this.callback = null), !u._hadError) {
                    var m = c - s;
                    if (a(m >= 0, `have should not go down`), m > 0) {
                        var h = u._buffer.slice(u._offset, u._offset + m);
                        u._offset += m, d ? u.push(h) : (f.push(h), p += h.length);
                    }
                    if ((s === 0 || u._offset >= u._chunkSize) && (c = u._chunkSize, u._offset = 0, u._buffer = t.allocUnsafe(u._chunkSize)), s === 0) {
                        if (l += i - o, i = o, !d) return !0;
                        var g = u._handle.write(n, e, l, i, u._buffer, u._offset, u._chunkSize);
                        g.callback = v, g.buffer = e;
                        return;
                    }
                    if (!d) return !1;
                    r();
                }
            }
        }, i.inherits(v, E), i.inherits(y, E), i.inherits(b, E), i.inherits(x, E), i.inherits(S, E), i.inherits(C, E), i.inherits(w, E);
    })), mt = t(((e, t)=>{
        let n = (e, t, n, r)=>function(...i) {
                let a = t.promiseModule;
                return new a((a, o)=>{
                    t.multiArgs ? i.push((...e)=>{
                        t.errorFirst ? e[0] ? o(e) : (e.shift(), a(e)) : a(e);
                    }) : t.errorFirst ? i.push((e, t)=>{
                        e ? o(e) : a(t);
                    }) : i.push(a);
                    let s = this === n ? r : this;
                    Reflect.apply(e, s, i);
                });
            }, r = new WeakMap;
        t.exports = (e, t)=>{
            t = {
                exclude: [
                    /.+(?:Sync|Stream)$/
                ],
                errorFirst: !0,
                promiseModule: Promise,
                ...t
            };
            let i = typeof e;
            if (!(e !== null && (i === `object` || i === `function`))) throw TypeError(`Expected \`input\` to be a \`Function\` or \`Object\`, got \`${e === null ? `null` : i}\``);
            let a = (e, n)=>{
                let i = r.get(e);
                if (i || (i = {}, r.set(e, i)), n in i) return i[n];
                let a = (e)=>typeof e == `string` || typeof n == `symbol` ? n === e : e.test(n), o = Reflect.getOwnPropertyDescriptor(e, n), s = o === void 0 || o.writable || o.configurable, c = (t.include ? t.include.some(a) : !t.exclude.some(a)) && s;
                return i[n] = c, c;
            }, o = new WeakMap, s = new Proxy(e, {
                apply (e, r, i) {
                    let a = o.get(e);
                    if (a) return Reflect.apply(a, r, i);
                    let c = t.excludeMain ? e : n(e, t, s, e);
                    return o.set(e, c), Reflect.apply(c, r, i);
                },
                get (e, r) {
                    let i = e[r];
                    if (!a(e, r) || i === Function.prototype[r]) return i;
                    let c = o.get(i);
                    if (c) return c;
                    if (typeof i == `function`) {
                        let r = n(i, t, s, e);
                        return o.set(i, r), r;
                    }
                    return i;
                }
            });
            return s;
        };
    })), ht = t(((e, t)=>{
        t.exports = import(`./brotli_wasm-DRMwhWRj.js`).then(async (m)=>{
            await m.__tla;
            return m;
        }), t.exports.default = t.exports, t.exports.BrotliWasmType = void 0;
    })), gt = t(((t, n)=>{
        n.exports = import(`./zstd-codec-DAMX-8_z.js`).then(async (m)=>{
            await m.__tla;
            return m;
        }).then((t)=>e(t.default)).then((e)=>({
                ZstdCodec: e.default.ZstdCodec
            }));
    })), _t = t(((e)=>{
        let t = Fe();
        r().Binding;
        let a = n(), o = i(), s = o.getClassName, c = o.toTypedArray, l = o.fromTypedArrayToBuffer, u = (e)=>{
            class n extends t.Transform {
                constructor(t, n, r){
                    super(r || {}), this.string_decoder = n, this.binding = new e.ZstdCompressStreamBinding, this.binding.begin(t || a.DEFAULT_COMPRESSION_LEVEL), this.callback = (e)=>{
                        this.push(l(e), `buffer`);
                    };
                }
                _transform(e, t, n) {
                    let r = c(e, t, this.string_decoder);
                    if (!r) {
                        let t = s(e) || typeof e;
                        n(Error(`unsupported chunk type: ${t}`));
                        return;
                    }
                    this.binding.transform(r, this.callback) ? n() : n(Error(`ZstdDecompressTransform: Error on _transform`));
                }
                _flush(e) {
                    this.binding.flush(this.callback) ? e() : e(Error(`ZstdDecompressTransform: Error on _flush`));
                }
                _final(e) {
                    this.binding.end(this.callback) ? e() : e(Error(`ZstdDecompressTransform: Error on _final`));
                }
            }
            class r extends t.Transform {
                constructor(t){
                    super(t || {}), this.binding = new e.ZstdDecompressStreamBinding, this.binding.begin(), this.callback = (e)=>{
                        this.push(l(e), `buffer`);
                    };
                }
                _transform(e, t, n) {
                    if (!c(e, t, this.string_decoder)) {
                        let t = s(e) || typeof e;
                        n(Error(`unsupported chunk type: ${t}`));
                        return;
                    }
                    this.binding.transform(e, this.callback) ? n() : n(Error(`ZstdDecompressTransform: Error on _transform`));
                }
                _flush(e) {
                    this.binding.flush(this.callback) ? e() : e(Error(`ZstdDecompressTransform: Error on _flush`));
                }
                _final(e) {
                    this.binding.end(this.callback) ? e() : e(Error(`ZstdDecompressTransform: Error on _final`));
                }
            }
            let i = {};
            return i.ZstdCompressTransform = n, i.ZstdDecompressTransform = r, i;
        };
        e.run = (e)=>r().run((t)=>{
                e(u(t));
            });
    })), vt = t(((e)=>{
        Object.defineProperty(e, "__esModule", {
            value: !0
        }), e.decodeBase64 = e.encodeBase64 = e.zstdDecompress = e.zstdCompress = e.brotliDecompress = e.brotliCompress = e.inflateRaw = e.inflate = e.deflateRaw = e.deflate = e.gunzip = e.gzip = void 0, e.createGzipStream = h, e.createGunzipStream = g, e.createDeflateStream = _, e.createInflateStream = v, e.createDeflateRawStream = y, e.createInflateRawStream = b, e.createBrotliCompressStream = x, e.createBrotliDecompressStream = S, e.createZstdCompressStream = E, e.createZstdDecompressStream = D, e.createBase64EncodeStream = N, e.createBase64DecodeStream = P, e.createDecodeStream = ne, e.createEncodeStream = re, e.decodeBuffer = z, e.decodeBufferSync = B, e.encodeBuffer = ie;
        let t = pt(), n = mt();
        e.gzip = n(t.gzip), e.gunzip = n(t.gunzip), e.deflate = n(t.deflate), e.deflateRaw = n(t.deflateRaw), e.inflate = n(t.inflate), e.inflateRaw = n(t.inflateRaw), e.brotliCompress = t.brotliCompress ? (async (e, n)=>new Promise((r, i)=>{
                t.brotliCompress(e, n === void 0 ? {} : {
                    params: {
                        [t.constants.BROTLI_PARAM_QUALITY]: n
                    }
                }, (e, t)=>{
                    e ? i(e) : r(t);
                });
            })) : (async (e, t)=>{
            let { compress: n } = await Promise.resolve().then(()=>ht());
            return n(e, {
                quality: t
            });
        }), e.brotliDecompress = t.brotliDecompress ? n(t.brotliDecompress) : (async (e)=>{
            let { decompress: t } = await Promise.resolve().then(()=>ht());
            return t(e);
        });
        let r, i = async ()=>t.zstdCompress && t.zstdDecompress ? {
                compress: (e, n)=>new Promise((r, i)=>{
                        let a = n === void 0 ? {} : {
                            [t.constants.ZSTD_c_compressionLevel]: n
                        };
                        t.zstdCompress(e, a, (e, t)=>{
                            e ? i(e) : r(t);
                        });
                    }),
                decompress: (e)=>new Promise((n, r)=>{
                        t.zstdDecompress(e, (e, t)=>{
                            e ? r(e) : n(t);
                        });
                    })
            } : (r ||= new Promise(async (e)=>{
                let { ZstdCodec: t } = await Promise.resolve().then(()=>gt());
                t.run((t)=>{
                    e(new t.Streaming);
                });
            }), await r);
        e.zstdCompress = async (e, t)=>(await i()).compress(e, t), e.zstdDecompress = async (e)=>(await i()).decompress(e);
        let a = typeof Buffer < `u` && typeof Buffer.from == `function`, o = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/`, s = new Uint8Array(64);
        for(let e = 0; e < 64; e++)s[e] = o.charCodeAt(e);
        let c = new Uint8Array(256).fill(255);
        for(let e = 0; e < 64; e++)c[o.charCodeAt(e)] = e;
        c[45] = 62, c[95] = 63, c[32] = 254, c[9] = 254, c[10] = 254, c[13] = 254, c[61] = 0;
        function l(e) {
            let t = e.length, n = Math.ceil(t / 3) * 4, r = new Uint8Array(n), i = 0, a = 0;
            for(; a + 2 < t; a += 3){
                let t = e[a], n = e[a + 1], o = e[a + 2];
                r[i++] = s[t >> 2], r[i++] = s[(t & 3) << 4 | n >> 4], r[i++] = s[(n & 15) << 2 | o >> 6], r[i++] = s[o & 63];
            }
            if (a < t) {
                let n = e[a];
                if (r[i++] = s[n >> 2], a + 1 < t) {
                    let t = e[a + 1];
                    r[i++] = s[(n & 3) << 4 | t >> 4], r[i++] = s[(t & 15) << 2];
                } else r[i++] = s[(n & 3) << 4], r[i++] = 61;
                r[i++] = 61;
            }
            return r;
        }
        function u(e) {
            let t = 0;
            for(let n = 0; n < e.length; n++){
                let r = c[e[n]];
                if (r < 64) t++;
                else if (r === 255) throw Error(`Invalid base64 character at position ${n}: ${e[n]}`);
            }
            let n = Math.floor(t / 4) * 3 + (t % 4 == 3 ? 2 : +(t % 4 == 2)), r = new Uint8Array(n), i = 0, a = 0, o = 0;
            for(let t = 0; t < e.length && i < n; t++){
                let n = c[e[t]];
                n >= 64 || (a = a << 6 | n, o += 6, o >= 8 && (o -= 8, r[i++] = a >> o & 255));
            }
            return r;
        }
        function d(e) {
            if (a) {
                let t = Buffer.from(e.buffer, e.byteOffset, e.byteLength).toString(`base64`);
                return Buffer.from(t, `utf8`);
            }
            return l(e);
        }
        function f(e) {
            if (a) {
                let t = Buffer.from(e.buffer, e.byteOffset, e.byteLength).toString(`utf8`);
                return Buffer.from(t, `base64`);
            }
            return u(e);
        }
        e.encodeBase64 = (e)=>Promise.resolve(d(e)), e.decodeBase64 = (e)=>Promise.resolve(f(e));
        let p, m = ()=>(p ||= Fe().Duplex, p);
        function h() {
            return typeof CompressionStream < `u` ? new CompressionStream(`gzip`) : m().toWeb(t.createGzip());
        }
        function g() {
            return typeof DecompressionStream < `u` ? new DecompressionStream(`gzip`) : m().toWeb(t.createGunzip());
        }
        function _() {
            return typeof CompressionStream < `u` ? new CompressionStream(`deflate`) : m().toWeb(t.createDeflate());
        }
        function v() {
            return typeof DecompressionStream < `u` ? new DecompressionStream(`deflate`) : m().toWeb(t.createInflate());
        }
        function y() {
            if (typeof CompressionStream < `u`) try {
                return new CompressionStream(`deflate-raw`);
            } catch  {}
            return m().toWeb(t.createDeflateRaw());
        }
        function b() {
            if (typeof DecompressionStream < `u`) try {
                return new DecompressionStream(`deflate-raw`);
            } catch  {}
            return m().toWeb(t.createInflateRaw());
        }
        function x() {
            if (typeof CompressionStream < `u`) try {
                return new CompressionStream(`br`);
            } catch  {}
            return t.createBrotliCompress ? m().toWeb(t.createBrotliCompress()) : w();
        }
        function S() {
            if (typeof DecompressionStream < `u`) try {
                return new DecompressionStream(`br`);
            } catch  {}
            return t.createBrotliDecompress ? m().toWeb(t.createBrotliDecompress()) : T();
        }
        let C = 1024 * 1024;
        function w() {
            let e, t, n = Promise.resolve().then(()=>ht());
            return new TransformStream({
                async start () {
                    e = await n, t = new e.CompressStream;
                },
                transform (n, r) {
                    let i = new Uint8Array(ArrayBuffer.isView(n) ? n.buffer : n, ArrayBuffer.isView(n) ? n.byteOffset : 0, n.byteLength), a = 0;
                    for(; a < i.length;){
                        let n = t.compress(i.subarray(a), C);
                        if (n.buf.length > 0 && r.enqueue(n.buf), a += n.input_offset, n.code === e.BrotliStreamResultCode.NeedsMoreInput) break;
                    }
                },
                flush (n) {
                    for(;;){
                        let r = t.compress(void 0, C);
                        if (r.buf.length > 0 && n.enqueue(r.buf), r.code !== e.BrotliStreamResultCode.NeedsMoreOutput) break;
                    }
                    t.free();
                }
            });
        }
        function T() {
            let e, t, n = Promise.resolve().then(()=>ht());
            return new TransformStream({
                async start () {
                    e = await n, t = new e.DecompressStream;
                },
                transform (n, r) {
                    let i = new Uint8Array(ArrayBuffer.isView(n) ? n.buffer : n, ArrayBuffer.isView(n) ? n.byteOffset : 0, n.byteLength), a = 0;
                    for(; a < i.length;){
                        let n = t.decompress(i.subarray(a), C);
                        if (n.buf.length > 0 && r.enqueue(n.buf), a += n.input_offset, n.code === e.BrotliStreamResultCode.NeedsMoreInput) break;
                    }
                },
                flush () {
                    t.free();
                }
            });
        }
        function E() {
            if (typeof CompressionStream < `u`) try {
                return new CompressionStream(`zstd`);
            } catch  {}
            return t.createZstdCompress ? m().toWeb(t.createZstdCompress()) : A();
        }
        function D() {
            if (typeof DecompressionStream < `u`) try {
                return new DecompressionStream(`zstd`);
            } catch  {}
            return t.createZstdDecompress ? m().toWeb(t.createZstdDecompress()) : j();
        }
        let O, k = ()=>(O ||= new Promise((e)=>{
                _t().run((t)=>{
                    e(t);
                });
            }), O);
        function A() {
            let e, t = k().then((t)=>{
                e = new t.ZstdCompressTransform;
            });
            return new TransformStream({
                async start () {
                    await t;
                },
                transform (t, n) {
                    let r = new Uint8Array(ArrayBuffer.isView(t) ? t.buffer : t, ArrayBuffer.isView(t) ? t.byteOffset : 0, t.byteLength);
                    return new Promise((t, i)=>{
                        let a = (e)=>{
                            n.enqueue(new Uint8Array(e));
                        };
                        e.once(`error`, i), e.on(`data`, a), e.write(r, (n)=>{
                            e.off(`data`, a), e.off(`error`, i), n ? i(n) : t();
                        });
                    });
                },
                flush (t) {
                    return new Promise((n, r)=>{
                        e.once(`error`, r), e.on(`data`, (e)=>{
                            t.enqueue(new Uint8Array(e));
                        }), e.once(`end`, n), e.end();
                    });
                }
            });
        }
        function j() {
            let e, t = k().then((t)=>{
                e = new t.ZstdDecompressTransform;
            });
            return new TransformStream({
                async start () {
                    await t;
                },
                transform (t, n) {
                    let r = new Uint8Array(ArrayBuffer.isView(t) ? t.buffer : t, ArrayBuffer.isView(t) ? t.byteOffset : 0, t.byteLength);
                    return new Promise((t, i)=>{
                        let a = (e)=>{
                            n.enqueue(new Uint8Array(e));
                        };
                        e.once(`error`, i), e.on(`data`, a), e.write(r, (n)=>{
                            e.off(`data`, a), e.off(`error`, i), n ? i(n) : t();
                        });
                    });
                },
                flush (t) {
                    return new Promise((n, r)=>{
                        e.once(`error`, r), e.on(`data`, (e)=>{
                            t.enqueue(new Uint8Array(e));
                        }), e.once(`end`, n), e.end();
                    });
                }
            });
        }
        let M = 1536 * 1024;
        function N() {
            let e = new Uint8Array, t = !0;
            return new TransformStream({
                async transform (n, r) {
                    let i = new Uint8Array(ArrayBuffer.isView(n) ? n.buffer : n, ArrayBuffer.isView(n) ? n.byteOffset : 0, n.byteLength), a = new Uint8Array(e.length + i.length);
                    a.set(e, 0), a.set(i, e.length);
                    let o = 0;
                    for(; o + 3 <= a.length;){
                        let e = Math.min(o + M, a.length), n = o + Math.floor((e - o) / 3) * 3;
                        if (n > o) {
                            let e = d(a.subarray(o, n));
                            r.enqueue(e), o = n, t || await new Promise((e)=>setTimeout(e, 0)), t = !1;
                        } else break;
                    }
                    e = a.subarray(o);
                },
                flush (t) {
                    if (e.length > 0) {
                        let n = d(e);
                        t.enqueue(n);
                    }
                }
            });
        }
        function P() {
            let e = new Uint8Array, t = !0;
            return new TransformStream({
                async transform (n, r) {
                    let i = new Uint8Array(ArrayBuffer.isView(n) ? n.buffer : n, ArrayBuffer.isView(n) ? n.byteOffset : 0, n.byteLength), a = new Uint8Array(e.length + i.length);
                    a.set(e, 0), a.set(i, e.length);
                    let o = 0;
                    for(let e = 0; e < a.length; e++){
                        let t = c[a[e]];
                        if (t < 64 || a[e] === 61) o++;
                        else if (t === 255) throw Error(`Invalid base64 character at position ${e}: ${a[e]}`);
                    }
                    let s = Math.floor(o / 4);
                    if (s > 0) {
                        let n = 0, i = 0;
                        for(let e = 0; e < a.length && n < s * 4; e++)(c[a[e]] < 64 || a[e] === 61) && (n++, i = e + 1);
                        let o = 0;
                        for(; o < i;){
                            let e = Math.min(o + M, i), n = f(a.subarray(o, e));
                            n.length > 0 && r.enqueue(n), o = e, !t && o < i && await new Promise((e)=>setTimeout(e, 0)), t = !1;
                        }
                        e = a.subarray(i);
                    } else e = a;
                },
                flush (t) {
                    if (e.length > 0) {
                        let n = 0;
                        for(let t = 0; t < e.length; t++)c[e[t]] < 64 && n++;
                        if (n > 0) {
                            let n = f(e);
                            n.length > 0 && t.enqueue(n);
                        }
                    }
                }
            });
        }
        function ee(e) {
            if (e.length === 1) return e[0];
            let t = e[0], n = e[e.length - 1];
            for(let t = 0; t < e.length - 1; t++)e[t].readable.pipeTo(e[t + 1].writable);
            return {
                writable: t.writable,
                readable: n.readable
            };
        }
        function F(e) {
            switch(e.toLowerCase()){
                case `gzip`:
                case `x-gzip`:
                    return g();
                case `deflate`:
                case `x-deflate`:
                    return v();
                case `br`:
                    return S();
                case `zstd`:
                    return D();
                case `base64`:
                    return P();
                default:
                    throw Error(`Unsupported encoding: ${e}`);
            }
        }
        function te(e) {
            switch(e.toLowerCase()){
                case `gzip`:
                case `x-gzip`:
                    return h();
                case `deflate`:
                case `x-deflate`:
                    return _();
                case `br`:
                    return x();
                case `zstd`:
                    return E();
                case `base64`:
                    return N();
                default:
                    throw Error(`Unsupported encoding: ${e}`);
            }
        }
        function I(e) {
            return e ? (Array.isArray(e) ? e : e.includes(`, `) ? e.split(`, `) : [
                e
            ]).filter((e)=>!R.includes(e.toLowerCase())) : [];
        }
        function ne(e) {
            let t = I(e);
            return t.length === 0 ? null : (t.reverse(), t.length === 1 ? F(t[0]) : ee(t.map((e)=>F(e))));
        }
        function re(e) {
            let t = I(e);
            return t.length === 0 ? null : t.length === 1 ? te(t[0]) : ee(t.map((e)=>te(e)));
        }
        let L = (e)=>Buffer.isBuffer(e) ? e : e instanceof ArrayBuffer ? Buffer.from(e) : Buffer.from(e.buffer, e.byteOffset, e.byteLength), R = [
            `identity`,
            `amz-1.0`,
            `none`,
            `text`,
            `binary`,
            `utf8`,
            `utf-8`
        ];
        async function z(t, n) {
            let r = L(t);
            if (Array.isArray(n) || typeof n == `string` && n.indexOf(`, `) >= 0) return (typeof n == `string` ? n.split(`, `).reverse() : n).reduce((e, t)=>e.then((e)=>z(e, t)), Promise.resolve(r));
            if (n = n ? n.toLowerCase() : `identity`, n === `gzip` || n === `x-gzip`) return (0, e.gunzip)(r);
            if (n === `deflate` || n === `x-deflate`) return (r[0] & 15) == 8 ? (0, e.inflate)(r) : (0, e.inflateRaw)(r);
            if (n === `br`) return L(await (0, e.brotliDecompress)(r));
            if (n === `zstd`) return L(await (0, e.zstdDecompress)(r));
            if (n === `base64`) return L(f(r));
            if (R.includes(n)) return L(r);
            throw Error(`Unsupported encoding: ${n}`);
        }
        function B(e, n) {
            let r = L(e);
            if (Array.isArray(n) || typeof n == `string` && n.indexOf(`, `) >= 0) return (typeof n == `string` ? n.split(`, `).reverse() : n).reduce((e, t)=>B(e, t), r);
            if (n = n ? n.toLowerCase() : `identity`, n === `gzip` || n === `x-gzip`) return t.gunzipSync(r);
            if (n === `deflate` || n === `x-deflate`) return (r[0] & 15) == 8 ? t.inflateSync(r) : t.inflateRawSync(r);
            if (n === `base64`) return L(f(r));
            if (R.includes(n)) return L(r);
            throw Error(`Unsupported encoding: ${n}`);
        }
        async function ie(t, n, r = {}) {
            let i = L(t), a = r.level ?? 4;
            if (n = n ? n.toLowerCase() : `identity`, n === `gzip` || n === `x-gzip`) return (0, e.gzip)(i, {
                level: a
            });
            if (n === `deflate` || n === `x-deflate`) return (0, e.deflate)(i, {
                level: a
            });
            if (n === `br`) return L(await (0, e.brotliCompress)(i, a));
            if (n === `zstd`) return L(await (0, e.zstdCompress)(i, a));
            if (n === `base64`) return L(d(i));
            if (R.includes(n)) return L(i);
            throw Error(`Unsupported encoding: ${n}`);
        }
    })), yt = t(((e, t)=>{
        function n(e) {
            this.__parent = e, this.__character_count = 0, this.__indent_count = -1, this.__alignment_count = 0, this.__wrap_point_index = 0, this.__wrap_point_character_count = 0, this.__wrap_point_indent_count = -1, this.__wrap_point_alignment_count = 0, this.__items = [];
        }
        n.prototype.clone_empty = function() {
            var e = new n(this.__parent);
            return e.set_indent(this.__indent_count, this.__alignment_count), e;
        }, n.prototype.item = function(e) {
            return e < 0 ? this.__items[this.__items.length + e] : this.__items[e];
        }, n.prototype.has_match = function(e) {
            for(var t = this.__items.length - 1; t >= 0; t--)if (this.__items[t].match(e)) return !0;
            return !1;
        }, n.prototype.set_indent = function(e, t) {
            this.is_empty() && (this.__indent_count = e || 0, this.__alignment_count = t || 0, this.__character_count = this.__parent.get_indent_size(this.__indent_count, this.__alignment_count));
        }, n.prototype._set_wrap_point = function() {
            this.__parent.wrap_line_length && (this.__wrap_point_index = this.__items.length, this.__wrap_point_character_count = this.__character_count, this.__wrap_point_indent_count = this.__parent.next_line.__indent_count, this.__wrap_point_alignment_count = this.__parent.next_line.__alignment_count);
        }, n.prototype._should_wrap = function() {
            return this.__wrap_point_index && this.__character_count > this.__parent.wrap_line_length && this.__wrap_point_character_count > this.__parent.next_line.__character_count;
        }, n.prototype._allow_wrap = function() {
            if (this._should_wrap()) {
                this.__parent.add_new_line();
                var e = this.__parent.current_line;
                return e.set_indent(this.__wrap_point_indent_count, this.__wrap_point_alignment_count), e.__items = this.__items.slice(this.__wrap_point_index), this.__items = this.__items.slice(0, this.__wrap_point_index), e.__character_count += this.__character_count - this.__wrap_point_character_count, this.__character_count = this.__wrap_point_character_count, e.__items[0] === ` ` && (e.__items.splice(0, 1), --e.__character_count), !0;
            }
            return !1;
        }, n.prototype.is_empty = function() {
            return this.__items.length === 0;
        }, n.prototype.last = function() {
            return this.is_empty() ? null : this.__items[this.__items.length - 1];
        }, n.prototype.push = function(e) {
            this.__items.push(e);
            var t = e.lastIndexOf(`
`);
            t === -1 ? this.__character_count += e.length : this.__character_count = e.length - t;
        }, n.prototype.pop = function() {
            var e = null;
            return this.is_empty() || (e = this.__items.pop(), this.__character_count -= e.length), e;
        }, n.prototype._remove_indent = function() {
            this.__indent_count > 0 && (--this.__indent_count, this.__character_count -= this.__parent.indent_size);
        }, n.prototype._remove_wrap_indent = function() {
            this.__wrap_point_indent_count > 0 && --this.__wrap_point_indent_count;
        }, n.prototype.trim = function() {
            for(; this.last() === ` `;)this.__items.pop(), --this.__character_count;
        }, n.prototype.toString = function() {
            var e = ``;
            return this.is_empty() ? this.__parent.indent_empty_lines && (e = this.__parent.get_indent_string(this.__indent_count)) : (e = this.__parent.get_indent_string(this.__indent_count, this.__alignment_count), e += this.__items.join(``)), e;
        };
        function r(e, t) {
            this.__cache = [
                ``
            ], this.__indent_size = e.indent_size, this.__indent_string = e.indent_char, e.indent_with_tabs || (this.__indent_string = Array(e.indent_size + 1).join(e.indent_char)), t ||= ``, e.indent_level > 0 && (t = Array(e.indent_level + 1).join(this.__indent_string)), this.__base_string = t, this.__base_string_length = t.length;
        }
        r.prototype.get_indent_size = function(e, t) {
            var n = this.__base_string_length;
            return t ||= 0, e < 0 && (n = 0), n += e * this.__indent_size, n += t, n;
        }, r.prototype.get_indent_string = function(e, t) {
            var n = this.__base_string;
            return t ||= 0, e < 0 && (e = 0, n = ``), t += e * this.__indent_size, this.__ensure_cache(t), n += this.__cache[t], n;
        }, r.prototype.__ensure_cache = function(e) {
            for(; e >= this.__cache.length;)this.__add_column();
        }, r.prototype.__add_column = function() {
            var e = this.__cache.length, t = 0, n = ``;
            this.__indent_size && e >= this.__indent_size && (t = Math.floor(e / this.__indent_size), e -= t * this.__indent_size, n = Array(t + 1).join(this.__indent_string)), e && (n += Array(e + 1).join(` `)), this.__cache.push(n);
        };
        function i(e, t) {
            this.__indent_cache = new r(e, t), this.raw = !1, this._end_with_newline = e.end_with_newline, this.indent_size = e.indent_size, this.wrap_line_length = e.wrap_line_length, this.indent_empty_lines = e.indent_empty_lines, this.__lines = [], this.previous_line = null, this.current_line = null, this.next_line = new n(this), this.space_before_token = !1, this.non_breaking_space = !1, this.previous_token_wrapped = !1, this.__add_outputline();
        }
        i.prototype.__add_outputline = function() {
            this.previous_line = this.current_line, this.current_line = this.next_line.clone_empty(), this.__lines.push(this.current_line);
        }, i.prototype.get_line_number = function() {
            return this.__lines.length;
        }, i.prototype.get_indent_string = function(e, t) {
            return this.__indent_cache.get_indent_string(e, t);
        }, i.prototype.get_indent_size = function(e, t) {
            return this.__indent_cache.get_indent_size(e, t);
        }, i.prototype.is_empty = function() {
            return !this.previous_line && this.current_line.is_empty();
        }, i.prototype.add_new_line = function(e) {
            return this.is_empty() || !e && this.just_added_newline() ? !1 : (this.raw || this.__add_outputline(), !0);
        }, i.prototype.get_code = function(e) {
            this.trim(!0);
            var t = this.current_line.pop();
            t && (t[t.length - 1] === `
` && (t = t.replace(/\n+$/g, ``)), this.current_line.push(t)), this._end_with_newline && this.__add_outputline();
            var n = this.__lines.join(`
`);
            return e !== `
` && (n = n.replace(/[\n]/g, e)), n;
        }, i.prototype.set_wrap_point = function() {
            this.current_line._set_wrap_point();
        }, i.prototype.set_indent = function(e, t) {
            return e ||= 0, t ||= 0, this.next_line.set_indent(e, t), this.__lines.length > 1 ? (this.current_line.set_indent(e, t), !0) : (this.current_line.set_indent(), !1);
        }, i.prototype.add_raw_token = function(e) {
            for(var t = 0; t < e.newlines; t++)this.__add_outputline();
            this.current_line.set_indent(-1), this.current_line.push(e.whitespace_before), this.current_line.push(e.text), this.space_before_token = !1, this.non_breaking_space = !1, this.previous_token_wrapped = !1;
        }, i.prototype.add_token = function(e) {
            this.__add_space_before_token(), this.current_line.push(e), this.space_before_token = !1, this.non_breaking_space = !1, this.previous_token_wrapped = this.current_line._allow_wrap();
        }, i.prototype.__add_space_before_token = function() {
            this.space_before_token && !this.just_added_newline() && (this.non_breaking_space || this.set_wrap_point(), this.current_line.push(` `));
        }, i.prototype.remove_indent = function(e) {
            for(var t = this.__lines.length; e < t;)this.__lines[e]._remove_indent(), e++;
            this.current_line._remove_wrap_indent();
        }, i.prototype.trim = function(e) {
            for(e = e === void 0 ? !1 : e, this.current_line.trim(); e && this.__lines.length > 1 && this.current_line.is_empty();)this.__lines.pop(), this.current_line = this.__lines[this.__lines.length - 1], this.current_line.trim();
            this.previous_line = this.__lines.length > 1 ? this.__lines[this.__lines.length - 2] : null;
        }, i.prototype.just_added_newline = function() {
            return this.current_line.is_empty();
        }, i.prototype.just_added_blankline = function() {
            return this.is_empty() || this.current_line.is_empty() && this.previous_line.is_empty();
        }, i.prototype.ensure_empty_line_above = function(e, t) {
            for(var r = this.__lines.length - 2; r >= 0;){
                var i = this.__lines[r];
                if (i.is_empty()) break;
                if (i.item(0).indexOf(e) !== 0 && i.item(-1) !== t) {
                    this.__lines.splice(r + 1, 0, new n(this)), this.previous_line = this.__lines[this.__lines.length - 2];
                    break;
                }
                r--;
            }
        }, t.exports.Output = i;
    })), bt = t(((e, t)=>{
        function n(e, t, n, r) {
            this.type = e, this.text = t, this.comments_before = null, this.newlines = n || 0, this.whitespace_before = r || ``, this.parent = null, this.next = null, this.previous = null, this.opened = null, this.closed = null, this.directives = null;
        }
        t.exports.Token = n;
    })), xt = t(((e)=>{
        var t = `\\x23\\x24\\x40\\x41-\\x5a\\x5f\\x61-\\x7a`, n = `\\x24\\x30-\\x39\\x41-\\x5a\\x5f\\x61-\\x7a`, r = `\\xaa\\xb5\\xba\\xc0-\\xd6\\xd8-\\xf6\\xf8-\\u02c1\\u02c6-\\u02d1\\u02e0-\\u02e4\\u02ec\\u02ee\\u0370-\\u0374\\u0376\\u0377\\u037a-\\u037d\\u0386\\u0388-\\u038a\\u038c\\u038e-\\u03a1\\u03a3-\\u03f5\\u03f7-\\u0481\\u048a-\\u0527\\u0531-\\u0556\\u0559\\u0561-\\u0587\\u05d0-\\u05ea\\u05f0-\\u05f2\\u0620-\\u064a\\u066e\\u066f\\u0671-\\u06d3\\u06d5\\u06e5\\u06e6\\u06ee\\u06ef\\u06fa-\\u06fc\\u06ff\\u0710\\u0712-\\u072f\\u074d-\\u07a5\\u07b1\\u07ca-\\u07ea\\u07f4\\u07f5\\u07fa\\u0800-\\u0815\\u081a\\u0824\\u0828\\u0840-\\u0858\\u08a0\\u08a2-\\u08ac\\u0904-\\u0939\\u093d\\u0950\\u0958-\\u0961\\u0971-\\u0977\\u0979-\\u097f\\u0985-\\u098c\\u098f\\u0990\\u0993-\\u09a8\\u09aa-\\u09b0\\u09b2\\u09b6-\\u09b9\\u09bd\\u09ce\\u09dc\\u09dd\\u09df-\\u09e1\\u09f0\\u09f1\\u0a05-\\u0a0a\\u0a0f\\u0a10\\u0a13-\\u0a28\\u0a2a-\\u0a30\\u0a32\\u0a33\\u0a35\\u0a36\\u0a38\\u0a39\\u0a59-\\u0a5c\\u0a5e\\u0a72-\\u0a74\\u0a85-\\u0a8d\\u0a8f-\\u0a91\\u0a93-\\u0aa8\\u0aaa-\\u0ab0\\u0ab2\\u0ab3\\u0ab5-\\u0ab9\\u0abd\\u0ad0\\u0ae0\\u0ae1\\u0b05-\\u0b0c\\u0b0f\\u0b10\\u0b13-\\u0b28\\u0b2a-\\u0b30\\u0b32\\u0b33\\u0b35-\\u0b39\\u0b3d\\u0b5c\\u0b5d\\u0b5f-\\u0b61\\u0b71\\u0b83\\u0b85-\\u0b8a\\u0b8e-\\u0b90\\u0b92-\\u0b95\\u0b99\\u0b9a\\u0b9c\\u0b9e\\u0b9f\\u0ba3\\u0ba4\\u0ba8-\\u0baa\\u0bae-\\u0bb9\\u0bd0\\u0c05-\\u0c0c\\u0c0e-\\u0c10\\u0c12-\\u0c28\\u0c2a-\\u0c33\\u0c35-\\u0c39\\u0c3d\\u0c58\\u0c59\\u0c60\\u0c61\\u0c85-\\u0c8c\\u0c8e-\\u0c90\\u0c92-\\u0ca8\\u0caa-\\u0cb3\\u0cb5-\\u0cb9\\u0cbd\\u0cde\\u0ce0\\u0ce1\\u0cf1\\u0cf2\\u0d05-\\u0d0c\\u0d0e-\\u0d10\\u0d12-\\u0d3a\\u0d3d\\u0d4e\\u0d60\\u0d61\\u0d7a-\\u0d7f\\u0d85-\\u0d96\\u0d9a-\\u0db1\\u0db3-\\u0dbb\\u0dbd\\u0dc0-\\u0dc6\\u0e01-\\u0e30\\u0e32\\u0e33\\u0e40-\\u0e46\\u0e81\\u0e82\\u0e84\\u0e87\\u0e88\\u0e8a\\u0e8d\\u0e94-\\u0e97\\u0e99-\\u0e9f\\u0ea1-\\u0ea3\\u0ea5\\u0ea7\\u0eaa\\u0eab\\u0ead-\\u0eb0\\u0eb2\\u0eb3\\u0ebd\\u0ec0-\\u0ec4\\u0ec6\\u0edc-\\u0edf\\u0f00\\u0f40-\\u0f47\\u0f49-\\u0f6c\\u0f88-\\u0f8c\\u1000-\\u102a\\u103f\\u1050-\\u1055\\u105a-\\u105d\\u1061\\u1065\\u1066\\u106e-\\u1070\\u1075-\\u1081\\u108e\\u10a0-\\u10c5\\u10c7\\u10cd\\u10d0-\\u10fa\\u10fc-\\u1248\\u124a-\\u124d\\u1250-\\u1256\\u1258\\u125a-\\u125d\\u1260-\\u1288\\u128a-\\u128d\\u1290-\\u12b0\\u12b2-\\u12b5\\u12b8-\\u12be\\u12c0\\u12c2-\\u12c5\\u12c8-\\u12d6\\u12d8-\\u1310\\u1312-\\u1315\\u1318-\\u135a\\u1380-\\u138f\\u13a0-\\u13f4\\u1401-\\u166c\\u166f-\\u167f\\u1681-\\u169a\\u16a0-\\u16ea\\u16ee-\\u16f0\\u1700-\\u170c\\u170e-\\u1711\\u1720-\\u1731\\u1740-\\u1751\\u1760-\\u176c\\u176e-\\u1770\\u1780-\\u17b3\\u17d7\\u17dc\\u1820-\\u1877\\u1880-\\u18a8\\u18aa\\u18b0-\\u18f5\\u1900-\\u191c\\u1950-\\u196d\\u1970-\\u1974\\u1980-\\u19ab\\u19c1-\\u19c7\\u1a00-\\u1a16\\u1a20-\\u1a54\\u1aa7\\u1b05-\\u1b33\\u1b45-\\u1b4b\\u1b83-\\u1ba0\\u1bae\\u1baf\\u1bba-\\u1be5\\u1c00-\\u1c23\\u1c4d-\\u1c4f\\u1c5a-\\u1c7d\\u1ce9-\\u1cec\\u1cee-\\u1cf1\\u1cf5\\u1cf6\\u1d00-\\u1dbf\\u1e00-\\u1f15\\u1f18-\\u1f1d\\u1f20-\\u1f45\\u1f48-\\u1f4d\\u1f50-\\u1f57\\u1f59\\u1f5b\\u1f5d\\u1f5f-\\u1f7d\\u1f80-\\u1fb4\\u1fb6-\\u1fbc\\u1fbe\\u1fc2-\\u1fc4\\u1fc6-\\u1fcc\\u1fd0-\\u1fd3\\u1fd6-\\u1fdb\\u1fe0-\\u1fec\\u1ff2-\\u1ff4\\u1ff6-\\u1ffc\\u2071\\u207f\\u2090-\\u209c\\u2102\\u2107\\u210a-\\u2113\\u2115\\u2119-\\u211d\\u2124\\u2126\\u2128\\u212a-\\u212d\\u212f-\\u2139\\u213c-\\u213f\\u2145-\\u2149\\u214e\\u2160-\\u2188\\u2c00-\\u2c2e\\u2c30-\\u2c5e\\u2c60-\\u2ce4\\u2ceb-\\u2cee\\u2cf2\\u2cf3\\u2d00-\\u2d25\\u2d27\\u2d2d\\u2d30-\\u2d67\\u2d6f\\u2d80-\\u2d96\\u2da0-\\u2da6\\u2da8-\\u2dae\\u2db0-\\u2db6\\u2db8-\\u2dbe\\u2dc0-\\u2dc6\\u2dc8-\\u2dce\\u2dd0-\\u2dd6\\u2dd8-\\u2dde\\u2e2f\\u3005-\\u3007\\u3021-\\u3029\\u3031-\\u3035\\u3038-\\u303c\\u3041-\\u3096\\u309d-\\u309f\\u30a1-\\u30fa\\u30fc-\\u30ff\\u3105-\\u312d\\u3131-\\u318e\\u31a0-\\u31ba\\u31f0-\\u31ff\\u3400-\\u4db5\\u4e00-\\u9fcc\\ua000-\\ua48c\\ua4d0-\\ua4fd\\ua500-\\ua60c\\ua610-\\ua61f\\ua62a\\ua62b\\ua640-\\ua66e\\ua67f-\\ua697\\ua6a0-\\ua6ef\\ua717-\\ua71f\\ua722-\\ua788\\ua78b-\\ua78e\\ua790-\\ua793\\ua7a0-\\ua7aa\\ua7f8-\\ua801\\ua803-\\ua805\\ua807-\\ua80a\\ua80c-\\ua822\\ua840-\\ua873\\ua882-\\ua8b3\\ua8f2-\\ua8f7\\ua8fb\\ua90a-\\ua925\\ua930-\\ua946\\ua960-\\ua97c\\ua984-\\ua9b2\\ua9cf\\uaa00-\\uaa28\\uaa40-\\uaa42\\uaa44-\\uaa4b\\uaa60-\\uaa76\\uaa7a\\uaa80-\\uaaaf\\uaab1\\uaab5\\uaab6\\uaab9-\\uaabd\\uaac0\\uaac2\\uaadb-\\uaadd\\uaae0-\\uaaea\\uaaf2-\\uaaf4\\uab01-\\uab06\\uab09-\\uab0e\\uab11-\\uab16\\uab20-\\uab26\\uab28-\\uab2e\\uabc0-\\uabe2\\uac00-\\ud7a3\\ud7b0-\\ud7c6\\ud7cb-\\ud7fb\\uf900-\\ufa6d\\ufa70-\\ufad9\\ufb00-\\ufb06\\ufb13-\\ufb17\\ufb1d\\ufb1f-\\ufb28\\ufb2a-\\ufb36\\ufb38-\\ufb3c\\ufb3e\\ufb40\\ufb41\\ufb43\\ufb44\\ufb46-\\ufbb1\\ufbd3-\\ufd3d\\ufd50-\\ufd8f\\ufd92-\\ufdc7\\ufdf0-\\ufdfb\\ufe70-\\ufe74\\ufe76-\\ufefc\\uff21-\\uff3a\\uff41-\\uff5a\\uff66-\\uffbe\\uffc2-\\uffc7\\uffca-\\uffcf\\uffd2-\\uffd7\\uffda-\\uffdc`, i = `\\u0300-\\u036f\\u0483-\\u0487\\u0591-\\u05bd\\u05bf\\u05c1\\u05c2\\u05c4\\u05c5\\u05c7\\u0610-\\u061a\\u0620-\\u0649\\u0672-\\u06d3\\u06e7-\\u06e8\\u06fb-\\u06fc\\u0730-\\u074a\\u0800-\\u0814\\u081b-\\u0823\\u0825-\\u0827\\u0829-\\u082d\\u0840-\\u0857\\u08e4-\\u08fe\\u0900-\\u0903\\u093a-\\u093c\\u093e-\\u094f\\u0951-\\u0957\\u0962-\\u0963\\u0966-\\u096f\\u0981-\\u0983\\u09bc\\u09be-\\u09c4\\u09c7\\u09c8\\u09d7\\u09df-\\u09e0\\u0a01-\\u0a03\\u0a3c\\u0a3e-\\u0a42\\u0a47\\u0a48\\u0a4b-\\u0a4d\\u0a51\\u0a66-\\u0a71\\u0a75\\u0a81-\\u0a83\\u0abc\\u0abe-\\u0ac5\\u0ac7-\\u0ac9\\u0acb-\\u0acd\\u0ae2-\\u0ae3\\u0ae6-\\u0aef\\u0b01-\\u0b03\\u0b3c\\u0b3e-\\u0b44\\u0b47\\u0b48\\u0b4b-\\u0b4d\\u0b56\\u0b57\\u0b5f-\\u0b60\\u0b66-\\u0b6f\\u0b82\\u0bbe-\\u0bc2\\u0bc6-\\u0bc8\\u0bca-\\u0bcd\\u0bd7\\u0be6-\\u0bef\\u0c01-\\u0c03\\u0c46-\\u0c48\\u0c4a-\\u0c4d\\u0c55\\u0c56\\u0c62-\\u0c63\\u0c66-\\u0c6f\\u0c82\\u0c83\\u0cbc\\u0cbe-\\u0cc4\\u0cc6-\\u0cc8\\u0cca-\\u0ccd\\u0cd5\\u0cd6\\u0ce2-\\u0ce3\\u0ce6-\\u0cef\\u0d02\\u0d03\\u0d46-\\u0d48\\u0d57\\u0d62-\\u0d63\\u0d66-\\u0d6f\\u0d82\\u0d83\\u0dca\\u0dcf-\\u0dd4\\u0dd6\\u0dd8-\\u0ddf\\u0df2\\u0df3\\u0e34-\\u0e3a\\u0e40-\\u0e45\\u0e50-\\u0e59\\u0eb4-\\u0eb9\\u0ec8-\\u0ecd\\u0ed0-\\u0ed9\\u0f18\\u0f19\\u0f20-\\u0f29\\u0f35\\u0f37\\u0f39\\u0f41-\\u0f47\\u0f71-\\u0f84\\u0f86-\\u0f87\\u0f8d-\\u0f97\\u0f99-\\u0fbc\\u0fc6\\u1000-\\u1029\\u1040-\\u1049\\u1067-\\u106d\\u1071-\\u1074\\u1082-\\u108d\\u108f-\\u109d\\u135d-\\u135f\\u170e-\\u1710\\u1720-\\u1730\\u1740-\\u1750\\u1772\\u1773\\u1780-\\u17b2\\u17dd\\u17e0-\\u17e9\\u180b-\\u180d\\u1810-\\u1819\\u1920-\\u192b\\u1930-\\u193b\\u1951-\\u196d\\u19b0-\\u19c0\\u19c8-\\u19c9\\u19d0-\\u19d9\\u1a00-\\u1a15\\u1a20-\\u1a53\\u1a60-\\u1a7c\\u1a7f-\\u1a89\\u1a90-\\u1a99\\u1b46-\\u1b4b\\u1b50-\\u1b59\\u1b6b-\\u1b73\\u1bb0-\\u1bb9\\u1be6-\\u1bf3\\u1c00-\\u1c22\\u1c40-\\u1c49\\u1c5b-\\u1c7d\\u1cd0-\\u1cd2\\u1d00-\\u1dbe\\u1e01-\\u1f15\\u200c\\u200d\\u203f\\u2040\\u2054\\u20d0-\\u20dc\\u20e1\\u20e5-\\u20f0\\u2d81-\\u2d96\\u2de0-\\u2dff\\u3021-\\u3028\\u3099\\u309a\\ua640-\\ua66d\\ua674-\\ua67d\\ua69f\\ua6f0-\\ua6f1\\ua7f8-\\ua800\\ua806\\ua80b\\ua823-\\ua827\\ua880-\\ua881\\ua8b4-\\ua8c4\\ua8d0-\\ua8d9\\ua8f3-\\ua8f7\\ua900-\\ua909\\ua926-\\ua92d\\ua930-\\ua945\\ua980-\\ua983\\ua9b3-\\ua9c0\\uaa00-\\uaa27\\uaa40-\\uaa41\\uaa4c-\\uaa4d\\uaa50-\\uaa59\\uaa7b\\uaae0-\\uaae9\\uaaf2-\\uaaf3\\uabc0-\\uabe1\\uabec\\uabed\\uabf0-\\uabf9\\ufb20-\\ufb28\\ufe00-\\ufe0f\\ufe20-\\ufe26\\ufe33\\ufe34\\ufe4d-\\ufe4f\\uff10-\\uff19\\uff3f`, a = `\\\\u[0-9a-fA-F]{4}|\\\\u\\{[0-9a-fA-F]+\\}`, o = `(?:` + a + `|[` + t + r + `])`, s = `(?:` + a + `|[` + n + r + i + `])*`;
        e.identifier = new RegExp(o + s, `g`), e.identifierStart = new RegExp(o), e.identifierMatch = RegExp(`(?:` + a + `|[` + n + r + i + `])+`), e.newline = /[\n\r\u2028\u2029]/, e.lineBreak = RegExp(`\r
|` + e.newline.source), e.allLineBreaks = new RegExp(e.lineBreak.source, `g`);
    })), St = t(((e, t)=>{
        function n(e, t) {
            this.raw_options = r(e, t), this.disabled = this._get_boolean(`disabled`), this.eol = this._get_characters(`eol`, `auto`), this.end_with_newline = this._get_boolean(`end_with_newline`), this.indent_size = this._get_number(`indent_size`, 4), this.indent_char = this._get_characters(`indent_char`, ` `), this.indent_level = this._get_number(`indent_level`), this.preserve_newlines = this._get_boolean(`preserve_newlines`, !0), this.max_preserve_newlines = this._get_number(`max_preserve_newlines`, 32786), this.preserve_newlines || (this.max_preserve_newlines = 0), this.indent_with_tabs = this._get_boolean(`indent_with_tabs`, this.indent_char === `	`), this.indent_with_tabs && (this.indent_char = `	`, this.indent_size === 1 && (this.indent_size = 4)), this.wrap_line_length = this._get_number(`wrap_line_length`, this._get_number(`max_char`)), this.indent_empty_lines = this._get_boolean(`indent_empty_lines`), this.templating = this._get_selection_list(`templating`, [
                `auto`,
                `none`,
                `angular`,
                `django`,
                `erb`,
                `handlebars`,
                `php`,
                `smarty`
            ], [
                `auto`
            ]);
        }
        n.prototype._get_array = function(e, t) {
            var n = this.raw_options[e], r = t || [];
            return typeof n == `object` ? n !== null && typeof n.concat == `function` && (r = n.concat()) : typeof n == `string` && (r = n.split(/[^a-zA-Z0-9_\/\-]+/)), r;
        }, n.prototype._get_boolean = function(e, t) {
            var n = this.raw_options[e];
            return n === void 0 ? !!t : !!n;
        }, n.prototype._get_characters = function(e, t) {
            var n = this.raw_options[e], r = t || ``;
            return typeof n == `string` && (r = n.replace(/\\r/, `\r`).replace(/\\n/, `
`).replace(/\\t/, `	`)), r;
        }, n.prototype._get_number = function(e, t) {
            var n = this.raw_options[e];
            t = parseInt(t, 10), isNaN(t) && (t = 0);
            var r = parseInt(n, 10);
            return isNaN(r) && (r = t), r;
        }, n.prototype._get_selection = function(e, t, n) {
            var r = this._get_selection_list(e, t, n);
            if (r.length !== 1) throw Error(`Invalid Option Value: The option '` + e + `' can only be one of the following values:
` + t + `
You passed in: '` + this.raw_options[e] + `'`);
            return r[0];
        }, n.prototype._get_selection_list = function(e, t, n) {
            if (!t || t.length === 0) throw Error(`Selection list cannot be empty.`);
            if (n ||= [
                t[0]
            ], !this._is_valid_selection(n, t)) throw Error(`Invalid Default Value!`);
            var r = this._get_array(e, n);
            if (!this._is_valid_selection(r, t)) throw Error(`Invalid Option Value: The option '` + e + `' can contain only the following values:
` + t + `
You passed in: '` + this.raw_options[e] + `'`);
            return r;
        }, n.prototype._is_valid_selection = function(e, t) {
            return e.length && t.length && !e.some(function(e) {
                return t.indexOf(e) === -1;
            });
        };
        function r(e, t) {
            var n = {};
            for(var r in e = i(e), e)r !== t && (n[r] = e[r]);
            if (t && e[t]) for(r in e[t])n[r] = e[t][r];
            return n;
        }
        function i(e) {
            var t = {}, n;
            for(n in e){
                var r = n.replace(/-/g, `_`);
                t[r] = e[n];
            }
            return t;
        }
        t.exports.Options = n, t.exports.normalizeOpts = i, t.exports.mergeOpts = r;
    })), Ct = t(((e, t)=>{
        var n = St().Options, r = [
            `before-newline`,
            `after-newline`,
            `preserve-newline`
        ];
        function i(e) {
            n.call(this, e, `js`);
            var t = this.raw_options.brace_style || null;
            t === `expand-strict` ? this.raw_options.brace_style = `expand` : t === `collapse-preserve-inline` ? this.raw_options.brace_style = `collapse,preserve-inline` : this.raw_options.braces_on_own_line !== void 0 && (this.raw_options.brace_style = this.raw_options.braces_on_own_line ? `expand` : `collapse`);
            var i = this._get_selection_list(`brace_style`, [
                `collapse`,
                `expand`,
                `end-expand`,
                `none`,
                `preserve-inline`
            ]);
            this.brace_preserve_inline = !1, this.brace_style = `collapse`;
            for(var a = 0; a < i.length; a++)i[a] === `preserve-inline` ? this.brace_preserve_inline = !0 : this.brace_style = i[a];
            this.unindent_chained_methods = this._get_boolean(`unindent_chained_methods`), this.break_chained_methods = this._get_boolean(`break_chained_methods`), this.space_in_paren = this._get_boolean(`space_in_paren`), this.space_in_empty_paren = this._get_boolean(`space_in_empty_paren`), this.jslint_happy = this._get_boolean(`jslint_happy`), this.space_after_anon_function = this._get_boolean(`space_after_anon_function`), this.space_after_named_function = this._get_boolean(`space_after_named_function`), this.keep_array_indentation = this._get_boolean(`keep_array_indentation`), this.space_before_conditional = this._get_boolean(`space_before_conditional`, !0), this.unescape_strings = this._get_boolean(`unescape_strings`), this.e4x = this._get_boolean(`e4x`), this.comma_first = this._get_boolean(`comma_first`), this.operator_position = this._get_selection(`operator_position`, r), this.test_output_raw = this._get_boolean(`test_output_raw`), this.jslint_happy && (this.space_after_anon_function = !0);
        }
        i.prototype = new n, t.exports.Options = i;
    })), wt = t(((e, t)=>{
        var n = RegExp.prototype.hasOwnProperty(`sticky`);
        function r(e) {
            this.__input = e || ``, this.__input_length = this.__input.length, this.__position = 0;
        }
        r.prototype.restart = function() {
            this.__position = 0;
        }, r.prototype.back = function() {
            this.__position > 0 && --this.__position;
        }, r.prototype.hasNext = function() {
            return this.__position < this.__input_length;
        }, r.prototype.next = function() {
            var e = null;
            return this.hasNext() && (e = this.__input.charAt(this.__position), this.__position += 1), e;
        }, r.prototype.peek = function(e) {
            var t = null;
            return e ||= 0, e += this.__position, e >= 0 && e < this.__input_length && (t = this.__input.charAt(e)), t;
        }, r.prototype.__match = function(e, t) {
            e.lastIndex = t;
            var r = e.exec(this.__input);
            return r && !(n && e.sticky) && r.index !== t && (r = null), r;
        }, r.prototype.test = function(e, t) {
            return t ||= 0, t += this.__position, t >= 0 && t < this.__input_length ? !!this.__match(e, t) : !1;
        }, r.prototype.testChar = function(e, t) {
            var n = this.peek(t);
            return e.lastIndex = 0, n !== null && e.test(n);
        }, r.prototype.match = function(e) {
            var t = this.__match(e, this.__position);
            return t ? this.__position += t[0].length : t = null, t;
        }, r.prototype.read = function(e, t, n) {
            var r = ``, i;
            return e && (i = this.match(e), i && (r += i[0])), t && (i || !e) && (r += this.readUntil(t, n)), r;
        }, r.prototype.readUntil = function(e, t) {
            var n = ``, r = this.__position;
            e.lastIndex = this.__position;
            var i = e.exec(this.__input);
            return i ? (r = i.index, t && (r += i[0].length)) : r = this.__input_length, n = this.__input.substring(this.__position, r), this.__position = r, n;
        }, r.prototype.readUntilAfter = function(e) {
            return this.readUntil(e, !0);
        }, r.prototype.get_regexp = function(e, t) {
            var r = null, i = `g`;
            return t && n && (i = `y`), typeof e == `string` && e !== `` ? r = new RegExp(e, i) : e && (r = new RegExp(e.source, i)), r;
        }, r.prototype.get_literal_regexp = function(e) {
            return RegExp(e.replace(/[-\/\\^$*+?.()|[\]{}]/g, `\\$&`));
        }, r.prototype.peekUntilAfter = function(e) {
            var t = this.__position, n = this.readUntilAfter(e);
            return this.__position = t, n;
        }, r.prototype.lookBack = function(e) {
            var t = this.__position - 1;
            return t >= e.length && this.__input.substring(t - e.length, t).toLowerCase() === e;
        }, t.exports.InputScanner = r;
    })), Tt = t(((e, t)=>{
        function n(e) {
            this.__tokens = [], this.__tokens_length = this.__tokens.length, this.__position = 0, this.__parent_token = e;
        }
        n.prototype.restart = function() {
            this.__position = 0;
        }, n.prototype.isEmpty = function() {
            return this.__tokens_length === 0;
        }, n.prototype.hasNext = function() {
            return this.__position < this.__tokens_length;
        }, n.prototype.next = function() {
            var e = null;
            return this.hasNext() && (e = this.__tokens[this.__position], this.__position += 1), e;
        }, n.prototype.peek = function(e) {
            var t = null;
            return e ||= 0, e += this.__position, e >= 0 && e < this.__tokens_length && (t = this.__tokens[e]), t;
        }, n.prototype.add = function(e) {
            this.__parent_token && (e.parent = this.__parent_token), this.__tokens.push(e), this.__tokens_length += 1;
        }, t.exports.TokenStream = n;
    })), Et = t(((e, t)=>{
        function n(e, t) {
            this._input = e, this._starting_pattern = null, this._match_pattern = null, this._until_pattern = null, this._until_after = !1, t && (this._starting_pattern = this._input.get_regexp(t._starting_pattern, !0), this._match_pattern = this._input.get_regexp(t._match_pattern, !0), this._until_pattern = this._input.get_regexp(t._until_pattern), this._until_after = t._until_after);
        }
        n.prototype.read = function() {
            var e = this._input.read(this._starting_pattern);
            return (!this._starting_pattern || e) && (e += this._input.read(this._match_pattern, this._until_pattern, this._until_after)), e;
        }, n.prototype.read_match = function() {
            return this._input.match(this._match_pattern);
        }, n.prototype.until_after = function(e) {
            var t = this._create();
            return t._until_after = !0, t._until_pattern = this._input.get_regexp(e), t._update(), t;
        }, n.prototype.until = function(e) {
            var t = this._create();
            return t._until_after = !1, t._until_pattern = this._input.get_regexp(e), t._update(), t;
        }, n.prototype.starting_with = function(e) {
            var t = this._create();
            return t._starting_pattern = this._input.get_regexp(e, !0), t._update(), t;
        }, n.prototype.matching = function(e) {
            var t = this._create();
            return t._match_pattern = this._input.get_regexp(e, !0), t._update(), t;
        }, n.prototype._create = function() {
            return new n(this._input, this);
        }, n.prototype._update = function() {}, t.exports.Pattern = n;
    })), Dt = t(((e, t)=>{
        var n = Et().Pattern;
        function r(e, t) {
            n.call(this, e, t), t ? this._line_regexp = this._input.get_regexp(t._line_regexp) : this.__set_whitespace_patterns(``, ``), this.newline_count = 0, this.whitespace_before_token = ``;
        }
        r.prototype = new n, r.prototype.__set_whitespace_patterns = function(e, t) {
            e += `\\t `, t += `\\n\\r`, this._match_pattern = this._input.get_regexp(`[` + e + t + `]+`, !0), this._newline_regexp = this._input.get_regexp(`\\r\\n|[` + t + `]`);
        }, r.prototype.read = function() {
            this.newline_count = 0, this.whitespace_before_token = ``;
            var e = this._input.read(this._match_pattern);
            if (e === ` `) this.whitespace_before_token = ` `;
            else if (e) {
                var t = this.__split(this._newline_regexp, e);
                this.newline_count = t.length - 1, this.whitespace_before_token = t[this.newline_count];
            }
            return e;
        }, r.prototype.matching = function(e, t) {
            var n = this._create();
            return n.__set_whitespace_patterns(e, t), n._update(), n;
        }, r.prototype._create = function() {
            return new r(this._input, this);
        }, r.prototype.__split = function(e, t) {
            e.lastIndex = 0;
            for(var n = 0, r = [], i = e.exec(t); i;)r.push(t.substring(n, i.index)), n = i.index + i[0].length, i = e.exec(t);
            return n < t.length ? r.push(t.substring(n, t.length)) : r.push(``), r;
        }, t.exports.WhitespacePattern = r;
    })), Ot = t(((e, t)=>{
        var n = wt().InputScanner, r = bt().Token, i = Tt().TokenStream, a = Dt().WhitespacePattern, o = {
            START: `TK_START`,
            RAW: `TK_RAW`,
            EOF: `TK_EOF`
        }, s = function(e, t) {
            this._input = new n(e), this._options = t || {}, this.__tokens = null, this._patterns = {}, this._patterns.whitespace = new a(this._input);
        };
        s.prototype.tokenize = function() {
            this._input.restart(), this.__tokens = new i, this._reset();
            for(var e, t = new r(o.START, ``), n = null, a = [], s = new i; t.type !== o.EOF;){
                for(e = this._get_next_token(t, n); this._is_comment(e);)s.add(e), e = this._get_next_token(t, n);
                s.isEmpty() || (e.comments_before = s, s = new i), e.parent = n, this._is_opening(e) ? (a.push(n), n = e) : n && this._is_closing(e, n) && (e.opened = n, n.closed = e, n = a.pop(), e.parent = n), e.previous = t, t.next = e, this.__tokens.add(e), t = e;
            }
            return this.__tokens;
        }, s.prototype._is_first_token = function() {
            return this.__tokens.isEmpty();
        }, s.prototype._reset = function() {}, s.prototype._get_next_token = function(e, t) {
            this._readWhitespace();
            var n = this._input.read(/.+/g);
            return n ? this._create_token(o.RAW, n) : this._create_token(o.EOF, ``);
        }, s.prototype._is_comment = function(e) {
            return !1;
        }, s.prototype._is_opening = function(e) {
            return !1;
        }, s.prototype._is_closing = function(e, t) {
            return !1;
        }, s.prototype._create_token = function(e, t) {
            return new r(e, t, this._patterns.whitespace.newline_count, this._patterns.whitespace.whitespace_before_token);
        }, s.prototype._readWhitespace = function() {
            return this._patterns.whitespace.read();
        }, t.exports.Tokenizer = s, t.exports.TOKEN = o;
    })), kt = t(((e, t)=>{
        function n(e, t) {
            e = typeof e == `string` ? e : e.source, t = typeof t == `string` ? t : t.source, this.__directives_block_pattern = RegExp(e + ` beautify( \\w+[:]\\w+)+ ` + t, `g`), this.__directive_pattern = / (\w+)[:](\w+)/g, this.__directives_end_ignore_pattern = RegExp(e + `\\sbeautify\\signore:end\\s` + t, `g`);
        }
        n.prototype.get_directives = function(e) {
            if (!e.match(this.__directives_block_pattern)) return null;
            var t = {};
            this.__directive_pattern.lastIndex = 0;
            for(var n = this.__directive_pattern.exec(e); n;)t[n[1]] = n[2], n = this.__directive_pattern.exec(e);
            return t;
        }, n.prototype.readIgnored = function(e) {
            return e.readUntilAfter(this.__directives_end_ignore_pattern);
        }, t.exports.Directives = n;
    })), At = t(((e, t)=>{
        var n = Et().Pattern, r = {
            django: !1,
            erb: !1,
            handlebars: !1,
            php: !1,
            smarty: !1,
            angular: !1
        };
        function i(e, t) {
            n.call(this, e, t), this.__template_pattern = null, this._disabled = Object.assign({}, r), this._excluded = Object.assign({}, r), t && (this.__template_pattern = this._input.get_regexp(t.__template_pattern), this._excluded = Object.assign(this._excluded, t._excluded), this._disabled = Object.assign(this._disabled, t._disabled));
            var i = new n(e);
            this.__patterns = {
                handlebars_comment: i.starting_with(/{{!--/).until_after(/--}}/),
                handlebars_unescaped: i.starting_with(/{{{/).until_after(/}}}/),
                handlebars: i.starting_with(/{{/).until_after(/}}/),
                php: i.starting_with(/<\?(?:[= ]|php)/).until_after(/\?>/),
                erb: i.starting_with(/<%[^%]/).until_after(/[^%]%>/),
                django: i.starting_with(/{%/).until_after(/%}/),
                django_value: i.starting_with(/{{/).until_after(/}}/),
                django_comment: i.starting_with(/{#/).until_after(/#}/),
                smarty: i.starting_with(/{(?=[^}{\s\n])/).until_after(/[^\s\n]}/),
                smarty_comment: i.starting_with(/{\*/).until_after(/\*}/),
                smarty_literal: i.starting_with(/{literal}/).until_after(/{\/literal}/)
            };
        }
        i.prototype = new n, i.prototype._create = function() {
            return new i(this._input, this);
        }, i.prototype._update = function() {
            this.__set_templated_pattern();
        }, i.prototype.disable = function(e) {
            var t = this._create();
            return t._disabled[e] = !0, t._update(), t;
        }, i.prototype.read_options = function(e) {
            var t = this._create();
            for(var n in r)t._disabled[n] = e.templating.indexOf(n) === -1;
            return t._update(), t;
        }, i.prototype.exclude = function(e) {
            var t = this._create();
            return t._excluded[e] = !0, t._update(), t;
        }, i.prototype.read = function() {
            var e = ``;
            e = this._match_pattern ? this._input.read(this._starting_pattern) : this._input.read(this._starting_pattern, this.__template_pattern);
            for(var t = this._read_template(); t;)this._match_pattern ? t += this._input.read(this._match_pattern) : t += this._input.readUntil(this.__template_pattern), e += t, t = this._read_template();
            return this._until_after && (e += this._input.readUntilAfter(this._until_pattern)), e;
        }, i.prototype.__set_templated_pattern = function() {
            var e = [];
            this._disabled.php || e.push(this.__patterns.php._starting_pattern.source), this._disabled.handlebars || e.push(this.__patterns.handlebars._starting_pattern.source), this._disabled.angular || e.push(this.__patterns.handlebars._starting_pattern.source), this._disabled.erb || e.push(this.__patterns.erb._starting_pattern.source), this._disabled.django || (e.push(this.__patterns.django._starting_pattern.source), e.push(this.__patterns.django_value._starting_pattern.source), e.push(this.__patterns.django_comment._starting_pattern.source)), this._disabled.smarty || e.push(this.__patterns.smarty._starting_pattern.source), this._until_pattern && e.push(this._until_pattern.source), this.__template_pattern = this._input.get_regexp(`(?:` + e.join(`|`) + `)`);
        }, i.prototype._read_template = function() {
            var e = ``, t = this._input.peek();
            if (t === `<`) {
                var n = this._input.peek(1);
                !this._disabled.php && !this._excluded.php && n === `?` && (e ||= this.__patterns.php.read()), !this._disabled.erb && !this._excluded.erb && n === `%` && (e ||= this.__patterns.erb.read());
            } else t === `{` && (!this._disabled.handlebars && !this._excluded.handlebars && (e ||= this.__patterns.handlebars_comment.read(), e ||= this.__patterns.handlebars_unescaped.read(), e ||= this.__patterns.handlebars.read()), this._disabled.django || (!this._excluded.django && !this._excluded.handlebars && (e ||= this.__patterns.django_value.read()), this._excluded.django || (e ||= this.__patterns.django_comment.read(), e ||= this.__patterns.django.read())), this._disabled.smarty || this._disabled.django && this._disabled.handlebars && (e ||= this.__patterns.smarty_comment.read(), e ||= this.__patterns.smarty_literal.read(), e ||= this.__patterns.smarty.read()));
            return e;
        }, t.exports.TemplatablePattern = i;
    })), jt = t(((e, t)=>{
        var n = wt().InputScanner, r = Ot().Tokenizer, i = Ot().TOKEN, a = kt().Directives, o = xt(), s = Et().Pattern, c = At().TemplatablePattern;
        function l(e, t) {
            return t.indexOf(e) !== -1;
        }
        var u = {
            START_EXPR: `TK_START_EXPR`,
            END_EXPR: `TK_END_EXPR`,
            START_BLOCK: `TK_START_BLOCK`,
            END_BLOCK: `TK_END_BLOCK`,
            WORD: `TK_WORD`,
            RESERVED: `TK_RESERVED`,
            SEMICOLON: `TK_SEMICOLON`,
            STRING: `TK_STRING`,
            EQUALS: `TK_EQUALS`,
            OPERATOR: `TK_OPERATOR`,
            COMMA: `TK_COMMA`,
            BLOCK_COMMENT: `TK_BLOCK_COMMENT`,
            COMMENT: `TK_COMMENT`,
            DOT: `TK_DOT`,
            UNKNOWN: `TK_UNKNOWN`,
            START: i.START,
            RAW: i.RAW,
            EOF: i.EOF
        }, d = new a(/\/\*/, /\*\//), f = /0[xX][0123456789abcdefABCDEF_]*n?|0[oO][01234567_]*n?|0[bB][01_]*n?|\d[\d_]*n|(?:\.\d[\d_]*|\d[\d_]*\.?[\d_]*)(?:[eE][+-]?[\d_]+)?/, p = /[0-9]/, m = /[^\d\.]/, h = `>>> === !== &&= ??= ||= << && >= ** != == <= >> || ?? |> < / - + > : & % ? ^ | *`.split(` `), g = `>>>= ... >>= <<= === >>> !== **= &&= ??= ||= => ^= :: /= << <= == && -= >= >> != -- += ** || ?? ++ %= &= *= |= |> = ! ? > < : / ^ - + * & % ~ |`;
        g = g.replace(/[-[\]{}()*+?.,\\^$|#]/g, `\\$&`), g = `\\?\\.(?!\\d) ` + g, g = g.replace(/ /g, `|`);
        var _ = new RegExp(g), v = `continue,try,throw,return,var,let,const,if,switch,case,default,for,while,break,function,import,export`.split(`,`), y = v.concat([
            `do`,
            `in`,
            `of`,
            `else`,
            `get`,
            `set`,
            `new`,
            `catch`,
            `finally`,
            `typeof`,
            `yield`,
            `async`,
            `await`,
            `from`,
            `as`,
            `class`,
            `extends`
        ]), b = RegExp(`^(?:` + y.join(`|`) + `)$`), x, S = function(e, t) {
            r.call(this, e, t), this._patterns.whitespace = this._patterns.whitespace.matching(`\\u00A0\\u1680\\u180e\\u2000-\\u200a\\u202f\\u205f\\u3000\\ufeff`, `\\u2028\\u2029`);
            var n = new s(this._input), i = new c(this._input).read_options(this._options);
            this.__patterns = {
                template: i,
                identifier: i.starting_with(o.identifier).matching(o.identifierMatch),
                number: n.matching(f),
                punct: n.matching(_),
                comment: n.starting_with(/\/\//).until(/[\n\r\u2028\u2029]/),
                block_comment: n.starting_with(/\/\*/).until_after(/\*\//),
                html_comment_start: n.matching(/<!--/),
                html_comment_end: n.matching(/-->/),
                include: n.starting_with(/#include/).until_after(o.lineBreak),
                shebang: n.starting_with(/#!/).until_after(o.lineBreak),
                xml: n.matching(/[\s\S]*?<(\/?)([-a-zA-Z:0-9_.]+|{[^}]+?}|!\[CDATA\[[^\]]*?\]\]|)(\s*{[^}]+?}|\s+[-a-zA-Z:0-9_.]+|\s+[-a-zA-Z:0-9_.]+\s*=\s*('[^']*'|"[^"]*"|{([^{}]|{[^}]+?})+?}))*\s*(\/?)\s*>/),
                single_quote: i.until(/['\\\n\r\u2028\u2029]/),
                double_quote: i.until(/["\\\n\r\u2028\u2029]/),
                template_text: i.until(/[`\\$]/),
                template_expression: i.until(/[`}\\]/)
            };
        };
        S.prototype = new r, S.prototype._is_comment = function(e) {
            return e.type === u.COMMENT || e.type === u.BLOCK_COMMENT || e.type === u.UNKNOWN;
        }, S.prototype._is_opening = function(e) {
            return e.type === u.START_BLOCK || e.type === u.START_EXPR;
        }, S.prototype._is_closing = function(e, t) {
            return (e.type === u.END_BLOCK || e.type === u.END_EXPR) && t && (e.text === `]` && t.text === `[` || e.text === `)` && t.text === `(` || e.text === `}` && t.text === `{`);
        }, S.prototype._reset = function() {
            x = !1;
        }, S.prototype._get_next_token = function(e, t) {
            var n = null;
            this._readWhitespace();
            var r = this._input.peek();
            return r === null ? this._create_token(u.EOF, ``) : (n ||= this._read_non_javascript(r), n ||= this._read_string(r), n ||= this._read_pair(r, this._input.peek(1)), n ||= this._read_word(e), n ||= this._read_singles(r), n ||= this._read_comment(r), n ||= this._read_regexp(r, e), n ||= this._read_xml(r, e), n ||= this._read_punctuation(), n ||= this._create_token(u.UNKNOWN, this._input.next()), n);
        }, S.prototype._read_word = function(e) {
            var t = this.__patterns.identifier.read();
            if (t !== ``) return t = t.replace(o.allLineBreaks, `
`), !(e.type === u.DOT || e.type === u.RESERVED && (e.text === `set` || e.text === `get`)) && b.test(t) ? (t === `in` || t === `of`) && (e.type === u.WORD || e.type === u.STRING) ? this._create_token(u.OPERATOR, t) : this._create_token(u.RESERVED, t) : this._create_token(u.WORD, t);
            if (t = this.__patterns.number.read(), t !== ``) return this._create_token(u.WORD, t);
        }, S.prototype._read_singles = function(e) {
            var t = null;
            return e === `(` || e === `[` ? t = this._create_token(u.START_EXPR, e) : e === `)` || e === `]` ? t = this._create_token(u.END_EXPR, e) : e === `{` ? t = this._create_token(u.START_BLOCK, e) : e === `}` ? t = this._create_token(u.END_BLOCK, e) : e === `;` ? t = this._create_token(u.SEMICOLON, e) : e === `.` && m.test(this._input.peek(1)) ? t = this._create_token(u.DOT, e) : e === `,` && (t = this._create_token(u.COMMA, e)), t && this._input.next(), t;
        }, S.prototype._read_pair = function(e, t) {
            var n = null;
            return e === `#` && t === `{` && (n = this._create_token(u.START_BLOCK, e + t)), n && (this._input.next(), this._input.next()), n;
        }, S.prototype._read_punctuation = function() {
            var e = this.__patterns.punct.read();
            if (e !== ``) return e === `=` ? this._create_token(u.EQUALS, e) : e === `?.` ? this._create_token(u.DOT, e) : this._create_token(u.OPERATOR, e);
        }, S.prototype._read_non_javascript = function(e) {
            var t = ``;
            if (e === `#`) {
                if (this._is_first_token() && (t = this.__patterns.shebang.read(), t) || (t = this.__patterns.include.read(), t)) return this._create_token(u.UNKNOWN, t.trim() + `
`);
                e = this._input.next();
                var n = `#`;
                if (this._input.hasNext() && this._input.testChar(p)) {
                    do e = this._input.next(), n += e;
                    while (this._input.hasNext() && e !== `#` && e !== `=`);
                    return e === `#` || (this._input.peek() === `[` && this._input.peek(1) === `]` ? (n += `[]`, this._input.next(), this._input.next()) : this._input.peek() === `{` && this._input.peek(1) === `}` && (n += `{}`, this._input.next(), this._input.next())), this._create_token(u.WORD, n);
                }
                this._input.back();
            } else if (e === `<` && this._is_first_token()) {
                if (t = this.__patterns.html_comment_start.read(), t) {
                    for(; this._input.hasNext() && !this._input.testChar(o.newline);)t += this._input.next();
                    return x = !0, this._create_token(u.COMMENT, t);
                }
            } else if (x && e === `-` && (t = this.__patterns.html_comment_end.read(), t)) return x = !1, this._create_token(u.COMMENT, t);
            return null;
        }, S.prototype._read_comment = function(e) {
            var t = null;
            if (e === `/`) {
                var n = ``;
                if (this._input.peek(1) === `*`) {
                    n = this.__patterns.block_comment.read();
                    var r = d.get_directives(n);
                    r && r.ignore === `start` && (n += d.readIgnored(this._input)), n = n.replace(o.allLineBreaks, `
`), t = this._create_token(u.BLOCK_COMMENT, n), t.directives = r;
                } else this._input.peek(1) === `/` && (n = this.__patterns.comment.read(), t = this._create_token(u.COMMENT, n));
            }
            return t;
        }, S.prototype._read_string = function(e) {
            if (e === "`" || e === `'` || e === `"`) {
                var t = this._input.next();
                return this.has_char_escapes = !1, e === "`" ? t += this._read_string_recursive("`", !0, "${") : t += this._read_string_recursive(e), this.has_char_escapes && this._options.unescape_strings && (t = C(t)), this._input.peek() === e && (t += this._input.next()), t = t.replace(o.allLineBreaks, `
`), this._create_token(u.STRING, t);
            }
            return null;
        }, S.prototype._allow_regexp_or_xml = function(e) {
            return e.type === u.RESERVED && l(e.text, [
                `return`,
                `case`,
                `throw`,
                `else`,
                `do`,
                `typeof`,
                `yield`
            ]) || e.type === u.END_EXPR && e.text === `)` && e.opened.previous.type === u.RESERVED && l(e.opened.previous.text, [
                `if`,
                `while`,
                `for`
            ]) || l(e.type, [
                u.COMMENT,
                u.START_EXPR,
                u.START_BLOCK,
                u.START,
                u.END_BLOCK,
                u.OPERATOR,
                u.EQUALS,
                u.EOF,
                u.SEMICOLON,
                u.COMMA
            ]);
        }, S.prototype._read_regexp = function(e, t) {
            if (e === `/` && this._allow_regexp_or_xml(t)) {
                for(var n = this._input.next(), r = !1, i = !1; this._input.hasNext() && (r || i || this._input.peek() !== e) && !this._input.testChar(o.newline);)n += this._input.peek(), r ? r = !1 : (r = this._input.peek() === `\\`, this._input.peek() === `[` ? i = !0 : this._input.peek() === `]` && (i = !1)), this._input.next();
                return this._input.peek() === e && (n += this._input.next(), n += this._input.read(o.identifier)), this._create_token(u.STRING, n);
            }
            return null;
        }, S.prototype._read_xml = function(e, t) {
            if (this._options.e4x && e === `<` && this._allow_regexp_or_xml(t)) {
                var n = ``, r = this.__patterns.xml.read_match();
                if (r) {
                    for(var i = r[2].replace(/^{\s+/, `{`).replace(/\s+}$/, `}`), a = i.indexOf(`{`) === 0, s = 0; r;){
                        var c = !!r[1], l = r[2];
                        if (!(r[r.length - 1] || l.slice(0, 8) === `![CDATA[`) && (l === i || a && l.replace(/^{\s+/, `{`).replace(/\s+}$/, `}`)) && (c ? --s : ++s), n += r[0], s <= 0) break;
                        r = this.__patterns.xml.read_match();
                    }
                    return r || (n += this._input.match(/[\s\S]*/g)[0]), n = n.replace(o.allLineBreaks, `
`), this._create_token(u.STRING, n);
                }
            }
            return null;
        };
        function C(e) {
            for(var t = ``, r = 0, i = new n(e), a = null; i.hasNext();)if (a = i.match(/([\s]|[^\\]|\\\\)+/g), a && (t += a[0]), i.peek() === `\\`) {
                if (i.next(), i.peek() === `x`) a = i.match(/x([0-9A-Fa-f]{2})/g);
                else if (i.peek() === `u`) a = i.match(/u([0-9A-Fa-f]{4})/g), a ||= i.match(/u\{([0-9A-Fa-f]+)\}/g);
                else {
                    t += `\\`, i.hasNext() && (t += i.next());
                    continue;
                }
                if (!a || (r = parseInt(a[1], 16), r > 126 && r <= 255 && a[0].indexOf(`x`) === 0)) return e;
                r >= 0 && r < 32 || r > 1114111 ? t += `\\` + a[0] : r === 34 || r === 39 || r === 92 ? t += `\\` + String.fromCharCode(r) : t += String.fromCharCode(r);
            }
            return t;
        }
        S.prototype._read_string_recursive = function(e, t, n) {
            var r, i;
            e === `'` ? i = this.__patterns.single_quote : e === `"` ? i = this.__patterns.double_quote : e === "`" ? i = this.__patterns.template_text : e === `}` && (i = this.__patterns.template_expression);
            for(var a = i.read(), s = ``; this._input.hasNext();){
                if (s = this._input.next(), s === e || !t && o.newline.test(s)) {
                    this._input.back();
                    break;
                } else s === `\\` && this._input.hasNext() ? (r = this._input.peek(), r === `x` || r === `u` ? this.has_char_escapes = !0 : r === `\r` && this._input.peek(1) === `
` && this._input.next(), s += this._input.next()) : n && (n === "${" && s === `$` && this._input.peek() === `{` && (s += this._input.next()), n === s && (e === "`" ? s += this._read_string_recursive(`}`, t, "`") : s += this._read_string_recursive("`", t, "${"), this._input.hasNext() && (s += this._input.next())));
                s += i.read(), a += s;
            }
            return a;
        }, t.exports.Tokenizer = S, t.exports.TOKEN = u, t.exports.positionable_operators = h.slice(), t.exports.line_starters = v.slice();
    })), Mt = t(((e, t)=>{
        var n = yt().Output, r = bt().Token, i = xt(), a = Ct().Options, o = jt().Tokenizer, s = jt().line_starters, c = jt().positionable_operators, l = jt().TOKEN;
        function u(e, t) {
            return t.indexOf(e) !== -1;
        }
        function d(e) {
            return e.replace(/^\s+/g, ``);
        }
        function f(e) {
            for(var t = {}, n = 0; n < e.length; n++)t[e[n].replace(/-/g, `_`)] = e[n];
            return t;
        }
        function p(e, t) {
            return e && e.type === l.RESERVED && e.text === t;
        }
        function m(e, t) {
            return e && e.type === l.RESERVED && u(e.text, t);
        }
        var h = [
            `case`,
            `return`,
            `do`,
            `if`,
            `throw`,
            `else`,
            `await`,
            `break`,
            `continue`,
            `async`
        ], g = f([
            `before-newline`,
            `after-newline`,
            `preserve-newline`
        ]), _ = [
            g.before_newline,
            g.preserve_newline
        ], v = {
            BlockStatement: `BlockStatement`,
            Statement: `Statement`,
            ObjectLiteral: `ObjectLiteral`,
            ArrayLiteral: `ArrayLiteral`,
            ForInitializer: `ForInitializer`,
            Conditional: `Conditional`,
            Expression: `Expression`
        };
        function y(e, t) {
            t.multiline_frame || t.mode === v.ForInitializer || t.mode === v.Conditional || e.remove_indent(t.start_line_index);
        }
        function b(e) {
            e = e.replace(i.allLineBreaks, `
`);
            for(var t = [], n = e.indexOf(`
`); n !== -1;)t.push(e.substring(0, n)), e = e.substring(n + 1), n = e.indexOf(`
`);
            return e.length && t.push(e), t;
        }
        function x(e) {
            return e === v.ArrayLiteral;
        }
        function S(e) {
            return u(e, [
                v.Expression,
                v.ForInitializer,
                v.Conditional
            ]);
        }
        function C(e, t) {
            for(var n = 0; n < e.length; n++)if (e[n].trim().charAt(0) !== t) return !1;
            return !0;
        }
        function w(e, t) {
            for(var n = 0, r = e.length, i; n < r; n++)if (i = e[n], i && i.indexOf(t) !== 0) return !1;
            return !0;
        }
        function T(e, t) {
            t ||= {}, this._source_text = e || ``, this._output = null, this._tokens = null, this._last_last_text = null, this._flags = null, this._previous_flags = null, this._flag_store = null, this._options = new a(t);
        }
        T.prototype.create_flags = function(e, t) {
            var n = 0;
            return e && (n = e.indentation_level, !this._output.just_added_newline() && e.line_indent_level > n && (n = e.line_indent_level)), {
                mode: t,
                parent: e,
                last_token: e ? e.last_token : new r(l.START_BLOCK, ``),
                last_word: e ? e.last_word : ``,
                declaration_statement: !1,
                declaration_assignment: !1,
                multiline_frame: !1,
                inline_frame: !1,
                if_block: !1,
                else_block: !1,
                class_start_block: !1,
                do_block: !1,
                do_while: !1,
                import_block: !1,
                in_case_statement: !1,
                in_case: !1,
                case_body: !1,
                case_block: !1,
                indentation_level: n,
                alignment: 0,
                line_indent_level: e ? e.line_indent_level : n,
                start_line_index: this._output.get_line_number(),
                ternary_depth: 0
            };
        }, T.prototype._reset = function(e) {
            var t = e.match(/^[\t ]*/)[0];
            this._last_last_text = ``, this._output = new n(this._options, t), this._output.raw = this._options.test_output_raw, this._flag_store = [], this.set_mode(v.BlockStatement);
            var r = new o(e, this._options);
            return this._tokens = r.tokenize(), e;
        }, T.prototype.beautify = function() {
            if (this._options.disabled) return this._source_text;
            var e, t = this._reset(this._source_text), n = this._options.eol;
            this._options.eol === `auto` && (n = `
`, t && i.lineBreak.test(t || ``) && (n = t.match(i.lineBreak)[0]));
            for(var r = this._tokens.next(); r;)this.handle_token(r), this._last_last_text = this._flags.last_token.text, this._flags.last_token = r, r = this._tokens.next();
            return e = this._output.get_code(n), e;
        }, T.prototype.handle_token = function(e, t) {
            e.type === l.START_EXPR ? this.handle_start_expr(e) : e.type === l.END_EXPR ? this.handle_end_expr(e) : e.type === l.START_BLOCK ? this.handle_start_block(e) : e.type === l.END_BLOCK ? this.handle_end_block(e) : e.type === l.WORD || e.type === l.RESERVED ? this.handle_word(e) : e.type === l.SEMICOLON ? this.handle_semicolon(e) : e.type === l.STRING ? this.handle_string(e) : e.type === l.EQUALS ? this.handle_equals(e) : e.type === l.OPERATOR ? this.handle_operator(e) : e.type === l.COMMA ? this.handle_comma(e) : e.type === l.BLOCK_COMMENT ? this.handle_block_comment(e, t) : e.type === l.COMMENT ? this.handle_comment(e, t) : e.type === l.DOT ? this.handle_dot(e) : e.type === l.EOF ? this.handle_eof(e) : (e.type, l.UNKNOWN, this.handle_unknown(e, t));
        }, T.prototype.handle_whitespace_and_comments = function(e, t) {
            var n = e.newlines, r = this._options.keep_array_indentation && x(this._flags.mode);
            if (e.comments_before) for(var i = e.comments_before.next(); i;)this.handle_whitespace_and_comments(i, t), this.handle_token(i, t), i = e.comments_before.next();
            if (r) for(var a = 0; a < n; a += 1)this.print_newline(a > 0, t);
            else if (this._options.max_preserve_newlines && n > this._options.max_preserve_newlines && (n = this._options.max_preserve_newlines), this._options.preserve_newlines && n > 1) {
                this.print_newline(!1, t);
                for(var o = 1; o < n; o += 1)this.print_newline(!0, t);
            }
        };
        var E = [
            `async`,
            `break`,
            `continue`,
            `return`,
            `throw`,
            `yield`
        ];
        T.prototype.allow_wrap_or_preserved_newline = function(e, t) {
            if (t = t === void 0 ? !1 : t, !this._output.just_added_newline()) {
                var n = this._options.preserve_newlines && e.newlines || t;
                if (u(this._flags.last_token.text, c) || u(e.text, c)) {
                    var r = u(this._flags.last_token.text, c) && u(this._options.operator_position, _) || u(e.text, c);
                    n &&= r;
                }
                if (n) this.print_newline(!1, !0);
                else if (this._options.wrap_line_length) {
                    if (m(this._flags.last_token, E)) return;
                    this._output.set_wrap_point();
                }
            }
        }, T.prototype.print_newline = function(e, t) {
            if (!t && this._flags.last_token.text !== `;` && this._flags.last_token.text !== `,` && this._flags.last_token.text !== `=` && (this._flags.last_token.type !== l.OPERATOR || this._flags.last_token.text === `--` || this._flags.last_token.text === `++`)) for(var n = this._tokens.peek(); this._flags.mode === v.Statement && !(this._flags.if_block && p(n, `else`)) && !this._flags.do_block;)this.restore_mode();
            this._output.add_new_line(e) && (this._flags.multiline_frame = !0);
        }, T.prototype.print_token_line_indentation = function(e) {
            this._output.just_added_newline() && (this._options.keep_array_indentation && e.newlines && (e.text === `[` || x(this._flags.mode)) ? (this._output.current_line.set_indent(-1), this._output.current_line.push(e.whitespace_before), this._output.space_before_token = !1) : this._output.set_indent(this._flags.indentation_level, this._flags.alignment) && (this._flags.line_indent_level = this._flags.indentation_level));
        }, T.prototype.print_token = function(e) {
            if (this._output.raw) {
                this._output.add_raw_token(e);
                return;
            }
            if (this._options.comma_first && e.previous && e.previous.type === l.COMMA && this._output.just_added_newline() && this._output.previous_line.last() === `,`) {
                var t = this._output.previous_line.pop();
                this._output.previous_line.is_empty() && (this._output.previous_line.push(t), this._output.trim(!0), this._output.current_line.pop(), this._output.trim()), this.print_token_line_indentation(e), this._output.add_token(`,`), this._output.space_before_token = !0;
            }
            this.print_token_line_indentation(e), this._output.non_breaking_space = !0, this._output.add_token(e.text), this._output.previous_token_wrapped && (this._flags.multiline_frame = !0);
        }, T.prototype.indent = function() {
            this._flags.indentation_level += 1, this._output.set_indent(this._flags.indentation_level, this._flags.alignment);
        }, T.prototype.deindent = function() {
            this._flags.indentation_level > 0 && (!this._flags.parent || this._flags.indentation_level > this._flags.parent.indentation_level) && (--this._flags.indentation_level, this._output.set_indent(this._flags.indentation_level, this._flags.alignment));
        }, T.prototype.set_mode = function(e) {
            this._flags ? (this._flag_store.push(this._flags), this._previous_flags = this._flags) : this._previous_flags = this.create_flags(null, e), this._flags = this.create_flags(this._previous_flags, e), this._output.set_indent(this._flags.indentation_level, this._flags.alignment);
        }, T.prototype.restore_mode = function() {
            this._flag_store.length > 0 && (this._previous_flags = this._flags, this._flags = this._flag_store.pop(), this._previous_flags.mode === v.Statement && y(this._output, this._previous_flags), this._output.set_indent(this._flags.indentation_level, this._flags.alignment));
        }, T.prototype.start_of_object_property = function() {
            return this._flags.parent.mode === v.ObjectLiteral && this._flags.mode === v.Statement && (this._flags.last_token.text === `:` && this._flags.ternary_depth === 0 || m(this._flags.last_token, [
                `get`,
                `set`
            ]));
        }, T.prototype.start_of_statement = function(e) {
            var t = !1;
            return t ||= m(this._flags.last_token, [
                `var`,
                `let`,
                `const`
            ]) && e.type === l.WORD, t ||= p(this._flags.last_token, `do`), t ||= !(this._flags.parent.mode === v.ObjectLiteral && this._flags.mode === v.Statement) && m(this._flags.last_token, E) && !e.newlines, t ||= p(this._flags.last_token, `else`) && !(p(e, `if`) && !e.comments_before), t ||= this._flags.last_token.type === l.END_EXPR && (this._previous_flags.mode === v.ForInitializer || this._previous_flags.mode === v.Conditional), t ||= this._flags.last_token.type === l.WORD && this._flags.mode === v.BlockStatement && !this._flags.in_case && !(e.text === `--` || e.text === `++`) && this._last_last_text !== `function` && e.type !== l.WORD && e.type !== l.RESERVED, t ||= this._flags.mode === v.ObjectLiteral && (this._flags.last_token.text === `:` && this._flags.ternary_depth === 0 || m(this._flags.last_token, [
                `get`,
                `set`
            ])), t ? (this.set_mode(v.Statement), this.indent(), this.handle_whitespace_and_comments(e, !0), this.start_of_object_property() || this.allow_wrap_or_preserved_newline(e, m(e, [
                `do`,
                `for`,
                `if`,
                `while`
            ])), !0) : !1;
        }, T.prototype.handle_start_expr = function(e) {
            this.start_of_statement(e) || this.handle_whitespace_and_comments(e);
            var t = v.Expression;
            if (e.text === `[`) {
                if (this._flags.last_token.type === l.WORD || this._flags.last_token.text === `)`) {
                    m(this._flags.last_token, s) && (this._output.space_before_token = !0), this.print_token(e), this.set_mode(t), this.indent(), this._options.space_in_paren && (this._output.space_before_token = !0);
                    return;
                }
                t = v.ArrayLiteral, x(this._flags.mode) && (this._flags.last_token.text === `[` || this._flags.last_token.text === `,` && (this._last_last_text === `]` || this._last_last_text === `}`)) && (this._options.keep_array_indentation || this.print_newline()), u(this._flags.last_token.type, [
                    l.START_EXPR,
                    l.END_EXPR,
                    l.WORD,
                    l.OPERATOR,
                    l.DOT
                ]) || (this._output.space_before_token = !0);
            } else {
                if (this._flags.last_token.type === l.RESERVED) this._flags.last_token.text === `for` ? (this._output.space_before_token = this._options.space_before_conditional, t = v.ForInitializer) : u(this._flags.last_token.text, [
                    `if`,
                    `while`,
                    `switch`
                ]) ? (this._output.space_before_token = this._options.space_before_conditional, t = v.Conditional) : u(this._flags.last_word, [
                    `await`,
                    `async`
                ]) ? this._output.space_before_token = !0 : this._flags.last_token.text === `import` && e.whitespace_before === `` ? this._output.space_before_token = !1 : (u(this._flags.last_token.text, s) || this._flags.last_token.text === `catch`) && (this._output.space_before_token = !0);
                else if (this._flags.last_token.type === l.EQUALS || this._flags.last_token.type === l.OPERATOR) this.start_of_object_property() || this.allow_wrap_or_preserved_newline(e);
                else if (this._flags.last_token.type === l.WORD) {
                    this._output.space_before_token = !1;
                    var n = this._tokens.peek(-3);
                    if (this._options.space_after_named_function && n) {
                        var r = this._tokens.peek(-4);
                        m(n, [
                            `async`,
                            `function`
                        ]) || n.text === `*` && m(r, [
                            `async`,
                            `function`
                        ]) ? this._output.space_before_token = !0 : this._flags.mode === v.ObjectLiteral ? (n.text === `{` || n.text === `,` || n.text === `*` && (r.text === `{` || r.text === `,`)) && (this._output.space_before_token = !0) : this._flags.parent && this._flags.parent.class_start_block && (this._output.space_before_token = !0);
                    }
                } else this.allow_wrap_or_preserved_newline(e);
                (this._flags.last_token.type === l.RESERVED && (this._flags.last_word === `function` || this._flags.last_word === `typeof`) || this._flags.last_token.text === `*` && (u(this._last_last_text, [
                    `function`,
                    `yield`
                ]) || this._flags.mode === v.ObjectLiteral && u(this._last_last_text, [
                    `{`,
                    `,`
                ]))) && (this._output.space_before_token = this._options.space_after_anon_function);
            }
            this._flags.last_token.text === `;` || this._flags.last_token.type === l.START_BLOCK ? this.print_newline() : (this._flags.last_token.type === l.END_EXPR || this._flags.last_token.type === l.START_EXPR || this._flags.last_token.type === l.END_BLOCK || this._flags.last_token.text === `.` || this._flags.last_token.type === l.COMMA) && this.allow_wrap_or_preserved_newline(e, e.newlines), this.print_token(e), this.set_mode(t), this._options.space_in_paren && (this._output.space_before_token = !0), this.indent();
        }, T.prototype.handle_end_expr = function(e) {
            for(; this._flags.mode === v.Statement;)this.restore_mode();
            this.handle_whitespace_and_comments(e), this._flags.multiline_frame && this.allow_wrap_or_preserved_newline(e, e.text === `]` && x(this._flags.mode) && !this._options.keep_array_indentation), this._options.space_in_paren && (this._flags.last_token.type === l.START_EXPR && !this._options.space_in_empty_paren ? (this._output.trim(), this._output.space_before_token = !1) : this._output.space_before_token = !0), this.deindent(), this.print_token(e), this.restore_mode(), y(this._output, this._previous_flags), this._flags.do_while && this._previous_flags.mode === v.Conditional && (this._previous_flags.mode = v.Expression, this._flags.do_block = !1, this._flags.do_while = !1);
        }, T.prototype.handle_start_block = function(e) {
            this.handle_whitespace_and_comments(e);
            var t = this._tokens.peek(), n = this._tokens.peek(1);
            this._flags.last_word === `switch` && this._flags.last_token.type === l.END_EXPR ? (this.set_mode(v.BlockStatement), this._flags.in_case_statement = !0) : this._flags.case_body ? this.set_mode(v.BlockStatement) : n && (u(n.text, [
                `:`,
                `,`
            ]) && u(t.type, [
                l.STRING,
                l.WORD,
                l.RESERVED
            ]) || u(t.text, [
                `get`,
                `set`,
                `...`
            ]) && u(n.type, [
                l.WORD,
                l.RESERVED
            ])) ? u(this._last_last_text, [
                `class`,
                `interface`
            ]) && !u(n.text, [
                `:`,
                `,`
            ]) ? this.set_mode(v.BlockStatement) : this.set_mode(v.ObjectLiteral) : this._flags.last_token.type === l.OPERATOR && this._flags.last_token.text === `=>` ? this.set_mode(v.BlockStatement) : u(this._flags.last_token.type, [
                l.EQUALS,
                l.START_EXPR,
                l.COMMA,
                l.OPERATOR
            ]) || m(this._flags.last_token, [
                `return`,
                `throw`,
                `import`,
                `default`
            ]) ? this.set_mode(v.ObjectLiteral) : this.set_mode(v.BlockStatement), this._flags.last_token && m(this._flags.last_token.previous, [
                `class`,
                `extends`
            ]) && (this._flags.class_start_block = !0);
            var r = !t.comments_before && t.text === `}`, i = r && this._flags.last_word === `function` && this._flags.last_token.type === l.END_EXPR;
            if (this._options.brace_preserve_inline) {
                var a = 0, o = null;
                this._flags.inline_frame = !0;
                do if (a += 1, o = this._tokens.peek(a - 1), o.newlines) {
                    this._flags.inline_frame = !1;
                    break;
                }
                while (o.type !== l.EOF && !(o.type === l.END_BLOCK && o.opened === e));
            }
            (this._options.brace_style === `expand` || this._options.brace_style === `none` && e.newlines) && !this._flags.inline_frame ? this._flags.last_token.type !== l.OPERATOR && (i || this._flags.last_token.type === l.EQUALS || m(this._flags.last_token, h) && this._flags.last_token.text !== `else`) ? this._output.space_before_token = !0 : this.print_newline(!1, !0) : (x(this._previous_flags.mode) && (this._flags.last_token.type === l.START_EXPR || this._flags.last_token.type === l.COMMA) && ((this._flags.last_token.type === l.COMMA || this._options.space_in_paren) && (this._output.space_before_token = !0), (this._flags.last_token.type === l.COMMA || this._flags.last_token.type === l.START_EXPR && this._flags.inline_frame) && (this.allow_wrap_or_preserved_newline(e), this._previous_flags.multiline_frame = this._previous_flags.multiline_frame || this._flags.multiline_frame, this._flags.multiline_frame = !1)), this._flags.last_token.type !== l.OPERATOR && this._flags.last_token.type !== l.START_EXPR && (u(this._flags.last_token.type, [
                l.START_BLOCK,
                l.SEMICOLON
            ]) && !this._flags.inline_frame ? this.print_newline() : this._output.space_before_token = !0)), this.print_token(e), this.indent(), !r && !(this._options.brace_preserve_inline && this._flags.inline_frame) && this.print_newline();
        }, T.prototype.handle_end_block = function(e) {
            for(this.handle_whitespace_and_comments(e); this._flags.mode === v.Statement;)this.restore_mode();
            var t = this._flags.last_token.type === l.START_BLOCK;
            this._flags.inline_frame && !t ? this._output.space_before_token = !0 : this._options.brace_style === `expand` ? t || this.print_newline() : t || (x(this._flags.mode) && this._options.keep_array_indentation ? (this._options.keep_array_indentation = !1, this.print_newline(), this._options.keep_array_indentation = !0) : this.print_newline()), this.restore_mode(), this.print_token(e);
        }, T.prototype.handle_word = function(e) {
            if (e.type === l.RESERVED && (u(e.text, [
                `set`,
                `get`
            ]) && this._flags.mode !== v.ObjectLiteral || e.text === `import` && u(this._tokens.peek().text, [
                `(`,
                `.`
            ]) || u(e.text, [
                `as`,
                `from`
            ]) && !this._flags.import_block || this._flags.mode === v.ObjectLiteral && this._tokens.peek().text === `:`) && (e.type = l.WORD), this.start_of_statement(e) ? m(this._flags.last_token, [
                `var`,
                `let`,
                `const`
            ]) && e.type === l.WORD && (this._flags.declaration_statement = !0) : e.newlines && !S(this._flags.mode) && (this._flags.last_token.type !== l.OPERATOR || this._flags.last_token.text === `--` || this._flags.last_token.text === `++`) && this._flags.last_token.type !== l.EQUALS && (this._options.preserve_newlines || !m(this._flags.last_token, [
                `var`,
                `let`,
                `const`,
                `set`,
                `get`
            ])) ? (this.handle_whitespace_and_comments(e), this.print_newline()) : this.handle_whitespace_and_comments(e), this._flags.do_block && !this._flags.do_while) if (p(e, `while`)) {
                this._output.space_before_token = !0, this.print_token(e), this._output.space_before_token = !0, this._flags.do_while = !0;
                return;
            } else this.print_newline(), this._flags.do_block = !1;
            if (this._flags.if_block) if (!this._flags.else_block && p(e, `else`)) this._flags.else_block = !0;
            else {
                for(; this._flags.mode === v.Statement;)this.restore_mode();
                this._flags.if_block = !1, this._flags.else_block = !1;
            }
            if (this._flags.in_case_statement && m(e, [
                `case`,
                `default`
            ])) {
                this.print_newline(), !this._flags.case_block && (this._flags.case_body || this._options.jslint_happy) && this.deindent(), this._flags.case_body = !1, this.print_token(e), this._flags.in_case = !0;
                return;
            }
            if ((this._flags.last_token.type === l.COMMA || this._flags.last_token.type === l.START_EXPR || this._flags.last_token.type === l.EQUALS || this._flags.last_token.type === l.OPERATOR) && !this.start_of_object_property() && !(u(this._flags.last_token.text, [
                `+`,
                `-`
            ]) && this._last_last_text === `:` && this._flags.parent.mode === v.ObjectLiteral) && this.allow_wrap_or_preserved_newline(e), p(e, `function`)) {
                (u(this._flags.last_token.text, [
                    `}`,
                    `;`
                ]) || this._output.just_added_newline() && !(u(this._flags.last_token.text, [
                    `(`,
                    `[`,
                    `{`,
                    `:`,
                    `=`,
                    `,`
                ]) || this._flags.last_token.type === l.OPERATOR)) && !this._output.just_added_blankline() && !e.comments_before && (this.print_newline(), this.print_newline(!0)), this._flags.last_token.type === l.RESERVED || this._flags.last_token.type === l.WORD ? m(this._flags.last_token, [
                    `get`,
                    `set`,
                    `new`,
                    `export`
                ]) || m(this._flags.last_token, E) || p(this._flags.last_token, `default`) && this._last_last_text === `export` || this._flags.last_token.text === `declare` ? this._output.space_before_token = !0 : this.print_newline() : this._flags.last_token.type === l.OPERATOR || this._flags.last_token.text === `=` ? this._output.space_before_token = !0 : !this._flags.multiline_frame && (S(this._flags.mode) || x(this._flags.mode)) || this.print_newline(), this.print_token(e), this._flags.last_word = e.text;
                return;
            }
            var t = `NONE`;
            this._flags.last_token.type === l.END_BLOCK ? this._previous_flags.inline_frame ? t = `SPACE` : m(e, [
                `else`,
                `catch`,
                `finally`,
                `from`
            ]) ? this._options.brace_style === `expand` || this._options.brace_style === `end-expand` || this._options.brace_style === `none` && e.newlines ? t = `NEWLINE` : (t = `SPACE`, this._output.space_before_token = !0) : t = `NEWLINE` : this._flags.last_token.type === l.SEMICOLON && this._flags.mode === v.BlockStatement ? t = `NEWLINE` : this._flags.last_token.type === l.SEMICOLON && S(this._flags.mode) ? t = `SPACE` : this._flags.last_token.type === l.STRING ? t = `NEWLINE` : this._flags.last_token.type === l.RESERVED || this._flags.last_token.type === l.WORD || this._flags.last_token.text === `*` && (u(this._last_last_text, [
                `function`,
                `yield`
            ]) || this._flags.mode === v.ObjectLiteral && u(this._last_last_text, [
                `{`,
                `,`
            ])) ? t = `SPACE` : this._flags.last_token.type === l.START_BLOCK ? t = this._flags.inline_frame ? `SPACE` : `NEWLINE` : this._flags.last_token.type === l.END_EXPR && (this._output.space_before_token = !0, t = `NEWLINE`), m(e, s) && this._flags.last_token.text !== `)` && (t = this._flags.inline_frame || this._flags.last_token.text === `else` || this._flags.last_token.text === `export` ? `SPACE` : `NEWLINE`), m(e, [
                `else`,
                `catch`,
                `finally`
            ]) ? (!(this._flags.last_token.type === l.END_BLOCK && this._previous_flags.mode === v.BlockStatement) || this._options.brace_style === `expand` || this._options.brace_style === `end-expand` || this._options.brace_style === `none` && e.newlines) && !this._flags.inline_frame ? this.print_newline() : (this._output.trim(!0), this._output.current_line.last() !== `}` && this.print_newline(), this._output.space_before_token = !0) : t === `NEWLINE` ? m(this._flags.last_token, h) || this._flags.last_token.text === `declare` && m(e, [
                `var`,
                `let`,
                `const`
            ]) ? this._output.space_before_token = !0 : this._flags.last_token.type === l.END_EXPR ? m(e, s) && this._flags.last_token.text !== `)` && this.print_newline() : (this._flags.last_token.type !== l.START_EXPR || !m(e, [
                `var`,
                `let`,
                `const`
            ])) && this._flags.last_token.text !== `:` && (p(e, `if`) && p(e.previous, `else`) ? this._output.space_before_token = !0 : this.print_newline()) : this._flags.multiline_frame && x(this._flags.mode) && this._flags.last_token.text === `,` && this._last_last_text === `}` ? this.print_newline() : t === `SPACE` && (this._output.space_before_token = !0), e.previous && (e.previous.type === l.WORD || e.previous.type === l.RESERVED) && (this._output.space_before_token = !0), this.print_token(e), this._flags.last_word = e.text, e.type === l.RESERVED && (e.text === `do` ? this._flags.do_block = !0 : e.text === `if` ? this._flags.if_block = !0 : e.text === `import` ? this._flags.import_block = !0 : this._flags.import_block && p(e, `from`) && (this._flags.import_block = !1));
        }, T.prototype.handle_semicolon = function(e) {
            this.start_of_statement(e) ? this._output.space_before_token = !1 : this.handle_whitespace_and_comments(e);
            for(var t = this._tokens.peek(); this._flags.mode === v.Statement && !(this._flags.if_block && p(t, `else`)) && !this._flags.do_block;)this.restore_mode();
            this._flags.import_block && (this._flags.import_block = !1), this.print_token(e);
        }, T.prototype.handle_string = function(e) {
            e.text.startsWith("`") && e.newlines === 0 && e.whitespace_before === `` && (e.previous.text === `)` || this._flags.last_token.type === l.WORD) || (this.start_of_statement(e) ? this._output.space_before_token = !0 : (this.handle_whitespace_and_comments(e), this._flags.last_token.type === l.RESERVED || this._flags.last_token.type === l.WORD || this._flags.inline_frame ? this._output.space_before_token = !0 : this._flags.last_token.type === l.COMMA || this._flags.last_token.type === l.START_EXPR || this._flags.last_token.type === l.EQUALS || this._flags.last_token.type === l.OPERATOR ? this.start_of_object_property() || this.allow_wrap_or_preserved_newline(e) : e.text.startsWith("`") && this._flags.last_token.type === l.END_EXPR && (e.previous.text === `]` || e.previous.text === `)`) && e.newlines === 0 ? this._output.space_before_token = !0 : this.print_newline())), this.print_token(e);
        }, T.prototype.handle_equals = function(e) {
            this.start_of_statement(e) || this.handle_whitespace_and_comments(e), this._flags.declaration_statement && (this._flags.declaration_assignment = !0), this._output.space_before_token = !0, this.print_token(e), this._output.space_before_token = !0;
        }, T.prototype.handle_comma = function(e) {
            this.handle_whitespace_and_comments(e, !0), this.print_token(e), this._output.space_before_token = !0, this._flags.declaration_statement ? (S(this._flags.parent.mode) && (this._flags.declaration_assignment = !1), this._flags.declaration_assignment ? (this._flags.declaration_assignment = !1, this.print_newline(!1, !0)) : this._options.comma_first && this.allow_wrap_or_preserved_newline(e)) : this._flags.mode === v.ObjectLiteral || this._flags.mode === v.Statement && this._flags.parent.mode === v.ObjectLiteral ? (this._flags.mode === v.Statement && this.restore_mode(), this._flags.inline_frame || this.print_newline()) : this._options.comma_first && this.allow_wrap_or_preserved_newline(e);
        }, T.prototype.handle_operator = function(e) {
            var t = e.text === `*` && (m(this._flags.last_token, [
                `function`,
                `yield`
            ]) || u(this._flags.last_token.type, [
                l.START_BLOCK,
                l.COMMA,
                l.END_BLOCK,
                l.SEMICOLON
            ])), n = u(e.text, [
                `-`,
                `+`
            ]) && (u(this._flags.last_token.type, [
                l.START_BLOCK,
                l.START_EXPR,
                l.EQUALS,
                l.OPERATOR
            ]) || u(this._flags.last_token.text, s) || this._flags.last_token.text === `,`);
            if (!this.start_of_statement(e)) {
                var r = !t;
                this.handle_whitespace_and_comments(e, r);
            }
            if (e.text === `*` && this._flags.last_token.type === l.DOT) {
                this.print_token(e);
                return;
            }
            if (e.text === `::`) {
                this.print_token(e);
                return;
            }
            if (u(e.text, [
                `-`,
                `+`
            ]) && this.start_of_object_property()) {
                this.print_token(e);
                return;
            }
            if (this._flags.last_token.type === l.OPERATOR && u(this._options.operator_position, _) && this.allow_wrap_or_preserved_newline(e), e.text === `:` && this._flags.in_case) {
                this.print_token(e), this._flags.in_case = !1, this._flags.case_body = !0, this._tokens.peek().type === l.START_BLOCK ? (this._flags.case_block = !0, this._output.space_before_token = !0) : (this.indent(), this.print_newline(), this._flags.case_block = !1);
                return;
            }
            var i = !0, a = !0, o = !1;
            if (e.text === `:` ? this._flags.ternary_depth === 0 ? i = !1 : (--this._flags.ternary_depth, o = !0) : e.text === `?` && (this._flags.ternary_depth += 1), !n && !t && this._options.preserve_newlines && u(e.text, c)) {
                var d = e.text === `:`, f = d && o, p = d && !o;
                switch(this._options.operator_position){
                    case g.before_newline:
                        this._output.space_before_token = !p, this.print_token(e), (!d || f) && this.allow_wrap_or_preserved_newline(e), this._output.space_before_token = !0;
                        return;
                    case g.after_newline:
                        this._output.space_before_token = !0, !d || f ? this._tokens.peek().newlines ? this.print_newline(!1, !0) : this.allow_wrap_or_preserved_newline(e) : this._output.space_before_token = !1, this.print_token(e), this._output.space_before_token = !0;
                        return;
                    case g.preserve_newline:
                        p || this.allow_wrap_or_preserved_newline(e), i = !(this._output.just_added_newline() || p), this._output.space_before_token = i, this.print_token(e), this._output.space_before_token = !0;
                        return;
                }
            }
            if (t) {
                this.allow_wrap_or_preserved_newline(e), i = !1;
                var y = this._tokens.peek();
                a = y && u(y.type, [
                    l.WORD,
                    l.RESERVED
                ]);
            } else if (e.text === `...`) this.allow_wrap_or_preserved_newline(e), i = this._flags.last_token.type === l.START_BLOCK, a = !1;
            else if (u(e.text, [
                `--`,
                `++`,
                `!`,
                `~`
            ]) || n) {
                if ((this._flags.last_token.type === l.COMMA || this._flags.last_token.type === l.START_EXPR) && this.allow_wrap_or_preserved_newline(e), i = !1, a = !1, e.newlines && (e.text === `--` || e.text === `++` || e.text === `~`)) {
                    var b = m(this._flags.last_token, h) && e.newlines;
                    b && (this._previous_flags.if_block || this._previous_flags.else_block) && this.restore_mode(), this.print_newline(b, !0);
                }
                this._flags.last_token.text === `;` && S(this._flags.mode) && (i = !0), this._flags.last_token.type === l.RESERVED ? i = !0 : this._flags.last_token.type === l.END_EXPR ? i = !(this._flags.last_token.text === `]` && (e.text === `--` || e.text === `++`)) : this._flags.last_token.type === l.OPERATOR && (i = u(e.text, [
                    `--`,
                    `-`,
                    `++`,
                    `+`
                ]) && u(this._flags.last_token.text, [
                    `--`,
                    `-`,
                    `++`,
                    `+`
                ]), u(e.text, [
                    `+`,
                    `-`
                ]) && u(this._flags.last_token.text, [
                    `--`,
                    `++`
                ]) && (a = !0)), (this._flags.mode === v.BlockStatement && !this._flags.inline_frame || this._flags.mode === v.Statement) && (this._flags.last_token.text === `{` || this._flags.last_token.text === `;`) && this.print_newline();
            }
            this._output.space_before_token = this._output.space_before_token || i, this.print_token(e), this._output.space_before_token = a;
        }, T.prototype.handle_block_comment = function(e, t) {
            if (this._output.raw) {
                this._output.add_raw_token(e), e.directives && e.directives.preserve === `end` && (this._output.raw = this._options.test_output_raw);
                return;
            }
            if (e.directives) {
                this.print_newline(!1, t), this.print_token(e), e.directives.preserve === `start` && (this._output.raw = !0), this.print_newline(!1, !0);
                return;
            }
            if (!i.newline.test(e.text) && !e.newlines) {
                this._output.space_before_token = !0, this.print_token(e), this._output.space_before_token = !0;
                return;
            } else this.print_block_commment(e, t);
        }, T.prototype.print_block_commment = function(e, t) {
            var n = b(e.text), r, i = !1, a = !1, o = e.whitespace_before, s = o.length;
            if (this.print_newline(!1, t), this.print_token_line_indentation(e), this._output.add_token(n[0]), this.print_newline(!1, t), n.length > 1) {
                for(n = n.slice(1), i = C(n, `*`), a = w(n, o), i && (this._flags.alignment = 1), r = 0; r < n.length; r++)i ? (this.print_token_line_indentation(e), this._output.add_token(d(n[r]))) : a && n[r] ? (this.print_token_line_indentation(e), this._output.add_token(n[r].substring(s))) : (this._output.current_line.set_indent(-1), this._output.add_token(n[r])), this.print_newline(!1, t);
                this._flags.alignment = 0;
            }
        }, T.prototype.handle_comment = function(e, t) {
            e.newlines ? this.print_newline(!1, t) : this._output.trim(!0), this._output.space_before_token = !0, this.print_token(e), this.print_newline(!1, t);
        }, T.prototype.handle_dot = function(e) {
            this.start_of_statement(e) || this.handle_whitespace_and_comments(e, !0), this._flags.last_token.text.match(`^[0-9]+$`) && (this._output.space_before_token = !0), m(this._flags.last_token, h) ? this._output.space_before_token = !1 : this.allow_wrap_or_preserved_newline(e, this._flags.last_token.text === `)` && this._options.break_chained_methods), this._options.unindent_chained_methods && this._output.just_added_newline() && this.deindent(), this.print_token(e);
        }, T.prototype.handle_unknown = function(e, t) {
            this.print_token(e), e.text[e.text.length - 1] === `
` && this.print_newline(!1, t);
        }, T.prototype.handle_eof = function(e) {
            for(; this._flags.mode === v.Statement;)this.restore_mode();
            this.handle_whitespace_and_comments(e);
        }, t.exports.Beautifier = T;
    })), Nt = t(((e, t)=>{
        var n = Mt().Beautifier, r = Ct().Options;
        function i(e, t) {
            return new n(e, t).beautify();
        }
        t.exports = i, t.exports.defaultOptions = function() {
            return new r;
        };
    })), Pt = t(((e, t)=>{
        var n = St().Options;
        function r(e) {
            n.call(this, e, `css`), this.selector_separator_newline = this._get_boolean(`selector_separator_newline`, !0), this.newline_between_rules = this._get_boolean(`newline_between_rules`, !0);
            var t = this._get_boolean(`space_around_selector_separator`);
            this.space_around_combinator = this._get_boolean(`space_around_combinator`) || t;
            var r = this._get_selection_list(`brace_style`, [
                `collapse`,
                `expand`,
                `end-expand`,
                `none`,
                `preserve-inline`
            ]);
            this.brace_style = `collapse`;
            for(var i = 0; i < r.length; i++)r[i] === `expand` ? this.brace_style = r[i] : this.brace_style = `collapse`;
        }
        r.prototype = new n, t.exports.Options = r;
    })), Ft = t(((e, t)=>{
        var n = Pt().Options, r = yt().Output, i = wt().InputScanner, a = kt().Directives, o = new a(/\/\*/, /\*\//), s = /\r\n|[\r\n]/, c = /\r\n|[\r\n]/g, l = /\s/, u = /(?:\s|\n)+/g, d = /\/\*(?:[\s\S]*?)((?:\*\/)|$)/g, f = /\/\/(?:[^\n\r\u2028\u2029]*)/g;
        function p(e, t) {
            this._source_text = e || ``, this._options = new n(t), this._ch = null, this._input = null, this.NESTED_AT_RULE = {
                page: !0,
                "font-face": !0,
                keyframes: !0,
                media: !0,
                supports: !0,
                document: !0
            }, this.CONDITIONAL_GROUP_RULE = {
                media: !0,
                supports: !0,
                document: !0
            }, this.NON_SEMICOLON_NEWLINE_PROPERTY = [
                `grid-template-areas`,
                `grid-template`
            ];
        }
        p.prototype.eatString = function(e) {
            var t = ``;
            for(this._ch = this._input.next(); this._ch;){
                if (t += this._ch, this._ch === `\\`) t += this._input.next();
                else if (e.indexOf(this._ch) !== -1 || this._ch === `
`) break;
                this._ch = this._input.next();
            }
            return t;
        }, p.prototype.eatWhitespace = function(e) {
            for(var t = l.test(this._input.peek()), n = 0; l.test(this._input.peek());)this._ch = this._input.next(), e && this._ch === `
` && (n === 0 || n < this._options.max_preserve_newlines) && (n++, this._output.add_new_line(!0));
            return t;
        }, p.prototype.foundNestedPseudoClass = function() {
            for(var e = 0, t = 1, n = this._input.peek(t); n;){
                if (n === `{`) return !0;
                if (n === `(`) e += 1;
                else if (n === `)`) {
                    if (e === 0) return !1;
                    --e;
                } else if (n === `;` || n === `}`) return !1;
                t++, n = this._input.peek(t);
            }
            return !1;
        }, p.prototype.print_string = function(e) {
            this._output.set_indent(this._indentLevel), this._output.non_breaking_space = !0, this._output.add_token(e);
        }, p.prototype.preserveSingleSpace = function(e) {
            e && (this._output.space_before_token = !0);
        }, p.prototype.indent = function() {
            this._indentLevel++;
        }, p.prototype.outdent = function() {
            this._indentLevel > 0 && this._indentLevel--;
        }, p.prototype.beautify = function() {
            if (this._options.disabled) return this._source_text;
            var e = this._source_text, t = this._options.eol;
            t === `auto` && (t = `
`, e && s.test(e || ``) && (t = e.match(s)[0])), e = e.replace(c, `
`);
            var n = e.match(/^[\t ]*/)[0];
            this._output = new r(this._options, n), this._input = new i(e), this._indentLevel = 0, this._nestedLevel = 0, this._ch = null;
            for(var a = 0, p = !1, m = !1, h = !1, g = !1, _ = !1, v = this._ch, y = !1, b, x, S; b = this._input.read(u), x = b !== ``, S = v, this._ch = this._input.next(), this._ch === `\\` && this._input.hasNext() && (this._ch += this._input.next()), v = this._ch, this._ch;)if (this._ch === `/` && this._input.peek() === `*`) {
                this._output.add_new_line(), this._input.back();
                var C = this._input.read(d), w = o.get_directives(C);
                w && w.ignore === `start` && (C += o.readIgnored(this._input)), this.print_string(C), this.eatWhitespace(!0), this._output.add_new_line();
            } else if (this._ch === `/` && this._input.peek() === `/`) this._output.space_before_token = !0, this._input.back(), this.print_string(this._input.read(f)), this.eatWhitespace(!0);
            else if (this._ch === `$`) {
                this.preserveSingleSpace(x), this.print_string(this._ch);
                var T = this._input.peekUntilAfter(/[: ,;{}()[\]\/='"]/g);
                T.match(/[ :]$/) && (T = this.eatString(`: `).replace(/\s+$/, ``), this.print_string(T), this._output.space_before_token = !0), a === 0 && T.indexOf(`:`) !== -1 && (m = !0, this.indent());
            } else if (this._ch === `@`) if (this.preserveSingleSpace(x), this._input.peek() === `{`) this.print_string(this._ch + this.eatString(`}`));
            else {
                this.print_string(this._ch);
                var E = this._input.peekUntilAfter(/[: ,;{}()[\]\/='"]/g);
                E.match(/[ :]$/) && (E = this.eatString(`: `).replace(/\s+$/, ``), this.print_string(E), this._output.space_before_token = !0), a === 0 && E.indexOf(`:`) !== -1 ? (m = !0, this.indent()) : E in this.NESTED_AT_RULE ? (this._nestedLevel += 1, E in this.CONDITIONAL_GROUP_RULE && (h = !0)) : a === 0 && !m && (g = !0);
            }
            else if (this._ch === `#` && this._input.peek() === `{`) this.preserveSingleSpace(x), this.print_string(this._ch + this.eatString(`}`));
            else if (this._ch === `{`) m && (m = !1, this.outdent()), g = !1, h ? (h = !1, p = this._indentLevel >= this._nestedLevel) : p = this._indentLevel >= this._nestedLevel - 1, this._options.newline_between_rules && p && this._output.previous_line && this._output.previous_line.item(-1) !== `{` && this._output.ensure_empty_line_above(`/`, `,`), this._output.space_before_token = !0, this._options.brace_style === `expand` ? (this._output.add_new_line(), this.print_string(this._ch), this.indent(), this._output.set_indent(this._indentLevel)) : (S === `(` ? this._output.space_before_token = !1 : S !== `,` && this.indent(), this.print_string(this._ch)), this.eatWhitespace(!0), this._output.add_new_line();
            else if (this._ch === `}`) this.outdent(), this._output.add_new_line(), S === `{` && this._output.trim(!0), m &&= (this.outdent(), !1), this.print_string(this._ch), p = !1, this._nestedLevel && this._nestedLevel--, this.eatWhitespace(!0), this._output.add_new_line(), this._options.newline_between_rules && !this._output.just_added_blankline() && this._input.peek() !== `}` && this._output.add_new_line(!0), this._input.peek() === `)` && (this._output.trim(!0), this._options.brace_style === `expand` && this._output.add_new_line(!0));
            else if (this._ch === `:`) {
                for(var D = 0; D < this.NON_SEMICOLON_NEWLINE_PROPERTY.length; D++)if (this._input.lookBack(this.NON_SEMICOLON_NEWLINE_PROPERTY[D])) {
                    y = !0;
                    break;
                }
                (p || h) && !(this._input.lookBack(`&`) || this.foundNestedPseudoClass()) && !this._input.lookBack(`(`) && !g && a === 0 ? (this.print_string(`:`), m || (m = !0, this._output.space_before_token = !0, this.eatWhitespace(!0), this.indent())) : (this._input.lookBack(` `) && (this._output.space_before_token = !0), this._input.peek() === `:` ? (this._ch = this._input.next(), this.print_string(`::`)) : this.print_string(`:`));
            } else if (this._ch === `"` || this._ch === `'`) {
                var O = S === `"` || S === `'`;
                this.preserveSingleSpace(O || x), this.print_string(this._ch + this.eatString(this._ch)), this.eatWhitespace(!0);
            } else if (this._ch === `;`) y = !1, a === 0 ? (m &&= (this.outdent(), !1), g = !1, this.print_string(this._ch), this.eatWhitespace(!0), this._input.peek() !== `/` && this._output.add_new_line()) : (this.print_string(this._ch), this.eatWhitespace(!0), this._output.space_before_token = !0);
            else if (this._ch === `(`) if (this._input.lookBack(`url`)) this.print_string(this._ch), this.eatWhitespace(), a++, this.indent(), this._ch = this._input.next(), this._ch === `)` || this._ch === `"` || this._ch === `'` ? this._input.back() : this._ch && (this.print_string(this._ch + this.eatString(`)`)), a && (a--, this.outdent()));
            else {
                var k = !1;
                this._input.lookBack(`with`) && (k = !0), this.preserveSingleSpace(x || k), this.print_string(this._ch), m && S === `$` && this._options.selector_separator_newline ? (this._output.add_new_line(), _ = !0) : (this.eatWhitespace(), a++, this.indent());
            }
            else if (this._ch === `)`) a && (a--, this.outdent()), _ && this._input.peek() === `;` && this._options.selector_separator_newline && (_ = !1, this.outdent(), this._output.add_new_line()), this.print_string(this._ch);
            else if (this._ch === `,`) this.print_string(this._ch), this.eatWhitespace(!0), this._options.selector_separator_newline && (!m || _) && a === 0 && !g ? this._output.add_new_line() : this._output.space_before_token = !0;
            else if ((this._ch === `>` || this._ch === `+` || this._ch === `~`) && !m && a === 0) this._options.space_around_combinator ? (this._output.space_before_token = !0, this.print_string(this._ch), this._output.space_before_token = !0) : (this.print_string(this._ch), this.eatWhitespace(), this._ch && l.test(this._ch) && (this._ch = ``));
            else if (this._ch === `]`) this.print_string(this._ch);
            else if (this._ch === `[`) this.preserveSingleSpace(x), this.print_string(this._ch);
            else if (this._ch === `=`) this.eatWhitespace(), this.print_string(`=`), l.test(this._ch) && (this._ch = ``);
            else if (this._ch === `!` && !this._input.lookBack(`\\`)) this._output.space_before_token = !0, this.print_string(this._ch);
            else {
                var A = S === `"` || S === `'`;
                this.preserveSingleSpace(A || x), this.print_string(this._ch), !this._output.just_added_newline() && this._input.peek() === `
` && y && this._output.add_new_line();
            }
            return this._output.get_code(t);
        }, t.exports.Beautifier = p;
    })), It = t(((e, t)=>{
        var n = Ft().Beautifier, r = Pt().Options;
        function i(e, t) {
            return new n(e, t).beautify();
        }
        t.exports = i, t.exports.defaultOptions = function() {
            return new r;
        };
    })), Lt = t(((e, t)=>{
        var n = St().Options;
        function r(e) {
            n.call(this, e, `html`), this.templating.length === 1 && this.templating[0] === `auto` && (this.templating = [
                `django`,
                `erb`,
                `handlebars`,
                `php`
            ]), this.indent_inner_html = this._get_boolean(`indent_inner_html`), this.indent_body_inner_html = this._get_boolean(`indent_body_inner_html`, !0), this.indent_head_inner_html = this._get_boolean(`indent_head_inner_html`, !0), this.indent_handlebars = this._get_boolean(`indent_handlebars`, !0), this.wrap_attributes = this._get_selection(`wrap_attributes`, [
                `auto`,
                `force`,
                `force-aligned`,
                `force-expand-multiline`,
                `aligned-multiple`,
                `preserve`,
                `preserve-aligned`
            ]), this.wrap_attributes_min_attrs = this._get_number(`wrap_attributes_min_attrs`, 2), this.wrap_attributes_indent_size = this._get_number(`wrap_attributes_indent_size`, this.indent_size), this.extra_liners = this._get_array(`extra_liners`, [
                `head`,
                `body`,
                `/html`
            ]), this.inline = this._get_array(`inline`, `a.abbr.area.audio.b.bdi.bdo.br.button.canvas.cite.code.data.datalist.del.dfn.em.embed.i.iframe.img.input.ins.kbd.keygen.label.map.mark.math.meter.noscript.object.output.progress.q.ruby.s.samp.select.small.span.strong.sub.sup.svg.template.textarea.time.u.var.video.wbr.text.acronym.big.strike.tt`.split(`.`)), this.inline_custom_elements = this._get_boolean(`inline_custom_elements`, !0), this.void_elements = this._get_array(`void_elements`, [
                `area`,
                `base`,
                `br`,
                `col`,
                `embed`,
                `hr`,
                `img`,
                `input`,
                `keygen`,
                `link`,
                `menuitem`,
                `meta`,
                `param`,
                `source`,
                `track`,
                `wbr`,
                `!doctype`,
                `?xml`,
                `basefont`,
                `isindex`
            ]), this.unformatted = this._get_array(`unformatted`, []), this.content_unformatted = this._get_array(`content_unformatted`, [
                `pre`,
                `textarea`
            ]), this.unformatted_content_delimiter = this._get_characters(`unformatted_content_delimiter`), this.indent_scripts = this._get_selection(`indent_scripts`, [
                `normal`,
                `keep`,
                `separate`
            ]);
        }
        r.prototype = new n, t.exports.Options = r;
    })), Rt = t(((e, t)=>{
        var n = Ot().Tokenizer, r = Ot().TOKEN, i = kt().Directives, a = At().TemplatablePattern, o = Et().Pattern, s = {
            TAG_OPEN: `TK_TAG_OPEN`,
            TAG_CLOSE: `TK_TAG_CLOSE`,
            CONTROL_FLOW_OPEN: `TK_CONTROL_FLOW_OPEN`,
            CONTROL_FLOW_CLOSE: `TK_CONTROL_FLOW_CLOSE`,
            ATTRIBUTE: `TK_ATTRIBUTE`,
            EQUALS: `TK_EQUALS`,
            VALUE: `TK_VALUE`,
            COMMENT: `TK_COMMENT`,
            TEXT: `TK_TEXT`,
            UNKNOWN: `TK_UNKNOWN`,
            START: r.START,
            RAW: r.RAW,
            EOF: r.EOF
        }, c = new i(/<\!--/, /-->/), l = function(e, t) {
            n.call(this, e, t), this._current_tag_name = ``;
            var r = new a(this._input).read_options(this._options), i = new o(this._input);
            if (this.__patterns = {
                word: r.until(/[\n\r\t <]/),
                word_control_flow_close_excluded: r.until(/[\n\r\t <}]/),
                single_quote: r.until_after(/'/),
                double_quote: r.until_after(/"/),
                attribute: r.until(/[\n\r\t =>]|\/>/),
                element_name: r.until(/[\n\r\t >\/]/),
                angular_control_flow_start: i.matching(/\@[a-zA-Z]+[^({]*[({]/),
                handlebars_comment: i.starting_with(/{{!--/).until_after(/--}}/),
                handlebars: i.starting_with(/{{/).until_after(/}}/),
                handlebars_open: i.until(/[\n\r\t }]/),
                handlebars_raw_close: i.until(/}}/),
                comment: i.starting_with(/<!--/).until_after(/-->/),
                cdata: i.starting_with(/<!\[CDATA\[/).until_after(/]]>/),
                conditional_comment: i.starting_with(/<!\[/).until_after(/]>/),
                processing: i.starting_with(/<\?/).until_after(/\?>/)
            }, this._options.indent_handlebars && (this.__patterns.word = this.__patterns.word.exclude(`handlebars`), this.__patterns.word_control_flow_close_excluded = this.__patterns.word_control_flow_close_excluded.exclude(`handlebars`)), this._unformatted_content_delimiter = null, this._options.unformatted_content_delimiter) {
                var s = this._input.get_literal_regexp(this._options.unformatted_content_delimiter);
                this.__patterns.unformatted_content_delimiter = i.matching(s).until_after(s);
            }
        };
        l.prototype = new n, l.prototype._is_comment = function(e) {
            return !1;
        }, l.prototype._is_opening = function(e) {
            return e.type === s.TAG_OPEN || e.type === s.CONTROL_FLOW_OPEN;
        }, l.prototype._is_closing = function(e, t) {
            return e.type === s.TAG_CLOSE && t && ((e.text === `>` || e.text === `/>`) && t.text[0] === `<` || e.text === `}}` && t.text[0] === `{` && t.text[1] === `{`) || e.type === s.CONTROL_FLOW_CLOSE && e.text === `}` && t.text.endsWith(`{`);
        }, l.prototype._reset = function() {
            this._current_tag_name = ``;
        }, l.prototype._get_next_token = function(e, t) {
            var n = null;
            this._readWhitespace();
            var r = this._input.peek();
            return r === null ? this._create_token(s.EOF, ``) : (n ||= this._read_open_handlebars(r, t), n ||= this._read_attribute(r, e, t), n ||= this._read_close(r, t), n ||= this._read_script_and_style(r, e), n ||= this._read_control_flows(r, t), n ||= this._read_raw_content(r, e, t), n ||= this._read_content_word(r, t), n ||= this._read_comment_or_cdata(r), n ||= this._read_processing(r), n ||= this._read_open(r, t), n ||= this._create_token(s.UNKNOWN, this._input.next()), n);
        }, l.prototype._read_comment_or_cdata = function(e) {
            var t = null, n = null, r = null;
            return e === `<` && (this._input.peek(1) === `!` && (n = this.__patterns.comment.read(), n ? (r = c.get_directives(n), r && r.ignore === `start` && (n += c.readIgnored(this._input))) : n = this.__patterns.cdata.read()), n && (t = this._create_token(s.COMMENT, n), t.directives = r)), t;
        }, l.prototype._read_processing = function(e) {
            var t = null, n = null, r = null;
            if (e === `<`) {
                var i = this._input.peek(1);
                (i === `!` || i === `?`) && (n = this.__patterns.conditional_comment.read(), n ||= this.__patterns.processing.read()), n && (t = this._create_token(s.COMMENT, n), t.directives = r);
            }
            return t;
        }, l.prototype._read_open = function(e, t) {
            var n = null, r = null;
            return (!t || t.type === s.CONTROL_FLOW_OPEN) && e === `<` && (n = this._input.next(), this._input.peek() === `/` && (n += this._input.next()), n += this.__patterns.element_name.read(), r = this._create_token(s.TAG_OPEN, n)), r;
        }, l.prototype._read_open_handlebars = function(e, t) {
            var n = null, r = null;
            return (!t || t.type === s.CONTROL_FLOW_OPEN) && (this._options.templating.includes(`angular`) || this._options.indent_handlebars) && e === `{` && this._input.peek(1) === `{` && (this._options.indent_handlebars && this._input.peek(2) === `!` ? (n = this.__patterns.handlebars_comment.read(), n ||= this.__patterns.handlebars.read(), r = this._create_token(s.COMMENT, n)) : (n = this.__patterns.handlebars_open.read(), r = this._create_token(s.TAG_OPEN, n))), r;
        }, l.prototype._read_control_flows = function(e, t) {
            var n = ``, r = null;
            if (!this._options.templating.includes(`angular`)) return r;
            if (e === `@`) {
                if (n = this.__patterns.angular_control_flow_start.read(), n === ``) return r;
                for(var i = +!!n.endsWith(`(`), a = 0; !(n.endsWith(`{`) && i === a);){
                    var o = this._input.next();
                    if (o === null) break;
                    o === `(` ? i++ : o === `)` && a++, n += o;
                }
                r = this._create_token(s.CONTROL_FLOW_OPEN, n);
            } else e === `}` && t && t.type === s.CONTROL_FLOW_OPEN && (n = this._input.next(), r = this._create_token(s.CONTROL_FLOW_CLOSE, n));
            return r;
        }, l.prototype._read_close = function(e, t) {
            var n = null, r = null;
            return t && t.type === s.TAG_OPEN && (t.text[0] === `<` && (e === `>` || e === `/` && this._input.peek(1) === `>`) ? (n = this._input.next(), e === `/` && (n += this._input.next()), r = this._create_token(s.TAG_CLOSE, n)) : t.text[0] === `{` && e === `}` && this._input.peek(1) === `}` && (this._input.next(), this._input.next(), r = this._create_token(s.TAG_CLOSE, `}}`))), r;
        }, l.prototype._read_attribute = function(e, t, n) {
            var r = null, i = ``;
            if (n && n.text[0] === `<`) if (e === `=`) r = this._create_token(s.EQUALS, this._input.next());
            else if (e === `"` || e === `'`) {
                var a = this._input.next();
                e === `"` ? a += this.__patterns.double_quote.read() : a += this.__patterns.single_quote.read(), r = this._create_token(s.VALUE, a);
            } else i = this.__patterns.attribute.read(), i && (r = t.type === s.EQUALS ? this._create_token(s.VALUE, i) : this._create_token(s.ATTRIBUTE, i));
            return r;
        }, l.prototype._is_content_unformatted = function(e) {
            return this._options.void_elements.indexOf(e) === -1 && (this._options.content_unformatted.indexOf(e) !== -1 || this._options.unformatted.indexOf(e) !== -1);
        }, l.prototype._read_raw_content = function(e, t, n) {
            var r = ``;
            if (n && n.text[0] === `{`) r = this.__patterns.handlebars_raw_close.read();
            else if (t.type === s.TAG_CLOSE && t.opened.text[0] === `<` && t.text[0] !== `/`) {
                var i = t.opened.text.substr(1).toLowerCase();
                this._is_content_unformatted(i) && (r = this._input.readUntil(RegExp(`</` + i + `[\\n\\r\\t ]*?>`, `ig`)));
            }
            return r ? this._create_token(s.TEXT, r) : null;
        }, l.prototype._read_script_and_style = function(e, t) {
            if (t.type === s.TAG_CLOSE && t.opened.text[0] === `<` && t.text[0] !== `/`) {
                var n = t.opened.text.substr(1).toLowerCase();
                if (n === `script` || n === `style`) {
                    var r = this._read_comment_or_cdata(e);
                    if (r) return r.type = s.TEXT, r;
                    var i = this._input.readUntil(RegExp(`</` + n + `[\\n\\r\\t ]*?>`, `ig`));
                    if (i) return this._create_token(s.TEXT, i);
                }
            }
            return null;
        }, l.prototype._read_content_word = function(e, t) {
            var n = ``;
            return this._options.unformatted_content_delimiter && e === this._options.unformatted_content_delimiter[0] && (n = this.__patterns.unformatted_content_delimiter.read()), n ||= t && t.type === s.CONTROL_FLOW_OPEN ? this.__patterns.word_control_flow_close_excluded.read() : this.__patterns.word.read(), n ? this._create_token(s.TEXT, n) : null;
        }, t.exports.Tokenizer = l, t.exports.TOKEN = s;
    })), zt = t(((e, t)=>{
        var n = Lt().Options, r = yt().Output, i = Rt().Tokenizer, a = Rt().TOKEN, o = /\r\n|[\r\n]/, s = /\r\n|[\r\n]/g, c = function(e, t) {
            this.indent_level = 0, this.alignment_size = 0, this.max_preserve_newlines = e.max_preserve_newlines, this.preserve_newlines = e.preserve_newlines, this._output = new r(e, t);
        };
        c.prototype.current_line_has_match = function(e) {
            return this._output.current_line.has_match(e);
        }, c.prototype.set_space_before_token = function(e, t) {
            this._output.space_before_token = e, this._output.non_breaking_space = t;
        }, c.prototype.set_wrap_point = function() {
            this._output.set_indent(this.indent_level, this.alignment_size), this._output.set_wrap_point();
        }, c.prototype.add_raw_token = function(e) {
            this._output.add_raw_token(e);
        }, c.prototype.print_preserved_newlines = function(e) {
            var t = 0;
            e.type !== a.TEXT && e.previous.type !== a.TEXT && (t = +!!e.newlines), this.preserve_newlines && (t = e.newlines < this.max_preserve_newlines + 1 ? e.newlines : this.max_preserve_newlines + 1);
            for(var n = 0; n < t; n++)this.print_newline(n > 0);
            return t !== 0;
        }, c.prototype.traverse_whitespace = function(e) {
            return e.whitespace_before || e.newlines ? (this.print_preserved_newlines(e) || (this._output.space_before_token = !0), !0) : !1;
        }, c.prototype.previous_token_wrapped = function() {
            return this._output.previous_token_wrapped;
        }, c.prototype.print_newline = function(e) {
            this._output.add_new_line(e);
        }, c.prototype.print_token = function(e) {
            e.text && (this._output.set_indent(this.indent_level, this.alignment_size), this._output.add_token(e.text));
        }, c.prototype.indent = function() {
            this.indent_level++;
        }, c.prototype.deindent = function() {
            this.indent_level > 0 && (this.indent_level--, this._output.set_indent(this.indent_level, this.alignment_size));
        }, c.prototype.get_full_indent = function(e) {
            return e = this.indent_level + (e || 0), e < 1 ? `` : this._output.get_indent_string(e);
        };
        var l = function(e) {
            for(var t = null, n = e.next; n.type !== a.EOF && e.closed !== n;){
                if (n.type === a.ATTRIBUTE && n.text === `type`) {
                    n.next && n.next.type === a.EQUALS && n.next.next && n.next.next.type === a.VALUE && (t = n.next.next.text);
                    break;
                }
                n = n.next;
            }
            return t;
        }, u = function(e, t) {
            var n = null, r = null;
            return t.closed ? (e === `script` ? n = `text/javascript` : e === `style` && (n = `text/css`), n = l(t) || n, n.search(`text/css`) > -1 ? r = `css` : n.search(/module|((text|application|dojo)\/(x-)?(javascript|ecmascript|jscript|livescript|(ld\+)?json|method|aspect))/) > -1 ? r = `javascript` : n.search(/(text|application|dojo)\/(x-)?(html)/) > -1 ? r = `html` : n.search(/test\/null/) > -1 && (r = `null`), r) : null;
        };
        function d(e, t) {
            return t.indexOf(e) !== -1;
        }
        function f(e, t, n) {
            this.parent = e || null, this.tag = t ? t.tag_name : ``, this.indent_level = n || 0, this.parser_token = t || null;
        }
        function p(e) {
            this._printer = e, this._current_frame = null;
        }
        p.prototype.get_parser_token = function() {
            return this._current_frame ? this._current_frame.parser_token : null;
        }, p.prototype.record_tag = function(e) {
            var t = new f(this._current_frame, e, this._printer.indent_level);
            this._current_frame = t;
        }, p.prototype._try_pop_frame = function(e) {
            var t = null;
            return e && (t = e.parser_token, this._printer.indent_level = e.indent_level, this._current_frame = e.parent), t;
        }, p.prototype._get_frame = function(e, t) {
            for(var n = this._current_frame; n && e.indexOf(n.tag) === -1;){
                if (t && t.indexOf(n.tag) !== -1) {
                    n = null;
                    break;
                }
                n = n.parent;
            }
            return n;
        }, p.prototype.try_pop = function(e, t) {
            var n = this._get_frame([
                e
            ], t);
            return this._try_pop_frame(n);
        }, p.prototype.indent_to_tag = function(e) {
            var t = this._get_frame(e);
            t && (this._printer.indent_level = t.indent_level);
        };
        function m(e, t, r, i) {
            this._source_text = e || ``, t ||= {}, this._js_beautify = r, this._css_beautify = i, this._tag_stack = null;
            var a = new n(t, `html`);
            this._options = a, this._is_wrap_attributes_force = this._options.wrap_attributes.substr(0, 5) === `force`, this._is_wrap_attributes_force_expand_multiline = this._options.wrap_attributes === `force-expand-multiline`, this._is_wrap_attributes_force_aligned = this._options.wrap_attributes === `force-aligned`, this._is_wrap_attributes_aligned_multiple = this._options.wrap_attributes === `aligned-multiple`, this._is_wrap_attributes_preserve = this._options.wrap_attributes.substr(0, 8) === `preserve`, this._is_wrap_attributes_preserve_aligned = this._options.wrap_attributes === `preserve-aligned`;
        }
        m.prototype.beautify = function() {
            if (this._options.disabled) return this._source_text;
            var e = this._source_text, t = this._options.eol;
            this._options.eol === `auto` && (t = `
`, e && o.test(e) && (t = e.match(o)[0])), e = e.replace(s, `
`);
            var n = e.match(/^[\t ]*/)[0], r = {
                text: ``,
                type: ``
            }, l = new h(this._options), u = new c(this._options, n), d = new i(e, this._options).tokenize();
            this._tag_stack = new p(u);
            for(var f = null, m = d.next(); m.type !== a.EOF;)m.type === a.TAG_OPEN || m.type === a.COMMENT ? (f = this._handle_tag_open(u, m, l, r, d), l = f) : m.type === a.ATTRIBUTE || m.type === a.EQUALS || m.type === a.VALUE || m.type === a.TEXT && !l.tag_complete ? f = this._handle_inside_tag(u, m, l, r) : m.type === a.TAG_CLOSE ? f = this._handle_tag_close(u, m, l) : m.type === a.TEXT ? f = this._handle_text(u, m, l) : m.type === a.CONTROL_FLOW_OPEN ? f = this._handle_control_flow_open(u, m) : m.type === a.CONTROL_FLOW_CLOSE ? f = this._handle_control_flow_close(u, m) : u.add_raw_token(m), r = f, m = d.next();
            return u._output.get_code(t);
        }, m.prototype._handle_control_flow_open = function(e, t) {
            var n = {
                text: t.text,
                type: t.type
            };
            return e.set_space_before_token(t.newlines || t.whitespace_before !== ``, !0), t.newlines ? e.print_preserved_newlines(t) : e.set_space_before_token(t.newlines || t.whitespace_before !== ``, !0), e.print_token(t), e.indent(), n;
        }, m.prototype._handle_control_flow_close = function(e, t) {
            var n = {
                text: t.text,
                type: t.type
            };
            return e.deindent(), t.newlines ? e.print_preserved_newlines(t) : e.set_space_before_token(t.newlines || t.whitespace_before !== ``, !0), e.print_token(t), n;
        }, m.prototype._handle_tag_close = function(e, t, n) {
            var r = {
                text: t.text,
                type: t.type
            };
            return e.alignment_size = 0, n.tag_complete = !0, e.set_space_before_token(t.newlines || t.whitespace_before !== ``, !0), n.is_unformatted ? e.add_raw_token(t) : (n.tag_start_char === `<` && (e.set_space_before_token(t.text[0] === `/`, !0), this._is_wrap_attributes_force_expand_multiline && n.has_wrapped_attrs && e.print_newline(!1)), e.print_token(t)), n.indent_content && !(n.is_unformatted || n.is_content_unformatted) && (e.indent(), n.indent_content = !1), !n.is_inline_element && !(n.is_unformatted || n.is_content_unformatted) && e.set_wrap_point(), r;
        }, m.prototype._handle_inside_tag = function(e, t, n, r) {
            var i = n.has_wrapped_attrs, o = {
                text: t.text,
                type: t.type
            };
            return e.set_space_before_token(t.newlines || t.whitespace_before !== ``, !0), n.is_unformatted ? e.add_raw_token(t) : n.tag_start_char === `{` && t.type === a.TEXT ? e.print_preserved_newlines(t) ? (t.newlines = 0, e.add_raw_token(t)) : e.print_token(t) : (t.type === a.ATTRIBUTE ? e.set_space_before_token(!0) : (t.type === a.EQUALS || t.type === a.VALUE && t.previous.type === a.EQUALS) && e.set_space_before_token(!1), t.type === a.ATTRIBUTE && n.tag_start_char === `<` && ((this._is_wrap_attributes_preserve || this._is_wrap_attributes_preserve_aligned) && (e.traverse_whitespace(t), i ||= t.newlines !== 0), this._is_wrap_attributes_force && n.attr_count >= this._options.wrap_attributes_min_attrs && (r.type !== a.TAG_OPEN || this._is_wrap_attributes_force_expand_multiline) && (e.print_newline(!1), i = !0)), e.print_token(t), i ||= e.previous_token_wrapped(), n.has_wrapped_attrs = i), o;
        }, m.prototype._handle_text = function(e, t, n) {
            var r = {
                text: t.text,
                type: `TK_CONTENT`
            };
            return n.custom_beautifier_name ? this._print_custom_beatifier_text(e, t, n) : n.is_unformatted || n.is_content_unformatted ? e.add_raw_token(t) : (e.traverse_whitespace(t), e.print_token(t)), r;
        }, m.prototype._print_custom_beatifier_text = function(e, t, n) {
            var r = this;
            if (t.text !== ``) {
                var i = t.text, a, o = 1, s = ``, c = ``;
                n.custom_beautifier_name === `javascript` && typeof this._js_beautify == `function` ? a = this._js_beautify : n.custom_beautifier_name === `css` && typeof this._css_beautify == `function` ? a = this._css_beautify : n.custom_beautifier_name === `html` && (a = function(e, t) {
                    return new m(e, t, r._js_beautify, r._css_beautify).beautify();
                }), this._options.indent_scripts === `keep` ? o = 0 : this._options.indent_scripts === `separate` && (o = -e.indent_level);
                var l = e.get_full_indent(o);
                if (i = i.replace(/\n[ \t]*$/, ``), n.custom_beautifier_name !== `html` && i[0] === `<` && i.match(/^(<!--|<!\[CDATA\[)/)) {
                    var u = /^(<!--[^\n]*|<!\[CDATA\[)(\n?)([ \t\n]*)([\s\S]*)(-->|]]>)$/.exec(i);
                    if (!u) {
                        e.add_raw_token(t);
                        return;
                    }
                    s = l + u[1] + `
`, i = u[4], u[5] && (c = l + u[5]), i = i.replace(/\n[ \t]*$/, ``), (u[2] || u[3].indexOf(`
`) !== -1) && (u = u[3].match(/[ \t]+$/), u && (t.whitespace_before = u[0]));
                }
                if (i) if (a) {
                    var d = function() {
                        this.eol = `
`;
                    };
                    d.prototype = this._options.raw_options;
                    var f = new d;
                    i = a(l + i, f);
                } else {
                    var p = t.whitespace_before;
                    p && (i = i.replace(RegExp(`
(` + p + `)?`, `g`), `
`)), i = l + i.replace(/\n/g, `
` + l);
                }
                s && (i = i ? s + i + `
` + c : s + c), e.print_newline(!1), i && (t.text = i, t.whitespace_before = ``, t.newlines = 0, e.add_raw_token(t), e.print_newline(!0));
            }
        }, m.prototype._handle_tag_open = function(e, t, n, r, i) {
            var o = this._get_tag_open_token(t);
            if ((n.is_unformatted || n.is_content_unformatted) && !n.is_empty_element && t.type === a.TAG_OPEN && !o.is_start_tag ? (e.add_raw_token(t), o.start_tag_token = this._tag_stack.try_pop(o.tag_name)) : (e.traverse_whitespace(t), this._set_tag_position(e, t, o, n, r), o.is_inline_element || e.set_wrap_point(), e.print_token(t)), o.is_start_tag && this._is_wrap_attributes_force) {
                var s = 0, c;
                do c = i.peek(s), c.type === a.ATTRIBUTE && (o.attr_count += 1), s += 1;
                while (c.type !== a.EOF && c.type !== a.TAG_CLOSE);
            }
            return (this._is_wrap_attributes_force_aligned || this._is_wrap_attributes_aligned_multiple || this._is_wrap_attributes_preserve_aligned) && (o.alignment_size = t.text.length + 1), !o.tag_complete && !o.is_unformatted && (e.alignment_size = o.alignment_size), o;
        };
        var h = function(e, t, n) {
            if (this.parent = t || null, this.text = ``, this.type = `TK_TAG_OPEN`, this.tag_name = ``, this.is_inline_element = !1, this.is_unformatted = !1, this.is_content_unformatted = !1, this.is_empty_element = !1, this.is_start_tag = !1, this.is_end_tag = !1, this.indent_content = !1, this.multiline_content = !1, this.custom_beautifier_name = null, this.start_tag_token = null, this.attr_count = 0, this.has_wrapped_attrs = !1, this.alignment_size = 0, this.tag_complete = !1, this.tag_start_char = ``, this.tag_check = ``, !n) this.tag_complete = !0;
            else {
                var r;
                this.tag_start_char = n.text[0], this.text = n.text, this.tag_start_char === `<` ? (r = n.text.match(/^<([^\s>]*)/), this.tag_check = r ? r[1] : ``) : (r = n.text.match(/^{{~?(?:[\^]|#\*?)?([^\s}]+)/), this.tag_check = r ? r[1] : ``, (n.text.startsWith(`{{#>`) || n.text.startsWith(`{{~#>`)) && this.tag_check[0] === `>` && (this.tag_check === `>` && n.next !== null ? this.tag_check = n.next.text.split(` `)[0] : this.tag_check = n.text.split(`>`)[1])), this.tag_check = this.tag_check.toLowerCase(), n.type === a.COMMENT && (this.tag_complete = !0), this.is_start_tag = this.tag_check.charAt(0) !== `/`, this.tag_name = this.is_start_tag ? this.tag_check : this.tag_check.substr(1), this.is_end_tag = !this.is_start_tag || n.closed && n.closed.text === `/>`;
                var i = 2;
                this.tag_start_char === `{` && this.text.length >= 3 && this.text.charAt(2) === `~` && (i = 3), this.is_end_tag = this.is_end_tag || this.tag_start_char === `{` && (!e.indent_handlebars || this.text.length < 3 || /[^#\^]/.test(this.text.charAt(i)));
            }
        };
        m.prototype._get_tag_open_token = function(e) {
            var t = new h(this._options, this._tag_stack.get_parser_token(), e);
            return t.alignment_size = this._options.wrap_attributes_indent_size, t.is_end_tag = t.is_end_tag || d(t.tag_check, this._options.void_elements), t.is_empty_element = t.tag_complete || t.is_start_tag && t.is_end_tag, t.is_unformatted = !t.tag_complete && d(t.tag_check, this._options.unformatted), t.is_content_unformatted = !t.is_empty_element && d(t.tag_check, this._options.content_unformatted), t.is_inline_element = d(t.tag_name, this._options.inline) || this._options.inline_custom_elements && t.tag_name.includes(`-`) || t.tag_start_char === `{`, t;
        }, m.prototype._set_tag_position = function(e, t, n, r, i) {
            if (n.is_empty_element || (n.is_end_tag ? n.start_tag_token = this._tag_stack.try_pop(n.tag_name) : (this._do_optional_end_element(n) && (n.is_inline_element || e.print_newline(!1)), this._tag_stack.record_tag(n), (n.tag_name === `script` || n.tag_name === `style`) && !(n.is_unformatted || n.is_content_unformatted) && (n.custom_beautifier_name = u(n.tag_check, t)))), d(n.tag_check, this._options.extra_liners) && (e.print_newline(!1), e._output.just_added_blankline() || e.print_newline(!0)), n.is_empty_element) n.tag_start_char === `{` && n.tag_check === `else` && (this._tag_stack.indent_to_tag([
                `if`,
                `unless`,
                `each`
            ]), n.indent_content = !0, e.current_line_has_match(/{{#if/) || e.print_newline(!1)), n.tag_name === `!--` && i.type === a.TAG_CLOSE && r.is_end_tag && n.text.indexOf(`
`) === -1 || (n.is_inline_element || n.is_unformatted || e.print_newline(!1), this._calcluate_parent_multiline(e, n));
            else if (n.is_end_tag) {
                var o = !1;
                o = n.start_tag_token && n.start_tag_token.multiline_content, o ||= !n.is_inline_element && !(r.is_inline_element || r.is_unformatted) && !(i.type === a.TAG_CLOSE && n.start_tag_token === r) && i.type !== `TK_CONTENT`, (n.is_content_unformatted || n.is_unformatted) && (o = !1), o && e.print_newline(!1);
            } else n.indent_content = !n.custom_beautifier_name, n.tag_start_char === `<` && (n.tag_name === `html` ? n.indent_content = this._options.indent_inner_html : n.tag_name === `head` ? n.indent_content = this._options.indent_head_inner_html : n.tag_name === `body` && (n.indent_content = this._options.indent_body_inner_html)), !(n.is_inline_element || n.is_unformatted) && (i.type !== `TK_CONTENT` || n.is_content_unformatted) && e.print_newline(!1), this._calcluate_parent_multiline(e, n);
        }, m.prototype._calcluate_parent_multiline = function(e, t) {
            t.parent && e._output.just_added_newline() && !((t.is_inline_element || t.is_unformatted) && t.parent.is_inline_element) && (t.parent.multiline_content = !0);
        };
        var g = `address.article.aside.blockquote.details.div.dl.fieldset.figcaption.figure.footer.form.h1.h2.h3.h4.h5.h6.header.hr.main.menu.nav.ol.p.pre.section.table.ul`.split(`.`), _ = [
            `a`,
            `audio`,
            `del`,
            `ins`,
            `map`,
            `noscript`,
            `video`
        ];
        m.prototype._do_optional_end_element = function(e) {
            var t = null;
            if (!(e.is_empty_element || !e.is_start_tag || !e.parent)) {
                if (e.tag_name === `body`) t ||= this._tag_stack.try_pop(`head`);
                else if (e.tag_name === `li`) t ||= this._tag_stack.try_pop(`li`, [
                    `ol`,
                    `ul`,
                    `menu`
                ]);
                else if (e.tag_name === `dd` || e.tag_name === `dt`) t ||= this._tag_stack.try_pop(`dt`, [
                    `dl`
                ]), t ||= this._tag_stack.try_pop(`dd`, [
                    `dl`
                ]);
                else if (e.parent.tag_name === `p` && g.indexOf(e.tag_name) !== -1) {
                    var n = e.parent.parent;
                    (!n || _.indexOf(n.tag_name) === -1) && (t ||= this._tag_stack.try_pop(`p`));
                } else e.tag_name === `rp` || e.tag_name === `rt` ? (t ||= this._tag_stack.try_pop(`rt`, [
                    `ruby`,
                    `rtc`
                ]), t ||= this._tag_stack.try_pop(`rp`, [
                    `ruby`,
                    `rtc`
                ])) : e.tag_name === `optgroup` ? t ||= this._tag_stack.try_pop(`optgroup`, [
                    `select`
                ]) : e.tag_name === `option` ? t ||= this._tag_stack.try_pop(`option`, [
                    `select`,
                    `datalist`,
                    `optgroup`
                ]) : e.tag_name === `colgroup` ? t ||= this._tag_stack.try_pop(`caption`, [
                    `table`
                ]) : e.tag_name === `thead` ? (t ||= this._tag_stack.try_pop(`caption`, [
                    `table`
                ]), t ||= this._tag_stack.try_pop(`colgroup`, [
                    `table`
                ])) : e.tag_name === `tbody` || e.tag_name === `tfoot` ? (t ||= this._tag_stack.try_pop(`caption`, [
                    `table`
                ]), t ||= this._tag_stack.try_pop(`colgroup`, [
                    `table`
                ]), t ||= this._tag_stack.try_pop(`thead`, [
                    `table`
                ]), t ||= this._tag_stack.try_pop(`tbody`, [
                    `table`
                ])) : e.tag_name === `tr` ? (t ||= this._tag_stack.try_pop(`caption`, [
                    `table`
                ]), t ||= this._tag_stack.try_pop(`colgroup`, [
                    `table`
                ]), t ||= this._tag_stack.try_pop(`tr`, [
                    `table`,
                    `thead`,
                    `tbody`,
                    `tfoot`
                ])) : (e.tag_name === `th` || e.tag_name === `td`) && (t ||= this._tag_stack.try_pop(`td`, [
                    `table`,
                    `thead`,
                    `tbody`,
                    `tfoot`,
                    `tr`
                ]), t ||= this._tag_stack.try_pop(`th`, [
                    `table`,
                    `thead`,
                    `tbody`,
                    `tfoot`,
                    `tr`
                ]));
                return e.parent = this._tag_stack.get_parser_token(), t;
            }
        }, t.exports.Beautifier = m;
    })), Bt = t(((e, t)=>{
        var n = zt().Beautifier, r = Lt().Options;
        function i(e, t, r, i) {
            return new n(e, t, r, i).beautify();
        }
        t.exports = i, t.exports.defaultOptions = function() {
            return new r;
        };
    })), Vt = t(((e, t)=>{
        var n = Nt(), r = It(), i = Bt();
        function a(e, t, a, o) {
            return a ||= n, o ||= r, i(e, t, a, o);
        }
        a.defaultOptions = i.defaultOptions, t.exports.js = n, t.exports.css = r, t.exports.html = a;
    })), Ht = t(((e, t)=>{
        function n(e, t, n) {
            var r = function(t, n) {
                return e.js_beautify(t, n);
            };
            return r.js = e.js_beautify, r.css = t.css_beautify, r.html = n.html_beautify, r.js_beautify = e.js_beautify, r.css_beautify = t.css_beautify, r.html_beautify = n.html_beautify, r;
        }
        typeof define == `function` && define.amd ? define([
            `./lib/beautify`,
            `./lib/beautify-css`,
            `./lib/beautify-html`
        ], function(e, t, r) {
            return n(e, t, r);
        }) : (function(e) {
            var t = Vt();
            t.js_beautify = t.js, t.css_beautify = t.css, t.html_beautify = t.html, e.exports = n(t, t, t);
        })(t);
    })), Ut = t(((e, t)=>{
        var n = ``, r;
        t.exports = i;
        function i(e, t) {
            if (typeof e != `string`) throw TypeError(`expected a string`);
            if (t === 1) return e;
            if (t === 2) return e + e;
            var i = e.length * t;
            if (r !== e || r === void 0) r = e, n = ``;
            else if (n.length >= i) return n.substr(0, i);
            for(; i > n.length && t > 1;)t & 1 && (n += e), t >>= 1, e += e;
            return n += e, n = n.substr(0, i), n;
        }
    })), Wt = t(((e, t)=>{
        var n = Ut(), r = function(e) {
            return e.split(/(<\/?[^>]+>)/g).filter(function(e) {
                return e.trim() !== ``;
            });
        }, i = function(e) {
            return /<[^>!]+>/.test(e);
        }, a = function(e) {
            return /<\?[^?>]+\?>/.test(e);
        }, o = function(e) {
            return /<\/+[^>]+>/.test(e);
        }, s = function(e) {
            return /<[^>]+\/>/.test(e);
        }, c = function(e) {
            return i(e) && !o(e) && !s(e) && !a(e);
        };
        t.exports = function(e, t) {
            var i = 0;
            t ||= `    `;
            var a = !1, s = [];
            return r(e).map(function(e) {
                if (e.trim().startsWith(`<![CDATA[`) && (a = !0), e.trim().endsWith(`]]>`)) {
                    a = !1, s.push(e);
                    var r = s.join(``);
                    return s = [], r;
                }
                if (a) return s.push(e), null;
                e = e.replace(/^\s+|\s+$/g, ``), o(e) && i--;
                var l = n(t, i) + e;
                return c(e) && i++, l;
            }).filter(function(e) {
                return e;
            }).join(`
`);
        };
    })), Gt = vt(), Kt = e(Ht(), 1), qt = e(Wt(), 1);
    const Jt = self, Yt = (e)=>e.buffer.slice(e.byteOffset, e.byteOffset + e.byteLength);
    async function Xt(e) {
        let { id: t, buffer: n, encodings: r } = e;
        return {
            id: t,
            inputBuffer: n,
            decodedBuffer: Yt(await (0, Gt.decodeBuffer)(n, r.length === 1 ? r[0] : r.join(`, `)))
        };
    }
    async function Zt(e) {
        let { id: t, buffer: n, encodings: r } = e;
        return {
            id: t,
            encodedBuffer: Yt(await r.reduce((e, t)=>e.then((e)=>(0, Gt.encodeBuffer)(e, t)), Promise.resolve(Buffer.from(n))))
        };
    }
    const Qt = 1024 * 1024 * 5;
    function $t(e, t) {
        let n = e.toString(`utf8`);
        switch(t){
            case `raw`:
                {
                    let t = e, n = t.length > Qt;
                    n && (t = t.subarray(0, Qt));
                    let r = t.toString(`hex`);
                    return n ? `${r}\n[-- Truncated to 5MB --]` : r;
                }
            case `json`:
                try {
                    return JSON.stringify(JSON.parse(n), null, 2);
                } catch  {
                    return n;
                }
            case `xml`:
                return (0, qt.default)(n, `  `);
            case `html`:
                return Kt.default.html(n, {
                    indent_size: 2
                });
            case `css`:
                return Kt.default.css(n, {
                    indent_size: 2
                });
            case `javascript`:
                return Kt.default.js(n, {
                    indent_size: 2
                });
            default:
                return n;
        }
    }
    async function en(e) {
        let t = $t(Buffer.from(e.buffer), e.format);
        return {
            id: e.id,
            formatted: t
        };
    }
    async function tn(e) {
        let { decodedBuffer: t } = e;
        return {
            id: e.id,
            encodingSizes: {
                br: (await (0, Gt.brotliCompress)(t)).length,
                zstd: (await (0, Gt.zstdCompress)(t)).length,
                gzip: (await (0, Gt.gzip)(t, {
                    level: 9
                })).length,
                deflate: (await (0, Gt.deflate)(t, {
                    level: 9
                })).length
            }
        };
    }
    Jt.addEventListener(`message`, async (e)=>{
        try {
            switch(e.data.type){
                case `decode`:
                    try {
                        let t = await Xt(e.data);
                        Jt.postMessage(t, [
                            t.inputBuffer,
                            t.decodedBuffer
                        ]);
                    } catch (t) {
                        let n = t;
                        typeof n == `string` && (n = Error(n)), Jt.postMessage({
                            id: e.data.id,
                            error: Object.assign(f(n), {
                                inputBuffer: e.data.buffer
                            })
                        }, [
                            e.data.buffer
                        ]);
                    }
                    break;
                case `encode`:
                    {
                        let t = await Zt(e.data);
                        Jt.postMessage(t, [
                            t.encodedBuffer
                        ]);
                        break;
                    }
                case `test-encodings`:
                    Jt.postMessage(await tn(e.data));
                    break;
                case `format`:
                    {
                        let t = await en(e.data);
                        Jt.postMessage(t);
                        break;
                    }
                default:
                    console.error(`Unknown worker event`, e);
            }
        } catch (t) {
            Jt.postMessage({
                id: e.data.id,
                error: f(t)
            });
        }
    });
})();
