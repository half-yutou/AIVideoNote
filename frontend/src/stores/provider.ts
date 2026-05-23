import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getProviders, createProvider, updateProvider, deleteProvider } from '@/api/provider'

export const useProviderStore = defineStore('provider', () => {
  const providers = ref<LLMProviderData[]>([])
  const loading = ref(false)

  async function fetchProviders() {
    loading.value = true
    try {
      const res = await getProviders()
      providers.value = res.data.data || []
    } finally {
      loading.value = false
    }
  }

  async function addProvider(req: ProviderRequest) {
    const res = await createProvider(req)
    providers.value.unshift(res.data.data!)
    return res.data.data!
  }

  async function saveProvider(id: string, req: ProviderUpdateRequest) {
    const res = await updateProvider(id, req)
    const idx = providers.value.findIndex((p) => p.id === id)
    if (idx !== -1) providers.value[idx] = res.data.data!
    return res.data.data!
  }

  async function removeProvider(id: string) {
    await deleteProvider(id)
    providers.value = providers.value.filter((p) => p.id !== id)
  }

  return { providers, loading, fetchProviders, addProvider, saveProvider, removeProvider }
})
