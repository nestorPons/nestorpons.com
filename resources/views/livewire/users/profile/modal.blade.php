@once
<style>
    .btn-close {
        position: absolute;
        width: 30px;
        height: 30px;
        top: -0.9rem;
        right: -0.9rem;
        background-color: #0069d9 !important;
        color: white !important;
        border: 1px solid #0069d9 !important;
        border-radius: 25px;
        font-size: 1rem;

    }

    .btn-close:hover {
        transform: scale(1.1);
    }
</style>
@endonce
<div wire:ignore.self class="modal fade" id="updateDataModal" data-backdrop="static" tabindex="-1" role="dialog"
    aria-labelledby="updateModalLabel" aria-hidden="true">
    <div class="modal-dialog  modal-dialog-centered" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="updateModalLabel">{{ __('Change password') }}</h5>

                    <button wire:click.prevent="cancel()" type="button" class="btn-close" data-dismiss="modal" aria-label="Close">
                        <i class="fa-solid fa-xmark"></i>
                    </button>
            </div>
            <form wire:submit.prevent="updatePassword">
                <div class="modal-body">
                    
                    <div class="form-group">
                        <label for="password">{{ __('Last password') }}</label>
                        <input wire:model="current_password" type="password" class="form-control" id="password"
                            name="password" autocomplete="off">
                    </div>
                    <div class="form-group">
                        <label for="new_password">{{ __('New password') }}</label>
                        <input wire:model="new_password" type="password" class="form-control" id="new_password"
                            name="new_password" autocomplete="off">

                    </div>
                    <div class="form-group">
                        <label for="password_confirmation">{{ __('Repeat new password') }}</label>
                        <input wire:model="password_confirmation" type="password" class="form-control"
                            id="password_confirmation" name="password_confirmation" autocomplete="off">
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary close-btn"
                        data-dismiss="modal">{{ __('Close') }}</button>
                    <button type="submit" class="btn btn-primary">{{ __('Change') }}</button>
                </div>
            </form>
        </div>
    </div>
</div>
