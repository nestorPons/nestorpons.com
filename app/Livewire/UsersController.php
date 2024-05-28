<?php

namespace App\Livewire;

use App\Models\User;
use Livewire\WithPagination;
use Spatie\Permission\Models\Role;
use Illuminate\Support\Facades\Auth;
use App\Livewire\Components\ListComponent;

class UsersController extends ListComponent
{
    use WithPagination;
    const UPDATE = 'update', STORE = 'store';

    protected $paginationTheme = 'bootstrap';

    public
        $roles,
        $permission,
        $selected_id,
        $name,
        $email,
        $two_factor_secret,
        $two_factor_recovery_codes,
        $current_team_id,
        $profile_photo_path;

    protected function config()
    {
        $this->view = "livewire.users.view";
        $this->icon = "fa-users";
        $this->section = "User";
        $this->buttons = ['delete', 'edit'];
        $this->model = 'app\Models\User';
        $this->fields = [
            'id' => null,
            'name' => 'string|required|max:255',
            'email' => 'email|required',
            'invited_by_name' => null,
            'roles' => 'array|exists:roles,id'
        ];

        $this->appends['short'] = ['id', 'name', 'email'];

        $user = Auth::user(); // Obtener el usuario autenticado

        // Aplicar las restricciones de rol a la consulta de roles
        $query = Role::select('id', 'name');
        if ($user->hasRole('super')) {
            // Usuario es 'super', puede ver todos los roles
            // No necesitas modificar la consulta
        } elseif ($user->hasRole('admin')) {
            // Usuario es 'admin', puede ver todos los roles excepto 'super'
            $query = $query->where('name', '!=', 'super');
        } elseif ($user->hasRole('user')) {
            // Usuario es 'user', puede ver 'user' y 'client'
            $query = $query->whereIn('name', ['user', 'client']);
        } else {
            // Usuario no cumple con las condiciones anteriores
            // Restricción que impide que se retorne cualquier rol
            $query = $query->where('id', 0); // Condición imposible para asegurar que no se devuelven datos
        }

        // Ahora aplicar mapWithKeys en el resultado de la consulta ajustada
        $this->appends['roles'] = $query->get()->mapWithKeys(function ($item) {
            return [$item['id'] => $item['name']];
        })->toArray();

    }

    public function render()
    {
        $user = Auth::user();
        if ($user->hasRole('super')) {
            // Usuario es 'super', puede ver todos los roles
            $this->roles = Role::pluck('id')->toArray();
        } elseif ($user->hasRole('admin')) {
            // Usuario es 'admin', puede ver todos los roles excepto 'super'
            $this->roles = Role::where('name', '!=', 'super')->pluck('id')->toArray();
        } elseif ($user->hasRole('user')) {
            // Usuario es 'user', puede ver 'user' y 'client'
            $this->roles = Role::whereIn('name', ['user', 'client'])->pluck('id')->toArray();
        } else {
            // Por defecto, o en caso de que no cumpla ninguna de las condiciones anteriores
            $this->roles = []; // Devuelve una colección vacía
        }
        $keyWord = '%' . $this->keyWord . '%';

        // Obtenemos el rol del usuario autenticado
        $userRole = $user->roles->pluck('name')->first(); // Asumiendo que un usuario tiene un solo rol
        $this->data = User::with('invitedBy')
            ->where(function ($query) use ($keyWord) {
                $query->orWhere('name', 'LIKE', $keyWord)
                    ->orWhere('email', 'LIKE', $keyWord)
                    ->orWhere('two_factor_secret', 'LIKE', $keyWord)
                    ->orWhere('two_factor_recovery_codes', 'LIKE', $keyWord)
                    ->orWhere('current_team_id', 'LIKE', $keyWord)
                    ->orWhere('profile_photo_path', 'LIKE', $keyWord);
            })
            ->when($userRole !== 'super', function ($query) use ($user, $userRole) {
                // Aplicamos la condición correcta para el rol 'admin'
                if ($userRole === 'admin') {
                    // Excluimos usuarios con rol 'super'
                    $query->whereDoesntHave('roles', function ($subQuery) {
                        $subQuery->where('name', 'super');
                    });
                } else {
                    // Para otros roles que no son 'super' o 'admin', limitamos a su propio usuario
                    $query->where('id', $user->id);
                }
            })
            ->orderBy($this->sortField, $this->sortDirection)
            ->paginate(15);

        return parent::render();
    }
}
