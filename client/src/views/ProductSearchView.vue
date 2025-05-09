<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AddressCard from '@/components/AddressCard.vue'
import InternetProductCard from '@/components/InternetProductCard.vue'
import ShareButton from '@/components/ShareButton.vue'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { AlertCircle } from 'lucide-vue-next'
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'
import type { QueryOptions } from '@/components/AddressCard.vue'
import type { Address, InternetProduct, CountryCode, InternetProductsResponse } from '@/api'
import { createConfiguration, InternetProductsApi, prod, dev } from '@/api'
import { getAverageMonthlyCost } from '@/types/ProductContext'

const { t } = useI18n()
const route = useRoute()

const address = ref<Address>({ street: '', houseNumber: '', postalCode: '', city: '', countryCode: 'De' as CountryCode })
const products = ref<InternetProduct[]>([])
const sortedProducts = ref<InternetProduct[]>([])
const version = ref('')
const requestCursor = ref<string | null>(null)
const loading = ref(false)
const fetching = ref(false)
const error = ref<string | null>(null)

// Unified query options
const queryOptions = ref<QueryOptions>({
  sort: { field: 'price', direction: 'asc' },
  filter: { age: undefined, providers: [], maxPrice: undefined, minSpeed: undefined, includeTV: undefined }
})

const addressDisplayMode = ref(false)

const itemsPerPage = 7
const currentPage = ref(1)
const highestPage = ref(1)
let toastId: string | number | undefined = undefined

// API client
const config = createConfiguration({ baseServer: import.meta.env.MODE === 'production' ? prod : dev })
const api = new InternetProductsApi(config)

const paginatedProducts = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  return sortedProducts.value.slice(start, start + itemsPerPage)
})

watch(currentPage, (v) => { if (v > highestPage.value) highestPage.value = v })

// apply filtering and sorting when products or options change
watch([products, queryOptions], () => {
  const { sort, filter } = queryOptions.value

  // unified filtering: age, providers, TV inclusion
  const age = filter.age
  const providers = filter.providers || []
  const includeTV = filter.includeTV === true

  let list = products.value.filter(p => {
    // age filter
    if (age != null) {
      const min = p.pricing.minAgeInYears ?? 0
      const max = p.pricing.maxAgeInJears ?? Infinity
      if (age < min || age > max) return false
    }
    // providers filter
    if (providers.length > 0 && !providers.includes(p.provider)) {
      return false
    }
    // TV inclusion filter
    if (includeTV && !Boolean(p.productInfo.tv)) {
      return false
    }
    return true
  })

  // price filter (max monthly cost in euros)
  const maxPrice = filter.maxPrice
  if (maxPrice != null) {
    list = list.filter(p => p.pricing.monthlyCostInCent <= maxPrice * 100)
  }

  // speed filter (min speed in Mbps)
  const minSpeedVal = filter.minSpeed
  if (minSpeedVal != null) {
    list = list.filter(p => p.productInfo.speed >= minSpeedVal)
  }

  // apply sort
  if (sort?.field === 'price') {
    list.sort((a, b) => {
      let costA = getAverageMonthlyCost(a).monthlyCostInCent
      let costB = getAverageMonthlyCost(b).monthlyCostInCent
      return sort.direction === 'asc' ? costA - costB : costB - costA
    })
  } else if (sort?.field === 'speed') {
    list.sort((a, b) => sort.direction === 'asc'
      ? a.productInfo.speed - b.productInfo.speed
      : b.productInfo.speed - a.productInfo.speed)
  }

  sortedProducts.value = list
  currentPage.value = 1
}, { immediate: true, deep: true })

watch(products, (newList, oldList) => {
  const added = newList.length - oldList.length
  if (added > 0 && oldList.length > 0) {
    const newPage = Math.floor(oldList.length / itemsPerPage) + 1
    if (newPage <= highestPage.value) {
      toast.dismiss(toastId)
      toastId = toast(t('notifications.newOffers'), {
        action: { label: t('actions.refresh'), onClick: () => { currentPage.value = newPage; highestPage.value = newPage; toast.dismiss(toastId!) } },
        duration: Infinity,
      })
    }
  }
}, { immediate: false })

