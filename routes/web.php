<?php

use App\Http\Controllers\DashboardController;
use App\Http\Controllers\LanguageController;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return view('home');
});

Route::get('lang/{lang}', [LanguageController::class, 'switchLang'])->name('lang');

Route::group(['prefix' => 'dashboard', 'middleware' => ['auth:sanctum', 'can:admin.home', 'verified']], function () {
    Route::get('/{view?}', [DashboardController::class, 'render'])->name('dashboard.index');
});