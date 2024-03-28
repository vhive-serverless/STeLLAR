// Transcrypt'ed from Python, 2023-09-20 06:42:19
import {AssertionError, AttributeError, BaseException, DeprecationWarning, Exception, IndexError, IterableError, KeyError, NotImplementedError, RuntimeWarning, StopIteration, UserWarning, ValueError, Warning, __JsIterator__, __PyIterator__, __Terminal__, __add__, __and__, __call__, __class__, __envir__, __eq__, __floordiv__, __ge__, __get__, __getcm__, __getitem__, __getslice__, __getsm__, __gt__, __i__, __iadd__, __iand__, __idiv__, __ijsmod__, __ilshift__, __imatmul__, __imod__, __imul__, __in__, __init__, __ior__, __ipow__, __irshift__, __isub__, __ixor__, __jsUsePyNext__, __jsmod__, __k__, __kwargtrans__, __le__, __lshift__, __lt__, __matmul__, __mergefields__, __mergekwargtrans__, __mod__, __mul__, __ne__, __neg__, __nest__, __or__, __pow__, __pragma__, __pyUseJsNext__, __rshift__, __setitem__, __setproperty__, __setslice__, __sort__, __specialattrib__, __sub__, __super__, __t__, __terminal__, __truediv__, __withblock__, __xor__, abs, all, any, assert, bool, bytearray, bytes, callable, chr, copy, deepcopy, delattr, dict, dir, divmod, enumerate, filter, float, getattr, hasattr, input, int, isinstance, issubclass, len, list, map, max, min, object, ord, pow, print, property, py_TypeError, py_iter, py_metatype, py_next, py_reversed, py_typeof, range, repr, round, set, setattr, sorted, str, sum, tuple, zip} from './org.transcrypt.__runtime__.js';
import {datetime} from './datetime.js';
var __name__ = '__main__';
export var handleRequest = function (request) {
	var incr_limit = 0;
	if (__in__ ('queryStringParameters', request) && __in__ ('IncrementLimit', request ['queryStringParameters'])) {
		var incr_limit = int (request ['queryStringParameters'].py_get ('IncrementLimit', 0));
	}
	else if (__in__ ('body', request) && JSON.parse (request ['body']) ['IncrementLimit']) {
		var incr_limit = int (JSON.parse (request ['body']) ['IncrementLimit']);
	}
	simulate_work (incr_limit);
	var response = JSON.stringify (dict ({'RequestID': 'cloudflare-does-not-specify', 'TimestampChain': [str (datetime.now ())]}));
	return new Response (response, dict ({'headers': dict ({'content-type': 'application/json'})}));
};
export var simulate_work = function (increment) {
	var num = 0;
	while (num < increment) {
		num++;
	}
};
addEventListener ('fetch', (function __lambda__ (event) {
	return event.respondWith (handleRequest (event.request));
}));

//# sourceMappingURL=index.map