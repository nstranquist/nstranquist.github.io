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

    function detailFor(row) {
      var btn = row.querySelector(".catalog-toggle");
      var id = btn && btn.getAttribute("aria-controls");
      return id ? document.getElementById(id) : null;
    }

    function close(row) {
      if (!row) {
        return;
      }
      var btn = row.querySelector(".catalog-toggle");
      var detail = detailFor(row);
      if (btn) {
        btn.setAttribute("aria-expanded", "false");
      }
      if (detail) {
        detail.hidden = true;
      }
      row.classList.remove("is-open");
      if (location.hash === "#" + row.id) {
        history.replaceState(null, "", location.pathname + location.search);
      }
      if (openRow === row) {
        openRow = null;
      }
    }

    function open(row, scrollIntoView) {
      if (!row) {
        return;
      }
      if (openRow && openRow !== row) {
        close(openRow);
      }
      var btn = row.querySelector(".catalog-toggle");
      var detail = detailFor(row);
      if (btn) {
        btn.setAttribute("aria-expanded", "true");
      }
      if (detail) {
        detail.hidden = false;
      }
      row.classList.add("is-open");
      if (location.hash !== "#" + row.id) {
        history.replaceState(null, "", "#" + row.id);
      }
      openRow = row;
      if (scrollIntoView) {
        var reduce = window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches;
        row.scrollIntoView({ block: "start", behavior: reduce ? "auto" : "smooth" });
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
        if (e.target.closest("a") || e.target.closest(".catalog-toggle")) {
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
