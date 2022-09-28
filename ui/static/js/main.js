tsParticles.load("tsparticles", {
  fullScreen: { enable: !1 },
  pauseOnBlur: { enable: !0 },
  pauseOnOutsideViewport: { enable: !0 },
  particles: {
    color: { value: ["#FFFFFF", "#FFd700"] },
    move: {
      direction: "bottom",
      enable: !0,
      outModes: { default: "out" },
      size: !0,
      speed: { min: 1, max: 3 },
    },
    number: { value: 100, density: { enable: !0, area: 800 } },
    opacity: {
      value: 1,
      animation: {
        enable: !1,
        startValue: "max",
        destroy: "min",
        speed: 0.3,
        sync: !0,
      },
    },
    rotate: {
      value: { min: 0, max: 0 },
      direction: "random",
      move: !1,
      animation: { enable: !0, speed: 60 },
      enable: !1,
    },
    tilt: {
      direction: "random",
      enable: !1,
      move: !0,
      value: { min: 0, max: 360 },
      animation: { enable: !0, speed: 60 },
    },
    shape: {
      type: "image",
      options: {
        image: [
          {
            src: "/static/images/js.svg",
            width: 32,
            height: 32,
            particles: { size: { value: { min: 3, max: 10 } } },
          },
          {
            src: "/static/images/html.svg",
            width: 32,
            height: 32,
            particles: { size: { value: { min: 3, max: 10 } } },
          },
          {
            src: "/static/images/css.svg",
            width: 32,
            height: 32,
            particles: { size: { value: { min: 3, max: 10 } } },
          },
          {
            src: "/static/images/gopher.svg",
            width: 32,
            height: 32,
            particles: { size: { value: { min: 3, max: 10 } } },
          },
        ],
      },
    },
    size: { value: { min: 2, max: 4 } },
    roll: {
      darken: { enable: !0, value: 30 },
      enlighten: { enable: !0, value: 30 },
      enable: !1,
      speed: { min: 15, max: 25 },
    },
    wobble: {
      distance: 30,
      enable: !0,
      move: !0,
      speed: { min: -15, max: 15 },
    },
  },
}),
  tsParticles.setOnClickHandler((e, a) => {});
const particles = tsParticles.domItem(0);

(function (e) {
  var a = e(window),
    i = e("body"),
    s = e("#header"),
    t = null,
    l = e("#nav");
  e("#wrapper");
  breakpoints({
    xlarge: ["1281px", "880px"],
    large: ["1025px", "1280px"],
    medium: ["737px", "1024px"],
    small: ["481px", "736px"],
    xsmall: [null, "480px"],
  }),
    a.on("load", function () {
      window.setTimeout(function () {
        i.removeClass("is-preload");
      }, 100);
    }),
    browser.canUse("object-fit") ||
      e(".image[data-position]").each(function () {
        var a = e(this),
          i = a.children("img");
        a
          .css("background-image", 'url("' + i.attr("src") + '")')
          .css("background-position", a.data("position"))
          .css("background-size", "cover")
          .css("background-repeat", "no-repeat"),
          i.css("opacity", "0");
      });
  var n = l.find("a");
  n
    .addClass("scrolly")
    .on("click", function () {
      var a = e(this);
      "#" == a.attr("href").charAt(0) &&
        (n.removeClass("active"),
        a.addClass("active").addClass("active-locked"));
    })
    .each(function () {
      var a = e(this),
        i = a.attr("href"),
        s = e(i);
      s.length < 1 ||
        s.scrollex({
          mode: "middle",
          top: "5vh",
          bottom: "5vh",
          initialize: function () {
            s.addClass("inactive");
          },
          enter: function () {
            s.removeClass("inactive"),
              0 == n.filter(".active-locked").length
                ? (n.removeClass("active"), a.addClass("active"))
                : a.hasClass("active-locked") && a.removeClass("active-locked");
          },
        });
    }),
    (t = e(
      '<div id="titleBar"><a href="#header" class="toggle"></a><span class="title">' +
        e("#logo").html() +
        "</span></div>"
    ).appendTo(i)),
    s.panel({
      delay: 500,
      hideOnClick: !0,
      hideOnSwipe: !0,
      resetScroll: !0,
      resetForms: !0,
      side: "right",
      target: i,
      visibleClass: "header-visible",
    }),
    e(".scrolly").scrolly({
      speed: 1e3,
      offset: function () {
        return breakpoints.active("<=medium") ? t.height() : 0;
      },
    });
})(jQuery);
