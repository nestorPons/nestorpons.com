<?php

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Seeder;
use Spatie\Permission\Models\Role;
use Illuminate\Support\Facades\Hash;

class UserSeeder extends Seeder
{

    public function run(): void
    {
        $superUser = User::firstOrCreate([
            'name' => 'Super',
            'email' => 'super@super',
            'password' => Hash::make('super'),
            'current_team_id' => 0
        ]);
        $superUser->assignRole(Role::firstOrCreate(['name' => 'super']));
        $superUser->givePermissionTo(['admin.home', 'admin.settings', 'admin.full']);

        $adminUser = User::firstOrCreate([
            'name' => 'Admin',
            'email' => 'admin@admin',
            'password' => Hash::make('admin'),
            'current_team_id' => 0
        ]);
        $adminUser->assignRole(Role::firstOrCreate(['name' => 'admin']));
        $adminUser->givePermissionTo( ['admin.home', 'admin.settings']);

        $staffUser = User::firstOrCreate([
            'name' => 'Staff',
            'email' => 'staff@staff',
            'password' => Hash::make('staff'),
            'current_team_id' => 0
        ]);
        $staffUser->assignRole(Role::firstOrCreate(['name' => 'staff']));
        $staffUser->givePermissionTo(['admin.home']);
        
        $clientUser = User::firstOrCreate([
            'name' => 'guest',
            'email' => 'guest@guest',
            'password' => Hash::make('guest'),
            'invited_by' => 3, 
            'current_team_id' => 0
        ]);
        $clientUser->assignRole(Role::firstOrCreate(['name' => 'guest']));
    }
}
