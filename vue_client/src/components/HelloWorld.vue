<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

interface User { id: number; name: string; email: string; age: number }

const users = ref<User[]>([])
const newUser = ref<User>({ id: 0, name: '', email: '', age: 0 })
const editUser = ref<User | null>(null)
const newUserEmailError = ref('')
const editUserEmailError = ref('')
const newUserAgeError = ref('')
const editUserAgeError = ref('')

function isValidEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

function isValidAge(age: number): boolean {
  return Number.isInteger(age)&& age > 0
}

function isUniqueEmail(email: string, excludeId: number): boolean {
  return !users.value.some(user => user.email === email&& user.id !== excludeId)
}

async function fetchUsers() {
  const res = await fetch('/api/users')
  users.value = await res.json()
}

onMounted(fetchUsers)

function startEditUser(user: User) {
  editUser.value = { ...user }
  editUserEmailError.value = ''
}

async function addUser() {
  newUserEmailError.value = ''
  newUserAgeError.value = ''
  if (!isValidEmail(newUser.value.email)) {
    newUserEmailError.value = 'Invalid email format'
    return
  }
  if (!isUniqueEmail(newUser.value.email, 0)) {
    newUserEmailError.value = 'Email already exists'
    return
  }
  if (!isValidAge(newUser.value.age)) {
    newUserAgeError.value = 'Age must be a positive integer'
    return
  }
  const res = await fetch('/api/users', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(newUser.value)
  })
  if (res.ok) {
    await fetchUsers()
    newUser.value = { id: 0, name: '', email: '', age: 0 }
  }
}

async function updateUser() {
  newUserEmailError.value = ''
  newUserAgeError.value = ''
  if (!editUser.value) return
  if (!isValidEmail(editUser.value.email)) {
    editUserEmailError.value = 'Invalid email format'
    return
  }
  if (!isUniqueEmail(editUser.value.email, editUser.value.id)) {
    editUserEmailError.value = 'Email already exists'
    return
  }
  if (!isValidAge(editUser.value.age)) {
    editUserAgeError.value = 'Age must be a positive integer'
    return
  }
  const res = await fetch(`/api/users/${editUser.value.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(editUser.value)
  })
  if (res.ok) {
    await fetchUsers()
    editUser.value = null
  }
}

async function deleteUser(id: number) {
  const res = await fetch(`/api/users/${id}`, {
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
    <p v-if="users.length === 0">No users found.</p>
    <ul v-else>
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
      <span v-if="newUserEmailError" style="color: red;">{{ newUserEmailError }}</span>

      <input v-model="newUser.age" placeholder="Age" type="number" required />
      <span v-if="newUserAgeError" style="color: red;">{{ newUserAgeError }}</span>
      <button type="submit">Add User</button>
    </form>

    <div v-if="editUser">
      <h3>Edit User</h3>
      <form @submit.prevent="updateUser">
        <input v-model="editUser.name" placeholder="Name" required />
        <input v-model="editUser.email" placeholder="Email" type="email" required />
        <span v-if="editUserEmailError" style="color: red;">{{ editUserEmailError }}</span>
        <input v-model="editUser.age" placeholder="Age" type="number" required />
        <span v-if="editUserAgeError" style="color: red;">{{ editUserAgeError }}</span>
        <button type="submit">Update User</button>
        <button type="button" @click="editUser = null">Cancel</button>
      </form>
    </div>
  </div>
</template>