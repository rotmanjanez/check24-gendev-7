<script setup lang="ts">
import { computed } from 'vue'
import { Button, type ButtonVariants } from '@/components/ui/button'
import { toast } from 'vue-sonner'
import { ExternalLink, Copy, Share } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '@/components/ui/context-menu'

const { t } = useI18n();

const props = defineProps<{
  hashRoute: string
  variant?: ButtonVariants['variant']
  prepare?: () => Promise<boolean>
}>()

/** build the full share URL */
const shareUrl = computed(() =>
  `${window.location.origin}${import.meta.env.BASE_URL}#/${props.hashRoute}`
)

const isShareSupported = computed(() => !!navigator.share)

async function onShare() {
  // prepare data if a prepare function is provided
  if (props.prepare) {
    try {
      if (!await props.prepare())
        return // preparation failed, do not proceed with sharing
    } catch (err) {
      console.error('Prepare function failed', err)
      return
    }
  }


  const payload: ShareData = {
    url: shareUrl.value,
  }

  if (navigator.share) {
    try {
      await navigator.share(payload)
    } catch (err) {
      console.warn('Share API failed or cancelled', err)
      toast(t('product.shareFailed'), { description: t('product.tryCopyManually') })
    }
  } else {
    // copy link
    await navigator.clipboard.writeText(shareUrl.value)
    toast(t('product.copyToClipboard'))
  }
}

async function onCopyLink() {
  // prepare data if a prepare function is provided
  if (props.prepare) {
    try {
      if (!await props.prepare())
        return // preparation failed, do not proceed with sharing
    } catch (err) {
      console.error('Prepare function failed', err)
      return
    }
  }

  try {
    await navigator.clipboard.writeText(shareUrl.value)
    toast(t('product.copyToClipboard'))
  } catch (err) {
    // Fallback for Safari/macOS: use execCommand('copy')
    try {
      const textarea = document.createElement('textarea')
      textarea.value = shareUrl.value
      textarea.setAttribute('readonly', '')
      textarea.style.position = 'absolute'
      textarea.style.left = '-9999px'
      document.body.appendChild(textarea)
      textarea.select()
      const successful = document.execCommand('copy')
      document.body.removeChild(textarea)
      if (successful) {
        toast(t('product.copyToClipboard'))
      } else {
        throw new Error('execCommand failed')
      }
    } catch (fallbackErr) {
      console.error('Failed to copy to clipboard', err, fallbackErr)
      toast(t('product.copyFailed'), { description: t('product.clipboardDenied') })
    }
  }
}

async function onShareNative() {
  // prepare data if a prepare function is provided
  if (props.prepare) {
    try {
      if (!await props.prepare())
        return // preparation failed, do not proceed with sharing
    } catch (err) {
      console.error('Prepare function failed', err)
      return
    }
  }

  const payload: ShareData = {
    url: shareUrl.value,
  }

  if (navigator.share) {
    try {
      await navigator.share(payload)
    } catch (err) {
      console.warn('Share API failed or cancelled', err)
      toast(t('product.shareFailed'), { description: t('product.tryCopyManually') })
    }
  } else {
    toast(t('product.shareNotSupported'), { description: t('product.useCopyLinkInstead') })
  }
}
</script>

<template>
  <ContextMenu>
    <ContextMenuTrigger>
      <Button :variant="variant ?? 'default'" class="flex items-center" @click="onShare">
        <ExternalLink />
        <slot></slot>
      </Button>
    </ContextMenuTrigger>
    <ContextMenuContent>
      <ContextMenuItem @click="onShareNative" :disabled="!isShareSupported">
        <Share class="mr-2 h-4 w-4" />
        Share
      </ContextMenuItem>
      <ContextMenuItem @click="onCopyLink">
        <Copy class="mr-2 h-4 w-4" />
        Copy link
      </ContextMenuItem>
    </ContextMenuContent>
  </ContextMenu>
</template>