import Lenis from 'lenis';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

/* ─── LENIS SMOOTH SCROLL ─── */
const lenis = new Lenis({
  duration: 1.2,
  easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
  orientation: 'vertical',
  smoothWheel: true,
});

lenis.on('scroll', ScrollTrigger.update);
gsap.ticker.add((time) => lenis.raf(time * 1000));
gsap.ticker.lagSmoothing(0);

ScrollTrigger.config({ ignoreMobileResize: true });

/* ─── HEADER SCROLL STATE ─── */
const header = document.querySelector('.site-header');
if (header) {
  ScrollTrigger.create({
    start: 'top -80',
    onUpdate: (self) => header.classList.toggle('scrolled', self.progress > 0),
  });
}

/* ─── SCROLL TO TOP ─── */
const scrollTopBtn = document.querySelector('.scroll-top');
if (scrollTopBtn) {
  ScrollTrigger.create({
    start: 'top -200',
    onUpdate: (self) => scrollTopBtn.classList.toggle('visible', self.progress > 0),
  });
  scrollTopBtn.addEventListener('click', () => lenis.scrollTo(0, { duration: 1.2 }));
}

/* ─── CANVAS PARTICLE NETWORK ─── */
;(function initCanvas() {
  const canvas = document.getElementById('particle-canvas');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  let W, H, particles = [];
  const COUNT = 80, CONNECT_DIST = 140, MOUSE_RADIUS = 120;

  let mouse = { x: -9999, y: -9999 };
  let mouseMoved = false;

  function resize() {
    const rect = canvas.parentElement.getBoundingClientRect();
    W = canvas.width = rect.width * window.devicePixelRatio;
    H = canvas.height = rect.height * window.devicePixelRatio;
    canvas.style.width = rect.width + 'px';
    canvas.style.height = rect.height + 'px';
    ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
    W = rect.width;
    H = rect.height;
  }

  class Particle {
    constructor() {
      this.x = Math.random() * W;
      this.y = Math.random() * H;
      this.vx = (Math.random() - 0.5) * 0.3;
      this.vy = (Math.random() - 0.5) * 0.3;
      this.radius = Math.random() * 2 + 0.5;
      this.isTeal = Math.random() > 0.5;
    }
    update() {
      this.x += this.vx;
      this.y += this.vy;
      if (this.x < 0 || this.x > W) this.vx *= -1;
      if (this.y < 0 || this.y > H) this.vy *= -1;

      if (mouseMoved) {
        const dx = this.x - mouse.x;
        const dy = this.y - mouse.y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        if (dist < MOUSE_RADIUS) {
          const force = (MOUSE_RADIUS - dist) / MOUSE_RADIUS;
          this.x += dx * force * 0.03;
          this.y += dy * force * 0.03;
        }
      }
    }
    draw() {
      ctx.beginPath();
      ctx.arc(this.x, this.y, this.radius, 0, Math.PI * 2);
      ctx.fillStyle = this.isTeal ? 'rgba(13,148,136,0.6)' : 'rgba(124,58,237,0.5)';
      ctx.fill();
    }
  }

  function initParticles() {
    particles = [];
    for (let i = 0; i < COUNT; i++) particles.push(new Particle());
  }

  function drawConnections() {
    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const dx = particles[i].x - particles[j].x;
        const dy = particles[i].y - particles[j].y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        if (dist < CONNECT_DIST) {
          const alpha = (1 - dist / CONNECT_DIST) * 0.3;
          const isCross = particles[i].isTeal !== particles[j].isTeal;
          ctx.beginPath();
          ctx.moveTo(particles[i].x, particles[i].y);
          ctx.lineTo(particles[j].x, particles[j].y);
          ctx.strokeStyle = isCross
            ? `rgba(124,58,237,${alpha * 0.7})`
            : particles[i].isTeal
              ? `rgba(13,148,136,${alpha})`
              : `rgba(124,58,237,${alpha})`;
          ctx.lineWidth = 0.5;
          ctx.stroke();
        }
      }
    }
  }

  function animate() {
    ctx.clearRect(0, 0, W, H);
    particles.forEach((p) => { p.update(); p.draw(); });
    drawConnections();
    requestAnimationFrame(animate);
  }

  resize();
  initParticles();
  animate();

  window.addEventListener('resize', () => {
    resize();
    initParticles();
  });

  canvas.addEventListener('mousemove', (e) => {
    const rect = canvas.getBoundingClientRect();
    mouse.x = e.clientX - rect.left;
    mouse.y = e.clientY - rect.top;
    mouseMoved = true;
  });
  canvas.addEventListener('mouseleave', () => {
    mouse.x = -9999; mouse.y = -9999;
    mouseMoved = false;
  });
})();

/* ─── GSAP HERO ENTRY ─── */
;(function heroEntry() {
  const tl = gsap.timeline({ defaults: { ease: 'power3.out' } });

  tl.from('.hero-tag', { opacity: 0, y: 20, duration: 0.6 })
    .from('.hero-headline .line', { opacity: 0, y: 40, duration: 0.8, stagger: 0.15 }, '-=0.3')
    .from('.hero-sub', { opacity: 0, y: 20, duration: 0.6 }, '-=0.4')
    .from('.hero-cta-group', { opacity: 0, y: 20, duration: 0.6 }, '-=0.3')
    .from('.hero-trust', { opacity: 0, y: 20, duration: 0.5 }, '-=0.2');
})();

