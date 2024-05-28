<?php

namespace App\Http\Controllers;

use Illuminate\View\View;
use Illuminate\Http\Request;

class DashboardController 
{
    public function render(Request $request, string $view = 'index'): View
    {
        switch ($view) {
            case 'profile': 
                $controller = 'App\\Livewire\\ProfileController';
                break;

            case 'users':
                $controller = 'App\\Livewire\\UsersController';
                break;
            
            default:    
                $controller = 'App\\Livewire\\ProfileController';
                break;
        }
        return view('dashboard',  [
            'controller' => $controller,
            'is_ajax' => $request->query('ajax')
        ]);
    }
}
