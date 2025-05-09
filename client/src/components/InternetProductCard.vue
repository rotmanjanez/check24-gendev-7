<script setup lang="ts">
import { ref } from 'vue'
import { CircleHelp, Cable, Tv, Phone, Smartphone, Calendar, ChevronDown, Percent, Coins, PiggyBank, Zap, Wrench, Clock, ArrowUp, Users } from 'lucide-vue-next'
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Collapsible,
  CollapsibleContent,
} from '@/components/ui/collapsible'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
} from '@/components/ui/table'

import { getAverageMonthlyCost, type ProductContext } from '@/types/ProductContext'
import { useI18n } from 'vue-i18n'

// Props
const props = defineProps<{
  productContext?: ProductContext
  loading?: boolean,
  alwaysShowDetails?: boolean
}>()

// Validate that product is provided when not loading
if (!props.loading && !props.productContext) {
  console.error('InternetProductCard: product is required when not loading')
}

// State
const isDetailsOpen = ref(props.alwaysShowDetails || false)

// Helper functions
const formatCurrency = (cents: number) => {
  return (cents / 100).toLocaleString('de-DE', {
    style: 'currency',
    currency: 'EUR'
  })
}

// Badge colors based on connection type
const connectionTypeBadgeColor = (type: string) => {
  switch (type) {
    case 'FIBER': return 'bg-green-500'
    case 'CABLE': return 'bg-blue-500'
    case 'DSL': return 'bg-orange-500'
    case 'MOBILE': return 'bg-purple-500'
    default: return 'bg-gray-500'
  }
}

const averageCost = getAverageMonthlyCost(props.productContext?.product)

const { t } = useI18n()

const isValidUnthrottledCapacity = (val: number | null | undefined) => {
  return typeof val === 'number' && val !== null && val !== undefined && val > 0 && val !== -2147483648;
};

