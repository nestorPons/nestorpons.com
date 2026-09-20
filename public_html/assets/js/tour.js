// Capa de anotaciones dibujada a mano: muestra de una vez todo lo que
// destaca la web, con flechas y tipografia manuscrita sobre un fondo oscuro.
(function () {
    var SVGNS = 'http://www.w3.org/2000/svg';

    var FEATURES = [
        {
            sel: '#lang-switch',
            side: 'bottom',
            title: 'Cambio de idioma',
            text: 'Uso de nuevas tecnologías, como las apis nativas de google que permiten traducir la web mediante IA.'
        },
        {
            sel: '.avatar',
            side: 'left',
            title: 'Avatar dinámico',
            text: 'No es una foto: es un vídeo que te saluda. Pasa el ratón por encima y lo repite.'
        },
        {
            sel: '.hero h2',
            side: 'left',
            title: 'La IA crea conmigo',
            text: 'Uso la IA como copiloto para acelerar el diseño y el código. La idea, la arquitectura y el criterio son míos: aquí no se construye una IA, se usa.'
        },
        {
            sel: '.hero p',
            side: 'right',
            title: 'De escribir código a orquestar',
            text: 'El nuevo paradigma del programador: ya no se trata de teclear cada línea, sino de dirigir. Ser orquestador es elegir tecnologías, conectar sistemas, coordinar IAs y personas para que todo encaje.'
        },
        {
            sel: '#ai-chat-log',
            side: 'right',
            title: 'Asistente a medida',
            text: 'Un chatbot propio entrenado con mi perfil. Pregúntale por proyectos, experiencia o servicios y responde en streaming.'
        },
        {
            sel: '[hx-get="/api/components/projects"]',
            side: 'left',
            title: 'HTMX, sin frameworks',
            text: 'El servidor envía fragmentos HTML y HTMX actualiza la página sin recargarla. Ligero y moderno.'
        },
        {
            sel: '.timeline-container',
            side: 'top',
            title: 'Línea de tiempo',
            text: 'Una vista cronológica de mi trayectoria profesional y los proyectos más importantes. Cada proyecto tiene su propia historia.'
        },
        {
            sel: '#skills .skills-grid',
            side: 'top-left',
            title: 'Habilidades y tecnologías',
            text: 'Tener un mapa mental amplio del ecosistema tecnológico me permite seleccionar e integrar la mejor solución en cada proyecto.'
        },
        {
            sel: 'footer .footer-tech',
            side: 'right',
            title: 'Backend en Go',
            text: 'Servidor escrito en Go: arranque instantáneo, poco consumo y respuestas ultrarrápidas.'
        }
    ];

    var state = {
        open: false,
        overlay: null,
        layer: null,
        svg: null,
        closeBtn: null,
        hint: null,
        items: [],
        onKey: null,
        onAfterSwap: null,
        waitTimer: null
    };
    var resizeTimer = null;

    function svgEl(tag, attrs) {
        var el = document.createElementNS(SVGNS, tag);
        for (var k in attrs) el.setAttribute(k, attrs[k]);
        return el;
    }

    // Rectángulo con lados ligeramente curvados y esquinas que se pasan,
    // para simular un trazo hecho a mano.
    function handRect(x, y, w, h) {
        var j = 2.5;
        return 'M' + (x + 3) + ' ' + (y + j) +
            ' Q' + (x + w / 2) + ' ' + (y - j) + ' ' + (x + w - 2) + ' ' + (y + 2) +
            ' Q' + (x + w + j) + ' ' + (y + h / 2) + ' ' + (x + w - 1) + ' ' + (y + h - 3) +
            ' Q' + (x + w / 2) + ' ' + (y + h + j) + ' ' + (x + 3) + ' ' + (y + h - 1) +
            ' Q' + (x - j) + ' ' + (y + h / 2) + ' ' + (x + 1) + ' ' + (y + 2);
    }

    // Punto del borde de un rectángulo en la direccion desde (sx, sy) hacia su centro.
    function boundaryPoint(sx, sy, rect) {
        var cx = rect.left + rect.width / 2;
        var cy = rect.top + rect.height / 2;
        var dx = cx - sx;
        var dy = cy - sy;
        if (!dx && !dy) return { x: cx, y: cy };
        var tx = dx ? (rect.width / 2) / Math.abs(dx) : Infinity;
        var ty = dy ? (rect.height / 2) / Math.abs(dy) : Infinity;
        var t = Math.min(tx, ty);
        return { x: cx - dx * t, y: cy - dy * t };
    }

    function clamp(v, min, max) {
        return Math.max(min, Math.min(max, v));
    }

    function drawPath(path, d, animate) {
        path.setAttribute('d', d);
        if (!animate) {
            path.style.strokeDasharray = 'none';
            path.style.strokeDashoffset = '0';
            return;
        }
        var len = path.getTotalLength ? path.getTotalLength() : 800;
        path.style.transition = 'none';
        path.style.strokeDasharray = len;
        path.style.strokeDashoffset = len;
        requestAnimationFrame(function () {
            requestAnimationFrame(function () {
                path.style.transition = '';
                path.style.strokeDashoffset = '0';
            });
        });
    }

    function createChrome() {
        var docW = document.documentElement.clientWidth;
        var docH = Math.max(document.documentElement.scrollHeight, window.innerHeight);

        var overlay = document.createElement('div');
        overlay.className = 'tour-overlay';

        var layer = document.createElement('div');
        layer.className = 'tour-layer';

        var svg = svgEl('svg', { width: docW, height: docH, viewBox: '0 0 ' + docW + ' ' + docH });
        layer.appendChild(svg);

        var closeBtn = document.createElement('button');
        closeBtn.type = 'button';
        closeBtn.className = 'tour-close';
        closeBtn.setAttribute('aria-label', 'Cerrar anotaciones');
        closeBtn.innerHTML = '\u2715 <span>Cerrar</span>';

        var hint = document.createElement('div');
        hint.className = 'tour-hint';
        hint.textContent = 'Desliza para ver cada detalle · Esc o toca fuera para cerrar';

        document.body.appendChild(overlay);
        document.body.appendChild(layer);
        document.body.appendChild(closeBtn);
        document.body.appendChild(hint);

        state.overlay = overlay;
        state.layer = layer;
        state.svg = svg;
        state.closeBtn = closeBtn;
        state.hint = hint;
        state.items = [];
    }

    // Crea la anotación de cada FEATURE cuyo elemento ya exista en el DOM.
    // Es idempotente: se puede invocar de nuevo tras cada carga HTMX para
    // añadir las que aparezcan más tarde (lazy load).
    function buildItems() {
        FEATURES.forEach(function (feature) {
            var already = state.items.some(function (it) { return it.feature.sel === feature.sel; });
            if (already) return;

            var target = document.querySelector(feature.sel);
            if (!target) return;

            var idx = FEATURES.indexOf(feature);

            var card = document.createElement('div');
            card.className = 'tour-note';
            card.style.animationDelay = (0.05 + idx * 0.12) + 's';

            var num = document.createElement('span');
            num.className = 'tour-note-num';
            num.textContent = String(idx + 1);
            var title = document.createElement('h3');
            title.textContent = feature.title;
            var desc = document.createElement('p');
            desc.textContent = feature.text;
            card.appendChild(num);
            card.appendChild(title);
            card.appendChild(desc);

            var focus = svgEl('path', { class: 'tour-focus' });
            var sketch = svgEl('path', { class: 'tour-focus-sketch' });
            var arrow = svgEl('path', { class: 'tour-arrow' });
            var head = svgEl('path', { class: 'tour-arrow-head' });
            state.svg.appendChild(focus);
            state.svg.appendChild(sketch);
            state.svg.appendChild(arrow);
            state.svg.appendChild(head);
            state.layer.appendChild(card);

            state.items.push({
                feature: feature,
                target: target,
                card: card,
                focus: focus,
                sketch: sketch,
                arrow: arrow,
                head: head
            });
        });
    }

    // ¿Están ya en el DOM todos los elementos objetivo?
    function allTargetsReady() {
        return FEATURES.every(function (f) { return !!document.querySelector(f.sel); });
    }

    // Espera a que terminen las cargas perezosas de HTMX antes de medir.
    // Si algún objetivo no aparece a tiempo, se construye con lo que haya.
    function waitForTargets(done) {
        var waited = 0;
        var step = 100;
        var max = 6000;

        function finish() {
            if (state.waitTimer) {
                clearInterval(state.waitTimer);
                state.waitTimer = null;
            }
            if (state.open) done();
        }

        if (allTargetsReady()) { finish(); return; }

        state.waitTimer = setInterval(function () {
            if (!state.open) { finish(); return; }
            waited += step;
            if (allTargetsReady() || waited >= max) finish();
        }, step);
    }

    function layout(animate) {
        var vw = document.documentElement.clientWidth;
        var docH = Math.max(document.documentElement.scrollHeight, window.innerHeight);
        var pad = 16;
        var gap = 30;
        var maxCard = 320;
        var scrollX = window.scrollX || window.pageXOffset || 0;
        var scrollY = window.scrollY || window.pageYOffset || 0;

        state.svg.setAttribute('width', vw);
        state.svg.setAttribute('height', docH);
        state.svg.setAttribute('viewBox', '0 0 ' + vw + ' ' + docH);

        state.items.forEach(function (item) {
            var r = item.target.getBoundingClientRect();
            var rect = {
                left: r.left + scrollX,
                top: r.top + scrollY,
                width: r.width,
                height: r.height
            };

            var cardW = Math.min(maxCard, vw - pad * 2);
            item.card.style.width = cardW + 'px';
            item.card.style.left = '-9999px';
            item.card.style.top = '0px';
            var cardH = item.card.offsetHeight || 130;

            var left;
            var top;
            if (vw < 720) {
                left = pad;
                var below = rect.top + rect.height + gap;
                top = (below + cardH < docH) ? below : Math.max(pad, rect.top - gap - cardH);
            } else {
                if (item.feature.side === 'right') {
                    left = rect.left + rect.width + gap;
                    if (left + cardW > vw - pad) left = rect.left - gap - cardW;
                } else if (item.feature.side === 'top') {
                    // Centrada encima del objetivo; si no cabe, se pone debajo.
                    left = rect.left + rect.width / 2 - cardW / 2;
                    left = clamp(left, pad, Math.max(pad, vw - cardW - pad));
                    top = rect.top - gap - cardH;
                    if (top < pad) top = rect.top + rect.height + gap;
                } else if (item.feature.side === 'top-left') {
                    // Arriba y a la izquierda del objetivo; si no cabe, se prueba
                    // al lado opuesto y, si tampoco, debajo.
                    left = rect.left - gap - cardW;
                    if (left < pad) left = rect.left + rect.width + gap;
                    top = rect.top - gap - cardH;
                    if (top < pad) top = rect.top + rect.height + gap;
                } else if (item.feature.side === 'bottom') {
                    // Centrada debajo del objetivo; si no cabe, se pone encima.
                    left = rect.left + rect.width / 2 - cardW / 2;
                    left = clamp(left, pad, Math.max(pad, vw - cardW - pad));
                    top = rect.top + rect.height + gap;
                    if (top + cardH > docH - pad) top = rect.top - gap - cardH;
                } else {
                    left = rect.left - gap - cardW;
                    if (left < pad) left = rect.left + rect.width + gap;
                }
                left = clamp(left, pad, Math.max(pad, vw - cardW - pad));
                var positioned = item.feature.side === 'top' ||
                    item.feature.side === 'top-left' ||
                    item.feature.side === 'bottom';
                if (!positioned) {
                    top = rect.top + rect.height / 2 - cardH / 2;
                }
            }
            top = clamp(top, pad, Math.max(pad, docH - cardH - pad));

            item.card.style.left = left + 'px';
            item.card.style.top = top + 'px';

            var cardRect = { left: left, top: top, width: cardW, height: cardH };
            var center = { x: rect.left + rect.width / 2, y: rect.top + rect.height / 2 };

            drawPath(item.focus, handRect(rect.left, rect.top, rect.width, rect.height), animate);
            drawPath(item.sketch, handRect(rect.left + 2, rect.top + 2, rect.width, rect.height), animate);

            var start = boundaryPoint(center.x, center.y, cardRect);
            var end = boundaryPoint(start.x, start.y, rect);
            var dx = end.x - start.x;
            var dy = end.y - start.y;
            var dist = Math.sqrt(dx * dx + dy * dy) || 1;
            var bend = Math.min(46, dist * 0.16);
            var mx = (start.x + end.x) / 2 - (dy / dist) * bend;
            var my = (start.y + end.y) / 2 + (dx / dist) * bend;
            drawPath(item.arrow, 'M' + start.x + ' ' + start.y + ' Q' + mx + ' ' + my + ' ' + end.x + ' ' + end.y, animate);

            var ang = Math.atan2(end.y - my, end.x - mx);
            var len = 14;
            var spread = 0.42;
            var b1x = end.x + Math.cos(ang + Math.PI - spread) * len;
            var b1y = end.y + Math.sin(ang + Math.PI - spread) * len;
            var b2x = end.x + Math.cos(ang + Math.PI + spread) * len;
            var b2y = end.y + Math.sin(ang + Math.PI + spread) * len;
            item.head.setAttribute('d', 'M' + b1x + ' ' + b1y + ' L' + end.x + ' ' + end.y + ' L' + b2x + ' ' + b2y);
        });
    }

    function forceLazyContent() {
        if (!window.htmx) return;
        document.querySelectorAll('[hx-trigger*="intersect"]').forEach(function (el) {
            try { window.htmx.trigger(el, 'intersect'); } catch (e) { /* ignorar */ }
        });
    }

    function open() {
        if (state.open) return;
        state.open = true;

        createChrome();

        // Pide ya los fragmentos HTMX perezosos y espera a que estén en el DOM
        // para poder apuntarles con flechas y recuadros.
        forceLazyContent();

        state.onKey = function (e) {
            if (e.key === 'Escape') close();
        };
        state.onAfterSwap = function () {
            if (!state.open) return;
            buildItems();
            layout(false);
        };

        document.addEventListener('keydown', state.onKey);
        document.body.addEventListener('htmx:afterSwap', state.onAfterSwap);
        window.addEventListener('resize', onResize);
        state.overlay.addEventListener('click', close);
        state.closeBtn.addEventListener('click', close);

        waitForTargets(function () {
            buildItems();
            layout(true);
        });
    }

    function close() {
        if (!state.open) return;
        state.open = false;

        if (state.waitTimer) {
            clearInterval(state.waitTimer);
            state.waitTimer = null;
        }

        document.removeEventListener('keydown', state.onKey);
        document.body.removeEventListener('htmx:afterSwap', state.onAfterSwap);
        window.removeEventListener('resize', onResize);

        [state.overlay, state.layer, state.closeBtn, state.hint].forEach(function (el) {
            if (el && el.parentNode) el.parentNode.removeChild(el);
        });
        state.items = [];
    }

    function onResize() {
        clearTimeout(resizeTimer);
        resizeTimer = setTimeout(function () {
            if (!state.open) return;
            close();
            open();
        }, 200);
    }

    function init() {
        var btn = document.getElementById('btn-tour');
        if (!btn) return;

        // El tour ya está operativo: retira el estado de carga del botón.
        btn.classList.remove('is-loading');
        btn.classList.add('is-ready');
        btn.setAttribute('aria-busy', 'false');
        btn.removeAttribute('aria-disabled');

        btn.addEventListener('click', function (e) {
            e.preventDefault();
            var cta = document.getElementById('tour-cta');
            if (cta) cta.classList.add('is-hidden');
            if (state.open) close(); else open();
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
