document.querySelectorAll('a[href^="#"]').forEach(link => {

    link.addEventListener('click', event => {

        const target = link.getAttribute('href');

        // Ignorar enlaces vacíos o controlados por JS (p. ej. #btn-tour).
        if (!target || target === '#') return;

        event.preventDefault();

        document.querySelector(target)?.scrollIntoView({
            behavior:'smooth'
        });

    });

});

particlesJS('particles-container', {
    particles: {
      number: { value: 60 },
      color: { value: '#ffffff' },
      shape: { type: 'circle' },
      opacity: { value: 0.3 },
      size: { value: 2.5 },
      line_linked: {
        enable: true,
        distance: 150,
        color: '#ffffff',
        opacity: 0.2,
        width: 1
      },
      move: { enable: true, speed: 1.5 }
    },
    interactivity: {
      detect_on: 'canvas',
      events: {
        onhover: { enable: true, mode: 'grab' },
        onclick: { enable: true, mode: 'push' }
      }
    }
});


// MENU (delegado: el header se inyecta por HTMX y no existe al cargar)
document.addEventListener('click', (e) => {
    if (e.target.closest('#hamburger')) {
        document.getElementById('hamburger')?.classList.toggle('open');
        document.querySelector('.menu')?.classList.toggle('open');
        return;
    }
    if (e.target.closest('.menu a')) {
        document.getElementById('hamburger')?.classList.remove('open');
        document.querySelector('.menu')?.classList.remove('open');
    }
});

window.addEventListener('scroll', function() {
    const menu = document.getElementById('main-menu');

    // Si el scroll vertical es mayor a 50px, añade la clase; si no, la quita
    if (window.scrollY > 50) {
        menu.classList.add('scrolled');
    } else {
        menu.classList.remove('scrolled');
    }
});

// AVATAR: reejecutar el vídeo de saludo al pasar el ratón
document.querySelector('.avatar')?.addEventListener('mouseenter', () => {
    const v = document.querySelector('.avatar video');
    if (!v) return;
    try {
        v.currentTime = 0;
        const p = v.play();
        if (p && p.catch) p.catch(() => {});
    } catch (e) { /* vídeo no listo, ignorar */ }
});