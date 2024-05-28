@once
<style>
    .cursor-pointer-none{
        cursor:default;
    }

    .table {
        tr {
            th {
                vertical-align: super;
            }

            td {

                vertical-align: middle;
            }
        }

        &.buttons-1 {
            .actions-buttons {
                width: 101px;

                @media (max-width: 1400px) {
                    width: 45px;
                }
            }
        }

        &.buttons-2 {
            .actions-buttons {
                width: 202px;

                @media (max-width: 1400px) {
                    width: 120px;
                }
            }
        }

        &.buttons-3 {
            .actions-buttons {
                width: 303px;

                @media (max-width: 1400px) {
                    width: 150px;
                }
            }
        }

        .actions-buttons {
            display: flex;
            justify-content: center;


            @media (max-width: 1400px) {
                span {
                    display: none;
                }
            }
        }

    }

    .form-modals {
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        gap: 1rem;

        .group-form-contols {
            width: 100%;

            span {
                width: 100%;
                font-size: large;
                font-weight: 500;
            }
        }
    }

    .modal-title {
        font-size: 2rem;
        font-weight: 500;

        padding: 1rem 2rem;
        text-transform: capitalize;
        background-image: linear-gradient(to left, #fff, #333);
        background-position: bottom;
        background-size: 100% 2px;
        /* Ancho y alto de la línea */
        background-repeat: no-repeat;
    }

    .error {
        border: 1px solid red;
    }

    .btn-close {
        position: absolute;
        width: 30px;
        height: 30px;
        top: -0.9rem;
        right: -0.9rem;
        background-color: #0069d9 !important;
        color: white !important;
        border: 1px solid #0069d9 !important;
        border-radius: 25px;
        font-size: 1rem;

    }

    .btn-close:hover {
        transform: scale(1.1);
    }

    .row:nth-child(even):not(.no-color)>* {
        background-color: #dcdcdc;
    }
</style>
@endonce

<div>
    <!-- TABLE -->
    <div class="card-header">
        <div class="d-flex align-items-baseline justify-content-between">
            <div class="float-left">
                <h4><i class="fa {{ $icon }}  text-info"></i>
                    {{ __($section) }}
                </h4>
            </div>

            <div>
                <input wire:model.live='keyWord' type="text" class="form-control" name="search" id="search" placeholder="{{ __('Search '. $section) }}">
            </div>
            <div class="btn btn-sm btn-info" data-toggle="modal" data-target="#dataModal" wire:click="showFormModal()">
                <i class="fa fa-plus"></i> {{ __('Add ' . $section) }}
            </div>
        </div>
    </div>
    <div class="card-body">
        <div class="table-responsive">
            <table class="table table-bordered table-sm {{ $widhtCellButons }}">
                <thead class="thead">
                    <tr>
                        @foreach(array_keys($fields) as $column)
                        @if(in_array($column, $appends['short']))
                        <th class="{{ $column == 'id'?'col-1':'col-3' }}">
                            <a wire:click.prevent="sortBy('{{ $column }}')" href="#" class="text-secondary">
                                {{ ucwords(__($column)) }}
                                @if ($sortField == $column && $sortDirection === 'asc')
                                <i class="fa-solid fa-caret-down ml-1"></i>
                                @else
                                <i class="fa-solid fa-caret-up ml-1"></i>
                                @endif
                            </a>
                        </th>
                        @else
                        <th class="{{ $column == 'id'?'col-1':'col-3' }}">
                            <a href="#" class="text-secondary cursor-pointer-none">
                                {{ ucwords(__($column)) }}
                            </a>
                        </th>     
                        @endif
                        @endforeach
                        @if ($widhtCellButons != "buttons-0")
                        <th class="col text-secondary text-center ">
                            {{ __('ACTIONS') }}
                        </th>
                        @endif
                    </tr>
                </thead>
                <tbody>
                    @forelse($data as $key => $row)
                    <tr>
                        @foreach(array_keys($fields) as $column)
                        <td>
                            @if( is_iterable($row->{$column}) )
                            @foreach($row->{$column} as $k => $obj )
                            {{ $appends[$column][$obj->id] }}
                            @if($k < (count($row->{$column}) - 1) )
                                ,
                                @endif
                                @endforeach
                                @else
                                {{ $row->{$column} }}
                                @endif
                        </td>
                        @endforeach
                        @if ($widhtCellButons != "buttons-0")
                        <td>
                            <div class="actions-buttons">
                                @if (in_array('edit', $buttons))
                                <x-buttons.modal :model="json_encode(get_class($row))" :id="$row->id" target="#dataModal" />
                                @endif
                                @if (in_array('delete', $buttons))
                                <x-buttons.delete function="delete({{$row->id}})" />
                                @endif

                            </div>
                        </td>
                        @endif
                    </tr>
                    @empty
                    <tr>
                        <td class="text-center" colspan="100%">{{ __('No data Found') }} </td>
                    </tr>
                    @endforelse
                </tbody>
            </table>
            <div class="float-end"></div>
        </div>
    </div>
    @include('components.forms.modal')
</div>
@once
<script>

</script>
@endonce