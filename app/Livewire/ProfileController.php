<?php

namespace App\Livewire;

use App\Models\User;
use Livewire\Component;
use Livewire\WithFileUploads;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Storage;
use Illuminate\Validation\ValidationException;

class ProfileController extends Component
{
    use WithFileUploads;
    public $icon, $section;
    public $image;
    public $user;
    public $new_password, $password, $password_confirmation, $current_password;
    public $name, $email;

    public function mount()
    {
        $this->user = User::find(Auth::id());
        $this->icon = "fa-user";
        $this->section  = "Profile";
        $this->name = $this->user->name;
        $this->email = $this->user->email;
        $this->image = $this->user->profile_photo_path;

    }

    public function render()
    {
        return view('livewire.users.profile.view');
    }

    public function cancel()
    {
    }

    public function rules()
    {
        return [
            'name' => 'required|string|max:255',  
            'email' => 'required|email|max:255|unique:users,email,' . Auth::id(),
        ];
    }
    public function messages()
{
    return [
        'name.required'     => __('The name field is required.'),
        'name.string'       => __('The name must be a string.'),
        'name.max'          => __('The name may not be greater than 255 characters.'),
        'email.required'    => __('The email field is required.'),
        'email.email'       => __('The email must be a valid email address.'),
        'email.max'         => __('The email may not be greater than 255 characters.'),
        'email.unique'      => __('This email has already been registered. Please use another one.'),
        'image.image'       => __('The file must be an image.'),
        'image.max'         => __('The image may not be larger than 1024 kilobytes.'),
    ];
}

    public function update()
    {
        $this->validate();
        $this->user->name = $this->name;
        $this->user->email = $this->email;

        if ($this->image) {
            $path = $this->image->store('images', 'public');
            // Obtener la URL accesible públicamente
            $url = Storage::url($path);
            // Guardar la ruta de la imagen en la base de datos

            $this->user->image_path = $url;
        }
        $this->user->save();

        $this->dispatch(
            'showMessageToast',
            title: __('Save image'),
            text: __("Updated"),
            icon: "success",
        );
        
    }

    private function resetInput()
    {
        $this->name = null;
        $this->email = null;
        $this->image = null;
    }

    public function updatePassword()
    {
        try {
            $this->validate([
                'current_password' => 'required',
                'new_password' => '',
                'password_confirmation' => 'required | required_with:new_password|same:new_password'
            ]);

            if (Hash::check($this->current_password, $this->user->password)) {
                $this->user->password = Hash::make($this->new_password);
                $this->user->save();

                $this->current_password = null;
                $this->new_password = null;
                $this->password_confirmation = null;
                $this->dispatch(
                    'showMessageToast',
                    title: __('Success'),
                    text: __("Updated Password"),
                    icon: "success",
                );
                $this->dispatch('closeModal');
            }else{
                $this->dispatch(
                    'showMessageToast',
                    title: __('Error'),
                    text: __("Updated Password"),
                    icon: "error",
                );
            }

        } catch (ValidationException $e) {
            $this->dispatch(
                'showMessageToast',
                title: __('Error'),
                text: __("Updated Password"),
                icon: "error",
            );
        }
    }
}
