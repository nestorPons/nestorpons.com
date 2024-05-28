@section('form')
<form>
    <div class="form-group">
        <label for="name">{{__('Name')}}</label>
        <input wire:model.live="name" type="text" class="form-control" id="name" placeholder="{{__('Name')}}">@error('name') <span class="error text-danger">{{ $message }}</span> @enderror
    </div>
    <div class="form-group">
        <label for="email">{{__('Email')}}</label>
        <input wire:model.live="email" type="text" class="form-control" id="email" placeholder="{{__('Email')}}">@error('email') <span class="error text-danger">{{ $message }}</span> @enderror
    </div>
    <div class="form-group">
        <label for="current_team_id">{{__('Team')}}</label>
        <input wire:model.live="current_team_id" type="text" class="form-control" id="current_team_id" placeholder="{{__('Current Team Id')}}">@error('current_team_id') <span class="error text-danger">{{ $message }}</span> @enderror
    </div>
    <div class="form-group">
        <label for="profile_photo_path">"{{__('Profile Photo Path')}}</label>
        <input wire:model.live="profile_photo_path" type="text" class="form-control" id="profile_photo_path" placeholder="{{__('Profile Photo Path')}}">@error('profile_photo_path') <span class="error text-danger">{{ $message }}</span> @enderror
    </div>
    <div class="form-group">
        <label for="permission">{{__('Permission')}}</label>
        <select wire:model.live="permission" class="form-control" id="permission" placeholder="{{__('Permission')}}">
            @foreach($roles as $role)
            <option value="{{ $role->id }}">{{ $role->name }}</option>
            @endforeach
        </select>
        @error('permission') <span class="error text-danger">{{ $message }}</span> @enderror
    </div>
</form>

@endsection


<!-- Add Modal -->
<div wire:ignore.self class="modal fade" id="createDataModal" data-backdrop="static" tabindex="-1" role="dialog" aria-labelledby="createDataModalLabel" aria-hidden="true">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="createDataModalLabel">{{__("Create New User")}}</h5>
                <button wire:click.prevent="cancel()" type="button" class="btn-close" data-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                @yield('form')
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary close-btn" data-dismiss="modal">{{__('Close')}}</button>
                <button type="button" wire:click.prevent="store()" class="btn btn-primary">{{__('Save')}}</button>
            </div>
        </div>
    </div>
</div>

<!-- Edit Modal -->
<div wire:ignore.self class="modal fade" id="updateDataModal" data-backdrop="static" tabindex="-1" role="dialog" aria-labelledby="updateModalLabel" aria-hidden="true">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="updateModalLabel">{{__("Update User")}}</h5>
                <button wire:click.prevent="cancel()" type="button" class="btn-close" data-dismiss="modal" aria-label="Close"></button>
            </div>
            <div class="modal-body">
                @yield('form')
            </div>
            <div class="modal-footer">
                <button type="button" wire:click.prevent="cancel()" class="btn btn-secondary" data-dismiss="modal">{{__("Close")}}</button>
                <button type="button" wire:click.prevent="update()" class="btn btn-primary">{{__("Save")}}</button>
            </div>
        </div>
    </div>
</div>
