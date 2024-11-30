// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import LoginPage from '@/components/LoginPage.vue';
import ProfilePage from '@/components/profile.vue';
import ResetPassword from '@/components/ResetPassword.vue';
import ChatBot from '@/components/ChatBot.vue';

// Définir les routes
const routes = [
  { path: '/', name: 'LoginPage', component: LoginPage },
  { path: '/profile', name: 'ProfilePage', component: ProfilePage },
  { path: '/reset-password', name: 'ResetPassword', component: ResetPassword },
  { path: '/chatbot', name: 'ChatBot', component: ChatBot },
];

// Créer le routeur
const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;

