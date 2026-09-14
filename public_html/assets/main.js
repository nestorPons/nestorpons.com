document.querySelectorAll('a[href^="#"]').forEach(link => {

    link.addEventListener('click', event => {

        event.preventDefault();

        document.querySelector(
            link.getAttribute('href')
        ).scrollIntoView({
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


// MENU
const hamburger = document.getElementById('hamburger');
const menu = document.querySelector('.menu');

hamburger.addEventListener('click', () => {
    hamburger.classList.toggle('open');
    menu.classList.toggle('open');
});

menu.querySelectorAll('a').forEach(link => {
    link.addEventListener('click', () => {
        hamburger.classList.remove('open');
        menu.classList.remove('open');
    });
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