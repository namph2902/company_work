<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface User {
  id: number
  name: string
  email: string
  age: number
}

const users = ref<User[]>([])

onMounted(async () => {
  const res = await fetch('/api/users')
  users.value = await res.json()
})
</script>

<template>
  <div>
    <h2>Users</h2>
    <ul>
      <li v-for="user in users" :key="user.id">
        {{ user.name }} ({{ user.email }}) - Age: {{ user.age }}
      </li>
    </ul>
  </div>
</template>