<template>
  <div class="w-full relative h-full pt-12 bg-zinc-50 dark:bg-zinc-950">
    <div class="relative z-20">
      <LazySwitchMode />
      <UContainer>
        <div v-motion-fade class="flex flex-col gap-2 items-center justify-center">
          <LazyPhotoProfile/>
          <LazyMyInfo/>
          <LazySocialNetwork/>
          <h3 class="text-2xl dark:text-white mt-5">Implementaciones para empresas</h3>
          <LazyGallery />
        </div>
      </UContainer>
    </div>
    <Footer />
    <img v-motion-fade-visible src="/shadow.png"
      class="object-cover object-center fixed w-full h-full left-0 top-0 z-10 dark:mix-blend-soft-light"
      alt="shadow-background">
  </div>
</template>

<script setup>
const img = useImage()
useSeoMeta({
  title: 'Néstor Pons | Portfolio',
  ogTitle: 'Néstor Pons | Portfolio',
  description:
    'Discover Hato, the sleek portfolio website perfect for freelancers, photographers, and artists made with Nuxt UI',
  ogDescription:
    'Discover Hato, the sleek portfolio website perfect for freelancers, photographers, and artists made with Nuxt UI',
  ogImage: "potrait.jpg",
  ogImageSecureUrl: 'potrait.jpg',
  ogUrl: 'potrait.jpg',
  twitterImage: "https://res.cloudinary.com/dkr1hluva/image/upload/v1713023084/og/og-hato_roec1l.jpg",
  twitterCard: "summary_large_image",
  twitterTitle: 'Hato | Minimal Portfolio Showcase',
  twitterDescription: 'Discover Hato, the sleek portfolio website perfect for freelancers, photographers, and artists made with Nuxt UI',
  ogImageWidth: 1280,
  ogImageHeight: 720,
});

useHead({
  htmlAttrs: {
    lang: 'es'
  },
  link: [
    {
      rel: 'icon',
      type: 'image/svg',
      href: '/favicon.svg'
    }
  ]
})
</script>
