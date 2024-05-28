import axios from 'axios';
window.axios = axios;

window.axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';

document.addEventListener('DOMContentLoaded', function () {
    const modalToggleButtons = document.querySelectorAll('[data-toggle="modal"]');
    const modalDismissButtons = document.querySelectorAll('[data-dismiss="modal"]');
    
    modalToggleButtons.forEach(function (button) {
        button.addEventListener('click', function () {
            const targetSelector = this.getAttribute('data-target');
            const modal = new bootstrap.Modal(document.querySelector(targetSelector));
            modal.show();
        });
    });

    modalDismissButtons.forEach(function (button) {
        button.addEventListener('click', function () {
            const modals = document.querySelectorAll('.modal');
            modals.forEach(function (modal) {
                const modalInstance = bootstrap.Modal.getInstance(modal);
                if (modalInstance) {
                    modalInstance.hide();
                }
            });
        });
    });
});

