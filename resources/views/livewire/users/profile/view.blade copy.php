
@extends('layouts.main')
<div class="" id="" data-backdrop="static" tabindex="-1" role="dialog" aria-labelledby="" aria-hidden="true">
    <div class="" role="document">
        <div class="">
            <div class="header">
                <h5 class="modal-title" id="updateModalLabel">{{ __('Profile') }}</h5>
            </div>

            <form action="#" method="POST" enctype="multipart/form-data" id="frmProfile">
                <div class="body">
                    @csrf
                    <input type="hidden" name="selected_id" value="{{ $selected_id }}">
                    <div class="form-group">
                        <label for="name">{{ __('Name') }}</label>
                        <input type="text" class="form-control" id="name" name="name" value="{{ $name }}" placeholder="">
                        @if ($errors->has('name'))
                            <span class="text-danger">{{ $errors->first('name') }}</span>
                        @endif
                    </div>
                    
                    <div class="form-group">
                        <label for="email">{{ __('Email') }}</label>
                        <input type="text" class="form-control" id="email" name="email" value="{{ $email }}">
                        @if ($errors->has('email'))
                            <span class="text-danger">{{ $errors->first('email') }}</span>
                        @endif
                    </div>
                    <div class="form-group">
                        <label for="current_team_id">{{ __('Current Team Id') }}</label>
                        <input type="text" class="form-control" id="current_team_id" name="current_team_id" value="{{ $current_team_id }}" placeholder="">
                    </div>
                    <div class="form-group">
                        <label for="image">{{ __('Select image') }}</label>
                        <input type="file" id="image" name="image">
                    </div>
                </div>
                <div class="footer">
                    <button type="button" class="btn btn-warning" data-target="#updateDataModal" data-toggle="modal">
                        {{ __('Change password') }}
                    </button>
                    <button type="submit" class="btn btn-primary">{{ __('Save') }}</button>
                </div>
            </form>
        </div>
    </div>
    @include('livewire.users.profile.modal')       
</div>

