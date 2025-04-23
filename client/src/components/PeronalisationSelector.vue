<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Icon } from '@iconify/vue'
import { useColorMode } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { watch, onMounted } from 'vue'


const mode = useColorMode()
const { t, availableLocales, locale } = useI18n()

onMounted(() => {
  const saved = localStorage.getItem('locale')
  if (saved && availableLocales.includes(saved))
    locale.value = saved
})

watch(locale, (newLocale) => {
  localStorage.setItem('locale', newLocale)
})

</script>

<template>
  <div class="flex rounded border overflow-hidden shadow-sm w-auto inline-flex h-10">
    <!-- Theme Toggle Section -->
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" class="rounded-none border-0 h-10 px-3">
          <Icon icon="radix-icons:moon"
            class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          <Icon icon="radix-icons:sun"
            class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          <span class="sr-only">{{ t('settings.toggleTheme') }}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem @click="mode = 'light'">
          {{ t('settings.themes.light') }}
        </DropdownMenuItem>
        <DropdownMenuItem @click="mode = 'dark'">
          {{ t('settings.themes.dark') }}
        </DropdownMenuItem>
        <DropdownMenuItem @click="mode = 'auto'">
          {{ t('settings.themes.system') }}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <!-- Vertical Divider -->
    <span class="w-[1px] bg-muted h-full"></span>

    <!-- Language Section -->
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" class="rounded-none border-0 h-10 px-3 flex items-center space-x-2">
          <Icon icon="radix-icons:globe" class="h-[1.2rem] w-[1.2rem] mx-0" />
          <span>{{ t('meta.languageCode') }}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem v-for="lang in availableLocales" :key="lang" @click="locale = lang" class="justify-between">
          <span>{{ t('meta.languageName', /* dummy */ 1, { locale: lang }) }}</span>
          <span class="opacity-60 text-xs">{{ t('meta.languageCode', /* dummy */ 1, { locale: lang }) }}</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
</template>
