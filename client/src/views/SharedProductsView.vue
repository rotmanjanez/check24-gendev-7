<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import InternetProductCard from '@/components/InternetProductCard.vue'
import AddressDisplay from '@/components/AddressDisplay.vue'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { Card, CardHeader, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { AlertCircle, Calendar, Undo2 } from 'lucide-vue-next'
import { deserialize, type ProductContext } from '@/types/ProductContext'

const { t, locale } = useI18n()

const route = useRoute()

// state
const productContext = ref<ProductContext | null>(null)
const error = ref<string | undefined>(undefined)

console.log('route', route)

// parse payload on mount
onMounted(() => {
  const productParam = route.query.product
  if (typeof productParam === 'string') {
    let data: ProductContext
    try {
      data = deserialize(productParam)
    } catch (e: unknown) {
      error.value = (e as Error).message || 'Failed to decode payload'
      return
    }

    if (data.version === '1.0.0') {
      productContext.value = data
    } else {
      error.value = `Unsupported payload version: ${data.version}`
    }
  }
})

// Formatted date components for better styling
const dateComponents = computed(() => {
  if (!productContext.value) return null

  const date = new Date(productContext.value.product.dateOffered)

  return {
    fullDate: date.toLocaleString(locale.value, {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      weekday: 'long'
    }),
  }
})
</script>

<template>
  <div class="flex justify-center items-start container mx-auto px-4">
    <div class="space-y-6 w-full max-w-md">
      <!-- Error Alert -->
      <Alert v-if="error" variant="destructive" class="mb-4">
        <AlertCircle class="w-4 h-4" />
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>{{ error }}</AlertDescription>
      </Alert>

      <!-- Enhanced Offer Date Card -->
      <Card v-if="productContext && dateComponents" class="w-full p-3 gap-0">
        <CardHeader class="p-0 w-full">
          <CardDescription class="text-sm flex items-center">
            <Calendar class="w-5 h-5 mr-2" />
            {{ t('product.offeredOn') }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="flex flex-col items-center text-center">
            <div class="mb-2 text-lg font-medium">{{ dateComponents.fullDate }}</div>
          </div>
        </CardContent>
      </Card>

      <AddressDisplay v-if="productContext" :address="productContext.address">
        <Button v-if="productContext">
          <!-- Icon -->
          <Undo2 class="w-4 h-4 mr-2" />
          {{ t('product.backToSearch') }}
        </Button>
      </AddressDisplay>

      <!-- Shared Product Card -->
      <div class="flex justify-center">
        <InternetProductCard v-if="productContext" :product-context="productContext" :always-show-details="true" />
      </div>
    </div>
  </div>
</template>