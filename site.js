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
