/* starbase runtime: navigation, math, and interactive widgets. Vanilla JS, no build step. */
(function () {
  "use strict";

  var staticWidgets = []; // {redraw} for charts/plots, repainted on theme change

  document.addEventListener("DOMContentLoaded", function () {
    initTheme();
    initSidebarScroll();
    initTOC();
    initWidgets();
    initMath();
  });

  /* ---------------- navigation filter ---------------- */
  window.sgFilterNav = function (q) {
    q = (q || "").trim().toLowerCase();
    document.querySelectorAll(".sg-nav li[data-name]").forEach(function (li) {
      var hit = !q || li.dataset.name.indexOf(q) !== -1 ||
        li.textContent.toLowerCase().indexOf(q) !== -1;
      li.classList.toggle("sg-hidden", !hit);
    });
  };

  function initTheme() {
    var saved = localStorage.getItem("sg-theme");
    if (saved) document.documentElement.setAttribute("data-theme", saved);
    var cur = document.documentElement.getAttribute("data-theme") || "dark";
    document.querySelectorAll(".sg-theme-select").forEach(function (sel) { sel.value = cur; });
  }
  window.sgSetTheme = function (name) {
    document.documentElement.setAttribute("data-theme", name);
    try { localStorage.setItem("sg-theme", name); } catch (e) {}
    document.querySelectorAll(".sg-theme-select").forEach(function (sel) { sel.value = name; });
    // Static visuals (charts/plots) are drawn once — repaint them in the new palette.
    // Animated sims/sketches read CSS variables every frame and adapt on their own.
    staticWidgets.forEach(function (w) { try { w.redraw(); } catch (e) {} });
  };

  /* scroll the sidebar so the current page sits in view among its neighbors */
  function initSidebarScroll() {
    var nav = document.querySelector(".sg-nav");
    if (!nav) return;
    var active = nav.querySelector('[aria-current="page"]');
    if (!active) return;
    var navRect = nav.getBoundingClientRect(), aRect = active.getBoundingClientRect();
    nav.scrollTop += (aRect.top - navRect.top) - nav.clientHeight / 2 + aRect.height / 2;
  }

  /* ---------------- audio: a tiny synth sketches can call ----------------
     Browsers require audio to begin from a user gesture, so we unlock the
     context on the first pointer/key event; sketches then call sgTone(freq). */
  function audioCtx() {
    var Ctx = window.AudioContext || window.webkitAudioContext;
    if (!Ctx) return null;
    if (!window.__sgAudio) window.__sgAudio = new Ctx();
    if (window.__sgAudio.state === "suspended") window.__sgAudio.resume();
    return window.__sgAudio;
  }
  document.addEventListener("pointerdown", audioCtx);
  document.addEventListener("keydown", audioCtx);
  window.sgTone = function (freq, dur, type) {
    try {
      var ac = audioCtx();
      if (!ac || !freq) return;
      dur = dur || 0.6;
      var osc = ac.createOscillator(), g = ac.createGain(), t0 = ac.currentTime;
      osc.type = type || "triangle";
      osc.frequency.value = freq;
      g.gain.setValueAtTime(0.0001, t0);
      g.gain.exponentialRampToValueAtTime(0.18, t0 + 0.012);
      g.gain.exponentialRampToValueAtTime(0.0001, t0 + dur);
      osc.connect(g); g.connect(ac.destination);
      osc.start(t0); osc.stop(t0 + dur + 0.02);
    } catch (e) {}
  };
  window.sgChord = function (freqs, dur, type) {
    (freqs || []).forEach(function (f) { window.sgTone(f, dur, type); });
  };

  /* ---------------- global topic search ---------------- */
  // The sidebar shows only the current page's connected component; search spans
  // the whole knowledge base so you can always jump to any "world".
  var sgIndex = null, sgLoading = false;
  function siteRoot() {
    var s = document.querySelector('script[src*="app.js"]');
    return (s ? s.getAttribute("src") : "").replace(/static\/app\.js.*$/, "");
  }
  function loadIndex(cb) {
    if (sgIndex) { cb(sgIndex); return; }
    if (sgLoading) return;
    sgLoading = true;
    fetch(siteRoot() + "search.json")
      .then(function (r) { return r.json(); })
      .then(function (j) { sgIndex = j; cb(j); })
      .catch(function () { sgIndex = []; });
  }
  function esc(s) {
    return (s || "").replace(/[&<>"]/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c];
    });
  }
  window.sgSearch = function (q) {
    q = (q || "").trim().toLowerCase();
    var box = document.getElementById("sg-search-results");
    if (!box) return;
    if (!q) { box.hidden = true; box.innerHTML = ""; return; }
    loadIndex(function (idx) {
      var hits = idx.filter(function (e) {
        return (e.t + " " + e.s + " " + (e.d || "")).toLowerCase().indexOf(q) !== -1;
      });
      hits.sort(function (a, b) {
        return (a.t.toLowerCase().indexOf(q) === 0 ? 0 : 1) - (b.t.toLowerCase().indexOf(q) === 0 ? 0 : 1);
      });
      hits = hits.slice(0, 14);
      var root = siteRoot();
      box.innerHTML = hits.length
        ? hits.map(function (e) {
          return '<a class="sg-search-hit" href="' + root + e.u + '"><span class="sg-search-t">' +
            esc(e.t) + "</span>" + (e.s ? '<span class="sg-search-s">' + esc(e.s) + "</span>" : "") + "</a>";
        }).join("")
        : '<div class="sg-search-empty">No matches</div>';
      box.hidden = false;
    });
  };
  window.sgSearchKey = function (ev) {
    var box = document.getElementById("sg-search-results");
    if (ev.key === "Escape") { if (box) box.hidden = true; ev.target.value = ""; return; }
    if (ev.key === "Enter" && box && !box.hidden) {
      var first = box.querySelector("a");
      if (first) window.location.href = first.getAttribute("href");
    }
  };
  document.addEventListener("click", function (e) {
    var box = document.getElementById("sg-search-results");
    if (box && !box.hidden && !e.target.closest(".sg-search")) box.hidden = true;
  });

  /* ---------------- table of contents scrollspy ---------------- */
  function initTOC() {
    var links = Array.prototype.slice.call(document.querySelectorAll(".sg-toc-rail .sg-toc a"));
    if (!links.length) return;
    var map = {};
    links.forEach(function (a) {
      var id = decodeURIComponent(a.getAttribute("href").slice(1));
      var el = document.getElementById(id);
      if (el) map[id] = a;
    });
    var obs = new IntersectionObserver(function (entries) {
      entries.forEach(function (e) {
        if (e.isIntersecting) {
          links.forEach(function (l) { l.classList.remove("sg-toc-active"); });
          var a = map[e.target.id];
          if (a) a.classList.add("sg-toc-active");
        }
      });
    }, { rootMargin: "0px 0px -75% 0px", threshold: 0 });
    Object.keys(map).forEach(function (id) { obs.observe(document.getElementById(id)); });
  }

  /* resolve a static asset URL relative to this script, carrying its ?v= */
  function assetRef(name) {
    var s = document.querySelector('script[src*="app.js"]');
    var src = s ? s.getAttribute("src") : "static/app.js";
    var base = src.replace(/app\.js.*$/, "");
    var v = (src.match(/[?&]v=([^&#]+)/) || ["", ""])[1];
    return base + name + (v ? "?v=" + v : "");
  }

  /* ---------------- math (KaTeX) ----------------
     Source is chosen by the build: data-katex holds a CDN base URL by default,
     or is empty when assets are vendored locally (--vendor) for offline use. */
  function initMath() {
    var nodes = document.querySelectorAll(".sg-math-inline, .sg-math-display");
    if (!nodes.length) return;
    var base = document.documentElement.getAttribute("data-katex");
    var cssHref = base ? base + "katex.min.css" : assetRef("katex/katex.min.css");
    var jsSrc = base ? base + "katex.min.js" : assetRef("katex/katex.min.js");
    var css = document.createElement("link");
    css.rel = "stylesheet"; css.href = cssHref;
    document.head.appendChild(css);
    var s = document.createElement("script");
    s.src = jsSrc;
    s.onload = function () {
      nodes.forEach(function (el) {
        try {
          window.katex.render(el.textContent, el, {
            displayMode: el.classList.contains("sg-math-display"),
            throwOnError: false,
          });
        } catch (e) { /* leave raw TeX visible */ }
      });
    };
    document.head.appendChild(s);
  }

  /* ---------------- widgets ---------------- */
  function cfg(el) {
    var id = el.getAttribute("data-config");
    if (!id) return {};
    var node = document.getElementById(id);
    if (!node) return {};
    try { return JSON.parse(node.textContent); } catch (e) { return {}; }
  }

  function initWidgets() {
    document.querySelectorAll('[data-widget="chart"]').forEach(function (el) {
      drawChart(el.querySelector("canvas"), cfg(el));
    });
    document.querySelectorAll('[data-widget="plot"]').forEach(function (el) {
      drawPlot(el.querySelector("canvas"), cfg(el));
    });
    document.querySelectorAll('[data-widget="sim"]').forEach(function (el) {
      var c = el.querySelector("canvas");
      var conf = cfg(el);
      var maker = SIMS[(conf.name || "").toLowerCase()];
      if (!maker) { fail(c, "unknown sim: " + conf.name); return; }
      runLoop(el, c, maker(conf));
    });
    document.querySelectorAll('[data-widget="sketch"]').forEach(function (el) {
      runSketch(el);
    });
    document.querySelectorAll('[data-widget="quiz"]').forEach(initQuiz);
  }

  function fail(canvas, msg) {
    var ctx = canvas.getContext("2d");
    ctx.fillStyle = "#ff6b81"; ctx.font = "13px system-ui";
    ctx.fillText(msg, 12, 22);
  }

  /* canvas with devicePixelRatio scaling; returns {ctx, W, H} via getters.
     onResize (optional) is called after a size change so static drawings can
     repaint — animated widgets repaint every frame and don't need it. */
  function hidpi(canvas, onResize) {
    var dpr = window.devicePixelRatio || 1;
    var first = true;
    function resize() {
      var r = canvas.getBoundingClientRect();
      canvas.width = Math.max(1, Math.round(r.width * dpr));
      canvas.height = Math.max(1, Math.round((r.height || 300) * dpr));
      if (!first && onResize) onResize();
      first = false;
    }
    resize();
    if (window.ResizeObserver) new ResizeObserver(resize).observe(canvas);
    var ctx = canvas.getContext("2d");
    return {
      ctx: ctx, dpr: dpr,
      get W() { return canvas.width / dpr; },
      get H() { return canvas.height / dpr; },
      clear: function () { ctx.setTransform(dpr, 0, 0, dpr, 0, 0); ctx.clearRect(0, 0, canvas.width, canvas.height); },
    };
  }

  function css(name) {
    return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || "#888";
  }

  /* ---- numeric data parsing for chart ---- */
  function parseSeries(data) {
    if (Array.isArray(data)) return [data.map(Number)];
    if (typeof data !== "string") return [[]];
    data = data.trim();
    if (data[0] === "[") { try { var j = JSON.parse(data); return Array.isArray(j[0]) ? j : [j]; } catch (e) {} }
    // "x:y" pairs, separated by commas, spaces, semicolons, or newlines.
    if (/[:]/.test(data)) {
      var pts = data.split(/[;,\s]+/).filter(Boolean).map(function (p) {
        var kv = p.split(":"); return [Number(kv[0]), Number(kv[1])];
      }).filter(function (p) { return !isNaN(p[0]) && !isNaN(p[1]); });
      return [pts];
    }
    return [data.split(/[\s,]+/).map(Number).filter(function (n) { return !isNaN(n); })];
  }

  function drawChart(canvas, c) {
    var view = hidpi(canvas, function () { render(); });
    function render() {
      var ctx = view.ctx, W = view.W, H = view.H;
      view.clear();
      var series = parseSeries(c.data);
      var pad = { l: 42, r: 14, t: 12, b: 28 };
      var pts = series[0] || [];
      var ys, xs;
      if (pts.length && Array.isArray(pts[0])) { xs = pts.map(function (p) { return p[0]; }); ys = pts.map(function (p) { return p[1]; }); }
      else { ys = pts; xs = pts.map(function (_, i) { return i; }); }
      if (!ys.length) return;
      var ymin = Math.min.apply(null, ys), ymax = Math.max.apply(null, ys);
      if (ymin > 0) ymin = 0; if (ymax < 0) ymax = 0;
      if (ymin === ymax) ymax = ymin + 1;
      var xmin = Math.min.apply(null, xs), xmax = Math.max.apply(null, xs);
      if (xmin === xmax) xmax = xmin + 1;
      var px = function (x) { return pad.l + (x - xmin) / (xmax - xmin) * (W - pad.l - pad.r); };
      var py = function (y) { return H - pad.b - (y - ymin) / (ymax - ymin) * (H - pad.t - pad.b); };

      // axes + gridlines
      ctx.strokeStyle = css("--border"); ctx.fillStyle = css("--text-faint");
      ctx.font = "11px " + css("--sans"); ctx.lineWidth = 1;
      for (var g = 0; g <= 4; g++) {
        var yv = ymin + (ymax - ymin) * g / 4, y = py(yv);
        ctx.globalAlpha = 0.4; ctx.beginPath(); ctx.moveTo(pad.l, y); ctx.lineTo(W - pad.r, y); ctx.stroke();
        ctx.globalAlpha = 1; ctx.fillText(fmt(yv), 4, y + 3);
      }
      var accent = c.color || css("--accent");
      if (c.type === "bar") {
        var bw = (W - pad.l - pad.r) / ys.length * 0.7;
        ys.forEach(function (y, i) {
          var x = px(xs[i]); ctx.fillStyle = accent;
          ctx.fillRect(x - bw / 2, py(Math.max(0, y)), bw, Math.abs(py(y) - py(0)));
        });
      } else if (c.type === "scatter") {
        ctx.fillStyle = accent;
        ys.forEach(function (y, i) { ctx.beginPath(); ctx.arc(px(xs[i]), py(y), 3, 0, 7); ctx.fill(); });
      } else {
        if (c.fill) {
          ctx.beginPath(); ctx.moveTo(px(xs[0]), py(0));
          ys.forEach(function (y, i) { ctx.lineTo(px(xs[i]), py(y)); });
          ctx.lineTo(px(xs[xs.length - 1]), py(0)); ctx.closePath();
          ctx.fillStyle = accent; ctx.globalAlpha = 0.12; ctx.fill(); ctx.globalAlpha = 1;
        }
        ctx.beginPath(); ctx.strokeStyle = accent; ctx.lineWidth = 2;
        ys.forEach(function (y, i) { var X = px(xs[i]), Y = py(y); i ? ctx.lineTo(X, Y) : ctx.moveTo(X, Y); });
        ctx.stroke();
      }
    }
    render();
    staticWidgets.push({ redraw: render });
  }

  function drawPlot(canvas, c) {
    var view = hidpi(canvas, function () { render(); });
    var fns = String(c.fn).split(";;").map(function (s) { return s.trim(); });
    var compiled = fns.map(function (src) { try { return new Function("x", "with (Math) { return (" + src + "); }"); } catch (e) { return null; } });
    function render() {
      var ctx = view.ctx, W = view.W, H = view.H; view.clear();
      var xmin = num(c.xmin, -6.283), xmax = num(c.xmax, 6.283), N = c.samples || 400;
      var series = compiled.map(function (f) {
        if (!f) return [];
        var arr = [];
        for (var i = 0; i <= N; i++) { var x = xmin + (xmax - xmin) * i / N; var y = f(x); if (isFinite(y)) arr.push([x, y]); }
        return arr;
      });
      var all = [].concat.apply([], series);
      if (!all.length) return;
      var ymin = c.ymin !== "" ? num(c.ymin) : Math.min.apply(null, all.map(function (p) { return p[1]; }));
      var ymax = c.ymax !== "" ? num(c.ymax) : Math.max.apply(null, all.map(function (p) { return p[1]; }));
      if (ymin === ymax) ymax += 1;
      var pad = 30;
      var px = function (x) { return pad + (x - xmin) / (xmax - xmin) * (W - pad - 10); };
      var py = function (y) { return H - pad + (y - ymin) / (ymax - ymin) * -(H - pad - 10); };
      // axes
      ctx.strokeStyle = css("--border"); ctx.lineWidth = 1; ctx.globalAlpha = 0.7;
      if (ymin < 0 && ymax > 0) { ctx.beginPath(); ctx.moveTo(px(xmin), py(0)); ctx.lineTo(px(xmax), py(0)); ctx.stroke(); }
      if (xmin < 0 && xmax > 0) { ctx.beginPath(); ctx.moveTo(px(0), py(ymin)); ctx.lineTo(px(0), py(ymax)); ctx.stroke(); }
      ctx.globalAlpha = 1;
      var palette = [css("--accent"), css("--accent-2"), css("--good"), css("--warn")];
      series.forEach(function (pts, k) {
        ctx.beginPath(); ctx.strokeStyle = palette[k % palette.length]; ctx.lineWidth = 2;
        pts.forEach(function (p, i) { var X = px(p[0]), Y = py(p[1]); i ? ctx.lineTo(X, Y) : ctx.moveTo(X, Y); });
        ctx.stroke();
      });
    }
    render();
    staticWidgets.push({ redraw: render });
  }

  /* ---- animation loop with play/pause/reset ---- */
  function runLoop(figure, canvas, sim) {
    var view = hidpi(canvas);
    var running = true, last = 0, t = 0, raf = null;
    var mouse = { x: 0, y: 0, down: false, clicked: false };
    var heldDown = false, clickLatch = false;
    canvas.addEventListener("pointermove", function (e) {
      var r = canvas.getBoundingClientRect(); mouse.x = e.clientX - r.left; mouse.y = e.clientY - r.top;
    });
    canvas.addEventListener("pointerdown", function (e) {
      var r = canvas.getBoundingClientRect(); mouse.x = e.clientX - r.left; mouse.y = e.clientY - r.top;
      heldDown = true; clickLatch = true;
    });
    window.addEventListener("pointerup", function () { heldDown = false; });
    if (sim.init) sim.init({ W: view.W, H: view.H, mouse: mouse });
    function frame(ts) {
      var dt = last ? Math.min((ts - last) / 1000, 0.05) : 0.016; last = ts;
      mouse.down = heldDown || clickLatch; mouse.clicked = clickLatch;
      if (running) { t += dt; if (sim.step) sim.step(dt, { W: view.W, H: view.H, mouse: mouse, t: t }); }
      view.clear();
      sim.draw(view.ctx, view.W, view.H, { t: t, mouse: mouse });
      clickLatch = false;
      raf = requestAnimationFrame(frame);
    }
    raf = requestAnimationFrame(frame);
    wireControls(figure, {
      toggle: function (btn) { running = !running; btn.textContent = running ? "⏸" : "▶"; last = 0; },
      reset: function () { t = 0; last = 0; if (sim.init) sim.init({ W: view.W, H: view.H, mouse: mouse }); },
    });
  }

  function runSketch(figure) {
    var canvas = figure.querySelector("canvas");
    var srcEl = document.getElementById(figure.getAttribute("data-config"));
    var src = srcEl ? srcEl.textContent : "";
    if (srcEl && srcEl.getAttribute("data-b64") === "1") {
      var b = src.trim().replace(/-/g, "+").replace(/_/g, "/");
      while (b.length % 4) b += "=";
      try { src = decodeURIComponent(escape(atob(b))); }
      catch (e) { try { src = atob(b); } catch (e2) { fail(canvas, "sketch decode error"); return; } }
    }
    var fn;
    try { fn = new Function("canvas", "ctx", "t", "dt", "W", "H", "mouse", "state", "frame", "with (Math) {" + src + "}"); }
    catch (e) { fail(canvas, "sketch error: " + e.message); return; }
    var view = hidpi(canvas), running = true, last = 0, t = 0, frameN = 0, state = {};
    var mouse = { x: 0, y: 0, down: false, clicked: false };
    // heldDown tracks the physical button; clickLatch guarantees a fast click
    // (down+up between two frames) is still seen for one frame.
    var heldDown = false, clickLatch = false;
    canvas.addEventListener("pointermove", function (e) {
      var r = canvas.getBoundingClientRect(); mouse.x = e.clientX - r.left; mouse.y = e.clientY - r.top;
    });
    canvas.addEventListener("pointerdown", function (e) {
      var r = canvas.getBoundingClientRect(); mouse.x = e.clientX - r.left; mouse.y = e.clientY - r.top;
      heldDown = true; clickLatch = true;
    });
    window.addEventListener("pointerup", function () { heldDown = false; });
    var errs = 0;
    function loop(ts) {
      var dt = last ? Math.min((ts - last) / 1000, 0.05) : 0.016; last = ts;
      if (running) {
        mouse.down = heldDown || clickLatch;
        mouse.clicked = clickLatch;
        t += dt; view.clear();
        try {
          fn(canvas, view.ctx, t, dt, view.W, view.H, mouse, state, frameN);
          frameN++;
        } catch (e) {
          // Tolerate transient errors (a single bad frame shouldn't kill the
          // animation); give up only if it keeps throwing.
          if (++errs > 8) { fail(canvas, "sketch error: " + e.message); running = false; return; }
          if (window.console) console.warn("sketch frame error:", e.message);
        }
        clickLatch = false; // consumed; true for exactly one frame after a press
      }
      requestAnimationFrame(loop);
    }
    requestAnimationFrame(loop);
    wireControls(figure, {
      toggle: function (btn) { running = !running; btn.textContent = running ? "⏸" : "▶"; last = 0; },
      reset: function () { t = 0; frameN = 0; last = 0; state = {}; errs = 0; },
    });
  }

  function wireControls(figure, handlers) {
    var bar = figure.querySelector(".sg-sim-controls");
    if (!bar) return;
    bar.addEventListener("click", function (e) {
      var btn = e.target.closest("[data-act]"); if (!btn) return;
      var h = handlers[btn.getAttribute("data-act")]; if (h) h(btn);
    });
  }

  function initQuiz(el) {
    var answer = parseInt(el.getAttribute("data-answer"), 10);
    var explain = el.querySelector(".sg-quiz-explain");
    el.querySelectorAll(".sg-quiz-opt").forEach(function (btn) {
      btn.addEventListener("click", function () {
        if (el.dataset.done) return;
        el.dataset.done = "1";
        var idx = parseInt(btn.getAttribute("data-idx"), 10);
        btn.classList.add(idx === answer ? "sg-correct" : "sg-wrong");
        if (idx !== answer) {
          var right = el.querySelector('.sg-quiz-opt[data-idx="' + answer + '"]');
          if (right) right.classList.add("sg-correct");
        }
        if (explain) explain.hidden = false;
      });
    });
  }

  /* ---------------- built-in simulations ---------------- */
  var SIMS = {
    pendulum: function (c) {
      var L = num(c.length, 1), g = num(c.gravity, 9.81), damp = num(c.damping, 0.0);
      var th0 = num(c.angle, 2.3), th, w;
      return {
        init: function () { th = th0; w = 0; },
        step: function (dt) { var a = -g / L * Math.sin(th) - damp * w; w += a * dt; th += w * dt; },
        draw: function (ctx, W, H) {
          var ox = W / 2, oy = H * 0.18, len = Math.min(H * 0.6, W * 0.4);
          var x = ox + len * Math.sin(th), y = oy + len * Math.cos(th);
          ctx.strokeStyle = css("--border"); ctx.lineWidth = 2;
          ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(x, y); ctx.stroke();
          ctx.fillStyle = css("--accent"); ctx.beginPath(); ctx.arc(x, y, 12, 0, 7); ctx.fill();
          ctx.fillStyle = css("--text-faint"); ctx.beginPath(); ctx.arc(ox, oy, 3, 0, 7); ctx.fill();
        },
      };
    },
    doublependulum: function (c) {
      var g = num(c.gravity, 9.81), m1 = 1, m2 = 1, L1 = 1, L2 = 1;
      var a1, a2, w1, w2, trail;
      function reset() { a1 = num(c.angle1, 2.4); a2 = num(c.angle2, 2.6); w1 = 0; w2 = 0; trail = []; }
      reset();
      return {
        init: reset,
        step: function (dt) {
          dt = Math.min(dt, 0.02);
          for (var s = 0; s < 3; s++) {
            var d = a1 - a2, den = 2 * m1 + m2 - m2 * Math.cos(2 * d);
            var a1a = (-g * (2 * m1 + m2) * Math.sin(a1) - m2 * g * Math.sin(a1 - 2 * a2)
              - 2 * Math.sin(d) * m2 * (w2 * w2 * L2 + w1 * w1 * L1 * Math.cos(d))) / (L1 * den);
            var a2a = (2 * Math.sin(d) * (w1 * w1 * L1 * (m1 + m2) + g * (m1 + m2) * Math.cos(a1)
              + w2 * w2 * L2 * m2 * Math.cos(d))) / (L2 * den);
            w1 += a1a * dt / 3; w2 += a2a * dt / 3; a1 += w1 * dt / 3; a2 += w2 * dt / 3;
          }
        },
        draw: function (ctx, W, H) {
          var ox = W / 2, oy = H * 0.36, sc = Math.min(H, W) * 0.2;
          var x1 = ox + sc * Math.sin(a1), y1 = oy + sc * Math.cos(a1);
          var x2 = x1 + sc * Math.sin(a2), y2 = y1 + sc * Math.cos(a2);
          trail.push([x2, y2]); if (trail.length > 220) trail.shift();
          ctx.lineWidth = 1.5;
          for (var i = 1; i < trail.length; i++) {
            ctx.globalAlpha = i / trail.length * 0.6; ctx.strokeStyle = css("--accent-2");
            ctx.beginPath(); ctx.moveTo(trail[i - 1][0], trail[i - 1][1]); ctx.lineTo(trail[i][0], trail[i][1]); ctx.stroke();
          }
          ctx.globalAlpha = 1; ctx.strokeStyle = css("--text-dim"); ctx.lineWidth = 2;
          ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(x1, y1); ctx.lineTo(x2, y2); ctx.stroke();
          ctx.fillStyle = css("--accent");
          [[x1, y1], [x2, y2]].forEach(function (p) { ctx.beginPath(); ctx.arc(p[0], p[1], 7, 0, 7); ctx.fill(); });
        },
      };
    },
    lorenz: function (c) {
      var sigma = num(c.sigma, 10), rho = num(c.rho, 28), beta = num(c.beta, 2.667);
      var x, y, z, pts;
      function reset() { x = 0.1; y = 0; z = 0; pts = []; }
      reset();
      return {
        init: reset,
        step: function (dt) {
          var h = 0.008;
          for (var i = 0; i < 4; i++) {
            var dx = sigma * (y - x), dy = x * (rho - z) - y, dz = x * y - beta * z;
            x += dx * h; y += dy * h; z += dz * h; pts.push([x, z]);
          }
          if (pts.length > 4000) pts.splice(0, pts.length - 4000);
        },
        draw: function (ctx, W, H) {
          ctx.strokeStyle = css("--accent"); ctx.lineWidth = 0.8;
          ctx.beginPath();
          for (var i = 0; i < pts.length; i++) {
            var px = W / 2 + pts[i][0] * (W / 60), py = H - pts[i][1] * (H / 55) - 8;
            i ? ctx.lineTo(px, py) : ctx.moveTo(px, py);
          }
          ctx.globalAlpha = 0.8; ctx.stroke(); ctx.globalAlpha = 1;
        },
      };
    },
    nbody: function (c) {
      var n = num(c.bodies, 5), G = num(c.g, 1), bodies;
      function reset() {
        bodies = [];
        var seed = 42;
        function rnd() { seed = (seed * 1103515245 + 12345) & 0x7fffffff; return seed / 0x7fffffff; }
        for (var i = 0; i < n; i++) bodies.push({ x: rnd(), y: rnd(), vx: (rnd() - 0.5) * 0.1, vy: (rnd() - 0.5) * 0.1, m: 0.5 + rnd() });
      }
      reset();
      return {
        init: reset,
        step: function (dt, env) {
          var W = env.W, H = env.H; dt = Math.min(dt, 0.03);
          for (var i = 0; i < bodies.length; i++) for (var j = i + 1; j < bodies.length; j++) {
            var a = bodies[i], b = bodies[j];
            var dx = (b.x - a.x) * W, dy = (b.y - a.y) * H, d2 = dx * dx + dy * dy + 400, d = Math.sqrt(d2);
            var f = G * a.m * b.m / d2;
            a.vx += f * dx / d / a.m * dt; a.vy += f * dy / d / a.m * dt;
            b.vx -= f * dx / d / b.m * dt; b.vy -= f * dy / d / b.m * dt;
          }
          bodies.forEach(function (p) {
            p.x += p.vx * dt / W * 60; p.y += p.vy * dt / H * 60;
            if (p.x < 0 || p.x > 1) p.vx *= -1; if (p.y < 0 || p.y > 1) p.vy *= -1;
            p.x = Math.max(0, Math.min(1, p.x)); p.y = Math.max(0, Math.min(1, p.y));
          });
        },
        draw: function (ctx, W, H) {
          bodies.forEach(function (p) {
            ctx.fillStyle = css("--accent"); ctx.beginPath();
            ctx.arc(p.x * W, p.y * H, 3 + p.m * 3, 0, 7); ctx.fill();
          });
        },
      };
    },
    life: function (c) {
      var cell = num(c.cell, 10), grid, cols, rows, acc = 0;
      function reset(env) {
        cols = Math.ceil((env ? env.W : 400) / cell); rows = Math.ceil((env ? env.H : 300) / cell);
        grid = []; var seed = 7;
        function rnd() { seed = (seed * 1103515245 + 12345) & 0x7fffffff; return seed / 0x7fffffff; }
        for (var i = 0; i < cols * rows; i++) grid.push(rnd() < 0.28 ? 1 : 0);
      }
      function idx(x, y) { return ((y + rows) % rows) * cols + ((x + cols) % cols); }
      return {
        init: function (env) { reset(env); },
        step: function (dt, env) {
          if (!grid || cols !== Math.ceil(env.W / cell)) reset(env);
          acc += dt; if (acc < 0.09) return; acc = 0;
          var next = grid.slice();
          for (var y = 0; y < rows; y++) for (var x = 0; x < cols; x++) {
            var nb = 0;
            for (var dy = -1; dy <= 1; dy++) for (var dx = -1; dx <= 1; dx++)
              if (dx || dy) nb += grid[idx(x + dx, y + dy)];
            var a = grid[idx(x, y)];
            next[idx(x, y)] = (a && (nb === 2 || nb === 3)) || (!a && nb === 3) ? 1 : 0;
          }
          grid = next;
        },
        draw: function (ctx, W, H) {
          if (!grid) return;
          ctx.fillStyle = css("--accent");
          for (var y = 0; y < rows; y++) for (var x = 0; x < cols; x++)
            if (grid[idx(x, y)]) ctx.fillRect(x * cell + 1, y * cell + 1, cell - 2, cell - 2);
        },
      };
    },
    vectorfield: function (c) {
      var fx = compile(c.fx || "y"), fy = compile(c.fy || "-x"), step = num(c.spacing, 34), sc = num(c.scale, 0.5);
      var SX = function (W) { return W / 8; }, SY = function (H) { return H / 8; };
      var parts = [];
      function spawnField(x, y) { parts.push({ x: x, y: y, trail: [] }); if (parts.length > 14) parts.shift(); }
      return {
        init: function () { parts = []; spawnField(2.6, 0.4); spawnField(-1.8, -1.4); },
        step: function (dt, env) {
          var W = env.W, H = env.H;
          if (env.mouse && env.mouse.clicked) {
            spawnField((env.mouse.x - W / 2) / SX(W), (H / 2 - env.mouse.y) / SY(H));
          }
          var h = Math.min(dt, 0.025) * 1.5;
          for (var i = parts.length - 1; i >= 0; i--) {
            var p = parts[i];
            var k1x = fx(p.x, p.y), k1y = fy(p.x, p.y); // RK2 midpoint
            var mx = p.x + k1x * h / 2, my = p.y + k1y * h / 2;
            p.x += fx(mx, my) * h; p.y += fy(mx, my) * h;
            p.trail.push([W / 2 + p.x * SX(W), H / 2 - p.y * SY(H)]);
            if (p.trail.length > 200) p.trail.shift();
            if (Math.abs(p.x) > 14 || Math.abs(p.y) > 14) parts.splice(i, 1);
          }
        },
        draw: function (ctx, W, H) {
          for (var py = step / 2; py < H; py += step) for (var px = step / 2; px < W; px += step) {
            var x = (px - W / 2) / SX(W), y = (H / 2 - py) / SY(H);
            var vx = fx(x, y), vy = fy(x, y), m = Math.hypot(vx, vy) || 1;
            var ux = vx / m, uy = -vy / m, L = Math.min(step * 0.45, m * sc * 8 + 4);
            ctx.strokeStyle = css("--text-faint"); ctx.globalAlpha = Math.min(0.55, 0.16 + m * 0.07); ctx.lineWidth = 1.2;
            ctx.beginPath(); ctx.moveTo(px, py); var ex = px + ux * L, ey = py + uy * L; ctx.lineTo(ex, ey);
            ctx.moveTo(ex, ey); ctx.lineTo(ex - ux * 4 - uy * 3, ey - uy * 4 + ux * 3);
            ctx.moveTo(ex, ey); ctx.lineTo(ex - ux * 4 + uy * 3, ey - uy * 4 - ux * 3);
            ctx.stroke();
          }
          ctx.globalAlpha = 1;
          for (var i = 0; i < parts.length; i++) {
            var tr = parts[i].trail; if (tr.length < 2) continue;
            ctx.strokeStyle = css("--accent"); ctx.lineWidth = 2; ctx.beginPath();
            for (var j = 0; j < tr.length; j++) j ? ctx.lineTo(tr[j][0], tr[j][1]) : ctx.moveTo(tr[j][0], tr[j][1]);
            ctx.stroke();
            var head = tr[tr.length - 1];
            ctx.fillStyle = css("--accent"); ctx.beginPath(); ctx.arc(head[0], head[1], 4, 0, 7); ctx.fill();
          }
        },
      };
    },
    wave: function (c) {
      var type = (c.type || "traveling"), n = num(c.harmonic, num(c.n, 3)),
        speed = num(c.speed, 1), amp = num(c.amplitude, 0.7), showNodes = (c.nodes !== "false");
      return {
        step: function () {},
        draw: function (ctx, W, H, env) {
          var t = env.t, pad = 24, midY = H / 2, A = amp * (H / 2 - pad);
          var yof = function (x) {
            return type === "standing"
              ? Math.sin(n * Math.PI * x) * Math.cos(2 * Math.PI * speed * t) * A
              : Math.sin(2 * Math.PI * (n * x - speed * t)) * A;
          };
          ctx.strokeStyle = css("--border"); ctx.lineWidth = 1;
          ctx.beginPath(); ctx.moveTo(pad, midY); ctx.lineTo(W - pad, midY); ctx.stroke();
          if (type === "standing") {
            ctx.strokeStyle = css("--text-faint"); ctx.globalAlpha = 0.5; ctx.setLineDash([3, 3]);
            for (var s = -1; s <= 1; s += 2) {
              ctx.beginPath();
              for (var px = pad; px <= W - pad; px++) {
                var xx = (px - pad) / (W - 2 * pad), Y = midY - s * Math.abs(Math.sin(n * Math.PI * xx)) * A;
                px === pad ? ctx.moveTo(px, Y) : ctx.lineTo(px, Y);
              }
              ctx.stroke();
            }
            ctx.setLineDash([]); ctx.globalAlpha = 1;
          }
          ctx.strokeStyle = css("--accent"); ctx.lineWidth = 2.5; ctx.beginPath();
          for (var p = pad; p <= W - pad; p++) {
            var x2 = (p - pad) / (W - 2 * pad), Y2 = midY - yof(x2);
            p === pad ? ctx.moveTo(p, Y2) : ctx.lineTo(p, Y2);
          }
          ctx.stroke();
          if (type === "standing" && showNodes) {
            ctx.fillStyle = css("--warn");
            for (var k = 0; k <= n; k++) {
              var nx = pad + (k / n) * (W - 2 * pad);
              ctx.beginPath(); ctx.arc(nx, midY, 3, 0, 7); ctx.fill();
            }
          }
        },
      };
    },
    interference: function (c) {
      var sep = num(c.separation, 0.34), wl = num(c.wavelength, 30), speed = num(c.speed, 2.4), cell = num(c.cell, 6);
      return {
        step: function () {},
        draw: function (ctx, W, H, env) {
          var t = env.t, m = env.mouse;
          var s = (m && m.x > 0) ? Math.max(0.08, Math.min(0.7, m.x / W)) : sep;
          var s1x = W * (0.5 - s / 2), s2x = W * (0.5 + s / 2), sy = H / 2, k = 2 * Math.PI / wl;
          var pos = rgb(css("--accent")), neg = rgb(css("--accent-2"));
          for (var y = 0; y < H; y += cell) for (var x = 0; x < W; x += cell) {
            var r1 = Math.hypot(x - s1x, y - sy), r2 = Math.hypot(x - s2x, y - sy);
            var v = Math.cos(k * r1 - speed * t) / (1 + r1 * 0.03) + Math.cos(k * r2 - speed * t) / (1 + r2 * 0.03);
            var a = Math.max(-1, Math.min(1, v)), col = a >= 0 ? pos : neg;
            ctx.fillStyle = "rgba(" + col[0] + "," + col[1] + "," + col[2] + "," + (Math.abs(a) * 0.85).toFixed(3) + ")";
            ctx.fillRect(x, y, cell, cell);
          }
          ctx.fillStyle = css("--warn");
          ctx.beginPath(); ctx.arc(s1x, sy, 4, 0, 7); ctx.fill();
          ctx.beginPath(); ctx.arc(s2x, sy, 4, 0, 7); ctx.fill();
        },
      };
    },
    wavepacket: function (c) {
      var groupv = num(c.groupv, 60), phasev = num(c.phasev, 160), k = num(c.k, 0.1), width = num(c.width, 70);
      return {
        step: function () {},
        draw: function (ctx, W, H, env) {
          var t = env.t, midY = H / 2, A = H * 0.32, span = W * 1.4,
            center = (W * 0.1 + groupv * t) % span, omega = phasev * k;
          ctx.strokeStyle = css("--text-faint"); ctx.globalAlpha = 0.6; ctx.setLineDash([4, 4]);
          for (var s = -1; s <= 1; s += 2) {
            ctx.beginPath();
            for (var px = 0; px < W; px++) {
              var e = Math.exp(-Math.pow((px - center) / width, 2)), Y = midY - s * e * A;
              px === 0 ? ctx.moveTo(px, Y) : ctx.lineTo(px, Y);
            }
            ctx.stroke();
          }
          ctx.setLineDash([]); ctx.globalAlpha = 1;
          ctx.strokeStyle = css("--accent"); ctx.lineWidth = 2; ctx.beginPath();
          for (var p = 0; p < W; p++) {
            var e2 = Math.exp(-Math.pow((p - center) / width, 2)), Y2 = midY - e2 * Math.cos(k * p - omega * t) * A;
            p === 0 ? ctx.moveTo(p, Y2) : ctx.lineTo(p, Y2);
          }
          ctx.stroke();
        },
      };
    },
  };

  /* ---------------- helpers ---------------- */
  function compile(src) { try { return new Function("x", "y", "with (Math) { return (" + src + "); }"); } catch (e) { return function () { return 0; }; } }
  function num(v, d) { var n = parseFloat(v); return isNaN(n) ? (d === undefined ? 0 : d) : n; }
  function rgb(s) {
    s = (s || "").trim();
    if (/^#[0-9a-f]{3}$/i.test(s)) return [s[1] + s[1], s[2] + s[2], s[3] + s[3]].map(function (h) { return parseInt(h, 16); });
    if (/^#[0-9a-f]{6}$/i.test(s)) return [s.substr(1, 2), s.substr(3, 2), s.substr(5, 2)].map(function (h) { return parseInt(h, 16); });
    var m = s.match(/(\d+)[,\s]+(\d+)[,\s]+(\d+)/);
    return m ? [+m[1], +m[2], +m[3]] : [120, 140, 200];
  }
  function fmt(v) { return Math.abs(v) >= 1000 || (v !== 0 && Math.abs(v) < 0.01) ? v.toExponential(1) : (Math.round(v * 100) / 100).toString(); }
})();
