<?php namespace App\Livewire\Components;

use Livewire\Component;

class Modals extends Component
{

    public $data, $model, $model_name;
    public $fill = [];

    protected $listeners = [
        'showFormModal' => 'handleLoadData'
    ];

    protected $paginationTheme = 'bootstrap';
    
    public function rules()
    {
        return [ 
            'fill.email' => ['nullable', 'email'], // Aplica la regla 'email' si el campo 'fill' es un correo electrónico válido
            'fill.text' => ['nullable', 'text'], // Aplica la regla 'email' si el campo 'fill' es un correo electrónico válido
        ];
    }

    public function handleLoadData($model_name, $id)
    {
        // Al recibir nuevos datos, también actualizamos $fill
        $model_instance = new $model_name();
        $this->model = $model_instance::find($id);
    }

    public function mount(object $model)
    {
        $this->model = $model;
        
    }
    
    public function render()
    {
        $this->model_name = class_basename($this->model);
        foreach ($this->model->getFillable() as $item) {
            $this->fill[$item] = $this->model->{$item}??null;
        }
        return view('livewire.components.modals');
    }

    public function save()
    {
        // Si necesitas acceder a los datos reactivos del formulario, puedes hacerlo a través de $this->fill
        dd($this->fill);
    }
}
