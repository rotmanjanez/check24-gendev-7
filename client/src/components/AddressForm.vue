<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useI18n } from 'vue-i18n'
import { CountryCode, type Address } from '@/api'

const { t } = useI18n()

z.setErrorMap((issue, _) => {
  // issue.code might be "too_small" or "too_big" or "invalid_type", etc.
  // issue.path will tell you which property failed (e.g. ["street"])
  // issue.params contains { minimum, maximum, … }
  const field = issue.path[0];      // e.g. "street"
  const code = issue.code;         // e.g. "too_small"

  if (field === undefined) {
    return { message: t(`errors.unknown`) }
  }

  const params: Record<string, unknown> = {};
  if ("minimum" in issue) {
    params.minimum = issue.minimum;
  }
  if ("maximum" in issue) {
    params.maximum = issue.maximum;
  }

  // build your translation key dynamically:
  const key = `addressForm.${field}.errors.${code}`;
  return { message: t(key, params) }
});

// key under which we store in localStorage
const STORAGE_KEY = 'addressForm'

// try to load saved values (or fall back to blanks)
function loadInitialValues() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return JSON.parse(raw)
  } catch {
    // ignore parse errors
  }
  return {
    street: '',
    houseNumber: '',
    postalCode: '',
    region: '',
    countryCode: '',
  }
}


// Define Zod schema for address form
const formSchema = toTypedSchema(z.object({
  street: z.string().min(2).max(100),
  houseNumber: z.string().min(1).max(10),
  postalCode: z.string().min(3).max(20),
  region: z.string().min(2).max(50),
  countryCode: z.enum([
    CountryCode.De,
    CountryCode.At,
    CountryCode.Ch,
  ]),
}))

// Initialize vee-validate form
const { handleSubmit, resetForm } = useForm({
  validationSchema: formSchema,
  initialValues: loadInitialValues(),
})

const emit = defineEmits<{
  (e: 'submit', values: Address): void
}>()

function saveValues(vals: any) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(vals))
}

// Submit handler
const submit = handleSubmit(values => {
  saveValues(values)
  // emit the values to the parent component
  emit('submit', values)
  resetForm({ values })
})
const countryCodes = [CountryCode.De, CountryCode.At, CountryCode.Ch]
</script>

<template>
  <form @submit.prevent="submit" class="space-y-4">
    <!-- Country and Region in one row on all screen sizes -->
    <div class="grid grid-cols-3 gap-2">
      <div class="col-span-2">
        <FormField v-slot="{ componentField }" name="countryCode">
          <FormItem>
            <FormLabel>{{ t('addressForm.country.title') }}</FormLabel>
            <FormControl>
              <Select v-bind="componentField">
                <SelectTrigger class="w-full">
                  <SelectValue :placeholder="t('addressForm.country.placeholder')" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem v-for="countryCode in countryCodes" :key="countryCode" :value="countryCode">
                      {{ t(`addressForm.countries.${countryCode}.name`) }}
                    </SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="col-span-1">
        <FormField v-slot="{ componentField }" name="postalCode">
          <FormItem>
            <FormLabel>{{ t('addressForm.postalCode.title') }}</FormLabel>
            <FormControl>
              <Input type="text" :placeholder="t('addressForm.postalCode.placeholder')" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </div>

    <!-- Region on its own row -->
    <div>
      <FormField v-slot="{ componentField }" name="region">
        <FormItem>
          <FormLabel>{{ t('addressForm.region.title') }}</FormLabel>
          <FormControl>
            <Input type="text" :placeholder="t('addressForm.region.placeholder')" v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>

    <!-- Street and House Number in one row -->
    <div class="grid grid-cols-4 gap-2">
      <div class="col-span-3">
        <FormField v-slot="{ componentField }" name="street">
          <FormItem>
            <FormLabel>{{ t('addressForm.street.placeholder') }}</FormLabel>
            <FormControl>
              <Input type="text" :placeholder="t('addressForm.street.placeholder')" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="col-span-1">
        <FormField v-slot="{ componentField }" name="houseNumber">
          <FormItem>
            <FormLabel>{{ t('addressForm.houseNumber.placeholder') }}</FormLabel>
            <FormControl>
              <Input type="text" :placeholder="t('addressForm.houseNumber.placeholder')" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </div>

    <Button type="submit" class="w-full">
      {{ t('addressForm.submit') }}
    </Button>
  </form>
</template>
