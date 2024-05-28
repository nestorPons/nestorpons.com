import "./bootstrap";
import '@fortawesome/fontawesome-free/js/all';
import Swal from 'sweetalert2'

  
const app = {
    isLoading: false, // Indicador de carga
    load: function() {
        const container = document.getElementsByClassName("sidebar")[0];
        const navLinks = container.getElementsByClassName("nav-bar-link");
        const navLinksAsync = document.getElementsByClassName("nav-link-async");

        // Convertir HTMLCollection o NodeList a arrays
        const navLinksArray = Array.from(navLinks);
        const navLinksAsyncArray = Array.from(navLinksAsync);
        
        // Combinar los dos arrays
        const combinedArray = navLinksArray.concat(navLinksAsyncArray);
        for (const element of combinedArray) {
            this.reactiveNavLinks(element)
        } 
    },
    reactiveNavLinks: function(link){
        link.addEventListener("click", event => {
            console.log("reactiveNavLinks")
            event.preventDefault();
            if (this.isLoading) {
                return; // Evitar peticiones adicionales si ya hay una en curso
            }

            const url = link.getAttribute("href");
            if (url) {
                const contentSection = document.querySelector(".content");
                this.isLoading = true; // Marcar que una petición está en curso
                // Verificar si la URL ya tiene parámetros
                const separator = url.includes('?') ? '&' : '?';

                // Añadir el parámetro 'ajax=1' a la URL
                const ajaxUrl = url + separator + 'ajax=1' + '&language=es';
                axios
                .get(ajaxUrl)
                .then((response) => {
                    contentSection.innerHTML = response.data;
                })
                .catch((error) => {
                    contentSection.innerHTML = error;
                    console.error("Error al cargar el contenido:", error);
                })
                .finally(() => {
                    this.isLoading = false; // Restablecer el indicador después de la petición
                });
            }
            history.pushState({ url }, "", "");
        });
    }
}

document.addEventListener("DOMContentLoaded", function() {
    app.load();
    // Colorea las filas de las tablas seleccionadas
    // Manejar el clic del botón que agrega la clase
    $('.content').on('click','.row .show-selected', function () {
        // Seleccionar todos los elementos padre con la clase .row y agregar la clase row-selected
        $('.row-selected').removeClass('row-selected')
        $(this).closest('.row').addClass('row-selected');
    });
        
    Livewire.on('showMessage', ({ title, text, icon }) => {
        Swal.fire({
            title: title,
            text: text,
            icon: icon,
            background: "#fff"}
            );
    });

    Livewire.on('showMessageToast', ({ text, icon }) => {
        Swal.fire({
            toast: true,
            position: "top-end",
            icon: icon,
            title: text,
            showConfirmButton: false,
            timer: 1500
        });
    });

    Livewire.on('showDeleteDialog', ({ title, text, confirmButtonText }) => {
        Swal.fire({
            title: title,
            text: text,
            icon: "warning",
            background: "#fff",
            showCancelButton: true,
            cancelButtonColor: "#d33",
            confirmButtonColor: "#3085d6",
            confirmButtonText: confirmButtonText,
            reverseButtons: true
          }).then((result) => {
                if (result.value) {
                    Livewire.dispatch('runDelete')
                    /*         
                    Swal.fire({
                        title: respondTitle,
                        text: respondText,
                        icon: "success"
                    }); 
                    */
                }
          });
    });

    Livewire.on('closeModal', (element) => {
        $("#dataModal").modal('hide')
        $(".modal").modal('hide')
        $(".modal-backdrop").remove()
    });
});

