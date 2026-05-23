import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getGroupList, createGroup, renameGroup, deleteGroup } from '@/api/group'

export const useGroupStore = defineStore('group', () => {
  const groups = ref<GroupData[]>([])
  const loading = ref(false)

  async function fetchGroups() {
    loading.value = true
    try {
      const res = await getGroupList()
      groups.value = res.data.data || []
    } catch {
      // ignore
    } finally {
      loading.value = false
    }
  }

  async function addGroup(name: string) {
    const res = await createGroup(name)
    const g = res.data.data!
    groups.value.push(g)
    return g
  }

  async function updateGroupName(id: string, name: string) {
    await renameGroup(id, name)
    const g = groups.value.find((x) => x.id === id)
    if (g) g.name = name
  }

  async function removeGroup(id: string) {
    await deleteGroup(id)
    groups.value = groups.value.filter((g) => g.id !== id)
  }

  return {
    groups,
    loading,
    fetchGroups,
    addGroup,
    updateGroupName,
    removeGroup,
  }
})