// Format unthrottled capacity in MB or GB
const formatUnthrottledCapacity = (mb: number) => {
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB`;
  }
  return `${mb} MB`;
};

const isValidSubsequentCost = (cost: any) => {
  return cost && typeof cost.monthlyCostInCent === 'number' && cost.monthlyCostInCent > 0;
};
</script>

<template>
  <Card
    class="w-full max-w-3xl overflow-hidden border shadow-sm hover:shadow-md transition-all gap-0 dark:border-accent/20"
    :class="{ 'cursor-pointer hover:bg-accent/10 dark:hover:bg-accent/20 dark:hover:border-accent/40 hover:ring-1 hover:ring-accent/30': !props.alwaysShowDetails && props.productContext }"
    @click="!props.alwaysShowDetails && props.productContext ? isDetailsOpen = !isDetailsOpen : null">
    <CardContent class="px-3 py-2 sm:px-4 sm:py-3">
      <!-- Compact 4-column layout in single row -->
      <div class="flex justify-between items-center gap-1">
        <!-- Column 1: Title and provider -->
        <div class="flex flex-col mr-1 flex-2 min-w-0">
          <CardTitle v-if="props.productContext" class="text-sm font-bold leading-tight sm:text-base">
            {{ props.productContext.product.name }}
          </CardTitle>
          <template v-else>
            <Skeleton class="h-5 w-24 mb-1 sm:h-6 sm:w-40" />
            <Skeleton class="h-3 w-20 sm:h-4 sm:w-32" />
          </template>
          <CardDescription v-if="props.productContext" class="text-xs mt-1">
            <span class="font-medium">{{ props.productContext.product.provider }}</span>
            <span v-if="props.productContext.product.description"
              class="text-muted-foreground italic ml-1 before:content-['•'] before:mx-1">
              {{ props.productContext.product.description }}
            </span>
          </CardDescription>
        </div>

        <!-- Column 2: Connection type and speed - aligned to the right -->
        <div class="flex flex-col items-center mr-2 flex-1 min-w-0">
          <div v-if="props.productContext" class="flex items-center gap-1 mb-1">
            <!-- Connection type icons based on type -->
            <Smartphone v-if="props.productContext.product.productInfo.connectionType === 'MOBILE'"
              class="h-3 w-3 text-primary sm:h-4 sm:w-4" />
            <Phone v-else-if="props.productContext.product.productInfo.connectionType === 'DSL'"
              class="h-3 w-3 text-primary sm:h-4 sm:w-4" />
            <Tv v-else-if="props.productContext.product.productInfo.connectionType === 'CABLE'"
              class="h-3 w-3 text-primary sm:h-4 sm:w-4" />
            <Cable v-else-if="props.productContext.product.productInfo.connectionType === 'FIBER'"
              class="h-3 w-3 text-primary sm:h-4 sm:w-4" />
            <CircleHelp v-else class="h-3 w-3 text-primary sm:h-4 sm:w-4" />

            <Badge :class="connectionTypeBadgeColor(props.productContext.product.productInfo.connectionType)"
              class="text-xs h-4 px-1 sm:h-5 sm:px-2">
              {{ props.productContext?.product.productInfo.connectionType }}
            </Badge>
          </div>
          <span v-if="props.productContext" class="font-semibold text-xs sm:text-sm">
            {{ props.productContext.product.productInfo.speed }} Mbps
          </span>
          <template v-else>
            <div class="flex items-center gap-1 mb-1 w-full justify-end sm:gap-2">
              <Skeleton class="h-3 w-3 rounded-full sm:h-4 sm:w-4" />
              <Skeleton class="h-4 w-16 rounded sm:h-5 sm:w-20" />
            </div>
            <Skeleton class="h-4 w-12 sm:h-5 sm:w-16" />
          </template>
        </div>

        <!-- Column 3: Pricing info -->
        <div class="flex flex-col items-end mr-2 flex-1 min-w-0">
          <template v-if="props.productContext">
            <div class="flex flex-col items-end">
              <div class="flex items-baseline flex-wrap justify-end">
                <span class="text-base font-bold text-primary sm:text-lg">
                  {{ formatCurrency(averageCost.monthlyCostInCent) }}
                </span>
                <span class="text-xs text-muted-foreground whitespace-normal ml-1">/{{ t('product.month') }}</span>
              </div>
              <span class="text-xs text-muted-foreground text-right leading-tight">
                {{ t('product.over') }} {{ averageCost.durationInMonths }} {{ t('product.months') }}
              </span>
            </div>
          </template>
          <template v-else>
            <Skeleton class="h-5 w-16 mb-1 sm:h-6 sm:w-24" />
            <Skeleton class="h-3 w-12 sm:h-4 sm:w-20" />
          </template>
        </div>

        <!-- Column 4: Dropdown icon -->
        <div class="flex items-center" v-if="props.productContext">
          <div class="rounded-full hover:bg-muted p-1" v-if="!alwaysShowDetails">
            <ChevronDown class="h-4 w-4 sm:h-5 sm:w-5" :class="{ 'transform rotate-180': isDetailsOpen }" />
          </div>
        </div>
        <template v-else>
          <div class="flex items-center">
            <Skeleton class="h-6 w-6 rounded-full sm:h-8 sm:w-8" />
          </div>
        </template>
      </div>

      <!-- Collapsible details section -->
      <Collapsible v-if="props.productContext" v-model:open="isDetailsOpen" class="w-full mt-2">
        <CollapsibleContent class="space-y-2 pt-2 sm:space-y-3" @click.stop>
          <Separator />

          <!-- Details as shadcn table -->
          <Table>
            <TableBody>
              <TableRow>
                <TableCell class="text-center w-8 sm:w-10">
                  <Coins class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ formatCurrency(props.productContext.product.pricing.monthlyCostInCent) }} / {{ t('product.month')
                  }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.contractDurationInMonths">
                <TableCell class="text-center w-8 sm:w-10">
                  <Calendar class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ props.productContext.product.pricing.contractDurationInMonths }} {{
                    t('product.contractDurationMonth') }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.absoluteDiscount">
                <TableCell class="text-center w-8 sm:w-10">
                  <PiggyBank class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ formatCurrency(props.productContext.product.pricing.absoluteDiscount.valueInCent) }} {{
                    t('product.bonus') || 'Bonus' }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.percentageDiscount">
                <TableCell class="text-center w-8 sm:w-10">
                  <Percent class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ props.productContext.product.pricing.percentageDiscount.percentage }}% {{ t('product.discount') ||
                  'Discount' }}
                </TableCell>
              </TableRow>
              <!-- New detail rows -->
              <TableRow v-if="props.productContext.product.productInfo.tv">
                <TableCell class="text-center w-8 sm:w-10">
                  <Tv class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ props.productContext.product.productInfo.tv }} {{ t('product.tv') }}
                </TableCell>
              </TableRow>
              <TableRow
                v-if="isValidUnthrottledCapacity(props.productContext.product.productInfo.unthrottledCapacityMb)">
                <TableCell class="text-center w-8 sm:w-10">
                  <Zap class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ formatUnthrottledCapacity(props.productContext.product.productInfo.unthrottledCapacityMb!) }} {{
                    t('product.unthrottledCapacity') || 'Unthrottled Speed' }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.installationServiceIncluded !== undefined">
                <TableCell class="text-center w-8 sm:w-10">
                  <Wrench class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ t('product.installation.title') || 'Installation' }}:
                  {{ props.productContext.product.pricing.installationServiceIncluded
                  ? (t('product.installation.included') || 'Included')
                  : (t('product.installation.notIncluded') || 'Not included') }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.minAgeInYears">
                <TableCell class="text-center w-8 sm:w-10">
                  <Users class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ t('product.minAge') || 'Minimum age' }}: {{ props.productContext.product.pricing.minAgeInYears }}
                  {{
                    t('product.years') || 'years' }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.maxAgeInJears">
                <TableCell class="text-center w-8 sm:w-10">
                  <Users class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ t('product.maxAge') || 'Maximum age' }}: {{ props.productContext.product.pricing.maxAgeInJears }}
                  {{
                    t('product.years') || 'years' }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.minContractDurationInMonths">
                <TableCell class="text-center w-8 sm:w-10">
                  <Clock class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ t('product.minContractDuration') || 'Minimum contract' }}: {{
                    props.productContext.product.pricing.minContractDurationInMonths }} {{ t('product.months') }}
                </TableCell>
              </TableRow>
              <TableRow v-if="isValidSubsequentCost(props.productContext.product.pricing.subsequentCosts)">
                <TableCell class="text-center w-8 sm:w-10">
                  <ArrowUp class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ t('product.subsequentCosts') || 'Price after' }} {{
                    props.productContext.product.pricing.subsequentCosts?.startMonth }} {{ t('product.months') }}:
                  {{ formatCurrency(props.productContext.product.pricing.subsequentCosts?.monthlyCostInCent ?? 0) }}/{{
                    t('product.month') }}
                </TableCell>
              </TableRow>
              <TableRow v-if="props.productContext.product.pricing.minOrderValueInCent">
                <TableCell class="text-center w-8 sm:w-10">
                  <Coins class="h-4 w-4 text-primary mx-auto sm:h-5 sm:w-5" />
                </TableCell>
                <TableCell class="text-sm sm:text-base">
                  {{ t('product.minOrderValue') || 'Minimum order value' }}: {{
                  formatCurrency(props.productContext.product.pricing.minOrderValueInCent) }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CollapsibleContent>
      </Collapsible>
    </CardContent>
  </Card>
</template>