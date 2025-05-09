<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'
import PersonalisationSelector from '@/components/PeronalisationSelector.vue'
import { Avatar, AvatarImage } from '@/components/ui/avatar'

import { useI18n } from 'vue-i18n'
const { t } = useI18n()

</script>

<template>
  <!-- mount Sonner toast container globally -->
  <Toaster />

  <PersonalisationSelector class="fixed bottom-4 right-4 z-50 bg-white dark:bg-secondary" />

  <div class="container mx-auto px-4 mt-6">
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-4">
        <!-- Avatar with link to home -->
        <router-link to="/">
          <Avatar class="size-16">
            <AvatarImage src="/check24-gendev-7/logo.svg" />
          </Avatar>
        </router-link>

        <div>
          <!-- Title with link to home -->
          <router-link to="/" class="hover:text-primary">
            <h1 class="text-2xl font-bold">{{ t('app.title') }}</h1>
          </router-link>
          <p class="text-gray-500">{{ t('app.subtitle') }}</p>
        </div>
      </div>
    </div>
  </div>

  <!-- render the current route's view (ProductsView or SharedProductsView) -->
  <router-view />
</template>