<script setup lang="ts">
import { computed } from 'vue'
import { VueMonacoDiffEditor } from '@guolao/vue-monaco-editor'

const props = withDefaults(
  defineProps<{
    original: string
    modified: string
    language?: string
    height?: string
  }>(),
  {
    language: 'yaml',
    height: '420px',
  },
)

const options = computed(() => ({
  automaticLayout: true,
  minimap: { enabled: false },
  renderSideBySide: true,
  readOnly: true,
  fontSize: 13,
  scrollBeyondLastLine: false,
}))
</script>

<template>
  <VueMonacoDiffEditor
    :original="props.original"
    :modified="props.modified"
    :language="props.language"
    :height="props.height"
    :options="options"
  />
</template>
