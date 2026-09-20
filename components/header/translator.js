let currentTranslator = null;
let currentLang = 'es'; // Idioma base en el HTML

// Código de idioma declarado en el onclick de cada botón (p. ej. "fr").
function langCodeOf(btn) {
    const match = (btn.getAttribute('onclick') || '').match(/setLanguage\('([^']+)'\)/);
    return match ? match[1] : null;
}

// ¿Puede la Translation API traducir del español a este idioma?
async function isLangAvailable(lang) {
    if (lang === 'es') return true;
    if (typeof window === 'undefined' || !('Translator' in window)) return false;
    try {
        const availability = await window.Translator.availability({
            sourceLanguage: 'es',
            targetLanguage: lang,
        });
        return !!availability && availability !== 'no';
    } catch (err) {
        return false;
    }
}

// Oculta los idiomas que la Translation API no puede ofrecer y, si no queda
// ninguno, oculta también el selector completo.
async function refreshLangSwitch() {
    const langSwitch = document.getElementById('lang-switch');
    if (!langSwitch || langSwitch.dataset.langChecked === '1') return;
    langSwitch.dataset.langChecked = '1';

    const buttons = langSwitch.querySelectorAll('.lang-btn');
    let available = 0;

    for (const btn of buttons) {
        const lang = langCodeOf(btn);
        const ok = lang ? await isLangAvailable(lang) : false;
        btn.hidden = !ok;
        if (ok && lang !== 'es') available++;
    }

    langSwitch.classList.toggle('is-supported', available > 0);
    if (available === 0) return;

    // Aplicar el idioma guardado solo si sigue estando disponible.
    const savedLang = localStorage.getItem('preferredLang');
    if (savedLang && savedLang !== 'es') {
        const savedBtn = Array.from(buttons).find(
            (btn) => langCodeOf(btn) === savedLang && !btn.hidden
        );
        if (savedBtn) setLanguage(savedLang);
    }
}

// Comprobar soporte e inicializar la visibilidad del selector al cargar el DOM
document.addEventListener('DOMContentLoaded', refreshLangSwitch);

// Cambiar idioma y actualizar la clase activa de los botones
async function setLanguage(lang) {
    if (currentLang === lang && currentTranslator) return;

    currentLang = lang;
    localStorage.setItem('preferredLang', lang);
    currentTranslator = await getTranslator(lang);

    const buttons = document.querySelectorAll('#lang-switch .lang-btn');
    buttons.forEach(btn => {
        const isMatch = btn.getAttribute('onclick')?.includes(`'${lang}'`);
        btn.classList.toggle('active', !!isMatch);
    });

    const elements = document.querySelectorAll('[data-i18n]');
    await translateElements(elements);

    const attrElements = document.querySelectorAll('[data-i18n-attr]');
    await translateAttributes(attrElements);
}

// Obtener o instanciar el traductor
async function getTranslator(targetLang) {
    if (typeof window === 'undefined' || !('Translator' in window)) return null;
    if (targetLang === 'es') return null;

    try {
        const availability = await window.Translator.availability({
            sourceLanguage: 'es',
            targetLanguage: targetLang,
        });

        if (availability === 'no') return null;

        return await window.Translator.create({
            sourceLanguage: 'es',
            targetLanguage: targetLang,
        });
    } catch (err) {
        console.error("Error inicializando Translator API:", err);
        return null;
    }
}

// Traducir elementos con data-i18n
async function translateElements(elements) {
    for (const el of elements) {
        if (!el.dataset.originalText) {
            el.dataset.originalText = el.innerText.trim();
        }

        if (currentLang === 'es') {
            el.innerText = el.dataset.originalText;
        } else if (currentTranslator) {
            try {
                const translated = await currentTranslator.translate(el.dataset.originalText);
                el.innerText = translated;
            } catch (error) {
                console.error("Error al traducir el elemento:", error);
            }
        }
    }
}

// Interceptar fragmentos inyectados por HTMX
document.body.addEventListener('htmx:afterSwap', async (evt) => {
    if (currentLang !== 'es' && currentTranslator) {
        const newElements = evt.detail.target.querySelectorAll('[data-i18n], h3, h4, p, span');
        await translateElements(newElements);

        const newAttrElements = evt.detail.target.querySelectorAll('[data-i18n-attr]');
        await translateAttributes(newAttrElements);
    }
});
// El header se inyecta por HTMX: filtrar los idiomas cuando aparezca.
document.body.addEventListener('htmx:afterSwap', refreshLangSwitch);
// Respaldo por si el swap ya ocurrió o el evento no llega.
window.addEventListener('load', () => setTimeout(refreshLangSwitch, 500));
// Cerrar dropdown al hacer clic fuera
document.addEventListener('click', (e) => {
    const langSwitch = document.getElementById('lang-switch');
    if (!langSwitch) return;
    if (!langSwitch.contains(e.target)) {
        langSwitch.classList.remove('open');
        langSwitch.querySelector('.lang-toggle')?.setAttribute('aria-expanded', 'false');
    }
});

// Toggle del dropdown
document.addEventListener('click', (e) => {
    const toggle = e.target.closest('.lang-toggle');
    if (!toggle) return;
    const langSwitch = toggle.closest('#lang-switch');
    if (!langSwitch) return;
    langSwitch.classList.toggle('open');
    toggle.setAttribute('aria-expanded',
        langSwitch.classList.contains('open').toString()
    );
});


// Traducir atributos (placeholder, title, aria-label, etc.) marcados con data-i18n-attr
async function translateAttributes(elements) {
    for (const el of elements) {
        const attrName = el.dataset.i18nAttr; // p.ej. "placeholder"
        if (!attrName) continue;

        if (!el.dataset.originalAttrText) {
            el.dataset.originalAttrText = el.getAttribute(attrName) || '';
        }

        if (currentLang === 'es') {
            el.setAttribute(attrName, el.dataset.originalAttrText);
        } else if (currentTranslator) {
            try {
                const translated = await currentTranslator.translate(el.dataset.originalAttrText);
                el.setAttribute(attrName, translated);
            } catch (error) {
                console.error("Error al traducir el atributo:", error);
            }
        }
    }
}