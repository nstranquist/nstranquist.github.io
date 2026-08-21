(function () {
  var KEY = "nicos-catalog-theme";
  var root = document.documentElement;
  var order = ["system", "dark", "light"];
  var labels = { system: "System", dark: "Dark", light: "Light" };

  function theme() {
    return root.getAttribute("data-theme") || "system";
  }

  function apply(next) {
    root.setAttribute("data-theme", next);
    var btn = document.querySelector("[data-theme-toggle]");
    if (btn) {
      btn.setAttribute("aria-label", "Color theme: " + labels[next] + ". Switch theme");
      btn.textContent = labels[next];
    }
  }

  apply(localStorage.getItem(KEY) || "system");

  var toggle = document.querySelector("[data-theme-toggle]");
  if (toggle) {
    toggle.addEventListener("click", function () {
      var next = order[(order.indexOf(theme()) + 1) % order.length];
      localStorage.setItem(KEY, next);
      apply(next);
    });
  }

  (function aggregateArrivals() {
    var endpoint = "https://profile-arrivals.nstranquist.workers.dev/e";
    var doNotTrack = navigator.doNotTrack || window.doNotTrack || navigator.msDoNotTrack;
    if (navigator.globalPrivacyControl === true || doNotTrack === "1" || doNotTrack === "yes") {
      return;
    }

    function send(payload) {
      var body = JSON.stringify(payload);
      if (navigator.sendBeacon && navigator.sendBeacon(endpoint, body)) {
        return;
      }
      if (window.fetch) {
        window.fetch(endpoint, {
          method: "POST",
          body: body,
          headers: { "content-type": "text/plain;charset=UTF-8" },
          keepalive: true,
          credentials: "omit",
          referrerPolicy: "no-referrer",
        }).catch(function () {});
      }
    }

    if (location.pathname === "/" || location.pathname === "/index.html") {
      send({ event: "pageview", surface: "home" });
    }

    document.addEventListener("click", function (event) {
      var link = event.target.closest && event.target.closest("[data-arrival-event]");
      if (!link) {
        return;
      }
      var payload = {
        event: link.getAttribute("data-arrival-event"),
        surface: link.getAttribute("data-arrival-surface"),
      };
      var product = link.getAttribute("data-arrival-product");
      if (product) {
        payload.product = product;
      }
      send(payload);
    });
  })();

  (function catalogAccordion() {
    var items = document.querySelectorAll(".catalog-item");
    if (!items.length) {
      return;
    }

    var openRow = null;
    var pending = new WeakMap();
    var CLOSE_MS = 280;

    function reduceMotion() {
      return window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    }

    function detailFor(row) {
      var btn = row.querySelector(".catalog-toggle");
      var id = btn && btn.getAttribute("aria-controls");
      return id ? document.getElementById(id) : null;
    }

    function revealFor(detail) {
      return detail ? detail.querySelector(".catalog-reveal") : null;
    }

    function clearPending(row) {
      var job = pending.get(row);
      if (!job) {
        return;
      }
      if (job.timer) {
        clearTimeout(job.timer);
      }
      if (job.node && job.onEnd) {
        job.node.removeEventListener("transitionend", job.onEnd);
      }
      pending.delete(row);
    }

    function syncHash(row, opening) {
      if (!row || !row.id) {
        return;
      }
      if (opening) {
        if (location.hash !== "#" + row.id) {
          history.replaceState(null, "", "#" + row.id);
        }
        return;
      }
      if (location.hash === "#" + row.id) {
        history.replaceState(null, "", location.pathname + location.search);
      }
    }

    function hideDetail(row, detail) {
      if (detail) {
        detail.hidden = true;
      }
      if (openRow === row) {
        openRow = null;
      }
    }

    function close(row) {
      if (!row) {
        return;
      }
      clearPending(row);
      var btn = row.querySelector(".catalog-toggle");
      var detail = detailFor(row);
      if (btn) {
        btn.setAttribute("aria-expanded", "false");
      }
      row.classList.remove("is-open");
      syncHash(row, false);

      if (!detail || detail.hidden || reduceMotion()) {
        hideDetail(row, detail);
        return;
      }

      var reveal = revealFor(detail);
      var done = false;
      function finish() {
        if (done) {
          return;
        }
        done = true;
        clearPending(row);
        hideDetail(row, detail);
      }
      function onEnd(e) {
        if (e.target === reveal && e.propertyName === "grid-template-rows") {
          finish();
        }
      }
      if (reveal) {
        reveal.addEventListener("transitionend", onEnd);
      }
      pending.set(row, {
        node: reveal,
        onEnd: onEnd,
        timer: setTimeout(finish, CLOSE_MS)
      });
    }

    function open(row, scrollIntoView) {
      if (!row) {
        return;
      }
      if (openRow && openRow !== row) {
        close(openRow);
      }
      clearPending(row);
      var btn = row.querySelector(".catalog-toggle");
      var detail = detailFor(row);
      if (btn) {
        btn.setAttribute("aria-expanded", "true");
      }
      if (detail) {
        detail.hidden = false;
      }
      function markOpen() {
        row.classList.add("is-open");
      }
      if (reduceMotion()) {
        markOpen();
      } else {
        requestAnimationFrame(markOpen);
      }
      syncHash(row, true);
      openRow = row;
      if (scrollIntoView) {
        row.scrollIntoView({ block: "start", behavior: reduceMotion() ? "auto" : "smooth" });
      }
    }

    function toggleRow(row) {
      if (row.classList.contains("is-open")) {
        close(row);
      } else {
        open(row, false);
      }
    }

    items.forEach(function (row) {
      var btn = row.querySelector(".catalog-toggle");
      if (btn) {
        btn.addEventListener("click", function (e) {
          e.stopPropagation();
          toggleRow(row);
        });
      }
      row.addEventListener("click", function (e) {
        if (e.target.closest("a")) {
          return;
        }
        if (e.target.closest(".catalog-toggle") || e.target.closest(".sort")) {
          return;
        }
        if (btn) {
          btn.click();
        }
      });
    });

    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape" && openRow) {
        close(openRow);
      }
    });

    function openFromHash() {
      var hash = location.hash.replace(/^#/, "");
      if (!hash) {
        return;
      }
      var row = document.getElementById(hash);
      if (row && row.classList.contains("catalog-item")) {
        open(row, true);
      }
    }

    window.addEventListener("hashchange", openFromHash);
    openFromHash();
  })();


  (function catalogSort() {
    var board = document.querySelector(".catalog-board");
    var listEl = document.querySelector(".catalog-list");
    if (!board || !listEl) {
      return;
    }
    var buttons = board.querySelectorAll(".sort");
    if (!buttons.length) {
      return;
    }

    var key = "";
    var dir = "none";

    function rows() {
      return Array.prototype.slice.call(listEl.querySelectorAll(".catalog-item"));
    }

    function versionParts(value) {
      var m = /^v?(\d+)\.(\d+)\.(\d+)/i.exec(value || "");
      if (!m) {
        return null;
      }
      return [parseInt(m[1], 10), parseInt(m[2], 10), parseInt(m[3], 10)];
    }

    function emptyProof(value) {
      var v = String(value || "").trim();
      return !v || /no public tag/i.test(v);
    }

    function cmpProof(a, b) {
      var lastA = emptyProof(a);
      var lastB = emptyProof(b);
      if (lastA !== lastB) {
        return lastA ? 1 : -1;
      }
      var va = versionParts(a);
      var vb = versionParts(b);
      if (va && vb) {
        for (var i = 0; i < 3; i++) {
          if (va[i] !== vb[i]) {
            return va[i] - vb[i];
          }
        }
        return 0;
      }
      if (va && !vb) {
        return -1;
      }
      if (!va && vb) {
        return 1;
      }
      return String(a || "").localeCompare(String(b || ""), undefined, { numeric: true, sensitivity: "base" });
    }

    function valueOf(row, sortKey) {
      return row.getAttribute("data-" + sortKey) || "";
    }

    function compare(sortKey, a, b) {
      var result;
      if (sortKey === "index") {
        result = (Number(a.getAttribute("data-index")) || 0) - (Number(b.getAttribute("data-index")) || 0);
      } else if (sortKey === "proof") {
        result = cmpProof(valueOf(a, "proof"), valueOf(b, "proof"));
      } else {
        result = valueOf(a, sortKey).localeCompare(valueOf(b, sortKey), undefined, {
          numeric: true,
          sensitivity: "base"
        });
      }
      if (result === 0) {
        return (Number(a.getAttribute("data-index")) || 0) - (Number(b.getAttribute("data-index")) || 0);
      }
      return result;
    }

    function headerOf(btn) {
      return btn.closest("[aria-sort]") || btn.parentElement;
    }

    function setAria(active, nextDir) {
      buttons.forEach(function (btn) {
        var header = headerOf(btn);
        if (!header) {
          return;
        }
        if (btn === active && nextDir === "asc") {
          header.setAttribute("aria-sort", "ascending");
        } else if (btn === active && nextDir === "desc") {
          header.setAttribute("aria-sort", "descending");
        } else {
          header.setAttribute("aria-sort", "none");
        }
      });
    }

    function apply(list) {
      var i;
      for (i = 0; i < list.length; i++) {
        var ref = listEl.children[i];
        if (ref !== list[i]) {
          listEl.insertBefore(list[i], ref);
        }
      }
    }

    buttons.forEach(function (btn) {
      btn.addEventListener("click", function (e) {
        e.stopPropagation();
        e.preventDefault();
        var next = btn.getAttribute("data-sort");
        if (!next) {
          return;
        }
        if (key === next) {
          dir = dir === "asc" ? "desc" : dir === "desc" ? "none" : "asc";
          if (dir === "none") {
            key = "";
          }
        } else {
          key = next;
          dir = "asc";
        }

        var list = rows();
        var sortKey = dir === "none" ? "index" : key;
        list.sort(function (a, b) {
          if (sortKey === "proof") {
            var emptyA = emptyProof(valueOf(a, "proof"));
            var emptyB = emptyProof(valueOf(b, "proof"));
            if (emptyA !== emptyB) {
              return emptyA ? 1 : -1;
            }
          }
          var result = compare(sortKey, a, b);
          return dir === "desc" ? -result : result;
        });
        setAria(btn, dir);
        apply(list);
      });
    });
  })();

  var nav = document.querySelectorAll("nav [data-section]");
  if (!nav.length || !("IntersectionObserver" in window)) {
    return;
  }

  var targets = [];
  nav.forEach(function (link) {
    var id = link.getAttribute("data-section");
    var el = id && document.getElementById(id);
    if (el) {
      targets.push({ id: id, el: el, link: link });
    }
  });
  if (!targets.length) {
    return;
  }

  var current = "";
  function setCurrent(id) {
    if (current === id) {
      return;
    }
    current = id;
    nav.forEach(function (link) {
      if (link.getAttribute("data-section") === id) {
        link.setAttribute("aria-current", "true");
      } else {
        link.removeAttribute("aria-current");
      }
    });
  }

  var observer = new IntersectionObserver(
    function (entries) {
      var visible = entries
        .filter(function (e) {
          return e.isIntersecting;
        })
        .sort(function (a, b) {
          return b.intersectionRatio - a.intersectionRatio;
        });
      if (visible[0] && visible[0].target && visible[0].target.id) {
        setCurrent(visible[0].target.id);
      }
    },
    { rootMargin: "-20% 0px -55% 0px", threshold: [0.15, 0.35, 0.6] }
  );

  targets.forEach(function (t) {
    observer.observe(t.el);
  });
})();
