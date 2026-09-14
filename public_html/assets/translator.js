let currentTranslator = null;
let currentLang = 'es'; // Idioma base en el HTML

// Comprobar soporte e inicializar la visibilidad del selector al cargar el DOM
document.addEventListener('DOMContentLoaded', async () => {
    const langSwitch = document.getElementById('lang-switch');
    if (!langSwitch) return;

    // Verificar si la Translation API (Chrome AI / On-Device Spec) está presente
    if (typeof window !== 'undefined' && 'Translator' in window) {
        try {
            const availability = await window.Translator.availability({
                sourceLanguage: 'es',
                targetLanguage: 'en',
            });

            // Solo mostrar si el soporte y la descarga/modelo están disponibles
            if (availability && availability !== 'no') {
                langSwitch.classList.add('is-supported');
            }
        } catch (err) {
            console.warn("Translation API no disponible en este navegador:", err);
        }
    }
});

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


// Comprobar soporte e inicializar la visibilidad del selector al cargar el DOM
document.addEventListener('DOMContentLoaded', async () => {
    const langSwitch = document.getElementById('lang-switch');
    if (!langSwitch) return;

    // Verificar si la Translation API (Chrome AI / On-Device Spec) está presente
    if (typeof window !== 'undefined' && 'Translator' in window) {
        try {
            const availability = await window.Translator.availability({
                sourceLanguage: 'es',
                targetLanguage: 'en',
            });

            // Solo mostrar si el soporte y la descarga/modelo están disponibles
            if (availability && availability !== 'no') {
                langSwitch.classList.add('is-supported');
                
                // Comprobar si hay un idioma guardado previamente
                const savedLang = localStorage.getItem('preferredLang');
                // Si existe y no es el idioma base ('es'), aplicarlo automáticamente
                if (savedLang && savedLang !== 'es') {
                    setLanguage(savedLang);
                }
            }
        } catch (err) {
            console.warn("Translation API no disponible en este navegador:", err);
        }
    }
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