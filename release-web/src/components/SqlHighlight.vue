<script setup lang="ts">
import { computed } from 'vue'

interface Issue {
  line: number
  level: string
  reason: string
}

const props = defineProps<{
  sql: string
  issues: Issue[]
}>()

const lineMap = computed(() => {
  const m = new Map<number, Issue[]>()
  props.issues.forEach((i) => {
    const list = m.get(i.line) || []
    list.push(i)
    m.set(i.line, list)
  })
  return m
})

const lines = computed(() => props.sql.split('\n'))

function rowClass(lineNo: number) {
  const issues = lineMap.value.get(lineNo)
  if (!issues || !issues.length) return ''
  const max = issues[0].level
  if (max.includes('P0') || max === 'P0') return 'row-p0'
  if (max.includes('P1') || max === 'P1') return 'row-p1'
  return 'row-p2'
}
</script>

<template>
  <div class="sql-box">
    <div v-for="(line, idx) in lines" :key="idx" class="line" :class="rowClass(idx + 1)">
      <span class="num">{{ idx + 1 }}</span>
      <code class="content">{{ line || ' ' }}</code>
    </div>
  </div>
</template>

<style scoped>
.sql-box {
  border: 1px solid rgba(175, 194, 222, 0.22);
  border-radius: 10px;
  overflow: auto;
  max-height: 460px;
  background: #091321;
}

.line {
  display: grid;
  grid-template-columns: 48px 1fr;
  align-items: start;
  min-height: 24px;
}

.num {
  color: #7f96be;
  padding: 3px 8px;
  text-align: right;
  border-right: 1px solid rgba(175, 194, 222, 0.22);
}

.content {
  white-space: pre-wrap;
  color: #d9e4ff;
  padding: 3px 10px;
}

.row-p0 {
  background: rgba(236, 95, 112, 0.24);
}

.row-p1 {
  background: rgba(242, 161, 63, 0.2);
}

.row-p2 {
  background: rgba(102, 182, 255, 0.16);
}
</style>