async function fetchProducts(reset = true) {
  fetching.value = true
  loading.value = true
  if (reset) { products.value = []; sortedProducts.value = []; currentPage.value = 1; highestPage.value = 1 }
  error.value = null
  toast.dismiss(toastId)
  try {
    const resp = await api.initiateInternetProductsQuery(address.value)
    version.value = resp.version || ''
    requestCursor.value = resp.nextCursor || null
    let cursor = resp.nextCursor || ''

    // give some time for the initial query to complete before starting continuation so that the UI does not update too quickly
    await new Promise(r => setTimeout(r, 3500))

    while (cursor) {
      await new Promise(r => setTimeout(r, 500))

      let resp: InternetProductsResponse | undefined
      resp = await api.continueInternetProductsQuery(cursor)
      if (resp === undefined) {
        // got a 202 response, wait for continuation
        continue;
      }
      products.value.push(...(resp.products || []))
      loading.value = false // Show products as they arrive
      cursor = resp.nextCursor || ''
    }
  } catch (e: any) {
    console.error('Failed to fetch products', e)
    error.value = e.code === 404 ? 'unknown' : 'network'
  } finally {
    loading.value = false
    fetching.value = false
  }
}

function retry() { fetchProducts(false) }

async function prepareShare() {
  if (!requestCursor.value) {
    toast.error(t('errors.noDataToShare'))
    return false
  }
  try {
    await api.shareInternetProducts(requestCursor.value)
    return true
  } catch {
    toast.error(t('errors.shareFailed'))
    return false
  }
}

onMounted(() => {
  const cursor = route.query.cursor as string

  if (cursor) {
    loading.value = true
    error.value = null
    addressDisplayMode.value = true
    api.getSharedInternetProducts(cursor)
      .then(resp => {
        if (resp.address) {
          address.value = resp.address
        }
        if (resp.products && resp.products.length > 0) {
          products.value = resp.products
          sortedProducts.value = [...products.value]
          version.value = resp.version || ''
          requestCursor.value = cursor
          console.log(`Loaded ${resp.products.length} shared products for cursor: ${cursor}`)
        } else {
          error.value = 'noResults'
        }
      })
      .catch(e => {
        console.error('Failed to load shared products', e)
        error.value = 'network'
      })
      .finally(() => {
        loading.value = false
      })
  }
})

function handleQueryUpdate(options: QueryOptions) {
  queryOptions.value = options
}

function resetOptions() {
  queryOptions.value = {
    sort: { field: 'price', direction: 'asc' },
    filter: { age: undefined, providers: [], maxPrice: undefined, minSpeed: undefined, includeTV: undefined }
  }
}
</script>

<template>
  <div class="px-6 max-w-3xl mx-auto">
    <AddressCard
      @submit="async (addr: Address) => { address = addr; await fetchProducts(); addressDisplayMode = false; }"
      @update-query="handleQueryUpdate" @reset="resetOptions" :product-count="products.length"
      :force-display-mode="addressDisplayMode" :query-options="queryOptions" class="mb-4">
      <ShareButton :hash-route="'?cursor=' + requestCursor" :prepare="prepareShare">
        {{ t('product.share') }}
      </ShareButton>
    </AddressCard>

    <Alert v-if="error" variant="destructive" class="mb-4">
      <template #icon>
        <AlertCircle class="w-4 h-4" />
      </template>
      <AlertTitle>{{ t(`errors.${error}`) }}</AlertTitle>
      <AlertDescription><Button @click="retry" variant="destructive">{{ t('actions.retry') }}</Button>
      </AlertDescription>
    </Alert>

    <div class="flex flex-col gap-4">
      <InternetProductCard v-if="loading && products.length === 0" v-for="n in itemsPerPage" :key="n" :loading="true" />
      <div v-if="!loading && paginatedProducts.length === 0 && requestCursor"
        class="flex flex-col items-center p-6 border rounded-lg shadow-sm">
        <AlertCircle class="w-8 h-8 text-muted-foreground mb-2" />
        <h3 class="font-medium text-lg">{{ t('product.noResults') }}</h3>
      </div>

      <InternetProductCard v-else-if="paginatedProducts.length > 0" v-for="product in paginatedProducts"
        :key="product.id" :product-context="{ address, product, version }" />

      <!-- Show loading indicators for additional products being fetched -->
      <InternetProductCard v-if="fetching && !loading && paginatedProducts.length < itemsPerPage"
        v-for="n in Math.min(itemsPerPage - paginatedProducts.length, 3)" :key="'loading-' + n" :loading="true" />
    </div>

    <Pagination v-if="!loading && products.length > itemsPerPage" :page="currentPage" :items-per-page="itemsPerPage"
      :total="sortedProducts.length" @update:page="page => currentPage = page" class="mt-4">
      <PaginationContent v-slot="{ items }">
        <PaginationPrevious />
        <template v-for="item in items" :key="item.type === 'page' ? item.value : Math.random()">
          <PaginationItem v-if="item.type === 'page'" :value="item.value" :is-active="item.value === currentPage">{{
            item.value }}</PaginationItem>
          <PaginationEllipsis v-else />
        </template>
        <PaginationNext />
      </PaginationContent>
    </Pagination>
  </div>
</template>
