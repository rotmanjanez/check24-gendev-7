<script setup lang="ts">
import PeronalisationSelector from '@/components/PeronalisationSelector.vue'

import { useI18n } from 'vue-i18n'

const { t } = useI18n()

</script>


<template>
  <div class="container mx-auto p-6">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold">{{ t('app.title') }}</h1>
        <p class="text-gray-400">
          {{ t('app.subtitle') }}
        </p>
      </div>
      <PeronalisationSelector />
    </div>
  </div>
</template>