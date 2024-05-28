<div>
    <div class="card-header">
        <div class="d-flex align-items-baseline justify-content-between">
            <div class="float-left">
                <h4><i class="fa {{ $icon }}  text-info mx-2"></i>
                    {{ __($section) }}
                </h4>
            </div>
    
            @yield('header-content')
        </div>
    </div>
    <div class="card-body">
        @yield('body-content')
    </div>
</div>
