import type { Ref } from 'vue'
import { ref, watch } from 'vue'

const STORAGE_KEY = 'release-platform-report'

export function useReportStore() {
  const report: Ref<any | null> = ref(null)

  try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    if (raw) {
      report.value = JSON.parse(raw)
    }
  } catch {
    report.value = null
  }

  watch(
    report,
    (val) => {
      if (!val) {
        sessionStorage.removeItem(STORAGE_KEY)
        return
      }
      sessionStorage.setItem(STORAGE_KEY, JSON.stringify(val))
    },
    { deep: true },
  )

  return {
    report,
    setReport(data: any) {
      report.value = data
    },
    clearReport() {
      report.value = null
    },
  }
}
