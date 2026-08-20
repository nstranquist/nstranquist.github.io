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

    function clearPending(row) {
      var job = pending.get(row);
      if (!job) {
        return;
      }
      if (job.timer) {
        clearTimeout(job.timer);
      }
      if (job.detail && job.onEnd) {
        job.detail.removeEventListener("transitionend", job.onEnd);
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

      if (!detail) {
        hideDetail(row, detail);
        return;
      }

      if (reduceMotion() || detail.hidden) {
        hideDetail(row, detail);
        return;
      }

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
        if (e.target === detail) {
          finish();
        }
      }
      detail.addEventListener("transitionend", onEnd);
      pending.set(row, {
        detail: detail,
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
        if (e.target.closest(".catalog-toggle")) {
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
