/*eslint-disable*/
export function  cdf (data)  {
    "use strict";
    var f, sorted, xs, ps, i, j, l, xx;
    if (Array.isArray(data) && (data.length > 0)) {
      for (i = 0, l = data.length; i < l; ++i) {
        if (typeof(data[i]) !== 'number') {
          throw new TypeError("cdf data must be an array of finite numbers, got:" + typeof(data[i]) + " at " + i);
        }
        if (!isFinite(data[i])) {
          throw new TypeError("cdf data must be an array of finite numbers, got:" + data[i] + " at " + i);
        }
      }
      sorted = data.slice().sort(function(a, b) {
        return +a - b;
      });
      xs = [];
      ps = [];
      j = 0;
      l = sorted.length;
      xs[0] = sorted[0];
      ps[0] = 1 / l;
      for (i = 1; i < l; ++i) {
        xx = sorted[i];
        if (xx === xs[j]) {
          ps[j] = (1 + i) / l;
        } else {
          j++;
          xs[j] = xx;
          ps[j] = (1 + i) / l;
        }
      }
      f = function(x) {
        if (typeof(x) !== 'number') throw new TypeError('cdf function input must be a number, got:' + typeof(x));
        if (Number.isNaN(x)) return Number.NaN;
        var left = 0,
          right = xs.length - 1,
          mid, midval, iteration;
        if (x < xs[0]) return 0;
        if (x >= xs[xs.length - 1]) return 1;
        iteration = 0;
        while ((right - left) > 1) {
          mid = Math.floor((left + right) / 2);
          midval = xs[mid];
          if (x > midval)
            left = mid;
          else if (x < midval)
            right = mid;
          else if (x === midval) {
            left = mid;
            right = mid;
          }
          ++iteration;
          if (iteration>40) throw new Error("cdf function exceeded 40 bisection iterations, aborting bisection loop");
        }
        return ps[left];
      };
      f.xs = function() {
        return xs;
      };
      f.ps = function() {
        return ps;
      };
    } else {
      // missing or zero length data
      throw new TypeError("cdf data must be an array of finite numbers, got: missing or empty array");
    }
    return f;
  };