<script setup lang="ts">
import { ref, watch, computed, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { MapPin, Edit, ArrowDownNarrowWide, ArrowUpNarrowWide, FilterX, Filter, ChevronDown, ChevronUp } from 'lucide-vue-next'
import type { Address, CountryCode } from '@/api'

// UI Components
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Slider } from '@/components/ui/slider'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/ui/tooltip'
import {
  Collapsible,
  CollapsibleContent,
} from '@/components/ui/collapsible'

// ------------- Props, Emits, and Types -------------
export type SortField = 'price' | 'speed'
export type SortDirection = 'asc' | 'desc'
export interface SortOptions {
  field: SortField
  direction: SortDirection
}

export interface FilterOptions {
  age?: number
  providers: string[]
  maxPrice?: number
  minSpeed?: number
  includeTV?: boolean
}

export interface QueryOptions {
  sort: SortOptions | null
  filter: FilterOptions
}

const props = defineProps<{
  productCount: number
  forceDisplayMode?: boolean
  queryOptions: QueryOptions
  address: Address | null
}>()

const emit = defineEmits<{
  (e: 'submit', address: Address): void
  (e: 'update-query', options: QueryOptions): void
  (e: 'reset'): void
}>()

watch(() => props.address, (newAddress) => {
  address.value = newAddress
})

// ------------- i18n -------------
const { t } = useI18n()

// ------------- State Management -------------
const isSubmitted = ref(false)
const address = ref<Address | null>(props.address)
const STORAGE_KEY = 'addressForm'

// ------------- Address Form Logic -------------
z.setErrorMap((issue, _) => {
  const field = issue.path[0];
  const code = issue.code;
  if (field === undefined) return { message: t(`errors.unknown`) }
  const params: Record<string, unknown> = {};
  if ("minimum" in issue) params.minimum = issue.minimum;
  if ("maximum" in issue) params.maximum = issue.maximum;

  return { message: t(`addressForm.${field}.errors.${code}`, params) };
});

const formSchema = toTypedSchema(z.object({
  street: z.string().min(2).max(110),
  postalCode: z.string().min(4).max(5),
  city: z.string().min(2).max(50),
  countryCode: z.enum(['DE', 'AT', 'CH']),
}))

function loadInitialValues() {
  if (props.address) {
    address.value = props.address;
  } else {
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      if (raw) {
        const parsed = JSON.parse(raw);
        address.value = parsed;
      }
    } catch { /* ignore errors */ }
  }
  if (address.value === null) {
    return { street: '', postalCode: '', city: '', countryCode: 'DE' as CountryCode };
  }
  // Merge street and houseNumber for the input field
  return { ...address.value, street: `${address.value.street}${address.value.houseNumber ? ' ' + address.value.houseNumber : ''}` };

}


const splitStreetAndHouseNumber = (fullStreet: string) => {
  const match = fullStreet.match(/^(.*?)(?:\s+(\S+))?$/)
  if (!match) return { street: fullStreet, houseNumber: '' }
  return { street: match[1]?.trim() || '', houseNumber: match[2] || '' }
}

const { handleSubmit, setValues } = useForm({
  validationSchema: formSchema,
  initialValues: loadInitialValues(),
})

const onFormSubmit = handleSubmit(formValues => {
  const { street, houseNumber } = splitStreetAndHouseNumber(formValues.street)
  const finalAddress: Address = {
    ...formValues,
    street,
    houseNumber,
    countryCode: formValues.countryCode as CountryCode,
  }
  address.value = finalAddress
  isSubmitted.value = true
  localStorage.setItem(STORAGE_KEY, JSON.stringify(finalAddress))
  emit('submit', finalAddress)
})

function handleEdit() {
  isSubmitted.value = false
  if (address.value) {
    setValues({
      ...address.value,
      street: `${address.value.street}${address.value.houseNumber ? ' ' + address.value.houseNumber : ''}`,
    })
  }
}

const countryCodes = ['DE', 'AT', 'CH']

// ------------- Filter Logic -------------
// Price and speed limits for sliders
const MAX_PRICE = 125 // euros per month
const MAX_SPEED = 250 // Mbps

// Advanced filter states
const showAdvancedFilters = ref(false)

// Local reactive state that mirrors props
const localQueryOptions = reactive<QueryOptions>({
  sort: props.queryOptions.sort,
  filter: { ...props.queryOptions.filter }
})

const availableProviders = ['ByteMe', 'ServusSpeed', 'PingPerfect', 'VerbynDich', 'WebWunder']

function updateQuery() {
  emit('update-query', {
    sort: localQueryOptions.sort,
    filter: { ...localQueryOptions.filter }
  })
}

