<button 
    type="button"
    class="btn btn-outline-{{$color}} {{ $class ?? '' }}" 
    @if (isset($closeModal) && $closeModal)
        data-dismiss="modal"
    @endif
    @if (isset($toggleModal) && $toggleModal)
        data-toggle="modal"
        data-target={{ $toggleModal }}
    @endif
    @if (isset($message))
        wire:confirm= "{{ $message }}"
    @endif
    @if (isset($collapse) && $collapse)
        data-toggle="collapse" 
        aria-expanded="true" 
        data-target="#collapse_{{$collapse}}" 
        aria-controls="collapse_{{$collapse}}"
    @endif

    @if (isset($onclick) && $onclick)
        onclick={{ $onclick }}
    @endif  
    
    @if (isset($function))
        wire:click={{ $function }}
    @endif    
    >
    <i class="fa fa-{{$icon}}"></i>
    @if (isset($label) && $label != 'none' )
        <span>
            @if ($label)
                {{__($label)}}
            @else
                {{ __('Default') }}
            @endif  
        </span>
    @endif
</button>