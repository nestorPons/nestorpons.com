<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use Spatie\Permission\Models\Role;
use Spatie\Permission\Models\Permission;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;

class RoleSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        // Eliminar todos los roles y permisos existentes
        Role::query()->delete();
        Permission::query()->delete();
        DB::table('model_has_roles')->delete();
        DB::table('model_has_permissions')->delete();
        // Reiniciar el contador autoincremental para MySQL
        DB::statement('ALTER TABLE roles AUTO_INCREMENT = 1;');
        DB::statement('ALTER TABLE permissions AUTO_INCREMENT = 1;');

        // Crear permisos
        Permission::firstOrCreate(['name' => 'admin.full',      'guard_name' => 'web']);
        Permission::firstOrCreate(['name' => 'admin.home',      'guard_name' => 'web']);
        Permission::firstOrCreate(['name' => 'admin.settings',  'guard_name' => 'web']);
        Permission::firstOrCreate(['name' => 'app.super',       'guard_name' => 'web']);
        Permission::firstOrCreate(['name' => 'app.admin',       'guard_name' => 'web']);
        Permission::firstOrCreate(['name' => 'app.staff',       'guard_name' => 'web']);
        Permission::firstOrCreate(['name' => 'app.guest',       'guard_name' => 'web']);
    
        // Crear roles
        $role1 = Role::create(["name" => "super"]);
        $role1->syncPermissions(['admin.home', 'admin.settings', 'admin.full', 'app.super', 'app.admin', 'app.staff', 'app.guest']);

        $role2 = Role::create(["name" => "admin"]);
        $role2->syncPermissions(['admin.home', 'admin.settings', 'app.admin', 'app.staff', 'app.guest']);

        $role3 = Role::create(["name" => "staff"]);
        $role3->syncPermissions(['admin.home', 'app.staff', 'app.guest']);

        $role4 = Role::create(["name" => "guest"]);
        $role4->syncPermissions(['app.guest']);    
    }
}
