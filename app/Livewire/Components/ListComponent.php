<?php

namespace App\Livewire\Components;

use Livewire\Attributes\On;

use Livewire\Component;

abstract class ListComponent extends Component
{
    public
        $widhtCellButons,
        $icon,
        $section,
        $buttons,
        $keyWord,
        $columns,
        $record_id,
        $model,
        $record,
        $sortField = 'id',
        $sortDirection = 'asc',
        $fields = [],
        $fields_data = [],
        $appends = [];

    protected
        $paginationTheme = 'bootstrap',
        $data,
        $view;

    abstract protected function config();

    public function render()
    {
        $this->config();
        $buttons = count($this->buttons);
        $this->widhtCellButons =   "buttons-" . $buttons;
        return view($this->view, ['data' => $this->data]);
    }

    public function sortBy(string $field): void
    {
        if ($this->sortField === $field) {
            $this->sortDirection = $this->sortDirection === 'asc' ? 'desc' : 'asc';
        } else {
            $this->sortField = $field;
            $this->sortDirection = 'asc';
        }
    }

    public function showFormModal(int $id = null): void
    {
        $this->record = null;
        $model = new $this->model();
        if ($id) {
            if ($this->loadModel($id)) {
                $this->record = $model::find($id);
                if ($this->record) {
                    foreach ($this->fields as $field_name => $value) {
                        if ($value) {
                            if (strpos($value, 'array') !== false) {
                                $this->fields_data[$field_name] = [];
                                foreach ($this->record->{$field_name} as $item) {
                                    $this->fields_data[$field_name][] = $item->id;
                                }
                            } else {
                                $this->fields_data[$field_name] = $this->record->{$field_name};
                            }
                        }
                    }
                }
            }
        } else {
            foreach ($this->fields as $key => $value) {
                if ($value) {
                    $type = explode('|', $value)[0];
                    $inicialized = $type == 'array' ? [] : "";
                    $this->fields_data[$key] = $inicialized;
                }
            }
        }
    }

    private function rulesValidate(): array
    {
        $fields = [];

        foreach ($this->fields_data as $key => $value) {
            if ($this->fields[$key]) {
                $fields['fields_data.' . $key] = $this->fields[$key];
            }
        }
        return $fields;
    }

    public function store(): void
    {
        $title = "";
        $rules = $this->rulesValidate();
        $validatedData = $this->validate($rules);
        if ($this->record) {
            $title = "update";
            $this->record->update($validatedData['fields_data']);
        } else {
            $title = "create";
            $this->record = $this->model::create($validatedData['fields_data']);
        }
        foreach ($this->fields as $field_name => $rules) {
            if (strpos($rules, 'array') !== false) {
                $this->record->{$field_name}()->sync($this->fields_data[$field_name]);
            }
        }
        $this->dispatch(
            'showMessageToast',
            icon: "success",
            title: __($title),
            showConfirmButton: false,
            text: __(ucfirst($this->section) . " " . $title . " successfully"),
        );

        $this->resetErrorBag();
        $this->dispatch('closeModal');
    }

    public function delete($id): void
    {
        $this->record_id = $id;
        $this->dispatch(
            'showDeleteDialog',
            title: __('Delete'),
            text: __("Are you sure?"),
            confirmButtonText: __('Accept'),
        );
    }

    public function deleteFromModal(): void
    {
        $this->resetErrorBag();
        $this->record->delete();
    }

    #[On('runDelete')]
    public function destroy(): void
    {
        if ($this->loadModel($this->record_id)) {
            $this->record->delete();
        }
    }

    private function loadModel(int $id): bool
    {
        $model = new $this->model();
        $this->record = $model::find($id);
        if (!$this->record) {
            session()->flash('error', 'Record no found');
            return false;
        }
        return true;
    }

    public function addRowModal(string $id, string $model_attr): void
    {
        $this->fields_data[$model_attr][] = $id;
    }

    public function deleteRowModal($id, $accessor)
    {
        $this->record;
        if ($this->record) {
            $this->record->{$accessor}()->detach($id);
            $this->refreshModal();
        }
    }

    private function refreshModal()
    {
        $this->showFormModal($this->record->id);
    }
}
