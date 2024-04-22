<template>
    <div class="grid grid-cols-4 gap-4 mt-6">
        <template v-for="({ image, name, url }, index) in instances" :key="index">
            <div class="col-span-full md:col-span-2 lg:col-span-1 group">
                <span class="font-sans font-light text-center  text-lg md:text-lg text-zinc-800/60 dark:text-zinc-50/50">{{ name }}</span>
                <NuxtLink :to="url">
                    <Image :src="image" />
                </NuxtLink>
            </div>
        </template>
    </div>
</template>

<script setup>
const app = useAppConfig();
const { data: { instances } } = app
</script>
