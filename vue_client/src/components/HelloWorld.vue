<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface User {
  id: number
  name: string
  email: string
  age: number
}

const users = ref<User[]>([])
const newUser = ref<User>({ id: 0, name: '', email: '', age: 0 })
const editUser = ref<User | null>(null)

async function fetchUsers() {
  const res = await fetch('/users/api')
  users.value = await res.json()
}

onMounted(fetchUsers)

function startEditUser(user: User) {
  editUser.value = { ...user }
}

async function addUser() {
  const res = await fetch('/users/api', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(newUser.value)
  })
  if (res.ok) {
    await fetchUsers()
    newUser.value = { id: 0, name: '', email: '', age: 0 }
  }
}

async function updateUser() {
  if (!editUser.value) return
  const res = await fetch(`/users/api/${editUser.value.id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(editUser.value)
  })
  if (res.ok) {
    await fetchUsers()
    editUser.value = null
  }
}

async function deleteUser(id: number) {
  const res = await fetch(`/users/api/${id}`, {
    method: 'DELETE'
  })
  if (res.ok) {
    await fetchUsers()
  }
}
</script>

<template>
  <div>
    <h2>Users</h2>
    <ul>
      <li v-for="user in users" :key="user.id">
        {{ user.name }} ({{ user.email }}) - Age: {{ user.age }}
        <button @click="startEditUser(user)">Edit</button>
        <button @click="deleteUser(user.id)">Delete</button>
      </li>
    </ul>
    <h3>Add User</h3>
    <form @submit.prevent="addUser">
      <input v-model="newUser.name" placeholder="Name" required />
      <input v-model="newUser.email" placeholder="Email" type="email" required />
      <input v-model="newUser.age" placeholder="Age" type="number" required />
      <button type="submit">Add User</button>
    </form>
    {{ editUser }}
    <div v-if="editUser.value">
      <h3>Edit User</h3>
      <form @submit.prevent="updateUser">
        <input v-model="editUser.value.name" placeholder="Name" required />
        <input v-model="editUser.value.email" placeholder="Email" type="email" required />
        <input v-model="editUser.value.age" placeholder="Age" type="number" required />
        <button type="submit">Update User</button>
        <button type="button" @click="editUser = null">Cancel</button>
      </form>
    </div>
  </div>
</template>