/* ─── STATS COUNTER ─── */
;(function statsCounter() {
  document.querySelectorAll('.stat-number').forEach((el) => {
    const target = parseInt(el.dataset.count, 10);
    if (target && !isNaN(target)) {
      ScrollTrigger.create({
        trigger: el,
        start: 'top 85%',
        onEnter: () => {
          gsap.to(el, {
            duration: 2,
            innerText: target,
            ease: 'power2.out',
            snap: { innerText: 1 },
            onUpdate: () => {
              const val = parseInt(el.textContent.replace(/,/g, ''), 10);
              if (val >= 1000) el.textContent = (val / 1000).toFixed(1) + 'k';
            },
          });
          ScrollTrigger.getById(el.dataset.statsId)?.kill();
        },
      });
    }
  });
})();

/* ─── FEATURES REVEAL ─── */
;(function featuresReveal() {
  const cards = gsap.utils.toArray('.feature-card');
  cards.forEach((card, i) => {
    gsap.from(card, {
      opacity: 0,
      y: 40,
      rotateX: gsap.utils.random(-3, 3),
      duration: 0.8,
      ease: 'power3.out',
      scrollTrigger: {
        trigger: card,
        start: 'top 88%',
        toggleActions: 'play none none none',
      },
    });
  });
})();

/* ─── WORKFLOW SCROLL PROGRESSION ─── */
;(function workflowSteps() {
  const items = gsap.utils.toArray('.workflow-visual-item');
  const counts = gsap.utils.toArray('.workflow-step-count');
  if (items.length === 0) return;

  items.forEach((item, i) => {
    ScrollTrigger.create({
      trigger: item,
      start: 'top 40%',
      end: 'bottom 40%',
      onEnter: () => {
        items.forEach((el) => el.classList.remove('active'));
        counts.forEach((el) => el.classList.remove('active'));
        item.classList.add('active');
        if (counts[i]) counts[i].classList.add('active');
      },
      onEnterBack: () => {
        items.forEach((el) => el.classList.remove('active'));
        counts.forEach((el) => el.classList.remove('active'));
        item.classList.add('active');
        if (counts[i]) counts[i].classList.add('active');
      },
    });
  });
})();

/* ─── COMMANDS REVEAL ─── */
;(function commandsReveal() {
  gsap.utils.toArray('.cmd-card').forEach((card, i) => {
    gsap.from(card, {
      opacity: 0,
      y: 30,
      duration: 0.6,
      ease: 'power3.out',
      delay: i * 0.05,
      scrollTrigger: {
        trigger: card,
        start: 'top 90%',
        toggleActions: 'play none none none',
      },
    });
  });
})();

/* ─── FAQ ACCORDION ─── */
;(function faq() {
  document.querySelectorAll('.faq-q').forEach((q) => {
    q.addEventListener('click', () => {
      const expanded = q.getAttribute('aria-expanded') === 'true';
      const answer = document.getElementById(q.getAttribute('aria-controls'));
      document.querySelectorAll('.faq-q').forEach((b) => {
        if (b !== q) {
          b.setAttribute('aria-expanded', 'false');
          const a = document.getElementById(b.getAttribute('aria-controls'));
          if (a) a.classList.remove('open');
        }
      });
      q.setAttribute('aria-expanded', expanded ? 'false' : 'true');
      if (answer) answer.classList.toggle('open', !expanded);
    });
  });
})();

/* ─── CUSTOM CURSOR ─── */
;(function cursor() {
  if (window.matchMedia('(pointer: coarse)').matches) return;
  const dot = document.querySelector('.cursor-dot');
  const ring = document.querySelector('.cursor-ring');
  if (!dot || !ring) return;

  let mouseX = -100, mouseY = -100;
  let ringX = -100, ringY = -100;

  document.addEventListener('mousemove', (e) => {
    mouseX = e.clientX;
    mouseY = e.clientY;
    dot.style.left = mouseX + 'px';
    dot.style.top = mouseY + 'px';
  });

  gsap.ticker.add(() => {
    ringX += (mouseX - ringX) * 0.12;
    ringY += (mouseY - ringY) * 0.12;
    ring.style.left = ringX + 'px';
    ring.style.top = ringY + 'px';
  });

  document.querySelectorAll('a, button, .feature-card, .cmd-card, .faq-q').forEach((el) => {
    el.addEventListener('mouseenter', () => ring.classList.add('hovering'));
    el.addEventListener('mouseleave', () => ring.classList.remove('hovering'));
  });

  document.addEventListener('mouseleave', () => ring.classList.add('hidden'));
  document.addEventListener('mouseenter', () => ring.classList.remove('hidden'));
})();

/* ─── COPY BUTTON ─── */
;(function copy() {
  document.querySelectorAll('.btn-copy').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const code = btn.dataset.cmd || '';
      try {
        await navigator.clipboard.writeText(code);
      } catch {
        const ta = document.createElement('textarea');
        ta.value = code;
        document.body.appendChild(ta);
        ta.select();
        document.execCommand('copy');
        document.body.removeChild(ta);
      }
      btn.textContent = 'Copied!';
      btn.classList.add('copied');
      setTimeout(() => {
        btn.textContent = 'Copy →';
        btn.classList.remove('copied');
      }, 2000);
    });
  });
})();

/* ─── LENIS ANCHOR LINKS ─── */
document.querySelectorAll('a[href^="#"]').forEach((a) => {
  a.addEventListener('click', (e) => {
    const target = document.querySelector(a.getAttribute('href'));
    if (target) {
      e.preventDefault();
      lenis.scrollTo(target, { duration: 1.2 });
    }
  });
});

/* ─── CLEANUP ON PAGE NAV ─── */
window.addEventListener('beforeunload', () => {
  ScrollTrigger.getAll().forEach((t) => t.kill());
  lenis.destroy();
});

/* ─── REDUCED MOTION ─── */
if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
  lenis.destroy();
  ScrollTrigger.getAll().forEach((t) => t.disable());
}
