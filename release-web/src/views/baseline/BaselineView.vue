<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getBaseline, updateBaseline } from '../../api/client'

const env = ref('prod')
const appName = ref('dis-modules-itsm')
const baselineText = ref('')

function toObject(text: string) {
  const out: Record<string, string> = {}
  text
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .forEach((line) => {
      const idx = line.indexOf('=')
      if (idx > 0) {
        out[line.slice(0, idx).trim()] = line.slice(idx + 1).trim()
      }
    })
  return out
}

function toText(values: Record<string, string>) {
  return Object.entries(values)
    .map(([k, v]) => `${k}=${v}`)
    .join('\n')
}

async function loadBaseline() {
  const res: any = await getBaseline(env.value, appName.value)
  if (!res.success) {
    ElMessage.error(res.message || '加载失败')
    return
  }
  baselineText.value = toText(res.data || {})
}

async function saveBaseline() {
  const values = toObject(baselineText.value)
  const res: any = await updateBaseline(env.value, appName.value, values)
  if (!res.success) {
    ElMessage.error(res.message || '保存失败')
    return
  }
  ElMessage.success('基线更新成功')
}

onMounted(loadBaseline)
</script>

<template>
  <div>
    <h2 class="page-title">环境基线管理</h2>
    <p class="page-tip">维护 dev/sit/uat/prod 基线配置，供配置比对引擎使用。</p>

    <div class="panel">
      <el-form inline>
        <el-form-item label="环境">
          <el-select v-model="env" style="width: 140px" @change="loadBaseline">
            <el-option label="dev" value="dev" />
            <el-option label="sit" value="sit" />
            <el-option label="uat" value="uat" />
            <el-option label="prod" value="prod" />
          </el-select>
        </el-form-item>
        <el-form-item label="应用">
          <el-input v-model="appName" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button @click="loadBaseline">加载</el-button>
          <el-button type="primary" @click="saveBaseline">保存</el-button>
        </el-form-item>
      </el-form>

      <el-input
        v-model="baselineText"
        type="textarea"
        :rows="20"
        placeholder="每行一个 key=value"
        style="margin-top: 12px"
      />
    </div>
  </div>
</template>
