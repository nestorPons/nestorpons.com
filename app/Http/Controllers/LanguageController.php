<?php

namespace App\Http\Controllers;

use Illuminate\Support\Facades\Session;

class LanguageController
{
    public function switchLang($lang){
        Session::put('locale', $lang);
        return redirect()->back()->with(['language'=>$lang]);
    }   
}