function updateSort(field: any) {
  if (typeof field === 'string' && (field === 'price' || field === 'speed')) {
    localQueryOptions.sort = {
      field: field as SortField,
      direction: localQueryOptions.sort?.direction || 'asc'
    }
    updateQuery()
  }
}

function toggleSortDirection() {
  if (localQueryOptions.sort) {
    localQueryOptions.sort.direction = localQueryOptions.sort.direction === 'asc' ? 'desc' : 'asc'
    updateQuery()
  }
}

function updateFilter(updates: Partial<FilterOptions>) {
  Object.assign(localQueryOptions.filter, updates)
  updateQuery()
}

function resetAllFilters() {
  Object.assign(localQueryOptions, {
    sort: { field: 'price', direction: 'asc' },
    filter: {
      age: undefined,
      providers: [],
      maxPrice: undefined,
      minSpeed: undefined,
      includeTV: undefined
    }
  })
  showAdvancedFilters.value = false
  emit('reset')
}

// Sync local state with props when they change
watch(() => props.queryOptions, (newOptions) => {
  // Copy the top‐level sort object (that’s okay to replace):
  localQueryOptions.sort = newOptions.sort ? { ...newOptions.sort } : null

  // But for filter, only copy over each nested property:
  Object.assign(localQueryOptions.filter, newOptions.filter)
}, { deep: true })

watch(() => props.forceDisplayMode, (val) => {
  if (val) isSubmitted.value = true
})

// Computed properties for easier template access
const currentSortDirection = computed(() => localQueryOptions.sort?.direction || 'asc')
const filterAge = computed({
  get: () => localQueryOptions.filter.age,
  set: (value: number | undefined) => updateFilter({ age: value })
})
const selectedProviders = computed(() => {
  return localQueryOptions.filter.providers
})
const priceRange = computed({
  get: () => [localQueryOptions.filter.maxPrice ?? MAX_PRICE],
  set: (value: number[]) => updateFilter({ maxPrice: value[0] })
})
const speedRange = computed({
  get: () => [10 + MAX_SPEED - (localQueryOptions.filter.minSpeed ?? 10)],
  set: (value: number[]) => updateFilter({ minSpeed: 10 + MAX_SPEED - value[0] })
})
const includeTV = computed({
  get: () => {
    return !!localQueryOptions.filter.includeTV
  },
  set: (value: boolean) => {
    updateFilter({ includeTV: value })
  }
})

// Computed for provider checkbox v-model
const providerCheckedMap = reactive({} as Record<string, boolean>)

// Sync providerCheckedMap with selectedProviders
watch(
  () => selectedProviders.value,
  (providers) => {
    for (const p of availableProviders) {
      providerCheckedMap[p] = providers.includes(p)
    }
  },
  { immediate: true }
)

// Watch for changes in providerCheckedMap and update selectedProviders
watch(
  providerCheckedMap,
  (map) => {
    const newProviders = availableProviders.filter(p => map[p])
    if (JSON.stringify(newProviders) !== JSON.stringify(selectedProviders.value)) {
      updateFilter({ providers: newProviders })
    }
  },
  { deep: true }
)

</script>

<style scoped>
.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-enter-active {
  transition: opacity 0.15s ease-out, transform 0.15s ease-out;
}

.fade-slide-leave-active {
  transition: opacity 0.1s ease-in, transform 0.1s ease-in;
}

/* Custom divider */
.divider {
  border-top: 1px solid hsl(var(--border));
  margin-top: 1rem;
  padding-top: 1rem;
}

/* Enhanced mobile styles */
@media (max-width: 640px) {
  .filter-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .compact-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 0.75rem;
  }
}

/* Filter badge for active filters */
.filter-badge {
  display: inline-flex;
  align-items: center;
  border-radius: 9999px;
  background-color: rgb(239 246 255);
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  color: rgb(29 78 216);
  box-shadow: 0 0 0 1px rgb(29 78 216 / 0.1);
}

.dark .filter-badge {
  background-color: rgb(59 130 246 / 0.1);
  color: rgb(96 165 250);
  box-shadow: 0 0 0 1px rgb(96 165 250 / 0.2);
}
</style>

