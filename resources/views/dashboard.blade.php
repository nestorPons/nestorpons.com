@if($is_ajax)
    @livewire($controller)
@else
<x-app-layout>
    @section('content')
        <div class="my-2">
            @livewire($controller)
        </div>
    @endsection
</x-app-layout>
@endif