<div>
    @section('form')

    @endsection


    <!-- Add Modal -->
    <div wire:ignore.self class="modal fade" id="dataModal" tabindex="-1" role="dialog" aria-labelledby="dataModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modader">
                    <h5 class="modal-title" id="createDataModalLabel">{{__($model_name)}}</h5>
                    <button type="button" class="btn-close" data-dismiss="modal" aria-label="Close">
                        <i class="fa-solid fa-xmark"></i>
                    </button>
                </div>
                <div class="modal-body">
                    <form wire:submit.prevent="save" class="form-modals">
                        @csrf
                        @foreach($fill as $key => $value)
                        <label class="w-100">
                            <span class="input-form">{{ucfirst(__($key))}}</span>
                                @if($key == 'email')
                                    <input type="email" wire:model="fill.{{ $key }}" class="form-control" placeholder="{{ucfirst(__($key))}}">
                                @elseif(is_int($value))
                                    <input type="number" wire:model="fill.{{ $key }}" class="form-control" placeholder="{{ucfirst(__($key))}}">
                                @else
                                    <input type="text" wire:model="fill.{{ $key }}" class="form-control" placeholder="{{ucfirst(__($key))}}">
                                @endif
                            @endforeach
                        </label>
                        <div class="d-flex justify-content-end w-100">
                            <button type="button" class="btn btn-secondary close-btn" data-dismiss="modal">{{__('Close')}}</button>
                            <button type="submit" wire:click.prevent="save()" class="btn btn-primary">{{__('Save')}}</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <!-- Edit Modal -->
    <div wire:ignore.self class="modal fade" id="updateDataModal" tabindex="-1" role="dialog" aria-labelledby="updateModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="updateModalLabel">{{__("Update User")}}</h5>
                    <button wire:click.prevent="cancel()" type="button" class="btn-close" data-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <form wire:submit.prevent="save" class="form-modals">
                        @csrf
                        @foreach($fill as $key => $value)
                        @if($key == 'email')
                        <input type="email" wire:model="fill.{{ $key }}" class="form-control" placeholder="{{ucfirst(__($key))}}">
                        @elseif(is_int($value))
                        <input type="number" wire:model="fill.{{ $key }}" class="form-control" placeholder="{{ucfirst(__($key))}}">
                        @else
                        <input type="text" wire:model="fill.{{ $key }}" class="form-control" placeholder="{{ucfirst(__($key))}}">
                        @endif
                        @endforeach

                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary close-btn" data-dismiss="modal">{{__('Close')}}</button>
                            <button type="submit" class="btn btn-primary">{{__('Save')}}</button>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button type="button" wire:click.prevent="cancel()" class="btn btn-secondary" data-dismiss="modal">{{__("Close")}}</button>
                    <button type="button" wire:click.prevent="update()" class="btn btn-primary">{{__("Save")}}</button>
                </div>
            </div>
        </div>
    </div>
</div>