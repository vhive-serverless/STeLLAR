(()=>{"use strict";var t="org.transcrypt.__runtime__",r={};function e(t,r,e){return t&&(t.hasOwnProperty("__class__")||"string"==typeof t||t instanceof String)?(e&&Object.defineProperty(t,e,{value:function(){var e=[].slice.apply(arguments);return r.apply(null,[t].concat(e))},writable:!0,enumerable:!0,configurable:!0}),function(){var e=[].slice.apply(arguments);return r.apply(null,[t.__proxy__?t.__proxy__:t].concat(e))}):r}r.interpreter_name="python",r.transpiler_name="transcrypt",r.executor_name=r.transpiler_name,r.transpiler_version="3.9.0";var n={__name__:"type",__bases__:[],__new__:function(t,r,e,n){for(var _=function(){var t=[].slice.apply(arguments);return _.__new__(t)},i=e.length-1;i>=0;i--){var o=e[i];for(var s in o)null!=(a=Object.getOwnPropertyDescriptor(o,s))&&Object.defineProperty(_,s,a);for(let t of Object.getOwnPropertySymbols(o)){let r=Object.getOwnPropertyDescriptor(o,t);Object.defineProperty(_,t,r)}}for(var s in _.__metaclass__=t,_.__name__=r.startsWith("py_")?r.slice(3):r,_.__bases__=e,n){var a=Object.getOwnPropertyDescriptor(n,s);Object.defineProperty(_,s,a)}for(let t of Object.getOwnPropertySymbols(n)){let r=Object.getOwnPropertyDescriptor(n,t);Object.defineProperty(_,t,r)}return _}};n.__metaclass__=n;var _={__init__:function(t){},__metaclass__:n,__name__:"object",__bases__:[],__new__:function(t){var r=Object.create(this,{__class__:{value:this,enumerable:!0}});return("__getattr__"in this||"__setattr__"in this)&&(r.__proxy__=new Proxy(r,{get:function(t,r){let e=t[r];return null==e?t.__getattr__(r):e},set:function(t,r,e){try{t.__setattr__(r,e)}catch(n){t[r]=e}return!0}}),r=r.__proxy__),this.__init__.apply(null,[r].concat(t)),r}};function i(t,r,e,n){return void 0===n&&(n=r[0].__metaclass__),n.__new__(n,t,r,e)}function o(t){return t.__kwargtrans__=null,t.constructor=Object,t}function s(t,r,e){t.hasOwnProperty(r)||Object.defineProperty(t,r,e)}function a(t){return t.startswith("__")&&t.endswith("__")||"constructor"==t||t.startswith("py_")}function u(t){if(null==t)return 0;if(t.__len__ instanceof Function)return t.__len__();if(void 0!==t.length)return t.length;var r=0;for(var e in t)a(e)||r++;return r}function l(t){if("inf"==t)return 1/0;if("-inf"==t)return-1/0;if("nan"==t)return NaN;if(isNaN(parseFloat(t))){if(!1===t)return 0;if(!0===t)return 1;throw C("could not convert string to float: '"+A(t)+"'",new Error)}return+t}function c(t){return 0|l(t)}function p(t){return!(null==(r=t)||!(["boolean","number"].indexOf(typeof r)>=0?r:r.__bool__ instanceof Function?r.__bool__()&&r:r.__len__ instanceof Function?0!==r.__len__()&&r:(r instanceof Function||0!==u(r))&&r));var r}function f(t){var r=typeof t;if("object"!=r)return"boolean"==r?p:"string"==r?A:"number"==r?t%1==0?c:l:null;try{return"__class__"in t?t.__class__:_}catch(t){return r}}function h(t,r){if(r instanceof Array){for(let e of r)if(h(t,e))return!0;return!1}try{var e=t;if(e==r)return!0;for(var n=[].slice.call(e.__bases__);n.length;){if((e=n.shift())==r)return!0;e.__bases__.length&&(n=[].slice.call(e.__bases__).concat(n))}return!1}catch(e){return t==r||r==_}}function y(t,r){try{return h("__class__"in t?t.__class__:f(t),r)}catch(e){return h(f(t),r)}}function g(t){try{return t.__repr__()}catch(i){try{return t.__str__()}catch(i){try{if(null==t)return"None";if(t.constructor==Object){var r="{",e=!1;for(var n in t)if(!a(n)){if(n.isnumeric())var _=n;else _="'"+n+"'";e?r+=", ":e=!0,r+=_+": "+g(t[n])}return r+"}"}return"boolean"==typeof t?t.toString().capitalize():t.toString()}catch(r){return"<object of type: "+typeof t+">"}}}}function m(t){this.iterable=t,this.index=0}function v(t){this.iterable=t,this.index=0}function d(t){return t?Array.from(t):[]}function b(t){let r=t?[].slice.apply(t):[];return r.__class__=b,r}function w(t){let r=[];if(t)for(let e=0;e<t.length;e++)r.add(t[e]);return r.__class__=w,r}function A(t){if("number"==typeof t)return t.toString();try{return t.__str__()}catch(r){try{return g(t)}catch(r){return String(t)}}}function O(t){return this.hasOwnProperty(t)}function x(){var t=[];for(var r in this)a(r)||t.push(r);return t}function S(){var t=[];for(var r in this)a(r)||t.push([r,this[r]]);return t}function j(t){delete this[t]}function k(){for(var t in this)delete this[t]}function P(t,r){var e=this[t];return null==e&&(e=this["py_"+t]),null==e?null==r?null:r:e}function E(t,r){var e=this[t];if(null!=e)return e;var n=null==r?null:r;return this[t]=n,n}function N(t,r){var e=this[t];if(null!=e)return delete this[t],e;if(void 0===r)throw L(t,new Error);return r}function U(){var t=Object.keys(this)[0];if(null==t)throw L("popitem(): dictionary is empty",new Error);var r=b([t,this[t]]);return delete this[t],r}function F(t){for(var r in t)this[r]=t[r]}function I(){var t=[];for(var r in this)a(r)||t.push(this[r]);return t}function W(t){return this[t]}function $(t,r){this[t]=r}function q(t){var r={};if(!t||t instanceof Array){if(t)for(var e=0;e<t.length;e++){var n=t[e];if(!(n instanceof Array)||2!=n.length)throw C("dict update sequence element #"+e+" has length "+n.length+"; 2 is required",new Error);var _=n[0],i=n[1];!(t instanceof Array)&&t instanceof Object&&(y(t,q)||(i=q(i))),r[_]=i}}else if(y(t,q)){var o=t.py_keys();for(e=0;e<o.length;e++)r[_=o[e]]=t[_]}else{if(!(t instanceof Object))throw C("Invalid type of object for dict creation",new Error);r=t}return s(r,"__class__",{value:q,enumerable:!1,writable:!0}),s(r,"__contains__",{value:O,enumerable:!1}),s(r,"py_keys",{value:x,enumerable:!1}),s(r,"__iter__",{value:function(){new m(this.py_keys())},enumerable:!1}),s(r,Symbol.iterator,{value:function(){new v(this.py_keys())},enumerable:!1}),s(r,"py_items",{value:S,enumerable:!1}),s(r,"py_del",{value:j,enumerable:!1}),s(r,"py_clear",{value:k,enumerable:!1}),s(r,"py_get",{value:P,enumerable:!1}),s(r,"py_setdefault",{value:E,enumerable:!1}),s(r,"py_pop",{value:N,enumerable:!1}),s(r,"py_popitem",{value:U,enumerable:!1}),s(r,"py_update",{value:F,enumerable:!1}),s(r,"py_values",{value:I,enumerable:!1}),s(r,"__getitem__",{value:W,enumerable:!1}),s(r,"__setitem__",{value:$,enumerable:!1}),r}r.executor_name=r.transpiler_name,l.__name__="float",l.__bases__=[_],c.__name__="int",c.__bases__=[_],p.__name__="bool",p.__bases__=[c],Math.abs,m.prototype.__next__=function(){if(this.index<this.iterable.length)return this.iterable[this.index++];throw T(new Error)},v.prototype.next=function(){return this.index<this.iterable.py_keys.length?{value:this.index++,done:!1}:{value:void 0,done:!0}},Array.prototype.__class__=d,d.__name__="list",d.__bases__=[_],Array.prototype.__iter__=function(){return new m(this)},Array.prototype.__getslice__=function(t,r,e){if(t<0&&(t=this.length+t),null==r?r=this.length:r<0?r=this.length+r:r>this.length&&(r=this.length),1==e)return Array.prototype.slice.call(this,t,r);let n=d([]);for(let _=t;_<r;_+=e)n.push(this[_]);return n},Array.prototype.__setslice__=function(t,r,e,n){if(t<0&&(t=this.length+t),null==r?r=this.length:r<0&&(r=this.length+r),null==e)Array.prototype.splice.apply(this,[t,r-t].concat(n));else{let _=0;for(let i=t;i<r;i+=e)this[i]=n[_++]}},Array.prototype.__repr__=function(){if(this.__class__==w&&!this.length)return"set()";let t=this.__class__&&this.__class__!=d?this.__class__==b?"(":"{":"[";for(let r=0;r<this.length;r++)r&&(t+=", "),t+=g(this[r]);return this.__class__==b&&1==this.length&&(t+=","),t+=this.__class__&&this.__class__!=d?this.__class__==b?")":"}":"]",t},Array.prototype.__str__=Array.prototype.__repr__,Array.prototype.append=function(t){this.push(t)},Array.prototype.py_clear=function(){this.length=0},Array.prototype.extend=function(t){this.push.apply(this,t)},Array.prototype.insert=function(t,r){this.splice(t,0,r)},Array.prototype.remove=function(t){let r=this.indexOf(t);if(-1==r)throw C("list.remove(x): x not in list",new Error);this.splice(r,1)},Array.prototype.index=function(t){return this.indexOf(t)},Array.prototype.py_pop=function(t){return null==t?this.pop():this.splice(t,1)[0]},Array.prototype.py_sort=function(){M.apply(null,[this].concat([].slice.apply(arguments)))},Array.prototype.__add__=function(t){return d(this.concat(t))},Array.prototype.__mul__=function(t){let r=this;for(let e=1;e<t;e++)r=r.concat(this);return r},Array.prototype.__rmul__=Array.prototype.__mul__,b.__name__="tuple",b.__bases__=[_],w.__name__="set",w.__bases__=[_],Array.prototype.__bindexOf__=function(t){t+="";let r=0,e=this.length-1;for(;r<=e;){let n=(r+e)/2|0,_=this[n]+"";if(_<t)r=n+1;else{if(!(_>t))return n;e=n-1}}return-1},Array.prototype.add=function(t){-1==this.indexOf(t)&&this.push(t)},Array.prototype.discard=function(t){var r=this.indexOf(t);-1!=r&&this.splice(r,1)},Array.prototype.isdisjoint=function(t){this.sort();for(let r=0;r<t.length;r++)if(-1!=this.__bindexOf__(t[r]))return!1;return!0},Array.prototype.issuperset=function(t){this.sort();for(let r=0;r<t.length;r++)if(-1==this.__bindexOf__(t[r]))return!1;return!0},Array.prototype.issubset=function(t){return w(t.slice()).issuperset(this)},Array.prototype.union=function(t){let r=w(this.slice().sort());for(let e=0;e<t.length;e++)-1==r.__bindexOf__(t[e])&&r.push(t[e]);return r},Array.prototype.intersection=function(t){this.sort();let r=w();for(let e=0;e<t.length;e++)-1!=this.__bindexOf__(t[e])&&r.push(t[e]);return r},Array.prototype.difference=function(t){let r=w(t.slice().sort()),e=w();for(let t=0;t<this.length;t++)-1==r.__bindexOf__(this[t])&&e.push(this[t]);return e},Array.prototype.symmetric_difference=function(t){return this.union(t).difference(this.intersection(t))},Array.prototype.py_update=function(){let t=[].concat.apply(this.slice(),arguments).sort();this.py_clear();for(let r=0;r<t.length;r++)t[r]!=t[r-1]&&this.push(t[r])},Array.prototype.__eq__=function(t){if(this.length!=t.length)return!1;this.__class__==w&&(this.sort(),t.sort());for(let r=0;r<this.length;r++)if(this[r]!=t[r])return!1;return!0},Array.prototype.__ne__=function(t){return!this.__eq__(t)},Array.prototype.__le__=function(t){if(this.__class__==w)return this.issubset(t);for(let r=0;r<this.length;r++){if(this[r]>t[r])return!1;if(this[r]<t[r])return!0}return!0},Array.prototype.__ge__=function(t){if(this.__class__==w)return this.issuperset(t);for(let r=0;r<this.length;r++){if(this[r]<t[r])return!1;if(this[r]>t[r])return!0}return!0},Array.prototype.__lt__=function(t){return this.__class__==w?this.issubset(t)&&!this.issuperset(t):!this.__ge__(t)},Array.prototype.__gt__=function(t){return this.__class__==w?this.issuperset(t)&&!this.issubset(t):!this.__le__(t)},Uint8Array.prototype.__add__=function(t){let r=new Uint8Array(this.length+t.length);return r.set(this),r.set(t,this.length),r},Uint8Array.prototype.__mul__=function(t){let r=new Uint8Array(t*this.length);for(let e=0;e<t;e++)r.set(this,e*this.length);return r},Uint8Array.prototype.__rmul__=Uint8Array.prototype.__mul__,String.prototype.__class__=A,A.__name__="str",A.__bases__=[_],String.prototype.__iter__=function(){new m(this)},String.prototype.__repr__=function(){return(-1==this.indexOf("'")?"'"+this+"'":'"'+this+'"').py_replace("\t","\\t").py_replace("\n","\\n")},String.prototype.__str__=function(){return this},String.prototype.capitalize=function(){return this.charAt(0).toUpperCase()+this.slice(1)},String.prototype.endswith=function(t){if(!(t instanceof Array))return""==t||this.slice(-t.length)==t;for(var r=0;r<t.length;r++)if(this.slice(-t[r].length)==t[r])return!0;return!1},String.prototype.find=function(t,r){return this.indexOf(t,r)},String.prototype.__getslice__=function(t,r,e){t<0&&(t=this.length+t),null==r?r=this.length:r<0&&(r=this.length+r);var n="";if(1==e)n=this.substring(t,r);else for(var _=t;_<r;_+=e)n=n.concat(this.charAt(_));return n},s(String.prototype,"format",{get:function(){return e(this,(function(t){var r=b([].slice.apply(arguments).slice(1)),e=0;return t.replace(/\{(\w*)\}/g,(function(t,n){if(""==n&&(n=e++),n==+n)return void 0===r[n]?t:A(r[n]);for(var _=0;_<r.length;_++)if("object"==typeof r[_]&&void 0!==r[_][n])return A(r[_][n]);return t}))}))},enumerable:!0}),String.prototype.isalnum=function(){return/^[0-9a-zA-Z]{1,}$/.test(this)},String.prototype.isalpha=function(){return/^[a-zA-Z]{1,}$/.test(this)},String.prototype.isdecimal=function(){return/^[0-9]{1,}$/.test(this)},String.prototype.isdigit=function(){return this.isdecimal()},String.prototype.islower=function(){return/^[a-z]{1,}$/.test(this)},String.prototype.isupper=function(){return/^[A-Z]{1,}$/.test(this)},String.prototype.isspace=function(){return/^[\s]{1,}$/.test(this)},String.prototype.isnumeric=function(){return!isNaN(parseFloat(this))&&isFinite(this)},String.prototype.join=function(t){return(t=Array.from(t)).join(this)},String.prototype.lower=function(){return this.toLowerCase()},String.prototype.py_replace=function(t,r,e){return this.split(t,e).join(r)},String.prototype.lstrip=function(){return this.replace(/^\s*/g,"")},String.prototype.rfind=function(t,r){return this.lastIndexOf(t,r)},String.prototype.rsplit=function(t,r){if(null==t||null==t){t=/\s+/;var e=this.strip()}else e=this;if(null==r||-1==r)return e.split(t);var n=e.split(t);if(r<n.length){var _=n.length-r;return[n.slice(0,_).join(t)].concat(n.slice(_))}return n},String.prototype.rstrip=function(){return this.replace(/\s*$/g,"")},String.prototype.py_split=function(t,r){if(null==t||null==t){t=/\s+/;var e=this.strip()}else e=this;if(null==r||-1==r)return e.split(t);var n=e.split(t);return r<n.length?n.slice(0,r).concat([n.slice(r).join(t)]):n},String.prototype.startswith=function(t){if(!(t instanceof Array))return 0==this.indexOf(t);for(var r=0;r<t.length;r++)if(0==this.indexOf(t[r]))return!0;return!1},String.prototype.strip=function(){return this.trim()},String.prototype.upper=function(){return this.toUpperCase()},String.prototype.__mul__=function(t){for(var r="",e=0;e<t;e++)r+=this;return r},String.prototype.__rmul__=String.prototype.__mul__,q.__name__="dict",q.__bases__=[_],s(Function.prototype,"__setdoc__",{value:function(t){return this.__doc__=t,this},enumerable:!1});var z=i("BaseException",[_],{__module__:t}),D=i("Exception",[z],{__module__:t,get __init__(){return e(this,(function(t){var r=q();if(arguments.length){var e=arguments.length-1;if(arguments[e]&&arguments[e].hasOwnProperty("__kwargtrans__")){var n=arguments[e--];for(var _ in n)"self"===_?t=n[_]:r[_]=n[_];delete r.__kwargtrans__}var i=b([].slice.apply(arguments).slice(1,e+1))}else i=b();t.__args__=i,null!=r.error?t.stack=r.error.stack:Error?t.stack=(new Error).stack:t.stack="No stack trace available"}))},get __repr__(){return e(this,(function(t){return u(t.__args__)>1?"{}{}".format(t.__class__.__name__,g(b(t.__args__))):u(t.__args__)?"{}({})".format(t.__class__.__name__,g(t.__args__[0])):"{}()".format(t.__class__.__name__)}))},get __str__(){return e(this,(function(t){return u(t.__args__)>1?A(b(t.__args__)):u(t.__args__)?A(t.__args__[0]):""}))}}),T=(i("IterableError",[D],{__module__:t,get __init__(){return e(this,(function(t,r){D.__init__(t,"Can't iterate over non-iterable",o({error:r}))}))}}),i("StopIteration",[D],{__module__:t,get __init__(){return e(this,(function(t,r){D.__init__(t,"Iterator exhausted",o({error:r}))}))}})),C=i("ValueError",[D],{__module__:t,get __init__(){return e(this,(function(t,r,e){D.__init__(t,r,o({error:e}))}))}}),L=i("KeyError",[D],{__module__:t,get __init__(){return e(this,(function(t,r,e){D.__init__(t,r,o({error:e}))}))}}),H=(i("AssertionError",[D],{__module__:t,get __init__(){return e(this,(function(t,r,e){r?D.__init__(t,r,o({error:e})):D.__init__(t,o({error:e}))}))}}),i("NotImplementedError",[D],{__module__:t,get __init__(){return e(this,(function(t,r,e){D.__init__(t,r,o({error:e}))}))}}),i("IndexError",[D],{__module__:t,get __init__(){return e(this,(function(t,r,e){D.__init__(t,r,o({error:e}))}))}}),i("AttributeError",[D],{__module__:t,get __init__(){return e(this,(function(t,r,e){D.__init__(t,r,o({error:e}))}))}}),i("py_TypeError",[D],{__module__:t,get __init__(){return e(this,(function(t,r,e){D.__init__(t,r,o({error:e}))}))}}),i("Warning",[D],{__module__:t})),M=(i("UserWarning",[H],{__module__:t}),i("DeprecationWarning",[H],{__module__:t}),i("RuntimeWarning",[H],{__module__:t}),function(t,r,e){if((void 0===r||null!=r&&r.hasOwnProperty("__kwargtrans__"))&&(r=null),(void 0===e||null!=e&&e.hasOwnProperty("__kwargtrans__"))&&(e=!1),arguments.length){var n=arguments.length-1;if(arguments[n]&&arguments[n].hasOwnProperty("__kwargtrans__")){var _=arguments[n--];for(var i in _)switch(i){case"iterable":t=_[i];break;case"key":r=_[i];break;case"reverse":e=_[i]}}}r?t.sort((function(t,e){if(arguments.length){var n=arguments.length-1;if(arguments[n]&&arguments[n].hasOwnProperty("__kwargtrans__")){var _=arguments[n--];for(var i in _)switch(i){case"a":t=_[i];break;case"b":e=_[i]}}}return r(t)>r(e)?1:-1})):t.sort(),e&&t.reverse()}),Z=i("__Terminal__",[_],{__module__:t,get __init__(){return e(this,(function(t){t.buffer="";try{t.element=document.getElementById("__terminal__")}catch(r){t.element=null}t.element&&(t.element.style.overflowX="auto",t.element.style.boxSizing="border-box",t.element.style.padding="5px",t.element.innerHTML="_")}))},get print(){return e(this,(function(t){var r=" ",e="\n";if(arguments.length){var n=arguments.length-1;if(arguments[n]&&arguments[n].hasOwnProperty("__kwargtrans__")){var _=arguments[n--];for(var i in _)switch(i){case"self":t=_[i];break;case"sep":r=_[i];break;case"end":e=_[i]}}var o=b([].slice.apply(arguments).slice(1,n+1))}else o=b();t.buffer="{}{}{}".format(t.buffer,r.join(function(){var t=[];for(var r of o)t.append(A(r));return t}()),e).__getslice__(-4096,null,1),t.element?(t.element.innerHTML=t.buffer.py_replace("\n","<br>").py_replace(" ","&nbsp"),t.element.scrollTop=t.element.scrollHeight):console.log(r.join(function(){var t=[];for(var r of o)t.append(A(r));return t}()))}))},get input(){return e(this,(function(t,r){if(arguments.length){var e=arguments.length-1;if(arguments[e]&&arguments[e].hasOwnProperty("__kwargtrans__")){var n=arguments[e--];for(var _ in n)switch(_){case"self":t=n[_];break;case"question":r=n[_]}}}t.print("{}".format(r),o({end:""}));var i=window.prompt("\n".join(t.buffer.py_split("\n").__getslice__(-8,null,1)));return t.print(i),i}))}})();Z.print,Z.input,addEventListener("fetch",(function(t){return t.respondWith((t.request,new Response("Python Worker hello world!",q({headers:q({"content-type":"text/plain"})}))))}))})();