<template>
  <Card>
    <CardContent class="p-x-6">
      <transition name="fade-slide" mode="out-in" appear>

        <div v-if="isSubmitted && address" :key="'display'" class="space-y-6">
          <!-- Address Display Section - More Compact -->
          <div class="bg-slate-50 dark:bg-slate-900/50 rounded-lg p-4 border border-slate-200 dark:border-slate-700">
            <div class="flex items-start justify-between">
              <div class="flex items-start space-x-3 flex-1">
                <div class="flex-shrink-0 mt-0.5">
                  <MapPin class="h-4 w-4 text-slate-500" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-slate-900 dark:text-slate-100 text-sm leading-tight">
                    {{ address.street }} {{ address.houseNumber }}
                  </div>
                  <div class="text-xs text-slate-600 dark:text-slate-400 mt-0.5">
                    {{ address.postalCode }}, {{ address.city }}, {{ address.countryCode }}
                  </div>
                </div>
              </div>
              <TooltipProvider :delayDuration=700>
                <Tooltip>
                  <TooltipTrigger>
                    <Button @click="handleEdit" variant="outline" size="icon" class="flex-shrink-0 p-2">
                      <Edit class="h-4 w-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    {{ t('addressForm.edit') }}
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </div>

          <!-- Product Filter Section -->
          <transition name="fade-slide">
            <div v-if="props.productCount > 0" class="space-y-4">
              <div class="h-px bg-gradient-to-r from-transparent via-slate-200 dark:via-slate-700 to-transparent"></div>

              <!-- Basic Sort Controls -->
              <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-3">
                <div class="flex-1 max-w-xs">
                  <div class="flex flex-col gap-1.5">
                    <Label for="sort-field" class="text-sm font-medium text-slate-700 dark:text-slate-300">
                      {{ t('productFilter.sortBy') }}
                    </Label>
                    <div class="flex items-center">
                      <Select :model-value="localQueryOptions.sort?.field || 'price'" @update:model-value="updateSort">
                        <SelectTrigger id="sort-field" class="h-9 rounded-r-none border-r-0">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="price">{{ t('productFilter.price') }}</SelectItem>
                          <SelectItem value="speed">{{ t('productFilter.speed') }}</SelectItem>
                        </SelectContent>
                      </Select>
                      <Button variant="outline" @click="toggleSortDirection" size="icon"
                        class="h-9 w-9 flex items-center justify-center rounded-l-none">
                        <component :is="currentSortDirection === 'asc' ? ArrowUpNarrowWide : ArrowDownNarrowWide"
                          class="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  </div>
                </div>

                <div class="flex items-center gap-1">
                  <Button variant="outline" @click="showAdvancedFilters = !showAdvancedFilters" size="sm"
                    class="flex items-center gap-1 px-2">
                    <Filter class="h-3.5 w-3.5" />
                    <component :is="showAdvancedFilters ? ChevronUp : ChevronDown" class="h-3.5 w-3.5" />
                  </Button>
                  <Button variant="ghost" @click="resetAllFilters" size="sm" class="flex items-center gap-1 px-2">
                    <FilterX class="h-3.5 w-3.5" />
                  </Button>
                  <!-- Slot for additional actions -->
                  <div class="flex items-center ml-1">
                    <slot></slot>
                  </div>
                </div>
              </div>

              <!-- Advanced Filters (Collapsible) -->
              <Collapsible v-model:open="showAdvancedFilters">
                <CollapsibleContent class="space-y-4">
                  <div
                    class="bg-slate-50 dark:bg-slate-900/30 rounded-lg p-4 border border-slate-200 dark:border-slate-700/50">
                    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

                      <!-- Age Input -->
                      <div class="space-y-3">
                        <Label class="text-sm font-medium text-slate-700 dark:text-slate-300">
                          {{ t('productFilter.ageLabel') }}
                        </Label>
                        <Input v-model.number="filterAge" type="number" :placeholder="t('productFilter.agePlaceholder')"
                          min="16" max="100" class="h-9" />
                        <p class="text-xs text-slate-500 dark:text-slate-400">
                          {{ t('productFilter.ageLabel') }}
                        </p>
                      </div>

                      <!-- Provider Selection -->
                      <div class="space-y-3">
                        <Label class="text-sm font-medium text-slate-700 dark:text-slate-300">
                          {{ t('productFilter.providers') }}
                        </Label>
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
                          <div v-for="provider in availableProviders" :key="provider"
                            class="flex items-center space-x-2">
                            <Checkbox :id="provider" v-model="providerCheckedMap[provider]" />
                            <Label :for="provider" class="text-sm font-normal cursor-pointer">
                              {{ provider }}
                            </Label>
                          </div>
                        </div>
                      </div>

                      <!-- Price Range Slider -->
                      <div class="space-y-3">
                        <div class="flex items-center justify-between">
                          <Label class="text-sm font-medium text-slate-700 dark:text-slate-300">
                            {{ t('productFilter.priceRange') }}
                          </Label>
                          <span class="text-xs text-slate-500 dark:text-slate-400">
                            {{ t('productFilter.maxPrice', { max: priceRange[0] }) }}
                          </span>
                        </div>
                        <Slider v-model="priceRange" :max="MAX_PRICE" :min="20" :step="5" class="w-full" />
                        <div class="flex justify-between text-xs text-slate-500 dark:text-slate-400">
                          <span>€20/{{ t('product.month') }}</span>
                          <span>€{{ MAX_PRICE }}/{{ t('product.month') }}</span>
                        </div>
                      </div>

                      <!-- Speed Range Slider -->
                      <div class="space-y-3">
                        <div class="flex items-center justify-between">
                          <Label class="text-sm font-medium text-slate-700 dark:text-slate-300">
                            {{ t('productFilter.speedRange') }}
                          </Label>
                          <span class="text-xs text-slate-500 dark:text-slate-400">
                            {{ t('productFilter.minSpeed', { min: localQueryOptions.filter.minSpeed ?? 10 }) }}
                          </span>
                        </div>
                        <Slider v-model="speedRange" :max="MAX_SPEED" :min="10" :step="10" dir="rtl" class="w-full" />
                        <div class="flex justify-between text-xs text-slate-500 dark:text-slate-400">
                          <span>10 Mbps</span>
                          <span>{{ MAX_SPEED }} Mbps</span>
                        </div>
                      </div>

                      <!-- TV Inclusion Checkbox -->
                      <div class="space-y-3">
                        <Label class="text-sm font-medium text-slate-700 dark:text-slate-300">
                          {{ t('productFilter.tvLabel') }}
                        </Label>
                        <div class="flex items-center space-x-2">
                          <Checkbox id="include-tv" v-model="includeTV" />
                          <Label for="include-tv" class="text-sm font-normal cursor-pointer">
                            {{ t('productFilter.includeTV') }}
                          </Label>
                        </div>
                      </div>

                    </div>

                    <!-- Filter Actions (Mobile) -->
                    <div
                      class="flex justify-between items-center mt-4 pt-4 border-t border-slate-200 dark:border-slate-700/50 sm:hidden">
                      <Button variant="ghost" @click="resetAllFilters" size="sm">
                        {{ t('productFilter.clearFilters') }}
                      </Button>
                      <Button @click="showAdvancedFilters = false" size="sm">
                        {{ t('productFilter.applyFilters') }}
                      </Button>
                    </div>
                  </div>
                </CollapsibleContent>
              </Collapsible>
            </div>
          </transition>
        </div>

        <!-- Form Section -->
        <form v-else :key="'form'" @submit.prevent="onFormSubmit" class="space-y-6">
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <FormField v-slot="{ componentField }" name="countryCode" class="sm:col-span-2">
              <FormItem>
                <FormLabel class="text-sm font-medium text-slate-700 dark:text-slate-300">
                  {{ t('addressForm.country.title') }}
                </FormLabel>
                <Select v-bind="componentField">
                  <FormControl>
                    <SelectTrigger autocomplete="country" class="h-10">
                      <SelectValue :placeholder="t('addressForm.country.placeholder')" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem v-for="code in countryCodes" :key="code" :value="code">
                        {{ t(`addressForm.countries.${code}.name`) }}
                      </SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            </FormField>

            <FormField v-slot="{ componentField }" name="postalCode" class="sm:col-span-1">
              <FormItem>
                <FormLabel class="text-sm font-medium text-slate-700 dark:text-slate-300">
                  {{ t('addressForm.postalCode.title') }}
                </FormLabel>
                <FormControl>
                  <Input type="text" :placeholder="t('addressForm.postalCode.placeholder')" v-bind="componentField"
                    autocomplete="postal-code" class="h-10" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>

          <FormField v-slot="{ componentField }" name="city">
            <FormItem>
              <FormLabel class="text-sm font-medium text-slate-700 dark:text-slate-300">
                {{ t('addressForm.city.title') }}
              </FormLabel>
              <FormControl>
                <Input type="text" :placeholder="t('addressForm.city.placeholder')" v-bind="componentField"
                  autocomplete="address-level2" class="h-10" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="street">
            <FormItem>
              <FormLabel class="text-sm font-medium text-slate-700 dark:text-slate-300">
                {{ t('addressForm.street.title') }}
              </FormLabel>
              <FormControl>
                <Input type="text" :placeholder="t('addressForm.street.placeholder')" v-bind="componentField"
                  autocomplete="street-address" class="h-10" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <div class="pt-2">
            <Button type="submit" class="w-full h-11 font-medium">
              {{ t('actions.compare') }}
            </Button>
          </div>
        </form>

      </transition>
    </CardContent>
  </Card>
</template>