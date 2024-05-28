    <!--modal fade -->
    <div wire:ignore.self class="modal fade" id="dataModal" tabindex="-1" role="dialog" aria-labelledby="dataModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modader">
                    <h5 class="modal-title" id="createDataModalLabel">{{__($section)}}</h5>
                    <button type="button" class="btn-close" data-dismiss="modal" aria-label="Close">
                        <i class="fa-solid fa-xmark"></i>
                    </button>
                </div>
                <div class="modal-body">
                    <form wire:submit.prevent="save" class="form-modals">
                        @csrf
                        @foreach($fields_data as $name_field => $value)
                            <div class="group-form-contols">
                            @if(is_iterable($value))
                                <span class="input-form">{{ucfirst(__($name_field))}}</span>
                                <div class="container border p-2 rounded">
                                    @foreach($value as $id_model)
                                    <div class="row p-2">
                                        <div class="col-10 p-2 px-3">
                                            {{ $id_model }} - {{ $appends[$name_field][$id_model] }}
                                            
                                        </div>
                                        <div class="col-2 text-right">
                                            <x-buttons.cross label="none" function="deleteRowModal({{ $id_model }},'{{ $name_field }}')" />
                                        </div>
                                    </div>
                                @endforeach

                                <div class="row no-color p-2">
                                    <select id="newRow{{$name_field}}" class="form-control" onchange="addRow('newRow{{$name_field}}','{{ $name_field }}')">
                                        <option value="">{{ __('Select option')}}</option>
                                        @foreach($appends[$name_field] as $id_append => $append_value)
                                        
                                        @if (!in_array($id_append, $value))
                                        <option value="{{  $id_append }}">{{ __($append_value) }}</option>
                                        @endif
                                        @endforeach
                                    </select>
                                </div>
                            </div>

                            @else

                            <span class="input-form">{{ucfirst(__($name_field))}}</span>
                            @if($name_field == 'email')
                            <input type="email" wire:model="fields_data.{{ $name_field }}" class="form-control {{ $errors->has('fields_data.' . $name_field) ? 'error' : '' }}" placeholder="{{ucfirst(__($name_field))}}" value="{{ $value }}">
                            @elseif(is_int($value))
                            <input type="number" wire:model="fields_data.{{ $name_field }}" class="form-control {{ $errors->has('fields_data.' . $name_field) ? 'error' : '' }}" placeholder="{{ucfirst(__($name_field))}}">
                            @else
                            <input type="text" wire:model="fields_data.{{ $name_field }}" class="form-control {{ $errors->has('fields_data.' . $name_field) ? 'error' : '' }}" placeholder="{{ucfirst(__($name_field))}}">
                            @endif
                            @error('fields_data.' . $name_field)<span class="text-danger">{{ $message }}</span>@enderror

                            @endif
                        </div>
                        @endforeach

                        <div class="d-flex justify-content-end w-100">
                            <button type="button" class="btn btn-secondary close-btn" data-dismiss="modal">{{__('Close')}}</button>
                            <button type="submit" wire:click.prevent="store()" class="btn btn-primary ">{{__('Save')}}</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
    @once
<script>
    function addRow(idSelect, modelName) {
        let selectedId = document.getElementById(idSelect).value;

        @this.call('addRowModal', selectedId, modelName);
    }
</script>
@endonce