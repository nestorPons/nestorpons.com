@extends('layouts.main')
@section('header-content')
@endsection
@section('body-content')


    <form  wire:submit="update">
        <div class="body">
            @csrf
            <div class="form-group">
                <label for="name">{{ __('Name') }}</label>
                <input type="text" id="name" class="form-control" wire:model="name">
                @if ($errors->has('name'))
                    <span class="text-danger">{{ $errors->first('name') }}</span>
                @endif
            </div>
            
            <div class="form-group">
                <label for="email">{{ __('Email') }}</label>
                <input type="text" id="email" class="form-control" wire:model="email">
                @if ($errors->has('email'))
                    <span class="text-danger">{{ $errors->first('email') }}</span>
                @endif
            </div>

            <div class="form-group">
                <label for="image">{{ __('Select image') }}</label>
                <input type="file" id="image" wire:model="image">
                @if ($errors->has('image'))
                    <span class="text-danger">{{ $errors->first('image') }}</span>
                @endif
            </div>
        </div>
        <div class="footer">
            <button type="button" class="btn btn-warning" data-target="#updateDataModal" data-toggle="modal">
                {{ __('Change password') }}
            </button>
            <button type="submit" class="btn btn-primary">{{ __('Save') }}</button>
        </div>
    </form>


    @include('livewire.users.profile.modal')       
@endsection
