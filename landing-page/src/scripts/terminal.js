document.addEventListener('DOMContentLoaded', () => {

  /* ── Typewriter ────────────────────────────────────────────── */
  ;(function typewriter() {
    const cmdEl = document.getElementById('install-command-text');
    if (!cmdEl) return;
    const fullCmd = cmdEl.textContent || '';
    cmdEl.textContent = '';
    cmdEl.style.opacity = '1';
    let i = 0;
    const typeInterval = 30;
    function tick() {
      if (i < fullCmd.length) {
        cmdEl.textContent += fullCmd.charAt(i);
        i++;
        setTimeout(tick, typeInterval);
      } else {
        cmdEl.dispatchEvent(new CustomEvent('typewriter:done'));
      }
    }
    const startDelay = 500;
    setTimeout(tick, startDelay);

    const progressFill = document.querySelector('.terminal-progress-fill');
    const progressText = document.querySelector('.terminal-progress-text');
    const outputLines = document.querySelectorAll('.terminal-line');
    const completeBadge = document.querySelector('.terminal-complete-badge');

    cmdEl.addEventListener('typewriter:done', () => {
      /* Progress bar animation */
      if (progressFill) {
        let pct = 0;
        const progInterval = setInterval(() => {
          pct += 2;
          if (progressFill) progressFill.style.width = Math.min(pct, 100) + '%';
          if (progressText) progressText.textContent = pct + '%';
          if (pct >= 100) {
            clearInterval(progInterval);
            /* Show output lines */
            outputLines.forEach((line, idx) => {
              setTimeout(() => {
                line.style.animationPlayState = 'running';
              }, idx * 150);
            });
            /* Show complete badge */
            if (completeBadge) {
              setTimeout(() => {
                completeBadge.style.opacity = '1';
                completeBadge.style.transform = 'scale(1)';
              }, outputLines.length * 150 + 200);
            }
          }
        }, 20);
      }
    });
  })();

  /* ── Scroll Reveal Observer ──────────────────────────────── */
  ;(function scrollReveal() {
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return;
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          observer.unobserve(entry.target);
        }
      });
    }, { threshold: 0.15, rootMargin: '0px 0px -60px 0px' });

    document.querySelectorAll('.terminal-section-header, .feature-terminal-item, .command-card').forEach(el => {
      observer.observe(el);
    });
  })();

  /* ── Sticky Workflow Steps ───────────────────────────────── */
  ;(function stickyWorkflow() {
    const steps = document.querySelectorAll('.workflow-step');
    if (steps.length === 0) return;
    const stepObserver = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const idx = Array.from(steps).indexOf(entry.target);
          steps.forEach((s, i) => {
            s.classList.toggle('active', i <= idx);
          });
        }
      });
    }, { threshold: 0.3 });

    steps.forEach(s => stepObserver.observe(s));
  })();

  /* ── Copy Button ──────────────────────────────────────────── */
  ;(function copyButton() {
    const btn = document.getElementById('cmd-copy-btn');
    if (!btn) return;
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
        btn.textContent = 'Copy';
        btn.classList.remove('copied');
      }, 2000);
    });
  })();

  /* ── FAQ Accordion ───────────────────────────────────────── */
  ;(function faqAccordion() {
    document.addEventListener('click', (e) => {
      const q = e.target.closest('.faq-terminal-q');
      if (!q) return;
      const expanded = q.getAttribute('aria-expanded') === 'true';
      const answer = document.getElementById(q.getAttribute('aria-controls'));
      document.querySelectorAll('.faq-terminal-q').forEach(b => {
        if (b !== q) {
          b.setAttribute('aria-expanded', 'false');
          const a = document.getElementById(b.getAttribute('aria-controls'));
          if (a) a.classList.remove('open');
        }
      });
      q.setAttribute('aria-expanded', expanded ? 'false' : 'true');
      if (answer) answer.classList.toggle('open', !expanded);
    });
  })();

  /* ── Dark Mode Toggle ────────────────────────────────────── */
  ;(function darkMode() {
    const btn = document.getElementById('theme-toggle-terminal');
    if (!btn) return;
    function apply(dark) {
      document.documentElement.classList.toggle('dark', dark);
      localStorage.setItem('theme', dark ? 'dark' : 'light');
    }
    function initDark() {
      const stored = localStorage.getItem('theme');
      const prefers = window.matchMedia('(prefers-color-scheme: dark)').matches;
      apply(stored === 'dark' || (!stored && prefers));
    }
    btn.addEventListener('click', () => {
      const isDark = !document.documentElement.classList.contains('dark');
      apply(isDark);
    });
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      if (!localStorage.getItem('theme')) apply(e.matches);
    });
    initDark();
  })();

});
